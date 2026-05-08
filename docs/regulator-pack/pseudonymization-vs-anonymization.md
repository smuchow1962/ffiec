# Pseudonymization vs Anonymization — Why It Matters for the Chain

> **What this doc is.** A working analysis of the difference between pseudonymized data and anonymized data as those terms are defined under GDPR Article 4(11), NIST SP 800-188 (Guidance for De-Identifying Government Datasets), and HIPAA 45 CFR 164.514. Written so the institution can confidently say what the chain's tokens are at any given point in their lifecycle and what regulatory rules apply.

> **Why this matters.** "Tokenized" sounds final, but it isn't. While the institution retains the privacy-store mapping that resolves tokens to original personal data, the tokens are pseudonymized — still personal data under GDPR, still subject to data-subject rights, still in scope for transfer rules and breach notification. Only when the mapping is deleted do the tokens become anonymized and exit the regulatory scope. Confusing the two leads to under-protection during retention and over-protection after deletion.

---

## 1. The two terms in plain definitions

### Pseudonymization

GDPR Article 4(5):

> "'pseudonymisation' means the processing of personal data in such a manner that the personal data can no longer be attributed to a specific data subject without the use of additional information, provided that such additional information is kept separately and is subject to technical and organisational measures to ensure that the personal data are not attributed to an identified or identifiable natural person."

Two requirements: separation of the additional information, and protection of that additional information.

### Anonymization

GDPR Recital 26:

> "The principles of data protection should... not apply to anonymous information, namely information which does not relate to an identified or identifiable natural person or to personal data rendered anonymous in such a manner that the data subject is not or no longer identifiable."

The bar is high: irreversibility. If anyone with reasonable means can re-identify the data subject, the data is not anonymized.

NIST SP 800-188 frames this in three categories — direct identifiers, quasi-identifiers, and sensitive attributes — and warns that even datasets with all direct identifiers removed often allow re-identification through quasi-identifier combinations or auxiliary-data linkage. NIST's guidance is that anonymization is achieved through tested techniques (k-anonymity, differential privacy, synthetic data) measured against re-identification risk thresholds, not by token replacement alone.

---

## 2. Where the chain's tokens sit

The chain captures HMAC-SHA-256 tokens of personal-information fields:

```
token = HMAC-SHA-256(privacy_key, original)
```

The `privacy_key` lives in the privacy-store, in custody separate from the chain (`token-vault-architecture.md`). The mapping table (token-to-original) lives in the privacy-store database.

| Lifecycle phase | State of the privacy-store mapping | Chain token status | GDPR scope | HIPAA scope |
|---|---|---|---|---|
| Day 1 — entry sealed | Mapping retained | **Pseudonymized** | Personal data | PHI (if healthcare deployment) |
| Day 2 to Day (7×365) | Mapping retained | Pseudonymized | Personal data | PHI |
| Day (7×365)+1 — retention period ends, mapping deleted | Mapping deleted | **Anonymized** | Not personal data | Eligible for Expert Determination de-identification with documentation |
| After erasure request granted (within retention period if exception applies) | Mapping deleted for that subject | Anonymized for that subject | Not personal data for that subject | De-identified (subject to Expert Determination documentation) |

This is the critical claim: **chain tokens are pseudonymized, not anonymized, for the entire 7-year retention period**. They become anonymized only when the privacy-store mapping is deleted. The transition is operational (a deletion event), not cryptographic (the token bytes are unchanged).

---

## 3. Why tokenization alone is not anonymization

### 3.1 The reversibility test

Anonymization requires that the data subject is not or no longer identifiable by anyone using reasonable means. The chain's tokens fail this test while the privacy-store mapping exists, because:

- The institution operates the privacy-store.
- The privacy-store contains the mapping.
- The mapping is by design queryable (for data-subject access, breach response, audit).
- "Reasonable means" includes the institution's own routine operations.

The data subject is identifiable — by the institution itself — for as long as the mapping exists.

### 3.2 The HIPAA Safe Harbor test

