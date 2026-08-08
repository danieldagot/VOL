# VOL

> **Vocabulary Optimized Language**

VOL is an experimental systems and backend programming language designed for humans and Large Language Models.

The project explores a simple idea:

> Write intent. Compile excellence.

Most programming languages were designed before code-generating Large Language Models existed. Their syntax, tooling, diagnostics, and abstractions were optimized primarily for human-written software and historical conventions.

VOL asks what a systems programming language should look like when LLMs are first-class programmers.

Its primary goal is to maximize **semantic density**: express more intent with fewer tokens while preserving deterministic structure, readability, safety, and native performance. The language should be efficient for an LLM to generate, understand, review, modify, and debug.

That goal is a research direction, not a measured property of the current
prototype. SF-3.1 freezes Supported syntax plus `@std`; it does **not** prove LLM optimization.
Hand-written source-density benches are not generate/repair evidence. Token
counts also depend on each model's tokenizer, so “fewer tokens” alone is not a
success metric. The intended long-term measure is task success relative to total
tokens consumed across generation, diagnostics, and repair — see
[`LLM_BENCHMARK.md`](LLM_BENCHMARK.md) and [`IDEAS.md`](IDEAS.md).

## Project Status

**VOL is currently a very early tree-walking interpreter prototype.**
There is no native compiler, no static type system, no ownership or borrow
checker, and no native backend. **Surface Freeze SF-3.1** (card
`vol_v3_1`) hardens the SF-3 `@std` + dict surface: namespaced imports,
`dict {…}` literals, multiline chains, and `.len`-only length. Assignment /
sharing today and future ownership intent: [`MEMORY_MODEL.md`](MEMORY_MODEL.md).
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

**Surface Freeze SF-3.1** (active foundation): Supported surface in
[`SPEC.md`](SPEC.md) — SF-3 `@std` (math, strings, fs, path, env, time, url,
json, yaml, http, process, db/SQLite) plus namespaced imports
(`import "@std/json"` → `json.parse(...)`), dict literals `dict { key: value }`
(ambient `dict()` / `dict("k", v, …)` kept), multiline expression continuation,
and `.len` as the only length form. Default harness card:
[`bench/llm/cards/vol_v3_1.md`](bench/llm/cards/vol_v3_1.md). Historical
[`vol_v3`](bench/llm/cards/vol_v3.md) / [`vol_v2`](bench/llm/cards/vol_v2.md) /
[`vol_v1`](bench/llm/cards/vol_v1.md) / [`vol_v0`](bench/llm/cards/vol_v0.md)
remain for published earlier tables. Language sugar (`|>`, enums, dual-return,
…) is **Planned (unscheduled)** in [`IDEAS.md`](IDEAS.md) — not a promised next
freeze.

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
- collection filtering with `.where(...)`, mapping with `.map(...)`, counting with `.count(pred)`, length via `.len` only (zero-arg `.count()` rejected), and numeric aggregation with `.sum()` (prefer side-effect-free predicates; purity checks Planned)
- Option values with `some(value)` / `none`; unwrap with if-let and `??` (`match` rejected) — see `examples/features/option.vol`
- Result error handling with `ok(value)` / `err(value)`, if-let, and postfix `?` in functions (bugs still trap) — see `examples/features/errors.vol` and `examples/projects/shop/`
- product `struct` types, named and positional construction, and field get/set
- namespaced `import "path"` / `import "@std/…"` (binds module basename; not remappable for `@std`) plus project-local aliases via nearest `vol.config.json`; live `export` across modules (see `examples/features/modules/aliases/` for `@lib`)
- dict values via `dict { key: value }`, ambient `dict()` / `dict("k", v, …)`, `d["k"]`, `.keys()` — see `examples/features/std/` and `examples/projects/hits/`
- multiline expression continuation (postfix `.`, calls, commas, binary ops)
- `print` with one or more values (space-joined display forms)
- string `+` coercion (`"n=" + 7`); keep `string(value)` for explicit convert
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
- a published compatibility policy and versioned conformance corpus (SF-0–SF-3.1 pins)
- broader LLM workflow results across more models and realistic backend tasks
- initial language-server support

### What Is Still an Open Experiment

- typed function signatures
- explicit type syntax
- ownership and borrowing (direction: local escape first, API contracts later —
  see [`IDEAS.md`](IDEAS.md); not implemented)
