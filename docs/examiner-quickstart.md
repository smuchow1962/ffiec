# Examiner quickstart

> **What this doc is.** Five-minute orientation for an examiner who has never run the chain verifier. From "what is this thing" to "I can produce a report" without reading the design docs.

## What the chain does, in three sentences

The chain captures every AI agent decision a bank makes and proves the record is authentic. The verifier you'll run reads the bank's ledger snapshot, the bank's published public key, and produces a pass/fail report per day in your examination scope. If the report is "pass," the AI decisions are integrity-bearing — no one altered them.

## What you need

Five things:

1. **The verifier binary.** A single ~12 MB executable. Allowlisted on your laptop ahead of the examination.
2. **The validator script.** `verifier-validate.sh` (or `.ps1` on Windows). Confirms the binary is signed and reproducible.
3. **The bank's ledger snapshot.** A file the bank produces for the examination. Bring on USB or read-only mount.
4. **The bank's tenant public key.** A small PEM file from the bank's key registry; the regulator already has the fingerprint on file.
5. **The bank's tenant ID.** A string like `tenant_acme_prod_us_east_1`.

## What you do

```
# 1. Validate the verifier (every time, before you run it)
./verifier-validate.sh ./verifier

# 2. Record the snapshot's hash for your working paper
sha256sum ./ledger-snapshot.dump

# 3. Run the verifier
./verifier verify \
  --ledger ./ledger-snapshot.dump \
  --root-key ./tenant.pub \
  --master-key ./tenant-ikm.bin \
  --tenant-id tenant_acme_prod_us_east_1 \
  --from 2026-04-01 \
  --to 2026-04-30 \
  --report ./report.pdf \
  --bundle ./bundle.tar.gz

# 4. Check exit status
echo $?
# 0 = overall pass; 1 = at least one day failed
```

The `--master-key` flag points at the institution's IKM file (32 raw bytes, file mode 0600). Without it the verifier performs structural verification only (chain links, Merkle, signature) and skips per-event HMAC equality — under `--strict` the absent IKM elevates to FAIL. For most examinations the institution provides the IKM under a documented shape (protective order or HSM-mediated derivation path); see `legal-disclosure.md` §"Court-ordered master-key disclosure" for the protocol.

That's it. The verifier produces the PDF report; the bundle is the working-paper artifact; the examiner reads the report.

## What the report tells you

The cover page summary tells you whether the period passed, passed with anomalies, or failed.

- **Pass.** Every day in scope verified completely. AI decisions during the period are integrity-bearing. Note the result; close as routine.
- **Pass with anomalies.** Verified, but with operational notes (sealing delays, late-binding events). Confirm the bank's incident log explains each anomaly.
- **Fail.** One or more days could not be verified. Read the failure section; investigate.

## Common failure modes (quick reference)

The verifier follows the spec §7 twelve-step procedure. Each step's failure produces a specific named reason. Map your verifier output's `step` and `reason` fields to this table:

| `step` | Verifier reason string | Meaning | Severity | IR scenario |
|---|---|---|---|---|
| 1 | `format_version not supported by this verifier` | Verifier-version skew (institution's SDK is newer than your verifier) | **Not an institutional finding** — obtain newer verifier from the project | (none — examiner-side) |
| 2 | `header HKDF inputs do not match running v1 inputs` | The file's HKDF inputs don't match v1 constants | Severe; format-construction defect | Scenario 1 |
| 3 | `header genesis_hash does not match v1 constant` | File's genesis bytes are not the v1 zero-bytes constant | Severe; format-construction defect | Scenario 1 |
| 4 | `cross-chain lift detected at seq N` | Event claims a tenant/run different from the file header | High; either mis-bundled snapshot or lift attempt | Scenario 1 + evidence-handling investigation |
| 5 | `format_version mismatch at seq N` | Mid-file format drift between header and entry | High; SDK version-handling defect | Scenario 1 |
| 6 | `chain link broken at seq N` / `seq out of order at seq N` | Structural chain integrity broken | High; possible tampering | Scenario 1 |
| 7 | `unknown key_version: no IKM for (tenant=T, key_version=V)` | Institution's IKM registry is missing a referenced generation | High; investigate IKM retention or provisioning | Scenario 8 |
| 8 | `key_fingerprint mismatch at seq N` | Looked-up IKM does not produce the recorded fingerprint. Identity mismatch (botched rotation, restored backup, cross-tenant drift), NOT content tampering. | Severe; investigate against IKM roster, not chain content | Scenario 7 |
| 9 | `payload_hash MAC mismatch at seq N` | Recomputed MAC differs from stored MAC | Severe; content tampering | Scenario 1 |
| 10 | `merkle root mismatch — ledger contents do not produce sealed root` | Day's events don't produce the sealed root | Severe; possible server-side history rewrite | Scenario 2 |
| 11 | `signature verification failed` (or `algorithm/key-type mismatch`) | Daily seal signature didn't verify | Severe; possible HSM key compromise | Scenario 3 |
| 12 | `cadence mismatch` | Recorded cadence doesn't match the institution's claimed cadence | Observation; control-description-accuracy issue | Scenario 5 variant |
| 12 | `dev-mode seal in production verification — refused` | Production seal signed with software key (compile-time exclusion bypassed) | Severe; misconfiguration | Scenario 6 |
| (file pre-flight) | `audit file ends mid-line — possible mid-write crash` | SDK process died mid-append (writer crash recovery) | Operational, NOT tampering — Medium severity | Scenario 9 |
| 12a | `gen_ai_model_identifier_missing at seq N: {field_name} required for chain entries representing model calls` | Chain entry represents a model call (carries any `gen_ai.*` attribute) but lacks `gen_ai.request.model` or `gen_ai.response.model`. Control-completeness for SR 11-7 reproducibility, NOT chain-integrity. | Medium under non-strict; FAIL under `--strict`. | (anomaly only — no IR; MRM-program finding under audit-procedures P-25) |

For each failure mode, the examination response language is in `regulator-pack/finding-language.md`. The single workflow doc `regulator-pack/examination-response-workflow.md` walks the path from a JSON failure record to a documented examination response in 60 seconds.

**Bundle sufficiency.** The bundle (`--bundle ./bundle.tar.gz`) is sufficient for a second examiner to re-verify the institution's chain output independently, without re-engaging the institution. Bundle contents per `07-verifier-design.md` §5.3.4: report.pdf, report.json, verifier.sha256, ledger.sha256, public_key.pem, metadata.json. A receiving examiner runs `bundleverify` on their own laptop, validates the verifier-binary signature, and re-runs against the bundled ledger snapshot. The bundle is the working-paper artifact; downstream re-verification is a property of the bundle alone, not of the institution's continued cooperation.

**Three specific points to note:**

- **`key_fingerprint mismatch` and `payload_hash MAC mismatch` are both "stop and call the bank" findings.** Both are Severe and warrant immediate institution contact. The distinction is the investigation path: `key_fingerprint mismatch` (step 8) is investigated through the IKM roster — the institution looked up an IKM that does not produce the recorded fingerprint (botched rotation, cross-tenant configuration drift, swapped backup). `payload_hash MAC mismatch` (step 9) is investigated through the chain content — the recomputed MAC does not match the stored MAC, indicating in-flight tampering or a defect. Neither is a verifier bug; both require institution response. The severity is identical; the investigation paths differ.
- **`audit file ends mid-line` is NOT a tampering finding.** It is a writer-side mid-write crash (the SDK process died between events). Severity is sealing-delay-equivalent operational, not high-severity integrity. Treat it as a crash-recovery case, not a tampering case.
- **For each failure, `examination-response-workflow.md` is the 60-second path** from the JSON failure record to a documented examination response. Do not reconstruct the workflow from this quickstart at the moment of decision — use the workflow doc.

## What you don't need to know

The cryptographic details. The chain uses HMAC-SHA-256, RFC 6962 Merkle trees, and Ed25519 signatures. The verifier handles all of this; you read the output. If you want to go deeper, the design docs at `docs/design/` are extensive.

## Where to go next

- **You're ready for your first verification.** Use the sample report in `regulator-pack/sample-report.md` as a reference; run the verifier on the bank's snapshot; produce your report.
- **You want full training.** `regulator-pack/examiner-training.md` is the 30-minute training session.
- **You want to understand the standard.** `regulator-pack/handbook-mapping.md` shows how the chain maps to FFIEC IT Handbook control objectives.

## If something goes wrong

- **The verifier won't run.** Re-run the validator. If it fails, the binary is corrupted or unsigned; obtain a fresh copy.
- **The verifier exits with an error before producing output.** The ledger snapshot or public key may be malformed. Check with the bank's chain-operations team.
- **The verifier produces a "fail" result.** That is the verifier doing its job — investigate per the failure-mode reference above.
- **You're stuck.** Contact the regulator's IT examination function for support. The verifier is a standard tool; they have seen this before.
