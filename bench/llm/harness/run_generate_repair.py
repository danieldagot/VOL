#!/usr/bin/env python3
"""
run_generate_repair.py — VOL LLM workflow benchmark harness (protocol v1.1).

See ../../LLM_BENCHMARK.md (repo root: LLM_BENCHMARK.md).

Modes:
  --provider ollama   Local Ollama OpenAI-compatible API (default for local work)
  --provider openai   OpenAI or any OpenAI-compatible cloud endpoint
  --provider gemini   Google Gemini via its OpenAI-compatible endpoint
  --dry-run           Skip the model; run frozen reference solutions from bench/tasks.
                      Use this to validate runners/summary wiring. NOT a benchmark result.

Protocol v1.1 notes:
  - Summaries report prompt vs completion tokens and cold vs warm (card-amortized) totals.
  - kind=repair with seed=starter: harness runs the broken starter first and attaches
    real exit code + diagnostics on the first model turn (diagnostic repair).

Environment:
  .env                            loaded from bench/.env or the repository root
  OLLAMA_HOST / OLLAMA_BASE_URL   default http://localhost:11434/v1
  OPENAI_API_KEY                  required for --provider openai
  OPENAI_BASE_URL                 optional cloud base URL
  OPENAI_MODEL / OLLAMA_MODEL     default model name override
  GEMINI_API_KEY                  required for --provider gemini
  GEMINI_BASE_URL / GEMINI_MODEL  optional Gemini overrides

Examples:
  uv run python llm/harness/run_generate_repair.py --dry-run
  uv run python llm/harness/run_generate_repair.py --dry-run --suite core --langs vol,python
  uv run python llm/harness/run_generate_repair.py --provider ollama --model qwen3:4b-instruct
  uv run python llm/harness/run_generate_repair.py --provider openai --model gpt-4.1-mini
  uv run python llm/harness/run_generate_repair.py --provider gemini --model gemini-3.6-flash
  uv run python llm/harness/run_generate_repair.py --provider gemini --suite core --langs vol,go
  # VOL-only re-run, reuse frozen Python rows from a prior JSONL:
  uv run python llm/harness/run_generate_repair.py --provider gemini --suite core --langs vol \\
    --baseline-jsonl llm/results/core_v2_live_gemini_…_python.jsonl
"""

from __future__ import annotations

import argparse
import json
import os
import re
import statistics
import subprocess
import sys
import tempfile
import time
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from datetime import date
from pathlib import Path

import tiktoken

# Card-size estimates use this encoding for cold/warm accounting. Absolute counts
# still come from each provider's usage fields; warm totals subtract this estimate.
CARD_TOKENIZER = "cl100k_base"

BENCH_ROOT = Path(__file__).resolve().parents[2]  # .../bench
VOL_REPO = BENCH_ROOT.parent
TASKS_DIR = BENCH_ROOT / "tasks"
LLM_ROOT = BENCH_ROOT / "llm"
CARDS_DIR = LLM_ROOT / "cards"
PROMPTS_DIR = LLM_ROOT / "tasks"
RESULTS_DIR = LLM_ROOT / "results"

SMOKE_TASKS = ["01-hello", "07-functions"]
# Historical / published tables — includes parity toys (fib, arrays). Keep frozen.
CORE_TASKS = [
    "05-arrays-each",
    "08-strings-assert",
    "10-fibonacci",
    "11-leaderboard",
    "13-temperatures",
]
# Intent-ops + real workflows only (no hello/fib parity). Prefer for language-use claims.
INTENT_TASKS = [
    "06-where-sum",
    "14-pipeline-stats",
    "16-map-filter",
    "08-strings-assert",
    "11-leaderboard",
]

# Cards bound to product surface freezes (SPEC.md §0). Next VOL bump is SF-2+.
# Default workflow baseline is Python (interpreted peer); Go remains optional.
LANG_META = {
    "vol": {
        "card": "vol_v1.md",
        "label": "VOL",
        "ext": ".vol",
        "prompt_lang": "VOL",
        "freeze": "SF-1",
    },
    "python": {
        "card": "python_v0.md",
        "label": "Python",
        "ext": ".py",
        "prompt_lang": "Python",
        "freeze": "SF-0",
    },
    "go": {
        "card": "go_v0.md",
        "label": "Go",
        "ext": ".go",
        "prompt_lang": "Go",
        "freeze": "SF-0",
    },
}

FENCE_RE = re.compile(r"```(?:[a-zA-Z0-9_+-]*)\n(.*?)```", re.DOTALL)
TIMEOUT_SEC = 5.0

DEFAULT_OLLAMA_BASE = "http://localhost:11434/v1"
DEFAULT_OLLAMA_MODEL = "qwen3:4b-instruct"
DEFAULT_OPENAI_BASE = "https://api.openai.com/v1"
DEFAULT_OPENAI_MODEL = "gpt-4.1-mini"
DEFAULT_GEMINI_BASE = "https://generativelanguage.googleapis.com/v1beta/openai"
DEFAULT_GEMINI_MODEL = "gemini-3.6-flash"
RETRYABLE_HTTP_CODES = {429, 500, 502, 503, 504}
MAX_API_RETRIES = 5


def load_dotenv(path: Path) -> None:
    """Load simple KEY=VALUE entries without overriding the process environment."""
    if not path.is_file():
        return
    for line_number, raw_line in enumerate(path.read_text(encoding="utf-8-sig").splitlines(), 1):
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        if line.startswith("export "):
            line = line[7:].lstrip()
        if "=" not in line:
            print(
                f"warning: skipping non KEY=VALUE line in {path}:{line_number}",
                file=sys.stderr,
            )
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip()
        if not re.fullmatch(r"[A-Za-z_][A-Za-z0-9_]*", key):
            print(
                f"warning: skipping invalid .env key in {path}:{line_number}: {key!r}",
                file=sys.stderr,
            )
            continue
        if len(value) >= 2 and value[0] == value[-1] and value[0] in "\"'":
            value = value[1:-1]
        os.environ.setdefault(key, value)


