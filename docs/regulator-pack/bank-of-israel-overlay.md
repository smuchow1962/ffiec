---
title: Bank of Israel Articulation Overlay — Directives 357 / 359 / 361 / 365 / 367 / 411 / 414, PPL Amendment 13, INCD coordination, Equal Opportunity Employment Law
status: informative
aligned-with:
  - Bank of Israel Proper Conduct of Banking Business Directive 357 (data management and third-party risk — תיעוד וניהול סיכון צד שלישי)
  - Bank of Israel Proper Conduct of Banking Business Directive 359 (cyber defense management — ניהול הגנה סייבר)
  - Bank of Israel Proper Conduct of Banking Business Directive 361 (cyber risk management — ניהול סיכוני סייבר)
  - Bank of Israel Proper Conduct of Banking Business Directive 365 (operational resilience — עמידות תפעולית)
  - Bank of Israel Proper Conduct of Banking Business Directive 367 (cloud computing and GenAI — חישוב ענן וינטליגנציה מלאכותית)
  - Bank of Israel Proper Conduct of Banking Business Directive 411 (incident reporting — דיווח על תקלות)
  - Bank of Israel Proper Conduct of Banking Business Directive 414 (third-party risk — סיכון צד שלישי)
  - Privacy Protection Law (PPL) Amendment 13 (effective August 2025)
  - Cyber Defense Law 5778-2018 (חוק הגנה סייבר)
  - Israeli National Cyber Directorate Defense Methodology for Organizations v2.0
  - Equal Opportunity in Employment Law (1988, amended 2022 for automated decisions)
  - Bank of Israel-Federal Reserve cross-border banking supervision MoU
date: 2026-05-07
version: 1.0.0
---

# Bank of Israel Articulation Overlay

> **What this doc is.** A single articulation overlay that maps the FFIEC chain-of-custody v1.0a specification onto the Israeli banking supervisory framework. Written so a competent authority — a Bank of Israel Banking Supervisor (פיקוח על הבנקים), an Israeli National Cyber Directorate (INCD) liaison, an Israeli Privacy Protection Authority examiner, or a Tier-1 Israeli megabank's CISO preparing for Bank of Israel examination — can read this document alongside the spec and confirm what the chain delivers, what it does not, and where supplementary measures are required for application in Israeli jurisdiction.

> **What this doc is NOT.** Not a normative extension to the v1.0a specification. Not a translation of the spec into Israeli supervisory language. Not an Israel-specific spec fork. The integrity primitives in v1.0a are framework-neutral cryptographic constructs; the supervisory translation is an articulation overlay an institution layers on top, not a change to the underlying specification. Spec §1.2 (epistemic scope) governs: the chain proves what was said and that the record was not tampered with after capture; it does not prove the substantive correctness of the captured content.

> **Threat model.** Israeli banks operate under the INCD assumption that a capable nation-state adversary will compromise bank infrastructure within 18 months. The overlay is written against this baseline, not against a theoretical-intrusion threat model. The chain's value under this baseline is detection and proof, not prevention; the overlay names the institutional controls that supplement the chain's primitives to produce a full operational-resilience and cyber-resilience posture under Bank of Israel and INCD oversight.

---

## 1. Scope and reading order

The overlay is consumed in five reading orders.

| Reader | Reading order |
|---|---|
| Bank of Israel Banking Supervisor (Pikuach HaBankim) | §3 (Directive 365 operational resilience) → §4 (Directive 367 cloud and PPL Amendment 13) → §5 (Directive 361 INCD coordination) → §7 (Directive 414 vendor audit) → §11 (BoI examination procedures) → §15 (translation table) |
| Israeli National Cyber Directorate liaison | §5 (Directive 361 INCD coordination) → §6 (Directive 359 Banking Cyber CERT) → §10 (nation-state threat model) → §13 (red-team readiness) |
| Tier-1 megabank CISO (Tal-style operational lens) | §10 (nation-state threat model) → §12 (supply-chain discipline) → §13 (red-team readiness) → §14 (incident-response procedures) → §11 (examination preparation) |
| Privacy Protection Authority examiner | §4 (Directive 367 + PPL Amendment 13) → §9 (PPL-A13 specifics) |
| Cross-border (USD correspondent) institution | §8 (BoI-Federal Reserve coordination) → §14 (multi-regulator IR) |

The overlay does not duplicate material in adjacent regulator-pack documents. Cross-references are exact — if a section sends the reader to `dora-articulation-overlay.md` §8, that section is the load-bearing source and this overlay names where it sits in the Israeli framework picture.

---

## 2. The Israeli supervisory architecture and the chain's place in it

### 2.1 The four principal authorities

| Authority | Mandate | Chain interaction |
|---|---|---|
| Bank of Israel Banking Supervision Department (פיקוח על הבנקים) | Issues directives binding on the five megabanks (Bank Hapoalim, Bank Leumi, Discount Bank, Mizrahi Tefahot, FIBI) and other supervised institutions; conducts examinations under Directives 357, 359, 361, 365, 367, 411, 414 | Examiner-side consumer of verifier output (spec §7); reviewer of CC8.1 documentation; coordinator with FFIEC examiners for institutions with USD correspondent operations |
| Israeli National Cyber Directorate (INCD) | National cybersecurity authority under Cyber Defense Law 5778-2018; operates Banking Cyber CERT under Directive 359; conducts annual red-team exercises | Receives chain-detected incident notifications per Directive 361 §5; participates in red-team scenarios that test chain resilience |
| Privacy Protection Authority (PPA) | Enforces Privacy Protection Law including Amendment 13 (effective August 2025); regulates cross-border transfer of sensitive personal information including AI inference logs | Reviews chain's data-handling posture under PPL-A13; receives breach notifications affecting privacy-store custody |
| Federal Reserve (for institutions with USD correspondent operations) | Examines US-side cross-border operations of Israeli megabanks under FFIEC standards | FFIEC examiner-side consumer of verifier output; coordinates with Bank of Israel under the BoI-Federal Reserve MoU |

