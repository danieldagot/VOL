# VOL Ideas and Future Work

This document collects ideas that may shape VOL in the future.

Items here are exploratory, not promises. Implemented language behavior belongs in
`SPEC.md`. Long-term vision without executable semantics stays here until it is
specified and tested.

## Near-Term Priority: Precise Core Before New Features

**Surface Freeze SF-1 is active** (see [`SPEC.md`](SPEC.md) §0). The vision-aligned
Supported surface is frozen. Harness card `bench/llm/cards/vol_v1.md` is the
`core_v2` task card (SF-1-bound subset — not a full SF-1 product tour). SF-0 /
`vol_v0` remains for historical harness tables. Do not expand the language until
foundations justify SF-2.

Frozen core (do not grow under SF-1 without a bump):

- values, variables (including multi-assign), named/anonymous/expression-body `fn`,
  product structs (named + positional literals)
- `if` / `elif` / `else`, if-let, `repeat`, `while`, `? :`, Option `??`, Result `?`
- arrays, `.each`, `.where`, `.map`, `.count`, `.sum()`, `some` / `none`, `ok` / `err`
- `import` / `export` with `vol.config.json` discovery
- built-ins already accepted by the interpreter (ambient tiny core)
- `match` is Rejected (`E153`)

Immediate documentation and design work:

- [x] Write `SPEC.md` covering lexical grammar, expression grammar, types/values,
      scopes, evaluation order, numeric behavior, arrays/strings, mutation,
      functions, control flow, and failure behavior for the implemented core.
- [x] Declare Surface Freeze SF-0 (SPEC + `vol_v0` language card).
- [x] Declare Surface Freeze SF-1 (vision-aligned surface; `vol_v1` = `core_v2` task card).
- [ ] Keep `SPEC.md` synchronized whenever interpreter behavior changes.
- [x] Error/result model — **Result values + if-let / `?` implemented** (SF-1);
      dual-return still Planned.
- [x] Option / optional-values — **implemented** (`some`/`none`/if-let/`??`); `?T` still Planned.
- [x] Modules — **implemented** (config + `path.vol` / `path/mod.vol`; ambient tiny core).
- [x] Phase-2 design directions #2–#11 (structs/Result/unwrap/density in SF-1;
      remaining ownership/alloc, parallel, build modes, pipelines).
- [ ] Foundations before SF-2: finish `vol fmt` rewriter (CLI stub + style rules
      drafted); richer std libraries behind imports (design only — aliases ≠ stdlib).
- [x] Write `LLM_BENCHMARK.md` with a falsifiable generate/repair protocol;
      harness + Gemini `core_v2` published — see checklist in that file.
- [ ] Keep Planned syntax out of `SPEC.md`; only Supported and Provisional forms belong there.
- [x] Publish frozen `core_v2` with default Python baseline (`--langs vol,python`).
- [ ] Publish ≥1 other model on `core_v2` before treating LLM results as stable.
- [x] Re-run / publish against `vol_v1` after source-check hygiene (`.count`/`.where`);
      published `20260808-041440` (`…040028` superseded); second model still needed.
- [x] **Density / unwrap surface (in SF-1)** — `.map` / `.count`, if-let /
      `??` / `?` shipped in SPEC. Shipped surface ≠ proven LLM workflow win;
      measure before further sugar (`LLM_BENCHMARK.md`).

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
- [x] Multi-assign / multi-declare (`a, b := …` / `a, b = …`) (SF-1). See `SPEC.md` §5.2.
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
- [x] `.map(_)` and `.count(_)` (SF-1).
- [ ] `.where` / `.map` / `.count` purity: intent guidance in `SPEC.md` §4.4;
      diagnostics / enforcement Planned (not decided-as-enforced).
- [ ] Purity diagnostics / checks for `.where` / `.map` / `.count` predicates.
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

### `.where` / `.map` / `.count` purity checks

