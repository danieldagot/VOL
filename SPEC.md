# VOL Language Specification (Prototype v0 / SF-2)

> Status: **Surface Freeze SF-2** (density dynamics; harness card `vol_v2`)
> Audience: humans and LLMs

> Source of truth for behavior: this file plus the tests in `internal/lang`

This is VOL's language specification. It describes **only what the tree-walking
interpreter does today**. If something is not written here, it is not part of
the language yet—even if it appears in vision docs.

You do not need prior language-design experience to maintain this file. When you
change the lexer, parser, resolver, or interpreter, update the matching section
here and add a test.

Related docs:

| Document | Role |
| --- | --- |
| [`README.md`](README.md) | Project status and examples |
| [`IDEAS.md`](IDEAS.md) | Planned features and open questions |
| [`AGENTS.md`](AGENTS.md) | Project vision and contribution rules |
| [`bench/llm/cards/vol_v2.md`](bench/llm/cards/vol_v2.md) | Default harness card (SF-2; density dynamics) |
| [`bench/llm/cards/vol_v1.md`](bench/llm/cards/vol_v1.md) | Historical SF-1 / `core_v2` task card |
| [`bench/llm/cards/vol_v0.md`](bench/llm/cards/vol_v0.md) | Historical `core_v2` card (SF-0) |

This file is the single source for implemented syntax, vocabulary, and semantics.
Planned words such as `parallel` belong only in [`IDEAS.md`](IDEAS.md).

---

## 0. How to read this spec

A programming-language specification answers four questions:

1. **Lexing** — which characters form tokens?
2. **Parsing** — which token sequences are programs?
3. **Static checks** — which programs are rejected before running?
4. **Runtime semantics** — what happens when a valid program runs?

VOL currently has:

- a lexer
- a parser
- a name/arity resolver
- an interpreter

It does **not** yet have a static type system, ownership rules, or a native
backend. Those belong in [`IDEAS.md`](IDEAS.md) until specified and implemented.

Notation:

- `a | b` means either `a` or `b`
- `a?` means optional
- `a*` means zero or more
- `a+` means one or more
- `"text"` means the literal characters `text`
- Names in *italics* are grammar nonterminals

**Supported** means implemented and covered by tests. **Provisional** means
implemented and tested, but spelling or meaning may change.

### Surface freeze SF-2

**SF-2** freezes the Supported / Provisional surface in this document and §11.
It includes the SF-1 vision-aligned surface plus **token-density dynamics**:
multi-arg `print`, string-context `+` coercion, and zero-arg `.count()` as
length. The default harness card is
[`bench/llm/cards/vol_v2.md`](bench/llm/cards/vol_v2.md). Historical published
`core_v2` / `intent_v1` tables that cite SF-1 / `vol_v1` stay tied to that card.

**Product freezes** (use these in status docs and LLM result tables; do not mix
freeze IDs in one table):

| Freeze | Card | Meaning |
| --- | --- | --- |
| SF-0 | `vol_v0.md` | First pin (historical `core_v2` results) |
| SF-1 | `vol_v1.md` | Vision-aligned surface before density dynamics |
| SF-2 | `vol_v2.md` | Density dynamics (multi-arg `print`, string `+` coercion, `.count()`) |
| SF-3+ | (when shipped) | Next expansion (`|>`, enums, dual-return, …) |

Only keep cards for product freezes that have (or will have) published harness
runs. Intermediate implementation drafts are not cards.

| Allowed under SF-2 | Requires bumping to SF-3 (or later) |
| --- | --- |
| Bug fixes that restore documented behavior | New keywords or operators beyond SF-2 |
| Clearer diagnostics / `fix` text for existing codes | `|>` pipelines, `=>`, dual-return sugar, lazy views |
| Tests, examples, and doc sync for existing forms | Enums/tagged unions, ownership, parallel |
| Card wording edits that do **not** add features | Any form that expands the Supported vocabulary further |

During SF-2, do **not** implement further Planned syntax from [`IDEAS.md`](IDEAS.md)
without a freeze bump, SPEC updates, tests, and a new card version.

### Syntax principles

- Use braces for every multi-line scope.
- Use newlines to separate statements; do not require semicolons.
- Do not require parentheses around conditions.
- Prefer type inference when enough information is available.
- Prefer one canonical representation for each distinct intent.
- Remove boilerplate without removing structural clarity.
- Optimize syntax for humans, language models, formatters, and static analysis.

`.each` (imperative iteration) and `.where` / `.map` / `.count` / `.sum()`
(filter / transform / count / reduce) are different intents. All are valid.

### Quick vocabulary

Familiar-equivalent columns relate VOL to C, Go, Rust, or Python concepts; they
do not claim identical implementation.

#### Keywords

| Word | Meaning | Familiar equivalent | Status |
| --- | --- | --- | --- |
| `and` | Both Boolean sides true | C `&&`, Python `and` | Supported |
| `else` | Final alternate `if` block | `else` | Supported |
| `elif` | Extra `if` branch | `else if` / `elif` | Supported |
| `false` | Boolean false | `false` | Supported |
| `const` | Opt-in immutable binding (shallow) | `const` / immutable `let` | Supported |
| `err` | Failed Result value (`err(x)`) | `Err(x)` | Supported |
| `fn` | Function declaration or anonymous function | `func` / `fn` | Supported |
| `export` | Make names public | export list / `pub` | Supported |
| `if` | Conditional statement | `if` | Supported |
| `import` | Load another module’s exports | `import` / `use` | Supported |
| `match` | Removed (SF-1); parse error `E153` | — | Rejected |
| `none` | Absent Option value | `None` / `nil` absence | Supported |
| `not` | Boolean negation | C `!`, Python `not` | Supported |
| `ok` | Successful Result value (`ok(x)`) | `Ok(x)` | Supported |
| `or` | Either Boolean side true | C logical or, Python `or` | Supported |
| `print` | Write one or more values (space-joined) | `println` | Supported |
| `repeat` | Run body N times | counting `for` | Supported |
| `return` | Exit function with value | `return` | Supported |
| `some` | Present Option value (`some(x)`) | `Some(x)` | Supported |
| `struct` | Product type declaration | `struct` / class fields | Supported |
| `true` | Boolean true | `true` | Supported |
| `while` | Loop while condition | `while` | Supported |

#### Collection words, built-ins, and symbols

