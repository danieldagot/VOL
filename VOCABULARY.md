# VOL Vocabulary

This is a quick guide to VOL's words and core syntax. The **familiar equivalent** column relates each form to concepts used in languages such as C, Go, Rust, and Python; it does not mean that VOL uses exactly the same implementation.

VOL is still experimental. **Supported** forms work in the current interpreter. **Provisional** forms work today, but their final spelling or behavior may change.

## Keywords

| Word | Meaning | Familiar equivalent | Status |
| --- | --- | --- | --- |
| `and` | True when both Boolean expressions are true. | `&&` | Supported |
| `else` | Runs an alternative block when an `if` condition is false. | `else` | Supported |
| `false` | The Boolean false value. | `false` | Supported |
| `fn` | Starts every function declaration. | `func`, `function`, or `fn` | Supported |
| `export` | Makes one or more module names public. | export list or `pub` | Supported |
| `if` | Runs a block only when its condition is true. | `if` | Supported |
| `not` | Reverses a Boolean value. | `!` | Supported |
| `or` | True when either Boolean expression is true. | `||` | Supported |
| `print` | Writes a value to standard output. | `printf`, `fmt.Println`, or `println!` | Supported |
| `repeat` | Runs a block a known number of times. | A counting `for` loop | Supported |
| `return` | Ends a function and gives a value back to its caller. | `return` | Supported |
| `true` | The Boolean true value. | `true` | Supported |
| `while` | Repeats a block while its condition is true. | `while` or a conditional `for` loop | Provisional |

## Collection and Value Words

These are contextual names rather than reserved keywords.

| Form | Meaning | Familiar equivalent | Status |
| --- | --- | --- | --- |
| `items.each item { ... }` | Runs the block once for every item in a collection. | `for item in items`, range loop, or iterator loop | Supported |
| `value.length` | Gets the number of elements in an array or characters in a string. | `len(value)`, `value.len()`, or `value.length` | Supported |
| `items.where(condition)` | Keeps items for which the condition is true; `_` is the current item. | `filter` | Supported |
| `items.sum` | Adds all numeric items. | `sum`, `reduce`, or iterator sum | Supported |

## Core Symbols

| Form | Meaning | Familiar equivalent | Status |
| --- | --- | --- | --- |
| `name := value` | Declares a variable and infers its type. | `let`, `var`, or `auto` declaration | Supported |
| `name = value` | Assigns a new value to an existing variable. | Assignment | Supported |
| `fn name(parameters) { ... }` | Declares a private function by default. | Unexported or private function | Supported |
| `export name, other` | Makes names public, regardless of declaration order. | Export list | Supported |
| `{ ... }` | Creates a block and lexical scope. | Block or scope | Supported |
| `[a, b, c]` | Creates an array. | Array literal | Supported |
| `items[index]` | Reads or replaces an array element at an index. | Array indexing | Supported |
| `// text` | Adds a comment that continues to the end of the line. | Line comment | Supported |

Arithmetic operators are `+`, `-`, `*`, and `/`. Comparison operators are `==`, `!=`, `<`, `<=`, `>`, and `>=`.

## Examples

Repeat a block three times:

```vol
repeat 3 {
    print "hello"
}
```

The familiar C-style idea is:

```c
for (int i = 0; i < 3; i++) {
    printf("hello\n");
}
```

Visit every value in a collection:

```vol
numbers.each number {
    print number
}
```

Choose between two blocks:

```vol
if score >= 50 {
    print "pass"
} else {
    print "fail"
}
```

Declare and call private and public functions:

```vol
fn double(value) {
    return value * 2
}

fn Double(value) {
    return value * 2
}

print double(4) // private function
print Double(4) // public function
```

Every function declaration starts with `fn`. Functions are private unless their
names appear in an `export` declaration. Capitalization has no visibility meaning.

## Planned Vocabulary

Some words shown in the project vision, such as `parallel`, are design ideas and are not accepted by the current interpreter. They should be added to this reference only when their syntax and meaning are defined clearly.
