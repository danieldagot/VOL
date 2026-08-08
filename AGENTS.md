# VOL — Vocabulary Optimized Language

> Write intent. Compile excellence.

## Vision

VOL is an experimental systems and backend programming language designed for humans and Large Language Models.

VOL targets the problem space served by languages such as Rust, Go, Zig, and C:

- backend services and APIs
- networking and distributed systems
- databases and storage engines
- operating-system tools and infrastructure
- command-line applications
- compilers and developer tooling
- embedded and low-level software
- security-sensitive software

The goal is to maximize **semantic density**:

- Fewer tokens when that improves task success.
- Less boilerplate.
- More meaning.
- Native performance.

The programmer describes intent. A future compiler should handle implementation
details that do not need to consume source tokens—**after** those details have
precise semantics. Inference is a design aspiration, not a description of the
current interpreter.

VOL aims to combine concise, intent-oriented source code with memory safety,
predictable execution, and native performance.

**Current reality:** VOL is a tree-walking interpreter prototype with a small
provisional syntax. **Surface Freeze SF-1** pins that vision-aligned Supported
surface ([`SPEC.md`](SPEC.md) §0). The harness card `bench/llm/cards/vol_v1.md`
is a `core_v2` **task card** (SF-1-bound subset — not a full SF-1 product tour).
SF-0 / `vol_v0` remains for historical tables. Do not add Planned syntax under
SF-1 without a freeze bump. SF-1 is a **syntax pin**, not proof of LLM workflow
superiority. Hand-written source-density numbers in `bench/` are not
generate/repair evidence. There is no native backend, static type system,
ownership checker, or broad standard library (ambient tiny core only). Ownership,
vectorization, and parallelization are vision/research only until specified and
tested. Keep vision and implementation status distinct in every document.

Settled prototype rules (source of truth: [`SPEC.md`](SPEC.md) §11 Decided):

- bindings are **mutable by default**; opt-in `const name := expression` is Supported (shallow; `S030`/`R030` on reassignment)
- multi-assign `a, b := …` / `a, b = …` is Supported (RHS fully evaluated before assigns)
- array assignment **shares** references; use `.copy()` (shallow) or `.deep_copy()` (recursive) for isolation
- integer overflow **traps** (`R028`) with a `fix` suggestion; wrapping modes are Planned
- `.where` predicates are **pure**; side effects belong in `.each`; `.map` / `.count` are Supported
- missing `return` yields `nothing`; discarding in a call statement is OK;
  assigning or using `nothing` as a value is `R029`
- Option uses `some` / `none` with if-let and `??`; Result uses `ok` / `err`
  with if-let and postfix `?` (functions only); both distinct from `nothing`;
  bugs still trap; `match` is Rejected (`E153`)
- anonymous `fn(params) { ... }` and expression-body `fn(params) expr` are Supported (`=>` is not)
- product `struct` with named and positional literals and `.` field access is Supported
- `import "path"` + `vol.config.json` discovery/aliases; `export` is live across modules
- `.len` is the length property (string: Unicode scalars); `.length` is rejected
- `while` is permanent Supported vocabulary (with `repeat` and `.each` for other intents)
- `if` is a statement (`elif` / `else`); value choice uses `? :` (not expression-`if`)

---

## Core Principles

### 1. Intent over implementation

The programmer should describe **what** should happen.

Where semantics allow, the compiler decides **how**.

Bad:

```c
for(int i=0;i<n;i++)
```

Good:

```vol
items.each item {
}
```

### 2. Every token matters

Every keyword must justify its existence.

Removing one token should never reduce readability.

Syntax should be optimized for both humans and LLMs.

Token efficiency is not the same as shortest source text. Prefer forms that
models generate reliably and that reduce total repair cost. Tokenizer differences
across models mean absolute token counts are never universal.

### 3. Zero Boilerplate

Avoid:

- `main()`
- imports for standard features
- explicit memory allocation when inference is sound
- repetitive types when inference is sound

Infer only when the language specification defines what inference means.

### 4. Compiler First (Research Goal — not implemented)

**Vision only.** A future compiler *might* perform analyses such as ownership,
borrow/lifetime inference, escape analysis, vectorization, parallelization, or
other optimization — **after** each has written semantic rules and tests. None
of these exist in the current interpreter. Do not document them as behavior,
diagnostics, or implied SF-1 capabilities. Local non-escaping ownership may be
tractable later; removing ownership contracts from public APIs is much harder.

### 5. Native Performance (research target — not implemented)

VOL is **intended** to become a compiled language. Targets below are aspirations,
not properties of today's tree-walker:

- comparable to C
- memory safety
- zero-cost abstractions
- predictable execution
- no mandatory garbage collector or heavyweight runtime

A future native backend *might* pursue stack allocation for non-escaping values,
specialization, DCE, vectorization/parallelization **where semantics allow**,
layout control, and selectable safety checks — only after those semantics are
specified. Do not imply any of that from the vision section alone.

### 6. AI Native

Tooling and diagnostics should be easy for humans and agents to consume. That is
a design goal; SF-1 surface + early harness runs do not prove workflow superiority.

