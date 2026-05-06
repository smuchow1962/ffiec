# Contributing to the FFIEC AI Chain-of-Custody project

Thank you for considering a contribution. This project is a **regulator-facing reference implementation**, which means contributions are evaluated against a higher correctness-and-security bar than typical open-source software. Please read this document before submitting changes.

## The auditor's lens

Every change is reviewed by the auditor's lens. That means a reviewer will ask, *"Could a Big Four auditor or a federal examiner walk this chain independently and verify it?"* If the answer is no, the change does not land. We treat the auditor as our advisor, not as a downstream consumer.

In practice, the auditor's lens means:

- **Cryptographic primitives** must be FIPS-approved or NIST-recommended. No experimental algorithms. No custom constructions.
- **Determinism** is non-negotiable. The same inputs produce the same outputs. Every time. On every platform.
- **Independence** is non-negotiable. The verifier must validate the chain without trusting the institution, the vendor, or any party except the FFIEC's published key registry.
- **Test vectors** accompany every primitive change. A change that does not extend the test corpus is not a complete change.
- **Documentation** is part of the deliverable. Code without a corresponding update to the spec or the design docs is incomplete.

## Reporting issues

- **Bugs** &mdash; open a [bug report](.github/ISSUE_TEMPLATE/bug_report.md). Include a reproducible test case.
- **Security vulnerabilities** &mdash; do **not** open a public issue. Follow [SECURITY.md](SECURITY.md).
- **Feature requests** &mdash; open a [feature request](.github/ISSUE_TEMPLATE/feature_request.md). Explain the regulatory or operational driver, not just the desired behavior.

## Pull requests

1. **Fork and branch.** Branch off `main` with a descriptive name (`fix/merkle-leaf-ordering`, `spec/v1.1-tenant-key-rotation`).
2. **Sign your commits.** All commits must be signed (`git commit -S`). Unsigned commits will not be merged.
3. **Run the conformance suite.** `go test ./...` must pass. The conformance test vectors in `spec/test-vectors/` must continue to pass.
4. **Update the spec or the design doc** if you change behavior. Behavior changes that are not reflected in the normative documents will be rejected.
5. **Open the PR with the auditor's-lens summary** &mdash; one paragraph explaining how a reviewer can verify the change without re-running it themselves.

## Code review

Two reviewers are required for changes that touch:

- The `core/chain/`, `core/merkle/`, or `core/hsm/` packages
- Anything in `spec/`
- Test vectors in `spec/test-vectors/`

One reviewer is sufficient for changes to `examples/`, `docs/` (other than `docs/design/`), GitHub workflows, or build infrastructure.

## Scope

This project does not accept:

- New AI capability features (the project is *capture and verify*, not *do AI things*)
- Vendor-specific integrations not aligned with OTLP semantic conventions (file an upstream OTel issue first)
- Cryptographic experiments or custom constructions (use FIPS-approved primitives)
- Telemetry or phone-home behavior of any kind
- Closed-source dependencies

This project does accept:

- Conformance test vectors that exercise edge cases
- Documentation improvements
- Performance improvements that preserve byte-for-byte compatibility
- Additional language SDKs in their own repositories that conform to the spec

## License

By contributing, you agree that your contribution is licensed under the Apache License 2.0 as set out in [LICENSE](LICENSE). No CLA is required; the Apache 2.0 inbound-equals-outbound model governs.
