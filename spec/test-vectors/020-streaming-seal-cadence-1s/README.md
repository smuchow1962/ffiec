# Case 020 — sign_payload v1.0b 12-line form (cadence = per_second)

## What this case verifies

Spec §10.27 normates `cadence = "per_second"` as a §4.2.1 enumerated value for streaming-mode 1-second cadence. The cadence value is bound under the §4.3 sign_payload form on its dedicated line; a tampered cadence value in the seal record is detected at signature verification because the verifier reconstructs sign_payload using the field as written.

This case pins the v1.0b 12-line sign_payload byte form for a streaming-mode 1-second-cadence seal so a clean-room implementation proves it emits the §10.27 streaming cadence value byte-identical to other implementations. The byte-distinct rule from §10.27 (`per_hour` and `hourly` produce the same one-hour interval but are NOT interchangeable in any chain) generalises here: the literal byte sequence `per_second` is what binds, not its "1-second" semantic meaning.

## Inputs

A streaming-mode seal-record for one 1-second interval, single key version, single KMS handle. The inputs are deliberately a degenerate-simple chain (one key version, one URI) because case 020 pins the **cadence** byte form, not the v1.0b key-versions / kms-handle-uris machinery (case 018 pins the mixed-version mixed-URI day; case 019 pins the empty-day variant).

| Field | Value |
|---|---|
| `sign_payload_version` | `v1.0b` |
| `algorithm` | `ed25519` |
| `format_version` | `v1` |
| `tenant_id` | `demo-bank-streaming-001` |
| `seal_date` | `2026-05-07` |
| `merkle_root` | `0001020304050607080910111213141516171819202122232425262728293031` (placeholder; the v1.0b form is independent of the merkle_root value, so any 64-char lowercase hex works for the canonical-bytes test) |
| `hkdf_inputs_digest` | `ffeeddccbbaa99887766554433221100ffeeddccbbaa99887766554433221100` (placeholder; same rationale) |
| `cadence` | `per_second` |
| `dev_mode` | `false` (serialized as the single ASCII byte `0`) |
| `key_versions` (distinct) | `[1]` |
| `kms_handle_uris` (distinct) | `aws-kms:arn:aws:kms:us-east-1:123456789012:key/streaming-mock-1` |

The `seal_period_start_utc` field (`2026-05-07T12:34:56.000000Z` per `input.json`) is institution-trusted ledger-side metadata per §4.2 — it is **NOT** part of the sign_payload byte form. Verifiers walking the chain check `seal_period_start_utc` continuity against the cadence-interval boundary discipline (§10.28) but do not bind it under the signature.

## Computed canonical fields

`key_versions_canon` is computed by sorting the distinct `key_version` integer values into ascending order and joining with single ASCII commas:

```
key_versions_canon = "1"
```

`kms_handle_uris_digest` is computed by sorting the distinct URI values into code-point ascending order, joining with single `\n` (0x0A) bytes (no trailing newline), and applying `SHA-256`:

```
canonical_uri_form = "aws-kms:arn:aws:kms:us-east-1:123456789012:key/streaming-mock-1"
canonical_uri_form_bytes_len = 63
SHA-256(canonical_uri_form) = ccbcf494d1dba1ced4d62e081f4b7b7aab88a21e15cd67a7be0cb5b6a78687f8
```

## Expected sign_payload (12 lines, byte-exact)

The expected byte form is captured in `expected_sign_payload.txt` (the canonical bytes verbatim, terminated with NO trailing newline) and `expected_sign_payload_sha256.txt` (the SHA-256 of those bytes, for quick verification without comparing 287 bytes by eye).

| Property | Value |
|---|---|
| sign_payload byte length | `287` |
| sign_payload SHA-256 | `45cd0389edbdea524d2b36f2e3b60cec4de502dcc3aac67269ed004ebcb9ac57` |

The byte sequence is the magic line followed by eleven fields, each separated by `0x0A`, with no trailing newline:

```
ffiec.chain-of-custody.v1\n   ← magic line
v1.0b\n                       ← sign_payload_version (line 2)
ed25519\n                     ← algorithm
v1\n                          ← format_version
demo-bank-streaming-001\n     ← tenant_id
2026-05-07\n                  ← seal_date
0001…3031\n                   ← hex(merkle_root)  (64 chars)
ffee…1100\n                   ← hex(hkdf_inputs_digest)  (64 chars)
per_second\n                  ← cadence (§10.27)
0\n                           ← dev_mode (single ASCII byte)
1\n                           ← key_versions_canon
ccbc…87f8                     ← hex(kms_handle_uris_digest)  (64 chars, terminal — no trailing \n)
```

## Conformance behavior

A conforming implementation:

1. Constructs the v1.0b 12-line `sign_payload` byte-for-byte identical to `expected_sign_payload.txt` for the inputs in `input.json`.
2. Computes a SHA-256 of the constructed bytes that matches `expected_sign_payload_sha256.txt`.
3. Refuses any seal record whose `cadence` field is not in the §10.27 enumeration with reason `cadence "X" is not in the §10.27 enumeration` (§7 step 12 sub-case).
4. Refuses any chain that mixes `per_hour` and `hourly` across adjacent seals with reason `cadence form changed mid-chain at seal_date {D}`.

## Negative cases this fixture supports

A v1.0b verifier presented with any of the following produces `signature verification failed` (because the reconstructed `sign_payload` no longer matches what was signed):

- Tampered `seal.cadence` value (e.g., the seal claims `daily` but the signed bytes used `per_second`) — under v1.0b the value bound into the cadence line produces a signature mismatch.
- A v1.0a verifier reading this v1.0b fixture without the v1.0b dispatch fails fast at §7 step 11 with `sign_payload_version "v1.0b" not supported by this verifier (running v1.0a)`.

## Cross-references

- Spec §10.27 streaming cadence enumeration (the `per_second` value)
- Spec §4.2.1 cadence field (the original enumeration §10.27 extends)
- Spec §4.3 amendment form (v1.0b, normative — the 12-line byte structure)
- Spec §7 step 11 verifier dispatch (sign_payload_version routing)
- Spec §7 step 12 cadence-and-dev-mode check (cadence-value-not-in-enum reason)
- Case 018 (`018-sign-payload-v1.0b/`) — v1.0b mixed key_versions / kms_handle_uris (cadence = `daily`)
- Case 019 (`019-sign-payload-v1.0b-empty-day/`) — v1.0b empty-day variant
- §10.28 streaming-mode IKM rotation discipline (the cadence-interval boundary §10.27 streaming values trigger)
- §10.29 streaming-mode verifier procedure (the verifier behaviour for streaming cadence)
- §10.30 trusted-time integration (normative for streaming cadence)

## Reproduction

The `_compute.py` script in this directory regenerates `input.json`, `expected_sign_payload.txt`, and `expected_sign_payload_sha256.txt` from the pinned inputs above. It uses only the Python standard library; it does not depend on any specific SDK. Re-run after any spec change that affects the v1.0b byte form:

```
python _compute.py
```

Cross-validation: the same inputs run through Herald.Py's production `herald._crypto.sign_payload.build_sign_payload` produce the same 287 bytes and the same SHA-256. A future implementation in any language that disagrees on these bytes is non-conformant.