| Form | Meaning | Familiar equivalent | Status |
| --- | --- | --- | --- |
| `items.each item { ... }` | Per-element block, in order | `for item in items` | Supported |
| `value.len` | Array length or string Unicode scalars | `len(value)` | Supported |
| `string.byte_len` | UTF-8 byte count of a string | `len(s)` in Go | Supported |
| `array.copy()` | Shallow copy of an array | `clone` / `[...arr]` | Supported |
| `array.deep_copy()` | Recursive deep clone of an array | deep clone | Supported |
| `items.where(condition)` | Eager filter; `_` is current item | eager `filter` | Supported |
| `items.map(transform)` | Eager map; `_` is current item | eager `map` | Supported |
| `items.count()` / `items.count(condition)` | Length, or eager count of matches | `len` / filter + length | Supported |
| `items.sum()` | Left-fold `+` from integer `0` | `sum` / `reduce` | Supported |
| `input()` / `input(prompt)` | Read one line | stdin / `readLine` | Supported |
| `assert(cond)` / `assert(cond, msg)` | Fail when false | assertion | Supported |
| `string(value)` | Display string | `toString` | Supported |
| `args` | CLI args after source file | `argv` | Supported |
| `const name := value` | Declare immutable binding (shallow) | `const` / immutable `let` | Supported |
| `name := value` / `a, b := x, y` | Declare mutable binding(s) | `var` / multi-declare | Supported |
| `name = value` / `a, b = x, y` | Assign binding(s) (RHS first) | assignment / swap | Supported |
| `cond ? a : b` | Expression conditional | JS ternary | Supported |
| `option ?? default` | Option coalesce (`some` → value) | nullish coalesce | Supported |
| `result?` | Result propagate (inside functions) | Rust `?` | Supported |
| `fn name(params) { ... }` | Named function (private by default) | unexported fn | Supported |
| `fn(params) { ... }` / `fn(params) expr` | Anonymous function (block or expr body) | lambda / closure | Supported |
| `some(value)` / `none` | Option present / absent | `Option` / `Some`/`None` | Supported |
| `ok(value)` / `err(value)` | Result success / failure | `Result` / `Ok`/`Err` | Supported |
| `if some x := opt { … } else { … }` | Option if-let unwrap | if-let | Supported |
| `if ok x := res { … } else err e { … }` | Result if-let unwrap | if-let | Supported |
| `struct Name { fields }` | Product type | struct | Supported |
| `Name { field: expr, … }` | Named struct construction | struct literal | Supported |
| `Name { expr, … }` | Positional struct construction | positional init | Supported |
| `import "path"` | Import module exports | import | Supported |
| `export name` | Public name | export list | Supported |
| `{ ... }` | Block / scope | block | Supported |
| `[a, b, c]` | Array literal | array | Supported |
| `items[index]` | Index read/write | indexing | Supported |
| `// text` | Line comment | `//` | Supported |

Capitalization has no visibility meaning. Only `export` makes a name public.

---

## 1. Program model

A VOL program is one source file containing a sequence of statements.

Execution order:

1. Lex the source into tokens.
2. Parse tokens into an AST.
3. Resolve names and known call arities.
4. Install built-ins and module-level functions.
5. Execute top-level statements in source order.

There is no `main`. Top-level statements run from top to bottom.

`export` lists are module metadata. In the current prototype they do not change
runtime behavior beyond marking functions public for later module systems.

---

## 2. Lexical grammar

Source is UTF-8 text.

### 2.1 Ignored input

Between tokens, the lexer skips:

- spaces and tabs (and other Unicode space except newline)
- line comments: `//` through the end of the line, not including the newline

Newlines are **not** ignored. They are tokens used as statement separators.

### 2.2 Tokens

```text
newline      = U+000A
identifier   = (letter | "_") (letter | digit | "_")*
integer      = digit+
float        = digit+ "." digit+
string       = '"' string-char* '"'
keyword      = "and" | "const" | "elif" | "else" | "err" | "export" | "false" | "fn"
             | "if" | "import" | "match" | "none" | "not" | "ok" | "or" | "print"
             | "repeat" | "return" | "some" | "struct" | "true" | "while"
operator     = ":=" | ":" | "??" | "?" | "=" | "==" | "!=" | "<" | "<=" | ">" | ">="
             | "+" | "-" | "*" | "/" | "." | ","
delimiter    = "{" | "}" | "[" | "]" | "(" | ")"
```

Notes:

- Keywords are reserved; they are never identifiers.
- `match` remains a reserved keyword but is **Rejected** at parse time (`E153`).
- Identifiers are case-sensitive. `Double` and `double` are different names.
- Capitalization has **no** visibility meaning.
- `!` alone is illegal; use `not` for Boolean negation (`E002`).
- `:` is a token used in `? :` and as part of `:=`.
- `??` is Option coalesce; lone `?` is either ternary (`? :`) or Result postfix
  propagate (see §4.1 / §4.2.2 / §5.10).
- Unexpected characters are `E003`.

### 2.3 Numeric literals

- Integer literals are decimal only. They must fit in signed 64-bit range.
- Float literals require digits on both sides of `.` (so `1.` and `.5` are illegal).
- Float literals must be finite `f64` values.
- Out-of-range numeric literals are rejected with `E006`.
- Negative numbers are produced by unary `-`, not by a negative literal token.

### 2.4 String literals

Strings are double-quoted and may not span source newlines.

Escapes:

| Escape | Meaning |
| --- | --- |
| `\n` | newline |
| `\r` | carriage return |
| `\t` | tab |
| `\\` | backslash |
| `\"` | quote |

Unknown escapes are `E005`. Unclosed strings are `E004`.

### 2.5 Line joining

Expressions do **not** continue across newlines.

Illegal today:

```vol
total := 1 +
2
```

```vol
value.foo
.bar()
```

Each simple statement must end before the next statement begins, usually with a
newline. A closing `}` is self-delimiting, so another statement may follow it on
the same line.

---

## 3. Values and types

The prototype is dynamically typed. Every runtime value has one of these kinds:

| Kind | Examples | Internal representation |
| --- | --- | --- |
| integer | `0`, `42`, `-3` | signed 64-bit |
| float | `3.14`, `-0.5` | 64-bit IEEE floating point |
| Boolean | `true`, `false` | Boolean |
| string | `"hi"` | UTF-8 text |
| array | `[1, true, "x"]` | ordered sequence of values |
| function | `fn add(a, b) { ... }` or `fn(x) { ... }` | callable closure |
| option | `some(1)`, `none` | present or absent value (see §3.4) |
| result | `ok(1)`, `err("x")` | success or failure value (see §3.5) |
| struct | `User { name: "Ada", age: 36 }` | product value (see §3.6) |
| struct type | `struct User { … }` | constructible type name |
| nothing | missing `return` result | absence of a return; not storable or usable as a normal value (see §5.8) |

`nothing` is not an Option or Result. Option (`some`/`none`) is for optional data.
Result (`ok`/`err`) is for operational success/failure values. Language bugs and
invariants still **trap** (§8) — Result does not replace traps.