- dual-return sugar over Result; ambient `input`/`assert` still trap
  (`@std` I/O already returns Result — see [`IDEAS.md`](IDEAS.md))
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

Near-term priority is **SF-3.1 foundation** — keep namespaced `@std`, dict
literals, multiline, and `.len`-only precise (tests, diagnostics, docs).
`vol fmt` rewriter remains parallel foundation work. Language sugar (`|>`,
enums, dual-return, …) and Postgres/MySQL / ORM / WebSockets are **Planned
(unscheduled)** — not the next freeze queue.

See [`SPEC.md`](SPEC.md) for vocabulary, syntax, and formal behavior of the
current interpreter, [`MEMORY_MODEL.md`](MEMORY_MODEL.md) for assignment /
sharing, and [`IDEAS.md`](IDEAS.md) for future work and open design questions.

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

On 21 equivalent small programs (SF-2 collection/report tasks plus SF-3 `@std`
stdlib tasks; `cl100k_base` medians):

- **~17% fewer tokens than Python** (all-suite median; ratio **0.829**)
- **~24% fewer than Python** on the **compression** tier (filter/map/count/sum; **0.764**)
- **~35% fewer than Python** on the **labeled** tier (multi-arg `print` / coercion; **0.651**)
- **~2% fewer than Python** on the **stdlib** tier (`17`–`21`; **0.976**) under
  SF-3.1 namespaced `@std` + `dict {…}` (namespaces trade density for clarity;
  historical SF-3 flat-import stdlib median was **0.929**)
- **~51% fewer tokens than Go** / **~48% fewer than Rust** / **~66% fewer than Zig**
  (all-suite medians)

(Nearly the same under `o200k_base`. Prefer **compression** for collection-intent
claims and **stdlib** for `@std` claims — see [`bench/README.md`](bench/README.md).)

That is source size only — hand-written programs that print the same output. It
does **not** measure whether an LLM generates correct VOL more easily, or how
many tokens generate/repair workflows consume. Do not cite these ratios as
workflow proof or as evidence that SF-3 is “LLM optimized.”

Full per-task numbers: [`bench/results/density.md`](bench/results/density.md).
How to run / regenerate: [`bench/README.md`](bench/README.md).
Working preset to push density further (dynamics-first): [`TOKEN_EFFICIENCY.md`](TOKEN_EFFICIENCY.md).

```text
cd bench && uv sync && uv run python harness/count_tokens.py
```

## LLM Workflow Benchmark

The default workflow baseline is **Python** (interpreted peer for the current
prototype). Go remains an optional compiled baseline (`--langs vol,go`).

Primary language-use table: **`intent_v1`** (7 tasks, includes SF-3 `@std`
strings/json) with **`vol_v3` / SF-3** vs Python (`gemini-3.5-flash-lite`,
temperature 0, 3 replicates, K=2). Card ~350 tokens (`vol_v3`) vs Python ~336.

| Language | First-try | Success @ K | Mean cold | Mean warm | Mean completion |
| --- | ---: | ---: | ---: | ---: | ---: |
| Python | 100% | 100% | 733.6 | 397.6 | 113.1 |
| VOL (SF-3) | 100% | 100% | 715.7 | 365.7 | 78.4 |

VOL matched Python on **first-try** / **success @ K**, used about **2.4% fewer
cold** and **8.0% fewer warm** workflow tokens, with ~31% smaller completions.
One-model result — do not treat as stable across models. Full summary:
[`bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-063341.md`](bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-063341.md).

Historical notes (freeze ids unchanged): pre-stdlib 5-task SF-3 `intent_v1`
([`…061355.md`](bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-061355.md));
SF-2 `intent_v1` / `vol_v2`
([`…051437.md`](bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-051437.md));
SF-1 `core_v2` (~+11% cold / ~−3% warm
[`…041440.md`](bench/llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-041440.md));
early SF-1 `intent_v1` 80% first-try (`.count()` arity)
([`…045310.md`](bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-045310.md)).

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
| `projects/hits/` | SQLite hit counter (`@std/db` + `@std/env`) |
| `projects/api/` | Mocked notes HTTP API (`@std/http` + `@std/db` + `@std/json`) |

For example:

```text
go run ./cmd/vol run ./examples/basics/functions.vol
go run ./cmd/vol run ./examples/basics/arguments.vol -- apple banana
go run ./cmd/vol run ./examples/projects/gradebook/main.vol
go run ./cmd/vol run ./examples/projects/api/main.vol
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
