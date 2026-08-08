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
| `README.md` | Project status, what works / missing / open, examples, bench summary |
| `SPEC.md` | Single source for implemented vocabulary, syntax, semantics, diagnostics |
| `IDEAS.md` | Planned work, open design questions, future features |
| `AGENTS.md` | Vision, principles, contribution / sync rules, decided prototype bullets |
| `LLM_BENCHMARK.md` | Falsifiable LLM generate/repair protocol + published-result pointers |
| `TOKEN_EFFICIENCY.md` | Working preset for density vs workflow iteration (SF-2 exhausted note) |
| `internal/lang/` | Lexer, parser, AST, resolver, interpreter, diagnostics, tests |
| `examples/` | Executable examples (`basics/`, `features/`, `projects/`) claimed by README |
| `cmd/vol/` | CLI (`vol run`, `--json` diagnostics) |
| `bench/` | Source token density + LLM harness (`bench/llm/`) |

Do **not** look for or cite `SYNTAX.md` / `VOCABULARY.md` as current sources.

## What the project actually is today (baseline for scoring)

Use this snapshot before roasting—do not pretend the repo is still pre-`SPEC.md`:

- **Implemented (SF-2):** tree-walking interpreter; resolver; structured diagnostics
  (human + JSON via `vol --json`); stable error codes with `Fix` on some codes;
  `i64`/`f64`/bool/string/arrays; `const` opt-in immutability; overflow traps
  (`R028`); `.where` / `.map` / `.count` / `.count()` / `.sum`; `.len` /
  `.byte_len`; `.copy` / `.deep_copy`; Option / Result; product structs;
  anonymous `fn`; multi-assign; multi-arg `print` + string `+` coercion;
  `import`/`export`; functions with `nothing` on missing return; CI + fuzz.
- **Specified in `SPEC.md`:** lexical tokens, line joining, expression precedence,
  evaluation order, array identity, collection / Option / Result semantics,
  failure model, conformance examples, **§11 Decided** core rules (mutability
  default, overflow, purity intent, `if` vs `? :`, permanent `while`, etc.).
- **Partially evidenced LLM claim:** `bench/` source density (~20% fewer than
  Python all-suite / ~24% compression) plus early Gemini `intent_v1`/`vol_v2`
  workflow runs — not multi-model proof or SF-2 “LLM optimized” evidence.
- **Still missing / vision-only:** static types, ownership/borrow, generics, native
  backend, formatter rewriter, enforced `.where` purity at runtime,
  parallel/lazy collection fusion, enums / `|>` / dual-return (SF-3 candidates).

Identity is **Vocabulary Optimized Language** in README and AGENTS.md. Do not
roast stale “Vector-Oriented” / “Vibe-Oriented” naming unless it reappears in
current docs.

## Before writing

Read the current files (do not rely on memory or an old roast):

1. `README.md`
2. `SPEC.md` (especially §4 precedence, §3.3 arrays, §4.4 collections, §8 failure, §11 decided)
3. `IDEAS.md`
4. `AGENTS.md`
5. `LLM_BENCHMARK.md` / `TOKEN_EFFICIENCY.md` when scoring LLM or density claims
6. Spot-check `internal/lang` + `examples/` + `bench/README.md` when docs claim behavior

Compare **claimed** behavior against **implemented** behavior.

[example-review.md](example-review.md) is a **tone and structure** reference only.
Its body reflects a point-in-time roast—verify every finding against today's tree;
do not copy outdated critiques (deleted grammar files, unfixed naming drift, “create
SPEC.md”, missing precedence tables, etc.).

## Audit checklist

Hunt these failure modes:

| Failure | What to check |
| --- | --- |
| Pitch deck wearing a parser | Vision language presented as shipped reality |
| Spec holes | Remaining gaps in the tiny core (types, purity enforcement, build modes) |
| Status lies | Planned syntax in Supported/Provisional; `SPEC.md` contradicts `IDEAS.md` / `README.md` / `AGENTS.md` |
| Spec/impl drift | `SPEC.md` claims that `internal/lang` or a one-liner probe falsifies |
| Aspirational purity | `.where` purity documented but not enforced; side effects “work” today |
| Identity drift | Conflicting expansions of “VOL”; name promises unmet semantics |
| Principle contradictions | “One canonical way” vs multiple intents (AGENTS now says one per **intent**—check if docs still contradict) |
| Overclaimed inference | Ownership/borrow/lifetime/vectorize/parallelize without semantic rules |
| LLM marketing overreach | Source density bench cited as “LLM optimized”; MPT undefined; no `LLM_BENCHMARK.md` |
| Doc hygiene | Broken tables, contradictory examples, provisional misused, stale verify commands |
| Feature greed | New syntax before the frozen core is boringly precise and tested |

Also note real progress since prior reviews (`SPEC.md`, §11 decided rules,
diagnostics JSON, `bench/`, resolver, honest README split)—then judge whether
precision and evidence are **enough** for the claims still being made.

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

Typical range guidance (adjust with evidence—do not copy blindly):

- **Specification quality** was ~3/10 pre-`SPEC.md`; a living draft with grammar,
  precedence, and §11 decided rules is often **5–7/10** unless you find major drift.
- **LLM optimization evidence** may be **2–4/10** if only `bench/` source density
  exists; still not **6+** without generate/repair benchmarks.
- **Systems semantics** stays low (**2–4/10**) until types, memory, and ABI rules
  exist—even if the prototype core is well specified.

Give a one-line harshest summary (memorable, accurate), then say it is fixable if true.

## Output structure

Use this shape (same spirit as the example):

1. **What you read** — short list (current paths above)
2. **Bottom line** — 2–4 sentences; acknowledge progress and remaining gap
3. **Score table**
4. **Harshest one-liner**
5. **The roast** — sectioned critique with concrete citations from current docs/code
6. **Biggest design issue** — usually: finish specifying/enforcing the tiny core before expanding
7. **What to do next** — ordered, boring, high-leverage actions (`LLM_BENCHMARK.md`, purity enforcement or honest “documented-only”, demote vision claims, close spec holes)
8. **Keep / throw out opinion** — e.g. keep familiar surface + semantic compression; throw out unsupported magic-inference claims until rules exist

## Voice rules

- Be specific: quote or paraphrase exact doc contradictions; name files and sections
- Prefer “what must be true for systems/LLM claims” over generic taste
- Distinguish Implemented / Provisional / Planned / Vision / Documented-but-not-enforced
- Recommend measurement for LLM claims: task success / total tokens (generate + errors + repair), tokenizer-aware; distinguish from `bench/` source density
- Do not invent implementation status; verify with tests or a quick probe
- Do not roast fixed problems unless they regressed
- End with how VOL could become genuinely interesting, not just cooler syntax

## Additional resources

- Gold-standard tone and shape: [example-review.md](example-review.md)
