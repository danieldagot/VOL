# VOL Ideas and Future Work

This document collects ideas that may shape VOL in the future.

Items here are exploratory, not promises. Implemented language behavior belongs in
`SPEC.md`. Long-term vision without executable semantics stays here until it is
specified and tested.

## Near-Term Priority: Precise Core Before New Features

**Surface Freeze SF-0 is active** (see [`SPEC.md`](SPEC.md) §0). The Supported
Prototype v0 surface and LLM card `bench/llm/cards/vol_v0.md` are frozen.
Do not expand the language until foundations below justify an SF-1 bump.

Frozen core (do not grow under SF-0):

- values, variables, functions
- `if` / `elif` / `else`, `repeat`, `while`, `? :`
- arrays, `.each`, `.where`, `.sum()`
- built-ins already accepted by the interpreter

Immediate documentation and design work:

- [x] Write `SPEC.md` covering lexical grammar, expression grammar, types/values,
      scopes, evaluation order, numeric behavior, arrays/strings, mutation,
      functions, control flow, and failure behavior for the implemented core.
- [x] Declare Surface Freeze SF-0 (SPEC + `vol_v0` language card).
- [ ] Keep `SPEC.md` synchronized whenever interpreter behavior changes.
- [x] Error/result model sketch (hybrid traps + Result; dual-return sugar later).
- [ ] Foundations before SF-1: smaller matched language card, formatter design,
      modules sketch (still Planned — not syntax).
- [x] Write `LLM_BENCHMARK.md` with a falsifiable generate/repair protocol;
      harness + Gemini `core_v2` published — see checklist in that file.
- [ ] Keep Planned syntax out of `SPEC.md`; only Supported and Provisional forms belong there.
- [x] Publish frozen `core_v2` with default Python baseline (`--langs vol,python`).
- [ ] Publish ≥1 other model on `core_v2` before treating LLM results as stable.

## Near-Term Foundation

- [x] Define tokens and lexical rules.
- [x] Implement the lexer.
- [x] Implement the parser.
- [x] Define the abstract syntax tree (AST).
- [x] Track precise source locations on every syntax node.
- [x] Create deterministic human-readable diagnostics.
- [x] Create machine-readable JSON diagnostics.
- [x] Implement lexical scopes.
- [x] Implement an interpreter.
- [x] Add table-driven parser, semantic, runtime, diagnostic, CLI, and example conformance tests.
- [x] Add parser fuzz seeds and CI checks for formatting, vet, races, and coverage.
- [ ] Add versioned diagnostic snapshots when the diagnostic presentation is stable.

## Initial Language Features

- [x] Integer, floating-point, Boolean, and string literals.
- [x] Variables with inferred types.
- [x] Mutability default: bindings are mutable (`:=` / `=`). See `SPEC.md` §5.2.
- [x] Opt-in immutable bindings with `const name := expression` (shallow). See `SPEC.md` §5.2.
- [x] Arithmetic, comparison, and Boolean operators.
- [x] Integer overflow traps by default (`R028`). See `SPEC.md` §4.3.
- [ ] Wrapping integer arithmetic and/or overflow build modes.
- [x] JSON diagnostic output from the `vol` CLI (`vol --json run <file.vol>` or `vol run --json <file.vol>`).
- [x] Blocks using braces.
- [x] `if` / `elif` / `else` statements.
- [x] Conditional operator `? :` (expression form; `if` stays a statement).
- [x] `repeat` loops.
- [x] Conditional loops with `while` (Supported; no alternate spellings).
- [x] Arrays, indexing, and bounds checking.
- [x] Array assignment shares references (see `SPEC.md` §3.3).
- [x] Explicit array clone: `a.copy()` (shallow) and `a.deep_copy()` (recursive). See `SPEC.md` §3.3.
- [x] Collection iteration with `.each`.
- [x] Filtering with `.where` and aggregation with `.sum()`.
- [x] `.where` predicates are pure by rule; side effects use `.each` (`SPEC.md` §4.4).
- [ ] Purity diagnostics / checks for `.where` predicates.
- [x] Functions, parameters, calls, and return values.
- [x] Missing return yields `nothing`; using it as a value is `R029` (`SPEC.md` §5.8).
- [ ] Typed void vs valued functions (require `return` on all paths when typed).
- [x] Comments.
- [x] Basic built-ins for printing, input, `.len`, assertions, conversion, and arguments.
- [x] `.len` is canonical (Unicode scalars for strings); `.length` rejected with fix hint.
- [x] String byte-count property `.byte_len` for systems I/O. See `SPEC.md` §3.2.

