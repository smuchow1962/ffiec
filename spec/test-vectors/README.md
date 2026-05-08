# FFIEC Chain-of-Custody — Conformance Test Vectors (v1)

> **What this directory is.** The byte-level conformance contract for FFIEC chain-of-custody v1. Every conforming SDK, ledger server, and verifier MUST reproduce the expected outputs in this corpus byte-for-byte from the same inputs. The corpus exists so that two independent implementations cannot disagree on what v1 means.
>
> **Primary fixture.** [`chain_vectors.json`](chain_vectors.json) carries the full byte-level fixture: HKDF inputs, expected session keys, expected `key_fingerprint` values, an `hkdf_inputs_digest`, two complete chains (a single-IKM 5-event chain and a rotation-mid-chain 5-event chain), and the resulting Merkle roots and signed payloads. This is the load-bearing artifact.

## Constants (FFIEC-conformance posture, locked at v1)

The constants below are the FFIEC-conformance values per spec §4.1. Implementations supporting vendor-flag mode (spec §4.1.2) substitute their vendor-namespaced byte values for `HKDF_SALT` and `HKDF_INFO_BASE` at SDK-construct time; chains produced under non-FFIEC constants are non-FFIEC chains and would fail this conformance corpus by design. The `hkdf_inputs_digest` published below is the FFIEC-posture digest for the test tenant; a vendor-mode digest computed from the vendor's own constants would differ.

```
HKDF_SALT          = b"ffiec.chain-of-custody.v1.salt"     // 30 bytes UTF-8 (FFIEC-conformance)
HKDF_INFO_BASE     = b"ffiec.chain-of-custody.v1.info"     // 30 bytes UTF-8 (FFIEC-conformance)
length             = 32                                    // HKDF output length
length_LE32        = (32).to_bytes(4, "little") = 0x20 0x00 0x00 0x00 (4 bytes)
genesis_hash       = 32 zero bytes                         // for seq=1
empty_day_merkle_root = SHA-256("") =
  e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

**Algorithm string table (locked at v1.0).** The signature-algorithm identifier strings used in `seal.algorithm`, in `signatures[].algorithm` (under dual-algorithm posture), and on the `algorithm` line of `sign_payload`:

```
ed25519           = "ed25519"      (lowercase ASCII; v1.0 single-algorithm)
dilithium3        = "dilithium3"   (lowercase ASCII; v1.x post-quantum candidate)
slh_dsa_shake_128s = "slh_dsa_shake_128s"  (lowercase ASCII; v1.x post-quantum candidate)
hmac_sha_256      = "HMAC-SHA-256" (uppercase + dashes; for ffiec.chain.algorithm per-entry forensic-only)
```

The chain-stamp algorithm name (`ffiec.chain.algorithm`, default `"HMAC-SHA-256"`) uses uppercase + dashes for legacy compatibility; the seal-record signature algorithm names use lowercase + underscores. The byte-form is locked because `algorithm` is bound into `sign_payload` and a different byte-form on different implementations would break cross-implementation seal verification.

`length_LE32` clarifies that the 32 in `hkdf_inputs_digest = SHA-256(salt || info_for_tenant || length_LE32)` is the 4-byte little-endian encoding of the integer 32, NOT 32 bytes of literal value.

These are normative per spec §4.1 + design 03 §3.1 (empty-day convention per RFC 6962 §2.1). The HKDF constants are the FFIEC-conformance values per spec §4.1.2; a vendor SDK supporting vendor-flag mode passes its own `HKDF_SALT` / `HKDF_INFO_BASE` byte values at construct time and produces non-FFIEC chains whose `hkdf_inputs_digest` differs from the value below. Changing the FFIEC constants themselves is a `format_version` bump (and would invalidate every FFIEC-posture chain ever produced under the old constants, which is the point of binding the constants into `hkdf_inputs_digest`). The `empty_day_merkle_root` is pinned here as a one-line conformance check: a v1 verifier reading a tenant-day with zero events recomputes the Merkle root as `SHA-256("")` and compares against this value byte-for-byte.

## Test-vector inputs

```
tenant_id   = "tenant-ffiec-test-1"
run_id      = "RUN-FFIEC-VECTOR"
ikm_v1_hex  = 6666...7321   // "ffiec-cross-language-test-vector-ikm-v1-32bytes!" UTF-8 (48 bytes)
ikm_v2_hex  = 6666...7321   // same shape with v2; for the rotation case
seal_date   = "2026-05-06"
```

These inputs are deliberately stable and human-readable so a conforming implementation can copy them verbatim. The full hex bytes are in `chain_vectors.json` under `inputs`.

## Expected outputs

```
session_key_v1     = HKDF(ikm_v1, salt=HKDF_SALT, info=info_for("tenant-ffiec-test-1"), length=32)
session_key_v1_hex = f061210167d307cffe4b91c2eaad9d386a8831c04957b458088e5bafd1f63ee9

