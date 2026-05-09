# Case 024 — Per-device HKDF derivation (§10.32)

## What this case verifies

Spec §10.32 normates the per-device session-key derivation extension to §4.1's tenant-bound HKDF. The extended derivation adds a third name segment binding a device identity:

```
info = HKDF_INFO_BASE || "|" || utf8(tenant_id) || "|" || utf8(device_id)
session_key = HKDF-SHA256(ikm, salt = HKDF_SALT, info = info, length = 32)
```

The `|` separator (byte `0x7c`) is the same fixed separator §4.1 uses; tenant_id or device_id values containing `|` are encoded as-is — both producer and verifier produce byte-identical info bytes from the same identifiers.

This case pins the exact info bytes, derived session keys, and key fingerprints for two devices under the same tenant + IKM, plus the §4.1-only baseline derivation. The byte-distinctness assertion is the central conformance bar:

- A `(ikm, tenant_id, device_id_A)` derivation produces a different 32-byte session key than `(ikm, tenant_id)` alone — a §4.1-only verifier walking a §10.32 chain fails its MAC checks.
- A `(ikm, tenant_id, device_id_A)` and `(ikm, tenant_id, device_id_B)` derivation produce different session keys — chain entries from device A and device B do not collide on session-key bytes.

## Inputs

| Field | Value |
|---|---|
| `tenant_id` | `tenant-ffiec-test-1` (matches central corpus) |
| `ikm_label` | `ffiec-cross-language-test-vector-ikm-v1-32bytes!` (UTF-8 → 48 bytes; matches central corpus IKM_V1) |
| `device_a` | `tablet-A1B2-C3D4` (field-tablet inventory tag) |
| `device_b` | `tpm-7eb3-c000-0001-aabb` (TPM 2.0 endorsement-key fingerprint) |

The two `device_id` values exercise distinct identity-authority shapes — a non-cryptographic inventory tag and a TPM-rooted hardware identifier. Both are valid §10.32 device_id values per the spec's "non-empty UTF-8 strings ≤ 256 chars" constraint; the institution's CC8.1 names the device-identity scheme it uses.

`tenant_id` and IKM match the central corpus values so a verifier reading both this case and `chain_vectors.json` can confirm the §4.1 baseline session key reproduces under both fixtures.

## Computed info bytes

| Recipe | Byte length | UTF-8 form |
|---|---|---|
| §4.1 base | 50 | `ffiec.chain-of-custody.v1.info\|tenant-ffiec-test-1` |
| §10.32 device A | 67 | `ffiec.chain-of-custody.v1.info\|tenant-ffiec-test-1\|tablet-A1B2-C3D4` |
| §10.32 device B | 74 | `ffiec.chain-of-custody.v1.info\|tenant-ffiec-test-1\|tpm-7eb3-c000-0001-aabb` |

Hex forms are in `expected.json`.

## Expected outputs

| Field | Value |
|---|---|
| §4.1 baseline session key | `f061210167d307cffe4b91c2eaad9d386a8831c04957b458088e5bafd1f63ee9` |
| §10.32 device A session key | `ae4bdffcf7917ea216370f82ae7d41fb16eccb622d3ee571467ed242e65af5fd` |
| §10.32 device B session key | `2b7f7cc14cd3f02300647fcaaceb465915f99f0a4948a38b95ffe4746f0ca869` |
| `key_fingerprint` (all three derivations) | `54c245a894b81d788d372badc563531f` |

The §4.1 baseline session key matches `session_key_v1_hex` published in the central `spec/test-vectors/README.md` — same `(ikm, tenant_id)` pair, same HKDF inputs. The case is anchored to the existing corpus.

The `key_fingerprint` is identical across all three derivations because the fingerprint binds tenant + IKM only; `device_id` is NOT part of the fingerprint per doc §3.5. This is intentional — the fingerprint exists to detect IKM-swap misconfigurations at lookup time, not device-binding drift. A device-binding misconfiguration (using device A's session key on device B's chain entries) is detected by MAC failure, not by fingerprint mismatch.

## Conformance behavior

A conforming implementation supporting §10.32:

1. Builds the extended HKDF `info` as `HKDF_INFO_BASE || 0x7c || utf8(tenant_id) || 0x7c || utf8(device_id)` and reproduces the `info_hex` byte-for-byte.
2. Derives the session key via HKDF-SHA256 with `salt=HKDF_SALT`, `info=<extended>`, `length=32` and reproduces the corresponding `section_10_32_device_*_session_key_hex` byte-for-byte.
3. Computes `key_fingerprint = SHA-256(utf8(tenant_id) || ikm)[:16]` per doc §3.5 — UNCHANGED from §4.1; device_id is NOT bound.
4. Refuses an empty `tenant_id` or empty `device_id` per §10.32 normative ("both MUST be non-empty UTF-8 strings").
5. Refuses a `device_id` exceeding 256 chars per §10.32 normative.

## Negative cases this fixture supports

- An implementation that uses a different separator byte (e.g., `:` instead of `|`) produces different info bytes and a different session key — the `info_hex` byte pin fails loudly.
- An implementation that NFC-normalises `device_id` UTF-8 (or applies any other Unicode normalisation) produces different info bytes for non-ASCII identifiers — §10.32 normates "encoded as-is" so normalisation is non-conformant.
- An implementation that includes `device_id` in the key fingerprint (e.g., `SHA-256(utf8(tenant_id) || utf8(device_id) || ikm)[:16]`) produces different fingerprint values per device — case 024 pins the device-independent fingerprint, surfacing this divergence.

## Cross-references

- Spec §10.32 per-device session-key derivation (the section this case pins)
- Spec §4.1 per-tenant HKDF binding (the base derivation §10.32 extends — the case's `section_4_1_baseline` outputs come from this)
- Spec §10.6 IKM minimum length (the IKM constraint §10.32 inherits unchanged)
- Doc §3.5 key_fingerprint recipe (UNCHANGED under §10.32 — no device_id binding)
- Central corpus `chain_vectors.json` (the §4.1 baseline session_key_v1 used here matches that file's `session_key_v1_hex`)
- Case 025 (`025-attestation-android-keystore/`) — §10.35 hardware-rooted attestation, the natural compositional pair when the institution operates both §10.32 + §10.35

## Reproduction

The `_compute.py` script in this directory regenerates `expected.json` from the pinned inputs above. It uses only the Python standard library; no external HKDF library is required (the script implements RFC 5869 directly to avoid library-version drift).

```
python _compute.py
```

Cross-validation: the §10.32 device-A session key reproduced by Herald.Py's `herald._crypto.keys.derive_session_key` from the same info bytes is byte-identical to `ae4bdf…f5fd`.