def load_environment() -> None:
    # A bench-local file is more specific; existing shell variables win over both.
    load_dotenv(BENCH_ROOT / ".env")
    load_dotenv(VOL_REPO / ".env")


@dataclass
class AttemptRecord:
    suite: str
    task: str
    task_kind: str
    language: str
    model: str
    temperature: float
    replicate: int
    attempt: int
    prompt_tokens: int
    completion_tokens: int
    exit_code: int | None
    outcome: str
    diagnostic_code: str | None = None
    stdout: str = ""
    stderr: str = ""
    source: str = ""
    dry_run: bool = False
    card_tokens_est: int = 0
    repair_seeded: bool = False


@dataclass
class TaskResult:
    task: str
    task_kind: str
    language: str
    replicate: int
    success: bool
    first_try: bool
    attempts: list[AttemptRecord] = field(default_factory=list)
    seed_diagnostic_code: str | None = None

    @property
    def total_tokens(self) -> int:
        """Cold total: measured prompt + completion across attempts."""
        return sum(a.prompt_tokens + a.completion_tokens for a in self.attempts)

    @property
    def prompt_tokens(self) -> int:
        return sum(a.prompt_tokens for a in self.attempts)

    @property
    def completion_tokens(self) -> int:
        return sum(a.completion_tokens for a in self.attempts)

    def warm_tokens(self, card_tokens: int) -> int:
        """Warm total: subtract estimated language-card tokens from every prompt."""
        return sum(
            max(0, a.prompt_tokens - card_tokens) + a.completion_tokens for a in self.attempts
        )


def die(msg: str, code: int = 1) -> None:
    print(msg, file=sys.stderr)
    raise SystemExit(code)


def get_vol_prefix() -> tuple[list[str], str | None]:
    vol_bin = os.environ.get("VOL_BIN")
    if vol_bin:
        return [vol_bin], None
    cmd_vol = VOL_REPO / "cmd" / "vol"
    if cmd_vol.exists():
        return ["go", "run", "./cmd/vol"], str(VOL_REPO)
    die("VOL_BIN not set and ../cmd/vol not found")


def run_cmd(
    cmd: list[str],
    *,
    cwd: str | None = None,
    timeout: float = TIMEOUT_SEC,
) -> tuple[str, str, int | None, bool]:
    """Returns stdout, stderr, returncode (None if timeout), timed_out."""
    try:
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            cwd=cwd,
            timeout=timeout,
        )
        return result.stdout, result.stderr, result.returncode, False
    except subprocess.TimeoutExpired as exc:
        out = exc.stdout if isinstance(exc.stdout, str) else (exc.stdout or b"").decode()
        err = exc.stderr if isinstance(exc.stderr, str) else (exc.stderr or b"").decode()
        return out, err, None, True


def extract_source(text: str) -> str | None:
    text = text.strip()
    if not text:
        return None
    blocks = [m.group(1).strip() for m in FENCE_RE.finditer(text)]
    if blocks:
        return max(blocks, key=len)
    # Raw source: reject obvious prose wrappers
    if text.startswith("{") and '"code"' in text[:80]:
        return None
    return text


def parse_vol_diagnostic_code(stderr: str) -> str | None:
    stderr = stderr.strip()
    if not stderr:
        return None
    try:
        data = json.loads(stderr.splitlines()[0] if stderr.startswith("{") else stderr)
        if isinstance(data, dict) and "code" in data:
            return str(data["code"])
    except json.JSONDecodeError:
        pass
    m = re.search(r"error\[([ERS]\d+)\]", stderr)
    return m.group(1) if m else None


def run_source(language: str, source: str, vol_prefix: list[str], vol_cwd: str | None) -> tuple[str, str, int | None, bool, str | None]:
    with tempfile.TemporaryDirectory(prefix="vol-llm-") as tmp:
        tmp_path = Path(tmp)
        if language == "vol":
            path = tmp_path / "main.vol"
            path.write_text(source, encoding="utf-8")
            cmd = vol_prefix + ["--json", "run", str(path)]
            stdout, stderr, rc, timed_out = run_cmd(cmd, cwd=vol_cwd)
            code = None if timed_out else parse_vol_diagnostic_code(stderr)
            return stdout, stderr, rc, timed_out, code
        if language == "go":
            path = tmp_path / "main.go"
            path.write_text(source, encoding="utf-8")
            stdout, stderr, rc, timed_out = run_cmd(["go", "run", str(path)])
            return stdout, stderr, rc, timed_out, None
        if language == "python":
            path = tmp_path / "main.py"
            path.write_text(source, encoding="utf-8")
            stdout, stderr, rc, timed_out = run_cmd([sys.executable, str(path)])
            return stdout, stderr, rc, timed_out, None
        die(f"unsupported language: {language}")


def load_card(language: str) -> str:
    path = CARDS_DIR / LANG_META[language]["card"]
    if not path.exists():
        die(f"missing language card: {path}")
    return path.read_text(encoding="utf-8")


def estimate_card_tokens(card: str) -> int:
    enc = tiktoken.get_encoding(CARD_TOKENIZER)
    return len(enc.encode(card))


def load_starter(task: str, language: str) -> str:
    path = PROMPTS_DIR / task / f"starter{LANG_META[language]['ext']}"
    if not path.exists():
        die(f"missing starter source: {path}")
    return path.read_text(encoding="utf-8")


def load_task_prompt(task: str, language: str) -> str:
    path = PROMPTS_DIR / task / "prompt.md"
    if not path.exists():
        die(f"missing task prompt: {path}")
    lang_name = LANG_META[language]["prompt_lang"]
    prompt = path.read_text(encoding="utf-8").replace("{{LANG}}", lang_name)
    starter_path = path.parent / f"starter{LANG_META[language]['ext']}"
    if "{{STARTER}}" in prompt:
        if not starter_path.exists():
            die(f"missing starter source: {starter_path}")
        prompt = prompt.replace("{{STARTER}}", starter_path.read_text(encoding="utf-8").rstrip())
    return prompt


