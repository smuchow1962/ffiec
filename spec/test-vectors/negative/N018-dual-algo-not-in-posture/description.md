# N018 — Seal carrying algorithm not on institution's declared posture list (case c)

## Status

**Stub case.** Byte-level fixture deferred to v1.x. Recipe and expected verifier outcome documented here.

## Purpose

Verify spec §7 step 11 case (c): the seal carries an algorithm that is not on the institution's declared algorithm-posture list. Under `--strict`: FAIL. Under non-strict: PASS-WITH-ANOMALY.

## Tampering recipe

Start from case 015. Modify the seal record's `signatures` to include a third algorithm not on the declared posture:

```
signatures:
  - algorithm: ed25519
    signature: <valid Ed25519 signature>
  - algorithm: dilithium3
    signature: <valid Dilithium3 signature>
  - algorithm: slh_dsa_shake_128s
    signature: <valid SLH-DSA signature>  # NOT on declared posture
```

The institution's declared algorithm posture is `[ed25519, dilithium3]`. The seal carries `slh_dsa_shake_128s` in addition.

## Expected verifier outcome

```
Status: --strict: FAIL
        non-strict: PASS-WITH-ANOMALY
Reason: algorithm not on institution's declared posture list at seal_date 2026-05-06
unknown_algorithms: [slh_dsa_shake_128s]
```

## What this case proves

The verifier cannot accept signatures from algorithms the institution has not declared in its posture configuration. The institution's IR program investigates why an unexpected algorithm appeared (SDK misconfiguration? unauthorized signing key?). The disposition under `--strict` is FAIL because production-grade examination requires the institution's declared posture to be authoritative.
