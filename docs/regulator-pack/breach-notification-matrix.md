# Breach Notification Matrix — Multi-Jurisdiction Timing

> **What this doc is.** The institution's binding decision matrix for when a chain-detected event triggers a breach-notification obligation, which clocks start at which moment, and how multi-jurisdiction deadlines are honored simultaneously. Written so the IR commander, privacy office, and legal team work from one shared timeline rather than three.

> **Why this matters.** Every privacy regime has a different clock. HIPAA gives 60 days. FFIEC gives 36 hours. GDPR gives 72 hours. DORA (EU) gives 24 hours. State data-breach laws vary. If the IR team manages each clock separately, the institution will miss the tightest deadline. This document collapses the four clocks onto a single decision tree.

---

## 1. The four primary clocks

| Regulation | Citation | Clock starts | Deadline | Recipient |
|---|---|---|---|---|
| **HIPAA Breach Notification Rule** | 45 CFR 164.400–414 | Discovery of unauthorized access, acquisition, use, or disclosure of unsecured PHI | **60 days** | Affected individuals; HHS; media (if breach affects 500+ individuals in a state or jurisdiction) |
| **FFIEC computer-security incident** | 12 CFR 53 (OCC), 12 CFR 225 (Fed), 12 CFR 304 (FDIC) | Determination that a "computer-security incident" has occurred and rises to "notification incident" thresholds | **36 hours** | Primary federal banking regulator |
| **GDPR Article 33** | GDPR Art. 33 (regulator) and Art. 34 (data subjects) | Awareness of a personal-data breach | **72 hours** to supervisory authority; without undue delay (typically within 60 days) to data subjects if high risk | EU Data Protection Authority; affected data subjects |
| **DORA (EU financial entities only)** | Regulation (EU) 2022/2554 Art. 19 | Classification of an ICT-related incident as "major" | **24 hours** initial notification; 72 hours intermediate; 1 month final | EU competent authority |

State data-breach statutes (California §1798.82, New York GBL §899-aa, etc.) add further clocks but generally run from "discovery" with deadlines between "without unreasonable delay" and 60 days. State clocks are coordinated through the institution's privacy team rather than the IR commander; they are not the binding floor here.

---

## 2. Per-event-class clock-start guidance

Different chain-detected events have different relationships to personal data. The institution determines per-event whether and which clocks start.

### 2.1 IKM (HMAC integrity-key) compromise

The IKM by itself is not personal data. An IKM compromise enables forgery of past entries but does not by itself disclose customer data. Whether a personal-data breach has occurred depends on what else was accessed.

| Sub-scenario | FFIEC clock | GDPR clock | HIPAA clock |
|---|---|---|---|
| IKM leaked, no concurrent access to chain entries or privacy-store | Starts at determination (computer-security incident) | Does not start (no personal data exposed) | Does not start (no PHI accessed) |
| IKM leaked + chain-entry access detected | Starts at determination | Starts at determination of personal-data breach (chain entries contain pseudonymized personal data — still GDPR-personal-data per `pseudonymization-vs-anonymization.md`) | Starts only if PHI was in the chain entries (healthcare deployments) |
| IKM leaked + privacy-store access detected | Starts at determination | Starts immediately — token-to-PII mapping plus chain entries together yield personal data | Starts if PHI was in the privacy-store |

### 2.2 Token-vault (privacy-store) compromise

The privacy-store contains the token-to-PII mapping. A compromise here is **always** a personal-data exposure — the tokens in the chain become reversible to the original PII.

| Sub-scenario | FFIEC clock | GDPR clock | HIPAA clock |
|---|---|---|---|
| Privacy-store accessed (read or copied) | Does not start by itself (FFIEC's computer-security-incident definition focuses on integrity, availability, and confidentiality of the institution's information systems; a vault breach qualifies if it meets the materiality threshold for "notification incident") | Starts immediately at awareness (Article 33: 72 hours to DPA; Article 34: notify data subjects without undue delay if high risk) | Starts at discovery if the vault contained PHI (60-day clock to individuals, HHS) |
| Privacy-store deleted (availability incident) | Starts at determination (availability impact on a critical system) | Starts only if the deletion deprived the institution of the ability to honor data-subject rights; consult the DPA | Starts only if PHI availability was affected for a covered entity's purposes |

### 2.3 Combined breach (IKM + privacy-store + chain entries)

The worst case. All three clocks start simultaneously and the GDPR 72-hour clock generally takes precedence as the tightest applicable to personal-data exposure.

### 2.4 Merkle integrity failure detected

The verifier reports a Merkle-root mismatch or a per-entry HMAC failure. Whether personal data was exposed depends on the cause.

| Sub-scenario | FFIEC clock | GDPR clock | HIPAA clock |
|---|---|---|---|
| Integrity failure caused by storage corruption (not adversarial) | Starts at determination (availability/integrity event) | Does not start (no personal data exposed) | Does not start |
| Integrity failure caused by tampering (entry modification or insertion) | Starts at determination | Starts only if the tampering was associated with adversarial access to chain entries (read access in addition to write) | Starts only if PHI access was associated with the tampering |

