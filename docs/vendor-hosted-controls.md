# Vendor-hosted topology controls

> **What this doc is.** Control distribution between institution and vendor in vendor-hosted topology. Closes the SOC vendor-hosted-controls partial.

## Topology recap

In vendor-hosted topology:

- The vendor operates the chain implementation (SDK, ledger, HSM, daily seal job)
- The institution holds the master HMAC key and the public key (or the vendor holds them with documented per-tenant isolation)
- The verifier remains operable by the regulator independently of the vendor

This differs from BYOC (institution operates in its own cloud) and self-hosted (institution operates in its own data center). In vendor-hosted, the vendor is the service organization for SOC purposes.

## Control distribution

### Institution-side controls (always operated by the institution)

Regardless of topology, the institution operates:

- The customer-correlation index (mapping customer to chain runs)
- The institution's own incident-response playbook for chain-detected events
- The institution's privacy-store (if privacy-by-design pattern is in use)
- The institution's regulator notifications
- The institution's audit-committee oversight
- The institution's internal-audit verification cadence

### Vendor-side controls (operated by the vendor in vendor-hosted)

In vendor-hosted, the vendor operates:

- HSM operations (PIN management, signing key, master key custody)
- Ledger operations (WAL, hot store, cold store)
- Daily seal job
- Operational-event logging
- Reconciliation cadence
- Backup and DR
- Verifier release pipeline (if vendor publishes the verifier)

### Shared / handoff controls

Some controls span the boundary:

| Control | Institution responsibility | Vendor responsibility |
|---|---|---|
| Master key custody | Choice: institution-held vs vendor-held with isolation | If vendor-held, document per-tenant isolation |
| Session-key handshake | Identity provider (institution or shared) | Custodian (vendor) |
| Public key registry | Hold the keys, share with regulator | Provide keys at issuance and rotation |
| Notification | Notify the institution's regulator | Notify the institution within agreed timeline |
| Discovery / legal disclosure | Receive subpoena, notify vendor | Produce data per institution's instruction |
| SOC report consumption | Consume vendor's SOC report; identify CUECs | Provide SOC report covering vendor-side controls |

## SOC reporting in vendor-hosted

The vendor is the service organization for the chain operations they perform. The vendor's SOC report covers:

- HSM operations (CC6.1 logical access; CC6.6 physical access; CC8.1 change management)
- Ledger operations (CC7.2 monitoring; A1.1 capacity; A1.3 backup/DR)
- Daily seal job (PI1.2 processing integrity)
- Verifier release pipeline (CC8.1 change management; supply chain)

The institution is the user entity. The institution's controls (CUECs in the user-entity sense) cover:

- Customer-correlation index
- Incident response
- Privacy-store (if applicable)
- Regulator notifications
- Internal-audit verification

The institution's SOC report (if applicable) covers institution-side controls; the vendor's SOC report is consumed as third-party-evidence for vendor-side controls.

## Vendor SOC report consumption

The institution's vendor-management procedures:

1. Receive vendor's SOC report annually
2. Review vendor's controls against the institution's expectations
3. Identify any CUECs the vendor's report assumes the institution operates
4. Confirm the institution operates those CUECs
5. Document the vendor's controls in the institution's third-party-risk-management framework

If the vendor's SOC report has qualifications or modified opinions, the institution evaluates whether the institution's risk posture is acceptable; the institution may add compensating controls or terminate the vendor relationship.

## CUEC distribution in vendor-hosted

In vendor-hosted topology, the CUECs from `docs/control-map/CUECs.md` redistribute:

| CUEC | Vendor-hosted distribution |
|---|---|
| CUEC-IAM-01..05 (RBAC, HSM separation) | Vendor operates HSM-related; institution operates institution-internal IAM |
| CUEC-CRY-01..05 (HSM, master key) | Vendor operates HSM PIN, signing key; institution holds master if institution-held variant |
| CUEC-OPS-01..05 (NTP, monitoring, backups, notifications) | Vendor operates infrastructure; institution operates monitoring of vendor |
| CUEC-VER-01..04 (verifier validation, internal audit verification) | Institution operates (regardless of topology) |
| CUEC-IR-01..04 (IR playbook) | Institution operates; vendor's IR playbook is institution-side input |
| CUEC-VND-01..04 (vendor management) | Institution operates this category fully (this is what makes it vendor-hosted) |
| CUEC-CFG-01..03 (configuration and change management) | Vendor operates configuration; institution operates change-management for institution-side changes |

## Examination considerations

For institutions in vendor-hosted topology:

- The Federal Reserve / OCC / FDIC examination of the institution covers institution-side controls
- The Fed / OCC / FDIC may also conduct vendor examinations under outsourcing-supervision authority (typically for material vendors)
- The examination uses the institution's verifier output (which the institution generates from vendor-provided ledger snapshots)
- The examiner consumes the vendor's SOC report as third-party evidence

The examination is comparable to non-vendor-hosted; the boundary just shifts which controls are vendor-side vs institution-side.

## Vendor-relationship transition

If the institution decides to leave the vendor:

1. Receive vendor's data export (ledger snapshots, public keys, master keys if institution-held variant)
2. Validate exports under the spec
3. Migrate to new vendor or to self-hosted/BYOC
4. Vendor retains legal-hold copies if applicable

The chain spec's open-standards property makes vendor-relationship transition tractable. Past data verifies under the institution's keys regardless of which implementation produced it.

