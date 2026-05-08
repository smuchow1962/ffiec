# Outside reviewer — SOC 2 / ISAE 3000 attestation partner

**Reviewer.** Marcus Brennan. Big-4 SOC 2 / ISAE 3000 attestation partner based in Chicago, with satellite teams in London (ISAE 3402) and Bangalore (ISAE 3000 / SOC 2). Background: fifteen years SOC 1 and SOC 2 attestation engagements for cloud-services providers, financial-services SaaS, and health-tech platforms; five years on the AICPA SOC for Cybersecurity technical committee; co-author of the AICPA practice aid on SOC for AI Systems (2025 draft) and the related ISAE 3000-AI proposal under IAASB consultation. Holds CISA, CISSP, CPA, and the AICPA Information Management & Technology Assurance specialization. Manages the firm's AI-attestation practice — currently engaged on three SOC 2 Type II engagements where the institution's scope includes the FFIEC chain-of-custody infrastructure as an in-scope subservice or in-house control set.

**Why this review matters.** Marcus's firm is engaged on three SOC 2 Type II attestation engagements where chain-of-custody infrastructure is in scope. His questions are practical: can the controls described in the spec be tested? Is the evidence the institution would produce sufficient for a SOC 2 opinion? Do the spec's prescribed control descriptions support a per-Trust-Service-Criteria mapping that withstands AICPA peer review? His sign-off is required for the engagement opinion; if the chain-of-custody controls aren't testable, the opinion carries scope limitations the institution does not want.

