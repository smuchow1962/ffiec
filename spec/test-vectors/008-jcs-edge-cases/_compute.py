# -*- coding: utf-8 -*-
"""
Compute JCS-canonical bytes for the FFIEC v1 test-vector case 008-jcs-edge-cases.

Uses the Python `jcs` package (RFC 8785 reference implementation by Anders Rundgren,
the RFC author). Reads `fixture.json`, produces `expected.json` next to it.

Each `expected.json` entry carries either:
  - "canonical_hex": the JCS-canonical bytes as lowercase hex
  - "canonical_utf8": the same bytes decoded as UTF-8 for readability
  - "error": the named error mode if RFC 8785 mandates rejection

The fixtures are organised by RFC 8785 edge-case category. See description.md.
"""

from __future__ import annotations

import hashlib
import json
import math
import os
import sys
from typing import Any

import jcs

HERE = os.path.dirname(os.path.abspath(__file__))


def _hash(b: bytes) -> str:
    return hashlib.sha256(b).hexdigest()


def _canonicalize_or_error(value: Any) -> dict:
    """Run jcs.canonicalize and capture either bytes or the rejection reason.

    Cognitive-complexity note: a single try/except keeps the dispatch simple.
    The library raises on NaN/Infinity inputs (RFC 8785 §3.2.2.3 / §4 — those
    values have no canonical form). We capture the exception class name so the
    fixture remains stable across patch versions of the library.
    """

    try:
        canonical = jcs.canonicalize(value)
    except (ValueError, TypeError) as exc:
        return {
            "error": type(exc).__name__,
            "error_message": str(exc),
        }

    return {
        "canonical_hex": canonical.hex(),
        "canonical_utf8": canonical.decode("utf-8"),
        "canonical_sha256": _hash(canonical),
        "canonical_length": len(canonical),
    }


