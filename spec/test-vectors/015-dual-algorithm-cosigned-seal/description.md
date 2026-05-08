# Case 015 — Dual-algorithm co-signed seal (v1.x post-quantum coexistence)

## Status

**Stub case.** Byte-level fixture deferred to the v1.x ship that introduces the second algorithm. v1.0 ships single-algorithm Ed25519 only; the structural recipe and expected verifier behavior are documented here so the v1.x implementer can populate the byte-level fixture from the same recipe.

## Purpose

Verify the dual-algorithm co-signed seal per spec §4.2 `signatures` list (Variant B) and §7 step 11 case (a) "both signatures present and both valid." Tests:

- The verifier reconstructs each algorithm's `sign_payload` independently (Variant B per spec §4.2: each entry's signature covers its own algorithm-bound `sign_payload`).
- The verifier validates each signature against the algorithm-specific public key.
- Both algorithms validate; the day reports PASS (co-signed).

## Tampering recipe

Start from `../010-tenant-ikm-rotation-mid-day/` chain (the rotation case provides multi-version chain entries). Construct a seal record:

```
tenant_id: tenant-ffiec-test-1
seal_date: 2026-05-06
spec_version: v1.0
format_version: v1
merkle_root: (computed from rotation chain payload_hashes)
algorithm: ed25519
key_versions: [1, 2]
hkdf_inputs_digest: 6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65
signature: <Ed25519 signature over Ed25519 sign_payload>
signatures:
  - algorithm: ed25519
    signature: <Ed25519 signature over Ed25519-bound sign_payload>
  - algorithm: dilithium3
    signature: <Dilithium3 signature over Dilithium3-bound sign_payload>
```

Each algorithm's `sign_payload` per spec §4.3 (consolidated v1.0a form — 2026-05-07 bound `cadence`, `dev_mode`, and `sign_payload_version` under the HSM signature):

```
ed25519 sign_payload:
  ffiec.chain-of-custody.v1\n
  v1.0a\n
  ed25519\n
  v1\n
  tenant-ffiec-test-1\n
  2026-05-06\n
  hex(merkle_root)\n
  hex(hkdf_inputs_digest)\n
  daily\n
  0

dilithium3 sign_payload:
  ffiec.chain-of-custody.v1\n
  v1.0a\n
  dilithium3\n
  v1\n
  tenant-ffiec-test-1\n
  2026-05-06\n
  hex(merkle_root)\n
  hex(hkdf_inputs_digest)\n
  daily\n
  0
```

Defaults for this corpus seal: `cadence = "daily"`, `dev_mode = false` (serialized as the single ASCII byte `"0"` per spec §4.3), `sign_payload_version = "v1.0a"`. No trailing newline after the `dev_mode` byte — nine `0x0A` separators total (1 magic-line terminator plus 8 inter-field terminators between the 9 fields that follow).

Note the third line differs by algorithm — Variant B is the load-bearing property closing algorithm-confusion attacks. Each algorithm's `sign_payload` independently carries the `sign_payload_version`, `cadence`, and `dev_mode` lines.

## Expected verifier outcome

```
Status: PASS (co-signed)
Spec §7 step 11 case: (a) both signatures present and both valid
per_algorithm_results:
  - ed25519: PASS
  - dilithium3: PASS
```

## v1.x population guidance

The v1.x implementer:
1. Generates Ed25519 + Dilithium3 (or SLH-DSA) test keypairs (clearly marked TEST USE ONLY).
2. Computes both algorithm-bound sign_payloads.
3. Signs each with the corresponding algorithm.
4. Publishes `inputs/`, `expected/`, including the per-algorithm public keys, sign_payloads (text and hex), and signatures.
5. Updates `chain_vectors.json` with `dual_algorithm_chain` section mirroring `single_chain` and `rotation_chain`.

The conformance contract: a v1.x verifier reproduces both signatures byte-for-byte from the published inputs.
