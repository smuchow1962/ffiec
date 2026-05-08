# GDPR Article 32 Security Mapping — Chain-of-Custody Controls

> **What this doc is.** A control-by-control mapping of GDPR Article 32 (security of processing) requirements to the chain-of-custody design. Written so a Data Protection Authority reviewer or a Chief Privacy Officer can confirm at a glance that the four Article 32 pillars — confidentiality, integrity, availability, resilience — are each addressed with named primitives and documented procedures.

> **Why this matters.** Article 32 is the operational pillar of GDPR. A controller that cannot show appropriate technical and organizational measures fails Article 32 even if every other GDPR provision is satisfied. The chain's integrity story is well known; the confidentiality, availability, and resilience story is less visible. This document fills that gap.

---

## 1. Article 32 in plain terms

GDPR Article 32(1) says the controller and processor "shall implement appropriate technical and organisational measures to ensure a level of security appropriate to the risk," and it names four specific properties as examples: pseudonymisation and encryption (32(1)(a)), confidentiality, integrity, availability and resilience (32(1)(b)), the ability to restore availability after an incident (32(1)(c)), and a process for regular testing of effectiveness (32(1)(d)).

The chain answers each of these. The mapping below is the institution's answer to the supervisory authority's standard Article 32 question.

---

## 2. The four pillars mapped to chain primitives

### 2.1 Integrity (Article 32(1)(b))

| Layer | Primitive | What it catches | Reference |
|---|---|---|---|
| Per-entry | HMAC-SHA-256 over canonicalized JSON | Modification of any single chain entry | Spec §4.1, NIST FIPS 198-1, RFC 2104 |
| Per-day | Merkle root over the day's entries | Insertion or deletion of entries within a sealed day | Spec §4.2, RFC 6962 |
| Per-seal | HSM-backed Ed25519 signature over the seal record | Forgery of the seal itself; un-attributed seal records | Spec §4.3, NIST FIPS 186-5, RFC 8032 |

Three independent detection layers. Hidden tampering would require simultaneous compromise of the IKM, the ledger storage, and the HSM — three separately-hardened systems with separate access controls. Verifier procedure (spec §7) is deterministic, public, and enumerates failure reasons per step. Test-vector corpus (spec §8) ships positive and negative vectors for over 20 attack scenarios.

### 2.2 Confidentiality (Article 32(1)(a) and (b))

| Property | Control | Reference |
|---|---|---|
| Encryption at rest | Chain entries and seal records encrypted using **AES-256-GCM** with a confidentiality key distinct from the HMAC IKM | Spec §10 (security overview); institution's encryption policy |
| Encryption in transit | TLS 1.3 between SDK and ledger; mTLS between ledger services | Spec §4.1.1 (session-key handshake security floor) |
| Key custody | Confidentiality key stored in HSM (FIPS 140-2 Level 3) or cloud-KMS with hardware-backed key protection; **separate from the chain's HMAC IKM and from the privacy-store key** | Spec §10.6, `cloud-hsm-guide.md`, `token-vault-architecture.md` |
| Pseudonymisation | Personal-information fields tokenized before canonicalization; tokens are HMAC-SHA-256 of original PII under a privacy-store key held in separate custody | `privacy-by-design.md`, `pseudonymization-vs-anonymization.md` |
| Public-disclosure mode | "De-key-ing" the chain (publishing verifier output without plaintext entries) satisfies confidentiality for any disclosure where the verifier transcript is the artifact | `litigation-support.md` § verifier transcript |

The institution holds three distinct keys with three distinct custody chains: the chain HMAC IKM, the chain confidentiality key (AES-GCM at rest), and the privacy-store key (token reversibility). A breach of any one does not automatically break the others. Key-rotation calendars are coordinated to prevent simultaneous rotation per `token-vault-architecture.md`.

### 2.3 Availability (Article 32(1)(b) and (c))

| Property | Control | Reference |
|---|---|---|
| Backup | Ledger snapshots taken daily; backup integrity verified by re-running the verifier on the snapshot | Spec §10.10 (backup-integrity procedures) |
| Disaster recovery | RPO and RTO defined per `06-ledger-server-design.md` §7.5; documented recovery sequence | `dr-and-resilience.md` |
| Recovery testing | Quarterly recovery drill: restore from backup, run verifier, confirm chain integrity. Drill outcome logged as `master.recovery_drill_completed` operational event | `incident-response-playbook.md` Scenario 10 |
| Capacity | Hot-path budget defined; streaming Merkle scales linearly with entry rate | Spec §5 (capacity), `02-chain-construction.md` §5 |
| Multi-region topology (where deployed) | Active-active or active-passive per institution's resilience policy; cross-region replication encrypted in transit and at rest | `byoc-deployment.md` topology section |

