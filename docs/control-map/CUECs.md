# Complementary User Entity Controls (CUECs)

> **What this doc is.** The controls the user entity (the institution using the chain) must operate for the chain's claims to hold. SOC 1 reports list these as CUECs; SOC 2 reports list them as institution responsibilities. Without these controls, the chain's integrity claims are weaker than the spec implies.

## Why CUECs matter

The chain is a substrate. Its integrity claims depend on the institution operating a set of controls around the substrate. If the institution does not operate these controls, the chain still produces output, but the output's defensibility against a determined adversary is reduced.

Each control below is the institution's responsibility. The institution documents how it operates each control; SOC and examination teams test that the documentation matches reality.

## Two axes: normative level vs operational tier

Each CUEC is annotated along two independent axes.

The **normative level** column maps each CUEC to the spec's RFC 2119 keyword that backs it: **MUST** when the spec text the CUEC operationalizes uses MUST or REQUIRED; **SHOULD** when the spec uses SHOULD or RECOMMENDED; **MAY** when the spec uses MAY; **N/A** when the CUEC operationalizes a non-normative best practice (operational hygiene, vendor-management discipline, change-management discipline) that is not gated by spec text. The normative level is a property of the SPEC, not of the institution's deployment maturity.

The **tiering** section that follows (Tier 1 / Tier 2 / Tier 3) is a different axis: production-readiness staging. Tier 1 CUECs are the ones an institution operates from day 1; Tier 2 in the first 6 months; Tier 3 within the first year. Tiering reflects ramp-up sequencing, not normative weight.