45 CFR 164.514(b)(2) Safe Harbor requires removal of 18 enumerated identifiers including (R) "any other unique identifying number, characteristic, or code." The HMAC-token is a unique identifying code. It does not pass Safe Harbor.

### 3.3 The HIPAA Expert Determination test

45 CFR 164.514(b)(1) Expert Determination requires a qualified expert to document that the risk of re-identification is very small. With the privacy-store mapping retained, the risk is not very small — it is bounded by the security of the privacy-store, not by the difficulty of re-identification given the released dataset alone. Expert Determination cannot be claimed during the retention period.

After the privacy-store mapping is deleted, the institution can engage an expert to document an Expert Determination. The expert's analysis would consider: residual quasi-identifiers in the chain, auxiliary-data risks, re-identification methods. If the expert concludes risk is very small, Expert Determination is met and the chain is HIPAA-de-identified going forward.

### 3.4 The NIST SP 800-188 test

NIST recommends measuring de-identification through tested techniques against re-identification thresholds. Token replacement is not on its own such a technique. NIST's guidance supports the conclusion that tokenization is pseudonymization until the reverse-mapping is irrecoverably destroyed.

---

## 4. Implications during the 7-year retention period

### 4.1 GDPR rules apply

Because tokens are personal data (pseudonymized) during retention, all GDPR rules apply to the chain:

| Rule | Application |
|---|---|
| Article 5 principles | Lawfulness, fairness, transparency, purpose limitation, data minimisation, accuracy, storage limitation, integrity and confidentiality, accountability — all apply to chain processing |
| Article 6 lawful basis | Required (typically legal obligation per FFIEC + legitimate interests) |
| Article 13/14 transparency | Privacy notice must disclose the chain's existence as a processing activity |
| Articles 15–22 data-subject rights | Apply (with exceptions per Article 17(3) for legal-obligation retention) |
| Articles 33–34 breach notification | A chain-entry compromise is a personal-data breach |
| Article 35 DPIA | Likely required for high-risk processing |
| Articles 44–50 international transfers | Apply (see `international-transfers.md`) |

### 4.2 HIPAA rules apply

For healthcare-related deployments:

| Rule | Application |
|---|---|
| Privacy Rule (164.500–534) | Use and disclosure rules, Minimum Necessary Standard apply (see `hipaa-minimum-necessary.md`) |
| Security Rule (164.300–318) | Administrative, physical, technical safeguards apply (see `hipaa-security-rule-mapping.md`) |
| Breach Notification Rule (164.400–414) | A chain-entry compromise is a PHI breach (60-day clock) |

### 4.3 CCPA/CPRA rules apply

| Rule | Application |
|---|---|
| Right to Know, Delete, Correct, Opt-Out | Apply to the consumer's personal information including chain-token-resolvable data (see `ccpa-cpra-rights.md`) |
| Disclosure obligations | Privacy notice discloses chain processing as a category |
| §1798.105(d)(7) exception | Available for legal-obligation retention |

---

## 5. Implications after de-mapping (anonymized phase)

When the privacy-store mapping is deleted (end of 7-year retention or post-erasure event), the chain tokens transition to anonymized status. The implications:

| Rule | Application |
|---|---|
| GDPR | Out of scope (Recital 26 — anonymous data is not personal data); no further data-subject rights, no transfer restrictions, no retention limits |
| HIPAA | Eligible for Expert Determination de-identification with documented expert analysis; if Expert Determination is achieved, the data is de-identified PHI and outside Privacy Rule disclosure rules |
| CCPA/CPRA | "De-identified" per §1798.140(m) requires that the business takes reasonable measures to ensure the information cannot be associated with a consumer; commits to maintain and use the information in de-identified form; contractually obligates recipients to comply with the same. The institution that deletes the privacy-store mapping AND commits not to attempt re-identification meets §1798.140(m) for the chain tokens |
| Other regulations | Generally outside scope; chain becomes a non-personal-data audit ledger |

The institution's records-retention schedule recognizes this transition: post-de-mapping chain entries can be retained indefinitely (subject to the institution's own data-management policy), but typically follow the 7-year deletion practice for storage-cost reasons.

