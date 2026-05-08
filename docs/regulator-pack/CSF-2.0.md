# NIST Cybersecurity Framework 2.0 mapping

> **What this doc is.** Mapping of chain-of-custody primitives to NIST CSF 2.0 functions, categories, and subcategories. Cybersecurity examiners working from a CSF assessment use this to confirm which subcategories the chain satisfies.

## Functions and categories addressed

| CSF 2.0 Function | Category | Subcategory | Chain primitive |
|---|---|---|---|
| GOVERN (GV) | Risk Management Strategy (RM) | GV.RM-06 (Risk responses) | Threat model `09-threat-model.md` documents accepted residual risks |
| GOVERN (GV) | Cybersecurity Supply Chain (SC) | GV.SC-01..10 | Supply chain `07-verifier-design.md` §8.4 (cosign, GPG, SBOM, reproducible builds) |
| IDENTIFY (ID) | Asset Management (AM) | ID.AM-08 (Systems are identified and managed throughout lifecycle) | Tenant key registry; per-tenant key isolation |
| IDENTIFY (ID) | Risk Assessment (RA) | ID.RA-09 (Authenticity of hardware and software is assessed) | Verifier supply chain (`07-verifier-design.md` §8.4) |
| PROTECT (PR) | Identity Management (AA) | PR.AA-05 (Access permissions are managed per least privilege) | RBAC defense-in-depth (`06-ledger-server-design.md` §3.5) |
| PROTECT (PR) | Data Security (DS) | **PR.DS-01 (Data at rest are protected)** | Encryption at rest is the institution's standard control; chain composes |
| PROTECT (PR) | Data Security (DS) | **PR.DS-02 (Data in transit are protected)** | TLS 1.3 minimum (`05-otlp-wire.md` §5.1) |
| PROTECT (PR) | Data Security (DS) | **PR.DS-06 (Integrity-checking mechanisms verify data integrity)** | **The chain itself.** HMAC chain + Merkle seal + HSM signature |
| PROTECT (PR) | Platform Security (PS) | PR.PS-01 (Platforms are configured per policy) | Memory-protection per platform (`02-chain-construction.md` §4.2) |
| PROTECT (PR) | Platform Security (PS) | PR.PS-06 (Secure software development) | Reproducible builds + cosign signing + GPG manifest |
| DETECT (DE) | Continuous Monitoring (CM) | **DE.CM-09 (Computing resources are monitored to find anomalous events)** | Verifier anomaly reporting; operational events list (`06-ledger-server-design.md` §7.3.1) |
| DETECT (DE) | Adverse Event Analysis (AE) | DE.AE-02 (Potentially adverse events are analyzed) | Verifier strict-mode + late-binding-rate threshold |
| DETECT (DE) | Adverse Event Analysis (AE) | DE.AE-03 (Information about adverse events is correlated) | Operational events emitted with correlation IDs |
| DETECT (DE) | Adverse Event Analysis (AE) | **DE.AE-03 (Anomalies and events / event correlation — routing-decision coverage)** | **`audit.routing.*` chain entries per spec §4.4.1.** Routing decisions are events whose absence-of-correlation (a multi-provider deployment with no routing chain entries, or LLM calls whose preceding routing entry is missing) is itself anomalous. The chain captures four routing event types (`attempt`, `success`, `failover`, `circuit_state_change`) and the cybersecurity examiner correlates them against the LLM-call chain entries per `audit-procedures.md` P-33 |
| GOVERN (GV) | Roles, Responsibilities, and Authorities (RR) | **GV.AC-05 (Roles, responsibilities, authorities — routing-policy governance)** | **Routing-policy change-management evidence per spec §4.4.1.** Who governs routing-policy changes is a documented authority-allocation question; the institution's CC8.1 control description names the routing-policy versioning procedure, the change-approval authority, and the chain decorator that emits the four routing event types. The `audit.routing.policy_version` attribute on each routing chain entry binds a captured decision to a specific policy version so MRM committees correlate routing-behavior change-points with policy-change records |
| RESPOND (RS) | Incident Management (MA) | RS.MA-01..05 | Incident response playbook `docs/incident-response-playbook.md` |
| RESPOND (RS) | Incident Analysis (AN) | RS.AN-03 (Analysis is performed to establish what happened) | Verifier `walk` and `diff` modes (`07-verifier-design.md` §7.2, §7.3) |
| RECOVER (RC) | Incident Recovery Plan Execution (RP) | RC.RP-04 (Critical mission functions are restored) | DR section `06-ledger-server-design.md` §7.5 |
| RECOVER (RC) | Communications (CO) | RC.CO-03 (Recovery activities are communicated to internal/external stakeholders) | 72-hour seal-delay notification (`04-hsm-custody.md` §5.2) |

