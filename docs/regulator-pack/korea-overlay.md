---
title: Korea Articulation Overlay — FSS / FSC examination, AI Basic Act 2025, PIPA, K-ISMS-P, EFSR §17, KISA coordination, AI ethics
status: informative
aligned-with:
  - Personal Information Protection Act (PIPA — 개인정보 보호법, March 2023 amendment)
  - Credit Information Use and Protection Act (CIUPA — 신용정보의 이용 및 보호에 관한 법률)
  - Electronic Financial Transactions Act (EFTA — 전자금융거래법)
  - Electronic Financial Supervisory Regulation (EFSR — 전자금융감독규정)
  - MyData regulation (효율적인 정보이용 및 보호에 관한 법률, effective January 2022)
  - AI Basic Act (AI 기본법, January 2025; effective January 2026)
  - FSC Innovative Financial Services Designation Program (regulatory sandbox)
  - FSS AI Guidelines for Financial Companies (issued July 2021, updated 2023)
  - K-ISMS-P (한국정보보호관리체계, mandatory 3-year recertification)
  - 망분리 (Network Separation Regulation)
  - KISA Cybersecurity Notification System (사이버 위협정보 공유 체계)
  - Cyber Defense Methodology and FSI Banking Cyber CERT
  - MSIT AI Ethics Standards (revised 2023)
  - AI Trustworthiness Verification System (안전한 AI 인증)
  - Equal Opportunity Employment Law (Korean fair-lending standards)
date: 2026-05-07
version: 1.0.0
---

# Korea Articulation Overlay

> **What this doc is.** A single articulation overlay that maps the FFIEC chain-of-custody v1.0a specification onto the Korean financial-sector supervisory framework and the Korean AI-ethics framework. Written so a competent authority — a Financial Supervisory Service examiner (FSS / 금융감독원), a Personal Information Protection Commission (PIPC / 개인정보보호위원회) examiner, a KISA Cybersecurity coordinator, a Korean megabank Chief Information Security Officer (정보보호최고책임자), or a Chief AI Ethics Officer accountable for AI Basic Act 2025 compliance — can read this document alongside the spec and confirm what the chain delivers, what it does not, and where supplementary measures are required for application in Korean jurisdiction.

> **What this doc is NOT.** Not a normative extension to the v1.0a specification. Not a translation of the spec into Korean. Not a Korea-specific spec fork. The integrity primitives in v1.0a are framework-neutral cryptographic constructs; the supervisory translation is an articulation overlay an institution layers on top, not a change to the underlying specification. Spec §1.2 (epistemic scope) governs: the chain proves what was said and that the record was not tampered with after capture; it does not prove the substantive correctness of the captured content.

> **Threat model.** Korean megabanks operate under the assumption of nation-state adversary compromise (Lazarus Group / 라자루스 그룹, APT37 / Reaper, Kimsuky) within 18 months. The overlay is written against this baseline. The Korean threat model parallels the INCD threat model in `bank-of-israel-overlay.md` §10; institutions operating under both Bank of Israel and FSS supervision read the two overlays together.

---

## 1. Scope and reading order

The overlay is consumed in six reading orders.

| Reader | Reading order |
|---|---|
| FSS examiner conducting comprehensive examination (종합검사) | §3 (FSS comprehensive examination workflow) → §6 (PIPA) → §7 (CIUPA) → §15 (FSS examination procedures) → §22 (translation table) |
| PIPC examiner | §6 (PIPA Article 37-2 right to refuse automated decision) → §8 (PIPC enforcement readiness) → §9 (PIPA §28-2 pseudonymisation) |
| Korean megabank CISO (정보보호최고책임자) | §10 (EFSR §17 24/7 SOC integration) → §11 (망분리 architecture) → §12 (K-ISMS-P control mapping) → §13 (KISA 2-hour clock) → §14 (nation-state threat model) → §16 (red-team) → §17 (incident response) |
| Chief AI Ethics Officer | §18 (AI Basic Act 2025 Articles 15-20) → §19 (right to explanation, 설명요구권) → §20 (fairness audit for Korean protected characteristics) → §21 (AI Trustworthiness Verification) |
| FBO Korean megabank operating in US | §3 → §15 → cross-reference to FFIEC examination procedures |
| MyData ecosystem participant | §6 (PIPA Article 26 cross-border) → §7 (MyData transmission requests) |

The overlay does not duplicate material in adjacent regulator-pack documents. Cross-references are exact — if a section sends the reader to `apac-overlay.md` §11.2, that section is the load-bearing source and this overlay names where it sits in the Korean framework picture.

---

## 2. The Korean supervisory architecture and the chain's place in it

### 2.1 Principal authorities

| Authority | Mandate | Chain interaction |
|---|---|---|
| Financial Services Commission (FSC / 금융위원회) | Sets sectoral financial regulation; issues AI Basic Act sectoral guidance for financial services; operates Innovative Financial Services Designation Program (regulatory sandbox) | Strategic policy oversight; the chain composes with FSC AI guidance and sandbox participation |
| Financial Supervisory Service (FSS / 금융감독원) | Supervises Korean financial institutions; conducts comprehensive examinations (종합검사); enforces FSS AI Guidelines (2021/2023) | Examiner-side consumer of verifier output (spec §7); reviewer of CC8.1 documentation; supervisor for FSS comprehensive examination |
| Personal Information Protection Commission (PIPC / 개인정보보호위원회) | Korea's data-protection authority; enforces PIPA including Article 37-2 (automated decision-making rights, strengthened March 2023); coordinates with FSS on AI-privacy enforcement | Examiner-side consumer of chain's data-handling posture; reviews PIPA §28-2 pseudonymisation, §26 cross-border, §37-2 automated decisioning rights |
| Korea Internet & Security Agency (KISA / 한국인터넷진흥원) | Operates Cybersecurity Notification System; receives 2-hour notifications for breaches affecting >1,000 customers; administers AI Trustworthiness Verification System under MSIT | Receives chain-detected incident notifications via the 2-hour clock; reviews AI Trustworthiness Verification submissions |
| Financial Security Institute (FSI / 금융보안원) | Operates Banking Cyber CERT; coordinates IOC sharing; conducts annual member-bank cyber drills | Receives anomaly reports from chain; participates in cyber drills involving chain scenarios |
| Ministry of Science and ICT (MSIT / 과학기술정보통신부) | National AI policy authority; administers AI Trustworthiness Verification; receives AI Basic Act incident notifications under Article 20 | Receives AI Basic Act §20 incident notifications; coordinates AI ethics committee oversight |

### 2.2 Regulatory landscape

| Statute / standard | Subject | Section in this overlay |
|---|---|---|
| PIPA (March 2023 amendment) | Personal data + automated decisioning rights | §6 |
| CIUPA | Credit-AI fairness and explainability | §7 |
| EFTA | Financial institution information-protection obligations | §10 |
| EFSR §17 | 24/7 SOC, dual approval, biennial red-team | §10, §16, §17 |
| MyData regulation | Consumer data-portability of AI-decision history | §6 |
| AI Basic Act 2025 | Risk-stratified AI governance | §18 |
| K-ISMS-P | National information-security baseline | §12 |
| 망분리 | Network separation | §11 |
| KISA Cybersecurity Notification System | 2-hour breach reporting | §13 |
| FSS AI Guidelines (2021/2023) | Five governance pillars | §18 |
| AI Trustworthiness Verification System | Voluntary certification | §21 |
| Equal Opportunity Employment Law / fair-lending | Korean adverse-action notice standards | §19 |

---

## 3. FSS comprehensive examination (종합검사) workflow

### 3.1 Pre-examination call

