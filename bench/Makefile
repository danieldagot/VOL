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
#   make llm-dry       — validate LLM harness (smoke reference solutions)
#   make llm-dry-intent — validate intent_v1 reference solutions
#   make llm-ollama    — live 2-task smoke suite using local Ollama
#   make llm-core      — live core_v2 (includes parity toys; historical tables)
#   make llm-intent    — live intent_v1 (filter/map/count + @std + repair/mod)

.PHONY: bench check count sync llm-dry llm-dry-intent llm-ollama llm-core llm-intent

bench: check count

sync:
	uv sync

check: sync
	uv run python harness/check_stdout.py

count: sync
	uv run python harness/count_tokens.py

llm-dry: sync
	uv run python llm/harness/run_generate_repair.py --dry-run --suite smoke

llm-dry-intent: sync
	uv run python llm/harness/run_generate_repair.py --dry-run --suite intent --langs vol,python

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

llm-intent: sync
	uv run python llm/harness/run_generate_repair.py --provider ollama --model $(MODEL) --suite intent --langs $(LANGS) --replicates 3 --tasks "$(TASKS)" $(if $(REQUEST_TIMEOUT),--request-timeout $(REQUEST_TIMEOUT),)
