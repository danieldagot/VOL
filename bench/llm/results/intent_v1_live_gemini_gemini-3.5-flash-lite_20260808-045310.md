# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `intent_v1`
- Protocol: v1.1
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Surface freeze: Python=SF-0, VOL=SF-1
- Cards: `python_v0.md` (SF-0), `vol_v1.md` (SF-1)
- Card tokens (est. `cl100k_base`): Python=336, VOL=436

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 745 | 756.8 ± 113.0 | 15 |
| VOL | 80.0 | 100.0 | 885 | 1057.5 ± 455.2 | 15 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 626.6 | 130.2 | 756.8 | 82.8% | 17.2% |
| VOL | 947.4 | 110.1 | 1057.5 | 89.6% | 10.4% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -15.4% |
| Prompt tokens | +51.2% |
| Total workflow tokens (cold) | +39.7% |
| Abs. prompt delta / task-replicate | +320.8 |
| Abs. completion delta / task-replicate | -20.1 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 756.8 | 420.8 | -336.0 | 336 |
| VOL | 1057.5 | 534.3 | -523.2 | 436 |

Warm VOL vs Python total: +27.0% (means 534.3 vs 420.8).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 684.7 | 575.7 | 109.0 | 9 |
| generation | VOL | 100.0 | 1141.2 | 1040.1 | 101.1 | 9 |
| modification | Python | 100.0 | 960.0 | 679.0 | 281.0 | 3 |
| modification | VOL | 100.0 | 979.3 | 770.0 | 209.3 | 3 |
| repair | Python | 100.0 | 770.0 | 727.0 | 43.0 | 3 |
| repair | VOL | 100.0 | 884.7 | 846.7 | 38.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 886 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 885 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 883 | success |
| 08-strings-assert | python | 1 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 771 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 769 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 06-where-sum | generation | vol | 1 | True | True | 1 | 718 | 662 | 56 | 282 | success |
| 06-where-sum | generation | vol | 2 | True | True | 1 | 718 | 662 | 56 | 282 | success |
| 06-where-sum | generation | vol | 3 | True | True | 1 | 718 | 662 | 56 | 282 | success |
| 14-pipeline-stats | generation | vol | 1 | True | False | 2 | 1961 | 1769 | 192 | 1089 | success |
| 14-pipeline-stats | generation | vol | 2 | True | False | 2 | 1941 | 1767 | 174 | 1069 | success |
| 14-pipeline-stats | generation | vol | 3 | True | False | 2 | 1944 | 1769 | 175 | 1072 | success |
| 16-map-filter | generation | vol | 1 | True | True | 1 | 757 | 690 | 67 | 321 | success |
| 16-map-filter | generation | vol | 2 | True | True | 1 | 757 | 690 | 67 | 321 | success |
| 16-map-filter | generation | vol | 3 | True | True | 1 | 757 | 690 | 67 | 321 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 886 | 848 | 38 | 450 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 885 | 847 | 38 | 449 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 883 | 845 | 38 | 447 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 983 | 770 | 213 | 547 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 972 | 770 | 202 | 536 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 983 | 770 | 213 | 547 | success |
| 06-where-sum | generation | python | 1 | True | True | 1 | 636 | 547 | 89 | 300 | success |
| 06-where-sum | generation | python | 2 | True | True | 1 | 637 | 547 | 90 | 301 | success |
| 06-where-sum | generation | python | 3 | True | True | 1 | 636 | 547 | 89 | 300 | success |
| 14-pipeline-stats | generation | python | 1 | True | True | 1 | 775 | 605 | 170 | 439 | success |
| 14-pipeline-stats | generation | python | 2 | True | True | 1 | 717 | 605 | 112 | 381 | success |
| 14-pipeline-stats | generation | python | 3 | True | True | 1 | 745 | 605 | 140 | 409 | success |
| 16-map-filter | generation | python | 1 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 16-map-filter | generation | python | 2 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 16-map-filter | generation | python | 3 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 08-strings-assert | repair | python | 1 | True | True | 1 | 770 | 727 | 43 | 434 | success |
| 08-strings-assert | repair | python | 2 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 08-strings-assert | repair | python | 3 | True | True | 1 | 769 | 726 | 43 | 433 | success |
| 11-leaderboard | modification | python | 1 | True | True | 1 | 958 | 679 | 279 | 622 | success |
| 11-leaderboard | modification | python | 2 | True | True | 1 | 964 | 679 | 285 | 628 | success |
| 11-leaderboard | modification | python | 3 | True | True | 1 | 958 | 679 | 279 | 622 | success |
