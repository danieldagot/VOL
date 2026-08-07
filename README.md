# VOL

> **Vector-Oriented Language**  
> Also known as the **Vibe-Oriented Language**.

VOL is an experimental systems and backend programming language designed for humans and Large Language Models.

The project explores a simple idea:

> Write intent. Compile excellence.

Most programming languages were designed before code-generating Large Language Models existed. Their syntax, tooling, diagnostics, and abstractions were optimized primarily for human-written software and historical conventions.

VOL asks what a systems programming language should look like when LLMs are first-class programmers.

Its primary goal is to maximize **semantic density**: express more intent with fewer tokens while preserving deterministic structure, readability, safety, and native performance. The language should be efficient for an LLM to generate, understand, review, modify, and debug.

Compiler inference supports that goal. The programmer or agent describes **what** should happen, and the compiler handles implementation details that do not need to consume source tokens. When correctness, performance, memory layout, or interoperability requires control, those decisions remain inspectable and constrainable.

VOL is designed around:

- semantic density: fewer tokens with more meaning
- syntax designed for reliable LLM generation and comprehension
- deterministic structure that is easy for tools and agents to transform
- structured human-readable and machine-readable diagnostics
- memory safety without a mandatory garbage collector
- native, predictable performance
- compiler-inferred ownership, lifetimes, allocation, and optimization
- batteries-included systems and backend capabilities
- deterministic formatting and compiler output
- explicit control when correctness, performance, or interoperability requires it

VOL is not intended to be another spelling of Rust, Go, or C. It is an AI-native language experiment targeting the same systems and backend problem space.

## Project Status

**VOL is currently a very early tree-walking interpreter prototype.**  
Most of the language described above does not exist yet. There is no compiler, no ownership checker, no borrow checker, no type system, no backend, and no standard library. What exists is a minimal interpreter that can execute small programs written in a provisional syntax.

This README describes both the current prototype and the long-term direction. Treat the vision as a design target, not a shipped reality.

### What Actually Works

- integer, floating-point, Boolean, and string literals
- inferred variable declarations with `:=`
- variable assignment with `=`
- arithmetic operators: `+`, `-`, `*`, `/`
- comparison operators: `==`, `!=`, `<`, `<=`, `>`, `>=`
- Boolean operators: `and`, `or`, `not`
- brace-delimited lexical scopes
- `if` and `else`
- `repeat` loops
- provisional `while` loops
- arrays and array indexing
- indexed array assignment
- array and string `.length`
- array iteration with `.each`
- `print`
- line comments beginning with `//`
- source-located parser and runtime diagnostics
- automated interpreter tests

### What Is Missing

- symbol resolution before execution
- static type checking
- explicit mutability rules
- functions, parameters, and return values
- improved diagnostic suggestions
- canonical VOL formatter
- JSON compiler diagnostics
- expanded conformance tests
- initial language-server support

### What Is Still an Open Experiment

- final function and call syntax
- explicit type syntax
- ownership and borrowing syntax
- error propagation
- optional and nullable values
- structs, methods, enums, and tagged unions
- modules, packages, and imports
- generics
- pattern matching
- asynchronous I/O syntax
- allocator and memory-layout constraints
- compile-time evaluation
- concurrency and automatic parallelization semantics
- bounds-checking and overflow behavior across build modes
- C and LLVM backend details

See [`SYNTAX.md`](SYNTAX.md) for the current syntax direction and [`IDEAS.md`](IDEAS.md) for future work and open design questions.

## Supported Examples

### Values and Variables

```vol
name := "VOL"
version := 1
ready := true

print name
print version
print ready
```

### Arithmetic and Assignment

```vol
total := 10 + 5 * 2
total = total - 4

print total
```

### Conditions

```vol
score := 75

if score >= 50 {
    print "pass"
} else {
    print "fail"
}
```

### Repeat

```vol
count := 0

repeat 3 {
    count = count + 1
}

print count
```

### Arrays and Iteration

```vol
numbers := [4, 7, 2, 9]
total := 0

numbers.each number {
    if number > 5 {
        total = total + number
    }
}

print total
```

### Indexing and Length

```vol
items := [1, 2, 3]
items[1] = 8

print items
print items.length
```

### Boolean Logic

```vol
ready := true
failed := false

if ready and not failed {
    print "start"
}
```

## Repository Layout

```text
cmd/vol          command-line entry point
internal/lang    lexer, parser, AST, diagnostics, and interpreter
examples         example VOL programs
```

## Requirements

- Go 1.24 or newer

## Run the Example

```text
go run ./cmd/vol ./examples/first.vol
```

Expected output:

```text
16
```

## Build the CLI

```text
go build ./cmd/vol
./vol ./examples/first.vol
```

On Windows:

```text
go build ./cmd/vol
.\vol.exe .\examples\first.vol
```

## Run the Tests

```text
go test ./...
```