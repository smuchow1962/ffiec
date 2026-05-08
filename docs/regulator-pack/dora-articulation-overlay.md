---
title: DORA / EBA / eIDAS / NIS2 Articulation Overlay
status: informative
aligned-with:
  - Regulation (EU) 2022/2554 (DORA)
  - EBA Guidelines on outsourcing arrangements (EBA/GL/2019/02)
  - Regulation (EU) 910/2014 (eIDAS) and Regulation (EU) 2024/1183 (eIDAS 2.0)
  - Directive (EU) 2022/2555 (NIS2)
  - Regulation (EU) 2016/679 (GDPR)
date: 2026-05-07
version: 1.0.0
---

# DORA / EBA / eIDAS / NIS2 Articulation Overlay

> **What this doc is.** A single articulation overlay that maps the FFIEC chain-of-custody v1.0a specification onto the European supervisory framework. Written so a competent authority — a national prudential supervisor under the Single Supervisory Mechanism (SSM), an ESA staff member, an EBA Joint Examination Team lead, or a national CSIRT contact under NIS2 — can read this document alongside the spec and confirm what the chain delivers, what it does not, and where supplementary measures are required for European application.

> **What this doc is NOT.** Not a normative extension to the v1.0a specification. Not a translation of the spec into European-framework language. Not an EU-specific spec fork. The integrity primitives in v1.0a are framework-neutral cryptographic constructs; the supervisory translation is an articulation overlay an institution layers on top, not a change to the underlying specification. Spec §1.2 (epistemic scope) is the discipline this document operates under: the chain proves what was said and that the record was not tampered with after capture; it does not prove the substantive correctness of the captured content. The overlay carries the same discipline — it maps integrity evidence to supervisory expectations; it does not claim the chain alone discharges the institution's substantive obligations under any of the regimes named here.

---

## 1. Scope and reading order

This overlay is consumed in three reading orders.

| Reader | Reading order |
|---|---|
| First-look EU competent authority (SREP / Joint Examination Team) | §2 (DORA controller/processor lattice) → §4 (DORA incident reporting clocks) → §14 (translation table) → §15 (bottom line) |
| Institution's DORA-readiness team | §2 → §3 (DORA Article 30(2) checklist) → §6 (RTO/RPO under multi-region) → §10 (concentration risk) |
| Institution's eIDAS / electronic-signature counsel | §8 (eIDAS qualification) |
| National CSIRT / NIS2 reporting team | §9 (NIS2 timing) → §4 (DORA lex specialis reconciliation) |
| Privacy office under cross-border deployment | §11 (Schrems II supplementary measures specific to the chain) → §12 (right-to-explanation under GDPR Article 22 / AI Act Article 86) → cross-references to `international-transfers.md` and `multi-jurisdiction-conflict.md` |
| Data-subject-rights team / DPO | §12 (right to explanation) |

The overlay does not duplicate material that already lives in adjacent regulator-pack documents. Cross-references are exact — if a section sends the reader to `international-transfers.md` §2.2, that section is the load-bearing source and this overlay names where it sits in the European-framework picture.

---

## 2. DORA Articles 28-30 — ICT third-party risk and the controller / processor lattice

The Digital Operational Resilience Act, Regulation (EU) 2022/2554 (hereafter DORA), Articles 28-30 govern the contractual posture between a Union financial entity and any provider of ICT services that supports critical or important functions. The chain-of-custody deployment necessarily places multiple parties in this lattice, and the institution's register of information under Article 28(3) must name each one. This section maps the v1.0a deployment shapes onto the DORA categories. *(addresses Lindqvist G-1, P-2)*

### 2.1 The four parties in a typical chain deployment

The chain is operated by some combination of four parties. A given deployment may collapse two or more of them into a single legal entity, but the lattice is the same.

| Party | Role in the chain | DORA category | GDPR category (cross-reference) |
|---|---|---|---|
| The institution | The Union financial entity operating the AI agent and consuming the chain's evidence | Financial entity (subject of DORA) | Controller |
| SDK provider | Authors and ships the chain SDK that the institution embeds in its application process | ICT third-party service provider supporting an important function under Article 3(21) | Processor (in vendor-shipped SDK; sub-processor where the SDK is shipped via another vendor) |
| Receiver / ledger provider | Operates the OTLP receiver, the ledger storage, the seal job, and the verifier endpoint | ICT third-party service provider supporting a critical function under Article 3(21) | Processor |
| HSM provider (cloud or on-prem) | Custodies the seal-signing key under FIPS 140-2 Level 3 or higher per spec §10.5 | ICT third-party service provider for cryptographic services | Sub-processor |
| LLM provider (Anthropic, OpenAI, Google, Mistral, etc.) | Provides the model the institution invokes; receives the prompt and returns the response captured by the chain | ICT third-party service provider supporting an important function under Article 3(21); a sub-contractor to the institution's AI workload under Article 30(2)(a) | Processor (where the institution sends personal data); sub-processor where the institution accesses through a SaaS aggregator |

The lattice is wider than the GDPR Article 28 lattice for two reasons. First, DORA imposes its register-of-information obligation regardless of whether personal data is transferred — an LLM provider that processes only synthetic prompts is still in scope. Second, DORA names "ICT services that support critical or important functions" as the classification trigger, and an integrity-bearing audit trail of automated decisions clearly qualifies for any decision class with regulatory consequence (credit decisions, fraud screening, KYC, payment routing).

### 2.2 The flows that warrant supervisory attention

Three flows in the lattice warrant explicit articulation in the institution's Article 28(3) register and CC8.1 control description.

| Flow | Concern | Required articulation |
|---|---|---|
| Application process → LLM provider | The prompt — typically containing the customer's matter — is sent to the LLM. The chain captures what was sent (per spec §4.4 `gen_ai.*` attributes) but the institution's register must name the LLM provider as a sub-contractor under Article 30(2)(a) and document the contractual flow-down of DORA obligations to that provider | Article 30(2) clauses applied to the LLM contract; institution's RoPA names the LLM provider as a sub-processor for the corresponding personal-data category |
| SDK process → receiver / ledger | Chain entries cross the wire (OTLP per §4.4) carrying integrity-bearing evidence and, depending on payload-redaction posture, possibly customer content | The receiver provider's contract names DORA Article 30(2) clauses, the EBA Guidelines on outsourcing audit-rights clause (per §7 below), and the institution's exit-strategy commitments per Article 28(8) |
| Receiver → HSM | The seal job dispatches signing operations to the HSM. The seal-signing key is non-extractable per spec §10.5; nevertheless the HSM provider remains a sub-processor for cryptographic services and a critical-tier ICT third-party | The HSM provider's contract names FIPS 140-2 Level 3 (or higher) custody, the separation-of-duties roster per spec §10.5, and the EBA-specified competent-authority audit rights |

The chain does NOT change the contractual obligations the institution holds with each provider; the chain produces the audit trail that lets the institution, its competent authority, and any examiner appointed by the authority verify that the contractual obligations are observed in operation.

### 2.3 The LLM provider as a sub-contractor — the central novelty

The LLM provider's position in the lattice is the central novelty of this workload. The articulation here is normative for the institution's register, not for the spec. The LLM provider is in scope of Article 30 contractual provisions even when accessed by API and even when the institution does not transfer Article 4(1) personal data — DORA's "supporting ICT services" classification under Article 3(21) applies to the operational dependency, not to the personal-data classification.

A practical corollary: the institution's contract with the LLM provider must include the Article 30(2)(a) through (h) sub-clauses (named in §3 below) regardless of whether the institution's GDPR analysis classifies the LLM provider as a processor. The two regimes are layered, not substitutable. The chain's `gen_ai.response.model` field (spec §4.4 normative) records the model identifier the vendor's API actually answered with — the load-bearing record for the institution's register-of-information entry naming the LLM provider, the model identifier, and the version. The chain's `audit.deployment.intent` field (spec §4.4.2) records whether the invocation was steady-state production, an A/B test, a canary, a multi-region deployment, or a vendor-side reroute observation; this metadata distinguishes deliberate vendor-managed model rollouts from silent vendor reroutes the institution did not authorize, which is the load-bearing distinction for the institution's vendor-management procedure under Article 28(2).

### 2.4 The SDK provider — embedded vs SaaS

A practical distinction the institution's register must articulate: whether the SDK provider ships an embedded library that runs in the institution's application process, or a SaaS that the institution invokes over a network boundary.

