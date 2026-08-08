# Task 17-strings-ops

Write a complete program in {{LANG}}.

## Goal

Using the language's standard string helpers (VOL: `import "@std/strings"` then
namespaced calls `strings.trim` / `strings.split` / `strings.join` /
`strings.has` / `strings.replace`), print **four lines**:

1. Trim surrounding spaces from `  vol  ` → `vol`
2. Split `a,b,c` on `,` and join with `-` → `a-b-c`
3. Whether `vocabulary` has substring `cab` → `true` or `false` (lowercase)
4. Replace every `-` in `a-a-a` with `+` → `a+a+a`

## Expected stdout

```text
vol
a-b-c
true
a+a+a
```

## Constraints

- Single source file; no network, files, arguments, or stdin
- Prefer std string helpers over hand-rolled character loops