Recovery from backup is tested rather than assumed. The drill is logged and the outcome is part of the institution's evidence file for examiners and DPAs.

### 2.4 Resilience (Article 32(1)(b) and (d))

| Property | Control | Reference |
|---|---|---|
| Incident-response playbook | 15 named scenarios, each with detection signal, containment action, and recovery procedure | `incident-response-playbook.md` |
| Annual red-team exercise | The chain operations team participates in an annual red-team exercise against a chain-deployment scenario; findings drive playbook updates | Institution's security-operations program |
| Continuous monitoring | Operational events (`master_key.retired`, `audit_file.truncation_detected`, `key_fingerprint.mismatch`, etc.) are alerted to the SOC | Spec §10.2 operational-event catalog |
| Verifier as continuous audit | Verifier runs are scheduled (daily or per-seal); failures escalate to the IR team | Spec §7, `audit-procedures.md` |
| Effectiveness testing (Article 32(1)(d)) | Quarterly verifier-procedure test using the public test-vector corpus; annual independent audit of chain controls | Spec §8 test vectors; institution's audit-procedures schedule |

The resilience story is not "we hope nothing breaks." It is "we exercise our recovery procedures on a published cadence, and we have a written record of each exercise."

---

## 3. The risk-appropriate test (Article 32(1) opening clause)

Article 32(1) opens with "taking into account the state of the art, the costs of implementation and the nature, scope, context and purposes of processing as well as the risk of varying likelihood and severity for the rights and freedoms of natural persons." The institution's risk-appropriate analysis:

| Risk dimension | Assessment |
|---|---|
| State of the art | HMAC-SHA-256, SHA-256, Ed25519, AES-256-GCM, RFC 6962 Merkle, TLS 1.3, FIPS 140-2 L3 HSM custody — all current per NIST guidance |
| Cost of implementation | Cryptographic primitives are FIPS-approved and widely available; HSM and KMS services are commercially available; the verifier is open-source. Cost is proportionate to a regulated-financial-services deployment |
| Nature of processing | Integrity-bearing audit trail of automated decisions affecting credit, fraud, pricing, and (in some deployments) healthcare-financed lending |
| Scope and context | Multi-year retention, multi-jurisdiction in some deployments, special-category data possible per `hipaa-minimum-necessary.md` |
| Risk to rights and freedoms | High where decisions have legal or similarly significant effect (Article 22); the chain's tamper-evidence directly supports the data subject's right to challenge an automated decision |

The institution concludes the chosen controls are appropriate to the risk. The conclusion is recorded in the DPIA and revisited annually.

---

## 4. The "regular testing" obligation (Article 32(1)(d))

Article 32(1)(d) requires "a process for regularly testing, assessing and evaluating the effectiveness of technical and organisational measures." The institution's testing program:

| Cadence | Test | Output |
|---|---|---|
| Daily | Verifier runs on previous day's seal | Pass/fail flag; failures escalate to IR |
| Weekly | Key-fingerprint reconciliation per spec §10.1 | `master.reconciliation_completed` event; mismatches escalate |
| Quarterly | Recovery-from-backup drill | `master.recovery_drill_completed` event; result recorded |
| Quarterly | Verifier-procedure test against public test-vector corpus | Pass/fail across 20+ attack vectors; record retained |
| Annually | Red-team exercise against chain deployment | Findings document; playbook updates |
| Annually | Independent audit of chain controls (internal audit or external) | Audit report; management response |

Each test produces an artifact that is retained per the institution's records-retention schedule. The artifact set is the institution's Article 32(1)(d) evidence pack.

---

## 5. Cross-references

- Spec §4 (chain construction), §7 (verifier), §10 (security overview, key custody, IR).
- `incident-response-playbook.md` — detection, containment, recovery procedures including red-team scenarios.
- `dr-and-resilience.md` — RPO, RTO, multi-region topology.
- `cloud-hsm-guide.md` — HSM and KMS custody for the chain HMAC IKM and seal signing key.
- `token-vault-architecture.md` — separate custody for the privacy-store reversibility key.
- `breach-notification-matrix.md` — escalation when an Article 32 control fails.
- `dpia-template.md` §3 — risk analysis incorporates this mapping.
- `ropa-template.md` "Security measures" column references this document.

---

## 6. Review cadence

This mapping is reviewed annually by the privacy office and the security-operations team. Triggers for early review:

- A new cryptographic algorithm enters or leaves NIST current guidance (e.g., post-quantum migration).
- A regulator issues updated Article 32 guidance.
- An incident reveals a control gap; the playbook is updated and this document is revisited.
- The institution adopts a new topology (e.g., moves from BYOC to vendor-hosted) that changes which entity operates which control.
