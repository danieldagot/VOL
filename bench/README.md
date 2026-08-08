# vol-llm-bench

Static source token density benchmark for the [VOL language](../README.md),
living inside the VOL repository at `bench/`.

`go test ./...` and other Go tooling ignore this directory naturally — it
contains no `.go` files. Python artifacts (`.venv/`, `__pycache__/`) are
listed in the root `.gitignore`.

Measures how many tokens equivalent VOL programs use relative to Python, Go,
Rust, and Zig under named OpenAI tokenizers. Ratios are computed per task and
summarized as a median across the suite.

---

## What this measures

> **Source token density of equivalent programs.**

For each task and tokenizer:

```
density_ratio(L) = vol_tokens / L_tokens
```

- ratio < 1.0 → VOL source is denser than L (fewer tokens)
- ratio = 1.0 → same count
- ratio > 1.0 → VOL uses more tokens than L

## What this does NOT measure

- LLM task-success rate
- Generate, compile, and repair round-trip cost
- Which language an LLM produces more correct code for

Those require a generate/repair harness. Protocol:
[`LLM_BENCHMARK.md`](../LLM_BENCHMARK.md). Harness:

```sh
make llm-dry                      # reference solutions only — not an LLM result
make llm-ollama                   # live 2-task smoke via local Ollama
make llm-core                     # live core_v2 (includes parity toys)
make llm-intent                   # live intent_v1 (filter/map/count + repair/mod)
make llm-dry-intent               # dry-run intent reference solutions
make llm-ollama MODEL=llama3.1:8b
# run one task, failing a model request after 30 seconds:
make llm-ollama MODEL=llama3.1:8b TASKS=06-where-sum REQUEST_TIMEOUT=30
# direct equivalent:
uv run python llm/harness/run_generate_repair.py --provider ollama --model llama3.1:8b \
  --suite smoke --tasks 06-where-sum --request-timeout 30
# cloud:
# OPENAI_API_KEY=… uv run python llm/harness/run_generate_repair.py --provider openai --suite smoke
# Gemini (loads GEMINI_API_KEY and optional GEMINI_MODEL from ../.env):
uv run python llm/harness/run_generate_repair.py --provider gemini --suite smoke
```

The workflow benchmark has:

- `smoke` (`smoke_v1`) — wiring only
- `core` (`core_v2`) — published continuity (includes fib/arrays parity)
- **`intent` (`intent_v1`)** — prefer for language-use claims (filter/map/count
  pipelines + repair + modification; no hello/fib toys)

Default languages are `vol,python`. Summaries report prompt vs completion and
cold vs warm (card-amortized) totals:

```sh
# Prefer for “how does an LLM use VOL?”
uv run python llm/harness/run_generate_repair.py --provider gemini --suite intent
# Historical core_v2 shape only when comparing to published tables:
uv run python llm/harness/run_generate_repair.py --provider gemini --suite core
```

**When to re-run Python:** only if the Python card, suite tasks, protocol, model,
temperature, or harness scoring change. VOL card tweaks, Fix-text diagnostics,
source-check hygiene, and SF-2 surface/card polish can use `--langs vol` with
`--baseline-jsonl` so the published JSONL/summary stays self-contained.

The separate static-density benchmark below has its own 16-task tiered suite.
Those density tasks are not the LLM workflow suite.

Primary published language-use table: Gemini `intent_v1` with `vol_v2` / SF-2
([`llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-051437.md`](llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-051437.md))
— both languages 100% first-try / success @ K; VOL about −10.0% cold / −13.1%
warm vs Python (card ~316 vs ~336). One model; not a broad claim.

Historical continuity: protocol-v1.1 `core_v2` with `vol_v1` / SF-1 is
[`llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-041440.md`](llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-041440.md)
(+11.2% cold / −3.3% warm; card ~436 vs ~336).

`--tasks` accepts a comma-separated list of task IDs (for example,
`01-hello,07-functions`). `--request-timeout SECONDS` caps each model API
request; omit it to use the provider default (600 seconds for Ollama, 120
seconds for cloud endpoints).

Tracked also in [`IDEAS.md`](../IDEAS.md) under *Compiler Metrics and LLM Evaluation*.

Tokenizer choice affects absolute counts. Numbers are always reported alongside
the tokenizer name. GPT, Claude, and other models may tokenize differently.

---

## Prerequisites

| Tool | Purpose |
|------|---------|
| Python ≥ 3.11 | harness scripts |
| `tiktoken` | token counting |
| `go` | run Go tasks |
| `rustc` | compile Rust tasks |
| `zig` | run Zig tasks (optional; skipped if not found) |
| `vol` binary or `go` + VOL source | run VOL tasks |

Install Python dependencies:

```sh
pip install -r requirements.txt
```

Or with [uv](https://github.com/astral-sh/uv):

```sh
uv add -r requirements.txt
# then run scripts with: uv run python harness/...
# or: .venv/bin/python harness/...
```

---

## Setup

This benchmark lives inside the VOL repo at `bench/`. No path configuration
needed — `check_stdout.py` automatically resolves the VOL interpreter from its
parent directory.

To build the Python environment:

```sh
cd bench
uv sync
```

---

## Running

```sh
cd bench

make check   # verify all programs produce correct stdout
make count   # count tokens → results/density*.csv + density.md
make bench   # run both
```

Or without make:

```sh
uv run python harness/check_stdout.py
uv run python harness/count_tokens.py
```

`VOL_BIN=/path/to/vol` can be set to use a pre-built binary instead of
`go run ./cmd/vol`.

Writes to `results/`:

| File | Contents |
|------|----------|
| `density_cl100k_base.csv` | per-task counts and ratios, cl100k_base tokenizer |
| `density_o200k_base.csv` | per-task counts and ratios, o200k_base tokenizer |
| `density_cl100k_base.md` | markdown table |
| `density_o200k_base.md` | markdown table |
| `density.md` | combined report (both tokenizers) |

Token counting reads source files directly and does not require any language
toolchain to be installed.

---

## Task suite

All tasks use only VOL features that are currently implemented in the
interpreter. Each task has five equivalent implementations producing identical
stdout. Tasks are tagged into tiers so print-label glue is not mistaken for
semantic-density wins (see `harness/count_tokens.py` and reports).

| Tier | Role |
|------|------|
| **parity** | Control surface near Python (hello, loops, fib, …) |
| **labeled** | Report-style string labels — sensitive to `print` / `string()` glue |
| **compression** | Bare numeric output; filter / map / count / sum intent |

| ID | Tier | Intent | VOL features exercised |
|----|------|--------|------------------------|
| 01-hello | parity | print variables | literals, `print`, `string()` |
| 02-arithmetic | parity | compute and print | `:=`, `+`, `*`, `-` |
| 03-conditions | parity | branch on bool + int | `? :`, `and` |
| 04-loops | parity | countdown + fixed repeats | `while`, `repeat` |
| 05-arrays-each | parity | index, length, iterate | arrays, `.len`, `.where`, `.each` |
| 06-where-sum | compression | filter + aggregate | `.where`, `.sum()`, `assert` |
| 07-functions | parity | two named functions | expression-body `fn` |
| 08-strings-assert | parity | string ops + assertion | `.len`, `+`, `assert` |
| 09-grade-report | labeled | multi-bucket report | `.count`, `.sum` |
| 10-fibonacci | parity | sequence | `repeat`, multi-assign |
| 11-leaderboard | labeled | compare aggregates | `.sum`, `.count`, `? :` |
| 12-revenue | labeled | filter + sum report | `.where`, `.count`, `.sum` |
| 13-temperatures | labeled | band report + asserts | `.count` |
| 14-pipeline-stats | compression | count/sum/map pipeline | `.count`, `.where`, `.map`, `.sum` |
| 15-band-counts | compression | bands, bare numbers | `.count`, `.sum`, `.len` |
| 16-map-filter | compression | map then count/sum | `.map`, `.count`, `.where` |

### Equivalence rules

- Same observable stdout (and exit 0).
- Required language boilerplate only — no golfing, no artificial padding.
- Idiomatic style in each language.
- Python: single-file `main.py` with `python3` (stdlib only).
- Rust: single-file `main.rs` with `rustc` (no Cargo.toml boilerplate).
- Zig: single-file `main.zig` with `zig run`.
- Go: `go run main.go` (no external modules, only `fmt`).

---

## Interpreting results

Reports include **median (all)** plus tier medians:

- **median (compression)** — best read on VOL collection intent vs loops/comprehensions.
- **median (labeled)** — moves a lot if `print` / string coercion changes; do not
  treat as pure semantic-density proof.
- **median (parity)** — control floor; often near Python.

A ratio below 1.0 on compression tasks reflects `.where` / `.map` / `.count` /
`.sum` versus list comprehensions (Python), explicit loops (Go/Zig), or iterator
chains (Rust).

Do not generalize from this small suite. It anchors the "denser syntax" claim to
measured numbers, not LLM workflow superiority. See also
[`TOKEN_EFFICIENCY.md`](../TOKEN_EFFICIENCY.md).

---

## Roadmap

- [ ] Extend to 20+ tasks once VOL has more settled features
- [ ] Add Hugging Face tokenizer counts for open-weight models
- [ ] Graduate to a generate/repair harness for task-success metrics (see `LLM_BENCHMARK.md` in VOL)
