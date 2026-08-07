#!/usr/bin/env python3
"""
check_stdout.py — verify that every language implementation for each task
produces the same stdout as expected.txt.

Usage:
    python harness/check_stdout.py

Environment:
    VOL_BIN   Path to a pre-built vol binary. If unset, falls back to
              go run ./cmd/vol inside the VOL repo (sibling of this repo).
"""

import os
import subprocess
import sys
import tempfile
from pathlib import Path

REPO_ROOT = Path(__file__).parent.parent.resolve()   # VOL/bench
TASKS_DIR = REPO_ROOT / "tasks"
VOL_REPO = REPO_ROOT.parent                          # VOL/


def get_vol_runner():
    """Return (cmd_prefix, cwd) for running a .vol file.

    The vol file path is appended to cmd_prefix by the caller.
    cwd is the working directory for the subprocess (needed so Go can find
    the module when falling back to 'go run ./cmd/vol').
    """
    vol_bin = os.environ.get("VOL_BIN")
    if vol_bin:
        return [vol_bin], None

    cmd_vol = VOL_REPO / "cmd" / "vol"
    if cmd_vol.exists():
        return ["go", "run", "./cmd/vol"], str(VOL_REPO)

    print(
        "ERROR: VOL_BIN not set and ../VOL/cmd/vol not found.\n"
        "  Set VOL_BIN=/path/to/vol or place vol-llm-bench next to the VOL repo.",
        file=sys.stderr,
    )
    sys.exit(1)


def tool_available(name: str) -> bool:
    try:
        subprocess.run([name, "--version"], capture_output=True, check=False)
        return True
    except FileNotFoundError:
        return False


def run(cmd: list[str], cwd: str | None = None) -> tuple[str, str, int]:
    """Returns (stdout, stderr, returncode)."""
    result = subprocess.run(cmd, capture_output=True, text=True, cwd=cwd)
    return result.stdout, result.stderr, result.returncode


def check_task(
    task_dir: Path,
    vol_prefix: list[str],
    vol_cwd: str | None,
    available: dict[str, bool],
) -> list[str]:
    expected = (task_dir / "expected.txt").read_text()
    errors: list[str] = []

    # VOL
    stdout, stderr, rc = run(vol_prefix + [str(task_dir / "vol" / "main.vol")], cwd=vol_cwd)
    if rc != 0:
        detail = stderr.strip().splitlines()[0] if stderr.strip() else "(no output)"
        errors.append(f"vol: exit {rc} — {detail}")
    elif stdout != expected:
        errors.append(f"vol: stdout mismatch\n    got:      {stdout!r}\n    expected: {expected!r}")

    # Go
    if available["go"]:
        stdout, stderr, rc = run(["go", "run", str(task_dir / "go" / "main.go")])
        if rc != 0:
            detail = stderr.strip().splitlines()[0] if stderr.strip() else "(no output)"
            errors.append(f"go: exit {rc} — {detail}")
        elif stdout != expected:
            errors.append(f"go: stdout mismatch\n    got:      {stdout!r}\n    expected: {expected!r}")
    else:
        errors.append("go: toolchain not found (skipped)")

    # Rust
    if available["rust"]:
        with tempfile.TemporaryDirectory() as tmpdir:
            bin_path = os.path.join(tmpdir, "task_bin")
            compile_result = subprocess.run(
                ["rustc", str(task_dir / "rust" / "main.rs"), "-o", bin_path],
                capture_output=True,
                text=True,
            )
            if compile_result.returncode != 0:
                first_err = compile_result.stderr.strip().splitlines()[0] if compile_result.stderr.strip() else ""
                errors.append(f"rust: compile failed — {first_err}")
            else:
                stdout, stderr, rc = run([bin_path])
                if rc != 0:
                    errors.append(f"rust: exit {rc}")
                elif stdout != expected:
                    errors.append(f"rust: stdout mismatch\n    got:      {stdout!r}\n    expected: {expected!r}")
    else:
        errors.append("rust: toolchain not found (skipped)")

    # Zig
    if available["zig"]:
        stdout, stderr, rc = run(["zig", "run", str(task_dir / "zig" / "main.zig")])
        if rc != 0:
            detail = stderr.strip().splitlines()[0] if stderr.strip() else "(no output)"
            errors.append(f"zig: exit {rc} — {detail}")
        elif stdout != expected:
            errors.append(f"zig: stdout mismatch\n    got:      {stdout!r}\n    expected: {expected!r}")
    else:
        errors.append("zig: toolchain not found (skipped)")

    return errors


def main() -> None:
    vol_prefix, vol_cwd = get_vol_runner()
    available = {
        "go": tool_available("go"),
        "rust": tool_available("rustc"),
        "zig": tool_available("zig"),
    }

    missing = [lang for lang, ok in available.items() if not ok]
    if missing:
        print(f"WARNING: toolchains not found, will skip: {', '.join(missing)}")

    tasks = sorted(d for d in TASKS_DIR.iterdir() if d.is_dir())
    if not tasks:
        print("No tasks found under tasks/", file=sys.stderr)
        sys.exit(1)

    all_ok = True
    for task in tasks:
        errors = check_task(task, vol_prefix, vol_cwd, available)
        skipped = [e for e in errors if "skipped" in e]
        failures = [e for e in errors if "skipped" not in e]
        if failures:
            print(f"FAIL  {task.name}")
            for e in failures:
                print(f"      {e}")
            all_ok = False
        elif skipped:
            skip_langs = ", ".join(e.split(":")[0] for e in skipped)
            print(f"SKIP  {task.name}  (no toolchain: {skip_langs})")
        else:
            print(f"OK    {task.name}")

    print()
    if all_ok:
        print("All implemented tasks pass.")
    else:
        print("Some tasks failed — see above.")
        sys.exit(1)


if __name__ == "__main__":
    main()
