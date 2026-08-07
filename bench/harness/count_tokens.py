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

Limitations:
    - Measures source token density of hand-written equivalent programs only.
    - Does NOT measure LLM task-success efficiency or generate/repair cost.
    - Ratios depend on the chosen tokenizer; different models tokenize differently.
    - Suite is small (8 tasks) and uses only VOL features that are currently
      implemented in the interpreter.
"""

import csv
import statistics
from pathlib import Path

import tiktoken

REPO_ROOT = Path(__file__).parent.parent.resolve()
TASKS_DIR = REPO_ROOT / "tasks"
RESULTS_DIR = REPO_ROOT / "results"

TOKENIZERS = ["cl100k_base", "o200k_base"]
LANGS = ["vol", "go", "rust", "zig"]
SOURCES = {
    "vol": "main.vol",
    "go": "main.go",
    "rust": "main.rs",
    "zig": "main.zig",
}


def count_tokens(text: str, enc: tiktoken.Encoding) -> int:
    return len(enc.encode(text))


def fmt_ratio(r: float) -> str:
    return f"{r:.3f}"


def build_table(rows: list[dict], tok_name: str) -> str:
    header = (
        f"## Token density — tokenizer: `{tok_name}`\n\n"
        "| Task | VOL | Go | Rust | Zig | VOL/Go | VOL/Rust | VOL/Zig |\n"
        "|------|-----|----|------|-----|--------|----------|---------|\n"
    )
    body_lines = []
    for r in rows:
        body_lines.append(
            f"| {r['task']} "
            f"| {r['vol']} "
            f"| {r['go']} "
            f"| {r['rust']} "
            f"| {r['zig']} "
            f"| {fmt_ratio(r['vol/go'])} "
            f"| {fmt_ratio(r['vol/rust'])} "
            f"| {fmt_ratio(r['vol/zig'])} |"
        )
    med_go = statistics.median(r["vol/go"] for r in rows)
    med_rust = statistics.median(r["vol/rust"] for r in rows)
    med_zig = statistics.median(r["vol/zig"] for r in rows)
    body_lines.append(
        f"| **median** | | | | | **{fmt_ratio(med_go)}** "
        f"| **{fmt_ratio(med_rust)}** | **{fmt_ratio(med_zig)}** |"
    )
    return header + "\n".join(body_lines) + "\n"


def main() -> None:
    RESULTS_DIR.mkdir(exist_ok=True)

    tasks = sorted(d for d in TASKS_DIR.iterdir() if d.is_dir())
    if not tasks:
        print("No tasks found.")
        return

    combined_sections: list[str] = [
        "# VOL Source Token Density\n",
        "Measures how many tokens equivalent VOL programs use relative to Go, Rust, and Zig.",
        "A ratio < 1.0 means VOL is denser (fewer tokens) under that tokenizer.\n",
        "> **What this measures:** source token density of hand-written equivalent programs.",
        "> **What this does not measure:** LLM task-success rate or generate/repair cost.",
        "> **Tokenizer note:** ratios depend on the tokenizer; GPT, Claude, and other models",
        "> tokenize differently. Numbers are reported per tokenizer.",
        "> **Suite size:** 8 tasks, using only currently-implemented VOL features.\n",
    ]

    for tok_name in TOKENIZERS:
        enc = tiktoken.get_encoding(tok_name)
        rows: list[dict] = []

        for task in tasks:
            row: dict = {"task": task.name}
            for lang in LANGS:
                src = task / lang / SOURCES[lang]
                text = src.read_text()
                row[lang] = count_tokens(text, enc)
            for baseline in ("go", "rust", "zig"):
                row[f"vol/{baseline}"] = row["vol"] / row[baseline]
            rows.append(row)

        # CSV
        csv_path = RESULTS_DIR / f"density_{tok_name}.csv"
        fieldnames = ["task"] + LANGS + ["vol/go", "vol/rust", "vol/zig"]
        with open(csv_path, "w", newline="") as f:
            writer = csv.DictWriter(f, fieldnames=fieldnames)
            writer.writeheader()
            for r in rows:
                writer.writerow(
                    {
                        "task": r["task"],
                        **{lang: r[lang] for lang in LANGS},
                        "vol/go": fmt_ratio(r["vol/go"]),
                        "vol/rust": fmt_ratio(r["vol/rust"]),
                        "vol/zig": fmt_ratio(r["vol/zig"]),
                    }
                )

        # Per-tokenizer markdown
        table_md = build_table(rows, tok_name)
        md_path = RESULTS_DIR / f"density_{tok_name}.md"
        md_path.write_text(table_md)

        combined_sections.append(table_md)

        med_go = statistics.median(r["vol/go"] for r in rows)
        med_rust = statistics.median(r["vol/rust"] for r in rows)
        med_zig = statistics.median(r["vol/zig"] for r in rows)
        print(f"[{tok_name}]  median VOL/Go={fmt_ratio(med_go)}  VOL/Rust={fmt_ratio(med_rust)}  VOL/Zig={fmt_ratio(med_zig)}")
        print(f"  Wrote {csv_path.relative_to(REPO_ROOT)}")
        print(f"  Wrote {md_path.relative_to(REPO_ROOT)}")

    combined_path = RESULTS_DIR / "density.md"
    combined_path.write_text("\n".join(combined_sections) + "\n")
    print(f"\nCombined report: {combined_path.relative_to(REPO_ROOT)}")


if __name__ == "__main__":
    main()
