# Case 010 — Tenant IKM rotation mid-day

## Purpose

The load-bearing rotation case. The same tenant's IKM rotates from `key_version=1` (with `ikm_v1`) to `key_version=2` (with `ikm_v2`). Events 1-3 are captured under `v1`; events 4-5 under `v2`. The same `run_id` spans the rotation. The chain must walk seamlessly: each entry's `key_version` and `key_fingerprint` resolve to the right IKM at verifier time.

This case exercises:

- The verifier's per-entry `(tenant_id, key_version) -> IKM` lookup
- The fingerprint check (step 8) succeeding for both halves
- The MAC compute using the right session key per half
- The seal record's `key_versions = [1, 2]` field
- The Merkle root over the day's full event set, regardless of which IKM signed each entry

## Inputs

From `../chain_vectors.json` `rotation_chain` array. Same tenant, same run_id. `ikm_v1` (key_version=1) for seq 1-3; `ikm_v2` (key_version=2) for seq 4-5.

## Expected outputs

| `seq` | `key_version` | `key_fingerprint_hex` | `payload_hash_hex` |
|---|---|---|---|
| 1 | 1 | `54c245a894b81d788d372badc563531f` | `e84168071562a9e866f28977a204b43abfdc959b7f84aad83c16fe790f63b049` |
| 2 | 1 | `54c245a894b81d788d372badc563531f` | `d9acb553eb1fc59911ebf04e78f9b97fd032c68b8b50bb9073ea4bff1a5de89a` |
| 3 | 1 | `54c245a894b81d788d372badc563531f` | `c6f45090813d8d81a8425e07dcc33318133539ddea9e53b0437060ad4ed732c6` |
| 4 | 2 | `1e05971d136fee9e6e9c1a96d2fca05d` | (computed under session_key_v2; see chain_vectors.json) |
| 5 | 2 | `1e05971d136fee9e6e9c1a96d2fca05d` | (computed under session_key_v2; see chain_vectors.json) |

Merkle root over the five payload_hashes: `88b968e7c106fb2dfe0d74825ec1442af5a7236d78e54d420d394c8166ade871`.

The seal record carries `key_versions = [1, 2]`.

## What this case proves

- The fingerprint check correctly distinguishes IKM generations: `54c245a894b81d788d372badc563531f` for v1, `1e05971d136fee9e6e9c1a96d2fca05d` for v2. The verifier asserts the looked-up IKM produces the recorded fingerprint at every entry.
- The MAC recompute uses the right session key per entry. A verifier that uses `session_key_v1` for seq=4 (wrong key for that key_version) gets a MAC mismatch and fails.
- The Merkle root is over the full day's events, not split per key_version. The seal covers everything.
- The seal record's `key_versions` field correctly captures the multi-version state.

## Failure mode

Most likely failure: the verifier doesn't switch IKMs at the rotation boundary and uses one IKM for the whole chain. This produces a fingerprint mismatch at `seq=4` (the looked-up `ikm_v1` does not produce the recorded `key_fingerprint` for v2). The verifier reports `key_fingerprint mismatch at seq 4` and the day fails. This is the correct behavior when the verifier's IKM lookup is wrong; the rotation case exercises the lookup-and-switch logic.

## Related auditor procedure

P-6 (key-fingerprint reconciliation) consumes operational events `master_key.rotated` and `master_key.rotation_observed` paired with this case's chain. The full audit evidence trail is:

1. KMS emits `master_key.rotated` when `ikm_v2` is provisioned
2. Ledger emits `master_key.rotation_observed` when it first sees an event under `key_version=2`
3. Weekly `master.reconciliation_completed` reports `key_versions_observed = [1, 2]` and `fingerprint_unmatched_count = 0`

The SOC team and FFIEC examiner cross-check these three events plus the chain entries plus the seal record's `key_versions` field. All four sources should agree.