def load_task_config(task: str) -> dict:
    path = PROMPTS_DIR / task / "task.json"
    if not path.exists():
        die(f"missing task metadata: {path}")
    try:
        config = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        die(f"invalid task metadata {path}: {exc}")
    if config.get("kind") not in {"generation", "repair", "modification"}:
        die(f"invalid task kind in {path}")
    if config["kind"] == "repair" and config.get("seed", "starter") != "starter":
        die(f"unsupported repair seed in {path}: {config.get('seed')!r}")
    return config


def load_expected(task: str) -> str:
    path = TASKS_DIR / task / "expected.txt"
    if not path.exists():
        die(f"missing expected.txt: {path}")
    return path.read_text(encoding="utf-8")


def load_reference(task: str, language: str) -> str:
    meta = LANG_META[language]
    path = TASKS_DIR / task / language / f"main{meta['ext']}"
    if not path.exists():
        die(f"missing reference solution: {path}")
    return path.read_text(encoding="utf-8")


def judge(
    stdout: str,
    rc: int | None,
    timed_out: bool,
    expected: str,
    source: str,
    source_checks: list[str],
) -> tuple[str, str | None]:
    if timed_out or rc is None:
        return "timeout", None
    if rc != 0:
        return "diag_error", None
    if stdout != expected:
        return "wrong_output", None
    for pattern in source_checks:
        if not re.search(pattern, source, re.MULTILINE):
            return "source_check_failed", f"Required source pattern was not found: {pattern}"
    return "success", None


SYSTEM_PROMPT = (
    "You are a careful programmer. Write a complete program that prints exactly "
    "the required stdout. Prefer raw source code with no markdown fences. "
    "Do not explain."
)


def format_failure_block(source: str, exit_code: int | None, stderr: str) -> str:
    return (
        "The previous program failed.\n\n"
        f"Exit code: {exit_code}\n\n"
        f"Diagnostics:\n{stderr.strip() or '(empty)'}\n\n"
        "Previous source:\n"
        f"```\n{source.rstrip()}\n```\n\n"
        "Return a full corrected program (not a diff)."
    )


def build_initial_messages(
    language: str,
    task: str,
    *,
    seed_failure: tuple[str, int | None, str] | None = None,
) -> list[dict[str, str]]:
    card = load_card(language)
    task_prompt = load_task_prompt(task, language)
    parts = [card.rstrip(), "---", task_prompt.rstrip()]
    if seed_failure is not None:
        source, exit_code, stderr = seed_failure
        parts.extend(["---", format_failure_block(source, exit_code, stderr)])
    else:
        parts.append("Return only the complete source file.")
    return [
        {"role": "system", "content": SYSTEM_PROMPT},
        {"role": "user", "content": "\n\n".join(parts)},
    ]


def build_repair_message(source: str, exit_code: int | None, stderr: str) -> dict[str, str]:
    return {"role": "user", "content": format_failure_block(source, exit_code, stderr)}


def normalize_openai_base(url: str) -> str:
    """Ensure base ends with /v1 for OpenAI-compatible chat completions."""
    url = url.rstrip("/")
    if url.endswith("/v1"):
        return url
    # Accept bare Ollama host like http://localhost:11434
    if url.rstrip("/").endswith(":11434") or url.endswith("localhost:11434"):
        return url + "/v1"
    return url


def retry_delay_seconds(exc: urllib.error.HTTPError, detail: str, retry: int) -> float:
    retry_after = exc.headers.get("Retry-After") if exc.headers else None
    if retry_after:
        try:
            return min(60.0, max(1.0, float(retry_after)))
        except ValueError:
            pass
    match = re.search(r'"retryDelay"\s*:\s*"([0-9.]+)s"', detail)
    if not match:
        match = re.search(r"retry in ([0-9.]+)s", detail, re.IGNORECASE)
    if match:
        return min(60.0, max(1.0, float(match.group(1)) + 1.0))
    return min(30.0, float(2**retry))


def chat_completion(
    *,
    base_url: str,
    api_key: str,
    model: str,
    temperature: float,
    messages: list[dict[str, str]],
    max_tokens: int,
    request_timeout: float,
) -> tuple[str, int, int]:
    url = normalize_openai_base(base_url) + "/chat/completions"
    payload = {
        "model": model,
        "temperature": temperature,
        "max_tokens": max_tokens,
        "messages": messages,
        "stream": False,
    }
    data = json.dumps(payload).encode("utf-8")
    headers = {
        "Content-Type": "application/json",
        # Ollama ignores the key but accepts the header; cloud requires it.
        "Authorization": f"Bearer {api_key or 'ollama'}",
    }
    body = None
    for retry in range(MAX_API_RETRIES + 1):
        req = urllib.request.Request(url, data=data, headers=headers, method="POST")
        try:
            with urllib.request.urlopen(req, timeout=request_timeout) as resp:
                body = json.loads(resp.read().decode("utf-8"))
            break
        except urllib.error.HTTPError as exc:
            detail = exc.read().decode("utf-8", errors="replace")
            if exc.code not in RETRYABLE_HTTP_CODES or retry >= MAX_API_RETRIES:
                die(f"API HTTP {exc.code}: {detail}")
            delay = retry_delay_seconds(exc, detail, retry)
            print(
                f"API HTTP {exc.code}; retrying in {delay:.1f}s "
                f"({retry + 1}/{MAX_API_RETRIES})",
                file=sys.stderr,
                flush=True,
            )
            time.sleep(delay)
        except urllib.error.URLError as exc:
            die(
                f"API request failed: {exc}\n"
                f"  URL: {url}\n"
                "  If using Ollama, is `ollama serve` running?"
            )

    if body is None:
        die("API request ended without a response")

    try:
        message = body["choices"][0]["message"]
        content = message.get("content")
    except (KeyError, IndexError, TypeError):
        die(f"unexpected API response: {body!r}")
    if content is None:
        die(f"empty model content in response: {body!r}")
    usage = body.get("usage") or {}
    # OpenAI: prompt_tokens / completion_tokens
    # Some Ollama builds also populate these via the OpenAI-compatible shim.
    prompt_tokens = int(
        usage.get("prompt_tokens")
        or usage.get("prompt_eval_count")
        or 0
    )
    completion_tokens = int(
        usage.get("completion_tokens")
        or usage.get("eval_count")
        or 0
    )
    return content, prompt_tokens, completion_tokens


