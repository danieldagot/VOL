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

## Motto

Write less.

Express more.

Trust the compiler.