**Planned (not enforced):** prefer side-effect-free predicates and transforms;
use `.each` for effects (`SPEC.md` §4.4). The prototype evaluates these
expressions eagerly and does not reject impurity — do not document purity as a
settled enforced rule until checks land.

Planned:

- diagnostics when a `.where` / `.map` / `.count` expression clearly has side
  effects (or cannot be proven pure), with a `fix` suggesting `.each` when
  appropriate
- allow clearly pure helper calls such as Boolean predicates over `_`
- treat purity as a prerequisite for any future fusion or parallel collection ops

### Wrapping integer arithmetic / overflow and build modes

Decided language rule (already in `SPEC.md` §4.3): integer overflow **traps** by
default (`R028`) with a human message and machine-readable `fix` field so LLMs
can help repair programs. Bounds checks likewise trap today (`R016`).

**Direction decided (not implemented):** keep trap as the default; later add
**build modes** and/or **explicit wrapping ops** in the same family:

- debug / default: trap on overflow and out-of-bounds (today’s behavior)
- optional release-wrap (or similar) mode for overflow when explicitly selected
- explicit wrapping operations for local opt-in without changing global mode
- bounds-checking policy across modes follows the same “safe by default, opt in
  to weaker checking” story — exact mode names TBD

Do not document silent wrap as Supported. JSON diagnostic output is implemented:
`vol --json run <file.vol>` emits a single JSON object on stderr (human form default).

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

#### Anonymous / arrow functions ✓ SF-1

**Implemented:** anonymous `fn(params) { ... }` and expression-body
`fn(params) expr` (see `SPEC.md` §5.8). Named `fn name(...)` and `_` in
`.where` / `.map` / `.count` remain Supported. Prefer anonymous `fn` over
JS-style `=>` (still rejected).

```vol
double := fn(x) { return x * 2 }
triple := fn(x) x * 3
items.where(_ > 5)
```

Open follow-ups:

- Canonical filter/map form: `_` vs `fn` in collection ops for the formatter/card

#### Pipeline sugar

Today: method chaining expresses collection intent (Supported):

```vol
nums.where(_ > 0).sum()
```

**Direction:** keep **method chains** under SF-1. If pipeline sugar is added at
a later freeze (SF-2+), use **`|>`** (first-arg / Elixir-like), not bare `|`.

```vol
// Planned — not accepted today
result := nums |> where(_ > 0) |> sum
```

Open follow-ups (if `|>` ships):

- Semantics: `x |> f` ⇒ `f(x)` vs special method-name RHS vs collection-only
- Eager vs lazy; interaction with `.where` / `.map` / `.count` purity and reduce
- Formatter canonical form when both chain and `|>` would be legal

### Parallel intent

**Direction decided (not implemented):** `parallel { ... }` remains Planned.
VOL claims **no parallelization guarantees** until scheduling, allocation,
ordering, and failure behavior are specified. Do not document automatic
parallelization as safe or deterministic today.

```vol
// Planned — not accepted today; semantics undefined until specified
parallel {
    process(requests)
}
```

### Error / result model ✓ SF-1 (values + if-let / `?`)

**Implemented:** Result values `ok(x)` / `err(x)` with if-let unwrap and postfix
`?` propagate inside functions (see `SPEC.md` §3.5 / §5.10). `match` removed
(`E153`). Hybrid rule remains:

1. **Traps** for bugs/invariants (overflow, bounds, `assert`, …).
2. **Result values** for operational success/failure data the program can inspect.
3. **No exceptions** as the primary model.

```vol
r := ok(7)
if ok n := r { print n } else err e { print e }
n := divide(10, 2)?
```

Still Planned:

- dual-return sugar over Result
- Unused-Result / ignored-`err` diagnostics
- Result-returning I/O (`input` stays trap `R024` for now)
- richer error types once tagged unions exist

### Optional / nullable values ✓ SF-1

**Implemented:** explicit Option with `some(x)` / `none`, if-let, and `??`
(see `SPEC.md` §3.4 / §5.10). Distinct from `nothing` (`R029`) and Result.
No universal `null`. Empty string / empty array remain ordinary data.
`match` removed (`E153`).

