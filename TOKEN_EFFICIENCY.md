# VOL Token Efficiency — Working Preset

> Continuation doc for maximizing VOL token efficiency.
> Last updated: 2026-08-08 (SF-3 `intent_v1` + `@std` tasks; juice loop stopped).
> Related: [`bench/`](bench/README.md), [`LLM_BENCHMARK.md`](LLM_BENCHMARK.md),
> [`IDEAS.md`](IDEAS.md), [`SPEC.md`](SPEC.md) §0 / SF-3.

## Goal

Maximize **token efficiency** on two tracks (keep them distinct):

1. **Source density** (`bench/` `count_tokens.py`) — hand-written equivalent
   programs; fewer tokens under named tokenizers.
2. **Workflow efficiency** (`LLM_BENCHMARK.md`) — task success / total tokens
   (card + generate + diagnostics + repair). Cold vs warm matter.

Do **not** cite density ratios as LLM workflow proof. Do **not** add features
only to juice the density table if they grow the language card and raise repair
cost.

**Constraint for this preset:** prefer **radical dynamics** (same surface, richer
omit/coercion/print/collection rules) over **radical syntax** (`|>`, f-string
glyphs, `=>`, dual-return sugar, etc.).

---

## Current scoreboard (source density)

Suite: **21 tasks** × VOL / Python / Go / Rust / Zig (SF-2 densified VOL + SF-3
`@std` stdlib tier), tiers:

| Tier | Tasks | Use when judging… |
| --- | --- | --- |
| **compression** | 06, 14, 15, 16 | collection intent (prefer this for density claims) |
| **labeled** | 09, 11, 12, 13 | report glue (multi-arg `print` / coercion) |
| **parity** | 01–05, 07–08, 10 | control near-Python floor |
| **stdlib** | 17–21 | SF-3 `@std` vs peer stdlibs (not collection density) |

Current medians (`cl100k_base`, 21 tasks): all VOL/Python **0.829**;
compression **0.764**; labeled **0.651**; stdlib **0.929** (`dict("k", v, …)`
closed the json>Python gap; env/strings remain near parity). See
[`bench/results/density.md`](bench/results/density.md). Do not mix stdlib
import surface with compression-tier claims.

**Target “50% fewer than Python” (ratio 0.5) on median (all) is still hard**
without suite skew; prefer honest compression gains, stdlib fairness notes, and
workflow wins over labeled-only golf.

---

## LLM workflow (separate track)

Published Gemini `intent_v1` (7 tasks incl. `@std` strings/json; `vol_v3` / SF-3
vs Python, `gemini-3.5-flash-lite`):

| Metric | VOL | Python |
| --- | ---: | ---: |
| First-try / success @ K | 100% / 100% | 100% / 100% |
| Mean cold | 715.7 | 733.6 (−2.4%) |
| Mean warm | 365.7 | 397.6 (−8.0%) |
| Mean completion | 78.4 | 113.1 (−30.7%) |
| Card est. (`cl100k_base`) | 350 | 336 |

Summary:
[`bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-063341.md`](bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-063341.md).

Historical 5-task SF-3 (no `@std` tasks, `…061355`): FT 100%, cold 685.3, warm
363.3, card 322. Historical SF-2 `intent_v1` (`vol_v2`, `…051437`): FT 100%,
cold 685.5, warm 369.5. Historical SF-1 `intent_v1` (`vol_v1`): first-try
**80%** (`.count()` arity) — see `…045310.md`. Historical `core_v2` SF-1: 100%
success with ~+11% cold / ~−3% warm vs Python.

Still need ≥1 other model before treating deltas as stable.

---

## Iteration log (this loop)