def prepare_repair_seed(
    *,
    task: str,
    language: str,
    expected: str,
    source_checks: list[str],
    vol_prefix: list[str],
    vol_cwd: str | None,
) -> tuple[str, int | None, str, str | None]:
    """Run the broken starter; it must fail so diagnostics reach the model."""
    starter = load_starter(task, language)
    stdout, stderr, rc, timed_out, diag = run_source(language, starter, vol_prefix, vol_cwd)
    outcome, validation_error = judge(stdout, rc, timed_out, expected, starter, source_checks)
    if validation_error:
        stderr = validation_error
    if outcome == "success":
        die(
            f"repair starter for {task}/{language} unexpectedly succeeds; "
            "seeded diagnostic repair requires a failing starter"
        )
    if timed_out or rc is None:
        die(f"repair starter for {task}/{language} timed out before producing diagnostics")
    return starter, rc, stderr, diag


def run_task_replicate(
    *,
    suite: str,
    task: str,
    language: str,
    model: str,
    temperature: float,
    replicate: int,
    max_repairs: int,
    max_tokens: int,
    dry_run: bool,
    api_key: str | None,
    base_url: str,
    request_timeout: float,
    vol_prefix: list[str],
    vol_cwd: str | None,
    card_tokens_est: int,
) -> TaskResult:
    expected = load_expected(task)
    task_config = load_task_config(task)
    task_kind = task_config["kind"]
    source_checks = task_config.get("source_checks", {}).get(language, [])
    result = TaskResult(
        task=task,
        task_kind=task_kind,
        language=language,
        replicate=replicate,
        success=False,
        first_try=False,
    )

    seed_failure: tuple[str, int | None, str] | None = None
    repair_seeded = False
    if task_kind == "repair":
        starter, seed_rc, seed_stderr, seed_diag = prepare_repair_seed(
            task=task,
            language=language,
            expected=expected,
            source_checks=source_checks,
            vol_prefix=vol_prefix,
            vol_cwd=vol_cwd,
        )
        seed_failure = (starter, seed_rc, seed_stderr)
        result.seed_diagnostic_code = seed_diag
        repair_seeded = True

    messages = build_initial_messages(language, task, seed_failure=seed_failure)

    for attempt in range(0, max_repairs + 1):
        if dry_run:
            # Attempt 0 uses the reference solution (validates runners).
            # For repair tasks the broken starter already failed in prepare_repair_seed.
            content = load_reference(task, language)
            prompt_tokens = completion_tokens = 0
        else:
            assert api_key is not None
            content, prompt_tokens, completion_tokens = chat_completion(
                base_url=base_url,
                api_key=api_key or "ollama",
                model=model,
                temperature=temperature,
                messages=messages,
                max_tokens=max_tokens,
                request_timeout=request_timeout,
            )

        source = extract_source(content)
        if source is None:
            rec = AttemptRecord(
                suite=suite,
                task=task,
                task_kind=task_kind,
                language=language,
                model=model,
                temperature=temperature,
                replicate=replicate,
                attempt=attempt,
                prompt_tokens=prompt_tokens,
                completion_tokens=completion_tokens,
                exit_code=None,
                outcome="extract_error",
                stdout="",
                stderr="could not extract source from model output",
                source=content,
                dry_run=dry_run,
                card_tokens_est=card_tokens_est,
                repair_seeded=repair_seeded and attempt == 0,
            )
            result.attempts.append(rec)
            if attempt >= max_repairs:
                break
            messages.append({"role": "assistant", "content": content})
            messages.append(
                {
                    "role": "user",
                    "content": (
                        "Could not extract a source file from your reply. "
                        "Return only the complete program source."
                    ),
                }
            )
            continue

        stdout, stderr, rc, timed_out, diag = run_source(language, source, vol_prefix, vol_cwd)
        outcome, validation_error = judge(stdout, rc, timed_out, expected, source, source_checks)
        if validation_error:
            stderr = validation_error
        rec = AttemptRecord(
            suite=suite,
            task=task,
            task_kind=task_kind,
            language=language,
            model=model,
            temperature=temperature,
            replicate=replicate,
            attempt=attempt,
            prompt_tokens=prompt_tokens,
            completion_tokens=completion_tokens,
            exit_code=rc,
            outcome=outcome,
            diagnostic_code=diag,
            stdout=stdout,
            stderr=stderr,
            source=source,
            dry_run=dry_run,
            card_tokens_est=card_tokens_est,
            repair_seeded=repair_seeded and attempt == 0,
        )
        result.attempts.append(rec)

        if outcome == "success":
            result.success = True
            result.first_try = attempt == 0
            return result

        if attempt >= max_repairs:
            break

        if dry_run:
            # Dry-run only exercises attempt 0 with references.
            break

        messages.append({"role": "assistant", "content": content})
        messages.append(build_repair_message(source, rc, stderr))

    # Exhausted attempts without success: last attempt keeps its real outcome
    # (diag_error / wrong_output / timeout / extract_error). TaskResult.success
    # stays False (= gave up after K repairs).
    return result


def _mean_sd(values: list[float]) -> tuple[str, str]:
    if not values:
        return "—", "0.0"
    mean = statistics.mean(values)
    spread = statistics.pstdev(values) if len(values) > 1 else 0.0
    return f"{mean:.1f}", f"{spread:.1f}"


def _pct_delta(vol: float, baseline: float) -> str:
    if baseline == 0:
        return "—"
    return f"{(vol - baseline) / baseline * 100:+.1f}%"