```vol
maybe := some("VOL")
if some name := maybe { print name } else { print "not found" }
print maybe ?? "guest"
```

Still Planned:

- light `?T` sugar (Option remains canonical)
- `or` / `or_else` / force-unwrap helpers
- nested Option flatten policy
- dual-return comma-ok for Option (not preferred)

### Structs ✓ SF-1; enums still Planned

**Implemented:** product `struct` with named-field and positional literals and
`.` access (see `SPEC.md` §3.6).

```vol
struct User {
    name
    age
}
u := User { name: "Ada", age: 36 }
v := User { "Ada", 36 }
```

Still Planned: methods, enums / tagged unions, pattern match on user tags.

### Modules and imports ✓ SF-1

**Implemented:** `import "path"` / `@alias`, nearest `vol.config.json`, resolve
to `path.vol` or `path/mod.vol`, live exports, cycle detection (see `SPEC.md`
§5.11).

```vol
import "services/users"
import "@db/models"
```

Still Planned: symbol-selection imports, multi-file folder packages beyond the
entry rule, sorted export-list formatter.

#### Standard-library import policy

**Direction decided (enforced for ambient surface):** ambient tiny core only —
`print`, `input`, `assert`, `string`, `args`, and collection methods in SPEC.
Everything else requires an **explicit import** when those libraries exist.

**Aliases ≠ stdlib:** `vol.config.json` `paths` (including example spellings like
`@std` or `@compiler` in this repo’s root config) are **project-local** alias
demos for the import feature. They are not a reserved standard-library root and
do not ship a product `std/` tree. A richer std behind imports remains Planned;
candidate first modules (when implemented): math helpers, strings, filesystem,
then net/serialization — always via explicit `import`, never ambient growth.

### Block comments

Block-comment spelling has not been decided.

## Formatter

- [x] CLI stub: `vol fmt [--check] <path>` parses and reports diagnostics; rewrite
      and format-equality check still exit with “not implemented” (see `cmd/vol`).
- [x] Draft canonical style rules (`SPEC.md` §12 + below).
- [ ] Build a syntax-aware rewriter using the compiler parser / AST.
- [ ] Make formatting deterministic and idempotent.
- [ ] Preserve comments reliably.
- [ ] Support formatting a directory (`.`).
- [ ] Support real `--check` format-equality for CI.
- [ ] Make formatting fast enough to run on save.

Canonical style (draft):

- four-space indent; K&R braces on `fn` / `if` / `while` / `repeat` / `.each` /
  `struct`
- one statement per line; spaces around binary ops and `:=` / `=`
- no space before `(` in calls; space after `,`
- prefer `_` over `fn` in `.where` / `.map` / `.count` when equivalent

Commands:

```text
vol fmt file.vol
vol fmt .
vol fmt --check file.vol
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

**Implemented (SF-1):** config discovery + `path.vol` / `path/mod.vol` + aliases.
Checklist:

- [x] Discover the nearest `vol.config.json` by walking from the source file toward the filesystem root.
- [x] Resolve imports relative to the configured project root.
- [x] Canonical resolution: `P.vol` else `P/mod.vol`.
- [x] Resolve path aliases such as `@db/*` without source-relative traversal.
- [x] Reject aliases and imports that escape the project root.
- [x] Detect import cycles and report their complete path deterministically.
- [ ] Collect exports written anywhere in a module and format one sorted export list at the top.

Proposed configuration (illustrative aliases only — not a shipped stdlib layout):

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

### Ownership direction (not implemented)

**Direction decided:** defer a full borrow checker. Prefer **local escape
analysis** first (stack when values do not escape); require **explicit API
contracts** for ownership across public boundaries later. Today’s shallow `const`
and shared array references (`SPEC.md` §3.3 / §5.2) stay as documented — they are
not ownership or move semantics.

Do not claim ownership, borrowing, or lifetime inference in docs or diagnostics
until those rules are written and tested.

### Allocation control (not implemented)

**Direction decided:** allocation policy is **unspecified** for the interpreter
prototype. Do **not** claim that the compiler infers allocation. Explicit
allocators / layout constraints stay research until a native backend exists.
When inference is specified later, expose control when programmers must constrain
layout, latency, or FFI.

Questions that must be answered before inference claims:

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
- [ ] Define selectable safety and optimization build modes (see wrapping /
      overflow section above).

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
[`LLM_BENCHMARK.md`](LLM_BENCHMARK.md). The protocol-v1 harness and published
Gemini `core_v2` results exist; after source-check hygiene the `vol_v1` task-card
run matches Python on first-try / success @ K with modest cold overhead (warm
slightly under Python once the card is amortized). That is still one model on a
tiny suite — not a broad LLM claim. Do not treat SF-1 itself as proven
optimization.

**Static source token density** (step 1 — implemented) is measured in
[`bench/`](bench/README.md): hand-written equivalent programs in VOL, Go, Rust,
and Zig are compared by token count under named OpenAI tokenizers
(`cl100k_base`, `o200k_base`). That benchmark measures source density only —
not LLM task-success efficiency. Never cite density ratios as workflow proof.

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
  **Planned:** prefer none — filter/transform intent vs `.each` for effects;
  not enforced today. See `SPEC.md` §4.4 and “`.where` purity checks” above.
- ~~Missing return value?~~ **Decided:** `nothing` if unused as a statement
  result; `R029` if assigned or used as a value. See `SPEC.md` §5.8.
- ~~How are fallible operations and errors represented?~~ **Implemented
  (SF-1):** hybrid traps + `ok`/`err` + if-let / `?`; no exception-primary
  model; dual-return sugar still Planned. See “Error / result model” above.
- ~~How are nullable or optional values represented?~~ **Decided and implemented
  (SF-1):** explicit Option (`some`/`none`/if-let/`??`); distinct from
  Result and `nothing`; optional `?T` sugar later. See `SPEC.md` §3.4 / §5.10.
- ~~How are structs, enums, and tagged unions declared?~~ **Structs implemented
  (SF-1);** enums/tagged unions still Planned. See “Structs” above.
- ~~How are modules and packages organized?~~ **Implemented (SF-1):**
  `import` + `vol.config.json` + `path.vol`/`mod.vol`. See “Modules and imports”.
- ~~Which standard features require explicit imports?~~ **Direction decided:**
  ambient tiny core only; everything else explicit import. See “Standard-library
  import policy” above.
- ~~Result / recoverable errors?~~ **Implemented (SF-1):** `ok`/`err` +
  if-let / `?`; dual-return still Planned. See “Error / result model” above.
- ~~How does a programmer constrain inferred allocation?~~ **Direction decided:**
  unspecified for now; do not claim inference. See “Allocation control” above.
- ~~What guarantees does automatic parallelization provide?~~ **Direction
  decided:** none until scheduling/failure are specified; `parallel` stays
  Planned. See “Parallel intent” above.
- ~~Which build modes control bounds and overflow checking?~~ **Direction
  decided (not implemented):** trap default; modes + explicit wrap ops later.
  See “Wrapping integer arithmetic / overflow and build modes” above.
- ~~Which ownership questions can be inferred locally, and which need API
  contracts?~~ **Direction decided (not implemented):** local escape first;
  contracts for public APIs later. See “Ownership direction” above.
- ~~Anonymous functions: anonymous `fn`, JS-style `=>`, or only `_` in `.where`?~~
  **Decided and implemented (SF-1):** anonymous `fn` + expression body; no
  `=>`; `_` in `.where` / `.map` / `.count` kept. See `SPEC.md` §5.8.
- ~~Pipeline sugar: method chains only, `|>`, or bare `|` — and what is `|` vs
  bitwise OR?~~ **Direction decided (not implemented):** chains under SF-1; if
  pipes later (SF-2+), `|>` not bare `|`. See “Pipeline sugar” above.
