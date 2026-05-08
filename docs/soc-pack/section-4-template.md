# SOC report — Section 4 description template

> **What this doc is.** Starter text for the Service Organization Description (Section 4) of a SOC report covering a chain-of-custody implementation. The institution adapts the language for its own deployment.

## How to use

Section 4 is the description of the system the SOC opinion attests to. The opinion attaches to the description. If the description is wrong, the opinion is wrong. Drafters take care.

The text below is starter language. Fill in `[bracketed]` items with institution-specific facts. Strike sections that don't apply.

---

## Sample Section 4 starter text

### A. Overview of services

`[Institution Name]` operates a chain-of-custody implementation for AI agent decisions made within `[scope: e.g., the Customer Service Routing platform]`. The implementation captures every AI agent decision as a structured event, secures the event with cryptographic primitives that defy retroactive alteration, and produces a daily integrity-bearing artifact (the daily Merkle seal) that an independent verifier can validate.

The implementation conforms to the FFIEC AI Chain-of-Custody Specification, version v1.0 (`spec/chain-of-custody-v1.md`).

### B. Principal service commitments and system requirements

The institution's principal service commitment for the chain is that AI agent decisions captured into the chain are integrity-bearing: any modification of a captured event after capture is detectable by the verifier. The system requirements that support this commitment are:

- HMAC chain construction at the moment of capture
- Daily Merkle aggregation of all events in the tenant-day
- HSM-rooted Ed25519 signature on the daily Merkle root
- An independent verifier that re-computes integrity from a ledger snapshot

### C. Components of the system used to provide the services

#### Infrastructure

- Application hosts running the chain-of-custody SDK in `[deployment region(s)]`
- A ledger server running in `[deployment region(s)]` with PostgreSQL-backed WAL and hot store, S3-backed cold store
- A `[AWS CloudHSM | Azure Managed HSM | Google Cloud HSM | on-prem HSM]` cluster in `[regions]` providing FIPS 140-2 Level `[3 | higher]` protection for cryptographic keys
- An OTel Collector for OTLP routing (where applicable)

#### Software

- `[Vendor name and version | self-developed]` ledger server software, conforming to `chain-of-custody-v1.0`
- `[Vendor name and version | self-developed]` SDK in `[Python | Java | Go | Node]`
- The reference verifier binary, version `[X.Y.Z]`, validated via `verifier-validate.sh` before each use

#### People

- The chain operations team (`[N]` engineers) operates the ledger and SDK deployments
- The HSM administrator role (`[N]` individuals) is held separately from the seal-job operator role (`[N]` individuals)
- Internal audit (`[N]` individuals) operates the verifier on the institution's defined cadence

#### Procedures

- Daily seal job runs at UTC 00:00 + 60 minutes per tenant (default cadence)
- HSM PIN rotation on `[quarterly | semi-annual]` cadence
- IKM rotation on `[institution-defined]` cadence with documented rotation procedure
- IKM retirement on `[institution-defined]` cadence per spec §10.9; retirement requires `chain_entries_referencing_remaining = 0` and is recorded as a `master_key.retired` operational event signed off by `[role]`
- Key-fingerprint reconciliation on `[weekly]` cadence per spec §10.1; recorded as `master.reconciliation_completed` operational event with `fingerprint_unmatched_count` field
- Mid-write-truncation recovery per IR Scenario 9; the SDK-local SQLite buffer and OTLP retention are the documented recovery sources; recovery outcome is recorded on the `audit_file.truncation_detected` operational event
- Verifier validation before each use; reproducible-build verification at least once per verifier release; SLSA provenance attestation validated against `slsa-verifier` at deployment time
- Verifier runs on `[institution-defined]` cadence; substantive SOC-evidence runs are invoked with `--strict`; anomaly-evaluation runs are invoked without
- Regulator-fingerprint reception procedure per `09-threat-model.md` §2.9 + audit-procedures P-28; emits `regulator_fingerprint.rotation_received` / `.rotation_validated` / `.installed` operational events; forged-notice failures route to IR Scenario 11 sub-variant
- Annual cold-DR-key dry-run attestation consumption per `supply-chain.md` §"Institution-side consumption of cold-DR-key dry-run attestation" + audit-procedures P-29; emits `cold_dr.dryrun_attestation_consumed` operational event
- Incident response per `docs/incident-response-playbook.md` (Scenarios 1-12 + 36-hour triage matrix + CIRCIA matrix + federal-regulator routing per charter)

#### Data

- Captured events: `[approximate volume]` per day
- Retention: `[X]` years for events, daily seals, and operational logs
- Storage: hot store is PostgreSQL; cold store is S3 with `[encryption-at-rest with KMS]`

### D. Boundaries of the system

The system covered by this report:

- Includes: the SDK, the ledger server, the daily seal subsystem, the HSM, the tenant key registry, the verifier, the operational events log
- Excludes: the AI agent code itself; the LLM infrastructure (provider-side); the institution's broader observability stack; the institution's IAM, network, and storage encryption controls (those are user-entity responsibilities; see Section 4F)

### E. The five principles applied

The institution claims the following Trust Services Criteria for this engagement:

- `[Common Criteria, Availability, Processing Integrity, Confidentiality, Privacy — pick those that apply]`

The mapping of chain primitives to TSC is documented at `docs/control-map/TSC-mapping.md`.

### F. Complementary user entity controls

The chain's claims depend on user entity controls listed in `docs/control-map/CUECs.md`. The user entity is responsible for operating these controls. The institution attests in this report to operating its share of the controls; the user entity is responsible for the remainder.

`[List the most important CUECs for the user entity reading this report.]`

### G. Significant changes during the period

`[List any material changes to the system during the reporting period, e.g., HSM cluster expansion, master key rotation, vendor change, spec version migration.]`

### H. Subservice organizations

`[List subservice organizations: cloud provider, HSM vendor, observability backend, etc. Identify whether each is in scope (carve-in) or out of scope (carve-out).]`

---

## Drafting notes

- The description must be true. Statements about controls that the institution does not actually operate are an audit failure.
- The description must be complete in scope. If the chain depends on a control, name the control or name the boundary that excludes it.
- The description must be clear. The reader (a regulator, a customer, a counterparty) reads this once. Ambiguity costs the institution credibility.

## Review checklist before issuance

- [ ] Section A names the service in plain English
- [ ] Section B states the principal commitment
- [ ] Sections C and D enumerate every meaningful component and boundary
- [ ] Section E lists the TSC the report covers
- [ ] Section F lists the CUECs without ambiguity
- [ ] Section G captures every change, even minor configuration changes
- [ ] Section H carves in/out subservice organizations correctly
- [ ] All `[bracketed placeholders]` are filled in or stricken
- [ ] All institution-specific facts are accurate as of the reporting period end date
