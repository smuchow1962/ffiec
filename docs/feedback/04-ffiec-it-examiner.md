# Round 9 — FFIEC IT Examiner review

> **Persona.** Robert Chen. Senior FFIEC IT Examiner with the OCC. Fifteen years examining bank IT controls, including a recent rotation through the FFIEC's Cybersecurity Specialist program. First-look reviewer who runs the verifier on examiner-laptop in real engagements.
>
> **Reading angle.** The examiner workflow — quickstart, training, sample report, finding language, deployment package. Will the examiner be able to run a verification, read the report, and write defensible findings?
>
> **Scope read.** `spec/chain-of-custody-v1.md` (§7, §13); `docs/examiner-quickstart.md`; `docs/regulator-pack/sample-report.md`; `docs/regulator-pack/finding-language.md`; `docs/regulator-pack/handbook-mapping.md`; `docs/regulator-pack/deployment-package.md`; `docs/regulator-pack/examiner-training.md`; `docs/regulator-pack/examiner-approval-template.md`; `docs/portfolio-comparison-procedures.md`; `docs/design/07-verifier-design.md`. Did not read `docs/feedback/historical-round-8-pre-hmac-rework/`.

## Headline

The package reads like a real examiner package, which I do not say about most regulator-facing material. The quickstart gets me from cold to "I have a defensible report" in five minutes. The deployment package answers what my IT shop will ask me before I get a binary on the laptop. The sample report shows me what a passing month looks like. The 30-minute training session lands on the right surface area for a brand-new examiner. I would adopt this package.

The structural problems I have are concentrated and fixable. **The examiner-facing docs (quickstart, sample report, finding language, training, handbook mapping) have not been updated to reflect the new failure modes that the spec §7 verification procedure produces.** The verifier ships with twelve named failure points; the consumer-facing docs still describe the four legacy ones. That mismatch is the issue that would bite me in the field — I would run the verifier, get back a `key_fingerprint mismatch at seq N` failure, and find no paragraph in `finding-language.md` to lift into my workpapers.

The other concentrated issue is rotation. The spec carries `key_version` and `key_fingerprint` per entry, the JSON report surfaces `key_versions_observed`, and the design doc shows a `master_key_rotation_observed` anomaly — but the sample examiner-facing report does not include a rotation day, and `finding-language.md` does not address the rotation case as a normal-operations event the examiner should expect to see.

I'd ship this package after closing the five questions below.

---

## Question 1 — Does the quickstart accurately describe the new failure modes the verifier produces?

**Status: Gap.**

`docs/examiner-quickstart.md` §"Common failure modes" is a six-row table that lists `Chain hash mismatch`, `Chain link broken`, `Merkle root mismatch`, `Signature verification failed`, sealing delay > 72h, and software-key in production. Spec §7 defines twelve verification steps, each with a specific named failure. The new failure modes that produce reasons the examiner has never seen before:

| Verifier reason | Spec §7 step | In quickstart? |
|---|---|---|
| `format_version not supported by this verifier (running v1)` | 1 | No |
| `header HKDF inputs do not match running v1 inputs` | 2 | No |
| `header genesis_hash does not match v1 constant` | 3 | No |
| `cross-chain lift detected at seq N` | 4 | No |
| `format_version mismatch at seq N` | 5 | No |
| `unknown key_version: no IKM for (tenant=T, key_version=V) at seq N` | 7 | No |
| `key_fingerprint mismatch at seq N` | 8 | No |
| `audit file ends mid-line — possible mid-write crash` | §4.1 truncation rule | No |
| `cadence mismatch` | 12 | No |
| `dev-mode seal in production verification — refused` | 12 | Partial — the row says "Software-key in production" but does not name the verifier reason string |

That is eight new reasons the examiner can encounter and the quickstart prepares them for none of them. The `Chain hash mismatch` row in the existing table maps to step 9 (`payload_hash MAC mismatch`) — the wording in the table no longer matches what the verifier prints, which means the examiner cannot grep `finding-language.md` for the verifier's actual reason string and find a hit.

The fix is mechanical. Replace the six-row table with a twelve-row table that uses the verifier's actual reason strings (one per spec §7 step) and the severity guidance from `finding-language.md`. Keep it on one page — the quickstart's value is that I can read it in five minutes.