| Step | What failed / gap | Change | Density | Workflow | Keep? |
| --- | --- | --- | --- | --- | --- |
| Baseline SF-2 | SF-1 had 80% FT on `.count()` | (already shipped SF-2 surface) | all 0.805 / comp 0.799 / lab 0.651 | FT 93.3%, cold 908.7, warm 426.6 (`…051110`) | baseline |
| Footgun | `16-map-filter` rep3: `.sum(pred)` → S003 | S003 Fix: use `.where(condition).sum()` | — | (repair aid) | **keep** |
| Card v1 | FT not 100%; card 452 | Clarify `.sum()` 0-arg; cut fluff → 421 | — | FT 100%, cold 808.3, warm 387.3 (superseded intermediate; artifact removed) | superseded |
| Card v2 | Cold still +card tax | Lean `vol_v2` → **316** (≤ Python 336) | — | FT 100%, cold **685.5**, warm **369.5** (`…051437`) | **keep** (historical) |
| Densify | 14/16 repeat pipelines | Bind shared `.where` / `.map` | all **0.804**, comp **0.764**, lab 0.651 | n/a (benches only) | **keep** |
| Stdlib `?` | stdlib VOL/Python **1.188** (`fn main` tax) | Top-level Result `?` (err → `R049`); densify tasks 19–21 | stdlib **0.936**; all 0.829 | n/a (density) | **keep** |
| Dict pairs | json task **1.071** (empty `dict()` + assign) | `dict("k", v, …)` + densify 20-json / examples; lean `vol_v3` **322** | stdlib **0.929**; json **0.929** | FT 100%, cold **685.3**, warm **363.3** (`…061355`, 5-task) | **keep** |
| Suite gap | `intent_v1` claimed SF-3 but never ran `@std` | Add LLM tasks `17`/`20` + source_checks | — | FT **71%** / success 86% (`…062641`) | diagnose |
| Footgun | `contains` / `strings.trim` / bare `parse` / `print dump` → `ok(…)` | S002 alias Fixes; R003 Result unwrap Fix; card `print dump(d)?`; wrong_output expected/got in harness | — | FT **100%**, cold **715.7**, warm **365.7** (`…063341`) | **keep** |

**Rejected / not pursued for juice:** `.sum(pred)` sugar (would be new surface),
shrinking `vol_v3` below ~350 while keeping `@std` first-try (quota blocked a
confirming re-run after a −1-token trim), `|>` / `=>` / enums / dual-return /
`{k:v}` literals (SF-4+), unfair Python padding, full `@std` API tour in the
task card (tour belongs in SPEC/examples).

---

## SF-3 exhausted — rationale

Stop criteria met for this model + expanded `intent_v1` under SF-3:

1. VOL first-try **100%** and success @ K **100%** (7 tasks incl. `@std`).
2. Cold and warm means **below Python** on the expanded suite (`…063341`); warm
   also below the historical 5-task SF-3 warm mean once card tax is amortized.
3. Remaining SF-3-legal levers were tried or ruled out:
   - **Diagnostics:** S002 std-alias Fixes + R003 Result-index Fix + harness
     wrong_output expected/got — cleared the `@std` footguns seen in `…062641`.
   - **Card:** `vol_v3` **350** (+14 vs Python) holds 100% first-try on `@std`
     reminders; further cuts risk FT and hit free-tier quota on re-measure.
   - **Bench densify:** stdlib density tasks already dense; env/strings near floor.
   - **In-freeze semantics:** `{k:v}` literals / changing Result `print` display
     are out of SF-3 spirit for table juice.
4. Further gains need **non-SF-3 language** work:
   - **`vol fmt`** rewriter (canonical `_` / bind style → smaller completions)
   - **Prompt caching / warm serving** (infrastructure, not language)
   - **Second model** confirmation
   - **SF-4+** only if measured necessary (`{k:v}`, `|>`, …) — not glyph theater

---

## Really measure how an LLM uses VOL

Density (`make count`) answers: “how short is good VOL a human already wrote?”  
It does **not** answer: “what does a model emit, break, and repair?”

Use [`LLM_BENCHMARK.md`](LLM_BENCHMARK.md) + `bench/llm/harness/run_generate_repair.py`.

### What to watch

| Signal | Meaning |
| --- | --- |
| **First-try %** | Model can write runnable VOL from card + task |
| **Success @ K** | Diagnostics + repair close the gap |
| **Completion tokens** | How verbose the model’s VOL is |
| **Cold total** | Full cost including re-sent language card |
| **Warm total** | Cost if the card is amortized / cached |
| **Failure mix** | parse / R0xx / wrong stdout / `source_check_failed` |
| **Idioms emitted** | `.count` vs `.where.len` vs loops — does the card stick? |

