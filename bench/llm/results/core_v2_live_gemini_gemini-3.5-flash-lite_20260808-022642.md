# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `core_v2`
- Protocol: v1.1
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Surface freeze: SF-0
- Cards: `python_v0.md`, `vol_v0.md` (SF-0)
- Card tokens (est. `cl100k_base`): Python=336, VOL=433

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 771 | 763.8 ± 146.2 | 15 |
| VOL | 100.0 | 100.0 | 874 | 839.9 ± 113.5 | 15 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 627.8 | 136.0 | 763.8 | 82.2% | 17.8% |
| VOL | 727.9 | 112.0 | 839.9 | 86.7% | 13.3% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -17.6% |
| Prompt tokens | +15.9% |
| Total workflow tokens (cold) | +10.0% |
| Abs. prompt delta / task-replicate | +100.1 |
| Abs. completion delta / task-replicate | -24.0 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 763.8 | 427.8 | -336.0 | 336 |
| VOL | 839.9 | 406.9 | -433.0 | 433 |

Warm VOL vs Python total: -4.9% (means 406.9 vs 427.8).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 695.3 | 577.3 | 118.0 | 9 |
| generation | VOL | 100.0 | 784.0 | 681.3 | 102.7 | 9 |
| modification | Python | 100.0 | 962.0 | 679.0 | 283.0 | 3 |
| modification | VOL | 100.0 | 973.0 | 759.0 | 214.0 | 3 |
| repair | Python | 100.0 | 771.0 | 728.0 | 43.0 | 3 |
| repair | VOL | 100.0 | 874.3 | 836.3 | 38.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 875 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 874 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 874 | success |
| 08-strings-assert | python | 1 | — | True | True | 1 | 771 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 771 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 771 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 05-arrays-each | generation | vol | 1 | True | True | 1 | 739 | 666 | 73 | 306 | success |
| 05-arrays-each | generation | vol | 2 | True | True | 1 | 739 | 666 | 73 | 306 | success |
| 05-arrays-each | generation | vol | 3 | True | True | 1 | 739 | 666 | 73 | 306 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 875 | 837 | 38 | 442 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 874 | 836 | 38 | 441 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 874 | 836 | 38 | 441 | success |
| 10-fibonacci | generation | vol | 1 | True | True | 1 | 678 | 635 | 43 | 245 | success |
| 10-fibonacci | generation | vol | 2 | True | True | 1 | 678 | 635 | 43 | 245 | success |
| 10-fibonacci | generation | vol | 3 | True | True | 1 | 678 | 635 | 43 | 245 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 978 | 759 | 219 | 545 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 978 | 759 | 219 | 545 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 963 | 759 | 204 | 530 | success |
| 13-temperatures | generation | vol | 1 | True | True | 1 | 935 | 743 | 192 | 502 | success |
| 13-temperatures | generation | vol | 2 | True | True | 1 | 937 | 743 | 194 | 504 | success |
| 13-temperatures | generation | vol | 3 | True | True | 1 | 933 | 743 | 190 | 500 | success |
| 05-arrays-each | generation | python | 1 | True | True | 1 | 630 | 562 | 68 | 294 | success |
| 05-arrays-each | generation | python | 2 | True | True | 1 | 630 | 562 | 68 | 294 | success |
| 05-arrays-each | generation | python | 3 | True | True | 1 | 630 | 562 | 68 | 294 | success |
| 08-strings-assert | repair | python | 1 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 08-strings-assert | repair | python | 2 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 08-strings-assert | repair | python | 3 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 10-fibonacci | generation | python | 1 | True | True | 1 | 575 | 531 | 44 | 239 | success |
| 10-fibonacci | generation | python | 2 | True | True | 1 | 575 | 531 | 44 | 239 | success |
| 10-fibonacci | generation | python | 3 | True | True | 1 | 575 | 531 | 44 | 239 | success |
| 11-leaderboard | modification | python | 1 | True | True | 1 | 964 | 679 | 285 | 628 | success |
| 11-leaderboard | modification | python | 2 | True | True | 1 | 958 | 679 | 279 | 622 | success |
| 11-leaderboard | modification | python | 3 | True | True | 1 | 964 | 679 | 285 | 628 | success |
| 13-temperatures | generation | python | 1 | True | True | 1 | 891 | 639 | 252 | 555 | success |
| 13-temperatures | generation | python | 2 | True | True | 1 | 861 | 639 | 222 | 525 | success |
| 13-temperatures | generation | python | 3 | True | True | 1 | 891 | 639 | 252 | 555 | success |
