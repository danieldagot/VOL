# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `full_v0`
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Cards: `vol_v0.md`, `go_v0.md`

> Live API run. Recompute from committed JSONL if numbers are quoted.

## Summary

| Language | First-try % | Success @ K % | Median tokens (success) | Mean tokens (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Go | 100.0 | 100.0 | 632 | 611.2 | 5 |
| VOL | 100.0 | 100.0 | 707 | 681.4 | 5 |

## Per task-replicate

| Task | Lang | Rep | Success | First-try | Attempts | Tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 01-hello | vol | 1 | True | True | 1 | 619 | success |
| 02-arithmetic | vol | 1 | True | True | 1 | 655 | success |
| 03-conditions | vol | 1 | True | True | 1 | 711 | success |
| 06-where-sum | vol | 1 | True | True | 1 | 707 | success |
| 07-functions | vol | 1 | True | True | 1 | 715 | success |
| 01-hello | go | 1 | True | True | 1 | 539 | success |
| 02-arithmetic | go | 1 | True | True | 1 | 586 | success |
| 03-conditions | go | 1 | True | True | 1 | 632 | success |
| 06-where-sum | go | 1 | True | True | 1 | 657 | success |
| 07-functions | go | 1 | True | True | 1 | 642 | success |
