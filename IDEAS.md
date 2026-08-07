# VOL Ideas and Future Work

This document collects ideas that may shape VOL in the future.

Items here are exploratory, not promises. Accepted language behavior belongs in `SYNTAX.md` and, eventually, the formal specification.

## Near-Term Foundation

- [x] Define tokens and lexical rules.
- [x] Implement the lexer.
- [x] Implement the parser.
- [x] Define the abstract syntax tree (AST).
- [x] Track precise source locations on every syntax node.
- [x] Create deterministic human-readable diagnostics.
- [ ] Create machine-readable JSON diagnostics.
- [x] Implement lexical scopes.
- [x] Implement an interpreter.
- [x] Add table-driven parser, semantic, runtime, diagnostic, CLI, and example conformance tests.
- [x] Add parser fuzz seeds and CI checks for formatting, vet, races, and coverage.
- [ ] Add versioned diagnostic snapshots when the diagnostic presentation is stable.

## Initial Language Features

- [x] Integer, floating-point, Boolean, and string literals.
- [x] Variables with inferred types.
- [ ] Explicit mutability rules.
- [x] Arithmetic, comparison, and Boolean operators.
- [x] Blocks using braces.
- [x] `if` and `else` statements.
- [x] `repeat` loops.
- [x] Conditional loops (provisional `while`).
- [x] Arrays, indexing, and bounds checking.
- [x] Collection iteration with `.each`.
- [x] Functions, parameters, calls, and return values.
- [x] Comments.
- [x] Basic built-ins for printing, input, length, assertions, conversion, and arguments.

## Formatter

- [ ] Build a syntax-aware formatter using the compiler parser.
- [ ] Define one canonical style.
- [ ] Make formatting deterministic and idempotent.
- [ ] Preserve comments reliably.
- [ ] Support formatting a file or directory.
- [ ] Support a check-only mode for continuous integration.
- [ ] Make formatting fast enough to run on save.

Planned commands:

```text
vol fmt file.vol
vol fmt .
vol fmt --check .
```

Formatting must only change presentation, never program structure. A separate
future simplifier may recognize explicit algorithms and suggest equivalent VOL
semantic operations, while preserving both source styles as valid language forms.

Possible workflow:

```text
vol simplify --check file.vol
vol simplify --diff file.vol
vol simplify --write file.vol
```

For example, it may suggest replacing an accumulator loop containing a filter
with `numbers.where(_ > 5).sum`. Rewrites must be optional, deterministic, and
proven safe with respect to mutation, side effects, ordering, and error behavior.

## Projects and Modules

- [ ] Discover the nearest `vol.config.json` by walking from the source file toward the filesystem root.
- [ ] Resolve imports relative to the configured project root.
- [ ] Treat folders as module namespaces with one canonical resolution rule.
- [ ] Resolve path aliases such as `@db/*` without source-relative traversal.
- [ ] Reject aliases and imports that escape the project root unless explicitly allowed.
- [ ] Detect import cycles and report their complete path deterministically.
- [ ] Collect exports written anywhere in a module and format one sorted export list at the top.

Proposed configuration:

```json
{
  "name": "vol",
  "root": ".",
  "paths": {
    "@compiler": "compiler",
    "@std": "std"
  }
}
```

## Language Server

The language server must share the compiler's lexer, parser, syntax tree, semantic model, diagnostics, and source-location system.

Initial features:

- [ ] Live diagnostics.
- [ ] Document formatting and format on save.
- [ ] Hover information.
- [ ] Go to definition.
- [ ] Find references.
- [ ] Rename symbol.
- [ ] Completion.
- [ ] Document symbols.
- [ ] Semantic highlighting.

Future VOL-specific features:

- [ ] Display inferred types.
- [ ] Display ownership and lifetime decisions.
- [ ] Show stack and heap allocation decisions.
- [ ] Warn about expensive or unbounded allocation.
- [ ] Explain compiler optimizations.
- [ ] Preview generated C or LLVM IR.
- [ ] Provide deterministic automatic fixes.

## Compiler Pipeline

```text
source
  -> lexer
  -> tokens
  -> parser
  -> AST
  -> semantic analysis
  -> HIR
  -> MIR
  -> optimization
  -> backend
```

Planned execution stages:

1. Tree-walking interpreter.
2. C backend.
3. LLVM backend.
4. Self-hosting compiler.

## Safety and Memory

- [ ] Infer ownership, borrowing, and lifetimes.
- [ ] Prefer stack allocation when values do not escape.
- [ ] Avoid a mandatory garbage collector.
- [ ] Avoid hidden unbounded allocation.
- [ ] Provide explicit control when inference is insufficient.
- [ ] Make inferred allocation and ownership inspectable.
- [ ] Define selectable safety and optimization build modes.

Guiding rule:

> Infer by default. Expose control when necessary.

## Optimization

- [ ] Compile-time specialization.
- [ ] Whole-program dead-code elimination.
- [ ] Escape analysis.
- [ ] Automatic vectorization.
- [ ] Safe automatic parallelization.
- [ ] Data-oriented memory-layout optimization.
- [ ] Profile-guided optimization.
- [ ] Incremental compilation.

## Batteries-Included Library

Build broad capabilities on small, orthogonal foundations. Unused functionality must not increase executable size or runtime cost.

Priority areas:

- [ ] Filesystems and processes.
- [ ] Networking, HTTP, and TLS.
- [ ] Serialization and structured data.
- [ ] Databases and storage.
- [ ] Synchronous and asynchronous I/O.
- [ ] Concurrency and parallelism.
- [ ] Cryptography, hashing, and compression.
- [ ] Logging, tracing, and metrics.
- [ ] Configuration and command-line parsing.
- [ ] Testing, benchmarking, and profiling.
- [ ] C interoperability.
- [ ] Cross-compilation.

## Compiler Metrics

Every build may report:

- compile time
- binary size
- estimated memory use
- optimization level
- source token count
- semantic density score
- Meaning per Token (MPT)

The definitions of semantic density and MPT must be objective, reproducible, and difficult to game.

## Open Design Questions

- What is VOL's exact mutability model?
- Is `if` an expression, a statement, or both?
- What syntax should replace or represent a conventional `while` loop?
- Are functions declared with a keyword or inferred from their form?
- How are fallible operations and errors represented?
- How are nullable or optional values represented?
- How are structs, enums, and tagged unions declared?
- How are modules and packages organized?
- Which standard features require explicit imports?
- How does a programmer constrain inferred allocation?
- What guarantees does automatic parallelization provide?
- Which build modes control bounds and overflow checking?
