# VOL Language Specification (Prototype v0)

> Status: **living draft for the current interpreter**  
> Audience: humans and LLMs  
> Source of truth for behavior: this file plus the tests in `internal/lang`

This is VOL's first real language specification. It describes **only what the
tree-walking interpreter does today**. If something is not written here, it is
not part of the language yet—even if it appears in vision docs.

You do not need prior language-design experience to maintain this file. When you
change the lexer, parser, resolver, or interpreter, update the matching section
here and add a test.

Related docs:

| Document | Role |
| --- | --- |
| [`README.md`](README.md) | Project status and examples |
| [`IDEAS.md`](IDEAS.md) | Planned features and open questions |
| [`AGENTS.md`](AGENTS.md) | Project vision and contribution rules |

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

### Syntax principles

- Use braces for every multi-line scope.
- Use newlines to separate statements; do not require semicolons.
- Do not require parentheses around conditions.
- Prefer type inference when enough information is available.
- Prefer one canonical representation for each distinct intent.
- Remove boilerplate without removing structural clarity.
- Optimize syntax for humans, language models, formatters, and static analysis.

`.each` (imperative iteration) and `.where` / `.sum` (filter / reduce) are
different intents. Both are valid.

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
| `fn` | Function declaration | `func` / `fn` | Supported |
| `export` | Make names public | export list / `pub` | Supported |
| `if` | Conditional statement | `if` | Supported |
| `not` | Boolean negation | C `!`, Python `not` | Supported |
| `or` | Either Boolean side true | C logical or, Python `or` | Supported |
| `print` | Write a value | `println` | Supported |
| `repeat` | Run body N times | counting `for` | Supported |
| `return` | Exit function with value | `return` | Supported |
| `true` | Boolean true | `true` | Supported |
| `while` | Loop while condition | `while` | Supported |

#### Collection words, built-ins, and symbols

| Form | Meaning | Familiar equivalent | Status |
| --- | --- | --- | --- |
| `items.each item { ... }` | Per-element block, in order | `for item in items` | Supported |
| `value.len` | Array length or string Unicode scalars | `len(value)` | Supported |
| `items.where(condition)` | Eager filter; `_` is current item | eager `filter` | Supported |
| `items.sum` | Left-fold `+` from integer `0` | `sum` / `reduce` | Supported |
| `input()` / `input(prompt)` | Read one line | stdin / `readLine` | Supported |
| `assert(cond)` / `assert(cond, msg)` | Fail when false | assertion | Supported |
| `string(value)` | Display string | `toString` | Supported |
| `args` | CLI args after source file | `argv` | Supported |
| `const name := value` | Declare immutable binding (shallow) | `const` / immutable `let` | Supported |
| `name := value` | Declare mutable binding | `var` / mutable `let` | Supported |
| `name = value` | Assign to existing binding | assignment | Supported |
| `cond ? a : b` | Expression conditional | JS ternary | Supported |
| `fn name(params) { ... }` | Private function by default | unexported fn | Supported |
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
keyword      = "and" | "elif" | "else" | "export" | "false" | "fn" | "if"
             | "not" | "or" | "print" | "repeat" | "return" | "true" | "while"
operator     = ":=" | ":" | "?" | "=" | "==" | "!=" | "<" | "<=" | ">" | ">="
             | "+" | "-" | "*" | "/" | "." | ","