**Reading angle.** Practical control-testing lens. Marcus reads the spec asking: "If I were the SOC auditor on this engagement, can I test this control? What evidence would I sample? What procedures would the institution write? Would my opinion stand up to AICPA peer review and to the institution's customers reading the SOC 2 report?" His threshold for a Partial is a spec prescription that says "the institution does X" without prescribing what evidence the auditor samples, what frequency, what sampling methodology, and which Trust Service Criterion the prescription maps to. He stays clear of cryptographic evaluation (Reuven, drop #04), legal-evidentiary framing (Diego, drop #03), regulatory architecture (Henrik EU drop #08, Hiroshi Japan drop #06, the broader APAC drops), privacy law (Elena, drop #05), and AI safety research (Tomás, drop #08/#09). His domain is the audit engagement itself: scoping, control testing, evidence sampling, opinion drafting, peer review.

**Why this is distinct from Maya (production-ops + controls-framework parity).** Maya is on the institution side — she runs the controls and produces the evidence. Marcus is the external party who tests the controls and samples the evidence. His findings cluster around testability and evidence-sufficiency, not around control design. Where Maya asks "can my team operate this control?" Marcus asks "can I attest to this control with reasonable assurance under AICPA AT-C 105 / 205?"

**Stopping criterion.** Eight to twelve questions covering: per-TSC mapping (CC1-CC9 plus A1, C1, PI1, P1-P9 where applicable); auditor-side procedure for verifying §10.6.1 RNG attestation; verifier output as audit evidence (institution-run vs auditor-run vs third-party-run); §10.1 reconciliation sampling methodology; §10.7 software-key adapter exclusion test procedure; §10.5 HSM custody evidence requirements; subservice-organization carve-out vs inclusive method for cloud HSM; SOC for AI Systems alignment; control-deficiency disposition guidance for chain-related findings. Marcus does not file Partials for cryptographic gaps or regulatory-architecture gaps. His threshold is auditor-engagement gaps — places where the spec leaves the audit engagement to write its own procedures, producing inconsistent attestation rigor across institutions.

---

## Headline

The chain's design is unusually attestation-friendly: byte-reproducible verifier output (§7), normative reconciliation procedures (§10.1), append-only enforcement (§4.1 inviolate properties), prescribed FIPS 140-2 Level 3 HSM custody (§10.5), and an exit-code contract for the verifier CLI (§10.12). These are the kinds of testable controls SOC engagements need. **My read pass found ten places where the spec prescribes a control without prescribing how to test it, or names "CC8.1" as a controls home when the prescription actually maps to multiple Trust Service Criteria.** None of these are integrity-bearing — they're audit-engagement gaps that, left unaddressed, produce inconsistent attestation rigor across institutions and across attestation firms. The gaps cluster into three categories: (a) per-TSC mapping precision, (b) auditor-side procedure templates, and (c) subservice-organization vs in-house control delineation.

**What the spec does well for SOC engagements.** The verifier output is byte-for-byte reproducible — that is uncommon for compliance-tooling outputs and translates directly into AICPA "tests of controls" evidence. The reconciliation procedures at §10.1 (P-6 fingerprint reconciliation, count reconciliation, anomaly detection) name the operational evidence the institution produces and the cadence on which it produces it. The §10.12 verifier exit-code contract (0 / 1 / 2 / 3) gives auditors a clean dispositional signal during evidence sampling. The §4.4 attribute schema is documented, allowing test-of-controls procedures to refer to attributes by name rather than by reverse-engineering. The supplementary `audit-procedures.md` (referenced in the change-log) appears to address the institution's audit-procedure design, though I haven't seen it referenced in scope here.

**What the spec needs for SOC engagements.** Per-TSC mapping table; auditor-side procedure templates for each spec-prescribed control; subservice-organization carve-out vs inclusive-method guidance for cloud HSM and KMS; sampling methodology and frequency for §10.1 reconciliations; SOC for AI Systems criterion alignment; control-deficiency disposition rubric (when does a chain-side finding rise to a "significant deficiency" or a "material weakness"); evidence-sufficiency guidance for §10.6.1 RNG attestation; and a posture statement on third-party verifier-run independence (institution-run vs auditor-run vs independent-third-party-run).

Each question below is a one-paragraph gap with a concrete close-out path. None of the gaps are integrity-bearing; all are audit-engagement quality items that would produce more consistent attestation rigor across institutions deploying the chain.

---

## Supporting questions

### MQ-1. Per-TSC mapping precision — "CC8.1" is too narrow for what the spec actually prescribes

**Question.** The spec frequently says "the institution's CC8.1 (or equivalent) control description names X." CC8.1 is one specific Trust Service Criterion (Change Management — "the entity authorizes, designs, develops, configures, documents, tests, approves, and implements changes to infrastructure, data, software, and procedures"). The spec's prescriptions actually map to multiple TSCs:

- §10.5 HSM custody → CC6.1 (logical and physical access), CC6.2 (system resources protection), CC6.3 (network communications), CC6.6 (encryption controls), CC6.7 (physical access)
- §10.6.1 RNG attestation → CC6.1 (cryptographic key management), CC6.6 (encryption controls)
- §10.1 reconciliation → CC4.1 (control monitoring), CC4.2 (deficiency remediation), CC7.2 (system monitoring for anomalies)
- §10.7 software-key adapter exclusion → CC8.1 (change management), CC6.1 (logical access — preventing dev-key path in production)
- §10.10 algorithm rotation → CC8.1 (change management — algorithm rotation IS a change)
- §7 verifier procedure → CC4.1 (monitoring), CC4.2 (deficiency remediation)
- §10.13 evidentiary artifact retention → CC9.1 (risk mitigation through information retention), C1.1 (confidentiality categorization)
- Append-only enforcement (§4.1 inviolate property) → CC6.1 (logical access controls), CC8.1 (change management — preventing unauthorized modification)

The flat "CC8.1" reference is auditor-friendly only at a high level; an actual SOC 2 engagement requires per-TSC mapping for each prescription. Different attestation firms map differently, producing inconsistent SOC 2 reports across institutions deploying the same chain.

**Status:** Partial.

**Closing language proposal.** A new informative annex titled "Trust Service Criteria mapping" (or addition to the supplementary `soc-pack/control-evidence-events.md`) listing each spec-prescribed control alongside the AICPA Trust Service Criteria it maps to, plus the equivalent ISAE 3000 / 3402 control objective for non-US engagements. The mapping is informative — institutions and attestation firms are free to refine it for their engagement scope — but it provides the baseline mapping every engagement would otherwise write itself, producing inconsistent results.

### MQ-2. Auditor-side procedure for §10.6.1 RNG attestation — what evidence does the auditor sample?

**Question.** §10.6.1 requires the institution to record the RNG type on the `master_key.generated` operational event and document the choice in CC8.1. From the auditor's side, how is this attested? The auditor cannot open the HSM and inspect its CSPRNG. The auditor's procedure is necessarily indirect: sample the operational-event log, sample the CC8.1 documentation, sample the vendor's FIPS validation certificate, sample the HSM's RNG self-test logs (if available), reconcile across these. But the spec doesn't prescribe the auditor's procedure, so each engagement writes its own — and the rigor varies. A weak procedure (just sample the operational event and stop there) produces an opinion that doesn't actually attest to anything testable; a strong procedure (sample event + CC8.1 documentation + vendor SOC 2 + RNG self-test logs + key-ceremony witness reports) is rigorous but expensive. The spec should provide an auditor-procedure template.

**Status:** Partial.

**Closing language proposal.** A new section in `audit-procedures.md` (or a new file `auditor-procedures.md` parallel to `audit-procedures.md` — note the naming distinction: `audit-procedures.md` is the institution's procedure, `auditor-procedures.md` is the external attestor's procedure) titled "RNG-source attestation procedure (auditor-side)" naming the sampling: (1) sample the most-recent `master_key.generated` event for each tenant in scope, (2) verify the institution's CC8.1 documentation names the same RNG source, (3) request the vendor's FIPS validation certificate (FIPS 140-2 Level 3 or higher per §10.5), (4) where the institution's posture is HSM-internal RNG, request the HSM's RNG self-test logs from the previous quarter; where OS-level CSPRNG, request the OS-distribution's CSPRNG documentation and the institution's host-hardening procedure; (5) reconcile across (1)-(4) — discrepancies are control deficiencies the institution remediates per CC4.2.

### MQ-3. Verifier output as audit evidence — institution-run vs auditor-run vs independent-third-party-run

**Question.** §7 verifier output is byte-for-byte reproducible (§7's normative output format, §10.12's exit-code contract). SOC auditors should be able to use it as evidence. But: who runs the verifier during the SOC engagement? Three postures, with different evidence-reliability implications:

- **Institution-run.** The institution runs the verifier and provides the output to the auditor. Cheapest. Lowest evidence reliability — the institution could (in principle) run a tampered verifier that always outputs PASS. AICPA AT-C 205 §A21 addresses management-prepared evidence: the auditor evaluates its reliability, often via re-performance.
- **Auditor-run.** The auditor runs the verifier directly against the chain artifacts the institution produces. Higher evidence reliability — the verifier is the auditor's own controlled tool. Operationally heavier — the auditor needs the verifier binary, the institution's tenant public key, and read access to the chain ledger.
- **Independent-third-party-run.** A neutral third party (e.g., a witness verifier per §7) runs the verifier and produces an attestation. Highest evidence reliability. Operationally heaviest. Used in the highest-stakes engagements (e.g., regulatory examination + SOC 2 opinion combined).

The spec doesn't address which posture the institution operates under, which posture the auditor selects, or how the auditor evaluates institution-run output evidence under AT-C 205 §A21. Without guidance, attestation firms write their own procedures, producing inconsistent rigor.

**Status:** Partial.

**Closing language proposal.** A new section in `audit-procedures.md` (or `auditor-procedures.md` per MQ-2) titled "Verifier output as evidence" naming the three postures, the AICPA AT-C 205 §A21 evaluation procedure for institution-run output (auditor selects a sample of institution-run outputs and re-performs the verification), and the institution's CC8.1 attribute declaring which posture the institution operates under. The institution's choice of posture is a control-design decision documented in CC8.1; the auditor's evaluation procedure is fixed by the spec. This produces consistent attestation rigor across institutions deploying the chain.

### MQ-4. §10.1 reconciliation sampling methodology — what's the sampling unit and frequency?

**Question.** §10.1 (operational reconciliation, including P-6 fingerprint reconciliation) prescribes that the institution reconciles event counts, fingerprints, and other operational data on cadences. SOC engagement testing of these reconciliations requires a sampling methodology: per-tenant-per-day random selection? Risk-based selection (high-volume tenants more frequently)? Stratified by `chain_kind`? Census of high-risk tenants plus sample of others? The spec doesn't prescribe the auditor's sampling. AICPA AT-C 205 §A47 names sampling considerations but the standard is generic — engagement-specific sampling unit and frequency should come from the spec or from the supplementary attestation guidance.

**Status:** Partial.

**Closing language proposal.** Section in `audit-procedures.md` titled "Reconciliation testing — sampling methodology" naming the per-control sampling unit (the tenant-day reconciliation record) and a tiered frequency: census for tenants flagged high-risk in the institution's risk register, statistical sample (typically 60 tenant-days per period for a 95% confidence / 5% tolerable misstatement) for the broader population, plus a backstop "all reconciliations from the most recent month before report-issuance" sample to confirm the cadence is currently operating. The methodology is informative — engagements with broader scope or higher risk tolerance refine it — but it provides a defensible baseline.

### MQ-5. §10.7 software-key adapter exclusion — auditor test procedure

**Question.** §10.7 says "Production builds ship without the adapter assembly, package, or module on disk. The institution's deployment pipeline asserts the adapter artifact is absent in the production image (registry-side check, cosign-verified content manifest, or deployment-gating control)." This is a strong control design but the auditor's testing procedure is unstated. Possibilities: (a) the auditor inspects the production image directly via container-registry pull (operationally heavy, requires production access), (b) the auditor inspects the institution's deployment pipeline output and the cosign-verified manifest (relies on the institution's tooling integrity), (c) the auditor performs a code-search against the production source tree (relies on the institution's source-control integrity), (d) the auditor relies on the institution's CC8.1 documentation alone (low reliability). Different procedures, different evidence reliability, different engagement cost.

**Status:** Partial.

**Closing language proposal.** Section in `audit-procedures.md` titled "Software-key adapter exclusion — auditor test procedure" naming the preferred procedure (cosign-verified content manifest from the institution's deployment pipeline, plus a sample of production-image pulls from the institution's registry) and the alternative procedures with their reliability assessment. The institution's CC8.1 control description names which exclusion pattern it operates (compile-time vs packaging vs equally-strict alternative); the auditor's procedure tests the corresponding evidence. Consistent procedure produces consistent attestation.

### MQ-6. §10.5 HSM custody — subservice-organization carve-out vs inclusive method

**Question.** When the institution uses AWS CloudHSM, Azure Managed HSM, or Google Cloud HSM, the cloud provider is a subservice organization under AICPA TSP. The institution's SOC 2 report scopes the cloud provider via either (a) the carve-out method (the cloud provider's controls are explicitly excluded from the SOC 2 scope; the report names the cloud provider's complementary user-entity controls and points to the cloud provider's own SOC 2 / ISAE 3402) or (b) the inclusive method (the cloud provider's controls are included in the institution's SOC 2 scope by reference). The choice has substantial implications for engagement scope, complementary user-entity controls, and the report's usability for downstream consumers. The spec doesn't address this — institutions and attestation firms decide independently, producing inconsistent SOC 2 reports across institutions using the same cloud HSM.