## Headline mapping

The chain is, primarily, a **PR.DS-06** control. PR.DS-06 reads: *"Integrity-checking mechanisms are used to verify hardware integrity, software integrity, firmware integrity, and information integrity."* The chain provides the information-integrity portion of this subcategory in a way that satisfies the "verify" requirement at FIPS-140-2-Level-3 protection.

The other subcategories listed are properties of the chain's surrounding deployment, not of the chain itself. Their satisfaction depends on the institution's broader cybersecurity posture; the chain composes with that posture rather than replacing it.

## Four attack approaches mapped to CSF 2.0 subcategories

The four attacks named in `00-overview.md` §2 each have a load-bearing CSF subcategory and a load-bearing operational evidence trail. The cybersecurity examiner constructing an evidence-to-subcategory matrix uses this table directly:

| Attack approach | Load-bearing CSF subcategory | Load-bearing operational evidence | IR scenario |
|---|---|---|---|
| **§2.1 Host-side chain tampering** | PR.DS-06 (information integrity) | `chain.verification_failure step=9` (payload_hash MAC mismatch) — verifier-emitted at ingest re-verification or examination-time | Scenario 1 (chain hash mismatch at ingest) |
| **§2.2 Cross-tenant or cross-version key confusion** | PR.DS-06 + ID.AM-08 (asset/key lifecycle integrity) | `chain.verification_failure step=8` (key_fingerprint mismatch); `master.reconciliation_completed.fingerprint_unmatched_count` weekly reconciliation event | Scenario 7 (key_fingerprint mismatch detected at verification) |
| **§2.3 Server-side or privileged-insider history rewrite** | PR.DS-06 + DE.CM-09 (continuous monitoring) | `chain.verification_failure step=10` (Merkle root mismatch); `chain.verification_failure step=11` (signature verification failed); `seal.job_completed` events with `merkle_root` for cross-correlation | Scenarios 2, 3 (Merkle root mismatch; signature verification failed) |
| **§2.4 Examination-time tooling subversion** | GV.SC (cybersecurity supply chain) + ID.RA-09 (authenticity of hardware and software) | Reproducible-build log (institution-archived); cosign + GPG validation logs; `bundleverify` output; institution's SLSA-verifier validation log | Scenario 3 (signature verification failed — binary-signature variant per `supply-chain.md`) |

The four-attack indexing lets the cybersecurity examiner construct an evidence matrix at scale. Each attack approach has a primary subcategory and a named operational artifact; the IR scenario is the institution-side response when the operational artifact fires. The information is the same as the per-primitive table above, indexed differently for the examiner's working pattern.

## How to use

In a CSF-aligned cybersecurity examination:

1. The examiner identifies which subcategories the institution claims to satisfy
2. For PR.DS-06, the chain is the headline evidence — the verifier output is the artifact
3. For the surrounding subcategories, the chain is supporting evidence — the institution's broader controls are the headline; the chain is one of several supporting controls

## Subcategories the chain does NOT satisfy

To set expectations explicitly, the chain does NOT satisfy:

- **PR.DS-01 (Data at rest)** by itself; the institution's storage encryption is the headline control
- **PR.AA (Identity management) at the principal level**; the chain captures tenant-level identity (`key_fingerprint`), not user identity (see `09-threat-model.md` §3.1.2)
- **DE.CM-01..08 (specific monitoring categories beyond CM-09)**; those depend on the institution's broader monitoring program

The boundary keeps the chain's claims defensible and lets the institution's other controls do their job without overlap.

## Operational events bound to CSF subcategories

The spec §10.2 operational events provide subcategory-level evidence at finer granularity than the headline mapping suggests. The institution's SOC team and the cybersecurity examiner consume these events as substantive evidence; the binding below resolves which subcategory each event most strongly supports.

| Operational event (spec §10.2) | Primary CSF subcategory | Secondary subcategories |
|---|---|---|
| `ledger.startup` | DE.CM-09 (resources monitored) | GV.PO-01 (policy adherence) |
| `ledger.hsm_session_opened` | PR.PS-01 (platforms configured per policy) | DE.CM-09 |
| `seal.job_started` | DE.CM-09 | RS.MA-04 (incidents categorized) |
| `seal.job_completed` | PR.DS-06 (information integrity) | DE.AE-03 (events correlated via `correlation_id`) |
| `seal.job_failed` | DE.AE-02 (potentially adverse events analyzed) | RS.MA-01..05 (Scenario 5 IR) |
| `chain.verification_failure` (any `step`) | PR.DS-06 (the verifier's failure is the integrity-check signal) | DE.AE-02; routes to RS.MA per spec §7 step → IR scenario |
| `chain.verification_failure step=8` (key_fingerprint mismatch) | PR.DS-06 + ID.AM-08 (asset/key lifecycle integrity) | RS.MA-04 (Scenario 7 IR) |
| `chain.verification_failure step=7` (unknown_key_version) | ID.AM-08 (lifecycle integrity) | RS.MA-04 (Scenario 8 IR) |
| `audit_file.truncation_detected` | RC.RP-04 (critical functions restored) | DE.AE-02 |
| `hsm.operation_success` | PR.PS-01 (platform integrity) | DE.CM-09 |
| `hsm.operation_failure` | DE.AE-02 | RS.MA-01..05 (Scenario 3 IR if signature failure) |
| `config.reload` | GV.SC-08 (cybersecurity supply chain mgmt — config changes tracked) | PR.PS-01 |
| `master_key.rotated` | PR.AA-05 (least privilege at IKM lookup) | ID.AM-08 |
| `master_key.rotation_observed` | DE.CM-09 | DE.AE-03 (correlate with `master_key.rotated`) |
| `master_key.retired` | ID.AM-08 (lifecycle) | DE.AE-02 (premature retirement detection) |
| `master.reconciliation_completed` | DE.CM-09 (continuous monitoring) | PR.AA-05; the headline P-6 audit-procedure evidence |
| `mirror.reconciliation_completed` (institution-defined; supply-chain mirror's signature-validation reconciliation per `supply-chain.md`) | DE.CM-09 + GV.SC-04 (cybersecurity supply chain — third-party assessments) | DE.AE-02 |
| SLSA-verifier validation log (institution-side, archived per supply-chain doc) | GV.SC-04 + ID.RA-09 (authenticity of hardware and software) | The deployment-gating control evidence; binding the L3 provenance attestation to the institution's deployment pipeline |
| `regulator_fingerprint.rotation_received` (institution-defined; trust-anchor lifecycle) | DE.CM-09 + GV.SC-04 | DE.AE-03 (correlate with regulator notification channel) |
| `regulator_fingerprint.rotation_validated` (institution-defined) | DE.AE-03 + GV.SC-04 | The two-channel validation evidence per `09-threat-model.md` §2.9 reception procedure |
| `regulator_fingerprint.installed` (institution-defined) | CC8.1 (change management) + GV.SC-04 | The trust-anchor-cache update is a change-managed event |
| Annual `KEY-DR-DRYRUN-{year}.asc` consumption log (institution-side, archived per supply-chain doc cold-DR-key lifecycle) | GV.SC-04 + ID.RA-09 | Evidence the institution exercises the cold-DR fallback path; pairs with IR Scenario 11 |

The cybersecurity examiner working a CSF assessment uses this binding to construct an evidence-to-subcategory matrix per institution, sample-tested against the institution's log store.

## GOVERN function alignment (one-page institution adoption)

CSF 2.0 made GOVERN a top-level function. Examiners now expect to see how integrity controls flow from risk-tolerance statements down to operational evidence. The chain composes with the institution's GOVERN function as follows:

**GV.OC (Organizational Context).** The institution's AI-decision audit-trail integrity tolerance is articulated in its information-security policy framework. The chain is the technical control implementing that tolerance for AI-driven decisions. The institution's risk-tolerance statement should name (a) the per-day Merkle seal as the load-bearing tampering-detection control, (b) the per-event HMAC chain as the load-bearing real-time-tampering-detection control, and (c) the regulator-held public key as the trust-anchor for the integrity claim.

**GV.RM (Risk Management Strategy).** The threat model in `09-threat-model.md` §6 documents the institution-side residual risks (R1–R13). The institution's risk-management framework should adopt the residual-risk register as the institution's accepted residual posture for AI-decision audit-trail integrity, OR document where the institution's posture differs (typical: a smaller institution may accept Scenario 9 (truncation) at lower severity than the playbook default). Risk responses (GV.RM-06) flow from the IR playbook scenarios.

**GV.RR (Roles, Responsibilities, and Authorities).** The IR playbook `Roles` table names the chain-specific roles; the institution maps these to its existing IR role assignments. The chain-ops team is typically a subset of the institution's broader cyber-ops team, with the chain-specific authority being "request HSM signing operations under the seal-job role" and "operate the IKM custodian's reconciliation procedure."

**GV.PO (Policy).** The institution's information-security policy framework references the spec (`spec/chain-of-custody-v1.md`) as the normative artifact for AI-decision audit-trail integrity. The institution's control description in its SOC opinion or examination response includes the chain's CUEC obligations (`docs/control-map/CUECs.md`).

**GV.OV (Oversight).** The MRM committee oversight (per `MRM-COMMITTEE-BRIEF.md`) and the audit committee oversight (per `audit-committee-summary.md`) are the institution's GV.OV evidence. Both committees receive verifier output and reconciliation reports on a documented cadence.

**GV.SC (Cybersecurity Supply Chain).** The cosign + GPG dual-signing trust path (`supply-chain.md`) is the institution's GV.SC evidence for the verifier's authenticity. The institution-side cadence for verifier-binary validation (per `regulator-pack/deployment-package.md`) is the institution's substantive operational evidence.

The institution copies this section's content into its own GV-alignment policy document, edits for institution-specific context, and presents the result as the institution's GOVERN-function evidence for the chain control.

## Risk-tolerance statement template

The GOVERN-function alignment above describes what the institution's risk-tolerance statement should name. The template below is the copy-pasteable form. The institution drops it into its information-security policy framework, fills in the placeholder fields, and presents the result as the institution's documented risk-tolerance statement for AI-decision audit-trail integrity. The language is anchored explicitly in CSF 2.0 GOVERN-function vocabulary so the alignment is visible to a CSF-aligned examiner without further translation.

Placeholders use angle-bracket form (`<institution name>`, `<primary regulator name>`, etc.) so a search-and-replace pass populates the document. Italicised parenthetical guidance is removal-on-adoption; the institution deletes those phrases when finalising the statement.

> **AI-decision audit-trail integrity — risk-tolerance statement**
>
> **Institution.** `<institution name>` (the "institution"), supervised by `<primary regulator name>` (the "primary regulator") and, where applicable, `<additional regulator name(s)>`. *(List each supervisor whose framework the institution is operating against.)*
>
> **Scope.** This statement governs the integrity of audit-trail records for AI-driven decisions executed by the institution's production AI agents. It applies to every tenant and every business line whose decisions are captured under the chain control, as enumerated in the institution's per-tenant scope register `<scope register location>`.
>
> **Risk-tolerance posture.** The institution accepts no tolerance for undetected tampering of AI-decision audit-trail records. Tampering that occurs MUST be detected at the next verifier execution and surfaced through the institution's incident-response program. The institution's tolerance for the time-to-detection is `<detection-window value>` (typical: at most one daily seal cycle plus one verifier-execution cadence).
>
> **Load-bearing controls.** The institution names three load-bearing controls implementing this tolerance. (a) The per-day Merkle seal is the load-bearing tampering-detection control: any modification to a sealed-day's records produces a Merkle-root mismatch the verifier reports at step 10. (b) The per-event HMAC chain is the load-bearing real-time-tampering-detection control: any modification to an in-flight event invalidates the payload-hash MAC the verifier reports at step 9. (c) The regulator-held public key is the load-bearing integrity trust anchor: the institution's signing public key is registered with `<primary regulator name>` so the integrity claim is verifiable independently of any institution-side or vendor-side infrastructure.
>
> **Custody locations.** The institution's Ed25519 signing keypair is held in HSM custody at `<HSM custody location, e.g., AWS CloudHSM us-east-1 cluster X>` under FIPS 140-2 Level 3 protection. The institution's per-tenant master HMAC key (IKM) is held in `<IKM custody location, e.g., HashiCorp Vault tenant-controlled namespace>` under the institution's documented key-custody procedure `<custody procedure document reference>`.
>
> **Trust anchor publication.** The institution publishes the SHA-256 fingerprint of its signing public key to `<primary regulator name>` through `<publication channel, e.g., the regulator's secure-correspondence portal>` and re-publishes on rotation per the institution's key-rotation schedule. The trust-anchor lifecycle is documented in `<trust-anchor procedure reference>`.
>
> **Residual risk acknowledgement.** The institution accepts the residual risks enumerated in the chain's threat model `09-threat-model.md` §6 (R1–R13), with institution-specific dispositions documented in `<residual-risk register reference>`. Residual risks not accepted by the institution are documented with mitigating controls in the same register.
>
> **Oversight.** The institution's MRM committee and audit committee receive verifier output, reconciliation reports, and incident-response evidence on the cadence documented in `<oversight charter reference>`. Both committees have authority to escalate findings under the institution's incident-management program.
>
> **Review and revision.** This statement is reviewed annually as part of the institution's information-security policy framework refresh, and on the occurrence of a chain-relevant material change (regulator-published guidance update, spec major version change, accepted-residual-risk register change). The CISO `<CISO name or title>` is the named owner; the responsible approval body is `<approval body, e.g., the institution's risk committee>`.

The template lands at the right level of CSF GOVERN evidence: it names the load-bearing controls, the custody locations, the trust-anchor publication path, the residual-risk posture, and the oversight cadence. Each item is something a CSF-aligned examiner expects to see in a risk-tolerance statement; each item ties back to a chain primitive the institution can produce evidence for. The institution shortens or extends individual subsections to fit its policy-framework house style without weakening the alignment.

## Parity analysis: CSF 2.0 - FFIEC IT Handbook

The chain has two parallel mappings. CSF 2.0 (this document) and FFIEC IT Handbook (`handbook-mapping.md`). Some institutions are examined under both frameworks — most commonly tier-1 banks whose IT examiner works the Handbook and whose cybersecurity examiner works CSF, or institutions whose internal audit function operates both frameworks at once. For those institutions the parity question matters: does a single body of chain evidence satisfy both frameworks equally, or are there subtle disagreements about which control is primary versus supporting?

The short answer is that the two mappings are coherent at every major subcategory the chain touches. Both frameworks agree on which chain primitive is the headline control and which is supporting. The remainder of this section walks the major subcategories, names the corresponding Handbook section, and confirms the agreement. Where a divergence in framing exists, it is named explicitly with the reasoning so the institution's compliance team is not surprised at examination.

### Subcategory-by-subcategory walk

**PR.DS-06 (Information integrity) - II.C.10 (Logging).** Both mappings name the chain itself as the headline control. The CSF mapping table calls PR.DS-06 the "headline mapping" and italicises it; the Handbook mapping calls II.C.10 the "headline mapping" and bolds it. The chain primitive named is identical in both: HMAC chain plus Merkle seal plus HSM signature. Both mappings name the verifier output as the artifact. Both extend to `master.reconciliation_completed` as the supporting weekly-discipline evidence per spec §10.1. No divergence. An institution producing verifier output and reconciliation evidence satisfies both PR.DS-06 and II.C.10 with the same artifact.

**PR.AA-05 (Access permissions per least privilege) - II.C.5 (Logical access).** The CSF mapping points PR.AA-05 at the ledger server's RBAC defense-in-depth (`06-ledger-server-design.md` §3.5). The Handbook mapping points II.C.5 at the session-key handshake security floor (`02-chain-construction.md` §4.1.1) — workload identity attestation, mTLS, HSM-issued tokens. Both controls are within the chain's surrounding deployment but are distinct primitives: PR.AA-05 is about who can call the ledger-server admin surface, II.C.5 is about how an SDK process authenticates to obtain a session key. The two are complementary, not divergent. An institution presenting both controls satisfies both subcategories; an institution presenting only one is partially covered against either framework. Naming the distinction explicitly avoids the appearance of double-counting.

**PR.IP-04 (predecessor — backups) and II.E (Change management).** CSF 2.0 retired PR.IP-04 in favour of subcategories under PR.PS (Platform Security) and ID.AM (Asset Management) lifecycle. The Handbook retains II.E as the change-management control objective. The chain's `format_version` field is the load-bearing change-management primitive in both framings: the Handbook mapping names it explicitly under II.E; the CSF mapping treats it as supporting evidence under ID.AM-08 (lifecycle integrity) and GV.SC-08 (config changes tracked). The framing differs (Handbook treats change management as a first-class control objective; CSF distributes the same evidence across two GOVERN/IDENTIFY subcategories) but the underlying evidence the institution produces is identical: documented change management for any `format_version` change, alignment between SDK-stamped values and ledger-stored values and verifier-version compatibility, plus algorithm-posture transition evidence per spec §7 step 11. No real divergence; a presentation difference an examiner working both frameworks navigates without effort.

**DE.AE-03 (Information about adverse events is correlated) - II.C.11 (Incident response, detection portion).** Both mappings name correlation IDs on operational events as the headline mechanism. The CSF mapping extends DE.AE-03 to the `seal.job_completed` plus `master_key.rotation_observed` correlation pattern; the Handbook mapping treats II.C.11 as covering the verifier anomaly reporting plus the dedicated IR scenarios for `key_fingerprint mismatch`, `unknown_key_version`, and `audit_file.truncation_detected`. Both are correct framings of the same evidence — the chain emits correlation-ID-bearing events, the verifier correlates them at examination time, the IR playbook documents the response. No divergence.

**GV.AC-05 (— note CSF 2.0 has no GV.AC-05; the closest cybersecurity-supply-chain access control is GV.SC) and II.C.13 (Cryptographic controls).** A clarification rather than a true subcategory mapping: CSF 2.0 does not have a GV.AC-05 subcategory; cybersecurity supply chain access discipline is under GV.SC-01..10, and access controls are PR.AA. The Handbook's II.C.13 maps to a CSF combination of PR.DS-06 (cryptographic integrity), GV.SC (supply chain authenticity), and ID.AM-08 (lifecycle integrity of cryptographic material). The chain's three rework primitives — per-entry `key_fingerprint`, IKM minimum 32 bytes, software-key adapter compile-time exclusion — appear under all three CSF subcategories as supporting evidence and under II.C.13 as headline cryptographic-controls evidence. The institution presents the same primitives in both framings; the headline-versus-supporting distinction does not change the evidence the institution produces.

**II.B (Risk Management) - GV.RM (Risk Management Strategy) and ID.RA (Risk Assessment).** Both mappings point at the threat model `09-threat-model.md` as the headline artifact. The Handbook treats II.B as a single control objective; the CSF mapping splits the same evidence across GV.RM-06 (risk responses) and ID.RA-09 (authenticity assessment). No divergence in evidence; a presentation split that follows CSF 2.0's GOVERN-function elevation.

### Mapping-summary table

| Chain primitive | CSF 2.0 primary | CSF 2.0 supporting | Handbook primary | Handbook supporting | Parity status |
|---|---|---|---|---|---|
| HMAC chain + Merkle seal + HSM signature | PR.DS-06 | DE.AE-02; DE.CM-09 | II.C.10 | III.B (Audit) | Agreed: both name the chain as headline |
| Per-entry `key_fingerprint` | PR.DS-06 | ID.AM-08; GV.SC | II.C.13 | III.B | Agreed: cryptographic-controls evidence in both |
| Weekly `master.reconciliation_completed` | DE.CM-09 | PR.AA-05 | II.C.10 | II.C.13 | Agreed: continuous-monitoring discipline in both |
| `format_version` change-management primitive | ID.AM-08 | GV.SC-08; PR.PS-01 | II.E | III.A (Operations) | Agreed: change-management evidence in both |
| RBAC defense-in-depth (ledger-server) | PR.AA-05 | DE.CM-09 | II.C.5 | III.B | Agreed: access-management control in both |
| Session-key handshake security floor | PR.AA-05 | PR.PS-01 | II.C.5 | II.C.13 | Agreed: workload-identity-and-cryptographic floor in both |
| Reproducible builds + cosign + GPG | GV.SC-01..10 | ID.RA-09; PR.PS-06 | II.C.20 | IV (Outsourcing) | Agreed: supply-chain authenticity in both |
| IR playbook scenarios | RS.MA-01..05 | RS.AN-03; DE.AE-02 | II.C.11 | II.C.10 | Agreed: incident-response evidence in both |
| Verifier output | PR.DS-06 (HQ) | RS.AN-03; DE.CM-09 | II.C.10 (HQ) | III.B; III.C; IV.B | Agreed: independent audit artifact in both |
| Threat model + residual-risk register | GV.RM-06 | ID.RA-09 | II.B | II.E | Agreed: risk-management posture in both |

"HQ" marks the headline control assignment in each framework. The agreement column reflects the analysis above: every major chain primitive is named as the headline control in one framework and as headline or substantive supporting in the other, with the underlying evidence identical.

### Conclusion: examiners working both frameworks rely on the same evidence

The two mappings agree on every major subcategory the chain touches. An institution whose IT examiner works the FFIEC IT Handbook and whose cybersecurity examiner works CSF 2.0 produces one body of chain evidence — verifier output, reconciliation events, IR playbook execution evidence, threat-model documentation, supply-chain validation logs — and that same body of evidence satisfies both frameworks. The presentation differs (the Handbook indexes by booklet section; CSF indexes by function and subcategory) but the underlying control set is the same.

Where the two frameworks split a single Handbook control objective into multiple CSF subcategories (notably II.E into ID.AM-08 plus GV.SC-08, and II.B into GV.RM plus ID.RA), the institution's compliance team prepares the evidence once and cross-references it under both framings. Where the two frameworks combine multiple Handbook objectives into a single CSF subcategory (notably PR.DS-06 absorbing portions of II.C.10 plus II.C.13), the same evidence carries both citations. No subcategory in the chain's coverage area is satisfied by one framework's view and not the other.

For an institution sequencing its examination calendar across a CSF assessment and an FFIEC IT examination, the practical implication is that the verifier-output binder, the reconciliation-event archive, and the IR-playbook execution log are the load-bearing artifacts in both engagements. The institution's compliance team prepares the binder once, presents it under PR.DS-06 + II.C.10 + III.B in the cybersecurity examination and under II.C.10 + II.C.13 + III.B in the IT examination, and the load-bearing chain evidence holds equally in both.
