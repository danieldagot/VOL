# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `smoke_v0`
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 0
- Dry-run: False
- Cards: `vol_v0.md`, `go_v0.md`

> Live API run. Recompute from committed JSONL if numbers are quoted.

## Summary

| Language | First-try % | Success @ K % | Median tokens (success) | Mean tokens (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| VOL | 100.0 | 100.0 | 619 | 619.0 | 1 |

## Per task-replicate

| Task | Lang | Rep | Success | First-try | Attempts | Tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 01-hello | vol | 1 | True | True | 1 | 619 | success |
