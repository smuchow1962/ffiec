# Token Vault Architecture — Privacy-Store Custody Requirements

> **What this doc is.** Architecture and operational requirements for the privacy-store, also known as the token vault — the system that holds the token-to-original-PII mapping that resolves chain tokens. Written so an architect designing a chain deployment, a security reviewer assessing the deployment, and an examiner reading the deployment package can confirm the privacy-store is custodied separately from the chain in a way that preserves pseudonymization (per `regulator-pack/pseudonymization-vs-anonymization.md`).

> **Why this matters.** The chain's privacy story rests on one architectural claim: that the additional information needed to re-identify chain tokens is "kept separately and... subject to technical and organisational measures" per GDPR Article 4(5). If the privacy-store is co-located with the chain, operated by the same team, and protected by the same key, the pseudonymization claim collapses. The chain becomes a personal-data ledger with stronger compliance obligations and weaker breach posture. This document is the institution's specification for keeping the privacy-store separate.

---

## 1. The role of the privacy-store

The privacy-store is the lookup that resolves a chain token back to the original PII. The chain captures `token = HMAC-SHA-256(privacy_key, original)`; the privacy-store carries the `(token, original)` pair so that authorized data-subject-rights workflows can map back from a token to a customer identity, an account, or a piece of personal data.

The privacy-store has three jobs:

1. Hold the token-to-original mapping.
2. Hold the privacy_key (the HMAC key used to produce the tokens).
3. Provide an authorized lookup path for data-subject-rights workflows (Right to Know, Right to Delete, breach response, audit).

Everything else — chain integrity verification, day-to-day operations, the verifier procedure — runs without privacy-store access. This separation is the architectural lever that makes the privacy-store a controllable, auditable, deletable component.

---

## 2. The seven custody requirements

### 2.1 Logical separation (required)

The privacy-store is a separate service or database from the chain ledger. A compromise of the chain ledger does **not** automatically grant access to the privacy-store. Specifically:

- Different network segment.
- Different authentication system or distinct credential set.
- Different database engine, instance, or schema.
- Different access-control policy (different IAM roles, different service accounts).

A chain-operations engineer with full access to the ledger cannot, by virtue of that access, query the privacy-store. The privacy-store is reachable only through the privacy-data-protection team's access path.

### 2.2 Physical separation (recommended, required for high-security deployments)

For institutions handling healthcare PHI, special-category GDPR data, or high-volume sensitive data, the privacy-store database is hosted on:

- A different physical server (not just a different VM on the same hypervisor).
- A different data center or availability zone.
- A different cloud provider account or, in the highest-security cases, a different cloud provider altogether.

Physical separation defeats single-host compromise. An adversary who roots the chain ledger host does not gain network reach to the privacy-store host. Cloud deployments achieve this through separate accounts/subscriptions with no shared IAM trust.

### 2.3 Encryption with a separate key (required)

Privacy-store contents are encrypted at rest using a key **distinct from**:

- The chain's HMAC IKM (used for per-entry HMAC).
- The chain's confidentiality key (used for AES-GCM at rest of chain entries per `regulator-pack/article-32-security-mapping.md`).
- Any other system's encryption key.

The privacy-store key is HSM-protected (FIPS 140-2 Level 3 or higher) or cloud-KMS-protected (with hardware-backed key storage). The key is generated within the HSM/KMS and never exported. Key custody is held by the privacy-data-protection team, not chain operations.

The institution maintains three independent key chains:

| Key | Purpose | Custody |
|---|---|---|
| Chain HMAC IKM | Per-entry chain integrity | Chain-operations team in HSM/KMS |
| Chain confidentiality key | AES-GCM encryption of chain at rest | Chain-operations team in HSM/KMS |
| Privacy-store key | Token reversibility (HMAC privacy_key) plus encryption of mapping table | Privacy-data-protection team in separate HSM/KMS |

A compromise of any one key does not compromise the others. Rotation calendars are coordinated to prevent simultaneous rotation (failure-mode-of-coincidence risk).

### 2.4 Access control by separate team (required)

Privacy-store access is restricted to the **privacy-data-protection team** — a team distinct from chain operations.

