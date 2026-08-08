# VOL Source Token Density

Measures how many tokens equivalent VOL programs use relative to Python, Go, Rust, and Zig.
A ratio < 1.0 means VOL is denser (fewer tokens) under that tokenizer.

> **What this measures:** source token density of hand-written equivalent programs.
> **What this does not measure:** LLM task-success rate or generate/repair cost.
> **Tokenizer note:** ratios depend on the tokenizer; GPT, Claude, and other models
> tokenize differently. Numbers are reported per tokenizer.
> **Suite size:** 16 tasks (parity / labeled / compression tiers).
> Prefer **median (compression)** when judging semantic-density ops;
> **median (labeled)** is sensitive to print/string glue;
> **median (parity)** is a control near-Python floor.

## Token density — tokenizer: `cl100k_base`

| Task | Tier | VOL | Python | Go | Rust | Zig | VOL/Python | VOL/Go | VOL/Rust | VOL/Zig |
|------|------|-----|--------|----|------|-----|------------|--------|----------|---------|
| 01-hello | parity | 28 | 28 | 45 | 36 | 73 | 1.000 | 0.622 | 0.778 | 0.384 |
| 02-arithmetic | parity | 25 | 37 | 52 | 58 | 102 | 0.676 | 0.481 | 0.431 | 0.245 |
| 03-conditions | parity | 33 | 38 | 57 | 51 | 90 | 0.868 | 0.579 | 0.647 | 0.367 |
| 04-loops | parity | 39 | 39 | 64 | 55 | 102 | 1.000 | 0.609 | 0.709 | 0.382 |
| 05-arrays-each | parity | 54 | 55 | 84 | 69 | 110 | 0.982 | 0.643 | 0.783 | 0.491 |
| 06-where-sum | compression | 43 | 57 | 112 | 71 | 108 | 0.754 | 0.384 | 0.606 | 0.398 |
| 07-functions | parity | 36 | 41 | 63 | 59 | 78 | 0.878 | 0.571 | 0.610 | 0.462 |
| 08-strings-assert | parity | 27 | 38 | 56 | 57 | 86 | 0.711 | 0.482 | 0.474 | 0.314 |
| 09-grade-report | labeled | 112 | 181 | 232 | 198 | 296 | 0.619 | 0.483 | 0.566 | 0.378 |
| 10-fibonacci | parity | 30 | 32 | 53 | 59 | 89 | 0.938 | 0.566 | 0.508 | 0.337 |
| 11-leaderboard | labeled | 144 | 185 | 277 | 218 | 337 | 0.778 | 0.520 | 0.661 | 0.427 |
| 12-revenue | labeled | 91 | 128 | 169 | 150 | 222 | 0.711 | 0.538 | 0.607 | 0.410 |
| 13-temperatures | labeled | 156 | 178 | 228 | 177 | 267 | 0.876 | 0.684 | 0.881 | 0.584 |
| 14-pipeline-stats | compression | 74 | 81 | 152 | 153 | 225 | 0.914 | 0.487 | 0.484 | 0.329 |
| 15-band-counts | compression | 73 | 110 | 145 | 159 | 244 | 0.664 | 0.503 | 0.459 | 0.299 |
| 16-map-filter | compression | 61 | 68 | 118 | 119 | 169 | 0.897 | 0.517 | 0.513 | 0.361 |
| **median (all)** | | | | | | | **0.872** | **0.529** | **0.606** | **0.380** |
| **median (compression, n=4)** | | | | | | | **0.826** | **0.495** | **0.498** | **0.345** |
| **median (labeled, n=4)** | | | | | | | **0.745** | **0.529** | **0.634** | **0.419** |
| **median (parity, n=8)** | | | | | | | **0.908** | **0.575** | **0.629** | **0.375** |

## Token density — tokenizer: `o200k_base`

| Task | Tier | VOL | Python | Go | Rust | Zig | VOL/Python | VOL/Go | VOL/Rust | VOL/Zig |
|------|------|-----|--------|----|------|-----|------------|--------|----------|---------|
| 01-hello | parity | 27 | 27 | 44 | 35 | 72 | 1.000 | 0.614 | 0.771 | 0.375 |
| 02-arithmetic | parity | 25 | 37 | 52 | 58 | 102 | 0.676 | 0.481 | 0.431 | 0.245 |
| 03-conditions | parity | 33 | 38 | 57 | 51 | 90 | 0.868 | 0.579 | 0.647 | 0.367 |
| 04-loops | parity | 38 | 38 | 63 | 54 | 101 | 1.000 | 0.603 | 0.704 | 0.376 |
| 05-arrays-each | parity | 54 | 55 | 84 | 69 | 110 | 0.982 | 0.643 | 0.783 | 0.491 |
| 06-where-sum | compression | 43 | 57 | 112 | 71 | 108 | 0.754 | 0.384 | 0.606 | 0.398 |
| 07-functions | parity | 36 | 41 | 63 | 59 | 78 | 0.878 | 0.571 | 0.610 | 0.462 |
| 08-strings-assert | parity | 27 | 38 | 56 | 58 | 86 | 0.711 | 0.482 | 0.466 | 0.314 |
| 09-grade-report | labeled | 111 | 180 | 225 | 197 | 294 | 0.617 | 0.493 | 0.563 | 0.378 |
| 10-fibonacci | parity | 30 | 32 | 53 | 59 | 89 | 0.938 | 0.566 | 0.508 | 0.337 |
| 11-leaderboard | labeled | 144 | 185 | 277 | 218 | 337 | 0.778 | 0.520 | 0.661 | 0.427 |
| 12-revenue | labeled | 91 | 128 | 169 | 150 | 222 | 0.711 | 0.538 | 0.607 | 0.410 |
| 13-temperatures | labeled | 156 | 178 | 228 | 177 | 267 | 0.876 | 0.684 | 0.881 | 0.584 |
| 14-pipeline-stats | compression | 74 | 81 | 152 | 153 | 225 | 0.914 | 0.487 | 0.484 | 0.329 |
| 15-band-counts | compression | 73 | 110 | 145 | 159 | 244 | 0.664 | 0.503 | 0.459 | 0.299 |
| 16-map-filter | compression | 61 | 68 | 118 | 119 | 169 | 0.897 | 0.517 | 0.513 | 0.361 |
| **median (all)** | | | | | | | **0.872** | **0.529** | **0.606** | **0.376** |
| **median (compression, n=4)** | | | | | | | **0.826** | **0.495** | **0.498** | **0.345** |
| **median (labeled, n=4)** | | | | | | | **0.745** | **0.529** | **0.634** | **0.419** |
| **median (parity, n=8)** | | | | | | | **0.908** | **0.575** | **0.629** | **0.371** |

