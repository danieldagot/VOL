# Task 13-temperatures

Write a complete program in {{LANG}}.

## Goal

Analyze the temperatures `22, 18, 25, 31, 29, 17, 24, 28, 20, 26`.
Compute the number of measurements, the integer average, and counts for these
non-overlapping categories:

- hot: at least 28
- mild: at least 20 and below 28
- cold: below 20

Print the report exactly as shown and include invariant checks that the three
category counts cover all measurements and that the average is 24.

## Expected stdout

```text
Days measured: 10
Average: 24
Hot days (28+): 3
Mild days: 5
Cold days (<20): 2
```

## Constraints

- Compute all values from the collection; do not hard-code report values
- Retain an assertion or explicit panic-style invariant check
- Single source file; no network, files, arguments, or stdin
