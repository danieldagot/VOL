# VOL Desired Syntax

This document describes the current desired direction for VOL syntax.

The syntax is provisional. Examples express design intent and may change as the lexer, parser, interpreter, and type system are implemented.

## Syntax Principles

- Use braces for every multi-line scope.
- Use newlines to separate statements.
- Do not require semicolons.
- Do not require parentheses around conditions.
- Prefer type inference when the compiler has enough information.
- Prefer one canonical way to express an operation.
- Remove boilerplate without removing structural clarity.
- Optimize syntax for humans, language models, formatters, and static analysis.

## Comments

```vol
// A line comment
```

Block-comment syntax has not been decided.

## Literals

```vol
42
3.14
true
false
"hello"
```

Initial literal types:

- integers
- floating-point numbers
- Booleans
- strings

The prototype accepts unsigned decimal source spelling for integers and decimal
floating-point values with digits on both sides of the dot. Unary `-` supplies
negative values. Integer literals must fit in a signed 64-bit value, and floating
literals must fit in a finite 64-bit floating-point value. An out-of-range literal
is rejected with diagnostic `E006`; it is never silently replaced or made infinite.

## Variables

Type inference:

```vol
name := "VOL"
count := 10
ready := true
```

Assignment:

```vol
count = count + 1
```

The exact mutability rules and explicit type syntax have not been decided.

## Operators

Arithmetic:

```vol
a + b
a - b
a * b
a / b
```

Comparison:

```vol
a == b
a != b
a < b
a <= b
a > b
a >= b
```

Boolean logic:

```vol
ready and active
ready or forced
not failed
```

Operator precedence must be formally defined and kept small and predictable.
Mixed integer and floating-point comparisons compare their numeric values, including
`==` and `!=`.

## Blocks and Scope

Braces define scope:

```vol
{
    value := 10
    print value
}
```

A standalone block is a statement. Its declarations shadow outer declarations and
leave scope at the closing brace.

Indentation is canonical formatting, but it does not define scope.

A newline is required between consecutive simple statements. A closing brace
self-delimits a block-bodied statement, so another statement may follow that brace
on the same line. Missing separators are rejected with diagnostic `E119`.

## Conditions

```vol
if score >= 50 {
    print "pass"
} else {
    print "fail"
}
```

No parentheses are required around the condition.

Whether `if` produces a value has not been decided.

## Repetition

Repeat a known number of times:

```vol
repeat 10 {
    print "hello"
}
```

Conditional loop, using provisional syntax:

```vol
while running {
    work()
}
```

The `while` keyword is not yet accepted as final VOL syntax.

## Arrays

```vol
numbers := [1, 2, 3]
```

Indexing and indexed assignment:

```vol
first := numbers[0]
numbers[1] = 10
```

Length:

```vol
count := numbers.length
```

Array bounds are checked by default. The behavior of optimized and explicitly unchecked builds remains to be designed.

## Collection Iteration

```vol
users.each user {
    print user
}
```

The compiler chooses the implementation of iteration. The programmer describes the intent.

VOL also accepts a concise semantic form for common collection operations:

```vol
total := numbers.where(_ > 5).sum
```

Inside `where`, `_` means the current collection item. The explicit loop and the
semantic form are both valid; programmers may choose the form that communicates
their intent most clearly.

## Functions

Functions begin with `fn`. Names are private unless listed in an `export`
declaration. An export may appear before or after a definition:

```vol
export add, double

fn add(a, b) {
    return a + b
}

fn double(value) {
    return value * 2
}

result := double(add(2, 3))
```

This is also valid:

```vol
fn start() {
    print "started"
}

export start
```

The future formatter will collect exports into one canonical declaration at the
top of the module. Parameter-type and return-type syntax are not final.

## Modules and Imports

VOL will use folders as module boundaries. Imports are project-root-relative or
resolved through aliases from `vol.config.json`:

```vol
import "services/users"
import "@db/models"
```

The exact symbol-selection syntax and module resolver are not implemented yet.

## Built-In Operations

Output remains a concise statement:

```vol
print "hello"
```

The initial interactive and checking operations use ordinary call syntax:

```vol
name := input("Name: ")
assert(name.length > 0, "Name cannot be empty.")
text := string(42)
```

`input()` reads one line from standard input and removes its line ending. Its
optional string argument is written as a prompt first. `assert` requires a
Boolean condition and accepts an optional string failure message. `string`
returns the display representation of one value.

The built-in `args` array contains command-line arguments passed after the source
file. For example, `vol run program.vol -- first second` provides
`["first", "second"]`. The separator is optional with the current CLI.

Built-in names cannot be redeclared at module scope.

## Name Resolution

Names are resolved before execution. Duplicate declarations in the same lexical
scope and references to undefined names are errors. Nested scopes may shadow
outer names. Function names are visible throughout their module, allowing a call
before the corresponding declaration. Function parameters and `.each` item names
belong to their function or iteration scope.

Calls to known VOL functions and built-ins are checked for the correct number of
arguments before execution.

## Parallel Intent

Provisional syntax:

```vol
parallel {
    process requests
}
```

Automatic parallelization must preserve correctness and predictable resource use. Its exact semantics have not been decided.

## Example Program

The first interpreter milestone should support a program similar to:

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

This example establishes the initial language foundation:

- literals
- variables
- arrays
- iteration
- conditions
- comparison
- arithmetic
- blocks
- assignment
- output

## Formatting

The formatter defines the canonical presentation of valid VOL code.

Required properties:

- deterministic output
- idempotent formatting
- four-space indentation unless later changed by specification
- opening braces on the declaration or control-flow line
- closing braces on their own line
- spaces around binary operators
- stable comment placement
- no mandatory semicolons

Canonical example:

```vol
if ready {
    start()
} else {
    wait()
}
```

## Undecided Syntax

The following must be designed before VOL has a stable grammar:

- explicit types
- immutable and mutable declarations
- typed function signatures
- structs and methods
- enums and tagged unions
- error propagation
- optional values
- modules and imports
- generics
- ownership constraints
- allocator constraints
- asynchronous operations
- pattern matching
- compile-time evaluation
- foreign-function declarations
