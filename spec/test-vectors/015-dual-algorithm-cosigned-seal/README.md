# Case 015 — Dual-algorithm co-signed seal (multi-algorithm placeholder)

> **Status — documented placeholder.** This case is a structural placeholder for the v1.x byte-level fixture covering FIPS 204 ML-DSA + Ed25519 dual-algorithm co-signed seals and FIPS 205 SLH-DSA + Ed25519 dual-algorithm co-signed seals. The byte-level keypair and signature fixture is **implementor-supplied, pending NIST FIPS 204 final test-vectors stabilization and FIPS-validated HSM availability for ML-DSA and SLH-DSA across two or more major vendors**. The structural recipe and expected verifier behavior are documented here so the v1.x implementer can populate the byte-level fixture from the same recipe.
>
> See `description.md` in this directory for the original recipe (which targeted Dilithium3 — superseded by the FIPS 204 ML-DSA-65 nomenclature). See `docs/cryptographic-agility-roadmap.md` §5 for the cryptographic-agility roadmap that frames this case's role in the v1.0a forward posture.

## What this case verifies

The dual-algorithm co-signed seal per spec §4.2 `signatures` list (Variant B) and §7 step 11 case (a) "both signatures present and both valid." The case exercises:

- The verifier reconstructs each algorithm's `sign_payload` independently (Variant B per spec §4.3 — each algorithm's signature covers its own algorithm-bound `sign_payload` with the third line differing by algorithm)
- The verifier validates each signature against the algorithm-specific public key resolved from the institution's tenant key registry
- Both algorithms validate; the day reports PASS (co-signed)
- The AND-security posture under §7 step 11 — a verifier presented with one valid + one invalid signature returns FAIL with the appropriate failure-reason string

## Two pairings the v1.x implementer ships

The v1.x amendment ships byte-level fixtures for two concrete pairings. Each pairing is independently testable; institutions may adopt either pairing based on their cryptographic-foundations posture.

### Pairing 1 — Ed25519 + ML-DSA-65 (FIPS 204)

ML-DSA-65 is the NIST-published name for Dilithium3 at security level 3 (192-bit classical / 192-bit quantum security under Module-LWE assumptions). The fixture under this pairing is `015a-ed25519-ml-dsa-65/` (subdirectory to be created by v1.x implementer).

Expected fixture-byte structure:

```
inputs/
  ed25519_keypair.txt        # TEST-USE-ONLY Ed25519 keypair (32-byte private, 32-byte public)
  ml_dsa_65_keypair.txt      # TEST-USE-ONLY ML-DSA-65 keypair (~4032-byte private, ~1952-byte public)
  ed25519_sign_payload.txt   # The 10-line Ed25519-bound sign_payload (canonical UTF-8)
  ml_dsa_65_sign_payload.txt # The 10-line ML-DSA-65-bound sign_payload (canonical UTF-8)
  chain_entries.json         # Source chain entries (typically reused from 010-tenant-ikm-rotation-mid-day/)
  seal_record.json           # The dual-algorithm seal record under construction

expected/
  ed25519_signature.hex      # 64-byte Ed25519 signature over Ed25519 sign_payload (128 hex chars)
  ml_dsa_65_signature.hex    # ~3293-byte ML-DSA-65 signature over ML-DSA-65 sign_payload
  merkle_root.hex            # 64 hex chars — the Merkle root both signatures cover
  verifier_output.txt        # PASS (co-signed); per-algorithm pass; spec §7 step 11 case (a)
  verifier_output.json       # Structured verifier output for programmatic conformance check
```

Approximate byte counts for the v1.x fixture:

| Artifact | Size |
|---|---|
| Ed25519 public key | 32 bytes |
| Ed25519 private key | 32 bytes |
| Ed25519 signature | 64 bytes |
| ML-DSA-65 public key | ~1952 bytes |
| ML-DSA-65 private key | ~4032 bytes |
| ML-DSA-65 signature | ~3293 bytes |
| Total per-seal signature payload | ~3357 bytes |
| Annual cost per tenant (daily cadence) | ~1.2 MB |

### Pairing 2 — Ed25519 + SLH-DSA-SHA2-192f (FIPS 205)

SLH-DSA-SHA2-192f is the NIST-published name for SPHINCS+ at security level 3 with SHA-2 underlying hash and 192-bit fast variant. The fixture under this pairing is `015b-ed25519-slh-dsa-sha2-192f/` (subdirectory to be created by v1.x implementer).