Tampering by an insider who already had legitimate read access to entries is a personal-data-disclosure scenario (the insider read the data). Tampering by an outsider who exploited a write path without read access is a different scenario; consult the IR forensics findings and the privacy office before starting GDPR/HIPAA clocks.

---

## 3. Multi-jurisdiction simultaneous-notification rule

When more than one clock applies, the institution operates **simultaneous notification** rather than sequential. The single incident-determination event triggers all applicable clocks at the same wall-clock moment. The IR commander, privacy office, and legal team consolidate the evidence package and dispatch notifications to each recipient on each recipient's timeline.

The decision tree:

1. **Incident determined.** IR commander records the determination timestamp `T0`.
2. **Personal-data assessment.** Privacy office determines (within hours) whether personal data was exposed, accessed, or affected. The assessment is recorded in writing.
3. **Clock activation.** Based on §2 above, the IR commander activates each applicable clock with `T0` as the start.
4. **Notification dispatch.** Notifications go out in **deadline order, tightest first**:
   1. DORA 24-hour (if applicable — EU financial entity with major ICT incident)
   2. FFIEC 36-hour (if applicable — federal banking regulator notification incident)
   3. GDPR Article 33 72-hour (if applicable — DPA notification)
   4. HIPAA 60-day (if applicable — affected individuals, HHS, media if 500+)
   5. State data-breach laws (if applicable — varies; coordinated by privacy office)
   6. GDPR Article 34 (notify data subjects) — typically follows DPA notification within days; high-risk threshold per Article 34(1)
5. **Consistent narrative.** All notifications use the same incident description, scope, and remediation summary. Inconsistencies between regulator notifications and customer notifications are themselves an enforcement risk.
6. **Coordination across regulators.** Where one recipient is competent for multiple regimes (e.g., a national authority that is both DPA and DORA competent), a single notification can satisfy multiple obligations if explicitly stated.

The institution's IR playbook (`incident-response-playbook.md`) carries this decision tree as Scenario 11 (multi-jurisdiction breach notification) and updates it annually.

---

## 4. The "discovery" vs "determination" distinction

Each regulation uses slightly different language for the clock-start moment. The institution operates on **earliest applicable trigger**:

| Regulation | Clock-start language | Operational meaning |
|---|---|---|
| HIPAA | "Discovery" — the day the breach is known or, by exercising reasonable diligence, would have been known | The day the IR team becomes aware of unauthorized PHI access |
| FFIEC | "Determination" — the institution has reason to believe a notification incident has occurred | The day the IR commander records the determination based on forensic evidence |
| GDPR | "Awareness" — when the controller becomes aware that a personal-data breach has occurred | The day the IR team identifies the breach scope (typically the same day as determination) |
| DORA | "Classification" — when the incident is classified as major | The day the institution applies the major-incident criteria |

In practice, "discovery," "determination," and "awareness" usually fall on the same day, with "classification" following within hours. The institution treats the earliest of these as `T0` and runs all clocks from that point. Where the earliest trigger is ambiguous, the institution operates the clocks on the earliest plausible day and documents the reasoning.

---

## 5. Notification-content checklist

Each notification carries the same core narrative, adapted to the recipient's required content. The checklist:

| Element | HIPAA | FFIEC | GDPR | DORA |
|---|---|---|---|---|
| Incident description | Required | Required | Required | Required |
| Categories of personal data / records affected | Required | Required | Required | Required |
| Approximate number of data subjects / individuals | Required (if known) | Encouraged | Required | Required |
| Likely consequences | Encouraged | Encouraged | Required | Required |
| Measures taken or proposed | Required | Required | Required | Required |
| DPO or contact point | n/a | Designated incident contact | Required | Required |
| Steps individuals can take | Required | n/a | Required (Article 34) | n/a |

The institution's IR templates (one per recipient) carry the matching content fields and are pre-approved by privacy office and legal.

---

## 6. Cross-references

- `incident-response-playbook.md` Scenario 11 — operational walkthrough of the simultaneous-notification rule.
- `dpia-template.md` §3 — risk analysis incorporates breach-notification timing as a residual-risk factor.
- `multi-jurisdiction-conflict.md` — when laws conflict (e.g., US compels disclosure GDPR forbids), this document handles the timing and the conflict document handles the lawfulness.
- `hipaa-security-rule-mapping.md` §164.308(a)(1) — security management process includes incident response with the timing rules above.
- `article-32-security-mapping.md` §2.4 — resilience controls include the IR playbook this matrix references.
- `pseudonymization-vs-anonymization.md` — explains why chain tokens are still personal data during the retention period and thus a chain-entry compromise is a GDPR breach.

---

## 7. Review cadence

This matrix is reviewed annually by the privacy office, IR commander, and legal team. Triggers for early review:

- A regulator updates timing or scope (e.g., FFIEC tightens to 24 hours; HIPAA modifies the unsecured-PHI definition).
- The institution enters a new jurisdiction whose clock is tighter than the existing matrix.
- An actual incident reveals a clock-start ambiguity not addressed here; the matrix is updated to remove the ambiguity.
- A new regulation creates a new clock (e.g., a federal cyber-incident law overlaying FFIEC).