**Two specific points the quickstart should call out by name:**

- **`key_fingerprint mismatch` is the new "stop and call the bank" finding.** It indicates the institution looked up an IKM that does not produce the recorded fingerprint — botched rotation, cross-tenant configuration drift, or a swapped backup. The examiner needs one sentence telling them this is severity-equivalent to chain-hash mismatch and is not a verifier bug.
- **`audit file ends mid-line` is not a tampering finding.** It is a writer-side mid-write crash (the SDK process died between events). The examiner needs to know this is severity-Medium operational, not severity-High integrity, so they don't escalate it as a tampering case when it is a crash-recovery case. The current quickstart does not warn the examiner about this distinction and a green examiner would treat it as a severe finding.

---

## Question 2 — Does the sample report cover the rotation case and the new anomaly types?

**Status: Partial.**

`docs/regulator-pack/sample-report.md` shows a healthy 30-day month with three sealing-delay anomalies and no rotation activity. The per-day detail blocks show no `key_versions` field, no `hkdf_inputs_digest`, no `dev_mode` field, and no `kms_handle_uri` — but `docs/design/07-verifier-design.md` §6 shows the JSON report carries `key_versions_observed`, `key_fingerprints_observed`, per-day `key_versions`, per-day `hkdf_inputs_digest`, per-day `dev_mode`, and a `master_key_rotation_observed` anomaly kind. The PDF structure described in §5.1 of the verifier-design doc does not reflect those new fields either, which means the sample-report.md is consistent with the verifier-design doc, but both are out of date relative to the spec.

What the sample report should show:

- **A rotation day.** At least one day in the 30-day excerpt should show `key_versions: [3, 4]`, the `master_key_rotation_observed` anomaly, and language explaining that the institution rotated mid-day per its incident log. The examiner needs to see that rotation is normal-operations PASS-with-anomaly, not a failure. As written, the first time an examiner sees a rotation day in production they will not know whether to escalate.
- **The per-day `key_fingerprint(s) observed` count.** When the institution operates one IKM per tenant, the count is 1 per day. When the institution rotated, the count is 2 (the old generation + the new). When the count is unexpected (3+ on a non-rotation day), that is an investigation trigger. The sample report needs to show what the normal value looks like so the examiner has a baseline.
- **A `kms_handle_uri` line per day.** The examiner should see `kms_handle_uri: aws-kms:arn:...` in production samples so they know what to look for. The presence of `plaintext-` prefix on a production day is a Severe finding under §10.7; the examiner needs the visual training to spot it.
- **A failure-section example.** The sample report shows only PASS-with-anomalies. The examiner has no model for what a FAIL day looks like in the report — the per-day block shape, the failure-record structure (with `step`, `reason`, `seq`, `expected_fingerprint`, `recorded_fingerprint`), and how the cover-page summary changes. A second short sample report showing one failed day would close this.

The methodology section on page N+2 still says "Chain hash" and "HMAC chain walk: PASS" in the old vocabulary. That should be updated to reference the spec §7 twelve-step procedure by step number, so the examiner reading the report can map the verifier's `step: 8` failure-record field to the methodology section.

---

## Question 3 — Does finding-language cover the new failure modes for the examiner's report?

**Status: Gap.**

This is the question that would bite me in the field. `docs/regulator-pack/finding-language.md` has the severity table and seven sample finding paragraphs. The severity table has nine rows but maps to the legacy verifier vocabulary (`Chain hash mismatch`, `Chain link broken`, etc.) — the verifier reason strings the examiner will actually encounter are not in the table, so the examiner cannot grep their verifier output against this doc and find a matching row.

The missing paragraphs:

- **`key_fingerprint mismatch`.** This is the load-bearing new failure mode (per spec §4.1 inviolate property 3). It indicates the institution's IKM lookup returned bytes that do not match the per-entry fingerprint — botched rotation, swapped tenant rows in the registry, restored backup pointed at the wrong tenant, or active cross-tenant configuration drift. The finding paragraph should distinguish it from `payload_hash MAC mismatch` (which is content tampering) — `key_fingerprint mismatch` is **identity** mismatch and is investigated against the institution's IKM roster and reconciliation log (§10.1), not against the chain content. Severity guidance: MRA, comparable to chain hash mismatch.