## Vendor-hosted exit master-key handoff procedure

The "Vendor-relationship transition" section above describes the high-level transition. For institutions in the vendor-held master-key variant — the vendor holds the IKM under documented per-tenant isolation rather than the institution holding it — the master-key handoff is the load-bearing transition step. Without an explicit exit procedure, the institution faces vendor lock-in: the vendor controls the IKM that protects 7 years of historical chain integrity, and the institution's choices reduce to abandoning the historical chain (loses audit-history continuity) or trusting the vendor to faithfully extract and transfer the IKM (places audit-history integrity in the vendor's hands at a moment when the institution has just decided the relationship is ending).

The two conformant approaches below address the lock-in. Both produce auditable evidence of the handoff and both leave the institution able to verify historical chain entries after exit. The institution chooses one approach per the vendor's HSM capabilities, the institution's own HSM capabilities, and the institution's tolerance for a chain-discontinuity disclosure.

### Approach A — HSM-to-HSM key wrap

The vendor's HSM exports the IKM as a wrapped key under the institution's HSM public-key wrapping key. The institution's HSM imports the wrapped key and unwraps it inside the institution's HSM. The IKM never appears in plaintext outside an HSM boundary. After successful import and verification, the vendor destroys the source IKM in vendor-side custody.

**Sequence.**