There is no static type checker yet. Type mistakes are usually runtime errors
(`Rxxx`), sometimes after resolver checks (`Sxxx`).

### 3.1 Equality and comparison

- `==` and `!=` compare values.
- For two numbers, comparison uses numeric value. An integer equals a float when
  the float is a whole number in range and matches exactly.
- Arrays compare by deep structural equality.
- Ordering operators `< <= > >=` require numbers.
- Mixed integer/float arithmetic and comparisons promote through floating-point
  evaluation when both sides are numbers and at least one is float.

### 3.2 Strings

- Strings are immutable values.
- `+` concatenates two strings, or a string and a displayable right-hand value
  (`"n=" + 7` → `"n=7"`) using the same display rules as `print` / `string()`.
- Non-string left + string right (for example `1 + "a"`) is a runtime error
  (`R013`). Either side `nothing` is `R029`.
- **String `.len` (decided):** counts Unicode scalar values (same idea as Go
  runes), not UTF-8 bytes and not grapheme clusters. The short spelling `.len`
  is canonical; `.length` is rejected with `R007` and a `fix` suggesting `.len`.
  Zero-arg `.count()` is also Supported as a length call (same result as `.len`).
- **`.byte_len`** returns the UTF-8 byte count of a string as an integer. For
  ASCII strings this equals `.len`; for multi-byte Unicode it is larger.
  Useful for systems I/O that measures bytes, not characters.
  Requires a string receiver; non-string → `R033`.

### 3.3 Arrays

- Array literals create a new array.
- Arrays hold mixed values.
- Indexing uses zero-based integer indices.
- Out-of-range indices are runtime errors (`R016`).
- Float indices are rejected (`R015`).
- Strings cannot be indexed (`R003`).
- **Array identity (decided):** assigning an array with `:=` or `=` copies the
  **reference**, not a deep clone. Indexed assignment mutates the shared array.
  The same sharing applies when an array is passed as a function argument: the
  callee receives the same array value, so indexed writes are visible to the
  caller.
- Rebinding a name to a different array (for example `b = [9]`) updates only that
  binding; other names that still reference the previous array are unchanged.
- `.where` always returns a **new** array. Mutating the result does not mutate
  the source.
- **`.copy()`** returns a shallow copy of the array: a new top-level `[]any` with
  the same elements. Mutations to the copy's indices do not affect the original,
  but nested arrays within the copy still share identity with the originals.
- **`.deep_copy()`** returns a recursive clone: all nested arrays are also cloned.
  Non-array values (integers, booleans, strings, functions) are copied by value
  semantics and are not affected by either clone operation.
- Both calls require an array receiver; non-array → `R031` (`.copy()`) or
  `R032` (`.deep_copy()`).
- Move semantics and ownership are out of scope until designed separately.

Example of shared mutation:

```vol
a := [1, 2]
b := a
b[0] = 9
print a // [9, 2]
```

Example of rebinding vs sharing:

```vol
a := [1, 2]
b := a
b = [9]
print a // [1, 2]
print b // [9]
```

### 3.4 Option values

Option represents optional data: a value may be present or absent.

```vol
found := some("VOL")
missing := none
```

Rules:

- `some(expression)` wraps a usable value (not `nothing` — `R029` applies).
- `none` is the absent Option.
- Display: `some(…)` with recursive display of the inner value, or `none`.
- Empty string / empty array are ordinary data, not `none`.
- Option is distinct from `nothing` (missing return) and from Result
  (operational failure).
- Unwrap with if-let or `??` (§5.10 / §4.2.2). Using Option where a bare `T` is
  required (e.g. arithmetic) fails with a type diagnostic as for other
  mismatched kinds.

Nested `some(some(x))` is allowed; there is no automatic flatten.

### 3.5 Result values

Result represents operational success or failure as a value.

```vol
r := ok(7)
e := err("missing")
```

Rules:

- `ok(expression)` / `err(expression)` wrap a usable value (not `nothing` — `R029`).
- Display: `ok(…)` or `err(…)` with recursive display of the inner value.
- Distinct from Option and from `nothing`.
- Unwrap with if-let or postfix `?` inside functions (§5.10). Using Result where
  a bare `T` is required fails with a type diagnostic as for other mismatched
  kinds.
- Dual-return sugar remains Planned ([`IDEAS.md`](IDEAS.md)).
- Built-ins such as `input` and `assert` still trap; they do not return Result yet.

### 3.6 Struct values

```vol
struct User {
    name
    age
}
u := User { name: "Ada", age: 36 }
v := User { "Ada", 36 }
print u.name
u.age = 37
```

Rules:

- `struct Name { field… }` declares a module-scoped product type (at least one field).
- Construction (every field required):
  - named: `Name { field: expression, … }` — unknown fields rejected (`R039`/`R040`)
  - positional: `Name { expression, … }` in declaration field order — arity must
    match (`R043`)
- Field access and assignment use `.` (`R036` for unknown fields).
- Struct instances share identity on assignment/arguments (like arrays).
- `const u` is shallow: the binding cannot be reassigned (`S030`/`R030`), but
  fields remain mutable.
- Equality compares fields deeply (via the same structural rules as other values).
- Display: `User { name: …, age: … }` in declaration field order.
- No methods; enums / tagged unions / pattern match on user tags are Planned.
- `export User` makes the type importable.

---

## 4. Expression grammar

```text
expression   = conditional
conditional  = coalesce (("?" conditional ":" conditional) | "?")?
coalesce     = or ("??" coalesce)?
or           = and ("or" and)*
and          = not ("and" not)*
not          = "not" not | equality
equality     = comparison (("==" | "!=") comparison)*
comparison   = term (("<" | "<=" | ">" | ">=") term)*
term         = factor (("+" | "-") factor)*
factor       = unary (("*" | "/") unary)*
unary        = "-" unary | postfix
postfix      = primary (call | index | property)*
call         = "(" args? ")"
index        = "[" expression "]"
property     = "." identifier
args         = expression ("," expression)*
primary      = literal | identifier | array | "(" expression ")"
             | "some" "(" expression ")" | "none"
             | "ok" "(" expression ")" | "err" "(" expression ")"
             | function-expr | struct-literal
array        = "[" args? "]"
literal      = integer | float | string | "true" | "false"
function-expr = "fn" "(" params? ")" (block | expression)
struct-literal = identifier "{" struct-fields? "}"
struct-fields  = named-fields | positional-fields
named-fields   = identifier ":" expression ("," identifier ":" expression)*
positional-fields = expression ("," expression)*
```

