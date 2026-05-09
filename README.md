# FFIEC AI Chain-of-Custody Reference Implementation

> A proposed standard for tamper-evident logging of AI-driven decisions in regulated financial institutions, with a working open-source reference implementation.

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8.svg)](https://go.dev)
[![Status](https://img.shields.io/badge/status-design%20phase-orange.svg)](docs/design/)

---

## What this is

The FFIEC IT Examination Handbook already mandates that audit logs be **tamper-evident, integrity-protected, immutable, and complete** (Information Security booklet, AIO booklet). The AIO booklet's §VII.D names the risks of AI in production but contains no logging or audit-trail procedure for AI activity. This project closes the gap.

Four primitives, language-neutral, vendor-neutral:

1. **HMAC chain at capture** &mdash; every AI decision event is hashed with a per-process session key derived via HKDF; each event includes the SHA-256 of the previous event in the same run.
2. **Daily Merkle seal** &mdash; events for each tenant-day are aggregated into a Merkle tree; the root is published in an append-only ledger.
3. **HSM-rooted root signature** &mdash; the daily Merkle root is signed in HSM custody (FIPS 140-2 Level 3).
4. **OpenTelemetry-native wire** &mdash; events ship over OTLP using standard OTel attributes plus the chain extension fields.

The four primitives compose into a chain of custody an examiner can independently walk &mdash; without trusting the institution&rsquo;s vendor or the institution&rsquo;s ops team.

## Repositories

This is a monorepo with three Go modules:

| Path | Purpose |
|---|---|
| [`core/`](core/) | The chain-of-custody primitive library. Pure Go. No I/O. Audited cryptography only. |
| [`ledger/`](ledger/) | Reference ingest server. Receives OTLP, verifies HMAC, writes append-only ledger, computes daily Merkle seal, signs in HSM. |
| [`verifier/`](verifier/) | Standalone offline CLI for examiners. Walks the chain, verifies every HMAC, validates the daily Merkle proof, validates the HSM signature. |

The chain-of-custody **specification**, conformance test vectors, governance, and all authoritative project documentation live in a separate repo, [`ffiec-public`](https://github.com/smuchow1962/ffiec-chain-of-custody). That is the canonical source.

The [`docs/`](docs/) tree in this repo is **deprecated** — preserved as a historical snapshot but no longer authoritative. See [`docs/INDEX.md`](docs/INDEX.md) for the deprecation banner.

## Quick start

```bash
# Build all three binaries
go build ./...

# Run the ledger server locally with default config
./ledger/cmd/ledger -config ./ledger/deploy/config.dev.yaml

# Verify a sample chain offline
./verifier/cmd/verifier -ledger ./testdata/sample-day.ledger -root-key ./testdata/sample.pub

# Run the conformance test suite against the spec
go test ./spec/test-vectors/...
```

## Status

**Design phase.** The specification and the design documents are the deliverable through Q3 2026. Implementation of `core/`, `ledger/`, and `verifier/` follows once the spec is locked.

The roadmap and design tracking is in [`docs/design/`](docs/design/).

## Audience

- **The FFIEC and its member agencies** &mdash; the chain-of-custody primitive is offered as a candidate standard for the next AIO booklet revision or a stand-alone AI examiner procedure.
- **Regulated financial institutions** &mdash; the reference implementation runs entirely inside the institution&rsquo;s perimeter; no telemetry phone-home.
- **AI-tooling vendors** &mdash; the spec is the conformance target. Compete on quality of implementation, not on the standard.
- **External auditors and internal audit teams** &mdash; the verifier is the independent verification path.

## Governance

This project is governed by the rules in [GOVERNANCE.md](GOVERNANCE.md). The intent is to transfer governance to a neutral foundation (CNCF, OpenSSF, or a banking-industry consortium) once the specification is mature and adoption is established.

## Licensing

Apache License 2.0. The Apache patent grant matters here: the chain-of-custody primitive cannot be hostage to a future patent claim by a contributor or a downstream vendor. See [LICENSE](LICENSE).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Security issues: see [SECURITY.md](SECURITY.md).

## Disclaimer

This project is not affiliated with, endorsed by, or sponsored by the FFIEC or any of its member agencies. The repository name reflects the regulatory framework this work is designed to satisfy. The FFIEC&rsquo;s adoption of any standard is solely the FFIEC&rsquo;s decision.
