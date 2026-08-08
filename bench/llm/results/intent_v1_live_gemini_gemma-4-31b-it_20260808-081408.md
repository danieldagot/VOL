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
- Card tokens (est. `cl100k_base`): Python=336, VOL=403

> Live API run. Recompute from committed JSONL if numbers are quoted.

> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card
> re-sent every request). **Warm** subtracts the estimated card tokens from each
> prompt (amortized / cached-card accounting).

## Summary

| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |
| --- | --- | --- | --- | --- | --- |
| Python | 66.7 | 100.0 | 818 | 1647.3 ± 1199.1 | 9 |
| VOL | 100.0 | 100.0 | 729 | 758.0 ± 52.7 | 9 |

## Prompt vs completion (cold, all attempts)

| Language | Mean prompt | Mean completion | Mean cold total | Prompt share | Completion share |
| --- | --- | --- | --- | --- | --- |
| Python | 1521.3 | 126.0 | 1647.3 | 92.4% | 7.6% |
| VOL | 701.0 | 57.0 | 758.0 | 92.5% | 7.5% |

## VOL vs Python token deltas (cold means)

| Metric | VOL vs Python |
| --- | --- |
| Generated completion tokens | -54.8% |
| Prompt tokens | -53.9% |
| Total workflow tokens (cold) | -54.0% |
| Abs. prompt delta / task-replicate | -820.3 |
| Abs. completion delta / task-replicate | -69.0 |

## Cold vs warm (card amortized)

| Language | Mean cold | Mean warm | Warm − cold | Card est. |
| --- | --- | --- | --- | --- |
| Python | 1647.3 | 1199.3 | -448.0 | 336 |
| VOL | 758.0 | 355.0 | -403.0 | 403 |

Warm VOL vs Python total: -70.4% (means 355.0 vs 1199.3).

## By workflow kind

| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |
| --- | --- | --- | --- | --- | --- | --- |
| generation | Python | 100.0 | 1647.3 | 1521.3 | 126.0 | 9 |
| generation | VOL | 100.0 | 758.0 | 701.0 | 57.0 | 9 |

## Per task-replicate

| Task | Kind | Lang | Rep | Success | First-try | Attempts | Cold | Prompt | Completion | Warm | Last outcome |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 20-json-fields | generation | vol | 1 | True | True | 1 | 832 | 779 | 53 | 429 | success |
| 20-json-fields | generation | vol | 2 | True | True | 1 | 832 | 779 | 53 | 429 | success |
| 20-json-fields | generation | vol | 3 | True | True | 1 | 832 | 779 | 53 | 429 | success |
| 22-ns-join | generation | vol | 1 | True | True | 1 | 713 | 668 | 45 | 310 | success |
| 22-ns-join | generation | vol | 2 | True | True | 1 | 713 | 668 | 45 | 310 | success |
| 22-ns-join | generation | vol | 3 | True | True | 1 | 713 | 668 | 45 | 310 | success |
| 23-pipeline-multiline | generation | vol | 1 | True | True | 1 | 729 | 656 | 73 | 326 | success |
| 23-pipeline-multiline | generation | vol | 2 | True | True | 1 | 729 | 656 | 73 | 326 | success |
| 23-pipeline-multiline | generation | vol | 3 | True | True | 1 | 729 | 656 | 73 | 326 | success |
| 20-json-fields | generation | python | 1 | True | True | 1 | 781 | 710 | 71 | 445 | success |
| 20-json-fields | generation | python | 2 | True | True | 1 | 781 | 710 | 71 | 445 | success |
| 20-json-fields | generation | python | 3 | True | True | 1 | 781 | 710 | 71 | 445 | success |
| 22-ns-join | generation | python | 1 | True | False | 2 | 3343 | 3267 | 76 | 2671 | success |
| 22-ns-join | generation | python | 2 | True | False | 2 | 3343 | 3267 | 76 | 2671 | success |
| 22-ns-join | generation | python | 3 | True | False | 2 | 3343 | 3267 | 76 | 2671 | success |
| 23-pipeline-multiline | generation | python | 1 | True | True | 1 | 818 | 587 | 231 | 482 | success |
| 23-pipeline-multiline | generation | python | 2 | True | True | 1 | 818 | 587 | 231 | 482 | success |
| 23-pipeline-multiline | generation | python | 3 | True | True | 1 | 818 | 587 | 231 | 482 | success |
