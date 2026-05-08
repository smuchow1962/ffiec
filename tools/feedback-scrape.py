#!/usr/bin/env python3
"""Scrape docs/feedback/ for reviewer findings.

Walks the entire feedback tree, parses every reviewer-shaped markdown file,
extracts Status counts and findings text, and produces:

  - A consolidated markdown report at docs/feedback/SCRAPE-REPORT.md
  - A machine-readable JSON summary at docs/feedback/scrape.json

Recognizes:
  - Spawned-reviewer files (round-N/0X-role.md)
  - Per-round roll-ups (round-N/99-gap-roll-up.md) — read separately
  - Outside-reviewer drops (outside/*.md OR round-N/outside-*.md)
  - Top-level README files (excluded from finding counts)

Skips:
  - historical-round-8-pre-hmac-rework/ (unless --include-historical)
  - Hidden state files (.watcher-state.json, .watcher-log.ndjson)

Run from anywhere:
  python tools/feedback-scrape.py [--include-historical] [--root <path>]
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import dataclass, field, asdict
from pathlib import Path

REPO_ROOT_DEFAULT = Path(__file__).resolve().parent.parent
FEEDBACK_DIR = "docs/feedback"

PERSONA_REVIEWER_DOT_RE = re.compile(r"\*\*Reviewer\.\*\*\s+(?P<name>(?:Dr\.\s+|Mr\.\s+|Ms\.\s+|Prof\.\s+)?[A-Z][^\n.]{2,80}?)\.", re.MULTILINE)
PERSONA_NAME_RE = re.compile(r"^\*\*(?:Reviewer|Persona|Name)\*\*[:.]?\s*(?:—|-|:)\s*\*?\*?(?P<name>[^\n*.]+?)\*?\*?[\n.]", re.MULTILINE)
PERSONA_HEADING_RE = re.compile(r"^##+\s+Persona\s*[—-]\s*(?P<name>[^\n]+?)\s*$", re.MULTILINE)
PERSONA_INLINE_RE = re.compile(r"persona[:.]?\s+\*\*(?P<name>[^*]+?)\*\*", re.IGNORECASE)

STATUS_LINE_RE = re.compile(
    r"^\s*\*\*Status\*\*[:.]?\s*(?P<status>Answered|Partial|Gap|Pass|Fail)\b",
    re.IGNORECASE | re.MULTILINE,
)
ROLLUP_TABLE_ROW_RE = re.compile(
    r"^\|\s*(?P<status>Answered|Partial|Gap)\s*\|\s*(?P<count>\d+)\s*\|",
    re.IGNORECASE | re.MULTILINE,
)
VERDICT_RE = re.compile(
    r"\*\*(?:Verdict|Disposition|Result)\.?\*\*\s+(?P<gaps>\d+)\s+gaps?\s*[,/]?\s+(?P<partials>\d+)\s+partials?",
    re.IGNORECASE,
)
VERDICT_FLIPPED_RE = re.compile(
    r"\*\*(?:Verdict|Disposition|Result)\.?\*\*\s+(?P<partials>\d+)\s+partials?\s*[,/]?\s+(?P<gaps>\d+)\s+gaps?",
    re.IGNORECASE,
)
VERDICT_SINGLE_PARTIAL_RE = re.compile(r"\*\*(?P<partials>\d+)\s+partials?\b", re.IGNORECASE)
VERDICT_SINGLE_GAP_RE = re.compile(r"\*\*(?P<gaps>\d+)\s+gaps?\b", re.IGNORECASE)
QUESTION_MARKER_RE = re.compile(
    r"^(?:\*\*Q[-\s]?\d+\b|### Q[-\s]?\d+\b|## Q[-\s]?\d+\b|\*\*Question[-\s]+\d+\b|### Question[-\s]+\d+\b|## Question[-\s]+\d+\b)",
    re.MULTILINE,
)


@dataclass
class ReviewerFile:
    path: str
    round: str | None
    role_slot: str | None
    persona: str | None
    is_outside: bool
    is_rollup: bool
    answered: int = 0
    partial: int = 0
    gap: int = 0
    rollup_answered: int | None = None
    rollup_partial: int | None = None
    rollup_gap: int | None = None
    headline: str | None = None


@dataclass
class RoundSummary:
    round: str
    files: list[ReviewerFile] = field(default_factory=list)
    answered: int = 0
    partial: int = 0
    gap: int = 0
    converged_count: int = 0
    total_count: int = 0


def detect_round(rel_path: Path) -> str | None:
    """Return 'round-N' if path is under a round directory, 'historical' for round-8 archive, else None."""
    parts = rel_path.parts
    if not parts:
        return None
    if parts[0].startswith("historical-round-"):
        return parts[0]
    if parts[0].startswith("round-"):
        return parts[0]
    if parts[0] == "outside":
        return "outside"
    return "top-level"


def detect_role_slot(filename: str) -> str | None:
    """Return '01-cryptographic-engineer' style slot from filename, or None."""
    stem = Path(filename).stem
    if re.match(r"^\d\d-", stem):
        return stem
    if stem.startswith("outside-"):
        return stem
    return None


def is_outside_file(filename: str, parent: str) -> bool:
    return parent == "outside" or filename.startswith("outside-")


def is_rollup_file(filename: str) -> bool:
    return filename.startswith("99-gap-roll-up") or filename == "SCRAPE-REPORT.md"


def is_excluded_file(filename: str) -> bool:
    if filename.startswith("."):
        return True
    if filename in ("README.md", "scrape.json"):
        return True
    return False


def extract_persona(text: str) -> str | None:
    for matcher in (PERSONA_REVIEWER_DOT_RE, PERSONA_HEADING_RE, PERSONA_NAME_RE, PERSONA_INLINE_RE):
        m = matcher.search(text)
        if m:
            name = m.group("name").strip().rstrip(",.")
            name = re.sub(r"\s+", " ", name)
            if name.lower() == "name":
                continue
            if 2 <= len(name) <= 80:
                return name
    return None


def extract_verdict(text: str) -> tuple[int, int] | None:
    """Return (gaps, partials) from a verdict line, or None if not present.

    Tries paired forms first (e.g. "0 gaps, 0 partials"), then falls back to
    single-side standalone forms (e.g. "**1 Partial**" or "**3 Gaps**" used
    on their own when only one disposition is non-zero).
    """
    for matcher in (VERDICT_RE, VERDICT_FLIPPED_RE):
        m = matcher.search(text)
        if m:
            return int(m.group("gaps")), int(m.group("partials"))

    p_match = VERDICT_SINGLE_PARTIAL_RE.search(text)
    g_match = VERDICT_SINGLE_GAP_RE.search(text)
    if p_match or g_match:
        gaps = int(g_match.group("gaps")) if g_match else 0
        partials = int(p_match.group("partials")) if p_match else 0
        return gaps, partials
    return None


def count_questions(text: str) -> int:
    return len(QUESTION_MARKER_RE.findall(text))


def count_per_question_statuses(text: str) -> tuple[int, int, int]:
    """Count Status: Answered / Partial / Gap lines (case-insensitive)."""
    answered = partial = gap = 0
    for m in STATUS_LINE_RE.finditer(text):
        s = m.group("status").lower()
        if s == "answered" or s == "pass":
            answered += 1
        elif s == "partial":
            partial += 1
        elif s == "gap" or s == "fail":
            gap += 1
    return answered, partial, gap


def parse_rollup_table(text: str) -> tuple[int | None, int | None, int | None]:
    """If the file ends with a rollup table 'Status | Count' rows, parse them."""
    answered = partial = gap = None
    for m in ROLLUP_TABLE_ROW_RE.finditer(text):
        s = m.group("status").lower()
        n = int(m.group("count"))
        if s == "answered":
            answered = n
        elif s == "partial":
            partial = n
        elif s == "gap":
            gap = n
    return answered, partial, gap


def extract_headline(text: str) -> str | None:
    """First non-empty paragraph after the first heading."""
    lines = text.splitlines()
    seen_heading = False
    for line in lines:
        if line.startswith("#"):
            seen_heading = True
            continue
        if seen_heading and line.strip() and not line.startswith(">"):
            stripped = line.strip()
            if len(stripped) > 30:
                return stripped[:200] + ("..." if len(stripped) > 200 else "")
    return None


def parse_reviewer_file(path: Path, rel_path: Path) -> ReviewerFile:
    text = path.read_text(encoding="utf-8", errors="replace")
    round_id = detect_round(rel_path)
    parent = rel_path.parts[0] if rel_path.parts else ""
    role_slot = detect_role_slot(path.name)
    is_outside = is_outside_file(path.name, parent)
    is_rollup = is_rollup_file(path.name)

    persona = extract_persona(text) if not is_rollup else None
    a, p, g = count_per_question_statuses(text) if not is_rollup else (0, 0, 0)
    rollup_a, rollup_p, rollup_g = parse_rollup_table(text) if not is_rollup else (None, None, None)

    # Precedence: rollup table > verdict line > per-question Status counts.
    if not is_rollup and rollup_a is not None:
        a, p, g = rollup_a, rollup_p or 0, rollup_g or 0
    elif not is_rollup:
        verdict = extract_verdict(text)
        if verdict is not None:
            v_gaps, v_partials = verdict
            q_count = count_questions(text)
            answered_from_verdict = max(0, q_count - v_gaps - v_partials) if q_count > 0 else 0
            a, p, g = answered_from_verdict, v_partials, v_gaps

    headline = extract_headline(text)

    return ReviewerFile(
        path=str(rel_path).replace("\\", "/"),
        round=round_id,
        role_slot=role_slot,
        persona=persona,
        is_outside=is_outside,
        is_rollup=is_rollup,
        answered=a,
        partial=p,
        gap=g,
        rollup_answered=rollup_a,
        rollup_partial=rollup_p,
        rollup_gap=rollup_g,
        headline=headline,
    )


def walk_feedback(feedback_root: Path, include_historical: bool) -> list[ReviewerFile]:
    files: list[ReviewerFile] = []
    for path in sorted(feedback_root.rglob("*.md")):
        rel = path.relative_to(feedback_root)
        if not include_historical and rel.parts and rel.parts[0].startswith("historical-"):
            continue
        if is_excluded_file(path.name):
            continue
        files.append(parse_reviewer_file(path, rel))
    return files


def aggregate_by_round(files: list[ReviewerFile]) -> dict[str, RoundSummary]:
    rounds: dict[str, RoundSummary] = {}
    for f in files:
        if f.is_rollup:
            continue
        key = f.round or "top-level"
        if key not in rounds:
            rounds[key] = RoundSummary(round=key)
        s = rounds[key]
        s.files.append(f)
        s.answered += f.answered
        s.partial += f.partial
        s.gap += f.gap
        s.total_count += 1
        if f.gap == 0 and f.partial == 0:
            s.converged_count += 1
    return rounds


def render_markdown_report(files: list[ReviewerFile], rounds: dict[str, RoundSummary]) -> str:
    lines: list[str] = []
    lines.append("# Feedback scrape report")
    lines.append("")
    lines.append("> Auto-generated by `tools/feedback-scrape.py`. Do not edit by hand; re-run the scraper.")
    lines.append("")

    # Outside drops first — these are the highest-priority items.
    outside = [f for f in files if f.is_outside and not f.is_rollup]
    if outside:
        lines.append("## Outside-reviewer drops")
        lines.append("")
        lines.append("| File | Persona | Answered | Partial | Gap |")
        lines.append("|---|---|---|---|---|")
        for f in outside:
            persona = f.persona or "(unknown)"
            lines.append(f"| `{f.path}` | {persona} | {f.answered} | {f.partial} | {f.gap} |")
        lines.append("")
        for f in outside:
            if f.headline:
                lines.append(f"- **`{f.path}`** — {f.headline}")
        lines.append("")
    else:
        lines.append("## Outside-reviewer drops")
        lines.append("")
        lines.append("(none — drop files into `docs/feedback/outside/` or `docs/feedback/round-N/outside-*.md` to surface them here)")
        lines.append("")

    lines.append("## Per-round summary")
    lines.append("")
    lines.append("| Round | Files | Converged 0/0 | Answered | Partial | Gap |")
    lines.append("|---|---|---|---|---|---|")
    for key in sorted(rounds.keys()):
        s = rounds[key]
        lines.append(f"| {s.round} | {s.total_count} | {s.converged_count}/{s.total_count} | {s.answered} | {s.partial} | {s.gap} |")
    lines.append("")

    lines.append("## Per-file detail")
    lines.append("")
    for key in sorted(rounds.keys()):
        s = rounds[key]
        lines.append(f"### {s.round}")
        lines.append("")
        lines.append("| File | Persona | A | P | G | Status |")
        lines.append("|---|---|---|---|---|---|")
        for f in s.files:
            persona = f.persona or "(unknown)"
            status = "✓ converged" if (f.gap == 0 and f.partial == 0) else "open"
            lines.append(f"| `{f.path}` | {persona} | {f.answered} | {f.partial} | {f.gap} | {status} |")
        lines.append("")

    lines.append("## Trajectory")
    lines.append("")
    lines.append("| Round | Gaps | Partials | Converged |")
    lines.append("|---|---|---|---|")
    for key in sorted(rounds.keys()):
        s = rounds[key]
        lines.append(f"| {s.round} | {s.gap} | {s.partial} | {s.converged_count}/{s.total_count} |")
    lines.append("")

    return "\n".join(lines) + "\n"


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description="Scrape docs/feedback/ for reviewer findings.")
    parser.add_argument("--root", default=str(REPO_ROOT_DEFAULT), help="Repo root (default: parent of tools/)")
    parser.add_argument("--include-historical", action="store_true", help="Include historical-round-8 archive")
    parser.add_argument("--quiet", action="store_true", help="Suppress per-file logging")
    args = parser.parse_args(argv)

    repo_root = Path(args.root).resolve()
    feedback_root = repo_root / FEEDBACK_DIR
    if not feedback_root.is_dir():
        print(f"error: feedback directory not found: {feedback_root}", file=sys.stderr)
        return 2

    files = walk_feedback(feedback_root, args.include_historical)
    rounds = aggregate_by_round(files)

    if not args.quiet:
        print(f"scanned {len(files)} reviewer files in {feedback_root}")
        for key in sorted(rounds.keys()):
            s = rounds[key]
            print(f"  {s.round}: {s.total_count} files, converged {s.converged_count}/{s.total_count}, A={s.answered} P={s.partial} G={s.gap}")

    md_path = feedback_root / "SCRAPE-REPORT.md"
    md_path.write_text(render_markdown_report(files, rounds), encoding="utf-8")

    json_path = feedback_root / "scrape.json"
    json_payload = {
        "files": [asdict(f) for f in files],
        "rounds": {k: {"round": v.round, "answered": v.answered, "partial": v.partial, "gap": v.gap, "converged": v.converged_count, "total": v.total_count} for k, v in rounds.items()},
    }
    json_path.write_text(json.dumps(json_payload, indent=2), encoding="utf-8")

    if not args.quiet:
        print(f"wrote {md_path.relative_to(repo_root)}")
        print(f"wrote {json_path.relative_to(repo_root)}")

    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
