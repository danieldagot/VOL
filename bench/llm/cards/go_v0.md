# Go language card (v0)

Single `main.go` with `package main` and `func main()`. Stdlib only (`fmt` is enough).

## Bindings and values

- `name := expr` short declare inside functions; `name = expr` assign
- Use `int`, `bool`, `string`, and slices such as `[]int`
- Slices are references; mutating a shared slice is visible to other variables

## Control

- `if cond { } else if cond { } else { }`
- `for i := 0; i < n; i++ { }`, `for cond { }`, `for _, item := range items { }`

## Operators

- `+ - * /` (int `/` truncates toward zero)
- `== != < <= > >=`
- Boolean: `&&` `||` `!`

## Slices / strings

- Literal `[]int{1, 2}`, `len(s)`, `append(s, x)`, `range`
- Filter/sum by writing loops (no special collection syntax)
- String concat with `+` when both sides are strings
- Numbers in strings: `fmt.Sprint(x)`

## Functions / I/O

```go
func add(a, b int) int {
	return a + b
}
```

- Typed params/results; call helpers from `main`
- `fmt.Println(...)` adds a newline
- For assertion-style checks, `panic("msg")` on failure is acceptable
- No stdin for these tasks

## Output

Print exactly the required lines. Prefer raw source (no markdown fences).