**Status:** Partial.

**Closing language proposal.** Section in `soc-pack/` titled "Subservice organization scoping — chain HSM custody" recommending the carve-out method as the default posture (cloud HSM provider's SOC 2 / ISAE 3402 / FedRAMP attestation is a separate report) plus the complementary user-entity controls the institution must operate (CKM policy, key-rotation procedure, HSM access logs review, subservice-organization SOC report review on receipt). The inclusive method is acceptable for institutions whose engagement scope includes the cloud HSM by explicit reference; this should be rare. Consistent posture produces consistent reports.

### MQ-7. SOC for AI Systems criterion alignment

**Question.** AICPA's SOC for AI Systems (2025 draft, currently in field-test under AICPA Audit and Attest Standards Board) introduces AI-specific control criteria covering data quality, model performance, model bias, model drift, model lifecycle management, and model-related access controls. The chain is designed precisely as logging infrastructure for AI systems — it captures the AI's inputs, outputs, model identity, decision rationale, and routing decisions. The chain SHOULD be attestable under SOC for AI Systems, but the spec's control prescriptions are written in FFIEC vocabulary, not SOC for AI Systems vocabulary. Each engagement consuming the chain for SOC for AI Systems opinion writes its own mapping; the mappings vary; the attestation rigor varies.

**Status:** Partial.

**Closing language proposal.** A new informative annex (parallel to the AI Act Article 12 mapping that Henrik proposed in his EU drop) titled "SOC for AI Systems criterion mapping" listing each of the AICPA SOC for AI Systems criteria and the chain attribute(s)/event(s) that fulfill each criterion. This depends on the SOC for AI Systems final publication — which is currently in 2025 field-test — so the annex's specifics evolve; but the mapping pattern is establishable now.

### MQ-8. Control-deficiency disposition rubric

**Question.** During the SOC 2 engagement, the auditor identifies findings — exceptions to the control design or operating effectiveness. Each finding is dispositioned as one of: control performed effectively, control deficiency (inconsequential), significant deficiency, or material weakness. AICPA AT-C 205 §A52 names the disposition criteria but the criteria are generic. For chain-of-custody specifically, what dispositions apply to:

- A single MAC-mismatch on a tenant-day during the engagement period (likely a control deficiency in §10.1 reconciliation if it wasn't surfaced by the institution; possibly a significant deficiency if multiple MAC-mismatches surfaced)
- A missing seal record for a tenant-day with zero events (per §4.2 every tenant-day requires a seal; missing empty-day seal is reported as `missing seal for tenant-day {D}` per §10.12 — control-completeness failure, but not chain-integrity failure)
- A `dev_mode=true` seal in production (this is severe — §10.7 says "double-protection regulator-visible line"; one occurrence likely a significant deficiency, recurring occurrences a material weakness)
- An algorithm rotation that crossed a tenant-day boundary without §10.10 procedure compliance (depends on operational impact)
- A late-binding entry above some threshold count or duration (§4.2.2 names late-binding as normal operations, not violation; high-volume late-binding may indicate operational deficiency)

Without a disposition rubric, attestation firms apply different thresholds, producing inconsistent reports across institutions.

**Status:** Partial.

**Closing language proposal.** Section in `auditor-procedures.md` titled "Disposition rubric — chain-related findings" listing each scenario above and the recommended disposition under AT-C 205 §A52, with the institution's CC8.1 control description naming the institution's response procedure for each scenario. The rubric is informative — engagement-specific facts may warrant different disposition — but it provides a defensible baseline.

### MQ-9. Evidence sufficiency for §10.13 evidentiary-artifact retention

**Question.** §10.13 (evidentiary artifact retention list, informative under Diego's drop #03) names the artifacts the institution retains for evidentiary purposes (chain files, seal records, KMS logs, RNG attestation, etc.). The auditor's testing of retention compliance is straightforward in design (sample the retention store, verify the artifacts are present, verify they are read-only / append-only, verify the retention period matches policy) but the spec doesn't prescribe the auditor's sampling methodology, frequency, or failure-disposition. Different engagements test differently.

**Status:** Nit.

**Closing language proposal.** Section in `auditor-procedures.md` titled "Evidentiary artifact retention — auditor test procedure" naming the sampling unit (one tenant-day's full evidence package), frequency (sample of 30 tenant-days per period, biased toward recent days for currency), and failure-disposition (a missing artifact is a control deficiency in CC9.1; multiple missing artifacts a significant deficiency).

### MQ-10. Multi-spec attestation — mapping across SOC 2 + ISAE 3402 + FedRAMP + StateRAMP

**Question.** Institutions deploying the chain in multi-jurisdiction commercial / public-sector contexts may need attestation under multiple frameworks: SOC 2 (commercial customers, US), ISAE 3402 (commercial customers, EU/UK/global), FedRAMP (US federal), StateRAMP (US state government), HITRUST CSF (healthcare), PCI DSS (payment processing). Each framework has its own control objectives and evidence requirements. The chain's attestation cost grows linearly with the number of frameworks unless the spec supplies a unified control-objective mapping.

**Status:** Nit.

**Closing language proposal.** A new informative annex (parallel to the per-TSC mapping in MQ-1) titled "Multi-framework control mapping" listing each chain-prescribed control alongside the equivalent control objective in SOC 2 (US), ISAE 3000 / 3402 (international), FedRAMP (US federal), StateRAMP (US state), HITRUST CSF (healthcare), PCI DSS (payments). The mapping is informative and supports unified evidence collection — the institution produces the same evidence once and the attestation under multiple frameworks reads the evidence through different criteria lenses.

---

## Per-role roll-up

**For the spec working group.** Eight Partials and two Nits. None affect the cryptographic-integrity claim. All affect the attestation engagement's quality and consistency. The pragmatic close-out path is two new supplementary documents: `auditor-procedures.md` (parallel to the existing institution-side `audit-procedures.md`) and a Trust Service Criteria mapping annex. Both are layered on the existing normative core; neither requires changes to the spec's normative text.

**For attestation firms engaged on chain-related SOC 2 / ISAE 3000 engagements.** The spec is unusually attestation-friendly at the technical level (byte-reproducible verifier output, normative reconciliation, prescribed control patterns) but underdeveloped at the engagement-procedure level. A firm without prior chain-engagement experience writes substantial engagement procedures from scratch on first engagement, then reuses them on subsequent engagements. The spec could absorb this engagement-procedure layer once and produce consistent attestation rigor across firms.

**For institutions adopting the chain expecting SOC 2 attestation.** The institution's CC8.1 (or equivalent) control description SHOULD include the chain-specific attributes named across MQ-2 through MQ-9: RNG source declaration, verifier-run posture, software-key adapter exclusion pattern, HSM custody scoping (carve-out vs inclusive), reconciliation cadence. Without these CC8.1 attributes, the SOC 2 engagement scopes them in the field, producing engagement variance.

---

## Where I would prioritize

1. **MQ-1 (per-TSC mapping table).** Highest leverage. Every SOC 2 engagement starts with a TSC mapping; the spec providing one once produces consistent attestation rigor across firms. One annex.

2. **MQ-3 (verifier output evidence — institution-run vs auditor-run posture).** Operationally common across engagements. Affects evidence reliability and engagement cost. One section in `audit-procedures.md` plus a CC8.1 attribute the institution declares.

3. **MQ-2, MQ-4, MQ-5 (auditor-procedure templates for RNG, reconciliation sampling, software-key adapter exclusion).** Three sections in a new `auditor-procedures.md`. High-leverage because every engagement writes these from scratch today.

4. **MQ-6 (subservice-organization scoping for HSM).** Substantively affects SOC 2 report shape. One section in `soc-pack/`.

5. **MQ-7, MQ-10 (SOC for AI Systems alignment + multi-framework mapping).** Forward-scope; depends on AICPA SOC for AI Systems final publication. Inform-scope addition when final guidance lands.

6. **MQ-8, MQ-9 (disposition rubric + retention testing).** Editorial polish; lower priority.

---

## Stopping criterion

This drop carries 8 Partial + 2 Nit findings against the v1.0-final-amendment spec, all clustered around audit-engagement quality and attestation-rigor consistency. None affect the spec's cryptographic-integrity claim or the implementation contract. All affect the SOC 2 / ISAE 3000 attestation engagement's procedural rigor and the consistency of attestation across firms and institutions. The pragmatic close-out path is supplementary `auditor-procedures.md` (auditor-side) plus per-TSC and multi-framework mapping annexes; neither disrupts the spec's normative core.

The chain's design is unusually attestation-friendly at the technical layer; the gap is the engagement-procedure layer that today is written by each attestation firm independently. Closing the gap would meaningfully improve the consistency of SOC 2 and ISAE 3000 attestation across institutions deploying the chain — and would reduce the engagement cost for institutions adopting the chain, because the engagement procedures are reused rather than rewritten.
