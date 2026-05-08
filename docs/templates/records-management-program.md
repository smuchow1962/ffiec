---
title: Records-Management Program (template)
status: normative-template
aligned-with:
  - Federal Records Act (44 U.S.C. § 3105) — records-officer designation model
  - 12 CFR Part 30 Appendix B — records-management standards for federally-supervised banks
  - NARA General Records Schedule 4.2 (Information Access and Protection Records)
  - NARA General Records Schedule 5.6 (Security Records)
  - SEC Rule 17a-4 (17 CFR 240.17a-4) — broker-dealer records
  - CFTC Rule 1.31 (17 CFR 1.31) — commodity records
  - FINRA Rule 4511 — general books and records
  - ISO 15489-1:2016 — records management, four-characteristic model
  - ISO 30300:2020 — records-management leadership accountability
  - ISO 14721:2012 — Open Archival Information System (OAIS) reference model
  - ISO 16363:2012 — trustworthy digital repository audit and certification
  - 5 USC §552 (FOIA) and exemptions (b)(4), (b)(8); 5 USC §552a (Privacy Act)
  - GLBA Title V — Safeguards Rule
date: 2026-05-07
version: 1.0.0
companion-templates:
  - docs/regulator-pack/retention-justification.md (legal-grounding analysis)
  - docs/templates/fre-902-certification.md (records-custodian declaration)
  - docs/regulator-pack/nydfs-part500-overlay.md (§500.13 records-disposition mapping)
---

# Records-Management Program (template)

> **What this template is.** The records-program scaffolding that wraps around the chain-of-custody cryptographic substrate. The chain delivers authenticity and integrity (two of the four ISO 15489 characteristics) at a level of cryptographic rigor most banking-records programs do not match. This template adds the records-management vocabulary and accountability around the substrate so the institution's records officer, the institution's senior officer, and the regulator's records inspector see the records program a federally-supervised bank actually operates under: a records officer with named accountability, a records schedule listing each record series with retention triggers and dispositions, a certificate of destruction at the disposition boundary, symmetric hold-release procedures, a preservation plan covering the cryptographic-primitive migration horizon, FOIA / Privacy Act / GLBA exemption posture, and a usability layer the records officer maintains independently of the chain's cryptographic operation.

> **What this template is NOT.** A substitute for the institution's existing records-management program. Where the institution already operates a records-management program with a designated records officer and a maintained records schedule, this template is a cross-walk that names how chain-of-custody record series fit into the existing program. Where the institution does not yet operate a records-management program, this template is the starter scaffold the records officer extends to the institution's full record holdings.

---

## 1. Scope and reading order

This template is consumed in three reading orders.

| Reader | Reading order |
|---|---|
| Institution adopting the chain for the first time | §2 (records-officer designation) → §3 (records schedule) → §4 (certificate of destruction) → §5 (symmetric hold-release) → §7 (preservation plan) |
| Records officer at an institution with an existing records-management program | §3 (records-schedule rows for chain record series) → §6 (retrieval and presentation — the usability layer) → §7 (preservation plan) |
| FOIA officer at a regulator receiving a FOIA request against the institution's chain records | §8 (FOIA / Privacy Act / GLBA exemption posture) → cross-reference to §500.06 audit-trail framing in `nydfs-part500-overlay.md` §2.1 |

The template does not duplicate material that already lives in `regulator-pack/retention-justification.md`. The retention-justification document is the legal-grounding analysis (the "why 7 years?" question); this template is the records-program scaffolding (the "who is the records officer, what does the records schedule look like, when does the certificate of destruction fire?" questions). The two documents pair — the records schedule references the retention-justification for legal grounding; the retention-justification stops carrying records-schedule freight and reverts to its proper role.

---

## 2. Records-officer designation

The Federal Records Act (44 U.S.C. § 3105) requires every federal agency to designate a records officer. For federally-supervised banks operating under 12 CFR Part 30 Appendix B and analogous prudential frameworks, the institution's records-management program with named accountability is the parallel discipline. The chain-of-custody specification reaches all the way to FIPS 140-2 Level 3 HSM custody but does not name a records officer; this template closes that gap. *(closes Park-Whittaker F-1)*

### 2.1 The records officer's authority

The institution's records officer holds five authorities the chain's existing operating model does not allocate to any other role.

| Authority | What the records officer does |
|---|---|
| (a) Records schedule custody | Maintains the records schedule at §3 below; reviews it annually; updates it when new record series are added or existing series change |
| (b) Disposition authorization | Signs the certificate of destruction at the 7-year boundary (or whatever the schedule names); signs disposition certifications for transfers to long-term archive |
| (c) Hold release authorization | Signs the hold-release document closing a legal hold; the symmetric counterpart to opening counsel's hold-opening signature |
| (d) Format-migration approval | Approves cryptographic-primitive migrations under the preservation plan; signs the migration certificate that becomes itself a record series |
| (e) FOIA / Privacy Act / GLBA exemption posture | Maintains the institution's FOIA-exemption posture; coordinates with regulator FOIA officers when chain records are subject to FOIA requests at the regulator |

The five authorities are independently necessary. A program that gives (a) and (b) to records management but leaves (c) implicit, (d) ambiguous, and (e) unstated is a partial program; the regulator inquiry will eventually find one of the unallocated authorities and the program will not have an answer.

### 2.2 Where the records officer fits in the operating model

The institution's existing operating-model roles named in `operator-guide.md`:

- Chain-operations team — runs the SDK, the receiver, the seal job
- Privacy / data-protection team — operates the privacy-store mapping; runs the deletion procedure under `retention-justification.md` §5
- SOC team — operates security monitoring; consumes the operational-event stream
- Legal team — issues legal holds; advises on subpoena response
- HSM administrator — operates the HSM custody

