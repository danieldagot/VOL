# VOL token-density benchmark
# Requires: uv (https://github.com/astral-sh/uv), go, rustc
# Optional:  zig (tasks skipped if not installed)
#
# First run — install Python deps:
#   uv sync
#
# Then:
#   make check   — verify all implementations produce the correct stdout
#   make count   — count tokens and write results/density*.csv + density.md
#   make bench   — run both in sequence

.PHONY: bench check count sync

bench: check count

sync:
	uv sync

check: sync
	uv run python harness/check_stdout.py

count: sync
	uv run python harness/count_tokens.py
