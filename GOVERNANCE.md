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

## Code of conduct

The project follows the [Contributor Covenant 2.1](CODE_OF_CONDUCT.md). Violations are reported to the maintainers via the contact in that document.