## Planned Syntax (Not Accepted Today)

These forms must not be documented as Supported or Provisional until implemented.

### String byte length ✓ implemented

Language rule (§3.2): `.len` counts Unicode scalar values; `.byte_len` counts
UTF-8 bytes. Grapheme-cluster counting is out of scope.

```vol
print "🙂".len       // 1  (Unicode scalar)
print "🙂".byte_len  // 4  (UTF-8 bytes)
```

### Typed void vs valued functions

Decided runtime rule (already in `SPEC.md` §5.8): falling off a function yields
`nothing`; call statements may discard it; `:=` / `=` / other value contexts
reject `nothing` with `R029`.

When static types exist, Planned refinements:

- valued functions must `return` on every path
- void / procedure functions may omit `return`
- better diagnostics naming the callee that produced `nothing`

### `.where` purity checks

Decided language rule (already in `SPEC.md` §4.4): `.where` predicates are pure
filters; use `.each` for side effects. The prototype still evaluates predicate
expressions eagerly and does not yet reject impurity.

Planned:

- diagnostics when a `.where` predicate clearly has side effects (or cannot be
  proven pure), with a `fix` suggesting `.each` when appropriate
- allow clearly pure helper calls such as Boolean predicates over `_`
- treat purity as a prerequisite for any future fusion or parallel `.where`

### Wrapping integer arithmetic / overflow modes

Decided language rule (already in `SPEC.md` §4.3): integer overflow **traps** by
default (`R028`) with a human message and machine-readable `fix` field so LLMs
can help repair programs.

Planned escapes (not accepted today):

- explicit wrapping operations, and/or
- build modes such as debug-trap (default today) vs release-wrap

JSON diagnostic output is now implemented: `vol --json run <file.vol>` emits a
single JSON object on stderr. The human form is still the default.

### Explicit array clone ✓ implemented

Language rule (§3.3): `:=` / `=` and argument passing share array identity.
Use `.copy()` for a shallow clone or `.deep_copy()` for a recursive clone:

```vol
a := [1, 2]
b := a.copy()
b[0] = 9
print a          // [1, 2]
print b          // [9, 2]

inner := [1, 2]
outer := [inner]
dc := outer.deep_copy()
dc[0][0] = 99
print inner      // [1, 2] — not affected
```

### Opt-in `const` bindings ✓ implemented

Language rule (§5.2): bindings are mutable by default; `const` is opt-in:

```vol
count := 1
count = 2          // OK

const limit := 10
limit = 11         // S030 / R030: Cannot assign to const binding

const a := [1, 2]
a = [3]            // S030 / R030: Cannot rebind
a[0] = 9           // OK (shallow const); array rules still apply
```

Remaining open items:

- `const` function parameters — Planned; parameters stay mutable for now.
- Deeper freeze / ownership of array elements — deferred to ownership design.

### Syntactic sugar from familiar languages (arrows + pipes)

**Blocked by Surface Freeze SF-0.** Desired surface sugar inspired by languages
used so far (JS arrow functions, PowerShell pipelines). Both are **Planned**
only — not accepted syntax. Prefer one canonical form per intent; do not ship
JS `=>` and PowerShell `|` and keep today’s `fn` / `.where` / `.sum` as three
competing spellings without a rule. Revisit only with an SF-1 bump.

#### Anonymous / arrow functions

Today: only named `fn name(params) { ... }`, plus `_` inside `.where` predicates.
Wanted: short first-class functions for callbacks and local helpers (JS feel).

Candidate shapes (pick one later; do not implement all):

