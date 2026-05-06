# Security policy

## Reporting a vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Email **security@tesseraseal.com** with:

- A description of the vulnerability and its impact
- Steps to reproduce, with a minimal test case if possible
- The affected version and platform
- Your name and contact information (optional &mdash; anonymous reports accepted)

We will acknowledge receipt within **2 business days** and provide a substantive response within **10 business days**.

## What we treat as a security issue

This is a regulator-facing chain-of-custody reference implementation. We treat the following as security issues, in addition to the obvious:

- **Any defect in the chain construction** that would allow an event to be inserted, removed, or altered without producing a verifiable break in the chain.
- **Any defect in the Merkle seal** that would allow the daily root to be reproduced for a tampered set of leaves.
- **Any defect in the HSM signing path** that would allow a forged daily root to be accepted by the verifier.
- **Any defect in the verifier** that would allow an invalid chain to pass verification, or a valid chain to fail.
- **Any deviation from the specification** that produces incompatible output between two conforming implementations.
- **Any cryptographic primitive choice** that is not FIPS-approved or NIST-recommended.
- **Any timing, fault-injection, or side-channel attack** that recovers session keys, master keys, or HSM-held keys.
- **Any memory-safety issue** in code paths that handle untrusted input (the verifier especially &mdash; it consumes adversarial ledgers from compromised institutions).

## What we do not treat as a security issue

- Performance regressions (file a regular issue)
- API ergonomics (file a regular issue)
- Documentation errors (file a regular issue, unless the error misrepresents a security property)
- Issues in third-party software unless the vulnerability arises from how we use it

## Coordinated disclosure

We follow industry-standard coordinated disclosure:

1. **Day 0** &mdash; report received, acknowledged
2. **Day 0&ndash;14** &mdash; triage, severity assessment, internal fix
3. **Day 14&ndash;90** &mdash; private fix release, embargo period
4. **Day 90+** &mdash; public disclosure with CVE assignment, advisory, and credit to reporter

We may shorten the embargo if the vulnerability is being exploited in the wild, or if the reporter publishes details independently.

## Bug bounty

This project does not currently offer a bug bounty. We acknowledge reporters in the security advisory and the release notes.

## PGP key

Our PGP key for encrypted reports is available at:

```
https://tesseraseal.com/.well-known/security.txt
```

(Will be populated when public release ships.)

## Hall of credit

Reporters who choose to be credited will be listed here once the first public release ships.
