# Case 021 — master.rotation.completed event payload (§10.28)

## What this case verifies

Spec §10.28 normates the `master.rotation.completed` operational event payload that the §10.29 streaming-mode verifier consumes. The payload's five fields are byte-locked across implementations; this case pins the JCS-canonical (RFC 8785) bytes for one rotation event so a clean-room implementation proves it produces the same bytes for the same inputs.

The case focuses on three byte-level invariants:

1. **Field order is canonicalisation-driven, not author-driven.** RFC 8785 sorts object keys lexicographically. The Python source dict in `_compute.py` is written in spec-doc order (`event`, `cadence`, `prior_key_version`, `new_key_version`, `rotation_at_utc`), but the canonical bytes emerge sorted (`cadence`, `event`, `new_key_version`, `prior_key_version`, `rotation_at_utc`).
2. **Microsecond precision is load-bearing.** `rotation_at_utc` MUST be formatted as `YYYY-MM-DDTHH:MM:SS.uuuuuuZ` with exactly 6 digits of microsecond precision and the literal trailing `Z`. Whole-second-resolution timestamps would erase the crossing-interval signal at sub-second cadence.
3. **Integer encoding has no quotes, no leading zeros, no decimal point.** `prior_key_version` and `new_key_version` serialize as bare JSON integers — not strings, not floats — per RFC 8785 §3.2.2.

## Inputs

| Field | Value |
|---|---|
| `event` | `master.rotation.completed` |
| `cadence` | `per_second` (§10.27 streaming-mode) |
| `prior_key_version` | `1` |
| `new_key_version` | `2` |
| `rotation_at_utc` | `2026-05-07T12:34:56.789012Z` (6-digit microsec, trailing Z) |

The microsecond value `.789012` is intentionally non-zero. An implementation that silently truncates to whole-second precision (a common bug under naive `datetime` formatting) produces a byte-different `rotation_at_utc` value and fails the SHA-256 pin loudly.

`prior_key_version` and `new_key_version` MUST differ (per §10.28); the values `1` and `2` are the simplest valid pair.

## Expected canonical bytes

The expected JCS-canonical bytes are captured in `expected_canonical.txt` (the bytes verbatim, no trailing newline) and `expected_canonical_sha256.txt` (the SHA-256 of those bytes, for quick verification without comparing 150 bytes by eye).

| Property | Value |
|---|---|
| canonical byte length | `150` |
| canonical SHA-256 | `be6adcb49bf4ce0fc9aa2d14b77a440e884f787a2ef09e79cc4243b6faee2270` |
| canonical text | `{"cadence":"per_second","event":"master.rotation.completed","new_key_version":2,"prior_key_version":1,"rotation_at_utc":"2026-05-07T12:34:56.789012Z"}` |

## Conformance behavior

A conforming implementation:

1. Constructs the event payload as a JSON object with the §10.28 five fields.
2. Canonicalises it per RFC 8785 (JCS) and produces bytes byte-identical to `expected_canonical.txt`.
3. Computes a SHA-256 of the canonical bytes that matches `expected_canonical_sha256.txt`.
4. Refuses any payload whose `rotation_at_utc` is not formatted with exactly 6 digits of microsecond precision plus trailing `Z`.
5. Refuses any payload whose `prior_key_version` equals `new_key_version` (§10.28 normative).

## Negative cases this fixture supports

- A producer that emits `rotation_at_utc` as `"2026-05-07T12:34:56Z"` (whole-second resolution) produces canonical bytes that diverge from `expected_canonical.txt`.
- A producer that emits `rotation_at_utc` as `"2026-05-07T12:34:56.789012+00:00"` (RFC 3339 numeric offset rather than `Z`) produces canonical bytes that diverge.
- A producer that emits `prior_key_version` as the JSON string `"1"` rather than the bare integer `1` produces canonical bytes that diverge (per JCS, integers and strings have distinct canonical forms).
- A producer that omits any of the five fields, or adds a sixth field, produces different canonical bytes by definition.

## Cross-references

- Spec §10.28 streaming-mode IKM rotation discipline (the section that normates this event payload)
- Spec §10.2 operational-events catalog (where `master.rotation.completed` is registered)
- Spec §10.29 streaming-mode verifier procedure (consumes these events as `rotation` inputs)
- RFC 8785 (the JCS canonicalisation rules this case applies)
- RFC 3339 (the timestamp form §10.28 mandates for `rotation_at_utc`)
- Case 008 (`008-jcs-edge-cases/`) — the JCS conformance bar; if 008 fails, 021 will fail in the same way

## Reproduction

The `_compute.py` script in this directory regenerates `input.json`, `expected_canonical.txt`, and `expected_canonical_sha256.txt` from the pinned inputs above. It depends on the Python `jcs` package (Anders Rundgren's reference RFC 8785 implementation, the same package case 008 uses).

```
python _compute.py
```
