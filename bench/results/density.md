# VOL Source Token Density

Measures how many tokens equivalent VOL programs use relative to Go, Rust, and Zig.
A ratio < 1.0 means VOL is denser (fewer tokens) under that tokenizer.

> **What this measures:** source token density of hand-written equivalent programs.
> **What this does not measure:** LLM task-success rate or generate/repair cost.
> **Tokenizer note:** ratios depend on the tokenizer; GPT, Claude, and other models
> tokenize differently. Numbers are reported per tokenizer.
> **Suite size:** 13 tasks, using only currently-implemented VOL features.

## Token density — tokenizer: `cl100k_base`

| Task | VOL | Go | Rust | Zig | VOL/Go | VOL/Rust | VOL/Zig |
|------|-----|----|------|-----|--------|----------|---------|
| 01-hello | 28 | 45 | 36 | 73 | 0.622 | 0.778 | 0.384 |
| 02-arithmetic | 37 | 52 | 58 | 102 | 0.712 | 0.638 | 0.363 |
| 03-conditions | 40 | 57 | 51 | 90 | 0.702 | 0.784 | 0.444 |
| 04-loops | 39 | 64 | 55 | 102 | 0.609 | 0.709 | 0.382 |
| 05-arrays-each | 57 | 84 | 69 | 110 | 0.679 | 0.826 | 0.518 |
| 06-where-sum | 51 | 112 | 71 | 108 | 0.455 | 0.718 | 0.472 |
| 07-functions | 44 | 63 | 59 | 78 | 0.698 | 0.746 | 0.564 |
| 08-strings-assert | 36 | 56 | 57 | 86 | 0.643 | 0.632 | 0.419 |
| 09-grade-report | 149 | 232 | 198 | 296 | 0.642 | 0.753 | 0.503 |
| 10-fibonacci | 45 | 53 | 59 | 89 | 0.849 | 0.763 | 0.506 |
| 11-leaderboard | 175 | 277 | 218 | 337 | 0.632 | 0.803 | 0.519 |
| 12-revenue | 97 | 169 | 150 | 222 | 0.574 | 0.647 | 0.437 |
| 13-temperatures | 169 | 228 | 177 | 267 | 0.741 | 0.955 | 0.633 |
| **median** | | | | | **0.643** | **0.753** | **0.472** |

## Token density — tokenizer: `o200k_base`

| Task | VOL | Go | Rust | Zig | VOL/Go | VOL/Rust | VOL/Zig |
|------|-----|----|------|-----|--------|----------|---------|
| 01-hello | 27 | 44 | 35 | 72 | 0.614 | 0.771 | 0.375 |
| 02-arithmetic | 37 | 52 | 58 | 102 | 0.712 | 0.638 | 0.363 |
| 03-conditions | 40 | 57 | 51 | 90 | 0.702 | 0.784 | 0.444 |
| 04-loops | 38 | 63 | 54 | 101 | 0.603 | 0.704 | 0.376 |
| 05-arrays-each | 57 | 84 | 69 | 110 | 0.679 | 0.826 | 0.518 |
| 06-where-sum | 51 | 112 | 71 | 108 | 0.455 | 0.718 | 0.472 |
| 07-functions | 44 | 63 | 59 | 78 | 0.698 | 0.746 | 0.564 |
| 08-strings-assert | 36 | 56 | 58 | 86 | 0.643 | 0.621 | 0.419 |
| 09-grade-report | 147 | 225 | 197 | 294 | 0.653 | 0.746 | 0.500 |
| 10-fibonacci | 45 | 53 | 59 | 89 | 0.849 | 0.763 | 0.506 |
| 11-leaderboard | 175 | 277 | 218 | 337 | 0.632 | 0.803 | 0.519 |
| 12-revenue | 97 | 169 | 150 | 222 | 0.574 | 0.647 | 0.437 |
| 13-temperatures | 169 | 228 | 177 | 267 | 0.741 | 0.955 | 0.633 |
| **median** | | | | | **0.653** | **0.746** | **0.472** |

