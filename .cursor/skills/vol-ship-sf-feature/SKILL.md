---
name: vol-ship-sf-feature
description: >-
  Ships a Planned VOL language feature as a surface-freeze bump: implements
  lexer/parser/AST/resolver/interpreter, tests, examples, SPEC/IDEAS/README/
  AGENTS sync, and a new LLM card. Use when the user asks to implement pipelines,
  enums, dual-return, or other IDEAS.md Planned syntax; to bump SF-1→SF-2 (or
  later); or to “ship the next SF feature.”
disable-model-invocation: true
---

# VOL Ship SF Feature

Turn one **decided** Planned direction from [`IDEAS.md`](IDEAS.md) into Supported
syntax under a new surface freeze. Complements `vol-close-spec-corner` (decide /
docs-only). Follow [`AGENTS.md`](AGENTS.md) sync rules.

## When to use

- User names a Planned feature to implement (imports, structs, Result, `|>`, …)
- User says ship / implement / SF-2 bump for language surface
- Not for pure design Q&A — use `vol-close-spec-corner` first if undecided

## Before anything

Read current state (do not trust memory):

1. `SPEC.md` §0 (active freeze + card) and §11
2. Matching `IDEAS.md` section (direction must already be decided)
3. `bench/llm/cards/` — latest `vol_vN.md` and harness `LANG_META` in
   `bench/llm/harness/run_generate_repair.py`
4. Spot-check `internal/lang/{token,lexer,parser,ast,resolver,interpreter}.go`

If the feature is still an open question in IDEAS, stop and run Phase 2 of
`vol-close-spec-corner` (or ask the user to pick) before coding.

## Scope rules

- **One coherent SF slice** per bump (active freeze is SF-2 / card `vol_v2`; next
  is SF-3 / `vol_v3`). Do not boil the ocean (ownership + parallel + stdlib in one pass).
- Prefer the IDEAS **recommended / direction decided** spelling. Do not invent a
  competing form without asking.
- Docs-only directions that cannot run yet (allocation “unspecified”, parallel
  “no guarantees”) stay docs — do not fake Supported syntax.
- Never label Supported without tests.
- Do not keep intermediate draft cards — only freeze-bound cards (`vol_v0`, `vol_v1`, …).
  Note: `vol_v1` is the `core_v2` task card (SF-1 subset), not a full SF-1 product tour.

## Workflow

### 1. Lock the slice

State briefly:

- Feature(s) shipping
- From freeze → to freeze (e.g. SF-1 → SF-2)
- New card path (`bench/llm/cards/vol_v2.md`, …)
- Out of scope for this bump

If the user said “the rest,” propose a **single** next slice (default ROI under
SF-1 leftovers: **formatter / std behind imports**, or **`|>` / enums /
dual-return** when ready) and proceed only on that slice unless they insist on more.

### 2. Implement in `internal/lang`

Typical order:

1. `token.go` + `lexer.go` (keywords / operators)
2. `ast.go` nodes + `ast_test.go`
3. `parser.go` + parse diagnostics (`E…`, next free codes)
4. `resolver.go` (scopes / arity)
5. `interpreter.go` (runtime values, `display` / `typeName`, `R…` codes)
6. Tests: lexer vocabulary, parser diagnostics, `lang_test.go` / interpreter
   tables, `examples_test.go` row if adding examples

Reuse patterns from the current SF-1 surface (Option / Result / structs / modules) when similar.

### 3. Examples

Add a small executable under `examples/features/` (or `examples/projects/` when
multi-file) when the feature is teachable. Wire it in
`internal/lang/examples_test.go` and [`examples/README.md`](../../../examples/README.md).
Run it with `go run ./cmd/vol run`.

### 4. Freeze + docs sync (mandatory)

| File | Update |
| --- | --- |
| `SPEC.md` | §0 freeze/card; vocabulary; grammar; semantics; diagnostic tables; §10/§11; conformance snippet |
| `IDEAS.md` | Mark shipped; keep remaining Planned honest; strike open questions |
| `README.md` | What works / open experiment / teaching examples |
| `AGENTS.md` | Current reality + settled bullets if behavior changed |
| `bench/llm/cards/vol_vN.md` | New card; leave old cards historical |
| Harness `LANG_META` | Point VOL at new card + freeze id |
| `LLM_BENCHMARK.md` | Card paths / checklist note if relevant |

Historical LLM tables must keep their old freeze/card ids — do not rewrite past results.

### 5. Verify

```text
gofmt -w <changed Go files>
go test ./internal/lang/ ./cmd/vol/
go run ./cmd/vol run ./examples/basics/first.vol
git diff --check
```

Prefer `go test ./...` when time allows; report if skipped.

### 6. Finish

- Summarize: freeze bump, syntax shipped, files touched
- Do **not** commit unless the user asks
- Do **not** claim LLM superiority without a new published harness run

## Default next-slice order (if user says “the rest”)

Implement one at a time in this order unless the user overrides (SF-1 already
has modules, structs, Result/Option unwrap, density sugar):

1. Formatter / richer std libraries behind imports (foundations)
2. `|>` pipelines (if measured)
3. Enums / tagged unions
4. Dual-return sugar over Result
5. Leave ownership / parallel / alloc / build modes as design until a backend or
   explicit user ask

## Anti-patterns

- Shipping Planned syntax without an SF bump + new card
- Mixing SF-N and SF-(N+1) claims in one LLM result table
- Implementing Result as exceptions/`try`/`catch` (rejected model)
- Ambient stdlib growth (ambient = tiny core only)
- Claiming Supported for incomplete parse-only stubs
- Editing `vol_v0` / `vol_v1` used in published tables instead of adding `vol_vN`
- Keeping draft cards for unfinished intermediate surfaces

## Relation to other skills

| Skill | Role |
| --- | --- |
| `vol-close-spec-corner` | Decide / docs-only Phase-2 corners |
| `vol-ship-sf-feature` | Implement + freeze bump after direction is decided |
| `vol-roast-review` | Harsh audit of honesty / vision vs reality |
