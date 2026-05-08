# Case 019 — sign_payload v1.0b 12-line form (empty-day edge case)

## What this case verifies

The v1.0b 12-line `sign_payload` byte form on an empty-day seal — a tenant-day with zero chain entries. Empty-day seal continuity is required per §4.2 (every tenant-day MUST receive a seal record, including tenant-days with zero events). Under v1.0b the two new terminal fields take their degenerate forms:

- `key_versions_canon` is the empty string (zero bytes between the two surrounding `\n` separators on lines 11 and 12)
- `hex(kms_handle_uris_digest)` is `SHA-256(b"")` = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`

Case 019 pins the byte form for both degenerate values so an implementation that mishandles the empty case (e.g., omits the empty `key_versions_canon` line entirely, producing an 11-line form, or substitutes a non-empty placeholder string) is detected.

## Inputs

| Field | Value |
|---|---|
| `algorithm` | `ed25519` |
| `format_version` | `v1` |
| `tenant_id` | `demo-bank-001` |
| `seal_date` | `2026-05-08` |
| `merkle_root` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (the empty-tree Merkle root per §4.2 — `SHA-256(b"")`) |
| `hkdf_inputs_digest` | `1122334455667788990011223344556677889900112233445566778899001122` (placeholder; same shape as case 018) |
| `cadence` | `daily` |
| `dev_mode` | `false` (serialized as `0`) |
| `key_versions` (distinct) | `[]` (empty — no chain entries on the day) |
| `kms_handle_uris` (distinct) | `[]` (empty) |

## Computed canonical fields

`key_versions_canon` is the empty string when no `key_version` values are present:

```
key_versions_canon = ""
```

`kms_handle_uris_digest` is `SHA-256(b"")` when no URIs are present:

```
canonical_uri_form = ""  (zero bytes)
SHA-256(canonical_uri_form) = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

## Expected sign_payload (12 lines, byte-exact)

The expected byte form is captured in `expected_sign_payload.txt`. Note that line 11 is empty — the bytes between the `\n` after `0` (line 10's terminator) and the `\n` before the digest (line 11's terminator) are exactly zero bytes. An implementation that normalizes empty lines to `<empty>`, `null`, or any non-zero placeholder produces a different byte sequence and a different signature.

| Property | Value |
|---|---|
| sign_payload byte length | `271` |
| sign_payload SHA-256 | `7db040483d03b62df1ba43bf6128fbbf4ddbbeb67b03551cecd6b63a1ac7ff72` |

## Conformance behavior

A conforming implementation:

1. Constructs the v1.0b 12-line `sign_payload` byte-for-byte identical to `expected_sign_payload.txt` for the empty-day inputs.
2. Computes a SHA-256 of the constructed bytes that matches `expected_sign_payload_sha256.txt`.
3. On verification, reconstructs `sign_payload` from the empty-day seal record using the same canonicalization and obtains the same byte sequence.

## Negative cases this fixture supports

A v1.0b verifier presented with any of the following fails (because the reconstructed `sign_payload` no longer matches what was signed):

- An empty-day seal that omits the empty `key_versions_canon` line (i.e., an 11-line form): the byte form differs by one `\n` byte; signature fails.
- An empty-day seal whose `key_versions_canon` is `"0"`, `"none"`, `"null"`, or any non-empty string: the byte form differs; signature fails.
- An empty-day seal whose `kms_handle_uris_digest` is anything other than `e3b0c44...b855`: signature fails.

## Cross-references

- Spec §4.2 empty-day Merkle root (`SHA-256(b"")` pinned)
- Spec §4.3 amendment form (v1.0b, normative) — empty-day discipline for `key_versions_canon` and `kms_handle_uris_digest`
- Spec §12 change log (v1.0b row records the empty-day binding)
- Case 018 (`018-sign-payload-v1.0b/`) covers a non-trivial mixed-version mixed-URI day for contrast