| Team | Access |
|---|---|
| Chain operations | Chain ledger, HMAC IKM, confidentiality key, verifier output. **No privacy-store access** |
| Privacy-data-protection team | Privacy-store, privacy-store key, mapping table, deletion authority. **No chain-ledger write access** |
| Incident response (chain-side) | Read access to chain ledger and operational events; coordinates with privacy-data-protection IR for combined incidents |
| Incident response (privacy-side) | Owned by privacy-data-protection team; coordinates with chain-side IR |
| Internal audit | Read-only access to both, scoped per audit engagement; access logged |
| Regulators (during exam) | Access mediated by audit and legal teams; never direct |

The access matrix is documented and reviewed quarterly. Role changes (joiners, movers, leavers) trigger access reviews.

### 2.5 Access logging (required)

Every privacy-store access produces an audit log entry:

| Field | Content |
|---|---|
| Timestamp | UTC, microsecond precision |
| Operator identifier | Unique user ID; service-account access uses the calling service identity |
| Operation | `read`, `write`, `delete`, `key_rotation`, `audit_export` |
| Scope | Token range, customer-ID range, or "full table" — whatever was accessed |
| Justification | Reference to a workflow or ticket (data-subject request ID, IR case ID, audit engagement ID) |
| Outcome | Success, denied, or partial |

Logs are retained per the institution's records-retention schedule (matched to the 7-year FFIEC floor per `regulator-pack/retention-justification.md`). Logs are reviewed monthly; anomalies (out-of-hours access, unjustified-justification fields, unusual scope) escalate to the Security Official.

### 2.6 Separate incident response (required)

A privacy-store compromise is a different incident from a chain compromise. The IR procedure reflects this:

| Incident | Primary IR team | Chain-operations involvement |
|---|---|---|
| Chain ledger compromise | Chain-side IR | Primary |
| HMAC IKM leak | Chain-side IR | Primary |
| Privacy-store compromise | Privacy-side IR | Coordinated awareness; chain-side does not take remediation action on the vault |
| Privacy-store key compromise | Privacy-side IR | Coordinated awareness |
| Combined (both compromised simultaneously) | Joint command | Both teams; legal and privacy office join |

Each IR team has its own playbook, its own escalation tree, its own external contacts. The chain operations team is **notified** of privacy-store incidents but does not perform remediation on the privacy-store (it lacks access).

### 2.7 Custody coordination (required)

Two custody chains independently maintained still need to coordinate on shared events. The coordination points:

| Event | Coordination |
|---|---|
| Chain HMAC IKM rotation | Privacy-data-protection team is notified; privacy-store key rotation is scheduled at least 30 days before or after to avoid simultaneous rotation |
| Privacy-store key rotation | Chain-operations team is notified; chain confidentiality-key rotation is similarly offset |
| Annual disaster-recovery drill | Both teams participate; restored ledger and restored privacy-store are both verified |
| Annual audit | Audit team coordinates with both teams; audit findings are shared |
| Cross-region replication setup | Joint architecture decision; replication topology documented for both systems |
| End-of-retention deletion | Privacy-data-protection team performs the deletion (per `regulator-pack/retention-justification.md` §5); chain operations confirms tokens become anonymized |

The institution maintains a shared calendar for these events to prevent collision and to ensure deliberate sequencing.

---

## 3. Anti-patterns to avoid

| Anti-pattern | Why it fails |
|---|---|
| Privacy-store as a column in the chain table | Kills logical separation; one compromise breaks both |
| Same encryption key for chain and privacy-store | Kills cryptographic separation; key compromise reverses both |
| Chain operations holds privacy-store credentials | Kills team separation; one operator can map any token |
| Privacy-store on same host, accessed via localhost | Kills physical separation; host compromise grants both |
| Same backup destination for both | Kills separation in the backup path; restore can collide |
| Shared rotation calendar | Kills failure-mode-of-coincidence protection; simultaneous rotation amplifies risk |
| Privacy-store as a SaaS the chain vendor also operates | Possible but requires contractual separation that mirrors team separation; default is to use a different vendor |

If any of these patterns are present, the institution cannot defensibly claim the chain's tokens are pseudonymized per GDPR Article 4(5). The compliance posture collapses.

---

## 4. Vendor-hosted multi-tenant deployments

