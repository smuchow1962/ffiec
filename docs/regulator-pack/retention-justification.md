# Retention Justification — Why 7 Years for the Chain

> **What this doc is.** The institution's written justification for the 7-year retention period applied to chain-of-custody records and the linked privacy-store mapping. Written so a GDPR supervisory authority, an FFIEC examiner, or a CCPA-enforcement counsel can read it once and see why 7 years is the storage-limitation answer rather than a longer or shorter window.

> **Why this matters for the privacy file.** GDPR Article 5(1)(e) is not satisfied by saying "7 years feels right." It is satisfied by naming the legal purpose, naming the statute or regulation, and showing that no shorter period satisfies the purpose. This document is the named-statute pack.

---

## 1. The retention commitment in plain terms

The chain ledger and the privacy-store mapping that resolves chain tokens to original PII are both retained for **7 years from the date the chain entry was sealed**, unless a documented legal hold extends the retention period for a specific subset of records. At the end of 7 years (and once any active legal hold is released), the privacy-store mapping for those records **MUST be deleted**. The chain ledger itself is not modified — the tokens remain, but they are no longer resolvable to personal data and thus transition from pseudonymized to anonymized status (see `pseudonymization-vs-anonymization.md`).

This document explains why 7 years is the answer.

---

## 2. The statutes and rules grounding the 7-year period

| Source | Citation | What it requires | Retention floor |
|---|---|---|---|
| ECOA — Equal Credit Opportunity Act | **12 CFR 1002.12(b)(1)** | Creditors must preserve records of credit applications and adverse-action notices for 25 months after the action; for business credit, 12 months — but supervisory record-retention practice extends through the supervisory cycle | 25 months minimum, longer in practice |
| FFIEC examination records | OCC Comptroller's Handbook, Federal Reserve SR letters, FDIC RMS Examination Manual | Supervisory exam files, MRA/MRIA evidence, and prior-period working papers retained through the exam-cycle window plus an additional buffer | **7 years** is the harmonised supervisory floor across federal banking regulators |
| Federal civil-fraud statute of limitations | **31 U.S.C. § 3731(b)** (False Claims Act) | Civil action must be brought within 6 years of violation, or 3 years after the government knows the facts, but **not more than 10 years after the violation** | Effective ceiling 10 years; 7 years covers the common-case window |
| CFPB supervisory record-retention | 12 CFR Part 1070 (CFPB confidential supervisory information rules) | CFPB-supervised entities preserve records relevant to consumer-financial-protection examination | Aligns with FFIEC 7-year floor |
| State-law consumer protection (varies) | e.g., California CCPA recordkeeping under §1798.130(a)(5); New York DFS Part 500 §500.06 audit-trail | Records of consumer requests, security events, and supervisory data | 3 to 6 years typically; some states extend to 7 |

The 7-year period is the **smallest window that simultaneously satisfies the longest applicable supervisory record-retention floor (FFIEC) and covers the common-case window of the federal civil-fraud statute of limitations**. Anything shorter forces the institution to default to multiple per-record retention rules and explain why a missing record is not a record-retention violation. Seven years is the documented uniform answer.

---

## 3. Why 7 years is not over-retention under GDPR Article 5(1)(e)

Article 5(1)(e) requires that personal data be kept "in a form which permits identification of data subjects for no longer than is necessary for the purposes for which the personal data are processed." The institution's purposes for the chain are:

1. **Adverse-action defense.** ECOA permits a customer to challenge an adverse-action decision; subsequent litigation may run 3 to 4 years from the original action. The chain must support defense through final judgment plus appeal.
2. **Examination evidence.** FFIEC examiners conduct multi-year supervisory cycles. Prior-period working papers must be retrievable for cross-cycle analysis.
3. **Regulatory inquiry response.** CFPB, OCC, Fed, and FDIC may open inquiries that look back several years; the chain provides the integrity-bound evidence that what the AI decided is what the records show.
4. **Civil-fraud and False Claims Act defense.** The federal civil-fraud SoL is up to 10 years (with the 6-year primary window). The chain must remain available as defense evidence within the typical inquiry window.

Each purpose independently requires a multi-year retention. The 7-year period is the smallest window that covers all four purposes together. Going shorter (e.g., 4 years) would mean the institution loses defense evidence for older but still-actionable claims. Going longer (e.g., 10 years uniformly) would over-retain for purposes 1 and 2. The institution's DPIA documents this proportionality test (see `dpia-template.md` §3 — risk and proportionality).

