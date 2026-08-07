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

## Blocks and Scope

Braces define scope:

```vol
{
    value := 10
    print value
}
```

Indentation is canonical formatting, but it does not define scope.

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

## Functions

Provisional concise syntax:

```vol
add a, b {
    return a + b
}

result := add 2, 3
```

Function declaration, parameter-type, return-type, and call syntax are not final.

## Built-In Operations

Desired concise form:

```vol
print "hello"
assert ready
```

Parenthesized calls remain under consideration where they improve clarity:

```vol
print("hello")
assert(ready)
```

VOL should select one canonical form before the syntax is stabilized.

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
- function signatures
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