---

## 6. Documenting the transition in the DPIA

The DPIA (`dpia-template.md`) explicitly addresses the pseudonymization-to-anonymization transition. The DPIA section states:

> "Chain tokens are pseudonymized data during the retention period. They are personal data subject to GDPR, HIPAA, and CCPA/CPRA in full. At the end of the retention period (or post-erasure for individual subjects), the privacy-store mapping is deleted. With the mapping deleted, the tokens are no longer reversible and become anonymized. The institution does not attempt re-identification post-deletion. The transition is documented in the deletion log per `retention-justification.md` §5."

The DPIA's residual-risk analysis treats this transition as a risk-reduction event: the regulatory exposure surface shrinks dramatically once the mapping is deleted.

---

## 7. The "additional information" safeguard test

GDPR Article 4(5) requires that the additional information (the privacy-store mapping) be "kept separately and... subject to technical and organisational measures to ensure that the personal data are not attributed to an identified or identifiable natural person."

The institution's evidence that this requirement is met:

| Requirement | Evidence |
|---|---|
| Kept separately | Privacy-store hosted on different physical infrastructure or different cloud account from the chain ledger; separate database; separate access path (`token-vault-architecture.md`) |
| Technical measures | AES-256-GCM encryption at rest with HSM/KMS-protected key; separate encryption key from chain confidentiality key; access logging |
| Organisational measures | Privacy-data-protection team owns the privacy-store; chain-operations team does not have privacy-store access; access reviews quarterly; training for privacy-store operators |
| Pseudonymization is preserved through operations | Routine operations (verifier runs, recovery drills, audits) do not require privacy-store access; only data-subject-rights workflows do |

If the institution failed the "additional information separately kept" test, the chain's tokens would not even qualify as pseudonymized — they would be unprotected personal data, with stronger compliance burden. The architecture choice in `token-vault-architecture.md` is what makes the pseudonymization claim defensible.

---

## 8. Common confusions to avoid

| Confusion | Reality |
|---|---|
| "Tokenized data is anonymized" | No. Tokenized data is pseudonymized while the mapping exists |
| "We're encrypted at rest, so we're anonymized" | No. Encryption is a confidentiality measure; the data is still personal data, decryptable by authorized operators |
| "Hashes are not personal data" | No. Hashes of personal data with the original retained are pseudonymized data |
| "GDPR doesn't apply because we tokenize" | No. GDPR applies during the retention period; it stops applying after de-mapping |
| "We can ignore data-subject rights for chain entries" | No. Rights apply to the chain via the privacy-store; the chain itself is preserved as evidence |
| "Anonymized once, anonymized forever" | Mostly. But if the institution restores the privacy-store from a backup, anonymization is undone for that backup-restored subset; backups are subject to the same retention rules |

---

## 9. Cross-references

- GDPR Article 4(5), 4(11), Recital 26 — operationalized here.
- NIST SP 800-188 — referenced as the technical standard for de-identification.
- 45 CFR 164.514(b) Safe Harbor and Expert Determination — referenced for HIPAA scope.
- `privacy-by-design.md` — SDK tokenization implementation.
- `token-vault-architecture.md` — separate custody supports the "kept separately" requirement.
- `retention-justification.md` — names the deletion procedure that triggers the transition.
- `hipaa-minimum-necessary.md` — references this document for the "tokenization is not de-identification" claim.
- `dpia-template.md` — documents the transition as a risk-reduction event.
- `ccpa-cpra-rights.md` — references this document for the §1798.140(m) de-identification analysis.
- `international-transfers.md` — references this document for the transfer-scope question.

---

## 10. Review cadence

This document is reviewed annually by the privacy office, security architecture, and legal counsel. Triggers for early review:

- A regulator publishes new guidance on pseudonymization or anonymization.
- A CJEU or supervisory-authority decision affects the pseudonymization analysis.
- NIST publishes an update to SP 800-188.
- An expert-determination opinion is sought for a specific chain-entry subset and the analysis informs this document.
- The institution changes the tokenization algorithm or key-custody model.
