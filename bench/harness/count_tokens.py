#!/usr/bin/env python3
"""
count_tokens.py — count source tokens for every task in each language,
compute VOL/baseline density ratios, and write results/density_*.csv and
results/density.md.

Usage:
    pip install -r requirements.txt
    python harness/count_tokens.py

Outputs (in results/):
    density_cl100k_base.csv
    density_o200k_base.csv
    density_cl100k_base.md
    density_o200k_base.md
    density.md              — combined report

What the numbers mean:
    density_ratio(L) = vol_tokens / L_tokens
    ratio < 1.0  →  VOL source is denser than L under this tokenizer
    ratio = 1.0  →  same token count
    ratio > 1.0  →  VOL source uses more tokens than L

Task tiers:
    parity       — control surface (hello, loops, fib, …); often near Python
    compression  — filter/map/count/sum intent; measures semantic density
    labeled      — report-style tasks with string labels (print-glue sensitive)
    stdlib       — SF-3 @std / library intent (strings, path, env, json, process)

Limitations:
    - Measures source token density of hand-written equivalent programs only.
    - Does NOT measure LLM task-success efficiency or generate/repair cost.
    - Ratios depend on the chosen tokenizer; different models tokenize differently.
    - Uses only VOL features that are currently implemented in the interpreter.
"""

import csv
import statistics
from pathlib import Path

import tiktoken

REPO_ROOT = Path(__file__).parent.parent.resolve()
TASKS_DIR = REPO_ROOT / "tasks"
RESULTS_DIR = REPO_ROOT / "results"

TOKENIZERS = ["cl100k_base", "o200k_base"]
LANGS = ["vol", "python", "go", "rust", "zig"]
BASELINES = ["python", "go", "rust", "zig"]
SOURCES = {
    "vol": "main.vol",
    "python": "main.py",
    "go": "main.go",
    "rust": "main.rs",
    "zig": "main.zig",
}

# Compression-tier tasks stress collection intent with little/no label glue.
COMPRESSION_TASKS = frozenset(
    {
        "06-where-sum",
        "14-pipeline-stats",
        "15-band-counts",
        "16-map-filter",
    }
)
# Labeled report tasks (string labels dominate; sensitive to print/coercion).
LABELED_TASKS = frozenset(
    {
        "09-grade-report",
        "11-leaderboard",
        "12-revenue",
        "13-temperatures",
    }
)
# SF-3 library surface (@std / idiomatic stdlib peers).
STDLIB_TASKS = frozenset(
    {
        "17-strings-ops",
        "18-path-parts",
        "19-env-default",
        "20-json-fields",
        "21-process-echo",
    }
)


def task_tier(name: str) -> str:
    if name in COMPRESSION_TASKS:
        return "compression"
    if name in LABELED_TASKS:
        return "labeled"
    if name in STDLIB_TASKS:
        return "stdlib"
    return "parity"


def count_tokens(text: str, enc: tiktoken.Encoding) -> int:
    return len(enc.encode(text))


def fmt_ratio(r: float) -> str:
    return f"{r:.3f}"


def median_ratios(rows: list[dict], baseline: str) -> float:
    return statistics.median(r[f"vol/{baseline}"] for r in rows)


def tier_median_line(rows: list[dict], tier: str) -> str:
    subset = [r for r in rows if r["tier"] == tier]
    if not subset:
        return f"| **median ({tier})** | | | | | | | — | — | — | — |"
    return (
        f"| **median ({tier}, n={len(subset)})** | | | | | | "
        f"| **{fmt_ratio(median_ratios(subset, 'python'))}** "
        f"| **{fmt_ratio(median_ratios(subset, 'go'))}** "
        f"| **{fmt_ratio(median_ratios(subset, 'rust'))}** "
        f"| **{fmt_ratio(median_ratios(subset, 'zig'))}** |"
    )


def build_table(rows: list[dict], tok_name: str) -> str:
    header = (
        f"## Token density — tokenizer: `{tok_name}`\n\n"
        "| Task | Tier | VOL | Python | Go | Rust | Zig "
        "| VOL/Python | VOL/Go | VOL/Rust | VOL/Zig |\n"
        "|------|------|-----|--------|----|------|-----"
        "|------------|--------|----------|---------|\n"
    )
    body_lines = []
    for r in rows:
        body_lines.append(
            f"| {r['task']} "
            f"| {r['tier']} "
            f"| {r['vol']} "
            f"| {r['python']} "
            f"| {r['go']} "
            f"| {r['rust']} "
            f"| {r['zig']} "
            f"| {fmt_ratio(r['vol/python'])} "
            f"| {fmt_ratio(r['vol/go'])} "
            f"| {fmt_ratio(r['vol/rust'])} "
            f"| {fmt_ratio(r['vol/zig'])} |"
        )
    body_lines.append(
        "| **median (all)** | | | | | | "
        f"| **{fmt_ratio(median_ratios(rows, 'python'))}** "
        f"| **{fmt_ratio(median_ratios(rows, 'go'))}** "
        f"| **{fmt_ratio(median_ratios(rows, 'rust'))}** "
        f"| **{fmt_ratio(median_ratios(rows, 'zig'))}** |"
    )
    for tier in ("compression", "labeled", "parity", "stdlib"):
        body_lines.append(tier_median_line(rows, tier))
    return header + "\n".join(body_lines) + "\n"


