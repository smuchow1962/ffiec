# Customer dispute procedures

> **What this doc is.** The customer-side path for disputing an AI-driven decision and the institution's response. Closes the model-risk customer-dispute partial.

## Customer initiates dispute

Customer believes the bank's AI-driven decision was incorrect, biased, or based on tampered records. The customer can:

1. **Direct dispute** — contact the bank's dispute resolution function (standard customer-service channel)
2. **Regulator complaint** — file with CFPB, OCC, FDIC, Fed, or state attorney general
3. **Litigation** — civil suit (typically last resort)

The chain provides the same evidence base regardless of dispute path.

## Bank's initial response

When a dispute arrives:

1. **Identify the affected run.** Using the customer's identifier (account, application ID, complaint reference), the bank's customer-correlation index (CUEC) maps to chain `(tenant_id, run_id)`.
2. **Pull the chain record.** The bank's chain-ops team pulls the run's events; produces a verifier `walk` output for the run.
3. **Review for accuracy.** The customer service representative or compliance specialist reviews the run's events: prompt, model output, decision.
4. **Respond to customer.** Direct disputes typically resolve within the bank's standard SLA; regulator complaints follow the regulator's response timeline.

## When the customer asks "How do I know the bank didn't tamper with this record?"

The customer's reasonable concern. The institution's response:

**Option A — Refer to the bank's SOC report.** If the bank has a SOC 2 or SOC 1 report covering the chain, the customer can request the report; an independent auditor has attested to the chain's controls.

**Option B — Refer to the regulator's verification.** The customer's primary regulator periodically examines the chain. The regulator's findings are public (in some cases) or available on request.

**Option C — Independent verification.** The customer's lawyer can hire a cryptographic expert to run the verifier on a bank-produced ledger snapshot, against the bank's public key (which the regulator already has). Independent verification is a design intent.

For most disputes, Option A or B is sufficient. Option C arises in litigation contexts.

## Adverse-action notices (ECOA)

