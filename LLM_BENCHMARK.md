# VOL LLM Generate / Repair Benchmark (Protocol)

> Status: **protocol + harness scaffolded; live generate/repair results not yet run**  
> Companion to: [`bench/`](bench/README.md) (source token density only)

This document defines how VOL measures whether it is actually better for
LLM workflows—not just denser when a human already wrote the answer.

**Static source density** lives in [`bench/`](bench/README.md).  
**This document** covers generate → run → repair efficiency.

---

## 1. Claim under test

Falsifiable claim (fill in numbers after runs):

> On suite **S**, model **M**, temperature **T**, max repair rounds **K**:  
> VOL achieved **X%** task success with **Y%** fewer **total tokens**  
> (prompt + generation + repair) than language **L**.

If VOL is denser in source tokens but costs more total workflow tokens, the
LLM-optimization claim fails for that suite. Report the loss honestly.

---

## 2. What this measures vs what it does not

| Measures | Does not measure |
| --- | --- |
| First-try success rate | Hand-written source density (see `bench/`) |
| Success within ≤ K repair rounds | Human readability preference |
| Total tokens to success or give-up | Native performance / binary size |
| How often structured diagnostics enable a fix | Ownership, types, or backend quality |
| Relative cost vs Go (v0 baseline) | “Meaning per Token” as a single magic score |

Tokenizer note: absolute token counts depend on the model’s tokenizer. Always
report model id + tokenizer (or API usage fields) with every table.

---

## 3. Primary metric

```text
workflow_efficiency(L) = successful_tasks(L) / total_tokens_consumed(L)
```

Where `total_tokens_consumed` includes, for every attempt and every task:

1. System / language-reference tokens in the prompt
2. Task-prompt tokens
3. Model completion tokens (generated source)
4. Tokens spent feeding diagnostics back on repair rounds
5. Repair completion tokens

Also report the three headline rates (do not collapse into one number only):

| Metric | Definition |
| --- | --- |
| **First-try success** | Share of tasks where attempt 0 stdout matches `expected.txt` |
| **Success @ K** | Share of tasks solved within ≤ K repairs (K default = 2) |
| **Tokens to success** | Mean / median total tokens among successful tasks; separately report mean tokens among failures (capped at give-up) |

Optional secondary metrics:

- Mean repair rounds among successes
- Failure mix: parse / resolve / runtime / wrong stdout / timeout
- Diagnostic usefulness: fraction of failed attempts whose next repair succeeds

---

## 4. Suite S (v0)

Reuse the existing density tasks. Do **not** invent a second task set until v0
has results.

### 4.1 Full suite (13 tasks)

| ID | Task dir | Intent exercised |
| --- | --- | --- |
| 01 | `bench/tasks/01-hello` | values, strings, `print`, `string()` |
| 02 | `bench/tasks/02-arithmetic` | `:=`, arithmetic, assignment |
| 03 | `bench/tasks/03-conditions` | `if` / `elif` / `else`, comparisons |
| 04 | `bench/tasks/04-loops` | `while`, `repeat` |
| 05 | `bench/tasks/05-arrays-each` | arrays, `.each` |
| 06 | `bench/tasks/06-where-sum` | `.where`, `.sum`, `assert` |
| 07 | `bench/tasks/07-functions` | `fn`, `return`, calls |
| 08 | `bench/tasks/08-strings-assert` | strings, `assert` |
| 09 | `bench/tasks/09-grade-report` | combined control flow |
| 10 | `bench/tasks/10-fibonacci` | functions + loops |
| 11 | `bench/tasks/11-leaderboard` | larger combined program |
| 12 | `bench/tasks/12-revenue` | aggregation-style logic |
| 13 | `bench/tasks/13-temperatures` | filtering / reporting |

### 4.2 Smoke suite (recommended first run)

Run these five before the full 13:

`01-hello`, `02-arithmetic`, `03-conditions`, `06-where-sum`, `07-functions`

### 4.3 Success criterion

Identical to [`bench/harness/check_stdout.py`](bench/harness/check_stdout.py):

