# N019 — Co-signed seal: one signature valid + one invalid (case e — load-bearing)

## Status

**Stub case.** Byte-level fixture deferred to v1.x. Recipe and expected verifier outcome documented here.

## Purpose

Verify spec §7 step 11 case (e): both signatures present in `signatures` list, one validates + the other does not. This is the load-bearing case: a co-signed seal failure indicates EITHER (i) one of the algorithms has been broken, OR (ii) one of the per-algorithm signing keys has been compromised. The un-broken algorithm's signature still provides integrity assurance under that algorithm's assumption.

The verifier MUST report BOTH the valid-algorithm validation AND the invalid-algorithm failure so the regulator and the institution have the full picture.

## Tampering recipe

Start from case 015. Tamper with the Dilithium3 signature (flip 1+ bytes):

```
signatures:
  - algorithm: ed25519
    signature: <valid Ed25519 signature>
  - algorithm: dilithium3
    signature: <Dilithium3 signature with one byte flipped>  # invalid
```

The Ed25519 signature still validates against the Ed25519 public key over the Ed25519-bound sign_payload. The Dilithium3 signature fails validation.

## Expected verifier outcome

```
Status: --strict: FAIL
        non-strict: PASS-WITH-ANOMALY
Reason: co-signed seal failure: algorithm ed25519 validated, algorithm dilithium3 did not
per_algorithm_results:
  - ed25519: PASS
  - dilithium3: FAIL
```

## What this case proves

The verifier correctly:
1. Reconstructs each algorithm's `sign_payload` independently (Variant B per spec §4.2).
2. Validates each signature against its algorithm-specific public key.
3. Reports per-algorithm results so the auditor sees both the validated half and the failed half.
4. Does NOT silently pass on one-valid-one-invalid — under `--strict` this is FAIL with a precise reason.

## Institution-side IR disposition

The institution's IR program (per IR Scenario 11 added in round 13 close-out) operates a three-branch clock-start determination:
- (i) Attributable to published algorithm break → no 36-hour clock; regulator coordination on migration timeline.
- (ii) Attributable to per-algorithm signing-key compromise → IR Scenario 4 (master-key compromise) starts the 36-hour clock.
- (iii) Under investigation → clock starts at investigation-conclusion determination.

## Negative case provenance

This case exists to lock in Variant B per-algorithm sign_payload encoding. A verifier that uses Variant A (single shared sign_payload) would mistakenly validate the Ed25519 signature against the Dilithium3 algorithm name (the algorithm-confusion attack the spec explicitly closes). Variant B's per-algorithm sign_payload binding makes that attack impossible by construction.