Expected fixture-byte structure (parallel to Pairing 1):

```
inputs/
  ed25519_keypair.txt
  slh_dsa_sha2_192f_keypair.txt  # TEST-USE-ONLY SLH-DSA keypair (~96-byte private, 48-byte public)
  ed25519_sign_payload.txt
  slh_dsa_sha2_192f_sign_payload.txt
  chain_entries.json
  seal_record.json

expected/
  ed25519_signature.hex
  slh_dsa_sha2_192f_signature.hex  # ~35664-byte SLH-DSA signature
  merkle_root.hex
  verifier_output.txt
  verifier_output.json
```

Approximate byte counts:

| Artifact | Size |
|---|---|
| Ed25519 public key | 32 bytes |
| Ed25519 private key | 32 bytes |
| Ed25519 signature | 64 bytes |
| SLH-DSA-SHA2-192f public key | 48 bytes |
| SLH-DSA-SHA2-192f private key | 96 bytes |
| SLH-DSA-SHA2-192f signature | ~35664 bytes |
| Total per-seal signature payload | ~35728 bytes |
| Annual cost per tenant (daily cadence) | ~13 MB |

The SLH-DSA fixture is significantly larger than the ML-DSA fixture but inherits SHA-256's security margin — under Grover, SLH-DSA degrades to roughly 2^96 against preimage attacks, comfortable for a 50-year horizon. SLH-DSA has no structured-lattice concerns (unlike ML-DSA, whose Module-LWE hardness is still under active cryptanalysis). Institutions choosing the SLH-DSA pairing accept the larger signatures in exchange for the more conservative cryptographic foundations.

## sign_payload byte structure for each algorithm

Each algorithm's `sign_payload` follows the v1.0a 10-line form. Concrete examples for tenant `tenant-ffiec-test-1`, seal-date `2026-05-06`, daily cadence, dev_mode=false:

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

ml-dsa-65 sign_payload:
  ffiec.chain-of-custody.v1\n
  v1.0a\n
  ml-dsa-65\n
  v1\n
  tenant-ffiec-test-1\n
  2026-05-06\n
  hex(merkle_root)\n
  hex(hkdf_inputs_digest)\n
  daily\n
  0

slh-dsa-sha2-192f sign_payload:
  ffiec.chain-of-custody.v1\n
  v1.0a\n
  slh-dsa-sha2-192f\n
  v1\n
  tenant-ffiec-test-1\n
  2026-05-06\n
  hex(merkle_root)\n
  hex(hkdf_inputs_digest)\n
  daily\n
  0
```

Note the third line differs by algorithm — Variant B's load-bearing property closing algorithm-confusion attacks. The other lines are identical because they describe the same seal-day's evidence.

Defaults for this corpus seal: `cadence = "daily"`, `dev_mode = false` (serialized as the single ASCII byte `"0"` per spec §4.3), `sign_payload_version = "v1.0a"`. No trailing newline after the `dev_mode` byte — nine `0x0A` separators total (1 magic-line terminator plus 8 inter-field terminators between the 9 fields that follow).

## Expected verifier outcome

```
Status: PASS (co-signed)
Spec §7 step 11 case: (a) both signatures present and both valid
per_algorithm_results:
  - ed25519: PASS
  - ml-dsa-65: PASS    (or slh-dsa-sha2-192f: PASS depending on pairing)