- Process exit code 0
- stdout **exactly** equals `bench/tasks/<id>/expected.txt`  
  (including newlines; no extra trailing blank line unless expected)

Wrong stdout with exit 0 is a failure (`wrong_output`).  
Non-zero exit is a failure (`diag_error`), with stderr/JSON captured for repair.

### 4.4 Languages (v0)

| Language | Runner | Notes |
| --- | --- | --- |
| VOL | `vol run <file.vol>` or `go run ./cmd/vol run <file.vol>` | Primary subject |
| Go | `go run main.go` | Baseline |

Rust / Zig are optional later. Density already covers them; generate/repair v0
optimizes for one clean baseline.

---

## 5. Model and decoding controls

Record for every published table:

| Field | Rule |
| --- | --- |
| Model id | Exact API / local id (e.g. `gpt-4.1-mini`, `claude-…`) |
| Temperature | Fixed; recommend `0` or `0.2` for v0 |
| Top-p / other | Fixed or “API default” (state which) |
| Max output tokens | Cap high enough for largest smoke task; same for all langs |
| Seed | Set when the API supports it |
| Date | ISO date of run |
| N replicates | ≥ 3 independent runs per (task, language); report mean ± spread |

Do not change prompts, temperature, or K between languages inside one table.

---

## 6. Prompt protocol

### 6.1 Structure (same for every language)

Each request has three parts:

1. **Role** — short system text: write a complete program that prints exactly the
   required output; no markdown fences unless the harness strips them.
2. **Language card** — compact reference for *that* language only (see §6.2).
3. **Task card** — natural-language description of behavior + the exact expected
   stdout (from `expected.txt`).

Never give VOL the full `SPEC.md` while Go gets nothing. Matched card length
budgets beat “dump the manual.”

### 6.2 Language cards (budget)

Target **≤ ~400 tokens** each (measure with the same tokenizer used for scoring,
or the API’s tokenizer). Cards must only describe **currently Supported** VOL
features from [`SPEC.md`](SPEC.md) / README “What Actually Works.”

VOL card must include at least:

- no `main`; top-level statements run
- `:=` / `=` / opt-in `const`
- `if` / `elif` / `else` (statement); `? :` for values
- `repeat`, `while`, `.each`
- arrays, `.len`, `.where(_ …)`, `.sum`
- `fn` / `return`; missing return is `nothing` (do not assign it)
- `print`, `string()`, `assert`
- integer overflow traps

Go card: equivalently dense stdlib-only reminders (`fmt`, slices, `for`), not a
tour of Effective Go.

Store frozen cards under:

```text
bench/llm/cards/vol_v0.md
bench/llm/cards/go_v0.md
```

Bump the version suffix when the card changes; never silently edit a card used
in a published result table.

### 6.3 Task cards

For each task, a frozen prompt under:

```text
bench/llm/tasks/<id>/prompt.md
```

Must include:

- Goal in plain language (what to compute / print)
- Exact expected stdout in a fenced `text` block (copy of `expected.txt`)
- Constraints: single file; no network; no stdin unless the task requires it
- Language placeholder: `Write the program in {{LANG}}.`

Do **not** paste the reference `main.vol` / `main.go` into the prompt. Those
files are oracles for density and for humans writing task cards—not for the model.

### 6.4 Repair turn

On failure, send one additional user message (same conversation or explicit
transcript) containing:

1. The previous source
2. The tool result: exit code + stderr (prefer `vol --json run` for VOL so the
   model sees `code`, `message`, `fix`)
3. Instruction: return a full corrected program, not a diff

Max repair rounds **K = 2** for v0 (attempts 0..2 = 3 total generations).  
After that: mark `gave_up` and stop spending tokens on that task replicate.

---

## 7. Harness

Layout:

```text
bench/
  llm/
    cards/           # frozen language cards
    tasks/           # frozen task prompts
    harness/
      run_generate_repair.py
    results/         # JSONL + summary markdown
```

Dry-run (validates runners + summary; uses reference solutions; **not** a benchmark):

```text
cd bench && uv run python llm/harness/run_generate_repair.py --dry-run
# or: make llm-dry
```

### Live with Ollama (recommended for local work)

