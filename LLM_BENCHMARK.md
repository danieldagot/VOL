# VOL LLM Workflow Benchmark (Protocol v1.1)

> Status: **protocol v1.1 harness implemented; Gemini `core_v2` VOL vs Python published; Go/`core_v1` historical**
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
| Relative cost vs Python (default baseline) and optional Go | “Meaning per Token” as a single magic score |

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
| **First-try success** | Share of tasks where attempt 0 stdout matches `expected.txt` (for diagnostic repair, attempt 0 already includes seed diagnostics) |
| **Success @ K** | Share of tasks solved within ≤ K repairs (K default = 2) |
| **Tokens to success** | Mean / median total tokens among successful tasks; separately report mean tokens among failures (capped at give-up) |

Always split **prompt** vs **completion** tokens in published summaries. Also report:

| Accounting | Definition |
| --- | --- |
| **Cold** | Sum of provider `prompt_tokens` + `completion_tokens` (language card re-sent every request) |
| **Warm** | Cold with estimated language-card tokens subtracted from each prompt (`cl100k_base` estimate of the frozen card). Models amortized / cached-card cost. |

Warm is an accounting view, not a separate API mode. Provider usage fields remain the source of cold totals.

Optional secondary metrics:

- Mean repair rounds among successes
- Failure mix: parse / resolve / runtime / wrong stdout / timeout
- Diagnostic usefulness: fraction of failed attempts whose next repair succeeds
- VOL vs baseline deltas (Python by default; Go optional) for prompt, completion, cold total, and warm total

---

## 4. Suite S (v1.1 / `core_v2`)

The workflow benchmark is deliberately smaller than the 13-task static-density
set. Fewer, richer tasks allow at least three replicates without spending most
of the budget on redundant syntax exercises.

Harness suite ids: smoke → `smoke_v1`; core → `core_v2` (protocol v1.1).
Historical Gemini results under `core_v1` used pre-diagnostic repair prompts and
must not be mixed into v1.1 tables.

### 4.1 Smoke suite (2 tasks)

| ID | Workflow | Intent exercised |
| --- | --- | --- |
| `01-hello` | generation | runner and basic output sanity check |
| `07-functions` | generation | functions, calls, strings, integer conversion |

Smoke validates wiring. It is not sufficient evidence for a language claim.

### 4.2 Core suite (5 tasks)

| ID | Workflow | Intent exercised |
| --- | --- | --- |
| `05-arrays-each` | generation | collection mutation, iteration, filtering |
| `08-strings-assert` | diagnostic repair | fix seeded failure using real runner diagnostics |
| `10-fibonacci` | generation | loop, mutable state, iterative computation |
| `11-leaderboard` | modification | preserve and extend a working program |
| `13-temperatures` | generation | aggregation, categories, multiple invariants |

Run each `(task, language)` at least three times. Report results separately by
workflow kind as well as in aggregate.

### 4.3 Success criterion

A task succeeds only when all checks pass:

- process exit code is 0;
- stdout exactly equals `bench/tasks/<id>/expected.txt`;
- every language-specific source constraint in `task.json` matches.

Source constraints enforce explicit requirements such as using a loop,
preserving starter collections, or retaining an invariant check. They must be
minimal and declared before a run. Failure is `source_check_failed` and its
message is eligible for a repair turn.

### 4.4 Languages (v1.1)

| Language | Runner | Notes |
| --- | --- | --- |
| VOL | `vol run <file.vol>` or `go run ./cmd/vol run <file.vol>` | Primary subject |
| Python | `python main.py` (CPython 3.11+) | **Default baseline** — interpreted peer for the current prototype |
| Go | `go run main.go` | Optional systems-language baseline (`--langs vol,go` or `vol,python,go`) |

While VOL is an interpreter prototype, prefer Python as the workflow baseline.
Go remains useful for a compiled-language comparison and for continuity with
historical `core_v1` / early `core_v2` tables. Rust / Zig stay optional; density
already covers them.

---

## 5. Model and decoding controls

Record for every published table:

| Field | Rule |
| --- | --- |
| Model id | Exact API / local id (e.g. `gpt-4.1-mini`, `claude-…`) |
| Temperature | Fixed; recommend `0` or `0.2` for v1 |
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

Never give VOL the full `SPEC.md` while a baseline gets nothing. Matched card
length budgets beat “dump the manual.”

### 6.2 Language cards (budget)

Target **≤ ~400 tokens** each (measure with the same tokenizer used for scoring,
or the API’s tokenizer). Cards must only describe **currently Supported** VOL
features from [`SPEC.md`](SPEC.md) / README “What Actually Works.”

VOL `core_v2` card (`vol_v1`) must include at least:

- no `main`; top-level statements run
- `:=` / `=` / opt-in `const`
- `if` / `elif` / `else` (statement); `? :` for values
- `repeat`, `while`, `.each`
- arrays, `.len`, `.where(_ …)` / `.count(_ …)`, `.sum()`
- `fn` / `return`; missing return is `nothing` (do not assign it)
- `print`, `string()`, `assert`
- integer overflow traps

