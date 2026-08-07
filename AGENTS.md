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
provisional syntax. There is no native backend, static type system, ownership
checker, or standard library. Keep vision and implementation status distinct in
every document.

Settled prototype rules (source of truth: [`SPEC.md`](SPEC.md) §11 Decided):

- bindings are **mutable by default**; opt-in `const` is Planned
- array assignment **shares** references; explicit clone is Planned
- integer overflow **traps** (`R028`) with a `fix` suggestion; wrapping modes are Planned
- `.where` predicates are **pure**; side effects belong in `.each`
- missing `return` yields `nothing`; discarding in a call statement is OK;
  assigning or using `nothing` as a value is `R029`
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

### 4. Compiler First (Research Goal)

A future compiler may perform analyses such as:

- ownership analysis
- borrow analysis
- lifetime inference
- escape analysis
- vectorization
- parallelization
- optimization

None of these are language semantics until specified. Local non-escaping
ownership inference may be tractable; removing ownership contracts from public
APIs is much harder. Do not document unimplemented inference as behavior.

### 5. Native Performance

VOL is intended to be a compiled language.

Performance targets:

- comparable to C
- memory safety
- zero-cost abstractions
- predictable execution

VOL programs should not require a garbage collector or heavyweight runtime.

A future compiler should provide:

- stack allocation whenever values do not escape
- compile-time specialization
- dead-code elimination across the standard library
- automatic vectorization and safe parallelization where semantics allow
- data-oriented memory layouts
- explicit control over allocation and representation when required
- optional safety checks according to build mode

### 6. AI Native

The compiler is designed to work with AI.

Goals:

- deterministic formatting
- deterministic diagnostics
- structured errors
- machine-readable compiler output
- falsifiable LLM generation and repair benchmarks

Errors should have both human and JSON formats. Diagnostics already carry
`code`, `message`, `file`, `position`, and optional `fix`. The CLI prints the
human form today; JSON emission from `vol` is Planned (see [`IDEAS.md`](IDEAS.md)).

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

Machine (`Diagnostic` JSON shape; CLI flag Planned):

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
[`SPEC.md`](SPEC.md). Prefer `const`-assignment errors once `const` is
implemented, not Rust-style “declare mutable” defaults.

---

## Language Goals

Priority:

1. Correctness
2. Readability
3. Token efficiency for successful generation and repair
4. Performance
5. Compile speed

Near-term project priority: freeze the implemented surface and specify it
precisely before adding major features. See `IDEAS.md`.

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
plus semantic compression (`.where(...).sum`) over inventing unfamiliar glyphs.

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

## Standard Library Philosophy

VOL is batteries included for systems and backend development.

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

## Runtime Philosophy

VOL provides high-level capabilities without hiding their cost.

Defaults should be safe and efficient. Systems programmers must be able to inspect and constrain compiler decisions when correctness, performance, memory layout, latency, or interoperability requires it.

The rule is:

> Infer by default. Expose control when necessary.

That rule applies only after inference has written semantic rules.

VOL must not impose:

- a mandatory garbage collector
- a heavyweight runtime
- hidden unbounded allocation
- unpredictable background work

Compiler reports and tooling should make inferred allocations, ownership, concurrency, and expensive operations visible.

---

## Compiler Metrics

Every build may report:

- compile time
- binary size
- memory estimate
- optimization level
- semantic density score
- token count for a documented tokenizer

Future metric:

Meaning per Token (MPT)

MPT is undefined until `IDEAS.md` / `LLM_BENCHMARK.md` give an objective,
reproducible definition. Prefer measuring task success against total tokens
consumed across generation and repair, not source length alone.

---

## Repository Structure

Current layout:

```text
cmd/vol
internal/lang
examples
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
- Interpreter

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
- `examples/*.vol`: Update or add a small executable example when a feature is
  important enough to teach users directly.
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
go run ./cmd/vol ./examples/first.vol
git diff --check
```

If a command cannot run, report that explicitly. A feature should not be described
as complete while its relevant tests are failing.

---

## Motto

Write less.

Express more.

Trust the compiler—once its decisions are specified.
