# VOL LLM generate/repair results

- Date: 2026-08-08
- Suite: `intent_v1`
- Protocol: v1.1
- Model: `gemma-4-31b-it`
- Temperature: 0.0
- Max repairs (K): 2
- Dry-run: False
- Surface freeze: Python=SF-0, VOL=SF-3.1
- Cards: `python_v0.md` (SF-0), `vol_v3_1.md` (SF-3.1)
- Card tokens (est. `cl100k_base`): Python=336, VOL=380

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 0.0 | 66.7 | 6150 | 5493.2 ± 1326.8 | 9 |
| VOL | 66.7 | 66.7 | 694 | 2523.4 ± 2587.2 | 9 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 5199.9 | 293.3 | 5493.2 | 94.7% | 5.3% |
| VOL | 2425.3 | 98.1 | 2523.4 | 96.1% | 3.9% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -66.6% |
| Prompt tokens | -53.4% |
| Total workflow tokens (cold) | -54.1% |
| Abs. prompt delta / task-replicate | -2774.6 |
| Abs. completion delta / task-replicate | -195.2 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 5493.2 | 4522.6 | -970.7 | 336 |
| VOL | 2523.4 | 1890.1 | -633.3 | 380 |

Warm VOL vs Python total: -58.2% (means 1890.1 vs 4522.6).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 66.7 | 5493.2 | 5199.9 | 293.3 | 9 |
| generation | VOL | 66.7 | 2523.4 | 2425.3 | 98.1 | 9 |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 20-json-fields | generation | vol | 1 | False | False | 3 | 6183 | 6007 | 176 | 5043 | diag_error |
| 20-json-fields | generation | vol | 2 | False | False | 3 | 6184 | 6007 | 177 | 5044 | diag_error |
| 20-json-fields | generation | vol | 3 | False | False | 3 | 6180 | 6004 | 176 | 5040 | diag_error |
| 22-ns-join | generation | vol | 1 | True | True | 1 | 686 | 641 | 45 | 306 | success |
| 22-ns-join | generation | vol | 2 | True | True | 1 | 686 | 641 | 45 | 306 | success |
| 22-ns-join | generation | vol | 3 | True | True | 1 | 686 | 641 | 45 | 306 | success |
| 23-pipeline-multiline | generation | vol | 1 | True | True | 1 | 702 | 629 | 73 | 322 | success |
| 23-pipeline-multiline | generation | vol | 2 | True | True | 1 | 702 | 629 | 73 | 322 | success |
| 23-pipeline-multiline | generation | vol | 3 | True | True | 1 | 702 | 629 | 73 | 322 | success |
| 20-json-fields | generation | python | 1 | False | False | 3 | 4796 | 4583 | 213 | 3788 | diag_error |
| 20-json-fields | generation | python | 2 | True | False | 3 | 4795 | 4582 | 213 | 3787 | success |
| 20-json-fields | generation | python | 3 | True | False | 2 | 2519 | 2377 | 142 | 1847 | success |
| 22-ns-join | generation | python | 1 | True | False | 3 | 6892 | 6762 | 130 | 5884 | success |
| 22-ns-join | generation | python | 2 | True | False | 3 | 6894 | 6764 | 130 | 5886 | success |
| 22-ns-join | generation | python | 3 | True | False | 3 | 6892 | 6762 | 130 | 5884 | success |
| 23-pipeline-multiline | generation | python | 1 | False | False | 3 | 5390 | 4829 | 561 | 4382 | diag_error |
| 23-pipeline-multiline | generation | python | 2 | False | False | 3 | 5854 | 5294 | 560 | 4846 | diag_error |
| 23-pipeline-multiline | generation | python | 3 | True | False | 3 | 5407 | 4846 | 561 | 4399 | success |