def _build_fixtures() -> dict:
    """The fixture corpus, keyed by RFC 8785 edge-case category.

    Each category has a list of named cases. We keep the data in code (rather
    than a separately-edited fixture.json) so the script is the source of
    truth and the JSON fixture file is regenerated mechanically. The fixture
    file IS still committed alongside the script — it is what conforming
    implementations consume; the script just records how it was produced.
    """

    fixtures: dict[str, dict[str, Any]] = {
        "float_canonicalization": {
            "description": (
                "RFC 8785 §3.2.2 mandates the ECMAScript Number-to-String "
                "canonical form. Integers print without a decimal point; "
                "values whose absolute magnitude is < 1e-6 or >= 1e21 use "
                "exponent notation; otherwise the shortest round-trip form."
            ),
            "cases": {
                "integer_value_one": {"input": {"v": 1.0}},
                "small_decimal": {"input": {"v": 0.1}},
                "decimal_one_third_approx": {"input": {"v": 1.0 / 3.0}},
                "ten_billion": {"input": {"v": 1e10}},
                "exponent_threshold_high": {"input": {"v": 1e21}},
                "below_exponent_threshold": {"input": {"v": 1e20}},
                "exponent_threshold_low": {"input": {"v": 1e-6}},
                "below_low_threshold": {"input": {"v": 1e-7}},
                "very_small_finite": {"input": {"v": 5e-324}},
                "very_large_finite": {"input": {"v": 1.7976931348623157e308}},
                "negative_zero": {"input": {"v": -0.0}},
                "positive_zero": {"input": {"v": 0.0}},
                "trailing_zero_strip": {"input": {"v": 1.50}},
                "long_mantissa_e_minus_100": {"input": {"v": 1.234567890123e-100}},
            },
        },
        "non_ascii_unicode": {
            "description": (
                "RFC 8785 §3.2.2.2 + RFC 8259 §7: serialise as UTF-8 with the "
                "minimal escape set (only the JSON-mandatory escapes). All "
                "other Unicode characters serialise as raw UTF-8 bytes."
            ),
            "cases": {
                "latin1_cafe": {"input": {"v": "café"}},
                "cjk_japanese": {"input": {"v": "日本語"}},
                "emoji_lock": {"input": {"v": "\U0001f512"}},
                "combining_mark_nfc": {"input": {"v": "é"}},
                "combining_mark_nfd": {"input": {"v": "é"}},
                "mixed_script": {"input": {"v": "Hello 世界 \U0001f30d"}},
                "rtl_arabic": {"input": {"v": "مرحبا"}},
                "zero_width_joiner": {"input": {"v": "a‍b"}},
            },
        },
        "surrogate_pairs": {
            "description": (
                "Characters in the supplementary plane (U+10000..U+10FFFF). "
                "RFC 8785 emits them as raw UTF-8 (4 bytes per character), "
                "not as a `\\uD83D\\uDE00`-style escape pair."
            ),
            "cases": {
                "emoji_smile": {"input": {"v": "\U0001f600"}},
                "emoji_flag_japan": {"input": {"v": "\U0001f1ef\U0001f1f5"}},
                "math_double_struck_a": {"input": {"v": "\U0001d538"}},
                "musical_g_clef": {"input": {"v": "\U0001d11e"}},
                "supplementary_cjk": {"input": {"v": "\U00020000"}},
            },
        },
        "nan_infinity_rejection": {
            "description": (
                "RFC 8785 §3.2.2.3 + §4 require rejection of NaN, +Infinity, "
                "and -Infinity — these have no canonical JSON form. The "
                "expected outcome is an error, not a byte sequence."
            ),
            "cases": {
                "nan": {"input": {"v": float("nan")}},
                "positive_infinity": {"input": {"v": float("inf")}},
                "negative_infinity": {"input": {"v": float("-inf")}},
            },
        },
        "very_long_strings": {
            "description": (
                "Strings approaching parser memory limits. The fixture sizes "
                "are conservative (1KB and 64KB) so the corpus stays small "
                "while still exercising the long-string code path. "
                "Implementations should not impose a smaller limit."
            ),
            "cases": {
                "string_1kb_ascii": {"input": {"v": "a" * 1024}},
                "string_64kb_ascii": {"input": {"v": "a" * 65536}},
                "string_1kb_unicode": {"input": {"v": "é" * 1024}},
            },
        },
        "deeply_nested_objects": {
            "description": (
                "Stack-recursion limits. Cases at 10 / 50 / 100 levels. "
                "1000-level nesting is omitted because Python's default "
                "recursion limit and most JSON parsers reject far below that. "
                "Implementations supporting deeper nesting MAY add their own "
                "extension fixtures; v1.0 conformance is satisfied at 100."
            ),
            "cases": {
                "nested_10": {"input": None},
                "nested_50": {"input": None},
                "nested_100": {"input": None},
            },
        },
        "control_characters": {
            "description": (
                "RFC 8259 §7 mandates escaping for U+0000..U+001F. Specific "
                "characters use the short escape forms (\\b \\t \\n \\f \\r \\\" "
                "\\\\); all others use \\u00XX. U+007F (DEL) is NOT escaped per "
                "RFC 8785."
            ),
            "cases": {
                "u0000_null":          {"input": {"v": chr(0x0000)}},
                "u0001_soh":           {"input": {"v": chr(0x0001)}},
                "u0008_backspace":     {"input": {"v": chr(0x0008)}},
                "u0009_tab":           {"input": {"v": chr(0x0009)}},
                "u000a_newline":       {"input": {"v": chr(0x000A)}},
                "u000b_vtab":          {"input": {"v": chr(0x000B)}},
                "u000c_formfeed":      {"input": {"v": chr(0x000C)}},
                "u000d_cr":            {"input": {"v": chr(0x000D)}},
                "u001f_us":            {"input": {"v": chr(0x001F)}},
                "u007f_del":           {"input": {"v": chr(0x007F)}},
                "quote_in_string":     {"input": {"v": "a" + chr(0x22) + "b"}},
                "backslash_in_string": {"input": {"v": "a" + chr(0x5C) + "b"}},
                "forward_slash_not_escaped": {"input": {"v": "a/b"}},
            },
        },
        "object_key_ordering": {
            "description": (
                "RFC 8785 §3.2.3 specifies key ordering by UTF-16 code-unit "
                "lexicographic comparison. For BMP characters this matches "
                "Unicode codepoint sort. For supplementary-plane characters "
                "the two orderings can diverge — the cases below exercise "
                "both regimes."
            ),
            "cases": {
                "ascii_mixed_case_digits": {
                    "input": {"b": 1, "B": 2, "a": 3, "A": 4, "0": 5, "_": 6}
                },
                "unicode_keys_bmp": {
                    "input": {"é": 1, "e": 2, "z": 3, "à": 4}
                },
                "unicode_keys_with_supplementary": {
                    "input": {
                        "z": 1,
                        "a": 2,
                        "\U0001f600": 3,
                        "é": 4,
                    }
                },
                "numeric_string_keys": {
                    "input": {"10": 1, "2": 2, "1": 3, "20": 4}
                },
            },
        },
        "numeric_edge_cases": {
            "description": (
                "Boundary cases for integer-valued numbers. JCS treats all "
                "JSON numbers as IEEE-754 doubles per RFC 8785 §3.2.2.3, so "
                "integers above 2^53 may lose precision. Values inside the "
                "safe-integer range are exact."
            ),
            "cases": {
                "integer_zero": {"input": {"v": 0}},
                "integer_negative_zero_int": {"input": {"v": -0}},
                "integer_one": {"input": {"v": 1}},
                "integer_negative_one": {"input": {"v": -1}},
                "max_safe_integer": {"input": {"v": 9007199254740991}},
                "min_safe_integer": {"input": {"v": -9007199254740991}},
                "above_safe_integer": {"input": {"v": 9007199254740992}},
                "small_negative_decimal": {"input": {"v": -0.0001}},
            },
        },
    }

    # Build the deeply-nested objects programmatically.
    def _nest(depth: int) -> dict:
        node: Any = {"leaf": True}
        for _ in range(depth):
            node = {"n": node}
        return node

    fixtures["deeply_nested_objects"]["cases"]["nested_10"]["input"] = _nest(10)
    fixtures["deeply_nested_objects"]["cases"]["nested_50"]["input"] = _nest(50)
    fixtures["deeply_nested_objects"]["cases"]["nested_100"]["input"] = _nest(100)

    return fixtures


