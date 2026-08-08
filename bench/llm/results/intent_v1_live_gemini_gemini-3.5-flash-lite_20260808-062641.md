# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `intent_v1`
- Protocol: v1.1
- Model: `gemini-3.5-flash-lite`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Surface freeze: Python=SF-0, VOL=SF-3
- Cards: `python_v0.md` (SF-0), `vol_v3.md` (SF-3)
- Card tokens (est. `cl100k_base`): Python=336, VOL=322

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 680 | 733.6 ± 103.8 | 21 |
| VOL | 71.4 | 85.7 | 710 | 1080.2 ± 701.8 | 21 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 620.5 | 113.1 | 733.6 | 84.6% | 15.4% |
| VOL | 982.7 | 97.5 | 1080.2 | 91.0% | 9.0% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -13.8% |
| Prompt tokens | +58.4% |
| Total workflow tokens (cold) | +47.2% |
| Abs. prompt delta / task-replicate | +362.2 |
| Abs. completion delta / task-replicate | -15.6 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 733.6 | 397.6 | -336.0 | 336 |
| VOL | 1080.2 | 620.2 | -460.0 | 322 |

Warm VOL vs Python total: +56.0% (means 620.2 vs 397.6).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 682.1 | 587.4 | 94.7 | 15 |
| generation | VOL | 80.0 | 1201.7 | 1108.7 | 92.9 | 15 |
| modification | Python | 100.0 | 954.7 | 679.0 | 275.7 | 3 |
| modification | VOL | 100.0 | 810.0 | 624.0 | 186.0 | 3 |
| repair | Python | 100.0 | 770.3 | 727.3 | 43.0 | 3 |
| repair | VOL | 100.0 | 743.3 | 711.3 | 32.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 742 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 744 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 744 | success |
| 08-strings-assert | python | 1 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 771 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 06-where-sum | generation | vol | 1 | True | True | 1 | 572 | 526 | 46 | 250 | success |
| 06-where-sum | generation | vol | 2 | True | True | 1 | 572 | 526 | 46 | 250 | success |
| 06-where-sum | generation | vol | 3 | True | True | 1 | 572 | 526 | 46 | 250 | success |
| 14-pipeline-stats | generation | vol | 1 | True | True | 1 | 677 | 584 | 93 | 355 | success |
| 14-pipeline-stats | generation | vol | 2 | True | True | 1 | 670 | 584 | 86 | 348 | success |
| 14-pipeline-stats | generation | vol | 3 | True | True | 1 | 677 | 584 | 93 | 355 | success |
| 16-map-filter | generation | vol | 1 | True | True | 1 | 623 | 554 | 69 | 301 | success |
| 16-map-filter | generation | vol | 2 | True | True | 1 | 623 | 554 | 69 | 301 | success |
| 16-map-filter | generation | vol | 3 | True | True | 1 | 623 | 554 | 69 | 301 | success |
| 17-strings-ops | generation | vol | 1 | True | False | 2 | 1492 | 1388 | 104 | 848 | success |
| 17-strings-ops | generation | vol | 2 | True | False | 2 | 1493 | 1389 | 104 | 849 | success |
| 17-strings-ops | generation | vol | 3 | True | False | 2 | 1492 | 1388 | 104 | 848 | success |
| 20-json-fields | generation | vol | 1 | False | False | 3 | 2590 | 2454 | 136 | 1624 | diag_error |
| 20-json-fields | generation | vol | 2 | False | False | 3 | 2761 | 2568 | 193 | 1795 | diag_error |
| 20-json-fields | generation | vol | 3 | False | False | 3 | 2588 | 2452 | 136 | 1622 | diag_error |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 742 | 710 | 32 | 420 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 744 | 712 | 32 | 422 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 744 | 712 | 32 | 422 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 810 | 624 | 186 | 488 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 810 | 624 | 186 | 488 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 810 | 624 | 186 | 488 | success |
| 06-where-sum | generation | python | 1 | True | True | 1 | 635 | 547 | 88 | 299 | success |
| 06-where-sum | generation | python | 2 | True | True | 1 | 613 | 547 | 66 | 277 | success |
| 06-where-sum | generation | python | 3 | True | True | 1 | 613 | 547 | 66 | 277 | success |
| 14-pipeline-stats | generation | python | 1 | True | True | 1 | 775 | 605 | 170 | 439 | success |
| 14-pipeline-stats | generation | python | 2 | True | True | 1 | 757 | 605 | 152 | 421 | success |
| 14-pipeline-stats | generation | python | 3 | True | True | 1 | 775 | 605 | 170 | 439 | success |
| 16-map-filter | generation | python | 1 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 16-map-filter | generation | python | 2 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 16-map-filter | generation | python | 3 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 17-strings-ops | generation | python | 1 | True | True | 1 | 680 | 595 | 85 | 344 | success |
| 17-strings-ops | generation | python | 2 | True | True | 1 | 680 | 595 | 85 | 344 | success |
| 17-strings-ops | generation | python | 3 | True | True | 1 | 680 | 595 | 85 | 344 | success |
| 20-json-fields | generation | python | 1 | True | True | 1 | 667 | 615 | 52 | 331 | success |
| 20-json-fields | generation | python | 2 | True | True | 1 | 673 | 615 | 58 | 337 | success |
| 20-json-fields | generation | python | 3 | True | True | 1 | 667 | 615 | 52 | 331 | success |
| 08-strings-assert | repair | python | 1 | True | True | 1 | 770 | 727 | 43 | 434 | success |
| 08-strings-assert | repair | python | 2 | True | True | 1 | 770 | 727 | 43 | 434 | success |
| 08-strings-assert | repair | python | 3 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 11-leaderboard | modification | python | 1 | True | True | 1 | 956 | 679 | 277 | 620 | success |
| 11-leaderboard | modification | python | 2 | True | True | 1 | 964 | 679 | 285 | 628 | success |
| 11-leaderboard | modification | python | 3 | True | True | 1 | 944 | 679 | 265 | 608 | success |
