# -*- coding: utf-8 -*-
"""
Inline minimal RFC 8785 JCS implementation. Used as a SECOND, independent
implementation to cross-validate the bytes the `jcs` library produced.

This is not a production JCS canonicaliser — it is a minimal reference
written from RFC 8785 directly, for the sole purpose of catching a bug in
the library that produced our `expected.json`.

Disagreements between this implementation and the library are reported and
must be resolved (either by fixing this reference or filing a finding for
the WG) before the fixture is considered authoritative.
"""

from __future__ import annotations

import json
import math
import os
import sys
from typing import Any

HERE = os.path.dirname(os.path.abspath(__file__))


# RFC 8259 §7 short-escape table. All other U+0000..U+001F use \u00XX.
_SHORT_ESCAPES = {
    0x08: "\\b",
    0x09: "\\t",
    0x0A: "\\n",
    0x0C: "\\f",
    0x0D: "\\r",
    0x22: "\\\"",
    0x5C: "\\\\",
}


def _encode_string(s: str) -> str:
    """RFC 8785 §3.2.2.2 string serialisation.

    Escape rules: only the JSON-mandatory escapes. Use the short form
    (\\b \\t \\n \\f \\r \\" \\\\) where defined; otherwise \\u00XX with
    LOWERCASE hex digits. All other characters serialise as raw UTF-8.
    """

    out: list[str] = ['"']
    for ch in s:
        cp = ord(ch)
        short = _SHORT_ESCAPES.get(cp)
        if short is not None:
            out.append(short)
        elif cp < 0x20:
            out.append(f"\\u{cp:04x}")
        else:
            out.append(ch)
    out.append('"')
    return "".join(out)


def _encode_number(n: float | int) -> str:
    """RFC 8785 §3.2.2.3: ECMAScript Number-to-String canonical form.

    Reject NaN and Infinity. Integer-valued doubles serialise without a
    decimal. The threshold for exponent notation is |n| < 1e-6 or |n| >= 1e21.
    For values in the normal range, emit the shortest round-trip form.

    We do NOT reimplement the full ECMA-262 7.1.12.1 algorithm here; we
    delegate to Python's repr/format. Python 3 implements shortest-round-trip
    repr for floats, which matches ECMA-262 for all values that fit a double.
    The exponent format requires translation: Python prints 1e+21 as `1e+21`
    which matches JCS, but 1e10 prints as `10000000000.0` which we strip.
    """

    if isinstance(n, bool):
        # bool is a subclass of int; encode separately at caller.
        raise TypeError("bool should not reach _encode_number")

    if isinstance(n, int):
        # Integer literal — but JCS treats all numbers as IEEE-754 doubles.
        # Round-trip through float to enforce that semantic.
        return _encode_number(float(n))

    if math.isnan(n) or math.isinf(n):
        raise ValueError(f"Invalid JSON number: {n}")

    if n == 0.0:
        return "0"

    abs_n = abs(n)
    # Python's repr() yields the shortest round-trip form.
    s = repr(n)

    # Decide between fixed and exponent form per ES2017 7.1.12.1.
    # Exponent form when n != 0 and (abs_n < 1e-6 or abs_n >= 1e21).
    use_exp = abs_n < 1e-6 or abs_n >= 1e21

    if use_exp:
        # Python may already give exponent form; if not, force it.
        if "e" not in s:
            s = "{:e}".format(n)
        # Normalise: drop leading zeros in exponent, keep sign.
        mantissa, _, exp_part = s.partition("e")
        # Strip trailing zeros after a decimal point, then trailing dot.
        if "." in mantissa:
            mantissa = mantissa.rstrip("0").rstrip(".")
        exp_sign = "+"
        if exp_part.startswith("+") or exp_part.startswith("-"):
            exp_sign = exp_part[0]
            exp_part = exp_part[1:]
        exp_part = exp_part.lstrip("0") or "0"
        return f"{mantissa}e{exp_sign}{exp_part}"

    # Fixed form.
    if "e" in s:
        # Python printed exponent form for an integer-valued or fixed-range
        # value (rare; can happen for very small fractions). Convert to fixed.
        s = "{:.17f}".format(n).rstrip("0").rstrip(".")
        if not s or s == "-":
            s = "0"
    if s.endswith(".0"):
        s = s[:-2]
    elif "." in s:
        s = s.rstrip("0").rstrip(".")
    return s