```vol
// A — anonymous fn (same keyword as named functions)
double := fn(x) { return x * 2 }

// B — expression arrow (JS-like)
double := (x) => x * 2
add := (a, b) => a + b

// C — keep `_` for collection predicates; named/anonymous fn elsewhere
items.where(_ > 5)
```

Open questions:

- Prefer anonymous `fn` (A) over `=>` (B) unless JS familiarity is worth a second glyph.
- Expression body vs block body; required parens for a single parameter?
- Does arrow/`fn` replace `_` in `.where`, or coexist? Coexistence needs a clear
  canonical form for the formatter and language card.
- Closures already exist for nested named `fn`; anonymous forms must share the
  same capture and `nothing` / `return` rules (`SPEC.md` §5.8).

#### Pipeline sugar

Today: method chaining expresses collection intent —

```vol
nums.where(_ > 0).sum()
```

Wanted: PowerShell-like left-to-right pipelines for “pass this into the next
step.” Bare `|` is attractive from shells but conflicts with a future bitwise OR
in a systems language — prefer `|>` (or keep method chains only) unless `|` is
explicitly reserved for pipes and bitwise uses another spelling.

Candidate shapes:

```vol
// Method chain (already Supported — default until pipes are specified)
nums.where(_ > 0).sum()

// First-arg / pipeline sugar (Elixir/F#-like), not shell stdout
result := nums |> where(_ > 0) |> sum

// Bare | (PowerShell-like) — only if bitwise OR is deferred or respelt
result := nums | where(_ > 0) | sum
```

Open questions:

- Semantics: value as first argument (`x |> f` ⇒ `f(x)`), special method-name
  RHS (`x |> where(...)`), or collection-only pipeline?
- Eager vs lazy; interaction with `.where` purity and future `.map` / reduce.
- Formatter canonical form: method chain vs `|>` when both would be legal.
- Do not document `|` as Supported until bitwise-or policy is decided.

### Parallel intent

```vol
parallel {
    process(requests)
}
```

Automatic parallelization must preserve correctness and predictable resource use.
Exact semantics, scheduling, allocation, and failure behavior are undecided.

### Error / result model (design direction — not syntax)

**Direction decided; blocked for implementation by Surface Freeze SF-0** (and by
missing types / stdlib). Today the interpreter only **aborts** on failure
(`SPEC.md` §8): no `try`/`catch`, no error values. That interim behavior stays
until an SF bump ships the model below.

**Hybrid rule:**

1. **Traps (language / programmer mistakes)** — keep aborting with diagnostics:
   overflow (`R028`), bounds, type mistakes, div-by-zero, failed `assert`, and
   similar invariant breaks. These are not meant to be handled as normal values.
2. **Result values (expected operational failure)** — when I/O, parse, net, DB,
   and similar APIs exist, they return a Result-like success/failure value the
   program can inspect or propagate. Prefer one canonical Result shape over
   exceptions.
3. **No exceptions** as the primary error model (`try`/`catch` is not planned
   as the default).

**Canonical form (Planned spelling — pick exact names with types/SF-1):**

```vol
// Illustrative only — not accepted today
result := open("config.json")   // Result[File, Error] (or equivalent tags)
match result {
    ok file { /* use file */ }
    err e { /* handle or return err */ }
}
```

**Future dual-return sugar (optional, not canonical alone):**

Go-style multi-assign may be added later as **sugar over Result**, not a second
competing model. Ignoring the error binding must remain a diagnosable mistake
(same discipline as discarding an unchecked Result).

```vol
// Illustrative sugar — not accepted today; desugars to Result
file, err := open("config.json")
if err != none {
    return err
}
```

Open follow-ups (do not invent surface until types exist):

- Exact Result / error type spelling and whether `?` propagate is worth tokens.
- Trap vs Result boundary for decode / parse helpers.
- Unused-Result and ignored-`err` diagnostics.
- Interaction with `nothing` / `R029` (separate: missing return ≠ error value).
- Whether built-in `input` stays trap (`R024`) or becomes Result when I/O grows.

### Modules and imports

VOL may use folders as module boundaries, with project-root-relative imports or
aliases from `vol.config.json`:

```vol
import "services/users"
import "@db/models"
```

Symbol-selection syntax and the module resolver are not implemented.

### Block comments