It is a **task card** for `core_v2`, bound to SF-1 Supported syntax, not a full
SF-1 product tour. Omit Option/Result/structs/`import` until a suite exercises
them (a fuller card is a separate, deliberate revision).

Python card: equivalently dense stdlib-only reminders (`print`, lists, `for`,
`assert`), not a tour of the standard library.

Go card (optional baseline): equivalently dense stdlib-only reminders (`fmt`,
slices, `for`), not a tour of Effective Go.

Store frozen cards under:

```text
bench/llm/cards/vol_v1.md      # core_v2 task card (SF-1-bound; harness default)
bench/llm/cards/vol_v0.md      # SF-0 (historical core_v2 tables; do not mix freeze IDs)
bench/llm/cards/python_v0.md   # matched Python baseline (default)
bench/llm/cards/go_v0.md       # matched Go baseline (optional)
```

Cards are **bound to** a product surface freeze ([`SPEC.md`](SPEC.md) §0) but
need not enumerate every Supported form. SF-0 → `vol_v0`; SF-1 → `vol_v1`
(`core_v2` subset). Do not keep intermediate draft cards for unfinished
surfaces. Bump the VOL card version suffix with a freeze bump (next: SF-2 →
`vol_v2.md`) or a deliberate card-only revision that does not add language
features. Never silently edit a card used in a published result table without
republishing; published summaries must name the freeze and card versions.
Source checks must accept every form the card teaches for the same intent
(e.g. `.count` and `.where` for counting).

### 6.3 Task artifacts

For each task, store a frozen prompt and machine-readable metadata:

```text
bench/llm/tasks/<id>/prompt.md
bench/llm/tasks/<id>/task.json
```

Must include:

- Goal in plain language (what to compute / print)
- Exact expected stdout in a fenced `text` block (copy of `expected.txt`)
- Constraints: single file; no network; no stdin unless the task requires it
- Language placeholder: `Write the program in {{LANG}}.`

`task.json` records the workflow kind (`generation`, `repair`, or
`modification`) and minimal language-specific source checks. Modification
tasks may embed working source via `{{STARTER}}`. Repair tasks use
`"seed": "starter"` and must **not** paste the broken starter into the prompt
without diagnostics—the harness attaches it after a real failed run.

Do **not** paste the reference `main.vol` / `main.py` / `main.go` into the
prompt. Those files are oracles for density and for humans writing task
cards—not for the model.

### 6.4 Repair turn (diagnostic)

**Seeded diagnostic repair** (`kind: repair`, `seed: starter`):

1. Harness runs `starter.vol` / `starter.py` / `starter.go` **before** any model call.
2. The starter **must fail** (exit ≠ 0, wrong stdout, or source-check failure).
3. Attempt 0’s user message includes the language card, task card, failed
   source, exit code, and stderr (`vol --json` for VOL so the model sees
   `code`, `message`, `fix`).
4. The model returns a full corrected program, not a diff.

This measures whether diagnostics help. Do **not** ask the model to rewrite an
obviously broken starter in the clear without first attaching runner output.

**Follow-up repairs** (generation / modification failures, or repair attempt > 0):
send one additional user message containing previous source, exit code, and
stderr, with the same “full corrected program” instruction.

Max repair rounds **K = 2** (attempts 0..2 = 3 total model generations).
After that, preserve the final failure outcome and stop spending tokens on that
task replicate.

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
# restrict a run to selected tasks and cap each model request (seconds)
make llm-ollama MODEL=llama3.1:8b TASKS=06-where-sum REQUEST_TIMEOUT=30
```

Optional env: `OLLAMA_HOST` / `OLLAMA_BASE_URL` (default `http://localhost:11434/v1`),
`OLLAMA_MODEL`. No API key required.

Use `--tasks ID[,ID...]` to select tasks explicitly. Use
`--request-timeout SECONDS` to override the provider's per-model-request
timeout; this is distinct from the maximum repair count.

### Live with OpenAI-compatible cloud

```text
cd bench
export OPENAI_API_KEY=…
uv run python llm/harness/run_generate_repair.py --provider openai --suite smoke \
  --model gpt-4.1-mini --temperature 0 --replicates 3 --max-repairs 2
```

### Live with Google Gemini

Copy `.env.example` to `.env`, then set your key without committing it:

```text
GEMINI_API_KEY=your-gemini-api-key
GEMINI_MODEL=gemini-3.6-flash
```

Run Gemini through Google's OpenAI-compatible endpoint:

```text
cd bench
uv run python llm/harness/run_generate_repair.py --provider gemini --suite core \
  --temperature 0 --replicates 3 --max-repairs 2
```

The harness loads `.env` from the repository root or `bench/.env`. Shell
variables take precedence. `GEMINI_BASE_URL` is optional.

HTTP 429 and transient 5xx responses are retried with a bounded provider-aware
delay. While a run is active, records use a `.partial.jsonl` suffix; the harness
renames the file to `.jsonl` only after every task-replicate finishes.

Each replicate writes a JSONL record approximately:

```json
{
  "suite": "core_v2",
  "task": "08-strings-assert",
  "task_kind": "repair",
  "language": "vol",
  "model": "…",
  "temperature": 0,
  "replicate": 1,
  "attempt": 0,
  "prompt_tokens": 0,
  "completion_tokens": 0,
  "card_tokens_est": 0,
  "prompt_tokens_warm": 0,
  "tokens_cold": 0,
  "tokens_warm": 0,
  "repair_seeded": true,
  "exit_code": 1,
  "outcome": "diag_error",
  "diagnostic_code": "R022",
  "stdout": "",
  "stderr": "…"
}
```

`repair_seeded` is true on attempt 0 of diagnostic-repair tasks (seed failure
already attached). Outcomes: `success` | `diag_error` | `wrong_output` |
`source_check_failed` | `timeout` | `extract_error` (could not recover source
from model output).

Strip markdown fences if present; if multiple files appear, take the largest
code block or fail `extract_error`.

---

## 8. Result tables (publish format)

### 8.1 Per-language summary

The harness markdown includes cold totals, prompt vs completion means, warm
(card-amortized) totals, and VOL vs each baseline deltas. Minimum headline table:

| Language | First-try % | Success @ K % | Mean cold | Mean prompt | Mean completion | Mean warm | N |
| --- | --- | --- | --- | --- | --- | --- | --- |
| VOL | | | | | | | |
| Python | | | | | | | |

### 8.2 Headlines to put in README later

Only after a frozen **`core_v2`** run that includes the languages being claimed:

- VOL vs Python success @ K (absolute points); optional VOL vs Go
- VOL vs baseline cold and warm total-token ratios
- Prompt vs completion deltas (so card cost is visible)
- Clear disclaimer: model, date, suite, protocol, K, card versions

Keep historical VOL-vs-Go `core_v1` / early `core_v2` results labeled; do not
mix them into Python-baseline tables without naming both baselines.

---

## 9. Anti-gaming rules

1. **No oracle leakage** — reference solutions in `bench/tasks/*/vol` are not in prompts.
2. **No post-hoc card tuning on the reported suite** — tune on a holdout or accept the loss.
3. **Frozen artifacts** — cards, prompts, model id, K, temperature committed with results.
4. **Matched context** — language cards within ~20% token budget of each other.
5. **Supported surface only** — do not score tasks that need Planned syntax.
6. **Same runner timeouts** for every language (recommend 5s CPU per attempt for v1).
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
- [x] Add frozen `bench/llm/cards/python_v0.md` (default interpreter baseline)
- [x] Add the 2-task smoke suite (`01-hello`, `07-functions`)
- [x] Add the 5-task core suite with generation, repair, and modification tasks
- [x] Add `task.json` source constraints and language-specific starter programs
- [x] Implement `run_generate_repair.py` (OpenAI-compatible API + local runners + `--dry-run`)
- [x] Wire VOL failures through `vol --json run` for repair context
- [x] Run core suite, N ≥ 3, VOL vs Go, one model (live API) — historical `core_v1`
- [x] Store JSONL + summary markdown under `bench/llm/results/`
- [x] Update README with numbers and disclaimers
- [x] Report prompt vs completion and cold vs warm (card-amortized) in summaries
- [x] Seeded diagnostic repair: failing starter runs before attempt 0 (`core_v2`)
- [x] Publish a frozen `core_v2` live run (Gemini) — VOL vs Go
- [x] Bind cards to Surface Freeze SF-0 (`SPEC.md` §0; `vol_v0` / `python_v0` / `go_v0`)
- [x] Publish frozen `core_v2` live run with default `--langs vol,python` (Gemini)
- [x] Surface Freeze SF-1; `vol_v1.md` = `core_v2` task card (SF-1-bound subset)
- [ ] Publish ≥1 other model on `core_v2` before further syntax optimization
- [x] Re-run / publish `core_v2` against `vol_v1`
      (`20260808-041440`; Python via `--baseline-jsonl` from `…040028`;
      `…040028` VOL numbers poisoned by `.where`-only source checks — superseded)
- [x] `11-leaderboard` / `13-temperatures` source checks accept `.count` or `.where`
- [x] Harness `--baseline-jsonl` for VOL-only republish with frozen Python rows

---

## 12. Change process

When changing the protocol:

1. Bump a protocol version note at the top of this file.
2. Invalidate or re-run results that depended on old cards/prompts.
3. Keep density (`bench/results/density*.md`) separate—never mix tables.

**Baseline re-run policy:** do not re-run Python (or Go) on every VOL iteration.
Use `--langs vol` when only the VOL card, diagnostics, source checks, or
Supported surface changed. Re-run the baseline language only when its card, the
task suite, protocol, model, temperature, or scoring changed. When publishing a
VOL-only run, merge frozen baseline rows with `--baseline-jsonl <prior.jsonl>`
(same model and suite) so the new JSONL/summary is self-contained and names the
artifact.

A generate/repair number is not official until transcripts and the summary table
are committed next to the frozen cards that produced them.