The records officer sits beside these roles, not under any of them. The records officer's accountability is to the senior officer or the board (parallel to the CISO under NYDFS §500.04 mapping). The records officer can be a designation on top of an existing privacy or compliance role at smaller institutions; the designation is what matters, not the headcount.

| Operating-model role | Records-officer interaction |
|---|---|
| Chain-operations team | Receives records-schedule guidance from records officer; routes record series correctly through retention buckets |
| Privacy / data-protection team | Coordinates with records officer on the §5 deletion procedure; obtains records-officer signature on the certificate of destruction |
| Legal team | Coordinates with records officer on hold openings (legal team signs) and hold releases (records officer counter-signs the closing) |
| HSM administrator | Coordinates with records officer on the format-migration plan; the migration certificate becomes a record series the records officer maintains |
| Senior officer / board | Receives the records officer's annual report (parallel to the CISO's §500.04(b) report); reviews the certificate-of-destruction log; reviews any preservation-plan triggers |

### 2.3 Records-officer designation document

The institution publishes a records-officer designation document. The document names:

- The individual designated as records officer (full legal name and title)
- The effective date of the designation
- The senior officer or board committee approving the designation
- The five authorities at §2.1 above (verbatim)
- The records officer's annual reporting cadence to the senior officer or the board
- The records officer's review and reaffirmation cadence (typically annually with the §500.17(b) certification cycle)
- The succession plan if the designated individual leaves the institution

The designation document is itself a record series — typically retained for the longer of 7 years or the institution's records-schedule retention for governance documents.

### 2.4 Records-officer attestation in the §10 control-attestation matrix

The institution's CC8.1 control description adds a records-officer attestation row to the §10 control-attestation matrix. The row names the records officer's authority over the §500.13 limitations-on-data-retention obligation and the chain's spec §10.9 IKM retention coupling. The records officer's attestation is the institution-side evidence in the §500.17(b) annual-certification evidence map (`nydfs-part500-overlay.md` §10.2 row for §500.13).

---

## 3. Records schedule

A records schedule is the formal document that names each record series, names the retention period for that series, names the legal authority for the retention period, and names the disposition action at the end of the period. The retention-justification document at `regulator-pack/retention-justification.md` carries the legal-grounding analysis (the "why 7 years?" question); this records schedule carries the records-series-by-records-series detail (the "what record series, what trigger, what disposition?" question). *(closes Park-Whittaker F-2)*

### 3.1 Record series the chain produces

The chain-of-custody system produces eight distinct record series. Each is governed independently with its own retention trigger, retention period, disposition action, and disposition authority. The records officer maintains the schedule below; it is reviewed annually and updated when the spec or the institution's operating model changes.

| Series | Description | Retention trigger | Retention period | Disposition action | Disposition authority |
|---|---|---|---|---|---|
| RS-1 Chain entries | Per-event JSON records of AI activity (spec §4.1, §4.4) | End of capture day | 7 years (composite per `retention-justification.md`) | Cryptographic erasure of privacy-store mapping; chain-entry storage may continue longer per institution policy | Records officer + privacy lead (joint) |
| RS-2 Daily seal records | Per-tenant-day Merkle root + HSM signature (spec §4.2, §4.3) | End of seal-day | 7 years (matched to RS-1) | Retain in cold storage; integrity-bound record of the tenant-day | Records officer |
| RS-3 IKM history | Tenant master keys and key versions (spec §10.6, §10.9) | End of last-entry-MAC'd-under-key retention | Longer of: (a) 7 years; (b) any held entry's hold release | Cryptographic erasure of IKM material | Records officer + HSM administrator (joint) |
| RS-4 Verifier output bundles | Per-period verifier PDF + JSON (spec §7) | Date of verifier run | 7 years (matched to underlying chain entries) | Retain in evidence-package archive | Records officer |
| RS-5 Operational events | `chain.verification_failure`, `master.*`, `incident.*`, etc. (spec §10.2) | End of capture day | 7 years (matched to RS-1) | Cryptographic erasure consistent with RS-1 | Records officer + chain-operations lead (joint) |
| RS-6 Key-fingerprint reconciliation evidence | Weekly key-fingerprint reconciliation outcomes (spec §10.1) | Date of reconciliation | 7 years | Retain alongside RS-2 | Records officer |
| RS-7 Disposition-documentation series | Certificates of destruction (§4 below); hold-release documents (§5 below); migration certificates (§7 below) | Date of disposition action | Longer of 7 years or institution's records-schedule retention for governance documents | Retain for the longer period | Records officer (sole authority) |
| RS-8 Hold documents | Legal-hold opening and release documents (§5 below) | Date of hold release | Longer of 7 years or the underlying matter's resolution plus 7 years | Retain for the longer period | Records officer + general counsel (joint) |
| RS-9 Redaction policy versions | Per-policy `audit.redaction.policy_id` + `audit.redaction.policy_version` artifacts (spec §10.22) | Date the policy version is superseded | Longer of: (a) 7 years; (b) the retention of any chain entry emitted under that `policy_version` | Retain in evidence-package archive (the policy text the chain entries reference) | Records officer + privacy lead (joint) [honors §10.22] |
| RS-10 Training-data shard manifests | Provider-side training-data shard hashes + manifest (spec §10.20, §10.21 `training_shard_manifest_sha256`) | Effective date of the model the shards trained | Longer of: (a) the longest active deployment window of any model trained on the data; (b) 7 years; **plus an investigation buffer of 60-90 days** | Cryptographic erasure of the underlying shard storage; manifest hash retained under RS-7 | Records officer + ML platform lead (joint) [honors §10.20, §10.21] |
| RS-11 Coverage-map versions | `chain.coverage_map_published` operational events + the named map document (spec §10.19) | Date the map version is superseded | Longer of 7 years or the retention of any chain entry covered by that map version | Retain alongside RS-7 | Records officer (sole authority) [honors §10.19] |
| RS-12 HSM ceremony attendance logs | Paper-and-PDF ceremony attendance logs whose hash is bound under `chain.partition_ceremony_attended` (spec §10.17) | Date of ceremony | Longer of 7 years or the retention of any chain entry MAC'd under the IKM whose custody the ceremony affected | Retain in compliance vault per institution's CC8.1 | Records officer + HSM administrator (joint) [honors §10.17] |

