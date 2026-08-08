# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `core_v2`
- Protocol: v1.1
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Surface freeze: Python=SF-0, VOL=SF-1
- Cards: `python_v0.md` (SF-0), `vol_v1.md` (SF-1)
- Card tokens (est. `cl100k_base`): Python=336, VOL=423

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 771 | 762.9 ± 142.9 | 15 |
| VOL | 60.0 | 100.0 | 870 | 1395.3 ± 782.1 | 15 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 627.8 | 135.1 | 762.9 | 82.3% | 17.7% |
| VOL | 1203.3 | 191.9 | 1395.3 | 86.2% | 13.8% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | +42.1% |
| Prompt tokens | +91.7% |
| Total workflow tokens (cold) | +82.9% |
| Abs. prompt delta / task-replicate | +575.5 |
| Abs. completion delta / task-replicate | +56.9 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 762.9 | 426.9 | -336.0 | 336 |
| VOL | 1395.3 | 803.1 | -592.2 | 423 |

Warm VOL vs Python total: +88.1% (means 803.1 vs 426.9).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 696.7 | 577.3 | 119.3 | 9 |
| generation | VOL | 100.0 | 1238.3 | 1069.8 | 168.6 | 9 |
| modification | Python | 100.0 | 953.3 | 679.0 | 274.3 | 3 |
| modification | VOL | 100.0 | 2391.0 | 1975.0 | 416.0 | 3 |
| repair | Python | 100.0 | 771.0 | 728.0 | 43.0 | 3 |
| repair | VOL | 100.0 | 870.3 | 832.3 | 38.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 870 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 870 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 871 | success |
| 08-strings-assert | python | 1 | — | True | True | 1 | 771 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 772 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 05-arrays-each | generation | vol | 1 | True | True | 1 | 735 | 662 | 73 | 312 | success |
| 05-arrays-each | generation | vol | 2 | True | True | 1 | 735 | 662 | 73 | 312 | success |
| 05-arrays-each | generation | vol | 3 | True | True | 1 | 735 | 662 | 73 | 312 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 870 | 832 | 38 | 447 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 870 | 832 | 38 | 447 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 871 | 833 | 38 | 448 | success |
| 10-fibonacci | generation | vol | 1 | True | True | 1 | 674 | 631 | 43 | 251 | success |
| 10-fibonacci | generation | vol | 2 | True | True | 1 | 672 | 631 | 41 | 249 | success |
| 10-fibonacci | generation | vol | 3 | True | True | 1 | 674 | 631 | 43 | 251 | success |
| 11-leaderboard | modification | vol | 1 | True | False | 2 | 2371 | 1965 | 406 | 1525 | success |
| 11-leaderboard | modification | vol | 2 | True | False | 2 | 2431 | 1995 | 436 | 1585 | success |
| 11-leaderboard | modification | vol | 3 | True | False | 2 | 2371 | 1965 | 406 | 1525 | success |
| 13-temperatures | generation | vol | 1 | True | False | 2 | 2272 | 1899 | 373 | 1426 | success |
| 13-temperatures | generation | vol | 2 | True | False | 2 | 2272 | 1899 | 373 | 1426 | success |
| 13-temperatures | generation | vol | 3 | True | False | 2 | 2376 | 1951 | 425 | 1530 | success |
| 05-arrays-each | generation | python | 1 | True | True | 1 | 630 | 562 | 68 | 294 | success |
| 05-arrays-each | generation | python | 2 | True | True | 1 | 630 | 562 | 68 | 294 | success |
| 05-arrays-each | generation | python | 3 | True | True | 1 | 630 | 562 | 68 | 294 | success |
| 08-strings-assert | repair | python | 1 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 08-strings-assert | repair | python | 2 | True | True | 1 | 770 | 727 | 43 | 434 | success |
| 08-strings-assert | repair | python | 3 | True | True | 1 | 772 | 729 | 43 | 436 | success |
| 10-fibonacci | generation | python | 1 | True | True | 1 | 575 | 531 | 44 | 239 | success |
| 10-fibonacci | generation | python | 2 | True | True | 1 | 587 | 531 | 56 | 251 | success |
| 10-fibonacci | generation | python | 3 | True | True | 1 | 575 | 531 | 44 | 239 | success |
| 11-leaderboard | modification | python | 1 | True | True | 1 | 958 | 679 | 279 | 622 | success |
| 11-leaderboard | modification | python | 2 | True | True | 1 | 944 | 679 | 265 | 608 | success |
| 11-leaderboard | modification | python | 3 | True | True | 1 | 958 | 679 | 279 | 622 | success |
| 13-temperatures | generation | python | 1 | True | True | 1 | 891 | 639 | 252 | 555 | success |
| 13-temperatures | generation | python | 2 | True | True | 1 | 861 | 639 | 222 | 525 | success |
| 13-temperatures | generation | python | 3 | True | True | 1 | 891 | 639 | 252 | 555 | success |