| SDK delivery model | DORA classification | Practical implication |
|---|---|---|
| Embedded library (the institution's application calls in-process; per spec §4.1 the chain construction occurs at the institution's host) | The SDK provider is an ICT third-party service provider whose service is the SDK code itself; the institution's data does not leave the host under SDK operation | Article 30(2)(b) location-of-processing names the institution's host environment; the SDK provider has no operational access to chain content |
| SaaS (the institution's application calls the SDK over a network boundary, and chain construction occurs at the SaaS provider's infrastructure) | The SDK provider is an ICT third-party service provider with operational access to chain content; processor under GDPR Article 28 | Article 30(2)(b) names the SaaS provider's data-residency posture; the SaaS provider's contract names the EBA/GL/2019/02 §13.2 audit-rights clause |

The institution's CC8.1 names which SDK delivery model it operates and the corresponding contractual posture. The chain's `service.name` and `service.version` Resource attributes (spec §4.4.3) identify the SDK build; the institution's evidentiary-artifacts retention (spec §10.13) preserves the SDK source-code hash and build reproducibility evidence so the JST examiner can verify the SDK in production matches the SDK named in the contract.

---

## 3. DORA Article 30(2) — contractual elements checklist

Article 30(2) of DORA enumerates the elements every contract for ICT services supporting critical or important functions must contain. The chain-of-custody outputs serve as the audit trail that lets the institution, its competent authority, and any appointed auditor verify the contract is observed in operation. This checklist names each sub-clause and the chain artifact that supports it. *(addresses Lindqvist G-1, P-2)*

| Article 30(2) sub-clause | Chain-of-custody evidence supporting verification |
|---|---|
| (a) Description of services and sub-contracting chain | The institution's register names each provider in §2.1 above; the SDK's `service.name` and `service.version` Resource attributes (spec §4.4.3) identify the in-process SDK; the seal record's `kms_handle_uri` identifies the HSM custody |
| (b) Locations where ICT services are provided and where data is processed | The OPTIONAL `ffiec.chain.region` attribute (spec §4.4) records the per-event capture region under MAC binding; the seal region is named in the institution's CC8.1 per spec §10.15 Pattern A; the LLM provider's processing location is contract-specified and the institution's RoPA records it |
| (c) Provisions on data availability, authenticity, integrity, confidentiality | Integrity is delivered by the three-layer construction (per-event MAC §4.1, daily Merkle seal §4.2, HSM Ed25519 signature §4.3); availability and resilience by the multi-region patterns of spec §10.15 (cross-reference §6 below); authenticity by the HSM-bound seal signature; confidentiality by the institution's encryption-at-rest posture documented in `article-32-security-mapping.md` |
| (d) Provisions on data return and deletion at end of service | The institution's exit strategy names the chain-data return format (NDJSON audit files per spec §6) and the IKM-retention coupling per spec §10.9 |
| (e) Service-level descriptions and quantitative performance targets | The institution's SLA names the receiver's seal-cadence delivery target per spec §4.2.1 (e.g., 60-minute publish window after seal-window end) and the verifier's exit-code contract per spec §10.12 |
| (f) Right of access, inspection, and audit by the institution and the competent authority | Spec §7's offline, no-network-call verifier is the load-bearing inspection mechanism; the EBA Guidelines on outsourcing §13 audit-rights overlay is in §7 of this overlay |
| (g) Obligation to assist with incident reporting | The institution's IR runbook produces evidence under spec §10.2 operational events (`chain.verification_failure`, `audit_file.truncation_detected`, `hsm.operation_failure`, etc.); the receiver provider's contract names the cooperation timeline aligned with the DORA 4-hour clock (cross-reference §4 below) |
| (h) Right to terminate with reasonable notice in cases of supervisory breach | The institution's contract names the termination-on-supervisory-breach clause; the chain's evidentiary-artifact retention per spec §10.13 supports the termination-handover record |
| (i) Cooperation with competent authorities (Article 30(2)(d) read with Article 30(3)) | The institution's CC8.1 names the procedure for surfacing the chain to a competent authority; the verifier's deterministic, single-binary, no-network posture (spec §7) is what makes this practicable at examination time |
| (j) Description of how the provider supports the institution's TLPT obligations | Cross-reference §5 below; the institution's contract names the provider's TLPT-cooperation commitment |
| (k) Notification of material changes | The institution's CC8.1 names the change-notification channel; spec §10.10 (rotation crossing the seal boundary) and §10.10.2 (within-day algorithm rotation) cover the in-band material changes the chain itself produces evidence for |
| (l) Insurance and indemnity provisions | Contract-specific; the chain does not produce evidence on this clause |
| (m) Access by competent authorities to the provider's premises | The EBA Guidelines on outsourcing §13.2 audit-rights clause (per §7 below) is the operational expression of this requirement |
| (n) Provider's commitment to security awareness, training, and ICT-risk management | The provider's SOC 2 Type II report and its FIPS validation certificates substantiate this clause; the institution's vendor-management procedure consumes them |

This checklist is the contractual side. The institution pairs it with the EBA outsourcing audit-rights template — once `docs/templates/eba-outsourcing-audit-rights.md` is published, this overlay cross-references that template for the load-bearing audit-rights clause text. Until that template is published, the institution's legal team drafts the audit-rights clause directly from EBA/GL/2019/02 §13.2 and references this overlay's §7 for the regulatory backing.

### 3.1 Three contractual postures the institution typically operates

A typical chain deployment carries three contractual instruments. The institution's register names which provider holds which instrument and how the chain's evidence supports the contract's enforcement.

| Instrument | Counterparty | Chain evidence |
|---|---|---|
| ICT services agreement under DORA Article 30 | Receiver / ledger provider | Verifier output (per spec §10.13) demonstrates the integrity SLA; operational events demonstrate the seal cadence delivery; receiver-policy minimum surface (per §10.4 of this overlay) demonstrates the audit-rights URI |
| ICT services agreement under DORA Article 30 (separate instrument) | LLM provider | `gen_ai.response.model` field demonstrates the model identifier the vendor's API actually answered with; `audit.deployment.intent` field distinguishes deliberate model rollouts from silent reroutes |
| Cryptographic-services agreement | HSM provider | `kms_handle_uri` field on every chain entry; HSM-signed seal records demonstrating the signing-key custody; `hsm.operation_failure` operational events demonstrating availability |

The institution's register-of-information entry under Article 28(3) names each provider, the contractual instrument, the institution's exit-strategy commitments, and the cooperation timeline aligned with the DORA 4-hour clock per §4.1 of this overlay.

---

## 4. DORA Articles 11-14 — incident classification and reporting

DORA Articles 11-14, taken with the Commission Delegated Regulation (EU) 2024/1772 (RTS on incident classification) and the Commission Implementing Regulation (EU) 2024/2956 (templates for major-incident reporting), set the incident-reporting clocks for Union financial entities. The DORA RTS finalized in 2024 sets the initial-notification clock at **4 hours** from major-incident classification, the intermediate-notification clock at 72 hours, and the final-report clock at 1 month. The earlier 24-hour figure that appeared in some pre-RTS materials is superseded; the 4-hour clock is the binding initial-notification clock for Union financial entities. *(addresses Lindqvist G-3, G-4)*

### 4.1 The 4-hour initial-notification clock and the seal cadence

The chain's daily seal cadence (spec §4.2.1) is too coarse for 4-hour initial reporting. A daily seal that publishes 60 minutes after midnight UTC delivers the day's integrity-bearing evidence on a 25-hour rolling cycle, which an institution cannot use to produce the 4-hour initial-notification dossier mechanically. The resolution has two components.

**Higher-frequency seal cadence under DORA-bound deployments.** Spec §4.2.1 names hourly cadence as a supported configuration. An institution operating under DORA SHOULD select hourly cadence for any tenant whose workload is in scope of major-incident classification. Hourly cadence binds the seal to the seal-payload at hourly granularity, so within any 4-hour window the institution holds three to four sealed-and-integrity-bound seal records covering the incident window. The institution's CC8.1 control description names the cadence binding and the seal-publish-latency target consistent with the 4-hour clock. Spec §10.10.1 documents the operational behavior of hourly cadence under master-key rotation and is the reference for institutions adopting this posture.

**Partial-disclosure verifier mode for the initial-notification dossier.** The chain's offline, single-binary verifier (spec §7) operates on the persisted seal-and-events. For 4-hour initial-notification, the institution produces a partial-disclosure verifier output covering the events the institution has classified into the incident scope at the moment of initial notification. The witness-verifier mode (spec §7) is the load-bearing artifact when the institution has not yet been authorized to disclose IKM-scope evidence to the supervisor; the witness-mode output establishes structural integrity (per-event walk, Merkle root, HSM signature) without exposing the IKM. The institution's IR runbook names the partial-disclosure procedure and the witness-verifier output as the initial-notification evidence shape.

