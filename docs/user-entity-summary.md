# User-entity summary

> **What this doc is.** Summary for the user entity reading a SOC 1 or SOC 2 report on a chain-of-custody implementation. The user entity may be the bank's external auditor, a downstream business partner, or any party that relies on the institution's chain controls.

## Audience

The user entity reading a SOC report on a chain-of-custody implementation typically falls into:

- **Bank's external financial-statement auditor** — relies on SOC 1 to confirm chain-related controls supporting financials
- **Downstream business partner** — depends on the bank's chain to confirm AI-driven decisions are integrity-bearing
- **Regulator** — consumes SOC 2 reports as part of supervisory assessment
- **Audit firm planning audit** — uses the report to plan reliance on service-organization controls

This doc orients the user entity quickly, without requiring full knowledge of the chain spec.

## What the chain-of-custody control does

The chain captures every AI agent decision the institution makes and produces an integrity-bearing daily seal that an independent verifier can validate. Three properties:

- **Authentic capture.** Captured events were produced by the institution's legitimate AI agent process at the claimed time.
- **Tamper evidence.** Modification of captured events is detectable.
- **Independent verifiability.** A regulator (or the user entity) can verify the integrity property without trusting the institution.

The control is documented in the FFIEC AI Chain-of-Custody Specification v1.0.

## What the SOC report attests to

The institution's operational evidence stream includes events for low-frequency, high-impact scenarios. Two specifically worth surfacing for a downstream user-entity audience:

- **`audit_file.truncation_detected`** — emitted when a writer-side process crashes mid-append and leaves the audit file with an incomplete tail. The verifier refuses to verify the truncated file (per spec §4.1) rather than silently passing a chain that lost its last entry. The institution's response procedure (IR Scenario 9) recovers the dropped events from the SDK-local buffer or upstream OTLP retention; the recovery outcome (`complete`, `partial`, `unrecoverable`) is recorded on the operational event.
- **`master_key.retired`** — emitted when an IKM is removed from the registry. Per spec §10.9 retention rule, the IKM MUST be retained as long as any chain entry referencing it is retained; premature retirement causes affected events to fail verification at spec §7 step 7. The operational event records `chain_entries_referencing_remaining` to demonstrate the retention rule was honored at retirement time.

The SOC report attests that the institution emits these events when the corresponding action occurs and that the institution's response procedure is documented. A downstream user entity relying on the chain has assurance through the SOC opinion that these scenarios are operationally accounted for.

For SOC 2 (TSC-aligned):

- **CC7.2 (System monitoring)** — verifier output is the monitoring evidence
- **CC8.1 (Change management)** — chain spec version is captured per event; configuration changes go through change management
- **PI1.1 / PI1.2 (Processing integrity)** — chain provides the headline processing-integrity property
- **A1.x (Availability)** — for institutions claiming Availability; DR/RPO/RTO are described

For SOC 1 (financial-reporting controls):

- The chain is relevant when AI decisions feed financial reporting (e.g., AI-assisted classification of transactions affecting financial statements)
- The SOC 1 description names the chain controls in plain language; the user entity's auditor maps them into the user entity's ICFR framework

The SOC report's Section 4 description identifies the specific criteria the institution claims; the user entity tests reliance on those criteria.

## What the user entity must operate (Complementary User Entity Controls)

The chain's claims depend on the institution operating supporting controls. If the user entity is the bank itself, the user entity is the institution. If the user entity is a downstream party, the institution is the service organization and the user entity inherits.

Key CUECs the user entity verifies (`docs/control-map/CUECs.md`):

- The institution's RBAC and HSM separation of duties
- The institution's session-key reconciliation cadence
- The institution's verifier validation cadence
- The institution's IR playbook for chain-detected events
- The institution's regulator-notification framework

The user entity's auditor tests CUECs as part of the broader engagement.

## How the user entity uses the SOC report

### For financial-statement auditors (SOC 1)

The financial-statement auditor:

1. Reads the SOC 1 report's description of the chain controls
2. Identifies which controls support assertions in the financial statements
3. Plans reliance on those controls (vs substantive testing)
4. Tests the institution's CUECs as part of the broader audit
5. Documents the reliance in the audit working papers

