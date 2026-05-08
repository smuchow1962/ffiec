# N020 — Algorithm / key-type mismatch at signature verification (spec §7 step 11)

## Purpose

Verify the spec §7 step 11 `algorithm/key-type mismatch` reason string is distinct from the generic `signature verification failed`. The verifier MUST report the more-specific reason when the seal record's `algorithm` field does not match the algorithm of the resolved public key (e.g. seal claims `"ed25519"` but `public_key_id` resolves to a Dilithium key).

A verifier that reports the generic `signature verification failed` for this case is non-conforming for the rework — the auditor needs the specific reason to investigate the public-key registry, not the chain content.

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs (v1.0 single-algorithm, computable today).

Tamper with the seal record's `algorithm` field OR the public-key registry resolution: configure the verifier such that `seal.algorithm = "ed25519"` (matches the actual Ed25519 signature) BUT the resolved public key for `seal.public_key_id` is of a different type (e.g. an RSA-3072 key, OR a v1.x post-quantum key). The signature was produced by a real Ed25519 key, so the cryptographic Ed25519.Verify operation against an RSA public key will refuse at the API level OR produce a signature-verification-failed result with the wrong-algorithm signal.

The simulation uses any of:
- The verifier's `public_key_id` registry contains an RSA-3072 key under the label `seal.public_key_id`.
- The institution's tenant key registry has been corrupted to point at a different algorithm's public key.
- A misconfigured registry that swapped Ed25519 and Dilithium key labels.

## Expected verifier outcome

```
Status: FAIL
Step:   11
Reason: algorithm/key-type mismatch at signature verification
        (seal claims algorithm "ed25519"; resolved public key for public_key_id
         is type "rsa-3072")
```

NOT the generic:

```
Reason: signature verification failed
```

The verifier MUST detect the algorithm/key-type mismatch at public-key resolution time (BEFORE the underlying cryptographic Verify call) AND report the more-specific reason. Implementations that pass the public key to the verify call without checking algorithm-type compatibility will produce the generic failure message; that is non-conforming.

## What this case proves

The verifier's algorithm-dispatch logic correctly distinguishes:

- Algorithm signature is invalid (cryptographic failure under the correct algorithm) → `signature verification failed`
- Algorithm field and resolved public key disagree (registry/configuration error) → `algorithm/key-type mismatch at signature verification`

The auditor reading the verifier output knows from the reason string which investigation path applies: `signature verification failed` → investigate signing-key compromise + chain content. `algorithm/key-type mismatch` → investigate public-key registry + tenant-key-publication discipline.

## Negative case provenance

This case fills the gap noted in cryptographic-engineer review round 13: spec §7 step 11 documents the distinction between the generic and specific reason strings, but no negative case exercised the distinction. Round 13 added the case; the recipe is computable under v1.0 single-algorithm posture and remains computable under v1.x dual-algorithm posture.
