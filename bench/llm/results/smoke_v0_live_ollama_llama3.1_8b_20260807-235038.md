# VOL LLM generate/repair results

- Date: 2026-08-07
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
| VOL | 0.0 | 0.0 | — | 2105.0 | 1 |

## Per task-replicate

| Task | Lang | Rep | Success | First-try | Attempts | Tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 01-hello | vol | 1 | False | False | 3 | 2105 | diag_error |
