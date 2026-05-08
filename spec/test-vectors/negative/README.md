# Negative test vectors

> **What this directory is.** Test cases where the verifier MUST report failure with a specific named reason. A verifier that produces a "pass" result on any of these cases is broken.

Each case takes a passing-case input (typically `001-single-event-empty-prev` or `002-multi-event-same-run`), tampers with one specific element, and documents the verifier's expected failure mode and reason string.

## Cases

Per `docs/design/08-test-vectors.md` §5.5:

| Case | What's tampered | Spec §7 step | Expected reason string |
|---|---|---|---|
| `N001-payload-hash-bit-flip/` | One byte in `payload_hash` flipped | 9 | `payload_hash MAC mismatch at seq N` |
| `N002-events-reordered/` | Two events swapped (without re-chaining) | 6 | `chain link broken at seq N` |
| `N003-merkle-root-altered/` | Seal's `merkle_root` set to garbage | 10 | `merkle root mismatch — ledger contents do not produce sealed root` |
| `N004-signature-garbage/` | Seal's `signature` set to garbage | 11 | `signature verification failed` |
| `N005-signature-wrong-tenant/` | Seal signed for a different tenant | 11 | `signature verification failed` (tenant_id is in sign_payload) |
| `N006-key-fingerprint-flipped/` | Entry's `key_fingerprint` flipped to arbitrary 16 bytes | 8 | `key_fingerprint mismatch at seq N: looked-up IKM does not match the entry's recorded fingerprint` (NO MAC compute) |
| `N007-unknown-key-version/` | Entry's `key_version` set to a generation not in the test IKM registry | 7 | `unknown key_version: no IKM for (tenant=T, key_version=V) at seq N` (NO MAC compute) |
| `N008-entry-format-version-mismatch/` | One entry's `format_version` set to `"v2"` | 5 | `format_version mismatch at seq N` |
| `N009-header-format-version-v2/` | `header.format_version` set to `"v2"` | 1 | `format_version v2 not supported by this verifier (running v1)` |
| `N010-header-hkdf-digest-flipped/` | `header.hkdf_inputs_digest` flipped | 2 | `header HKDF inputs do not match running v1 inputs` |
| `N011-header-genesis-nonzero/` | `header.genesis_hash` set to non-zero bytes | 3 | `header genesis_hash does not match v1 constant` |
| `N012-cross-chain-tenant-mismatch/` | Event's `tenant_id` differs from header's | 4 | `cross-chain lift detected at seq N (event.tenant_id mismatch)` |
| `N013-mid-write-truncation/` | Audit file's last byte is not `\n` | (file pre-flight) | `audit file ends mid-line — possible mid-write crash` |
| `N014-botched-rotation/` | `key_version=1` re-used for a different IKM (same tenant) | 8 | `key_fingerprint mismatch at seq N` (load-bearing rotation defence) |
| `N015-prev-hash-substituted/` | Entry's `prev_hash` substituted, `payload_hash` left alone | 6 | `chain link broken at seq N` (structural walk catches it) |
| `N016-prev-hash-and-payload-recomputed-by-attacker/` | Attacker substitutes `prev_hash` AND recomputes `payload_hash` (without IKM) | 6 | `chain link broken at seq N` (structural walk catches the substitution; even if the structural check were relaxed, the MAC step would catch the attacker's MAC because it would not match the verifier's MAC computed with the real IKM via `expected_prev_hash`) |
| `N022-format-version-v1-1/` | `header.format_version` set to `"v1.1"` (unrecognized minor within v1 family) | 1 | `format_version v1.1 not supported by this verifier (running v1)` |
| `N023-format-version-case-variant/` | `header.format_version` set to `"V1"` (uppercase V; case-variant of the recognized lowercase string) | 1 | `format_version "V1" not supported by this verifier (running v1)` |

## Conformance test

For each case, the verifier MUST:

1. Load the input
2. Run the spec §7 procedure
3. Report failure
4. The failure reason MUST match the expected reason string (or a syntactically equivalent variant — the substantive parts must be identical)
5. The failure step MUST match the expected step number

A verifier that returns `ok=True` on any case is non-conforming. A verifier that fails for the wrong reason (e.g. reports a MAC mismatch when the expected failure is a fingerprint mismatch) is non-conforming and indicates that the verifier ran step 9 before step 8 — the order is load-bearing per spec §7.

## Implementation note

Each case directory's `description.md` documents the exact tampering recipe. The reference Go runner produces the tampered inputs from the passing-case fixtures and asserts the expected failure shape.

The MOST IMPORTANT cases for the rework are:

- **N006** — flipped `key_fingerprint`. Catches the post-rework defense that the verifier checks the fingerprint BEFORE computing any MAC. A verifier that runs MAC compute first will pass N006 incorrectly (the MAC compute happens to fail because the IKM is wrong, but the failure reason will be `payload_hash MAC mismatch`, not `key_fingerprint mismatch`).
- **N014** — botched rotation. The load-bearing case for the per-tenant `key_fingerprint` defence. Catches the operational mistake that a malicious or careless operator could otherwise mask.
- **N013** — mid-write truncation. The verifier rejects rather than silently passes a chain that lost its last entry.

These three are the rework's signature contributions. A verifier that passes the v1.0-pre-rework corpus but fails N006/N013/N014 is non-conforming for the v1.0-rework spec.
