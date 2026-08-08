# VOL Token Efficiency — Working Preset

> Continuation doc for maximizing VOL token efficiency.
> Last updated: 2026-08-08 (SF-2 exhausted for this loop).
> Related: [`bench/`](bench/README.md), [`LLM_BENCHMARK.md`](LLM_BENCHMARK.md),
> [`IDEAS.md`](IDEAS.md), [`SPEC.md`](SPEC.md) §0 / SF-2.

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

Suite: **16 tasks** × VOL / Python / Go / Rust / Zig (SF-2 densified VOL), tiers:

| Tier | Tasks | Use when judging… |
| --- | --- | --- |
| **compression** | 06, 14, 15, 16 | collection intent (prefer this for density claims) |
| **labeled** | 09, 11, 12, 13 | report glue (multi-arg `print` / coercion) |
| **parity** | 01–05, 07–08, 10 | control near-Python floor |

Post-loop medians (`cl100k_base`): all VOL/Python **0.804** (~20% fewer);
compression **0.764** (~24% fewer); labeled **0.651** (~35% fewer). See
[`bench/results/density.md`](bench/results/density.md).

**Target “50% fewer than Python” (ratio 0.5) on median (all) is still hard**
without suite skew or SF-3 surface; prefer honest compression gains and
workflow wins over labeled-only golf.

---

## LLM workflow (separate track)

Published Gemini `intent_v1` (`vol_v2` / SF-2 vs Python, `gemini-3.5-flash-lite`):

| Metric | VOL | Python |
| --- | ---: | ---: |
| First-try / success @ K | 100% / 100% | 100% / 100% |
| Mean cold | 685.5 | 761.3 (−10.0%) |
| Mean warm | 369.5 | 425.3 (−13.1%) |
| Mean completion | 85.7 | 134.5 (−36.3%) |
| Card est. (`cl100k_base`) | 316 | 336 |

Summary:
[`bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-051437.md`](bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-051437.md).

Historical SF-1 `intent_v1` (`vol_v1`): first-try **80%** (all misses: zero-arg
`.count()` before SF-2) — see `…045310.md`. Historical `core_v2` SF-1: 100%
success with ~+11% cold / ~−3% warm vs Python.

Still need ≥1 other model before treating deltas as stable.

---

## Iteration log (this loop)

| Step | What failed / gap | Change | Density | Workflow | Keep? |
| --- | --- | --- | --- | --- | --- |
| Baseline SF-2 | SF-1 had 80% FT on `.count()` | (already shipped SF-2 surface) | all 0.805 / comp 0.799 / lab 0.651 | FT 93.3%, cold 908.7, warm 426.6 (`…051110`) | baseline |
| Footgun | `16-map-filter` rep3: `.sum(pred)` → S003 | S003 Fix: use `.where(condition).sum()` | — | (repair aid) | **keep** |
| Card v1 | FT not 100%; card 452 | Clarify `.sum()` 0-arg; cut fluff → 421 | — | FT 100%, cold 808.3, warm 387.3 (`…051408`) | superseded |
| Card v2 | Cold still +card tax | Lean `vol_v2` → **316** (≤ Python 336) | — | FT 100%, cold **685.5**, warm **369.5** (`…051437`) | **keep** |
| Densify | 14/16 repeat pipelines | Bind shared `.where` / `.map` | all **0.804**, comp **0.764**, lab 0.651 | n/a (benches only) | **keep** |

**Rejected / not pursued for juice:** `.sum(pred)` sugar (would be new surface /
SF-3), further card cuts below ~316 (risk first-try), `|>` / `=>` / enums /
dual-return, unfair Python padding.

---

## SF-2 exhausted — rationale

Stop criteria met for this model + `intent_v1`:

1. VOL first-try **100%** and success @ K **100%**.
2. Cold and warm means **below** prior SF-2 JSONL (`…051110` / `…051408`) and
   **below Python** on both cold and warm (not merely “card tax excuse”).
3. Remaining SF-2-legal levers were tried or ruled out:
   - **Diagnostics:** `.sum` arity Fix shipped; no other high-frequency repair
     footguns in the latest JSONL (all first-try).
   - **Card:** already **≤ Python** while holding 100% first-try; further cuts
     are predicted to hurt first-try more than they save cold tokens.
   - **Bench densify:** labeled/compression already use multi-arg `print`,
     coercion, `.count()` / `.count(pred)`, and shared binds; further shaves are
     micro-golf or need new collection sugar.
   - **In-freeze semantics:** accepting `.sum(pred)` would expand Supported
     call shapes beyond the SF-2 pin spirit → defer to SF-3 if ever measured
     necessary.
4. Further gains need **non-SF-2** work:
   - **`vol fmt`** rewriter (canonical `_` / bind style → smaller completions)
   - **Prompt caching / warm serving** (infrastructure, not language)
   - **Second model** confirmation
   - **SF-3** only if a measured footgun demands new surface (not for density theater)

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
  --baseline-jsonl llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-051110.jsonl
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

### D5 — Canon + tooling (next, outside SF-2 language juice)

- `vol fmt`: prefer `_` over `fn` in collection ops; one bind/loop/value style
- Move remaining “don’t write X” from card → diagnostic `Fix` when new footguns appear
- Prompt caching for cold totals in real deployments

---

## Explicit non-goals for density juice

Do **not** pursue these *for density/workflow juice* until a measured SF-3 need:

- `|>` pipelines, `?T` Option sugar, `=>` arrows, dual-return sugar
- Parallel / ownership inference
- Renaming keywords for 1-token shaves that grow the card
- `.sum(pred)` unless a multi-model footgun proves card+Fix insufficient

---

## Ranked next actions

### A. Language dynamics (SF-2) ✓ exhausted for this loop

- [x] String `+` coercion; multi-arg `print`; `.count()` length; densified benches
- [x] S003 Fix for `.sum(pred)` → `.where(condition).sum()`
- [x] Lean `vol_v2` card (316 ≤ Python 336) with explicit `.sum()` 0-arg rule
- [x] Bind shared pipelines in compression tasks 14/16

### B. Workflow / tooling (beyond SF-2 surface)

- [ ] Finish `vol fmt` rewriter with canonical table
- [ ] Publish ≥1 other model on `intent_v1` + `vol_v2`
- [ ] Optional: prompt-cache accounting experiments (warm ≈ production)

### C. Suite / hygiene

- [x] Compression + labeled + parity tier reporting
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
print "High score:", score
```

---

## One-line north star

**Omit ceremony for display and aggregation; keep one familiar spelling per intent;
measure density and workflow separately; beat Python on intent-heavy programs,
not on hello-world ties.**