The institution's IR commander operates the 4-hour clock as the binding clock for any DORA-classified major incident. The chain's evidence is assembled inside the 4-hour window using the partial-disclosure procedure; the intermediate-notification (72 hours) and final-report (1 month) evidence is assembled from the institution's standard verifier output once the seal-publish cycle has completed for all affected seal-windows.

**Worked example — 4-hour clock operating against an integrity anomaly.** The institution detects a `chain.verification_failure` operational event at 09:14 UTC on a Tuesday. The IR commander opens the incident at 09:18 UTC. Classification is committed at 09:42 UTC under the institution's CRITICAL/HIGH/MEDIUM/LOW policy: HIGH (a chain-integrity anomaly affecting credit-decisioning evidence). The 4-hour clock starts at 09:42 UTC. The institution's IR runbook produces:

- **At 09:48 UTC** — the tenant-day's most recent hourly seal record (08:00-09:00 UTC sealed, with `chain.verification_failure` entries inside it) is exported under witness-mode verifier
- **At 10:15 UTC** — the witness-verifier output for the affected hour is bound to the institution's incident classification dossier under spec §10.2 operational events; the dossier carries the `incident.classification.*` attribute set per §4.2 of this overlay
- **At 12:30 UTC** — the dossier is delivered to the institution's competent authority on the DORA reporting template (Implementing Regulation (EU) 2024/2956)
- **At 13:42 UTC** — the 4-hour clock closes; the institution has discharged the initial-notification obligation

The institution's CC8.1 names the cadence binding (hourly), the partial-disclosure procedure, and the dossier-assembly target latency. The 4-hour budget is comfortably met because the chain's evidence is bound under integrity at hourly granularity and the verifier produces deterministic output in single-binary execution time bounded by the seal's event count.

### 4.2 RTS Article 18 incident-classification fields

The RTS Article 18 (Commission Delegated Regulation (EU) 2024/1772) requires the institution's classification dossier to carry structured fields for impact, duration, geographic spread, and data-loss extent. The chain's primary outputs do not produce these fields directly — the chain proves what was said and when. The institution's IR system produces the classification dossier and consumes the chain's evidence. *(addresses Lindqvist G-3)*

The institution's IR system operates a supplementary attribute set that the chain integrity-binds when the IR system emits chain entries for incident-related operational events. The recommended attribute set is named below; the chain integrity-binds whatever the institution emits, so the field set is institution-side and outside the spec's normative scope.

| Attribute | Purpose | RTS Article 18 mapping |
|---|---|---|
| `incident.classification.affected_jurisdictions` | List of EEA Member States whose data subjects, customers, or operations are affected | RTS Article 18 geographic-spread criterion |
| `incident.classification.continuous_duration_hours` | Continuous duration of service degradation in hours | RTS Article 18 duration criterion |
| `incident.classification.criticality` | One of `CRITICAL` / `HIGH` / `MEDIUM` / `LOW` per the institution's CRITICAL/HIGH/MEDIUM/LOW classification policy | RTS Article 18 criticality of services affected |
| `incident.classification.data_loss_extent` | Estimated count of records or data-subjects whose data is impaired | RTS Article 18 data-loss-extent criterion |
| `incident.classification.economic_impact_band` | One of the RTS-defined economic-impact bands | RTS Article 18 economic-impact criterion |
| `incident.classification.recovery_objective_state` | One of `WITHIN_RTO` / `RTO_BREACHED` / `RECOVERED` | RTS Article 18 recovery-objective criterion |

The IR system emits these attributes on chain entries representing operational events under spec §10.2 (e.g., `incident.opened`, `incident.classification_committed`, `incident.recovery_committed`). The chain integrity-binds the classification dossier under the daily Merkle seal, so the supervisor's working paper carries an integrity-bound classification record alongside the integrity-bound captured-event evidence.

### 4.3 Lex specialis reconciliation

DORA Article 1(2) makes DORA the lex specialis for Union financial entities — where DORA imposes a stricter or more specific obligation than another EU regime, DORA prevails. The institution's IR commander operates the **tightest applicable clock first** under the simultaneous-notification rule. *(addresses Lindqvist P-3)*

| Regime | Initial / first-notification clock | Position vs DORA |
|---|---|---|
| FFIEC (US, where the institution is dual-supervised) | 36 hours | Looser than DORA; DORA prevails for the EU operation |
| DORA (Article 19) | 4 hours | Tightest applicable clock for the EU financial entity |
| GDPR (Article 33) | 72 hours | DORA's 4-hour clock prevails for major incidents that are also personal-data breaches under DORA-classification scope |
| HIPAA (US, healthcare-financed lending) | 60 days | Looser than DORA; DORA prevails for the EU operation |
| NIS2 (Article 23, where institution is also an essential entity) | 24 hours early warning + 72 hours notification | DORA Article 1(2) lex specialis: DORA prevails for Union financial entities; institution may still notify the national CSIRT under NIS2 in parallel where its national transposition requires |

The simultaneous-notification rule from `breach-notification-matrix.md` §3 and `multi-jurisdiction-conflict.md` §6 governs the operational sequencing — the IR commander notifies the tightest deadline first, then the next-tightest, with a consistent narrative across all notifications.

### 4.4 The DORA reporting flow

```mermaid
sequenceDiagram
    participant det as Detection layer
    participant ir as IR commander
    participant chain as Chain hourly seal
    participant ver as Witness verifier
    participant ca as Competent authority

    det-&gt;&gt;ir: chain.verification_failure operational event
    ir-&gt;&gt;ir: open incident, classify under RTS Article 18
    Note over ir: 4-hour clock starts at classification commit
    ir-&gt;&gt;chain: request affected hourly seals
    chain-&gt;&gt;ver: persisted seal-and-events
    ver-&gt;&gt;ir: PASS-STRUCTURALLY witness output
    ir-&gt;&gt;ca: initial notification under Implementing Reg 2024/2956
    Note over ir,ca: Inside 4-hour window
    ir-&gt;&gt;ca: intermediate notification at 72 hours
    ir-&gt;&gt;ca: final report at 1 month
```

The flow assumes hourly cadence per spec §4.2.1. Daily-cadence deployments cannot meet the 4-hour clock mechanically; the institution's cadence selection is the load-bearing operational decision under DORA-bound deployments.

---

## 5. DORA Articles 26-27 — threat-led penetration testing (TLPT)

DORA Articles 26-27 mandate threat-led penetration testing for significant Union financial entities every three years. The TIBER-EU framework, originally a voluntary methodology of the European Central Bank, is the operational reference for TLPT execution; CBEST methodology (Bank of England) is the recognized cross-border equivalent. The chain MUST record red-team activity as audit events without leaking attack TTPs into the production retention. *(addresses Lindqvist G-5)*

### 5.1 The TLPT problem — chain entries become attacker playbook

A red-team engagement against an AI agent will produce thousands of chain entries documenting attack patterns: prompt-injection attempts, tool-call probing, attempts to subvert policy guardrails, attempts to elicit unauthorized actions. If those entries are retained at production retention, they become an information-leak path: a future insider, a future legal-hold-driven discovery, a future breach of the ledger storage all expose the institution's attack surface to a recipient who has no need to know.

### 5.2 The conformant resolution — segregated retention via attribute-flagged routing

The chain captures every event under integrity binding regardless of the event's purpose. The TLPT segregation operates at the receiver layer, not at the SDK layer.

**Recommended attribute pattern.** The institution's red-team engagement runs under a `tenant.testing_mode = "tlpt"` Resource attribute (per spec §4.4.3 OTLP transport identification) AND a corresponding `audit.tlpt.engagement_id` per-event attribute. The receiver's policy routes traffic carrying these attributes to a sealed-but-segregated red-team chain — same integrity construction, same verifier procedure, separate storage with separate access control and a separate retention window aligned with the TLPT engagement's documented test-data-handling policy. The SDK does not need to know it is operating under TLPT mode; the institution's deployment runs the red-team engagement against a labeled tenant whose receiver routes accordingly.

**Integrity binding under TLPT.** The TLPT chain is integrity-bound the same way as any production chain. The TLPT lead AND the institution's blue team can both verify the chain offline using the verifier without exposing the engagement's content to the broader operations team — the witness-verifier mode (spec §7) confirms structural integrity without IKM access for parties whose role does not extend to key custody.

