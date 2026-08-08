# VOL

> **Vocabulary Optimized Language**

VOL is an experimental systems and backend programming language designed for humans and Large Language Models.

The project explores a simple idea:

> Write intent. Compile excellence.

Most programming languages were designed before code-generating Large Language Models existed. Their syntax, tooling, diagnostics, and abstractions were optimized primarily for human-written software and historical conventions.

VOL asks what a systems programming language should look like when LLMs are first-class programmers.

Its primary goal is to maximize **semantic density**: express more intent with fewer tokens while preserving deterministic structure, readability, safety, and native performance. The language should be efficient for an LLM to generate, understand, review, modify, and debug.

That goal is a research direction, not a measured property of the current
prototype. SF-1 freezes Supported syntax; it does **not** prove LLM optimization.
Hand-written source-density benches are not generate/repair evidence. Token
counts also depend on each model's tokenizer, so “fewer tokens” alone is not a
success metric. The intended long-term measure is task success relative to total
tokens consumed across generation, diagnostics, and repair — see
[`LLM_BENCHMARK.md`](LLM_BENCHMARK.md) and [`IDEAS.md`](IDEAS.md).

## Project Status

**VOL is currently a very early tree-walking interpreter prototype.**
There is no native compiler, no static type system, no ownership or borrow
checker, no broad standard library (ambient tiny core only), and no backend.
Ownership, vectorization, and parallelization are not language semantics today.
What exists is a minimal interpreter that can execute small programs written in
a provisional syntax.

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

**Surface Freeze SF-1:** vision-aligned Supported syntax in [`SPEC.md`](SPEC.md)
— Option/Result (if-let, `??`, postfix `?`; `match` removed), modules, product
structs, anonymous and expression-body `fn`, multi-assign, `.map`/`.count`, and
the rest of the Implemented core. Further vocabulary (`|>`, enums, dual-return
sugar, …) waits for an SF-2 bump. Harness card
[`bench/llm/cards/vol_v1.md`](bench/llm/cards/vol_v1.md) is the `core_v2` **task
card** (SF-1-bound subset — not a full SF-1 product tour). SF-0 /
[`vol_v0`](bench/llm/cards/vol_v0.md) remains for historical tables.

- integer, floating-point, Boolean, and string literals
- inferred variable declarations with `:=` (including multi-declare `a, b := …`)
- variable assignment with `=` (including multi-assign `a, b = …`; RHS first)
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
- collection filtering with `.where(...)`, mapping with `.map(...)`, counting with `.count(...)`, and numeric aggregation with `.sum()` (prefer side-effect-free predicates; purity checks Planned)
- Option values with `some(value)` / `none`; unwrap with if-let and `??` (`match` rejected)
- Result values with `ok(value)` / `err(value)`; unwrap with if-let and postfix `?` in functions (bugs still trap)
- product `struct` types, named and positional construction, and field get/set
- `import "path"` / project-local aliases via nearest `vol.config.json`; live `export` across modules (see `examples/features/modules/aliases/` for `@lib`; future `@std` spelling is vision-only in [`IDEAS.md`](IDEAS.md), not a shipped stdlib)
- `print`
- interactive text input with `input()` or `input(prompt)`
- runtime checks with `assert(condition)` or `assert(condition, message)`
- value-to-string conversion with `string(value)`
- command-line arguments through the built-in `args` array
- named functions with `fn name(...)`, anonymous `fn(...)` expressions (block or expression body), parameters, calls, and `return` (missing return is `nothing`; using it as a value is `R029`)
- module export lists that may appear before or after definitions
- line comments beginning with `//`
- source-located parser and runtime diagnostics with human and JSON output (`vol --json run <file.vol>`)
- table-driven lexer, parser, resolver, interpreter, diagnostic, CLI, and executable-example tests
- fuzz seed coverage ensuring arbitrary parser input does not panic
- continuous integration with formatting, vet, race-detection, and coverage checks

### What Is Missing

- static type checking
- richer diagnostic suggestions across more error codes
- canonical VOL formatter rewriter (`vol fmt` CLI stub parses only today)
- a published compatibility policy and versioned conformance corpus (SF-0 / SF-1 pins)
- broader LLM workflow results across more models and realistic backend tasks
- initial language-server support

### What Is Still an Open Experiment

- typed function signatures
- explicit type syntax
- ownership and borrowing (direction: local escape first, API contracts later —
  see [`IDEAS.md`](IDEAS.md); not implemented)
- dual-return sugar over Result; Result-returning I/O APIs
  (`input`/`assert` still trap — see [`IDEAS.md`](IDEAS.md))
- enums / tagged unions / methods on structs (product `struct` is Supported)
- symbol-selection imports; multi-file folder packages beyond `path.vol` /
  `path/mod.vol`
- generics
- richer pattern matching / tagged-union match (binary unwrap uses if-let / `??` / `?`)
- asynchronous I/O syntax
- allocator and memory-layout constraints (unspecified; do not claim inference)
- compile-time evaluation
- concurrency / `parallel` (Planned; no guarantees until scheduling is specified)
- overflow wrapping ops / build modes (default is already trap; modes + wrap ops
  Planned — see `SPEC.md` §4.3 and [`IDEAS.md`](IDEAS.md))
- bounds-checking behavior across build modes (same family as overflow modes)
- `|>` pipeline sugar / `=>` arrows (anonymous `fn` is Supported; see [`IDEAS.md`](IDEAS.md))
- optional `?T` sugar for Option (canonical form is `some`/`none`)
- C and LLVM backend details

Near-term priority is to keep SF-1 precise and finish foundations (`vol fmt`
rewriter; richer std **behind imports** when designed — not ambient growth)
before an SF-2 feature bump. Remaining directions live in [`IDEAS.md`](IDEAS.md).

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