Block-comment spelling has not been decided.

## Formatter

- [ ] Build a syntax-aware formatter using the compiler parser.
- [ ] Define one canonical style.
- [ ] Make formatting deterministic and idempotent.
- [ ] Preserve comments reliably.
- [ ] Support formatting a file or directory.
- [ ] Support a check-only mode for continuous integration.
- [ ] Make formatting fast enough to run on save.

Planned commands:

```text
vol fmt file.vol
vol fmt .
vol fmt --check .
```

Formatting must only change presentation, never program structure. A separate
future simplifier may recognize explicit algorithms and suggest equivalent VOL
semantic operations, while preserving both source styles as valid language forms
when they express different or equivalent intents under proven rules.

Possible workflow:

```text
vol simplify --check file.vol
vol simplify --diff file.vol
vol simplify --write file.vol
```

For example, it may suggest replacing an accumulator loop containing a filter
with `numbers.where(_ > 5).sum()`. Rewrites must be optional, deterministic, and
proven safe with respect to mutation, side effects, ordering, and error behavior.

## Projects and Modules

- [ ] Discover the nearest `vol.config.json` by walking from the source file toward the filesystem root.
- [ ] Resolve imports relative to the configured project root.
- [ ] Treat folders as module namespaces with one canonical resolution rule.
- [ ] Resolve path aliases such as `@db/*` without source-relative traversal.
- [ ] Reject aliases and imports that escape the project root unless explicitly allowed.
- [ ] Detect import cycles and report their complete path deterministically.
- [ ] Collect exports written anywhere in a module and format one sorted export list at the top.

Proposed configuration:

```json
{
  "name": "vol",
  "root": ".",
  "paths": {
    "@compiler": "compiler",
    "@std": "std"
  }
}
```

## Language Server

The language server must share the compiler's lexer, parser, syntax tree, semantic model, diagnostics, and source-location system.

Initial features:

- [ ] Live diagnostics.
- [ ] Document formatting and format on save.
- [ ] Hover information.
- [ ] Go to definition.
- [ ] Find references.
- [ ] Rename symbol.
- [ ] Completion.
- [ ] Document symbols.
- [ ] Semantic highlighting.

Future VOL-specific features:

- [ ] Display inferred types.
- [ ] Display ownership and lifetime decisions once those semantics exist.
- [ ] Show stack and heap allocation decisions.
- [ ] Warn about expensive or unbounded allocation.
- [ ] Explain compiler optimizations.
- [ ] Preview generated C or LLVM IR.
- [ ] Provide deterministic automatic fixes.

## Compiler Pipeline

```text
source
  -> lexer
  -> tokens
  -> parser
  -> AST
  -> semantic analysis
  -> HIR
  -> MIR
  -> optimization
  -> backend
```

Planned execution stages:

1. Tree-walking interpreter.
2. C backend.
3. LLVM backend.
4. Self-hosting compiler.

## Safety and Memory

Ownership, borrowing, and lifetimes are research goals, not current language
semantics. Before claiming inference, the specification must answer at least:

- Who owns returned memory?
- Can references escape a scope?
- Can two mutable aliases exist?
- What happens when inference is ambiguous?
- When is a value copied, moved, or reference-counted?
- How does foreign-function memory work?

Checklist:

- [ ] Define ownership and aliasing rules before inference claims.
- [ ] Prefer stack allocation when values do not escape.
- [ ] Avoid a mandatory garbage collector.
- [ ] Avoid hidden unbounded allocation.
- [ ] Provide explicit control when inference is insufficient.
- [ ] Make inferred allocation and ownership inspectable.
- [ ] Define selectable safety and optimization build modes.

Guiding rule once semantics exist:

> Infer by default. Expose control when necessary.

## Optimization

- [ ] Compile-time specialization.
- [ ] Whole-program dead-code elimination.
- [ ] Escape analysis.
- [ ] Automatic vectorization of bulk operations with defined semantics.
- [ ] Safe automatic parallelization.
- [ ] Data-oriented memory-layout optimization.
- [ ] Profile-guided optimization.
- [ ] Incremental compilation.

## Batteries-Included Library

