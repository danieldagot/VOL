# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `core_v2`
- Protocol: v1.1
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Cards: `vol_v0.md`, `go_v0.md`
- Card tokens (est. `cl100k_base`): Go=341, VOL=433

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Go | 93.3 | 100.0 | 821 | 934.1 ± 510.0 | 15 |
| VOL | 93.3 | 100.0 | 874 | 941.5 ± 418.4 | 15 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Go | 738.1 | 195.9 | 934.1 | 79.0% | 21.0% |
| VOL | 814.0 | 127.5 | 941.5 | 86.5% | 13.5% |

## VOL vs Go token deltas (cold means)

| Metric | VOL vs Go |
| --- | --- |
| Generated completion tokens | -34.9% |
| Prompt tokens | +10.3% |
| Total workflow tokens (cold) | +0.8% |
| Abs. prompt delta / task-replicate | +75.9 |
| Abs. completion delta / task-replicate | -68.5 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Go | 934.1 | 570.3 | -363.7 | 341 |
| VOL | 941.5 | 479.6 | -461.9 | 433 |

Warm VOL vs Go total: -15.9% (means 479.6 vs 570.3).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Go | 100.0 | 941.6 | 741.3 | 200.2 | 9 |
| generation | VOL | 100.0 | 955.6 | 825.0 | 130.6 | 9 |
| modification | Go | 100.0 | 1025.0 | 720.0 | 305.0 | 3 |
| modification | VOL | 100.0 | 966.7 | 759.0 | 207.7 | 3 |
| repair | Go | 100.0 | 820.7 | 746.7 | 74.0 | 3 |
| repair | VOL | 100.0 | 874.0 | 836.0 | 38.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 873 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 875 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 874 | success |
| 08-strings-assert | go | 1 | — | True | True | 1 | 821 | success |
| 08-strings-assert | go | 2 | — | True | True | 1 | 820 | success |
| 08-strings-assert | go | 3 | — | True | True | 1 | 821 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 05-arrays-each | generation | vol | 1 | True | True | 1 | 739 | 666 | 73 | 306 | success |
| 05-arrays-each | generation | vol | 2 | True | True | 1 | 739 | 666 | 73 | 306 | success |
| 05-arrays-each | generation | vol | 3 | True | True | 1 | 739 | 666 | 73 | 306 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 873 | 835 | 38 | 440 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 875 | 837 | 38 | 442 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 874 | 836 | 38 | 441 | success |
| 10-fibonacci | generation | vol | 1 | True | True | 1 | 678 | 635 | 43 | 245 | success |
| 10-fibonacci | generation | vol | 2 | True | True | 1 | 678 | 635 | 43 | 245 | success |
| 10-fibonacci | generation | vol | 3 | True | True | 1 | 678 | 635 | 43 | 245 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 961 | 759 | 202 | 528 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 978 | 759 | 219 | 545 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 961 | 759 | 202 | 528 | success |
| 13-temperatures | generation | vol | 1 | True | True | 1 | 931 | 743 | 188 | 498 | success |
| 13-temperatures | generation | vol | 2 | True | False | 2 | 2451 | 2036 | 415 | 1585 | success |
| 13-temperatures | generation | vol | 3 | True | True | 1 | 967 | 743 | 224 | 534 | success |
| 05-arrays-each | generation | go | 1 | True | True | 1 | 670 | 565 | 105 | 329 | success |
| 05-arrays-each | generation | go | 2 | True | True | 1 | 658 | 565 | 93 | 317 | success |
| 05-arrays-each | generation | go | 3 | True | True | 1 | 670 | 565 | 105 | 329 | success |
| 08-strings-assert | repair | go | 1 | True | True | 1 | 821 | 747 | 74 | 480 | success |
| 08-strings-assert | repair | go | 2 | True | True | 1 | 820 | 746 | 74 | 479 | success |
| 08-strings-assert | repair | go | 3 | True | True | 1 | 821 | 747 | 74 | 480 | success |
| 10-fibonacci | generation | go | 1 | True | True | 1 | 604 | 534 | 70 | 263 | success |
| 10-fibonacci | generation | go | 2 | True | True | 1 | 604 | 534 | 70 | 263 | success |
| 10-fibonacci | generation | go | 3 | True | True | 1 | 609 | 534 | 75 | 268 | success |
| 11-leaderboard | modification | go | 1 | True | True | 1 | 1019 | 720 | 299 | 678 | success |
| 11-leaderboard | modification | go | 2 | True | True | 1 | 1037 | 720 | 317 | 696 | success |
| 11-leaderboard | modification | go | 3 | True | True | 1 | 1019 | 720 | 299 | 678 | success |
| 13-temperatures | generation | go | 1 | True | False | 2 | 2750 | 2091 | 659 | 2068 | success |
| 13-temperatures | generation | go | 2 | True | True | 1 | 947 | 642 | 305 | 606 | success |
| 13-temperatures | generation | go | 3 | True | True | 1 | 962 | 642 | 320 | 621 | success |
