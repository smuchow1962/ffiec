# Multi-Jurisdiction Privacy-Law Conflict Resolution

> **What this doc is.** The institution's procedure for handling situations where two or more privacy regimes give conflicting directives — most commonly when a US legal process compels disclosure of data that GDPR or another regime forbids transferring or disclosing. Written so the legal team, privacy office, and IR commander can act on consistent principles rather than ad-hoc decisions during a high-pressure event.

> **Why this matters.** A bank operating in EU, US, Singapore, and Brazil is exposed to four different privacy regimes. When they disagree — and they do, regularly — the institution needs a documented decision tree that produces a defensible response under each regime simultaneously. Without the procedure, the institution either over-discloses (breaching GDPR or APAC privacy laws) or under-discloses (contempt of US court) or improvises in a way that creates evidentiary or reputational damage. The procedure is the institution's pre-decided answer.

---

## 1. Conflict scenarios at a glance

| Conflict | Frequency | Severity | Primary regimes |
|---|---|---|---|
| **EU-US disclosure conflict** | Common | High | GDPR Article 44/48 vs FISA, EO 12333, CLOUD Act |
| **US HIPAA vs EU GDPR healthcare data** | Occasional | High | HIPAA Privacy Rule (permissive disclosures) vs GDPR Article 44 (transfer prohibition) |
| **State-law data-localization vs federal supervisory access** | Occasional | Medium | State data-residency laws vs federal banking-regulator demand |
| **APAC localization vs cross-border ops** | Increasing | Medium | China PIPL, India DPDP, Singapore PDPA — varies |
| **Brazil LGPD vs US discovery** | Occasional | Medium | LGPD Article 33 transfer rules vs US litigation discovery |
| **Conflicting breach-notification-deadline directives** | Common | Low (timing only) | See `breach-notification-matrix.md` for resolution |

Each scenario has a default resolution path in this document. Specific incidents may require deviation; deviations are documented and approved by legal counsel and the DPO.

---

## 2. The EU-US disclosure conflict — most common

### 2.1 The conflict

A US legal process — federal subpoena, FISA §702 directive, CLOUD Act warrant, civil-litigation discovery order — demands disclosure of chain entries or privacy-store records that contain or reference personal data of EU data subjects. GDPR Article 48 prohibits transfer of personal data based solely on a third-country court or regulator decision that is not based on an international agreement. GDPR Article 44 prohibits transfer outside the EEA without a valid transfer mechanism.

The institution faces:
- Comply with US order → breach GDPR.
- Refuse US order → contempt or sanctions in US court.

### 2.2 The resolution sequence

