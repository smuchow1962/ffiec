# FFIEC IT Examination Handbook mapping

> **What this doc is.** Mapping of chain-of-custody primitives to specific control objectives in the FFIEC IT Examination Handbook. IT examiners working from the Handbook use this to confirm which control objectives the chain satisfies.

## Information Security booklet (IS, September 2016)

| IS Section | Topic | Chain coverage |
|---|---|---|
| II.B (Risk Management) | Information security risk identification | Threat model `09-threat-model.md` enumerates adversaries and accepted residual risks |
| II.C.5 (Logical access) | Authentication and authorization | Session-key handshake security floor (`02-chain-construction.md` §4.1.1) — workload identity attestation, mTLS, HSM-issued tokens; spec §4.1.1 properties 1-4 |
| **II.C.10 (Logging)** | **Activity logging and integrity** | **The chain itself.** HMAC chain + Merkle seal + HSM signature provide logging integrity that defies retroactive alteration. **Plus weekly key-fingerprint reconciliation per spec §10.1**: institutions match every observed `(tenant_id, key_version, key_fingerprint)` triple against the IKM-roster's expected fingerprint, surfacing botched rotation / cross-tenant configuration drift / restored-from-wrong-backup as a high-priority alert. The `master.reconciliation_completed` operational event is the audit-evidence artifact (audit-procedures P-6). |
| II.C.11 (Incident response) | Detection and response | Verifier anomaly reporting (`07-verifier-design.md` §5.1); IR playbook including dedicated scenarios for `key_fingerprint mismatch` (Scenario 7), `unknown_key_version` (Scenario 8), and `audit_file.truncation_detected` (Scenario 9). |
| **II.C.13 (Cryptographic controls)** | **Encryption and key management** | Ed25519 in HSM custody (`04-hsm-custody.md`); HMAC-SHA-256 chain construction. **Plus three rework primitives**: (a) per-entry `key_fingerprint = SHA-256(utf8(tenant_id) \|\| ikm)[:16]` is a public identity binding the verifier checks before any MAC compute (spec §7 step 8), going beyond key custody to integrity-of-key-derivation; (b) IKM minimum 32 bytes per spec §10.6 (RFC 4868), closing offline-grinding attacks against the public fingerprint on low-entropy IKMs; (c) software-key adapter compile-time exclusion per spec §10.7, structurally preventing dev-mode key material from shipping to production. |
| II.C.20 (Software integrity) | Software assurance | Reproducible builds + cosign + GPG manifest (`07-verifier-design.md` §8.4); supply chain doc covers SLSA-equivalent properties. |
| III.B (Audit) | Audit program | Verifier provides independent audit artifact; conformance corpus at `spec/test-vectors/` with byte-level fixture (`chain_vectors.json`) including the rotation-mid-day case (`010-tenant-ikm-rotation-mid-day/`). |

The headline mapping is **II.C.10 (Logging)**. The IS booklet's logging section requires that activity logging be sufficient to support incident detection and forensic investigation. The chain extends this requirement: the logs are not just present, they are integrity-bearing — verifiable by an independent examiner without trusting the institution's logging infrastructure.

The v1.0-rework deepens the **II.C.13 (Cryptographic controls)** mapping by introducing per-entry key identity provenance (`key_fingerprint`). Examiners working II.C.13 should confirm: (1) the institution's IKM length is at least 32 bytes (spec §10.6); (2) the institution's chain entries carry valid `key_fingerprint` bytes that the verifier confirms at step 8; (3) the institution's production build does NOT include the software-key adapter (no chain entries with `kms_handle_uri = "plaintext-dev"` in production); (4) the institution operates weekly key-fingerprint reconciliation per spec §10.1. These are cryptographic-controls evidence beyond "HMAC + Ed25519 in HSM custody."

## Architecture, Infrastructure, and Operations booklet (AIO, June 2021)

| AIO Section | Topic | Chain coverage |
|---|---|---|
| II.A.1 (Architecture documentation) | System architecture documentation | Design docs `docs/design/` |
| II.B (Resilience) | Availability and continuity | DR/RPO/RTO (`06-ledger-server-design.md` §7.5) |
| II.C (Capacity management) | Capacity planning | Hot-path budget (`02-chain-construction.md` §5); streaming Merkle for scale |
| **II.E (Change management)** | **Configuration management and change control** | Spec version per event (`ffiec.chain.spec`); chain-stamp `format_version` per entry AND per file header AND per signed seal payload. **The `format_version` field is the load-bearing change-management primitive**: when the institution upgrades SDKs across a `format_version` boundary (v1 → v2 in a hypothetical future), the verifier dispatches on the field and refuses unrecognised versions most-specific-first (spec §7 step 1). Examiners should confirm the institution operates documented change management for any `format_version` change, including alignment between SDK-stamped values, ledger-stored values, and verifier-version compatibility. **Algorithm change management (dual-algorithm transitional period)**: when the institution adopts a post-quantum algorithm alongside Ed25519 (a v1.x feature), the institution operates change management on (a) the `algorithm` field stamped on each seal record, (b) the `signatures` list (when co-signing), (c) the institution's declared algorithm-posture configuration the verifier consults to dispatch on `signatures` per spec §7 step 11 cases (a)-(e). Examiners confirm: the institution's change-management process covers algorithm-posture transitions; the institution's IR program has documented response procedures for spec §7 step 11 dual-algorithm failure modes; the institution coordinates with the regulator on multi-year migration timelines per `09-threat-model.md` §2.8 dual-algorithm transitional period guidance. |
| III.A (Operations) | Operational monitoring | Operational events (`06-ledger-server-design.md` §7.3.1); metrics (§7.2). The §10.2 operational-event catalog now includes `audit_file.truncation_detected` and `master_key.retired` for full coverage of chain-detected and IKM-lifecycle events. |
| IV (Outsourcing dependency) | Vendor management | BYOC and vendor-hosted topologies (`00-overview.md` §5); supply chain (`supply-chain.md`); IAM permission matrix (`byoc-deployment.md`) |

## Audit booklet (April 2012)

| Audit Section | Topic | Chain coverage |
|---|---|---|
| II (Audit program) | Independent assessment | Verifier is operable by internal audit, external audit, or examiner |
| III (Audit testing) | Effective challenge | Conformance corpus + verifier provides testable basis |
| IV (Audit reporting) | Audit findings | Verifier produces standardized PDF + JSON |

## Outsourcing Technology Services booklet (June 2021)

| OTS Section | Topic | Chain coverage |
|---|---|---|
| II (Risk assessment) | Vendor risk | Topology-neutral spec; vendor inherits standard third-party-risk-management |
| III (Selection) | Due diligence | Conformance corpus enables vendor evaluation against an objective standard |
| IV.B (Contracting) | Right to audit | Verifier is the audit artifact; institution-held public key enables independent verification |
| V (Ongoing monitoring) | Vendor performance | Verifier output over time |

## Information Security – AIO crosswalk

The chain primarily serves the IS booklet's logging integrity requirement (II.C.10) and the AIO booklet's monitoring and resilience requirements. It supports the Audit booklet by providing an artifact independent audit can use, and the OTS booklet by allowing examination of vendor-hosted implementations through the same verifier.

## How to use

In an IT examination:

1. The examiner identifies which Handbook sections are in scope
2. For II.C.10, the chain is the headline evidence — verifier output is the artifact
3. For supporting sections, the institution's broader controls are the headline; the chain provides one of several pieces of evidence
4. The examiner combines chain output with other examination evidence to form the overall finding

The chain does not replace the Handbook; it satisfies a specific portion of it.