- **`unknown key_version`.** The verifier looked up `(tenant_id, key_version)` and the IKM registry returned null. Cause set: institution rotated and did not register the new generation in the verifier's IKM registry; the institution withdrew an old IKM from the registry while events stamped with that `key_version` still exist; the verifier's IKM registry is stale relative to the institution's actual operations. This is an institutional operations finding — the institution must keep its IKM roster complete for the retention window of the events. Severity guidance: MRA.

- **`format_version not supported by this verifier`.** The institution upgraded its SDK to a future spec version and the examiner is running an older verifier. This is not a finding against the institution; it is a verifier-version mismatch the examiner resolves by obtaining the newer verifier from the project. The finding-language doc should explicitly say "this is not an institutional finding" so the examiner does not write up the bank for a verifier-version skew.

- **`audit file ends mid-line — possible mid-write crash`.** The SDK process died mid-append. The chain content is intact through the last complete entry; the verifier refuses to verify a truncated tail because permissive verification would silently lose events. The finding paragraph should classify this as Operational (sealing-delay-equivalent severity) and direct the institution to its incident log and crash-recovery procedure. This is not a tampering finding.

- **`cross-chain lift detected at seq N`.** An event in the file claims a different `tenant_id` or `run_id` than the file header. This is severe — either the institution mis-bundled snapshots (operational) or an attacker attempted to lift an event from one chain into another (integrity). The finding paragraph should require the institution to produce evidence of the event's original chain provenance.

- **`cadence mismatch`.** The seal record's recorded cadence does not match the institution's claimed cadence in its control description. This is a control-description-accuracy finding — the institution's documentation does not match its operations. Severity guidance: Observation, escalating to MRA on repeat.

The severity table at the top of the doc should be replaced with a twelve-row table keyed by the spec §7 step number, so the examiner can map verifier output (`{"step": 8, "reason": "key_fingerprint mismatch ..."}`) directly to a row.

---

## Question 4 — Does the IT Handbook mapping reflect the new control surface around per-entry IKM identity binding?

**Status: Partial.**

`docs/regulator-pack/handbook-mapping.md` correctly maps the chain to II.C.10 (Logging) as the headline mapping and II.C.13 (Cryptographic controls) for HMAC + Ed25519. But the new per-entry IKM identity binding — `key_version`, `key_fingerprint`, `hkdf_inputs_digest`, `kms_handle_uri` stamped on every entry, plus the §10.1 weekly fingerprint reconciliation — is a substantively new control surface that the mapping does not address.

What the IT Handbook mapping should add:

- **II.C.13 (Cryptographic controls) coverage of key identity provenance.** The per-entry `key_fingerprint` is a public identity binding the verifier checks before any MAC compute (spec §7 step 8). This is a cryptographic control beyond "HMAC + Ed25519 in HSM custody" — it is an integrity proof that the IKM in the registry matches the IKM that signed the event. The Handbook mapping should say so by name, because II.C.13 examiners will look for cryptographic controls on key identity, not just key custody.