### 2.2 Directive landscape

The Bank of Israel issues directives as binding operational mandates — not guidance, not frameworks, but rules with teeth. Banks that do not comply face capital surcharges, operations restrictions, or board-level supervisory action. The seven directives load-bearing for chain deployment are enumerated below; each is articulated in detail in the corresponding section of this overlay.

| Directive | Subject | Section in this overlay |
|---|---|---|
| 357 (תיעוד וניהול סיכון צד שלישי) | Data management and third-party risk | §7 |
| 359 (ניהול הגנה סייבר) | Cyber defense management; INCD Banking Cyber CERT integration | §6 |
| 361 (ניהול סיכוני סייבר) | Cyber risk management; INCD coordination on incidents | §5 |
| 365 (עמידות תפעולית) | Operational resilience; ≤2-hour recovery for critical operations | §3 |
| 367 (חישוב ענן וינטליגנציה מלאכותית) | Cloud computing and GenAI; data localisation | §4 |
| 411 (דיווח על תקלות) | Incident reporting | §14 |
| 414 (סיכון צד שלישי) | Third-party risk and annual vendor audit | §7 |

---

## 3. Directive 365 — operational resilience and ≤2-hour recovery

### 3.1 What Directive 365 requires

Directive 365 §3 (עמידות תפעולית — operational resilience) mandates that Israeli banks recover from the failure of any critical operational system within 2 hours. The failure case includes both tampering detection (the chain's daily seal mismatches) and unavailability (the HSM is offline, the ledger is corrupted, the seal job has crashed). Directive 365 §4 mandates annual operational-resilience drills.

### 3.2 Conformant deployment patterns

The spec's daily Merkle seal cadence produces a 24-hour detection window. Three deployment patterns satisfy Directive 365's 2-hour recovery target.

| Pattern | Detection window | Recovery time | Operational cost | Best suited for |
|---|---|---|---|---|
| Hourly seal cadence | 1 hour | ~2 hours (IKM rotation at hour N, seals issued at hour N+1) | ~24x daily seal HSM operational load | Tier-1 megabanks with high-value operations |
| Multi-region active-active failover (spec §10.15 Pattern A) | 1 hour (per-region) | ~30 minutes (failover from primary HSM region to secondary) | Two HSMs (primary + secondary); cross-region replication | Tier-1 megabanks with multi-region BCP |
| Hot-standby HSM with 15-minute synchronisation | Daily at default (institutionally tunable) | ~15 minutes failover | Two HSMs with active synchronisation | Tier-1 megabanks with multi-region BCP and the DR posture to support hot-standby |

The institution's CC8.1 names the chosen pattern. Sole reliance on default daily cadence with 24-hour recovery is non-compliant with Directive 365 §3; institutions defaulting to daily cadence file a formal exception request to the Bank of Israel that names the specific compensating controls (e.g., out-of-band integrity monitoring, ledger-storage SIEM coverage, real-time HSM health monitoring) that bring the effective recovery time within 2 hours under realistic incident scenarios.

### 3.3 Directive 365 §4 annual operational-resilience drill — audit procedure P-12

The institution conducts an annual operational-resilience drill for the chain. The drill scenario is selected from: primary HSM failure, ledger storage corruption, seal-job process crash, regional infrastructure failure (data center power loss, network partition).

| Drill step | Activity | Evidence captured |
|---|---|---|
| 1. Pre-drill | Scenario announced to operational team 1 week in advance; team prepares recovery procedures but does not execute | Drill plan filed in Directive 365 evidence repository |
| 2. Drill execution | At drill start time T0, operational team executes the failure simulation | T0 timestamp recorded |
| 3. Recovery start | At T1, recovery procedures begin (IKM rotation, HSM failover, seal-job restart) | T1 timestamp recorded |
| 4. First seal post-recovery | At T2, the first seal is successfully signed and published | T2 timestamp recorded; recovery time = T2 - T0 |
| 5. Pass/fail determination | Recovery time ≤ 2 hours: PASS. Recovery time > 2 hours: FAIL with remediation plan and follow-up drill within 30 days | Drill report filed in Directive 365 evidence repository |
| 6. Annual report | Institution files annual Directive 365 report with the Bank of Israel naming drill date, scenario, recovery time, pass/fail, and any procedure improvements | Report filed with Bank of Israel |

The Bank of Israel supervisor reviews the past 3 years of drill results during examination. All drills showing recovery time ≤ 2 hours indicates the control is operationally effective under Directive 365. Any failure is investigated; repeated failures (>1 in past 3 years) indicate the institution's recovery-time RTO is not achievable with the current architecture and require redesign.

### 3.4 Pattern A vs Pattern B for Directive 365

Spec §10.15 Pattern A (active-active per-region with seal-region pinning) is the preferred pattern for Directive 365 compliance. Each region has independent IKM and ledger, so region failure is contained; failover involves starting a new run in the alternate region rather than transferring cross-region state; failover time is typically 5 minutes for process restart and 30 minutes to full operational capacity.

Pattern B (per-region tenant separation) requires global IKM coordination and cross-region ledger coordination; failover involves tenant re-routing which risks state inconsistency. Pattern B is operationally complex for Directive 365 recovery and is not recommended for Israeli megabanks unless the institution has specific multi-region architecture constraints that force the choice. Pattern B selection requires Bank of Israel approval and may require additional compensating controls (real-time cross-region ledger replication and cryptographic proof of consistency).

---

## 4. Directive 367 cloud localisation + PPL Amendment 13 cross-border transfer

### 4.1 The compounding constraint

Directive 367 §2 (חישוב ענן וינטליגנציה מלאכותית — cloud computing and GenAI) requires that AI decision logs and the data that feeds into AI decisions remain in Israeli jurisdiction unless the bank obtains explicit Bank of Israel approval. PPL Amendment 13 (August 2025) goes further: it treats AI inference logs and AI-generated decisions as "sensitive personal information" (זיכרון רגיש) — the same classification as biometric or genetic data — and prohibits their transfer outside Israel without explicit Privacy Protection Authority (PPA) approval.

The constraint is compounding: an Israeli bank operating chain-of-custody must satisfy both Directive 367 (Bank of Israel approval) and PPL Amendment 13 (PPA approval) for any cross-border transfer of chain content.

### 4.2 Compliant deployment topologies

| Topology | Directive 367 conformance | PPL Amendment 13 conformance |
|---|---|---|
| **Self-hosted in Israeli data center / colocation** (on-prem Thales Luna, Entrust nShield, or Utimaco SecurityServer) | Conformant by default; the bank retains full institutional custody | Conformant by default; chain content does not cross Israeli border |
| **Vendor-hosted in Israeli jurisdiction only** (AWS Tel Aviv region, Azure Israel region if available, or vendor's Israeli subsidiary's infrastructure) | Conformant if the vendor commits in contract to Israeli localisation | Conformant if the vendor commits in contract to Israeli localisation; sensitive personal information does not cross Israeli border |
| **Vendor-hosted in non-Israeli jurisdiction with PPA + Bank of Israel approvals** | Requires Bank of Israel exception per Directive 367 | Requires PPA approval for transfer of sensitive personal information; PPA approval typically requires SCC + supplementary safeguards (encryption key held by bank, vendor cannot decrypt, audit access) |
| **Vendor-hosted multi-tenant US/EU cloud** (without Israeli-region commitment) | Non-conformant without explicit Bank of Israel approval | Non-conformant without explicit PPA approval |

### 4.3 Vendor contract requirements (Directive 367 + Directive 357)

If using a vendor (either Israeli-hosted or approved non-Israeli with PPA clearance), the bank's contract with the vendor names:

- Data localisation commitment (Israeli storage, with attestation).
- Audit rights (bank can audit vendor's storage and access logs on demand).
- Key-holding commitment (if encryption is required, bank holds keys).
- Regulatory coordination (vendor commits to INCD incident notification within 1 hour and Bank of Israel examination access).
- Subcontractor management (if vendor uses third-party HSM provider, vendor confirms HSM is in Israeli jurisdiction).
- Background checks for vendor personnel with access to HSM, ledger, or IKM.

### 4.4 Privacy-store custody under PPL Amendment 13

The token-to-PII mapping (the privacy-store) MUST remain in Israeli jurisdiction under all circumstances; it cannot be replicated to off-shore backup. The privacy-store is a high-risk processing activity under PPL Amendment 13.

| Privacy-store control | PPL Amendment 13 requirement |
|---|---|
| Custody | In Israeli jurisdiction; encrypted at rest |
| Access controls | Only designated personnel (privacy team, customer service for disputes); access logs retained for 1 year |
| Retention | Per institution's stated retention period (typically 7 years for dispute resolution) or upon customer erasure request, whichever is sooner |
| Customer notification | Customer notified at time of decision: "your AI-decision log is retained for [purpose] for [period]" |
| Data subject rights | Access, rectification, erasure, objection — fulfilled within 30 days |
| Breach notification | PPA within 48 hours; affected individuals within 60 days |

---

## 5. Directive 361 — cyber risk management and INCD coordination

### 5.1 The 1-hour INCD clock

Directive 361 §5 requires Israeli banks to coordinate cyber-incident response with the Israeli National Cyber Directorate in real time. For a suspected or confirmed nation-state cyber incident — APT compromise, credential exfiltration, HSM firmware tampering, ledger tampering — the bank must notify the INCD within 1 hour of determination.

The chain's primary outputs do not produce the determination directly. The determination is made by the bank's incident-response team based on available evidence (SOC alerts, intrusion-detection systems, anomalous logs, user reports, initial forensics). The chain's evidence (verifier output, daily seal verification result, fingerprint reconciliation) supplements the determination but does not gate it. The bank notifies INCD on preliminary determination; chain evidence follows in the 24-hour follow-up.

### 5.2 Severity taxonomy

Directive 361 §5 classifies incidents into four severity tiers. The bank's IR runbook names the determination procedure for each tier.

| Tier | Description | INCD notification | Bank of Israel notification |
|---|---|---|---|
| **Red (אדום)** | Suspected nation-state compromise: HSM firmware tampering detected, IKM exfiltration suspected, ledger tampering detected by verifier | Within 1 hour | Within 1 hour (parallel) |
| **Orange (כתום)** | Significant cyber incident: vendor breach affecting HSM or ledger, supply-chain compromise detected, OTLP collector tampering suspected | Within 4 hours | Within 4 hours (parallel) |
| **Yellow (צהוב)** | Operational incident with cyber implications: HSM unavailable due to maintenance/hardware failure, ledger storage corruption, key rotation procedural error | Optional | Within 24 hours |
| **Green (ירוק)** | Minor operational issues: single verifier failure traced to known transaction-ordering issue, seal-age delay <75 minutes, anomalous but explained event | None | Internal incident log only |

### 5.3 Multi-regulator incident notification for cross-border banks

Israeli megabanks holding USD correspondent accounts are subject to FFIEC notification under 12 CFR §53 (36-hour clock) in addition to Directive 361 (1-hour INCD clock) and Directive 411 (Bank of Israel notification). The bank's IR commander operates the **tightest applicable clock first** under the simultaneous-notification rule.

| Regulator | Trigger | Clock | Notification format |
|---|---|---|---|
| INCD | Red or Orange chain-related incident | 1 hour (Red) or 4 hours (Orange) | Cyber Defense Liaison phone + formal incident report |
| Bank of Israel | Red, Orange, or Yellow chain-related incident | 1 hour, 4 hours, or 24 hours | Banking Supervisor on-call number; formal operational-incident report |
| Federal Reserve (for institutions with USD correspondent operations) | Computer-security incident affecting correspondent operations | 36 hours per 12 CFR §53 | FFIEC Computer-Security Incident form |
| PPA | Privacy-store breach | 48 hours per PPL Amendment 13 | PPA breach notification form |

The bank's IR runbook names the parallel notification procedure: a single triage produces three (or four, including PPA) parallel notifications with the same core content formatted per each regulator's requirements.

### 5.4 Key rotation post-compromise (Scenario 4 from incident-response playbook)

Directive 361's compromise-recovery posture composes with the spec's Scenario 4 (Master-key compromise or suspected compromise) procedure. The institution's CC8.1 names:

- The IKM rotation procedure (issue new key_version, deploy to HSM, fingerprint-registry update synchronously).
- The re-attestation period (typically 7 or 30 days post-rotation) during which entries using the new key are treated as "high confidence, pending full attestation."
- The INCD coordination procedure during the re-attestation period.
- The forensic evidence preservation procedure (HSM audit logs, ledger snapshots, network captures).

---

## 6. Directive 359 — INCD Banking Cyber CERT integration

### 6.1 The 24-hour CERT reporting clock

Directive 359 (ניהול הגנה סייבר — cyber defense management) requires Israeli banks to participate in the INCD's Banking Cyber CERT (Coordination and Emergency Response Team) and to share indicators of compromise (IOCs) and threat intelligence with the CERT within 24 hours of detection.

Chain-detected anomalies that trigger CERT reporting:

- Fingerprint mismatch outside a documented rotation.
- Seal-age delay suggesting HSM interference.
- Unexplained verifier failures.
- OTLP collector tampering signatures.
- Supply-chain compromise indicators (verifier binary SHA-256 mismatch, SDK build non-reproducible).

### 6.2 CERT report structure

The institution's Cyber Defense Liaison files CERT reports in the structure below.

| Field | Content |
|---|---|
| Anomaly type | Fingerprint mismatch / seal-age delay / verifier failure / OTLP tampering / supply-chain mismatch |
| First-detection timestamp | RFC 3339 UTC |
| Scope | Number of entries affected, dates, tenants, regions |
| Preliminary root-cause hypothesis | Plain-text explanation |
| Interim remediation taken | Isolation, key rotation, vendor escalation |
| IOCs if available | Unauthorised IKM access pattern, HSM audit log anomaly, network traffic signature |

### 6.3 CERT threat intelligence consumption

The INCD CERT may respond with intelligence about coordinated attacks across multiple Israeli banks, threat-actor TTPs (per MITRE ATT&CK), or recommended mitigations. The institution incorporates CERT guidance into its incident-response procedure and documents the incorporation. The Chief Information Security Officer reports quarterly to the Bank of Israel on CERT notifications received and actions taken.

---

## 7. Directive 414 + Directive 357 — vendor audit and data-controller posture

### 7.1 Annual vendor audit (Directive 414)

Directive 414 (סיכון צד שלישי) requires Israeli banks to audit any critical third party — including a vendor providing chain-of-custody service — annually. The audit scope includes control effectiveness, security posture, compliance with Bank of Israel requirements, and incident-response integration.

| Audit area | Scope |
|---|---|
| Operational control design | HSM custody, ledger append-only enforcement, seal-job automated execution, incident-response procedures — verified against spec §4.1-§4.3 |
| Cryptographic control testing | Vendor provides 30 days of seal records; bank re-computes Merkle root and verifies signature against vendor's public key; 100% pass rate required |
| Change-management and access control | Vendor's change-control log for seal-job logic, HSM firmware, ledger code; vendor audit log showing which vendor employees accessed HSM in past 12 months |
| Incident-response integration | Vendor's incident-response procedures for the chain; SLA verification for incident notification timing |
| Supply-chain and build reproducibility | Vendor provides SHA-256 hashes of released SDKs and ledger binaries signed with release key; bank rebuilds locally and verifies match |
| Third-party subcontractors | If vendor uses third-party HSM providers (AWS CloudHSM, Azure Managed HSM), audit subcontractor's controls independently |

Findings are reported to the bank's Board Audit Committee and the Bank of Israel within 30 days; remediation plans are negotiated with the vendor and tracked to closure.

### 7.2 Data-controller posture (Directive 357)

Directive 357 §3 requires that data-processor contracts name specific audit and breach-notification obligations. The IKM is the most sensitive piece of data in the chain; a foreign vendor holding the IKM (even in an HSM) means a foreign party has influence over the bank's ability to re-sign the chain.

| Contract clause | Directive 357 §3 requirement |
|---|---|
| Data localisation | Vendor commits to storing all chain data exclusively in Israeli jurisdiction; no replication to non-Israeli sites without bank approval |
| Audit rights | Bank can audit vendor's chain-of-custody operations on demand (at least quarterly) |
| Subcontractor management | Vendor names HSM provider and confirms HSM is in Israeli jurisdiction; vendor takes responsibility for subcontractor compliance |
| Incident notification | Vendor commits to 4-hour notification of incidents; INCD notification within 1 hour for suspected cyber-attack incidents |
| Data retention and deletion | Vendor retains data for institution's specified retention window; upon contract termination, vendor delivers data and destroys replicas with cryptographic proof of destruction |
| Regulatory cooperation | Vendor commits to cooperating with Bank of Israel examinations and INCD investigations |
| Background checks | Vendor confirms personnel with HSM, ledger, IKM access have passed Israeli security vetting |
| SLAs | Seal-publishing SLA (≤60 minutes from seal window end); >2 SLA breaches per quarter triggers re-negotiation or vendor replacement |

---

## 8. Bank of Israel — Federal Reserve coordination for cross-border banks

### 8.1 The MoU framework

The Bank of Israel and the Federal Reserve maintain a Memorandum of Understanding on cross-border banking supervision. Israeli megabanks holding USD correspondent accounts (all five megabanks do) are examined by both regulators on overlapping but distinct timelines. The MoU provides for coordinated examinations at the institution's request and for examiner-to-examiner notification of significant findings.

### 8.2 Coordinated evidence production

To minimise duplication, the institution maintains one audit trail that satisfies both frameworks.

| Evidence | FFIEC examination | Bank of Israel examination |
|---|---|---|
| Daily verifier output | Examiner consumes verifier output per FFIEC P-1 procedure | Supervisor consumes verifier output per Bank of Israel parallel procedure (see §11) |
| Weekly anomaly review | Procedure P-5/P-6 | Composes with Directive 361 §5 INCD coordination evidence |
| Annual vendor audit | TPRM evidence | Directive 414 evidence |
| Operational-resilience testing | FFIEC operational-risk program | Directive 365 §4 drill (see §3.3 of this overlay) |

### 8.3 Examiner-to-examiner notification

When the Federal Reserve issues a finding on the chain (e.g., MRA on fingerprint-reconciliation gaps), the institution notifies the Bank of Israel that the FFIEC examination found a defect; the Bank of Israel may incorporate the finding into its Directive-compliance assessment. Conversely, if the Bank of Israel identifies a Directive 365 recovery-time deficiency, the institution notifies the Federal Reserve that a compensating control may be needed to meet the FFIEC examination standard.

### 8.4 Conflict resolution

If a requirement under Bank of Israel directives is stricter than the FFIEC spec (e.g., Directive 365's 2-hour recovery vs the spec's daily seal cadence), the institution satisfies the stricter requirement. Bank of Israel directives take precedence for Israeli-supervised activities.

---

## 9. Privacy Protection Authority (PPA) coordination

### 9.1 PPA jurisdiction

The Privacy Protection Authority enforces the Privacy Protection Law including Amendment 13. Amendment 13 specifically addresses AI applications and treats AI inference logs and AI-driven decisions as sensitive personal information.

### 9.2 PPA enforcement triggers

Triggers for PPA audit of an institution's chain deployment:

- Customer complaints about automated decision-making (PPL §3 right to human review) that institutions fail to fulfill within 30 days.
- Data-subject requests for access or erasure that institutions fail to fulfill within 30 days.
- Incident notifications of privacy-store breaches (PPL §7).
- Random audits — PPA conducts periodic examinations of critical-sector organisations including banks.

### 9.3 PPA examination preparation

If the PPA initiates an audit, the institution produces:

- Privacy policy and notices (customer-facing).
- Data Processing Agreement (DPA) with any vendor.
- Breach notification logs for the past 3 years.
- Access logs for the privacy-store.
- Incident reports and remediation plans.
- Documentation of data-subject-rights fulfillment (sample of 10 recent erasure requests with confirmation of completion).

### 9.4 PPA enforcement posture

The PPA may issue directives to cease or modify processing if it determines processing violates PPL §1-§7. Failure to comply is an administrative fine of up to 5 percent of global turnover under PPL §12. For a Tier-1 megabank this can be tens of millions of shekels.

---

## 10. Nation-state threat model — INCD-aligned articulation

### 10.1 The INCD baseline

The Israeli National Cyber Directorate threat model assumes a capable nation-state adversary (Iran IRGC cyber, Lazarus-equivalent operators) with 18 months of persistence inside the bank's infrastructure. This is not a theoretical risk; it is the baseline posture for Israeli megabank cybersecurity.

### 10.2 Adversary capabilities

The adversary can:

- Compromise the SDK supply chain (inject backdoor in released binaries; poison npm/pip/Maven dependencies).
- Compromise the OTLP collector (tamper with in-flight chain entries; inject forged entries).
- Compromise the HSM firmware (modify signing logic to sign arbitrary roots without verification; leak signing key after N uses).
- Exfiltrate the IKM if the HSM is network-connected (decrypt past session keys, forge historical entries with old keys).
- Simultaneously compromise detection mechanisms (disable audit logging, suppress anomaly alerts).
- Persist through backup/disaster-recovery infrastructure if not independently secured.

### 10.3 Chain defense under the INCD scenario

The chain's three-layer composition provides sequential defense.

| Layer | Compromise scenario | Detection mechanism |
|---|---|---|
| 1. SDK / per-event MAC | Adversary forges session key, generates events | Verifier detects forged events via HMAC chain mismatch |
| 2. Ledger / Merkle seal | Adversary rewrites events, recomputes root | Seal-job recomputed root mismatches the previously signed root from prior days |
| 3. HSM / Ed25519 signature | Adversary compromises HSM firmware, signs arbitrary roots | Detection requires firmware-attestation or out-of-band signing-pattern monitoring; this is residual risk under the INCD baseline |

### 10.4 Detection-speed and recovery asymmetry

The daily seal cadence detects HSM-layer compromise within 24 hours. Recovery (IKM rotation, HSM firmware restoration, forensic analysis) takes 24-48 hours for a prepared institution. This asymmetry is acceptable under Directive 365 (2-hour recovery requirement) only if the institution uses hourly seals or active-active failover (see §3.2).

### 10.5 Institutional mitigations beyond the spec

Against the INCD scenario, the institution implements:

- Verifier supply-chain discipline (rebuild verifier locally, confirm SHA-256 match; see §12).
- Vendor-compromise detection (build reproducibility, firmware attestation, SBOM audits; see §12).
- Anomaly detection tuned for 18-month-persistence timelines (daily anomaly review with nation-state behavioral signatures).
- INCD coordination and red-team exercises (annual INCD red-team scenario; see §13).

---

## 11. Bank of Israel examination procedures

### 11.1 Examination objectives

During a routine examination of an Israeli bank deploying chain-of-custody, the Bank of Israel Banking Supervisor evaluates:

| Objective | Evidence |
|---|---|
| Directive 365 compliance | Operational-resilience drill results (past 3 years); recovery time ≤ 2 hours |
| Directive 367 compliance | Chain data, ledger storage, HSM, backups all in Israeli jurisdiction or in INCD-approved non-Israeli jurisdiction |
| Directive 361 compliance | Documented INCD coordination procedures; incident-response drill results or real-incident records |
| Directive 414 compliance | Annual vendor audit conducted; findings remediated or documented remediation plan |
| Directive 357 compliance | Data-processor contract with specific custody, audit, and breach-notification obligations |

### 11.2 Examination work-papers

The supervisor requests:

- Daily seal records for the examination period (12 months typical); spot-checks 5 dates for signed seal with timestamp within 60 minutes of UTC midnight.
- Operational-resilience drill results (past 3 years); confirms all drills show recovery time ≤ 2 hours.
- Incident logs for the examination period; for any chain-detected incidents, confirms classification per Directive 361 severity taxonomy and notification per timeline.
- Vendor audit report (if applicable); confirms annual vendor audit was conducted and findings remediated.
- IKM rotation log; confirms rotations are documented and fingerprint registries updated synchronously.
- Hebrew-language disclosure evidence (if lending or HR AI is in use); spot-checks 5 customers who received adverse decisions and confirms disclosure was provided in Hebrew within 24 hours.

### 11.3 Test procedures (parallel to FFIEC P-1 through P-10)

| Procedure | Activity |
|---|---|
| Verifier procedure (parallel to P-1) | Supervisor downloads reference verifier, runs on sample of 100 chain entries from random dates; confirms 100% pass rate |
| Seal integrity (parallel to P-8) | Sample 12 seal records (one per month); for each, recompute Merkle root and verify signature against institution's public key |
| Anomaly review (parallel to P-5, P-6) | Request 12 months of weekly anomaly-review meeting minutes; spot-check 4 months for IT operations / IT risk / audit attendance, anomaly identification, and CRO/CISO escalation |
| Directive 361 §5 incident-response | Review institution's incident-response log; confirm INCD timeline (Red within 1 hour, Orange within 4 hours), Bank of Israel timeline (24 hours for Yellow), and INCD coordination on forensics and remediation |
| Directive 414 vendor audit | Review past 3 years of annual vendor audits |

### 11.4 Findings and MRA thresholds

| Finding | Trigger |
|---|---|
| **Control operating effectively (PASS)** | All sampled entries pass verifier; all seal records verify; anomaly review documented monthly; operational-resilience drill ≤ 2 hours; no Directive 365/361/367/414 deficiencies |
| **MRIA (Matters Requiring Immediate Attention)** | Recovery-time drill fails 2-hour target; chain-detected incident not escalated to INCD per Directive 361 timeline; vendor audit not conducted or significant findings not remediated. Immediate action plan and corrective action within 30 days |
| **MRA (Matters Requiring Attention)** | Anomaly review not documented monthly; vendor audit findings exist but under remediation within agreed timeline; IKM rotation documented but fingerprint registry out of sync for a period (root cause identified and remediated). Corrective action plan and remediation within 90 days |

---

## 12. Supply-chain discipline and verifier escrow

### 12.1 SDK reproducible-build commitment

Israeli megabanks operating under the INCD threat model verify the SDK and verifier supply chains explicitly. Herald.Py and the reference Go verifier publish SHA-256 hashes of all released packages with GPG signatures from the release-signing key; the institution rebuilds locally and verifies match.

### 12.2 Verifier escrow for tier-1 deployments

The institution maintains a "golden copy" of the verifier binary in escrow, independent from the vendor-supplied release. Three escrow patterns:

- **Physical escrow:** CD-ROM or USB drive, sealed envelope, stored in the bank's vault with access controls (CISO + Audit chair approval to remove).
- **Cryptographic escrow:** Verifier binary and SHA-256 hash signed by the institution's CISO using an Ed25519 key held in the institution's own HSM (separate from the chain's signing HSM).
- **Digital escrow:** Verifier binary stored on an air-gapped computer held by the bank's internal audit department, powered off except during examination use.

### 12.3 Examination-time joint verification

When the Bank of Israel supervisor or Federal Reserve examiner arrives for examination, the examiner and the institution jointly:

- Retrieve the escrowed verifier binary from vault or air-gapped storage.
- Verify the binary's SHA-256 against the CISO's signed manifest.
- Verify the SHA-256 matches the published release SHA-256 (the examiner brings a pre-computed list of trusted SHA-256 values).
- Run the verifier on the ledger export.

This procedure ensures the examiner is using a verifier that the institution and the examiner jointly trust, not a binary that could have been substituted between release and examination.

### 12.4 OTLP transport hardening

Spec §5.1 mandates TLS 1.3 for OTLP encryption. For the INCD threat model, the institution adopts:

- TLS 1.3 strict enforcement (TLS 1.2 disabled at the receiver endpoint).
- Certificate pinning (HPKP or pin-on-first-use) for tier-1 deployments.
- Mutual TLS with client-certificate verification per institution and per tenant.
- Network segmentation: OTLP ingress firewalled to accept connections only from named SDK sources.
- SIEM monitoring for TLS anomalies (cipher-suite downgrades, certificate-chain mismatches, unexpected client-certificate identifiers).

---

## 13. Red-team exercise readiness

### 13.1 The INCD red-team program

The INCD conducts annual red-team exercises (cyber defense drills) for tier-1 banks. The exercise scope may include the chain-of-custody system; the institution should be prepared for an exercise scenario testing the chain's resilience under simultaneous SDK + OTLP + ledger + HSM compromise.

### 13.2 Exercise scenarios

| Scenario | Red-team objective | Expected chain outcome |
|---|---|---|
| SDK compromise | Forge a chain event going forward | Verifier detects forged events via HMAC chain mismatch when re-verifying; red team fails |
| OTLP tampering | Inject, modify, or replay OTLP messages | Ledger receiver re-verifies HMAC chain on ingest (per spec §7 step 5-6); forged messages rejected; red team fails |
| Ledger-server compromise | Rewrite events and recompute Merkle root | Next daily seal-verification recomputes root from modified events; root does not match the prior signed root; tampering detected; red team fails |
| HSM firmware compromise | Sign forged root under bank's key | Residual risk under spec; institution's compensating control (firmware attestation if available) detects firmware version mismatch; otherwise out-of-band detection via network anomalies, unauthorized SSH sessions, physical security alerts |
| Simultaneous compromise (all four) | Forge a history that appears valid | Only succeeds if all three trust boundaries (application, ledger, HSM) are simultaneously compromised. If any one component is defended, verifier detects the forgery or real-time detection alerts during the exercise |

### 13.3 Pre-exercise institutional readiness

Before the INCD red-team exercise, the institution confirms:

- Verifier binary is held securely (air-gapped, independently rebuilt; per §12.2).
- Public key is registered with the institution, the Bank of Israel, and the Federal Reserve.
- Daily verification run (sample-checking entries against the public key) is operating to detect forged signatures.
- Incident-response procedures are live and testable.

### 13.4 Post-exercise reporting

The red team documents which attacks succeeded and why; which failed and why; time-to-detection per scenario; and remediation recommendations. The CISO maps findings to control improvements; if SDK compromise detection was slow, strengthen process memory protection; if OTLP tampering succeeded due to TLS downgrade, enforce TLS 1.3 strictly. The exercise report is filed with the Bank of Israel as evidence of cyber-resilience under the nation-state threat model.

---

## 14. Incident-response procedures (Directive 411 + Directive 361 §5)

### 14.1 Two-phase response procedure

Phase 1 (0-2 hours): Preliminary determination. The CISO and IR lead triage the incident based on available evidence. Decision point: does this breach likely affect >1,000 customers or involve nation-state compromise indicators? If yes or unclear, trigger INCD notification and Bank of Israel notification per the Directive 361 §5 severity taxonomy.

Phase 2 (2-24 hours): Full investigation. The IR team preserves evidence, performs forensics, and engages the chain-of-custody system. If the incident involves potential chain tampering, the IR team isolates the affected ledger-server or HSM, runs the verifier on ledger exports from the suspected compromise window, and documents results.

### 14.2 Incident severity classification

The IR runbook names the determination procedure for each severity tier per §5.2. The decision is documented with timestamp, classifying officer signature, and the evidence basis for the classification.

### 14.3 Incident scenarios specific to chain operation

| Scenario | Severity | INCD timing | Bank of Israel timing |
|---|---|---|---|
| Master-key (IKM) compromise suspected | Red | 1 hour | 1 hour |
| HSM firmware tampering detected | Red | 1 hour | 1 hour |
| Ledger tampering detected by verifier | Red | 1 hour | 1 hour |
| Vendor breach affecting HSM or ledger | Orange | 4 hours | 4 hours |
| Supply-chain compromise (verifier SHA-256 mismatch) | Orange | 4 hours | 4 hours |
| OTLP collector tampering suspected | Orange | 4 hours | 4 hours |
| HSM unavailable (maintenance, hardware failure) | Yellow | Optional | 24 hours |
| Ledger storage corruption | Yellow | Optional | 24 hours |
| Key rotation procedural error (rotation logged but registry not updated) | Yellow | Optional | 24 hours |
| Single verifier failure traced to known transaction-ordering issue | Green | None | Internal log |
| Seal-age delay <75 minutes | Green | None | Internal log |
| Anomalous but explained event (clock skew on SDK host, NTP corrected) | Green | None | Internal log |

### 14.4 Key recovery ceremony (Scenario 4b — recovering from suspected compromise)

The institution's CC8.1 names the key recovery ceremony procedure including:

- Participants: CISO (chair), Chief of Operations (technical execution), DBA Lead (database consistency), Internal Audit Observer (witness), Bank Executive (CFO or CRO; ultimate approver), External Auditor (optional, recommended for tier-1 banks).
- Pre-ceremony: 48 hours notice; backup IKM prepared in physical escrow (sealed envelope held by two escrow agents); ledger backup prepared.
- Ceremony execution: quorum verification, backup opening with seal-intact verification, HSM provisioning, fingerprint verification, test signing, activation, log documentation.
- Post-ceremony: Operations chief report to CISO within 24 hours; Internal Audit ceremony report; INCD notification if the recovery was due to compromise.

---

## 15. Translation table — Bank of Israel directives to chain artifacts

| Directive | Subject | Chain artifact |
|---|---|---|
| 357 §3 (data-processor contracts) | Data localisation, audit rights, sub-contractor management | DORA-overlay §2 lattice (institution + SDK + receiver + HSM + LLM); contractual flow-down per §7.2 of this overlay |
| 359 (cyber defense management) | INCD Banking Cyber CERT participation | §6 of this overlay; spec §10.2 operational events |
| 361 §5 (cyber risk management) | INCD coordination on incidents | §5 of this overlay; severity taxonomy in §5.2 |
| 365 §3 (operational resilience) | ≤2-hour recovery for critical operations | §3 of this overlay; hourly cadence or Pattern A active-active |
| 365 §4 (operational-resilience drills) | Annual drill | Audit procedure P-12 in §3.3 of this overlay |
| 367 §2 (cloud and GenAI localisation) | Israeli-jurisdiction storage | §4 of this overlay; vendor contract requirements in §4.3 |
| 411 (incident reporting) | Bank of Israel notification | §5 + §14 of this overlay |
| 414 (third-party risk) | Annual vendor audit | §7.1 of this overlay |
| PPL Amendment 13 (sensitive personal information) | Cross-border transfer of AI inference logs | §4 + §9 of this overlay |
| Cyber Defense Law 5778-2018 | INCD notification of cyber incidents | §5 + §6 of this overlay |
| Equal Opportunity in Employment Law (1988, amended 2022) | Hebrew-language adverse-action disclosure within 24 hours | `audit.adverse_action.regulatory_basis = il-eea-1988-2022`; institution's disclosure-generation tooling pulls from chain entry; see §16 of this overlay |

---

## 16. Equal Opportunity Employment Law adverse-action workflow

### 16.1 What the law requires

Israel's Equal Opportunity in Employment Law (1988, amended 2022 for automated decisions) requires that any adverse employment action or credit decision made (wholly or partly) by automated means be disclosed to the affected person in Hebrew within 24 hours, along with the right to human review.

### 16.2 Chain-driven disclosure workflow

| Step | Activity | Chain attribute |
|---|---|---|
| 1. AI agent makes adverse decision | Decision captured per spec §4.4 | `gen_ai.response.text`; `audit.routing.*` |
| 2. Disclosure system queries chain | Query by run_id, captured_at, customer ID | Standard query against ledger |
| 3. Extract decision evidence | Pull `gen_ai.request.messages`, `gen_ai.response.text`, model parameters | Standard chain attributes |
| 4. Construct Hebrew disclosure | Disclosure-generation tool produces plain-Hebrew notice naming model, inputs, decision, confidence, right to review | Institution's tooling |
| 5. Send disclosure within 24 hours | Notice sent to affected person | `audit.adverse_action.regulatory_basis = il-eea-1988-2022`; `audit.adverse_action.notice_language = he-IL`; `audit.adverse_action.delivery_timestamp` |
| 6. Test workflow quarterly | Select 10 past adverse decisions; reconstruct from chain; verify disclosure provided to customer matches chain record | Compliance test record |

### 16.3 Regulatory evidence

If the Privacy Protection Authority or the Labor Court requests evidence the bank disclosed the decision correctly, the bank produces the chain entry (verifier-signed, integrity-bearing) as proof the disclosure was accurate and timely. The chain's cryptographic binding ensures the bank cannot later claim the disclosure was different from what the model actually output.

---

## 17. Operational checklist for an Israeli bank deploying the chain

The institution working through Bank of Israel + INCD + PPA + Federal Reserve compliance produces the following artefacts. This checklist names the minimum set; the institution's compliance department extends it as their specific exposure requires.

| Artefact | Source | Reviewed by |
|---|---|---|
| Directive 365 operational-resilience drill report (annual) | Institution's IT operations + CC8.1 | Bank of Israel Banking Supervisor |
| Directive 367 deployment topology declaration (Israeli jurisdiction or PPA-approved non-Israeli) | Institution's CC8.1 | Bank of Israel + PPA |
| PPL Amendment 13 privacy-store custody documentation | Institution's CC8.1 + RoPA | PPA |
| Directive 361 §5 incident-response runbook with INCD coordination | Institution's IR runbook | INCD + Bank of Israel |
| Directive 359 INCD CERT integration procedure | Institution's CC8.1 | INCD |
| Directive 414 annual vendor audit | Institution's internal audit | Bank of Israel + Board Audit Committee |
| Directive 357 vendor contract with Israeli localisation | Institution's procurement | Bank of Israel |
| BoI-Federal Reserve coordinated evidence package | Institution's compliance | Bank of Israel + Federal Reserve |
| Hebrew-language adverse-action disclosure tooling + chain integration | Institution's customer-service + CC8.1 | PPA + Labor Court (on demand) |
| Verifier escrow procedure | Institution's CC8.1 + Internal Audit | Internal Audit |
| Annual INCD red-team exercise readiness | Institution's CC8.1 | INCD |
| Key recovery ceremony procedure | Institution's CC8.1 | Internal Audit + CISO |
| `audit.adverse_action.regulatory_basis` enumeration including `il-eea-1988-2022` | Institution's CC8.1 | Institution's consumer-protection counsel |
| Pattern A or hourly-cadence selection for Directive 365 compliance | Institution's CC8.1 | Bank of Israel |

---

## 18. Bottom line

The chain v1.0a is conformant with the Bank of Israel supervisory architecture and operationally defensible under the INCD threat model when the institution adopts the compensating controls articulated in this overlay. The conformance is operational: the institution's CC8.1 names the seal cadence (hourly or daily-with-Pattern-A), the deployment topology (Israeli jurisdiction or PPA-approved non-Israeli with documented safeguards), the multi-regulator IR procedure (INCD + Bank of Israel + Federal Reserve + PPA), the annual vendor-audit posture, and the Hebrew-language disclosure workflow.

The chain's integrity primitives are NIST/IETF-standardised and Israel-conformance-friendly: HMAC-SHA-256, HKDF-SHA-256, Ed25519, RFC 8785, RFC 6962 Merkle, FIPS 140-2 Level 3 HSM custody. None of these primitives require modification for Israeli deployment; what differs is the institutional CC8.1 documentation, the per-jurisdiction tenant naming (with `bank-il` for Israel-jurisdiction tenants and HSM custody pinned to AWS Tel Aviv or Azure Israel region or Israeli on-premises HSM), and the multi-regulator coordination procedure.

For Israeli megabanks under nation-state threat baseline, the chain produces detection and proof. The chain does not prevent nation-state compromise; it bounds tampering to single-day windows (default daily seal) or single-hour windows (hourly seal) and makes tampering detectable and non-repudiable. The institutional discipline articulated in §10 through §13 produces the full operational and cyber-resilience posture under Bank of Israel and INCD oversight.

The Bank of Israel supervisor consumes the verifier output and the institution's CC8.1 jointly during examination; the FFIEC examiner (where applicable for cross-border banks) consumes the same evidence under the Bank of Israel-Federal Reserve MoU coordinated procedure. The institution's single audit trail satisfies both frameworks, with the institutional CC8.1 documenting the per-regulator framing without producing duplicate evidence.

For institutions evaluating the chain for the first time under Israeli regulatory mandate, the recommended posture is: self-host the ledger and verifier in an Israeli data center or colocation, adopt hourly seals or Pattern A active-active for Directive 365 compliance, integrate with INCD coordination per §5 and §6, conduct annual vendor audits per §7, prepare for INCD red-team exercises per §13, and document the Hebrew-language adverse-action disclosure workflow per §16. With these compensating controls, the chain materially strengthens the bank's cyber-resilience and AI-governance posture under Bank of Israel supervision.

---

## Document control

| Field | Value |
|---|---|
| Document | Bank of Israel Articulation Overlay |
| Version | 1.0.0 |
| Status | Informative — institution-side articulation; no normative spec change |
| Aligned with | Bank of Israel Directives 357, 359, 361, 365, 367, 411, 414; PPL Amendment 13; Cyber Defense Law 5778-2018; Equal Opportunity in Employment Law 1988 (amended 2022); INCD Defense Methodology v2.0; BoI-Federal Reserve MoU |
| Date | 2026-05-07 |