Struct literals are recognized when `{` is followed by `field:`, a value expression,
or similar (so `if cond {` remains a block). Bare `{ }` is never a struct literal.

A trailing `?` after a coalesce expression is Result postfix propagate when it is
not the start of `? :` (see §5.10). `??` is Option coalesce (§4.2.2).

### 4.1 Precedence

Lowest to highest:

1. `? :` and postfix Result `?` (disambiguated: `?` then `:` is ternary;
   otherwise postfix when allowed — §5.10)
2. `??` (right-associative)
3. `or`
4. `and`
5. `not`
6. `==` `!=`
7. `<` `<=` `>` `>=`
8. `+` `-`
9. `*` `/`
10. unary `-`
11. call, index, `.property`

Binary operators at the same level associate left-to-right.
Unary operators associate right-to-left.
`? :` associates right-to-left: `a ? b : c ? d : e` means `a ? b : (c ? d : e)`.
`??` associates right-to-left: `a ?? b ?? c` means `a ?? (b ?? c)`.

Examples:

| Source | Meaning |
| --- | --- |
| `a + b * c` | `a + (b * c)` |
| `not a == b` | `not (a == b)` |
| `true or false and x` | `true or (false and x)` |
| `a ? b : c ? d : e` | `a ? b : (c ? d : e)` |
| `a ?? b ?? c` | `a ?? (b ?? c)` |

### 4.2 Boolean operators

- `and` and `or` require Boolean operands.
- They short-circuit:
  - `false and x` does not evaluate `x`
  - `true or x` does not evaluate `x`
- `not` requires a Boolean operand.
- `not` binds after comparisons and before `and` / `or`; use parentheses to
  negate a wider expression, such as `not (ready and allowed)`.

### 4.2.1 Conditional operator `? :`

```vol
ready ? "yes" : "no"
```

- The condition must be Boolean (`R004`).
- Only the selected branch is evaluated (short-circuit).
- Both branches must produce a usable value (`R029` rejects `nothing`).
- This is the **expression** form for choosing a value. `if` remains a
  statement (see §5.3). Expression-`if` (`x := if ...`) is not accepted.
- A `?` that starts a ternary requires `:`; missing `:` is `E001` when the `?`
  is not Result postfix propagate (§5.10).

### 4.2.2 Option coalesce `??`

```vol
print maybe ?? "missing"
```

- Left operand must be an Option (`R041` otherwise).
- Result values are rejected (`R042`) — use if-let or postfix `?` for Result.
- If left is `some(x)`, yield `x` and do **not** evaluate the right operand.
- If left is `none`, evaluate and yield the right operand.
- Both sides must produce a usable value (`R029` rejects `nothing`).

### 4.3 Arithmetic

- Integer op integer stays integer, except that `/` is truncating integer division.
- If either operand is float, both are treated as floats.
- **Integer overflow (decided):** signed 64-bit integer `+`, `-`, `*`, `/`, and
  unary `-` **trap** on overflow with runtime error `R028`. The diagnostic
  includes a machine-readable `fix` suggestion so humans and LLMs can repair
  the program. Silent wrapping is **not** the default.
- Wrapping arithmetic and selectable build modes (for example debug trap /
  release wrap) are **Planned**—see [`IDEAS.md`](IDEAS.md).
- Division by zero is `R014` (separate from overflow).
- Floating-point arithmetic follows IEEE rules and is not covered by `R028`.

### 4.4 Properties and collection forms

Known properties:

| Form | Meaning |
| --- | --- |
| `array.len` | element count as integer |
| `string.len` | Unicode scalar count as integer (not bytes) |
| `string.byte_len` | UTF-8 byte count as integer (`R033` on non-string) |
| `array.copy()` | shallow copy of the array (`R031` on non-array) |
| `array.deep_copy()` | recursive deep clone of the array (`R032` on non-array) |
| `array.sum()` | left fold of `+` starting from integer `0` |
| `array.where(condition)` | eager filter; see below |
| `array.map(transform)` | eager map; see below |
| `array.count()` | element count (same as `.len`) |
| `array.count(condition)` | eager count of matches; see below |
| `string.count()` | Unicode scalar count (same as `.len`) |

Unknown properties are `R007`. Using `.length` instead of `.len` is `R007` with
a `fix` pointing at `.len`. Using `.each` as a property or call (for example
`.each(fn...)`) is `R007` with a `fix` pointing at statement form
`items.each item { ... }`.

#### `.where`

```vol
items.where(_ > 5)
```

Semantics:

1. Evaluate `items`; it must be an array (`R021`).
2. Require exactly one argument (`R020` / resolver `S003`).
3. For each element in order:
   - bind `_` to that element in a fresh scope
   - evaluate the condition
   - condition must be Boolean (`R022`); a `fn` value is not Boolean — `R022`
     includes a `fix` suggesting a `_` expression such as `.where(_ > 5)`
   - if true, append the original element to a new result array
4. Return the new array.

The condition may read enclosing variables. Nested `.where` calls each bind their
own `_`.

**Purity (Planned):** the intended intent split is that `.where` / `.map` /
`.count` express filtering, transforms, and counting, while `.each` is for
imperative per-item work (`print`, mutation, and similar). Prefer side-effect-free
predicates and transforms.

This prototype does **not** reject impure predicates: expressions are evaluated
eagerly left to right, so helpers with effects may appear to “work.” Do not
depend on that. Future purity checks (and any fusion or parallelization that
assumes purity) are Planned in [`IDEAS.md`](IDEAS.md); until they exist, purity
is guidance, not an enforced language rule.

#### `.sum()`

```vol
items.sum()
```

Semantics:

1. Require an array (`R019`).
2. Start from integer `0`.
3. Add each item left to right with the normal `+` rules.
4. Empty array yields `0`.
5. Non-numeric items fail with `R013`.

#### `.map`

```vol
items.map(_ * 2)
```

Semantics:

1. Evaluate `items`; it must be an array (`R021`).
2. Require exactly one argument (`R020`).
3. For each element in order:
   - bind `_` to that element in a fresh scope
   - evaluate the transform expression
   - append the result to a new array
4. Return the new array.

#### `.count`

```vol
items.count()
items.count(_ > 2)
```

Semantics:

**Zero arguments** (length):

1. Evaluate the receiver; it must be an array or string (`R021`).
2. Return the same integer as `.len` (array element count or string Unicode
   scalars).

**One argument** (match count):

1. Evaluate `items`; it must be an array (`R021`).
2. For each element in order:
   - bind `_` to that element in a fresh scope
   - evaluate the condition (must be Boolean — `R022`, with the same `_`
     expression `fix` as `.where` when the value is not Boolean)
   - if true, increment a counter
