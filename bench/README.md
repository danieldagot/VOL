# vol-llm-bench

Static source token density benchmark for the [VOL language](../README.md),
living inside the VOL repository at `bench/`.

`go test ./...` and other Go tooling ignore this directory naturally — it
contains no `.go` files. Python artifacts (`.venv/`, `__pycache__/`) are
listed in the root `.gitignore`.

Measures how many tokens equivalent VOL programs use relative to Go, Rust, and
Zig under named OpenAI tokenizers. Ratios are computed per task and summarized
as a median across the suite.

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
make llm-core                     # live 5-task core, 3 replicates
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

The workflow benchmark has a 2-task `smoke` suite for wiring checks and a
5-task `core` suite (`core_v2`, protocol v1.1) for reported comparisons. Core
covers generation, diagnostic-seeded repair, and modification, and defaults to
three replicates. Default languages are `vol,python` (interpreter baseline);
pass `--langs vol,go` or `--langs vol,python,go` for Go. Summaries report prompt
vs completion and cold vs warm (card-amortized) totals:

```sh
uv run python llm/harness/run_generate_repair.py --provider gemini --suite core
# optional compiled baseline:
uv run python llm/harness/run_generate_repair.py --provider gemini --suite core --langs vol,go
```

The separate static-density benchmark below still contains 13 tasks. Those 13
are not the LLM workflow suite.

The published protocol-v1.1 (`core_v2`) Gemini run with the default Python
baseline is
[`llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-022642.md`](llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-022642.md).
Both languages reached 100% first-try and success @ K; VOL was about +10.0% cold
tokens and about −4.9% warm tokens vs Python. Earlier VOL vs Go:
[`llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-021122.md`](llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-021122.md).
Historical non-diagnostic `core_v1`:
[`llm/results/core_v1_live_gemini_gemini-3.5-flash-lite_20260808-014539.md`](llm/results/core_v1_live_gemini_gemini-3.5-flash-lite_20260808-014539.md).

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

All 13 tasks use only VOL features that are currently implemented in the
interpreter (as of the VOL prototype). Each task has four equivalent
implementations producing identical stdout.

| ID | Intent | VOL features exercised |
|----|--------|------------------------|
| 01-hello | print variables | literals, `print`, `string()` |
| 02-arithmetic | compute and print | `:=`, `=`, `+`, `*`, `-` |
| 03-conditions | branch on bool + int | `if`/`else`, `and` |
| 04-loops | countdown + fixed repeats | `while`, `repeat` |
| 05-arrays-each | index, length, iterate | arrays, `.len`, `.each` |
| 06-where-sum | filter + aggregate | `.where`, `.sum()`, `assert` |
| 07-functions | two named functions | `fn`, `return` |
| 08-strings-assert | string ops + assertion | `.len`, `+`, `assert` |

### Equivalence rules

- Same observable stdout (and exit 0).
- Required language boilerplate only — no golfing, no artificial padding.
- Idiomatic style in each language.
- Rust: single-file `main.rs` with `rustc` (no Cargo.toml boilerplate).
- Zig: single-file `main.zig` with `zig run`.
- Go: `go run main.go` (no external modules, only `fmt`).

---

## Interpreting results

A ratio below 1.0 on task 06-where-sum reflects VOL's `.where(...).sum()`
compression compared with explicit filter loops in Go/Zig or iterator chains in
Rust.

Tasks 01–04 and 07–08 exercise constructs that are syntactically similar across
all four languages. Differences there reflect keyword weight, required
boilerplate (`fn main`, `package main`, `println!`, etc.), and type annotation
requirements.

Do not generalize from 13 tasks. This suite exists to anchor the "denser syntax"
claim to measured numbers, not to prove overall LLM workflow superiority.

---

## Roadmap

- [ ] Extend to 20+ tasks once VOL has more settled features
- [ ] Add Hugging Face tokenizer counts for open-weight models
- [ ] Graduate to a generate/repair harness for task-success metrics (see `LLM_BENCHMARK.md` in VOL)
