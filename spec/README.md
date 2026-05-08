# Specification

> **The spec is the standard.** Implementations conform to the spec. The reference implementation in `core/`, `ledger/`, and `verifier/` is one example of a conforming implementation; it is not the standard itself.

## Versioning

Specifications follow the form `chain-of-custody-vN.M.md`:

- `N` major version &mdash; incompatible wire-format or primitive changes
- `M` minor version &mdash; additive changes that an older verifier can ignore safely

## Current version

| Spec | Version | Status |
|---|---|---|
| [`chain-of-custody-v1.md`](chain-of-custody-v1.md) | v1.0-final | Issued |

## Conformance

A conforming implementation:

1. Produces output that passes every test vector in [`test-vectors/`](test-vectors/) for its declared spec version.
2. Accepts as input any output produced by another conforming implementation of the same spec version.
3. Documents which spec version it implements in its release notes and `--version` output.

The reference verifier in [`../verifier/`](../verifier/) reports the spec version it validates against. An institution running an SDK from one vendor and a ledger from another vendor and a verifier from this repository should see byte-for-byte agreement on every chain artifact.

## Test vectors

[`test-vectors/`](test-vectors/) contains the canonical conformance corpus. Every primitive change ships test vectors that demonstrate the change. The corpus is the discriminator between conforming and non-conforming implementations.