Goals:

- deterministic formatting (Planned)
- deterministic diagnostics (human + JSON — Supported today)
- structured errors with stable codes (Supported today)
- machine-readable compiler/interpreter output (`vol --json` — Supported)
- falsifiable LLM generation and repair benchmarks (protocol + early runs;
  not a proven VOL win yet)

Errors have both human and JSON formats. Diagnostics carry `code`, `message`,
`file`, `position`, and optional `fix`. The CLI prints the human form by default;
pass `--json` anywhere in the command (`vol --json run <file.vol>` or
`vol run --json <file.vol>`) to receive a single JSON object on stderr instead.

Use examples that match **decided** language rules. Do not imply
immutable-by-default bindings or silent integer wrapping.

Example (implemented overflow trap):

Human:

```text
error[R028] demo.vol:1:27

Integer overflow in `+`.

   1 | print 9223372036854775807 + 1
     |                           ^

Suggestion:
Use smaller values, switch to floating-point, or wait for planned wrapping arithmetic / build modes.
```

Machine (`Diagnostic` JSON shape; `vol --json`):

```json
{
  "code": "R028",
  "message": "Integer overflow in `+`.",
  "file": "demo.vol",
  "position": {
    "Offset": 26,
    "Line": 1,
    "Column": 27
  },
  "fix": "Use smaller values, switch to floating-point, or wait for planned wrapping arithmetic / build modes."
}
```

Ownership or borrow diagnostics belong only after those semantics exist in
[`SPEC.md`](SPEC.md). Prefer `const`-assignment errors (`S030`/`R030`) for
immutable bindings, not Rust-style “declare mutable” defaults.

---

## Language Goals

Priority:

1. Correctness
2. Readability
3. Token efficiency for successful generation and repair
4. Performance
5. Compile speed

Near-term project priority: **SF-1 is active** — keep the frozen surface precise
(tests, diagnostics, docs) and build foundations (formatter; richer std via
imports) before an SF-2 feature bump. Remaining directions (enums, dual-return
sugar, ownership/alloc, build modes, `|>`, parallel) are sketched in
[`IDEAS.md`](IDEAS.md). See `SPEC.md` §0 / §11.

---

## Syntax Philosophy

Prefer:

```vol
repeat 10 {
}
```

Instead of:

```c
for(...)
```

Prefer:

```vol
users.each user {
}
```

Instead of iterator boilerplate.

Prefer familiar, trainable surface syntax (`:=`, braces, `fn`, `if`, arrays)
plus semantic compression (`.where(...).sum()`) over inventing unfamiliar glyphs.

Keep about one canonical representation for each distinct intent. Imperative
iteration and pure filtering or aggregation are different intents; both may exist.

Planned concurrency intent (not implemented):

```vol
parallel {
}
```

Instead of manual thread creation—only after scheduling, allocation, and failure
semantics are defined.

---

## Standard Library Philosophy (vision — not shipped)

**Aspiration:** batteries included for systems and backend development. Today only
an ambient tiny core exists; richer libraries require explicit imports when added.

The foundations should be small, orthogonal, and predictable while the available capabilities should be broad.

Primary capabilities include:

- filesystems, processes, and operating-system APIs
- networking, HTTP, and TLS
- serialization and structured data
- databases and storage
- concurrency, parallelism, and asynchronous I/O
- cryptography, hashing, and compression
- logging, tracing, and metrics
- configuration and command-line parsing
- testing, benchmarking, and profiling
- C interoperability and cross-compilation

Prefer one canonical library approach for each problem. Distinct intents may still
have distinct forms in the language.

Unused library functionality must not increase executable size or runtime cost.
Only required functionality should be linked into the final program.

The library may be broad, but using it should remain token-efficient and conceptually small.

---

## Runtime Philosophy (research target)

**Aspiration:** high-level capabilities without hiding their cost. Today's
interpreter does not infer ownership, allocation, or parallelism.

Defaults should be safe and efficient once a native runtime exists. Systems
programmers must be able to inspect and constrain compiler decisions when
correctness, performance, memory layout, latency, or interoperability requires it.

The rule (only after inference has written semantic rules) is:

> Infer by default. Expose control when necessary.

VOL must not impose:

- a mandatory garbage collector
- a heavyweight runtime
- hidden unbounded allocation
- unpredictable background work

When (and only when) inference rules exist, tooling should make inferred
allocations, ownership, concurrency, and expensive operations visible.

---

## Compiler Metrics

Every build may report:

- compile time
- binary size
- memory estimate
- optimization level
- semantic density score
- token count for a documented tokenizer

**Static source token density** (implemented today) is measured in
[`bench/`](bench/README.md): equivalent hand-written programs in VOL, Go, Rust,
and Zig are compared by token count under named OpenAI tokenizers. That measures
source size only. Do **not** cite density ratios as LLM workflow proof, SF-1
success, or task-success efficiency. See [`README.md`](README.md) for current
ratios.