| Step | Action | Owner |
|---|---|---|
| 1 | Legal team receives the US legal process; flags it as a multi-jurisdiction matter | Legal team |
| 2 | Privacy office and DPO are notified within 24 hours | Legal team |
| 3 | Scope assessment: identify the EU-data-subject component of the demanded data | Privacy office + chain operations |
| 4 | EU regulator engagement: legal team consults the lead supervisory authority (typically the institution's headquartered DPA or the relevant national DPA for the affected data subjects) for guidance | Legal team + DPO |
| 5 | US court engagement: legal team files a notice of the GDPR conflict in the issuing court, typically a motion to quash or modify, or a request for a protective order limiting the disclosure scope | Outside US counsel |
| 6 | Determine the response option (per §2.3 below) | Legal team + DPO |
| 7 | Execute the response; document every step | Legal team + privacy office |
| 8 | Post-event reporting to the EU DPA confirming the actions taken | DPO |

The sequence is paced; step 4 (EU regulator engagement) typically begins within days of receiving the order, not at the eleventh hour.

### 2.3 The four response options

| Option | When usable | Output |
|---|---|---|
| **A. Protective order** | US court is willing to limit disclosure to a tightly scoped subset, under seal, with foreign-data-subject anonymization | Court order narrowing the scope; institution discloses only the protected subset |
| **B. Local-law translation** | The data can be translated into a non-identifying form (e.g., aggregate statistics, synthetic data, anonymized chain entries post-de-mapping) that satisfies the US court's evidentiary need without disclosing personal data | Translated artifact disclosed; original personal data retained in EU |
| **C. Refuse and challenge** | The US order is overly broad and the institution accepts the risk of contempt to defend GDPR compliance | Motion to quash; if denied, contempt finding accepted; appellate challenge with amicus support from EU institutions |
| **D. Comply with EU regulator approval** | The EU DPA reviews the matter and provides written acknowledgment that the disclosure can proceed under specific conditions | Conditional disclosure; DPA's letter is part of the institution's record |

The institution's preferred sequence is A → B → D → C. Option C is the last resort and is taken only with full DPO and Board awareness.

### 2.4 The CLOUD Act and EU subsidiary route

The US CLOUD Act extends US legal process to data held by US providers regardless of storage location. For institutions using US cloud providers, the CLOUD Act is a recurring exposure. Mitigations:

| Mitigation | How |
|---|---|
| EU-held keys (Schrems II safeguard) | Per `international-transfers.md` Scenario B; US cloud provider stores ciphertext only |
| EU subsidiary as data controller | The EU subsidiary holds the data; the US parent does not. CLOUD Act reaches the US provider but not the subsidiary's data unless the parent has access |
| Data-localization architecture | EU data stays in EU region; US data stays in US region; no cross-region replication |

These mitigations reduce — but do not eliminate — the conflict. The institution maintains the resolution sequence regardless.

---

## 3. US HIPAA vs EU GDPR healthcare data

### 3.1 The conflict

HIPAA permits certain disclosures (e.g., treatment, payment, healthcare operations under 164.506; legal proceedings under 164.512(e); law enforcement under 164.512(f)) without individual authorization. GDPR Article 44 may forbid the corresponding cross-border transfer. The institution holds healthcare-related chain data for an EU patient and faces a HIPAA-permitted-disclosure scenario.

### 3.2 The resolution

The institution applies the **more restrictive standard**. HIPAA's permitted-disclosure rules are permissive — they allow disclosure but do not require it (with exceptions like required-by-law disclosures under 164.512(a)). GDPR's transfer rules are prohibitive — they forbid transfer absent a valid mechanism.

| Sub-scenario | Resolution |
|---|---|
| HIPAA permits disclosure to a US covered entity for treatment | If the data subject is an EU resident, the institution applies GDPR Article 44/49: disclosure may proceed only under a transfer mechanism (SCCs, BCRs, or derogation). For treatment, Article 49(1)(f) (vital interests of the data subject who is physically or legally incapable of giving consent) may apply for emergency treatment |
| HIPAA permits disclosure for payment or healthcare operations | The institution evaluates whether GDPR Article 44 transfer mechanism is in place; if not, the institution does not transfer |
| HIPAA requires disclosure (e.g., a court order under 164.512(e)) | The institution applies the EU-US disclosure-conflict procedure (§2 above) |
| HIPAA permits disclosure to law enforcement; the request is from US law enforcement | The institution treats this as a US-legal-process matter; §2 procedure applies |

The institution's legal team documents the conflict analysis for each disclosure decision. Audit findings of routine disclosures missing a documented analysis trigger remediation.

---

## 4. Conflict avoidance through storage topology

The most effective conflict-resolution strategy is conflict avoidance. The institution designs the storage topology so that the data subject's natural jurisdiction matches the storage location:

| Data subject jurisdiction | Storage location |
|---|---|
| EU data subject | EU region only; no replication outside EEA except under documented exception |
| US data subject | US region (default); EU region permitted under transfer mechanism |
| APAC data subject | APAC region(s) per applicable local-law residency requirements |
| Brazil data subject | Brazil region per LGPD residency analysis |

When data must be replicated across regions for resilience or operational reasons, the replication is encrypted with regional keys such that the receiving region cannot decrypt without explicit cooperation. This is the same EU-held-keys safeguard as `international-transfers.md` §2.2.

The privacy-store custody pattern in `token-vault-architecture.md` reinforces this: each region has its own privacy-store with its own key. A US legal process served on US infrastructure does not reach the EU privacy-store.

---

## 5. Legal hold across jurisdictions

When litigation, regulatory inquiry, or examination opens against records that span multiple jurisdictions, the legal hold itself can create conflicts.

| Sub-scenario | Resolution |
|---|---|
| US litigation opens; hold extends to EU records | Hold is documented; deletion is suspended; EU DPA is informed; the hold-driven retention is recorded as a §1798.105(d)(7)-equivalent retention exception under each applicable regime |
| EU regulatory inquiry opens; data exists in US backups | Hold extends to US backups; US-side custodians are notified; cross-border replication of the held data is suspended |
| Contradictory hold and deletion order | Legal counsel assesses; typically the hold prevails (preserving evidence) but the data subject is notified that retention is extended due to the hold |

The legal team maintains a unified hold register that names the affected data scope, the issuing authority, the effective dates, and the cross-jurisdictional implications.

---

## 6. Incident-response simultaneous-notification rule

When a breach is detected, the institution applies the simultaneous-notification rule from `breach-notification-matrix.md` §3. The rule extends to multi-jurisdiction conflicts:

| Conflict | Resolution |
|---|---|
| Two regulators' deadlines conflict | Notify the tightest deadline first; notify the next-tightest before its deadline; consistent narrative across all notifications |
| Two regulators want different notification content | Provide each its required content; do not omit fields one requires merely because another does not |
| One regulator demands public disclosure; another demands confidentiality | Coordinate with both; typically a joint regulator-mediated approach yields a narrowly scoped public notice plus a fuller confidential disclosure |
| Notification to a data subject in country A versus a regulator in country B | Both occur on their respective timelines; content is coordinated |

The IR commander is the operational owner of the simultaneous-notification rule; legal counsel and the DPO advise on content; privacy office maintains the notification log.

---

## 7. Documentation and after-action review

Every multi-jurisdiction-conflict event produces a written record:

| Record | Content |
|---|---|
| Conflict identification memo | Date the conflict was identified; the regimes involved; the data scope |
| Engagement record | Engagement with each affected regulator or court; dates, recipients, communications |
| Decision memo | The resolution option chosen; the reasoning; the approver (typically DPO + legal lead + Board if material) |
| Execution log | What was disclosed, what was withheld, what was anonymized, to whom, on what date |
| Post-event report | Outcome; lessons learned; any updates to this document |

The records are retained per the institution's records-retention schedule and are part of the file produced for any subsequent regulatory examination.

---

## 8. Pre-positioning for likely conflicts

The institution does not wait for a conflict to begin planning. Pre-positioning steps:

| Pre-position | Owner |
|---|---|
| Identify each jurisdiction the institution operates in; map data flows to jurisdictions | Privacy office |
| Maintain TIAs (Transfer Impact Assessments) for each non-adequate-country flow | Privacy office |
| Maintain SCCs, DPAs, and BCRs for each transfer relationship | Legal team |
| Maintain relationships with each relevant DPA and supervisory authority (named contacts; periodic touchpoints) | DPO |
| Maintain outside counsel relationships for each major jurisdiction (US, EU, UK, APAC) | Legal team |
| Maintain incident-response templates per jurisdiction | IR commander |
| Tabletop exercises annually covering EU-US disclosure conflict scenarios | IR commander + legal team |

Pre-positioning is the difference between a 14-day resolution and a 90-day crisis.

---

## 9. Cross-references

- GDPR Articles 44, 48, 49 — operationalized here.
- HIPAA 45 CFR 164.506, 164.512 — disclosure permissions.
- US CLOUD Act, FISA §702, EO 12333 — US legal process exposure.
- Brazil LGPD Article 33 — transfer rules.
- Various APAC regimes (China PIPL, India DPDP, Singapore PDPA) — regional residency.
- `international-transfers.md` — Schrems II compliance and TIA methodology; the conflict procedure here is the runtime counterpart.
- `breach-notification-matrix.md` — multi-jurisdiction notification timing.
- `token-vault-architecture.md` — regional privacy-store custody supports topology-based conflict avoidance.
- `litigation-support.md` — handles US-litigation discovery; this document handles the multi-jurisdiction overlay.
- `hipaa-minimum-necessary.md` — HIPAA scope analysis that interacts with cross-border healthcare data.
- `dpia-template.md` §3 — risk analysis incorporates multi-jurisdiction conflict as a residual risk.

---

## 10. Review cadence

This document is reviewed annually by the legal team, privacy office, DPO, and IR commander. Triggers for early review:

- A CJEU decision affecting transfer mechanisms or third-country surveillance analysis.
- US legislation or executive action affecting CLOUD Act, FISA, or related powers.
- A new APAC or LATAM regime entering force.
- An actual multi-jurisdiction-conflict event whose post-event report identifies a procedural gap.
- The institution adds or removes a deployment region.
- An EU adequacy decision is granted or revoked.
