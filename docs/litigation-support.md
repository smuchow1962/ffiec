# Litigation Support Guide — FFIEC Chain-of-Custody v1.0

> **What this doc is.** The institution's playbook for moving FFIEC chain-of-custody evidence from a cryptographic record into a courtroom. The chain is mathematically sound. This document is the bridge from mathematics to FRE 901 authentication, FRCP 26 discovery, FRE 702 expert qualification, FRCP 37(e) spoliation defense, and the international comparable standards (UK CPR, EU eIDAS, Canadian SCC).

> **What this doc is not.** A substitute for retained counsel. The procedures here describe the institution's posture; counsel applies them to specific litigation.

---

## 1. Audience and scope

This guide is written for two readers. The first is the institution's internal counsel, records-custodian, IT director, and incident-response lead — the people who must be ready to produce evidence, lay foundation, and answer subpoenas. The second is the institution's outside counsel and retained expert, who use this as the procedural backbone for a Daubert hearing or a deposition.

The scope is U.S. federal litigation under the Federal Rules of Evidence and the Federal Rules of Civil Procedure, with comparable-standard appendices for UK, EU, and Canadian jurisdictions.

**The single most important framing.** The chain proves what the AI said and proves no one tampered with the record of what it said. The chain does not prove the AI was correct, the AI was unbiased, or the AI's output complied with policy. Those are separate evidentiary claims that require separate proof. Section 14 returns to this in detail; every other section assumes the reader has internalized it.

Cross-references to the spec use the form "spec §N" and resolve to `spec/chain-of-custody-v1.md`.

---

## 2. Daubert four-factor grounding

A federal trial court applies the Daubert standard (*Daubert v. Merrell Dow Pharmaceuticals*, 509 U.S. 579 (1993); codified at FRE 702) to admit expert testimony based on scientific methodology. The four factors are testability, peer review, known error rate, and general acceptance. The chain's design satisfies each.

### 2.1 Testability

Spec §7 defines a 12-step verification procedure with an explicit failure reason for each step. The reference test-vector corpus (spec §8) ships positive and negative vectors for over 20 attack scenarios — entry tampering, fingerprint mismatch, root mismatch, signature forgery, IKM reuse, witness-mode boundary cases, and more. Each vector is independently runnable against the open-source verifier.

The institution's expert testifies: "The chain's integrity claim is testable because the verifier procedure is deterministic, the test vectors are public, and the failure modes are enumerated. Anyone with the ledger snapshot, the public key, and the verifier binary can reproduce the verification result."

### 2.2 Peer review and publication

The spec was developed under FFIEC working-group review with cross-disciplinary input (banking supervision, cryptography, digital forensics, model-risk management). The reference implementation is Apache-2.0 licensed and has been reviewed by outside cryptographic and forensic experts whose feedback drops are preserved in `docs/feedback/outside/`.

The underlying primitives are not novel:

- HMAC-SHA-256 is specified in NIST FIPS 198-1 and RFC 2104.
- SHA-256 is specified in NIST FIPS 180-4.
- Ed25519 is specified in NIST FIPS 186-5 and RFC 8032.
- The RFC 6962 Merkle-tree construction has been deployed in Certificate Transparency since 2013, surviving sustained adversarial scrutiny.
- JCS canonical form is specified in RFC 8785.

The institution's expert testifies: "The cryptographic substrate is decades-old, FIPS-approved, and deployed at planetary scale in Certificate Transparency. The novelty is the integration into a regulator-aligned ledger schema, and that integration was reviewed by the FFIEC working group with public feedback drops."

### 2.3 Known error rate

Daubert's third factor asks for a quantified false-negative rate. The chain's threat model (spec §9) lists the conditions under which tampering would go undetected:

The chain's integrity claim relies on three independent detection layers. Per-entry HMAC catches modification of any single entry. Daily Merkle root catches insertion or deletion of entries within a sealed day. HSM-backed Ed25519 signature catches forgery of the seal record itself. Hidden tampering would require simultaneous compromise of the IKM, the ledger storage, and the HSM — three separately-hardened systems with separate access controls.

The institution's expert testifies: "Under the documented threat model, no known practical attack defeats all three detection layers. The false-negative rate is bounded by the cryptographic strength of HMAC-SHA-256 and Ed25519, both of which are NIST-current and not known to be broken."

### 2.4 General acceptance

The integration is new; the components are not. Each component is named in NIST FIPS standards (180-4, 186-5, 198-1) and is deployed in production systems used by the U.S. federal government, the financial-services industry, and global internet infrastructure. The verifier is open-source. The test vectors are public.

The institution's expert testifies: "Each cryptographic primitive in this chain is generally accepted in the relevant scientific community. The NIST endorsement is dispositive on this factor."

### 2.5 Hash-collision risk framing

A predictable cross-examination challenge: "SHA-256 collisions are theoretically possible — how do you know this entry wasn't swapped with a colliding one?"

The institution's expert answers: "SHA-256 collision resistance is current per NIST guidance. No practical SHA-256 collision is known. A collision attack against the chain would require finding a second JSON document that produces the same canonical bytes and the same MAC under an HMAC key the attacker does not possess. The probability is computationally negligible under current technology. NIST monitors this risk and will issue transition guidance if the algorithm weakens; until then, we rely on the algorithm's integrity per NIST's current risk assessment."

### 2.6 Algorithm rollover

Chain entries created in 2026 may be evidence in cases tried in 2032. If SHA-256 weakens during the retention window, the institution must be ready.

The institution's posture, on the record before any dispute arises:

| Time horizon | Action |
|---|---|
| 2026 — present | Operate under v1.0 with SHA-256 + HMAC-SHA-256 + Ed25519. Monitor NIST guidance. |
| Pre-break warning | NIST issues a transition advisory. Institution begins v1.x migration planning. |
| Post-break | Re-verify or migrate affected entries under the new algorithm per published NIST transition guidance. Old verifier output remains the contemporaneous evidence of the chain's integrity at the time of capture. |

Spec §4.3.2 already adds dual-algorithm support (Ed25519 + Dilithium) for seals as a forward-compatibility path. Per-entry MAC migration is a v1.x or v2.0 work item.

The institution's expert testifies: "The chain was created under an algorithm considered secure at the time and current standards still endorse that algorithm. If standards change, the institution will migrate per NIST guidance. The chain entry's evidentiary weight at the time of capture is what matters; the contemporaneous verifier output is the load-bearing artifact."

---

## 3. FRE 901 authentication

FRE 901 requires the proponent of evidence to produce evidence "sufficient to support a finding that the item is what the proponent claims it is." Subdivision (b)(9) speaks directly to system-generated records: "evidence describing a process or system used to produce a result and showing that the process or system produces an accurate result."

The chain authenticates two things. The verifier output authenticates the **result** — that the chain entries verify and have not been tampered with. The institution authenticates the **process** — the deployment that produced the entries.

### 3.1 Documentary evidence the institution must retain

The institution preserves the following artifacts for the length of any litigation hold plus the standard retention period:

| Artifact | Purpose | Where it lives |
|---|---|---|
| SDK version manifest | Names the SDK build that produced entries during the period | Build system; release notes |
| SDK source-code hash | Proves the SDK binary corresponds to the named source | Reproducible-build attestation |
| HSM configuration snapshot | Documents the HSM model, firmware, FIPS level, signing-key slot | Vendor management; HSM admin console |
| Signing-key rotation history | Documents which signing key was active for each seal | Key management system |
| Daily seal-job logs | Confirms successful signing for each seal record | Operational logs |
| Change-management records | Documents any configuration change touching the chain | Change-management system |
| Verifier output for the affected period | Shows PASS for the period under dispute | Verifier bundle |
| `verifier.run_completed` operational events | Documents who ran the verifier, when, and on which host | Operational events store (per soc-pack) |

The IT witness lays foundation from these artifacts without re-engineering the system at deposition time.

### 3.2 Foundation testimony — the witness's seven sentences

The IT witness needs a tight narrative thread from capture to courtroom. Section 15 expands this into a full deposition outline with sample Q&A. The seven-sentence summary is here for quick reference:

1. "Our SDK captured the AI's decision at the moment of decision."
2. "The SDK computed an HMAC-SHA-256 over the entry, using a master key stored in our HSM."
3. "The ledger stored the entry and the MAC on persistent storage."
4. "At the end of each day, our seal job computed a Merkle root over all entries for the day, and the HSM signed that root with an Ed25519 signature."
5. "When we verify the chain, the verifier checks 12 steps including format, structure, key identity, content integrity, Merkle root, and signature."
6. "The verifier output attached to this production is the result of those 12 checks. PASS means all 12 succeeded."
7. "I have not modified the entry after capture. My evidence is the operational logs, the database audit trails, and the verifier output showing no tampering."

### 3.3 FRE 902 self-authentication

A 902(11) custodian certification or 902(13) (records generated by an electronic process) certification can self-authenticate the chain output. The institution's records-custodian affidavit names:

- The records (verifier bundle, ledger snapshot)
- Confirms regular-course-of-business generation
- Identifies the responsible custodian
- Provides the bundle's SHA-256 hash, the ledger snapshot's SHA-256, the tenant public-key fingerprint
- Is signed under penalty of perjury per FRE 902(11)

Self-authentication is the institution's first option; it removes the need for live foundation testimony at trial. The witness remains available for cross-examination.

---

## 4. FRE 1001-1004 best evidence

The best-evidence rule (FRE 1001-1004) requires the original of a writing, recording, or photograph, with FRE 1003 admitting duplicates unless authenticity is genuinely disputed.

The chain has two forms of the captured AI response:

**The original captured JSON.** The bytes the application captured at the moment of decision. This is the content-bearing form — what a human reads to understand what the AI said.

**The canonical bytes (RFC 8785 / JCS).** A deterministic re-serialization of the captured JSON. This is the integrity-bearing form — what the HMAC covers.

Both forms must be produced in discovery. The canonical bytes are the load-bearing evidence for the integrity claim because the MAC is computed over them. The original captured JSON is the load-bearing evidence for the content claim because that is what the application actually saw.

The institution's posture on a best-evidence challenge: the canonical form is computed deterministically from the original per RFC 8785, so any verifier can re-derive the canonical bytes from the original and confirm the MAC. The relationship between the two forms is reproducible and inspectable. There is no "lost original" problem; both forms are produced.

Spec §5 fixes the canonical-form algorithm. A note in spec §5 (added by a parallel agent) names the best-evidence posture explicitly.

---

## 5. FRE 803(6) hearsay and the reliability exception

The AI's output is hearsay on its face — a statement offered for the truth of what the AI asserted. Two exceptions are relevant.

### 5.1 FRE 803(6) — records of regularly conducted activity

FRE 803(6) admits records made at or near the time, by someone with knowledge, kept in the regular course of a business activity, when the source and method indicate trustworthiness.

The chain entry is a record made at or near the time (the SDK captures the entry at the moment of decision; `captured_at` records the timestamp). The chain entry is kept in the regular course of business (the institution operates the chain as part of its model-deployment infrastructure). The chain entry is made by a process the institution operates (the SDK invoked by the AI agent process).

The wrinkle: the AI did not make the record; the institution's system did. The record contains a statement *by* the AI. This is double hearsay — the chain is the record, the AI's output is the statement within the record. FRE 803(6) admits the record (the chain entry); the statement within the record (what the AI said) requires its own basis.

### 5.2 The reliability framing for AI output

The institution's argument is not that the AI's statement is true; it is that the chain proves what the AI said. The chain is admissible under FRE 803(6) as a business record. The AI's statement within the record is admissible to show what the AI said, not necessarily for the truth of what the AI asserted.

For ECOA adverse-action disputes, the customer's required disclosure includes the model's stated reasoning. The chain proves the model's stated reasoning is what was actually generated. Whether the reasoning is *correct* requires the institution's human-review record (which is itself a chain entry per the customer-dispute-procedures.md translation-step requirement).

The institution's customer-service team must understand this distinction: the chain proves the AI said X. Whether X is true or compliant is a separate evidentiary claim.