**Generate/repair workflow** protocol is in [`LLM_BENCHMARK.md`](LLM_BENCHMARK.md)
(task success vs total tokens across generation, diagnostics, and repair). Early
published `core_v2` runs exist; they do **not** establish VOL superiority. The
current Gemini table (`vol_v1` task card) matches Python on first-try / success
@ K with modest cold overhead after source-check hygiene. Meaning per Token
(MPT) remains undefined as a single scalar. More models and realistic tasks are
required before treating results as stable.

---

## Repository Structure

Current layout:

```text
cmd/vol          CLI entry point
internal/lang    lexer, parser, AST, resolver, interpreter, diagnostics
examples         executable VOL programs
bench            source token density benchmark (VOL vs Go/Rust/Zig)
```

Aspirational layout as the project grows:

```text
/compiler
    lexer
    parser
    ast
    hir
    mir
    optimizer
    backend

/std

/docs

/examples

/tests

/spec

/tools
```

---

## Backend Roadmap

### Phase 1

- Lexer
- Parser
- AST
- Resolver
- Interpreter
- Structured diagnostics (human + JSON)

### Phase 2

- C backend

### Phase 3

- LLVM backend

### Phase 4

- Self-hosting compiler

---

## Non-Goals

VOL is **not**:

- another Rust
- another Go
- another C replacement
- a game-focused language
- a scripting language with a mandatory virtual machine

VOL explores a new programming paradigm: programming by intent.

Its primary domain is safe, high-performance systems and backend software.

---

## Design Rule

Whenever adding a feature, ask:

> Can the compiler infer this?

If yes, and the specification defines how, infer it by default.

Expose control only when the programmer needs to constrain behavior for correctness, performance, memory layout, latency, or interoperability.

If inference cannot yet be specified, do not claim it.

---

## Documentation and Feature Synchronization

Language changes are incomplete until the implementation, tests, examples, and
documentation agree. Whenever syntax, a keyword, or language behavior is added,
removed, renamed, or changed, inspect and update every applicable file below.

### Required documentation updates

- `README.md`: Update the project-status lists and supported examples. Describe
  only behavior that works in the current implementation. Keep vision clearly
  labeled as not implemented.
- `SPEC.md`: Single source for vocabulary, syntax, formal lexical rules, grammar,
  value semantics, evaluation order, and failure behavior. Update it for every
  implemented change. Mark Provisional when spelling or meaning may change. Do
  not include Planned syntax here.
- `IDEAS.md`: Record future work, unresolved design questions, formatter behavior,
  and functionality that has been designed but not implemented. Move completed
  behavior into `SPEC.md` instead of leaving contradictory plans.
- `AGENTS.md`: Keep vision examples and “current reality” bullets consistent with
  [`SPEC.md`](SPEC.md) decided rules. Do not show diagnostic samples that imply
  undecided or rejected defaults (for example immutable-by-default or silent
  integer wrap).
- `examples/`: Update or add a small executable example when a feature is
  important enough to teach users directly. Prefer `examples/basics/` (tour),
  `examples/features/` (one-topic demos), or `examples/projects/` (multi-file
  apps). Keep [`examples/README.md`](examples/README.md) in sync.
- `vol.config.json`: Update only when project discovery, roots, aliases, or other
  project-level configuration changes.

### Required implementation checks

For a new or changed keyword or syntax form, inspect all relevant compiler layers:

- `internal/lang/token.go` for token kinds.
- `internal/lang/lexer.go` for keyword recognition and lexical rules.
- `internal/lang/parser.go` for grammar and deterministic diagnostics.
- `internal/lang/ast.go` for syntax-tree representation.
- `internal/lang/interpreter.go` for current prototype behavior.
- `internal/lang/lang_test.go` for successful cases, failure cases, source locations,
  and interaction with existing syntax.

When the compiler gains HIR, MIR, semantic analysis, formatters, module resolution,
or native backends, update those layers as well. Do not treat parser acceptance as
complete feature support when later stages cannot execute or compile the construct.

### Status and consistency rules

- Never label a form Supported unless it is implemented and covered by tests.
- Use Provisional for implemented behavior whose spelling or semantics may change.
- Keep Planned behavior in `IDEAS.md`; do not present it as executable syntax in
  `SPEC.md`.
- Prefer one canonical representation for each distinct intent. When multiple forms
  exist for different intents, say so explicitly.
- When multiple equivalent forms are supported, document whether they are
  semantically equivalent and which form the future formatter will emit.
- Keep diagnostic codes stable. Add tests for every new error code and suggestion.
- Preserve unrelated user changes when synchronizing documentation.
- Do not claim ownership, borrowing, lifetimes, vectorization, parallelization, or
  LLM token superiority without written semantics or measured evidence.

### Verification before completion

Run:

```text
gofmt -w <changed Go files>
go test ./...
go run ./cmd/vol run ./examples/basics/first.vol
git diff --check
```

If a command cannot run, report that explicitly. A feature should not be described
as complete while its relevant tests are failing.

---

## Motto

Write less.

Express more.

Trust the compiler—once its decisions are specified.
