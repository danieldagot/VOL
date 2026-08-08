# VOL language card (v3 / SF-3)

`.vol` file(s). No `main`. Top-level runs in order. Newlines separate statements. No `match`.

- `:=` declare; `=` assign
- `if`/`elif`/`else`; values: `cond ? a : b`
- `repeat n { }`, `while cond { }`, `items.each item { ... }`
- Ops: `+ - * /` (int `/` truncates); `== != < <= > >=`; `and` `or` `not`
- `"n=" + 7` ok; `1 + "x"` error
- `[1, 2]`, `a[i]`, `.len`/`.count()`; `_` in `.where`/`.map`/`.count`
- `.sum()` is **0-arg** — `.where(_ > 5).sum()`, not `.sum(pred)`
- `fn add(a, b) { return a + b }`; missing `return` → `nothing`
- `print expr` or `print "label:", value`; `assert(cond)`
- `dict()` / `dict("k", v, …)`; `d["k"]`; `.keys()` (no `{k:v}`)
- `import "@std/…"` → flat `trim`/`has`/`parse` (not `strings.trim` / `json.parse`)
- Result/`?`: `v := parse(s)?`; `print dump(d)?`; `has` not `contains`; `get(k) ?? d`

Print exactly the required lines. Prefer raw source (no markdown fences).
