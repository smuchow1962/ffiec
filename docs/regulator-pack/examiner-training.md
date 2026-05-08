# Examiner training — chain verification

> **What this doc is.** Onboarding a new examiner to chain-of-custody verification. From "I have a ledger snapshot and a public key" to "I have a defensible report" in 30 minutes.

## Prerequisites

Before this training:

- The examiner has read [`examiner-quickstart.md`](../examiner-quickstart.md) (5-minute orientation)
- The examiner has the verifier binary on their machine (or knows where to obtain it)
- The examiner has cosign and a copy of the project's cosign public key

If any of these are missing, start with the quickstart and the deployment package.

## Training session structure (30 minutes)

### Module 1 — Concepts (5 minutes)

The chain has three integrity layers:

1. **HMAC chain at capture.** Each event carries a hash bound to the previous event in its run plus a per-process key. An external attacker cannot forge an event without the key.
2. **Daily Merkle seal.** Each tenant-day's events feed an RFC 6962 Merkle tree; the root is HSM-signed. Modifying any past event changes the root, which the verifier catches.
3. **Ed25519 signature in HSM custody.** The signing key is in a FIPS 140-2 L3 device; even a fully compromised application stack cannot forge a daily root.

The verifier runs all three checks per day in the examination range and produces a pass/fail per day plus an overall pass/fail.

### Module 2 — Inputs the examiner needs (5 minutes)

| Input | Source | Notes |
|---|---|---|
| Ledger snapshot | The institution; format is implementation-flexible | Bring on a USB or read-only mount |
| Tenant public key | The institution's tenant key registry | A PEM-encoded Ed25519 public key |
| Tenant ID | The institution's identifier | Matches the ledger and the public key |
| Verifier binary | Project releases (verified via cosign) | Same binary every time; trust anchor cached |
| Date range | The examination scope | Optional; defaults to the snapshot's full range |

The examiner brings these to a dedicated audit machine. The audit machine has the verifier and `verifier-validate.sh` allowlisted.

### Module 3 — Run the verifier (10 minutes — hands-on)

Step by step:

```
# 1. Validate the verifier binary itself
./verifier-validate.sh ./verifier
# Expected: cosign PASS + GPG manifest PASS

# 2. Inspect the ledger snapshot's hash (record for working paper)
sha256sum ./ledger-snapshot.dump

# 3. Inspect the public key (record for working paper)
sha256sum ./tenant-public.pem

# 4. Run the verification
./verifier verify \
  --ledger ./ledger-snapshot.dump \
  --root-key ./tenant-public.pem \
  --master-key ./tenant-ikm.bin \
  --tenant-id <tenant_id> \
  --from <from-date> \
  --to <to-date> \
  --report ./report.pdf \
  --json-report ./report.json \
  --bundle ./working-paper-bundle.tar.gz

# 5. Check exit status
echo $?
# 0 = pass; 1 = at least one day failed
```

The `--master-key` flag points at the institution's IKM file (32 raw bytes; file mode 0600 on POSIX, ACL restricted to the examiner account on Windows). Without it the verifier performs structural verification only and skips per-event HMAC equality (spec §7 fail-closed degradation under `--master-key absent`); under `--strict` the absent IKM elevates to FAIL. The institution provides the IKM at examination time per a documented disclosure shape (typically protective-order disclosure or HSM-mediated derivation per `legal-disclosure.md` §"Court-ordered master-key disclosure" and `customer-dispute-procedures.md` §"IKM access for customer-side verification").

The bundle is the working-paper artifact. Its SHA-256 goes into IT-EX with the examination metadata.

### Module 4 — Read the report (5 minutes)

A passing report has:

- **Cover page** — overall result, summary counts
- **Per-day detail** — each day in range with PASS/FAIL plus per-check status
- **Anomalies section** — non-fatal observations (late-binding, sealing delays, clock skew)
- **Methodology section** — algorithms and rules used (always identical for the same spec version)

A failing report has the same structure plus a **Failures section** detailing which days failed and why.

The examiner scans the cover page first, reviews any failures, then reviews anomalies for context.

### Module 5 — Common patterns (5 minutes)

The verifier follows the spec §7 twelve-step procedure. Each step's failure produces a specific named reason. Map the JSON failure record's `step` field to the table below; the full taxonomy is in `examiner-quickstart.md` and `finding-language.md`.

| Pattern (JSON `step` + `reason`) | What it means | Examiner action |
|---|---|---|
| All pass, no anomalies | Healthy chain | Note the result; close as routine |
| All pass with `master_key_rotation_observed` anomaly (documented in incident log) | Normal-operations IKM rotation | Note as PASS-with-anomaly; no finding |
| All pass with sealing-delay anomalies | HSM availability issue | Confirm institution's incident log explains |
| All pass with late-binding above threshold | Operational/network issue | Note for institution's next control review |
| `step: 8` `key_fingerprint mismatch` | **The new "stop and call the bank" finding.** Identity mismatch (botched rotation, restored backup, cross-tenant drift), NOT content tampering | Severe; investigate against IKM roster (NOT chain content); IR Scenario 7 |
| `step: 9` `payload_hash MAC mismatch` | Content tampering or SDK defect | Severe; require institution's root-cause analysis; IR Scenario 1 |
| `step: 10` `merkle root mismatch` | Server-side history-rewrite signal | Severe; investigate immediately; IR Scenario 2 |
| `step: 11` `signature verification failed` | Possible HSM key compromise OR registry mismatch | Severe; require key rotation and notification; IR Scenario 3 |
| `step: 11 (dual-algo)` `co-signed seal failure` | One algorithm validated, one didn't (transitional period) | Severe; coordinate with regulator on migration timeline; IR Scenario 4 + Scenario 3 |
| `step: 12` `dev-mode seal in production verification — refused` | Misconfiguration in production (compile-time exclusion bypassed) | Severe; require remediation; IR Scenario 6 |
| `step: 7` `unknown key_version` | Institution's IKM roster missing a referenced generation | High; investigate IKM retention or provisioning; IR Scenario 8 |
| `(file pre-flight)` `audit file ends mid-line` | Writer-side mid-write crash (NOT tampering) | Operational; recovery from SDK-local SQLite or upstream OTLP; IR Scenario 9 |

