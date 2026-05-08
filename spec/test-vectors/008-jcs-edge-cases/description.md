# Case 008 — JCS edge cases (RFC 8785 conformance)

## Purpose

Spec §5 says the canonical-JSON form used inside `payload_hash` MUST follow RFC 8785, and that two implementations of v1 MUST produce byte-identical canonical bytes for the same logical event. The basic fixtures (`001-` through `015-`) exercise integers, ASCII strings, and base64-encoded byte strings. They do not exercise the parts of RFC 8785 that historically cause silent disagreement between implementations: floats, non-ASCII Unicode, NaN/Infinity rejection, surrogate pairs, deeply nested objects, control characters, and non-trivial key ordering.

This case closes that gap. Two implementations passing every other fixture in the corpus can still disagree on the first event whose attributes contain a float or a non-ASCII string. That disagreement breaks chain integrity silently — the writer's `payload_hash` is correct under its canonicaliser; the verifier's recompute mismatches; the verifier reports `payload_hash MAC mismatch`; the cause is misattributed to a key or chain bug. We treat this as a v1.0 conformance failure and gate it here.

## Conformance bar

An implementation is v1.0-conformant only if it reproduces every byte sequence in `expected.json` for the corresponding input in `fixture.json`, AND rejects every input in the `nan_infinity_rejection` category with an error. Implementations that pass `001-` through `015-` but fail any case in `008-` are NOT v1.0-conformant.

## Categories exercised

The fixture is organised by RFC 8785 edge-case category. Each category targets one specific class of disagreement we have seen between independent JCS implementations.

### Float canonicalization

RFC 8785 §3.2.2.3 mandates the ECMAScript Number-to-String form. The rules that catch implementations:

- Integer-valued doubles (`1.0`) print without a decimal (`1`), not `1.0`.
- Exponent notation kicks in when `|n| < 1e-6` or `|n| >= 1e21`.
- Inside the fixed-form range, the output is the shortest round-trip decimal — not Python's default `repr`, not Java's `Double.toString`, not Go's `%g`.
- Negative zero serialises as `0`, not `-0`.
- Trailing zeros after the decimal point are stripped (`1.50` -> `1.5`).
- The exponent uses lowercase `e`, an explicit sign, and no leading zeros in the exponent value.

Cases: integer-valued double, small decimal, repeating-decimal approximation, ten billion, the high exponent threshold (`1e21` and `1e20`), the low exponent threshold (`1e-6` and `1e-7`), the smallest positive double (`5e-324`), the largest finite double (`1.7976931348623157e308`), positive and negative zero, trailing-zero strip, long mantissa with `e-100` exponent.

### Non-ASCII Unicode strings

RFC 8785 §3.2.2.2 + RFC 8259 §7 together: serialise as UTF-8 with the minimal escape set. Only the JSON-mandatory escapes are emitted (`\b \t \n \f \r \" \\` and `\u00XX` for U+0000..U+001F). Every other character serialises as raw UTF-8 bytes. Implementations that follow the looser RFC 8259 latitude and emit `\uXXXX` for non-ASCII characters produce different bytes from a conforming JCS canonicaliser.

Cases: Latin-1 (`café`), CJK ideographs (`日本語`), an emoji (`🔒`), the same character in NFC and NFD form (`é` as a single codepoint vs `e` + combining acute), mixed scripts, RTL Arabic, zero-width joiner.

The NFC vs NFD case is deliberate: RFC 8785 explicitly does NOT normalise. The two forms produce different canonical bytes, and a v1.0 implementation MUST preserve that distinction. Application-layer normalisation is the application's choice and happens before canonicalisation.

### Surrogate pairs

Characters in the supplementary plane (U+10000..U+10FFFF) are encoded as 4-byte UTF-8 sequences in the canonical output, NOT as a `😀`-style surrogate pair escape. Cases: emoji, the Japanese flag (which is a regional-indicator surrogate pair under the hood), mathematical alphanumeric symbols, musical symbols, supplementary CJK.

### NaN and Infinity rejection

RFC 8785 §3.2.2.3 requires that NaN, +Infinity, and -Infinity be rejected — they have no canonical JSON form. The expected outcome is an error, not a byte sequence. The fixture encodes these inputs as `{"__special__": "NaN"}`, `{"__special__": "Infinity"}`, and `{"__special__": "-Infinity"}` because JSON cannot carry the raw float values; conforming implementations decode the markers to native non-finite floats before invoking their canonicaliser. The expected.json records `error` and the error class name (`ValueError` for the Python reference).

### Very long strings

Strings near typical parser limits exercise the long-string code path. The fixture uses 1KB ASCII, 64KB ASCII, and 1KB Unicode. Larger sizes (1MB, 10MB) are not committed to the fixture for repository-size reasons; the `canonical_sha256` of those sizes can be added as out-of-band conformance checks if a future implementation reports problems.