### 5.3 The residual exception (FRE 807)

Where 803(6) is contested, the residual exception (FRE 807) admits hearsay with sufficient guarantees of trustworthiness. The chain's cryptographic integrity is exactly the trustworthiness guarantee FRE 807 contemplates. The institution should be prepared to argue the residual exception in the alternative.

---

## 6. FRCP discovery — scope, format, negative evidence, multi-tenant isolation

### 6.1 Proportionality (FRCP 26(b)(1))

A subpoena reading "all AI decisions and associated metadata for plaintiff John Doe" can produce 50,000 chain entries. FRCP 26(b)(1) allows proportionality objections — discovery must be proportional to the needs of the case.

The institution's proportionality response, on the record:

> The chain authenticates the process and proves the AI's decision was faithfully recorded without tampering. That is proportional to any dispute about what the AI decided. The chain does not prove the decision was correct or compliant with policy; that requires additional evidence (the policy text, the model's training data, the decision's factual inputs). We will produce the chain entries for the affected customer for the relevant period, plus the verifier output proving the entries are unaltered, plus (if requested) the model's input context and output reasoning. We will not produce every entry for every customer; that is overbroad given the scope of damages and the relevance of the chain to authentication only.

This response gives the court a clear reasoning path: the chain is proportional to authentication, not to merits.

### 6.2 Production format (FRCP 34, Sedona Conference)

The Sedona Conference guidelines recommend native-format production for metadata-bearing electronic records. The chain is fundamentally NDJSON (one JSON object per line, per spec §6).

| Format | Verdict | Reason |
|---|---|---|
| Native NDJSON + load file | Recommended | Preserves all metadata; opposing party can re-verify |
| CSV/Excel of selected fields | Acceptable as supplement | Loses metadata; useful for deposition reference |
| TIFF/PDF of entries | Not recommended | Strips metadata; opposing party cannot re-verify; Sedona objection |
| Native NDJSON only, no load file | Acceptable | If the receiving party has technical capacity |

**Recommended production package:**

1. Chain entries in native NDJSON form for the requested `(customer_id, date_range)` scope.
2. All metadata fields per spec §4.4 (no field-stripping).
3. Load file in CSV or Excel mapping `(run_id, seq)` to the source file's line number for deposition reference.
4. Verifier output for the produced period showing PASS (or documenting any failures, with the IR record per Section 8 of this guide).
5. Tenant public key and key fingerprint for the period.
6. Records-custodian certification under FRCP 34(b) and FRE 902(11).

### 6.3 Negative evidence — proving an AI decision did NOT occur

A customer might claim "the AI never made a decision on my application; you're denying credit without giving me an AI reason." The institution's response relies on the chain's completeness claim.

This claim is admissible only if it is documented before the dispute arises. The institution's CC8.1 control description (per soc-pack) names the decision-completeness scope. Examples:

- "Every AI decision on a credit application is captured in the chain."
- "All decisions flagged as fraud-related are captured."
- "All adverse-action decisions are captured."

When the dispute arises, the IT witness testifies from this documented claim. The verifier output proves the chain's structural integrity (no entry was deleted). The completeness claim plus the structural integrity together establish that the absence of an entry for the customer's application means the AI did not make a decision on that application.

If the institution did not document the completeness claim before the dispute, it cannot rely on the absence of an entry as proof of non-occurrence. Audit-procedures P-33 (extended by a parallel agent) reinforces this requirement.

### 6.4 Multi-tenant isolation in subpoena scope

In a vendor-hosted or multi-tenant deployment, a subpoena to Tenant A should not compel production of Tenant B's chain. Each tenant's `key_fingerprint` is computed as `SHA-256(utf8(tenant_id) || ikm)[:16]` per spec §4.1, which proves each tenant's entries are integrity-bound under a separate IKM.

The IT witness's testimony for a multi-tenant proportionality dispute:

> The verifier output shows each entry's tenant_id and key_fingerprint. Tenant A's entries have one set of fingerprints; Tenant B's have a different set. Each tenant's entries are verified independently under their respective master keys. I can produce Tenant A's entries and Tenant A's verifier output. Tenant B's entries are subject to a separate access control under our vendor agreement and a separate confidentiality objection.

The cryptographic isolation argument is grounded in the spec, not in policy. That is a stronger position than a policy-only argument; opposing counsel cannot override cryptography.

---

## 7. Privilege and work-product

### 7.1 Attorney-AI consultation entries

Some chain entries may record attorney consultation — for example, "the institution's legal counsel reviewed the AI's adverse-action notice before sending it to the customer." If those entries are in the chain, they are not automatically privileged by virtue of being chain entries.

The institution must assert privilege on a per-entry basis in discovery. Entries under the `audit.legal_review.*` namespace are flagged in the institution's privilege log with the entry's `(run_id, seq)` pair.

Suggested attribute names (the institution may extend):

- `audit.legal_review.attorney_id` — bar member or internal counsel identifier
- `audit.legal_review.opinion_summary` — brief description of the legal review
- `audit.legal_review.privilege_basis` — `attorney_client` | `work_product` | `dual`

Privilege-log production for each withheld entry includes:

| Field | Source |
|---|---|
| Entry identifier | `(run_id, seq)` pair |
| Date | `captured_at` |
| Author / recipient | From the entry's audit attributes |
| Privilege basis | From `audit.legal_review.privilege_basis` |
| Subject summary | Counsel-drafted; does not reveal the privileged content |

### 7.2 Redaction and chain integrity

The MAC covers the canonical bytes of the entry. Redaction modifies the production copy, not the underlying canonical bytes; the chain's integrity verification still passes against the original canonical bytes.

The production posture: redact the privileged portion of the entry's content before producing, retain the original (un-redacted) canonical bytes in the institution's custody, produce the MAC and verifier output against the original. Provide a privilege-log summary explaining what was redacted and why.

This avoids the false dilemma "either waive privilege or break chain integrity." Neither happens. The chain's integrity is preserved because the MAC is computed over the original bytes; the privilege is preserved because the production copy is redacted.

---

## 8. FRCP 37(e) spoliation defense and adverse-inference

### 8.1 The spoliation rule

FRCP 37(e) allows sanctions for "failure to preserve electronically stored information" if the failure is willful or reckless and the loss is prejudicial. Subdivision (e)(1) authorizes curative measures; (e)(2) authorizes adverse-inference instructions and dismissal where the failure was intentional.

Two failure modes for chain evidence:

- **Scenario A.** The chain detects a violation (Scenario 1, 2, or 3 per the IR playbook). The institution must preserve evidence of the violation, not panic-repair it.
- **Scenario B.** Chain entries are missing for legitimate reasons (system downtime, network failure, application crash). The institution must document the failure as an operational event, not as destruction.

Both scenarios have an affirmative defense; both require evidence the institution operated in good faith.

### 8.2 Affirmative defense — good-faith preservation

When a chain-detected event occurs or is anticipated, the institution's preservation steps:

1. **Extend chain-data retention beyond the standard schedule.** Place a litigation hold on the affected ledger period.
2. **Freeze IKM rotation.** Or document the post-incident IKM generation with explicit before-and-after provenance, so the chain's pre- and post-incident periods are independently verifiable.
3. **Log all access to the chain data and the IKM during the incident and remediation window.** Every read, every export, every verifier run.
4. **Retain verifier output showing the incident and remediation.** The verifier output is the contemporaneous record of what was wrong and what was fixed.

These steps satisfy FRCP 37(e)(1) good-faith preservation. The institution's affirmative defense:

> We operated the chain in good faith. When a failure occurred, we documented it, investigated it, and recovered what we could. The missing entries (or detected violations) are the result of a documented system failure, not destruction of evidence. We took immediate action; we preserved the evidence of the problem and the response. There is no prejudice to the opposing party because the chain's integrity claim stands for the unaffected period and the affected period is fully documented.

### 8.3 Anti-spoliation forensic discipline for Critical scenarios

When Scenario 1 (hash mismatch), Scenario 2 (Merkle root mismatch), or Scenario 3 (signature verification failed) is detected, the institution applies forensic discipline before remediation:

| Step | Action | Why |
|---|---|---|
| 1 | Capture full memory dump of the affected verifier process if it is still running | Preserves in-memory state at the moment of detection |
| 2 | Snapshot the chain-entry file or database as it was at the moment of detection | Preserves the on-disk state — do not apply any repairs yet |
| 3 | Preserve the verifier output showing the failure | Contemporaneous record of what was detected |
| 4 | Preserve the operational event log for the affected period | Records the events the verifier ran, the alerts it generated, the IR team's response |
| 5 | Document the preservation date, time, and custodian | Standard chain-of-custody for the preserved artifacts |
| 6 | Begin remediation | Only after preservation is complete |

These artifacts are evidence of the incident for internal investigation and, if litigation arises, for the spoliation defense. Skipping step 1 or 2 to "fix the problem fast" hands the opposing party an adverse-inference argument.

### 8.4 Adverse-inference defense (FRCP 37(e)(2))

The most severe sanction is an adverse-inference jury instruction: "If the institution failed to produce evidence, you may assume that evidence would have been unfavorable to the institution." The defense requires evidence the institution operated in good faith.

The institution's evidence package:

1. Operational event documenting the failure (per soc-pack/control-evidence-events.md).
2. Incident-response record showing the institution immediately investigated and attempted recovery.
3. Root-cause explanation (hardware failure, power loss, software defect).
4. Recovery effort documentation (pulling from backup, replay from log, etc.).
5. Resolution record (entries recovered or gap documented).
6. Verifier output showing PASS for the unaffected periods bracketing the gap.

The institution's argument to the court: "Your Honor, we operated our chain in good faith. The missing entries reflect a documented system failure. We detected, investigated, recovered where possible, and documented where recovery was not possible. There is no intentional destruction; the (e)(2) standard is not met."

---

## 9. Witness qualification and evidentiary weight

### 9.1 FRE 702 expert qualification

FRE 702 requires the expert to have specialized knowledge that helps the jury understand the evidence. A verifier operator must be qualified to (a) explain the verifier's 12-step procedure, (b) interpret the failure-reason strings, and (c) connect the verifier output to the institutional controls.

Required training for the verifier operator:

- Spec §7 verifier procedure — every step, what it tests, what failure means.
- Failure-reason taxonomy — the meaning of each failure-reason string.
- Institutional control mapping — why step 8 (fingerprint mismatch) means "check the IKM roster," not "the cryptography is broken."
- IT-examiner training (per regulator-pack/examiner-training.md) is the institution's internal training baseline.

The institution documents training completion as part of the operator's personnel records. An IT witness who has completed that training and can articulate the verifier's logic under cross-examination is qualified per FRE 702.

The opposing party will challenge: "You're a database administrator, not a cryptographer; how do you know what the verifier's output means?"

The witness's answer: "I am qualified to operate the verifier and interpret its output because I have completed [name the training], I understand each of the 12 steps the verifier performs, I can explain what each failure-reason means, and I can map each failure to the institutional control that addresses it. I do not claim to be a cryptographer. I claim to be a qualified operator of the institution's verification tool, and I rely on the verifier's output as the institution's IT witness."

### 9.2 Witness-verifier mode evidentiary weight

Spec §7 allows a "witness verifier" run without the IKM (`--master-key` omitted), which produces `PASS-STRUCTURALLY, key-bound verification skipped`. This is a lower confidence than full verification. The chain's structure proves no entries were inserted or deleted; the per-entry MAC was not checked.

The witness-mode result is admissible under FRE 901(b)(9). The expert must explicitly tell the jury what was tested and what was not:

> The chain's structure was verified, but I could not verify the individual entry integrity (MAC) because I did not have access to the master key. This means I can prove no entry was added or removed, but I cannot independently verify the entry's content was not modified within the chain — that would require the master key. The institution's IT staff certifies they have the master key and verified the full chain; my structural verification is consistent with the structural evidence I have.

