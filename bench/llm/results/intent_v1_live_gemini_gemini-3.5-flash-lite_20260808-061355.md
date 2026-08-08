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
- Baseline reuse: live langs=[vol]; frozen rows from `bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-051437.jsonl`

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 771 | 761.3 ± 113.3 | 15 |
| VOL | 100.0 | 100.0 | 676 | 685.3 ± 85.7 | 15 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 626.8 | 134.5 | 761.3 | 82.3% | 17.7% |
| VOL | 599.9 | 85.4 | 685.3 | 87.5% | 12.5% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -36.5% |
| Prompt tokens | -4.3% |
| Total workflow tokens (cold) | -10.0% |
| Abs. prompt delta / task-replicate | -26.9 |
| Abs. completion delta / task-replicate | -49.1 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 761.3 | 425.3 | -336.0 | 336 |
| VOL | 685.3 | 363.3 | -322.0 | 322 |

Warm VOL vs Python total: -14.6% (means 363.3 vs 425.3).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 693.2 | 575.7 | 117.6 | 9 |
| generation | VOL | 100.0 | 623.1 | 554.7 | 68.4 | 9 |
| modification | Python | 100.0 | 956.0 | 679.0 | 277.0 | 3 |
| modification | VOL | 100.0 | 813.7 | 624.0 | 189.7 | 3 |
| repair | Python | 100.0 | 771.0 | 728.0 | 43.0 | 3 |
| repair | VOL | 100.0 | 743.3 | 711.3 | 32.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | python | 1 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 772 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 771 | success |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 744 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 742 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 744 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 06-where-sum | generation | python | 1 | True | True | 1 | 614 | 547 | 67 | 278 | success |
| 06-where-sum | generation | python | 2 | True | True | 1 | 636 | 547 | 89 | 300 | success |
| 06-where-sum | generation | python | 3 | True | True | 1 | 636 | 547 | 89 | 300 | success |
| 08-strings-assert | repair | python | 1 | True | True | 1 | 770 | 727 | 43 | 434 | success |
| 08-strings-assert | repair | python | 2 | True | True | 1 | 772 | 729 | 43 | 436 | success |
| 08-strings-assert | repair | python | 3 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 11-leaderboard | modification | python | 1 | True | True | 1 | 944 | 679 | 265 | 608 | success |
| 11-leaderboard | modification | python | 2 | True | True | 1 | 950 | 679 | 271 | 614 | success |
| 11-leaderboard | modification | python | 3 | True | True | 1 | 974 | 679 | 295 | 638 | success |
| 14-pipeline-stats | generation | python | 1 | True | True | 1 | 775 | 605 | 170 | 439 | success |
| 14-pipeline-stats | generation | python | 2 | True | True | 1 | 775 | 605 | 170 | 439 | success |
| 14-pipeline-stats | generation | python | 3 | True | True | 1 | 787 | 605 | 182 | 451 | success |
| 16-map-filter | generation | python | 1 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 16-map-filter | generation | python | 2 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 16-map-filter | generation | python | 3 | True | True | 1 | 672 | 575 | 97 | 336 | success |
| 06-where-sum | generation | vol | 1 | True | True | 1 | 572 | 526 | 46 | 250 | success |
| 06-where-sum | generation | vol | 2 | True | True | 1 | 572 | 526 | 46 | 250 | success |
| 06-where-sum | generation | vol | 3 | True | True | 1 | 572 | 526 | 46 | 250 | success |
| 14-pipeline-stats | generation | vol | 1 | True | True | 1 | 676 | 584 | 92 | 354 | success |
| 14-pipeline-stats | generation | vol | 2 | True | True | 1 | 677 | 584 | 93 | 355 | success |
| 14-pipeline-stats | generation | vol | 3 | True | True | 1 | 670 | 584 | 86 | 348 | success |
| 16-map-filter | generation | vol | 1 | True | True | 1 | 623 | 554 | 69 | 301 | success |
| 16-map-filter | generation | vol | 2 | True | True | 1 | 623 | 554 | 69 | 301 | success |
| 16-map-filter | generation | vol | 3 | True | True | 1 | 623 | 554 | 69 | 301 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 744 | 712 | 32 | 422 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 742 | 710 | 32 | 420 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 744 | 712 | 32 | 422 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 806 | 624 | 182 | 484 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 810 | 624 | 186 | 488 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 825 | 624 | 201 | 503 | success |