### Option if-let and `??`

```vol
found := some("VOL")
if some name := found {
    print "Hello, " + name
} else {
    print "missing"
}
print found ?? "missing"
```

### Anonymous functions

```vol
double := fn(x) {
    return x * 2
}
print double(21)

triple := fn(x) x * 3
print triple(7)
```

### Result if-let and `?`

```vol
r := ok(7)
if ok n := r {
    print n
} else err msg {
    print msg
}

fn divide(a, b) {
    if b == 0 {
        return err("zero")
    }
    return ok(a / b)
}
n := divide(10, 2)?
```

### Multi-assign and collections

```vol
a, b := 0, 1
a, b = b, a + b
nums := [1, 2, 3, 4]
print nums.map(_ * 2)
print nums.count(_ > 2)
```

### Structs

```vol
struct User {
    name
    age
}
u := User { name: "Ada", age: 36 }
v := User { "Ada", 36 }
print u.name
```

### Imports

```vol
import "examples/features/modules/math"
print add(21, 21)
```

Executable samples live under [`examples/`](examples/) — see
[`examples/README.md`](examples/README.md) for the full index (`basics/`,
`features/`, and multi-file `projects/` such as gradebook, contacts, and shop).

## Repository Layout

```text
cmd/vol          command-line entry point
internal/lang    lexer, parser, AST, diagnostics, and interpreter
examples         example VOL programs (basics, features, projects)
bench            source token density benchmark (VOL vs Python/Go/Rust/Zig)
```

## Source Token Density Benchmark

On 16 equivalent small programs, VOL currently uses about:

- **~13% fewer tokens than Python** (all-suite median)
- **~17% fewer than Python** on the **compression** tier (filter/map/count/sum,
  bare numeric output — best read for semantic density)
- **~47% fewer tokens than Go**
- **~39% fewer tokens than Rust**
- **~62% fewer tokens than Zig**

(medians; nearly the same under both `cl100k_base` and `o200k_base`. Reports also
split **labeled** vs **parity** tiers — see [`bench/README.md`](bench/README.md).)

That is source size only — hand-written programs that print the same output. It
does **not** measure whether an LLM generates correct VOL more easily, or how
many tokens generate/repair workflows consume. Do not cite these ratios as
workflow proof or as evidence that SF-1 is “LLM optimized.”

Full per-task numbers: [`bench/results/density.md`](bench/results/density.md).
How to run / regenerate: [`bench/README.md`](bench/README.md).
Working preset to push density further (dynamics-first): [`TOKEN_EFFICIENCY.md`](TOKEN_EFFICIENCY.md).

```text
cd bench && uv sync && uv run python harness/count_tokens.py
```

## LLM Workflow Benchmark

The default workflow baseline is **Python** (interpreted peer for the current
prototype). Go remains an optional compiled baseline (`--langs vol,go`).

One protocol-v1.1 (`core_v2`) run is published for `gemini-3.5-flash-lite`,
temperature 0, three replicates, and at most two repair rounds — **VOL
(`vol_v1` / SF-1-bound `core_v2` task card) vs Python**. The published JSONL is
self-contained (VOL live + frozen Python rows from the prior same-model run).
The suite has five tasks covering generation,
diagnostic-seeded repair, and modification. Summaries split prompt vs completion
and report cold totals (card re-sent every request) plus warm totals (estimated
language-card cost amortized away).

| Language | First-try | Success @ K | Mean cold | Mean warm | Mean completion |
| --- | ---: | ---: | ---: | ---: | ---: |
| Python | 100% | 100% | 762.9 | 426.9 | 135.1 |
| VOL | 100% | 100% | 848.6 | 412.6 | 109.9 |

VOL matched Python on **first-try** and **success @ K**, used about **11.2% more
cold** workflow tokens, and about **3.3% fewer warm** tokens once the card is
amortized. Completions were about 18.7% smaller; prompts about 17.7% larger.
Card size is ~436 tokens (`vol_v1`) vs Python ~336.

This small synthetic run still does **not** prove real-world LLM superiority or
runtime performance. A second model on `core_v2` is needed before treating
results as stable. Do not mix these workflow numbers with the hand-written
density table above.

Full result: [`bench/llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-041440.md`](bench/llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-041440.md).
Protocol and limitations: [`LLM_BENCHMARK.md`](LLM_BENCHMARK.md).

## Requirements

- Go 1.24 or newer

## Run the Example

```text
go run ./cmd/vol run ./examples/basics/first.vol
```

Expected output:

```text
16
```

## Example Programs

Organized under [`examples/`](examples/) — full index in
[`examples/README.md`](examples/README.md):

| Path | Demonstrates |
| --- | --- |
| `basics/` | Language tour: values, control flow, arrays, functions, I/O |
| `features/` | Option, Result, structs, anonymous `fn`, modules |
| `projects/gradebook/` | Class roster (structs + `.map` / `.count` / `.where`) |
| `projects/contacts/` | Contact lookup (Option if-let / `??`) |
| `projects/shop/` | Cart checkout (Result + postfix `?`) |
| `projects/fibonacci/` | Multi-assign sequence |

For example:

```text
go run ./cmd/vol run ./examples/basics/functions.vol
go run ./cmd/vol run ./examples/basics/arguments.vol -- apple banana
go run ./cmd/vol run ./examples/projects/gradebook/main.vol
```

## Build the CLI

```text
go build ./cmd/vol
./vol ./examples/basics/first.vol
```

On Windows:

```text
go build ./cmd/vol
.\vol.exe run .\examples\basics\first.vol
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
