# HIPAA Security Rule Mapping — Chain-of-Custody Safeguards

> **What this doc is.** A control-by-control mapping of the HIPAA Security Rule (45 CFR 164.300+) to the chain's design and operational procedures. Written so an OCR auditor, a Covered Entity's privacy officer, or a Business Associate's compliance team can confirm at a glance which Security Rule sub-section is satisfied by which chain primitive or institutional procedure.

> **Why this matters.** HIPAA's Security Rule is the operational counterpart to the Privacy Rule. When PHI flows through the chain (see `hipaa-minimum-necessary.md` for when this happens), the Security Rule prescribes the technical, administrative, and physical safeguards. This document is the line-by-line answer.

---

## 1. The Security Rule structure

The HIPAA Security Rule (45 CFR Subpart C of Part 164) organizes safeguards into three categories:

- **Administrative safeguards** (164.308) — policies, training, sanctions, contingency planning
- **Physical safeguards** (164.310) — facility access, workstation security, device and media controls
- **Technical safeguards** (164.312) — access control, audit controls, integrity, transmission security

Plus organizational requirements (164.314) covering Business Associate Agreements, group health plans, and other relationship structures.

Each standard within a category has either "required" implementation specifications (must be implemented) or "addressable" specifications (must be implemented if reasonable and appropriate; if not, an alternative measure must be implemented and documented). The institution's documentation states which posture applies to each addressable specification.

---

## 2. The full mapping

