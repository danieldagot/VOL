# VOL Examples

Executable programs for the Supported SF-1 surface. Run any entry file with:

```text
go run ./cmd/vol run ./examples/<path>.vol
```

## Layout

| Folder | Purpose |
| --- | --- |
| [`basics/`](basics/) | Short language tour (values, control flow, arrays, functions, I/O) |
| [`features/`](features/) | One-topic demos for Option, Result, structs, anonymous `fn`, modules |
| [`projects/`](projects/) | Small multi-file “real enough” apps that combine features |

## Basics

| File | Demonstrates |
| --- | --- |
| `basics/hello.vol` | Values, strings, `print`, `string` |
| `basics/first.vol` | `.each` and a running total |
| `basics/conditions.vol` | `if` / `and` |
| `basics/loops.vol` | `while`, `repeat` |
| `basics/arrays.vol` | Indexing, `.len`, `.copy` / `.deep_copy` |
| `basics/collections.vol` | `.where`, `.sum()`, `assert` |
| `basics/functions.vol` | Named `fn`, parameters, `return` |
| `basics/scope.vol` | Lexical scopes and shadowing |
| `basics/arguments.vol` | Built-in `args` |
| `basics/interaction.vol` | `input` and `assert` |

```text
go run ./cmd/vol run ./examples/basics/first.vol
go run ./cmd/vol run ./examples/basics/arguments.vol -- apple banana
```

## Features

| File | Demonstrates |
| --- | --- |
| `features/anonymous.vol` | Anonymous `fn` values |
| `features/option.vol` | `some` / `none`, if-let, `??` |
| `features/result.vol` | `ok` / `err`, if-let |
| `features/result_helpers.vol` | Returning Result from helpers |
| `features/option_result.vol` | Option vs Result side by side |
| `features/struct.vol` | Product structs, named + positional literals |
| `features/struct_nested.vol` | Nested structs and shared identity |
| `features/const_struct.vol` | Shallow `const` on struct bindings |
| `features/modules/main.vol` | `import` + `export` |
| `features/modules/import_struct.vol` | Importing a struct type |

## Projects

| Project | Entry | Story |
| --- | --- | --- |
| [`projects/gradebook/`](projects/gradebook/) | `main.vol` | Class roster: structs, `.map` / `.count` / `.where` |
| [`projects/contacts/`](projects/contacts/) | `main.vol` | Contact lookup with Option if-let / `??` |
| [`projects/shop/`](projects/shop/) | `main.vol` | Cart checkout with Result and postfix `?` |
| [`projects/fibonacci/`](projects/fibonacci/) | `main.vol` | Sequence with multi-assign |

```text
go run ./cmd/vol run ./examples/projects/gradebook/main.vol
go run ./cmd/vol run ./examples/projects/shop/main.vol
```

These are covered by `TestExamplesRemainExecutable` in `internal/lang/examples_test.go`.
