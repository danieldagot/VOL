# VOL

> **Vocabulary Optimized Language**

VOL is an experimental systems and backend programming language designed for humans and Large Language Models.

The project explores a simple idea:

> Write intent. Compile excellence.

Most programming languages were designed before code-generating Large Language Models existed. Their syntax, tooling, diagnostics, and abstractions were optimized primarily for human-written software and historical conventions.

VOL asks what a systems programming language should look like when LLMs are first-class programmers.

Its primary goal is to maximize **semantic density**: express more intent with fewer tokens while preserving deterministic structure, readability, safety, and native performance. The language should be efficient for an LLM to generate, understand, review, modify, and debug.

That goal is a research direction, not a measured property of the current prototype. Token counts also depend on each model's tokenizer, so “fewer tokens” alone is not a success metric. The intended long-term measure is task success relative to total tokens consumed across generation, diagnostics, and repair. See [`IDEAS.md`](IDEAS.md).

## Project Status

**VOL is currently a very early tree-walking interpreter prototype.**
There is no native compiler, no static type system, no ownership or borrow checker, no standard library, and no backend. What exists is a minimal interpreter that can execute small programs written in a provisional syntax.

This README separates **what works today** from **long-term design targets**. Treat vision language as a target, not a shipped reality.

### Design Targets (Not Implemented)

- semantic density with measurable LLM workflow metrics
- deterministic structure that is easy for tools and agents to transform
- structured human-readable and machine-readable diagnostics
- memory safety without a mandatory garbage collector
- native, predictable performance
- compiler-assisted ownership, allocation, and optimization where semantics can be specified
- batteries-included systems and backend capabilities
- deterministic formatting and compiler output
- explicit control when correctness, performance, or interoperability requires it

### What Actually Works

**Surface Freeze SF-0:** the Supported Prototype v0 syntax in [`SPEC.md`](SPEC.md)
is frozen. New vocabulary (arrows, pipes, `.count`, etc.) waits for an SF-1 bump.
LLM card: [`bench/llm/cards/vol_v0.md`](bench/llm/cards/vol_v0.md).

- integer, floating-point, Boolean, and string literals
- inferred variable declarations with `:=`
- variable assignment with `=`
- opt-in immutable bindings with `const name := expression` (shallow; rebind blocked with `S030`/`R030`; indexed writes on array bindings still allowed)
- arithmetic operators: `+`, `-`, `*`, `/` (integer overflow traps with `R028`)
- comparison operators: `==`, `!=`, `<`, `<=`, `>`, `>=`
- Boolean operators: `and`, `or`, `not`
- brace-delimited lexical scopes
- `if` / `elif` / `else` statements
- conditional operator `cond ? a : b`
- `repeat` loops
- `while` loops
- arrays and array indexing
- indexed array assignment (assignment and arguments share array identity; use `.copy()` for a shallow clone or `.deep_copy()` for a recursive clone)
- array `.len` and string `.len` (Unicode scalar values) and string `.byte_len` (UTF-8 byte count)
- array iteration with `.each` (imperative / side effects)
- collection filtering with `.where(...)` (eager new array; pure predicates) and numeric aggregation with `.sum()`
- `print`
- interactive text input with `input()` or `input(prompt)`
- runtime checks with `assert(condition)` or `assert(condition, message)`
- value-to-string conversion with `string(value)`
- command-line arguments through the built-in `args` array
- functions declared with `fn`, parameters, calls, and `return` (missing return is `nothing`; using it as a value is `R029`)
- module export lists that may appear before or after definitions
- line comments beginning with `//`
- source-located parser and runtime diagnostics with human and JSON output (`vol --json run <file.vol>`)
- table-driven lexer, parser, resolver, interpreter, diagnostic, CLI, and executable-example tests
- fuzz seed coverage ensuring arbitrary parser input does not panic
- continuous integration with formatting, vet, race-detection, and coverage checks

### What Is Missing

- static type checking
- richer diagnostic suggestions across more error codes
- canonical VOL formatter
- a published compatibility policy and versioned conformance corpus (SF-0 is the first surface pin)
- broader LLM workflow results across more models and realistic backend tasks
- initial language-server support

### What Is Still an Open Experiment

- typed function signatures
- explicit type syntax
- ownership and borrowing semantics (local inference vs API contracts)
- error propagation / Result values (direction: hybrid traps + Result, optional
  dual-return sugar later — see [`IDEAS.md`](IDEAS.md); not implemented)
- optional and nullable values
- structs, methods, enums, and tagged unions
- modules, packages, and imports
- generics
- pattern matching
- asynchronous I/O syntax
- allocator and memory-layout constraints
- compile-time evaluation
- concurrency and automatic parallelization semantics
- overflow wrapping ops / build modes (default is already trap; see `SPEC.md` §4.3)
- bounds-checking behavior across build modes
- C and LLVM backend details

Near-term priority is to freeze the current surface and make its semantics precise before adding new language features. See [`IDEAS.md`](IDEAS.md).