1. The institution publishes its HSM wrapping public key to the vendor under documented authorization (the institution's CISO and the vendor's contractually-authorized custodian sign the wrapping-key delivery record).
2. The vendor's HSM performs the key-wrap operation per the vendor's HSM provider's wrap mechanism (PKCS#11 `CKM_RSA_AES_KEY_WRAP` or equivalent). The wrap operation is authorized inside the vendor's HSM by the same dual-control posture the vendor uses for IKM rotation.
3. The vendor delivers the wrapped IKM artifact to the institution under the vendor's standard sensitive-data delivery channel. The institution receives the artifact and validates the chain-of-custody record.
4. The institution's HSM imports the wrapped artifact and unwraps it. The institution's HSM verifies the unwrap operation succeeded (the resulting IKM produces the expected fingerprint under the chain's `key_fingerprint` derivation).
5. **Fingerprint match verification.** The institution computes the fingerprint of the freshly-imported IKM and compares it byte-for-byte against the fingerprint the vendor previously published in the chain's seal records. A mismatch is a failed handoff and is investigated as a potential vendor-side substitution before any historical chain operation proceeds.
6. The vendor performs IKM destruction in vendor-side custody under the vendor's HSM provider's destroy mechanism. The vendor captures `master_key.exit_handoff_destroyed` (institution-defined operational event) with the institution's identifier, the affected tenant identifiers, and the destruction evidence (HSM destroy-operation log reference).
7. The institution captures `master_key.exit_handoff_imported` with the vendor's identifier, the affected tenant identifiers, the wrap-mechanism reference, and the fingerprint-match-verification result. The two events bracket the handoff for the institution's CC8.1 control evidence.

**Conformance posture.** Approach A preserves the chain's continuity. Historical chain entries verify under the same IKM after handoff because the IKM is byte-identical (the wrap-and-unwrap is mathematically lossless). The institution's chain spans both the vendor's custody window and the post-exit institution-held window without a chain-discontinuity event.

**SOC and FFIEC implications.** The vendor's SOC report covers the export-and-destruction operations under CC6.1 logical access and CC8.1 change management. The institution's SOC report (or the institution's CC8.1 control description if the institution does not produce a SOC report) covers the import operation and the fingerprint-match verification. The FFIEC examiner reviews both sides during the institution's next examination cycle. The institution's control description names Approach A as the chosen exit procedure and references the vendor's authorization-record format.

**Worked example.**

> Acme Bank exits a vendor-hosted relationship after seven years. The vendor's HSM holds the IKM for tenants `acme_consumer_lending` and `acme_servicing_collections`. Acme Bank's CISO publishes Acme's HSM wrapping public key to the vendor's contractually-authorized custodian under the wrapping-key-delivery record dated 2026-08-01. On 2026-08-15, the vendor's HSM performs `CKM_RSA_AES_KEY_WRAP` on the two tenant IKMs and delivers the two wrapped artifacts to Acme via the vendor's secure-data delivery channel. Acme's HSM imports and unwraps both artifacts on 2026-08-16. The IKM fingerprints match the seal-record fingerprints from the vendor-custody period (2019-2026) byte-for-byte. The vendor destroys both source IKMs on 2026-08-17 and captures `master_key.exit_handoff_destroyed` for each tenant with the destruction-log reference. Acme captures `master_key.exit_handoff_imported` for each tenant. Acme's verifier runs against the historical chain under the imported IKM and produces PASS for the seven-year period. Acme's chain continuity is preserved.

### Approach B — Fork-and-rotate with explicit chain-discontinuity disclosure

The institution stands up a new IKM in institution-controlled custody and rotates the chain to the new IKM at the exit boundary. The historical chain remains under the vendor's custody with a documented continuation-of-evidence agreement; the post-exit chain operates under the institution's new IKM. The institution captures the chain-discontinuity event and discloses it to the regulator as part of the exit notification.

**Sequence.**

1. The institution generates a new IKM in institution-controlled custody per the standard tenant-onboarding sequence (per `operator-guide.md` "Tenant onboarding" §1) — HSM-resident, 32-byte minimum, non-extractable in Model B or bounded-extractable in Model A.
2. The institution registers the new public-key fingerprint with the regulator per the same regulator-fingerprint reception procedure used during onboarding. The regulator's acknowledgement reference is captured.
3. At a documented rotation boundary (typically a UTC midnight that aligns with a seal boundary so the chain-discontinuity event aligns cleanly with a seal record), the institution rotates the chain to the new IKM. Events captured before the boundary remain under the vendor's IKM; events after the boundary use the institution's IKM.
4. The institution captures `chain.discontinuity_recorded` (institution-defined operational event) with the affected tenant identifiers, the discontinuity timestamp, the previous IKM fingerprint, the new IKM fingerprint, and the disclosure reference.
5. The institution and the vendor execute a continuation-of-evidence agreement under which the vendor retains the historical chain and the historical IKM under the vendor's standard 7-year retention. The institution retains read access to the historical chain for examination purposes per the agreement. The agreement names the IKM-handoff fallback if the vendor terminates retention or becomes unable to produce the historical chain (typical fallback is Approach A executed during the agreement's wind-down period).
6. The institution discloses the chain-discontinuity to the regulator in the next regulatory communication (typically the institution's next examination's evidence package or a standalone notification under the institution's regulatory-relations cadence).

**Conformance posture.** Approach B does NOT preserve chain continuity. The verifier sees two distinct chains — the historical chain under the vendor's IKM and the post-exit chain under the institution's IKM — separated by an explicit discontinuity event. Both chains are individually verifiable; the institution's audit-history claim is two-segment rather than single-segment. The discontinuity is a recorded fact rather than a verification failure; the verifier's disposition for events spanning the boundary is "verified under their respective IKMs with documented discontinuity at `[timestamp]` per `chain.discontinuity_recorded` event reference `[id]`".

**SOC and FFIEC implications.** Approach B requires explicit documentation in the institution's SOC report and the institution's CC8.1 control description naming the discontinuity, the disclosure to the regulator, and the continuation-of-evidence agreement with the former vendor. The FFIEC examiner treats the discontinuity as a known, documented event rather than a finding; the institution's evidence package includes the regulator-disclosure record and the vendor-side continuation-of-evidence agreement. The examiner MAY require additional substantive testing of the chain-spanning-discontinuity edge cases (e.g., events captured close to the boundary timestamp) to confirm the boundary semantic is operationally clean.

**Worked example.**

> Beta Bank exits a vendor-hosted relationship after three years and chooses Approach B because the vendor's HSM provider does not support the wrap mechanism Beta's HSM provider expects (Approach A is not feasible). Beta generates new IKMs in Beta's HSM on 2026-08-01 and registers the new fingerprints with the OCC. The OCC acknowledges receipt on 2026-08-08. Beta and the vendor agree on a UTC-midnight rotation boundary of 2026-09-01T00:00:00Z. On 2026-09-01, Beta rotates the chain to the new IKMs and captures `chain.discontinuity_recorded` for each tenant. The vendor retains the 2023-2026 historical chain under continuation-of-evidence agreement until 2033 (matching Beta's standard 7-year retention). Beta discloses the discontinuity to the OCC in the September 2026 examination communication. The OCC accepts the discontinuity as documented and Beta's chain continues operation under the new IKMs from September 2026 forward. The vendor delivers historical-chain evidence on demand for examinations covering the 2023-2026 period until 2033.

### CC8.1 and CUEC additions

For either approach, the institution's CC8.1 control description names the chosen exit procedure, the operational events captured, and the regulator-disclosure timing. The institution's CUEC inventory adds a per-vendor-relationship CUEC referencing the chosen approach:

| CUEC | Description |
|---|---|
| CUEC-VND-EXIT-A (Approach A) | Institution maintains an HSM wrapping-key publication procedure, validates fingerprint matches on import, and retains import/destroy operational events for the vendor-relationship's full retention period |
| CUEC-VND-EXIT-B (Approach B) | Institution captures and retains the chain-discontinuity event, retains the regulator-disclosure record, and operates the continuation-of-evidence agreement with the former vendor for the historical-chain retention period |

The CUEC text is institution-adaptable; the structural content is the auditable handoff evidence.

### Failure mode — the vendor refuses to cooperate

If the vendor refuses to cooperate with either approach — refuses to wrap-and-deliver the IKM under Approach A, refuses to enter a continuation-of-evidence agreement under Approach B, or refuses to acknowledge the institution's exit notice — the institution's options are legal-and-regulatory rather than technical.

The escalation path:

1. **Contractual escalation.** The institution invokes the contract's exit-cooperation clause. If the contract lacks such a clause, the institution invokes the contract's general dispute-resolution clause and prepares for a contractual remedy.
2. **Regulatory notification.** The institution notifies its primary regulator that the vendor relationship is in non-cooperative termination. The notification names the vendor, the affected chain operations, and the institution's planned interim posture (typically Approach B with the regulator's understanding that the vendor is not cooperating with the continuation-of-evidence agreement).
3. **IR scenario invocation.** The institution treats the non-cooperation as a vendor-induced operational risk per IR Scenario 12 (vendor-supply-chain-driven incident) or the institution's vendor-failure scenario. The IR playbook documents the non-cooperation, the chain-discontinuity decision, and the regulator-notification timing.
4. **Legal escalation.** The institution's legal function evaluates breach-of-contract and ancillary remedies. If the vendor is regulated, the institution MAY refer the non-cooperation to the vendor's primary regulator under outsourcing-supervision authority. The institution's legal function leads the escalation.
5. **Forced fork-and-rotate without continuation-of-evidence agreement.** As a last resort, the institution executes Approach B without the vendor's signed agreement. The historical chain becomes inaccessible to the institution if the vendor terminates retention; the institution's audit-history claim for the vendor-custody period reduces to "the vendor-held chain was operated under the spec for `[period]` and the chain-of-custody is documented in vendor-side records the institution does not control after `[non-cooperation date]`". This is the operationally weakest exit posture; the institution's risk function and legal function jointly approve before invocation.

