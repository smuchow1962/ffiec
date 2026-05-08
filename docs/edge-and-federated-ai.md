# Edge and federated AI deployment patterns

> **What this doc is.** Deployment guidance for institutions running AI agents in distributed contexts: edge devices, federated learning architectures, on-device inference, multi-agent systems with autonomous coordination. The chain accommodates these without modification; this doc articulates the patterns.

## Edge AI

### What edge AI looks like

The institution runs AI agents on edge hardware: branch ATMs with on-device fraud-detection models, mobile devices with on-device decision-support, in-car telematics for auto-finance, IoT-style deployments for risk-sensing. The agent runs locally; communication with the institution's data center is intermittent.

### Chain operation on edge

The SDK runs in the AI agent process on the edge device:

1. Agent makes a decision; SDK constructs the chain hash locally
2. SDK persists to local storage (SQLite or equivalent) with synchronous=FULL for audit events
3. SDK exports over OTLP when connectivity is available
4. The institution's ledger receives events when connectivity is restored

The chain integrity is unaffected: the HMAC chain is computed at capture; persistence is local-first; export is asynchronous. Disconnection delays export but does not weaken integrity.

### Operational considerations

- **Local storage capacity.** Edge devices have limited storage. The SDK's local SQLite must be sized for the disconnection window plus a safety margin. Configurable per-device.
- **Memory protection.** Edge devices may have less robust memory protection than cloud hosts. The institution's threat model considers the edge attacker's increased capability against in-memory keys.
- **Session-key handshake.** Edge devices require a handshake mechanism that works over intermittent connectivity. SPIFFE/SPIRE supports edge through SPIRE Agent on the device; alternatively, hardware-issued one-shot tokens at device commissioning.
- **Reconciliation.** Session-key-id reconciliation on edge devices may show "unmatched" entries during disconnection windows; the institution baselines the unmatched count for the edge fleet separately from cloud.

### Master-key custody at edge

For edge deployments, two patterns:

**Pattern A — Edge-derived session keys.** Each edge device derives its own session keys from a per-device master held in a hardware security element (TPM, secure enclave). The chain's master is held by the institution; per-device session keys derive at the device.

**Pattern B — Bulk session-key issuance.** The institution's master-key custodian issues batches of session keys to edge devices at commissioning; the device uses the batch over its operational lifetime; rotation triggers re-commissioning.

Pattern A is more secure (no batch leak); Pattern B is simpler operationally for fleets without secure-enclave capability. Institutions choose based on the device fleet's hardware.

### Pattern B conformance for v1.0

Pattern B (bulk session-key issuance to edge devices without secure-enclave hardware) is conformant under v1.0 with the compensating-control set named in this section. Without these compensating controls, Pattern B is non-conformant for v1.0 production deployments.

The threat model's R12 entry (`design/09-threat-model.md` §6) names hardware-attested key custody as the v1.1 forward commitment for edge-device physical compromise. The compensating-control set below is the v1.0 posture that closes the practical exposure while v1.1 hardware support matures. Institutions that adopt Pattern B for v1.0 operate the full set; the SOC team confirms each control during the engagement; the FFIEC examiner expects the set to be named in the institution's CC8.1 control description.

The v1.0 compensating controls for Pattern B:

1. **Per-device IKM rotation cadence aligned with the device-class compromise window.** The institution rotates the per-device session-key set on a cadence tuned to the device class's exposure profile. Default cadences:
   - **Branch-office tooling** (banker laptops, in-branch ATM-class devices in physically controlled premises): monthly
   - **Higher-risk fleet** (mobile devices, in-vehicle telematics, devices in less-controlled environments, vendor-returned hardware in the field): weekly
   - **Federated-learning nodes participating in cross-institution training**: weekly minimum, with the cadence floor named in the consortium agreement
   The cadence is documented in the institution's CC8.1 control description and operated as a recurring procedure; rotation evidence is preserved per the institution's standard retention.

2. **IR Scenario 14 readiness with documented per-device-class playbook.** The institution operates `incident-response-playbook.md` Scenario 14 (Edge-device physical compromise — Pattern B in-service device) with a documented per-device-class playbook that names the containment, remediation, and notification steps for each device class in scope. The playbook is reviewed at least annually and is referenced in the CC8.1 control description.