When the bank's AI denies credit, ECOA requires the bank to provide an adverse-action notice with specific reasons. The chain captures the AI's reasoning (the model's response). The bank's compliance procedure translates:

1. The chain's record provides the AI's response (which may include the model's stated reasons)
2. The bank's compliance team translates the AI's response into ECOA-compliant adverse-action language
3. **The translation is itself a chain entry**, with the AI's original response as the parent (per spec §4.4 `parent_run_id` / `parent_seq` linkage) and `chain_kind = "translation"` per spec §3. Institutions MUST chain the translation step. Logging the translation outside the chain leaves the customer-side repudiation gap the chain exists to close — a non-conformant posture for ECOA-supervised activity. The requirement is asymmetric in cost: institutions that already chain the translation gain nothing from the change; institutions that do not chain the translation close the customer-side repudiation gap they currently expose at no incremental cost beyond the existing chain pipeline.
4. The notice is sent to the customer per ECOA timing requirements

The chain does not generate the adverse-action notice; it provides the evidence that supports the notice's reasons.

### Translation chain-entry schema (normative per spec §10.11)

Spec §10.11 normates the attribute schema for the translation chain entry. The entry's `chain_kind` is `"translation"` per spec §3. The attributes are part of the canonical bytes — the chain MAC covers them, same as any other `audit.*` namespace.

| Attribute | Required | Notes |
|---|---|---|
| `audit.ecoa.translation.target_language` | yes | BCP 47 language code the translation produced (e.g., `es-US`, `zh-CN`) |
| `audit.ecoa.translation.source_language` | RECOMMENDED | BCP 47 language code of the original AI response |
| `audit.ecoa.translation.translator_kind` | yes | One of `human` \| `llm` \| `translation_api` \| `glossary_lookup` |
| `audit.ecoa.translation.translator_id` | conditional | REQUIRED when `translator_kind = "llm"` (the model identifier; matches `gen_ai.response.model` shape). Optional otherwise |
| `audit.ecoa.translation.glossary_version` | conditional | When the institution operates a regulated-language glossary the translation conforms to: the glossary version |
| `audit.ecoa.translation.output_hash` | yes | Lowercase hex SHA-256 of the customer-facing translated text. The text itself MAY be customer-PII; the hash binds the translation under the chain without binding the PII |
| `audit.ecoa.translation.delivery_method` | RECOMMENDED | One of `mail` \| `secure_message` \| `email` \| `phone` \| `in_person` |
| `audit.ecoa.translation.delivery_timestamp` | RECOMMENDED | RFC 3339 UTC time at delivery |

The translation entry binds to the AI's original response via `parent_run_id` / `parent_seq` (spec §4.4). The original response and the translation form a pair: the original is integrity-bound by its own chain entry; the translation is integrity-bound by the translation entry; the link between them is integrity-bound by the canonical-bytes inclusion of `parent_run_id` / `parent_seq` (spec §5).

**Delivery confirmation.** The delivery event itself (the "this letter was mailed" event) is institution-tracked separately and is not required to be a chain entry. `audit.ecoa.translation.delivery_method` and `audit.ecoa.translation.delivery_timestamp` on the translation entry are the chain-of-custody record of the institution's delivery commitment; the operational evidence that confirms delivery occurred (postal tracking, secure-message receipt, email logs) is institution-side.

### Worked example — chained ECOA cycle

Consider an ECOA adverse-action cycle where the AI generates the original adverse-action response in English and the institution translates it to Spanish for a customer whose preferred language is `es-US`. The chain captures the cycle as two linked entries:

**Entry 1 — original AI response (the AI's adverse-action reasoning).** This is the chain entry the AI's decision pipeline emits when the model produces the adverse-action content. The entry's `chain_kind = "model_call"` per spec §3 (the entry represents the LLM invocation). Selected attributes:

```
ffiec.chain.run_id          = "r_ecoa_2026_05_07_a3f29b71c"
ffiec.chain.seq             = 14
ffiec.chain.chain_kind      = "model_call"
gen_ai.request.model        = "claude-opus-4-7"
gen_ai.response.model       = "claude-opus-4-7"
audit.ecoa.adverse_action.reasons = "[institution's structured reasons payload]"
audit.ecoa.adverse_action.language = "en-US"
```

The entry's `payload_hash` is the HMAC-SHA-256 over `prev_hash || canonical_bytes` per spec §4.1.

**Entry 2 — translation chain entry (the institution's translation step).** The institution's compliance pipeline emits this entry when the translation completes. The entry's `chain_kind = "translation"`; the entry binds to entry 1 via `parent_run_id` / `parent_seq`. Selected attributes:

```
ffiec.chain.run_id                              = "r_ecoa_translation_2026_05_07_b7c3f8a2d"
ffiec.chain.seq                                 = 1
ffiec.chain.chain_kind                          = "translation"
ffiec.chain.parent_run_id                       = "r_ecoa_2026_05_07_a3f29b71c"
ffiec.chain.parent_seq                          = 14
audit.ecoa.translation.target_language          = "es-US"
audit.ecoa.translation.source_language          = "en-US"
audit.ecoa.translation.translator_kind          = "llm"
audit.ecoa.translation.translator_id            = "claude-opus-4-7"
audit.ecoa.translation.output_hash              = "8a4b7c5e2f9d1a3c6e8b4f7d2c5a9e1b3d7f5c8a2e4b6d9c1f3a5e7b9d2c4f6a"
audit.ecoa.translation.delivery_method          = "mail"
audit.ecoa.translation.delivery_timestamp       = "2026-05-08T14:30:00Z"
```

The translation entry's `payload_hash` is the HMAC-SHA-256 over `prev_hash || canonical_bytes` (its own canonical bytes, including the `parent_run_id` / `parent_seq` linkage per spec §5). Because `parent_run_id` and `parent_seq` are integrity-bound by the canonical bytes, an attacker cannot rewrite the linkage post-capture without breaking the MAC.

**The CFPB inquiry answered from the chain alone.** A CFPB / ECOA examiner asking "did this customer receive the adverse-action notice in their preferred language within the regulatory window?" answers from the two chain entries:

- **Original response language** — entry 1's `audit.ecoa.adverse_action.language = "en-US"` confirms the AI produced its reasoning in English.
- **Translation produced** — entry 2 exists, with `audit.ecoa.translation.target_language = "es-US"`. The customer's preferred language was Spanish; the translation step produced Spanish content.
- **Translation provenance** — `translator_kind = "llm"` plus `translator_id = "claude-opus-4-7"` records the institution used an LLM to translate; `output_hash` binds the translated text to the chain entry so the institution can produce the exact translated text on request and the examiner can confirm the bound hash matches the produced text.
- **Delivery posture** — `delivery_method = "mail"` and `delivery_timestamp = "2026-05-08T14:30:00Z"` record the institution's delivery commitment (the actual mailing record is institution-side per spec §10.11; the chain-recorded commitment lets the examiner cross-check the institution's mailroom records against the chain's stated delivery time).
- **Regulatory-window check** — the time delta between entry 1's `captured_at` and entry 2's `delivery_timestamp` answers the within-window question. ECOA requires the adverse-action notice within 30 days of the adverse-action decision; the chain entries record both anchor points.
- **Parent linkage** — entry 2's `parent_run_id` / `parent_seq` resolves to entry 1; the verifier walking the disputed run produces both entries and the linkage is integrity-bound. The customer-side narrative is "AI denied credit on date D in English; institution translated to Spanish and delivered by mail on date D+1" — answered from the chain entries alone, without the institution surfacing supplementary evidence.

The audit-procedure shape that operationalizes the translation-completeness check is `audit-procedures.md` P-35.

### Reproducibility and the customer's "would the AI have said the same thing on a different day" question

A customer's advocate (or the CFPB on the customer's behalf) is increasingly likely to ask not just "what did the AI say" but "would the AI have said the same thing on a different day" — a reproducibility question, not a logging-integrity question.

The chain captures the model's decoding parameters (where the institution has configured `gen_ai_parameters` capture per spec §4.4) and supports the institution's ability to re-run the decision under the same parameters, the same model version (request-side AND response-side, captured per OTel `gen_ai.request.model` / `gen_ai.response.model`), the same system prompt content (captured in `audit.*` per institution convention), and the same retrieval context. The institution's MRM team performs the reproduction; the chain integrity-binds the inputs that the reproduction needs.

The institution's response to the customer's reproducibility question:

> The decision's full input set is recorded in our integrity-bearing audit trail, including the AI model's version, the decoding parameters, the system prompt, and any retrieved context. We can re-run the decision under the same inputs to evaluate whether the AI would produce the same output. Significant variation across re-runs (a phenomenon called model nondeterminism) is itself information our model-risk-management committee evaluates; consistent variation in your favor would prompt us to revisit the original decision.

The institution's customer-service team adapts the language; the MRM committee establishes the institution's reproduction posture.

### MRM committee variation-threshold guidance

The MRM committee establishes — in advance, not at the moment of dispute — what variation across re-runs constitutes "significant" for the institution's reproduction posture. The threshold is tier-specific:

- **Decision-equivalence threshold.** Two re-runs are decision-equivalent when they produce the same final decision (approve/deny/route-to-X). Variation BELOW the decision-equivalence threshold is operational nondeterminism; variation AT or ABOVE the threshold (different final decision) triggers MRM review.
- **Per-decision-class threshold.** For high-stakes decisions (credit denial, fraud flag, loan-pricing adjustment), the threshold is decision-equivalence on 95%+ of re-runs (5 re-runs minimum). For medium-stakes decisions (account routing with material customer impact, fee determination, transaction categorization), the threshold is 90%+ decision-equivalence. For low-stakes decisions (routing without material customer impact, advisory recommendations, internal task assignment), 80%+ decision-equivalence is acceptable. Variation below the applicable threshold triggers documented model-validation review.
- **Customer-impact tier sensitivity.** The threshold tightens for higher customer-impact tiers per the institution's MRM-defined tier definition. The chain-captured `audit.*` payload should include a customer-impact-tier indicator so the MRM committee can stratify the variation analysis.

The committee documents the threshold in the institution's MRM policy framework; the customer-dispute desk references the policy when responding to a customer's variation question. Without a pre-established threshold, the institution's variation analysis becomes ad-hoc and the MRM committee's oversight loses calibration.

### Reproduction posture: chain-or-not + result schema + SOC test

The MRM committee establishes three additional load-bearing decisions before the institution conducts its first dispute reproduction:

**(1) Chain-or-not posture.** The reproduction is itself a chain entry. The recommended posture is to chain the reproduction with the original decision's `(run_id, seq)` as the reproduction's `parent_run_id` / `parent_seq` (per spec §4.4 linkage fields, integrity-bound under the per-event MAC per spec §5). A reproduction logged outside the chain is repudiable in the same way a vendor's logs are repudiable, which is the gap the chain exists to close. The MRM committee opines on whether the institution chains reproduction entries; the recommended decision is to chain.

**(2) Reproduction-result schema.** When the reproduction is chained, the chain entry's `audit.*` payload carries a documented schema for the reproduction result. RECOMMENDED schema fields under `audit.reproduction.*`:

| Field | Type | Notes |
|---|---|---|
| `original_run_id` | string | The disputed decision's run_id |
| `original_seq` | int64 | The disputed decision's seq |
| `iterations` | int | Number of re-runs the institution conducted |
| `decision_equivalent_count` | int | Number of re-runs producing the same final decision |
| `decision_equivalent_threshold` | float | The institution's MRM-defined tier-specific threshold (e.g. 0.95 for high-stakes) |
| `outcome` | enum | `decision_consistent` (100% equivalence), `variation_within_threshold` (≥ threshold), `variation_exceeds_threshold` (< threshold) |
| `mrm_committee_review` | string \| null | Reference to MRM committee review if `outcome = variation_exceeds_threshold` |

The institution adapts the schema to its MRM program; the chain integrity-binds whatever the institution puts in `audit.*`.

**(3) SOC procedure for testing reproduction evidence.** The audit-procedures doc carries P-27 (customer-dispute reproduction evidence completeness) — the SOC team confirms reproductions during the reporting period have the four-piece evidence package (chain integration, parent linkage, documented schema, MRM-review link for variation-exceeds-threshold cases). Without P-27, the SOC engagement cannot mechanically test reproduction evidence; the MRM committee's reproduction posture is unaccountable across reporting periods.

The institution's reproduction posture is documented in its MRM policy framework with these three decisions named explicitly. The customer-dispute desk references the documented posture when responding to customer variation questions; the SOC team consumes the chain entries with `audit.reproduction.*` fields against P-27 during the engagement.

## Hallucination cross-check

Reproducibility answers "would the AI have said the same thing on a different day." It does not answer the harder question the customer is increasingly likely to raise: "the AI cited fact X about my account in its denial reasoning; was X actually in the bank's records, or did the AI hallucinate it?" The chain captures what the AI said with cryptographic integrity. Whether what the AI said was true is a different question that requires a cross-check against the institution's authoritative source-of-truth records (account systems, document repositories, customer-service note systems).

The spec scopes itself to integrity-of-AI-output. Hallucination detection is the institution's customer-dispute procedure, layered on top of the chain's integrity guarantee. The institution operates the cross-check; the chain integrity-binds the cross-check itself so the customer-side disclosure produces a coherent narrative.

### Procedure

When the customer disputes an AI-cited fact:

1. **Identify the disputed fact.** The customer-dispute desk records the specific factual claim from the AI's reasoning (e.g., "the AI's denial reasoning states that my account had three NSF events in the prior 90 days"). The desk references the chain entry containing the AI's reasoning by its `(run_id, seq)`.
2. **Cross-check against authoritative sources.** The institution queries its authoritative source-of-truth records for the disputed fact. The query targets the system of record for the asserted fact (core banking system for account events; document-management system for documentary claims; customer-relationship system for customer-history claims). The cross-check is a substantive records-review step, not a paperwork formality.
3. **Record the disposition.** The cross-check produces one of three outcomes: the fact is `confirmed` (the authoritative record agrees with the AI's claim), `contradicted` (the authoritative record disagrees with the AI's claim — the AI hallucinated), or `unverifiable` (the authoritative record is silent or ambiguous on the fact). The disposition is the basis for the institution's response to the customer.
4. **Chain the cross-check.** The institution emits a chain entry of type `audit.fact_verification.*` recording the cross-check, with the disputed decision's `(run_id, seq)` as `parent_run_id` / `parent_seq` (per spec §4.4 linkage fields). The cross-check entry is integrity-bound under the per-event MAC per spec §5; the parent linkage binds the cross-check to the specific disputed decision so the customer-side disclosure produces a coherent narrative (decision → cross-check → disposition).

### Schema

The fact-verification entry's `audit.*` payload carries a documented schema:

| Field | Type | Notes |
|---|---|---|
| `audit.fact_verification.disputed_fact` | string | The specific factual claim from the AI's reasoning, recorded verbatim as the customer or the institution's reviewer identified it |
| `audit.fact_verification.authoritative_source` | string | The system of record consulted for the cross-check (e.g., `"core-banking-system:fdr-v3"`, `"document-management:opentext-v16"`); the institution names the source per its system inventory |
| `audit.fact_verification.disposition` | enum | One of `confirmed`, `contradicted`, `unverifiable` |
| `audit.fact_verification.evidence_id` | string | Reference to the authoritative-source record used for the cross-check (e.g., the core-banking transaction-id, the document-management object-id); enables the SOC team and the customer's expert to retrieve the underlying evidence |

The institution adapts the schema to its records inventory; the chain integrity-binds whatever the institution puts in `audit.*`.

### Why parent linkage matters

The fact-verification entry's `parent_run_id` / `parent_seq` bind to the disputed decision's chain entry. The customer-side narrative is "the AI said X in decision D; the institution cross-checked X against authoritative source S; the disposition was Y." Without the parent linkage, the cross-check is an unattached record the customer cannot tie to the specific decision they disputed. With the linkage, the verifier walking the disputed run produces the cross-check entry as a child of the decision, and the customer's expert can mechanically reconstruct the narrative from the chain alone.

The MRM committee opines on the institution's fact-verification posture (which AI-asserted facts trigger cross-check, which authoritative sources are consulted for which fact classes, the SLA for cross-check completion). The customer-dispute desk references the MRM policy when responding to a customer's hallucination question. The SOC team's audit-procedures sample confirms that fact-verification entries during the period have the parent linkage, the documented schema, and the authoritative-source reference.

## FRE 803(6) hearsay exception and the chain's reliability framing

When a customer dispute escalates to litigation (the customer files suit, or an enforcement action turns into a contested hearing), the chain entries become evidence under federal rules of evidence. The customer-service team and the dispute investigators do not need to argue evidence law in the dispute itself — that is the institution's litigation counsel's work — but the team's documentation choices during the dispute affect what the chain proves later. This section names how the chain composes with the federal hearsay rules so the dispute team's records are written for the legal posture the institution will eventually take.

**The chain's content is hearsay on its face.** The AI's output ("the customer had three NSF events in the prior 90 days") is an out-of-court assertion offered to prove a fact about the institution's behavior. Under FRE 802 hearsay is generally inadmissible; the institution's chain evidence is admitted under exceptions to the hearsay rule rather than as non-hearsay.

**FRE 803(6) — business records exception — covers operational events.** The institution's system observed timestamps, system state, the routing decision, the verifier output. These are records of regularly conducted activity, made by a system the institution operates in the ordinary course, kept as part of the institution's standard record-keeping. They satisfy the FRE 803(6) elements directly: a record made at or near the time of the event, by someone with knowledge (the system), kept in the course of regularly conducted activity. The chain's seal-job records, the verifier output, the operational-events stream, the chain entries' metadata fields (timestamps, tenant identifiers, run identifiers) — these are clearly business records.

**Entries that record the AI's assertions need a reliability framing beyond 803(6).** The chain entries that capture model outputs and reasoning chains record what the AI asserted, not what the institution observed. The institution's litigation posture is that these entries are admissible under a residual reliability exception: the chain proves no post-hoc tampering, which is the foundation for admitting the AI's output as a reliable record of what the AI said. The chain's MAC + Merkle + HSM coverage prove the chain entry has not been altered since capture; that integrity property is the reliability hook the residual exception turns on. The institution's IT witness lays the foundation by walking the verifier procedure (12 steps per spec §7) and showing the entry passed all 12.

**The truth-vs-integrity distinction the customer-service team must understand.** The chain proves the model said X. The chain does NOT prove that X is true. This distinction is the load-bearing point in customer-facing disclosures and in dispute-investigation records.

For ECOA adverse-action disputes, the customer's required disclosure under 12 CFR 1002.9 includes the model's stated reasoning (what the AI said) plus the institution's human review of whether the reasoning was factually sound. The dispute team's investigation produces both: the AI's reasoning is retrieved from the chain entry; the human review is the institution's substantive cross-check (the hallucination cross-check above; the MRM committee's review of the model's reasoning quality; the institution's policy compliance review). The customer's adverse-action notice carries both — what the AI said AND the institution's review of whether what the AI said was a sound basis for the decision.

The chain entry containing the AI's reasoning is admissible as a reliable record of what the AI said; the institution's human-review record is admissible as a separate business record under FRE 803(6) covering the reviewer's substantive analysis. Together they answer the two questions the customer's adverse-action posture asks: what the AI decided and on what basis, and whether the institution's human reviewer concurred with the AI's basis.

**Practical guidance for the dispute team.** Write dispute-investigation records as if they will be admitted in litigation. Reference the chain entries by their `(run_id, seq)` so the litigation team can pull them later; record the institution's human review as substantive analysis, not as paperwork sign-off; preserve the verifier output for the dispute period so the IT witness can lay foundation without reconstruction at deposition time. The dispute team does not write evidence law into the customer correspondence — but the team's records should support the institution's eventual evidence posture without rewriting.

**Cross-reference.** `docs/litigation-support.md` (FRE 803(6), residual reliability exception, foundation testimony for the IT witness), spec §7 (verifier procedure the IT witness walks under direct examination), §"Hallucination cross-check" above (the human-review record that complements the AI-output record).

## GDPR Article 22 — automated decision-making and human review

GDPR Article 22 grants data subjects a right not to be subject to decisions based solely on automated processing — including profiling — that produce legal effects or similarly significant effects. Credit denial, fraud flagging, and loan-pricing adjustments routinely fall in scope. The right's substantive content is **meaningful human review**: when a customer objects, the institution must provide a human reviewer with the authority and the information to make an independent judgment about the decision. The chain does not satisfy Article 22 by itself, and it does not replace the human reviewer. The chain serves the institution's Article 22 posture in four discrete ways.

**(1) The chain proves what the AI decided and on what basis.** The chain entry for the disputed decision binds the AI's response, the model version (request-side and response-side), the decoding parameters, the retrieved context, and the institution's `audit.*` decision-class payload. The reviewer reading the chain entry sees what the AI actually produced — not a summary the dispute team typed into a ticket, but the integrity-bound capture under spec §4.1 MAC and §4.3 daily seal. The reviewer can therefore evaluate the AI's actual reasoning rather than a rendered version of it.

**(2) The chain supports the human review by providing tamper-evident reasoning.** The reviewer's substantive analysis depends on the AI's reasoning being trustworthy as evidence of what the AI said. The chain's MAC + Merkle + HSM coverage gives the reviewer the foundation to trust the captured reasoning as the AI's actual output, not a post-hoc reconstruction or a redaction. The reviewer's independent judgment is built on the reasoning the chain integrity-binds; without the integrity property, the reviewer cannot distinguish an authentic AI reasoning from a sanitized version of it.

**(3) The human reviewer makes the independent judgment.** Article 22's load-bearing requirement is that the review is meaningful — the reviewer has the authority to overturn the AI's decision and the information to make the call on independent grounds. The chain entry is a primary input the reviewer consults; the reviewer's judgment is not the chain's judgment. The reviewer's analysis covers the AI's reasoning quality (was the reasoning sound?), the factual basis (did the cited facts match authoritative records, per the §"Hallucination cross-check" procedure above?), and the institutional policy (does the decision align with the institution's documented credit, fraud, or pricing policy?). The reviewer's record of that analysis is itself a chain entry per the institution's posture (recommended: chained as a child of the disputed decision via `parent_run_id` / `parent_seq` per spec §4.4) so the reviewer's substantive work is integrity-bound and admissible alongside the AI's reasoning.

**(4) The chain entry is the foundation for the institution's dispute response.** When the institution responds to the customer under Article 22, the response describes the AI's decision, the reviewer's independent analysis, and the disposition (the AI's decision stands, is reversed, or is escalated). The chain entries — the AI's reasoning, the human reviewer's analysis, any fact-verification cross-check, any reproduction record — together form the evidence package that supports the response. The customer's right to challenge the response is preserved; the institution's response is grounded in integrity-bound evidence rather than reconstructed narrative.

**Customer-notification procedure.** The institution's privacy notice and the institution's adverse-action / Article 21-22 notices disclose the right to human review and the contact path for invoking it. The notices name: that the decision involved automated processing within the meaning of Article 22, that the customer has a right to obtain human intervention, that the contact path is the institution's dispute desk (telephone, secure-message, written request), and that the institution will respond within the regulatory window. The notice does not require disclosure of the chain itself — the chain is institution-internal evidence — but the notice does commit the institution to the substantive review and to a documented response.

**Practical guidance for the dispute desk.** When a customer invokes Article 22, the dispute desk pulls the affected `(tenant_id, run_id, seq)` via the customer-correlation index, retrieves the chain entry for the AI's decision, and routes the matter to the human reviewer. The reviewer's analysis is recorded as a chain entry under the institution's `audit.*` review schema (see §"Hallucination cross-check" above for the parent-linkage shape). The institution's response to the customer references the reviewer's disposition and is delivered per the institution's documented response window. The chain provides the substrate; the human reviewer provides the meaningful review Article 22 requires.

**Cross-reference.** `docs/privacy-by-design.md` (lawful-basis analysis under Article 6 and Article 9 that covers the chain's processing of automated-decision data), §"Hallucination cross-check" above (the fact-verification step the human reviewer consults), §"Adverse-action notices (ECOA)" above (the chained notice cycle the Article 22 disposition feeds into when the decision is an adverse-action under ECOA).

## State-law considerations

State laws (Virginia VCDPA, Colorado CPA, California CCPA, etc.) impose various consumer-rights regarding AI-driven decisions:

- Right to know the decision was AI-driven (some states)
- Right to request human review (some states)
- Right to deletion of personal information used (CCPA)
- Right to access the data used (CCPA)

The chain accommodates these:

- The bank's notice procedure incorporates state-specific disclosure requirements
- The chain provides the integrity-bearing record of the decision
- The privacy-store responds to deletion requests (`docs/privacy-by-design.md`)
- The chain provides the data-source for access requests via the customer-correlation index

State counsel advises on state-specific requirements; the chain provides the substrate.

## CFPB-specific procedures

When CFPB inquiries arrive:

1. The bank's chain-ops team produces the affected runs
2. The bank's compliance team prepares the response per CFPB's standard format
3. The chain output (verifier walk for the affected runs) is included as evidence
4. The CFPB independently verifies if needed

CFPB examiners use the same verifier the FFIEC examiners use; the artifact is identical. The CFPB-specific examiner-orientation framing lives at `docs/regulator-pack/cfpb-overlay.md` (added per Round-17 CFPB-N2 close-out — the close-out directive collapsed v1.1 candidates into v1.0a companion docs, so the previously-deferred CFPB section is part of the regulator-pack alongside `nydfs-part500-overlay.md`).

### CID-class production and consumer-keyed retrieval (Round-17 CFPB-G2)

A CFPB Civil Investigative Demand commonly reads "produce all adverse-action decisions for consumers in [ZIP X] during [period Y]." The chain itself is keyed by `(tenant_id, run_id, seq)` per spec §3 and §4.1 — correct for integrity but not directly answerable as a consumer-keyed query. The institution's customer-correlation index (CUEC) is the retrieval substrate that maps the consumer-keyed question into the chain-keyed answer.

At v1.0a the CUEC was institution-internal and unchained. Round-17 CFPB-G2 surfaced the gap: a CID response built off an unchained, institution-controlled index leaves a gap between what the chain proves (integrity over the entries the institution chose to surface) and what the CID requires (a complete enumeration of consumers matching the criteria). Spec §10.23 closes the gap by defining two acceptable shapes for CUEC integrity, both of which make consumer-keyed retrieval testable from the chain alone.

**Shape 1 — Chain-anchored index (recommended).** Each CUEC entry is a chain entry under `chain_kind = "operational"` carrying `consumer_index.consumer_id_hash` (SHA-256 of the canonicalized consumer identifier per institution policy), `consumer_index.run_id`, `consumer_index.seq`, `consumer_index.relationship`. The Bureau's verifier reconstructs the index at any historical moment by replaying the operational events through the period in question.

**Shape 2 — Index attestation (acceptable for high-volume institutions).** The institution emits a daily `consumer_index.attestation` operational event under spec §10.2 carrying `index_snapshot_sha256`, `consumer_count`, `coverage_period_start_utc`, `coverage_period_end_utc`. The Bureau's verifier independently recomputes the index hash from the chain and compares against the attestation.

The institution names which shape it operates in CC8.1 per spec §10.18. A CFPB CID response under either shape is testable from the chain alone — the Bureau does not have to take the institution's index-construction discipline on trust. Audit-procedure P-6 (anomaly review) surfaces any discrepancy between the institution's CUEC output and the chain's recomputation as a control failure. Cross-reference spec §10.23 for the binding section.

## Customer-side verification path

For customer-side independent verification (typically engaged by the customer's counsel):

1. The customer's counsel requests a ledger snapshot covering the affected runs
2. The bank produces the snapshot per discovery procedure (`docs/legal-disclosure.md`)
3. The customer's expert verifier runs against the snapshot using the bank's public key
4. The expert's report describes integrity findings

The bank should not resist this independent verification; the design intent is exactly that customers (with counsel and expert) can verify. The institution's response is to facilitate, not impede, independent verification.

### IKM access for customer-side verification

The verifier's full per-event integrity check requires the IKM (or an HSM-mediated derivation path) per spec §7 step 7-9. Without it, the verifier degrades to structural verification (chain links, Merkle, signature) and reports `structurally consistent, key-bound verification skipped` per spec §7 fail-closed semantics. For a customer-dispute context, the institution operates one of two acceptable shapes for IKM access:

1. **Protective-order disclosure (litigation context).** Per `docs/legal-disclosure.md` §"Court-ordered master-key disclosure," the IKM is disclosed under a court-issued protective order limiting use to the verification activity, requiring secure handling, and binding the customer's expert to non-disclosure. This is the conservative posture; appropriate when the dispute is in litigation.
2. **HSM-mediated verification (operational context).** The institution provides the customer's expert with an HSM-mediated verification path (the expert's verifier dispatches HMAC operations through the institution's HSM API for the affected runs without ever holding the IKM bytes). The expert's report attests to the verification using the institution's HSM as the cryptographic delegate. Appropriate for operational disputes (CFPB inquiries, regulatory complaints) where the institution wants to preserve IKM custody.

The institution's customer-dispute procedure MUST name which shape it operates and the threshold at which it shifts (typically: operational shape for CFPB inquiries and regulator complaints; protective-order shape only when the dispute escalates to litigation). Without naming the shape, the institution risks the customer's expert encountering `unknown_key_version` at step 7 mid-verification, which leaves the verification incomplete and the dispute unresolved.

The MRM committee opines on the institution's chosen shape per its overall risk tolerance for IKM disclosure. The institution documents the decision in its control description and reviews it annually.

### Bureau-mediated verification (normative — Round-17 CFPB-P3)

Round-17 CFPB-P3 surfaced an operational fragility in the consumer-side IKM access posture: a consumer's expert running witness mode and reporting `Status: PASS-STRUCTURALLY, key-bound verification skipped` does not know whether the per-event MAC would also verify, and the consumer's lawyer reading that output should not have to litigate to find out. The two shapes above (protective-order disclosure, HSM-mediated verification) cover most cases, but a consumer in an active dispute often does not have standing to invoke either shape unilaterally. The Bureau-mediated path closes the gap.

**The Bureau-mediated shape.** The Bureau (CFPB Office of Supervision Examinations or, where the dispute has been referred to a federal court, a court-appointed master) holds the IKM under a standing protective order and runs the full §7 verification on the consumer's behalf when witness-mode `PASS-STRUCTURALLY` is insufficient for the dispute. The Bureau's role is delegated cryptographic verification — the consumer's expert prepares the chain artifacts and the verifier invocation; the Bureau executes the IKM-bound steps (§7 steps 7, 8, 9) the witness-mode procedure skips; the Bureau reports the per-event MAC result to both parties. The institution's role is providing the IKM under the protective order and the verifier binary the Bureau runs.

**When the Bureau-mediated shape applies.** The Bureau-mediated shape applies when ALL of the following hold: (a) the consumer has filed a complaint with the CFPB or a federal court has referred the matter; (b) the consumer's expert has run witness mode and reported a structural PASS; (c) the consumer's representative believes the per-event MAC verification is material to the dispute resolution; (d) the protective-order shape is unavailable or impractical (the institution declines to disclose IKM directly to the expert and the HSM-mediated shape is unavailable for the affected runs).

**Composition with the existing two shapes.** The Bureau-mediated shape does NOT replace the protective-order or HSM-mediated shapes — it is a third path the institution's CC8.1 control description names alongside the existing two. The institution names the threshold at which the Bureau-mediated shape becomes the operative path; for most institutions the threshold is "active CFPB enforcement action or federal-court referral with the consumer represented by counsel." Disputes below that threshold continue under the existing two shapes.

**Why this matters at v1.0b.** The chain's `Status: PASS-STRUCTURALLY` output is real integrity property under §7 witness mode, but a non-technical consumer or a court evaluating the witness-mode output cannot tell from the output alone whether the per-event MAC would also verify. A Bureau-mediated full verification produces `Status: PASS` (or `Status: FAIL` with a named §7 step) — the same output an institutional self-verification would produce — and disambiguates the structural PASS from a key-bound PASS without requiring the consumer to obtain IKM access directly. The Bureau's standing protective order provides the IKM-custody discipline; the consumer's expert remains the party preparing the chain artifacts and interpreting the verification output for the dispute resolution.

**Cross-references.** Spec §7 witness-verifier mode; spec §10.12 verifier CLI exit-code contract; `docs/legal-disclosure.md` §"Court-ordered master-key disclosure" for the protective-order discipline; `docs/regulator-pack/cfpb-overlay.md` for the Bureau-side examiner-orientation framing.

## Documentation for dispute response

The bank's chain-related response materials:

- The chain's `walk` output for the affected run(s)
- The institution's customer-correlation-index entry mapping customer to runs
- The institution's standard response language (compliance team)
- Any chain-specific incident records during the affected period
- The institution's SOC report (if applicable)
- The institution's regulator examination findings (when public)

## Customer-side communication template

For a routine direct dispute, the bank's response language might be:

> Thank you for raising your concern about the [decision] made on [date]. Our records show that the decision was made by our [AI agent / automated system] based on [factors]. Our internal records of this decision are tamper-evident; we have verified that the record is authentic and unaltered.
>
> [If the bank determines the decision was correct:] Based on our review, the decision was consistent with our policies and procedures.
>
> [If the bank determines the decision should be revisited:] We have reviewed the decision and will [specific action: re-evaluate, escalate to human review, etc.].
>
> If you have additional concerns, you may [next step: escalate to senior management, file with regulator, etc.].

The institution's customer-service team adapts the language; the chain provides the evidence base.

## Audit considerations

For audit teams reviewing dispute response:

- Sample disputes from the period; confirm the bank's chain-correlation-index produced the right runs
- Confirm the bank's compliance procedure for adverse-action and similar consumer notices
- Confirm the bank's response timing meets regulatory SLAs
- Confirm any CFPB or other regulator inquiries during the period had timely response

The chain is one input to the dispute response; the institution's broader compliance procedures complete the response.

---

## EU consumer-rights extensions — disclosure form, verification access, and asymmetric-evidence balancing

The procedures above cover U.S. dispute response. EU consumer-protection law imposes substantive obligations the U.S. framework does not — GDPR Article 15 right of access, AI Act Article 86 right to explanation, Charter Article 47 effective remedy, and the CJEU's consumer-procedural jurisprudence on evidentiary fairness in disputes between consumers and traders. The institution operating in EU consumer markets handles disputes under the framework below in addition to the U.S. procedures above. This section is the consumer-side mirror of the institution-side compliance posture in `docs/privacy-by-design.md`; the two documents together form the symmetric pair the institution's product counsel, the consumer's litigator, and a national consumer-protection authority each need.

### Data-subject disclosure form under GDPR Article 15

GDPR Article 15(1) gives the data subject the right to obtain confirmation of processing, access to the personal data, and meaningful information about automated decision-making logic. Article 15(3) adds the right to a copy of the personal data undergoing processing. SCHUFA Holding (CJEU C-634/21, 7 December 2023) held that automated credit scoring is itself an Article 22(1) decision, and the corresponding Article 15(1)(h) obligation requires meaningful information about the procedure and principles concretely applied — not a generic algorithm description.

When an institution operating the chain receives an Article 15 request from a data subject whose AI-driven decision is captured in the chain, the institution's disclosure includes:

| Disclosure component | What the institution provides | Purpose |
|---|---|---|
| Captured JSON form | Every chain entry concerning the data subject in its captured JSON form per spec §5.2 best-evidence posture | The content-bearing form. The data subject's expert reads the AI's prompt, response, decision-class payload, and decoding parameters as captured. |
| Canonical bytes | The canonical-bytes form of the same entries (per spec §5) | The integrity-bearing form. The data subject's expert can independently reproduce the per-event MAC compute when the institution makes its IKM verification accessible. |
| Seal record(s) | The seal record(s) covering the tenant-day(s) the disclosed entries belong to (`merkle_root`, `algorithm`, `public_key_id`, `signature`, `signed_at`) | Closes the seal layer of the integrity claim. The data subject's expert verifies the seal's HSM signature against the institution's published public-key registry. |
| Public-key resolution | The resolution of `public_key_id` to the Ed25519 public-key bytes, or a pointer to the institution's published tenant-key registry | Lets the data subject's expert verify the seal signature without trusting the institution's claim about the key. |
| Verifier-run output | The institution's verifier output for the disclosed period, showing `Status: PASS` (or the documented anomaly disposition with the `regulator-pack/finding-language.md` row reference) | The institution's own demonstration that the disclosed entries pass the verifier procedure. |
| Witness-runnable form | The same artifacts in a form the §7 witness-verifier can run against — the on-disk NDJSON form of the entries, the seal records, and the public-key resolution | Lets the data subject's independent expert run §7 witness mode without IKM access and without institutional veto. |

The disclosure is delivered within Article 12(3)'s one-month timeline (extendable by two months for complex requests with notification to the data subject). The first copy is free of charge per Article 15(3) and CJEU FT (C-307/22, 26 October 2023); subsequent copies may carry a reasonable fee based on actual administrative cost, NOT a fee that effectively makes verification unaffordable. The institution's CC8.1 (or EU-equivalent) control description names the disclosure form, the cost framework, and the procedure that surfaces the artifacts within the regulatory window.

The institution's customer-dispute desk produces the disclosure from the customer-correlation index — same retrieval path as the U.S. dispute procedure above — extended to assemble the witness-runnable form rather than only the institution-internal verifier output.

### Independent consumer-side verification — making §7 witness mode accessible

Spec §7 defines a witness-verifier mode operating without `--master-key`: the verifier executes the structural and seal-verification steps and skips the IKM-bound steps, producing `Status: PASS-STRUCTURALLY, key-bound verification skipped`. The mode confirms no insertion, deletion, reordering, format-version mismatch, header-tampering, cross-chain lift, or seal-signature forgery — and confirms the daily Merkle root recomputation matches the HSM-signed seal. This is real integrity property that a competent consumer-side expert can verify without ever holding the institution's IKM.

The institution's commitment, recorded in its CC8.1 control description: **the institution disclosing chain entries to a data subject under Article 15 includes the artifacts required to run §7 witness mode against the disclosed entries, and the institution does not gate access to the verifier binary or the §7 procedure**.

The reference verifier ships under Apache 2.0 from the project repository. The institution does not need to provide it; the data subject's expert obtains it independently. The institution's role is to provide the artifacts in a runnable shape — the on-disk NDJSON entries (with canonical bytes preserved), the seal records, and the public-key resolution. With those three inputs the expert runs witness-mode verification and produces an independent integrity report admissible as expert testimony in EU national courts under each country's procedural code.

The institution publishes (or links to in its disclosure cover letter) a "consumer-side verification guide" naming the §7 procedure, the witness-mode invocation, the expected `Status: PASS-STRUCTURALLY, key-bound verification skipped` output, and the meaning of any failure-reason string. The guide reduces the consumer-side cost of verification from "hire a cryptographer to read the spec" to "run a documented procedure against the institution's disclosed bundle."

### AI Act Article 86 right-to-explanation mapping

EU AI Act Article 86(1) gives any affected person subject to a decision taken by a deployer on the basis of a high-risk AI system's output (Annex III, including credit-scoring under point 5(b)) the right to obtain clear and meaningful explanations of the role of the AI system in the decision-making procedure and the main elements of the decision taken. The institution's response under Article 86 maps to chain attributes as follows:

| Article 86 obligation | Chain attribute(s) | Institution's narrative role |
|---|---|---|
| Role of the AI system in the decision-making procedure | `audit.routing.provider_chosen`, `audit.routing.policy_version`, `audit.deployment.intent`, `audit.deployment.policy_version` | The institution's documented description of the AI's role under its CC8.1 control description, surfaced to the data subject in plain language. |
| Main elements of the decision taken | The entry's `audit.*` decision-class payload, `gen_ai.response.model`, `gen_ai_parameters` (per spec §4.4 RECOMMENDED contents) | The factual record of what the AI decided, on what model version, under what decoding parameters. |
| Human review (when invoked under Article 22) | The chained human-reviewer entry (via `parent_run_id` / `parent_seq` per the GDPR Article 22 section above) | The substantive analysis the institution's reviewer produced. Integrity-bound under the per-event MAC; admissible alongside the AI's reasoning. |
| Hallucination cross-check (when invoked) | The `audit.fact_verification.*` chained child entry (per the §"Hallucination cross-check" section above) | The institution's cross-check disposition (`confirmed`, `contradicted`, `unverifiable`) bound to the disputed AI claim. |

The dispute desk uses the mapping as the baseline content of its Article 86 response. The narrative is plain-language; the chain entries are the integrity-bound substrate the data subject's expert (or counsel's expert) can verify. The institution does not over-claim — the chain proves what the AI decided; the institution's separate reasoning-quality and policy-compliance evidence proves whether the decision was correct.

### SCHUFA C-634/21 — score-producing decisions in the chain

SCHUFA Holding (CJEU C-634/21, 7 December 2023) held that an automated probability value concerning a person's ability to meet payment obligations is itself an Article 22(1) automated decision when third parties draw on the value for consequential downstream decisions. The judgment changed what counts as an Article 22 decision in EU credit-decision systems: not only the bank's final credit-grant or credit-deny is an Article 22 decision, but the upstream score the bank consumed.

The institution's chain captures both the score-producing AI call (when the institution operates the score-producing pipeline) and the bank's decision based on it (when the institution operates the consuming pipeline). When a data subject invokes Article 22 on a credit denial:

1. The dispute desk identifies the bank-side decision entry via the customer-correlation index.
2. The desk identifies any score-producing entry chained as an upstream parent (`parent_run_id` / `parent_seq` may link to the score-producing entry when the institution operates a multi-process pipeline per spec §4.4 cross-run linkage).
3. The desk surfaces both entries in the data-subject's disclosure under Article 15.
4. The human-reviewer record per the GDPR Article 22 section above is required for both — the score-producing decision and the bank-side decision are separate Article 22 decisions, each with a separate meaningful-review obligation.

When a score-producing AI pipeline is operated by an upstream vendor (a credit bureau, a third-party scoring service), the upstream vendor is the controller for the score-producing decision under SCHUFA Holding paragraph 49. The bank's chain captures the bank's consumption of the score; the bank's data subject's Article 15 request to the bank produces the bank-side chain entries. The upstream score-producing entries are obtained from the upstream vendor under a separate Article 15 request to the upstream controller. The institution's customer-dispute desk informs the data subject of both paths when the institution's pipeline involves an upstream vendor's scoring service.

### Charter Article 47 — effective remedy and asymmetric evidence

Charter of Fundamental Rights Article 47 gives every person whose rights guaranteed by EU law are violated the right to an effective remedy and to a fair hearing. The CJEU has consistently read Article 47 as imposing substantive procedural-fairness obligations on the EU and Member States, particularly in proceedings between consumers and traders (Profi Credit Polska C-176/17; EOS KSI Slovensko C-448/17; Banco Santander C-105/17). Effective remedy in this context means the consumer can in fact challenge the institution's decision; an evidentiary regime that requires every consumer to retain a cryptographic expert to verify the chain artifacts the institution produces effectively forecloses the remedy for low-value disputes.

The institution's commitment under Article 47:

1. **Witness-runnable disclosure as default.** The institution disclosing chain artifacts to a data subject does so in a form the §7 witness-verifier procedure can run against, without requiring the data subject to negotiate IKM access or receive only a printed summary. The CC8.1 control description names the disclosure shape; deviation requires documented justification.

2. **Court-appointed-expert cost commitment.** When a national court orders a court-appointed expert under each Member State's procedural code to verify the chain artifacts in a consumer dispute (a standard mechanism under Profi Credit Polska's effective-remedy doctrine), the institution bears the cost of that expert. The institution does not contest the appointment on cost grounds. The expert runs the witness-verifier and produces an independent integrity report.

3. **Procedural-language translation.** When the proceeding is in a non-English EU forum, the institution provides the verifier output and the spec §7 step explanations translated to the consumer's procedural language (BCP 47 codes per spec §10.11). The translation does not modify the chain entries — the chain stays canonical — but the consumer-side reading of the verifier output is in the language of the proceeding.

4. **Unfair-terms posture.** The institution's terms-of-service drafting does NOT include clauses purporting to require the consumer to bear cryptographic-expert costs in a chain dispute, render the chain's records "conclusive evidence," or waive the consumer's right to challenge the records' integrity, accuracy, or completeness. Such clauses are presumptively unfair under Directive 93/13/EEC Article 3(1) and the indicative list at point 1(q) of the Annex (excluding or hindering the consumer's right to take legal action). The institution's product-counsel review under EU consumer-market deployment screens out such clauses.

5. **Cooperation with national consumer-protection authorities.** The institution provides the verifier binary and the §7 procedure documentation to qualified consumer-protection authorities (DGCCRF, Verbraucherzentrale, Test-Achats, BEUC, UFC Que Choisir, Altroconsumo, equivalents) under the cooperation procedures those authorities operate. National authorities can run the verifier on aggregated complaints data within their own infrastructure without case-by-case access to the institution.

### Brussels I bis cross-border consumer jurisdiction

Brussels I bis Regulation (Reg 1215/2012) Article 18(1) gives EU consumers the right to sue the institution in the consumer's home Member State court. Article 19 limits parties' freedom to derogate — a forum clause in a consumer contract is essentially unenforceable against the consumer. A French consumer suing a German bank gets jurisdiction in France; a Dutch consumer suing a Belgian fintech gets jurisdiction in the Netherlands.

Spec §10.15 Pattern B (per-region `tenant_id`) supports the Brussels I bis posture by aligning each EU consumer's chain entries with a `tenant_id` whose seal region, IKM custody, and ledger storage are within the consumer's home Member State. The institution's production of chain artifacts in the consumer's forum:

- Pulls only the `tenant_id` corresponding to the consumer's home Member State.
- Avoids triggering Article 44 GDPR cross-border transfer questions on the chain disclosure itself — the disclosed chain is already pinned to the consumer's home jurisdiction.
- Meets the consumer's home forum's procedural timeline (typically 14-60 days under each Member State's civil procedure rules).

The institution's CC8.1 control description names the per-Member-State `tenant_id` allocation, the seal-region pinning per Member State, and the institution's commitment to produce chain artifacts within the procedural timeline of the consumer's home forum. Pattern A (single seal region) is conformant for EU institutions whose risk posture admits cross-Member-State seal aggregation, but the institution's litigation-support team must address cross-border production timeliness and Schrems II-style transfer analysis when the consumer's home forum is in a Member State whose data did not reach the seal region in compliant form.

### Representative actions under Directive 2020/1828

The Representative Actions Directive (Dir 2020/1828, in force since 25 June 2023) enables qualified consumer-protection associations to bring representative actions on behalf of large groups of consumers seeking injunctive or compensatory measures across borders. Representative actions in financial-services contexts often involve thousands of consumers affected by the same allegedly unlawful AI-decision pattern (a credit-scoring model with disparate-impact characteristics, a fraud-flag system with high false-positive rates against a protected class).

The institution's response posture for representative actions:

1. **Sample-aligned production.** The chain is keyed by `(tenant_id, run_id, seq)`, not by consumer identity. The institution's customer-correlation index resolves a sample of N consumers to chain entries. The institution produces the sampled entries plus the daily seal records covering them — typically via the partial-disclosure mode (`docs/selective-production-and-sampling.md`).

2. **Witness-mode verifier output for the sampled corpus.** The institution provides a verifier-run output across the sampled extracts. The qualified entity's expert runs the witness-mode verifier independently per the asymmetric-evidence-balancing posture above; the qualified entity does not need the IKM and is not gated by institutional veto.

3. **Aggregated statistics on the pattern at issue.** The institution produces — under court order and per the qualified entity's expert's specification — aggregated statistics on the pattern (decision-distribution by demographic cohort, false-positive rate, model-version cohort breakdown) that preserve the chain's integrity claim while supporting the representative-action argument about institutional behavior at population scale.

4. **Cross-jurisdictional production.** When the representative action is brought across Member States (Article 6 of the Directive on cross-border representative actions), the institution coordinates production across the relevant `tenant_id`s under Pattern B; each Member State's affected consumers are produced from their home `tenant_id`.

5. **Cost commitment.** The institution provides the artifacts in witness-runnable form; the qualified entity runs the verifier within its own infrastructure. Representative-action production does NOT require the consumer association to bear the cost of the verifier or the cryptographic-expert engagement.

The institution's regulator-pack maintains a representative-actions production guide; the FFIEC working group's outreach extends to the EU consumer-association network so the production posture is co-authored with the qualified entities under the Directive.

### Counterfactual evidence — what the chain does NOT capture

Spec §1.2 epistemic scope is explicit that the chain proves what the AI said and that the record was not tampered with; the chain does not prove the AI's statement is factually accurate, policy-compliant, or unbiased. EU consumer-protection proceedings require an additional non-claim acknowledgment, important for damages quantification:

**The chain does NOT capture the counterfactual decision** — what a non-AI process or an alternative AI process would have produced. In consumer-protection proceedings where damages quantification turns on the counterfactual (a representative action alleging disparate impact; a customer claim that "had your bank used a different model my loan would have been approved"), the counterfactual evidence is produced from the institution's separate evidence regimes:

- Model-validation records — what the institution's MRM committee's testing showed about alternative models.
- Fairness-audit records — what the institution's disparate-impact testing showed about the deployed model versus alternatives.
- Customer-dispute investigation records — what the institution's human reviewer concluded when reviewing the disputed decision.

The chain composes alongside these regimes without subsuming them. The consumer's expert testifying from the chain alone cannot opine on the counterfactual — the chain does not adjudicate it. The institution's expert testifying from the chain alone cannot defeat a counterfactual claim — the chain does not adjudicate that either. Both sides bring counterfactual evidence from outside the chain.

### Article 15(3) cost framework

GDPR Article 15(3) third sentence: the controller "shall provide a copy of the personal data undergoing processing. For any further copies requested by the data subject, the controller may charge a reasonable fee based on administrative costs." Recital 63 elaborates that the right of access should not adversely affect the rights or freedoms of others, including business secrets or intellectual property. CJEU FT (C-307/22, 26 October 2023) further clarified the "free first copy" rule and the reasonableness limit on subsequent copies.

The institution's cost framework for chain disclosure:

| Disclosure occasion | Charge |
|---|---|
| First copy of chain artifacts associated with the data subject's personal data — captured JSON, canonical bytes, seal records, public-key resolution, verifier-run output | Free of charge per Article 15(3). |
| Subsequent copy (e.g., re-disclosure six months later) | Reasonable fee based on actual administrative cost; capped at the institution's documented per-disclosure cost. NOT a fee that effectively makes verification unaffordable. |
| Verifier-run output as a separate service | The verifier run is automated and the marginal administrative cost is near zero; it is part of the Article 15(3) free-of-charge entitlement on first copy. The institution does not charge a marginal fee for the verifier output alone. |
| Translation to the data subject's procedural language | When the proceeding is in a non-English EU forum, included in the institution's first-copy disclosure per the Article 47 commitment above. |

The institution's CC8.1 (or EU-equivalent) control description names the cost framework, including the institution's posture on subsequent-copy fees and the threshold above which a fee is reasonable. A national consumer-protection authority reviewing the institution's framework consults the institution's CC8.1 and the institution's actual subsequent-copy fee history for proportionality.

### Dual-algorithm transitional disclosure to data subjects

Spec §7 step 11 case (e) — dual-algorithm seal where one signature verifies and the other does not — is a Severe finding regardless of the verifier's PASS-WITH-ANOMALY bracket under non-strict mode. When this case is reported in a verifier-output disclosure to a data subject under Article 15 or to a qualified entity under Directive 2020/1828:

1. The disclosure includes the case (e) "Severe MRA" finding-language reference verbatim from `regulator-pack/finding-language.md` row 11, so a non-cryptographer reading the disclosure understands the seal is in a Severe-finding state regardless of the PASS-WITH-ANOMALY bracket.
2. The institution's customer-dispute desk explanation accompanying the disclosure uses plain-language summary: *"Two signature algorithms cover this seal. One verifies; the other does not. This is a Severe finding our regulator treats as material; we are investigating and will report results to the relevant supervisory authority and to you under the regulatory window."*
3. The non-strict PASS-WITH-ANOMALY bracket reflects the spec's posture that the un-broken algorithm still provides integrity assurance for downstream consumers; it does NOT diminish the finding's severity for examination, regulatory, or consumer-protection purposes.

The institution does not lean on the PASS-WITH-ANOMALY bracket to elide the Severe finding in a customer-side disclosure. Doing so would be an asymmetric-evidence move the CJEU's consumer-procedural jurisprudence would not tolerate.

### Cross-reference

`docs/privacy-by-design.md` — institution-side compliance posture for GDPR Article 6, 9, 22, 35 (DPIA), and the lawful-basis analysis for chain processing of automated-decision data. `docs/regulator-pack/charter-art-47-effective-remedy.md` and `docs/regulator-pack/eu-unfair-terms-and-the-chain.md` — the regulator-pack documents national consumer-protection authorities consult alongside the institution's CC8.1. `docs/regulator-pack/eu-representative-actions.md` — the representative-actions production guide for qualified entities under Directive 2020/1828. `docs/litigation-support.md` — the U.S. litigation-side framing; the EU consumer-side procedure here composes with the U.S. discovery posture for institutions operating in both markets.