def _serialise_fixture_input(value: Any) -> Any:
    """Encode the input for fixture.json.

    NaN/Infinity are not representable in JSON, so we encode them as marker
    strings. The expected.json result records the rejection. Conforming
    implementations must decode the markers to native float values before
    feeding them to their JCS canonicaliser.
    """

    if isinstance(value, float):
        if math.isnan(value):
            return {"__special__": "NaN"}
        if math.isinf(value):
            return {"__special__": "Infinity" if value > 0 else "-Infinity"}
    if isinstance(value, dict):
        return {k: _serialise_fixture_input(v) for k, v in value.items()}
    if isinstance(value, list):
        return [_serialise_fixture_input(v) for v in value]
    return value


def _decode_fixture_input(value: Any) -> Any:
    """Inverse of _serialise_fixture_input."""

    if isinstance(value, dict):
        marker = value.get("__special__") if len(value) == 1 else None
        if marker == "NaN":
            return float("nan")
        if marker == "Infinity":
            return float("inf")
        if marker == "-Infinity":
            return float("-inf")
        return {k: _decode_fixture_input(v) for k, v in value.items()}
    if isinstance(value, list):
        return [_decode_fixture_input(v) for v in value]
    return value


def main() -> int:
    fixtures = _build_fixtures()

    fixture_out: dict[str, Any] = {
        "_about": (
            "FFIEC v1 test-vector case 008-jcs-edge-cases. "
            "Inputs for RFC 8785 JCS conformance. "
            "NaN and Infinity are encoded as {\"__special__\": \"NaN\"|"
            "\"Infinity\"|\"-Infinity\"} because JSON cannot carry the raw "
            "float values. Conforming implementations decode these markers "
            "to native non-finite floats before invoking their canonicaliser."
        ),
        "categories": {},
    }

    expected_out: dict[str, Any] = {
        "_about": (
            "Expected outputs for fixture.json. canonical_hex is the lowercase "
            "hex of the JCS-canonical bytes; canonical_utf8 is the same bytes "
            "decoded as UTF-8 for human inspection; canonical_sha256 is the "
            "SHA-256 of the canonical bytes; rejection cases carry an error "
            "field instead. See computation_method.md for provenance."
        ),
        "categories": {},
    }

    for cat_name, cat in fixtures.items():
        cat_in: dict[str, Any] = {
            "description": cat["description"],
            "cases": {},
        }
        cat_out: dict[str, Any] = {
            "description": cat["description"],
            "cases": {},
        }
        for case_name, case in cat["cases"].items():
            inp = case["input"]
            cat_in["cases"][case_name] = {
                "input": _serialise_fixture_input(inp),
            }
            cat_out["cases"][case_name] = _canonicalize_or_error(inp)
        fixture_out["categories"][cat_name] = cat_in
        expected_out["categories"][cat_name] = cat_out

    fixture_path = os.path.join(HERE, "fixture.json")
    expected_path = os.path.join(HERE, "expected.json")

    with open(fixture_path, "w", encoding="utf-8", newline="\n") as fh:
        json.dump(fixture_out, fh, ensure_ascii=False, indent=2, sort_keys=False)
        fh.write("\n")
    with open(expected_path, "w", encoding="utf-8", newline="\n") as fh:
        json.dump(expected_out, fh, ensure_ascii=False, indent=2, sort_keys=False)
        fh.write("\n")

    # Round-trip self-check: re-read fixture.json, decode, re-canonicalise,
    # confirm bytes match expected.json. Catches any silent encoding drift.
    with open(fixture_path, encoding="utf-8") as fh:
        fixture_loaded = json.load(fh)
    with open(expected_path, encoding="utf-8") as fh:
        expected_loaded = json.load(fh)

    mismatches: list[str] = []
    for cat_name, cat in fixture_loaded["categories"].items():
        for case_name, case in cat["cases"].items():
            inp = _decode_fixture_input(case["input"])
            try:
                produced = jcs.canonicalize(inp).hex()
            except (ValueError, TypeError) as exc:
                produced_err = type(exc).__name__
                exp_err = (
                    expected_loaded["categories"][cat_name]["cases"][case_name]
                    .get("error")
                )
                if produced_err != exp_err:
                    mismatches.append(
                        f"{cat_name}.{case_name}: error {produced_err} vs {exp_err}"
                    )
                continue
            exp_hex = (
                expected_loaded["categories"][cat_name]["cases"][case_name]
                .get("canonical_hex")
            )
            if produced != exp_hex:
                mismatches.append(
                    f"{cat_name}.{case_name}: bytes drift after round-trip"
                )

    if mismatches:
        print("ROUND-TRIP SELF-CHECK FAILED:", file=sys.stderr)
        for m in mismatches:
            print("  -", m, file=sys.stderr)
        return 1

    print(f"Wrote {fixture_path}")
    print(f"Wrote {expected_path}")
    print("Round-trip self-check: OK")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