This framing is more credible than a hand-wave. The jury understands the limitation without thinking the chain is weak. The structural-only result is exactly what the chain is designed to support for adverse-verification scenarios where the institution cannot disclose the IKM.

### 9.3 Sample voir dire Q&A

The opposing party may seek voir dire on the IT witness's qualification. Prepared answers:

```
Q: What is your role at the institution?
A: I am [title], responsible for [scope]. My direct involvement with the
   chain-of-custody system includes [scope].

Q: What training have you received on the chain-of-custody system?
A: I completed the IT-examiner training documented in our regulator-pack.
   The training covers the spec's 12-step verification procedure, the
   failure-reason taxonomy, and the institutional controls that map to
   each verification step.

Q: Are you a cryptographer?
A: No. I am a qualified operator of the institution's verification tool.
   The cryptographic primitives are NIST-standardized and are documented
   in the spec and in published NIST publications. My role is to operate
   the verifier and interpret its output; the cryptography is the
   subject of the spec and the test-vector corpus.

Q: How do you know the verifier itself was not tampered with?
A: The verifier binary's SHA-256 hash is recorded in the institution's
   build attestation. Each verifier run produces a verifier.run_completed
   operational event recording the verifier version, the binary hash,
   the operator's identity, and the host. Those records are themselves
   audit trails for the verifier's own custody.
```

---

## 10. Operational litigation support

### 10.1 Litigation hold

When the institution receives notice of litigation or anticipates a subpoena, the IT team triggers a litigation hold on the chain data:

| Step | Action |
|---|---|
| 1 | Extend chain-data retention beyond the standard schedule for the affected scope |
| 2 | Freeze IKM rotation for the affected ledger period (or document the rotation with explicit before-and-after IKM custody records) |
| 3 | Prevent any deletion or modification of chain entries or seal records for the affected scope |
| 4 | Preserve verifier output for the affected period |
| 5 | Notify the operations team that backup retention for the affected period is extended |
| 6 | Document the litigation-hold in the institution's control-evidence repository |

The hold runs until counsel releases it.

### 10.2 Subpoena response timing

FRCP 34 sets a 14-day response window (extendable by court order). The institution's discovery timeline:

| Day | Action |
|---|---|
| 0 | Subpoena served. Counsel reviews scope. Litigation hold issued. |
| 1-3 | IT team identifies the relevant `(tenant_id, run_id)` set using the customer-correlation index. |
| 3-7 | IT team retrieves chain entries for the affected scope. Verifier runs across the affected period. |
| 7-10 | Counsel reviews for privilege; redactions applied where needed. |
| 10-12 | Records-custodian certification drafted. Production package assembled per Section 6.2. |
| 12-14 | Production delivered. Verifier output included in the package. |

Plan for expedited verification during the discovery window. A multi-month verification is not credible at a deposition; the institution should be able to verify a year of entries in hours, not days.

### 10.3 Chain-of-custody for the chain itself — `verifier.run_completed`

A predictable Daubert challenge: "How do we know the verifier output you're relying on is itself authentic and unaltered?" The answer is the meta-evidence: a record of every verifier invocation.

The institution emits a `verifier.run_completed` operational event for each verifier run. Suggested fields:

| Field | Purpose |
|---|---|
| `run_at` | RFC 3339 UTC timestamp of the run |
| `verifier_version` | The verifier software version |
| `verifier_binary_sha256` | Cryptographic hash of the verifier executable for reproducible-build verification |
| `command_line_args` | Flags used (`--strict`, `--posture`, `--master-key` or `--witness`) |
| `run_by` | The operator's identity |
| `run_on_host` | The hostname confirming the run was on institutional infrastructure |
| `result` | PASS / FAIL / PASS-STRUCTURALLY |
| `affected_tenant_ids` | Which tenants were verified |
| `coverage_period` | Start and end of the ledger window verified |

This event is the chain's own audit trail. The IT witness testifies from these events when cross-examined on the verifier's custody.

The soc-pack/control-evidence-events.md (extended by a parallel agent) defines the event schema; this document names its evidentiary purpose.

### 10.4 Forensic handoff procedure

When the institution hands off chain evidence to an examiner — law enforcement, defense counsel, neutral third party — the handoff requires forensic discipline. The handoff documentation becomes the foundation for the examiner's testimony that the chain was received in integrity state and remained unaltered throughout the examination.

**Physical transfer of master-key material (rare; only when adversarial verification with the IKM is required):**

1. Tamper-evident envelope for the IKM material.
2. Chain-of-custody form signed by both parties naming: the IKM identifier, SHA-256 hash of the key bytes, date, time, signatures of both parties.
3. Two-person rule on the institution's side — the transferring custodian and a witness.
4. Receipt confirmation — the recipient signs that they received the sealed evidence in integrity state.

**Digital transfer (HSM-mediated verification without IKM disclosure):**

1. The institution provides the examiner with an HSM API endpoint and documented credentials.
2. The examiner's verification tool queries the HSM for signing operations without ever holding the IKM bytes.
3. The institution logs every HSM call and provides the log to the examiner as part of the handoff package.
4. Two-person rule applies on the institution's side for credential issuance.

Either method is acceptable. The institution documents which method was used and why. The handoff documentation includes:

- Date, time, custodians (institution side, recipient side)
- Method (physical or digital)
- Materials transferred (IKM bytes, HSM access, ledger snapshot, verifier binary)
- Hashes of all transferred materials
- Receipt confirmation
- Examiner's intended scope of examination

The operator-guide (extended by a parallel agent) covers the operational procedure; this document covers the evidentiary framing.

---

## 11. Cross-border evidence

A global institution faces a layered problem: U.S. discovery may compel production of evidence stored under EU data-protection law, UK procedural law, or Canadian rules of evidence. Each jurisdiction has its own admissibility standard. The institution's posture must be planned in advance.

### 11.1 GDPR + Schrems II

If chain entries are stored on EU servers and contain customer personal data (name, account ID, address), U.S. discovery production may conflict with GDPR Articles 44-50 (cross-border data transfer) and the Schrems II decision (Court of Justice of the European Union, *Data Protection Commissioner v. Facebook Ireland*, 2020) restricting transfers to countries without adequacy.

Institution posture:

| Step | Action |
|---|---|
| 1 | Flag PII-bearing entries with an `audit.pii_category` field (`name` \| `account_id` \| `address` \| etc.) so the institution can segregate during discovery |
| 2 | When U.S. court compels production, counsel seeks a protective order limiting distribution to the court and the parties' counsel before producing PII-bearing entries |
| 3 | For entries on EU servers, counsel may pursue MLAT or letter rogatory to compel EU custodial testimony before producing — alternative to direct U.S. discovery |
| 4 | Counsel may invoke GDPR Article 49(1)(e) (necessary for legal claims) as the lawful basis for the transfer, where applicable |

The institution's control description names its chosen posture before any dispute arises. Reactive Schrems II decisions in the middle of discovery are expensive.

### 11.2 UK CPR Part 31

The UK Civil Procedure Rules Part 31 governs disclosure (the UK equivalent of U.S. discovery). The CPR Part 31 standard for electronic disclosure requires authentic, reliable records.

The chain satisfies the standard by demonstrating:

- **Contemporaneous capture** per spec §4.1 — entries are captured at the moment of decision.
- **Integrity verification** per spec §7 — the 12-step verifier procedure is documented and reproducible.
- **HSM-backed signing** per spec §4.3 — the seal record's signature is HSM-protected.
- **Audit trails** per soc-pack operational events — the chain's operation is itself recorded.

UK courts recognize cryptographic chains under the same evidential principles as U.S. courts. The institution's UK counsel relies on the same foundation testimony with UK-specific procedural framing.

### 11.3 EU eIDAS

The EU eIDAS Regulation (910/2014) defines qualified electronic signatures (QES) and qualified trust services (QTS). The chain's Ed25519 signatures meet the QES standard if the HSM has FIPS 140-2 Level 3 (or eIDAS-equivalent) certification.

For EU jurisdictions, the institution documents:

- HSM model and certification level (FIPS 140-2 Level 3 or higher; eIDAS-recognized)
- Signing-key generation procedures
- Seal-record audit trail per spec §4.3

The institution may also choose to engage a Qualified Trust Service Provider (QTSP) for additional eIDAS qualification — a v1.x extension worth tracking.

### 11.4 Canadian standard — *R. v. Hape*

The Supreme Court of Canada's decision in *R. v. Hape*, 2007 SCC 26, sets the standard for foreign electronic evidence. The three-factor reliability test:

1. Are there reasonable grounds to believe the evidence is reliable?
2. Is there a principled reason to admit it?
3. Would reception undermine the administration of justice?

The chain satisfies factor 1 through its cryptographic integrity claim and the verifier output. Factor 2 is satisfied where the evidence is necessary to the dispute. Factor 3 is supported by the open-source reference implementation and the public test-vector corpus — the technical basis for the chain's reliability is open to scrutiny, which Canadian courts treat as a strong indicator of reliability.

The institution's Canadian counsel points the court to the reference implementation (Apache 2.0) and the test-vector corpus (publicly available) as evidence of open technical scrutiny.

### 11.5 Quick reference

| Jurisdiction | Standard | Institution's evidence |
|---|---|---|
| U.S. federal | Daubert / FRE 901 | Spec §7 procedure, test vectors, foundation testimony, verifier output |
| UK | CPR Part 31 | Same plus contemporaneous-capture documentation |
| EU | eIDAS QES/QTS | Same plus HSM FIPS/eIDAS certification |
| Canada | *R. v. Hape* | Same plus reference to open-source reference implementation |

---

## 12. Trusted timestamps — RFC 3161

The chain's timestamps come from the institution's system clock disciplined by NTP (per spec §10.5 / audit-procedures P-7). For most disputes, NTP discipline is sufficient — the institution testifies to its clock-synchronization posture as the foundation for timestamp reliability.

For high-stakes disputes, the timestamp is challenged: "How do we know the timestamp is accurate? The institution could have set its system clock back and backdated the entry."

RFC 3161 (the Time-Stamp Protocol) provides an independent attestation. A trusted third-party Time Stamping Authority (TSA) issues a TimeStampToken — a digitally signed attestation that a particular hash existed at a particular time. Integrating RFC 3161 into the chain creates an external time anchor that does not depend on the institution's clock.

### 12.1 v1.0 baseline — NTP

Spec v1.0 does not require RFC 3161. The institution's timestamp foundation is:

- NTP synchronization to a stratum-2 or higher source (per spec §10.5)
- Audit-procedures P-7 documents the NTP discipline
- The institution's IT witness testifies to NTP-source configuration, drift monitoring, and the operational events that record any clock anomalies

This is sufficient for the typical dispute. Most courts accept NTP discipline as reliable for timestamps in business records.

### 12.2 v1.x optional — RFC 3161 token

For institutions anticipating maximum-stakes litigation (large-class adverse-action disputes, criminal investigations, regulatory enforcement actions with sanctions), RFC 3161 is a recommended but not required upgrade. A future v1.x extension may define:

- `audit.timestamp.rfc3161_token` — Base64-encoded RFC 3161 TimeStampToken for the `captured_at` value
- `audit.timestamp.tsa_url` — the TSA endpoint that issued the token
- `audit.timestamp.tsa_cert_chain_sha256` — hash of the TSA's certificate chain at the time of issuance

The verifier validates the token alongside the entry's MAC. Failure of the token validation does not invalidate the chain (the chain's integrity is independent), but successful token validation gives the timestamp an external anchor.

### 12.3 Daubert framing for timestamp credibility

The institution's expert testifies on timestamps using a layered argument:

1. The system clock is disciplined by NTP per spec §10.5 and audit-procedures P-7.
2. The institution monitors NTP drift through documented operational events.
3. (If RFC 3161 is in use) The independent TSA token provides an external time anchor that does not rely on the institution's clock.
4. Backdating an entry would require simultaneous compromise of the institution's NTP discipline, the chain's MAC, and (where applicable) a forged TSA token.