**Retention discipline.** The TLPT chain's retention is named in the institution's CC8.1 control description, aligned with TIBER-EU §3.4 (test-data handling) and the RTS on TLPT under DORA Article 26(11). A typical retention is the engagement's working-paper retention period plus the regulatory minimum for TLPT evidence — substantially shorter than the production chain's 7-year retention, and bounded by the institution's red-team data-protection program.

**Cross-tenant isolation.** A TLPT chain runs under its own `tenant_id`; the per-tenant HKDF binding (spec §4.1) means a TLPT chain entry cannot mechanically be lifted into a production chain (cross-chain-lift detection at spec §7 step 4). This is the cryptographic floor that makes the segregation robust: an operator who accidentally points the TLPT receiver at a production storage path cannot produce a chain that verifies under the production tenant's IKM.

### 5.3 TIBER-EU and CBEST cross-references

TIBER-EU §3.4 names the test-data-handling expectations the institution's TLPT contract must observe. The institution's CC8.1 names the TLPT chain's retention bucket, the access-control separation between the TLPT lead and the production operations team, and the engagement-completion procedure that closes the TLPT chain.

### 5.4 The TLPT engagement evidence flow

```mermaid
sequenceDiagram
    participant rt as Red team
    participant app as Application under test
    participant sdk as SDK in TLPT mode
    participant rcv as Receiver routing
    participant tlpt as TLPT segregated chain
    participant prod as Production chain
    participant lead as TLPT lead and blue team

    rt-&gt;&gt;app: prompt-injection probe under engagement_id E1
    app-&gt;&gt;sdk: emit event with audit.tlpt.engagement_id=E1
    sdk-&gt;&gt;rcv: OTLP with tenant.testing_mode=tlpt resource attribute
    rcv-&gt;&gt;tlpt: route to segregated retention bucket
    rcv--&gt;&gt;prod: production chain unaffected
    lead-&gt;&gt;tlpt: offline witness-verifier confirms structural integrity
    Note over tlpt,prod: Cross-tenant isolation enforced by per-tenant HKDF binding per spec §4.1
```

The flow keeps the production chain free of attack TTPs while the TLPT chain remains integrity-bound and verifiable. The TLPT lead and the institution's blue team verify the chain offline using the witness-mode verifier; the production operations team has no read access to the TLPT chain.

---

## 6. DORA Articles 6-9 — ICT risk management framework, RTO and RPO

DORA Articles 6-9 require the institution to maintain an ICT risk-management framework with documented recovery-time objective (RTO) and recovery-point objective (RPO) targets. Article 12 (ICT business-continuity policy) operationalizes the framework. The chain's multi-region resilience patterns (spec §10.15) compose with the institution's RTO/RPO targets. *(complements `dr-and-resilience.md`)*

### 6.1 The chain's resilience contribution

The chain itself is not the institution's resilience layer; the institution's standard DR program is. The chain contributes integrity evidence under the institution's resilience program. The reference RTO and RPO targets named in `dr-and-resilience.md` §1 are not normative for the chain; they are the institution-side targets the chain composes with.

| Layer | Reference RPO | Reference RTO | DORA Article 12 mapping |
|---|---|---|---|
| Chain WAL (PostgreSQL streaming replication) | < 1 second | < 15 minutes | Recovery from primary-database failure |
| Hot store | Derivable from WAL | < 1 hour | Re-derive from WAL on standby |
| Cold store | < 1 hour | < 4 hours | Cross-region replication |
| Daily seals | Same as WAL | Same as WAL | Stored in WAL-backed table |
| HSM | Sub-seconds (active-active cluster) | < 5 minutes | HSM cluster failover |

The institution declares its own targets in its CC8.1 control description under Article 9 (operational resilience) and Article 12 (business-continuity policy).

### 6.2 Pattern A under DORA

Spec §10.15 Pattern A (active-active with seal-region pinning) is the conformant multi-region posture for institutions whose risk tolerance admits cross-region replication trust. The institution's CC8.1 names the seal region per tenant; the per-region event-count reconciliation (Pattern A invariant 5) satisfies Article 12's expectation that the institution operate verifiable replication completeness.

For DORA-scoped EU operations, the seal region MUST be in an adequate jurisdiction absent a documented Article 49 derogation or an EDPB-approved supplementary-measure architecture. The data-residency-by-default posture is in §11 below. *(addresses Lindqvist G-10)*

### 6.3 Pattern B under DORA

Spec §10.15 Pattern B (per-region `tenant_id`) is RECOMMENDED for EU tenants whose risk posture treats US-cloud-provider access as an unacceptable trust boundary. National prudential frameworks that mandate in-region key custody — BaFin's BAIT, ACPR's Notice on Outsourcing, and Banca d'Italia Circular 285 — are common drivers of Pattern B selection for German, French, and Italian institutions respectively. The institution's CC8.1 names the national framework that drives the Pattern B choice and the in-region key-custody arrangement that satisfies it. *(addresses Lindqvist P-5)*

---

## 7. EBA Guidelines on outsourcing — competent-authority audit rights

The EBA Guidelines on outsourcing arrangements (EBA/GL/2019/02) operationalize the ICT outsourcing requirements for credit institutions, payment institutions, and investment firms. Section 13 (rights of access and audit) and Section 14 (security of data and systems) are the load-bearing provisions for the chain. *(addresses Lindqvist G-8)*

### 7.1 Section 13.2 — competent-authority and appointed-auditor access

EBA/GL/2019/02 §13.2 requires the contract between the institution and the ICT service provider to grant the institution's competent authority (and any auditor appointed by it) **unrestricted access rights** to the service provider's premises and to all relevant business premises, devices, systems, networks, information, and data. The clause is non-negotiable for any outsourcing arrangement supporting a critical or important function.

The receiver-policy discovery endpoint mentioned in spec §4 is the surface where this access right surfaces in the institution's deployment. A vendor-hosted ledger that does not contractually grant the EBA-specified access right to the institution's competent authority is non-conformant under the Guidelines. The institution's contract with the receiver provider MUST include the §13.2 clause; the receiver-policy discovery endpoint MUST surface the access-right URI for the examiner's working paper.

When `docs/templates/eba-outsourcing-audit-rights.md` is published, this overlay cross-references that template for the load-bearing clause text. Until then, the institution's legal team drafts the clause directly from §13.2 and references this overlay's §7.1 for the regulatory backing.

### 7.2 Section 14 — security of data and systems

EBA/GL/2019/02 §14 requires the provider to apply security measures appropriate to the risk. The chain's three-layer integrity construction (spec §4.1, §4.2, §4.3) and the FIPS 140-2 Level 3 (or higher) HSM custody (spec §10.5) are the operational expression of §14 for the integrity dimension. The confidentiality, availability, and resilience dimensions are mapped in `article-32-security-mapping.md`; that mapping addresses GDPR Article 32 and EBA §14 in parallel.

### 7.3 Cloud outsourcing — EBA's 2017 recommendations and the 2023 update

The EBA's recommendations on outsourcing to cloud service providers (originally EBA/REC/2017/03, updated through subsequent guidance) require the institution to operate cloud arrangements under additional discipline. §4.6 of the recommendations names the in-region storage expectation for cloud arrangements; institutions whose risk posture admits cross-region replication still document the supplementary measures per `international-transfers.md` §2.2 (Schrems II safeguard) and this overlay's §11.

---

## 8. eIDAS — qualification of the seal signature

The eIDAS Regulation, Regulation (EU) 910/2014 (and Regulation (EU) 2024/1183 introducing eIDAS 2.0), governs the legal effect of electronic signatures in the Union. Articles 25-29 distinguish three signature levels: simple electronic signatures (Article 3(10)), advanced electronic signatures (AdES, Article 3(11) and Article 26 criteria), and qualified electronic signatures (QES, Article 3(12) and Article 25(2) equivalent legal effect). The v1.0a specification does not declare a qualification level; this overlay provides the explicit declaration. *(addresses Lindqvist G-6, P-4)*

### 8.1 Default qualification — Advanced Electronic Signature (AdES)

The seal signature in the chain is produced by an HSM at FIPS 140-2 Level 3 or higher (spec §10.5) using Ed25519 (FIPS 186-5, RFC 8032). This signature meets the four Article 26 criteria for an Advanced Electronic Signature without further qualification:

| Article 26 criterion | How the chain satisfies it |
|---|---|
| (a) It is uniquely linked to the signatory | The HKDF binding (spec §4.1) and the per-tenant key registry (spec §10.1) bind the seal-signing key to the institution-as-signatory; the seal record's `public_key_id` resolves the signing entity through the tenant key registry |
| (b) It is capable of identifying the signatory | The seal record's `public_key_id` and the institution's CC8.1 control description identify the signing entity |
| (c) It is created using electronic-signature-creation data that the signatory can use under their sole control with a high level of confidence | FIPS 140-2 Level 3 HSM custody with separation-of-duties roster (spec §10.5); the signing-key is non-extractable; the seal-job operator role grants `sign` only |
| (d) It is linked to the data signed in such a way that any subsequent change is detectable | The seal payload (`sign_payload` per spec §4.3 amendment form) binds the Merkle root, the algorithm, the format version, the tenant_id, the seal date, and the cadence; tampering with any of these surfaces as signature-verification failure at spec §7 step 11 |

The default declaration for the v1.0a chain is therefore: **the seal signature is an Advanced Electronic Signature (AdES) within the meaning of eIDAS Article 26**.

### 8.2 Qualified Electronic Signature (QES) checklist

Article 25(2) gives a Qualified Electronic Signature the equivalent legal effect of a handwritten signature throughout the Union. A QES requires three additional conditions on top of the AdES baseline:

| QES condition | What the institution must arrange |
|---|---|
| The signing device is a Qualified Signature Creation Device (QSCD) on the EU Trust List | The institution selects an HSM model that appears on the EU Trust List of QSCDs maintained by the European Commission. Examples include certain configurations of Thales Luna, Entrust nShield, and Utimaco SecurityServer when deployed under their QSCD-listed configurations. AWS CloudHSM, Azure Managed HSM, and Google Cloud HSM are FIPS 140-2 Level 3 conformant for AdES but are not (at the time of this overlay) on the EU Trust List of QSCDs; institutions seeking QES status select an on-prem or vendor-hosted HSM that IS on the Trust List |
| The public-key certificate is issued by a Qualified Trust Service Provider (QTSP) under eIDAS Article 28 | The institution contracts a QTSP (e.g., one of the QTSPs listed on the European Trust List) to issue the certificate binding the seal-signing key to the institution as the legal entity |
| The signing-key generation occurred inside the QSCD | The institution's IKM-equivalent for the signing key — the seal private key — MUST be generated inside the QSCD by the QSCD's internal CSPRNG. Spec §10.6.1 names HSM-internal RNG as the highest-assurance RNG pattern; under QES, this is the only conformant pattern |

The institution declares its qualification posture in its CC8.1 control description with one of three values: `signature_qualification = "AdES"` (default for v1.0a deployments), `signature_qualification = "AdES-QC"` (Advanced Electronic Signature with Qualified Certificate, an intermediate posture available under eIDAS), or `signature_qualification = "QES"` (Qualified Electronic Signature, Article 25(2) full equivalent legal effect).

The qualification declaration is institution-side; the spec is qualification-neutral. An institution operating under `AdES` for one tenant and `QES` for another is conformant — the chain construction is identical at the byte level; what differs is the contractual stack and the certificate's QTSP status.

### 8.3 FIPS / Common Criteria equivalence under EN 419 221-5

Many European institutions operate HSMs evaluated under BSI AIS 31 (Germany) or ANSSI CSPN (France), or certified to Common Criteria EAL4+ under EN 419 221-5 (server-signing QSCDs) and EN 419 211 (smartcard-class SSCDs). The spec's FIPS 140-2 Level 3 baseline (§10.5) is recognized as equivalent for AdES purposes; the institution's CC8.1 names the certification regime its HSM operates under. *(addresses Lindqvist P-4)*

### 8.4 Qualified electronic time stamps (forward note)

eIDAS Article 41 gives qualified electronic time stamps a presumption of accuracy on date and time, shifting the burden of proof in any subsequent dispute. The chain's `mac_computed_at_utc` field (spec §4.4) is wall-clock and forensic-only; it is NOT a qualified time stamp. Spec §10.14 names RFC 3161 trusted-timestamp integration as a v1.x extension candidate. Institutions seeking the eIDAS Article 41 presumption operate the v1.x extension once it is published; until then, the institution's NTP-discipline evidence (spec §10.4) is the timestamp foundation, and the institution's IT witness testifies to NTP synchronization at deposition or examination time.

---

## 9. NIS2 Directive — incident-notification thresholds for essential entities

Directive (EU) 2022/2555 (NIS2), transposed into national law in 2024 across the Member States, classifies entities as "essential" or "important" and imposes incident-notification obligations to the national CSIRT (or competent authority). Many Union financial entities are also essential entities under NIS2; the institution's CC8.1 names the dual classification where applicable. *(addresses Lindqvist G-7)*

### 9.1 The three NIS2 clocks

NIS2 Article 23 sets three clocks for significant incidents:

| Clock | Purpose | Notification target |
|---|---|---|
| 24 hours | Early warning — the institution has detected a significant incident and provides preliminary information | National CSIRT or competent authority |
| 72 hours | Incident notification — the institution provides the assessment of the incident, including its severity and impact | National CSIRT or competent authority |
| 1 month | Final report — the institution provides a detailed description of the incident, the root cause, and the mitigation measures | National CSIRT or competent authority |

The Commission Implementing Regulation (EU) 2024/2690 specifies the technical and methodological requirements for incident notification under NIS2.

### 9.2 Reconciliation with DORA — lex specialis

For Union financial entities classified as essential entities under NIS2, DORA Article 1(2) makes DORA the lex specialis for ICT-related incidents. The institution's IR commander operates the DORA 4-hour clock as the binding initial-notification clock and notifies the national CSIRT under NIS2 in parallel where the national transposition requires it. Most Member States' transpositions accept the DORA notification as discharging the NIS2 obligation for ICT-related incidents at financial entities; a minority require parallel notification. The institution's CC8.1 names the national-transposition position and the parallel-notification procedure where applicable.

### 9.3 Same partial-disclosure verifier mode

The partial-disclosure verifier mode used for the DORA 4-hour clock (per §4.1 above) supports the NIS2 24-hour early warning identically. The institution's IR runbook produces the same witness-verifier output for the NIS2 early warning as for the DORA initial notification; the dossier content differs (NIS2 requires preliminary information; DORA requires categorization criteria) but the chain-evidence shape is the same.

---

## 10. DORA Article 29 — ICT concentration risk

DORA Article 29 requires the institution to assess concentration risk at the level of ICT third-party service providers. Where a single receiver provider serves N institutions across the Union, the failure of that provider becomes a Union-level systemic event. Article 31 allows the European Supervisory Authorities (ESAs) to designate a critical ICT third-party service provider, which materially changes the supervisory regime applicable to the provider. *(addresses Lindqvist G-2)*

### 10.1 The systemic dimension

The chain's vendor-hosted multi-tenant topology (per `international-transfers.md` Scenario D) accumulates institutions on shared infrastructure. The institution running on a vendor-hosted chain MUST track the concentration metrics the vendor publishes and SHOULD coordinate exit-strategy testing under Article 28(8) with other institutions on the same infrastructure where a coordinated approach is operationally tractable.

### 10.2 The receiver-policy concentration tier — articulation pattern

The receiver-policy discovery endpoint (referenced in spec §4) is the natural surface for concentration-tier disclosure. The articulation pattern recommended for vendor-hosted deployments is for the receiver-policy endpoint to expose a `concentration.tier` field with values aligned to the ESA designation thresholds:

| `concentration.tier` value | Meaning | ESA polling expectation |
|---|---|---|
| `tier-1-critical-candidate` | The provider's market share within the Union approaches the Article 31 threshold for designation as a critical ICT third-party service provider | ESAs poll for register-of-information compliance under DORA Title V Chapter I |
| `tier-2-significant` | The provider's market share is significant but below the Article 31 threshold | Institution-side concentration-risk monitoring |
| `tier-3-standard` | Standard ICT third-party arrangement | Institution-side vendor-management procedure |

This articulation is NOT a normative extension to v1.0a; the spec normates the receiver-policy security floor (the discovery endpoint exists and is observed on the wire) but not the policy's content. The articulation pattern is institution-side; the institution's CC8.1 names the polling cadence and the vendor's commitment to publish the concentration-tier field. ESAs operating their register-of-information procedures under DORA Title V Chapter I poll the field when the vendor has committed to publish it; otherwise the ESA falls back to bilateral-disclosure procedures.

### 10.3 Per-institution-segregation evidence

A vendor-hosted multi-tenant deployment must provide per-institution-segregation evidence that the vendor's other customers' data is not accessible from the institution's tenant. The chain's per-tenant HKDF binding (spec §4.1) and the IKM-registry uniqueness enforcement (spec §10.1) provide the cryptographic floor: two tenants whose IKMs are accidentally swapped derive verifiably different session keys, and the per-entry `key_fingerprint` check (spec §7 step 8) catches the swap before any MAC compute. The cryptographic floor does not absolve the vendor from the contractual segregation commitments under EBA/GL/2019/02 §4.7 (sub-outsourcing); the institution's vendor-management procedure consumes both layers.

