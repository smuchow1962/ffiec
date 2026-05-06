# `core` — chain-of-custody primitive library

Pure Go implementation of the four primitives defined in [`spec/chain-of-custody-v1.md`](../spec/). No I/O. No telemetry. Audited cryptography only (Go standard library `crypto/hmac`, `crypto/sha256`, `crypto/ed25519`).

## Packages

| Package | Purpose |
|---|---|
| [`chain/`](chain/) | HMAC chain construction at capture. Per-process session key via HKDF; per-event payload hashed with SHA-256. |
| [`merkle/`](merkle/) | Daily Merkle seal. Aggregates events into a binary Merkle tree; root is the daily seal. |
| [`hsm/`](hsm/) | HSM signing interface. Abstract over PKCS#11, AWS CloudHSM, Azure Key Vault HSM. Reference impl uses an in-memory Ed25519 keypair for testing only. |
| [`otlp/`](otlp/) | Encode and decode OTLP envelopes carrying chain extension fields. Aligns with OpenTelemetry GenAI semantic conventions. |

## Status

Design phase. The package directories exist; implementation lands once the spec is at v1.0-draft.

## Build

```bash
go build ./...
go test -race ./...
```

## Auditor's-lens guarantees

- **Determinism.** Given the same input, every primitive produces byte-for-byte identical output across platforms.
- **No I/O.** This package never opens a file, makes a network call, or accesses a secret store. All inputs are explicit function parameters.
- **No global state.** Every operation takes its keys, prev-hashes, and clock as explicit parameters.
- **Stdlib crypto only.** No third-party cryptographic dependencies. The auditor verifies against `crypto/*` upstream, not against this repository.