session_key_v2     = HKDF(ikm_v2, salt=HKDF_SALT, info=info_for("tenant-ffiec-test-1"), length=32)
session_key_v2_hex = bbb3a221fcf8ee4320817a04bf5a6801ac61fef3d09ab4d4db851faa3403477e

key_fingerprint_v1 = SHA-256(utf8("tenant-ffiec-test-1") || ikm_v1)[:16]
key_fingerprint_v1_hex = 54c245a894b81d788d372badc563531f

key_fingerprint_v2 = SHA-256(utf8("tenant-ffiec-test-1") || ikm_v2)[:16]
key_fingerprint_v2_hex = 1e05971d136fee9e6e9c1a96d2fca05d

hkdf_inputs_digest = SHA-256(HKDF_SALT || info_for("tenant-ffiec-test-1") || (32).to_bytes(4, "little"))
hkdf_inputs_digest_hex = 6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65
```

Per-event `payload_hash` values are in `chain_vectors.json` under `single_chain` and `rotation_chain`.

## Conformance test

A conforming implementation MUST:

1. Read `chain_vectors.json`.
2. Derive `session_key_v1`, `session_key_v2`, `key_fingerprint_v1`, `key_fingerprint_v2`, and `hkdf_inputs_digest`. Each MUST match the published hex byte-for-byte.
3. For each event in `single_chain` and `rotation_chain`: recompute `payload_hash = HMAC-SHA-256(session_key_for_kv, prev_hash || event_canonical)`. Each MUST match the published hex byte-for-byte.
4. Compute `merkle_root` over each chain's ordered `payload_hash` values per RFC 6962. Each MUST match `merkle_root_single_hex` / `merkle_root_rotation_hex`.
5. Construct `sign_payload` per spec §4.3 (amendment form, v1.0a) from `sign_payload_version="v1.0a"`, `algorithm="ed25519"`, `format_version="v1"`, `tenant_id`, `seal_date`, `merkle_root`, `hkdf_inputs_digest`, `cadence="daily"`, and `dev_mode=false` (the corpus seal defaults; published as `sign_payload_seal_defaults` in `chain_vectors.json`). The byte sequence is the magic line followed by nine fields, each separated by `0x0A`, with no trailing newline:

   ```
   ffiec.chain-of-custody.v1\n   ← magic line
   v1.0a\n                       ← sign_payload_version (line 2)
   ed25519\n                     ← algorithm
   v1\n                          ← format_version
   tenant-ffiec-test-1\n         ← tenant_id
   2026-05-06\n                  ← iso8601_date(tenant_day)
   hex(merkle_root)\n            ← 64 lowercase hex chars
   hex(hkdf_inputs_digest)\n     ← 64 lowercase hex chars
   daily\n                       ← cadence
   0                             ← dev_mode (single ASCII byte, no trailing \n)
   ```

   Ten fields total. Nine `0x0A` separators (one magic-line terminator plus eight inter-field terminators between the nine fields that follow it). The byte sequence MUST match `sign_payload_single_hex` / `sign_payload_rotation_hex`. The text form is also published as `sign_payload_single_text` / `sign_payload_rotation_text` for human inspection.

If any byte disagrees, the implementation is non-conforming. Investigate (a) JCS canonical-form mismatch, (b) HKDF input byte-encoding differences, (c) UTF-8 normalization drift, (d) different hash-output endianness or truncation.

**JCS edge-case conformance (mandatory).** The basic fixtures above exercise integers, ASCII strings, and base64-encoded byte strings. They do NOT exercise the RFC 8785 edge cases that historically cause silent disagreement between implementations: floats, non-ASCII Unicode, NaN/Infinity rejection, surrogate pairs, control characters, deeply nested objects, very long strings, and non-trivial object-key ordering. Case [`008-jcs-edge-cases/`](008-jcs-edge-cases/description.md) closes that gap. Implementations that pass the basic fixtures (`001-` through `015-`) but fail any of the `008-jcs-edge-cases/` fixtures are NOT v1.0-conformant. The JCS edge cases are part of the conformance bar.

## Per-case directories

Each numbered directory carries a `description.md` documenting the case's purpose, the inputs (which subset of `chain_vectors.json` to consume), and the expected outputs (the subset of `chain_vectors.json` to compare against). The reference Go runner (`runner.go`, ships with a future implementation) walks each directory and produces pass/fail per case.

Cases:

| Case | What it exercises |
|---|---|
| `001-single-event-empty-prev/` | The smallest case: one event with `seq=1`, `prev_hash` is 32 zero bytes |
| `002-multi-event-same-run/` | Five events in one run; chain links across multiple events |
| `003-multi-run-same-day/` | Two runs in the same tenant-day. Run 2's first event resets `prev_hash` to 32 zero bytes (cross-run isolation). Daily Merkle aggregates all events in `(run_id, seq)` order. |
| `008-jcs-edge-cases/` | RFC 8785 JCS conformance: floats, non-ASCII Unicode, surrogate pairs, NaN/Infinity rejection, very long strings, deeply nested objects, control characters, object-key ordering, numeric edge cases. Cross-validated across two independent implementations. |
| `010-tenant-ikm-rotation-mid-day/` | Rotation: events 1-3 under `key_version=1` with `ikm_v1`, events 4-5 under `key_version=2` with `ikm_v2`. Same tenant. Verifier walks both halves. |
| `015-dual-algorithm-cosigned-seal/` | Dual-algorithm cosigned seal record. |
| `016-non-power-of-2-merkle/` | RFC 6962 §2.1 odd-leaf right-promote balancing. Three sub-fixtures: 3-leaf, 5-leaf, 7-leaf trees with pinned roots. |
| `negative/` | Tampering cases the verifier MUST report as failure with a specific named reason. See `negative/README.md`. |

Future cases (per `docs/design/08-test-vectors.md` §5):

- `004-empty-day/` — empty-tree convention
- `005-out-of-order-input/` — events received in reverse seq; ordering applied at Merkle time
- `006-late-binding/` — late-arriving event included in the next day's seal
- `007-ten-thousand-events-perf/` — bulk Merkle streaming
- `009-session-key-rotation-mid-process/` — same run, different processes, same IKM
- `011-cross-spec-version/` — v1 verifier reading a hypothetical v2 file (refuses)
- `012-cross-day-run/` — run that spans UTC midnight
- `013-dag-multi-process-flow/` — DAG-shaped run with `dag_parents`
- `014-multi-master-version-seal/` — same as `010` but exercises the seal record's `key_versions` list

## Versioning

The corpus is append-only. Existing test cases are never modified after a spec version locks. A bug in a test case is fixed by adding a new case that supersedes the buggy one; the old case is marked deprecated but not deleted.

**`sign_payload` v1.0a consolidation (2026-05-07).** Spec §4.3 was extended on 2026-05-07 in two consecutive same-day cycles, consolidated into the single canonical v1.0a form. The morning cycle bound two seal-record fields under the HSM signature: `cadence` (`"hourly"` | `"daily"` | `"weekly"`) and `dev_mode` (single ASCII byte `"1"` if true, `"0"` if false or absent). The afternoon cycle added `sign_payload_version` as the new second line of the structure (immediately after the magic line, before `algorithm`); binding the version identifier itself into the second line means a tampered `sign_payload_version` value on a seal record is detected at signature verification because the verifier reconstructs `sign_payload` using the field as written, and a mismatch between the field value and the byte form actually used by the signer produces a signature failure.

The progression of the wire form during 2026-05-07:

| Form | Fields | Separators | Status |
|---|---|---|---|
| Original v1.0 (pre-amendment) | 7 (magic + algorithm + format_version + tenant_id + date + merkle_root + hkdf_inputs_digest) | 6 `0x0A` | Frozen for pre-amendment chains; verified under the pre-amendment 6-line dispatch when `sign_payload_version` is absent |
| Intermediate (morning cycle, 2026-05-07) | 9 (above + cadence + dev_mode) | 8 `0x0A` | Transitional; never shipped as conformant on its own |
| **Canonical v1.0a (afternoon cycle, 2026-05-07)** | **10 (above + sign_payload_version on line 2)** | **9 `0x0A`** | **Current conformant form** |

Older fixtures that pin either the original 6-line form OR the intermediate 9-line form are NOT v1.0a-conformant. The `chain_vectors.json` `sign_payload_*_hex` and `sign_payload_*_text` values, case 003's `sign_payload_*` values, and case 015's per-algorithm structures were regenerated under the consolidated v1.0a form. Per-event `payload_hash`, Merkle roots, `hkdf_inputs_digest`, session keys, and `key_fingerprint` values are UNCHANGED — the wire bumps are at the seal-record signature layer only; per-event MACs and Merkle aggregation are unaffected. A pre-amendment chain (no `sign_payload_version` field on its seal record) continues to verify cleanly under amendment-aware verifiers per spec §4.3's verifier dispatch — pre-amendment chains use the 6-line form; v1.0a chains use the 10-line form.

## Where the byte values come from

`chain_vectors.json` was computed by a script that exercises the spec §4.1 construction directly (HKDF-SHA-256, HMAC-SHA-256, RFC 6962 Merkle, and the §4.3 sign_payload). The script is in `tools/compute-test-vectors.py` (ships with future implementations); regenerating is deterministic and reproducible.

## Note on signature

`chain_vectors.json` does NOT carry an Ed25519 signature today. Generating Ed25519 signatures is non-deterministic in some implementations (the Go stdlib uses deterministic Ed25519 per RFC 8032; some libraries use randomized variants). To keep the corpus implementation-portable, the conformance test stops at "the implementation builds the byte-identical `sign_payload`" and lets each implementation produce its own signature against a test private key. A future iteration may add a deterministic-Ed25519 expected signature.
