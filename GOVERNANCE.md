# Project governance

This project is a **proposed open standard**. Governance reflects that posture &mdash; decisions about the specification are made transparently, with the expectation that the project will eventually be transferred to a neutral foundation.

## Roles

- **Maintainers.** Day-to-day stewardship, code review, release management. Listed in [.github/CODEOWNERS](.github/CODEOWNERS).
- **Spec editors.** Authority over the specification in `spec/`. A spec change requires unanimous spec-editor approval.
- **Security response team.** Triage and remediation of security reports. Members listed privately; contact via [SECURITY.md](SECURITY.md).
- **Contributors.** Anyone who has had a pull request merged.

## Decision-making

| Change | Decision rule |
|---|---|
| Bug fix in implementation, no spec change | One maintainer approval |
| Implementation behavior change, spec-conforming | Two maintainer approvals |
| Non-normative documentation change | One maintainer approval |
| Specification change (spec/) | **Unanimous spec-editor approval** + 14-day public comment period |
| Security release | Two maintainer approvals + security response team approval |
| Governance change (this document) | Unanimous maintainer approval + 30-day public comment period |
| Foundation transfer | Unanimous maintainer approval + announced 90 days in advance |

## Specification changes

Specification changes are the most consequential. They affect every conforming implementation. The process:

1. **Proposal** &mdash; opened as an issue tagged `spec-proposal`, with rationale citing FFIEC handbook language, regulatory drivers, or threat-model gaps.
2. **Discussion** &mdash; minimum 14 days. All parties (vendors, banks, examiners, contributors) may comment.
3. **Pull request** &mdash; the proposal is converted to a PR against `spec/`. Includes draft test vectors that demonstrate the change.
4. **Review** &mdash; spec editors review for normative correctness and backward compatibility. Reviewers from the security response team review for cryptographic and integrity implications.
5. **Public comment** &mdash; 14 days for the broader community to comment on the PR.
6. **Merge** &mdash; unanimous spec-editor approval merges the change. The spec version increments. Release notes call out the change explicitly.

Specification changes are **versioned**. The specification at HEAD on `main` may differ from the most recent published spec version. Implementations are pinned to a spec version.

## Foundation transfer

The intent of this project is to be transferred to a neutral foundation once:

- The specification has reached version 1.0 with stable adoption
- At least three independent conforming implementations exist
- The FFIEC has named the chain-of-custody primitive in published guidance, or a sister-agency has done so

Likely targets, in order of fit:

1. **Open Source Security Foundation (OpenSSF)** under the Linux Foundation &mdash; security-aligned, regulator-friendly governance
2. **Cloud Native Computing Foundation (CNCF)** &mdash; if the project anchors more in OTel and cloud-native infrastructure
3. **A new Linux Foundation working group on regulated AI** &mdash; if such a group is formed
4. **A banking-industry consortium** (FS-ISAC, BPI) &mdash; if the FFIEC adoption is well-defined and bank-specific governance is appropriate

The foundation transfer is the project&rsquo;s **end state for governance**. The maintainers&rsquo; job is to make that transfer easy, by keeping the project open, well-documented, and free of vendor lock-in.

## What governance does not do

- **Endorse vendors.** Conforming implementations are listed in `docs/conforming-implementations.md` (when it exists) without ranking.
- **Certify compliance.** This project produces a specification and a reference implementation. Conformance certification, if any, is the FFIEC&rsquo;s decision, not ours.
- **Accept paid-for spec changes.** Specification changes are made on technical merit and regulatory alignment alone.

## Vendor-conformance attestation registry

The project operates a public registry of vendor-conformance attestations &mdash; signed attestations from vendors that their chain-of-custody implementations pass the FFIEC conformance corpus published in `spec/test-vectors/`. The registry is the project-side trust mechanism that complements vendor SOC reporting (which covers vendor operational controls but does NOT attest implementation conformance against the corpus). Institutions consume the registry as part of CC8.1 vendor-management evidence; the procedure is documented in [`docs/vendor-conformance-attestation.md`](docs/vendor-conformance-attestation.md).

### Sub-committee operation

A **vendor-conformance sub-committee** within the project's working group operates the registry. The sub-committee comprises at minimum three project maintainers, separated from the spec-editor role to keep decision-making distributed. Sub-committee membership rotates per the project's standard maintainer-rotation discipline. Decisions are logged in public issue tracking with rationale published.

### Project-side commitments

The project commits to the following operational cadences for the vendor-conformance attestation procedure:

| Commitment | Cadence |
|---|---|
| Registry review (sub-committee confirms each registry entry's URLs are reachable, pending re-attestations are tracked, the revocation log is current) | Quarterly |
| Submission acknowledgement after a vendor submits an attestation | 5 business days |
| Submission review completion after acknowledgement | 30 calendar days |
| Revocation publication after sub-committee revocation decision | 14 calendar days |
| Vendor re-attestation grace period after a new corpus version publishes | 90 calendar days |
| Quarterly registry-review summary publication (signed by release-management role) | Within 14 calendar days of quarter end |

These commitments are normative project-side governance. The sub-committee operates them as standing obligations; institutions consuming the registry can rely on the cadences when scheduling vendor-management evidence cycles. Failure to meet a commitment is escalated to the maintainer group per the project's standard escalation discipline; persistent failure is a foundation-transfer-readiness concern noted in the project's standard public reporting.

### Corpus-version update coordination

The FFIEC conformance corpus updates per the "Specification changes" process above. When a new corpus version publishes, the working group announces the corpus version with the 90-day vendor re-attestation grace period. During the grace period, prior-corpus attestations remain "Active" in the registry with a "re-attestation pending" annotation; after the grace period expires, the sub-committee revokes the prior-corpus attestation per the documented revocation procedure if no re-attestation has been submitted. The 90-day grace period is calibrated against typical vendor product-release cadence and the test-execution time required to run the corpus against a new product version.

### Registry public access

The registry is public, signed by the working group's release-management role using the project's standard cosign and GPG trust paths. Institutions consume the registry without authentication. A machine-readable feed (JSON Lines) supports institutions integrating registry consumption into vendor-management automation. Registry-document signing follows the same trust-path discipline as binary and spec PDF signing; the institution's trust anchors validate the registry as they validate the binary.

## Code of conduct

The project follows the [Contributor Covenant 2.1](CODE_OF_CONDUCT.md). Violations are reported to the maintainers via the contact in that document.