def _baseline_languages(languages: list[str]) -> list[str]:
    return [lang for lang in languages if lang != "vol"]


def summarize(
    results: list[TaskResult],
    *,
    suite: str,
    model: str,
    temperature: float,
    max_repairs: int,
    dry_run: bool,
    card_tokens: dict[str, int],
    baseline_note: str | None = None,
    live_languages: list[str] | None = None,
) -> str:
    languages = sorted({r.language for r in results})
    lines: list[str] = [
        "# VOL LLM generate/repair results",
        "",
        f"- Date: {date.today().isoformat()}",
        f"- Suite: `{suite}`",
        f"- Protocol: v1.1",
        f"- Model: `{model}`",
        f"- Temperature: {temperature}",
        f"- Max repairs (K): {max_repairs}",
        f"- Dry-run: {dry_run}",
        (
            "- Surface freeze: "
            + ", ".join(
                f"{LANG_META[lang]['label']}={LANG_META[lang]['freeze']}"
                for lang in languages
            )
        ),
        (
            "- Cards: "
            + ", ".join(
                f"`{LANG_META[lang]['card']}` ({LANG_META[lang]['freeze']})"
                for lang in languages
            )
        ),
        (
            f"- Card tokens (est. `{CARD_TOKENIZER}`): "
            + ", ".join(
                f"{LANG_META[lang]['label']}={card_tokens[lang]}"
                for lang in languages
            )
        ),
    ]
    if baseline_note:
        live = ", ".join(live_languages or [])
        lines.append(
            f"- Baseline reuse: live langs=[{live}]; frozen rows from `{baseline_note}`"
        )
    lines.extend(
        [
            "",
            (
                "> Dry-run uses reference solutions and is **not** an LLM benchmark result."
                if dry_run
                else "> Live API run. Recompute from committed JSONL if numbers are quoted."
            ),
            "",
            "> **Cold** totals use provider `prompt_tokens` + `completion_tokens` (language card",
            "> re-sent every request). **Warm** subtracts the estimated card tokens from each",
            "> prompt (amortized / cached-card accounting).",
            "",
            "## Summary",
            "",
            "| Language | First-try % | Success @ K % | Median cold (success) | Mean ± SD cold (all) | N task-replicates |",
            "| --- | --- | --- | --- | --- | --- |",
        ]
    )

    for language in languages:
        subset = [r for r in results if r.language == language]
        n = len(subset)
        first = sum(1 for r in subset if r.first_try) / n * 100 if n else 0.0
        ok = sum(1 for r in subset if r.success) / n * 100 if n else 0.0
        success_tokens = [r.total_tokens for r in subset if r.success]
        all_tokens = [r.total_tokens for r in subset]
        med = statistics.median(success_tokens) if success_tokens else float("nan")
        mean_s, spread_s = _mean_sd(all_tokens)
        med_s = f"{med:.0f}" if success_tokens else "—"
        label = LANG_META[language]["label"]
        lines.append(
            f"| {label} | {first:.1f} | {ok:.1f} | {med_s} | {mean_s} ± {spread_s} | {n} |"
        )

    lines.extend(["", "## Prompt vs completion (cold, all attempts)", ""])
    lines.append(
        "| Language | Mean prompt | Mean completion | Mean cold total | "
        "Prompt share | Completion share |"
    )
    lines.append("| --- | --- | --- | --- | --- | --- |")
    mix: dict[str, tuple[float, float, float]] = {}
    for language in languages:
        subset = [r for r in results if r.language == language]
        if not subset:
            continue
        mean_p = statistics.mean(r.prompt_tokens for r in subset)
        mean_c = statistics.mean(r.completion_tokens for r in subset)
        mean_t = statistics.mean(r.total_tokens for r in subset)
        mix[language] = (mean_p, mean_c, mean_t)
        p_share = (mean_p / mean_t * 100) if mean_t else 0.0
        c_share = (mean_c / mean_t * 100) if mean_t else 0.0
        lines.append(
            f"| {LANG_META[language]['label']} | {mean_p:.1f} | {mean_c:.1f} | "
            f"{mean_t:.1f} | {p_share:.1f}% | {c_share:.1f}% |"
        )

    if "vol" in mix:
        vp, vc, vt = mix["vol"]
        for baseline in _baseline_languages(list(mix)):
            if baseline not in mix:
                continue
            bp, bc, bt = mix[baseline]
            blabel = LANG_META[baseline]["label"]
            lines.extend(
                [
                    "",
                    f"## VOL vs {blabel} token deltas (cold means)",
                    "",
                    f"| Metric | VOL vs {blabel} |",
                    "| --- | --- |",
                    f"| Generated completion tokens | {_pct_delta(vc, bc)} |",
                    f"| Prompt tokens | {_pct_delta(vp, bp)} |",
                    f"| Total workflow tokens (cold) | {_pct_delta(vt, bt)} |",
                    f"| Abs. prompt delta / task-replicate | {vp - bp:+.1f} |",
                    f"| Abs. completion delta / task-replicate | {vc - bc:+.1f} |",
                ]
            )

    lines.extend(["", "## Cold vs warm (card amortized)", ""])
    lines.append(
        "| Language | Mean cold | Mean warm | Warm − cold | Card est. |"
    )
    lines.append("| --- | --- | --- | --- | --- |")
    for language in languages:
        subset = [r for r in results if r.language == language]
        if not subset:
            continue
        card = card_tokens.get(language, 0)
        cold = [float(r.total_tokens) for r in subset]
        warm = [float(r.warm_tokens(card)) for r in subset]
        mean_cold = statistics.mean(cold)
        mean_warm = statistics.mean(warm)
        lines.append(
            f"| {LANG_META[language]['label']} | {mean_cold:.1f} | {mean_warm:.1f} | "
            f"{mean_warm - mean_cold:+.1f} | {card} |"
        )

    if "vol" in languages:
        vol_sub = [r for r in results if r.language == "vol"]
        vol_warm = statistics.mean(r.warm_tokens(card_tokens["vol"]) for r in vol_sub)
        for baseline in _baseline_languages(languages):
            base_sub = [r for r in results if r.language == baseline]
            if not base_sub:
                continue
            base_warm = statistics.mean(
                r.warm_tokens(card_tokens[baseline]) for r in base_sub
            )
            blabel = LANG_META[baseline]["label"]
            lines.extend(
                [
                    "",
                    f"Warm VOL vs {blabel} total: {_pct_delta(vol_warm, base_warm)} "
                    f"(means {vol_warm:.1f} vs {base_warm:.1f}).",
                ]
            )

    lines.extend(["", "## By workflow kind", ""])
    lines.append(
        "| Kind | Language | Success @ K % | Mean cold | Mean prompt | Mean completion | N |"
    )
    lines.append("| --- | --- | --- | --- | --- | --- | --- |")
    for kind in sorted({r.task_kind for r in results}):
        for language in languages:
            subset = [r for r in results if r.task_kind == kind and r.language == language]
            if not subset:
                continue
            success = sum(1 for r in subset if r.success) / len(subset) * 100
            mean_t = statistics.mean(r.total_tokens for r in subset)
            mean_p = statistics.mean(r.prompt_tokens for r in subset)
            mean_c = statistics.mean(r.completion_tokens for r in subset)
            lines.append(
                f"| {kind} | {LANG_META[language]['label']} | {success:.1f} | "
                f"{mean_t:.1f} | {mean_p:.1f} | {mean_c:.1f} | {len(subset)} |"
            )

    repair_rows = [r for r in results if r.task_kind == "repair"]
    if repair_rows:
        lines.extend(
            [
                "",
                "## Diagnostic repair notes",
                "",
                "- Repair tasks seed a failing starter **before** attempt 0.",
                "- Attempt 0 already includes exit code + diagnostics (not a blind rewrite).",
                "- `First-try` on repair means the model fixed the seed in one diagnostic turn.",
                "",
            ]
        )
        lines.append("| Task | Lang | Rep | Seed diag | Success | First-try | Attempts | Cold tokens | Last outcome |")
        lines.append("| --- | --- | --- | --- | --- | --- | --- | --- | --- |")
        for r in repair_rows:
            last = r.attempts[-1].outcome if r.attempts else "—"
            seed = r.seed_diagnostic_code or "—"
            lines.append(
                f"| {r.task} | {r.language} | {r.replicate} | {seed} | {r.success} | "
                f"{r.first_try} | {len(r.attempts)} | {r.total_tokens} | {last} |"
            )

    lines.extend(["", "## Per task-replicate", ""])
    lines.append(
        "| Task | Kind | Lang | Rep | Success | First-try | Attempts | "
        "Cold | Prompt | Completion | Warm | Last outcome |"
    )
    lines.append("| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |")
    for r in results:
        last = r.attempts[-1].outcome if r.attempts else "—"
        card = card_tokens.get(r.language, 0)
        lines.append(
            f"| {r.task} | {r.task_kind} | {r.language} | {r.replicate} | {r.success} | "
            f"{r.first_try} | {len(r.attempts)} | {r.total_tokens} | {r.prompt_tokens} | "
            f"{r.completion_tokens} | {r.warm_tokens(card)} | {last} |"
        )
    lines.append("")
    return "\n".join(lines)