- **II.C.10 (Logging) coverage of weekly reconciliation.** Spec §10.1 requires institutions to reconcile every observed `(tenant_id, key_version, key_fingerprint)` triple against the IKM roster at no greater than weekly cadence. That is a logging-integrity control the II.C.10 examiner will want to see evidence of (the `master.reconciliation_completed` operational event, the institution's roster, the unmatched-count metric). The Handbook mapping does not mention §10.1 today.

- **II.C.13 coverage of IKM minimum length and software-adapter exclusion.** Spec §10.6 (32-byte IKM minimum, RFC 4868) and §10.7 (compile-time exclusion of the software-key adapter from production builds) are cryptographic-controls findings the II.C.13 examiner needs to know about. Today the Handbook mapping table mentions Ed25519 in HSM custody but not the IKM-side floor or the dev-adapter exclusion rule.

- **II.E (Change management) coverage of `format_version`.** The spec carries `format_version` per entry and binds it into the seal's signed payload (spec §4.3). When the institution upgrades SDKs across a `format_version` boundary, that is a configuration change with an integrity contract — the Handbook mapping should mention it under II.E because change-management examiners will ask "how do you know your verifier version is compatible with your SDK version" and the mapping should point them at the `format_version` field.

The existing II.C.5 (Logical access) mention of the session-key handshake is fine. The gap is that II.C.13 and II.C.10 entries do not yet reflect the per-entry identity-binding control surface.

---

## Question 5 — Does the examiner have a clear path from a `key_fingerprint` mismatch finding to a documented examination response?

**Status: Gap.**

This is the workflow question. The examiner gets a `key_fingerprint mismatch at seq N` failure in the verifier's JSON output. What do they do next? The package today gives the examiner three pieces — the spec §7 normative description, the verifier-design §4.1 procedure step, and the JSON failure-record shape. None of these are examination-workflow documents.

The missing documented path:

1. **`finding-language.md` should have a `key_fingerprint mismatch` paragraph** (covered in Q3 above) that names the failure, sets severity, and points to the IR playbook entry.

2. **The IR playbook should have a "fingerprint-mismatch" scenario.** The handbook-mapping doc references `docs/incident-response-playbook.md` but I did not see a fingerprint-mismatch scenario called out by name in the read scope. The examiner should be able to direct the institution to a specific IR scenario when the finding lands. (Note: I cannot confirm the IR playbook contents without reading that doc; this question may be answerable inside that doc — but the examiner-facing docs do not signpost the path.)

3. **The sample report should show the institution's response language.** When the institution receives a `key_fingerprint mismatch` finding, what does their response document look like? The sample report shows the examiner's side; a paired institution-response sample would close the loop and let the examiner see the full back-and-forth.

4. **The reconciliation operational event (§10.1) should be a documented input.** When the examiner gets a fingerprint-mismatch finding, the first question they ask the institution is "what does your most recent `master.reconciliation_completed` event say about this `(tenant_id, key_version, key_fingerprint)` triple?" That workflow step should be in the finding-language paragraph or in a short examination-procedure note.

The examiner needs about 60 seconds of doc to know what to do when they see this failure. Right now they would need to read the spec, the verifier design, the operational-events catalog, and the IR playbook to assemble the response. That is not a workflow.

---

## Smaller observations

**Quickstart "What you need" item 4 mentions `the regulator already has the fingerprint on file`** — but the rest of the package does not document how the regulator obtains and stores the public-key fingerprint. The deployment package §"Trust anchors" mentions trust anchors are cached out-of-band; the same paragraph should mention tenant-public-key fingerprints. A green examiner reading the quickstart will not know where their copy of the bank's public key fingerprint came from.

**Sample report "anomalies" section uses the operational story format very well** — the 2026-04-15 entries are exactly the right level of detail for a working paper. Keep this format when adding the rotation-day and failure-day examples.

**Training Module 5 "Common patterns" table is good** but uses the legacy vocabulary. Update in lockstep with the quickstart.

**Examiner-approval template is solid.** The reversal commitment in §7 is the right shape (any chain-verification failure during the relaxation period reverses the relaxation). I have no findings here.

**Portfolio comparison procedures** are appropriate for the manual-process state. The note that v1.1 `verifier portfolio` will scale beyond ~50 institutions is honest. I have no findings here.

**Deployment package "Operational behavior" section** ("No telemetry. No network access. Read-only ledger access.") is exactly what my IT shop will want in writing. Keep.

---

## Per-role roll-up

| Question | Status |
|---|---|
| Q1. Quickstart accurately describes new failure modes | Gap |
| Q2. Sample report covers rotation case and new anomaly types | Partial |
| Q3. Finding-language covers new failure modes | Gap |
| Q4. Handbook mapping reflects per-entry IKM identity binding | Partial |
| Q5. Documented examination-response path from `key_fingerprint` mismatch | Gap |

**Adoption-readiness call.** Approve with conditions. The package is structurally sound and the workflow shape is right. Three of five gaps are concentrated in one mechanical change — replacing the legacy four-failure vocabulary with the spec §7 twelve-step vocabulary across quickstart, sample report, finding language, training, and handbook mapping. The other two gaps (rotation case in sample report; documented response path for fingerprint mismatch) require small additions to the same documents. None of these are deep design issues. Estimated effort: one focused editing pass across six docs, plus a paired sample-report addendum showing a rotation day and a failed day.

I would not block on the v1.1 `verifier portfolio` tooling. Manual portfolio comparison is acceptable at present scale.
