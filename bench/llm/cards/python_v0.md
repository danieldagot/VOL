# Python language card (v0)

Single `main.py` script. Top-level statements run. Stdlib only (no third-party packages).

## Bindings and values

- `name = expr` assign; no separate declare keyword
- Use `int`, `bool`, `str`, and lists such as `[1, 2]`
- Lists are references; mutating a shared list is visible through other names

## Control

- `if cond: … elif cond: … else: …` (indent body; Boolean conditions)
- `for i in range(n):`, `while cond:`, `for item in items:`

## Operators

- `+ - * //` (`//` truncates toward −∞ for ints; prefer non-negative ints here)
- `== != < <= > >=`
- Boolean: `and` `or` `not`

## Lists / strings

- Literal `[1, 2]`, index `a[i]`, `len(a)`, `a.append(x)`
- Filter/sum with loops or comprehensions (no special collection syntax required)
- String concat with `+` when both sides are strings
- Numbers in strings: `str(x)` or f-strings

## Functions / I/O

```python
def add(a, b):
    return a + b
```

- Call helpers from top level
- `print(...)` adds a newline
- For assertion-style checks, `assert cond` or `assert cond, "msg"` is acceptable
- No stdin for these tasks

## Output

Print exactly the required lines. Prefer raw source (no markdown fences).
