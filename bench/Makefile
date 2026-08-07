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
#   make llm-dry    — validate LLM harness runners with reference solutions (not a benchmark)
#   make llm-ollama — live smoke vs Go using local Ollama (default model qwen3:4b-instruct)

.PHONY: bench check count sync llm-dry llm-ollama

bench: check count

sync:
	uv sync

check: sync
	uv run python harness/check_stdout.py

count: sync
	uv run python harness/count_tokens.py

llm-dry: sync
	uv run python llm/harness/run_generate_repair.py --dry-run --suite smoke

# Override model: make llm-ollama MODEL=llama3.1:8b
MODEL ?= qwen3:4b-instruct
llm-ollama: sync
	uv run python llm/harness/run_generate_repair.py --provider ollama --model $(MODEL) --suite smoke --replicates 1
