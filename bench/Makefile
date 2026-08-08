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
#   make llm-ollama — live 2-task smoke suite using local Ollama
#   make llm-core   — live 5-task core suite with 3 replicates

.PHONY: bench check count sync llm-dry llm-ollama llm-core

bench: check count

sync:
	uv sync

check: sync
	uv run python harness/check_stdout.py

count: sync
	uv run python harness/count_tokens.py

llm-dry: sync
	uv run python llm/harness/run_generate_repair.py --dry-run --suite smoke

# Optional overrides:
#   make llm-ollama MODEL=llama3.1:8b
#   make llm-ollama TASKS=06-where-sum REQUEST_TIMEOUT=30
MODEL ?= qwen3:4b-instruct
TASKS ?=
REQUEST_TIMEOUT ?=
llm-ollama: sync
	uv run python llm/harness/run_generate_repair.py --provider ollama --model $(MODEL) --suite smoke --replicates 1 --tasks "$(TASKS)" $(if $(REQUEST_TIMEOUT),--request-timeout $(REQUEST_TIMEOUT),)

# LANGS=vol skips Python when iterating VOL-only (reuse last published Python numbers).
LANGS ?= vol,python
llm-core: sync
	uv run python llm/harness/run_generate_repair.py --provider ollama --model $(MODEL) --suite core --langs $(LANGS) --replicates 3 --tasks "$(TASKS)" $(if $(REQUEST_TIMEOUT),--request-timeout $(REQUEST_TIMEOUT),)
