# N017 — Partial-coverage seal under declared dual-algorithm posture (case b)

## Status

**Stub case.** Byte-level fixture deferred to v1.x. Recipe and expected verifier outcome documented here.

## Purpose

Verify spec §7 step 11 case (b): the institution operates dual-algorithm posture but a seal carries only one algorithm's signature. The verifier reports PASS-WITH-ANOMALY (control-completeness, NOT chain-integrity); under `--strict` STILL PASS-WITH-ANOMALY (the seal is integrity-bearing under the present algorithm; the institution's posture commitment is incomplete on this seal-day).

## Tampering recipe

Start from case 015 (dual-algorithm co-signed seal). Modify the seal record to remove the Dilithium3 entry from `signatures`:

```
signatures:
  - algorithm: ed25519
    signature: <Ed25519 signature over Ed25519-bound sign_payload>
  # dilithium3 entry removed
```

The institution's declared algorithm posture (out-of-band configuration the verifier consults) is `[ed25519, dilithium3]`. The seal carries only `ed25519`.

## Expected verifier outcome

```
Status: PASS-WITH-ANOMALY (regardless of --strict)
Reason: partial-coverage seal: single-algorithm signature during institution's declared dual-algorithm posture
per_algorithm_results:
  - ed25519: PASS
  - (dilithium3: NOT PRESENT)
```

## What this case proves

The verifier correctly classifies a partial-coverage seal as a control-completeness anomaly, NOT a chain-integrity failure. The institution's IR program treats this as a posture-completeness investigation (why did the SDK not co-sign? was the institution's PQ-signing key unavailable on this seal-day?), not as a tampering investigation.
