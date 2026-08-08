# VOL language card (v3.1 / SF-3.1 foundation)

`.vol` file(s). No `main`. Top-level runs in order. No `match`.

- `:=` declare; `=` assign
- `if`/`elif`/`else`; values: `cond ? a : b`
- `repeat n { }`, `while cond { }`, `items.each item { ... }`
- Ops: `+ - * /` (int `/` truncates); `== != < <= > >=`; `and` `or` `not`
- Strings: `"..."` only (no `'...'`); escape JSON as `"{\"n\":3}"`; `"n=" + 7` ok; `1 + "x"` error
- `[1, 2]`, `a[i]`, `.len` (not `.count()` / `.length`); `_` in `.where`/`.map`/`.count(pred)`
- `.sum()` is **0-arg** — `.where(_ > 5).sum()`, not `.sum(pred)`
- Pipelines may span lines before `.` — `a\n.where(_ > 0)\n.map(_ * 2)`
- `fn add(a, b) { return a + b }`; missing `return` → `nothing`
- `print expr` or `print "label:", value`; `assert(cond)`
- `dict { k: v }` or `dict("k", v, …)`; `d["k"]`; `.keys()`
- `import "@std/…"` binds module name: `json.parse`, `strings.trim` (not flat `parse`)
- Result/`?`: `v := json.parse(s)?`; `print json.dump(d)?`; `strings.has` not `contains`; `env.get(k) ?? d`

Print exactly the required lines. Prefer raw source (no markdown fences).
