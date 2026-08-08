# VOL language card (v2 / SF-2)

Single `.vol` file. No `main`. Top-level statements run in order.
Newlines separate statements (not `print a print b`). No `match`.

## Core

- `name := expr` declare; `name = expr` assign
- `if` / `elif` / `else`; values: `cond ? a : b`
- `repeat n { }`, `while cond { }`, `items.each item { ... }` (not C `for`)
- Ops: `+ - * /` (int `/` truncates); `== != < <= > >=`; `and` `or` `not`
- `"n=" + 7` ok; `1 + "x"` error

## Collections / I/O

- `[1, 2]`, `a[i]`, `a.len` / `a.count()` (not `.length`)
- `_` expr: `.where(_ > 5)`, `.map(_ * 2)`, `.count(_ > 5)`
- `.sum()` is **0-arg** — write `.where(_ > 5).sum()`, not `.sum(pred)`
- `fn add(a, b) { return a + b }`; missing `return` → `nothing`
- `print expr` or `print "label:", value`; `assert(cond)` / `assert(cond, "msg")`

## Output

Print exactly the required lines. Prefer raw source (no markdown fences).
