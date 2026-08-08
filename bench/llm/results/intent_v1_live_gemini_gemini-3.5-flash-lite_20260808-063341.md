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
- Card tokens (est. `cl100k_base`): Python=336, VOL=350
- Baseline reuse: live langs=[vol]; frozen rows from `bench/llm/results/intent_v1_live_gemini_gemini-3.5-flash-lite_20260808-062641.jsonl`

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 680 | 733.6 ± 103.8 | 21 |
| VOL | 100.0 | 100.0 | 718 | 715.7 ± 71.7 | 21 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 620.5 | 113.1 | 733.6 | 84.6% | 15.4% |
| VOL | 637.3 | 78.4 | 715.7 | 89.0% | 11.0% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -30.7% |
| Prompt tokens | +2.7% |
| Total workflow tokens (cold) | -2.4% |
| Abs. prompt delta / task-replicate | +16.8 |
| Abs. completion delta / task-replicate | -34.7 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 733.6 | 397.6 | -336.0 | 336 |
| VOL | 715.7 | 365.7 | -350.0 | 350 |

Warm VOL vs Python total: -8.0% (means 365.7 vs 397.6).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 682.1 | 587.4 | 94.7 | 15 |
| generation | VOL | 100.0 | 679.8 | 612.8 | 67.0 | 15 |
| modification | Python | 100.0 | 954.7 | 679.0 | 275.7 | 3 |
| modification | VOL | 100.0 | 837.0 | 655.0 | 182.0 | 3 |
| repair | Python | 100.0 | 770.3 | 727.3 | 43.0 | 3 |
| repair | VOL | 100.0 | 774.0 | 742.0 | 32.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | python | 1 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 771 | success |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 773 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 774 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 775 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 06-where-sum | generation | python | 1 | True | True | 1 | 635 | 547 | 88 | 299 | success |
| 06-where-sum | generation | python | 2 | True | True | 1 | 613 | 547 | 66 | 277 | success |
| 06-where-sum | generation | python | 3 | True | True | 1 | 613 | 547 | 66 | 277 | success |
| 08-strings-assert | repair | python | 1 | True | True | 1 | 770 | 727 | 43 | 434 | success |
| 08-strings-assert | repair | python | 2 | True | True | 1 | 770 | 727 | 43 | 434 | success |
| 08-strings-assert | repair | python | 3 | True | True | 1 | 771 | 728 | 43 | 435 | success |
| 11-leaderboard | modification | python | 1 | True | True | 1 | 956 | 679 | 277 | 620 | success |
| 11-leaderboard | modification | python | 2 | True | True | 1 | 964 | 679 | 285 | 628 | success |
| 11-leaderboard | modification | python | 3 | True | True | 1 | 944 | 679 | 265 | 608 | success |
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
| 06-where-sum | generation | vol | 1 | True | True | 1 | 603 | 557 | 46 | 253 | success |
| 06-where-sum | generation | vol | 2 | True | True | 1 | 603 | 557 | 46 | 253 | success |
| 06-where-sum | generation | vol | 3 | True | True | 1 | 603 | 557 | 46 | 253 | success |
| 14-pipeline-stats | generation | vol | 1 | True | True | 1 | 718 | 615 | 103 | 368 | success |
| 14-pipeline-stats | generation | vol | 2 | True | True | 1 | 714 | 615 | 99 | 364 | success |
| 14-pipeline-stats | generation | vol | 3 | True | True | 1 | 728 | 615 | 113 | 378 | success |
| 16-map-filter | generation | vol | 1 | True | True | 1 | 654 | 585 | 69 | 304 | success |
| 16-map-filter | generation | vol | 2 | True | True | 1 | 654 | 585 | 69 | 304 | success |
| 16-map-filter | generation | vol | 3 | True | True | 1 | 654 | 585 | 69 | 304 | success |
| 17-strings-ops | generation | vol | 1 | True | True | 1 | 686 | 634 | 52 | 336 | success |
| 17-strings-ops | generation | vol | 2 | True | True | 1 | 686 | 634 | 52 | 336 | success |
| 17-strings-ops | generation | vol | 3 | True | True | 1 | 686 | 634 | 52 | 336 | success |
| 20-json-fields | generation | vol | 1 | True | True | 1 | 736 | 673 | 63 | 386 | success |
| 20-json-fields | generation | vol | 2 | True | True | 1 | 736 | 673 | 63 | 386 | success |
| 20-json-fields | generation | vol | 3 | True | True | 1 | 736 | 673 | 63 | 386 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 773 | 741 | 32 | 423 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 774 | 742 | 32 | 424 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 775 | 743 | 32 | 425 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 837 | 655 | 182 | 487 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 837 | 655 | 182 | 487 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 837 | 655 | 182 | 487 | success |
