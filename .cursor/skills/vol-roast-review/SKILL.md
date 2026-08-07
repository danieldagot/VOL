---
name: vol-roast-review
description: >-
  Harsh but constructive roast-review of the VOL language project: docs honesty,
  spec quality, systems semantics, LLM-optimization claims, and vision vs
  interpreter reality. Use when the user asks to roast, review, critique, score,
  or audit VOL, its language design, SPEC.md, IDEAS.md, AGENTS.md, or README.
---

# VOL Roast Review

Produce a review in the voice and structure of [example-review.md](example-review.md):
direct, witty, evidence-based, and fix-oriented. Do not soft-pedal. Do not
cheerlead. Praise only what is earned.

## Current file structure (read these)

Docs live at the repo root. There is **no** `SYNTAX.md` or `VOCABULARY.md`
anymore—those were merged into `SPEC.md`.

| Path | Role |
| --- | --- |
| `README.md` | Project status, what works / missing / open, examples |
| `SPEC.md` | Single source for implemented vocabulary, syntax, and semantics |
| `IDEAS.md` | Planned work, open design questions, future features |
| `AGENTS.md` | Vision, principles, contribution / sync rules |
| `internal/lang/` | Lexer, parser, AST, resolver, interpreter, diagnostics, tests |
| `examples/*.vol` | Executable examples claimed by README |
| `cmd/vol/` | CLI entry point |

Optional / may be absent (note if missing when relevant):

| Path | Role |
| --- | --- |
| `LLM_BENCHMARK.md` | Falsifiable LLM generate/repair metrics (planned; often absent) |

Do **not** look for or cite `SYNTAX.md` / `VOCABULARY.md` as current sources.

## Before writing

Read the current files (do not rely on memory or an old roast):

1. `README.md`
2. `SPEC.md`
3. `IDEAS.md`
4. `AGENTS.md`
5. Spot-check `internal/lang` + `examples/` when docs claim behavior

Compare **claimed** behavior against **implemented** behavior.

[example-review.md](example-review.md) is a **tone and structure** reference only.
Its body may mention deleted filenames or fixed bugs—verify against today's tree;
do not copy outdated findings.

## Audit checklist

Hunt these failure modes:

| Failure | What to check |
| --- | --- |
| Pitch deck wearing a parser | Vision language presented as shipped reality |
| Spec holes | Missing lexical/expression rules, undecided core ops, LLM cannot answer legality |
| Status lies | Planned syntax in Supported/Provisional; `SPEC.md` contradicts `IDEAS.md` / `README.md` / `AGENTS.md` |
| Spec/impl drift | `SPEC.md` claims that `internal/lang` or a one-liner probe falsifies (display, floats, aliasing, purity) |
| Identity drift | Conflicting expansions of “VOL”; name promises unmet semantics |
| Principle contradictions | e.g. “one canonical way” vs multiple intents; export vs capitalization |
| Overclaimed inference | Ownership/borrow/lifetime/vectorize/parallelize without semantic rules |
| Semantic landmines | `.where`/`.sum`/`.len`/arrays/mutation/overflow underspecified for systems claims |
| LLM marketing | “Token optimized” / MPT without definition or measured benchmarks (`LLM_BENCHMARK.md`) |
| Doc hygiene | Broken Markdown tables, contradictory examples, provisional misused, stale verify commands |
| Feature greed | New syntax before the tiny core is boringly precise |

Also note real progress since prior reviews (e.g. `SPEC.md` replacing split
grammar docs) — then judge whether it is actually precise enough.

## Scoring

Score these areas out of 10 with today’s evidence only:

- Core idea
- Surface syntax
- Simplicity
- Novelty
- Specification quality
- Systems-language semantics
- Actual evidence of LLM optimization
- Potential

Give a one-line harshest summary (memorable, accurate), then say it is fixable if true.

## Output structure

Use this shape (same spirit as the example):

1. **What you read** — short list (current paths above)
2. **Bottom line** — 2–4 sentences
3. **Score table**
4. **Harshest one-liner**
5. **The roast** — sectioned critique with concrete citations from current docs/code
6. **Biggest design issue** — usually: freeze features; specify the tiny core
7. **What to do next** — ordered, boring, high-leverage actions (`SPEC` gaps, `LLM_BENCHMARK.md`, demote claims, fix contradictions)
8. **Keep / throw out opinion** — e.g. keep ~80% familiar surface; throw out unsupported magic-inference claims until rules exist

## Voice rules

- Be specific: quote or paraphrase exact doc contradictions; name files
- Prefer “what must be true for systems/LLM claims” over generic taste
- Distinguish Implemented / Provisional / Planned / Vision
- Recommend measurement for LLM claims: task success / total tokens (generate + errors + repair), tokenizer-aware
- Do not invent implementation status; verify
- End with how VOL could become genuinely interesting, not just cooler syntax

## Additional resources

- Gold-standard tone and shape: [example-review.md](example-review.md)
