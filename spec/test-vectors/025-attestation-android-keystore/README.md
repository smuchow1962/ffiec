# Case 025 — Attestation envelope, Android Keystore platform (§10.35)

## What this case verifies

Spec §10.35 normates the `ffiec.chain.attestation` attribute on chain entries plus a delegate-dispatch validator pattern that platform-specific attestation parsers plug into at chain-walk time. The attestation envelope schema is byte-locked across implementations:

```json
{
  "platform": "<one of the §10.35 enumeration>",
  "attestation_document_b64": "<base64-encoded platform-specific attestation, ≤ 64 KiB>",
  "root_certificate_id": "<institution-named identifier of the platform-vendor root, ≤ 256 chars>",
  "public_key_fingerprint": "<optional 64-char lowercase hex SHA-256 of device public key>"
}
```

This case pins the JCS-canonical (RFC 8785) bytes for an Android Keystore attestation envelope so a clean-room implementation proves it constructs the §10.35 envelope identically to other implementations.

## What this case does NOT verify

The attestation **document body** is a synthetic placeholder. Real Android Keystore attestations are CBOR-encoded and X.509-rooted at Google's Hardware-Backed Keys CA; parsing them is the platform validator's job (the delegate the dispatcher invokes per §10.35). Case 025 pins the **envelope** byte form — the wrapper around the document — not the document's internal structure. A future case (or a Herald-side test) can pin the platform-validator dispatch behaviour with a real attestation chain.

## Inputs

| Field | Value |
|---|---|
| `platform` | `android_keystore` (§10.35 enumeration) |
| `attestation_document_b64` | 344-char base64 (256 bytes of synthetic document; recipe in `input.json`) |
| `root_certificate_id` | `google.hardware-backed-keys.root.2023` |
| `public_key_fingerprint` | `SHA-256(device_pubkey)` — 64 lowercase hex chars |

The synthetic 256-byte attestation document is deterministic — the recipe is `seed = SHA-256(utf8(label)); body = seed || SHA-256(seed) || SHA-256(SHA-256(seed)) || …` truncated to 256 bytes. Any reader can recompute the exact bytes from the label `025-android-keystore-synthetic-document-v1`.

## Expected canonical bytes

| Property | Value |
|---|---|
| envelope canonical byte length | `559` |
| envelope canonical SHA-256 | `1c12aa015bd3534f45e1e093fb484dab9c86584c1ca1580dd7a236a712b3016d` |

`expected_canonical.txt` carries the bytes verbatim; `expected_canonical_sha256.txt` is the SHA-256 for quick comparison.

JCS canonicalisation sorts object keys lexicographically. The envelope's keys appear in canonical order: `attestation_document_b64`, `platform`, `public_key_fingerprint`, `root_certificate_id`.

## Conformance behavior

A conforming implementation supporting §10.35:

1. Constructs the attestation envelope as a JSON object with the four fields above.
2. Canonicalises per RFC 8785 (JCS) and produces bytes byte-identical to `expected_canonical.txt`.
3. Refuses any envelope whose `platform` is outside the §10.35 enumeration with a control-completeness failure (the institution's CC8.1 names the supported platforms; an envelope outside that set indicates a configuration drift).
4. Bounds `attestation_document_b64` at ≤ 64 KiB and `root_certificate_id` at ≤ 256 chars per §10.35 normative.
5. When `public_key_fingerprint` is present, exactly 64 lowercase hex chars (a SHA-256 of the device's public key bytes).

## Delegate-dispatch behaviour (informative)

§10.35 normates that platform validators are registered against the `platform` field at construct time. The verifier dispatcher reads `envelope["platform"]`, looks up the registered validator, and invokes it with `(attestation_document_b64, root_certificate_id)`. Case 025 does NOT pin the dispatcher — that's a control-flow contract and varies by implementation language. The case pins the **envelope bytes** the dispatcher's caller produces and the validator's caller consumes.

## Negative cases this fixture supports

- A producer that emits `platform` outside the §10.35 enumeration (e.g., `"my_custom_attestation"`) produces a structurally-valid envelope whose dispatcher lookup fails — case 025 surfaces the byte-form divergence; the dispatcher behaviour is institution-specific.
- A producer that adds extra envelope fields (e.g., a custom `vendor_extension`) produces canonical bytes that diverge from `expected_canonical.txt` because JCS includes the extra field.
- A producer that emits `public_key_fingerprint` with uppercase hex chars produces canonical bytes that diverge — §10.35 normates "lowercase hex".
- A producer that emits `attestation_document_b64` with whitespace in the base64 (PEM-style line wrapping) produces canonical bytes that diverge — base64 in JSON is whitespace-free.

## Cross-references

- Spec §10.35 edge-attestation primitive (the section this case pins)
- Spec §10.32 per-device session-key derivation (the natural compositional pair — §10.32 binds device_id in the HKDF info; §10.35 binds the device's hardware identity to the chain entry)
- Spec §10.29 streaming-mode verifier (the delegate-dispatch validator pattern §10.35 mirrors)
- RFC 8785 (JCS canonicalisation)
- Case 008 (`008-jcs-edge-cases/`) — the JCS conformance bar
- Case 024 (`024-per-device-derivation/`) — the §10.32 device-binding compositional partner

## Reproduction

The `_compute.py` script in this directory regenerates `input.json`, `expected_canonical.txt`, and `expected_canonical_sha256.txt` from the pinned inputs above. It depends on the Python `jcs` package (Anders Rundgren's reference RFC 8785 implementation, the same package case 008 uses).

```
python _compute.py
```