See [`SPEC.md`](SPEC.md) for vocabulary, syntax, and formal behavior of the
current interpreter, and [`IDEAS.md`](IDEAS.md) for future work and open design
questions.

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

if score >= 90 {
    print "great"
} elif score >= 50 {
    print "pass"
} else {
    print "fail"
}

label := score >= 50 ? "pass" : "fail"
print label
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

Filtering and summing are a different intent from imperative iteration:

```vol
total := numbers.where(_ > 5).sum()
print total
```

### Indexing and Length

```vol
items := [1, 2, 3]
items[1] = 8

print items
print items.len
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
bench            source token density benchmark (VOL vs Go/Rust/Zig)
```

## Source Token Density Benchmark

On 13 equivalent small programs, VOL currently uses about:

- **~36% fewer tokens than Go**
- **~25% fewer tokens than Rust**
- **~53% fewer tokens than Zig**

(median across the suite; nearly the same under both `cl100k_base` and
`o200k_base`)

That is source size only — hand-written programs that print the same output. It
does **not** measure whether an LLM generates correct VOL more easily, or how
many tokens generate/repair workflows consume.

Full per-task numbers: [`bench/results/density.md`](bench/results/density.md).
How to run / regenerate: [`bench/README.md`](bench/README.md).

```text
cd bench && uv sync && uv run python harness/count_tokens.py
```

## LLM Workflow Benchmark

The default workflow baseline is **Python** (interpreted peer for the current
prototype). Go remains an optional compiled baseline (`--langs vol,go`).

One protocol-v1.1 (`core_v2`) run is published for `gemini-3.5-flash-lite`,
temperature 0, three replicates, and at most two repair rounds — **VOL vs
Python**. The suite has five tasks covering generation, diagnostic-seeded
repair, and modification. Summaries split prompt vs completion and report cold
totals (card re-sent every request) plus warm totals (estimated language-card
cost amortized away).

| Language | First-try | Success @ K | Mean cold | Mean warm | Mean completion |
| --- | ---: | ---: | ---: | ---: | ---: |
| Python | 100% | 100% | 763.8 | 427.8 | 136.0 |
| VOL | 100% | 100% | 839.9 | 406.9 | 112.0 |

VOL used about 10.0% more **cold** workflow tokens and about 4.9% fewer **warm**
tokens than Python. Completions were about 17.6% smaller; prompts were about
15.9% larger (language-card teaching cost). Every replicate succeeded on the
first try, including diagnostic repair (`R007` for VOL).

An earlier VOL-vs-Go `core_v2` run is kept for comparison. This small synthetic
suite does not establish real-world superiority for any language and does not
measure runtime performance. A second model on `core_v2` is still needed before
treating results as stable.

Full result (VOL vs Python): [`bench/llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-022642.md`](bench/llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-022642.md).
Earlier VOL vs Go: [`bench/llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-021122.md`](bench/llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-021122.md).
Historical non-diagnostic `core_v1`: [`bench/llm/results/core_v1_live_gemini_gemini-3.5-flash-lite_20260808-014539.md`](bench/llm/results/core_v1_live_gemini_gemini-3.5-flash-lite_20260808-014539.md).
Protocol and limitations: [`LLM_BENCHMARK.md`](LLM_BENCHMARK.md).

## Requirements

- Go 1.24 or newer

## Run the Example

```text
go run ./cmd/vol run ./examples/first.vol
```

Expected output:

```text
16
```

## Example Programs

The `examples` folder contains small programs that can be run independently:

| File | Demonstrates |
| --- | --- |
| `hello.vol` | Values, strings, printing, and conversion |
| `conditions.vol` | Comparisons and Boolean conditions |
| `loops.vol` | `while`, `repeat`, and assignment |
| `arrays.vol` | Arrays, indexing, length, and `.each` |
| `collections.vol` | `.where`, `.sum()`, and `assert` |
| `functions.vol` | Function declarations, calls, parameters, and returns |
| `scope.vol` | Lexical scopes and shadowing |
| `interaction.vol` | Interactive input and assertions |
| `arguments.vol` | Command-line arguments |

For example:

```text
.\vol.exe run .\examples\functions.vol
.\vol.exe run .\examples\arguments.vol -- apple banana
```

## Build the CLI

```text
go build ./cmd/vol
./vol ./examples/first.vol
```

On Windows:

```text
go build ./cmd/vol
.\vol.exe run .\examples\first.vol
```

The command-line interface uses `vol run <file.vol>`. Arguments following the
file are available to the program in the built-in `args` array. An optional `--`
separates VOL options from program arguments:

```text
.\vol.exe run .\program.vol -- first second
```

The original shorthand
`vol <file.vol>` is also accepted for compatibility.

## Run the Tests

```text
go test ./...
```

For the same deeper checks used in continuous integration:

```text
gofmt -l .
go vet ./...
go test -race -coverprofile=coverage.out ./...
```

The test suite covers successful behavior, stable diagnostic codes, source locations,
invalid types and arities, collection and numeric boundaries, lexical scoping,
short-circuiting, built-in I/O failures, CLI exit behavior, and every checked-in
example program.
