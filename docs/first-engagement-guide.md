# First chain examination guide

> **What this doc is.** Bank-side and EIC-side guide for the institution's first chain-of-custody examination. Both audiences benefit from coordinated expectations.

## Bank-side preparation

### Before the examination scheduling notice

Six months out, the institution's chain operations team should already have:

- Chain operating in production for at least one quarterly cycle
- Internal audit has run the verifier at least once
- The institution's control description is current
- The IR playbook is documented and exercised at table-top level
- The SOC engagement (if applicable) is in progress or completed

### When the examination is scheduled

When the regulator notifies the institution of an upcoming examination:

1. **Identify the examiner contact.** The regulator names an EIC. The institution's chain-ops team identifies its primary contact.
2. **Pre-examination communication.** The institution and the EIC discuss scope. For a first chain examination, both sides benefit from explicit clarification: which tenant_ids, which date range, what verifier-output the institution will provide pre-examination.
3. **Pre-examination verifier run.** The institution runs the verifier on the production ledger snapshot for the examination period. The institution provides the bundle to the EIC pre-examination so the EIC knows what to expect.
4. **Pre-examination documentation handoff.** The institution provides the regulator pack documents the institution has adopted (typically: control description, recent SOC report if applicable, the verifier-output bundle, the public-key fingerprint).

### During the examination

- The institution provides the EIC with a workspace to run the verifier
- The institution's chain-ops team is available for questions about the verifier output
- The institution's MRM committee chair, audit committee chair, and CISO are available for follow-up if the examiner has questions
- The institution's IR playbook is available for the examiner's review

### After the examination

- The institution receives the examination report
- The institution responds to any findings per its standard examination response process
- The institution updates its control description if changes were identified

## EIC-side preparation

### Before the examination

The EIC's preparation:

1. **Confirm the institution operates the chain.** Receive the institution's control description; confirm the chain is in scope.
2. **Familiarize with the institution's tenant_ids.** A first-time chain examination often involves the institution explaining the tenant scheme.
3. **Receive the institution's pre-examination verifier output.** Read the bundle. The institution's pre-examination run is helpful for both sides.
4. **Review the institution's IR playbook.** Look for chain-specific scenarios; the institution should have these.
5. **Review the institution's SOC report (if applicable).** The SOC report covers the chain operations; the EIC's examination is independent but consumes the same evidence base.
6. **Allocate examination time.** A first chain examination typically takes 2–4 days as part of a larger IT examination. Subsequent examinations are faster.

### During the examination

The EIC's examination steps:

1. **Validate the verifier binary.** Use `verifier-validate.sh` against the binary the EIC will run.
2. **Run the verifier independently.** Take the institution's ledger snapshot, the public key, run the verifier; produce the bundle as the working-paper artifact.
3. **Compare to the institution's pre-examination run.** Identical results confirm the institution's pre-examination was accurate.
4. **Review anomalies.** If anomalies are present, sample-test the institution's anomaly-evaluation records (`docs/anomaly-documentation-template.md`).
5. **Examine CUECs.** Sample-test the institution's CUEC operation (`docs/control-map/CUECs.md`).
6. **Examine operational events.** Pull the institution's operational events log; confirm the events match the institution's claims.
7. **Examine IR playbook usage.** Confirm any chain-detected events during the period had documented response.

### After the examination

- The EIC writes the examination report. Findings, if any, use language from `regulator-pack/finding-language.md`.
- The institution responds.
- The EIC closes findings or schedules follow-up examinations.

## Common first-engagement issues

### Issue: Institution's tenant_id naming is unclear

**Cause.** The institution didn't think through tenant naming before deployment.

**Resolution.** The institution documents the tenant_id naming scheme (typically `tenant_<bank>_<environment>_<region>` for clarity). Future tenant_ids follow the scheme.

### Issue: Pre-examination verifier run differs from EIC's run

**Cause.** Either the institution or the EIC has a different ledger snapshot or public key.

**Resolution.** Both sides reconcile the snapshot and key. Typically the issue is that the institution has updated the ledger after producing the snapshot for the EIC; the EIC re-receives a fresh snapshot.

### Issue: Anomalies that the institution did not document

**Cause.** The institution's monitoring identified the anomaly but the anomaly-evaluation record was not produced.

**Resolution.** The institution produces the record retrospectively for the period. Ongoing, the institution updates its monitoring procedure to produce records contemporaneously.

### Issue: CUEC operation is incomplete

**Cause.** The institution adopted minimum-viable deployment but didn't document the planned ramp-up.

**Resolution.** The institution provides its CUEC ramp-up plan. The EIC accepts the plan with documented timeline; reviews progress at the next examination.

### Issue: First chain finding

**Cause.** The verifier reports a chain-related issue (chain hash mismatch, sealing delay, etc.) during the examination period.

**Resolution.** The institution invokes the IR playbook (perhaps for the first time). The institution documents the response and produces the supporting evidence. The examination report uses the appropriate finding language.

## Pre-examination checklist

For the institution:

- [ ] Chain operating in production for at least one quarterly cycle
- [ ] Internal audit has run the verifier at least once
- [ ] Control description current
- [ ] CUECs operational (or planned with timeline)
- [ ] IR playbook documented and exercised
- [ ] Customer-correlation index operational
- [ ] Reconciliation operational
- [ ] Pre-examination verifier run completed; bundle provided to EIC
- [ ] Workspace prepared for the EIC

For the EIC:

- [ ] Verifier binary on examiner laptop, validated with `verifier-validate.sh`
- [ ] Institution's control description received and reviewed
- [ ] Institution's pre-examination verifier output reviewed
- [ ] Institution's SOC report received (if applicable)
- [ ] Examination time allocated (2–4 days for first chain examination)
- [ ] Specific findings to investigate identified

## Post-engagement learning loop

After the first engagement, both sides review:

- What worked: the institution's preparation and the EIC's familiarity with the chain
- What didn't: gaps in either side's preparation, surprising findings, communication issues
- What changes for next time: institution updates its preparation procedures; the regulator updates its training materials

This guide is updated as institutions and regulators accumulate first-engagement experience.

## Resources during the engagement

The institution's chain-ops team has these resources for in-engagement support:

- `docs/operator-guide.md` — operations
- `docs/incident-response-playbook.md` — IR scenarios
- `docs/regulator-pack/` — examiner-facing material
- `docs/control-map/` — controls
- The institution's primary regulator contact

The EIC has these resources:

- `docs/regulator-pack/examiner-quickstart.md` — orientation
- `docs/regulator-pack/examiner-training.md` — training
- `docs/regulator-pack/sample-report.md` — sample
- `docs/regulator-pack/finding-language.md` — finding language
- `docs/regulator-pack/handbook-mapping.md` — Handbook alignment

Both sides reference the same materials; the chain spec is shared infrastructure.