3. Return the count as an integer.

Other arities are `S003` / `R020` (`0 or 1` arguments). `.len` remains the
canonical length **property**; `.count()` is the zero-arg call form.

`.where` / `.map` / `.count` / `.sum()` are not parallel and not lazy in this
prototype. Future fusion or parallelization of collection pipelines depends on
Planned purity checks (treat `.map` transforms and `.count` predicates like
`.where` for that future rule).

---

## 5. Statement grammar

```text
program     = statement*
statement   = block-stmt
            | export-stmt
            | import-stmt
            | struct-decl
            | function-decl
            | return-stmt
            | print-stmt
            | if-stmt
            | if-let-stmt
            | repeat-stmt
            | while-stmt
            | declaration
            | multi-declaration
            | assignment
            | multi-assignment
            | each-stmt
            | expression-stmt

block       = "{" statement* "}"
block-stmt  = block
declaration = "const"? identifier ":=" expression
multi-declaration = "const"? identifier ("," identifier)+ ":=" expression ("," expression)+
assignment  = (identifier | index-expr | property-expr) "=" expression
multi-assignment = identifier ("," identifier)+ "=" expression ("," expression)+
each-stmt   = expression ".each" identifier block
if-stmt     = "if" expression block ("elif" expression block)* ("else" block)?
if-let-stmt = option-if-let | result-if-let
option-if-let = "if" "some" identifier ":=" expression block "else" block
result-if-let = "if" "ok" identifier ":=" expression block "else" "err" identifier block
struct-decl = "struct" identifier "{" identifier+ "}"
import-stmt = "import" string
repeat-stmt = "repeat" expression block
while-stmt  = "while" expression block
print-stmt  = "print" expression ("," expression)*
return-stmt = "return" expression
function-decl = "fn" identifier "(" params? ")" (block | expression)
export-stmt = "export" identifier ("," identifier)*
params      = identifier ("," identifier)*
```

`match` is not a statement. Writing `match …` is parse error `E153`.

Statement separation:

- After a simple statement, expect newline, end of file, or a following `}`
  context as implemented by the parser.
- Missing separators are `E119`.

### 5.1 Blocks and scopes

`{ ... }` creates a lexical scope.

- Declarations in an inner scope may shadow outer names.
- Leaving a block discards its local declarations.
- A standalone block is a valid statement.

### 5.2 Variables

```vol
name := expression
const name := expression
name = expression
a, b := 0, 1
a, b = b, a + b
```

- `:=` declares a new **mutable** binding in the current scope.
- Multi-declare: `a, b := x, y` (also `const a, b := …`) — name count must equal
  value count (`E158`).
- Redeclaring the same name in the same scope is an error (`S001` / `R001`).
- `=` assigns to an existing variable or array index.
- Multi-assign: `a, b = x, y` — name count must equal value count (`E159`). All
  RHS expressions are fully evaluated (left to right) before any assignment, so
  swap forms such as `a, b = b, a` work.
- **Mutability model (decided):** bindings are mutable by default. Programmers do
  not mark ordinary variables as mutable. There is no `mut` keyword.
- **Opt-in immutability:** `const name := expression` declares a binding that
  cannot be reassigned with `=` (`S030` from resolver, `R030` at runtime).
  Semantics are **shallow**: rebinding is forbidden; mutating an array value
  through that binding (for example `a[0] = 9`) still follows §3.3. `const`
  applies only at declaration; there is no `const` on a bare `=`.
  Shadowing follows the same scope rules as mutable bindings. Function
  parameters stay mutable (no const parameter syntax yet).
- Binding mutability (can this name be reassigned?) is separate from value
  mutability (strings are immutable values; arrays share identity on assignment—
  see §3.3).

### 5.3 Conditions

```vol
if condition {
    ...
} elif other {
    ...
} else {
    ...
}
```

- **Decided:** `if` is a **statement**, not an expression. To choose a value in
  an expression, use `? :` (§4.2.1).
- Each condition must evaluate to Boolean (`R004`).
- Zero or more `elif` branches may appear after the `if` body.
- `else` is optional and must be last when present.
- Branches are tried in order; the first true condition runs its block.

### 5.4 `repeat`

```vol
repeat count {
    ...
}
```

- `count` must be a non-negative integer (`R005`).
- Float counts are rejected.
- The body runs `count` times.
- Each iteration uses a fresh body scope.

### 5.5 `while`

```vol
while condition {
    ...
}
```

- Same Boolean rule as `if`.
- Repeats until the condition is false.
- **Decided vocabulary:** `while` is permanent. There is no alternate or legacy
  spelling (`until`, `loop`, etc.). Use `repeat` for counted loops and `.each`
  for array iteration—those are different intents.

### 5.6 `.each`

```vol
items.each item {
    ...
}
```

- `items` must be an array (`R006`).
- Visits elements in order.
- Binds `item` in a fresh scope for each iteration.
- Intended for imperative per-item work, including mutation and `print`.
- The form is a **statement**: `expression ".each" identifier block`. After
  `.each`, an item name is required (`E103` with a `fix` suggesting
  `items.each item { ... }`). Writing `.each(...)` as a call is not this
  statement; it fails as unknown property `each` (`R007` with the same `fix`).

`.each` and `.where` express different intents. Both are valid: prefer `.where`
for filtering (and `.sum()` for reduction); use `.each` when the body has side
effects. Purity of `.where` predicates is Planned guidance — see §4.4.

### 5.7 `print`

```vol
print expression
print "label:", value
```

Writes one or more values, space-joined, followed by a newline. Arguments are
evaluated left to right; each must be a value (`R029` rejects `nothing`).

Display rules (per argument):

- integers, floats, Booleans, strings: ordinary text
- arrays: `[a, b, c]` with recursive display
- options: `some(…)` or `none`
- results: `ok(…)` or `err(…)`
- structs: `Type { field: …, … }` in field-declaration order
- no quotes are added around strings in display

Prefer `print "A grades:", n` over `"A grades: " + string(n)` when labels are
needed.

### 5.8 Functions

```vol
fn add(a, b) {
    return a + b
}

double := fn(x) {
    return x * 2
}

triple := fn(x) x * 3
```

Rules:

- Named functions are declared with `fn name(params) { ... }` or
  `fn name(params) expression` (expression body = implicit `return`).
- Anonymous functions are expressions: `fn(params) { ... }` or `fn(params) expr`
  (same keyword; no name). They may be bound, passed, or called immediately:
  `fn(x) { return x }(1)`.
- Parameters are local names in the function body.
- Module-level named functions are visible throughout the module, including before
  their declaration text.
