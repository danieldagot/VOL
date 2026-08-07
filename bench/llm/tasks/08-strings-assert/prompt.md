# Task 08-strings-assert — diagnostic repair

Write a complete {{LANG}} program that preserves this intent: compute the length
of `hello`, duplicate the word, print both values, and fail an assertion-style
check unless the length is five.

The harness runs a broken starter **before** the first model turn and attaches
that starter's exit code and diagnostics. Repair using those diagnostics.
Return the complete corrected source, not a diff.

## Expected stdout

```text
5
hellohello
```

## Constraints

- Correct the failed program rather than replacing it with hard-coded output
- Retain an assertion or explicit panic-style invariant check
- Single source file; no network, files, arguments, or stdin
