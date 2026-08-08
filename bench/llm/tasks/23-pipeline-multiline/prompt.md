# Task 23-pipeline-multiline

Write a complete program in {{LANG}}.

## Goal

From the list `[3, 8, 1, 12, 5, 9]`, print **three lines**:

1. Count of values greater than 5
2. Sum of values greater than 5
3. Sum of (values greater than 5, each doubled)

VOL: prefer a collection pipeline with `.where` / `.map` / `.sum` / `.len` or
`.count(pred)`. Pipelines may span multiple lines before `.`.

## Expected stdout

```text
3
29
58
```

## Constraints

- Single source file; no network, files, arguments, or stdin
- Do not use zero-arg `.count()` for length in VOL — use `.len`