- Nested named functions are installed when execution reaches their declaration.
- Functions (named and anonymous) capture their enclosing environment (closures).
- `return expression` exits the current function with a value.
- `return` is only legal inside a function (`E115`).
- **Missing return (decided):** If execution falls off the end without `return`,
  the call result is `nothing`.
- **Using `nothing`:** A call used as a **statement** may discard `nothing`
  (procedure-style functions). Binding or using `nothing` as a value is a
  runtime error (`R029`) with a `fix` suggestion—this includes `:=`, `=`,
  `print`, operators, conditions, indexes, call arguments, and array elements.
  Variables always expect a real value when initialized or assigned.
- Propagating `nothing` with `return callee()` is allowed; the outer caller is
  still subject to the use rules above.
- Argument count must match parameter count (`S003` / `R018`).
- Only functions/builtins are callable (`R017`).
- Typed void-vs-value functions (requiring `return` on all paths for valued
  functions) are Planned—see [`IDEAS.md`](IDEAS.md).
- JS-style `=>` is not Supported; prefer anonymous `fn`.

Visibility:

```vol
export add

fn add(a, b) {
    return a + b
}
```

- Names are private unless listed in `export`.
- Export may appear before or after the definition.
- Unknown or duplicate exports are parse errors (`E118`, `E117`).
- Capitalization does not affect visibility.

### 5.9 Built-ins

| Name | Form | Behavior |
| --- | --- | --- |
| `print` | statement | see above |
| `input` | `input()` or `input(prompt)` | read one line; trim trailing `\n` / `\r\n`; EOF with no data yields `""` |
| `assert` | `assert(cond)` or `assert(cond, message)` | fail with `R027` when `cond` is false |
| `string` | `string(value)` | convert value to display string |
| `args` | value | array of CLI arguments after the source file |

### 5.10 Option / Result unwrap (if-let, `??`, `?`)

`match` was removed in SF-1 (`E153`). Unwrap with if-let, Option `??`, or Result
postfix `?`.

#### Option if-let

```vol
if some name := maybe {
    print name
} else {
    print "missing"
}
```

Rules:

- Scrutinee must be an Option (`R034`).
- `else` block is required (`E154`).
- `some` branch binds `name` in a fresh scope when present.
- `else` runs when the value is `none`.

#### Result if-let

```vol
if ok value := result {
    print value
} else err message {
    print message
}
```

Rules:

- Scrutinee must be a Result (`R037`).
- `else err <name>` is required (`E155` / `E156` / `E157`).
- Plain `else { … }` without `err` binding is rejected.
- `ok` / `err` branches each bind their payload in a fresh scope.

#### Option coalesce

See §4.2.2 (`??`).

#### Result postfix `?`

```vol
fn twice(a, b) {
    n := divide(a, b)?
    return ok(n * 2)
}
```

Rules:

- Operand must be a Result (`R044`).
- Legal only **inside a function** (`S041` at top level).
- If `ok(x)`, the expression yields `x`.
- If `err(e)`, the enclosing function returns `err(e)` immediately (propagate).
- Does not apply to Option — use `??` or if-let.

### 5.11 Modules and `import`

```vol
import "examples/features/modules/math"
export add
```

Aliases are project-local via `vol.config.json` `paths` (see
`examples/features/modules/aliases/`):

```vol
import "@lib/math"   // expands only if the project defines "@lib" → some dir
```

Rules:

- `vol run <entry.vol>` discovers the nearest `vol.config.json` by walking from
  the entry file’s directory toward the filesystem root.
- Config fields used: `root` (project root relative to the config file directory;
  default `.`) and `paths` (alias prefix → subdirectory under the project root).
- Missing config: project root is the entry file’s directory; no aliases.
- After alias expansion, `import "P"` resolves to `{root}/P.vol` if that file
  exists, else `{root}/P/mod.vol`. Otherwise `S031`.
- Paths must not escape the project root (`S032`); `..` segments are rejected.
- Import cycles report `S033` with a deterministic path list.
- Only names listed in `export` are visible to importers (functions, bindings,
  and struct types). Flat import installs those names into the importer’s module
  scope; collisions are `S034`.
- Dependency modules execute (top-level side effects) in topological order before
  the entry module.
- Ambient tiny core is unchanged; there is no automatic stdlib. Alias names such
  as `@std` are **not reserved or magic** — they work only when a project’s
  `paths` maps them to a real directory that contains the imported `.vol` file.
  A richer standard library behind imports remains Planned
  ([`IDEAS.md`](IDEAS.md)).

Built-in names cannot be redeclared at module scope.

---

## 6. Name resolution

Before execution, VOL resolves names.

Static errors:

| Code | Meaning |
| --- | --- |
| `S001` | duplicate declaration in the same scope |
| `S002` | use of an undefined name |
| `S003` | wrong argument count for a known function/builtin/`.where`/`.map`/`.count` |
| `S030` | assignment to a `const` binding |
| `S041` | Result postfix `?` outside a function |

Additional resolution rules:

- Module functions are declared before body checking so forward calls work.
- Ordinary locals are order-sensitive: use after declare in that scope.
- `.where` / `.map` / `.count` expressions see `_` as a declared name.
- `.each` item names belong to the loop scope.
- Function parameters belong to the function scope.
- If-let bindings belong to their branch scopes.

---

## 7. Evaluation order

Unless an operator short-circuits:

1. Evaluate operands left to right.
2. For calls: evaluate callee, then arguments left to right, then invoke.
3. For arrays: evaluate elements left to right.
4. For statements: execute in source order.

Short-circuiting applies to `and`, `or`, `? :` (only the taken branch runs), and
`??` (right side skipped when left is `some`).

---

## 8. Failure model

VOL diagnostics have:

- a stable code (`E` lex, `E` parse, `S` resolve, `R` runtime)
- a message
- a source location
- an optional fix suggestion

Pipeline:

1. Lex errors abort before parse.
2. Parse errors abort before resolve.
3. Resolve errors abort before execute.
4. Runtime errors abort program execution immediately.

There is no `try` / `catch`. **Hybrid failure model (Supported):** traps for
programmer/invariant failures (`Rxxx`), plus Result values (`ok` / `err`) for
expected operational failure data, unwrapped with if-let or postfix `?`.
Dual-return sugar remains Planned — see [`IDEAS.md`](IDEAS.md).

The CLI emits diagnostics as human-readable text by default. Pass `--json`
anywhere in the command to receive a single JSON object on stderr instead:

```
vol --json run <file.vol>
vol run --json <file.vol>
```

JSON shape: `{"code":"…","message":"…","file":"…","position":{"Offset":…,"Line":…,"Column":…},"fix":"…"}`.
`fix` is omitted when absent. I/O failures (file not found) remain plain text.