> **Technical Notes — what the Round-17 record series add.** RS-9 through RS-12 land four record series the v1.0a template predated. RS-9 keeps the redaction policy text the chain entry references via `audit.redaction.policy_id` / `policy_version` (§10.22) so a CFPB examiner reading a chain entry years later can still resolve the policy version named in the entry. RS-10 keeps training-data shard manifests for the longest deployment window plus an investigation buffer, so a regression detected six months into deployment can still be traced to a specific shard set (§10.20). RS-11 keeps every published version of the chain-coverage map (§10.19) so an 18-month-lookback acquirer-side IT due-diligence team can confirm the map version in force on any given date. RS-12 keeps the paper attendance logs whose hashes are bound under HSM partition-ceremony events (§10.17), pairing the cryptographic attestation with the dispute-resolution paper record.

### 3.2 Retention triggers — the load-bearing detail

A retention period without a clearly named trigger is ambiguous. Each series above names a specific trigger. The legal-grounding analysis at `retention-justification.md` §4 names the "longer of" rule — the trigger date is whichever is later among the applicable supervisory-floor triggers and the civil-statute-of-limitations triggers.

For RS-1 chain entries and RS-5 operational events, the trigger is "end of capture day" with the 7-year clock running from that date. The "end of capture day" is unambiguous because the chain itself records the capture timestamp under integrity binding; the records officer's audit confirms the trigger date by reading the chain.

For RS-3 IKM history, the trigger is more complex — the "longer of" rule fires against any chain entry MAC'd under that IKM. As long as any chain entry is retained under hold or under standard retention, the IKM cannot be erased because re-verification would fail. The records officer's quarterly retention audit (cross-reference `retention-justification.md` §8) confirms the trigger date by checking whether any chain entry under the IKM is still held.

For RS-8 hold documents, the trigger is the date of hold release. Hold documents retained for "matter resolution plus 7 years" cover the case where the underlying matter (litigation, regulator inquiry, criminal investigation) takes longer than 7 years; the documents survive the matter for the institution's standard governance-document retention.

### 3.3 Disposition actions — what happens at the trigger

| Action | What it means | Evidence produced |
|---|---|---|
| Cryptographic erasure | Destruction of the IKM (or other key material) such that the underlying records become cryptographically unrecoverable | Certificate of destruction (§4 below) |
| Retain in cold storage | Move from active retention to archival storage; no destruction | Migration log entry; cold-storage manifest |
| Retain in evidence-package archive | Continue active retention as evidence | No additional disposition; the record continues active retention |
| Transfer to NARA | Apply only to federal-records authorities; not the typical action for federally-supervised banks | NARA-specific accession documentation |

The chain's primary disposition action is cryptographic erasure under RS-3 IKM history. The chain-entry storage (RS-1) may continue beyond 7 years if the institution chooses, but with the IKM erased the entries become cryptographically unverifiable — the institution can show the entry bytes but cannot prove they have not been altered. This is by design: the 7-year retention is the rebut-discoverable horizon for the underlying business records; the integrity binding is removed at the disposition boundary.

### 3.4 Records-schedule maintenance

The records officer maintains the schedule under a documented annual review cycle. The review confirms:

- Each record series above is still produced by the chain
- The retention period reflects the current legal grounding (the retention-justification document is checked for changes)
- The disposition action and authority are still appropriate
- New record series produced by spec amendments since the last review have been added
- Record series that have been retired (no longer produced by the current spec or operating model) are documented with their final-retention treatment

Changes to the records schedule are themselves a record series — the schedule's change history is retained alongside RS-7 disposition documentation.

### 3.5 Cross-reference to NARA General Records Schedules

For institutions operating against NARA's federal records schedules as a reference model:

| NARA GRS | Chain record series correspondence |
|---|---|
| GRS 4.2 (Information Access and Protection Records) | RS-7 disposition documentation; RS-8 hold documents |
| GRS 5.6 (Security Records) | RS-2 seal records; RS-3 IKM history; RS-4 verifier output |

NARA GRSs are not normative for federally-supervised banks (which operate under their prudential supervisor's records expectations rather than NARA's), but the GRS structure is the documented federal model the institution can cite when explaining its records-schedule shape to a regulator.

---

## 4. Certificate of destruction

When records reach end of retention and the disposition is destruction, the institution's records officer issues a certificate of destruction. The certificate is the records-officer-signed attestation that disposition was authorized, lawful, and complete. It is distinct from the deletion log — the deletion log captures the operator-side fact ("operator X deleted N rows on date Y"); the certificate is the authority-side attestation. *(closes Park-Whittaker F-3)*

### 4.1 Certificate-of-destruction template

The template below is filled in by the records officer at the disposition action and signed by both the records officer and a named witness (typically the privacy lead or the operator who performed the destruction). The signed certificate is retained under RS-7 disposition documentation.

