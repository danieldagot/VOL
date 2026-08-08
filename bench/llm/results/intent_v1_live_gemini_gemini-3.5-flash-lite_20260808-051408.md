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
- Card tokens (est. `cl100k_base`): Python=336, VOL=421
- Baseline reuse: live langs=[vol]; frozen rows from `bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-051110.jsonl`

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 771 | 761.3 ± 113.3 | 15 |
| VOL | 100.0 | 100.0 | 796 | 808.3 ± 84.1 | 15 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 626.8 | 134.5 | 761.3 | 82.3% | 17.7% |
| VOL | 721.9 | 86.4 | 808.3 | 89.3% | 10.7% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -35.8% |
| Prompt tokens | +15.2% |
| Total workflow tokens (cold) | +6.2% |
| Abs. prompt delta / task-replicate | +95.1 |
| Abs. completion delta / task-replicate | -48.1 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 761.3 | 425.3 | -336.0 | 336 |
| VOL | 808.3 | 387.3 | -421.0 | 421 |

Warm VOL vs Python total: -8.9% (means 387.3 vs 425.3).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 693.2 | 575.7 | 117.6 | 9 |
| generation | VOL | 100.0 | 746.1 | 676.7 | 69.4 | 9 |
| modification | Python | 100.0 | 956.0 | 679.0 | 277.0 | 3 |
| modification | VOL | 100.0 | 931.7 | 746.0 | 185.7 | 3 |
| repair | Python | 100.0 | 771.0 | 728.0 | 43.0 | 3 |
| repair | VOL | 100.0 | 871.3 | 833.3 | 38.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | python | 1 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 772 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 771 | success |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 871 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 873 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 870 | success |

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
| 06-where-sum | generation | vol | 1 | True | True | 1 | 702 | 648 | 54 | 281 | success |
| 06-where-sum | generation | vol | 2 | True | True | 1 | 702 | 648 | 54 | 281 | success |
| 06-where-sum | generation | vol | 3 | True | True | 1 | 694 | 648 | 46 | 273 | success |
| 14-pipeline-stats | generation | vol | 1 | True | True | 1 | 788 | 706 | 82 | 367 | success |
| 14-pipeline-stats | generation | vol | 2 | True | True | 1 | 796 | 706 | 90 | 375 | success |
| 14-pipeline-stats | generation | vol | 3 | True | True | 1 | 798 | 706 | 92 | 377 | success |
| 16-map-filter | generation | vol | 1 | True | True | 1 | 745 | 676 | 69 | 324 | success |
| 16-map-filter | generation | vol | 2 | True | True | 1 | 745 | 676 | 69 | 324 | success |
| 16-map-filter | generation | vol | 3 | True | True | 1 | 745 | 676 | 69 | 324 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 871 | 833 | 38 | 450 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 873 | 835 | 38 | 452 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 870 | 832 | 38 | 449 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 928 | 746 | 182 | 507 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 943 | 746 | 197 | 522 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 924 | 746 | 178 | 503 | success |