| Security Rule | Citation | Implementation status | Chain mapping |
|---|---|---|---|
| **Security management process** | 164.308(a)(1) | Required | Risk analysis per `dpia-template.md` (treats HIPAA risk as one input); risk management per the institution's information security program; sanction policy per HR; information system activity review via verifier daily runs and operational events |
| **Assigned security responsibility** | 164.308(a)(2) | Required | Named Security Official; chain-operations team reports to Security Official for chain-relevant responsibilities |
| **Workforce security — authorization and supervision** | 164.308(a)(3)(ii)(A) | Addressable | Only designated chain-operations staff may access the IKM, manage rotations, or access the privacy-store. Access decisions follow the institution's role-based access control policy. Per `token-vault-architecture.md`, privacy-store access is restricted to the privacy-data-protection team — separate from chain operations |
| **Workforce security — clearance procedure** | 164.308(a)(3)(ii)(B) | Addressable | Background checks per HR policy for staff with PHI-access roles |
| **Workforce security — termination procedure** | 164.308(a)(3)(ii)(C) | Addressable | HR coordinates with Security Official to revoke chain-operations and privacy-store access on termination; revocation is logged as an operational event |
| **Information access management — access authorization** | 164.308(a)(4)(ii)(B) | Addressable | Documented access matrix listing who may access what (chain operators, privacy-data-protection team, IR responders, audit) |
| **Information access management — access establishment and modification** | 164.308(a)(4)(ii)(C) | Addressable | Access changes follow change-management; quarterly review of access list |
| **Security awareness and training** | 164.308(a)(5) | Required | Annual workforce training on chain operations, Privacy Rule scope, Minimum Necessary, breach reporting; security-awareness reminders quarterly; password and login monitoring per institution policy; protection against malware per IT operations |
| **Security incident procedures** | 164.308(a)(6) | Required | `incident-response-playbook.md` covers detection, response, reporting, and mitigation. 15 named scenarios cover IKM compromise, privacy-store compromise, integrity failure, and others. Documented response and outcome of each incident |
| **Contingency planning — data backup plan** | 164.308(a)(7)(ii)(A) | Required | Chain ledger and privacy-store backups per spec §10.10 backup-integrity procedures |
| **Contingency planning — disaster recovery plan** | 164.308(a)(7)(ii)(B) | Required | `dr-and-resilience.md` defines RPO, RTO, and recovery sequence |
| **Contingency planning — emergency mode operation plan** | 164.308(a)(7)(ii)(C) | Required | `incident-response-playbook.md` Scenario 10 (recovery from backup) plus emergency-mode operating procedures (degraded service while ledger is being restored) |
| **Contingency planning — testing and revision** | 164.308(a)(7)(ii)(D) | Addressable | Quarterly recovery drill; outcome logged as `master.recovery_drill_completed`; drill report reviewed annually |
| **Contingency planning — applications and data criticality analysis** | 164.308(a)(7)(ii)(E) | Addressable | Chain is classified as critical for AI-decision audit; recovery prioritization documented |
| **Evaluation** | 164.308(a)(8) | Required | Annual evaluation of Security Rule compliance covering chain controls; output incorporated into the institution's broader Security Rule evaluation report |
| **Business Associate Contracts and Other Arrangements** | 164.308(b) | Required | BAA in place for vendor-hosted deployments; BAA covers Security Rule obligations, sub-BA authorization, breach notification per `breach-notification-matrix.md` |
| **Facility access controls — contingency operations** | 164.310(a)(2)(i) | Addressable | Procedures to access facility during emergency for data restoration |
| **Facility access controls — facility security plan** | 164.310(a)(2)(ii) | Addressable | HSM and ledger servers in physically secured data center; FIPS 140-2 Level 3 HSM physical-tamper-detection; locked racks for ledger servers; badge-and-MFA access |
| **Facility access controls — access control and validation** | 164.310(a)(2)(iii) | Addressable | Visitor logs; contractor-access procedures; physical access reviewed quarterly |
| **Facility access controls — maintenance records** | 164.310(a)(2)(iv) | Addressable | Records of physical-security-component repairs, modifications |
| **Workstation use** | 164.310(b) | Required | Policies governing workstation access to chain operations; restricted to dedicated-purpose workstations for IKM management |
| **Workstation security** | 164.310(c) | Required | Physical safeguards for workstations with chain-operations access; secured rooms or full-disk encryption |
| **Device and media controls — disposal** | 164.310(d)(2)(i) | Required | Cryptographic erasure of HSMs at end-of-life; secure disposal of storage media holding chain or privacy-store data per institution media-disposal policy |
| **Device and media controls — media re-use** | 164.310(d)(2)(ii) | Required | Sanitization before media re-use |
| **Device and media controls — accountability** | 164.310(d)(2)(iii) | Addressable | Chain-of-custody records for hardware (HSMs, storage media) holding PHI |
| **Device and media controls — data backup and storage** | 164.310(d)(2)(iv) | Addressable | Backup media stored in physically secured location; encrypted backups |
| **Access control — unique user identification** | 164.312(a)(2)(i) | Required | Each chain-operations user has unique identifier; shared accounts prohibited |
| **Access control — emergency access procedure** | 164.312(a)(2)(ii) | Required | Emergency-access procedure for IR scenarios; emergency access logged |
| **Access control — automatic logoff** | 164.312(a)(2)(iii) | Addressable | Session timeouts on chain-operations consoles per institution policy |
| **Access control — encryption and decryption** | 164.312(a)(2)(iv) | Addressable | Per `article-32-security-mapping.md` §2.2: AES-256-GCM encryption at rest with HSM/KMS-protected keys; chain confidentiality key separate from HMAC IKM and privacy-store key |
| **Audit controls** | 164.312(b) | Required | Verifier daily runs (spec §7); operational events logged per spec §10.2 (`master.reconciliation_completed`, `audit_file.truncation_detected`, `key_fingerprint.mismatch`, etc.); access logging on chain and privacy-store; logs reviewed per audit-procedures schedule |
| **Integrity — mechanism to authenticate ePHI** | 164.312(c)(2) | Addressable | The chain's integrity story IS the authentication mechanism: per-entry HMAC, daily Merkle, HSM-signed seal. Verifier reproduces the integrity check; failure escalates to IR |
| **Person or entity authentication** | 164.312(d) | Required | Multi-factor authentication for chain-operations users; HSM-issued tokens for SDK-to-ledger authentication per spec §4.1.1; mTLS for service-to-service |
| **Transmission security — integrity controls** | 164.312(e)(2)(i) | Addressable | TLS 1.3 ensures transmission integrity; chain HMAC adds independent integrity over the canonicalized payload |
| **Transmission security — encryption** | 164.312(e)(2)(ii) | Addressable | TLS 1.3 between SDK and ledger; mTLS between ledger services |
| **Business Associate Contracts** | 164.314(a) | Required | BAA Schedule A names Security Rule compliance requirements; vendor's SOC 2 Type II report attests to controls |
| **Requirements for Group Health Plans** | 164.314(b) | Required | Applies to health-plan deployments; plan documents amended to incorporate Security Rule obligations |
| **Policies and procedures** | 164.316(a) | Required | All policies and procedures referenced in this document are documented in writing |
| **Documentation** | 164.316(b) | Required | Documentation retained 6 years from creation or last effective date (matches FFIEC 7-year floor; institution applies 7-year retention to align) |

