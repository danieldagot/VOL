# Task 15-band-counts

Write a complete program in {{LANG}}.

## Goal

Analyze the values `22, 18, 25, 31, 29, 17, 24, 28, 20, 26`.

Print **five lines**, each a single integer (no labels):

1. Number of values
2. Integer average (truncating division)
3. Count of values ≥ `28`
4. Count of values ≥ `20` and `< 28`
5. Count of values `< 20`

Prefer collection count / sum operations over hand-rolled loops when the
language supports them.

## Expected stdout

```text
10
24
3
5
2
```

## Constraints

- Compute from the collection; do not hard-code the five results as literals alone
- Single source file; no network, files, arguments, or stdin