Build broad capabilities on small, orthogonal foundations. Unused functionality must not increase executable size or runtime cost.

Priority areas:

- [ ] Filesystems and processes.
- [ ] Networking, HTTP, and TLS.
- [ ] Serialization and structured data.
- [ ] Databases and storage.
- [ ] Synchronous and asynchronous I/O.
- [ ] Concurrency and parallelism.
- [ ] Cryptography, hashing, and compression.
- [ ] Logging, tracing, and metrics.
- [ ] Configuration and command-line parsing.
- [ ] Testing, benchmarking, and profiling.
- [ ] C interoperability.
- [ ] Cross-compilation.

## Compiler Metrics and LLM Evaluation

Every future build may report:

- compile time
- binary size
- estimated memory use
- optimization level
- source token count for a documented tokenizer
- semantic density score
- Meaning per Token (MPT)

MPT is not defined yet. A useful definition must be objective, reproducible,
tokenizer-aware, and hard to game. Raw source token count is insufficient because
tokenizers differ across models and shorter source can increase repair cost.

Preferred experimental metric direction:

```text
task success / total tokens consumed
```

including generated code, compiler or runtime diagnostics, repair prompts, and
revisions across a fixed task suite. Protocol is defined in
[`LLM_BENCHMARK.md`](LLM_BENCHMARK.md). The protocol-v1 harness and one
Gemini core-suite result are implemented; more models and realistic backend
tasks are still required before making broad claims.

**Static source token density** (step 1 — implemented) is measured in
[`bench/`](bench/README.md): hand-written equivalent programs in VOL, Go, Rust,
and Zig are compared by token count under named OpenAI tokenizers
(`cl100k_base`, `o200k_base`). That benchmark measures source density only —
not LLM task-success efficiency.

## Open Design Questions

- ~~What is VOL's exact mutability model?~~ **Decided and implemented:** mutable by default;
  opt-in `const name := expr` (shallow; `S030`/`R030`). See `SPEC.md` §5.2.
- ~~Array assignment: shared, copy, or move?~~ **Decided and implemented:** shared references;
  `.copy()` (shallow) and `.deep_copy()` (recursive) are available. See `SPEC.md` §3.3.
- ~~Is `if` an expression, a statement, or both?~~ **Decided:** statement with
  `elif`/`else`; use `? :` for expression values. See `SPEC.md` §5.3 / §4.2.1.
- ~~What syntax should replace or represent a conventional `while` loop?~~
  **Decided:** keep `while` as Supported. See `SPEC.md` §5.5.
- ~~Should integer overflow trap, wrap, or use another rule outside the prototype?~~
  **Decided:** trap by default (`R028` + `fix`); wrapping/modes Planned.
  See `SPEC.md` §4.3.
- ~~Should string `.length` remain Unicode scalars, become bytes, or offer both?~~
  **Decided and implemented:** `.len` = Unicode scalars; `.byte_len` = UTF-8 bytes.
  See `SPEC.md` §3.2.
- ~~Which side effects are allowed inside `.where` predicates in a future compiler?~~
  **Decided:** none relied upon — pure filter; `.each` for effects; checks Planned.
  See `SPEC.md` §4.4.
- ~~Missing return value?~~ **Decided:** `nothing` if unused as a statement
  result; `R029` if assigned or used as a value. See `SPEC.md` §5.8.
- ~~How are fallible operations and errors represented?~~ **Direction decided
  (not implemented):** hybrid — traps for bugs/invariants; Result for expected
  operational failure; no exception-primary model; optional Go-style dual-return
  sugar later over Result. See “Error / result model” above.
- How are nullable or optional values represented?
- How are structs, enums, and tagged unions declared?
- How are modules and packages organized?
- Which standard features require explicit imports?
- How does a programmer constrain inferred allocation?
- What guarantees does automatic parallelization provide?
- Which build modes control bounds and overflow checking?
- Which ownership questions can be inferred locally, and which need API contracts?
- Anonymous functions: anonymous `fn`, JS-style `=>`, or only `_` in `.where`?
- Pipeline sugar: method chains only, `|>`, or bare `|` — and what is `|` vs bitwise OR?
