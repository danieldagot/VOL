---
name: vol-close-spec-corner
description: >-
  Facilitates VOL Phase A open language decisions from SPEC.md §11: builds an
  Already done / Still open status table with options and suggested picks, helps
  the user decide, then updates SPEC.md, IDEAS.md, README, and implementation
  or tests when behavior changes. Use when the user asks about Phase A, SPEC
  open decisions, closing language corners, mutability/overflow/while/if/.where
  decisions, or wants a weekly open-questions list.
disable-model-invocation: true
---

# VOL Close Spec Corner (Phase A)

Help the user close one boring language corner at a time. Prefer precision over
new features. Follow `AGENTS.md` documentation and implementation sync rules.

## Before anything

Read current files (do not rely on memory or an old table):

1. `SPEC.md` — especially §11 Open decisions, Decided, and the sections each issue touches
2. `IDEAS.md` — Planned syntax and Open Design Questions
3. `README.md` — status lists users see
4. Spot-check `internal/lang` when “Today” behavior matters

Regenerate the status list from **today’s** SPEC. Do not paste a stale table.

## Phase 1 — Status list

Unless the user already picked an issue, produce this structure:

### Already done

One bullet per decided item from SPEC §11 Decided (and related SPEC sections), e.g.:

```text
Mutability — mutable by default; const planned (shallow).
```

### Still open (decide one per week)

Markdown table:

| # | Issue | Today | Options | Suggested pick for VOL |
| --- | --- | --- | --- | --- |
| … | … | current interpreter/SPEC behavior | real alternatives | one opinionated recommendation aligned with VOL vision |

Cover every unresolved item in SPEC §11. Number stably for the session.

### Recommended order (ROI)

Ordered list of the open issues by leverage for a systems + LLM-friendly core.
Default heuristic (reorder if SPEC evidence says otherwise):

1. Array assignment — aliases vs systems story
2. Overflow — correctness
3. `.where` purity — intent split vs `.each`
4. Missing return — quiet bugs
5. String `.len` — easy to specify
6. `while` — quick keep/rename/drop
7. `if` expression — nicety, not blocking

Then ask which issue to decide (or proceed if they already named one).

## Phase 2 — Help decide

For the chosen issue only:

1. **What it means** — today’s rule in plain language + a tiny VOL example
2. **Options** — pros/cons table per option
3. **Vision fit** — how each option matches `AGENTS.md` (intent over boilerplate, LLM-familiar surface, systems honesty, no fake inference)
4. **Recommendation** — one clear pick, with what can stay Planned in `IDEAS.md` vs what must be true in the interpreter today
5. **Edge cases** — list what tests or SPEC paragraphs must cover if they accept

Do not implement or edit docs until the user confirms the decision (or explicitly says to apply a stated pick).

## Phase 3 — Close the corner

After the user confirms, do this loop:

1. Write the rule in `SPEC.md` (real section + §11: move to Decided, remove from open list)
2. Update `IDEAS.md` (checklists, Planned syntax, strike/answer open questions)
3. Change the interpreter / lexer / parser / resolver only if behavior or accepted syntax changes
4. Add tests for the edge cases from Phase 2
5. Update `README.md` only if users would write different code or status lists change
6. Update `examples/` (basics / features / projects) when the feature is teachable;
   sync `examples/README.md`

### Spec hygiene (mandatory)

- Never label Supported unless implemented and tested
- Planned syntax belongs in `IDEAS.md`, not as executable grammar in `SPEC.md`
- Keep vision vs prototype distinct
- Preserve unrelated user edits
- Prefer one canonical form per distinct intent

### Verification when Go or examples change

```text
gofmt -w <changed Go files>
go test ./...
go run ./cmd/vol run ./examples/basics/first.vol
git diff --check
```

If a command cannot run, say so. Do not claim the corner closed while relevant tests fail.

Docs-only decisions (behavior already matches): still update SPEC/IDEAS/README as needed; skip interpreter changes.

## Output style

- Concise; lead with the list or the decision summary
- One issue at a time after the status list
- When applying changes, summarize what was decided and which files changed
- Do not commit unless the user asks

## Anti-patterns

- Inventing new syntax while open corners remain vague
- Putting Planned forms in SPEC vocabulary as Supported
- Claiming ownership/borrow/parallel/fusion without rules
- Closing multiple unrelated §11 items in one pass unless the user asks
