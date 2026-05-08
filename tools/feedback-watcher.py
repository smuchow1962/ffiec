#!/usr/bin/env python3
"""Polling watcher for docs/feedback/ — detects file drops without `watchdog` dep.

Runs in the foreground, polls the feedback directory, prints a log line per
detected change, and (with --on-change scrape) re-runs the scraper.

Usage:
  python tools/feedback-watcher.py
  python tools/feedback-watcher.py --interval 10 --on-change scrape
  python tools/feedback-watcher.py --once    # one-shot scan; no loop

State and log files (auto-created, hidden):
  docs/feedback/.watcher-state.json    # last-seen file -> mtime/size
  docs/feedback/.watcher-log.ndjson    # one JSON line per detected event

Stdlib-only. No external dependencies. Cross-platform.
"""

from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
import time
from dataclasses import dataclass
from pathlib import Path

REPO_ROOT_DEFAULT = Path(__file__).resolve().parent.parent
FEEDBACK_DIR = "docs/feedback"
STATE_FILENAME = ".watcher-state.json"
LOG_FILENAME = ".watcher-log.ndjson"


@dataclass
class FileSnapshot:
    path: str
    mtime: float
    size: int


SCRAPER_OUTPUT_FILES = {"SCRAPE-REPORT.md", "scrape.json"}


def snapshot_dir(root: Path) -> dict[str, FileSnapshot]:
    snap: dict[str, FileSnapshot] = {}
    for path in root.rglob("*"):
        if path.is_dir():
            continue
        if path.name.startswith("."):
            continue
        if path.name in SCRAPER_OUTPUT_FILES:
            continue
        try:
            st = path.stat()
        except OSError:
            continue
        rel = str(path.relative_to(root)).replace("\\", "/")
        snap[rel] = FileSnapshot(path=rel, mtime=st.st_mtime, size=st.st_size)
    return snap


def load_state(state_path: Path) -> dict[str, FileSnapshot]:
    if not state_path.is_file():
        return {}
    try:
        raw = json.loads(state_path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError):
        return {}
    return {p: FileSnapshot(**v) for p, v in raw.items()}


def save_state(state_path: Path, snap: dict[str, FileSnapshot]) -> None:
    payload = {p: {"path": s.path, "mtime": s.mtime, "size": s.size} for p, s in snap.items()}
    state_path.write_text(json.dumps(payload, indent=2), encoding="utf-8")


def diff(prev: dict[str, FileSnapshot], curr: dict[str, FileSnapshot]) -> list[dict[str, str]]:
    events: list[dict[str, str]] = []
    prev_keys = set(prev.keys())
    curr_keys = set(curr.keys())

    for added in sorted(curr_keys - prev_keys):
        s = curr[added]
        events.append({"event": "added", "path": s.path, "size": str(s.size)})
    for removed in sorted(prev_keys - curr_keys):
        events.append({"event": "removed", "path": removed})
    for shared in sorted(prev_keys & curr_keys):
        p, c = prev[shared], curr[shared]
        if p.mtime != c.mtime or p.size != c.size:
            events.append({"event": "modified", "path": shared, "size": str(c.size)})
    return events


def append_log(log_path: Path, events: list[dict[str, str]]) -> None:
    if not events:
        return
    ts = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
    with log_path.open("a", encoding="utf-8") as f:
        for e in events:
            f.write(json.dumps({"ts": ts, **e}) + "\n")


def trigger_scrape(repo_root: Path) -> None:
    script = repo_root / "tools" / "feedback-scrape.py"
    if not script.is_file():
        print(f"  [warn] scraper not found at {script}", file=sys.stderr)
        return
    try:
        subprocess.run([sys.executable, str(script), "--quiet"], check=False, cwd=str(repo_root))
    except OSError as exc:
        print(f"  [warn] scraper invocation failed: {exc}", file=sys.stderr)


def print_event(e: dict[str, str]) -> None:
    ts = time.strftime("%H:%M:%S")
    tag = e["event"].upper().ljust(8)
    print(f"[{ts}] {tag} {e['path']}", flush=True)


def run_loop(args: argparse.Namespace) -> int:
    repo_root = Path(args.root).resolve()
    feedback_root = repo_root / FEEDBACK_DIR
    if not feedback_root.is_dir():
        print(f"error: feedback directory not found: {feedback_root}", file=sys.stderr)
        return 2

    state_path = feedback_root / STATE_FILENAME
    log_path = feedback_root / LOG_FILENAME

    print(f"watching {feedback_root}")
    print(f"  state: {state_path.relative_to(repo_root)}")
    print(f"  log:   {log_path.relative_to(repo_root)}")
    print(f"  poll:  every {args.interval}s")
    if args.on_change:
        print(f"  on-change: {args.on_change}")
    print("press Ctrl+C to stop")
    print()

    state = load_state(state_path)
    initial = not state
    if initial:
        print("  [init] no prior state; baselining current files")

    while True:
        curr = snapshot_dir(feedback_root)
        events = diff(state, curr)

        if events:
            for e in events:
                print_event(e)
            append_log(log_path, events)
            if args.on_change == "scrape" and not initial:
                print("  [scrape] running tools/feedback-scrape.py")
                trigger_scrape(repo_root)

        state = curr
        save_state(state_path, state)
        initial = False

        if args.once:
            return 0

        try:
            time.sleep(args.interval)
        except KeyboardInterrupt:
            print()
            print("watcher stopped")
            return 0


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description="Polling watcher for docs/feedback/")
    parser.add_argument("--root", default=str(REPO_ROOT_DEFAULT), help="Repo root (default: parent of tools/)")
    parser.add_argument("--interval", type=float, default=30.0, help="Poll interval in seconds (default 30)")
    parser.add_argument("--on-change", choices=["scrape"], default=None, help="Action on each detected change")
    parser.add_argument("--once", action="store_true", help="One-shot scan and exit")
    args = parser.parse_args(argv)
    return run_loop(args)


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
