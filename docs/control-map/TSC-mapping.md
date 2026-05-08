# SOC 2 Trust Services Criteria mapping

> **What this doc is.** Mapping of chain-of-custody primitives to AICPA Trust Services Criteria 2017 (with 2022 points of focus updates). SOC 2 practitioners use this to translate the design docs into TSC language.

## Common Criteria (CC)

The Common Criteria apply to every SOC 2 engagement.

### CC1 — Control environment

| TSC | Mapping |
|---|---|
| CC1.1 (Integrity and ethical values) | Reflected in the institution's broader control environment; chain inherits |
| CC1.4 (Commitment to attract, develop, retain individuals) | Examiner training program; HSM operator separation of duties |

### CC2 — Communication and information

| TSC | Mapping |
|---|---|
| CC2.1 (Internally generated information) | Operational events (`06-ledger-server-design.md` §7.3.1) |
| CC2.2 (Information communicated externally) | Verifier output (PDF + JSON + bundle); SOC report |

### CC3 — Risk assessment

| TSC | Mapping |
|---|---|
| CC3.1 (Specifies suitable objectives) | Spec scope statement (`spec/chain-of-custody-v1.md` §1) |
| CC3.2 (Identifies risks) | Threat model (`09-threat-model.md`) |
| CC3.4 (Assesses changes) | Spec version change-control; algorithm rotation provisions |

### CC5 — Control activities

| TSC | Mapping |
|---|---|
| CC5.1 (Selects and develops control activities) | The four primitives plus surrounding controls |
| CC5.2 (Selects and develops technology controls) | HSM, append-only storage, deterministic verifier |
| CC5.3 (Deploys through policies and procedures) | Examiner approval workflow; PIN rotation; reconciliation cadence |

### CC6 — Logical and physical access controls

| TSC | Mapping |
|---|---|
| CC6.1 (Logical access controls — administer) | RBAC defense-in-depth (`06-ledger-server-design.md` §3.5); HSM operator separation; **`master_key.retired` operational event recording IKM-roster lifecycle changes (spec §10.9 retention rule)** |
| CC6.2 (Authentication of internal/external users) | Session-key handshake (`02-chain-construction.md` §4.1); SPIFFE/SPIRE recommendation |
| CC6.3 (Authorization) | HSM `sign`-only operator role; tenant-isolated keys |
| CC6.6 (Physical access — HSM) | FIPS 140-2 L3 hardware tamper resistance |
| CC6.7 (Restricts movement of data) | BYOC IAM matrix; vendor support telemetry routing; per-tenant `key_fingerprint` identity binding (spec §4.1) |
| CC6.8 (Prevents/detects unauthorized software) | Cosign + GPG + reproducible-build verification; **verifier `--strict` refusal of `dev_mode=true` seals and `kms_handle_uri = "plaintext-*"` chain entries (spec §10.7 + §7 step 12)** |

### CC7 — System operations

| TSC | Mapping |
|---|---|
| **CC7.1 (Detect new vulnerabilities)** | SBOM publication (CycloneDX); vulnerability scan reports per release |
| **CC7.2 (Monitor system components)** | **Verifier anomaly reporting; operational events list; metrics. Includes the new `audit_file.truncation_detected` operational event (verifier-emitted on §4.1 mid-write refusal) and `master.reconciliation_completed` (weekly fingerprint reconciliation per spec §10.1).** |
| CC7.3 (Identify, evaluate security events) | Chain-detected event flow into SIEM; integration with institution's incident response |
| CC7.4 (Respond to security incidents) | IR playbook for chain-detected events (Scenarios 1-10 + 36-hour triage matrix) |
| CC7.5 (Recover from security incidents) | DR/RPO/RTO (`06-ledger-server-design.md` §7.5); IR Scenario 9 (audit_file truncation recovery); IR Scenario 10 (backup-integrity verifier failure) |

### CC8 — Change management (selected mappings)

| TSC | Mapping |
|---|---|
| CC8.1 (Change management procedures) | Spec version change-control; `format_version` per chain entry (spec §4.4) and per file header bound under signed `sign_payload` (§4.3); **`master_key.retired` and `master_key.rotated` operational events documenting IKM-lifecycle change events**; institution's IKM-retirement procedure documented per spec §10.9; **`regulator_fingerprint.installed` operational event documenting trust-anchor-cache change events**; **`cold_dr.dryrun_attestation_consumed` operational event documenting institution's annual cold-DR-key dry-run consumption per `supply-chain.md`** |