### 10.4 Receiver-policy minimum surface for examiner inspection

The receiver-policy discovery endpoint is the natural surface for an examiner inspection visit and for the EBA/GL/2019/02 §13.2 audit-rights assertion. A per-vendor protocol forces the examiner to learn each implementation; a baseline minimum surface lets the examiner perform a competent-authority inspection without per-vendor training. *(addresses Lindqvist P-6)*

The minimum surface recommended for examiner inspection is below. The articulation is institution-side; the spec normates the receiver-policy security floor (the discovery endpoint exists and is observed on the wire per spec §4) but not the policy's content shape. Vendors operating receiver-policy endpoints SHOULD adopt this minimum surface so the institution's competent authority recognizes the response shape without vendor-specific training.

| Field | Type | Purpose |
|---|---|---|
| `receiver.spec_version` | string | The chain spec version the receiver verifies; e.g. `"v1.0a"` |
| `receiver.posture` | string | One of `"ffiec"` (FFIEC-conformance) or `"vendor:<name>"` (vendor-flag mode); per spec §4.1.2 |
| `receiver.tenant_id` | string | The tenant identifier the response covers |
| `receiver.seal_region` | string | The seal region per spec §10.15 Pattern A; or `"per-region"` for Pattern B |
| `receiver.replication_regions` | string[] | The replication regions under Pattern A; empty for Pattern B and single-region |
| `receiver.cadence` | string | One of `"daily"` / `"hourly"` / `"weekly"` per spec §4.2.1 |
| `receiver.algorithm_posture` | string[] | The institution's declared signing-algorithm posture (e.g., `["ed25519"]` for single-algorithm; `["ed25519", "dilithium3"]` during dual-algorithm transitional period per spec §10.10.2) |
| `receiver.audit_rights_uri` | URI | The institution's reception URI for an EBA/GL/2019/02 §13.2 competent-authority inspection request |
| `receiver.concentration_tier` | string | The vendor-published concentration tier per §10.2 of this overlay |
| `receiver.signature_qualification` | string | One of `"AdES"` / `"AdES-QC"` / `"QES"` per §8.2 of this overlay |
| `receiver.data_residency` | object | Per-data-class residency declaration: `{"chain_data": "EEA", "tenant_config": "EEA", "ikm_custody": "EEA-QSCD"}` for the strongest Schrems II posture |

The institution's competent authority polls the endpoint at examination time and reads the response into the working paper. A response that omits any of the required fields above is non-conformant for examiner-inspection purposes; the institution's contract with the receiver provider names the minimum surface in operational language. The minimum surface composes with the spec's wire-or-on-disk observation rule (§4 of the spec) — the endpoint's response is observed on the wire and is the citable form.

---

## 11. Schrems II supplementary measures — chain-of-custody specifics

The Schrems II jurisprudence (CJEU C-311/18 and the related C-362/14 antecedent) and the EDPB Recommendations 01/2020 require the controller to apply supplementary measures where third-country protection is not essentially equivalent to Union standards. The general Schrems II treatment for the chain lives in `international-transfers.md`; this section names the chain-of-custody-specific supplementary measures. *(addresses Lindqvist P-1, G-10; complements `international-transfers.md`)*

### 11.1 Cross-reference posture

`international-transfers.md` is the load-bearing document for Schrems II. It enumerates four deployment scenarios (EU on-prem, EU data on US cloud, multi-region active-active, vendor-hosted multi-tenant), the SCC modules, the six-step Transfer Impact Assessment, and the Article 28 DPA. This overlay does not duplicate that material; this section adds the chain-of-custody-specific supplementary measures the institution operates on top of the base Schrems II treatment.

### 11.2 EU-held keys — promoted to MUST for vendor-hosted EU controllers

For any vendor-hosted deployment serving an EU controller, the institution's risk-tolerance posture under Schrems II treats US-cloud-provider access as a trust boundary. The EU-held-keys posture — the chain's HMAC IKM, the chain's confidentiality key (AES-GCM at rest), and the seal-signing key all custodied in HSMs physically located in the EU under the operational control of an EU subsidiary — is the cryptographic floor that satisfies EDPB Recommendations 01/2020 §85 (Use Case 1: data storage for backup and other purposes that do not require access to data in the clear).

This overlay promotes EU-held keys from "preferred" (per `international-transfers.md` §2.2) to **MUST** for vendor-hosted EU controllers absent a documented EDPB-approved alternative. The institution's CC8.1 names the EU subsidiary that holds the keys, the HSM custody location, and the operational-control evidence (separation-of-duties roster per spec §10.5, key-rotation procedure per spec §10.10).

The promotion is scoped to vendor-hosted deployments serving EU controllers. Pattern A (multi-region active-active) and Pattern B (per-region tenant_id) under spec §10.15 have their own per-region key-custody postures; this promotion does not change those postures, it constrains the vendor-hosted deployment shape. *(addresses Lindqvist P-1)*

### 11.3 TenantConfig and PII-pattern stores in the EU

For institutions operating tenant-side PII redaction patterns (the vendor-shipped TenantConfig store, where applicable), the configuration store MUST be hosted in the EU. This is the natural extension of the EU-held-keys posture: configuration that drives redaction decisions on customer content is itself a confidentiality control, and storing it in a US-cloud-provider exposes the redaction decisions to the same FISA §702 / EO 12333 / CLOUD Act surveillance risk Schrems II identifies.

### 11.4 Transit encryption and Article 49 derogations

Cross-region replication encryption (spec §10.15 Pattern A) MUST be verified end-to-end under the institution's TIA. Where the institution operates Article 49 derogations (typically Article 49(1)(b) contract-necessity or Article 49(1)(d) important-reasons-of-public-interest), the derogation's documentary basis is part of the institution's TIA per `international-transfers.md` §4. The chain's role in the derogation evidence is to bind the institution's transit-encryption posture and the actual cross-region replication events under integrity — the operational events `master.cross_region_replication_completed` (spec §10.2) carry per-region replication evidence the institution's TIA consumes.

### 11.5 Vendor-hosted EU controller — the QTSP-issued seal certificate

For vendor-hosted EU controllers seeking the strongest Schrems II posture, the seal certificate SHOULD be issued by an EU-established QTSP under eIDAS Article 28. The QTSP's certificate root is in the European Trust List; the institution's seal-signing key is then bound to a trust anchor that operates under European supervisory law. Combined with the EU-held-keys posture (§11.2), this closes the trust chain to within the Union for institutions whose risk posture demands it.

### 11.6 The CLOUD Act exposure and the receiver topology

The US CLOUD Act extends US legal process to data held by US providers regardless of storage location. For chain deployments using US-headquartered cloud providers — even when the storage region is in the EU — the CLOUD Act exposure is the recurring trust-boundary concern. `multi-jurisdiction-conflict.md` §2.4 names the operational mitigations; this overlay's §11.2 names the cryptographic floor (EU-held keys). The institution's contract with the receiver provider names the cooperation commitment under a CLOUD Act event — typically a notice obligation to the institution within 24 hours of receiving US legal process and a cooperation commitment in any motion to quash or modify the order under the procedure in `multi-jurisdiction-conflict.md` §2.3.

The chain's role under a CLOUD Act event is to produce the integrity-bound record of what the institution disclosed and what it withheld; the institution's IR commander operates the disclosure procedure under legal counsel and the institution's DPO. The chain does not change the institution's substantive disclosure obligation; it produces the audit trail the institution's competent authority and the EU DPA review after the event.

---

## 12. GDPR Article 22 and EU AI Act Article 86 — right to explanation

The chain proves what the AI said. GDPR Article 22 gives a data subject who is the subject of a decision based solely on automated processing the right to obtain meaningful information about the logic involved, the significance of the processing, and the consequences. Regulation (EU) 2024/1689 (EU AI Act) Article 86, applicable to deployers of high-risk AI systems, gives an affected person the right to obtain from the deployer a clear and meaningful explanation of the role of the AI system in the decision-making procedure and the main elements of the decision taken. The chain is the obvious source of the explanation evidence. *(addresses Lindqvist G-11)*

### 12.1 What the chain delivers for the explanation right