def main() -> None:
    RESULTS_DIR.mkdir(exist_ok=True)

    tasks = sorted(d for d in TASKS_DIR.iterdir() if d.is_dir())
    if not tasks:
        print("No tasks found.")
        return

    combined_sections: list[str] = [
        "# VOL Source Token Density\n",
        "Measures how many tokens equivalent VOL programs use relative to Python, Go, Rust, and Zig.",
        "A ratio < 1.0 means VOL is denser (fewer tokens) under that tokenizer.\n",
        "> **What this measures:** source token density of hand-written equivalent programs.",
        "> **What this does not measure:** LLM task-success rate or generate/repair cost.",
        "> **Tokenizer note:** ratios depend on the tokenizer; GPT, Claude, and other models",
        "> tokenize differently. Numbers are reported per tokenizer.",
        f"> **Suite size:** {len(tasks)} tasks (parity / labeled / compression / stdlib tiers).",
        "> Prefer **median (compression)** when judging semantic-density ops;",
        "> **median (labeled)** is sensitive to print/string glue;",
        "> **median (parity)** is a control near-Python floor;",
        "> **median (stdlib)** is SF-3 `@std` / library intent vs peer stdlibs.\n",
    ]

    for tok_name in TOKENIZERS:
        enc = tiktoken.get_encoding(tok_name)
        rows: list[dict] = []

        for task in tasks:
            row: dict = {"task": task.name, "tier": task_tier(task.name)}
            for lang in LANGS:
                src = task / lang / SOURCES[lang]
                text = src.read_text()
                row[lang] = count_tokens(text, enc)
            for baseline in BASELINES:
                row[f"vol/{baseline}"] = row["vol"] / row[baseline]
            rows.append(row)

        # CSV
        csv_path = RESULTS_DIR / f"density_{tok_name}.csv"
        fieldnames = ["task", "tier"] + LANGS + [f"vol/{b}" for b in BASELINES]
        with csv_path.open("w", newline="", encoding="utf-8") as f:
            writer = csv.DictWriter(f, fieldnames=fieldnames, lineterminator="\n")
            writer.writeheader()
            for r in rows:
                writer.writerow(
                    {
                        "task": r["task"],
                        "tier": r["tier"],
                        **{lang: r[lang] for lang in LANGS},
                        **{f"vol/{b}": fmt_ratio(r[f"vol/{b}"]) for b in BASELINES},
                    }
                )

        table_md = build_table(rows, tok_name)
        md_path = RESULTS_DIR / f"density_{tok_name}.md"
        md_path.write_text(table_md, encoding="utf-8")
        combined_sections.append(table_md)

        def fmt_tier(tier: str) -> str:
            subset = [r for r in rows if r["tier"] == tier]
            if not subset:
                return f"{tier}=—"
            return (
                f"{tier} VOL/Python={fmt_ratio(median_ratios(subset, 'python'))}"
            )

        print(
            f"[{tok_name}]  all VOL/Python={fmt_ratio(median_ratios(rows, 'python'))}"
            f"  VOL/Go={fmt_ratio(median_ratios(rows, 'go'))}"
            f"  VOL/Rust={fmt_ratio(median_ratios(rows, 'rust'))}"
            f"  VOL/Zig={fmt_ratio(median_ratios(rows, 'zig'))}"
        )
        print(
            f"  tiers: {fmt_tier('compression')}; {fmt_tier('labeled')}; "
            f"{fmt_tier('parity')}; {fmt_tier('stdlib')}"
        )
        print(f"  Wrote {csv_path.relative_to(REPO_ROOT)}")
        print(f"  Wrote {md_path.relative_to(REPO_ROOT)}")

    combined_path = RESULTS_DIR / "density.md"
    combined_path.write_text("\n".join(combined_sections) + "\n", encoding="utf-8")
    print(f"\nCombined report: {combined_path.relative_to(REPO_ROOT)}")


if __name__ == "__main__":
    main()
