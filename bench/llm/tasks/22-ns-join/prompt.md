# Task 22-ns-join

Write a complete program in {{LANG}}.

## Goal

Print **two lines** using path join and string join from the standard library:

1. Join path parts `a`, `b`, `c.txt` with the language's path joiner → `a/b/c.txt`
   (forward slashes)
2. Join the string list `["x", "y", "z"]` with separator `-` → `x-y-z`

VOL must `import "@std/path"` and `import "@std/strings"` and call
`path.join` / `strings.join` (namespaced — both modules in one file).

## Expected stdout

```text
a/b/c.txt
x-y-z
```

## Constraints

- Single source file; no network, files, arguments, or stdin