The chain entry representing the AI's response to the data subject's matter contains the integrity-bound record of: the prompt the institution sent (under `gen_ai.request.messages` in OTel GenAI semconv); the parameters the institution used (under `ffiec.chain.gen_ai_parameters` per spec §4.4 — temperature, top_p, top_k, seed, max_tokens, sampling implementation, system-prompt content or content-hash); the model identifier the vendor's API actually answered with (under `gen_ai.response.model` per spec §4.4 normative); the response itself; the routing decision that selected the provider (under spec §4.4.1 `audit.routing.*` chain entries); and the deployment intent that classified the invocation (under spec §4.4.2 `audit.deployment.*` attributes). For ECOA-style adverse-action notices the institution operates outside the EU, the translation entry's chain (spec §10.11) extends the record into the customer-language disclosure.

### 12.2 What the institution's explanation procedure adds

The chain's record is integrity evidence, not the explanation itself. The explanation is a customer-facing artifact the institution constructs from the chain's evidence plus the institution's substantive explanation logic — the role of the AI in the decision-making procedure, the main parameters, and the consequences. Article 86's "main elements of the decision taken" requires the institution to translate the chain's integrity-bound record into customer-meaningful language; the chain does not generate the customer-language artifact mechanically.

The institution's data-subject explanation procedure operates the following sequence:

1. The data subject submits an Article 22 / AI Act Article 86 request naming the decision and the data subject's identity
2. The institution's privacy office locates the chain entry representing the decision (typically by `run_id` and the customer-facing decision identifier)
3. The institution's verifier runs against the affected tenant-day, producing PASS or FAIL for the chain segment containing the decision
4. The privacy office assembles the explanation from the chain's record plus the institution's explanation library (the role of the AI, the main parameters in customer-meaningful language, the consequences)
5. The institution's redaction procedure removes information not pertinent to the data subject (other customers' data, internal model artifacts, security-sensitive fields) before disclosure
6. The institution delivers the explanation within the regulatory window — Article 22 names "appropriate measures" without a fixed window; AI Act Article 86 names the deployer's reasonable-time obligation; the institution's CC8.1 names a target window aligned with its data-subject-rights program

### 12.3 The redaction discipline before disclosure

The chain entry's `gen_ai.request.messages` field will frequently contain customer-content, system-prompt content, and retrieval context that includes other customers' data or institution-internal information. The institution's redaction procedure removes those fields before disclosure to the data subject. The redaction operates on a copy; the original chain entry remains integrity-bound and unchanged. The institution's redaction discipline is documented in the institution's data-subject-rights procedure and consumes the chain's evidentiary artifacts (per spec §10.13) without modifying the chain itself.

### 12.4 Cross-reference and forward note

A future companion document `regulator-pack/ai-act-article-86.md` will operationalize the AI Act Article 86 explanation procedure in detail. Until that companion is published, this overlay's §12 names the chain's role and the institution's procedure shape. The institution's CC8.1 names the explanation procedure under both GDPR Article 22 and AI Act Article 86, recognizing that the two regimes overlap in substance for high-risk AI systems making decisions affecting natural persons.

---

## 13. ECB SREP / SSM expectations on operational resilience

The European Central Bank's Single Supervisory Mechanism (SSM) operates SREP (the Supervisory Review and Evaluation Process) for significant credit institutions in the euro area. The SREP guide §5.5 (operational risk, ICT risk) and the SSM's expectations on operational resilience for clearing-relevant systems set the supervisory posture for the chain's evidence in capital-allocation review. *(addresses Lindqvist G-9)*

### 13.1 ICAAP / ILAAP residual-risk articulation

The institution's ICAAP (Internal Capital Adequacy Assessment Process) and ILAAP (Internal Liquidity Adequacy Assessment Process) carry the institution's residual-risk allocation under Pillar 2. Chain-integrity failures (per spec §10.2 operational-event catalog: `chain.verification_failure`, `audit_file.truncation_detected`, `key_fingerprint.mismatch`, `hsm.operation_failure`) are operational-risk scenarios the institution names in its ICAAP residual-risk register. Where the chain documents a clearing-relevant decision (correspondent-banking, payment-screening, KYC at clearing scale) an integrity failure could propagate to TARGET2 reconciliation; the SSM's expectations on operational resilience treat such propagation as a Pillar 2 capital-allocation matter.

The institution's ICAAP team consumes the chain's evidence — the verifier output for the period showing PASS for each tenant-day, plus any anomaly documented under the institution's IR record (per spec §10.13 evidentiary artifacts). The chain's role is to produce the evidence; the analysis is the institution's responsibility. The overlay does not extend the chain's integrity scope into the substantive ICAAP analysis.

### 13.2 Joint Supervisory Team (JST) review

For significant institutions under direct ECB supervision, the Joint Supervisory Team conducts the SREP review. This articulation overlay is the single-document landing point for the JST's ICT third-party assessment under DORA Article 28-30. The verifier's deterministic, single-binary, no-network posture (spec §7) is what makes the assessment practicable inside the JST's working-paper budget — the JST runs the verifier offline against the institution's evidence, observes PASS or FAIL deterministically, and reads this overlay alongside the spec to interpret the result.

The JST's typical reading order:

1. The institution's CC8.1 control description names the postures the institution operates (cadence, signature qualification, multi-region pattern, EU-held-keys posture for vendor-hosted deployments)
2. This overlay's §13 (translation table) maps each posture to the European-framework articles
3. The institution's verifier output for the SREP period evidences each tenant-day's PASS or FAIL
4. The institution's evidentiary-artifact bundle (spec §10.13) substantiates the foundation of the verifier output
5. The JST's working paper records the assessment result and any management response on residual-risk findings

The bundle is the single working-paper input the JST consumes; the institution's CC8.1 and this overlay close the articulation gap so the JST does not need to reverse-engineer the chain's evidence model from the spec alone.

### 13.3 TARGET2 / T2S clearing-relevant decisions

Where the chain documents an AI-driven decision in correspondent banking, payment screening, KYC at clearing scale, or any other workflow whose output enters TARGET2 reconciliation or T2S settlement, an integrity failure on that chain becomes systemically relevant. The SSM's expectations on operational resilience for clearing-relevant systems treat such propagation as a Pillar 2 capital-allocation matter and as an incident-classification trigger under DORA Article 18.

The institution's chain-of-custody program names the clearing-relevant tenants in its CC8.1 with an elevated criticality classification. The institution's IR runbook treats `chain.verification_failure` events on a clearing-relevant tenant as DORA major-incident candidates by default; the IR commander operates the 4-hour clock from the moment of detection rather than from a downstream classification step.

### 13.4 Capital-allocation evidence flow

The institution's ICAAP cycle consumes the chain's evidence as one input to the operational-risk component of Pillar 2. The flow:

| Step | Owner | Input | Output |
|---|---|---|---|
| Annual ICAAP refresh | ICAAP team | Verifier output for the period (PASS/FAIL per tenant-day per spec §10.13) | Operational-risk scenario inventory |
| Quarterly residual-risk review | ICAAP team + IR commander | `chain.verification_failure` operational events (spec §10.2); incident-classification dossiers per §4.2 of this overlay | Updated residual-risk capital allocation |
| Joint Supervisory Team dialogue | JST + institution | Annual residual-risk register; verifier outputs for any SREP-period incident | JST working paper closing the operational-risk dimension |

The chain's role is to produce integrity evidence the ICAAP and the JST consume; the chain does not perform the capital allocation. The institution retains responsibility for the analytical work and for the residual-risk register's content.

---

## 14. Translation table — v1.0a sections to European framework articles

The translation table below maps each load-bearing v1.0a section to the European framework articles it satisfies, partially satisfies, leaves outside scope, or for which it requires a supplementary measure. The status values:

- **ALIGNED** — the v1.0a section satisfies the European-framework article without supplementary measure
- **PARTIAL** — the v1.0a section satisfies the article in part; the institution's CC8.1 control description completes the articulation
- **OUTSIDE-SCOPE** — the article is outside the chain's scope; the institution's other systems satisfy it
- **SUPPLEMENTARY-MEASURE-REQUIRED** — the v1.0a section satisfies the article only when a named supplementary measure is in place

