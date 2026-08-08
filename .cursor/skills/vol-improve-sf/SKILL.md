---
name: vol-improve-sf
description: >-
  Iterates the active VOL Surface Freeze for token efficiency and honesty:
  density benches, live LLM intent_v1, diagnostic Fix text, card cuts, denser
  examples/docs — without new unscheduled Planned glyphs. Use when the user asks
  to improve SF-3.1/any SF, maximize token efficiency, tighten vol_vN cards,
  re-run Gemini intent, densify benches/examples, or finish an SF juice loop.
  Active freeze is SF-3.1 (`vol_v3_1`); SF-3 juice is exhausted in
  TOKEN_EFFICIENCY.md — next is SF-3.1 foundation, then tooling / second model;
  do not default-ship Planned sugar.
disable-model-invocation: true
---

# VOL Improve Any SF

Improve the **active** surface freeze in place (SF-3.1 / `vol_v3_1`). Complements
`vol-ship-sf-feature` (new freeze / new syntax). Do **not** invent `|>`, `=>`,
enums, dual-return, or other Planned sugar for juice — ship those only via the
ship skill + freeze bump. SF-3 juice is exhausted; prefer SF-3.1 foundation
precision over sugar.

Follow [`AGENTS.md`](AGENTS.md) sync rules. Keep density and LLM workflow tracks
**distinct**.

## When to use

- “Improve SF-N”, “SF juice”, “token efficiency”, tighten card / Fix / examples
- Re-measure Gemini `intent_v1` after card or diagnostic changes
- Densify hand-written benches under already-Supported forms
- Exhaust / document why further in-freeze gains need tooling or SF+1

## Before anything

Read current state (do not trust memory):

1. `SPEC.md` §0 — active freeze id + default card (`vol_vN`)
2. `TOKEN_EFFICIENCY.md` — scoreboard, decision rule, exhausted notes
3. `bench/llm/cards/vol_vN.md` + harness `LANG_META` in
   `bench/llm/harness/run_generate_repair.py`
4. Latest `bench/llm/results/intent_v1_live_gemini_*.md` for the same model
5. Spot-check `examples/` and `examples/README.md` if teaching surface changed

Confirm `GEMINI_API_KEY` in env / `.env` before live runs.

## Two tracks (never mix claims)

| Track | Measure | Command / artifact |
| --- | --- | --- |
| **Source density** | median(all), median(compression), median(labeled) | `cd bench && make check && make count` → `bench/results/density.md` |
| **LLM workflow** | first-try %, success @ K, cold/warm means vs Python | `uv run python llm/harness/run_generate_repair.py --provider gemini --suite intent --langs vol,python --replicates 3` |

Do **not** cite density ratios as workflow proof. Historical result tables keep
their freeze/card ids — never rewrite past freezes in published tables.

## Decision rule

> Ship a change only if **first-try % stays high** and **cold or warm tokens to
> success** drop vs the previous same-suite JSONL — not because hand-written
> density alone looked better.

Keep Fix-text-only changes that clearly cut a known repair round even when
generation is unchanged, then re-measure when possible.

## Constraints

- Active freeze + card from `SPEC.md` §0 — no freeze bump unless evidence demands
  **new** Supported surface and the user agrees (then use `vol-ship-sf-feature`)
- Prefer dynamics, card cuts, diagnostic `Fix`, denser idioms under frozen surface
- Card target: ≤ ~Python card (`python_v0.md`) **without** losing first-try
- Ambient builtins stay tiny; `@std` growth only if already in the freeze pin
- Needs live API for workflow claims; report if skipped

## Loop (repeat until stop)

### 1) Baseline / measure

```sh
cd bench
make check && make count
uv run python llm/harness/run_generate_repair.py \
  --provider gemini --suite intent --langs vol,python --replicates 3
```

Read the new `bench/llm/results/intent_v1_live_gemini_*.md` + JSONL. Note density
medians. Prefer `intent_v1` over `core_v2` for language-use claims.

### 2) Diagnose failures (from JSONL)

Classify every non-first-try / non-success:

`parse` / `S0xx` / `R0xx` / wrong stdout / `source_check_failed`

Rank footguns by (frequency × repair token cost). Note idioms emitted
(`.count()` vs `.count(pred)` vs `.where.len`, multi-arg `print` vs `string()` glue,
`@std` names if SF-3+).

### 3) Improve — **one primary lever** per iteration

Preference order:

1. **Diagnostic `Fix` text** that prevents a known repair round (no new syntax)
2. **Shrink/clarify** `bench/llm/cards/vol_vN.md` (cut tokens; teach failing idiom)
3. **Densify** hand-written VOL benches/examples using already-Supported forms
4. **Docs sync** (`README`, `TOKEN_EFFICIENCY`, `examples/README`, SPEC pointers)
5. Only if evidence demands: tiny semantics still inside freeze spirit — SPEC +
   tests + docs per AGENTS.md (else stop and propose ship skill + bump)

### 4) Verify

```text
gofmt -w <changed Go files>
go test ./...
go run ./cmd/vol run ./examples/basics/first.vol
# plus any new/changed feature example
cd bench && make check && make count
```

Re-run Gemini (VOL-only OK if Python unchanged):

```sh
uv run python llm/harness/run_generate_repair.py \
  --provider gemini --suite intent --langs vol --replicates 3 \
  --baseline-jsonl <prior_same_model_intent_jsonl_with_python>
```

Update `README.md` / `TOKEN_EFFICIENCY.md` honestly. Do not rewrite historical
SF tables’ freeze ids.

### 5) Keep / revert

Keep only if metrics improve under the decision rule. Otherwise revert and try
the next lever. Log: what failed → what changed → density delta → workflow delta.

## Stop criteria

Stop the loop when **all** are true:

1. Live `intent_v1` VOL first-try = 100% and success @ K = 100% (same model/protocol)
2. Cold and warm VOL means ≤ previous best JSONL for this suite+model
   **or** warm ≤ Python and cold gap is only irreducible card tax
3. No remaining freeze-legal change (card, Fix, densify, in-freeze dynamics)
   predicted to improve metrics without growing the card or lowering first-try
4. Document why further juice is exhausted in `TOKEN_EFFICIENCY.md`

## Done output

When stopping, produce:

- Final density medians (all / compression / labeled)
- Final `intent_v1` Gemini summary path
- Changes kept vs rejected
- Explicit “SF-N exhausted” rationale (what was tried; what needs `vol fmt`,
  second model, prompt caching, or SF+1)

## Anti-patterns

- Claiming workflow win from density alone
- Growing the card to teach one rare idiom without measuring cold impact
- Editing historical cards used in published tables without republishing
- Shipping SF+1 sugar “for density” without the ship skill + freeze bump
- Mixing freeze IDs in one result table
- Skipping live re-measure after card/Fix changes when claiming workflow gains

## Relation to other skills

| Skill | Role |
| --- | --- |
| `vol-improve-sf` | Iterate active freeze (this skill) |
| `vol-ship-sf-feature` | New syntax / freeze bump |
| `vol-close-spec-corner` | Decide undecided design corners |
| `vol-roast-review` | Harsh honesty audit |