### Deeply nested objects

Stack-recursion limits. The fixture nests at 10, 50, and 100 levels. 1000-level nesting is omitted because Python's default recursion limit and most production JSON parsers reject far below that level. Implementations supporting deeper nesting may add their own extension fixtures; v1.0 conformance is satisfied at 100 levels.

### Control characters

RFC 8259 §7 mandates escaping for U+0000..U+001F. Specific characters use the short escape forms (`\b \t \n \f \r \" \\`); all others use `\u00XX` with lowercase hex digits per RFC 8785 §3.2.2.2. U+007F (DEL) is NOT escaped — it sits above U+001F and JCS preserves it raw. The forward slash (`/`) is NOT escaped under JCS even though RFC 8259 permits the escape — the canonical form forbids unnecessary escapes.

Cases: each of the named short-escape characters, U+0000 (the parser-killer), U+001F (the boundary), U+007F (DEL — must NOT be escaped), embedded quote, embedded backslash, embedded forward slash.

### Object key ordering

RFC 8785 §3.2.3 specifies key ordering by UTF-16 code-unit lexicographic comparison. For BMP-only keys this matches Unicode codepoint sort. For supplementary-plane keys the two orderings can diverge — the supplementary-plane case below specifically exercises that divergence. Implementations that sort by UTF-8 bytes, by Python codepoint, or by locale-aware comparison fail this case.

Cases: ASCII mixed-case + digits + underscore (verifies digits sort before uppercase, uppercase before lowercase, and underscore sorts after uppercase), BMP Unicode keys, keys mixing BMP and supplementary plane (the test that catches UTF-8-byte-sort and codepoint-sort implementations), numeric-string keys (the trap that `"10"` sorts before `"2"` lexicographically).

### Numeric edge cases

Integer boundary cases that interact with the IEEE-754-double semantic in RFC 8785 §3.2.2.3. JavaScript's `Number.MAX_SAFE_INTEGER` is `2^53 - 1`; integers above that lose precision when round-tripped through a double. The fixture exercises both ends of the safe-integer range and one value just above it, plus zero, signed zero, and a small negative decimal.

## Inputs and expected outputs

`fixture.json` carries the inputs, organised by category. `expected.json` carries the expected outputs (canonical bytes as lowercase hex, plus a UTF-8 decode for human inspection, plus the SHA-256 of the canonical bytes for quick byte-equality checks). For the rejection cases, `expected.json` records `error` and the error class name instead of bytes.

The `canonical_sha256` field on each expected case is a fast equality check: an implementation can compute its own canonical bytes, take the SHA-256, and compare against the published hash before comparing the full hex.

## Computing party

`expected.json` was produced by the Python `jcs` package (version 0.2.1, by Anders Rundgren — the RFC 8785 author), and cross-validated against an independent inline implementation written from RFC 8785 directly. See `computation_method.md` for the full provenance.

## Failure modes

If an implementation disagrees on a case, suspect (in order):

1. **Float case**: implementation uses its language's default float-to-string rather than the ECMAScript shortest-round-trip form. Python `repr`, Java `Double.toString`, and Go `strconv.FormatFloat` all need post-processing.
2. **Non-ASCII case**: implementation escapes non-ASCII as `\uXXXX` instead of emitting raw UTF-8. Common in libraries that target RFC 8259 generally rather than RFC 8785 specifically.
3. **NFC vs NFD case**: implementation normalises Unicode before canonicalising. RFC 8785 forbids this — normalisation is the application's responsibility.
4. **NaN/Infinity case**: implementation emits `null` or `"NaN"` instead of rejecting. Python's stdlib `json.dumps(float('nan'))` produces `NaN` (invalid JSON) by default; a JCS canonicaliser MUST reject.
5. **Key ordering case**: implementation sorts by UTF-8 bytes or by locale. JCS sorts by UTF-16 code units.
6. **Control character case**: implementation uses uppercase hex (``) instead of lowercase (``), or escapes U+007F (DEL) which JCS does not.
7. **Forward slash case**: implementation escapes `/` as `\/`. JCS does not — the JSON spec permits the escape but the canonical form forbids unnecessary escapes.

## Why this case is part of the conformance bar

The whole point of binding `payload_hash` to a canonical form is that two independent verifiers — a regulator's verifier and the institution's own verifier — recompute the same bytes from the same logical event and reach the same conclusion. If the conformance corpus only exercises ASCII and integers, that property holds for ASCII-and-integers chains and silently fails for everything else. Audit events in production carry float values (durations in seconds), Unicode strings (user names, prompts), and control characters (in arbitrary attribute payloads). This case turns "v1.0-conformant" into a property an implementation can actually claim with confidence.