| v1.0a section | DORA article | EBA/GL/2019/02 | eIDAS article | NIS2 article | GDPR article | Status |
|---|---|---|---|---|---|---|
| §1.1 Daubert grounding | n/a (informative) | n/a | n/a | n/a | n/a | OUTSIDE-SCOPE (US Daubert; EU analogue is eIDAS Article 35-37 and the e-Evidence Regulation (EU) 2023/1543) |
| §1.2 Epistemic scope | Article 6 (general principles) | n/a | n/a | n/a | Article 5(2) accountability | ALIGNED |
| §1.3 Security definitions | Article 9 (operational resilience) | §14 | Article 26 | Article 21 | Article 32(1)(b) | ALIGNED |
| §3 Tenant binding | Article 28(2) sub-processor articulation | §4.7 | n/a | n/a | Article 28 | PARTIAL (CC8.1 names the registry uniqueness enforcement) |
| §4.1 Per-event MAC | Article 9 | §14 | Article 26(d) | Article 21 | Article 32(1)(a) | ALIGNED |
| §4.2 Daily Merkle seal | Article 9 | §14 | Article 26(d) | Article 21 | Article 32(1)(b) | ALIGNED |
| §4.3 HSM signature | Article 9; Article 30(2)(c) | §14 | Article 26 (AdES) / Article 25(2) (QES) | Article 21 | Article 32(1)(a) | PARTIAL (default AdES; QES requires QTSP-issued certificate per §8.2) |
| §4.4 OTLP wire | Article 9; Article 30(2)(c) | §14 | n/a | Article 21 | Article 32(1)(b) | ALIGNED |
| §4.4.1 Routing decisions | Article 12 (BCP) | §14 | n/a | Article 21 | Article 22 (automated decision-making) | PARTIAL (institution's CC8.1 names routing-policy versioning) |
| §4.4.2 Deployment intent | Article 6; Article 12 | §14 | n/a | Article 21 | Article 22; AI Act Article 86 | PARTIAL (institution's CC8.1 names deployment-policy versioning) |
| §4.4.3 OTLP transport identification | Article 30(2)(a) | §13 | n/a | Article 21 | n/a | ALIGNED |
| §5 Wire format / canonical encoding | Article 9 | §14 | Article 26(d) | Article 21 | Article 32(1)(b) | ALIGNED |
| §5.1 Transport encryption | Article 9; Article 30(2)(c) | §14 | n/a | Article 21 | Article 32(1)(a) | ALIGNED |
| §6 Storage | Article 9 | §14 | Article 26(d) | Article 21 | Article 32(1)(b), (c) | ALIGNED |
| §7 Verification | Article 28(7) audit; Article 30(2)(f) | §13.2 | Article 26(d) | Article 21 | Article 32(1)(d) | ALIGNED |
| §8 Conformance test vectors | Article 28(7) | §14 | n/a | Article 21 | Article 32(1)(d) | ALIGNED |
| §10.1 Reconciliation | Article 8 (ICT risk) | §14 | n/a | Article 21 | Article 32(1)(d) | ALIGNED |
| §10.2 Operational events | Article 11-14 (incident reporting) | §14 | n/a | Article 23 | Article 33 | PARTIAL (DORA-bound deployments operate higher-frequency cadence per §4.1 of this overlay) |
| §10.3 Append-only enforcement | Article 9 | §14 | n/a | Article 21 | Article 32(1)(b) | ALIGNED |
| §10.4 Time synchronization | Article 9 | §14 | Article 41 (forward) | Article 21 | Article 32(1)(b) | PARTIAL (qualified time stamps per §8.4 of this overlay are a v1.x extension) |
| §10.5 HSM custody | Article 9; Article 30(2)(c) | §14 | Article 29 (QSCD) | Article 21 | Article 32(1)(a) | PARTIAL (FIPS 140-2 L3 baseline; QSCD requires Trust List entry per §8.2) |
| §10.6 IKM minimum length | Article 9 | §14 | n/a | Article 21 | Article 32(1)(a) | ALIGNED |
| §10.7 Software-key adapter exclusion | Article 9 | §14 | n/a | Article 21 | Article 32(1)(b) | ALIGNED |
| §10.9 IKM registry retention | Article 28(8) exit strategy | §13 | n/a | Article 21 | Article 5(1)(e) | ALIGNED |
| §10.10 Rotation across seal boundary | Article 9 | §14 | n/a | Article 21 | Article 32(1)(d) | ALIGNED |
| §10.11 ECOA translation | Article 9 | §14 | n/a | n/a | Article 22; AI Act Article 86 | OUTSIDE-SCOPE (US ECOA; EU analogue is GDPR Article 22 and AI Act Article 86) |
| §10.13 Evidentiary artifacts | Article 28(3) register | §13.2 | Article 33 (validation) | Article 21 | Article 5(2) | ALIGNED |
| §10.14 Trusted-time integration | Article 9 | §14 | Article 41 | Article 21 | Article 32(1)(b) | SUPPLEMENTARY-MEASURE-REQUIRED (RFC 3161 v1.x extension for Article 41 presumption) |
| §10.15 Multi-region resilience | Article 12 (BCP); Article 9 | §13.2 (audit rights cross-region); §14 | n/a | Article 21 | Article 44; Article 32(1)(c) | SUPPLEMENTARY-MEASURE-REQUIRED for vendor-hosted EU controllers (EU-held keys per §11.2 of this overlay) |

The table is the dense reference for the JST's working paper. Each row names the load-bearing article in each regime; the status column tells the JST whether to expect supplementary documentation from the institution.

---

## 15. Bottom line for an EU competent authority

In an SREP review or a Joint Examination Team review, this articulation overlay is the single document the supervisor reads alongside the v1.0a specification. The spec's American-framework primary alignment does not impede European obligations — the integrity primitives in §4 are framework-neutral cryptographic constructs and the FFIEC origin shows in the framing, not in the primitives. The overlay closes the articulation gap: DORA Articles 28-30 lattice positioning, the corrected DORA 4-hour clock with partial-disclosure verifier mode, the eIDAS qualification declaration, the NIS2 reconciliation, the EBA Guidelines on outsourcing audit-rights backing, the Schrems II supplementary measures specific to the chain, the ICT concentration-risk articulation, and the right-to-explanation interface under GDPR Article 22 / EU AI Act Article 86 are all named in one place. Subject to the institution's CC8.1 control description naming the postures the overlay calls out, the chain's evidence is directly referenceable in a Union joint examination team's working file.

The integrity primitives carry their weight without supplementary measure. The supervisory translation layer — the contractual stack, the qualification declarations, the data-residency posture, the cadence binding for the 4-hour clock, the EU-held-keys posture for vendor-hosted EU controllers — is institution-side work that this overlay names in operational language. Subject to that supplementary work being in place, the chain is directly referenceable in a Union joint examination team's working file under DORA Article 28 ICT third-party risk, EBA/GL/2019/02 §13 audit rights, eIDAS Article 26 (or Article 25(2) for QES-postured deployments), NIS2 Article 23 incident notification, and GDPR Article 32 security of processing. The cryptographic substrate is not the obstacle; the supervisory translation layer that this overlay provides is what closes the file.

---

## 16. Cross-references

- `international-transfers.md` — Schrems II treatment for the chain; this overlay's §11 layers the chain-specific supplementary measures on top
- `article-32-security-mapping.md` — GDPR Article 32 controls; this overlay's §7.2 names the parallel EBA/GL/2019/02 §14 backing
- `multi-jurisdiction-conflict.md` — multi-jurisdiction handling including DORA / GDPR / NIS2 deadline reconciliation
- `breach-notification-matrix.md` — operational matrix of notification clocks; this overlay's §4.3 names the lex specialis position
- `gdpr-controller-vs-processor.md` — GDPR Article 28 lattice; this overlay's §2 names the parallel DORA Article 28-30 lattice
- `dr-and-resilience.md` — RTO / RPO targets and multi-region patterns; this overlay's §6 names the DORA Article 12 mapping
- Spec §1.2 (epistemic scope), §3 (tenant binding), §4 (the four primitives), §5 (wire format), §7 (verifier procedure), §10.5 (HSM custody), §10.13 (evidentiary artifacts), §10.14 (trusted-time integration), §10.15 (multi-region resilience)
- `docs/templates/eba-outsourcing-audit-rights.md` — load-bearing audit-rights clause text (forward reference; template to be published)

---

## 17. Review cadence

This overlay is reviewed annually by the institution's legal team, privacy office, DPO, and competent-authority liaison. Triggers for early review:

- A CJEU decision affecting transfer mechanisms, third-country surveillance analysis, or the Schrems jurisprudence
- An updated Commission Delegated or Implementing Regulation under DORA, NIS2, or eIDAS
- An eIDAS 2.0 implementing-act publication that changes the QSCD or QTSP regime
- An EBA guideline revision affecting outsourcing arrangements or ICT third-party risk
- An ESA designation of the institution's receiver provider as a critical ICT third-party service provider under DORA Article 31
- The institution adds or removes a deployment region in a way that changes the data-residency posture
- A material change in the institution's contractual stack with the SDK provider, receiver provider, HSM provider, or LLM provider