delimiter    = "{" | "}" | "[" | "]" | "(" | ")"
```

Notes:

- Keywords are reserved; they are never identifiers.
- Identifiers are case-sensitive. `Double` and `double` are different names.
- Capitalization has **no** visibility meaning.
- `!` alone is illegal; use `not` for Boolean negation (`E002`).
- `:` is a token used in `? :` and as part of `:=`.
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
| function | `fn add(a, b) { ... }` | callable closure |
| nothing | missing `return` result | absence of a value; not storable or usable as a normal value (see §5.8) |

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
- `+` concatenates two strings.
- `"a" + 1` is a runtime error (`R013`).
- **String `.len` (decided):** counts Unicode scalar values (same idea as Go
  runes), not UTF-8 bytes and not grapheme clusters. The short spelling `.len`
  is canonical; `.length` is rejected with `R007` and a `fix` suggesting `.len`.
- A separate byte-count property (for example `.byte_len`) is **Planned** for
  systems I/O—see [`IDEAS.md`](IDEAS.md).

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
- Explicit cloning is **Planned** (for example `a.copy` or `copy(a)`); spelling
  is not accepted until implemented—see [`IDEAS.md`](IDEAS.md). Move semantics
  and ownership are out of scope until designed separately.

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

---

## 4. Expression grammar

```text
expression  = conditional
conditional = or ("?" conditional ":" conditional)?
or          = and ("or" and)*
and         = equality ("and" equality)*
equality    = comparison (("==" | "!=") comparison)*
comparison  = term (("<" | "<=" | ">" | ">=") term)*
term        = factor (("+" | "-") factor)*
factor      = unary (("*" | "/") unary)*
unary       = ("not" | "-") unary | postfix
postfix     = primary (call | index | property)*
call        = "(" args? ")"
index       = "[" expression "]"
property    = "." identifier
args        = expression ("," expression)*
primary     = literal | identifier | array | "(" expression ")"
array       = "[" args? "]"
literal     = integer | float | string | "true" | "false"
```

### 4.1 Precedence

Lowest to highest:

1. `? :` (right-associative)
2. `or`
3. `and`
4. `==` `!=`
5. `<` `<=` `>` `>=`
6. `+` `-`
7. `*` `/`
8. unary `not` `-`
9. call, index, `.property`

Binary operators at the same level associate left-to-right.
Unary operators associate right-to-left.
`? :` associates right-to-left: `a ? b : c ? d : e` means `a ? b : (c ? d : e)`.

Examples:

| Source | Meaning |
| --- | --- |
| `a + b * c` | `a + (b * c)` |
| `not a == b` | `(not a) == b` |
| `true or false and x` | `true or (false and x)` |
| `a ? b : c ? d : e` | `a ? b : (c ? d : e)` |

### 4.2 Boolean operators

- `and` and `or` require Boolean operands.
- They short-circuit:
  - `false and x` does not evaluate `x`
  - `true or x` does not evaluate `x`
- `not` requires a Boolean operand.

### 4.2.1 Conditional operator `? :`

```vol
ready ? "yes" : "no"
```

- The condition must be Boolean (`R004`).
- Only the selected branch is evaluated (short-circuit).
- Both branches must produce a usable value (`R029` rejects `nothing`).
- A `?` without a following `:` is `E001`.
- This is the **expression** form for choosing a value. `if` remains a
  statement (see §5.3). Expression-`if` (`x := if ...`) is not accepted.

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
| `array.sum` | left fold of `+` starting from integer `0` |
| `array.where(condition)` | eager filter; see below |

Unknown properties are `R007`. Using `.length` instead of `.len` is `R007` with
a `fix` pointing at `.len`.

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
   - condition must be Boolean (`R022`)
   - if true, append the original element to a new result array
4. Return the new array.

The condition may read enclosing variables. Nested `.where` calls each bind their
own `_`.

**Purity (decided):** `.where` expresses **pure filtering**. The predicate must
not rely on observable side effects (mutation of enclosing state, I/O, or other
effects). Read-only use of `_` and enclosing variables is allowed. Use `.each`
for imperative per-item work (`print`, mutation, and similar).

In this prototype the predicate expression is still evaluated eagerly left to
right, so impure helpers may appear to “work.” Such programs are
**non-conforming**: future compilers may assume purity, reorder or fuse
`.where` pipelines, or reject impure predicates. Do not depend on side effects
inside `.where`.

#### `.sum`

```vol
items.sum
```

Semantics:

1. Require an array (`R019`).
2. Start from integer `0`.
3. Add each item left to right with the normal `+` rules.
4. Empty array yields `0`.
5. Non-numeric items fail with `R013`.

`.where` / `.sum` are not parallel and not lazy in this prototype. Future
fusion or parallelization of `.where` depends on the purity rule above.

---

## 5. Statement grammar

```text
program     = statement*
statement   = block-stmt
            | export-stmt
            | function-decl
            | return-stmt
            | print-stmt
            | if-stmt
            | repeat-stmt
            | while-stmt
            | declaration
            | assignment
            | each-stmt
            | expression-stmt

block       = "{" statement* "}"
block-stmt  = block
declaration = identifier ":=" expression
assignment  = (identifier | index-expr) "=" expression
each-stmt   = expression ".each" identifier block
if-stmt     = "if" expression block ("elif" expression block)* ("else" block)?
repeat-stmt = "repeat" expression block
while-stmt  = "while" expression block
print-stmt  = "print" expression
return-stmt = "return" expression
function-decl = "fn" identifier "(" params? ")" block
export-stmt = "export" identifier ("," identifier)*
params      = identifier ("," identifier)*
```

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
name = expression
```

- `:=` declares a new **mutable** binding in the current scope.
- Redeclaring the same name in the same scope is an error (`S001` / `R001`).
- `=` assigns to an existing variable or array index.
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

`.each` and `.where` express different intents. Both are valid: use `.where` for
pure filters (and `.sum` for reduction); use `.each` when the body has side
effects. See §4.4.

### 5.7 `print`

```vol
print expression
```

Writes the display form of the value followed by a newline.

Display rules:

- integers, floats, Booleans, strings: ordinary text
- arrays: `[a, b, c]` with recursive display
- no quotes are added around strings in display

### 5.8 Functions

```vol
fn add(a, b) {
    return a + b
}
```

Rules:

- Functions are declared with `fn`.
- Parameters are local names in the function body.
- Module-level functions are visible throughout the module, including before their
  declaration text.
- Nested functions are installed when execution reaches their declaration.
- Functions capture their enclosing environment (closures).
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

Built-in names cannot be redeclared at module scope.

