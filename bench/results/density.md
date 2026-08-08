# VOL Source Token Density

Measures how many tokens equivalent VOL programs use relative to Python, Go, Rust, and Zig.
A ratio < 1.0 means VOL is denser (fewer tokens) under that tokenizer.

> **What this measures:** source token density of hand-written equivalent programs.
> **What this does not measure:** LLM task-success rate or generate/repair cost.
> **Tokenizer note:** ratios depend on the tokenizer; GPT, Claude, and other models
> tokenize differently. Numbers are reported per tokenizer.
> **Suite size:** 21 tasks (parity / labeled / compression / stdlib tiers).
> Prefer **median (compression)** when judging semantic-density ops;
> **median (labeled)** is sensitive to print/string glue;
> **median (parity)** is a control near-Python floor;
> **median (stdlib)** is SF-3 `@std` / library intent vs peer stdlibs.

## Token density — tokenizer: `cl100k_base`

| Task | Tier | VOL | Python | Go | Rust | Zig | VOL/Python | VOL/Go | VOL/Rust | VOL/Zig |
|------|------|-----|--------|----|------|-----|------------|--------|----------|---------|
| 01-hello | parity | 26 | 28 | 45 | 36 | 73 | 0.929 | 0.578 | 0.722 | 0.356 |
| 02-arithmetic | parity | 25 | 37 | 52 | 58 | 102 | 0.676 | 0.481 | 0.431 | 0.245 |
| 03-conditions | parity | 33 | 38 | 57 | 51 | 90 | 0.868 | 0.579 | 0.647 | 0.367 |
| 04-loops | parity | 39 | 39 | 64 | 55 | 102 | 1.000 | 0.609 | 0.709 | 0.382 |
| 05-arrays-each | parity | 48 | 55 | 84 | 69 | 110 | 0.873 | 0.571 | 0.696 | 0.436 |
| 06-where-sum | compression | 40 | 57 | 112 | 71 | 108 | 0.702 | 0.357 | 0.563 | 0.370 |
| 07-functions | parity | 34 | 41 | 63 | 59 | 78 | 0.829 | 0.540 | 0.576 | 0.436 |
| 08-strings-assert | parity | 27 | 38 | 56 | 57 | 86 | 0.711 | 0.482 | 0.474 | 0.314 |
| 09-grade-report | labeled | 97 | 181 | 232 | 198 | 296 | 0.536 | 0.418 | 0.490 | 0.328 |
| 10-fibonacci | parity | 30 | 32 | 53 | 59 | 89 | 0.938 | 0.566 | 0.508 | 0.337 |
| 11-leaderboard | labeled | 131 | 185 | 277 | 218 | 337 | 0.708 | 0.473 | 0.601 | 0.389 |
| 12-revenue | labeled | 76 | 128 | 169 | 150 | 222 | 0.594 | 0.450 | 0.507 | 0.342 |
| 13-temperatures | labeled | 139 | 178 | 228 | 177 | 267 | 0.781 | 0.610 | 0.785 | 0.521 |
| 14-pipeline-stats | compression | 67 | 81 | 152 | 153 | 225 | 0.827 | 0.441 | 0.438 | 0.298 |
| 15-band-counts | compression | 73 | 110 | 145 | 159 | 244 | 0.664 | 0.503 | 0.459 | 0.299 |
| 16-map-filter | compression | 58 | 68 | 118 | 119 | 169 | 0.853 | 0.492 | 0.487 | 0.343 |
| 17-strings-ops | stdlib | 49 | 47 | 65 | 58 | 236 | 1.043 | 0.754 | 0.845 | 0.208 |
| 18-path-parts | stdlib | 56 | 73 | 79 | 99 | 157 | 0.767 | 0.709 | 0.566 | 0.357 |
| 19-env-default | stdlib | 50 | 48 | 105 | 87 | 147 | 1.042 | 0.476 | 0.575 | 0.340 |
| 20-json-fields | stdlib | 41 | 42 | 106 | 217 | 205 | 0.976 | 0.387 | 0.189 | 0.200 |
| 21-process-echo | stdlib | 32 | 33 | 91 | 66 | 165 | 0.970 | 0.352 | 0.485 | 0.194 |
| **median (all)** | | | | | | | **0.829** | **0.492** | **0.563** | **0.342** |
| **median (compression, n=4)** | | | | | | | **0.764** | **0.466** | **0.473** | **0.321** |
| **median (labeled, n=4)** | | | | | | | **0.651** | **0.461** | **0.554** | **0.366** |
| **median (parity, n=8)** | | | | | | | **0.871** | **0.569** | **0.612** | **0.361** |
| **median (stdlib, n=5)** | | | | | | | **0.976** | **0.476** | **0.566** | **0.208** |

