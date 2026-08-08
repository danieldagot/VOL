# Task 20-json-fields

Write a complete program in {{LANG}}.

## Goal

Parse the JSON object `{"n":3,"name":"vol"}` with the language's JSON helpers
(VOL: `import "@std/json"`), then print **three lines**:

1. The integer field `n`
2. The string field `name`
3. Compact JSON for an object with only `n` set to `3` (no spaces): `{"n":3}`

On VOL, build that object with `dict("n", 3)` (or equivalent pairs) and
`print dump(...)?` (dump returns Result — unwrap with `?` so stdout is bare JSON).

## Expected stdout

```text
3
vol
{"n":3}
```

## Constraints

- Single source file; no network, files, arguments, or stdin
- Prefer JSON parse/dump helpers over printing hard-coded field lines alone
- VOL: use `dict("k", v, …)` pairs when building the dump object (no `{k:v}` literals)
- VOL: unwrap Result from `parse` / `dump` with `?` (do not print `ok(...)`)
