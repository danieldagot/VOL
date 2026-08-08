# VOL Token Efficiency — Working Preset

> Continuation doc for maximizing VOL token efficiency.
> Last updated: 2026-08-08.
> Related: [`bench/`](bench/README.md), [`LLM_BENCHMARK.md`](LLM_BENCHMARK.md),
> [`IDEAS.md`](IDEAS.md), [`SPEC.md`](SPEC.md) §0 / SF-1.

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

Suite: **16 tasks** × VOL / Python / Go / Rust / Zig, with tiers:

| Tier | Tasks | Use when judging… |
| --- | --- | --- |
| **compression** | 06, 14, 15, 16 | collection intent (prefer this for density claims) |
| **labeled** | 09, 11, 12, 13 | report glue / future print·coercion impact |
| **parity** | 01–05, 07–08, 10 | control near-Python floor |

Re-run `cd bench && make count` and read tier medians in
[`bench/results/density.md`](bench/results/density.md). Headline README %s use
**median (all)**; quote **median (compression)** when claiming semantic density.

**Target “50% fewer than Python” (ratio 0.5) on median (all) is not realistic**
without print/coercion dynamics + more compression tasks. Prefer improving
**median (compression)** first; do not juice labeled tasks alone and call it a
language win.

---

## LLM workflow (separate track)

Published Gemini `core_v2` (`vol_v1` vs Python): 100% first-try / success @ K;
completions ~18.7% smaller; **cold** ~+11% (card tax); **warm** ~−3.3%.
Card ~436 vs Python ~336.

For workflow ROI (no new language features):

1. Shrink `bench/llm/cards/vol_v1.md` toward ≤ Python card size.
2. Ship `vol fmt` with one canonical form per intent (`_` in collection ops).
3. Move “don’t write X” from card → diagnostic `fix`.
4. Re-measure after every card cut; keep first-try %.

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

Protocol always gives a **matched language card** (not full `SPEC.md`). That is
the fair “LLM without memorizing VOL” setup: task + compact card only.

### Commands

```sh
cd bench

# Wiring only (no API) — reference solutions must pass
make llm-dry

# Live smoke (2 tasks) — cheap sanity
uv run python llm/harness/run_generate_repair.py \
  --provider gemini --suite smoke --langs vol,python --replicates 1

# Prefer intent_v1 for “how does an LLM use VOL?” (no fib/hello toys)
uv run python llm/harness/run_generate_repair.py \
  --provider gemini --suite intent --langs vol,python --replicates 3

# Historical published table shape (includes parity — keep for continuity only)
uv run python llm/harness/run_generate_repair.py \
  --provider gemini --suite core --langs vol,python --replicates 3

# Dry-run intent reference solutions (no API)
make llm-dry-intent
```

Needs `GEMINI_API_KEY` (or another provider) in env / `.env`. Results land under
`bench/llm/results/`.

### Suites

| Suite | Freeze id | Use for |
| --- | --- | --- |
| `smoke` | `smoke_v1` | Wiring only |
| `core` | `core_v2` | Published continuity (includes fib / arrays parity) |
| **`intent`** | **`intent_v1`** | **Real language-use claims** |

`intent_v1` tasks:

| ID | Kind | Why it is here |
| --- | --- | --- |
| `06-where-sum` | generation | filter + aggregate |
| `14-pipeline-stats` | generation | count / where / map / sum pipeline |
| `16-map-filter` | generation | map then count/sum |
| `08-strings-assert` | repair | diagnostic repair workflow |
| `11-leaderboard` | modification | edit a working program |

Optional add-on (not in default intent): `15-band-counts` via `--tasks`.
Do **not** mix `intent_v1` rows into `core_v2` tables.

### Second model before claiming a win

One Gemini flash-lite table is not enough. Re-run the same suite/protocol on ≥1
other model before treating VOL vs Python workflow deltas as stable.

### Decision rule

> Ship language or card changes only if **first-try % stays high** and
> **cold or warm tokens to success** drop vs the previous published JSONL —
> not because hand-written density looked better.

---


## Strategy: radical dynamics (not radical syntax)

Same look; change **what you can omit** and **how values flow**.

### D1 — Coercion in string context (highest language ROI)

Today: `"a" + 1` → `R013`; pay `string(...)`.

Desired dynamic:

```vol
print "Sum: " + total    // string + displayable → concat
```

Keep `1 + "a"` / ambiguous numeric+string rejected if needed.
Still need `string(v)` for explicit convert in non-concat contexts.

Simulated “strip `string(...)`” on current suite only moved median
~0.868 → ~0.83 — necessary but not sufficient alone.

### D2 — Multi-arg `print` (biggest label killer)

Today: `print expression` (one value, one newline).