def canonicalize(value: Any) -> bytes:
    """RFC 8785 entry point. Returns UTF-8 canonical bytes."""

    return _encode(value).encode("utf-8")


def _encode(value: Any) -> str:
    """Recursive encoder. Cognitive complexity kept low by per-type dispatch."""

    if value is None:
        return "null"
    if value is True:
        return "true"
    if value is False:
        return "false"
    if isinstance(value, str):
        return _encode_string(value)
    if isinstance(value, (int, float)):
        return _encode_number(value)
    if isinstance(value, list):
        return "[" + ",".join(_encode(v) for v in value) + "]"
    if isinstance(value, dict):
        # RFC 8785 §3.2.3: sort keys by UTF-16 code-unit lexicographic order.
        # For BMP-only keys this matches Unicode codepoint sort.
        sorted_items = sorted(value.items(), key=lambda kv: _utf16_codeunits(kv[0]))
        return "{" + ",".join(_encode_string(k) + ":" + _encode(v) for k, v in sorted_items) + "}"
    raise TypeError(f"Unsupported type: {type(value).__name__}")


def _utf16_codeunits(s: str) -> tuple[int, ...]:
    """Return the UTF-16 code-unit sequence for sort comparison."""

    units: list[int] = []
    for ch in s:
        cp = ord(ch)
        if cp <= 0xFFFF:
            units.append(cp)
        else:
            # Encode as surrogate pair.
            cp -= 0x10000
            units.append(0xD800 | (cp >> 10))
            units.append(0xDC00 | (cp & 0x3FF))
    return tuple(units)


def _decode_input(value: Any) -> Any:
    if isinstance(value, dict):
        marker = value.get("__special__") if len(value) == 1 else None
        if marker == "NaN":
            return float("nan")
        if marker == "Infinity":
            return float("inf")
        if marker == "-Infinity":
            return float("-inf")
        return {k: _decode_input(v) for k, v in value.items()}
    if isinstance(value, list):
        return [_decode_input(v) for v in value]
    return value


def main() -> int:
    fixture_path = os.path.join(HERE, "fixture.json")
    expected_path = os.path.join(HERE, "expected.json")
    with open(fixture_path, encoding="utf-8") as fh:
        fixture = json.load(fh)
    with open(expected_path, encoding="utf-8") as fh:
        expected = json.load(fh)

    disagreements: list[str] = []
    matched = 0
    rejected_match = 0

    for cat_name, cat in fixture["categories"].items():
        for case_name, case in cat["cases"].items():
            inp = _decode_input(case["input"])
            exp = expected["categories"][cat_name]["cases"][case_name]
            try:
                produced = canonicalize(inp).hex()
            except (ValueError, TypeError) as exc:
                if "error" in exp:
                    rejected_match += 1
                else:
                    disagreements.append(
                        f"{cat_name}.{case_name}: inline rejected ({exc}) but library accepted"
                    )
                continue
            if "error" in exp:
                disagreements.append(
                    f"{cat_name}.{case_name}: inline accepted ({produced}) but library rejected"
                )
                continue
            if produced != exp["canonical_hex"]:
                disagreements.append(
                    f"{cat_name}.{case_name}:\n"
                    f"    inline:  {produced}\n"
                    f"    library: {exp['canonical_hex']}"
                )
            else:
                matched += 1

    print(f"Matched: {matched}")
    print(f"Rejection match: {rejected_match}")
    if disagreements:
        print(f"Disagreements: {len(disagreements)}")
        for d in disagreements:
            print(" -", d)
        return 1
    print("All cases agree between jcs library and inline reference.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
