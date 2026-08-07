# VOL language card (v0)

Single `.vol` file. No `main`. Top-level statements run in order.

**Newlines separate statements.** Put each statement on its own line.
Do not write `print a print b` on one line — that is a parse error (`E119`).

## Bindings and values

- `name := expr` mutable declare; `name = expr` assign; `const name := expr` no rebind
- Dynamic values: integer (i64), float, bool, string, array
- Arrays share identity on assign/args; `.copy()` / `.deep_copy()` to clone

## Control

- `if cond { } elif cond { } else { }` — statement; Boolean conditions only
- Values: `cond ? a : b` (no expression-`if`)
- `repeat n { }`, `while cond { }`, `items.each item { }`

## Operators

- `+ - * /` (int `/` truncates; int overflow traps)
- `== != < <= > >=`
- Boolean words: `and` `or` `not` (not `&&` `||` `!`)
- `not` applies after comparisons: `not a == b` means `not (a == b)`

## Collections / strings

- `[1, 2]`, index `a[i]`, `a.len`
- `items.where(_ > 5)` eager filter; `items.sum()` fold `+` from `0`
- String `+`; `.len` Unicode scalars; `.byte_len` UTF-8 bytes; `string(v)`

## Functions / I/O

```vol
fn add(a, b) {
    return a + b
}
```

- Fall off without `return` → `nothing` (OK as call statement; error if assigned/printed)
- `print expr`, `assert(cond)` or `assert(cond, "msg")`
- No stdin for these tasks

## Output

Print exactly the required lines. Prefer raw source (no markdown fences).