When a vendor hosts the chain for multiple institutional customers, the privacy-store custody question becomes:

| Option | Where the privacy-store lives | Vendor access |
|---|---|---|
| **Option A — Customer-held privacy-store (recommended)** | Each customer hosts its own privacy-store on its own infrastructure; customer's SDK tokenizes before transmitting; vendor never sees originals | Vendor sees only tokens; vendor cannot map tokens to PII |
| **Option B — Vendor-hosted privacy-store with customer-held key** | Vendor hosts the database; the privacy-store key is held by the customer in customer's HSM; vendor cannot decrypt without customer cooperation | Vendor stores ciphertext; customer-held key required for any read |
| **Option C — Vendor-hosted privacy-store with vendor-held key** | Vendor hosts the database and the key | Vendor has full reversibility access; pseudonymization claim collapses for the customer's data |

Option A is the institution's preferred posture. Option B is acceptable with appropriate contractual safeguards and key-custody attestation. Option C is **not acceptable** for institutions claiming GDPR-compliant pseudonymization.

`regulator-pack/international-transfers.md` Scenario D covers the multi-tenant case in more detail; the architecture choices here align with the transfer analysis there.

---

## 5. Disaster recovery and the privacy-store

The privacy-store has its own backup, recovery, and resilience requirements distinct from the chain's. Critical points:

- **Backup encryption**: Privacy-store backups use the same separate-key model. Backups are encrypted with the privacy-store key; restoring requires the same key custody.
- **Recovery target**: Privacy-store recovery should not lag chain recovery by more than the institution's RPO/RTO; otherwise the chain is operational but unable to honor data-subject rights.
- **Quarterly drill**: A privacy-store recovery drill is performed quarterly. The drill confirms backup integrity, key availability, and team readiness.
- **Geographic placement**: Privacy-store backups follow the same jurisdiction rules as the chain. EU data → EU-stored privacy-store backups (see `regulator-pack/international-transfers.md`).

---

## 6. Lifecycle management

| Phase | Action |
|---|---|
| Initial deployment | Privacy-store provisioned, key generated in HSM/KMS, access matrix established, team training completed |
| Ongoing operation | Daily access logs monitored; weekly key-fingerprint health check; quarterly access review; annual key rotation |
| Key rotation | New privacy_key generated; existing tokens remain valid (HMAC keys are versioned; old version retained for lookup); rotation logged |
| End-of-life mapping deletion (per record) | At end of 7-year retention or on erasure request: mapping row deleted; deletion logged; chain tokens become anonymized |
| End-of-life vault decommission (full system) | Cryptographic erasure of the privacy-store key; physical destruction of storage media if required by data-disposal policy |

---

## 7. Cross-references

- `privacy-by-design.md` — SDK-side tokenization that produces the tokens stored here.
- `regulator-pack/pseudonymization-vs-anonymization.md` — explains why this architecture is what makes pseudonymization defensible.
- `regulator-pack/article-32-security-mapping.md` §2.2 — confidentiality controls; this document is the operational detail.
- `regulator-pack/retention-justification.md` §5 — describes the end-of-retention deletion procedure performed against this vault.
- `regulator-pack/hipaa-security-rule-mapping.md` — physical and technical safeguards that apply equally to the privacy-store.
- `regulator-pack/ccpa-cpra-rights.md` — Right to Delete deletion procedure operates against this vault.
- `regulator-pack/international-transfers.md` Scenario D — multi-tenant vendor-hosted deployment options.
- `cloud-hsm-guide.md` — HSM and KMS guidance applies to the privacy-store key as well as the chain keys.
- `incident-response-playbook.md` — privacy-side IR scenarios.
- Spec §10 (security overview) — chain-side security, distinct from this document.

---

## 8. Review cadence

This document is reviewed annually by the privacy office, security architecture, and chain operations. Triggers for early review:

- The institution adopts a new HSM, KMS, or database technology for the privacy-store.
- A new deployment topology is added (e.g., multi-region active-active for the privacy-store).
- A regulator publishes guidance affecting pseudonymization-store custody.
- An incident reveals a custody gap; the architecture is updated.
- A new vendor relationship requires re-evaluation of multi-tenant options.