---

## 6. Name resolution

Before execution, VOL resolves names.

Static errors:

| Code | Meaning |
| --- | --- |
| `S001` | duplicate declaration in the same scope |
| `S002` | use of an undefined name |
| `S003` | wrong argument count for a known function/builtin/`.where` |

Additional resolution rules:

- Module functions are declared before body checking so forward calls work.
- Ordinary locals are order-sensitive: use after declare in that scope.
- `.where` conditions see `_` as a declared name.
- `.each` item names belong to the loop scope.
- Function parameters belong to the function scope.

---

## 7. Evaluation order

Unless an operator short-circuits:

1. Evaluate operands left to right.
2. For calls: evaluate callee, then arguments left to right, then invoke.
3. For arrays: evaluate elements left to right.
4. For statements: execute in source order.

Short-circuiting applies to `and`, `or`, and `? :` (only the taken branch runs).

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

There is no `try` / `catch` and no error values yet.

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
| `E103` | missing `.each` item name |
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

### 8.3 Runtime codes

| Code | Meaning |
| --- | --- |
| `R001` | duplicate variable in scope |
| `R002` | unknown variable |
| `R003` | index into non-array |
| `R004` | non-Boolean condition |
| `R005` | invalid `repeat` count |
| `R006` | `.each` on non-array |
| `R007` | unknown property |
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
| `R019` | `.sum` on non-array |
| `R020` | `.where` wrong arity |
| `R021` | `.where` on non-array |
| `R022` | `.where` condition not Boolean |
| `R023` | `input` prompt not string |
| `R024` | input read failure |
| `R025` | `assert` condition not Boolean |
| `R026` | `assert` message not string |
| `R027` | assertion failed |
| `R028` | integer overflow |
| `R029` | expected a value, got `nothing` |
| `R999` | internal unsupported expression |

---

## 9. Conformance examples

These programs must keep working. They are also covered by `examples/` and tests.

### 9.1 Precedence

```vol
print 1 + 2 * 3 // 7
print true or false and false // true
print not false == true // true
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
print large.sum // 28
```

### 9.4 Functions

```vol
fn square(n) {
    return n * n
}

print square(6) // 36
```

---

## 10. Explicitly out of scope

Do not treat these as specified just because vision docs mention them:

- static types and typed signatures
- `const` parameter syntax (function parameters stay mutable; Planned)
- ownership, borrowing, lifetimes
- structs, enums, methods
- generics
- packages/imports beyond local `export` metadata
- `parallel`, async, channels
- wrapping integer arithmetic and overflow build modes (default is trap; see §4.3)
- native memory layout / allocators
- C or LLVM backends

When one of these is designed, add a numbered section here **and** implement
tests before calling it Supported.

---

## 11. Open decisions (do not pretend these are settled)

Track progress in [`IDEAS.md`](IDEAS.md).

#### Decided

- **Mutability default (§5.2):** bindings are **mutable by default**. Opt-in
  immutability uses `const name := expression` with shallow semantics
  (implemented; `S030`/`R030` on reassignment).
- **Array assignment (§3.3):** assignment and argument passing **share** the
  array reference. Explicit clone is Planned; move/ownership are not implied.
- **Integer overflow (§4.3):** overflow **traps** (`R028` with a `fix`
  suggestion). Wrapping ops / build modes are Planned.
- **`.where` purity (§4.4):** predicates are **pure filtering**; side effects
  belong in `.each`. Impure predicates are non-conforming; purity checks Planned.
- **Missing return (§5.8):** fall-off yields `nothing`; discarding in a call
  statement is OK; assigning or using `nothing` as a value is `R029`.
- **String / array `.len` (§3.2, §4.4):** canonical short property; string `.len`
  is Unicode scalars; byte count is Planned.
- **`while` (§5.5):** permanent Supported vocabulary; no alternate spellings.
- **`if` (§5.3):** statement only, with `elif` / `else`. Value choice uses
  `? :` (§4.2.1). Expression-`if` is not part of the language.

There are no remaining open core decisions in this section. Track new questions
in [`IDEAS.md`](IDEAS.md). Implementers and LLMs must follow the concrete
behavior in sections 1–9.

---

## 12. Intended formatting (not implemented)

The formatter does not exist yet. Planned presentation rules:

- deterministic, idempotent output
- four-space indentation unless later changed
- opening braces on the declaration or control-flow line
- closing braces on their own line
- spaces around binary operators
- stable comment placement
- no mandatory semicolons

Details and commands live in [`IDEAS.md`](IDEAS.md).

---

## 13. Change process

When changing language behavior:

1. Update this specification (including the quick vocabulary tables when forms change).
2. Add or adjust tests in `internal/lang`.
3. Update `README.md` examples if users should learn the change.
4. Move completed plans out of `IDEAS.md`.
5. Run:

```text
go test ./...
go run ./cmd/vol run ./examples/first.vol
git diff --check
```

A behavior is not official until it is specified here and covered by a test.
