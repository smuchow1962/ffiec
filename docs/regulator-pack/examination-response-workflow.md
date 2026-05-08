# Examination response workflow

> **What this doc is.** The single workflow doc that walks an examiner from a `chain.verification_failure` JSON failure record to a documented examination response. ~60 seconds of reading.

## The five-step path

When the verifier produces a failure, you have a JSON record like this in `--json-report` output:

```json
{
  "step": 8,
  "reason": "key_fingerprint mismatch at seq 17: looked-up IKM does not match the entry's recorded fingerprint",
  "run_id": "r_a3f29b71c",
  "seq": 17,
  "tenant_id": "tenant_acme_prod",
  "key_version": 1,
  "expected_fingerprint": "b94c1a77b40bf5106c66ca6c1c1b4989",
  "recorded_fingerprint": "2eede65f0f764c97eaf3b3f306a48537"
}
```

Walk it through these five steps:

### 1. Map `step` → severity and finding paragraph

`step` is the spec §7 procedure step. Look it up in `regulator-pack/finding-language.md` severity table to get the severity guidance and the candidate finding paragraph. The table is keyed by `step`; one lookup, one row.

For the example above (`step: 8`), the row says:

> Severe MRA — the load-bearing new failure mode. Investigated against the IKM roster and §10.1 reconciliation log, NOT against chain content.

The matching finding paragraph in the same doc is "Key fingerprint mismatch (spec §7 step 8)."

### 2. Map `step` → IR scenario

The institution's IR playbook (`docs/incident-response-playbook.md`) has a scenario per chain-detected failure mode. Map `step` to scenario:

| `step` | IR scenario |
|---|---|
| 1 | (Not institutional — verifier-version skew, no IR) |
| 2-3 | Scenario 1 (chain hash mismatch / construction defect) |
| 4 | Scenario 1 + evidence-handling investigation (mis-bundled snapshot OR lift attempt) |
| 5 | Scenario 1 (SDK version drift) |
| 6 | Scenario 1 (chain hash mismatch) |
| 7 | Scenario 8 (unknown key_version) |
| 8 | Scenario 7 (key_fingerprint mismatch) |
| 9 | Scenario 1 (chain hash mismatch) |
| 10 | Scenario 2 (Merkle root mismatch) |
| 11 (single-algorithm) | Scenario 3 (signature verification failed) |
| 11 (dual-algo: partial-coverage seal) | Scenario 5 variant (control-completeness — institution's declared dual-algorithm posture is incomplete on this seal-day; not chain-integrity) |
| 11 (dual-algo: algorithm not on declared posture list) | Scenario 5 variant (control-description-accuracy — posture documentation and SDK configuration disagree) |
| 11 (dual-algo: co-signed seal failure, one algorithm valid + one invalid) | Scenario 4 (master-key compromise on the failed algorithm) + Scenario 3 (signature failure path); regulator coordination on migration timeline if attributable to a published algorithm break |
| 12 (cadence) | Scenario 5 (sealing delay) variant |
| 12 (dev-mode) | Scenario 6 (software-key fallback in production) |
| (file pre-flight truncation) | Scenario 9 (audit file truncation detected) |

For the example (`step: 8`), the IR scenario is Scenario 7. Direct the institution to that scenario.

### 3. Pull the institution's reconciliation evidence

For `step: 7` and `step: 8` failures, the institution's most recent `master.reconciliation_completed` operational event is the first piece of evidence to ask for. The event carries:

```
period, key_versions_observed, key_fingerprints_observed, fingerprint_unmatched_count
```

Two questions for the institution:
- "What does your most recent `master.reconciliation_completed` say about `(tenant_id=T, key_version=V, key_fingerprint=F)`?"
- "Is `recorded_fingerprint` from the failure record present in `key_fingerprints_observed`?"

If yes: the reconciliation already flagged it; ask for the institution's investigation log.
If no: the reconciliation hasn't yet run on the affected period; ask the institution to run the reconciliation against the evidence period and produce the result.

For `step: 9` (payload_hash MAC mismatch), the relevant evidence is the institution's `chain.verification_failure` operational event for the affected `(run_id, seq)`; ask for the matching event from the institution's log store.

For `step: 10` (Merkle root mismatch), the relevant evidence is the institution's seal-job log for the affected `seal_date` (`seal.job_started`, `seal.job_completed`, `seal.job_failed` events).

For `step: 11` (signature verification failed), the relevant evidence is the institution's `hsm.operation_*` events for the seal-signing operation, plus the institution's `master_key.rotated` events for the period.

### 4. Apply the severity / finding paragraph from `finding-language.md`

Lift the matching paragraph from `regulator-pack/finding-language.md` into the examination report; edit for institution-specific facts (substitute the actual `tenant_id`, `key_version`, `seq`, dates, counts).

For `key_fingerprint mismatch` specifically, the finding paragraph emphasises that this is an **identity** mismatch (investigated against the IKM roster), not a **content** tampering. Keep that distinction in the report; it is what makes the rework's defensive primitive examinable cleanly.

### 5. Confirm the institution's response

The institution's response to a chain-detected finding should include:

- The root cause (one of the IR scenario's documented root causes).
- The remediation (per the IR scenario).
- The post-remediation re-verification (the verifier's subsequent run on the institution's corrected state).
- The control update (if the root cause indicates a control gap, the institution updates its control description).

For `step: 8` `key_fingerprint mismatch` specifically, the institution's response should include:
- Identification of the IKM-roster row that was wrong (typically `(tenant_id, key_version)`).
- The corrective action (typically restoring the correct IKM bytes from KMS history or backup).
- The re-run of the verifier against the corrected roster, demonstrating PASS on the affected period.
- A change-management record for the IKM-roster correction.

If the institution cannot produce these four items, the finding remains open until they are provided.

## Second worked example — `step: 10` Merkle root mismatch

For comparison: a structurally different failure type. The JSON record:

```json
{
  "step": 10,
  "reason": "merkle root mismatch — ledger contents do not produce sealed root",
  "tenant_id": "tenant_acme_prod",
  "seal_date": "2026-04-15",
  "computed_root": "5e1ac73fa9b8c7d6e5f4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2",
  "recorded_root": "9f8e7d6c5b4a3210fedcba9876543210fedcba9876543210fedcba98765432ff"
}
```

Walk the same five-step path:

1. **Map `step` → severity.** Look up `step: 10` in `regulator-pack/finding-language.md` severity table. Row says: "Severe MRA; the day's events do not produce the sealed root. Possible server-side history rewrite." Matching finding paragraph: "Merkle root mismatch (single day)."
2. **Map `step` → IR scenario.** Per the table in step 2 above, `step: 10` maps to Scenario 2 (Merkle root mismatch at verification). Direct the institution to Scenario 2.
3. **Pull the institution's seal-job evidence.** For `step: 10` failures, the relevant evidence is the institution's seal-job log for the affected `seal_date`: `seal.job_started`, `seal.job_completed`, `seal.job_failed` events from 2026-04-15 (the seal day). Two questions for the institution:
   - "What does your `seal.job_completed` event for 2026-04-15 say about `merkle_root`? Does it match the verifier's `recorded_root` (5e1a...) or `computed_root` (9f8e...)?"
   - "Has the ledger been recovered from a backup since the seal was originally produced? If so, was the backup from before the affected day or after?"

   The verifier's recorded_root is what the seal record contains; the verifier's computed_root is what the ledger contents now produce. The mismatch direction tells the auditor whether the seal record was altered (computed correct, recorded wrong) or the ledger contents were altered or recovered from a backup (recorded correct, computed wrong). Most often the latter — Scenario 2's "recovery from older backup without gap-fill" branch.
4. **Apply the severity / finding paragraph.** Lift the "Merkle root mismatch (single day)" paragraph from `finding-language.md` and edit for institution-specific facts (substitute the actual `tenant_id`, date, root hashes).
5. **Confirm the institution's response.** Same four-item checklist (root cause, remediation, post-remediation re-verification, control update). For Scenario 2 specifically, the post-remediation re-verification means the institution replays the gap from any available source (SDK buffers, downstream OTLP) and re-runs the seal job for the affected day; the verifier's next pass should produce a matching recomputed root.

The workflow generalizes across structurally different failure types: identity mismatch (step 8), content tampering (step 9), Merkle ledger-content mismatch (step 10), signature compromise (step 11). Each maps to a specific IR scenario, a specific evidence path, a specific finding paragraph.

## When this workflow does NOT apply

- **`step: 1` failures** are verifier-version skew, not institutional findings. Obtain the matching verifier from the project supply chain (per `regulator-pack/deployment-package.md`) and re-run.
- **PASS-with-anomaly results** (late-binding events, sealing delays under threshold, clock-skew, `master_key_rotation_observed` for documented rotations) are operational signals, not findings. They land in the examination report under "Anomalies noted" with a sentence each.
- **PASS results** require no examination action beyond noting the result.

## Where this workflow comes from

The five-step path mirrors the spec §7 verifier procedure and the §10.2 operational-event taxonomy. The mapping in step 2 lifts directly from `incident-response-playbook.md` scenario list. The reconciliation evidence in step 3 lifts from `audit-procedures.md` P-6 (key-fingerprint reconciliation). The severity table in step 1 lifts from `finding-language.md`.

This doc exists so the examiner does not have to assemble the workflow from four other docs at the moment of finding evaluation.