FSS conducts comprehensive examinations annually or every 18 months. The examination is 4-6 weeks on-site, with a pre-exam call 60 days before. During the pre-exam call, the FSS asks the institution's risk and compliance leadership to produce specific artefacts.

### 3.2 Pre-examination evidence package

The institution prepares the following 10 business days before the FSS comprehensive examination opens.

| Control area | Artefact |
|---|---|
| Governance | Board-approved AI governance policy (보드 결의 AI 정책) naming the chain; AI Governance Committee minutes (last 12 months); AI Basic Act §15 governance structure documentation |
| Key management | IKM rotation log; HSM access log (anonymised); KMS attestation; dual-approval evidence per EFSR §17(3) |
| Daily operations | Seal records with timestamps (365 files for 1-year period); verifier output (daily PASS/FAIL summary); daily seal-verification SIEM alerts |
| Incidents | IR log for any chain-detected anomalies; remediation records; KISA notification records (if any); INCD-equivalent FSI Banking Cyber CERT notifications |
| Configuration | Posture declaration (FFIEC vs vendor-namespaced per spec §4.1.2); HKDF salt and info parameter documentation; format-version consistency checks |
| AI governance | AI Impact Assessment (AI 영향평가) for each high-impact AI system; fairness-audit reports (quarterly or semi-annual); customer-dispute log |
| Vendor management | Annual vendor audit per Directive-414-equivalent procedure; SOC 2 Type II report for vendor-hosted topology; vendor SLA breach record |
| Network separation | 망분리 architecture diagram; FSS approval letter for cross-boundary flows (if any); firewall rules for OTLP ingress |
| K-ISMS-P | Most recent K-ISMS-P recertification audit report; surveillance audit findings; control implementation guide for controls 2.10, 3.2, 9.1-9.3, 13 |

### 3.3 Examination duration and scope

The examination examines governance, risk management, compliance, and operational controls. For chain-of-custody specifically, the examiner verifies the chain operates as documented; the verifier output passes; the institution's incident-response procedures are testable; and the documented control framework satisfies the FSS's expectations under FSS AI Guidelines, K-ISMS-P, and EFSR §17.

### 3.4 Examination findings and 검사확인서 (examination confirmation letter)

After the examination, the FSS issues a 검사확인서 to the institution's CEO. The letter names findings, classifies severity, and prescribes remediation timelines. For chain-of-custody control, three classification tiers apply.

| Tier | Trigger | Remediation |
|---|---|---|
| Control operating effectively | All sampled entries pass verifier; all seal records verify; anomaly review documented monthly; operational-resilience drill within target; no Directive-equivalent deficiencies | None |
| Material control deficiency | Recovery-time drill fails target; chain-detected incident not escalated to KISA per timing; vendor audit not conducted | Immediate action plan and corrective action within 30 days |
| Control deficiency requiring attention | Anomaly review not documented monthly; vendor audit findings exist but under remediation; IKM rotation documented but fingerprint registry out of sync for a period (root cause identified and remediated) | Corrective action plan and remediation within 90 days |

---

## 4. FBO dual supervision — FSS + Federal Reserve

### 4.1 FBO context

Some Korean megabanks operate US subsidiaries or branches (Foreign Banking Organisations under Federal Reserve supervision per 12 CFR Part 211). When a Korean bank operates an FBO in the US and deploys a global AI system, the system is subject to both FSS supervision (of the Korean parent) and Federal Reserve supervision (of the US subsidiary).

### 4.2 Single technical control, dual-regulator consumption

The chain's design satisfies both FFIEC examination standards and FSS comprehensive-examination expectations. Both regulators independently verify the chain without re-testing.

| Examination | Regulator | Scope | Evidence |
|---|---|---|---|
| FFIEC examination | Federal Reserve | US subsidiary's operational-risk control | Verifier output per spec §7; audit procedures P-1 through P-10 |
| FSS comprehensive examination | FSS | Korean parent's AI governance per FSS AI Guidelines and AI Basic Act | Verifier output; CC8.1 documentation; AI Impact Assessment; fairness-audit reports |

### 4.3 Cross-regulator findings

If one regulator issues an MRA (e.g., FSS issues MRA on fairness-testing procedures), the institution remediates and the other regulator (Federal Reserve) is informed per supervisory-coordination protocols. The institution's CC8.1 explicitly names dual-jurisdiction purpose.

---

## 5. Conglomerate consolidated supervision (통합 감독)

### 5.1 Chaebol structure and circular shareholding

Many Korean megabanks are part of larger financial conglomerates with circular-shareholding structures (KB Financial Group, Shinhan Financial Group, Hana Financial Group). The FSS conducts consolidated supervision of financial groups; during a group-level comprehensive examination, the FSS examines AI-decision consistency across the group's affiliates.

### 5.2 Group-level governance and per-affiliate operational independence

| Layer | Responsibility |
|---|---|
| Group AI policy | Board-approved policy names group-level AI Governance Committee, fairness-testing requirements, incident-response escalation, regulator coordination |
| Entity-level implementation | Each affiliate (bank, securities, insurance) implements the chain per group policy framework with separate `tenant_id` and HSM custody per entity (operational independence, separation of duties) |
| Group-level audit | Group internal audit conducts annual audits across all affiliates testing consistency (verifier version, fingerprint-reconciliation procedure) |
| Consolidated examination coordination | Group's compliance function produces coordinated evidence for cross-affiliate examination |

### 5.3 Cross-affiliate query under data-sharing agreement

