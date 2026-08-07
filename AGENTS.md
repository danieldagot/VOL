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

- Fewer tokens.
- Less boilerplate.
- More meaning.
- Native performance.

The compiler is responsible for implementation details.
The programmer describes intent.

VOL combines concise, intent-oriented source code with memory safety, predictable execution, and native performance.

---

## Core Principles

### 1. Intent over implementation

The programmer should describe **what** should happen.

The compiler decides **how**.

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

### 3. Zero Boilerplate

Avoid:

- `main()`
- imports for standard features
- explicit memory allocation
- repetitive types

The compiler should infer whenever possible.

### 4. Compiler First

The compiler performs:

- ownership analysis
- borrow analysis
- lifetime inference
- escape analysis
- vectorization
- parallelization
- optimization

These should not burden the programmer.

### 5. Native Performance

VOL is a compiled language.

Performance targets:

- comparable to C
- memory safety
- zero-cost abstractions
- predictable execution

VOL programs should not require a garbage collector or heavyweight runtime.

The compiler should provide:

- stack allocation whenever values do not escape
- compile-time specialization
- dead-code elimination across the standard library
- automatic vectorization and safe parallelization
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

Errors should have both human and JSON formats.

Example:

Human:

```text
error[E001]

Cannot borrow immutable value as mutable.

Suggestion:
Declare variable as mutable.
```

Machine:

```json
{
  "code": "E001",
  "severity": "error",
  "message": "Cannot borrow immutable value.",
  "fix": "Declare variable mutable."
}
```

---

## Language Goals

Priority:

1. Correctness
2. Readability
3. Token efficiency
4. Performance
5. Compile speed

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

Prefer:

```vol
parallel {
}
```

Instead of thread creation.

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

There should not be multiple competing ways to solve the same problem.

Unused library functionality must not increase executable size or runtime cost.
Only required functionality should be linked into the final program.

The library may be broad, but using it should remain token-efficient and conceptually small.

---

## Runtime Philosophy

VOL provides high-level capabilities without hiding their cost.

Defaults should be safe and efficient. Systems programmers must be able to inspect and constrain compiler decisions when correctness, performance, memory layout, latency, or interoperability requires it.

The rule is:

> Infer by default. Expose control when necessary.

VOL must not impose:

- a mandatory garbage collector
- a heavyweight runtime
- hidden unbounded allocation
- unpredictable background work

Compiler reports and tooling should make inferred allocations, ownership, concurrency, and expensive operations visible.

---

## Compiler Metrics

Every build reports:

- compile time
- binary size
- memory estimate
- optimization level
- semantic density score
- token count

Future metric:

Meaning per Token (MPT)

---

## Repository Structure

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

If yes, infer it by default.

Expose control only when the programmer needs to constrain behavior for correctness, performance, memory layout, latency, or interoperability.

---

## Documentation and Feature Synchronization

Language changes are incomplete until the implementation, tests, examples, and
documentation agree. Whenever syntax, a keyword, or language behavior is added,
removed, renamed, or changed, inspect and update every applicable file below.

### Required documentation updates

- `README.md`: Update the project-status lists and supported examples. Describe
  only behavior that works in the current implementation.
- `SYNTAX.md`: Define the grammar, semantics, restrictions, and canonical examples.
  Mark designs as provisional when their behavior is not stable.
- `VOCABULARY.md`: Add or update every keyword, contextual word, property, operator,
  and core symbol. Keep its Supported and Provisional status accurate.
- `IDEAS.md`: Record future work, unresolved design questions, formatter behavior,
  and functionality that has been designed but not implemented. Move completed
  behavior into the supported documentation instead of leaving contradictory plans.
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
- Keep Planned behavior in `IDEAS.md`; do not present it as executable syntax.
- Use one canonical spelling in documentation unless multiple forms are an explicit
  language feature.
- When multiple forms are supported, document whether they are semantically
  equivalent and which form the future formatter will emit.
- Keep diagnostic codes stable. Add tests for every new error code and suggestion.
- Preserve unrelated user changes when synchronizing documentation.

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

Trust the compiler.
