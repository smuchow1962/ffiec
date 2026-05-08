# Case 001 — Single event, empty `prev_hash`

## Purpose

The smallest case. Verifies that `seq=1` correctly uses 32 zero bytes as `prev_hash`, that HKDF derives the expected session key, that the per-event MAC produces the expected `payload_hash`, and that the file header carries the right `hkdf_inputs_digest`.

A conforming implementation that fails this case has a fundamental construction bug.

## Inputs

From `../chain_vectors.json` `inputs` block:

```
HKDF_SALT      = b"ffiec.chain-of-custody.v1.salt"
HKDF_INFO_BASE = b"ffiec.chain-of-custody.v1.info"
tenant_id      = "tenant-ffiec-test-1"
run_id         = "RUN-FFIEC-VECTOR"
ikm_v1         = bytes from ikm_v1_hex (48 bytes)
seal_date      = "2026-05-06"
```

The single event is `single_chain[0]` from `chain_vectors.json`:

- `seq = 1`
- `prev_hash = 0x00 * 32`
- `event payload`: the canonical bytes are `single_chain[0].event_canonical_hex` (decoded from hex)
- `key_version = 1`

## Expected outputs

| Field | Expected value (from `chain_vectors.json`) |
|---|---|
| `session_key_v1_hex` | `f061210167d307cffe4b91c2eaad9d386a8831c04957b458088e5bafd1f63ee9` |
| `key_fingerprint_v1_hex` | `54c245a894b81d788d372badc563531f` |
| `hkdf_inputs_digest_hex` | `6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65` |
| `payload_hash` for seq=1 | `e84168071562a9e866f28977a204b43abfdc959b7f84aad83c16fe790f63b049` |

## What this case proves

- HKDF-SHA-256 over the FFIEC constants produces the expected session key.
- The `info` parameter binds tenant_id correctly (`info = info_base || "|" || utf8(tenant_id)`).
- The `key_fingerprint` is the per-tenant SHA-256 truncation.
- The first event's `prev_hash` is 32 zero bytes.
- The MAC input is `prev_hash || canonical_bytes` with `prev_hash` fixed at 32 bytes.
- The MAC output is the persisted `payload_hash` byte-for-byte (no SHA-256-of-payload substitution).

## Failure mode

If `payload_hash` differs from the expected, suspect (in order):

1. Wrong HKDF inputs (check salt, info_for_tenant, length).
2. Wrong canonical-form encoding (JCS sorted keys, no whitespace, UTF-8).
3. Chain-stamp fields included in the canonical bytes (must be excluded).
4. `prev_hash` not exactly 32 zero bytes for `seq=1`.
5. HMAC-SHA-256 implementation bug (rare; stdlib in every target language).