def record_to_json(rec: AttemptRecord) -> dict:
    warm_prompt = max(0, rec.prompt_tokens - rec.card_tokens_est)
    return {
        "suite": rec.suite,
        "task": rec.task,
        "task_kind": rec.task_kind,
        "language": rec.language,
        "model": rec.model,
        "temperature": rec.temperature,
        "replicate": rec.replicate,
        "attempt": rec.attempt,
        "prompt_tokens": rec.prompt_tokens,
        "completion_tokens": rec.completion_tokens,
        "card_tokens_est": rec.card_tokens_est,
        "prompt_tokens_warm": warm_prompt,
        "tokens_cold": rec.prompt_tokens + rec.completion_tokens,
        "tokens_warm": warm_prompt + rec.completion_tokens,
        "repair_seeded": rec.repair_seeded,
        "exit_code": rec.exit_code,
        "outcome": rec.outcome,
        "diagnostic_code": rec.diagnostic_code,
        "stdout": rec.stdout,
        "stderr": rec.stderr,
        "source": rec.source,
        "dry_run": rec.dry_run,
    }


def attempt_from_json(row: dict) -> AttemptRecord:
    return AttemptRecord(
        suite=str(row["suite"]),
        task=str(row["task"]),
        task_kind=str(row["task_kind"]),
        language=str(row["language"]),
        model=str(row["model"]),
        temperature=float(row["temperature"]),
        replicate=int(row["replicate"]),
        attempt=int(row["attempt"]),
        prompt_tokens=int(row["prompt_tokens"]),
        completion_tokens=int(row["completion_tokens"]),
        exit_code=row.get("exit_code"),
        outcome=str(row["outcome"]),
        diagnostic_code=row.get("diagnostic_code"),
        stdout=str(row.get("stdout") or ""),
        stderr=str(row.get("stderr") or ""),
        source=str(row.get("source") or ""),
        dry_run=bool(row.get("dry_run", False)),
        card_tokens_est=int(row.get("card_tokens_est") or 0),
        repair_seeded=bool(row.get("repair_seeded", False)),
    )


