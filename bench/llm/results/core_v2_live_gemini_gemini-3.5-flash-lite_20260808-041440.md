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
- Card tokens (est. `cl100k_base`): Python=336, VOL=436
- Baseline reuse: live langs=[vol]; frozen rows from `bench/llm/results/core_v2_live_gemini_gemini-3.5-flash-lite_20260808-040028.jsonl`

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 100.0 | 100.0 | 771 | 762.9 ± 142.9 | 15 |
| VOL | 100.0 | 100.0 | 884 | 848.6 ± 110.5 | 15 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 627.8 | 135.1 | 762.9 | 82.3% | 17.7% |
| VOL | 738.7 | 109.9 | 848.6 | 87.1% | 12.9% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -18.7% |
| Prompt tokens | +17.7% |
| Total workflow tokens (cold) | +11.2% |
| Abs. prompt delta / task-replicate | +110.9 |
| Abs. completion delta / task-replicate | -25.2 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 762.9 | 426.9 | -336.0 | 336 |
| VOL | 848.6 | 412.6 | -436.0 | 436 |

Warm VOL vs Python total: -3.3% (means 412.6 vs 426.9).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 696.7 | 577.3 | 119.3 | 9 |
| generation | VOL | 100.0 | 795.9 | 692.3 | 103.6 | 9 |
| modification | Python | 100.0 | 953.3 | 679.0 | 274.3 | 3 |
| modification | VOL | 100.0 | 970.7 | 770.0 | 200.7 | 3 |
| repair | Python | 100.0 | 771.0 | 728.0 | 43.0 | 3 |
| repair | VOL | 100.0 | 884.7 | 846.7 | 38.0 | 3 |

## Diagnostic repair notes

- Repair tasks seed a failing starter **before** attempt 0.
- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).
- `First-try` on repair means the model fixed the seed in one diagnostic turn.

| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 08-strings-assert | python | 1 | — | True | True | 1 | 771 | success |
| 08-strings-assert | python | 2 | — | True | True | 1 | 770 | success |
| 08-strings-assert | python | 3 | — | True | True | 1 | 772 | success |
| 08-strings-assert | vol | 1 | R007 | True | True | 1 | 884 | success |
| 08-strings-assert | vol | 2 | R007 | True | True | 1 | 884 | success |
| 08-strings-assert | vol | 3 | R007 | True | True | 1 | 886 | success |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
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
| 05-arrays-each | generation | vol | 1 | True | True | 1 | 750 | 677 | 73 | 314 | success |
| 05-arrays-each | generation | vol | 2 | True | True | 1 | 750 | 677 | 73 | 314 | success |
| 05-arrays-each | generation | vol | 3 | True | True | 1 | 750 | 677 | 73 | 314 | success |
| 08-strings-assert | repair | vol | 1 | True | True | 1 | 884 | 846 | 38 | 448 | success |
| 08-strings-assert | repair | vol | 2 | True | True | 1 | 884 | 846 | 38 | 448 | success |
| 08-strings-assert | repair | vol | 3 | True | True | 1 | 886 | 848 | 38 | 450 | success |
| 10-fibonacci | generation | vol | 1 | True | True | 1 | 694 | 646 | 48 | 258 | success |
| 10-fibonacci | generation | vol | 2 | True | True | 1 | 688 | 646 | 42 | 252 | success |
| 10-fibonacci | generation | vol | 3 | True | True | 1 | 688 | 646 | 42 | 252 | success |
| 11-leaderboard | modification | vol | 1 | True | True | 1 | 972 | 770 | 202 | 536 | success |
| 11-leaderboard | modification | vol | 2 | True | True | 1 | 968 | 770 | 198 | 532 | success |
| 11-leaderboard | modification | vol | 3 | True | True | 1 | 972 | 770 | 202 | 536 | success |
| 13-temperatures | generation | vol | 1 | True | True | 1 | 939 | 754 | 185 | 503 | success |
| 13-temperatures | generation | vol | 2 | True | True | 1 | 965 | 754 | 211 | 529 | success |
| 13-temperatures | generation | vol | 3 | True | True | 1 | 939 | 754 | 185 | 503 | success |