The institution's contract with the vendor SHOULD include an exit-cooperation clause naming Approach A or Approach B as the contracted exit procedure. Institutions adopting vendor-hosted topology MUST review the vendor contract's exit terms before signature; institutions whose contracts predate this guidance SHOULD seek contract amendment at the next renewal cycle.

## Evaluating vendors

Institutions evaluating chain-of-custody vendors look at:

- **Conformance.** Does the vendor pass the conformance corpus?
- **SOC report.** Does the vendor provide a SOC 2 Type II report on chain-relevant controls?
- **HSM operation.** Is the HSM provider's FIPS validation acceptable to the institution's regulator?
- **Master-key custody options.** Does the vendor support institution-held master with HSM-mediated handshake?
- **Verifier accessibility.** Does the vendor support institution running the verifier independently?
- **Pricing model.** Is the pricing aligned with the institution's volume and budget?
- **Deployment model.** Vendor-hosted vs BYOC vs self-hosted as options?
- **Vendor transition support.** Does the vendor commit to data-export and key-handoff procedures?

The vendor's conformance is the ground floor; the rest is institution-specific.

## Pricing considerations

Vendor-hosted pricing typically combines:

- Per-tenant base fee
- Per-event volume tier
- HSM operating costs (passed through or bundled)
- Optional add-ons (multi-region resilience, hourly cadence, additional support)

Vendor-hosted is generally cheaper than self-hosted for low-volume institutions because HSM costs are amortized across multiple tenants. Self-hosted becomes more economical at higher volumes.

The cost model (`docs/cost-model.md`) addresses self-hosted; vendor-hosted institutions get pricing from their vendor and use the vendor's pricing model in budget planning.

## Common vendor-hosted issues

- **Vendor changes its implementation without institution notice.** Defense: institution requires the vendor to declare conformance for each new version; the institution validates against the corpus.
- **Vendor's SOC report has qualifications affecting chain controls.** Defense: institution reviews the qualifications and decides whether to add compensating controls.
- **Vendor's pricing increases.** Defense: contract terms; institution's choice to renew or migrate.
- **Vendor's outage affects chain operations.** Defense: vendor's SLAs; institution's IR playbook composes; institution may require multi-region or multi-vendor for resilience.
- **Vendor's master-key custody differs from institution's expectation.** Defense: contract terms; institution's choice to use institution-held master.

## Audit-committee oversight in vendor-hosted

The institution's audit committee oversees the vendor relationship. Key questions:

- Has the vendor's SOC report been received and reviewed? Any concerns?
- Are the institution's CUECs operating as expected?
- Has any chain-detected event involved vendor coordination? How was it handled?
- Are vendor pricing and contract terms current?
- Is there a vendor-transition plan if needed?

The vendor relationship is one of the institution's third-party relationships; chain-related concerns flow through standard third-party-risk-management.

---

## TPRM lifecycle alignment with OCC Bulletin 2013-29 and the June 2023 Interagency Guidance

OCC Bulletin 2013-29 and the June 6, 2023 Interagency Guidance on Third-Party Relationships establish the documented life-cycle the institution applies to every third-party arrangement: planning, due diligence, contracting, ongoing monitoring, and termination. The chain spec is technically sound; the TPRM program is what wraps the chain vendor in that lifecycle. This section extends the existing vendor-hosted controls with the TPRM-specific surfaces a $10B+ regional bank's risk team needs.

### Vendor risk tier classification

For a vendor operating any part of the chain-of-custody infrastructure (hosting the ledger, running the seal job, holding the HSM, providing the SDK), the minimum risk tier is **Critical**. The chain is regulatory evidence (the FFIEC examiner relies on it during examination), forensic evidence (litigation and incident investigation depend on it), and an operational control gating the institution's MRM and AI-decision logging posture. Failure of the chain vendor blocks the institution's ability to demonstrate AI-decision integrity to the examiner.

Critical tier implies: 99.9% availability SLA (no more than 8.7 hours of downtime per year), 1-hour incident-notification SLA, 4-hour RTO for service disruption, annual SOC 2 Type II attestation with chain-integrity controls scoped, board-level escalation for any incident exceeding 4 hours of downtime, quarterly risk review at the Chief Risk Officer level, vendor financial-stability assessment requiring the vendor to be profitable with at least 2 years of capitalization. Startups or poorly capitalized vendors are unacceptable for Critical tier without an institution-side compensating control (typically source-code escrow plus a documented in-house operation plan).

Vendors providing components (SDK, routing logic, model-decision functions) that integrate with the chain but do not host the chain itself may be classified as **High** (non-substitutable component, failure affects chain data quality) or **Medium** (commodity component, substitutable). The institution documents the risk tier in the vendor risk register and adjusts annually based on performance and financial stability.

### Vendor onboarding scoring

The institution's onboarding scoring covers six dimensions, each weighted by the institution's risk team:

| Dimension | What it scores | Evidence the institution consumes |
|---|---|---|
| Spec conformance | Does the vendor's implementation pass the FFIEC corpus? | Vendor-conformance attestation per `docs/vendor-conformance-attestation.md` |
| SOC 2 / ISAE 3000 quality | Does the vendor's SOC scope cover chain-integrity controls? | SOC 2 Type II report with the custom chain-integrity scope below |
| Financial stability | Is the vendor profitable, diversified, and capitalized? | Vendor financials, customer-diversification evidence, funding-round history |
| Incident-response maturity | Can the vendor meet the 1-hour escalation SLA? | IR playbook walkthrough, prior-incident references |
| Data-portability and exit | Can the institution exit cleanly with a forensically defensible export? | Documented export procedure, sample export, escrow agreement (if applicable) |
| AI governance (for AI vendors) | Model lineage, training-data governance, fairness audits, RAI program | AI-vendor DDQ responses (see below) |

Each dimension is scored 1-5; the institution's risk team sets a minimum aggregate score for tier admission. Vendors below the minimum are rejected, conditionally accepted (with named compensating controls), or escalated to the Chief Risk Officer for case-by-case decision.

### AI-vendor DDQ extensions

The institution's standard DDQ (covering security, access control, encryption, incident response) is a baseline. AI vendors introduce risks the standard DDQ does not address. The DDQ extension covers:

**Model governance.** What is the source of the model the vendor uses (vendor-trained, OpenAI, Anthropic, Cohere, open-source)? Is the model fine-tuned on customer data? If so, is the bank's data included? What is the data-retention and deletion policy for fine-tuning data? When is the underlying model updated, and does the vendor notify customers in advance? Are model changes tested for regression on the bank's use cases before deployment?

**Training-data governance.** What data was used to train the base model? Is the bank's data or competitors' data included? For fine-tuned models, what is the training dataset composition? Does the vendor have a data-retention agreement with the upstream model provider (does OpenAI retain and use the bank's prompts for model improvement, for example)?

**Fairness and bias.** Has the vendor conducted fairness audits? Documented evidence of audits for disparate impact on protected classes? Does the vendor provide model-output explainability? Does the vendor have a process for detecting and remediating model drift?

**Responsible AI program.** Does the vendor have a Responsible AI team or function? A model-review or governance committee approving models for use in regulated industries? Output-anomaly monitoring?

**Subprocessor dependencies.** Which LLM providers does the vendor depend on? For each subprocessor, what is the data-handling agreement? Does the subprocessor use the bank's data to improve its model? If the vendor switches LLM providers, how does that affect the chain's continuity?

**AI-related incident response.** Escalation procedure if a model produces biased or harmful outputs. Notification timeline for model-related incidents (a model that discriminates against a protected class, for example). Remediation obligations if model changes cause a security or compliance issue.

The DDQ responses become input to the onboarding scoring. The Shared Assessments SIG Lite/Core questionnaire framework can be extended to carry these questions; the institution's risk team and Shared Assessments coordinate on mainstream framework alignment.

### MSA audit-rights template

The institution's MSA includes audit-rights language extending the existing vendor-hosted controls. Recommended text:

> If the institution uses a vendor-hosted chain implementation, the vendor MUST grant the institution audit rights to retrieve, at minimum: (1) daily seal records (the seal file per spec §4.2), (2) the institution's chain entries in the vendor's ledger, exportable in the wire format specified in spec §5 (OTLP or JCS-canonical JSON), (3) read-only access to the verifier output log produced by the vendor's daily verification procedure (spec §7 steps 1-9 run by the vendor at ingest), (4) access to the vendor's incident-response logs for any chain-detected anomalies. The vendor must provide this access to the institution on demand within 5 business days and on a standing schedule (at least quarterly). The vendor must NOT commingle the institution's data with other tenants' data in the export; if technical constraints require multi-tenant exports, the vendor must redact other tenants' entries or provide cryptographically isolated subsets.

This closes the gap where a vendor claims "proprietary ledger, read-only" and the bank loses audit visibility. The recommended language is templatized; institutions modify wording to fit their MSA structure but preserve the substantive obligations.

### SOC 2 Type II scope for chain-integrity controls

The vendor's annual SOC 2 Type II report must include a custom Trust Service Criteria section covering chain-of-custody controls. The institution's contract requires the vendor to scope the SOC engagement explicitly. The expected scope:

**Control objective.** The vendor produces chain entries, computes per-entry HMACs, aggregates entries into daily ledgers, computes daily Merkle roots, and produces daily HSM-signed seal records, all per FFIEC chain-of-custody v1.0 specification §§1-4. The vendor's controls are designed to prevent or detect tampering with chain entries, unauthorized modification of ledger aggregation, or deviation from the deterministic HMAC-SHA-256 / Merkle / HSM-signature process.

**Control activities tested.** (a) Daily verifier execution: the vendor runs spec §7 steps 1-9 on a sample of recent entries and confirms zero defects; per-entry MAC recompute using the institution's IKM confirms the persisted `payload_hash` matches the recomputed MAC. (b) Merkle root reproducibility: the vendor recomputes the daily Merkle root independently using a separate codebase and confirms it matches the seal record's root. (c) Signature reproducibility: the vendor verifies the daily seal record's signature using the public key corresponding to the HSM's private key. (d) IKM non-extractability: the vendor confirms via HSM audit logs that the institution's IKM was never exported. (e) Incident-response logs: the vendor provides a sample of any anomalies detected during the period and confirms they were escalated per the spec's 36-hour notification clock.

**Test results.** Pass/fail per control activity. If any test fails, the vendor provides corrective action within 30 days. Persistent failures require IR-Scenario coordination with the institution.

**Scope limitations.** The SOC auditor does NOT test the upstream LLM vendor's correctness; that is outside chain-of-custody scope. The SOC auditor tests only the vendor's chain-construction and custody procedures. The upstream LLM vendor's own SOC report covers the LLM's controls separately.

The SOC engagement scope language is templatized; institutions hand the language to the vendor and the vendor's SOC firm at engagement-planning time.

### Fourth-party (subprocessor) governance

The chain introduces nested third-party dependencies: bank → AI vendor → cloud hosting (AWS, Azure, GCP) → upstream LLM providers (OpenAI, Anthropic, Cohere) → infrastructure components. The 2023 Interagency Guidance treats subprocessor relationships as in-scope: the institution's TPRM program governs the upstream chain.

**Dependency mapping.** The institution's vendor-management procedure maps the chain's upstream dependencies at contract time: which vendors the AI vendor depends on for ledger hosting, which vendors for LLM services, which vendors the hosting vendor depends on for infrastructure. The map updates whenever the AI vendor's subprocessor list changes (the AI vendor's contractual obligation to notify under the change-of-subprocessor clause).

