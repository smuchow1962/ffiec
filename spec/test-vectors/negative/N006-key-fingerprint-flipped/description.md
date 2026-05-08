# N006 — `key_fingerprint` flipped (load-bearing rework case)

## Purpose

The load-bearing negative case for the v1.0-rework's defining defensive primitive. The verifier MUST detect the flipped `key_fingerprint` at spec §7 step 8, BEFORE any MAC compute, and report the failure with the precise named reason. A verifier that runs MAC compute first (and reports `payload_hash MAC mismatch` instead of `key_fingerprint mismatch`) is non-conforming for the rework: the spec ordering is load-bearing because the fingerprint check is the cheap rejection, the MAC compute is the expensive one, and the reason-string distinction is what tells the auditor where to investigate (IKM roster, not chain content).

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs. Modify the chain entry's `key_fingerprint` field to an arbitrary 16 bytes that differ from the correct value. Specifically:

```
correct:   key_fingerprint = 54c245a894b81d788d372badc563531f
tampered:  key_fingerprint = aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
```

Leave every other field intact, including `payload_hash` (which is still the valid HMAC under the correct IKM and session_key). The verifier's IKM lookup at step 7 returns the correct IKM; the fingerprint recompute at step 8 produces `54c245a894b81d788d372badc563531f`; the entry's stamped fingerprint is `aaaaa...`; constant-time compare fails.

## Expected verifier outcome

```
Status: FAIL
Step:   8
Reason: key_fingerprint mismatch at seq 1: looked-up IKM does not match the entry's recorded fingerprint
```

Specifically:

- The verifier MUST execute steps 1, 2, 3, 4, 5, 6, 7 (in order) before reaching step 8.
- The verifier MUST NOT execute step 9 (no MAC compute happens after fingerprint mismatch).
- The reason string MUST identify the failure as a fingerprint mismatch (NOT a MAC mismatch); the auditor's investigation path is the IKM roster, not the chain content.

## What this case proves

A verifier that passes this case — by reporting `key_fingerprint mismatch at seq 1` and stopping at step 8 — has the right ordering of cryptographic operations to defend the load-bearing rework property: a botched rotation that re-uses `key_version=1` with a different IKM is detected at lookup time with a precise message, not buried in a MAC-mismatch storm. A verifier that reports `payload_hash MAC mismatch at seq 1` (because it skipped or reordered step 8) is non-conforming.

## Negative case provenance

This case exists to prevent the regression where a future maintainer reorders steps 8 and 9 for "performance reasons." The fingerprint check costs ~100 nanoseconds (one SHA-256 + constant-time compare); the MAC compute costs ~1 microsecond (HKDF + HMAC-SHA-256). The ordering is correct as specified.