The two axes are independent. Tier 1 includes both MUST-level and SHOULD-level CUECs (a SHOULD-level control can be load-bearing enough on day 1 that the institution operates it from initial production). A SHOULD-level CUEC can also be Tier 3 (low operational priority on the institution's ramp-up path but still a normative recommendation in the spec).

**How institutions consume the two axes.**

- An institution running a v1.0-conformant deployment operates ALL **MUST**-level CUECs without exception. Failing to operate a MUST-level CUEC is a conformance break; the institution's CC8.1 control description names the specific MUST-level CUECs it operates.
- **SHOULD**-level CUECs are operated unless the institution documents a compensating control. The institution's control description names the SHOULD-level CUECs it operates AND any SHOULD-level CUECs it does NOT operate (with the compensating-control rationale and the SOC/examiner-acceptance disposition). A SHOULD-level CUEC operated without exception is the default conformant posture; a SHOULD-level CUEC declined without compensating control is non-conformant.
- **MAY**-level CUECs are operated by institutions whose context (risk profile, regulatory posture, deployment topology) indicates the relevant value. The institution's control description lists which MAY-level CUECs it elects to operate.
- **N/A**-level CUECs are operated by institutions with the relevant operational context (e.g., vendor management for institutions consuming a vendor-hosted topology). The N/A label means the spec is silent on the discipline, NOT that the discipline is optional — it means the discipline is institution-defined operational hygiene.

The normative-level annotation is the answer to the examiner's question "which of these is mandatory for v1.0 conformance?" The tiering annotation is the answer to the operations team's question "which of these does the institution operate first?"

## Identity and access management

| CUEC | Normative level | Description |
|---|---|---|
| **CUEC-IAM-01** | SHOULD | The database role used by the ledger writer is granted INSERT and SELECT only on the events and daily_seals tables; UPDATE, DELETE, and TRUNCATE are revoked. (Spec §10.3 names the database-role SHOULD; CUEC-IAM-01 is the institution-side control description naming how the SHOULD is operated.) |
| **CUEC-IAM-02** | SHOULD | DBAs operate under a role with the same restriction; schema migrations run under a privileged role exercised through change management. (Tracks §10.3 SHOULD; institution operationalization.) |
| **CUEC-IAM-03** | MUST | The HSM seal-job operator role grants `sign` only; `extract`, `delete`, and `import` are revoked or require dual control. (Spec §10.5 names the seal-job-operator MUST.) |
| **CUEC-IAM-04** | MUST | The HSM administrator role is held by a different individual than the seal-job operator (separation of duties). At small institutions where this is impractical, dual-control is the documented compensating control. (Spec §10.5 names separation-of-duties REQUIRED where operationally feasible; the dual-control compensating-control posture is named in the spec text.) |
| **CUEC-IAM-05** | SHOULD | Application processes that capture chain events authenticate to the master-key custodian using a credential bound to the host's identity (mTLS, SPIFFE, or HSM-issued token). Anonymous handshakes are not accepted. (Spec §4.1.1 establishes the session-key handshake security floor; CUEC-IAM-05 is the institution-side discipline.) |

## Cryptographic controls

| CUEC | Normative level | Description |
|---|---|---|
| **CUEC-CRY-01** | N/A | The HSM PKCS#11 PIN is sourced from the institution's secret-management system at process start. The PIN is not embedded in configuration files. (Operational hygiene; the spec is silent on PIN management. Institution-defined best practice.) |
| **CUEC-CRY-02** | N/A | The HSM PIN is rotated on a documented cadence (typically quarterly). Each rotation is logged in the institution's control-evidence repository. (Operational hygiene; institution-defined.) |
| **CUEC-CRY-03** | MUST | The institution's tenant master HMAC key is held in tenant-controlled storage at FIPS 140-2 L3 protection. The master never reaches application hosts. (Spec §10.5 / §3 IKM definition / §10.6 — IKM custody is MUST.) |
| **CUEC-CRY-04** | SHOULD | The institution operates **key-fingerprint reconciliation** at a documented cadence (no more than weekly per spec §10.1). The reconciliation matches every `(tenant_id, key_version, key_fingerprint)` triple observed on captured events against the institution's IKM-roster's expected fingerprint per `(tenant_id, key_version)` pair; mismatches are high-priority alerts (botched rotation, restored backup pointed at wrong tenant, cross-tenant key swap). The corresponding operational event is `master.reconciliation_completed` (audit-procedures P-6). (Spec §10.1 names the SHOULD.) |
| **CUEC-CRY-05** | SHOULD | Key rotation procedures are documented and tested at least annually. The institution retains historical public keys with their validity windows in the tenant key registry. (Spec §10.9 SHOULD on retention coupling; rotation-procedure documentation is institution-defined operational discipline backing the spec's normative posture.) |
| **CUEC-CRY-06** | SHOULD | The institution operates a documented regulator-fingerprint reception procedure per `09-threat-model.md` §2.9 Adversary I. The procedure emits three operational events (`regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, `regulator_fingerprint.installed`) with the field schemas in `docs/soc-pack/control-evidence-events.md`. Forged-rotation-notice failures route to IR Scenario 11 sub-variant. (Threat-model §2.9 names the reception procedure as the institution's defense; the SOC-pack schema makes it mechanically testable.) |
| **CUEC-IR-05** | SHOULD | The institution consumes the project's annual `KEY-DR-DRYRUN-{year}.asc` cold-DR-key dry-run attestation within 30 days of publication. The consumption is recorded as a `cold_dr.dryrun_attestation_consumed` operational event; the archived attestation is retained per the institution's chain-event retention. Missed dry-run window triggers IR Scenario 11. (supply-chain.md "Institution-side consumption" section names the discipline; IR Scenario 11 names the failure mode.) |

## Operational controls

| CUEC | Normative level | Description |
|---|---|---|
| **CUEC-OPS-01** | SHOULD | Time synchronization (NTP or equivalent) is operational on all hosts capturing chain events and on the ledger server. (Spec §10.4 names the SHOULD.) |
| **CUEC-OPS-02** | MUST | Operational events emitted by the ledger (`06-ledger-server-design.md` §7.3.1) are retained at least as long as the chain events themselves. (Spec §10.2 names the retention MUST.) |
| **CUEC-OPS-03** | N/A | Backups of the WAL are continuous (streaming replication or equivalent). Backup integrity is tested on the institution's standard cadence. (Operational hygiene; institution-defined.) |
| **CUEC-OPS-04** | N/A | The institution operates monitoring on the seal-age metric (`ledger_seal_age_seconds`) and alerts when the value exceeds the institution's defined threshold. (Spec §4.3.1 names a 72-hour notification SHOULD; CUEC-OPS-04 is the institution-defined monitoring discipline that produces the evidence the SHOULD relies on.) |
| **CUEC-OPS-05** | SHOULD | The institution notifies its primary regulator if the daily seal is delayed beyond 72 hours. (Spec §4.3.1 names the SHOULD.) |

## Verifier and audit

| CUEC | Normative level | Description |
|---|---|---|
| **CUEC-VER-01** | N/A | Internal audit (or an equivalent independent function) runs the verifier on a documented cadence (at minimum, before each examination). (Operational discipline; the spec normates the verifier procedure but is silent on internal-audit cadence.) |
| **CUEC-VER-02** | N/A | The verifier binary is validated with `verifier-validate.sh` (cosign + GPG manifest) before each use. (Operational discipline; supply-chain.md names the validation; the spec is silent on per-use cadence.) |
| **CUEC-VER-03** | N/A | Independent reproducible-build verification is performed at least once per verifier release. The institution retains the rebuild log as control evidence. (Operational discipline; supply-chain.md describes the procedure; the spec is silent on cadence.) |
| **CUEC-VER-04** | N/A | Verifier output (PDF, JSON, bundle) is retained alongside the examination working papers for the institution's standard retention period. (Examination-readiness discipline; institution-defined retention.) |
| **CUEC-VER-05** | SHOULD | The institution validates the SLSA provenance attestation (`verifier.intoto.jsonl`) against `slsa-verifier` (or equivalent SLSA-aware tooling) at deployment time as a deployment-gating control. The validation MUST complete with a PASS disposition before the binary reaches production; a binary whose provenance attestation does not validate MUST NOT be deployed. The `slsa-verifier` output log is archived for the binary's deployment lifetime (typically 7 years, aligned with chain-event retention). The institution's CC8.1 procedure tests that the deployment gate is operative — a sample of deployment events shows the SLSA-verifier output captured, archived, and pre-condition to deployment success. The institution's trust-anchor table includes the `slsa-verifier` public key and binary SHA-256 cached out-of-band per `supply-chain.md` "Trust anchors". (supply-chain.md "Institutional consumption of the SLSA L3 attestation (MUST-tier)" section; the spec's normative-level posture on SLSA L3 is SHOULD because the spec text positions provenance attestation as a recommended supply-chain control rather than a chain-of-custody primitive — institutions adopting v1.0 operate the discipline unless they document a compensating supply-chain control.) |

## Incident response

| CUEC | Normative level | Description |
|---|---|---|
| **CUEC-IR-01** | N/A | The institution operates an incident-response playbook for chain-detected events (`docs/incident-response-playbook.md`). (Institution-defined IR discipline.) |
| **CUEC-IR-02** | N/A | When the verifier reports a chain-integrity failure, the institution executes the playbook and documents the response. (Institution-defined IR discipline.) |
| **CUEC-IR-03** | N/A | When a master-key compromise is suspected, the institution rotates the key, notifies the regulator, and preserves forensic evidence per the playbook. (Institution-defined IR discipline; the regulator-notification rule is the FFIEC computer-security incident notification framework, NOT spec text.) |
| **CUEC-IR-04** | N/A | When the verifier cannot produce a passing report due to backup corruption or recovery from older backup, the institution treats this as an integrity-control failure and follows the cyber-incident notification framework. (Institution-defined IR discipline.) |

## Vendor and supply chain (BYOC and vendor-hosted topologies)

| CUEC | Normative level | Description |
|---|---|---|
| **CUEC-VND-01** | SHOULD | The vendor's container image is pulled through a bank-controlled mirror; cosign signatures are verified before deployment. The institution adopts the signature-preserving mirror pattern as the v1.0 recommended posture (per `supply-chain.md` "Project recommendation"); institutions on the re-signing pattern operate the documented bridging controls (WORM 7-year audit log, `mirror.reconciliation_completed` continuous monitoring, bridge-invariant documentation, IR Scenario 13). (supply-chain.md project-level recommendation.) |
| **CUEC-VND-02** | N/A | The vendor's IAM access to bank cloud accounts is the minimum required; the institution can revoke unilaterally. (Institution-defined vendor-management discipline.) |
| **CUEC-VND-03** | N/A | Vendor support telemetry (when applicable) flows through a bank-controlled router that applies a documented redaction policy. (Institution-defined vendor-management discipline.) |
| **CUEC-VND-04** | N/A | Vendor SOC reports are obtained on the institution's standard third-party-risk-management cadence. (Institution-defined third-party-risk-management discipline.) |
| **CUEC-VND-06** | SHOULD | The institution validates the vendor's FFIEC conformance attestation as part of vendor selection and ongoing review. The attestation is separate from, and complementary to, the vendor's SOC report — the SOC report covers vendor-side operational controls, while the conformance attestation establishes that the vendor's implementation passes the FFIEC conformance corpus. The institution: (a) consults the project's public vendor-conformance registry at vendor selection and on its ongoing-review cadence (typically annual, with ad-hoc re-validation on corpus-version change, vendor product-version change affecting chain behavior, or registry revocation); (b) fetches the vendor's published attestation public key from the documented vendor-controlled URL and caches it out-of-band; (c) validates the attestation signature against the cached public key; (d) confirms the corpus version named in the attestation matches the corpus version the institution itself runs against (or commits to run against in its deployment timeframe); (e) reviews the attestation's exception list against the institution's deployment posture. The institution retains as evidence: the signed attestation document, the signature artifact, the vendor's public key, the validation log (validating individual, date, validation outcome, registry-current-state cross-check), and the registry-consultation record (date, vendor entries reviewed, status of those entries on the consultation date). Retention is the longer of (a) the institution's chain-event retention period or (b) seven years from validation. The institution's CC8.1 control description names the attestation evidence consumed in vendor selection, the cadence of re-validation (annually plus ad-hoc trigger conditions), the validation procedure operated, and the failure-scenario disposition (vendor refuses to attest → control gap escalation per the institution's vendor-management procedure; attestation revocation → vendor-management re-evaluation per the institution's standard policy). The institution's SOC team confirms the attestation-validation evidence as part of the vendor-management control test: sample vendor relationships in scope, confirm the evidence repository contains the attestation, signature, public key, and validation log, confirm validation occurred within the documented cadence, confirm signature validation, corpus-version match, and exception-list review are recorded, confirm registry revocation or supersession events during the period triggered the institution's documented response. (`docs/vendor-conformance-attestation.md` documents the full procedure; `GOVERNANCE.md` §"Vendor-conformance attestation registry" names the project-side governance commitments — sub-committee operation, quarterly registry review, 14-day revocation publication, 90-day corpus-version re-attestation grace period.) |

## Configuration and change management

| CUEC | Normative level | Description |
|---|---|---|
| **CUEC-CFG-01** | N/A | Changes to the chain configuration (cadence, retention, HSM endpoints, public keys) flow through the institution's change management process with appropriate approvals. (Institution-defined change-management discipline.) |
| **CUEC-CFG-02** | N/A | Spec version changes are documented in the institution's change log with the spec version, the rationale, and the migration plan. (Institution-defined change-management discipline.) |
| **CUEC-CFG-03** | N/A | Cadence relaxation requests follow `regulator-pack/examiner-approval-template.md`. (Institution-defined operational-discipline; the spec §4.2.1 admits cadence configurability without normating the relaxation procedure.) |

## How to use

In a SOC 1 report, list the CUECs in Section 4 (or the equivalent section). The institution's control description claims these CUECs are operating; the SOC team tests them; the SOC opinion attests to their effectiveness.

In a SOC 2 report, the same content lives in the institution's part of the description. The SOC 2 framework treats them as user-entity responsibilities since the institution is also the entity providing the service to its own customers.

In an FFIEC examination, the IT examiner reviews the same controls as standard examination evidence.

## Tiering

For institutions ramping up to full operation, the CUECs are tiered:

### Tier 1 — Critical (operate from day 1)

Without these, the chain's claims are weaker than the spec implies. Institutions adopting the chain must operate Tier 1 from initial production.

- **CUEC-IAM-01..05** — RBAC; HSM separation of duties (or dual-control at small institutions)
- **CUEC-CRY-01..05** — HSM PIN management; master key custody; key rotation
- **CUEC-OPS-01** — NTP synchronization
- **CUEC-OPS-04** — Seal-age monitoring
- **CUEC-OPS-05** — Regulator notification on 72-hour seal delay
- **CUEC-VER-01** — Run the verifier at minimum quarterly
- **CUEC-IR-01..04** — Incident response playbook

### Tier 2 — Important (operate within first 6 months)

These strengthen the chain's claims and align with full SOC and examination coverage.

- **CUEC-VER-02..04** — Verifier validation; reproducible-build verification
- **CUEC-VER-05** — SLSA L3 institution-side validation (deployment-gating per `supply-chain.md`); operate before the institution's first production binary deployment so the deployment gate is exercised on the very first release
- **CUEC-VND-01..04** — Vendor management (if vendor implementation in use)
- **CUEC-VND-06** — Vendor-conformance attestation validation (if vendor implementation in use); operate from initial vendor selection so the attestation evidence is in the institution's vendor-management evidence base before first deployment
- **CUEC-OPS-02..03** — Operational event retention; backup integrity testing

### Tier 3 — Operational hygiene (operate within first year)

These are good practice and align with mature operations.

- **CUEC-CFG-01..03** — Configuration and change management

The institution declares its tier of operation in its control description; SOC and examination teams test against the declared scope.

## Mapping to chain primitives

| Chain primitive | Most relevant CUECs |
|---|---|
| HMAC chain at capture | IAM-05, CRY-03, CRY-04 |
| Daily Merkle seal | OPS-01, OPS-04, OPS-05 |
| HSM-rooted root signature | IAM-03, IAM-04, CRY-01, CRY-02, CRY-05 |
| OpenTelemetry-native wire | (Configuration; no specific CUEC) |
| Verifier independent check | VER-01..05, IR-01..04 |
| Supply-chain trust path (binary, container, provenance) | VER-02, VER-03, VER-05, VND-01, IR-05 |