When customer disputes cross affiliates (e.g., the credit score used by a securities firm differs from the bank's score), the bank and securities company exchange chain entries via a data-sharing agreement (정보 공유 약정) under PIPA Article 18 (consent-based sharing) or the Financial Group's Internal Information Sharing Protocol (금융그룹 내부 정보 공유 규정).

---

## 6. PIPA Articles 23, 26, 28-2, 37-2 and MyData

### 6.1 PIPA Article 37-2 — right to refuse automated decision-making (이의제기권)

The March 2023 amendment substantially strengthened PIPA Article 37-2. A customer denied credit by an AI system can demand human review. The chain captures the AI's preliminary decision; the customer's invocation of Article 37-2; the human review event; the human's decision; whether the human's decision differed from the AI's.

| Step | Chain attribute |
|---|---|
| Customer objection received | `audit.pipa.article_37_2.objection_received_at` (RFC 3339 UTC); `audit.pipa.article_37_2.original_ai_decision_seq` (parent entry seq); `audit.pipa.article_37_2.objection_reason` (Korean-language customer explanation) |
| Human reviewer assigned | `audit.pipa.article_37_2.human_review_assigned_to` (reviewer identifier); `audit.pipa.article_37_2.review_deadline` (typically 14 days per PIPA Article 37-4) |
| Human review outcome | `audit.pipa.article_37_2.reviewer_decision` (approve/modify/reverse); `audit.pipa.article_37_2.decision_rationale` (Korean text); `audit.pipa.article_37_2.decision_matches_original_ai` (boolean) |
| Customer notification | `audit.pipa.article_37_2.notification_sent_at_utc`; `audit.pipa.article_37_2.notification_language = ko-KR` |

### 6.2 설명 요구권 (right to demand explanation) vs 이의제기권 (right to object)

PIPA Article 37-2 establishes two distinct statutory rights with distinct timelines.

| Right | Timeline | Chain pathway |
|---|---|---|
| 이의제기권 (right to object) | Human review begins within 14 days | Triggers human-review pathway per §6.1 |
| 설명 요구권 (right to demand explanation) | Explanation provided "promptly" — typically 5 business days | Institution extracts decision rationale from chain entry; translates to Korean plain language; provides explanation |

### 6.3 PIPA Article 26 — cross-border transfer

PIPA Article 26 restricts transfers of personal data to foreign jurisdictions lacking equivalent-protection status. The US does not have equivalent-protection status per PIPC's 2023 guidance (unlike the EU). Compliant deployment topologies for Korean institutions:

| Topology | PIPA Article 26 conformance |
|---|---|
| **De-mapping before US transfer (strongest):** Korean institution retains the privacy-store mapping on Korean servers; US system receives only chain tokens, without reversal capability | Conformant; institution files transfer-impact assessment (개인정보 이전영향평가) with PIPC naming the de-mapping as supplementary safeguard |
| **MyData consent-mediated transfer:** Consumers explicitly consent under MyData regulation; chain tokens flow to the intermediary with consumer consent | Conformant if consent scope is documented; institution records `audit.mydata.consent_id` and `audit.mydata.consent_scope` |
| **Korean-only retention:** Institution retains chain within Korean-only infrastructure | Conformant by default |
| **Direct transfer without supplementary safeguards** | Non-conformant |

### 6.4 PIPA Article 28-2 — pseudonymisation (가명처리)

PIPA Article 28-2 (amended March 2023) defines pseudonymisation as processing of personal data such that the data cannot be attributed to a data subject without additional information held separately. Under PIPA §28-2, pseudonymised data is treated more leniently than identified personal data.

The institution chooses one of two postures.

| Posture | Treatment under PIPA |
|---|---|
| **True pseudonymisation:** Customer identifier is one-way hashed or irreversibly transformed; the institution does NOT retain a mapping that allows reversal | Pseudonymised under PIPA §28-2; chain entries can be more flexibly shared, analyzed, and disclosed under broader PIPA exemptions |
| **Encryption-with-reversible-key:** Customer identifier is encrypted under a key the institution holds and controls; institution maintains the encryption key separately | NOT pseudonymisation; remains personal-data-under-encryption; all PIPA personal-data restrictions apply |

The institution's CC8.1 explicitly names the chosen approach and the rationale. FSS examiners will verify the choice is documented and consistently applied.

### 6.5 PIPA Article 23 — sensitive personal information

AI inference logs that describe identifiable individuals are personal data under PIPA §1 even when tokenised, if the token-to-original mapping exists. The institution's CC8.1 names the privacy-store custody, the access controls, and the breach-notification procedure (PIPA §7: 72-hour notification to PIPC; affected individuals within 60 days).

### 6.6 MyData transmission requests including AI-decision history

Under the MyData regulation, consumers can demand transmission of their financial data including AI-decision history. The institution's procedure:

1. Customer submits MyData transmission request (API call or web request).
2. Institution extracts chain entries for the customer over the requested period (typically last 1 year), filtered by `chain_kind = audit` and customer identifier.
3. Institution prepares MyData-format response (JSON per Korean MyData standard schema) including consumer identifier, AI decision type, decision timestamp, decision inputs (tokenised), decision output, and (if available) any human-review override.
4. Institution transmits via authorised MyData intermediary (인증된 데이터 중개기관) to the downstream recipient.
5. Institution records the transmission in the chain as a new entry with `chain_kind = audit` and `audit.mydata.transmission_request_id`, `audit.mydata.transmitted_entries_seq_range`, `audit.mydata.recipient_id`.

---

## 7. CIUPA Article 15 fairness testing for credit AI

### 7.1 What CIUPA requires

CIUPA Article 15 requires institutions to conduct fairness testing for credit-scoring AI, measure disparate impact across protected characteristics (gender, age, national origin, credit-worthy status), and explain individual credit decisions citing the factors that drove the decision.

### 7.2 Population-level fairness testing

The chain captures the AI's inputs and outputs per decision (supporting individual explainability per PIPA §37-2). For population-level fairness testing required under CIUPA Article 15, the institution maintains a separate, governance-controlled de-identification mapping (not part of the chain; not subject to chain retention rules; governed by privacy-by-design.md logic).

| CIUPA control | Implementation |
|---|---|
| Raw demographic fields available for fairness testing | Stored in separate encrypted-but-decryptable form available only to the institution's fairness-testing team under strict access controls |
| Chain entries themselves | Tokenised for external verification and customer-dispute purposes |
| Fairness-audit dataset | Separate institutional artifact derived from chain data; not itself on the chain; used for demographic-disparity analysis |

### 7.3 CIUPA dispute procedures

When a customer disputes a credit-scoring AI decision, the customer can file a complaint with: (1) the issuing bank (직접 민원 — direct complaint), or (2) the Korea Credit Information Services Association (KCISA / 신용정보협회) or the credit-information agency that holds the score record.

| Step | Activity |
|---|---|
| Bank receives complaint | Compliance team pulls the chain entry for the disputed decision (indexed by customer ID and decision date), verifies integrity using the verifier, produces a compliance report |
| Agency receives complaint | If customer filed with the credit agency, the agency requests the chain entry from the bank and conducts independent review (chain entry is independently verifiable per the spec's design) |
| Institution response | Bank submits written response (이의 제기 처리 내용 — dispute-handling report) to the agency or customer within 30 days, citing chain evidence |
| Agency determination | Agency determines whether the complaint is valid; may order the bank to remove inaccurate credit information from the customer's file |

---

## 8. PIPC examination readiness

### 8.1 PIPC jurisdiction and triggers

The PIPC may conduct on-site examinations of AI-decision-making systems triggered by:

- Customer complaints about automated decision-making (PIPA §3 right to human review).
- Data-subject requests for access or erasure that institutions fail to fulfill within 30 days.
- Incident notifications of privacy-store breaches (PIPA §7).
- Random audits of critical-sector organisations (banks fall in this category).

### 8.2 PIPC examination evidence package

| Evidence area | Content |
|---|---|
| Governance documentation | Board policy naming AI systems subject to PIPA §37-2; AI Governance Committee charter; AI Impact Assessment (AI 영향평가) documents |
| Safeguards description | Bank describes the chain as a technical control: "Every AI decision is captured in an integrity-bound audit trail; the audit trail enables fairness audits and customer dispute investigation" |
| Customer-rights procedures | Right to refuse automated decision-making (opt-out procedure); right to explanation (Q&A pathway with chain extraction); right to manual review (per §6.1) |
| Fairness audit evidence | Verifier reports for recent periods; fairness-audit reports showing demographic decision distributions and any detected discrimination |
| Incident and complaint log | If PIPC has received complaints from customers about specific AI decisions, the bank extracts chain entries (de-identified per privacy law) to explain to PIPC what happened |

### 8.3 PIPC enforcement posture

PIPC may issue directives to cease or modify processing if it determines processing violates PIPA. Failure to comply is an administrative fine of up to 5 percent of global turnover under PIPL §12. For a Tier-1 megabank this can be tens of billions of won.

### 8.4 FSS-PIPC coordination

The FSS and PIPC coordinate on AI governance: FSS focuses on financial governance, operations, and risk management; PIPC focuses on personal-data protection, cross-border-transfer rules, and individual rights. The institution prepares an integrated evidence package satisfying both regulators; the chain artefacts and CC8.1 are unified across the two examinations.

---

## 9. PIPA §28-2 and pseudonymisation declaration

The institution's CC8.1 explicitly names the chosen pseudonymisation posture (per §6.4 of this overlay), supports the choice with technical evidence (the tokenisation algorithm, the key-custody arrangement), and documents the FSS-PIPC examiner walkthrough demonstrating the choice is consistently applied across the chain's lifetime.

---

## 10. EFSR §17 — 24/7 SOC integration and dual approval

### 10.1 EFSR §17(2) real-time monitoring

EFSR §17(2) (전자금융감독규정 제17조 제2항) mandates real-time monitoring — the institution's 24/7 SOC must detect cyber incidents on production systems within minutes. The chain's daily Merkle seal provides 24-hour detection latency, which does NOT satisfy real-time monitoring on its own. The institution supplements the chain with real-time detection mechanisms.

| Real-time mechanism | Purpose |
|---|---|
| Host-based detection on ledger server | Monitor for unexpected process spawning, file modifications, OTLP receiver changes |
| OTLP ingest log anomaly detection | SIEM with ML-based anomaly detection flags unusual patterns (1000x spike in rejections, zero events for >1 minute, out-of-sequence arrival) |
| Ledger-storage append-only monitoring | Append-only storage monitored for unexpected UPDATE/DELETE attempts; immediate alert on append-only invariant violation |
| Daily seal-verification result alerting | Daily seal-verification job is integrated with SOC's incident-management system; FAIL result is immediately escalated as Severity-1 incident |

The chain is the long-tail evidence archive; real-time detection is the immediate-response mechanism. The institution's CC8.1 names all real-time detection mechanisms in the control description and provides SIEM screenshot evidence showing the chain as one component of a multi-layered real-time posture.

### 10.2 EFSR §17(3) dual approval (이중 승인)

EFSR §17(3) mandates dual-approval for production changes affecting customer data or security controls. Two authorised personnel from different organisational units must approve before the change is executed; approvers sign with cryptographic evidence.

| Operation | Approvers |
|---|---|
| IKM generation (initial provisioning) | CISO + Chief of Operations |
| IKM rotation (routine quarterly or emergency post-compromise) | CISO + Chief of Operations |
| Daily seal-job execution | SOC manager + Operations manager (each enters one-time PIN from hardware token) |
| Ledger-server maintenance | Chief of Operations + DBA Lead |
| Verifier-binary update for examiner testing | CISO + Chief of Internal Audit |
| Key recovery (restoring IKM from backup after disaster) | CISO + Chief of Operations + Bank Executive (CFO or CRO) — three-person ceremony |

Approval evidence retention: 7 years per K-ISMS-P control 9.3 (see §12). Logs include timestamp, operation type, approver 1 and 2 (name, user ID, digital signature), operation outcome, and witness (if applicable).

### 10.3 EFSR §17(8) biennial red-team

See §16 of this overlay.

### 10.4 EFSR §17(1) 24/7 SOC and ledger-server resilience

The chain ledger and HSM operate within the institution's 24/7 SOC posture. Operational-resilience patterns:

| Pattern | RPO | RTO |
|---|---|---|
| Primary-backup | ~1 hour | ~30 minutes |
| Active-active | ~1 minute | Near-zero (failover transparent to SDK) |
| Hot-standby HSM with 15-minute synchronisation | ~15 minutes | ~15 minutes |

If the ledger is offline, the SDK's local SQLite buffer continues to capture events for up to N days; once the ledger recovers, the SDK exports the buffered events; no events are lost.

---

## 11. 망분리 (network separation) architecture

### 11.1 What 망분리 requires

Korean network-separation regulation (망분리 지침) requires logical and (for tier-1 banks) physical separation between internal banking systems and the internet. Tier-1 banks cannot allow customer-identified data to cross the 망분리 boundary into the internet or into vendor-hosted cloud without explicit FSS approval and a separate network-access control document (정보보호관리방침 추가 사항).

### 11.2 Compliant deployment topologies

| Topology | 망분리 conformance |
|---|---|
| **On-premises ledger, on-premises SDK:** Both on internal network; air-gapped from internet except outbound API calls to LLM provider | Conformant |
| **On-premises ledger, cloud SDK with transit encryption:** SDK runs in cloud-hosted GenAI service (e.g., AWS SageMaker in Seoul region, approved per cloud framework MOU); OTLP wire encrypted (TLS 1.3, mTLS with client certificate pinning); ledger remains on-premises | Conformant |
| **FSS-approved cross-boundary topology:** Vendor-hosted in non-Korean jurisdiction; institution submits 정보보호관리방침 추가 사항 naming vendor, jurisdiction, data minimisation (tokenisation), encryption-key custody (institution holds keys, vendor cannot decrypt) | Conformant subject to FSS approval |
| **Vendor-hosted multi-tenant US/EU cloud processing Korean customer data** | Non-conformant without explicit FSS approval |

### 11.3 OTLP transport hardening for the 망분리 environment

| Control | Specification |
|---|---|
| TLS version | TLS 1.3 strict; TLS 1.2 disabled at the receiver endpoint |
| Certificate pinning | Recommended for tier-1 banks: HPKP or pin-on-first-use; backup public key in pin configuration |
| Mutual TLS | Receiver requires mTLS; client certificate signed by trusted CA; certificate's tenant_id matches request's Resource attribute `ffiec.chain.tenant_id`; mismatch triggers rejection |
| Network segmentation | OTLP ingress firewalled to accept connections only from named SDK sources; rejected-traffic logging triggers SIEM alert at >10 rejections in 5 minutes |
| SIEM monitoring | Detect TLS 1.2 connection attempts; certificate validation failures; mTLS rejections; cipher-suite downgrades |

### 11.4 Korean BCP multi-region — Seoul/Busan failover

Korean tier-1 institutions operate active-active BCP topology per FSS expectations.

| Region | Role |
|---|---|
| Seoul (primary) | Ledger, HSM with master key, seal job, SDK deployment |
| Busan (secondary) | Standby ledger, failover HSM with master key, standby seal job, live SDK deployment |

Cross-region key-rotation coordination ensures both regions' keys are derived from the same IKM via HKDF with same `key_version`; entries from both regions use the same `key_fingerprint`. Verifiers validate seals signed by either region using the same public key.

---

## 12. K-ISMS-P control mapping

### 12.1 K-ISMS-P recertification context

K-ISMS-P (한국정보보호관리체계) is mandatory for all Korean financial companies. 3-year recertification with annual surveillance audits. Failure results in operational restrictions, fines up to 500M KRW, and reputational damage.

### 12.2 Control 3.2 (암호화) — cryptographic key management

| Audit-trail artifact | Specification |
|---|---|
| IKM generation logs | Timestamp (RFC 3339 UTC); operator (user ID, name); method (HSM key-generation function); algorithm; key length; `key_version` assigned; tenant_id; `key_fingerprint` computed; HSM serial number. Retention 7 years; read-only by CISO and internal audit |
| IKM rotation logs | Start and completion timestamps; operator; reason (scheduled quarterly, emergency); old `key_version`; new `key_version`; ceremony type; witnesses (EFSR §17(3) signatories); success/failure status. Retention 7 years |
| IKM backup/restore logs | Timestamp; operator; destination (secure storage location); verification checksum; witnesses; success/failure status. Retention 7 years |
| IKM access logs | Timestamp; operator; API call (sign, derive, rotate); parameters; success/failure; error code; source IP/certificate. Retention 1 year operational; 7 years for annual internal audit samples |

### 12.3 Control 9.1 (로깅) — audit logging completeness

| Audit-trail artifact | Specification |
|---|---|
| Ledger-ingest logs | Every OTLP event received, accepted, or rejected: timestamp; tenant_id; run_id; seq; event_kind; payload_hash; key_fingerprint; HMAC verification result; error code if failed; source IP. Retention 7 years |
| Seal-operation logs | Daily seal job execution: start timestamp; tenant_id; event count sealed; Merkle root hash; HSM signing request/response; seal signature; completion timestamp; status. Retention 7 years |
| Verifier-run logs | Auditor or examiner runs: timestamp; operator; ledger file hashed; verifier binary SHA-256; public key used; verification result; events processed; events that failed; duration. Retention 7 years |

### 12.4 Control 9.2 (감시) — real-time monitoring integration

Chain-of-custody events integrated into SOC/SIEM (Splunk, ArcSight, or equivalent): HMAC verification failure → medium-severity alert; seal-verification failure → critical alert; key-rotation event → informational alert (tracked for quarterly access review). Alerts correlated with other security events; HMAC failure + unauthorised ledger access logs within 5 minutes → Severity-1 incident.

### 12.5 Control 9.3 (정보보호 기록의 보관) — log integrity and retention

All logs (control 3.2 IKM logs and control 9.1 ingest/seal/verifier logs) are stored in append-only ledger or separate append-only audit log. Logs are protected from tampering via the same HMAC + daily Merkle seal. Backup copies maintained in separate secure location with access controls. Annual internal audit samples logs to verify completeness and integrity.

### 12.6 Control 2.10 (접근제어) — personal-data flow documentation

Quarterly personal-data-flow documentation maintained in designated tool (지정 도구). Chain-of-custody system entry includes:

- System name: "Chain-of-Custody Ledger (FFIEC v1.0a)".
- Purpose: audit trail for AI-agent decision events.
- Personal-data types: customer identifiers, transaction details, decision inputs, decision outputs.
- Retention period: 7 years per control 9.3.
- Access controls: read-only by CISO, internal audit, compliance team, examiners during examination, legal team for litigation response. Write-only by SDK host (via OTLP). No update/delete (append-only).
- Data-handling operator: ledger-server operations team; HSM custody team; verifier operators.
- Location of processing: on-premises ledger server in Seoul colocation facility (or Busan backup).
- Data-transfer mechanism: OTLP/gRPC over TLS 1.3; network-segmented (internal network only).

Quarterly review and sign-off by Chief Data Governance Officer or equivalent.

### 12.7 Control 13 (사건대응) — incident response

The institution's IR runbook integrates with K-ISMS-P control 13. See §17 of this overlay.

### 12.8 K-ISMS-P recertification audit checklist

The bank's auditor will ask for evidence in each of the above control areas. The bank provides:

- Seal-operation logs export filtered to seal events.
- Verifier output for the audit period (Merkle seal verification covering those events).
- Dual-approval records (names of operators per EFSR §17(3)).
- Personal-data inventory and data-flow diagram (PNDA — 개인정보 처리 흐름도).

If any log is missing, incomplete, or tampered with, the control fails and the audit outcome is "Non-compliant; remediation required."

---

## 13. KISA Cybersecurity Notification System — 2-hour clock

### 13.1 The 2-hour clock

The KISA Cybersecurity Notification System (사이버 위협정보 공유 체계) requires that if a breach affects more than 1,000 customers, the bank must notify KISA within 2 hours of determining the breach occurred. The chain's daily seal detection window is 24 hours; the bank cannot wait 24 hours. Determination is based on available evidence (anomalous logs, SOC alerts, user reports, initial forensics), not final forensic conclusion.

### 13.2 Two-phase response

| Phase | Activity | Chain integration |
|---|---|---|
| Phase 1 (0-2 hours): Preliminary determination | Detect via SOC alerts, IDS, user reports; triage within 30 minutes; if breach affects >1,000 customers, notify KISA within 2 hours of triage completion | Chain evidence not yet available; KISA notification clearly states preliminary status |
| Phase 2 (2-24 hours): Full investigation | Daily seal-verification job runs at next UTC midnight; if root mismatch, confirms compromise was real; bank files follow-up to KISA with chain evidence | Chain evidence (verifier output, seal verification result) attached to follow-up |
| Seal passes (no tampering detected) | Bank's follow-up clarifies: "24-hour chain-of-custody verification found NO evidence of ledger tampering. This does NOT prove the incident did not occur; it proves that IF the incident occurred, it did not involve tampering with our audit chain. We continue forensic investigation via network logs, endpoint logs, and third-party forensics" | |

### 13.3 KISA notification template

The institution's IR runbook includes a preliminary KISA notification template naming: incident type; affected customer count (estimated); discovery timestamp; incident-start timestamp (estimated); response actions initiated; preliminary impact; the bank's incident commander contact. The template is reviewed quarterly to align with KISA's evolving guidance.

---

## 14. Nation-state threat model — Lazarus Group, APT37, Kimsuky

### 14.1 The Korean threat baseline

Korean tier-1 banks operate under the assumption of nation-state adversary compromise within 18 months. The threat is not theoretical; it is operational baseline. Lazarus Group has demonstrated capability to compromise banking infrastructure for months with no detection (2014 SWIFT hack, 2016 Bangladesh Bank heist), forge transaction records and authentication logs simultaneously, and exploit supply-chain vulnerabilities in SWIFT software and ATM networks.

### 14.2 Adversary capabilities

The adversary controls:

- Supply-chain access to HSM firmware vendors, OTLP library maintainers, or SDK build pipelines.
- Network position (ability to intercept or inject OTLP traffic on the bank's internal network).
- Process-level compromise (rootkit implant on the ledger server, SDK host, or HSM appliance).
- Long persistence window (18+ months before detection).

### 14.3 Chain defense

The chain's three-layer composition provides sequential defense (parallel to `bank-of-israel-overlay.md` §10.3). Detection latency: 24 hours under default daily cadence; 1 hour under hourly cadence.

### 14.4 Residual risk gaps and compensating controls

| Residual risk | Compensating control |
|---|---|
| HSM firmware-supply-chain compromise | Annual firmware-integrity attestation via HSM vendor (where supported) or independent auditor sign-off; HSM vendor switching threshold |
| Vendor-code compromise | Annual code reviews of Herald.Compliance source code; quarterly build reproducibility verification; verifier binary held in secure escrow (see §17.3) |
| OTLP network-position compromise | Network segmentation (VLAN isolation, microsegmentation); network-based IDS for OTLP traffic; mTLS with strict cipher-suite filtering |

---

## 15. FSS examination test procedures

### 15.1 Verifier procedure

Examiner downloads reference verifier (or FFIEC-conformant verifier), runs on sample of 100 chain entries from random dates, confirms 100% pass rate.

### 15.2 Seal integrity

Sample 12 seal records (one per month). For each seal, recompute Merkle root from sealed events and verify signature against the institution's public key. All 12 must verify successfully.

### 15.3 Anomaly review

Request 12 months of weekly anomaly-review meeting minutes. Spot-check 4 months for IT operations / IT risk / audit attendance, anomaly identification, severity classification, escalation to CRO or CISO when warranted.

### 15.4 Incident-response integration

Review institution's incident-response log. Confirm KISA notification timeline (within 2 hours for breaches affecting >1,000 customers); confirm FSS notification timeline; confirm institution coordinated with KISA on forensics and remediation.

### 15.5 Vendor audit (if vendor-hosted)

Review past 3 years of annual vendor audits. Confirm audits covered HSM custody, ledger integrity, seal-job automation; findings remediated or remediation plan executed; vendor's SOC 2 report obtained and reviewed.

---

## 16. Red-team exercise readiness (EFSR §17(8))

### 16.1 Biennial scope

EFSR §17(8) mandates red-team exercises every 2 years. Exercises coordinated with FSS and (for some banks) National Intelligence Service (국정원). Scope includes external network attack surface, internal network, and critical systems. The chain-of-custody system may be in scope.

### 16.2 Exercise scenarios

| Scenario | Red-team objective | Expected chain outcome |
|---|---|---|
| SDK compromise | Forge a chain event going forward | Verifier detects forged events via HMAC chain mismatch; red team fails |
| OTLP tampering | Inject, modify, or replay OTLP messages | Ledger receiver re-verifies HMAC chain on ingest; forged messages rejected; red team fails |
| Ledger-server compromise | Rewrite events and recompute Merkle root | Next daily seal-verification recomputes root from modified events; root mismatches the prior signed root; tampering detected; red team fails |
| HSM firmware compromise | Sign forged root under bank's key | Residual risk; firmware-attestation compensating control detects firmware version mismatch; otherwise out-of-band detection |
| Simultaneous compromise | Forge a history that appears valid | Only succeeds if all three trust boundaries are simultaneously compromised |

### 16.3 Pre-exercise readiness

Verifier binary held securely (air-gapped, independently rebuilt). Public key registered with the institution, FSS, and (where applicable) Federal Reserve. Daily verification run operating to detect forged signatures. Incident-response procedures live and testable.

### 16.4 Post-exercise reporting

Red team documents which attacks succeeded, time-to-detection, remediation recommendations. CISO maps findings to control improvements. Exercise report filed with FSS as evidence of cyber-resilience.

---

## 17. Incident-response procedures (Directive-equivalent)

### 17.1 Multi-regulator coordination

Korean banks face simultaneous incident-notification clocks: KISA 2-hour clock for >1,000-customer breaches; FSS notification (and Federal Reserve for FBO operations); FSI Banking Cyber CERT 24-hour clock; PIPC 72-hour personal-data breach notification; MSIT AI Basic Act §20 notification (timeline pending guidance, anticipated 5 business days).

### 17.2 Severity taxonomy

| Tier | Trigger | KISA timing | FSS timing |
|---|---|---|---|
| Level 1 — breach affecting >1,000 customers | 2 hours | 36 hours | |
| Level 2 — breach affecting 100-999 customers | Optional | 36 hours | |
| Level 3 — breach affecting <100 customers or internal-only data | None | Internal log; K-ISMS-P control 13 incident registry | |

### 17.3 Verifier escrow procedure

| Escrow pattern | Description |
|---|---|
| Physical escrow | CD-ROM or USB drive, sealed envelope, stored in bank's vault with access controls (CISO + Audit chair approval to remove) |
| Cryptographic escrow | Verifier binary and SHA-256 hash signed by CISO using Ed25519 key in bank's own HSM (separate from chain's signing HSM) |
| Digital escrow | Verifier binary on air-gapped computer held by internal audit; powered off except during examination use |

### 17.4 Examination-time joint verification

When FSS examiner arrives, examiner and bank jointly retrieve escrowed verifier, verify SHA-256 against CISO's signed manifest, verify match against published release SHA-256, run verifier on ledger export.

### 17.5 Key recovery ceremony

Three-person ceremony per EFSR §17(3) for key recovery (per `bank-of-israel-overlay.md` §14.4 structure). Korean variant: CISO chairs; Chief of Operations executes; DBA Lead observes; Internal Audit Observer witnesses; Bank Executive (CFO or CRO) signs as ultimate approver. External Auditor optional but recommended for tier-1 banks.

---

## 18. AI Basic Act 2025 — Articles 15-20 high-impact AI obligations

### 18.1 Risk-stratified governance

The AI Basic Act 2025 (effective January 2026) introduces risk-stratified governance: 일반 (general) / 고영향 (high-impact) / 금지 (prohibited). Credit-scoring, lending, hiring, insurance underwriting, and law-enforcement AI are 고영향 (high-impact).

### 18.2 Article 15-20 obligations and chain mapping

| Article | Obligation | Chain role | Bank responsibility |
|---|---|---|---|
| 15 | Governance structure; AI committee | Chain is one technical control within bank's broader AI Governance Committee oversight; chain does not establish governance | Establish AI Governance Committee, name members, define escalation procedures |
| 16 | AI Impact Assessment before high-impact AI deployment | Chain captures data integrity; AI 영향평가 is a separate document. Chain's evidence (post-deployment fairness audits, customer-dispute logs) supports the bank's assessment | Commission AI Impact Assessment from in-house or external evaluator before deployment |
| 17 | Transparency and explainability to users; right to know AI decision | Customer-dispute-procedures.md operationalises explanation; chain provides audit trail (inputs, reasoning, output) | Disclose to customer at time of decision; maintain procedures for explanation requests |
| 18 | Human oversight and review | Chain captures human-review events; chain does not enforce human review | Define decision classes requiring human review; record review per `audit.routing.reviewer_identity` and parent-linkage |
| 19 | Fairness monitoring and discrimination prevention; regular audits | Chain's post-capture fairness audit consumes chain data | Commission fairness audits quarterly or semi-annually; integrate findings with model retraining |
| 20 | Post-market monitoring and incident reporting to MSIT | Incident-response-playbook.md covers FFIEC 36-hour notification; AI Basic Act §20 has separate trigger criteria | Establish incident-reporting criteria; notify MSIT per Article 20 timeline |

### 18.3 FSS AI Guidelines five governance pillars

| Pillar | Korean term | Chain support |
|---|---|---|
| Governance structure and board-level oversight | 지배구조 | Chain named in board policy; AI Governance Committee oversees chain operation |
| Explainability and transparency to customers | 설명가능성 | Chain captures decision inputs, parameters, reasoning chains (`audit.*` namespace) |
| Fairness and bias mitigation | 공정성 | Chain provides data integrity for fairness audits |
| Robust risk management | 위험관리 | Chain detects integrity failures; incident-response-playbook.md covers FSS notification |
| Cybersecurity and resilience | 사이버보안 | HMAC-SHA-256 + HSM controls satisfy K-ISMS-P baseline expectations per EFSR §17 |

### 18.4 AI Impact Assessment integration

| Phase | Activity | Chain role |
|---|---|---|
| Pre-deployment | Commission independent AI Impact Assessment; assess fairness risk, discrimination risk (especially region of origin, age, disability, marital status), user impact, transparency/explainability readiness, human-oversight capability, incident-response readiness; signed by chief risk officer, chief compliance officer, AI Governance Committee | Chain not yet operational |
| Deployment | System goes live; fairness baseline established | Chain begins capturing decision data |
| Post-deployment monitoring | Monthly or quarterly fairness audits; customer-dispute log; incident log | Chain provides per-decision data for audit; customer-interaction records for dispute analysis; system-event logs for incident investigation |
| Annual AI Impact Assessment review | Revisit fairness risk and discrimination risk; update assessment based on post-deployment evidence; re-certify or identify needed model retraining | Chain-derived fairness audit, customer-dispute summary, incident summary inform the annual update |

### 18.5 AI Governance Committee oversight

The bank's AI Governance Committee (per AI Basic Act §15) reviews quarterly:

- Operational health: verifier pass/fail rate; incident summary.
- Fairness audit results: demographic decision distributions for the period; identified bias; remediation actions.
- Incident reporting: any MSIT notifications under Article 20; PIPC complaints or investigations.
- Governance decisions: approval of new high-impact AI systems; decisions to retrain or adjust models; response to fairness-audit exceptions.

Committee minutes document decisions, satisfying AI Basic Act accountability obligation.

---

## 19. Right to explanation (설명요구권) and adverse-action workflow

### 19.1 What the right requires

AI Basic Act Article 17 grants users the right to demand explanation of an AI decision. The bank provides the explanation in plain language, in Korean, within 5 business days. The chain captures what the AI decided and why.

### 19.2 Customer-facing explanation workflow

| Step | Activity | Chain attribute |
|---|---|---|
| 1. Customer submits explanation request | "Explain why you denied my loan application" | Standard customer-service intake |
| 2. Bank's dispute team extracts chain entry | By customer ID, run_id, captured_at | Standard query |
| 3. Bank's AI explainability function translates `audit.*` fields | `audit.fraud_score=45, audit.fraud_threshold=40` becomes "귀하의 신용도 점수가 40% 기준을 초과하여 승인되지 않았습니다" | Institution's tooling |
| 4. Bank sends explanation within 5 business days | Plain-Korean notice | `audit.pipa.right_to_explanation.explanation_provided_at_utc`; `audit.pipa.right_to_explanation.explanation_text_hash` |
| 5. Customer disputes or requests further detail | Bank escalates to human-review process per AI Basic Act Article 18 | Triggers human-review pathway per §6.1 |

### 19.3 안내 의무 (duty of guidance) — proactive notification

The duty of guidance is stronger than GDPR's transparency requirement — it is affirmative guidance. When AI makes a credit-denial decision, the bank actively explains decision factors and offers a reconsideration path.

| Step | Activity | Chain attribute |
|---|---|---|
| 1. Generate notification within 2 business days of AI decision | Korean-language notice naming decision, date, factors (in language customer can understand), right to request human review (Article 37-2 이의제기권), reconsideration process | `audit.korean_guidance_duty.notification_event_id`; `audit.korean_guidance_duty.customer_contact_method` (SMS, email, phone call); `audit.korean_guidance_duty.notification_sent_at_utc`; `audit.korean_guidance_duty.customer_acknowledgment_received` (boolean); `audit.korean_guidance_duty.decision_factors_disclosed` (Korean-language text) |
| 2. Reconsideration path | If customer requests reconsideration | `audit.korean_guidance_duty.reconsideration_request_received_at`; `audit.korean_guidance_duty.human_review_assigned`; `audit.korean_guidance_duty.human_review_outcome` |

### 19.4 Korean adverse-action notice — chain attribute schema

For Korean adverse-action notices (per FSC fair-lending standards), use `audit.adverse_action.regulatory_basis = kr-fsc-fair-lending` per the schema in `eu-articulation-extension.md` §4.2.

| Notice content | Chain attribute |
|---|---|
| Decision and date | `audit.adverse_action.delivery_timestamp`; chain entry's `received_at` |
| Decision factors enumeration (extracted from `audit.*` fields, plain Korean) | `audit.adverse_action.decision_factors` (array of strings) |
| Right to reconsideration under Article 37-2 (이의제기권) | `audit.adverse_action.consumer_rights_disclosed` includes `right_to_human_review`, `right_to_object`, `right_to_explanation` |
| Reconsideration process and timeline | Institution's tooling; chain captures `audit.korean_adverse_action.reconsideration_requested` (boolean) |

---

## 20. Fairness audit for Korean protected characteristics

### 20.1 Korean protected characteristics

Korean law protects: gender, age, disability, marital status, region of origin (출신지역). Region-of-origin discrimination is a live concern reflecting Seoul vs provincial divides.

### 20.2 Fairness-audit procedure

The institution maintains a separate, governance-controlled de-identification mapping (not on the chain, not subject to chain retention). Fairness-audit procedure:

| Step | Activity |
|---|---|
| 1. Compliance/MRM team requests fairness report for [period, AI system] | Quarterly or semi-annual cadence |
| 2. Data team uses de-identification mapping to stratify chain-captured decisions by demographic cohort | Age groups <25 / 25-40 / 40-65 / >65; Seoul / 경기도 / 인천 / 강원도 / 대전 / 세종 / 충청남도 / 충청북도 / 전주 / 전라남도 / 경상남도 / 경상북도 / 부산 / 울산 / 제주도; with/without disability; marital status |
| 3. Report shows decision distribution per cohort, differential impact analysis (statistical significance), root-cause analysis | Approval rate; chi-square test for significance; data bias vs model bias diagnosis |
| 4. If discrimination detected, AI Governance Committee directs remediation | Retrain model; adjust thresholds; add fairness constraints; add human review for affected cohort |
| 5. Document finding and remediation in AI Impact Assessment; report to MSIT if required | AI Basic Act §20 incident reporting where applicable |

### 20.3 Region-of-origin (출신지역) testing

If Seoul's approval rate is 70 percent and a province's rate is 60 percent, chi-square test confirms whether the difference is statistically significant. If statistically significant: data bias (region X has worse credit profiles) vs model bias (model approves region X at lower rates with similar profiles).

### 20.4 Age-group testing

Approval rate for age groups <25, 25-40, 40-65, >65. For >65, special scrutiny — older adults are sometimes more vulnerable; bank policy may exclude >65 from fully automated decisions, requiring human review.

### 20.5 Disability-status testing

Disability-status information is sensitive; tokenised before chain capture; mapping kept separate. Bank's AI model development excludes disability status and related proxies (accessibility accommodation requests, medical claims, healthcare provider names). Fairness audit tests for differential impact on customers with known disability status; remediation involves model retraining, human review, or feature exclusion.

### 20.6 Marital status testing

Approval rate for single vs married applicants (controlling for income, credit history, debt-to-income ratio). If married applicants approved at significantly higher rates, model has learned to discriminate. Remediation: retrain with marital-fairness constraint; add human review for single-applicant edge cases.

---

## 21. AI Trustworthiness Verification (안전한 AI 인증)

### 21.1 Voluntary certification context

The Trusted AI Certification program, administered by KISA under MSIT, provides voluntary certification that AI systems meet specified safety, fairness, and explainability standards. Voluntary but increasingly expected by FSC for institutions in the regulatory sandbox.

### 21.2 Chain as supporting evidence

| Certification criterion | Chain support |
|---|---|
| Governance and oversight | Chain named in bank's AI Governance Committee charter and board-level AI policy |
| Auditability and transparency | Chain's verifier (spec §7) is independent, reproducible audit mechanism |
| Fairness and non-discrimination monitoring | Bank's fairness-audit procedures consume chain data |
| Security and integrity | Chain's HMAC + Merkle + HSM design satisfies cryptographic security expectations |

Institutions pursuing certification reference the chain implementation in their submission and ensure verifier-friendly and audit-evidence-friendly deployment.

### 21.3 Social legitimacy (사회적 정당성)

Korean governance places weight on social legitimacy — AI systems must be perceived as legitimate, fair, and trustworthy by the broader public, not just legally compliant. The bank's communication strategy:

- Customer-facing disclosure at decision time: "Your application was evaluated by AI. You have the right to ask for an explanation and to request human review."
- AI governance page on the bank's website explaining AI Governance Committee, AI Impact Assessment, fairness monitoring, and chain-of-custody audit trail.
- Annual AI governance report (voluntary, recommended) publishing fairness audit results, incident summary, and committee meeting topics.
- Engagement with civil-society groups on fairness procedures and remediation.
- Media communication: spokesperson can confidently cite the chain as evidence of integrity and oversight.

---

## 22. Translation table — Korean regulatory expectations to chain artifacts

| Spec section | PIPA | CIUPA | EFTA / EFSR | AI Basic Act | K-ISMS-P | Conformance status |
|---|---|---|---|---|---|---|
| §4.1 Per-event MAC | §29 (security measures) | §16 (security) | EFSR §17 (cryptographic) | §17 (transparency) | 3.2 | CONFORMANT |
| §4.2 Daily Merkle seal | §29 | §16 | EFSR §17 (audit log) | §20 (incident reporting) | 9.1, 9.3 | CONFORMANT |
| §4.3 HSM signature | §29 | §16 | EFSR §17 | §17 | 3.2 | CONFORMANT |
| §4.4.1 Routing decisions | §37-2 (right to refuse / right to explain) | §15 (transparency) | EFSR §17 (audit log) | §17 (transparency); §18 (human oversight) | 9.1 | PARTIAL (institution's CC8.1 names routing-policy versioning) |
| §4.4.2 Deployment intent | §37-2 | §15 | EFSR §17 | §16 (AI Impact Assessment); §18 | 9.1 | PARTIAL (institution's CC8.1 names deployment-policy versioning) |
| §6 Deployment topologies | §29; §26 (cross-border) | §16 | EFSR §17; 망분리 | §15 (governance) | 2.10 | CONFORMANT (망분리 may require on-premises ledger; FSS approval for cross-boundary) |
| §7 Verifier | §29 | §16 | EFSR §17 | §16; §17; §19 | 9.1 | CONFORMANT |
| §10.5 HSM custody | §29 | §16 | EFSR §17 | §17 | 3.2 | CONFORMANT (FIPS 140-2 Level 3 or K-ISMS-P certified KMS) |
| §10.9 Retention | §29 | §16 | FEFTA 7-year minimum; FIEA 5-year minimum | §20 (incident reporting) | 9.3 | CONFORMANT |
| §10.11 Adverse-action notice | §37-2 | §15 (explainability) | EFSR §17 | §17 (transparency) | 9.1 | PARTIAL (institution uses `audit.adverse_action.regulatory_basis = kr-fsc-fair-lending` and `notice_language = ko-KR`) |
| §10.14 Trusted-time integration | §29 | §16 | EFSR §17 | §17 | 9.3 | PARTIAL (institution adopts RFC 3161 with a TSA in Korean jurisdiction or KISA-recognised TSA) |
| §10.15 Multi-region resilience | §29 | §16 | EFSR §17(1); BCP | §15 | 2.10 (BCP) | CONFORMANT (Seoul/Busan active-active recommended; cross-affiliate isolation in conglomerate context per §5.2) |

---

## 23. Operational checklist for a Korean institution deploying the chain

| Artefact | Source | Reviewed by |
|---|---|---|
| FSS comprehensive examination evidence package (per §3.2) | Institution's compliance | FSS examiner |
| FBO dual-supervision evidence (if applicable) | Institution's CC8.1 | FSS + Federal Reserve |
| Group-level AI policy (for conglomerate-affiliated banks) | Institution's board + group AI Governance Committee | FSS consolidated supervisor |
| PIPA Article 37-2 right-to-refuse and right-to-explain procedures | Institution's customer-service + CC8.1 | PIPC + FSS |
| PIPA Article 26 cross-border posture | Institution's CC8.1 + RoPA + transfer-impact assessment (개인정보 이전영향평가) | PIPC |
| PIPA Article 28-2 pseudonymisation declaration | Institution's CC8.1 | PIPC + FSS |
| MyData transmission procedure | Institution's CC8.1 | FSC + MyData intermediary |
| CIUPA Article 15 fairness-testing program | Institution's MRM + CC8.1 | FSS + KCISA |
| EFSR §17 dual-approval procedures (per §10.2) | Institution's CC8.1 | FSS examiner |
| 망분리 architecture diagram + FSS approval (if cross-boundary) | Institution's CC8.1 | FSS |
| K-ISMS-P control implementation guide for controls 2.10, 3.2, 9.1-9.3, 13 | Institution's CC8.1 | K-ISMS-P certification audit |
| Quarterly personal-data-flow documentation (per §12.6) | Institution's data-governance team + CC8.1 | K-ISMS-P auditor |
| KISA Cybersecurity Notification System procedure (2-hour clock) | Institution's IR runbook | KISA |
| FSI Banking Cyber CERT integration procedure | Institution's CC8.1 | FSI |
| Verifier escrow procedure | Institution's CC8.1 + Internal Audit | Internal Audit |
| EFSR §17(8) biennial red-team exercise readiness | Institution's CC8.1 | FSS + (where applicable) NIS |
| AI Basic Act §16 AI Impact Assessment | Institution's risk + AI ethics + CC8.1 | MSIT + AI Governance Committee |
| AI Basic Act §15 AI Governance Committee charter and meeting cadence | Institution's board | MSIT |
| AI Basic Act §17 right-to-explanation procedure | Institution's customer-service + CC8.1 | PIPC + MSIT |
| AI Basic Act §19 fairness-audit program (Korean protected characteristics) | Institution's MRM | MSIT + AI Governance Committee |
| AI Basic Act §20 MSIT incident reporting procedure | Institution's IR runbook | MSIT |
| AI Trustworthiness Verification submission (if pursued) | Institution's submission to KISA | KISA + MSIT |

---

## 24. Bottom line

The chain v1.0a is conformant with Korean financial-sector supervisory architecture and operationally defensible under the Korean nation-state threat model when the institution adopts the compensating controls articulated in this overlay. The conformance is operational: the institution's CC8.1 names the seal cadence, the deployment topology (Korean jurisdiction with 망분리 compliance or FSS-approved cross-boundary), the multi-regulator IR procedure (KISA + FSS + FSI + PIPC + MSIT + Federal Reserve for FBOs), the K-ISMS-P control mapping, the AI Basic Act compliance posture, and the fairness-audit program for Korean protected characteristics.

The chain's integrity primitives are NIST/IETF-standardised and Korea-conformance-friendly: HMAC-SHA-256, HKDF-SHA-256, Ed25519, RFC 8785, RFC 6962 Merkle. Korean K-ISMS-P certified KMS is acceptable equivalence for spec §10.5; FIPS 140-2 Level 3 also applies. None of the primitives require modification for Korean deployment; what differs is institutional CC8.1 documentation, per-jurisdiction tenant naming (`bank-kr` for Korea-jurisdiction tenants per `apac-overlay.md` §11.2), and the multi-regulator coordination procedure.

For Korean megabanks under nation-state threat baseline, the chain produces detection and proof. The chain bounds tampering to single-day windows (default daily seal) or single-hour windows (hourly seal under EFSR §17(2) real-time monitoring composition with SIEM); the institutional discipline articulated in §10 through §17 produces the full operational and cyber-resilience posture under FSS, PIPC, KISA, and MSIT oversight.

For institutions pursuing AI Trustworthiness Verification or operating under the AI Basic Act 2025, the chain is strong supporting evidence for governance, auditability, fairness, and security criteria. Korean institutions deploying the chain should engage with MSIT and KISA early in the deployment cycle to confirm the chain's role in their certification or sandbox participation.

For institutions evaluating the chain for the first time under Korean regulatory mandate, the recommended posture is: self-host the ledger and verifier in a Korean data center or colocation (Seoul primary, Busan secondary for BCP); adopt hourly seals or active-active failover for EFSR §17 real-time monitoring composition; integrate with KISA per §13; conduct annual vendor audits; prepare for FSS comprehensive examination per §3 and §15; document the right-to-explanation and 안내 의무 workflows per §19; deploy fairness-audit procedures for Korean protected characteristics per §20; document AI Basic Act compliance per §18. With these compensating controls, the chain materially strengthens the institution's AI-governance and cyber-resilience posture under Korean supervision.

For multi-jurisdiction Korean institutions also operating in other APAC jurisdictions, this overlay composes with `apac-overlay.md` (Pattern B per-jurisdiction tenant naming, APEC CBPR cross-jurisdiction posture). For Korean megabanks operating FBOs in the US, this overlay composes with the FFIEC examination procedures (audit-procedures.md, examination-response-workflow.md). For Korean institutions also operating in EU jurisdictions, this overlay composes with `dora-articulation-overlay.md` and `eu-articulation-extension.md`.

함께 신뢰받는 AI를 만들어 갑시다.

---

## Document control

| Field | Value |
|---|---|
| Document | Korea Articulation Overlay |
| Version | 1.0.0 |
| Status | Informative — institution-side articulation; no normative spec change |
| Aligned with | PIPA (March 2023) / CIUPA / EFTA / EFSR §17 / MyData / AI Basic Act 2025 / FSS AI Guidelines / K-ISMS-P / 망분리 / KISA Notification System / MSIT AI Ethics Standards / AI Trustworthiness Verification System / Korean fair-lending standards |
| Date | 2026-05-07 |