**Contractual flow-down.** The AI vendor's MSA includes flow-down language requiring the AI vendor to extend chain-of-custody obligations to the hosting vendor. Specifically: (a) the hosting vendor commits to non-extraction of the institution's IKM (the cryptographic key remains in the hosting provider's HSM, never exported to the AI vendor's process); (b) the hosting vendor's SOC scope covers encryption, access control, and incident response for the institution's chain data; (c) if the hosting vendor experiences a breach or security incident affecting the institution's chain data, the hosting vendor notifies the AI vendor within 2 hours and the AI vendor notifies the institution within 4 hours (a 6-hour total escalation window).

**LLM upstream dependencies.** If the AI vendor depends on an upstream LLM provider, the AI vendor documents this in the contract and in chain routing-event entries (per spec §4.4.1 routing-decision capture). The institution's examiner has visibility into which LLM providers are in the call path. The AI vendor does not need to flow chain-of-custody requirements to the LLM provider (the LLM provider does not custody the chain); but the LLM provider's incident-notification flows to the AI vendor within 2 hours and from the AI vendor to the institution per the standard escalation.

**Concentration risk.** The institution's vendor-management procedure monitors concentration: if more than 50% of the institution's chain-of-custody traffic is hosted by a single cloud provider, a regional outage affects the majority of the institution's ledgers. The institution either (a) accepts the risk with board approval, (b) diversifies across multiple cloud providers, or (c) requires the vendor to operate multi-region failover. Industry-wide concentration (one vendor exceeding 30% market share) escalates to the institution's primary regulator if institutional exposure is high.

### Vendor incident-notification SLA

The spec's incident-response playbook mandates a 36-hour notification clock to the examiner. The vendor's contractual obligation to notify the institution must close within that window. The MSA mandates the vendor's escalation timeline:

**Detection-to-notification (1 hour).** The vendor notifies the institution within 1 hour of detection by phone or email to the institution's 24/7 security operations center. The notification includes: incident summary; affected chain(s) (tenant ID, run IDs, date range); preliminary root-cause hypothesis; interim remediation steps taken; estimated timeline to full remediation.

**Written notice (4 hours).** Within 4 hours of detection, the vendor provides a written incident report (email or portal) covering: incident description and technical details; forensic findings (if available); impact assessment (entries affected, days of sealing delayed); remediation plan and timeline; commitment to provide detailed root-cause analysis within 24 hours.

**Root-cause analysis (24 hours).** Within 24 hours of initial detection, the vendor provides a complete root-cause analysis: what went wrong, why, when it started, how long undetected, preventive measures.

**Forensic handoff (48 hours).** If the incident involves suspected security breach or tampering, the vendor provides all forensic evidence (logs, snapshots, audit trails) within 48 hours to the institution's IR team. The institution's team may run parallel investigations.

**Triggering events.** Any of: chain-entry verification failure (spec §7 step 1-9 failing), seal-job delay exceeding 60 minutes past the UTC day boundary, IKM custody anomaly (attempted extraction, unexplained access, HSM audit-log gap), cryptographic-material theft or unauthorized-access suspicion, breach affecting ledger integrity or confidentiality, scheduled maintenance that would delay seal-job execution beyond 60 minutes.

The 1-hour-to-institution notification is a contractual requirement, not an optional best-effort. The institution's 36-hour clock starts when the institution learns of the incident; the vendor's escalation must occur within the institution's 36-hour window to allow the institution to meet the regulatory deadline.

### Ongoing monitoring KPIs and KRIs

The institution monitors the vendor's chain operations through KPIs (operational effectiveness) and KRIs (risk indicators).

**KPIs (operational effectiveness).** Seal-publication SLA achievement (target >99.9%, alert <99%). Daily verifier pass rate (target 100%, alert <99.9%). Incident-response time (target <1 hour, alert >4 hours). SOC report currency (target: current within 12 months). System availability (target 99.9%, alert <99%).