| `step: 12a` `gen_ai_model_identifier_missing` | Chain entry represents a model call but lacks `gen_ai.request.model` or `gen_ai.response.model`. Control-completeness for SR 11-7 reproducibility, NOT chain-integrity. | Medium under non-strict; FAIL under `--strict`. Investigate the institution's SDK-instrumentation for AI model calls; MRM-program finding (audit-procedures P-25). |

Steps 1-6 (header pre-flight + cross-chain lift + structural walk) are the format-construction and structural defenses; their failure modes are listed in `examiner-quickstart.md` and `finding-language.md` with IR scenario mapping. The verifier output is precise. The examiner translates the precision into report language using `finding-language.md` and follows the five-step path in `examination-response-workflow.md`.

## Verifier-operator qualification standard

The 30-minute training above gets a new examiner from "I have a ledger snapshot and a public key" to "I have a defensible report." That is the operational qualification. The evidentiary qualification is harder. When a chain entry becomes evidence in a contested proceeding — a federal trial, an enforcement hearing, a litigation deposition — the operator who ran the verifier may be called as an expert witness under FRE 702. The institution's IT staff are not cryptographers by trade, and opposing counsel will challenge the operator's qualification: "you're a database administrator, not a cryptographer; how do you know what the verifier's output means?" The institution's posture is that the operator's training, completion of this examiner-training module, and demonstrated ability to articulate the verifier's logic under cross-examination together qualify the operator under FRE 702.

A FRE 702 expert witness for the chain must be trained to the point of understanding:

1. **The verifier's 12-step procedure and what each step tests.** The §7 walk-through in the spec covers each step in order; the operator should be able to name each step, the data inputs it consumes, and the integrity property it verifies. Module 5 above introduces the steps; the operator's deeper preparation comes from working the spec's §7 text and the test-vector corpus.
2. **The meaning of each failure-reason string.** The §7 normative reason strings are not interchangeable. `key_fingerprint mismatch` (step 8) means an identity problem with the IKM roster; `payload_hash MAC mismatch` (step 9) means content tampering or an SDK defect. The operator must explain which reason maps to which institutional control under cross-examination.
3. **The relationship between the verifier's output and institutional controls.** A `step 8 fingerprint mismatch` does not mean "the crypto is broken." It means "check the IKM roster" — the operator routes to the institution's key-management control, not to a cryptographic-algorithm-failure response. The operator must articulate the routing under cross-examination because opposing counsel will ask "your verifier said integrity failed; isn't that the same as saying the cryptography is broken?" and the answer is "no — the failure surface is the institution's key-roster discipline, which is a different control."

The institution's IT-examiner training (this document) becomes the institution's internal training baseline. An IT witness who has completed this training, can articulate the verifier's logic under cross-examination, and has demonstrated the qualification across past examinations or mock-deposition exercises is qualified per FRE 702.

**Documentation of training completion.** The institution records the operator's training completion in personnel records — a training certificate signed by the institution's chain-operations lead, an audit-trail entry in the institution's training system, or an equivalent record. The record names the operator's identity, the date of training completion, the version of the spec the training was conducted against, and any recurring-training requirements (annual refresh, refresh on spec version change, refresh on verifier-binary major-version change). The record is part of the institution's witness-qualification evidence package; opposing counsel deposing the operator can request the record under FRE 702 foundation discovery, and the institution produces it as part of the standard discovery package per the litigation-hold posture in `operator-guide.md` §"Litigation-hold and subpoena-response timing."

**Mock-deposition exercises.** Tier-1 institutions with frequent litigation may conduct mock-deposition exercises with the operator, walking through cross-examination scripts on the verifier's logic. The exercises rehearse the operator's articulation of the 12-step procedure, the failure-reason taxonomy, and the institutional-control mapping under hostile questioning. Smaller institutions may rely on the training certificate alone; the witness-qualification posture scales with the institution's litigation exposure.

**Cross-reference.** `docs/litigation-support.md` (FRE 702 expert-witness qualification, foundation testimony for the IT witness), spec §7 (verifier procedure and normative failure reasons), `examination-response-workflow.md` (the operational five-step path the operator follows during examinations).

## Beyond the 30 minutes

Topics for deeper training:

- **Strict mode** — when to use, when not. See `07-verifier-design.md` §5.4.
- **`verifier walk`** — for incident investigation; walks one run end to end.
- **`verifier diff`** — comparing two snapshots; useful for re-examination scenarios.
- **`verifier consolidate`** — merging multi-region reports into a single view (when shipping).
- **NIST CSF mapping** — how chain output feeds CSF-aligned assessments. See `CSF-2.0.md`.
- **FFIEC Handbook mapping** — how chain output feeds Handbook-aligned examinations. See `handbook-mapping.md`.

## Sample exercises for new examiners

Available in `docs/regulator-pack/exercises/` (separately distributed):

1. **A clean-pass exercise** — verify a healthy chain; produce the bundle.
2. **A failure-mode exercise** — verify a deliberately-tampered ledger; identify the failure mode and produce findings.
3. **An anomaly-only exercise** — verify a ledger with operational anomalies; explain how to evaluate them.
4. **A multi-region exercise** — verify two regions of the same institution and consolidate the result.

Examiners complete the exercises in a sandbox before working examinations.