The expert's bottom line: the timestamp is contemporaneous within the documented NTP discipline; the chain does not allow silent backdating.

---

## 13. Plain-language explainability for jurors

A jury does not understand HMAC-SHA-256 or Merkle trees. The institution's expert must translate. The chain's design is unusually juror-friendly because every cryptographic primitive has a physical-world analogy.

### 13.1 The four core analogies

**Analogy 1 — The record card with a fingerprint.**
"Every time the AI makes a decision, the institution records it on a card. The card has the AI's decision, the time it happened, and a special fingerprint computed from the card's contents. If anyone changes the card later, the fingerprint won't match. The fingerprint is computed by a one-way mathematical formula that is easy to compute one way but practically impossible to reverse — like a recipe that's easy to follow but hard to deduce from the finished cake."

**Analogy 2 — The pyramid of fingerprints (Merkle tree).**
"At the end of each day, the institution takes all the cards from that day and builds a pyramid of fingerprints. At the bottom, the cards' fingerprints. Above that, fingerprints made by combining pairs of fingerprints. At the top, a single fingerprint that depends on every card below it. Change any card, anywhere in the pyramid, and the top fingerprint changes."

**Analogy 3 — The locked box (HSM).**
"The institution has a special locked box — a piece of hardware called a Hardware Security Module — that holds a cryptographic key. The locked box is designed so that nobody, not even the institution, can take the key out of the box. At the end of each day, the institution shows the top fingerprint of the day's pyramid to the locked box, and the box signs it. The signature proves the day's pyramid was real and was signed at the end of the day. If anyone tries to forge a day's pyramid, they can't reproduce the signature without the key inside the locked box, which they cannot get."

**Analogy 4 — The tamper detector.**
"If anyone tampers with a single card, three things break at once: the card's own fingerprint won't match the card. The pyramid's top fingerprint won't match the day's pyramid. The locked box's signature won't match the top fingerprint. Three independent tamper detectors, each working separately. To hide tampering, an attacker would have to break all three at once — which would mean simultaneously compromising the institution's records, the institution's fingerprinting computer, and the institution's locked box. The institution operates these as three separate systems with three separate protections."

### 13.2 Sample demonstrative exhibits

A demonstrative exhibit is a visual aid the jury sees but does not take into deliberation. Three exhibits suffice:

**Exhibit A — A single chain entry as a record card.**
A card-shaped diagram showing:
- Top: "AI Decision Card"
- Middle row: "What the AI said: [text of the decision]"
- Middle row: "When: [timestamp]"
- Bottom: "Fingerprint: a1b2c3d4..."

**Exhibit B — The day's pyramid (Merkle tree).**
A pyramid-shaped diagram with:
- Bottom row: 8 small boxes labeled "Card 1" through "Card 8"
- Middle row: 4 boxes, each labeled "Combined fingerprint of [pair below]"
- Above that: 2 boxes, each labeled "Combined fingerprint of [pair below]"
- Top: 1 box labeled "Top fingerprint of the day"

**Exhibit C — The signing process.**
A flow diagram:
- Start: "Top fingerprint of the day"
- Arrow to: "Institution's locked box (HSM)"
- Arrow out: "Signed top fingerprint"
- Caption: "The locked box signs the top fingerprint. The key never leaves the box."

If the spec ships SVG-friendly diagrams under `docs/assets/images/svg/`, the institution may use those directly as exhibits with appropriate captions. If not, simple hand-drawn or computer-drawn versions are sufficient — the analogies are simpler than any actual diagram.

### 13.3 The narrative the jury hears

The expert's closing summary in plain language:

> The institution has a system that, every time the AI makes a decision, takes a fingerprint of the decision and stores it. Every day, it builds a pyramid of fingerprints and signs the top. If anyone later tries to change what the AI said, the fingerprints won't match — three different ways, all at once. The institution runs a verification tool that checks the fingerprints. The tool's output for this case is what you have in front of you. The tool says PASS, which means none of the fingerprints were broken. The chain proves what the AI said and proves no one changed the record. It does not prove the AI was right; that's a separate question.

---

## 14. Epistemic scope — what the chain proves

The chain's most important boundary is also the most easily misunderstood. The institution's expert and counsel must keep this distinction crisp at every step of the case.

### 14.1 What the chain proves

- **What the AI said.** The exact content of the AI's output at the moment of the decision.
- **When the AI said it.** The contemporaneous timestamp, disciplined by NTP and (optionally) RFC 3161.
- **That no one tampered with the record.** The cryptographic integrity layer detects any post-hoc modification.

That is it. Three claims, all testable, all defensible at Daubert.

### 14.2 What the chain does NOT prove

- **That the AI's statement is factually accurate.** The chain captures hallucinations the same way it captures correct statements.
- **That the AI's statement is policy-compliant.** Compliance is a separate review, not an integrity check.
- **That the AI's statement is unbiased.** Bias testing requires statistical analysis on populations of decisions, not integrity verification on individual records.
- **That the AI was the right tool for the decision.** Suitability is a model-risk question, not a chain-of-custody question.

### 14.3 Why the distinction matters in litigation

A juror who hears "the chain proves what the AI said" can drift into "the chain proves what the AI said is true." The expert's job is to prevent that drift.

The expert's framing line:

> The chain proves our system said X about your account. Whether X is true is a separate question, and we answer it with separate evidence — the policy text, the model's input data, the institution's human-review record, the regulator's audit findings.

For ECOA adverse-action disputes, this distinction translates directly to the customer's required disclosure. The chain proves the model's stated reasoning. The institution's separate review (also chain-recorded per the customer-dispute-procedures.md translation requirement) proves whether the reasoning was factually sound and ECOA-compliant.

For AI-hallucination disputes, the distinction is the institution's defense: "Yes, the AI made a statement that turned out to be inaccurate. The chain proves the institution's record of what the AI said is unaltered. Whether the AI's statement was accurate is a model-risk question; the institution's customer-dispute procedure addresses that question through its hallucination cross-check (per customer-dispute-procedures.md). The chain is one layer of evidence; the model-risk record is another."

Spec §1 (extended by a parallel agent) names this scope explicitly. This document operationalizes the framing for litigation use.

---

## 15. Foundation testimony — sample direct examination

This section is the deposition outline for the institution's IT witness. It is structured as a sample direct-examination Q&A, with notes on what counsel should be prepared to follow up on.

### 15.1 Witness preparation checklist

Before the deposition or trial:

- [ ] Witness has completed IT-examiner training documented in the regulator-pack.
- [ ] Witness has reviewed the chain-of-custody spec, especially §4 (primitives), §5 (wire format), and §7 (verification).
- [ ] Witness has run the verifier on the affected period and reviewed the output.
- [ ] Witness has reviewed the `verifier.run_completed` events for the affected period.
- [ ] Witness has reviewed the change-management records for the affected period.
- [ ] Witness understands the four core juror analogies (Section 13.1) for cross-examination "explain it to a 12-year-old" questions.
- [ ] Witness has reviewed Section 14 (epistemic scope) and can recite the framing line under pressure.

### 15.2 Sample direct examination

```
Q: Please state your name and current position.
A: [Name], [title] at [institution].

Q: How long have you held that position?
A: [Duration].

Q: What are your responsibilities?
A: I am responsible for [scope, including chain-of-custody system operation].

Q: Are you familiar with the institution's chain-of-custody system for AI decisions?
A: Yes.

Q: How did you become familiar with it?
A: I completed the institution's IT-examiner training, which is documented
   in our regulator-pack. I have operated the verifier as part of my
   responsibilities for [duration]. I have reviewed the FFIEC chain-of-
   custody specification, particularly the sections on the four
   cryptographic primitives, the wire format, and the verification procedure.

Q: Please describe, in plain language, what the chain-of-custody system does.
A: The system records every AI decision the institution makes. At the
   moment of decision, our SDK captures the AI's output and computes a
   cryptographic fingerprint over the entry, using a master key stored in
   our Hardware Security Module — our HSM, which is a tamper-resistant
   piece of hardware. The entry and the fingerprint are stored in our
   ledger. At the end of each day, our seal job builds a tree of
   fingerprints over all entries for the day, and the HSM signs the
   top of the tree. The signature proves that the day's entries are
   complete and unaltered.

Q: How does the institution verify that the chain has not been tampered with?
A: We run a verifier tool that performs 12 checks against the ledger.
   The 12 checks are documented in spec §7. They include format validation,
   structural validation, key-identity validation, content integrity using
   the per-entry fingerprint, Merkle-root validation, and signature validation
   against the HSM's public key. If all 12 checks pass, the verifier emits
   PASS. If any check fails, the verifier emits a specific failure reason.

Q: Have you run the verifier on the chain entries in this case?
A: Yes. I ran the verifier on [date] covering the period [date range].
   The verifier output is exhibit [exhibit number]. The result was PASS.

Q: What does PASS mean for purposes of this case?
A: PASS means that for the entries in the affected period, the chain's
   integrity is intact. No entry was inserted, deleted, or modified after
   capture. The fingerprints, the Merkle root, and the HSM signature are
   all consistent with the original entries.

Q: How do we know your verifier itself was not tampered with?
A: Each verifier run produces a verifier.run_completed operational event
   that records the verifier's version, the SHA-256 hash of the verifier
   binary, my identity, the host the verifier ran on, and the result.
   The institution's build attestation records the verifier binary's
   hash for each release. The hash in the operational event matches the
   hash in the build attestation. That cross-check confirms the verifier
   binary I ran is the institution's released verifier.

Q: Mr./Ms. [witness], does the chain prove that the AI's decision in this
   case was correct?
A: No. The chain proves what the AI said and proves no one tampered with
   the record. Whether the AI's decision was correct is a separate question
   that requires separate evidence — the policy that governed the decision,
   the data that fed the decision, and the institution's human-review record.

Q: Does the chain prove that the AI's decision was unbiased?
A: No, for the same reason. The chain is an integrity record, not a
   correctness or fairness record.

Q: Then what is the chain's purpose in this case?
A: The chain is the institution's contemporaneous record of what the AI
   said and when it said it. Its evidentiary value is to prove that the
   record we are producing today is the same record that existed at the
   moment of the decision. Without the chain, the institution would have
   no way to prove that to the court or to the customer.
```

### 15.3 Anticipated cross-examination

The opposing party's lawyer will probe. Prepared answers:

```
Q: You're not a cryptographer, are you?
A: Correct. I am a qualified operator of the institution's verification
   tool. The cryptographic primitives are NIST-standardized and documented
   in the spec. The cryptography is the subject of the spec and the
   test-vector corpus, both of which are public.

Q: How do you know SHA-256 hasn't been broken?
A: SHA-256 is current per NIST guidance. NIST monitors cryptographic
   primitives and issues transition advisories when they weaken. No such
   advisory has issued for SHA-256. Until NIST issues guidance otherwise,
   the institution relies on the algorithm's integrity per NIST's
   current risk assessment.

Q: Couldn't the institution have backdated this entry?
A: The chain's design prevents silent backdating. Every entry is captured
   at the moment of decision; the timestamp is recorded under the integrity
   layer; the institution's clock is disciplined by NTP per audit-procedure
   P-7; backdating an entry would require simultaneous compromise of the
   master key, the HSM, the ledger storage, and the NTP discipline. The
   institution's IR records show no such compromise.

Q: What if the verifier output you're showing me is fake?
A: The verifier.run_completed operational event records every verifier
   run with the verifier binary's hash, the host, the operator, and the
   result. The output I'm showing you is consistent with that operational
   record. Both records are themselves part of the chain or the operational
   event store, with their own integrity protections.

Q: What if the AI's decision was wrong?
A: The chain proves the AI's decision was faithfully recorded; it does
   not prove the decision was correct. The institution's separate review
   procedures address correctness. The chain is the integrity layer; the
   review procedure is the correctness layer.
```

### 15.4 Refusal-to-speculate phrases

The witness should not extrapolate beyond competence. Prepared refusal phrases:

- "I am not qualified to opine on [matter outside the verifier and chain operation]; that question is for [the appropriate expert]."
- "The chain does not address that question. The institution's [other procedure / record] addresses it."
- "I do not have personal knowledge of [matter]; I would refer you to [the appropriate witness or record]."

Refusing to speculate is a sign of credibility, not weakness. A witness who answers every question — including questions outside their competence — invites cross-examination on the speculative answers.

---

## 16. Cross-references

This document is the litigation-side framing. The substantive material is distributed across the spec and the operational documents.

| Topic | Where the substance lives |
|---|---|
| Daubert grounding text in the spec | spec §1 (added by parallel agent) |
| Best-evidence note on canonical form | spec §5 (added by parallel agent) |
| Evidentiary artifacts subsection | spec §10 (added by parallel agent) |
| RFC 3161 optional extension | spec §10 (added by parallel agent) |
| Forensic preservation for IR scenarios | incident-response-playbook.md (extended by parallel agent) |
| P-33 completeness claim | audit-procedures.md (extended by parallel agent) |
| Multi-tenant isolation procedure | audit-procedures.md (extended by parallel agent) |
| Verifier exit codes | operator-guide.md (extended by parallel agent) |
| Forensic handoff procedure | operator-guide.md (extended by parallel agent) |
| Hearsay framing for ECOA | customer-dispute-procedures.md (extended by parallel agent) |
| Verifier-operator qualification | examiner-training.md (extended by parallel agent) |
| `verifier.run_completed` event schema | soc-pack/control-evidence-events.md (extended by parallel agent) |
| Court-ordered key disclosure | legal-disclosure.md (existing) |

The cross-references resolve once the parallel agents' edits land. Section numbers are stable; specific paragraph references inside sections may shift as the parallel agents fold their material in.

---

## 17. Operational checklist

When litigation is anticipated:

- [ ] Counsel issues litigation hold (Section 10.1)
- [ ] IT extends chain-data retention for the affected scope
- [ ] IKM rotation frozen or documented with before-and-after provenance
- [ ] Verifier run on the affected period; output preserved
- [ ] `verifier.run_completed` operational events preserved
- [ ] Witness identified and prepared per Section 15.1 checklist
- [ ] Production package assembled per Section 6.2
- [ ] Privilege screen completed (Section 7); privilege log drafted
- [ ] Cross-border posture confirmed if applicable (Section 11)
- [ ] Records-custodian certification drafted under FRE 902(11)
- [ ] Production delivered within FRCP 34 timeline
- [ ] Disclosure recorded in the institution's control-evidence repository

When a chain-detected event occurs and litigation is possible:

- [ ] Forensic preservation steps complete before any remediation (Section 8.3)
- [ ] Memory dump captured if verifier was running
- [ ] Ledger snapshot preserved as-detected
- [ ] Verifier output preserved
- [ ] Operational event log preserved
- [ ] IR record documents detection, root cause, recovery, resolution
- [ ] Affirmative-defense package assembled per Section 8.2
- [ ] Counsel notified
- [ ] Litigation hold issued if litigation is reasonably anticipated

When a subpoena arrives:

- [ ] Counsel reviews scope; proportionality objection prepared if applicable (Section 6.1)
- [ ] Records-custodian engaged
- [ ] Affected `(tenant_id, run_id)` set identified via customer-correlation index
- [ ] Multi-tenant isolation argument prepared if applicable (Section 6.4)
- [ ] Production format follows Sedona / FRCP 34 (Section 6.2)
- [ ] Witness preparation completed per Section 15
- [ ] If cross-border: protective order pursued before producing (Section 11.1)
- [ ] If master-key compelled: per legal-disclosure.md procedure
- [ ] Production delivered with verifier output and custodian certification

---

## 18. Adversarial-balance — anticipating opposing counsel

The chain's cryptography passes Daubert. The institutional procedures around it are where opposing counsel has angles. This section names the four most-pressed adversarial angles, the institution's documented response posture for each, and the operational discipline that closes each angle before the cross-examination starts. The section is written for the institution's litigation-support team to read alongside Sections 2 and 15; counsel preparing for adversarial litigation reviews this section together with the FRE 902 certification's adversarial-balance addendum (`templates/fre-902-certification.md` Addendum A1-A5). *(closes Hartwell G-3, G-4, G-5, G-7, G-8)*

### 18.1 The partial-disclosure-as-cherry-pick angle

Opposing counsel will press: "The defendant produced a partial-disclosure bundle naming the entries the defendant chose to disclose. The cryptography authenticates inclusion of those entries in the seal-day's Merkle root. It does not assert that the defendant disclosed every responsive entry. The defendant's `--strict` basis field is procedural, not cryptographic. The defendant controls what is produced AND the count of what existed. A 100 KB full-day seal is reduced to a 15 KB partial-disclosure bundle without any cryptographic discipline naming what was withheld." *(closes Hartwell G-3)*

The angle is honest about a real property of the partial-disclosure mode — the mode is scope-narrowing by design and does not assert completeness. The selective-production document at `selective-production-and-sampling.md` §6 acknowledges this verbatim ("the consumer plaintiff assesses completeness through the order's specificity and through the receiver's compliance certification, not through the cryptography").

The institution's response posture has three elements.

**Element 1 — Total-count disclosure on the record.** The institution's partial-disclosure bundle is accompanied by a sworn statement of (a) the total entry count for the seal-day, (b) the count produced, and (c) the proportion withheld. The statement is on the record before the bundle is offered into evidence. Cross-reference to `selective-production-and-sampling.md` §10 (the institution's response procedure when partial-disclosure is challenged). The total-count disclosure is integrity-bound — the seal record's count of entries is a property the seal commits to via the Merkle root construction, even though the count is not directly readable from the audit path alone. The institution can produce the count from its ledger and stipulate to it.

**Element 2 — Disclosure-basis specificity.** The disclosure manifest's `basis` field is upgraded from procedural to load-bearing in adversarial contexts. The institution's manifest names the specific search criteria the institution ran against the ledger to identify responsive entries (the customer identifier, the date range, the run_id range, the topic-classification filter). Opposing counsel can examine the basis criteria, identify entries that should have matched the criteria but were not produced, and challenge the criteria as too narrow.

**Element 3 — Full-day seal output as separate evidence.** When opposing counsel disputes completeness, the institution offers the full-day seal output for the seal-day(s) covering the partial-disclosure bundle. The full-day seal does not reveal the bodies of unrelated entries (the verifier output names entry counts and seal metadata, not entry bodies), but it establishes the day's total entry count and the institution's claimed completeness. Opposing counsel can compare the partial-disclosure bundle's count against the full-day seal's count and the basis criteria; any divergence is a documented question.

The combined effect of the three elements: opposing counsel pressing the cherry-pick angle receives a partial-disclosure bundle plus a total-count statement plus a basis-specificity declaration plus a full-day seal output. The institution's production has not changed — the responsive entries are still the entries the institution chose to produce — but the procedural record now names the institution's choices explicitly. Opposing counsel arguing cherry-picking has to argue against documented criteria rather than against unstated ones.

### 18.2 The FRCP 37(e) policy-driven-exclusion angle

Opposing counsel will press: "The institution's data-classification policy at the time of capture excluded substantive prompt and response content from the chain. The chain is integrity-bound for what it captured — but what it captured is the SHA-256 hash, not the substantive content. The institution chose this exclusion under a policy it controlled, during a period when class-member harms were foreseeable. Under FRCP 37(e)(2), policy-driven exclusion of substantive content from a record-keeping system otherwise designed to preserve evidence is a candidate for an adverse-inference instruction." *(closes Hartwell G-4)*

The angle is the most viable spoliation theory against an institution operating hash-only retention. The institutional response posture has four elements.

**Element 1 — Data-classification policy preservation under RS-7.** The institution's data-classification policy at the time of capture is itself preserved as a record series under RS-7 disposition documentation (`templates/records-management-program.md` §3.1). When opposing counsel seeks the data-classification policy in discovery, the institution produces it from the records-management archive. The policy, the change history during the period at issue, and the senior-officer approval records are all preserved in integrity-bound form (the institution's governance archive integrity-binds the records-management documents under the same disciplines as the chain itself).

**Element 2 — Documented contemporaneous business rationale.** The institution's policy decision to capture hash-only rather than full content was made under a documented business rationale — typically privacy minimization under GLBA Title V, GDPR Article 5(1)(c) data minimization, or the institution's own NPI-minimization posture. The rationale is part of the policy's preservation record. When opposing counsel argues the policy was driven by litigation-anticipation, the institution rebuts with the contemporaneous rationale that names privacy minimization as the driver.

**Element 3 — Substantive content provenance from upstream sources.** The chain captures the hash; the substantive content existed at the LLM provider during the institution's call to the provider's API. The institution's vendor contract under §500.11 / DORA Article 30 (cross-reference `dora-articulation-overlay.md` §3) names the LLM provider's retention discipline. When opposing counsel seeks the substantive content, the institution coordinates with the LLM provider to determine whether the provider retained the content and whether the provider's retention is subject to compelled disclosure. The institution does not control the LLM provider's retention; the chain's hash-only posture does not preclude the provider's retention from being a discoverable source.

**Element 4 — Foreseeable-harms framing rebuttal.** The plaintiff's argument under FRCP 37(e)(2) requires showing intent to deprive. Foreseeable-harms-plus-policy-choice is a viable factual basis but is not the same as intent. The institution's response: the data-classification policy was institution-wide, not litigation-targeted; the policy applies to all chain entries, not just to entries that might appear in litigation; the policy's purpose was privacy minimization, not evidence destruction; the policy's effect on later litigation was an incidental consequence of a privacy decision, not the driver of the decision. The policy preservation record (Element 1) and the contemporaneous rationale (Element 2) establish the rebuttal.

The combined effect: opposing counsel's adverse-inference theory rests on a four-element factual base the institution's policy preservation, contemporaneous rationale, vendor-side content provenance, and foreseeable-harms framing rebuttal address. The bench can credit the rebuttal or not; the institution's posture is documented and consistent. The institution's witness deposing on this question has the policy in front of them, the rationale memo, the vendor-relationship paperwork, and the privacy-minimization framing — all preserved long before the litigation arose.

### 18.3 The hash-only certification scope-limitation angle

Opposing counsel will press: "The §C hash-only certification authenticates that the institution received bytes whose hash is recorded. It explicitly disclaims attestation of factual accuracy, policy compliance, and freedom from bias. Read as a defense instrument, this is precisely scoped — which is what makes it dangerous. A class-action plaintiff is left arguing that the defendant's AI made discriminatory adverse-action decisions with no record of what the AI actually said — only a hash that the defendant claims matches bytes the defendant claims it received." *(closes Hartwell G-5)*

The angle is a Daubert-flavored 702 challenge against the certification's evidentiary value rather than against the cryptography itself. The institution's response posture has three elements.

**Element 1 — Hash-only certification's true scope.** The §C certification is a hash-only certification by design. It authenticates bytes-as-received; it does not authenticate substantive AI behavior. The institution's witness on the stand acknowledges this scope explicitly. The §C certification is offered to authenticate the bytes, not to authenticate what the AI said. The institution's case on what the AI said depends on either the substantive content captured under §A or §B (where the institution's data-classification policy admitted full content capture), or on the LLM provider's contemporaneous logs (where the institution's policy was hash-only).

**Element 2 — Disclosure of substantive content in alternate systems.** The institution discloses to the requesting party whether substantive content was retained under a different system: the LLM provider's logs (subject to vendor-side retention policy), downstream systems that consumed the AI output (institution-side internal systems), or regulator-side examination archives (where the substantive content was reviewed during examination). The institution's records officer maintains the documentation of which alternate systems hold substantive content; the records officer's annual report (cross-reference `templates/records-management-program.md` §10) names which alternate systems are operative.