```

A v1.x verifier shipping FIPS 204 (or FIPS 205) support reproduces both signatures byte-for-byte from the published inputs.

## Negative-case fixtures (v1.x extension)

The v1.x implementer should also ship negative-case fixtures exercising spec §7 step 11 cases (b) through (e):

| Case | Scenario | Expected verifier outcome |
|---|---|---|
| (a) | Both signatures present and both valid | PASS (co-signed) |
| (b) | Single algorithm in seal record (Ed25519 only); single-algorithm verifier path | PASS (single-algorithm) |
| (c) | Both signatures present; first valid, second invalid | FAIL with reason `algorithm 2 signature verification failed` |
| (d) | Both signatures present; first invalid, second valid | FAIL with reason `algorithm 1 signature verification failed` |
| (e) | Two signatures named in `signatures` list but `algorithm` field on seal record names a different one | FAIL with reason `algorithm/key-type mismatch at signature verification` |

Each negative case is a separate `015c-*` through `015f-*` subdirectory under the implementer's v1.x work.

## Why the byte-level fixture is deferred

Three independent factors drive the byte-level fixture deferral:

1. **NIST FIPS 204 / FIPS 205 final-test-vectors stabilization.** The final NIST test vectors for ML-DSA and SLH-DSA were published in 2024 but interoperability evidence across implementations is still maturing in 2026. A test-vector fixture shipped before cross-implementation interop converges risks showing one implementation's quirks rather than the algorithm's actual conformance.
2. **HSM-vendor PQ-readiness availability.** ML-DSA and SLH-DSA HSM support is shipping in 2026-2028 across the major vendors (AWS CloudHSM, Google Cloud KMS, Azure Key Vault HSM, Thales Luna, Entrust nShield). The v1.x amendment that lands the fixture is gated on FIPS-validated HSM support being available from at least two major vendors — the chain's seal-signing key MUST be HSM-bound per spec §10.5, so a fixture without HSM-validatable byte structure is operationally premature.
3. **Algorithm-dispatch plumbing in v1.0a is sufficient for the v1.x amendment.** The structural plumbing (per-algorithm `sign_payload` Variant B, the `signatures` list, the verifier's §7 step 11 dispatch) is in place in v1.0a. The v1.x amendment adds the byte-level fixture; the v1.0a verifier already dispatches correctly when given a dual-algorithm seal — it just doesn't have a fixture to validate against.

## What the structural placeholder accomplishes for v1.0a

The structural placeholder is the v1.0a close-out for Vasiliev P-1 (hybrid signatures structurally admitted but no shipped fixture). The placeholder:

- Documents the byte structure expected of the v1.x fixture
- Names the two recommended pairings (ML-DSA-65 and SLH-DSA-SHA2-192f) with concrete byte counts
- Specifies the sign_payload form for each algorithm
- Specifies the expected verifier outcome
- Names the negative-case fixtures the v1.x implementer should also ship
- Documents why the byte-level fixture is deferred and what the deferral is gated on

The structural readiness is the v1.0a accomplishment. The byte-level fixture is the v1.x amendment work. A v1.0a verifier built today against the spec text plus this structural placeholder is ready to accept dual-algorithm seals when the v1.x fixture lands.

## v1.x implementer checklist

For the v1.x implementer producing the byte-level fixture:

- [ ] Generate TEST-USE-ONLY ML-DSA-65 keypair under FIPS 204
- [ ] Generate TEST-USE-ONLY SLH-DSA-SHA2-192f keypair under FIPS 205
- [ ] Generate TEST-USE-ONLY Ed25519 keypair (or reuse from existing fixture)
- [ ] Compute the three algorithm-bound sign_payloads for the chosen seal-day (typically reused from `010-tenant-ikm-rotation-mid-day/` chain)
- [ ] Sign each sign_payload with the corresponding algorithm
- [ ] Publish `015a-ed25519-ml-dsa-65/inputs/` and `015a-ed25519-ml-dsa-65/expected/` directories
- [ ] Publish `015b-ed25519-slh-dsa-sha2-192f/inputs/` and `015b-ed25519-slh-dsa-sha2-192f/expected/` directories
- [ ] Update `chain_vectors.json` with `dual_algorithm_chain` section mirroring `single_chain` and `rotation_chain`
- [ ] Ship negative-case fixtures `015c-*` through `015f-*` per the table above
- [ ] Document the v1.x amendment that lifts case 015 from structural placeholder to normative
- [ ] Validate the fixtures against an independently-implemented verifier (cross-implementation conformance)

The conformance contract: a v1.x verifier reproduces both signatures byte-for-byte from the published inputs, and the negative-case verifier outcomes match the expected failure-reason strings.

## Cross-references

- `docs/cryptographic-agility-roadmap.md` §5 — the full hybrid signature variant B specification with concrete primitive selections
- `docs/cryptographic-agility-roadmap.md` §10 — the HNDL response and the dated dual-algorithm seal mandate (2030-01-01)
- `spec/chain-of-custody-v1.md` §4.2 — the `signatures` list specification
- `spec/chain-of-custody-v1.md` §4.3 — the per-algorithm `sign_payload` Variant B specification
- `spec/chain-of-custody-v1.md` §7 step 11 — the dispatch table covering cases (a) through (e)
- `description.md` (in this directory) — the original recipe targeting Dilithium3 (superseded by ML-DSA-65 under FIPS 204 nomenclature)
