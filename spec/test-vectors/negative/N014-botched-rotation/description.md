# N014 — Botched rotation (`key_version=1` re-used for a different IKM, load-bearing rework case)

## Purpose

The load-bearing operational scenario the rework's per-entry `key_fingerprint` defense exists to catch. The verifier MUST detect that `key_version=1` is mapped to a different IKM than the one that originally signed the chain entries — at spec §7 step 8, BEFORE any MAC compute, with the precise reason that points the auditor at the IKM roster rather than the chain content.

This case is operationally significant: institutions occasionally reuse `key_version` values in their KMS/HSM after a deletion (deliberately or by mistake during a rotation rollout). Without the per-entry fingerprint check, the verifier would compute a MAC under the wrong IKM, get a mismatch, and report `payload_hash MAC mismatch` — leaving the auditor to investigate the chain content (no useful signal there) rather than the IKM roster (where the actual problem lives).

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs. The chain entry was originally signed with `ikm_v1` (`key_fingerprint = 54c245a894b81d788d372badc563531f`). Configure the verifier's IKM lookup such that `(tenant_id, key_version=1)` returns `ikm_v2` instead (the correct IKM exists in the registry under a different label, but the lookup table for `key_version=1` was incorrectly remapped to v2's bytes).

The simulation:
- Lookup: `ikm_lookup("tenant-ffiec-test-1", 1)` returns `ikm_v2` bytes.
- Verifier recomputes fingerprint: `SHA-256(utf8("tenant-ffiec-test-1") || ikm_v2)[:16] = 1e05971d136fee9e6e9c1a96d2fca05d`.
- Entry's recorded fingerprint: `54c245a894b81d788d372badc563531f` (the v1 fingerprint).
- Constant-time compare: FAIL.

## Expected verifier outcome

```
Status: FAIL
Step:   8
Reason: key_fingerprint mismatch at seq 1: looked-up IKM does not match the entry's recorded fingerprint
```

Specifically:

- The verifier MUST execute steps 1-7 in order, then fail at step 8.
- The verifier MUST NOT execute step 9 (no MAC compute happens — this is the load-bearing distinction; if the verifier ran the MAC compute, it would also fail, but the reason string would be wrong).
- The reason string MUST point the auditor at the IKM roster ("looked-up IKM does not match the entry's recorded fingerprint"), NOT at the chain content.

## What this case proves

A verifier that passes this case is robust against the most common operational failure mode the rework exists to detect: an IKM-roster misconfiguration. The auditor reading the verifier output knows to investigate the institution's IKM-roster `(tenant_id, key_version=1)` row, not to suspect chain tampering. The IR Scenario 7 (key_fingerprint mismatch) triage tree directly applies.

A verifier that reports `payload_hash MAC mismatch at seq 1` instead is non-conforming for the rework; it indicates the verifier ran step 9 before step 8, defeating the load-bearing defense.

## Negative case provenance

This case exists to lock in the spec §7 step ordering. Per spec §4.1 inviolate property #3: "Per-entry fingerprint, checked before MAC compute. The verifier asserts the looked-up IKM produces this fingerprint BEFORE computing any MAC." A regression that reorders steps 8 and 9 fails this case.