### 8.1 Lexical codes

| Code | Meaning |
| --- | --- |
| `E001` | `?` without `:` |
| `E002` | `!` without `=` |
| `E003` | unexpected character |
| `E004` | unclosed string |
| `E005` | unknown string escape |
| `E006` | numeric literal out of range |

### 8.2 Parse codes

| Code | Meaning |
| --- | --- |
| `E101` | expected expression |
| `E102` | invalid assignment target |
| `E103` | missing `.each` item name (`fix`: use `items.each item { ... }`) |
| `E104` | missing `{` |
| `E105` | missing `}` |
| `E106` | missing `]` after index |
| `E107` | missing property name |
| `E108` | missing `)` in group |
| `E109` | missing `]` in array literal |
| `E110` | missing function name |
| `E111` | missing `(` after function name |
| `E112` | missing parameter name |
| `E113` | missing `)` after parameters |
| `E114` | missing `)` after call arguments |
| `E115` | `return` outside function |
| `E116` | missing export name |
| `E117` | duplicate export |
| `E118` | export of unknown name |
| `E119` | missing statement separator newline |
| `E120` | expected a name after `const` / in a name list |
| `E121` | expected `:=` after const name (or `:=`/`=` after names) |
| `E129` | expected `(` after `some` |
| `E130` | expected `)` after `some` value |
| `E131` | expected string path after `import` |
| `E132` | expected type name after `struct` |
| `E133` | expected `{` after struct name |
| `E134` | expected field name in `struct` |
| `E135` | duplicate struct field |
| `E136` | expected `}` to close `struct` |
| `E137` | struct with no fields |
| `E144`–`E148` | struct literal punctuation / fields |
| `E149`–`E152` | `ok` / `err` call punctuation |
| `E153` | `match` removed — use if-let / `??` / Result `?` |
| `E154` | Option if-let missing `else` |
| `E155` | Result if-let missing `else err` |
| `E156` | Result if-let `else` without `err` binding |
| `E157` | expected binding name after `err` |
| `E158` | multi-declare name/value count mismatch |
| `E159` | multi-assign name/value count mismatch |

### 8.2.1 Module / resolve codes (loader)

| Code | Meaning |
| --- | --- |
| `S031` | module not found / unreadable / imports without loader |
| `S032` | import escapes project root or contains `..` |
| `S033` | import cycle |
| `S034` | imported name collides |
| `S035` | invalid / unreadable `vol.config.json` or unknown alias |
| `S036` | export missing at link time |
| `S037` | nested `struct` declaration |
| `S038` | construct non-struct type |

### 8.3 Runtime codes

| Code | Meaning |
| --- | --- |
| `R001` | duplicate variable in scope |
| `R002` | unknown variable |
| `R003` | index into non-array |
| `R004` | non-Boolean condition |
| `R005` | invalid `repeat` count |
| `R006` | `.each` on non-array |
| `R007` | unknown property (`.length` / call-form `.each` include `fix`) |
| `R008` | `not` on non-Boolean |
| `R009` | unary `-` on non-number |
| `R010` | `and` left operand not Boolean |
| `R011` | `or` left operand not Boolean |
| `R012` | Boolean operator right operand not Boolean |
| `R013` | operator type mismatch |
| `R014` | division by zero |
| `R015` | non-integer array index |
| `R016` | array index out of bounds |
| `R017` | call of non-function |
| `R018` | wrong function argument count |
| `R019` | `.sum()` on non-array |
| `R020` | `.where` / `.map` / `.count` wrong arity |
| `R021` | `.where` / `.map` / `.count(pred)` on non-array; `.count()` on non-array/non-string |
| `R022` | `.where` / `.count` condition not Boolean (`fix`: use `_` expression) |
| `R023` | `input` prompt not string |
| `R024` | input read failure |
| `R025` | `assert` condition not Boolean |
| `R026` | `assert` message not string |
| `R027` | assertion failed |
| `R028` | integer overflow |
| `R029` | expected a value, got `nothing` |
| `R030` | assignment to a `const` binding |
| `R031` | `.copy()` on non-array |
| `R032` | `.deep_copy()` on non-array |
| `R033` | `.byte_len` on non-string |
| `R034` | Option if-let scrutinee is not an Option |
| `R035` | assign through `.` on non-struct |
| `R036` | unknown struct field |
| `R037` | Result if-let scrutinee is not a Result |
| `R038` | unknown / non-struct type in struct literal |
| `R039` | struct literal missing field |
| `R040` | struct literal unknown field |
| `R041` | `??` left operand is not an Option |
| `R042` | `??` applied to a Result |
| `R043` | positional struct literal arity mismatch |
| `R044` | postfix `?` operand is not a Result |
| `R999` | internal unsupported expression |

---

## 9. Conformance examples

These programs must keep working. They are also covered by `examples/` and tests.

### 9.1 Precedence

```vol
print 1 + 2 * 3 // 7
print true or false and false // true
print not 1 == 1 // false
```

### 9.2 Scope

```vol
value := 1
{
    value := 2
    print value // 2
}
print value // 1
```

### 9.3 Collections

```vol
numbers := [4, 7, 2, 9, 12]
large := numbers.where(_ > 5)
print large // [7, 9, 12]
print large.sum() // 28
print numbers.map(_ * 2) // [8, 14, 4, 18, 24]
print numbers.count(_ > 5) // 3
print numbers.count() // 5
print "Sum:", large.sum() // Sum: 28
```

### 9.4 Functions

```vol
fn square(n) {
    return n * n
}

print square(6) // 36

double := fn(x) {
    return x * 2
}
print double(21) // 42

triple := fn(x) x * 3
print triple(7) // 21
```

### 9.5 Option if-let and `??`

```vol
maybe := some(7)
if some n := maybe {
    print n // 7
} else {
    print "missing"
}
print maybe ?? 0 // 7
print none // none
```

### 9.6 Result if-let and `?`

```vol
r := ok(7)
if ok n := r {
    print n // 7
} else err msg {
    print msg
}

fn divide(a, b) {
    if b == 0 {
        return err("zero")
    }
    return ok(a / b)
}
fn twice(a, b) {
    n := divide(a, b)?
    return ok(n * 2)
}
```

### 9.7 Structs

```vol
struct User {
    name
    age
}
u := User { name: "Ada", age: 36 }
v := User { "Ada", 36 }
print u.name // Ada
print v.age // 36
```

### 9.8 Multi-assign

```vol
a, b := 0, 1
a, b = b, a + b
print a // 1
print b // 1
```

### 9.9 Imports

```vol
import "examples/features/modules/math"
print add(21, 21) // 42
```

---

## 10. Explicitly out of scope

