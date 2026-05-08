# HIPAA Minimum Necessary Analysis for the Chain

> **What this doc is.** A working analysis of the HIPAA Privacy Rule's Minimum Necessary Standard (45 CFR 164.502(b) and 164.514(d)) applied to the chain when it is deployed against healthcare-related decisions — healthcare-financed lending, medical-payment management, health-plan administration, or any AI system that processes Protected Health Information (PHI).

> **Why this matters.** A Covered Entity that captures more PHI than necessary into an integrity-bound, 7-year-retained ledger has a Privacy Rule problem before it has a Security Rule answer. The Minimum Necessary Standard is the door the institution must walk through before the chain's controls become relevant. This document is the institution's written walk-through.

---

## 1. The Minimum Necessary Standard in plain terms

45 CFR 164.502(b)(1) requires that, "When using or disclosing protected health information or when requesting protected health information from another covered entity or business associate, a covered entity or business associate must make reasonable efforts to limit protected health information to the minimum necessary to accomplish the intended purpose of the use, disclosure, or request."

45 CFR 164.514(d) operationalizes the standard: covered entities must identify those persons or classes of persons who need access to PHI to carry out their duties, and limit access accordingly; develop policies and procedures that restrict use and disclosure to the minimum necessary; and reasonably limit requests for PHI from other entities.

The standard does not apply to:

- Disclosures to or requests by a healthcare provider for treatment;
- Disclosures to the individual who is the subject of the PHI;
- Disclosures pursuant to an authorization;
- Disclosures required by law (45 CFR 164.512(a));
- Disclosures to HHS for compliance investigations.

The chain's purpose — integrity-bound audit trail of AI decisions — does **not** fall into these exemptions. The Minimum Necessary Standard applies in full.

---

## 2. Tokenization is not de-identification under HIPAA

A common misconception: "We tokenize PHI, therefore the chain doesn't have PHI in it." HIPAA does not recognize tokenization as de-identification.

45 CFR 164.514(b) recognizes two paths to de-identification:

| Path | Citation | Requirement |
|---|---|---|
| **Safe Harbor** | 164.514(b)(2) | Removal of 18 enumerated identifiers (names, geographic subdivisions smaller than state, dates more granular than year, telephone, fax, email, SSN, MRN, health plan numbers, account numbers, certificate/license, vehicle identifiers, device identifiers, URLs, IP, biometric identifiers, full-face photos, "any other unique identifying number, characteristic, or code") AND the covered entity has no actual knowledge that the information could be used alone or in combination to identify the subject |
| **Expert Determination** | 164.514(b)(1) | A person with appropriate knowledge of and experience with statistical and scientific methods determines, and documents, that the risk of re-identification is very small |

A token produced by HMAC-SHA-256 is **a unique identifying code** — it falls inside Safe Harbor's catch-all 164.514(b)(2)(i)(R). It is not Safe-Harbor de-identified. It can become Expert-Determination de-identified only if a qualified expert documents the analysis; the institution's standard tokenization configuration does not pass that bar because the privacy-store mapping is retained and reversibility is operationally available.

**Conclusion.** Tokenized PHI in the chain remains PHI for HIPAA purposes during the retention period. The chain is a use and disclosure of PHI subject to the Minimum Necessary Standard. After the privacy-store mapping is deleted (end of retention or post-Article-17 erasure for GDPR-overlapping cases), the tokens become irreversible — at that point, with appropriate documentation, they may qualify as Expert-Determination de-identified.

This conclusion mirrors the GDPR analysis in `pseudonymization-vs-anonymization.md`: tokenization is pseudonymization, not anonymization, until the mapping is deleted.

---

## 3. The Covered Entity (CE) perspective

When the deploying institution is itself a Covered Entity (a hospital, health plan, or healthcare provider), the CE has primary Minimum Necessary responsibility.

### 3.1 The CE's Minimum Necessary determination

The CE asks four questions before deploying the chain on PHI:

1. **What is the AI system's documented business need?** The CE writes down the AI's purpose: e.g., "approve or deny medical-payment plans based on credit history and current payment status." The business need defines the scope of PHI the AI must see.
2. **What is the minimum PHI that supports the business need?** The CE identifies the smallest set of PHI fields that, given the AI's documented purpose, supports the decision. This may be narrower than the AI's technical capability; the AI can ingest more, but the CE limits its access.
3. **What is the chain's role?** The chain captures what the AI saw and decided. The chain's PHI scope therefore mirrors the AI's PHI access scope. The chain does not add PHI; it records PHI access already authorized.
4. **What tokenization applies?** PHI fields are tokenized at the SDK boundary per `privacy-by-design.md`. Tokens are still PHI per §2 above, but tokenization mitigates exposure if the chain is breached without simultaneous privacy-store breach.

The CE's Minimum Necessary documentation lives in the policies and procedures required by 164.514(d) and is reviewed annually.

### 3.2 The CE's audit obligation

Under 164.530(c), the CE must implement appropriate administrative, technical, and physical safeguards to protect PHI confidentiality, integrity, and availability. The chain's controls (`hipaa-security-rule-mapping.md`) satisfy the technical and administrative safeguards. The CE's audit role is to confirm the chain's PHI scope matches the documented Minimum Necessary determination — i.e., that the SDK is configured to capture only what the determination authorizes.

---

## 4. The Business Associate (BA) perspective

When the deploying institution is a Business Associate of a CE — for example, a vendor providing an AI-driven payment-management service to a hospital — the BA's Minimum Necessary scope is **set by the CE, not the BA**.

### 4.1 The BAA defines the scope

45 CFR 164.504(e) requires a Business Associate Agreement (BAA) that specifies the BA's permitted uses and disclosures. The BA may not use PHI beyond the BAA's scope. For a chain deployed by a BA:

- The BAA names the AI system and its purpose.
- The BAA names the PHI categories the BA may receive and process.
- The BAA may name specific Minimum Necessary limitations (e.g., "only diagnosis codes relevant to payment-eligibility determination, not full clinical history").
- The BAA addresses the chain explicitly: the CE acknowledges the chain's audit-trail purpose and authorizes PHI capture in the chain to the extent it mirrors the AI's authorized PHI access.

### 4.2 The BA's implementation

The BA configures the SDK such that the captured fields match the BAA's authorization. Fields outside the authorization are not captured, even if the AI sees them transiently. The BA's chain-operations procedures document this configuration and the BAA-mapping behind it.

If the CE's Minimum Necessary determination is narrower than the AI's technical capability, the BA's configuration must enforce the narrower scope. The chain captures what the AI uses, not what the AI could use.

### 4.3 Sub-processors

If the BA uses sub-processors (cloud providers, downstream vendors), each sub-processor must sign a sub-BAA. The chain operations team confirms each sub-processor with PHI access has executed a sub-BAA before any PHI flows.

---

## 5. Institution-side guidance for healthcare AI deployments

### 5.1 Pre-deployment checklist

Before the chain is enabled on any AI system processing PHI, the institution confirms:

| Item | Owner | Evidence |
|---|---|---|
| AI purpose documented | AI product team | AI design document |
| Minimum Necessary determination written | Privacy officer + clinical informatics | Policy document under 164.514(d) |
| BAA in place (if BA deployment) or CE-internal authorization (if CE deployment) | Legal team | Signed BAA or internal authorization memo |
| SDK configuration matches Minimum Necessary scope | Chain operations + privacy office | Configuration file under change control; reviewed by privacy office |
| Tokenization applied to PHI fields | Chain operations | SDK config tokenization rules; test cases verifying tokens not PHI |
| Privacy-store custody separate from chain | Privacy data protection team | Per `token-vault-architecture.md` |
| Workforce training completed | HR + privacy office | Training records |
| Incident-response playbook covers PHI breach | IR commander | Updated playbook including HIPAA notification timing per `breach-notification-matrix.md` |

