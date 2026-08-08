# Task 20-json-fields

Write a complete program in {{LANG}}.

## Goal

Parse the JSON object `{"n":3,"name":"vol"}` with the language's JSON helpers
(VOL: `import "@std/json"` then `json.parse` / `json.dump`), then print **three lines**:

1. The integer field `n`
2. The string field `name`
3. Compact JSON for an object with only `n` set to `3` (no spaces): `{"n":3}`

On VOL, use double-quoted strings with escapes (not `'...'`), e.g.
`json.parse("{\"n\":3,\"name\":\"vol\"}")?`; build the dump object with
`dict { n: 3 }` and `print json.dump(...)?` (Result — unwrap with `?`).

## Expected stdout

```text
3
vol
{"n":3}
```

## Constraints

- Single source file; no network, files, arguments, or stdin
- Prefer JSON parse/dump helpers over printing hard-coded field lines alone
- VOL: prefer `dict { k: v }` when building the dump object (ambient `dict("k", v)` still works)
- VOL: unwrap Result from `json.parse` / `json.dump` with `?` (do not print `ok(...)`)