### Commands

```sh
cd bench

make check && make count

# Prefer intent_v1 for “how does an LLM use VOL?”
uv run python llm/harness/run_generate_repair.py \
  --provider gemini --suite intent --langs vol,python --replicates 3

# VOL-only iteration against frozen Python rows
uv run python llm/harness/run_generate_repair.py \
  --provider gemini --suite intent --langs vol --replicates 3 \
  --baseline-jsonl llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-063341.jsonl
```

Needs `GEMINI_API_KEY` (or another provider) in env / `.env`.

### Decision rule

> Ship language or card changes only if **first-try % stays high** and
> **cold or warm tokens to success** drop vs the previous published JSONL —
> not because hand-written density looked better.

---

## Strategy: radical dynamics (not radical syntax)

### D1–D3 — SF-2 density dynamics ✓

- string `+` coercion (`"n=" + 7`); `1 + "a"` stays `R013`
- multi-arg `print` (space-join)
- `.count()` ≡ `.len`; `.count(pred)` for filtered counts; `.sum()` remains 0-arg
  (filter with `.where` first — taught by card + S003 Fix)

### D4 — Suite dynamics (honest scoreboard)

Do not fake 50% by padding Python or stripping VOL assert messages.

### D5 — Dict pairs (SF-3) ✓

- `dict("k", v, …)` alternating string keys/values; `dict()` still empty
- Odd arity → `R018` with Fix; non-string key → `R045`
- `{k:v}` literals stay SF-4+

### D6 — Canon + tooling (next, outside SF-3 language juice)

- `vol fmt`: prefer `_` over `fn` in collection ops; one bind/loop/value style
- Move remaining “don’t write X” from card → diagnostic `Fix` when new footguns appear
- Prompt caching for cold totals in real deployments

---

## Explicit non-goals for density juice

Do **not** pursue these *for density/workflow juice* as SF-3 language sugar
(SF-3 is `@std` usability; syntax sugar is SF-4+):

- `|>` pipelines, `?T` Option sugar, `=>` arrows, dual-return sugar, `{k:v}`
- Parallel / ownership inference
- Renaming keywords for 1-token shaves that grow the card
- `.sum(pred)` unless a multi-model footgun proves card+Fix insufficient
- Full `@std` API enumeration in the task card

---

## Ranked next actions

### A. Language dynamics (SF-2 / SF-3) ✓ exhausted for this loop

- [x] String `+` coercion; multi-arg `print`; `.count()` length; densified benches
- [x] S003 Fix for `.sum(pred)` → `.where(condition).sum()`
- [x] Lean `vol_v2` / `vol_v3` cards (≤ Python) with explicit `.sum()` 0-arg rule
- [x] Bind shared pipelines in compression tasks 14/16
- [x] `dict("k", v, …)` + densified stdlib json / examples; published `vol_v3` `intent_v1`

### B. Workflow / tooling (beyond SF-3 surface)

- [ ] Finish `vol fmt` rewriter with canonical table
- [ ] Publish ≥1 other model on `intent_v1` + `vol_v3`
- [ ] Optional: prompt-cache accounting experiments (warm ≈ production)

### C. Suite / hygiene

- [x] Compression + labeled + parity + stdlib tier reporting
- [x] Split density vs workflow; cold vs warm in README
- [ ] Never claim 50% vs Python until median (compression) ≤ 0.50 on a
      documented suite (not labeled-only)

---

## Idioms (hand-written / formatter target)

```vol
print a + b
print cond ? "yes" : "no"
repeat 8 { … }
fn square(n) n * n
scores.count(_ >= 90)
g := nums.where(_ > 5)
print g.count()
print g.sum()
print "Sum:", total
print dump(dict("n", 3))?
```

---

## One-line north star

**Omit ceremony for display and aggregation; keep one familiar spelling per intent;
measure density and workflow separately; beat Python on intent-heavy programs,
not on hello-world ties.**