### 5.2 Ongoing operations

| Task | Cadence | Owner |
|---|---|---|
| Configuration review (SDK still matches Minimum Necessary scope) | Annually + on AI version change | Privacy office |
| Workforce access review (only authorized roles can access chain or privacy-store) | Annually | Privacy office + HR |
| Workforce training refresher | Annually | HR + privacy office |
| Risk analysis update (per Security Rule 164.308(a)(1)) | Annually + on material system change | Security office |
| Verifier run + recovery drill | Quarterly | Chain operations |
| BAA review (BA deployments) | Annually + on contract renewal | Legal |

### 5.3 Common deployment patterns

| Pattern | Minimum Necessary posture |
|---|---|
| Healthcare-financed lending (CE = healthcare financer) | CE configures SDK to capture credit-relevant data plus the minimum diagnosis/treatment data needed for payment decisions; tokenize SSN, MRN, account numbers; retain only the diagnosis-code categories used by the AI |
| Medical-payment management (BA serving hospitals) | BA configures per BAA; typically CPT/ICD codes plus payment status; full clinical history is excluded |
| Health-plan administration (CE = health plan) | CE configures for member-decision auditing; tokenize member ID, retain decision-relevant claims data; full PHI access is rare |
| Telehealth credentialing (CE = telehealth provider) | CE limits chain to credential-decision data (license status, training records); does not extend to clinical data |

Each pattern has a written Minimum Necessary memo. The memos are reviewed annually and updated when the AI's purpose, the data flows, or the regulatory environment changes.

---

## 6. Special considerations

### 6.1 Genetic information (GINA + HIPAA)

The Genetic Information Nondiscrimination Act layers additional restrictions on genetic data. If the AI system uses genetic information for decision-making, the institution applies:

- A separate Minimum Necessary memo for genetic data.
- Stricter SDK tokenization (genetic identifiers are tokenized regardless of other configuration).
- Confirmation that the AI's use of genetic information is permitted by GINA (which generally prohibits use of genetic information for employment and health-insurance decisions).

### 6.2 Substance-use disorder records (42 CFR Part 2)

Records of SUD treatment are subject to 42 CFR Part 2, which is more restrictive than HIPAA. If the AI processes Part 2 records, the institution applies the stricter Part 2 rules, not HIPAA Minimum Necessary alone. The chain's SDK configuration must reflect the Part 2 redisclosure restrictions; routine chain operations are not within the patient-consent scope unless the consent explicitly authorizes audit-trail logging.

### 6.3 Mental-health records (state law)

Several states (California, New York, Illinois) extend protection to mental-health records beyond HIPAA. The institution's configuration applies the strictest applicable rule.

---

## 7. Cross-references

- `hipaa-security-rule-mapping.md` — administrative, physical, and technical safeguards (the controls side; this document is the privacy/scope side).
- `breach-notification-matrix.md` — HIPAA 60-day notification and how it coordinates with FFIEC and GDPR clocks.
- `pseudonymization-vs-anonymization.md` — why tokenization is pseudonymization and tokens remain PHI.
- `token-vault-architecture.md` — separate custody for the token-to-PHI mapping.
- `privacy-by-design.md` — SDK tokenization configuration.
- `dpia-template.md` §2 — Data Categories incorporates the Minimum Necessary scope.
- 45 CFR 164.502(b), 164.514(b) and (d), 164.504(e) — operationalized here.

---

## 8. Review cadence

This document is reviewed annually by the privacy officer, legal counsel, and clinical-informatics lead. Triggers for early review:

- HHS issues new HIPAA guidance affecting the Minimum Necessary Standard.
- A new AI system enters scope (new clinical purpose, new data category).
- A BAA is renegotiated and changes the BA's authorized scope.
- An OCR (HHS Office for Civil Rights) audit or enforcement action highlights a Minimum Necessary issue at a peer institution.
