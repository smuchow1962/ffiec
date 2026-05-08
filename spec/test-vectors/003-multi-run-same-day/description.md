# Case 003 — Multi-run, same tenant-day

## Purpose

Two runs within the same tenant-day. Verifies cross-run chain isolation: each run's chain starts with `prev_hash = 32 zero bytes`, and the daily Merkle root aggregates events from both runs in `(run_id, seq)` ordering per spec §4.2.

The case answers a load-bearing operational question: when a tenant captures multiple runs in the same UTC day, does the chain link cross run boundaries (no — each run is a fresh chain) and does the Merkle seal cover everything (yes — one seal per tenant-day, regardless of run count)?

A verifier that walks the chain naively and treats run B's first event as continuing run A's chain produces a chain-link error at the seam; a verifier that skips run B from the daily Merkle entirely produces a Merkle-root mismatch when the seal is checked. Both failures are caught by this fixture.

## Inputs

```
HKDF_SALT      = b"ffiec.chain-of-custody.v1.salt"
HKDF_INFO_BASE = b"ffiec.chain-of-custody.v1.info"
tenant_id      = "tenant-ffiec-test-1"
run_id_a       = "RUN-FFIEC-VECTOR-A"
run_id_b       = "RUN-FFIEC-VECTOR-B"
ikm_v1         = bytes from ikm_v1_hex (mirrors the central corpus)
seal_date      = "2026-05-06"
algorithm      = "ed25519"
format_version = "v1"
cadence        = "daily"
dev_mode       = false
```

Run A has 3 events under `key_version=1`. Run B has 3 events under `key_version=1`. Both runs use the same tenant IKM, so both derive the same `session_key_v1` and the same `key_fingerprint_v1`. The full byte-level inputs are pinned in `fixture.json`.

## Expected outputs

| Run | `seq` | `prev_hash_hex` | `payload_hash_hex` |
|---|---|---|---|
| A | 1 | `0000...0000` | (pinned in `fixture.json` `run_a_chain`) |
| A | 2 | run A seq 1 `payload_hash` | pinned |
| A | 3 | run A seq 2 `payload_hash` | pinned |
| B | 1 | `0000...0000` (NOT run A's last) | pinned |
| B | 2 | run B seq 1 `payload_hash` | pinned |
| B | 3 | run B seq 2 `payload_hash` | pinned |

The daily Merkle root is computed over the 6 `payload_hash` values in `(run_id, seq)` ASC order: `RUN-FFIEC-VECTOR-A` events 1, 2, 3 followed by `RUN-FFIEC-VECTOR-B` events 1, 2, 3. The exact root is pinned in `expected.json` as `merkle_root_daily_hex`.

The `sign_payload` text (post-2026-05-07 spec §4.3 form) and its hex encoding are pinned in `expected.json`. The seal binds `cadence = "daily"` and `dev_mode = "0"`.

## What this case proves

- **Cross-run chain isolation.** Run B's first event has `prev_hash = 32 zero bytes`, NOT run A's last `payload_hash`. The chain MAC for that event is computed against the genesis prev — confirming chains never link across run boundaries even when both runs share a tenant-day and an IKM.
- **Merkle daily aggregation.** The daily Merkle root sees all 6 events from both runs. A verifier that processes only one run produces a different root and the seal check fails.
- **`(run_id, seq)` ordering.** The leaves enter the Merkle tree sorted lexicographically by `run_id` ASC, then `seq` ASC. `RUN-FFIEC-VECTOR-A` precedes `RUN-FFIEC-VECTOR-B`. A verifier that uses different ordering (e.g., timestamp ASC, or per-run subtrees) computes a different root.
- **Same IKM, two chains.** Both runs derive the same `session_key_v1` from the same `(ikm_v1, tenant_id)` HKDF inputs. The session-key cache is per-tenant, not per-run.

## Conformance test

A v1 verifier consuming this fixture MUST:

1. Derive `session_key_v1` and `key_fingerprint_v1` for the tenant. Both must match `expected.json`.
2. Derive `hkdf_inputs_digest`. Must match `expected.json`.
3. For run A and run B independently, walk events in `seq` order. Each `prev_hash` (after `seq=1`) must equal the previous event's `payload_hash`. Each `payload_hash = HMAC-SHA-256(session_key_v1, prev_hash || event_canonical)` must match the pinned hex byte-for-byte. Run B's `seq=1` `prev_hash` MUST be 32 zero bytes.
4. Concatenate the run-A `payload_hash` list with the run-B `payload_hash` list in that order (lexicographic by `run_id` ASC). Compute the RFC 6962 Merkle root. Must match `merkle_root_daily_hex`.
5. Construct `sign_payload` per spec §4.3 with `cadence="daily"` and `dev_mode=false`. Must match `sign_payload_hex` byte-for-byte.

## Failure modes

If the verifier reports a chain-link error at run B `seq=1`, it incorrectly carried run A's last `payload_hash` into run B; fix the run-boundary handling to reset `prev_hash` to genesis at the start of each run.

If the daily Merkle root mismatches, suspect the leaf ordering — the spec requires `run_id` ASC, then `seq` ASC, not insertion order or timestamp order.

If the `sign_payload` mismatches but the merkle root matches, suspect the post-2026-05-07 extension: the seal now binds `cadence` and `dev_mode` lines (no trailing newline after `dev_mode`).

## Provenance

`fixture.json` and `expected.json` are produced by `_compute.py` (HKDF-SHA-256 + HMAC-SHA-256 + RFC 6962 Merkle + sorted-keys JSON canonicalisation). Re-running the script reproduces the bytes deterministically. Run with `python _compute.py`.
