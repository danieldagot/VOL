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

Suite: 13 tasks × VOL / Python / Go / Rust / Zig. Regenerated after densifying
VOL idioms (`.count`, `repeat`, expression-body `fn`, ternary, fewer temps) and
adding Python baselines.

Median ratios (`cl100k_base` ≈ `o200k_base`):

| Baseline | VOL / baseline | ≈ fewer tokens |
| --- | ---: | ---: |
| Python | **0.868** | **~13%** |
| Go | **0.566** | **~43%** |
| Rust | **0.610** | **~39%** |
| Zig | **0.384** | **~62%** |

Full tables: [`bench/results/density.md`](bench/results/density.md).

### Per-task VOL/Python (why 50% is hard)

| Task | VOL/Python | Notes |
| --- | ---: | --- |
| 01-hello | 1.000 | Tied — label + concat tax |
| 04-loops | 1.000 | Tied |
| 05-arrays-each | 0.982 | Near tie |
| 10-fibonacci | 0.938 | Already `repeat` + multi-assign |
| 13-temperatures | 0.876 | Labels dominate |
| 03-conditions | 0.868 | |
| 07-functions | 0.878 | |
| 11-leaderboard | 0.778 | |
| 06-where-sum | 0.754 | Collection win |
| 08-strings-assert | 0.711 | |
| 12-revenue | 0.711 | |
| 02-arithmetic | 0.676 | |
| 09-grade-report | **0.619** | Best — `.count` pipelines |

**Target “50% fewer than Python” (ratio 0.5) on this suite is not realistic**
without either (a) strong display/print dynamics + suite rebalance, or
(b) unfair golf. Zero tasks are at ≤0.5 today; best is 0.619.

Realistic near-term vs Python after dynamics: **~20–35% fewer** on a
compression-heavy mix — not 50% on the current 13.

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
- [ ] Publish ≥1 other model on `core_v2` before further sugar.

### C. Suite (optional, honest)

- [ ] Add compression-heavy tasks (multi-filter reports) if claiming vs-Python density.
- [ ] Keep a small “parity” set (hello/loops) so ties stay visible.

### D. Measurement hygiene

- [ ] Always report tokenizer name + median + per-task.
- [ ] Always split density vs workflow; cold vs warm for LLM runs.
- [ ] Never claim 50% vs Python until a measured median ≤ 0.50 on a documented suite.

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