Requires [`ollama`](https://ollama.com) running (`ollama serve`) and a pulled model.

```text
cd bench
# default provider is ollama; default model qwen3:4b-instruct
uv run python llm/harness/run_generate_repair.py --provider ollama --model qwen3:4b-instruct
# or:
make llm-ollama
make llm-ollama MODEL=llama3.1:8b
```

Optional env: `OLLAMA_HOST` / `OLLAMA_BASE_URL` (default `http://localhost:11434/v1`),
`OLLAMA_MODEL`. No API key required.

### Live with OpenAI-compatible cloud

```text
cd bench
export OPENAI_API_KEY=…
uv run python llm/harness/run_generate_repair.py --provider openai --suite smoke \
  --model gpt-4.1-mini --temperature 0 --replicates 3 --max-repairs 2
```

Each replicate writes a JSONL record approximately:

```json
{
  "suite": "smoke_v0",
  "task": "06-where-sum",
  "language": "vol",
  "model": "…",
  "temperature": 0,
  "replicate": 1,
  "attempt": 0,
  "prompt_tokens": 0,
  "completion_tokens": 0,
  "exit_code": 1,
  "outcome": "diag_error",
  "diagnostic_code": "R022",
  "stdout": "",
  "stderr": "…"
}
```

Outcomes: `success` | `diag_error` | `wrong_output` | `timeout` | `gave_up` |
`extract_error` (could not recover source from model output).

Strip markdown fences if present; if multiple files appear, take the largest
code block or fail `extract_error`.

---

## 8. Result tables (publish format)

### 8.1 Per-language summary

| Language | First-try % | Success @ K % | Median tokens (success) | Mean tokens (all) | N |
| --- | --- | --- | --- | --- | --- |
| VOL | | | | | |
| Go | | | | | |

### 8.2 Headlines to put in README later

Only after a frozen run:

- VOL vs Go success @ K (absolute points)
- VOL vs Go total-token ratio on successes (and on all attempts)
- Clear disclaimer: model, date, suite, K, card versions

Until then, README must continue to say generate/repair is **not** measured.

---

## 9. Anti-gaming rules

1. **No oracle leakage** — reference solutions in `bench/tasks/*/vol` are not in prompts.
2. **No post-hoc card tuning on the reported suite** — tune on a holdout or accept the loss.
3. **Frozen artifacts** — cards, prompts, model id, K, temperature committed with results.
4. **Matched context** — language cards within ~20% token budget of each other.
5. **Supported surface only** — do not score tasks that need Planned syntax.
6. **Same runner timeouts** for every language (recommend 5s CPU per attempt for v0).
7. **Report failures** — a VOL loss is a valid scientific result.

---

## 10. Relation to MPT

**Meaning per Token (MPT)** remains undefined as a single scalar.

Until a better definition exists, treat §3 `workflow_efficiency` and the three
headline rates as the operational substitute. Do not invent an MPT formula that
cannot be recomputed from committed JSONL.

---

## 11. Implementation checklist

- [x] Add frozen `bench/llm/cards/vol_v0.md` and `go_v0.md`
- [x] Add frozen `bench/llm/tasks/<id>/prompt.md` for smoke suite
      (`01-hello`, `02-arithmetic`, `03-conditions`, `06-where-sum`, `07-functions`)
- [ ] Add prompts for remaining full-suite tasks (04, 05, 08–13)
- [x] Implement `run_generate_repair.py` (OpenAI-compatible API + local runners + `--dry-run`)
- [x] Wire VOL failures through `vol --json run` for repair context
- [ ] Run smoke suite, N ≥ 3, VOL vs Go, one model (live API)
- [ ] Commit JSONL + summary markdown under `bench/llm/results/`
- [ ] Update README with numbers and disclaimers
- [ ] Only then expand models / Rust / full suite

---

## 12. Change process

When changing the protocol:

1. Bump a protocol version note at the top of this file.
2. Invalidate or re-run results that depended on old cards/prompts.
3. Keep density (`bench/results/density*.md`) separate—never mix tables.

A generate/repair number is not official until transcripts and the summary table
are committed next to the frozen cards that produced them.
