# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `core_v1`
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Cards: `vol_v0.md`, `go_v0.md`

> Live API run. Recompute from committed JSONL if numbers are quoted.

## Summary

| Language | First-try % | Success @ K % | Median tokens (success) | Mean ± SD tokens (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Go | 100.0 | 100.0 | 685 | 789.3 ± 170.7 | 15 |
| VOL | 100.0 | 100.0 | 739 | 805.5 ± 118.4 | 15 |

## By workflow kind

| Kind | Language | Success @ K % | Mean ± SD tokens | N |
| --- | --- | --- | --- | --- |
| generation | Go | 100.0 | 745.0 ± 155.6 | 9 |
| generation | VOL | 100.0 | 783.8 ± 109.3 | 9 |
| modification | Go | 100.0 | 1026.3 ± 10.4 | 3 |
| modification | VOL | 100.0 | 962.0 ± 1.4 | 3 |
| repair | Go | 100.0 | 685.0 ± 0.0 | 3 |
| repair | VOL | 100.0 | 714.0 ± 0.0 | 3 |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 05-arrays-each | generation | vol | 1 | True | True | 1 | 739 | success |
| 05-arrays-each | generation | vol | 2 | True | True | 1 | 739 | success |
| 05-arrays-each | generation | vol | 3 | True | True | 1 | 739 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 714 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 714 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 714 | success |
| 10-fibonacci | generation | vol | 1 | True | True | 1 | 678 | success |
| 10-fibonacci | generation | vol | 2 | True | True | 1 | 678 | success |
| 10-fibonacci | generation | vol | 3 | True | True | 1 | 678 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 960 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 963 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 963 | success |
| 13-temperatures | generation | vol | 1 | True | True | 1 | 933 | success |
| 13-temperatures | generation | vol | 2 | True | True | 1 | 933 | success |
| 13-temperatures | generation | vol | 3 | True | True | 1 | 937 | success |
| 05-arrays-each | generation | go | 1 | True | True | 1 | 670 | success |
| 05-arrays-each | generation | go | 2 | True | True | 1 | 664 | success |
| 05-arrays-each | generation | go | 3 | True | True | 1 | 670 | success |
| 08-strings-assert | repair | go | 1 | True | True | 1 | 685 | success |
| 08-strings-assert | repair | go | 2 | True | True | 1 | 685 | success |
| 08-strings-assert | repair | go | 3 | True | True | 1 | 685 | success |
| 10-fibonacci | generation | go | 1 | True | True | 1 | 604 | success |
| 10-fibonacci | generation | go | 2 | True | True | 1 | 604 | success |
| 10-fibonacci | generation | go | 3 | True | True | 1 | 609 | success |
| 11-leaderboard | modification | go | 1 | True | True | 1 | 1019 | success |
| 11-leaderboard | modification | go | 2 | True | True | 1 | 1019 | success |
| 11-leaderboard | modification | go | 3 | True | True | 1 | 1041 | success |
| 13-temperatures | generation | go | 1 | True | True | 1 | 941 | success |
| 13-temperatures | generation | go | 2 | True | True | 1 | 990 | success |
| 13-temperatures | generation | go | 3 | True | True | 1 | 953 | success |
