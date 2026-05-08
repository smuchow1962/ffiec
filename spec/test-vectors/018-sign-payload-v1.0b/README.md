# Case 018 — sign_payload v1.0b 12-line form (mixed key_versions, mixed kms_handle_uris)

## What this case verifies

The v1.0b 12-line `sign_payload` byte form per spec §4.3 amendment form (v1.0b, normative). The case exercises the two new terminal lines that close NIST-G1 (`kms_handle_uri` provenance binding) and NIST-G2 (`key_versions` cross-check binding):

- `key_versions_canon` — the deterministic comma-separated ascending decimal encoding of the day's distinct `key_version` integer values
- `hex(kms_handle_uris_digest)` — the 64-character lowercase hex of the 32-byte SHA-256 over the canonical sorted-distinct-URI form

The fixture covers a non-trivial mixed-version, mixed-URI day so that both new fields carry substantive content (rather than degenerate-empty values, which case 019 covers separately).

## Inputs

The seal-day under construction has three distinct `key_version` values (1, 2, 3) across its chain entries and two distinct `kms_handle_uri` values. The remaining fields use the same shape as v1.0a fixtures.

| Field | Value |
|---|---|
| `algorithm` | `ed25519` |
| `format_version` | `v1` |
| `tenant_id` | `demo-bank-001` |
| `seal_date` | `2026-05-07` |
| `merkle_root` | `aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899` (placeholder; the v1.0b form is independent of the merkle_root value, so any 64-char lowercase hex works for the canonical-bytes test) |
| `hkdf_inputs_digest` | `1122334455667788990011223344556677889900112233445566778899001122` (placeholder; same rationale) |
| `cadence` | `daily` |
| `dev_mode` | `false` (serialized as the single ASCII byte `0`) |
| `key_versions` (distinct) | `{1, 2, 3}` |
| `kms_handle_uris` (distinct) | `aws-kms:arn:aws:kms:us-east-1:123456789012:key/4f2a-mock-1`, `aws-kms:arn:aws:kms:us-east-1:123456789012:key/9b3c-mock-2` |

## Computed canonical fields

`key_versions_canon` is computed by sorting the distinct `key_version` integer values into ascending order and joining with single ASCII commas:

```
key_versions_canon = "1,2,3"
```

`kms_handle_uris_digest` is computed by sorting the distinct URI values into code-point ascending order, joining with single `\n` (0x0A) bytes (no trailing newline), and applying `SHA-256`:

```
canonical_uri_form = "aws-kms:arn:aws:kms:us-east-1:123456789012:key/4f2a-mock-1\naws-kms:arn:aws:kms:us-east-1:123456789012:key/9b3c-mock-2"
canonical_uri_form_bytes_len = 117
SHA-256(canonical_uri_form) = 9e5b3b8f18b428e7cda7b3270673eb2a7dd044a52059102d2d38d16a8d1972b7
```

## Expected sign_payload (12 lines, byte-exact)

The expected byte form is captured in `expected_sign_payload.txt` (the canonical bytes verbatim, terminated with NO trailing newline) and `expected_sign_payload_sha256.txt` (the SHA-256 of those bytes, for quick verification without comparing 276 bytes by eye).

| Property | Value |
|---|---|
| sign_payload byte length | `276` |
| sign_payload SHA-256 | `f20fb9942bc8336c89647ac2977d06c942659290703eae66ee40f4067924cbed` |

## Conformance behavior

A conforming implementation:

1. Constructs the v1.0b 12-line `sign_payload` byte-for-byte identical to `expected_sign_payload.txt` for the inputs in `input.json`.
2. Computes a SHA-256 of the constructed bytes that matches `expected_sign_payload_sha256.txt`.
3. Signs `sign_payload` with the institution's HSM-held Ed25519 key (signature bytes are not pinned in this fixture because Ed25519 is not deterministic across implementations under all libraries; the `sign_payload` byte form is the conformance witness).
4. On verification, reconstructs `sign_payload` from the seal record and the day's events using the same canonicalization and obtains the same byte sequence and SHA-256.

## Negative cases this fixture supports

A v1.0b verifier presented with any of the following produces `signature verification failed` (because the reconstructed `sign_payload` no longer matches what was signed):

- Tampered `seal.key_versions` value (e.g., the seal claims `[1, 2]` but the day's entries carry versions `{1, 2, 3}`) — under v1.0b the tampered value bound into `key_versions_canon` produces a signature mismatch; the cross-check (§7 step 11 normative cross-check paragraph) also catches this case as defense-in-depth.
- Tampered `kms_handle_uri` value on any per-event entry (a flip from the legitimate URI to a different URI changes the canonical sorted-distinct-URI form, the digest, and the signature input).
- A v1.0a verifier reading this v1.0b fixture without the v1.0b dispatch fails fast at §7 step 11 with `sign_payload_version "v1.0b" not supported by this verifier (running v1.0a)` rather than reconstructing the wrong byte form and reporting a generic signature failure.

## Cross-references

- Spec §4.3 amendment form (v1.0b, normative)
- Spec §7 step 11 verifier dispatch (v1.0b case added)
- Spec §12 change log (v1.0b row records the wire-form extension)
- Case 019 (`019-sign-payload-v1.0b-empty-day/`) covers the empty-day edge case where both new fields take their degenerate forms