```
CERTIFICATE OF DESTRUCTION

Institution: [INSTITUTION:fillin — institution legal name]
Records-management program reference: [INSTITUTION:fillin — program reference number]

DESTRUCTION DETAILS

Record series: [INSTITUTION:fillin — record-series identifier from §3.1, e.g., RS-1, RS-3, RS-5]
Description: [INSTITUTION:fillin — brief description of the records destroyed]

Tenant identifier(s): [INSTITUTION:fillin — list of tenant_id values]
Volume destroyed: [INSTITUTION:fillin — count of records or count of IKMs]
Date range covered: [INSTITUTION:fillin — earliest capture date] through [INSTITUTION:fillin — latest capture date]

Destruction method: [INSTITUTION:fillin — one of:
  (a) Cryptographic erasure — IKM material destroyed under HSM administrator procedure
  (b) Cryptographic erasure — privacy-store mapping rows deleted under privacy lead procedure
  (c) Physical destruction of media containing the records
  (d) Other (describe)]

Destruction date: [INSTITUTION:fillin — UTC date and time]

AUTHORITY

Authorized under: Records schedule [INSTITUTION:fillin — schedule reference], record series [INSTITUTION:fillin — series identifier], disposition action [INSTITUTION:fillin — action]

Legal grounding: [INSTITUTION:fillin — reference to retention-justification.md §N or specific statute]

Hold review: I confirm that no legal hold, examination hold, regulator inquiry, or other preservation
obligation applies to the records destroyed. The hold review was performed on
[INSTITUTION:fillin — date] by [INSTITUTION:fillin — name and role].

VERIFICATION

The destruction was performed by [INSTITUTION:fillin — operator name and role] on
[INSTITUTION:fillin — destruction date].

Post-destruction verification was performed [INSTITUTION:fillin — date] (typically 30 days after
destruction per `retention-justification.md` §5 step 6) confirming the records are no longer
recoverable. Verification result: [INSTITUTION:fillin — pass or fail].

RECORDS-OFFICER ATTESTATION

I, [INSTITUTION:fillin — records officer name], serving as the institution's designated records officer,
certify that the destruction described above was authorized under the institution's records schedule,
that no preservation obligation applies to the destroyed records, and that the destruction was performed
in accordance with the institution's records-management program.

Date: ____________________
Signature: ____________________
[INSTITUTION:fillin — records officer name and title]

WITNESS ATTESTATION

I witnessed the destruction described above and confirm it was performed in accordance with the
institution's records-management program.

Date: ____________________
Signature: ____________________
[INSTITUTION:fillin — witness name and role]

RETENTION

This certificate of destruction is retained for the longer of 7 years or the institution's
records-schedule retention for disposition documentation, under record series RS-7.
```

### 4.2 Integration with the §5 deletion procedure

The deletion procedure at `retention-justification.md` §5 describes the operational steps. The certificate of destruction adds step 7 to that procedure:

| Step | Action | Owner |
|---|---|---|
| 1 | Records-management team identifies the batch of records reaching end of retention | Records officer |
| 2 | Legal team confirms no hold applies | General counsel |
| 3 | Privacy / data-protection team performs the deletion (cryptographic erasure of mapping rows) | Privacy lead |
| 4 | Deletion log records the operator-side fact (operator X, N rows, date Y) | Privacy lead |
| 5 | 30-day post-deletion verification confirms permanent deletion | Privacy lead |
| 6 | Verification outcome recorded in deletion log | Privacy lead |
| 7 | **Records officer issues certificate of destruction (§4.1 template); retains under RS-7** | **Records officer** |

The deletion log captures the operational fact; the certificate captures the authority. Both are retained.

### 4.3 IKM-erasure certificate variant

When the disposition action is IKM destruction (RS-3), the certificate has a slightly different shape — the destruction method names cryptographic erasure of HSM-stored key material, the operator is the HSM administrator, the post-destruction verification confirms the IKM is no longer present in any HSM partition the institution operates. The records officer's signature is the same; the witness is the HSM administrator's peer (from the separation-of-duties roster per spec §10.5).

### 4.4 Certificate-of-destruction log

The records officer maintains a log of certificates of destruction issued during the year. The log is part of the records officer's annual report to the senior officer or the board. The log is reviewed during the §500.17(b) annual-certification cycle (`nydfs-part500-overlay.md` §10.2 row for §500.13).

---

## 5. Symmetric hold-release procedure

The litigation-hold operations described in `operator-guide.md` lines 498-520 and §6 of `retention-justification.md` are nearly complete on the hold-opening side. A hold opens with a written document carrying hold identifier, scope, source, estimated duration, review date, and issuing counsel's signature. What is partial is the hold-release leg — `retention-justification.md` §6 says "when the hold is released, the records-management team resumes the standard quarterly retention audit," which is operationally right but procedurally thin. *(closes Park-Whittaker F-4)*

A hold release is a records-program event analogous to a hold opening: counsel signs the release, the records officer counter-signs, and the records-program log records the release with the same fields as the hold opening. Without a parallel structure on release, an opposing party in later litigation can argue the hold was implicitly released without authorization, or that records were destroyed before authorization to release.

### 5.1 Hold-release document template

```
LITIGATION HOLD RELEASE

Hold reference: [INSTITUTION:fillin — original hold identifier]
Hold opened: [INSTITUTION:fillin — original hold opening date]
Hold scope (as opened): [INSTITUTION:fillin — verbatim from original hold document]

RELEASE DETAILS

Release identifier: [INSTITUTION:fillin — unique release identifier]
Effective date: [INSTITUTION:fillin — UTC date]
Source closing the hold: [INSTITUTION:fillin — settlement reference, court order reference, regulator
  closeout letter reference, examination conclusion reference, or other authority]

Scope released: [INSTITUTION:fillin — full scope released, OR partial release naming specific
  tenant_id / run_id / date-range subsets]

Records affected: [INSTITUTION:fillin — record series and approximate volume returning to standard
  retention treatment]

ISSUING COUNSEL

I confirm the hold identified above is released as of the effective date. The records returning to
standard retention treatment may proceed through the institution's records-schedule disposition
process from this date forward.

Date: ____________________
Signature: ____________________
[INSTITUTION:fillin — issuing counsel name and title]

RECORDS-OFFICER ACKNOWLEDGEMENT

I, [INSTITUTION:fillin — records officer name], acknowledge the hold release and confirm the
institution's records-management program has been updated to reflect the release. The records-program
log entry naming this release has been written.

Date: ____________________
Signature: ____________________
[INSTITUTION:fillin — records officer name and title]

RETENTION

This hold-release document is retained under record series RS-8 for the longer of 7 years or the
underlying matter's resolution plus 7 years.
```