**Element 3 — Section 1.2 epistemic-scope discipline.** The chain's §1.2 epistemic-scope statement names exactly what the chain proves and what it does not. Spec §1.2: "The chain proves what was said and that the record was not tampered with after capture; it does not prove the substantive correctness of the captured content." The institution's witness cites §1.2 as the institution's own honest framing — the institution does not claim the chain proves substantive correctness, and the §C certification operates within that scope. Opposing counsel pressing the §C certification beyond its scope is pressing against the institution's own scope-limitation, not against an institutional overreach.

The combined effect: opposing counsel's argument that the §C certification is scope-limited and therefore zero-value is met with the institution's documented framing that the §C certification IS scope-limited and operates within that scope. The institution's case on what the AI said rests on alternate evidence (vendor logs, internal systems, examination archives) — not on the §C certification. Opposing counsel's 702 challenge against using the §C certification as substantive-content evidence is met with the institution's stipulation that the §C certification is not being used for substantive-content evidence.

### 18.4 The verifier-exit-codes-as-cross-exam-target angle

Opposing counsel will press: "The verifier's exit code and reason strings are normative byte-for-byte. The IT witness cannot soften a FAIL into a 'minor inconsistency.' Producing a misleading PASS certification when the verifier reported FAIL is perjury under 28 USC §1746. This pins the bank to the exit code — good for both sides, but particularly good for me when the verifier reports FAIL on the produced records." *(closes Hartwell G-7)*

This angle is opposing counsel acknowledging the verifier exit-code contract works in their favor when FAIL occurs. It is a confirmation rather than a gap, but it shapes the institution's posture: the institution does not soften, deflect, or narrate FAIL outcomes. The institution's response posture has three elements.

**Element 1 — FAIL outcomes disclosed verbatim.** When the verifier reports FAIL on records produced in litigation, the institution's certification reports the FAIL exit code, the step number, and the normative reason string verbatim per the FRE 902 template's §F. The institution does not produce a soft narrative; the certification is the load-bearing document and the certification's truthfulness is binding.

**Element 2 — FAIL outcomes paired with IR record.** A FAIL outcome is paired with the institution's incident-response record showing the institution's response to the FAIL — the incident classification, the root-cause analysis, the recovery procedure, the post-event review. The IR record's integrity binding under the chain itself (cross-reference `incident-response-playbook.md`) means opposing counsel cannot challenge the IR record's authenticity. The institution's posture is "FAIL occurred, here is how we responded, here is the institution-side documentation of what happened." The bench evaluates the institution's response on its merits.

**Element 3 — Spoliation-defense framing for FAIL outcomes.** A FAIL outcome can become the basis for an FRCP 37(e) spoliation argument by opposing counsel. The institution's response posture under spoliation: the FAIL outcome occurred, the institution preserved the records under FAIL (the chain's append-only enforcement at the database level prevents tampering by the institution after the fact), the institution's IR record documents the FAIL's disposition, and opposing counsel can examine the FAIL records under the certification or under court-controlled discovery. The institution's defense rests on the chain's tamper-evidence — even a chain in FAIL state is integrity-bound up to the FAIL point, and the FAIL itself is documented evidence.

The combined effect: opposing counsel pressing the FAIL angle finds an institution that has not softened the FAIL, has paired it with documented IR response, and is offering it as evidence of the institution's tamper-evidence discipline rather than hiding it. Opposing counsel's argument shifts from "the bank is hiding a FAIL" to "the bank had a FAIL and responded to it," which is a different argument with different stakes.

### 18.5 The customer-correlation-index-unbound angle

Opposing counsel will press: "The chain captures `(tenant_id, run_id, seq)` triples. The institution maintains a customer-correlation index that maps a customer dispute to the affected runs. That index is not in the chain, not signed by the HSM, not covered by the per-event MAC. It is the institution's internal lookup table. If the index is wrong — accidentally or otherwise — the chain entries the defendant produces are integrity-bearing for what they are, but they may not be the entries the index should have surfaced." *(closes Hartwell G-8)*

The angle is the most surgical of the five — it does not challenge the chain's integrity, only the institution's procedural completeness in producing responsive entries. The institutional response posture has four elements.

**Element 1 — Customer-correlation index documented as a load-bearing artifact.** The institution's records officer documents the customer-correlation index as a load-bearing artifact for production completeness (`templates/records-management-program.md` §6.2). The index's documentation includes the procedural document governing how a dispute reference resolves to a run set, the index's audit log, the index's change history during the period at issue, and the institution's reconciliation discipline against the chain.

**Element 2 — Index integrity-binding through periodic reconciliation.** The records officer's quarterly review reconciles the index against the chain — for a sampled customer, the index's run-set is compared against the chain's run-set for that customer to confirm the index is accurate. The reconciliation outcome is recorded as an operational event in the chain itself, becoming an integrity-bound record of the index's reconciled state. Opposing counsel pressing "the index could be wrong" receives the integrity-bound reconciliation evidence showing the index was reconciled against the chain at named dates during the period.

**Element 3 — Independent re-derivation availability.** The institution's chain's append-only structure means opposing counsel's expert (under court-controlled IKM escrow per FRE 902 Addendum A2.1) can independently re-derive the run set for a given customer by querying the chain directly — bypassing the institution's customer-correlation index entirely. The independent re-derivation produces the chain-side ground truth; comparing against the institution's index-side production reveals any divergence. The institution stipulates to permitting this independent re-derivation as part of the protective-order arrangement.

**Element 4 — Adverse-inference exposure under documented procedure.** When the institution cannot produce a contemporaneous, integrity-bound version of the customer-correlation index covering the period at issue, the institution acknowledges this as an exposure. The institution's records officer's annual report (`templates/records-management-program.md` §10) names whether the index was integrity-bound during the period and whether reconciliation was performed. An institution that did not maintain integrity-binding during the period accepts that opposing counsel may seek an adverse-inference instruction; the institution's defense rests on whether the index was nonetheless accurate (provable through the reconciliation evidence) and on the institution's institutional motivation (privacy minimization, operational practice, etc.).

The combined effect: opposing counsel's argument that "the right entries may not have reached production" is met with the institution's documented reconciliation evidence, the offer to permit independent re-derivation, and the records officer's annual reporting on the index's discipline. Opposing counsel's argument shifts from "the institution may have suppressed entries through index error" to "the institution should have maintained better index discipline" — a procedural argument the institution either has answered (with documented reconciliation) or has not (in which case the records officer's annual report acknowledges the gap).

### 18.6 Pre-deposition checklist for adversarial litigation

When the institution's IT custodian or records officer is being deposed in adversarial litigation, the pre-deposition checklist:

- [ ] FRE 902 certification's adversarial-balance addendum reviewed and adopted (cross-reference `templates/fre-902-certification.md` Addendum A1-A5)
- [ ] Independent verifier-runner identified and named in the certification (Addendum A1.1)
- [ ] External verifier-runner posture chosen if applicable (Addendum A1.2: SOC 2 audit firm, court-appointed neutral, regulator-side examination team)
- [ ] Court-controlled IKM escrow arrangement negotiated if applicable (Addendum A2.1)
- [ ] Three-layer rebuttal articulation reviewed by the witness (Addendum A3.1)
- [ ] Total-count disclosure prepared for any partial-disclosure bundle (§18.1 Element 1)
- [ ] Basis-specificity declaration prepared for any partial-disclosure bundle (§18.1 Element 2)
- [ ] Data-classification policy at time of capture preserved and producible (§18.2 Element 1)
- [ ] Contemporaneous business rationale documentation reviewed (§18.2 Element 2)
- [ ] Vendor-side content provenance investigated (§18.2 Element 3)
- [ ] Section 1.2 epistemic-scope statement reviewed by the witness (§18.3 Element 3)
- [ ] FAIL records (if any) paired with IR record (§18.4 Element 2)
- [ ] Customer-correlation index reconciliation evidence retrieved (§18.5 Element 2)
- [ ] Records officer's annual reports for the period at issue retrieved (§18.5 Element 4)

The checklist takes the institution's documented postures and ensures the deposed witness has them in front of them. Opposing counsel's questions are anticipated; the institution's answers are the documented postures rather than improvisation.

### 18.7 The defense's read of the same posture

The honesty cuts both ways. The institution's response postures at §18.1 through §18.5 are not theatre — they are the institution's honest acknowledgment of where the chain's evidentiary value has limits and where the institution's procedures around the chain are load-bearing. A bench evaluating the institution's posture sees an institution that has documented the limits and built procedures around them. Opposing counsel evaluating the same posture sees an institution that has anticipated the cross-examination and is unlikely to produce the surprises that make adversarial discovery profitable.