def load_baseline_results(
    path: Path,
    *,
    suite: str,
    model: str,
    skip_languages: set[str],
) -> tuple[list[TaskResult], dict[str, int], str]:
    """Load frozen attempt rows for languages not being re-run.

    Returns (task results, card_tokens_est by language, baseline path note).
    """
    if not path.is_file():
        die(f"baseline JSONL not found: {path}")
    by_key: dict[tuple[str, str, int], list[AttemptRecord]] = {}
    card_tokens: dict[str, int] = {}
    models: set[str] = set()
    suites: set[str] = set()
    with path.open(encoding="utf-8") as fh:
        for line_no, line in enumerate(fh, 1):
            line = line.strip()
            if not line:
                continue
            try:
                row = json.loads(line)
            except json.JSONDecodeError as exc:
                die(f"baseline JSONL {path}:{line_no}: {exc}")
            lang = str(row.get("language", ""))
            if lang in skip_languages:
                continue
            if lang not in LANG_META:
                die(f"baseline JSONL {path}:{line_no}: unsupported language {lang!r}")
            rec = attempt_from_json(row)
            models.add(rec.model)
            suites.add(rec.suite)
            card_tokens[lang] = rec.card_tokens_est
            key = (rec.language, rec.task, rec.replicate)
            by_key.setdefault(key, []).append(rec)
    if not by_key:
        die(f"baseline JSONL {path} has no rows outside --langs {sorted(skip_languages)}")
    if suite not in suites:
        die(
            f"baseline JSONL suite(s) {sorted(suites)} do not include current suite {suite!r}"
        )
    if model not in models:
        die(
            f"baseline JSONL model(s) {sorted(models)} do not include current model {model!r}"
        )
    results: list[TaskResult] = []
    for (language, task, replicate), attempts in sorted(by_key.items()):
        attempts.sort(key=lambda a: a.attempt)
        success = any(a.outcome == "success" for a in attempts)
        first_try = bool(attempts) and attempts[0].outcome == "success"
        seed_diag = None
        if attempts and attempts[0].repair_seeded:
            seed_diag = attempts[0].diagnostic_code
        results.append(
            TaskResult(
                task=task,
                task_kind=attempts[0].task_kind,
                language=language,
                replicate=replicate,
                success=success,
                first_try=first_try,
                attempts=attempts,
                seed_diagnostic_code=seed_diag,
            )
        )
    note = str(path)
    try:
        note = str(path.resolve().relative_to(VOL_REPO.resolve()))
    except ValueError:
        pass
    return results, card_tokens, note


def resolve_tasks(suite: str, only: list[str] | None) -> list[str]:
    if suite == "smoke":
        tasks = list(SMOKE_TASKS)
    elif suite == "core":
        tasks = list(CORE_TASKS)
    elif suite == "intent":
        tasks = list(INTENT_TASKS)
    else:
        die(f"unknown suite: {suite}")
    if only:
        tasks = list(only)
        for t in tasks:
            if not (PROMPTS_DIR / t / "prompt.md").exists():
                die(f"no prompt for task {t}")
            if not (TASKS_DIR / t / "expected.txt").exists():
                die(f"no expected.txt for task {t}")
    return tasks


def suite_freeze_name(suite: str) -> str:
    return {
        "smoke": "smoke_v1",
        "core": "core_v2",
        "intent": "intent_v1",
    }.get(suite, suite)


def resolve_endpoint(provider: str, model_arg: str | None) -> tuple[str, str, str, float]:
    """Return (base_url, api_key, model, request_timeout)."""
    if provider == "ollama":
        host = (
            os.environ.get("OLLAMA_BASE_URL")
            or os.environ.get("OLLAMA_HOST")
            or DEFAULT_OLLAMA_BASE
        )
        # OLLAMA_HOST is often http://localhost:11434 without /v1
        base_url = normalize_openai_base(host)
        model = (
            model_arg
            or os.environ.get("OLLAMA_MODEL")
            or os.environ.get("OPENAI_MODEL")
            or DEFAULT_OLLAMA_MODEL
        )
        return base_url, os.environ.get("OPENAI_API_KEY") or "ollama", model, 600.0
    if provider == "openai":
        base_url = os.environ.get("OPENAI_BASE_URL", DEFAULT_OPENAI_BASE)
        model = model_arg or os.environ.get("OPENAI_MODEL") or DEFAULT_OPENAI_MODEL
        api_key = os.environ.get("OPENAI_API_KEY")
        if not api_key:
            die("OPENAI_API_KEY is required for --provider openai")
        return base_url, api_key, model, 120.0
    if provider == "gemini":
        base_url = os.environ.get("GEMINI_BASE_URL", DEFAULT_GEMINI_BASE)
        model = model_arg or os.environ.get("GEMINI_MODEL") or DEFAULT_GEMINI_MODEL
        api_key = os.environ.get("GEMINI_API_KEY")
        if not api_key:
            die("GEMINI_API_KEY is required for --provider gemini (set it in .env or the shell)")
        return base_url, api_key, model, 120.0
    die(f"unknown provider: {provider}")