3. **Device-fleet inventory monitored against the IKM roster.** The institution maintains a fleet inventory that pairs each commissioned device with its issued IKM generation. The inventory is reconciled against the IKM roster on the institution's documented reconciliation cadence (typically aligned with §10.1 weekly reconciliation). A device that produces chain entries under an IKM generation it was not issued, OR an IKM generation that does not appear on a known commissioned device, is a fleet-anomaly alert that the institution investigates as a potentially-rogue device.

4. **Compensating-control documentation in CC8.1.** The institution's CC8.1 control description names the device classes in scope, the per-class rotation cadence, the IR Scenario 14 playbook reference, the fleet-inventory reconciliation procedure, and the residual-risk acceptance language for the v1.0-vs-v1.1 gap. The documentation is the load-bearing audit-evidence shape; the SOC team and the FFIEC examiner consume the description and confirm the operational evidence matches the description.

The four-control set provides the v1.0 conformance posture for Pattern B. v1.1 will introduce hardware-attested key custody as the forward mitigation per `design/09-threat-model.md` R12; institutions running Pattern B under v1.0 plan the v1.1 migration when hardware support and spec text mature.

### Delayed-upload thresholds

Edge devices accumulate chain entries locally during disconnection windows; the entries upload when connectivity returns. The chain integrity holds across the gap (per "Chain operation on edge" above), but the institution's audit-evidence record shows a per-entry gap between `captured_at` (when the SDK chained the entry on the edge device) and `received_at` (when the institution's ledger received the entry over OTLP). Without a documented threshold for that gap, examiners and the institution's own SOC team have no normative reference for when delayed upload becomes a finding versus normal operational variance.

**The institution's per-device-class upload cadence is named in the CC8.1 control description.** The control description names the expected upload-cadence per device class — for example, intermittent-laptop fleets uploading on connection (typical: hourly to daily), branch-office tooling uploading on a documented schedule (typical: daily to weekly), federated-learning nodes uploading on a coordinated cadence (typical: weekly), and disconnected-fleet devices uploading on commissioning windows (typical: monthly). The cadence values are institution-defined per the device fleet's operational profile.

**The spec normates a default ceiling of 30 days.** A delayed upload exceeding 30 days from the entry's `captured_at` to the entry's `received_at` triggers examiner notification under the analog of `chain-of-custody-v1.md` §4.3.1 (HSM unavailability and notification). Institutions MAY operate stricter thresholds when their device-class operational profile supports them; institutions MAY operate looser thresholds only with documented examiner approval (analogous to the cadence-relaxation pattern in spec §4.2.1). The 30-day default ceiling is the failure-mode disposition for institutions that have not negotiated a different threshold with their primary regulator.

**Per-entry monitoring.** The institution monitors `received_at - captured_at` per chain entry as a metric. The monitoring is operated through the institution's standard observability stack; the metric is exposed as a histogram or per-device-class distribution so threshold breaches are visible in real time. Threshold breaches generate alerts that route to the institution's IR program for disposition: typical disposition for breaches that fall within the device-class upload cadence is no action; breaches that exceed the device-class threshold but remain under the 30-day ceiling are operational-variance observations the institution records; breaches that exceed the 30-day ceiling generate examiner notification per the institution's documented procedure.

**Evidence shape.** The chain entry's `captured_at` and `received_at` fields are the load-bearing audit evidence. The verifier reads both fields per spec §4.4 and `design/03-merkle-seal.md`; the institution's SOC team and the FFIEC examiner sample-test the per-entry gap against the institution's documented thresholds via audit-procedure P-32 (Delayed-upload threshold testing).

### Operational cost of edge

Per-device cost is small; fleet cost scales with the fleet size. The HSM operating cost stays at the institution; edge devices use derived session keys without per-device HSM cost.

## Federated learning

### What federated learning looks like

The institution trains models across distributed data sources: branch-level customer data, regional aggregations, partner-network data. The training algorithm aggregates model updates without centralizing the training data. AI agents using the trained model run wherever the use case demands.

### Chain operation in federated learning

The chain captures decisions, not training. Federated learning is the model-training phase; the chain operates on the inference / decision-making phase.

For inference:

- AI agents using the federated-trained model capture decisions normally
- The model identifier (`gen_ai.request.model`) records which version of the federated-trained model was used
- The chain is the same as for any other AI agent

For training:

- The chain doesn't directly capture training events (out of scope of the spec)
- The institution's broader MLOps program documents the training procedure
- The institution captures the resulting model version in chain events when the model is used

### What the chain doesn't cover

The chain doesn't capture:

- Federated learning aggregation rounds
- Per-data-source contributions to model training
- Model-update deltas across rounds

These are MLOps controls outside the chain's scope. The institution's MLOps program operates them; the chain composes with the resulting model deployment.

## Multi-agent systems with autonomous coordination

### What autonomous coordination looks like

Multiple AI agents coordinate to make decisions: one agent does triage, another does deeper analysis, a third makes the final decision. The agents communicate without a central orchestrator; coordination is emergent.

### Chain operation with autonomous coordination

The chain handles autonomous-coordination patterns through DAG semantics (`02-chain-construction.md` §8.1.3):

- Each agent is a chain participant with its own session key
- Each agent captures events for its own runs
- Agent-to-agent communication is captured as cross-run references (parent-child or DAG)
- The verifier reconstructs the topology at examination time

The pattern accommodates arbitrary coordination topology — directed acyclic graphs, recursive flows, multi-step handoffs.

### Operational considerations

- **Cross-run correlation.** The institution's monitoring tools may need to reconstruct the agent-coordination topology; chain events make this tractable.
- **Per-agent session keys.** Each agent has its own session key; the chain captures the boundary between agents.
- **Handoff events.** When one agent hands control to another, the chain captures the handoff as a chain event with `parent_run_id` linkage.

### Examination considerations

For examination, the verifier walks each run independently. Cross-run topology is the institution's analytical responsibility (typically reconstructed in the institution's observability stack).

## On-device inference (mobile, IoT)

### What on-device inference looks like

The AI model runs on a mobile device, IoT device, or other on-device target. The institution's mobile app or device firmware includes the inference logic.

### Chain operation on-device

Same pattern as edge AI:

- SDK runs in the on-device process
- Local persistence buffers chain events
- OTLP export when connectivity is available

### Privacy considerations

On-device inference may process sensitive data. The chain's privacy-by-design pattern (`docs/privacy-by-design.md`) applies: tokenize PII before canonicalization; the chain captures tokens, not personal data; the privacy-store at the institution holds the mapping.

For on-device contexts, the privacy-store handshake is important: the device can submit chain events with tokens; the device may or may not have access to the privacy-store directly. The institution decides based on device-trust posture.

## Cross-institution federated patterns

### Multi-institution AI

Some emerging patterns involve AI across institutions: shared fraud-detection models across consortium banks, multi-bank syndication AI, joint-decision-making AI. These are rare but real.

### Chain operation across institutions

Each institution operates its own chain. Cross-institution AI captures events in each institution's chain at the events that institution observes. The institutions may correlate analytically; the chain doesn't model cross-institution coordination directly.

For consortium-style deployments, the consortium may operate a shared chain (with consortium-level master key custody) or separate per-member chains (with per-member master keys). The choice affects the chain's tenant_id structure but not the substance.

## Operational guidance

For institutions deploying chain on edge / federated / on-device:

1. **Start with a controlled pilot.** Edge deployments are operationally complex; start with one device type, one use case.
2. **Size local storage.** Disconnection windows drive local storage requirements; budget conservatively.
3. **Plan handshake mechanism.** SPIFFE/SPIRE works for many edge contexts; alternative mechanisms exist for hardware-restricted devices.
4. **Document the topology in the control description.** The institution's control description names the edge / federated / on-device pattern in scope.
5. **Adapt the IR playbook.** Edge devices have different attack surface; the IR playbook adapts accordingly.
6. **Monitor disconnect-reconnect telemetry.** The institution's observability stack tracks edge devices' connection state; chain export delays correlate with connection issues.

## Future-state considerations

As edge / federated / on-device AI patterns mature, the spec may add:

- **`ffiec.chain.deployment_context`** — captures the deployment context (cloud, edge, on-device, federated)
- **Edge-specific telemetry attributes** — connectivity state, local-buffer depth
- **Federated-learning-specific attributes** — model-version provenance for federated-trained models

These are v1.1 candidates; v1.0 supports edge / federated / on-device use cases without modification.

## Related documents

- [`design/02-chain-construction.md`](design/02-chain-construction.md) — multi-process patterns
- [`operator-guide.md`](operator-guide.md) — runtime operations
- [`incident-response-playbook.md`](incident-response-playbook.md) — IR for chain-detected events
- [`privacy-by-design.md`](privacy-by-design.md) — privacy patterns
- [`at-scale-operations.md`](at-scale-operations.md) — fleet-scale operations
