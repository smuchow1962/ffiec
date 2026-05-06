# `verifier` — offline standalone examiner CLI

Single static binary an examiner runs against a ledger to verify chain integrity, Merkle proofs, and HSM root signatures &mdash; entirely offline, with no dependencies on the institution or its vendor.

## Build

```bash
CGO_ENABLED=0 go build -o /tmp/verifier ./cmd/verifier

# Cross-compile for examiner laptops
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/verifier.exe ./cmd/verifier
GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -o /tmp/verifier-mac ./cmd/verifier
GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/verifier-linux ./cmd/verifier
```

## Usage

```bash
verifier verify \
  --ledger      /path/to/tenant-day.ledger \
  --root-key    /path/to/tenant.pub \
  --report      /path/to/report.pdf

verifier walk \
  --ledger      /path/to/tenant-day.ledger \
  --run-id      r_a3f29b71c

verifier diff \
  --before      /path/to/before.ledger \
  --after       /path/to/after.ledger
```

## Auditor's-lens guarantees

- **Offline.** The verifier never makes a network call. No DNS, no HTTPS, no telemetry. Verified by network-policy at build time (no `net/http`, no `net/url` against external hosts).
- **Single binary.** No dynamic linking. No external configuration files except those passed on the command line. The binary plus its arguments is the entire trusted-computing-base for verification.
- **Deterministic output.** A given ledger plus a given root-key produces a byte-for-byte identical PDF report. Two examiners run the verifier independently and compare reports by hash.
- **Failure surfaces.** Every failure is loud, specific, and actionable. The verifier never silently passes a partial verification. If it cannot reach a conclusion, it fails closed.
- **Independent of the institution.** The verifier reads only the ledger and the public-key file. It does not call the institution's vendor, the institution's HSM, or any third party. The institution's only role is to hand the examiner the ledger.