Do not treat these as specified just because vision docs mention them:

- static types and typed signatures
- `const` parameter syntax (function parameters stay mutable; Planned)
- ownership, borrowing, lifetimes (direction: local escape first, API contracts
  later — [`IDEAS.md`](IDEAS.md); not implemented)
- allocation inference / explicit allocators (unspecified; do not claim inference)
- enums, tagged unions, methods on structs (product `struct` is Supported; tags
  Planned — [`IDEAS.md`](IDEAS.md))
- generics
- symbol-selection imports (`import { x } from …`); multi-file folder packages
  beyond `path.vol` / `path/mod.vol`
- ambient growth of the prelude beyond today’s tiny built-ins
- `parallel`, async, channels (no parallelization guarantees until specified)
- wrapping integer arithmetic and overflow/bounds build modes (default is trap;
  see §4.3; modes + wrap ops Planned in [`IDEAS.md`](IDEAS.md))
- dual-return sugar over Result; converting `input`/`assert` to Result
  (Result values + if-let / `?` are Supported; traps remain for bugs — §8)
- `=>` arrow functions and `|>` pipeline sugar (directions in [`IDEAS.md`](IDEAS.md);
  anonymous / expression-body `fn` and method chains are Supported)
- richer pattern matching / `match` (removed in SF-1 surface; if-let is the Supported unwrap)
- native memory layout / allocators
- C or LLVM backends

When one of these is designed, add a numbered section here **and** implement
tests before calling it Supported.

---

## 11. Decided core rules

These core rules are settled for the prototype. Track remaining open questions
and Planned work in [`IDEAS.md`](IDEAS.md).

- **Mutability default (§5.2):** bindings are **mutable by default**. Opt-in
  immutability uses `const name := expression` with shallow semantics
  (implemented; `S030`/`R030` on reassignment).
- **Array assignment (§3.3):** assignment and argument passing **share** the
  array reference. Use `.copy()` for a shallow clone or `.deep_copy()` for a
  recursive clone. Move/ownership are not implied.
- **Integer overflow (§4.3):** overflow **traps** (`R028` with a `fix`
  suggestion). Wrapping ops / build modes are Planned.
- **`.where` / `.map` / `.count` purity (§4.4):** **Planned** — prefer
  side-effect-free predicates/transforms and use `.each` for effects; the
  prototype does not enforce purity today (see [`IDEAS.md`](IDEAS.md)).
- **Missing return (§5.8):** fall-off yields `nothing`; discarding in a call
  statement is OK; assigning or using `nothing` as a value is `R029`.
- **String / array `.len` (§3.2, §4.4):** canonical short property; string `.len`
  is Unicode scalars; byte count is `.byte_len` (Supported, `R033` on non-string).
- **`while` (§5.5):** permanent Supported vocabulary; no alternate spellings.
- **`if` (§5.3):** statement only, with `elif` / `else`. Value choice uses
  `? :` (§4.2.1). Expression-`if` is not part of the language.
- **Error model (§3.5, §5.10, §8):** hybrid — traps for bugs/invariants;
  Result values (`ok`/`err`) with if-let and postfix `?` (functions only) for
  operational failure data. No exception-primary model. Dual-return sugar
  Planned. `input`/`assert` still trap. `match` is Rejected (`E153`).
- **Option (§3.4, §5.10):** explicit Option with `some` / `none`; unwrap with
  if-let and `??`; distinct from Result and `nothing`; no null. Optional `?T`
  sugar remains Planned.
- **Anonymous / expression-body `fn` (§5.8):** Supported; `=>` rejected; `_` in
  `.where` / `.map` / `.count` unchanged. Pipeline `|>` remains Planned.
- **Multi-assign (§5.2):** `a, b := …` / `a, b = …` Supported (RHS evaluated
  before assigns).
- **Collection map/count (§4.4):** `.map(_)` and `.count(_)` Supported;
  `.count()` (0 args) is length (same as `.len`) on arrays and strings (SF-2).
- **Multi-arg `print` (§5.7):** `print a, b` space-joins display forms (SF-2).
- **String `+` coercion (§3.2):** `string + displayable` concatenates via
  display rules; `non-string + string` stays `R013` (SF-2).
- **Structs (§3.6):** product `struct` with named and positional literals
  Supported. Tagged unions/enums and methods remain Planned.
- **Modules (§5.11):** `import` + `vol.config.json` discovery/aliases Supported;
  ambient tiny core; explicit path for everything else.
- **Ownership / allocation direction (Planned):** local escape analysis first;
  API contracts later; allocation unspecified — do not claim inference.
- **Parallel direction (Planned):** `parallel` stays Planned; no guarantees
  until scheduling and failure are specified.
- **Build modes direction (Planned):** trap default for overflow/bounds; modes
  and/or explicit wrap ops later.
- **Pipeline direction (Planned):** method chains Supported; if pipes later,
  `|>` not bare `|`.

Implementers and LLMs must follow the concrete behavior in sections 1–9.

---

## 12. Intended formatting (not implemented)

`vol fmt` exists as a **CLI stub**: it parses the file and reports diagnostics,
but does **not** rewrite source or enforce format equality. Canonical
presentation rules (for a future rewriter):

- deterministic, idempotent output
- four-space indentation
- opening braces on the declaration or control-flow line (`fn`, `if`, `while`,
  `repeat`, `.each`, `struct`, …); closing braces on their own line
- one statement per line; blank line between top-level declarations when helpful
- spaces around binary operators and around `:=` / `=`
- no space before `(` in calls; space after `,`
- prefer `_` over `fn` in `.where` / `.map` / `.count` when equivalent
- stable comment placement (exact algorithm TBD with the AST rewriter)
- no mandatory semicolons

Commands (rewrite / `--check` format-equality still Planned):

```text
vol fmt file.vol
vol fmt .
vol fmt --check file.vol
```

Details live in [`IDEAS.md`](IDEAS.md).

---

## 13. Change process

When changing language behavior:

1. If the change adds, removes, or renames Supported surface, **bump the surface
   freeze** (SF-2 → SF-3, …), add a new language card (`vol_v3.md`, …), and note
   the bump in [`IDEAS.md`](IDEAS.md) / README. Pure bugfixes and doc sync stay
   on the active freeze without a bump. Do not keep intermediate draft cards.
2. Update this specification (including the quick vocabulary tables when forms change).
3. Add or adjust tests in `internal/lang`.
4. Update `README.md` examples if users should learn the change.
5. Move completed plans out of `IDEAS.md`.
6. Run:

```text
go test ./...
go run ./cmd/vol run ./examples/basics/first.vol
git diff --check
```

A behavior is not official until it is specified here and covered by a test.
