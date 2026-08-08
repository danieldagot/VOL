# VOL Examples

Executable programs for the Supported SF-3 surface. Run any entry file with:

```text
go run ./cmd/vol run ./examples/<path>.vol
```

## Layout

| Folder | Purpose |
| --- | --- |
| [`basics/`](basics/) | Short language tour (values, control flow, arrays, functions, I/O) |
| [`features/`](features/) | One-topic demos for Option, Result, structs, modules, `@std` |
| [`projects/`](projects/) | Small multi-file “real enough” apps that combine features |

## Basics

| File | Demonstrates |
| --- | --- |
| `basics/hello.vol` | Values, strings, multi-arg `print`, string `+` |
| `basics/first.vol` | `.each` and a running total |
| `basics/conditions.vol` | `if` / `and` |
| `basics/loops.vol` | `while`, `repeat` |
| `basics/arrays.vol` | Indexing, `.len`, `.copy` / `.deep_copy` |
| `basics/collections.vol` | `.where`, `.count`, `.sum()`, multi-arg `print`, `assert` |
| `basics/functions.vol` | Named `fn`, parameters, `return` |
| `basics/scope.vol` | Lexical scopes and shadowing |
| `basics/arguments.vol` | Built-in `args` |
| `basics/interaction.vol` | `input`, multi-arg `print`, `assert` |

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
| `features/errors.vol` | **Error handling tour**: if-let, postfix `?`, labeled failures |
| `features/option_result.vol` | Option vs Result side by side |
| `features/struct.vol` | Product structs, named + positional literals |
| `features/struct_nested.vol` | Nested structs and shared identity |
| `features/const_struct.vol` | Shallow `const` on struct bindings |
| `features/print_multi.vol` | Multi-arg `print`, string `+` coercion, `.count()` |
| `features/modules/main.vol` | `import` + `export` |
| `features/modules/import_struct.vol` | Importing a struct type |
| `features/modules/aliases/main.vol` | `@lib` path alias via local `vol.config.json` |
| `features/std/dict.vol` | `dict("k", v, …)`, string-key index, `.keys()` |
| `features/std/math_strings.vol` | `@std/math` + `@std/strings` |
| `features/std/fs_path.vol` | `@std/fs` + `@std/path` |
| `features/std/env_time.vol` | `@std/env` + `@std/time` |
| `features/std/url.vol` | `@std/url` parse |
| `features/std/json.vol` | `@std/json` parse/dump |
| `features/std/yaml.vol` | `@std/yaml` parse |
| `features/std/process.vol` | `@std/process` run |
| `features/std/db.vol` | `@std/db` SQLite |
| `features/std/http_fetch.vol` | `@std/http` reply helper |

```text
go run ./cmd/vol run ./examples/features/errors.vol
go run ./cmd/vol run ./examples/features/std/math_strings.vol
```

Note: `@std/strings` and `@std/path` both export `join`; `@std/url` /
`@std/json` / `@std/yaml` all export `parse`. Import one module per colliding
name in a given file.

## Projects

| Project | Entry | Story |
| --- | --- | --- |
| [`projects/gradebook/`](projects/gradebook/) | `main.vol` | Class roster: structs, `.map` / `.count` / `.where` |
| [`projects/contacts/`](projects/contacts/) | `main.vol` | Contact lookup with Option if-let / `??` |
| [`projects/shop/`](projects/shop/) | `main.vol` | Cart checkout with Result and postfix `?` |
| [`projects/fibonacci/`](projects/fibonacci/) | `main.vol` | Sequence with multi-assign |
| [`projects/hits/`](projects/hits/) | `main.vol` | SQLite hit counter via `@std/db` + `@std/env` |
| [`projects/api/`](projects/api/) | `main.vol` | Mocked notes HTTP API (`@std/http` + `@std/db` + `@std/json`); see folder README for `server.vol` |

```text
go run ./cmd/vol run ./examples/projects/gradebook/main.vol
go run ./cmd/vol run ./examples/projects/hits/main.vol
go run ./cmd/vol run ./examples/projects/api/main.vol
```

These are covered by `TestExamplesRemainExecutable` in `internal/lang/examples_test.go`.