## Token density — tokenizer: `o200k_base`

| Task | Tier | VOL | Python | Go | Rust | Zig | VOL/Python | VOL/Go | VOL/Rust | VOL/Zig |
|------|------|-----|--------|----|------|-----|------------|--------|----------|---------|
| 01-hello | parity | 25 | 27 | 44 | 35 | 72 | 0.926 | 0.568 | 0.714 | 0.347 |
| 02-arithmetic | parity | 25 | 37 | 52 | 58 | 102 | 0.676 | 0.481 | 0.431 | 0.245 |
| 03-conditions | parity | 33 | 38 | 57 | 51 | 90 | 0.868 | 0.579 | 0.647 | 0.367 |
| 04-loops | parity | 38 | 38 | 63 | 54 | 101 | 1.000 | 0.603 | 0.704 | 0.376 |
| 05-arrays-each | parity | 48 | 55 | 84 | 69 | 110 | 0.873 | 0.571 | 0.696 | 0.436 |
| 06-where-sum | compression | 40 | 57 | 112 | 71 | 108 | 0.702 | 0.357 | 0.563 | 0.370 |
| 07-functions | parity | 34 | 41 | 63 | 59 | 78 | 0.829 | 0.540 | 0.576 | 0.436 |
| 08-strings-assert | parity | 27 | 38 | 56 | 58 | 86 | 0.711 | 0.482 | 0.466 | 0.314 |
| 09-grade-report | labeled | 96 | 180 | 225 | 197 | 294 | 0.533 | 0.427 | 0.487 | 0.327 |
| 10-fibonacci | parity | 30 | 32 | 53 | 59 | 89 | 0.938 | 0.566 | 0.508 | 0.337 |
| 11-leaderboard | labeled | 131 | 185 | 277 | 218 | 337 | 0.708 | 0.473 | 0.601 | 0.389 |
| 12-revenue | labeled | 76 | 128 | 169 | 150 | 222 | 0.594 | 0.450 | 0.507 | 0.342 |
| 13-temperatures | labeled | 139 | 178 | 228 | 177 | 267 | 0.781 | 0.610 | 0.785 | 0.521 |
| 14-pipeline-stats | compression | 67 | 81 | 152 | 153 | 225 | 0.827 | 0.441 | 0.438 | 0.298 |
| 15-band-counts | compression | 73 | 110 | 145 | 159 | 244 | 0.664 | 0.503 | 0.459 | 0.299 |
| 16-map-filter | compression | 58 | 68 | 118 | 119 | 169 | 0.853 | 0.492 | 0.487 | 0.343 |
| 17-strings-ops | stdlib | 49 | 47 | 66 | 58 | 239 | 1.043 | 0.742 | 0.845 | 0.205 |
| 18-path-parts | stdlib | 56 | 73 | 80 | 99 | 157 | 0.767 | 0.700 | 0.566 | 0.357 |
| 19-env-default | stdlib | 50 | 48 | 105 | 87 | 147 | 1.042 | 0.476 | 0.575 | 0.340 |
| 20-json-fields | stdlib | 42 | 42 | 108 | 216 | 206 | 1.000 | 0.389 | 0.194 | 0.204 |
| 21-process-echo | stdlib | 32 | 33 | 93 | 66 | 165 | 0.970 | 0.344 | 0.485 | 0.194 |
| **median (all)** | | | | | | | **0.829** | **0.492** | **0.563** | **0.342** |
| **median (compression, n=4)** | | | | | | | **0.764** | **0.466** | **0.473** | **0.321** |
| **median (labeled, n=4)** | | | | | | | **0.651** | **0.461** | **0.554** | **0.366** |
| **median (parity, n=8)** | | | | | | | **0.871** | **0.567** | **0.612** | **0.357** |
| **median (stdlib, n=5)** | | | | | | | **1.000** | **0.476** | **0.566** | **0.205** |

