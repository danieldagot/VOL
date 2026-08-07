# Example roast (tone and structure reference)

The following is a gold-standard example of voice, structure, and severity.
When reviewing, match this energy against **current** docs — do not copy outdated
findings if the repo has fixed them.

**Current doc layout (use this, not historical filenames in older roasts):**
`README.md`, `SPEC.md`, `IDEAS.md`, `AGENTS.md`, plus spot-checks of
`internal/lang/`, `examples/`, and `bench/`. `SYNTAX.md` and `VOCABULARY.md` were
removed; implemented grammar and vocabulary live in `SPEC.md` only.

The sample below reflects the project **after** `SPEC.md` became the single
source of truth, `bench/` added source token density, and §11 settled core
prototype rules. Re-verify everything before citing it in a live review.

---

I read `README.md`, `SPEC.md`, `IDEAS.md`, `AGENTS.md`, spot-checked
`internal/lang/`, and skimmed `bench/README.md`.

**Bottom line:** VOL is no longer “manifesto plus vibes.” You actually shipped a
living prototype spec, structured diagnostics, tests, CI, and a honest split
between what works and what is vision. That is real progress. You are still not
a systems language—you are a well-documented interpreter prototype with
research ambitions and one partial measurement of source token density. The gap
between AGENTS.md compiler fantasies and §10 “explicitly out of scope” is still
the story.

**My score today**

| Area | Score |
| --- | --- |
| Core idea | 8/10 |
| Surface syntax | 7/10 |
| Simplicity | 8/10 |
| Novelty | 5/10 |
| Specification quality | 6/10 |
| Systems-language semantics | 3/10 |
| Actual evidence of LLM optimization | 3/10 |
| Potential | 9/10 |

**The harshest version of the review is:**

VOL finally wrote the spec—now it needs to stop treating research slides as semantics.

And that's completely fixable at this stage.

## The roast

### You fixed the biggest doc sin

`SPEC.md` exists. It has lexical rules, expression precedence, evaluation order,
array identity, `.where` / `.sum` semantics, a failure model with stable codes,
conformance examples, and **§11 Decided** rules that actually answer questions
the old split grammar docs punted on: mutable-by-default with opt-in `const`,
overflow traps, `if` as statement vs `? :` for values, permanent `while`, `.len`
vs rejected `.length`, and eager `.where` returning a new array. An LLM can now
answer whether `not false == true` is legal—it is, and it prints `true` (§9.1).

That moves you from “language design cosplay” to “prototype with rules.” Credit
where due.

### But “systems language” is still mostly a wardrobe

`SPEC.md` §10 lists the elephant herd: static types, ownership, structs, modules,
generics, backends—all explicitly out of scope. README admits the same with
admirable honesty. Meanwhile AGENTS.md still talks about stack allocation when
values don't escape, automatic vectorization, and ownership analysis. The docs
*mostly* label inference as aspiration, but the **tone** still sells a compiler
that does not exist. For a systems claim, i64/f64/bool/string/`[]any` is a toy
value model, not a platform.

### `.where` is specified—and still a future landmine

This is much better than before:

```vol
total := numbers.where(_ > 5).sum
```

`SPEC.md` §4.4 now says: eager new array, ordered, `_` binding rules, integer
fold for `.sum`, overflow via normal `+` and `R028`. String `.len` is Unicode
scalars; bytes are `.byte_len`. Good.

Two problems remain:

1. **Purity is decided in prose but not enforced.** The spec says impure `.where`
   predicates are non-conforming and may break under future fusion—yet the
   interpreter still evaluates them eagerly left-to-right. You are asking users
   and LLMs to follow a rule the runtime does not police. That is spec/impl drift
   with a time bomb.

2. **Performance story is still hand-wavy.** “Not parallel and not lazy in this
   prototype” is honest—but AGENTS.md still dangles vectorization. Collection
   syntax is not a vector programming model until fusion/SIMD semantics exist.

### LLM optimization: half a datapoint is not a thesis

README and `bench/` report ~36% fewer source tokens than Go (median, 13 tasks,
two tokenizers). That is **real** and **labeled correctly** as source density
only—not generate/repair success. IDEAS.md still tracks `LLM_BENCHMARK.md` as
todo. AGENTS.md's Meaning per Token remains undefined.

So the fair roast is not “you have zero evidence.” It is: **you have one metric
of a multi-metric claim.** Shortest source still loses if models produce 40%
more `R028`/`S003`/`R029` and burn tokens fixing them. The next document that
matters is not another syntax idea—it is falsifiable generate/repair numbers.

### Documentation hygiene is mostly healed

Naming is stable: **Vocabulary Optimized Language** in README and AGENTS.md.
The “one canonical way” principle was updated to **one canonical representation
per distinct intent**—which legitimately allows `.each` vs `.where(...).sum`.
Planned `parallel { }` lives in IDEAS.md, not executable grammar in SPEC.

Roast what's left: vision sections in AGENTS.md are long relative to the tiny
implemented core; new contributors could still skim the top and overestimate
readiness. Keep the “current reality” banner loud.

### Engineering substance improved

Resolver, JSON diagnostics (`vol --json`), stable codes with `fix` on some
errors, table-driven tests, example conformance, fuzz non-panic, CI with race
detector—this is how a language prototype earns trust. The interpreter is still
tree-walking Go, not a backend, but the **process** looks serious now.

## Biggest design issue

Do not add features. **Close the loop on the core you already claim.**

Pick the frozen surface from README (“What Actually Works”) and make it boringly
complete:

- enforce or statically warn on `.where` purity—or demote purity to IDEAS until
  you can enforce it
- finish diagnostic `fix` coverage for the codes LLMs hit most
- ship `LLM_BENCHMARK.md` with a tiny task suite (generate → compile/run → repair
  token totals)
- keep §11 synchronized with tests on every change

## What to do next

1. **`LLM_BENCHMARK.md`** — even 20 tasks with public prompts beats another keyword.
2. **Purity** — reject impure `.where` in the interpreter *or* mark purity as
   Planned-only in SPEC until then.
3. **Demote vision tone** — shorten AGENTS compiler sections or label every bullet
   “not specified.”
4. **Types/memory** — when ready, one numbered SPEC section at a time with tests;
   no stealth semantics.

## Keep / throw out

**Keep** ~80% of the surface: `:=`, braces, `fn`, `if`, arrays, structured
errors, semantic compression (`.where(...).sum`).

**Throw out** (for now) treating compiler inference, vectorization, and ownership
elimination as implied behavior. Write the semantic rules first; let the optimizer
be a later chapter.

That's where VOL becomes interesting: **the first small systems language whose
docs and benchmarks let you falsify the LLM workflow claim—not just count tokens
in hand-written examples.**