### 5.2 Operational event integration

When a hold is released, the institution's chain-operations team emits a `hold.released` operational event under spec §10.2. The event carries the release identifier, the effective date, the scope released, and the records-officer attribution. The operational event is integrity-bound under the chain itself — the hold release becomes part of the institution's auditable record. The corresponding `hold.opened` event was already part of the operational-event taxonomy; the symmetric pair closes the loop.

### 5.3 Partial release handling

A partial hold release covers the case where the hold's scope is narrowed but not fully closed. The release document names the specific subset returning to standard retention; the original hold continues to apply to the unreleased subset. The records officer's log tracks the partial release as a hold-modification event distinct from a full release.

### 5.4 Multi-hold handling

A record subject to multiple overlapping holds (a single chain entry held under two or more litigation matters) is released only when all applicable holds release. The records officer maintains the hold-stack discipline — the entry returns to standard retention only when the last hold releases. The release document names the specific hold being released; the entry remains held under any other applicable hold until that hold also releases.

---

## 6. Retrieval and presentation — the usability layer

ISO 15489-1:2016 names four characteristics a record must have: authenticity, reliability, integrity, and usability. The chain covers three of these explicitly — authenticity from the per-tenant IKM and per-event HMAC, reliability from the spec §10.4 NTP synchronization and the deterministic verifier output, integrity from the HMAC chain plus daily Merkle plus HSM signature. Usability — the property that says a record must be "located, retrieved, presented, and interpreted" by a competent reviewer over the retention horizon — is the records-program layer this section names. *(closes Park-Whittaker F-6)*

### 6.1 Retrieval — finding the record

The chain captures `(tenant_id, run_id, seq)` triples. The institution maintains an index over the chain that lets the records officer (or counsel, or the customer-dispute desk) find a specific record by customer identifier, by date range, by tenant, or by other operationally-relevant dimensions.

| Index | Purpose | Custodian |
|---|---|---|
| Customer-correlation index | Maps customer identifier to affected `(tenant_id, run_id)` set | Privacy / data-protection lead |
| Date-range index | Maps date range to affected `(tenant_id, seal_date)` set | Chain-operations team |
| Tenant index | Maps tenant identifier to all affected `(run_id, seq)` ranges | Chain-operations team |
| Operational-event index | Maps operational-event class to affected events | SOC team |

The records officer does not maintain the indexes directly but holds the operational responsibility for ensuring they exist, are integrity-bound (or are documented as institution-side lookup tables that are NOT integrity-bound), and have a retrieval SLA the institution can answer FRCP 34 production requests against. The retrieval SLA is the records-program standard "produce within 14 days of request" matching FRCP 34.

### 6.2 The customer-correlation index — load-bearing for production completeness

The customer-correlation index is the institution's lookup from a customer dispute, an examiner inquiry, or a discovery production request to the affected runs. Adversarial litigation context (cross-reference `litigation-support.md` §"Adversarial-balance") establishes the index as a load-bearing artifact for production completeness: chain integrity proves that the entries the institution produces are unaltered, but it says nothing about whether the institution surfaced the right entries.

The records officer's discipline for the customer-correlation index:

| Discipline | Detail |
|---|---|
| Documented procedure | The institution publishes the procedure governing how a dispute reference resolves to a run set |
| Audit log | Every lookup against the index is logged with the operator, the timestamp, and the dispute reference |
| Change history | Changes to the index (additions, corrections, schema changes) are logged with the operator, the timestamp, and the rationale |
| Periodic reconciliation | The records officer's quarterly review reconciles the index against the chain — for a sampled customer, the index's run-set is compared against the chain's run-set for that customer to confirm the index is accurate |

### 6.3 Presentation — making the record interpretable

The chain entry's `gen_ai.completion.text` is captured as the full text or as a SHA-256 hash, depending on the institution's data-classification policy. The presentation layer is what the records officer presents to a reviewer — examiner, customer-dispute desk, opposing counsel, the records officer's own quarterly review.

| Presentation form | Use case |
|---|---|
| NDJSON | Examiner consumption; the chain entry's canonical bytes one record per line |
| PDF working paper | Litigation production; the chain entry rendered in human-readable form with the verifier output |
| Verifier-output PDF | Examiner-and-litigation evidence package; deterministic and re-producible |
| Customer-facing summary | Customer-dispute response; the chain entry's substantive content rendered in plain language by the dispute desk |

The institution's CC8.1 control description names the presentation forms the institution supports and the procedure for producing each. The records officer maintains the procedure; the chain-operations team executes the production.

### 6.4 Hash-only presentation — the institution's data-classification choice

When the institution's data-classification policy excluded substantive prompt or response content from the chain (only the SHA-256 hash is retained), the presentation layer is constrained — the records officer can present the hash but not the content. This is by design when the institution chose hash-only retention; it is a constraint the records officer notes in the production with reference to the institution's data-classification policy at the time of capture.

The records officer preserves the data-classification policy that was in effect at the time of capture as part of RS-7 disposition documentation. When a later request seeks the substantive content, the records officer can produce the policy that explains why the content was excluded. This is the load-bearing record for the institution's defense if a later party argues policy-driven exclusion was inappropriate.

### 6.5 Retrieval SLA and production timing

The institution's retrieval SLA aligns with FRCP 34's response timing — 30 days from request, with the institution's internal target typically 14 days to allow legal review and privilege screening. The records officer maintains the SLA evidence:

| Stage | Target |
|---|---|
| Request received | Day 0 |
| Records officer acknowledges | Day 1 |
| Indexing query produces run set | Day 1-3 |
| Chain entries retrieved and audit paths computed | Day 3-7 |
| Legal review and privilege screening | Day 7-12 |
| Production assembled | Day 12-14 |
| Delivered to requester | Day 14 |

A request that exceeds the SLA produces a documented exception in the records officer's log; chronic exceptions trigger a process review.

---

## 7. Preservation plan — surviving the cryptographic-primitive migration horizon

Federally-supervised banks routinely operate records under hold for 10, 12, even 15 years. The chain ledger itself, even when its privacy-store mapping is deleted at 7 years, continues to exist as an integrity-bound record for as long as the institution chooses to retain it. The cryptographic primitives the chain depends on — Ed25519, HMAC-SHA-256, SHA-256, JCS canonicalization — have known and unknown obsolescence horizons. The preservation plan is the institution's answer to the records-officer question "will this still verify in 2046?" *(closes Park-Whittaker F-5)*

### 7.1 OAIS / ISO 16363 alignment

ISO 14721:2012 (the OAIS reference model) §5.1 and ISO 16363:2012 (trustworthy digital repository audit) §5.4 both require a documented preservation plan that names format-migration triggers, format-migration procedures, and format-migration evidence preservation. The preservation plan below is structured to satisfy these.

### 7.2 Cryptographic-primitive obsolescence monitoring

The institution monitors the relevant standards bodies for deprecation announcements affecting any of the chain's primitives. The monitoring program is not the records officer's direct responsibility (the chief information security officer typically owns it) but the records officer is informed and is the authority on whether a deprecation announcement triggers a migration.

| Body | What the institution monitors |
|---|---|
| NIST | FIPS publications and SP 800-131A deprecation announcements; FIPS 186-5 (Ed25519), FIPS 198-1 (HMAC), FIPS 180-4 (SHA-256), FIPS 204 (ML-DSA), FIPS 205 (SLH-DSA) |
| BSI (Bundesamt für Sicherheit in der Informationstechnik) | TR-02102 cryptographic-primitive recommendations; relevant for the European leg of multi-jurisdiction operation |
| NCSC (National Cyber Security Centre, UK) | Cryptography guidance updates |
| IETF CFRG | Working-group recommendations affecting RFC 8032 (Ed25519), RFC 5869 (HKDF), RFC 6962 (Merkle), RFC 8785 (JCS) |
| ETSI | Quantum-safe cryptography guidance (relevant for the harvest-now-decrypt-later horizon) |

Cross-reference to `cryptographic-agility-roadmap.md` for the technical detail of the migration mechanism.

### 7.3 Migration triggers

The institution names specific events that trigger migration consideration. A trigger does NOT automatically initiate migration — it initiates the records officer's review of whether migration is warranted under the institution's specific risk posture.

| Trigger | Action |
|---|---|
| NIST SP 800-131A deprecation announcement for any chain primitive | Records officer reviews; CISO assesses migration scope; senior officer approves migration plan |
| NIST FIPS publication superseding a primitive currently in use | Same |
| Credible peer-reviewed cryptanalytic result against any primitive | Records officer reviews; CISO assesses urgency; senior officer approves accelerated migration if warranted |
| Quantum-computer milestone affecting the harvest-now-decrypt-later threat model | Cross-reference to `cryptographic-agility-roadmap.md` §"HNDL response"; records officer reviews; senior officer approves dual-algorithm seal mandate trigger |
| HSM vendor end-of-support announcement | Records officer reviews; HSM administrator assesses replacement; senior officer approves replacement plan |
| Institution's specific policy (e.g., 5-year primitive review) | Records officer initiates review; full plan refresh |

### 7.4 Migration procedure

When the records officer authorizes migration, the procedure has five steps:

1. **Migration scoping.** The records officer, the CISO, and the chain-operations lead identify the affected record series and the migration target (the new primitives). The migration scope is documented.
2. **Re-seal under new primitives.** The chain-operations team produces new daily seals under the new primitives covering the in-scope record series. The original seals are preserved as integrity-bound parallel records — the migration does NOT destroy the original artifacts. The new seals are integrity-bound under the new primitives.
3. **Verification under both old and new primitives.** During the transitional period, the verifier validates both seal sets — the original seals under the old primitives, the migrated seals under the new primitives. Spec §7 step 11 dual-algorithm dispatch handles this case for signature primitives; for hash-function migration, cross-reference to `cryptographic-agility-roadmap.md` §"Hash-function agility".
4. **Migration certificate.** The records officer issues a migration certificate that becomes itself a record series under RS-7. The certificate names the in-scope record series, the migration target, the dates of the original sealing and the migrated sealing, the verifier-output evidence demonstrating both sets validate, and the records officer's signature. The certificate is retained for the longer of 7 years or the underlying record series' retention.
5. **Original-artifact retention.** The original seals are retained under their existing record-series treatment. They become a historical record of the chain's integrity at the original sealing time; the migrated seals become the going-forward integrity binding. A future verifier can validate either set independently.

### 7.5 Migration certificate template

