# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `smoke_v0`
- Model: `llama3.1:8b`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Cards: `vol_v0.md`, `go_v0.md`

> Live API run. Recompute from committed JSONL if numbers are quoted.

## Summary

| Language | First-try % | Success @ K % | Median tokens (success) | Mean tokens (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Go | 100.0 | 100.0 | 602 | 578.8 | 5 |
| VOL | 60.0 | 60.0 | 659 | 1379.4 | 5 |

## Per task-replicate

| Task | Lang | Rep | Success | First-try | Attempts | Tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 01-hello | vol | 1 | False | False | 3 | 2255 | wrong_output |
| 02-arithmetic | vol | 1 | True | True | 1 | 609 | success |
| 03-conditions | vol | 1 | True | True | 1 | 659 | success |
| 06-where-sum | vol | 1 | False | False | 3 | 2713 | diag_error |
| 07-functions | vol | 1 | True | True | 1 | 661 | success |
| 01-hello | go | 1 | True | True | 1 | 518 | success |
| 02-arithmetic | go | 1 | True | True | 1 | 550 | success |
| 03-conditions | go | 1 | True | True | 1 | 602 | success |
| 06-where-sum | go | 1 | True | True | 1 | 613 | success |
| 07-functions | go | 1 | True | True | 1 | 611 | success |