**KRIs (risk indicators).** Fingerprint-mismatch rate (target 0%, alert >0.1%). IKM custody anomalies (target 0, alert ≥1 per quarter). Seal-age distribution (target all seals <65 minutes, alert any seal >80 minutes). Entry-verification failure categories (track by category — MAC mismatch, seq error, prev_hash error, fingerprint error). Vendor financial health (track quarterly). Vendor concentration risk (institutional <50%, industry-wide <30%).

**Monitoring cadence.** Daily: vendor publishes brief operational status. Weekly: institution's operations team reviews seal-age and verification failures. Monthly: institution's risk team reviews KPIs/KRIs and meets with vendor's ops team. Quarterly: Chief Risk Officer reviews vendor performance and continuation decision; escalate to board if any KRI threshold is breached. Annually: re-assess contract terms, SLA performance, and risk tier.

The institution's GRC platform (RSA Archer, ServiceNow GRC, OneTrust) carries the dashboard. Inputs: vendor's published operational status, institution's independent verification (institution's own verifier run), vendor's SOC report and audit findings.

### Exit, data portability, and forensic defensibility

When terminating the vendor or non-renewing the contract, the institution must receive the chain data in forensically defensible format within a defined window. The MSA includes:

**Exit window.** 30 days from contract termination for full data export. Expedited (14 days) is available for a reasonable fee not exceeding the vendor's labor cost.

**Export contents.** All chain entries (full ledger dump) in spec §5 wire format (OTLP or JCS-canonical JSON); all daily seal records (Merkle root + signature); the vendor's verifier output log showing all entries verified per spec §7 steps 1-9; the institution's IKM rotation history (key versions, fingerprints, rotation dates, no raw IKM bytes).