Example: if the bank's AI is classifying revenue transactions, the financial-statement auditor relies on the chain's integrity property (CC7.2 + PI1.1) to support assertions about transaction completeness and accuracy.

### For downstream business partners (SOC 2)

The downstream business partner:

1. Reads the SOC 2 report's description and the controls in scope
2. Confirms the controls match the partner's reliance posture
3. Confirms the report's opinion is unmodified (or evaluates modifications)
4. Maintains the report on file as part of the partner's vendor-management evidence

Example: a fintech partner relying on the bank's chain to confirm authenticity of AI-driven approvals uses the SOC 2 report as evidence of the bank's controls.

### For regulators

Regulators consuming SOC reports:

1. Use the report as one input alongside their own examinations
2. Compare the institution's claims in the report to the regulator's findings
3. Identify any misalignment between claims and observed behavior
4. Use the report to inform the regulator's risk assessment of the institution

The SOC report does not replace regulatory examination; it complements.

### For audit firms planning engagements

A new audit firm assigned to the user entity:

1. Reads the SOC report to understand the institution's controls
2. Reviews the institution's CUECs
3. Plans audit testing based on reliance vs substantive
4. Documents the reliance approach in the audit plan

## What the user entity should ask

When evaluating a SOC report on a chain-of-custody implementation:

1. **Is the report unmodified?** Modified or qualified opinions warrant deeper review.
2. **Are the criteria the institution claims aligned with what the user entity needs?** A SOC 2 with Processing Integrity is the typical core; Availability and Confidentiality are additions if relevant.
3. **Is the reporting period aligned with the user entity's needs?** SOC 2 Type II covers a period; Type I covers a point in time.
4. **Are the CUECs the institution lists ones the user entity operates?** Without the user entity operating its share, the controls don't hold.
5. **Are there any anomalies or exceptions in the report that warrant follow-up?**
6. **Has the institution committed to re-issuance of subsequent SOC reports?** Continuity of the assurance is the user entity's ongoing assurance.

## Where the chain control composes with the user entity's broader posture

The chain is a service-organization control. It composes with the user entity's broader posture:

- The user entity's IAM identifies who is authorized to consume the report
- The user entity's vendor management retains the report
- The user entity's audit program determines what reliance to place on the report
- The user entity's risk function evaluates the institution's chain claims against the user entity's risk profile

The chain is one of several controls supporting the user entity's confidence; it is the integrity-bearing layer for AI-driven decisions specifically.

## Comparison to alternative

Without the chain, the user entity relying on the institution's AI decisions would have:

- Institution's word that AI decisions were as the institution describes
- Vendor logs (if applicable) without integrity guarantees
- Limited recourse in disputes

With the chain (and the SOC report attesting to its operation), the user entity has:

- Integrity-bearing records the user entity (or its auditor) can verify independently
- Defensible reliance on AI-driven decisions reflected in financials or business processes
- Standardized evidence pattern across institutions

The shift from "trust the vendor's word" to "verify the chain's integrity" is the headline value for user entities.

## What the chain does NOT do for the user entity

The user entity should not expect the chain to:

- Validate that the AI's decision was correct (chain proves the decision was made; the institution's MRM program validates correctness)
- Provide consumer-protection compliance (chain is integrity; consumer-protection is the institution's broader compliance)
- Replace the user entity's own controls (chain is one input; the user entity's broader posture remains)

The boundary is articulated to set expectations; the chain provides specific value, not all assurance the user entity needs.

## Glossary (one paragraph)

The chain uses standard cryptographic terminology: HMAC (a keyed hash), Merkle tree (a hash-of-hashes structure), HSM (a hardware-protected signing device), Ed25519 (the digital signature algorithm). The full glossary is at `docs/design/10-glossary.md`. The user entity does not need to memorize cryptographic detail; the SOC report describes the controls in plain language.

## Update cadence

This summary is updated:

- Annually with each SOC report cycle
- When the chain implementation materially changes (vendor change, spec version bump, deployment topology change)
- When user-entity audience considerations evolve (new state laws, new regulator expectations)
