# VOL language card (v2 / SF-2)

Single `.vol` file. No `main`. Top-level statements run in order.
**Newlines separate statements** — not `print a print b` (`E119`). No `match`.

## Bindings and values

- `name := expr` / `a, b := x, y` declare; `name = expr` / `a, b = x, y` assign (RHS first)
- `const name := expr` no rebind
- Values: i64, float, bool, string, array; arrays share identity (`.copy()` / `.deep_copy()`)

## Control

- `if` / `elif` / `else`; values: `cond ? a : b`
- `repeat n { }`, `while cond { }`
- Loop: `items.each item { ... }` — not `.each(fn...)` or C-style `for`

## Operators

- `+ - * /` (int `/` truncates; overflow traps); `== != < <= > >=`
- Boolean infix: `and` `or` `not` — write `_ >= 20 and _ < 28`, not `&&`
- `"n=" + 7` ok; `1 + "x"` error

## Collections / strings

- `[1, 2]`, `a[i]`, `a.len` / `a.count()` (not `.length`)
- `_` expr: `items.where(_ > 5)`, `items.map(_ * 2)`, `items.count(_ > 5)`, `items.sum()`
- Prefer `.count(pred)` for filtered counts; `string(v)` if needed outside concat

## Functions / I/O

```vol
fn add(a, b) {
    return a + b
}
```

- No `return` → `nothing` (OK as call statement; error if assigned/printed)
- `print expr` or `print "label:", value` (space-join); `assert(cond)` / `assert(cond, "msg")`

## Output

Print exactly the required lines. Prefer raw source (no markdown fences).