The institution's outside counsel's posture: adversarial litigation against an institution operating this chain with this addendum and these procedures is likely to settle on factual disputes (was the AI's decision actually discriminatory?) rather than on evidentiary-foundation disputes (can the institution authenticate what the AI said?). The chain's cryptography plus the institution's adversarial-balance procedures move the litigation onto its merits, which is where the institution wants it.

---

## 18. Plaintiff-side balance — adversarial discovery posture for ECOA, FCRA, and state AI-fairness statutes

Sections 1 through 17 frame the chain as defense-side evidence. The institution producing chain artifacts under FRCP 34 is the proponent; the institution's witness lays the foundation under FRE 901(b)(9); the institution's counsel argues admissibility under FRE 902(13), 902(14), and 803(6). That framing is correct as far as it goes, but it is not the only framing the institution will encounter. A plaintiff-side class-action lawyer arriving after a consent order, a state-AG action, or a class-certification motion will read the chain as evidence she has subpoenaed against the institution, not evidence the institution offers in its own defense.

This section is the symmetric pair of the first 17 sections. It names the plaintiff-side discovery posture under ECOA, FCRA, TILA, ADA Title III, state UDAP statutes (California UCL/CLRA, Massachusetts 93A, New York GBL §349), and emerging state AI-fairness statutes (Colorado SB 24-205, Illinois HB 3773, NYC Local Law 144, California AB 331). The institution that closes the questions below enters class-action discovery with documented procedures, pre-positioned proportionality anchors, deposition-prepared witnesses, and a clear evidentiary boundary that prevents over-claiming. The institution that does not close them faces discovery disputes over proportionality, FRE 902 self-authentication, FRE 803(6) business-records foundation, adverse-inference arguments, and multi-tenant confidentiality, all reactive rather than prepared.

### 18.1 Discovery scope and proportionality under FRCP 26(b)(1)

When a plaintiff files a class action on behalf of 50,000 borrowers alleging ECOA disparate-impact bias, the discovery demand will read: "All documents and electronically stored information (ESI) evidencing the institution's AI decision, reasoning, inputs, and outputs for each class member." Under FRCP 26(b)(1), discovery is proportional to "the needs of the case." The institution's prepared response:

> The institution will produce, within FRCP 26(b)(1) proportionality bounds: (a) the full chain extract for every named plaintiff; (b) a representative sample of chain extracts for unnamed class members, sized per the plaintiff's statistical expert's specified power needs (typically 5%-10% of a large class, or all applicants from a randomly selected week or month); (c) the daily Merkle seals and HSM signatures covering the produced extracts (proving the entries are unaltered and complete within their seal); (d) the verifier output across the produced extracts (institution-run, with the binary's SHA-256 documented per FRE 902(13)); (e) the institution's customer-correlation-index entries resolving the named class members. The institution will not produce every entry for every class member on the theory that the chain "is transparent." Producing 50,000 full chains is disproportionate under FRCP 26(b)(1); a 5% statistically valid sample plus all named plaintiffs supports the plaintiff's expert's disparate-impact analysis without imposing an undue burden.

The proportionality anchor is the FRE 901(b)(9) authentication property — the chain authenticates the institution's process of producing AI decisions, and a representative sample plus the seal records prove the process produced unaltered outputs across the full population. The plaintiff who insists on every entry is asking for forensic-audit-at-the-institution's-expense rather than the proportional discovery the rule contemplates.

The institution's CC8.1 control description names the standard sampling strategy for class-action discovery; the records-custodian's deposition preparation references it.

### 18.2 Rule of Completeness (FRE 106) and the "context" the chain entry must carry

FRE 106 says when a party introduces a writing or recording, the other party can require production of "any other part… that in fairness should be considered." A plaintiff's expert on the stand with the chain entry showing "AI decided: DENY" will turn to the jury and say: *"Notice what's NOT in this chain entry — the model's stated reason, the model's confidence score, the customer's FICO score, the underwriter's manual review."* The defense argument that "the spec doesn't require those fields" loses if the missing fields are the substantive context FRE 106 reads as load-bearing.

The institution's pre-positioned posture: **chain entries representing model_call or decision events MUST include the model's stated reason in the canonical form**. The institution's standard `audit.*` schema for credit decisions includes:

| Field | Purpose | Litigation role |
|---|---|---|
| `audit.decision.principal_reason` | The reason given in the FCRA §615 adverse-action notice | Establishes the institution disclosed a reason; defends against "the AAN was a template" attacks. |
| `audit.decision.factors_considered` | The material factors the AI model considered, ordered by weight where available | Closes the FRE 106 completeness objection; defends against "you produced the decision but hid the reasoning." |
| `audit.decision.model_output.confidence` | The model's confidence in the decision (when available) | Closes the gap between "the AI was confident" and "the AI was uncertain"; affects damages quantification. |
| `audit.decision.notice_date_sent` | When the FCRA §615 notice was sent | Defends against "the AAN was late" attacks. |

Recording these fields preempts the FRE 106 challenge. An institution that omits them will face discovery disputes about whether the chain entry is incomplete and whether the omission is itself evidence of bias. The institution's records-custodian deposition preparation includes the answer to "where is the model's reasoning?" — pointing to the captured `audit.decision.factors_considered` payload bound under the chain MAC.

### 18.3 Causation and model-feature visibility

In an ECOA disparate-impact case, the plaintiff will argue: *"The AI made a race-proxy decision (using ZIP code as a predictor that correlates with race), and the institution can't explain why."* The institution's defense — "here's the chain entry showing what the model decided" — does not answer the question. The chain proves what the institution got OUT of the model; the plaintiff is asking what the model SAW.

The institution's pre-positioned posture: **chain entries representing model_call events SHOULD include the model's input features**, when feature recording is technically and contractually feasible. The recommended `audit.model_call.*` schema:

| Field | Type | Notes |
|---|---|---|
| `audit.model_call.inputs` | object | Structured record of the features the model observed: credit score, income, debt-to-income ratio, employment status, location-derived features, and any other inputs to the decision. |
| `audit.model_call.feature_importance` | object \| null | When available from the model: SHAP values, attention scores, or similar feature-importance signals. Null when the model is a black box that does not expose feature importance. |
| `audit.model_call.input_source` | string | The institution's identifier for the input pipeline (e.g., `credit-bureau-feed-v3`, `internal-account-history-v2`) — supports the plaintiff's expert auditing whether the inputs themselves were biased. |
| `audit.model_call.model_id` | string | Versioned model identifier (e.g., `credit-underwriting-model:v2.3.1`) — supports class-period segmentation by model version (per §18.4 below). |
| `audit.model_call.deployment_timestamp` | string | RFC 3339 UTC time when this model version was deployed to production. |

Recording these fields:

- Defends against the "black box" attack. The institution can point to the recorded features and rebut "you can't explain the decision."
- Supports the plaintiff's expert's causation analysis without forcing depositions of model-development engineers as a discovery substitute.
- Documents the institution's actual feature usage so emerging state AI-fairness statutes (Colorado SB 24-205, Illinois HB 3773, NYC Local Law 144) have an integrity-bound substrate to audit against.

The institution that omits these fields makes a deliberate trade-off the CC8.1 control description must justify — typically vendor-confidentiality of the model's internals or a specific privacy commitment to the data subject. The MRM committee's policy on feature recording is documented; deviation is documented; the records-custodian deposition prep covers the institution's reasons.

### 18.4 Model-version cohorts and the class-period partition

A plaintiff's statistical expert will argue: *"The model was trained on biased data, so all decisions from June 2024 to January 2025 are presumptively biased."* The institution's defense — "we retrained the model in October 2024 with debiasing techniques" — only works if the chain proves which model version produced which decision.

The `audit.model_call.model_id` and `audit.model_call.deployment_timestamp` fields above support this defense. The institution's class-action posture:

1. The institution's expert pulls all chain entries for the class period.
2. The expert partitions the entries by `model_id`.
3. For each model_id cohort, the expert reports decision distribution, disparate-impact ratios, and the deployment timestamp.
4. The plaintiff's expert receives the partition and conducts independent disparate-impact analysis per cohort.

The institution that doesn't record model versions faces discovery demands for "model genealogy" documentation and depositions of model-development engineers as a substitute for the chain's silence. Recording the model version in every model_call entry closes this gap preemptively.

### 18.5 Linked fair-AI audit records — the chain plus the disparate-impact audit

Emerging state AI-fairness statutes require institutions to document training-data composition and to audit for disparate impact on protected classes. The chain proves what decisions the AI made; the chain does NOT prove what data the AI was trained on or whether the institution tested for bias. A plaintiff's expert will argue: *"The chain is useless without the training-data documentation; I can't assess bias without seeing the model's inputs."*

The institution's pre-positioned posture: **maintain a linked record between chain entries and the institution's disparate-impact audit results**. Each chain entry's model_id resolves (via an institution-side audit-record index, NOT in the chain) to:

- The training-data composition for that model version (percentage of positive outcomes by race, ethnicity, sex, where the institution legally collects this).
- The institution's disparate-impact testing results for that model version (statistical disparate-impact ratio, four-fifths-rule analysis, demographic-parity metric).
- The institution's remediation steps (if any) when disparate impact was found.

The audit record is discoverable as a business record under FRE 803(6). The plaintiff demands the audit (which has direct factual relevance) without demanding the training data itself (which may be vendor-confidential under a protective order, and which the audit makes substantively unnecessary for the plaintiff's analysis).

The institution that doesn't maintain this linkage will face discovery battles over whether training-data composition is proportional to produce. The institution that does maintain it gives the plaintiff what the plaintiff actually needs (the disparate-impact analysis) without exposing vendor-confidential training data.

### 18.6 Routing and provider-selection logic

In a multi-LLM deployment, the institution routes different decision types to different vendors. The chain's `chain_kind = 'routing'` entries record what was routed where; the institution's `audit.routing.*` payload should also record WHY the routing decision was made. A plaintiff's expert will argue: *"The institution chose Vendor X knowing Vendor X's model has a documented bias problem; the routing decision is the real evidence of intent."*

The institution's `audit.routing.*` schema, extending spec §4.4.1:

| Field | Notes |
|---|---|
| `audit.routing.selection_reason` | E.g., `task type = mortgage, provider = vendor-a per institution policy`, or `random-ab-test: cohort = treatment`, or `failover: primary timeout`. |
| `audit.routing.provider_audit_status` | E.g., `last-disparate-impact-audit: vendor-a PASS 2024-03-01` — the institution's most recent audit of the selected provider. |
| `audit.routing.policy_version` | The institution's routing policy at decision time. |

The selection reason lets the plaintiff's expert audit the institution's routing logic for bias; the audit status lets the plaintiff cross-check whether the institution knowingly selected a biased vendor. Omission exposes the institution to adverse inference (the jury may assume the worst if the selection logic isn't documented).

### 18.7 Litigation hold under FRCP 37(e)

FRCP 37(e) requires "reasonable steps to preserve" ESI once litigation is "reasonably anticipated." The institution's standard chain-data retention is 7 years (per spec §10.13); litigation hold extends retention beyond the standard schedule. The institution's litigation-hold procedure:

1. **Trigger.** Counsel issues a litigation hold notice when litigation or a class-action filing names the AI system as a defendant, or when notice of a regulatory investigation arrives that is reasonably anticipated to escalate to litigation.

2. **Scope.** The hold covers: (a) chain entries for the affected period and class members; (b) IKM rotation history for the affected period (the institution does NOT rotate keys during the hold without explicit before-and-after custody records); (c) verifier output for the affected period; (d) operational events covering the period (the `master.*`, `seal.*`, `verifier.*` events in `soc-pack/control-evidence-events.md`); (e) any internal AI-bias tests, disparate-impact audits, or adverse-action-notice audits conducted during or covering the period; (f) change-management records (deployments, configuration changes, model retraining events).

3. **Suspension of automatic destruction.** The institution's automatic chain-data destruction schedule is suspended immediately for the affected scope. The IT team confirms in writing to legal that the hold is in place and the data cannot be deleted by the standard retention process. The hold persists until the case is resolved or the court orders otherwise.

4. **Litigation-hold memo.** A written memo names the data subject to the hold, the reason, the affected `(tenant_id, run_id, seq)` ranges (when known) or the affected period (when not), and the institution personnel responsible for confirming compliance.

5. **Hold release.** When the hold is lifted (case resolved, settlement approved, court order releasing), the institution documents the lifting, restores the standard destruction schedule, and the IT team confirms in writing that the hold is released.

Failure to issue a timely hold is spoliation under FRCP 37(e)(1) and can result in sanctions or adverse-inference instructions to the jury. A plaintiff who asks for entries older than 7 years (before the institution adopted the chain) gets a clean answer: those entries don't exist because the institution didn't operate the chain at that time. A plaintiff who asks for entries that should have been preserved under a hold but were destroyed gets an FRCP 37(e) sanction motion.

### 18.8 Records-custodian deposition preparation

The plaintiff's attorney will depose the institution's IT manager (the records custodian) and ask the questions below. The institution's deposition-preparation script:

| Question | Prepared answer |
|---|---|
| "How do you know this chain entry is authentic?" | "The entry is authenticated by a per-event MAC computed under a per-tenant key, plus a daily Merkle seal signed by our HSM. The full procedure for verifying authenticity is in spec §7 — a 12-step procedure that examines structural integrity, MAC, seal signature, and the chain's link to the prior entry. We run that procedure regularly; the verifier output is preserved as a `verifier.run_completed` operational event in the chain itself." |
| "What would cause this entry to be rejected by the verifier?" | "Any of the 12 failure modes in spec §7. The negative test vectors N001 through N022 in our test corpus exercise each failure mode. The vectors are public; anyone with the verifier and the test corpus can confirm the procedure rejects every documented attack scenario." |
| "Have you ever seen a chain entry fail verification?" | (If yes) "Yes — I can describe the incident, the root cause, and the remediation per our incident-response playbook §Scenario [X]. The IR record documents detection, root cause, and resolution; I can produce it under proportional discovery." (If no) "No. Our institutional controls — the spec §7 verifier procedure, the daily seal job, the IKM-fingerprint reconciliation — are designed to detect failures, and we have not observed one in the period under discussion." |
| "Can you explain this in plain English for a jury?" | "The chain proves what the AI said and proves no one altered the record of what it said. The chain does not prove the AI was correct, the AI was unbiased, or the AI's output complied with policy. Those are separate questions the institution answers with separate evidence — the institution's policy documents, the bias-audit results, the human-review records." |
| "When you rotate the master key, what happens to old chain entries?" | "The IKM history is preserved in our key registry. Old entries remain verifiable under their original IKM; new entries derive from the rotated IKM. The verifier resolves the right IKM by `(tenant_id, key_version)` per the registry. We do not rotate during a litigation hold without explicit before-and-after custody records." |
| "When the HSM goes offline, what happens?" | "The seal job cannot complete on a day the HSM is unavailable. The chain entries for that day are captured normally (the per-event MACs do not depend on the HSM); the seal record is delayed until the HSM is back online. The delayed seal is itself a documented operational event; the IR procedure for HSM unavailability is in our incident-response playbook." |
| "When a chain entry fails MAC verification, what's your process?" | "The failure surfaces in the next verifier run as a §7 step 9 rejection. The IR procedure: pull the failing entry, identify the root cause (data corruption, IKM-registry mismatch, malicious tampering), document the finding, and remediate per the relevant IR scenario. The IR record is preserved." |

The institution's IT team is trained on these answers. The team's ability to point to supporting documentation — the test vectors, the incident logs, the HSM certification, the spec text — closes cross-examination on the verifier's mechanics. The team's discipline of refusing to speculate beyond the verifier and the IT operations protects the institution's credibility.

### 18.9 Verifier-credibility defense — the "in-house software" attack

A defendant's expert will sometimes attack a plaintiff-side verifier as "home-grown software written by the plaintiff's IT person." The plaintiff's response — and the institution's pre-positioned posture in case the institution finds itself on the same side of the table — is verifier-validation documentation:

1. **The verifier was built against the public FFIEC chain-of-custody-v1.0a spec.** The spec defines a 12-step verification procedure; the verifier implements those steps.

2. **The verifier passes all positive and negative test vectors in the public test corpus.** Test vectors 001, 002, 003, 008, 010, 015, 016, and N001 through N022 are public; any verifier implementation that passes them implements the spec correctly.

3. **The verifier's output format matches spec §7's normative output** — the three-line `Status: PASS/FAIL/PASS-STRUCTURALLY/PASS-WITH-ANOMALY`, `Step: N`, `Reason: <text>` form. The format is byte-deterministic; two correct verifiers produce identical output on identical input.

4. **An independent review confirms the verifier's results** — a cryptographer or compliance expert reviews a sample of entries and confirms the verifier's verdicts. A second-implementation cross-check (running a different verifier built independently from the same spec) is the strongest version of this confirmation.

5. **The reference verifier is Apache 2.0 open source** with reproducible builds; cosign signatures; SLSA build provenance. Anyone can validate that the binary they hold corresponds to the public source.

The "in-house software" attack fails when the verifier is a publicly documented procedure with public test vectors that anyone can re-run. The institution's expert testimony on its own verifier is grounded in the same five points; the plaintiff's expert testimony on a plaintiff-side verifier follows the same template.

### 18.10 Evidentiary boundaries — what the chain proves and what it does not

The institution's expert testimony must explicitly tell the jury the chain's limitations. Spec §1.2 epistemic scope is the load-bearing text; the expert's testimony reads it directly:

> The chain proves: (1) what decision the AI made on this customer's application; (2) that no one altered the decision after the fact. The chain does NOT prove: (3) the decision is factually correct; (4) the decision complies with the bank's policy; (5) the decision is free of bias; (6) the decision was the same one a non-AI process or an alternative AI would have produced.
>
> The decision's correctness, policy compliance, and fairness are separate questions answered with other evidence — the customer's actual income (to check if "insufficient income" was accurate), the bank's written policy (to check if the decision complies), statistical testing (to check if the AI is biased), and the institution's MRM-committee records (to check if the institution validated the model). If the AI claimed the customer had insufficient income but the customer's pay stubs show high income, the chain proves the AI made the claim, but the customer's pay stubs prove the claim is false.

The institution's expert NEVER says "the chain proves the decision was fair" or "the chain proves the institution operated lawfully." Either statement is an evidentiary overreach the plaintiff's cross-examination will dismantle. The chain's credibility for its real purpose (authenticity) survives only when the institution does not over-claim. The plaintiff's expert benefits from the same boundary discipline — the chain is evidence of authenticity, not evidence of fairness, and the plaintiff's expert testimony respects the same limit when relevant to the plaintiff's case.

### 18.11 Hearsay framing — business records exception versus AI assertions

When the institution produces a chain entry showing "AI Model said: DENY because insufficient income," the defendant will sometimes argue the AI's statement is hearsay. The plaintiff's response — and the institution's pre-positioned response when the plaintiff demands the chain — is the FRE 803(6) business-records framing already covered in `customer-dispute-procedures.md` §"FRE 803(6) hearsay exception and the chain's reliability framing":

- The chain entry's metadata fields (timestamps, tenant identifiers, run identifiers, the verifier output) are clearly business records under FRE 803(6).
- The AI's stated reasoning is admissible under a residual reliability exception grounded in the chain's MAC + Merkle + HSM coverage. The chain proves no post-hoc tampering; the IT witness lays the foundation by walking the verifier procedure.
- The AI's statement is not offered for the truth of the asserted fact ("the customer had insufficient income"); it is offered to prove the institution's decision was made on that basis. The institution's separate substantive evidence — the customer's actual income records — proves whether the asserted fact was accurate.

The institution's hearsay framing in a class-action production matches its framing in a customer dispute. Consistency across the two contexts is itself credibility evidence.

### 18.12 Multi-tenant production scope and confidentiality

In a vendor-hosted deployment with multiple institutions using the same ledger, a plaintiff might argue: *"We need the entire ledger to ensure you didn't hide our client's entries or swap them with another customer's."* The institution's response anchored in the chain's multi-tenant cryptographic isolation (per spec §4.1):

> The chain's multi-tenant design isolates each tenant's entries cryptographically: the per-event `key_fingerprint = SHA-256(utf8(tenant_id) || ikm)[:16]` binds every entry to its tenant's IKM. The Merkle seal is computed per-tenant per-day; one tenant's seal does not cover any other tenant's entries. The verifier's spec §7 step 4 cross-chain-lift detection rejects any entry whose tenant_id does not match the chain's header.
>
> We will produce: (a) the chain entries for our tenant covering the affected period and class members; (b) the daily Merkle seals covering those entries; (c) the verifier output across the produced extracts. The Merkle seal cryptographically proves the integrity of all entries for our tenant-day without disclosing other customers' or other institutions' data. The entries for other tenants are subject to separate confidentiality agreements and are not discoverable in this case.

The multi-tenant cryptographic isolation is the technical foundation for the confidentiality objection. The plaintiff's verifier-side validation runs against the produced tenant's entries and the produced seals; the plaintiff cannot argue "you cherry-picked which entries to produce" because the seal record covers the full tenant-day and would surface any deletion as a Merkle-root mismatch.

### 18.13 Cross-border discovery — GDPR and Schrems II

When a U.S. court compels production of chain entries containing EU customer personal data, the institution faces a conflict between the U.S. discovery order and EU data-protection law. The institution's prepared response:

1. **Seek a protective order from the U.S. court.** Limit distribution to the parties' counsel and the court; not the public. Cite Schrems II (Case C-311/18) and the EU's privacy framework's interaction with U.S. discovery.

2. **Assert a GDPR conflict in the discovery objection.** Article 44 GDPR restricts cross-border transfers of personal data to third countries. A U.S. discovery order producing EU customer chain entries to a U.S. court may trigger Article 44 review depending on the safeguards in place.

3. **Offer secure-facility discovery.** Conduct the discovery in a facility on the institution's servers in the EU under the control of a neutral third party (a special master or a court-appointed examiner) rather than producing paper or digital copies to the U.S. lawyers.

4. **Tokenize before production.** The chain entries contain customer-PII tokens (per `docs/privacy-by-design.md`) rather than direct identifiers. The institution's production includes the tokenized form; the institution does NOT produce the customer-correlation-index resolving tokens to identifiers without separate legal authorization. The plaintiff's expert can conduct disparate-impact analysis on tokenized chain entries when the institution provides cohort-level demographic indicators (per the EU consumer-procedure section in `customer-dispute-procedures.md`).

5. **Pattern B production.** When the institution operates Pattern B (per-region `tenant_id` per spec §10.15), the U.S. discovery is bounded to U.S.-tenant entries; EU-tenant entries are produced under EU procedural rules in EU forums. The institution's CC8.1 names the per-region tenant allocation and the production posture per region.

The institution's outside counsel briefs on Schrems II and Article 44 before litigation is anticipated; the cross-border posture is documented in the institution's litigation-hold procedure so the records-custodian does not produce EU-tenant entries to a U.S. court without the protective-order workflow above.

### 18.14 Examiner-acceptance versus courtroom-admissibility

A plaintiff might argue: *"The FFIEC examiner accepted this chain and certified it as conforming to the spec; that proves the chain is reliable and admissible in court."* The institution's response: regulatory acceptance is not the same as judicial admissibility. The FFIEC examiner's job is compliance oversight; Daubert standards are different.

The institution preparing for litigation obtains a separate Daubert-focused validation from a third-party expert or from outside counsel, addressing the four Daubert factors:

1. **Testability** — spec §7 + the test-vector corpus (this document §2.1).
2. **Peer review** — the FFIEC working group, outside reviewers, academic vetting (this document §2.2).
3. **Known error rate** — the threat model's three independent integrity layers (this document §2.3).
4. **General acceptance** — HMAC-SHA-256 and Merkle-RFC-6962 are NIST-standardized and widely deployed (this document §2.4).

The Daubert validation is separate from the examiner's conformance review. The institution presents BOTH at trial: regulatory approval (the examiner accepted the chain as conforming) AND independent evidentiary credibility (the Daubert factors are satisfied). The plaintiff's expert must address the Daubert factors in challenging the chain's admissibility; "the examiner accepted it" is not enough on either side.

### 18.15 Batch verification at scale

When a plaintiff's expert receives 50,000 chain entries, running the verifier on each entry one-by-one is impractical. The expert's batch-verification procedure:

1. **Compute the Merkle root for each day's entries** in the produced extract; spot-check against the institution's daily seal record.
2. **Spot-check 5-10% of random entries** under the full §7 verifier procedure (per-entry MAC, fingerprint, chain link) for randomly selected entries from each day.
3. **Aggregate statistics** — count total entries per day, per run, per decision type; compare against the institution's reported statistics for the production period.
4. **Identify any unexplained gaps** between produced extracts and the institution's reported statistics; flag for follow-up discovery.

A 5% random sample identifies tampering in 95% of cases under the Merkle-walked detection model; 10% covers 99% of cases. The expert documents the sample size and the verified entry count in the expert's report. The institution's expert does the same in its own validation.

### 18.16 The institution's ECOA / FCRA pre-positioning checklist

When ECOA / FCRA / state-AI-fairness-statute litigation is reasonably anticipated:

- [ ] CC8.1 control description names: feature recording posture (§18.3), audit-record linkage to chain entries (§18.5), routing-selection-reason recording (§18.6), litigation-hold trigger and scope (§18.7).
- [ ] Standard `audit.*` schema includes: `audit.decision.principal_reason`, `audit.decision.factors_considered`, `audit.decision.model_output.confidence`, `audit.decision.notice_date_sent`, `audit.model_call.model_id`, `audit.model_call.deployment_timestamp`.
- [ ] Records-custodian deposition prep covers the §18.8 question/answer template.
- [ ] Verifier-credibility validation documented per §18.9.
- [ ] Daubert-focused validation obtained independent of the FFIEC examiner's conformance review (§18.14).
- [ ] Cross-border production posture documented if the institution operates in EU markets (§18.13).
- [ ] Multi-tenant confidentiality argument prepared if the institution uses a vendor-hosted ledger (§18.12).
- [ ] Statistical-sampling protocol for class-action discovery established with statistical justification (§18.1).
- [ ] Litigation-hold procedure tested in tabletop exercise; IT team's hold-toggle procedure documented (§18.7).

When the litigation arrives, the institution's response is reactive only at the operational level (which `(tenant_id, run_id, seq)` ranges, which timing) — the substantive posture is pre-positioned. The plaintiff's first-draft motion to compel meets a documented institution response, not a scrambled one.

### 18.17 Cross-reference

`customer-dispute-procedures.md` — the consumer-side dispute procedure (the same chain entries the plaintiff demands in litigation are the entries the institution discloses to consumers under GDPR Article 15). Maintaining a single posture across both contexts is itself credibility evidence — the institution is not producing one shape to consumers and a different shape under FRCP 34. `regulator-pack/finding-language.md` — the verifier-output finding-language references the IT witness uses in deposition. `legal-disclosure.md` — the master-key disclosure procedure when court-ordered. `incident-response-playbook.md` — the IR records preserved under litigation hold. `docs/AI-safety-evaluation-overlay.md` §10 — the chain's role in capability and OOD evaluation, relevant when the plaintiff's expert argues the model was deployed without adequate testing.
