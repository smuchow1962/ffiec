# `ledger` — reference ingest server

Receives OTLP from `core/`-compatible SDKs, verifies HMAC chains on ingest, writes to an append-only ledger, computes the daily Merkle seal, and signs it in HSM custody.

## Architecture

```
OTLP receiver (gRPC + HTTP)
    →  HMAC chain verifier
    →  Append-only ledger writer
    →  Daily Merkle seal job
    →  HSM signing
    →  Published signed root
```

## Build

```bash
go build -o /tmp/ledger ./cmd/ledger
./tmp/ledger -config ./deploy/config.dev.yaml
```

## Configuration

See [`deploy/`](deploy/) for sample configurations. The server reads:

- **OTLP endpoints** &mdash; gRPC port and HTTP port
- **Ledger storage** &mdash; append-only directory or PostgreSQL DSN
- **HSM configuration** &mdash; PKCS#11 module path or cloud-HSM credentials
- **Tenant key registry** &mdash; HMAC master key location for chain verification

## Auditor's-lens guarantees

- **Append-only.** Ledger storage is append-only by design. No DELETE or UPDATE statements anywhere in the writer path. Verified by linting and runtime assertion.
- **Chain re-verification at ingest.** Every event's HMAC is re-verified at the receiver before it is committed to the ledger. Failed verifications are recorded but do not corrupt the ledger.
- **Deterministic Merkle.** The daily seal computation is byte-for-byte deterministic. Two implementations of the spec producing the same set of leaves produce the same root.
- **HSM signing only.** The reference implementation refuses to sign with a software-backed key in production mode. Software keys exist for development and CI only and produce a clearly-marked &ldquo;dev-mode&rdquo; signature.