---

## 3. Required vs addressable — the institution's posture

The institution treats every "required" specification as mandatory and every "addressable" specification as implemented to a documented standard. Where an addressable specification's standard implementation is reasonable and appropriate for the institution's environment, the standard implementation is adopted. Where it is not, an alternative measure is documented along with the reasoning.

For chain deployments, no addressable specification is currently implemented at less than its standard expectation. The chain's design — HSM custody, separate privacy-store custody, FIPS-approved primitives, multi-layer integrity — meets or exceeds the addressable bar across the board.

---

## 4. The "scalability" provision

164.306(b) acknowledges that the Security Rule's flexibility allows covered entities and business associates to use security measures appropriate to size, complexity, capabilities, technical infrastructure, security capabilities, and the costs of security measures. The institution's posture:

- Scale: large multinational financial-services entity; high transaction volume; multi-jurisdiction deployment.
- Complexity: chain integrates with multiple AI systems, multiple data sources, multiple regulatory regimes.
- Capabilities: HSM custody, KMS integration, formal security operations, dedicated privacy-data-protection team.
- Cost: chain controls are commodity (FIPS-approved primitives, standard HSM, standard KMS) — cost is proportionate to the regulatory environment.

The institution chose chain controls at the higher end of the reasonable spectrum because the deployment scenario (regulated financial services + healthcare data in some flows + multi-jurisdiction) places the institution at the higher end of the risk spectrum.

---

## 5. Audit evidence package

When OCR or an internal auditor reviews HIPAA Security Rule compliance for the chain, the institution presents:

| Evidence | Source |
|---|---|
| Risk analysis | `dpia-template.md` instance + institution's broader risk analysis |
| Security management policies | Institution's security-management policy library |
| Workforce training records | HR / training system |
| Incident response playbook | `incident-response-playbook.md` |
| BAA(s) | Legal repository |
| SDK configuration files | Configuration management system |
| Verifier output (sample) | Chain operations |
| Operational event logs | SIEM / chain operational-event store |
| Recovery drill reports | Chain operations |
| Annual evaluation report | Security office |
| Access matrix and review records | IAM + privacy office |
| Physical-security records | Facilities + security office |
| HSM / KMS attestations | Vendor SOC 2; HSM FIPS 140-2 certification |

Each evidence item is retained per the institution's records-retention schedule (matched to `retention-justification.md` 7-year floor).

---

## 6. Cross-references

- `hipaa-minimum-necessary.md` — the privacy-side scope question; this document is the security-side controls answer.
- `breach-notification-matrix.md` — HIPAA breach notification timing alongside FFIEC, GDPR, DORA.
- `article-32-security-mapping.md` — GDPR Article 32 mapping; many controls are shared between the two regulations and this document references them.
- `token-vault-architecture.md` — separate custody for the privacy-store / token vault.
- `incident-response-playbook.md` — security incident procedures (164.308(a)(6)).
- `dr-and-resilience.md` — contingency planning (164.308(a)(7)).
- `cloud-hsm-guide.md` — HSM and KMS custody for the chain.
- Spec §4 (chain construction), §7 (verifier), §10 (security overview).

---

## 7. Review cadence

This mapping is reviewed annually by the Security Official, privacy office, and legal counsel. Triggers for early review:

- HHS publishes a Security Rule update or guidance.
- An OCR enforcement action affects a peer institution and reveals a control gap.
- The institution adopts a new technology (e.g., a new HSM model, a new KMS) requiring re-mapping.
- A new BA relationship is established and the BAA needs current Security Rule references.
- An incident reveals that a control is not operating as documented; the mapping is updated to reflect actual practice and the practice is updated to match the documented standard.
