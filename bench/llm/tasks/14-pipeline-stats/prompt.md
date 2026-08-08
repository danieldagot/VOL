# Task 14-pipeline-stats

Write a complete program in {{LANG}}.

## Goal

Start from the integer list `3, 8, 1, 12, 5, 9, 4, 15, 2, 11`.

Compute and print **four lines**, each a single integer (no labels):

1. Count of values strictly greater than `5`
2. Sum of values strictly greater than `5`
3. Sum of those same values each multiplied by `2`
4. Count of values strictly less than `5`

Prefer collection filter / map / count / sum operations over hand-rolled
loops when the language supports them.

## Expected stdout

```text
5
55
110
4
```

## Constraints

- Compute from the list; do not hard-code the four results as literals alone
- Single source file; no network, files, arguments, or stdin