Desired dynamic (familiar, not new glyphs):

```vol
print "Class average:", avg
print "A grades:", scores.count(_ >= 90)
```

Space-join args, one trailing newline. Removes most `"…" + string(x)` glue
without inventing interpolation syntax.

### D3 — Collection intent ops (already shipping; lean harder)

Keep / teach as canonical:

- `.where(_)` / `.map(_)` / `.count(_)` / `.sum()`
- Prefer `.count(pred)` over `.where(pred).len`
- Chain filters; use `.each` only for effects

Future **fusion** (single-pass chains) is runtime — does not cut source tokens.
Source wins come from expressing more in one chain (already strong on 09).

### D4 — Suite dynamics (honest scoreboard)

If median vs Python is the goal, **what you measure** matters:

- More filter / aggregate / report tasks → VOL pulls ahead.
- More hello / while / fib → ties Python.

Do not fake 50% by padding Python or stripping VOL assert messages.

### D5 — Canon + tooling (SF-1 leftovers, density + workflow)

- `vol fmt`: prefer `_` over `fn` in collection ops; one bind/loop/value style.
- Examples emit formatter output only.
- Rejected aliases stay rejected (`.length` → `fix` for `.len`).

---

## Explicit non-goals for this preset

Do **not** pursue these *for density juice* until D1–D2 + card work land and
are measured:

- `|>` pipelines
- `?T` Option sugar
- `=>` arrows
- Dual-return sugar
- Parallel / ownership inference (not source-token wins today)
- Renaming keywords for 1-token shaves that grow the card

---

## Ranked next actions

### A. Language dynamics (SF-2 candidate — specify before ship)

- [ ] Spec draft: string `+` coercion rules + failure cases (what stays `R013`).
- [ ] Spec draft: multi-arg `print` (separator, newline, `nothing` handling).
- [ ] Implement + tests + examples; bump freeze only when SPEC/tests/docs agree.
- [ ] Rewrite density VOL tasks with new dynamics; `make check && make count`.
- [ ] Update README density bullets from new medians.

### B. Workflow tokens (can do under SF-1)

- [ ] Cut `vol_v1` card; re-run `core_v2` (`--langs vol` + baseline JSONL).
- [ ] Finish `vol fmt` rewriter with canonical table.
- [ ] Publish ≥1 model on **`intent_v1`** (prefer over `core_v2` for language-use claims).
- [ ] Publish ≥1 other model before treating workflow deltas as stable.

### C. Suite (honest measurement)

- [x] Add compression-heavy tasks (`14-pipeline-stats`, `15-band-counts`, `16-map-filter`).
- [x] Tier reporting in `count_tokens.py` (compression / labeled / parity).
- [x] Keep parity set (hello/loops/fib) so ties stay visible.
- [ ] After D1/D2 land, re-densify **labeled** VOL tasks and compare labeled median.

### D. Measurement hygiene

- [x] Report tokenizer name + median (all) + tier medians + per-task.
- [ ] Always split density vs workflow; cold vs warm for LLM runs.
- [ ] Never claim 50% vs Python until a measured **median (compression)** ≤ 0.50
      on a documented suite (not labeled-only).
- [x] Lang tests lock current `print` / `string` / `R013` / `R029` so D1–D2 must
      update SPEC + tests deliberately (`TestPrintDisplayForms`,
      `TestStringConcatTypeMismatches`, `TestPrintRejectsNothing`).

---

## Idioms already used in densified `bench/tasks/*/vol`

Use these as the hand-written / formatter target style:

```vol
// arithmetic — print expressions, no dead temps
print a + b

// value choice
print cond ? "yes" : "no"

// fixed repeats over while+counter
repeat 8 { … }

// expression-body functions
fn square(n) n * n

// count over where.len
scores.count(_ >= 90)

// chain
[4, 7, 2, 9, 12].where(_ > 5).sum()

// effects after filter
scores.where(_ >= 80).each score { print "High score: " + string(score) }
```

After D1/D2, prefer:

```vol
print "High score:", score
print "Sum:", total
```

---

## How to resume

```sh
# Source density
cd bench && make check && make count
# → bench/results/density.md

# LLM workflow (after card edits)
cd bench
uv run python llm/harness/run_generate_repair.py --provider gemini --suite core \
  --langs vol --baseline-jsonl llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-041440.jsonl
```

Decision rule for any change:

> Does it reduce **total tokens to success** (or source tokens without growing
> the card / repair rate)? If it only shortens hand-written golf, skip it.

---

## One-line north star

**Omit ceremony for display and aggregation; keep one familiar spelling per intent;
measure density and workflow separately; beat Python on intent-heavy programs,
not on hello-world ties.**
