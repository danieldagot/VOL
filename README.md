# VOL

> **Vector-Oriented Language**  
> Also known as the **Vibe-Oriented Language**.

VOL is an experimental systems and backend programming language built around intent, semantic density, safety, and native performance.

The repository currently contains the first tree-walking interpreter prototype. The syntax is provisional and documented in [`SYNTAX.md`](SYNTAX.md).

## Run the Example

Requirements:

- Go 1.24 or newer

Run directly:

```text
go run ./cmd/vol ./examples/first.vol
```

Expected output:

```text
16
```

Build the CLI:

```text
go build ./cmd/vol
./vol ./examples/first.vol
```

Run the tests:

```text
go test ./...
```

## Implemented Syntax

- integer, floating-point, Boolean, and string literals
- inferred variable declarations with `:=`
- assignment with `=`
- arithmetic, comparison, and Boolean expressions
- arrays, indexing, indexed assignment, and `.length`
- brace-delimited lexical scopes
- `if` and `else`
- `repeat` and provisional `while` loops
- array iteration with `.each`
- `print`
- line comments beginning with `//`
- source-located parser and runtime diagnostics

Functions, static type checking, ownership analysis, the formatter, the language server, and native backends are not implemented yet.

## Layout

```text
cmd/vol          command-line entry point
internal/lang    lexer, parser, AST, diagnostics, and interpreter
examples         example VOL programs
```