| TSC | Mapping |
|---|---|
| **CC8.1 (Authorize, design, develop, configure)** | Spec version change-control; deterministic builds; conformance corpus |

### CC9 — Risk mitigation

| TSC | Mapping |
|---|---|
| CC9.1 (Identifies, selects, develops risk mitigation) | Threat model + accepted residual risks |
| CC9.2 (Vendor management) | BYOC and vendor-hosted topologies; CUECs around vendor IAM |

## Additional Criteria

The institution selects which of the additional criteria apply based on its claims.

### Availability (A)

| TSC | Mapping |
|---|---|
| A1.1 (Maintain capacity) | Streaming Merkle for scale; hot-path budget |
| A1.2 (Environmental protections) | Inherited from cloud provider; HSM cluster failover |
| A1.3 (Backup and disaster recovery) | DR/RPO/RTO targets; multi-region workaround |

### Confidentiality (C)

| TSC | Mapping |
|---|---|
| C1.1 (Identifies confidential information) | The institution identifies in its data classification |
| C1.2 (Disposes of confidential information) | Retention policy; institution's standard data-disposal |

The chain itself does not provide confidentiality. The institution operates encryption-at-rest on the storage layer if confidentiality is in scope. The chain composes — its outputs (`payload_hash`, `prev_hash`, etc.) contain no key material and are safe to route to any backend.

### Processing Integrity (PI)

| TSC | Mapping |
|---|---|
| **PI1.1 (Inputs are complete and accurate)** | **Chain integrity property** |
| **PI1.2 (Inputs are processed completely)** | **Chain catches dropped events as a verification failure** |
| PI1.3 (Outputs are complete) | Verifier reports per-day; missing seals are detectable |
| PI1.4 (Inputs/processing are authorized) | Session-key handshake authenticates capture |
| PI1.5 (Outputs are distributed to authorized parties) | Verifier produces PDF/JSON; institution distributes per its policy |

The chain is fundamentally a **Processing Integrity** control. PI1.1 and PI1.2 are the headline.

### Privacy (P)

For institutions claiming the Privacy criterion, criterion-level breakdown of how the chain composes:

| Privacy Criterion | Chain composition |
|---|---|
| **P1 (Notice and communication)** | Institution-side; chain composes by being part of the institution's notice about AI use |
| **P2 (Choice and consent)** | Institution-side; chain neutral |
| **P3 (Collection)** | The chain captures tokens of personal data, not the data itself; tokenization happens at the SDK before canonicalization (see `docs/privacy-by-design.md`) |
| **P4 (Use, retention, disposal)** | The institution's privacy-store retention vs the chain's retention are separate; chain retention is integrity-bearing-record retention |
| **P5 (Access)** | Privacy-store responds to access requests; chain's integrity property protects the records' authenticity |
| **P6 (Disclosure to third parties)** | Privacy-store responds; the chain's integrity protects against unauthorized disclosure of integrity-bearing records |
| **P7 (Security for privacy)** | Standard security applies to privacy-store; the chain provides the integrity layer for records describing the privacy-store activity |
| **P8 (Quality)** | The chain provides integrity-bearing records of decisions affecting the privacy-store |

The chain is not a Privacy control by itself; it composes with the institution's Privacy program through the privacy-by-design pattern (`docs/privacy-by-design.md`).

## How to use

When drafting a SOC 2 control description:

1. Pick the criteria the institution claims (always Common Criteria; Availability and Processing Integrity are typical for chain-of-custody implementations)
2. For each criterion, write the institution's control description using the language above as a starting point
3. Identify the testing approach (how the SOC team will gather evidence)
4. Identify the CUECs (`CUECs.md`)

## Headline mapping (one-paragraph summary)

The chain is a **PI1.1 / PI1.2 / CC7.2** control with strong **CC8.1** properties (deterministic builds, conformance corpus). Surrounding controls (CC6.1–6.8, CC7.1–7.5) are the institution's responsibility, with the chain providing supporting evidence. The chain does not satisfy Confidentiality or Privacy criteria by itself; those are the institution's separate controls.