**Data integrity.** The vendor cryptographically signs the data-export manifest (using the vendor's release-signing key for software exports; HSM-backed for HSM-custody exports). The institution validates the signature and confirms authenticity.

**Forensic defensibility.** The exported data is suitable for forensic use (institution's law firm or examiner can ingest it into an independent verifier and produce a verification report). The vendor must NOT export in a proprietary format that only the vendor's custom tools can read. If proprietary tooling is required, the contract includes either (a) perpetual license to the vendor's export tool, (b) source-code escrow for the export tool with documented release triggers, or (c) the vendor publishes the proprietary-format specification so the institution can independently rebuild the export tool.

**Retained audit rights.** During the 30-day exit window (or longer if litigation/examination depends on the chain), the institution retains read-only access to the ledger to confirm the export is complete. The vendor must NOT delete the ledger until the institution explicitly confirms the export is complete and verified.

**Retention.** The institution retains the exported chain data for the full 7-year retention period per spec §10.9, regardless of vendor relationship status. Multi-decade retention applies for healthcare deployments per `docs/regulator-pack/healthcare-overlay.md` §H4.

### Termination and offboarding procedure

The full offboarding procedure has six phases:

1. **Notice and planning (6 months before contract end).** Institution notifies vendor of non-renewal or early termination. Risk team and IT team meet with vendor to plan offboarding. Vendor confirms data-export timeline (30 days from contract end). Institution identifies target environment for the exported data.
2. **Knowledge transfer (3 months before contract end).** Vendor documents the chain's configuration and operational procedures (HSM setup, seal-job configuration, key-rotation schedule). Vendor's SRE team trains institution's IT team on chain operations post-vendor.
3. **Data export preparation (2 months before contract end).** Vendor prepares data-export specification. Institution validates the format against spec §5. Vendor produces a test export on a historical subset for validation.
4. **Final data export (within 30 days post-contract).** Vendor freezes the ledger to read-only. Vendor produces full export in spec §5 format. Institution validates: all entries present (no gaps); spec §7 verifier runs on a sample of exported entries with PASS results; export matches the data-count the vendor reported during monitoring.
5. **Post-export validation (within 30 days post-export).** Institution's IT team runs a full verifier pass. If verification fails, the institution escalates to the vendor (residual obligation to fix export issues within 60 days of contract end). Audit/compliance team confirms export completeness.
6. **Archival and relationship closeout.** Institution archives the exported data with contract-end and export-completion dates labeled. Institution retains the data for the full retention horizon. Institution confirms all invoices paid; vendor confirms institution's data deleted from vendor-controlled systems. Legal team archives the MSA.

### Change-of-control and merger/acquisition provisions

The MSA includes change-of-control language:

**Notification.** Vendor notifies the institution within 10 business days of any change-of-control event (acquisition, merger, bankruptcy, significant investor change, sale of critical assets).

**Assumption of obligations.** The acquiring entity assumes 100% of the vendor's obligations: all SLAs, data-protection and confidentiality, audit and attestation requirements, incident-response and breach-notification. The acquiring entity has no unilateral right to renegotiate.

**Audit rights.** The institution audits the acquiring entity's chain operations within 30 days of the announcement. Audit covers: acquiring entity's SOC report or equivalent; verification that the chain ledger was not tampered with during the transition; confirmation that IKM custody was not disrupted.

**Data integrity verification.** Within 60 days, the institution runs spec §7 verifier on a sample of historic chain entries (entries from before the change-of-control) to confirm the chain ledger remains unmodified. Verification failure invokes IR Scenario 4 (suspected compromise).

**Termination right.** If the acquiring entity is a competitor or conflicting entity, the institution may terminate without penalty and receive the standard 30-day data-export window.

**Subprocessor changes.** Acquiring-entity-driven changes to chain hosting provider, data-center location, or HSM vendor are subject to institutional audit and approval. The acquiring entity has no right to migrate the institution's chain data without explicit written approval from the institution's Chief Risk Officer.

### Vendor business continuity and financial stability

**Financial stability assessment.** At vendor selection and annually in ongoing-monitoring review, the institution assesses: vendor profitability; customer diversification (the vendor is not dependent on one bank); capitalization (sufficient for at least 2 years of operations); public stress signals (layoffs, funding rounds indicating distress).

**Continuity-of-service commitment.** The MSA requires: in the event the vendor ceases operations or becomes unable to provide the service, the vendor maintains the ledger in a read-only state for at least 90 days while the institution arranges migration; the vendor provides daily backup exports of the institution's chain data to cloud storage under the institution's control; the vendor does not delete the institution's data until the institution explicitly authorizes deletion in writing or 1 year has elapsed since the contract end date, whichever is earlier.

**Source-code escrow.** For Critical-tier vendors with proprietary code, the institution requires source-code escrow with release triggers: vendor bankruptcy, vendor breach of MSA without cure within 30 days, vendor failure-to-cure within 30 days of a material service-level breach. The escrow agent releases the vendor's code to the institution under the documented triggers.

**Cold-start recovery.** The vendor documents a cold-start procedure: if the ledger is unavailable, the institution can rebuild from (a) backup exports, (b) SDK local SQLite files per spec §6, (c) exported OTLP messages. The institution validates the procedure with the vendor's SRE team before going live.

### Insurance and indemnification

The MSA requires the vendor to maintain:

**Cyber liability / data breach insurance.** $5M-$25M coverage depending on vendor size and data sensitivity. Covers breach notification costs, forensic investigation, credit monitoring (where applicable), business interruption, regulatory fines and penalties (where applicable), legal liability.

**Errors and omissions insurance.** $2M-$10M coverage. Covers professional liability for MSA breaches, negligence in implementing chain controls, SLA failures.

**Directors and officers insurance.** Standard good-governance indicator for corporate vendors.

**Insurance certificates.** Annually and upon request. Institution named as additional insured on cyber-liability and E&O policies. Policies are primary (pay before institution's own insurance pays). 30-day cancellation notice to the institution.

**Indemnification.** Vendor indemnifies the institution for losses from MSA breach, vendor negligence or willful misconduct, vendor failure to meet audit or incident-notification obligations, vendor's subprocessor breaches. Indemnification capped at the greater of (a) vendor's annual fees or (b) vendor's insurance limits. Standard exclusions apply (institution's misconfiguration, force majeure with narrow definition, institution-side modifications introducing vulnerabilities).

### Data-residency and geopolitical posture

The MSA names the data-residency commitment: the institution's chain data resides only in [specified regions, e.g., AWS US regions only: us-east-1, us-west-2, us-gov-cloud]. Prohibition on data transfer to other regions without the institution's written consent. Mechanism for the institution to audit the vendor's residency (institution requests logs showing which data-center region holds the institution's ledger).

For institutions operating under data-residency requirements (US-only, EU-only, country-of-origin restrictions), the MSA includes the residency commitment as a non-negotiable. If the vendor cannot commit to the institution's residency requirement, the institution rejects the vendor.

### Contracting timeline expectations

A representative timeline for vendor onboarding under the templates and procedures above:

- Phase 1 — Vendor selection and RFP: 3-6 weeks.
- Phase 2 — Due diligence (DDQ + AI-vendor extension): 4-8 weeks.
- Phase 3 — Contract negotiation: 4-8 weeks if the vendor accepts the templates; 12+ weeks if the vendor objects to material terms.
- Phase 4 — Executive sign-off: 1-2 weeks.
- Phase 5 — Implementation and onboarding: 2-4 weeks (SDK integration, IKM provisioning, seal-job configuration, initial verifier run).

Total: typically 4-6 months under standardized terms; 9+ months when material terms are disputed. Institutions using the spec-recommended templates achieve the shorter cycle. Institutions writing custom audit-rights, SOC-scope, and exit-procedure language from scratch face the longer cycle.

### Documentation checklist for FFIEC examinations

The institution's vendor-management file (assembled during onboarding, updated through ongoing monitoring) carries:

1. **Vendor selection.** RFP and vendor-response summary; vendor-evaluation matrix; selection decision and approval.
2. **Vendor due diligence.** Completed DDQ (with AI-vendor extension if applicable); vendor risk-tier assessment; vendor financial-stability assessment; vendor insurance certificates.
3. **MSA and key contracts.** Signed MSA with audit-rights, SOC-scope, incident-notification, exit, and termination clauses highlighted; DPA if applicable; any amendments.
4. **Vendor operating controls.** Vendor's SOC 2 Type II report (current, with chain-integrity controls scoped); vendor's IR procedure and escalation contacts; vendor's BCP/DR plan; vendor's change-management procedure.
5. **Ongoing monitoring evidence.** Quarterly or monthly vendor performance metrics; incident logs; vendor-management review meeting notes; vendor risk register entry showing current risk tier and any escalated findings.
6. **Third-party audit evidence.** External audit findings on the chain control; vendor's SOC audit findings; any FFIEC examination findings or supervisory letters mentioning the chain.

The file is retained for the full 7-year chain-retention period (longer if institution's document-retention policy or healthcare overlays require). The examiner requests any of these documents during examination; the institution produces within 5 business days.