---

## 4. State-law variations and the "longer of" rule

A handful of states extend retention obligations beyond the federal floor:

| Jurisdiction | Extension trigger | Effective retention |
|---|---|---|
| California (CCPA/CPRA) | Consumer-rights request records, automated-decision opt-out logs | 24 months for request audit trail (§1798.130); the underlying chain records continue to follow the 7-year federal floor |
| New York (DFS Part 500) | Cybersecurity-event audit trails | 5 years minimum; extends to 7 to align with federal floor |
| Texas, Illinois (consumer-credit) | State-licensed lender record-retention | 4 to 6 years; covered by federal 7-year floor |
| Massachusetts (201 CMR 17.00) | Personal-information protection program records | 5 years; covered by federal 7-year floor |

The institution applies the **"longer of"** rule: where state law extends beyond 7 years for a specific record type, the institution honors the longer period for that record type only. Where state law is shorter than 7 years, the institution applies the federal 7-year floor. The institution's records-retention schedule (maintained by the records-management team, separate from this document) lists per-record-type exceptions.

---

## 5. The end-of-retention deletion procedure

At the 7-year boundary (or end of any active legal hold), the institution deletes the privacy-store mapping for the affected records. The procedure is:

1. **Quarterly retention audit.** The records-management team identifies chain entries whose seal date is 7 years prior.
2. **Legal-hold check.** For each batch, the legal team confirms no active hold applies. Active holds extend retention for the held subset only.
3. **Privacy-store deletion.** The privacy-data-protection team (separate from chain operations per `token-vault-architecture.md`) deletes the token-to-PII mapping rows for the batch. This is the operative privacy action — it irreversibly severs the chain tokens from the original data.
4. **Chain ledger untouched.** The chain entries remain. Tokens are now anonymized per `pseudonymization-vs-anonymization.md`.
5. **Deletion log.** The privacy-data-protection team records: batch identifier, count of mappings deleted, date of deletion, name of deleting operator, legal-hold status confirmed. The deletion log is itself retained for 7 years to demonstrate the institution honored Article 5(1)(e).
6. **Verification.** A sample of deleted mappings is re-queried 30 days later to confirm permanent deletion (no soft-delete recovery).

The chain ledger is **never modified**. Erasure is achieved by removing the resolution path, not by altering the integrity-bound record.

---

## 6. Documented legal-hold extension

When litigation, regulatory inquiry, or examination opens against records that would otherwise reach the 7-year boundary, the legal team issues a **written legal hold** identifying the affected scope. The hold suspends end-of-retention deletion for the held records. The hold document carries:

- Hold identifier and effective date
- Scope (tenant, run-id range, customer-identifier range, or date range)
- Source of the hold (case caption, regulator inquiry number, examination cycle)
- Estimated duration and review date
- Issuing counsel's signature

When the hold is released, the records-management team resumes the standard quarterly retention audit for the released records and applies the deletion procedure in §5. Hold issuance and release are both logged in the privacy-log to demonstrate continuous control over retention boundaries.

---

## 7. Cross-references

- GDPR Article 5(1)(e) storage limitation — operationalized here.
- DPIA template `dpia-template.md` §3 — names this document as the proportionality analysis.
- RoPA template `ropa-template.md` — retention column references this document.
- `pseudonymization-vs-anonymization.md` — explains the anonymization transition that occurs at end of retention.
- `token-vault-architecture.md` — defines the team that performs the deletion (separate custody from chain operations).
- Spec §10.9 (IKM retention) — the institution's IKM-retention policy aligns with this document's 7-year window.
- `incident-response-playbook.md` litigation-hold scenario — describes how a hold is opened and operated.

---

## 8. Review cadence

This document is reviewed annually by the privacy office, the records-management team, and the legal team. Triggers for early review:

- A new statute or regulation extends a retention floor beyond 7 years for a record type the chain captures.
- The institution enters a new jurisdiction whose retention rules require recalibration.
- A regulator issues guidance changing the supervisory record-retention expectation.
- The institution's product mix changes such that healthcare data, children's data, or other special-category data flows through the chain (see `hipaa-minimum-necessary.md`, `childrens-data-safeguards.md`).

The annual review either reaffirms the 7-year period or documents the change in writing with revised legal grounding.