```
MIGRATION CERTIFICATE

Institution: [INSTITUTION:fillin — institution legal name]
Records-management program reference: [INSTITUTION:fillin — program reference number]

MIGRATION SCOPE

Original primitives: [INSTITUTION:fillin — e.g., HMAC-SHA-256 + Ed25519 + SHA-256 Merkle]
Migration target: [INSTITUTION:fillin — e.g., HMAC-SHA-256 + Ed25519 + ML-DSA-65 dual-algorithm]

Record series in scope: [INSTITUTION:fillin — list of record series identifiers, e.g., RS-1, RS-2, RS-5]
Tenant identifier(s) in scope: [INSTITUTION:fillin — list of tenant_id values]
Date range covered: [INSTITUTION:fillin — earliest seal-date] through [INSTITUTION:fillin — latest seal-date]

Original sealing dates: [INSTITUTION:fillin — date range of original seals]
Migration sealing date: [INSTITUTION:fillin — date the migrated seals were produced]

MIGRATION TRIGGER

Trigger event: [INSTITUTION:fillin — NIST SP 800-131A announcement, vendor EOL, etc.]
Trigger date: [INSTITUTION:fillin — date of trigger event]
Records-officer review date: [INSTITUTION:fillin — date the review concluded]
Senior-officer approval: [INSTITUTION:fillin — name and date]

VERIFICATION EVIDENCE

Original-set verifier output: [INSTITUTION:fillin — reference to verifier output bundle]
Migrated-set verifier output: [INSTITUTION:fillin — reference to verifier output bundle]
Cross-validation result: [INSTITUTION:fillin — confirmation that both sets validate]

RECORDS-OFFICER ATTESTATION

I, [INSTITUTION:fillin — records officer name], serving as the institution's designated records officer,
certify that the migration described above was authorized under the institution's preservation plan,
that the original artifacts are preserved alongside the migrated artifacts, and that the migration was
performed in accordance with the institution's records-management program.

Date: ____________________
Signature: ____________________
[INSTITUTION:fillin — records officer name and title]

RETENTION

This migration certificate is retained for the longer of 7 years or the underlying record series'
retention, under record series RS-7.
```

### 7.6 Format-migration evidence preservation

ISO 16363:2012 §5.4 expects the institution to preserve evidence of every format migration as part of the trustworthy-digital-repository discipline. The institution's preservation-plan compliance evidence:

- Migration certificates under RS-7 (per §7.5 above)
- Verifier output for both original and migrated sets (RS-4)
- Migration-procedure documentation (the procedure at §7.4 above, retained as a governance document)
- Senior-officer approval records (institution governance archive)
- Original-artifact integrity binding (RS-2 seal records preserved indefinitely or until the institution explicitly destroys them under a documented disposition action)

### 7.7 Long-horizon outlook

Will the chain still verify in 2046? The records officer's answer:

- **2026-2030: yes, with current primitives.** Ed25519 + HMAC-SHA-256 + SHA-256 + JCS are NIST-current and have no known practical attacks affecting v1.0 chains.
- **2030-2035: yes, with monitoring.** The cryptographic-agility roadmap names migration triggers in this window. The institution's preservation plan is exercised at first credible trigger.
- **2035-2046: yes, with documented migration.** At least one full migration cycle has occurred by this horizon; the original artifacts and the migrated artifacts are both retained under RS-7-bound migration certificates. A 2046 verifier validates the original artifacts under the original primitives (which still mathematically work, though may be computationally weaker) and the migrated artifacts under the then-current primitives.

The 2046 verifier may find the 2026-vintage Ed25519 signatures susceptible to a 2046 attacker; the migrated 2030-vintage post-quantum signatures will be the load-bearing integrity binding. The records officer's discipline preserves the historical record alongside the going-forward record so a 2046 reviewer sees both.

---

## 8. FOIA, Privacy Act, and GLBA exemption posture

When a federally-supervised bank produces records to a regulator, the records frequently fall under FOIA exemptions and Privacy Act / GLBA protections. When a regulator subsequently faces a FOIA request for the bank's records, the regulator's FOIA officer determines which exemptions apply and consults with the originating institution. The institution's documented FOIA-exemption posture lets the regulator's FOIA officer act consistently. *(closes Park-Whittaker F-7)*

### 8.1 Applicable FOIA exemptions

5 USC §552(b) lists the FOIA exemptions. Three are typically relevant to chain records:

| Exemption | Statutory text (paraphrased) | Application to chain records |
|---|---|---|
| (b)(4) | Trade secrets and commercial or financial information obtained from a person and privileged or confidential | Chain entries carry confidential commercial information — model identifiers, prompt content, pricing, vendor information, routing decisions. The institution claims (b)(4) for chain records produced to a regulator. |
| (b)(8) | Examination, operating, or condition reports prepared by, on behalf of, or for the use of an agency responsible for the regulation or supervision of financial institutions | Where the regulator's examination working papers incorporate the chain records, those working papers are exempt under (b)(8). |
| (b)(7)(A)-(F) | Records or information compiled for law enforcement purposes | Where chain records are part of a 2703(d) production or other law-enforcement context, the (b)(7) exemptions may apply. |

### 8.2 Privacy Act (5 USC §552a) protections

Chain entries carrying customer identifiers (or identifiers that can be cross-referenced to customers) are subject to Privacy Act protection where they are part of a regulator's system of records. The institution's FOIA-exemption posture states:

- Chain entries containing customer identifiers are protected by the Privacy Act where they are part of the regulator's system of records.
- The regulator's Privacy Act notice (Federal Register publication) governs the disclosure rules.
- The institution's Privacy Act consultation procedure is the same as its general FOIA-consultation procedure.

### 8.3 GLBA Title V Safeguards Rule and §6802 disclosure restrictions

Gramm-Leach-Bliley Title V (15 USC §6801 et seq.) restricts the institution's disclosure of customer nonpublic personal information to nonaffiliated third parties. Chain entries are subject to GLBA where they contain NPI under the GLBA definition. The institution's FOIA-exemption posture:

- Chain entries containing NPI are subject to GLBA disclosure restrictions.
- Disclosure to a regulator under examination authority is permitted under §6802(e)(8) (disclosure to comply with federal, state, or local laws, rules, and other applicable legal requirements).
- Onward disclosure by the regulator is governed by the regulator's own GLBA-compliance procedures and, where applicable, the FOIA-exemption mechanism.

### 8.4 Institution's FOIA-exemption posture statement

The institution's records officer publishes the FOIA-exemption posture as a section of the institution's records-management program. The posture statement names:

- The institution's claimed FOIA exemptions for chain records (typically (b)(4) and (b)(8))
- The institution's Privacy Act and GLBA Title V positions
- The institution's FOIA-consultation contact at the records-officer level
- The procedure for the regulator's FOIA officer to consult before responding to a FOIA request affecting chain records

### 8.5 Regulator FOIA-officer consultation procedure

When a regulator receives a FOIA request that may affect institution chain records:

1. The regulator's FOIA officer identifies the affected records.
2. The FOIA officer notifies the institution's records officer using the institution's published consultation contact.
3. The institution's records officer reviews the affected records and identifies applicable exemptions.
4. The institution's records officer responds to the regulator's FOIA officer within 10 business days (or per the consultation procedure's named SLA).
5. The regulator's FOIA officer applies the exemption analysis and responds to the FOIA requester.

The consultation procedure is part of the institution's records-management program and is reviewed annually.

---

## 9. ISO 15489 four-characteristic mapping

ISO 15489-1:2016 §5.2.2 names four characteristics a record must have. The chain plus this records-management program covers all four.

| Characteristic | What it requires | Chain artifact + records-program scaffolding |
|---|---|---|
| Authenticity | The record is what it purports to be | Per-tenant IKM (spec §10.6); HKDF-bound session key (spec §4.1); per-event HMAC-SHA-256 (spec §4.1); records-officer attribution on operational events |
| Reliability | The record can be relied on to be a complete and accurate representation | Spec §10.4 NTP synchronization; spec §4.2.2 receive-timestamp authoritativeness; deterministic verifier output (spec §7); records-officer's quarterly retention audit (cross-reference `retention-justification.md` §8) |
| Integrity | The record is complete and unaltered | Three-layer construction — per-event HMAC + daily Merkle seal + HSM Ed25519 signature; verifier's twelve-step procedure (spec §7) |
| Usability | The record can be located, retrieved, presented, and interpreted | This template's §6 (retrieval and presentation); the customer-correlation index discipline; the institution's CC8.1 names the presentation forms; the records officer maintains the procedure |

The first three characteristics are delivered by the cryptographic substrate. The fourth — usability — is delivered by this records-management program.

---

## 10. Records-officer annual report

The records officer produces an annual report to the senior officer or the board, parallel to the CISO's §500.04(b) report. The report covers:

| Element | Detail |
|---|---|
| Records-schedule status | The schedule is current; reviewed during the year; any changes are documented |
| Disposition activities | List of certificates of destruction issued during the year (RS-7 log); volumes destroyed; any exceptions |
| Hold openings and releases | List of legal holds opened during the year; list of holds released during the year; any partial releases |
| Migration activities | Any preservation-plan migration triggers reviewed; any migrations executed; migration certificates issued |
| FOIA / Privacy Act / GLBA activities | Number of FOIA-officer consultations during the year; any exemption-determination disputes |
| Records-program metrics | Retrieval SLA performance; quarterly-retention-audit findings; sampling-based reconciliation results |
| Material events | Any material records-program events during the year (major hold openings, migration triggers, regulator inquiries) |
| Records-program forward plan | Any planned changes to the records schedule, the preservation plan, or the records-program operating model |

The annual report is itself a record series — typically retained under RS-7 disposition documentation for the longer of 7 years or the institution's records-schedule retention for governance documents.

---

## 11. Cross-references

- `regulator-pack/retention-justification.md` — the legal-grounding companion to this records-management program; the "longer of" rule, the proportionality test, the §5 deletion procedure (this template's §4 adds step 7 to that procedure).
- `regulator-pack/nydfs-part500-overlay.md` — the §500.13 limitations-on-data-retention mapping uses this records-management program as the institution-side artifact; the §500.17(b) annual-certification evidence map's §500.13 row points here.
- `templates/fre-902-certification.md` — the records-custodian declaration template; the records officer is one possible declarant.
- `cryptographic-agility-roadmap.md` — the technical detail behind §7 (preservation plan); the migration mechanism the records officer authorizes.
- `litigation-support.md` — the customer-correlation index discipline at §6.2 of this template is load-bearing for production completeness; cross-reference to the adversarial-balance section.
- `operator-guide.md` — the operational context the records officer fits into; the chain-operations team, privacy lead, HSM administrator, and legal team are named there.
- `audit-procedures.md` — the audit procedures (P-1 through P-30) that test the institution's records-management program through chain-integrity testing.
- `customer-dispute-procedures.md` — the customer-dispute desk uses the customer-correlation index this template names at §6.2.

---

## 12. Implementation checklist

For an institution adopting this records-management program for the first time:

- [ ] Designate the records officer in writing (§2.3 above) with senior-officer or board approval
- [ ] Add the records-officer role to the operating model in `operator-guide.md`
- [ ] Publish the records schedule (§3 above) listing the eight record series the chain produces
- [ ] Add the certificate-of-destruction template (§4.1) to the institution's templates directory
- [ ] Update the deletion procedure at `retention-justification.md` §5 to include step 7 (records-officer issues certificate of destruction)
- [ ] Add the hold-release document template (§5.1) to the institution's templates directory
- [ ] Add the `hold.released` operational event to the institution's chain-operations schema
- [ ] Document the customer-correlation index disciplines (§6.2)
- [ ] Document the retrieval SLA (§6.5)
- [ ] Publish the preservation plan (§7) including the cryptographic-primitive obsolescence-monitoring procedure, the migration triggers, the migration procedure, and the migration-certificate template
- [ ] Publish the FOIA-exemption posture (§8.4) and the regulator FOIA-officer consultation procedure (§8.5)
- [ ] Add the records-officer annual report to the institution's senior-officer reporting cycle
- [ ] Update the §500.17(b) annual-certification evidence map (`nydfs-part500-overlay.md` §10.2) to name the records officer in the §500.13 row

The implementation closes inside 90 days of records-officer designation. The chain is close enough to a complete records program that the close-out is documentation work, not engineering work.
