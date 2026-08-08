# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `intent_v1`
- Protocol: v1.1
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Surface freeze: Python=SF-0, VOL=SF-2
- Cards: `python_v0.md` (SF-0), `vol_v2.md` (SF-2)
- Card tokens (est. `cl100k_base`): Python=336, VOL=452

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 771 | 761.3 ± 113.3 | 15 |
| VOL | 93.3 | 100.0 | 830 | 908.7 ± 256.0 | 15 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 626.8 | 134.5 | 761.3 | 82.3% | 17.7% |
| VOL | 818.0 | 90.7 | 908.7 | 90.0% | 10.0% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -32.6% |
| Prompt tokens | +30.5% |
| Total workflow tokens (cold) | +19.4% |
| Abs. prompt delta / task-replicate | +191.2 |
| Abs. completion delta / task-replicate | -43.8 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 761.3 | 425.3 | -336.0 | 336 |
| VOL | 908.7 | 426.6 | -482.1 | 452 |

Warm VOL vs Python total: +0.3% (means 426.6 vs 425.3).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 693.2 | 575.7 | 117.6 | 9 |
| generation | VOL | 100.0 | 895.4 | 816.7 | 78.8 | 9 |
| modification | Python | 100.0 | 956.0 | 679.0 | 277.0 | 3 |
| modification | VOL | 100.0 | 956.3 | 777.0 | 179.3 | 3 |
| repair | Python | 100.0 | 771.0 | 728.0 | 43.0 | 3 |
| repair | VOL | 100.0 | 901.0 | 863.0 | 38.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 901 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 901 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 901 | success |
| 08-strings-assert | python | 1 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 772 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 771 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 06-where-sum | generation | vol | 1 | True | True | 1 | 738 | 679 | 59 | 286 | success |
| 06-where-sum | generation | vol | 2 | True | True | 1 | 730 | 679 | 51 | 278 | success |
| 06-where-sum | generation | vol | 3 | True | True | 1 | 738 | 679 | 59 | 286 | success |
| 14-pipeline-stats | generation | vol | 1 | True | True | 1 | 830 | 737 | 93 | 378 | success |
| 14-pipeline-stats | generation | vol | 2 | True | True | 1 | 825 | 737 | 88 | 373 | success |
| 14-pipeline-stats | generation | vol | 3 | True | True | 1 | 830 | 737 | 93 | 378 | success |
| 16-map-filter | generation | vol | 1 | True | True | 1 | 774 | 707 | 67 | 322 | success |
| 16-map-filter | generation | vol | 2 | True | True | 1 | 774 | 707 | 67 | 322 | success |
| 16-map-filter | generation | vol | 3 | True | False | 2 | 1820 | 1688 | 132 | 916 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 901 | 863 | 38 | 449 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 901 | 863 | 38 | 449 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 901 | 863 | 38 | 449 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 959 | 777 | 182 | 507 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 955 | 777 | 178 | 503 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 955 | 777 | 178 | 503 | success |
| 06-where-sum | generation | python | 1 | True | True | 1 | 614 | 547 | 67 | 278 | success |
| 06-where-sum | generation | python | 2 | True | True | 1 | 636 | 547 | 89 | 300 | success |
| 06-where-sum | generation | python | 3 | True | True | 1 | 636 | 547 | 89 | 300 | success |
| 14-pipeline-stats | generation | python | 1 | True | True | 1 | 775 | 605 | 170 | 439 | success |
| 14-pipeline-stats | generation | python | 2 | True | True | 1 | 775 | 605 | 170 | 439 | success |
| 14-pipeline-stats | generation | python | 3 | True | True | 1 | 787 | 605 | 182 | 451 | success |
| 16-map-filter | generation | python | 1 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 16-map-filter | generation | python | 2 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 16-map-filter | generation | python | 3 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 08-strings-assert | repair | python | 1 | True | True | 1 | 770 | 727 | 43 | 434 | success |
| 08-strings-assert | repair | python | 2 | True | True | 1 | 772 | 729 | 43 | 436 | success |
| 08-strings-assert | repair | python | 3 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 11-leaderboard | modification | python | 1 | True | True | 1 | 944 | 679 | 265 | 608 | success |
| 11-leaderboard | modification | python | 2 | True | True | 1 | 950 | 679 | 271 | 614 | success |
| 11-leaderboard | modification | python | 3 | True | True | 1 | 974 | 679 | 295 | 638 | success |