def main() -> None:
    load_environment()
    parser = argparse.ArgumentParser(description="VOL LLM workflow benchmark harness")
    parser.add_argument(
        "--suite",
        choices=["smoke", "core", "intent"],
        default="smoke",
        help="smoke=wiring; core=published core_v2 (includes parity); intent=filter/map/count + repair/mod",
    )
    parser.add_argument(
        "--langs",
        default="vol,python",
        help="Comma-separated languages: vol,python,go (default: vol,python)",
    )
    parser.add_argument("--tasks", default="", help="Comma-separated task ids (optional filter)")
    parser.add_argument(
        "--provider",
        choices=["ollama", "openai", "gemini"],
        default=os.environ.get("VOL_LLM_PROVIDER", "ollama"),
        help="API provider (default: ollama)",
    )
    parser.add_argument(
        "--model",
        default=None,
        help=(
            f"Model id (ollama: {DEFAULT_OLLAMA_MODEL}; openai: {DEFAULT_OPENAI_MODEL}; "
            f"gemini: {DEFAULT_GEMINI_MODEL})"
        ),
    )
    parser.add_argument("--temperature", type=float, default=0.0)
    parser.add_argument("--replicates", type=int, default=3)
    parser.add_argument("--max-repairs", type=int, default=2, dest="max_repairs")
    parser.add_argument("--max-tokens", type=int, default=2048, dest="max_tokens")
    parser.add_argument(
        "--request-timeout",
        type=float,
        default=None,
        metavar="SECONDS",
        help="Maximum time for one model request; provider default when omitted",
    )
    parser.add_argument("--dry-run", action="store_true", help="Use reference solutions; no API")
    parser.add_argument(
        "--baseline-jsonl",
        type=Path,
        default=None,
        help=(
            "Prior frozen JSONL: merge attempt rows for languages not listed in --langs "
            "(same model + suite). Use with --langs vol to reuse published Python/Go."
        ),
    )
    parser.add_argument("--out-dir", type=Path, default=None, help="Results directory")
    args = parser.parse_args()

    if args.request_timeout is not None and args.request_timeout <= 0:
        parser.error("--request-timeout must be greater than zero")
    if args.replicates <= 0:
        parser.error("--replicates must be greater than zero")
    if args.max_repairs < 0:
        parser.error("--max-repairs cannot be negative")

    languages = [p.strip() for p in args.langs.split(",") if p.strip()]
    for lang in languages:
        if lang not in LANG_META:
            die(f"unsupported language: {lang}")

    only = [t.strip() for t in args.tasks.split(",") if t.strip()] or None
    tasks = resolve_tasks(args.suite, only)
    # Freeze ids for result tables (protocol v1.1). Do not rename published core_v2.
    suite_name = suite_freeze_name(args.suite)

    if args.dry_run:
        base_url, api_key, model_name, request_timeout = "", None, "dry-run/reference", 120.0
    else:
        base_url, api_key, model_name, request_timeout = resolve_endpoint(args.provider, args.model)
    if args.request_timeout is not None:
        request_timeout = args.request_timeout

    if "go" in languages and subprocess.run(["go", "version"], capture_output=True).returncode != 0:
        die("go toolchain not found")
    if "python" in languages:
        probe = subprocess.run(
            [sys.executable, "-c", "import sys; raise SystemExit(0 if sys.version_info >= (3, 11) else 1)"],
            capture_output=True,
        )
        if probe.returncode != 0:
            die(f"Python 3.11+ required for python runner (using {sys.executable})")

    card_tokens = {lang: estimate_card_tokens(load_card(lang)) for lang in languages}

    vol_prefix, vol_cwd = get_vol_prefix()
    out_dir = args.out_dir or RESULTS_DIR
    out_dir.mkdir(parents=True, exist_ok=True)
    stamp = time.strftime("%Y%m%d-%H%M%S")
    mode = "dryrun" if args.dry_run else f"live_{args.provider}"
    # Sanitize model for filenames: qwen3:4b-instruct -> qwen3_4b-instruct
    model_slug = re.sub(r"[^\w.-]+", "_", model_name)
    jsonl_path = out_dir / f"{suite_name}_{mode}_{model_slug}_{stamp}.jsonl"
    partial_jsonl_path = jsonl_path.with_suffix(".partial.jsonl")
    summary_path = out_dir / f"{suite_name}_{mode}_{model_slug}_{stamp}.md"

    model_label = model_name
    results: list[TaskResult] = []
    baseline_note: str | None = None
    baseline_results: list[TaskResult] = []

    if args.baseline_jsonl is not None:
        baseline_results, baseline_cards, baseline_note = load_baseline_results(
            args.baseline_jsonl,
            suite=suite_name,
            model=model_label,
            skip_languages=set(languages),
        )
        for lang, est in baseline_cards.items():
            card_tokens.setdefault(lang, est)
        results.extend(baseline_results)

    print(
        f"suite={suite_name} mode={mode} provider={args.provider if not args.dry_run else 'none'} "
        f"model={model_label} base={base_url or '-'} tasks={len(tasks)} langs={languages} "
        f"card_tokens={card_tokens}"
        + (f" baseline={baseline_note}" if baseline_note else "")
    )

    with partial_jsonl_path.open("w", encoding="utf-8") as jsonl:
        for tr in baseline_results:
            for attempt in tr.attempts:
                jsonl.write(json.dumps(record_to_json(attempt), ensure_ascii=False) + "\n")
        for language in languages:
            for task in tasks:
                for replicate in range(1, args.replicates + 1):
                    print(f"  {language} {task} rep={replicate} ...", flush=True)
                    tr = run_task_replicate(
                        suite=suite_name,
                        task=task,
                        language=language,
                        model=model_label,
                        temperature=args.temperature,
                        replicate=replicate,
                        max_repairs=args.max_repairs,
                        max_tokens=args.max_tokens,
                        dry_run=args.dry_run,
                        api_key=api_key,
                        base_url=base_url,
                        request_timeout=request_timeout,
                        vol_prefix=vol_prefix,
                        vol_cwd=vol_cwd,
                        card_tokens_est=card_tokens[language],
                    )
                    results.append(tr)
                    for attempt in tr.attempts:
                        jsonl.write(json.dumps(record_to_json(attempt), ensure_ascii=False) + "\n")
                    status = "OK" if tr.success else "FAIL"
                    last = tr.attempts[-1].outcome if tr.attempts else "?"
                    warm = tr.warm_tokens(card_tokens[language])
                    seed = f" seed_diag={tr.seed_diagnostic_code}" if tr.seed_diagnostic_code else ""
                    print(
                        f"    {status} attempts={len(tr.attempts)} last={last} "
                        f"cold={tr.total_tokens} warm={warm}{seed}"
                    )

    partial_jsonl_path.replace(jsonl_path)

    summary = summarize(
        results,
        suite=suite_name,
        model=model_label,
        temperature=args.temperature,
        max_repairs=args.max_repairs,
        dry_run=args.dry_run,
        card_tokens=card_tokens,
        baseline_note=baseline_note,
        live_languages=languages,
    )
    summary_path.write_text(summary, encoding="utf-8")
    print()
    print(summary)
    print(f"Wrote {jsonl_path}")
    print(f"Wrote {summary_path}")

    if args.dry_run and not all(r.success for r in results):
        die("dry-run expected all reference solutions to succeed", code=1)


if __name__ == "__main__":
    main()
