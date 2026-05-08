# Incident response playbook — chain-detected events

> **What this doc is.** The institution's playbook for events the chain catches. Adapted by the institution to fit its existing IR framework, role assignments, and notification thresholds.

## Scope

This playbook covers events the chain produces or detects:

- Chain hash mismatch (in-flight or at verification) — spec §7 step 9 (`payload_hash MAC mismatch`)
- Chain link broken (sequence integrity failed) — spec §7 step 6
- Merkle root mismatch (recomputed root differs from signed root) — spec §7 step 10
- Signature verification failed (Ed25519 signature does not match) — spec §7 step 11
- **Key fingerprint mismatch (Scenario 7)** — spec §7 step 8 (the load-bearing rework primitive)
- **Unknown key_version (Scenario 8)** — spec §7 step 7
- **Audit file truncation detected (Scenario 9)** — spec §4.1 mid-write truncation refusal
- **Backup integrity failure causing verifier failure (Scenario 10)** — RC.RP-04
- Sealing delay beyond 72 hours
- Master key compromise (suspected or confirmed)
- Software-key fallback observed in production

Events in the bank's broader infrastructure (network compromise, application-host compromise) are out of scope here; the bank's main IR framework handles them. The chain composes with that framework — this playbook is the chain-specific overlay.

## Severity classification

| Severity | Examples | Response time |
|---|---|---|
| **Critical** | Signature verification failure; suspected master compromise; software-key in production | 15 min to acknowledge; 1 hour to contain |
| **High** | Chain hash mismatch; merkle root mismatch | 30 min to acknowledge; 4 hours to contain |
| **Medium** | Sealing delay >72h; backup-recovery integrity failure | 1 hour to acknowledge; 24 hours to remediate |
| **Low** | Late-binding rate elevated; clock skew anomalies | 24 hours to acknowledge; remediate at next business day |

## Roles

| Role | Responsibility |
|---|---|
| Incident Commander (IC) | Coordinates response; communicates with leadership |
| Chain Operations Lead | Investigates the chain-specific aspects |
| HSM Administrator | Investigates HSM-related aspects; performs key rotation when needed |
| CISO Delegate | Decides on regulator notification and external communications |
| Legal Counsel | Advises on regulatory obligations and external statements |
| Communications Lead | Manages internal and external messaging |

## Common scenarios

### Scenario 1 — Chain hash mismatch detected at ingest

**Trigger.** Operational event `chain.verification_failure` with `reason=payload_hash_mismatch`.

**Investigation.**

1. Identify the affected `(tenant_id, run_id, seq)`.
2. Pull the affected event from the ingest queue (it was rejected; not in the WAL).
3. Pull the run's previous events; confirm chain integrity up to the affected seq.
4. Identify the source: which application host produced the event, under which `(tenant_id, key_version)` pair (and recorded `key_fingerprint`).
5. Confirm the host's process and the recorded `key_fingerprint` are on the institution's authorized roster (cross-check against the most recent `master.reconciliation_completed` event).

**Most common root causes.**

- SDK defect (chain construction bug)
- Network corruption (very rare with TLS)
- Tampering (most concerning case)

**Containment.**

- If SDK defect suspected: pause new ingest from the affected SDK version; pin the rollback target.
- If tampering suspected: isolate the affected host; preserve forensic evidence; revoke the session key.

**Remediation.**

- SDK defect: fix and re-deploy; replay any lost events from local SQLite buffers if available.
- Tampering: full incident response per the institution's main IR framework; preserve forensic evidence.

**Notification.**

- SDK defect: standard change-management; no regulator notification unless customer impact.
- Tampering: regulator notification per the cyber-incident notification rule.

### Scenario 2 — Merkle root mismatch at verification

**Trigger.** Verifier reports `merkle root mismatch` for a specific tenant-day.

**Investigation.**

1. Identify the affected `(tenant_id, day)`.
2. Compare the recomputed root to the recorded root.
3. Walk the day's events; identify which event(s) produced the mismatch.
4. Determine whether events were added, removed, or altered relative to the original computation.

**Most common root causes.**

- Storage corruption (rare with replicated storage)
- Recovery from older backup without gap-fill
- Tampering (most concerning case)

**Containment.**

- If recovery scenario: stop further verifier runs on the affected period; preserve the WAL.
- If tampering suspected: full incident-response activation; isolate database administrators not on the response team.

**Remediation.**

- Recovery scenario: replay the gap from any available source (SDK buffers, downstream OTLP); re-run the seal job; verify.
- Storage corruption: investigate the storage layer; restore from a known-good backup; replay the gap.
- Tampering: forensic investigation; rotate signing key as a precautionary measure; notify regulator.

### Scenario 3 — Signature verification failed

**Trigger.** Verifier reports `signature verification failed` for a specific tenant-day.

This is among the most severe scenarios because it suggests one of:

- Public-key registry mismatch
- Seal record corruption
- Substitution of the seal by an unauthorized party
- HSM signing-key compromise

**Investigation.**

1. Confirm the public key the verifier used matches the published key for the affected day's `public_key_id`.
2. Re-check the seal record; confirm it is byte-identical to the original at the time of recording.
3. If the public key matches and the seal record is unaltered, the only remaining explanation is that the original signer was not the authorized HSM key.

**Containment.**

- Treat the affected day's seal as untrusted.
- Activate the HSM-compromise response: stop the seal job, isolate the seal-job operator role.
- Notify the CISO and the regulator immediately.

**Remediation.**

- Rotate the HSM signing key.
- Re-sign all seals signed under the suspect key (if the underlying ledger contents are intact).
- Publish the new public key to the tenant key registry; coordinate with the regulator on the verification approach for the affected period.

**Notification.**

- Regulator notification within 36 hours per the cyber-incident notification rule.
- Public disclosure if customer impact is identified (per the institution's customer notification policy).

### Forensic preservation for Critical-severity scenarios

This subsection applies to Scenarios 1, 2, and 3 above (chain hash mismatch, Merkle-root mismatch, signature verification failed). When a chain-integrity violation is detected, the institution MUST preserve forensic evidence of the violation BEFORE attempting remediation. Remediation that runs first overwrites the very state a later examiner needs to see.

**Required preservation steps.** The IR team performs each of these before any repair, key rotation, or ledger modification:

1. Capture a full memory dump of the verifier or affected ledger process if it is still running. The dump is stored read-only with a SHA-256 hash recorded at capture time.
2. Snapshot the chain-entry file or database table at the moment of detection. The snapshot is taken read-only and never written back.
3. Preserve the verifier output showing the failure — the full stdout and stderr, the exit code, and any structured-output file the verifier wrote.
4. Preserve the operational event log for the affected period. Operational events relevant here include `chain.verification_failure`, `audit_file.truncation_detected`, and the day's `seal.job_completed`.
5. Preserve the verifier-invocation history per the soc-pack `verifier.run_completed` event. This shows when the verifier was run, by whom, with what flags, and on what host — the meta-evidence that the verifier output is itself authentic.

**Order matters.** Only after forensic preservation is complete may the institution begin remediation (key rotation, ledger repair, replay from backup, etc.). A repair that runs before preservation destroys the evidence the institution needs to defend itself later.

**Custody documentation.** Record the preservation date, time, and custody chain for each artifact: who handled it, when, where it is stored, and the access-control restrictions placed on it. The chain-of-custody form names each artifact by its SHA-256 hash so a later examiner can confirm the artifact has not been altered since preservation.

**Why this matters.** This evidence package is the institution's affirmative defense against later spoliation claims. If the incident leads to litigation or a regulatory examination, the institution's IT witness testifies from this preserved evidence: the violation was detected, the state at detection was preserved before any remediation, and the preserved state is integrity-bound by hash. Without this discipline, opposing counsel argues "the institution detected a problem and then altered the system before anyone could examine it" — a spoliation claim under FRCP 37(e).

Cross-reference: `docs/litigation-support.md` (spoliation defense, evidence preservation), spec §7 (verifier procedure, what each step tests).

### FRCP 37(e) spoliation defense

When a chain-detected event occurs and the institution anticipates litigation (customer dispute, regulatory action, civil suit), the institution's response includes preservation steps designed to satisfy FRCP 37(e)(1)'s "good-faith preservation" standard. FRCP 37(e) allows sanctions for failure to preserve electronically stored information when the failure is willful or reckless and the loss prejudices the opposing party. The institution's affirmative posture is that it operated the chain in good faith, detected the problem, and took immediate action — no prejudice to the opposing party because the chain's integrity claim stands.

**Required steps when chain-detected events occur during anticipated litigation:**

1. Extend chain-data retention beyond the institution's standard schedule. The litigation-hold extension covers chain entries, seal records, IKM history, verifier output, and operational events for the affected period.
2. Freeze IKM rotation OR document the post-incident IKM generation with explicit before-and-after provenance. If a rotation is operationally required, the institution captures the IKM custody chain before and after the rotation so the chain's verifiability survives the rotation.
3. Log all access to chain data and IKM during the incident and remediation window. Every read, every verifier run, every administrative action against the affected period is recorded.
4. Retain verifier output showing the incident and remediation. The "before" verifier output documents the failure; the "after" output documents the remediation succeeded.

**Affirmative defense narrative.** These steps constitute good-faith preservation per FRCP 37(e)(1). The institution's litigation posture is:

> "We operated our chain in good faith. When a violation was detected, we preserved evidence of the violation before any remediation, extended retention on all affected data, froze key rotation during the response, and logged every action taken during the response. The preserved evidence and the verifier output prove both the problem and our response. The opposing party suffers no prejudice because the chain's integrity claim is verifiable from the preserved evidence."

This language protects the institution against the "spoliation" claim that the chain failure means the institution destroyed evidence.

Cross-reference: `docs/litigation-support.md` (litigation hold, FRCP 37(e) procedural posture), spec §10.9 (IKM retention rules).

### Scenario 4 — Master key compromise (suspected)

**Trigger.** Could be:

- Key-fingerprint reconciliation (spec §10.1) produces `fingerprint_unmatched_count > 0` on the `master.reconciliation_completed` event
- Out-of-band intelligence (vendor security advisory, threat-intel feed)
- Anomalous handshake patterns at the master-key custodian

**Investigation.**

1. Confirm the suspect `(tenant_id, key_version, key_fingerprint)` triples in the ledger.
2. Identify the events captured under the suspect ids; bound the compromise window.
3. Determine whether any captured events were business-critical (decisions affecting customers).

**Containment.**

- Rotate the master key immediately at the master-key custodian.
- Force all application processes to re-handshake; new sessions derive from the new master.
- Mark events captured during the suspected compromise window as "potentially repudiable" in the institution's records.

**Remediation.**

- Forensic investigation of the master-key custodian.
- Update the institution's incident-response framework if a new attack vector is identified.
- Re-validate any business decisions made during the compromise window using non-chain evidence.

**Notification.**

- Regulator notification within 36 hours.
- Customer notification if specific customer-impacting decisions occurred.

### Scenario 5 — Sealing delay beyond 72 hours

**Trigger.** Operational event `seal.job_failed` persists beyond the 72-hour notification threshold.

**Investigation.** Standard operational triage — HSM availability, network connectivity, seal-job health.

**Containment.** Continue capturing events; the chain construction is independent of the seal job.

**Remediation.** Restore HSM availability or network path; the seal job runs on the next retry; the seal is recorded with `signed_at` later than `seal_date + 60 minutes` and the verifier reports the delay.

**Notification.** Regulator notification per `04-hsm-custody.md` §5.2 SHOULD clause.

### Scenario 6 — Software-key fallback in production

**Trigger.** The verifier reports a seal carrying `dev-mode: true` for a production tenant, OR a chain entry carrying `kms_handle_uri` with the `"plaintext-"` prefix.

This indicates the software-key fallback was active in production — a misconfiguration severity. Per spec §10.7 the production build is supposed to exclude the software-key adapter at compile time; this scenario fires when that exclusion is bypassed (faulty build, vendor distribution gap, or an institution that built from source incorrectly).

**Investigation.**

1. Identify how the fallback was enabled (faulty build pipeline, distribution gap, source-build error).
2. Identify the affected tenant-days.
3. Confirm the production build's compile-time exclusion is in force going forward.

**Containment.** Disable the fallback; restart the seal job under HSM-backed signing; rebuild from a known-good source with the exclusion in place.

**Remediation.** Re-issue the affected seals from the HSM-backed key. Update the institution's deployment-validation checks to detect the fallback in production.

**Notification.** This is a control failure, NOT a 36-hour cyber-incident. Regulator notification per the institution's standard control-failure framework. **Clock-start trigger:** does not start the 36-hour clock unless the fallback's presence is associated with confirmed unauthorized signing.

### Scenario 7 — Key fingerprint mismatch (the rework's load-bearing detection)

**Trigger.** The verifier reports `key_fingerprint mismatch at seq N: looked-up IKM does not match the entry's recorded fingerprint`. The verifier emits this at spec §7 step 8 — BEFORE any MAC compute. The corresponding operational event is `chain.verification_failure` with `step=8`. The reconciliation procedure (spec §10.1) may also flag the condition through `master.reconciliation_completed` with `fingerprint_unmatched_count > 0`.

**Investigation triage tree.** A `key_fingerprint` mismatch is *almost always* an operational misconfiguration, not a security incident. Triage in this order:

1. **Was a rotation in flight on the affected day?** Pull `master_key.rotated` and `master_key.rotation_observed` events for the period. If yes: route to change-management owner; the mismatch is a rotation-rollout incident, severity P2.
2. **Was a backup restored to the IKM-roster recently?** Pull the institution's KMS / IKM-roster change log. If a restore is correlated, the restore likely brought back stale or wrong-tenant IKM bytes; route to backup-recovery procedure, severity P2.
3. **Was a tenant row in the registry restored from a different period?** Same as above, scoped to the registry layer; route to registry-recovery procedure, severity P2.
4. **None of the above?** Treat as suspected unauthorized key substitution, severity P1; route to Scenario 4 (master compromise suspected) and start the 36-hour clock.

**Containment.**

- For triage cases 1-3 (operational misconfig): pause new ingest from the affected `(tenant_id, key_version)` pair; stop the verifier from processing further events under the wrong fingerprint until the IKM roster is corrected.
- For triage case 4 (suspected unauthorized substitution): full activation of Scenario 4.

**Remediation.**

- For operational misconfig: identify the IKM-roster row that was wrong; restore the correct IKM bytes from KMS history or backup; document the change-management record; re-run the verifier against the corrected roster.
- For suspected unauthorized substitution: per Scenario 4 — rotate the master, mark the compromise window, forensic investigation.

**Tie-in to spec §10.1 reconciliation.** If the fingerprint mismatch fires OUTSIDE the weekly reconciliation cadence, escalate to this scenario. If it fires DURING the reconciliation, route to the reconciliation owner first (triage cases 1-3 typically surface there).

**Notification.** **Clock-start trigger:** triage cases 1-3 do NOT start the 36-hour clock (control failure, not security incident). Triage case 4 starts the clock at the moment the institution determines the mismatch is not attributable to documented operations.

### Scenario 8 — Unknown key_version

**Trigger.** The verifier reports `unknown key_version: no IKM for (tenant=T, key_version=V) at seq N`. The verifier emits this at spec §7 step 7 — also BEFORE any MAC compute. The corresponding operational event is `chain.verification_failure` with `step=7`.

**Three plausible root causes, each with different IR triage:**

1. **Premature IKM retirement.** The institution decommissioned an IKM and discarded the bytes while chain entries that reference it still exist. Per spec §10.9 retention rule, IKM MUST be retained for the duration of any chain entry referencing it; this is a retention-control failure. Route to retention-control owner; restore the IKM from KMS history if available; if not, the affected events are unverifiable and the institution documents the gap.
2. **Provisioning gap.** The institution rotated the IKM at the SDK side but never registered the new generation in the verifier-side registry. Operational; route to provisioning-control owner; register the IKM and re-run the verifier.
3. **Tampered key_version field.** The entry's `key_version` was altered after capture to a value that has no corresponding IKM. Tampering; route to Scenario 1 (chain hash mismatch) for the integrity-side investigation; the per-event MAC will also fail at step 9 once the right key is found, so cross-check both signals.

**Containment.** For cases 1 and 2 (operational): no immediate containment beyond fixing the registry. For case 3 (tampering): treat as Scenario 1 with the addition that the attacker chose to alter the `key_version` field — investigate why; the change suggests the attacker had ledger-write access.

**Remediation.** Case-specific per the triage above.

**Notification.** **Clock-start trigger:** cases 1 and 2 do NOT start the 36-hour clock (control failures). Case 3 starts the clock per Scenario 1.

### Scenario 9 — Audit file truncation detected

**Trigger.** The verifier refuses to verify an audit file because the file's last byte is not `\n`. The corresponding operational event is `audit_file.truncation_detected`.

This is almost always a writer-side mid-write crash (the SDK process died between events). Per spec §4.1, the verifier refuses truncated files rather than silently passing chains that lost their last entry — this is a crash-recovery surface, not a tampering surface.

**Investigation.**

1. Identify the writer process that crashed (PID, host, time).
2. Pull the SDK's local SQLite buffer for the affected period; recover the dropped event(s) if possible.
3. Pull any upstream OTLP backend retention for the affected period; recover from there if SDK buffer is unavailable.
4. Document the gap if recovery is not possible.

**Containment.** No containment needed; the chain is intact through the last complete entry, and new chain entries continue to be captured.

**Remediation.** Recover the dropped events; replay into the ledger if a replay path exists; document the gap in the institution's incident log.

**Notification.** **Clock-start trigger:** does NOT start the 36-hour clock unless the truncation pattern recurs across hosts in a way that suggests deliberate corruption (rare; investigate carefully before assuming).

### Scenario 11 — Project-side trust-anchor degradation (cold-DR fallback)

**Trigger.** The institution observes one of: (a) the project's annual cold-DR-key dry-run attestation (`KEY-DR-DRYRUN-{year}.asc`) failed to validate against the cached cold-DR public key; (b) the project announces an emergency cold-DR-key activation due to dual-compromise of cosign + standard-GPG; (c) the project misses its annual dry-run window without explanation.

This is a project-side governance event the institution does not control but must respond to operationally.

**Investigation.**

1. Confirm the trigger via the institution's release-validation procedure (the procedure consumes `KEY-DR-DRYRUN-{year}.asc` per `supply-chain.md` §"Cold-disaster-recovery key lifecycle").
2. Determine which trigger applies: failed dry-run, dual-compromise activation, or missed dry-run window.
3. Cross-check against out-of-band channels (FFIEC working group notification list, the project's announcements channel).

**Containment.**

- For (a) failed dry-run: pause new-version verifier deployments; continue using the historical verifier against historical seals (the historical binary's signatures remain validatable against the pre-degradation cached public keys).
- For (b) dual-compromise activation: emergency response per `supply-chain.md` §"Dual-compromise edge case"; coordinate directly with the institution's primary federal regulator on examination posture during the multi-week recovery.
- For (c) missed dry-run window: institution opens a project-side governance ticket via the FFIEC working group; treats the residual as elevated until resolved.

**Remediation.** Per the project's published recovery procedure. Institution-side action: when the project completes the new-key authentication chain, install the new keys per the institution's standard trust-anchor-update procedure (per the new Adversary I reception procedure in `09-threat-model.md` §2.9).

**Notification.** **Clock-start trigger:** does NOT start the 36-hour FFIEC clock by itself (this is project-side governance, not an institution-side cyber-security incident). May trigger institution-side regulatory communications on examination-posture stability if the institution is mid-engagement; counsel decides per institution's communication policy.

#### Scenario 13 — Mirror-registry signature-validation gap

**Trigger.** The institution's mirror-registry continuous-monitoring control (per `supply-chain.md` §"Mirror continuous-monitoring control (DE.CM-09)") fires `mirror.reconciliation_completed` with a discrepancy: the mirror's audit-log-validation set does not include all images pulled in the period. The mirror has been pulling images without recording project-side cosign validation — the bridge invariant (institution's IR playbook bridge documentation per `supply-chain.md` §"Mirror-registry signature handling") cannot be confirmed for the affected images.

**Investigation.** Three plausible root causes:

1. **Operator misconfiguration at the mirror.** Validation step was disabled or misconfigured during a mirror-side change.
2. **Registry-side bug.** The mirror's audit-log emission is broken (logs to wrong destination, log-rotation deleted entries).
3. **Vendor regression.** A mirror-software update broke the validation pipeline.

**Containment.** Pause new deployments from the mirror. Rebuild mirror-side images from project-side source against fresh cosign validation; confirm validation log is being written. The historical binary deployed against unvalidated images remains in service if no integrity signal has fired against it; the institution treats deployed unvalidated images as a known residual until the validation chain is rebuilt.

**Remediation.** Mirror-side validation rebuild; full re-validation of all images pulled during the affected period; document the remediation in the institution's control-evidence repository. The next `mirror.reconciliation_completed` event MUST show 100% match.

**Notification.** **Clock-start trigger:** does NOT start the 36-hour FFIEC clock by itself (this is a control-completeness issue, not a security incident). Conditional on whether deployed unvalidated images are subsequently confirmed compromised.

### Scenario 14 — Edge-device physical compromise (Pattern B in-service device)

**Trigger.** Per `edge-and-federated-ai.md` Pattern B, edge devices are commissioned with bulk session-key issuance from the institution's master-key custodian. Trigger: institution's edge-fleet inventory shows a physical device unaccounted for (lost, stolen, vendor-returned without decommissioning), OR the device's tamper-detection mechanism (TPM, secure enclave) reports an intrusion event.

**Investigation.**

1. Identify the device's last-known commissioning batch and the bulk session-key set it received.
2. Pull the device's last-known operational events from the institution's fleet-monitoring system.
3. Determine whether chain entries were captured under the device's session keys after the last-known communication; flag those entries as potentially-compromised.

**Containment.** Revoke the device's session-key set at the master-key custodian (the new IKM generation invalidates the device's stale keys). Rotate the institution's IKM if the compromise is suspected to extend to the IKM derivation path (e.g. multiple devices in the same batch). Mark the affected period's edge events as "potentially compromised" in the institution's records.

**Remediation.** Re-commission the affected fleet under a new IKM generation; verify the chain entries from the affected period under the new IKM (which will fail key_fingerprint check, surfacing the compromise scope). The institution's broader fleet-IR procedure handles the physical-device aspects (recovery, forensics).

**Notification.** **Clock-start trigger:** STARTS the 36-hour clock at the institution's determination that the device's session keys may have been used to forge chain entries. Edge-device physical compromise is per spec §10.1 a key-fingerprint reconciliation alert; per the FFIEC rule it is a key-compromise event for safety/soundness purposes when the affected events are material.

### Scenario 11 sub-variant — Trust-anchor reception failure (forged or invalid rotation notice)

**Trigger.** During the institution's regulator-fingerprint reception procedure (`09-threat-model.md` §2.9 step 2), the rotation notice fails authenticity validation: the notice's GPG signature does not validate against the regulator's published key, OR the cross-channel parallel notification on the regulator's authenticated-domain channel disagrees with the received notice. The institution emits `regulator_fingerprint.rotation_validated` with `overall_validation: FAIL`.

**Investigation.** Three plausible root causes:

1. **Forged notice.** An attacker attempted to substitute a malicious public-key fingerprint into the institution's verifier configuration. Severe; likely associated with a broader compromise attempt.
2. **Operational error at the regulator.** The regulator's signing identity has rotated and the institution's cached signing-key record is stale; OR the regulator's communication channel had an integrity failure (truncation, transcription error in a paper notice).
3. **Institution-side validation tooling defect.** The institution's GPG-validation script or cached key has drifted; the notice is legitimate but the institution can't verify it.

**Containment.** Do NOT install the new fingerprint. The institution's verifier configuration retains the previous fingerprint as the active trust anchor until validation succeeds.

**Remediation.** Out-of-band cross-check: contact the regulator's primary IT examination liaison via a separate authenticated channel (phone call, in-person visit, regulator's encrypted-email system if separate from the notice channel). Confirm or refute the notice's authenticity. If forged: full IR activation per Adversary I attack-2.4 framing; coordinate with the regulator on whether the forgery attempt indicates a broader attack against the institution. If operational error or tooling defect: the regulator re-issues a corrected notice OR the institution updates its validation tooling, then re-runs the reception procedure.

**Notification.** **Clock-start trigger:** for case 1 (forged notice) — STARTS the 36-hour clock at the institution's determination that the notice is forged. For cases 2 and 3 (operational): does NOT start the clock; resolve operationally.

### Scenario 12 — Co-signed seal failure (spec §7 step 11 case (e))

**Trigger.** Verifier reports `co-signed seal failure: algorithm X validated, algorithm Y did not` per spec §7 step 11 case (e). One algorithm in the institution's declared dual-algorithm posture validated; the other did not.

**Investigation triage tree.** Three plausible root-cause branches with different IR dispositions:

1. **Branch (i): attributable to a published algorithm break.** A NIST-published or FIPS-published practical attack on algorithm Y was announced in the period. The un-broken algorithm's signature still provides integrity assurance under that algorithm's assumption. The institution coordinates with its primary federal regulator on the migration timeline (per `09-threat-model.md` §2.8 dual-algorithm transitional period). **Clock-start: NO 36-hour clock** (this is a known cryptographic-deprecation event the regulator is also aware of).
2. **Branch (ii): attributable to per-algorithm signing-key compromise.** The Y-algorithm signing key has been compromised (key material exfiltrated, HSM compromised, etc.). This is IR Scenario 4 (master-key compromise) for the Y-algorithm key specifically; the X-algorithm key is unaffected. **Clock-start: STARTS** at the institution's determination that the compromise is credible (Scenario 4's clock-start trigger applies).
3. **Branch (iii): under investigation.** The root cause is not yet attributable. **Clock-start: STARTS** at the institution's investigation-conclusion determination — the institution must commit to a determination within a bounded window (typically 48-72 hours) to avoid the clock-start being indefinitely deferred.

**Containment.** For (i): no immediate containment beyond pausing new co-signed seals under the broken algorithm; continue accepting Y-algorithm-only seals during the transition. For (ii) and (iii): per Scenario 4 — rotate the affected algorithm's key, mark the compromise window, forensic investigation.

**Remediation.** Branch-specific per the triage above.

**Notification.** Per the clock-start trigger above. The institution's working paper records both the valid-algorithm validation result and the invalid-algorithm failure (per spec §7 step 11 examiner working-paper convention) so the regulator has the full picture.

### Scenario 10 — Backup integrity failure causing verifier failure

**Trigger.** The verifier returns FAIL on a recovered ledger because the day's events do not produce the sealed root (verifier reports `merkle root mismatch`). This is the Scenario 2 triage tree's "recovery from older backup without gap-fill" branch, escalated to its own scenario for the RC.RP-04 evidence requirement.

**Investigation.**

1. Confirm the recovery scenario: was a recent recovery run on the affected ledger?
2. Identify the gap: which events captured between the backup point and the failure point are missing?
3. Identify gap-fill sources: SDK local SQLite buffers, downstream OTLP backends that retained the events, vendor-held copies.

**Containment.** Stop further verifier runs on the affected period until the gap is filled or formally documented as unverifiable.

**Remediation.**

- Replay the gap from any available source.
- Re-run the seal job for the affected days; with the gap filled, the recomputed root matches the original seal.
- The verifier's next pass succeeds; the recovery event is documented in the institution's incident log.
- If the gap cannot be filled, the affected days remain unverifiable. The institution treats this as an integrity-control failure that triggers cyber-incident notification.

**Notification.** **Clock-start trigger:** does NOT start the 36-hour clock if the gap is filled and the verifier passes within the same recovery window. STARTS the clock when the institution determines the gap is unfillable and the affected days remain unverifiable. The institution documents the integrity-control failure for RC.RP-04 evidence.

### Scenario 15 — Routing decision missing during incident reconstruction

**Trigger.** Regulator inquiry, customer dispute, or model-substitution audit requires reconstruction of the institution's behavior at time T. The institution pulls the chain entries for the affected `(tenant_id, run_id)` window and discovers that the routing chain entries (`audit.routing.attempt`, `audit.routing.success`, `audit.routing.failover`, `audit.routing.circuit_state_change` per spec §4.4.1) are absent or sparse. The LLM-call chain entries are present and the chain itself verifies clean, but the *decision path* leading to each call is invisible. The trigger may also surface from the SOC team's P-33 routing-completeness sample testing (`docs/audit-procedures.md` P-33) when the period's gap rate exceeds the institution's documented threshold.

This is a control-program reconstruction scenario, not a chain-integrity scenario. The chain's MAC + Merkle + HSM coverage of whatever was captured is intact; the gap is in what the chain captured during the window, not in how it was captured. The institution's response is to determine which of three failure-modes produced the gap and route the investigation accordingly.

**Clock-start.** Scenario 15 does NOT by itself start the 36-hour FFIEC clock. The scenario is a control-completeness reconstruction event, not a security incident. If the investigation under failure-mode (c) below subsequently identifies a vendor-side reroute that constitutes an unauthorized substitution of provider, the institution's IR Commander evaluates whether Scenario 4 (master-key compromise — broader integrity framing) or a separate vendor-incident path applies; in those subordinate paths the clock-start triggers per the subordinate scenario, not per Scenario 15 itself.

**Investigation — three failure-modes.**

1. **Failure-mode (a): chain decorator was not wired into the router (deployment misconfiguration).** The institution's chain decorator is supposed to subscribe to the router's routing-event hook per `docs/design/02-chain-construction.md` §11. In this failure-mode the subscription was never wired in the affected deployment — the router is operating its policy and emitting hook events, but no chain entry is produced because nothing is listening. The chain captures every LLM call (the LLM-call decorator is wired separately) but the routing decisions are absent across the entire window. **Disposition: control-program finding, NOT a security incident.** The institution remediates by wiring the decorator into the affected router, redeploys, and reruns P-33 against a post-remediation period to confirm capture is restored. Document the gap in the institution's control-evidence repository; the affected window's routing decisions are unrecoverable but the LLM calls themselves remain auditable for the limited reconstruction the chain still supports.

2. **Failure-mode (b): routing entries are present but sparse or inconsistent during the window (operational signal degradation).** The chain decorator IS wired and routing entries appear before and after the window, but the window itself shows degraded coverage — some LLM calls have routing predecessors and others do not, OR the coverage drops below the institution's baseline for the period. Likely root causes: (i) the OTel collector was dropping routing events during the window (collector-side back-pressure, queue overflow, deliberate filter rule the institution did not intend); (ii) the router's hook emission was intermittently failing during the window (router-side bug, capacity event, hot-reload race); (iii) the chain decorator's subscription was operating but a downstream pipeline component was filtering on a stale event-type allowlist. **Disposition: investigate router or OTel collector for the affected window.** The institution pulls the OTel collector logs for the window, the router's operational logs, and the chain decorator's emission logs; cross-reference to identify which component was losing events. Remediate the affected component. The window's routing-decision evidence remains partial; the institution's reconstruction narrative names the partial-coverage explicitly when responding to the regulator inquiry / customer dispute / model-substitution audit.

3. **Failure-mode (c): routing entries are present but inconsistent with downstream LLM calls (provider-side reroute the institution did not record).** The chain decorator captured routing entries that name Provider A as `audit.routing.provider_chosen`, but the LLM call that follows landed at Provider B. Likely root causes: (i) the institution's router selected Provider A and emitted the routing event, but the underlying HTTP client (or vendor SDK) silently rerouted to Provider B at the network or DNS layer — vendor-side load-balancing across regions, a CDN-level provider substitution, or an institution-side service-mesh rule the chain decorator does not see; (ii) the institution's routing policy was changed mid-window in a way that the chain decorator captured under the old `policy_version` but the LLM client respected the new version; (iii) genuine vendor reroute the institution did not authorize (vendor-relationship investigation). **Disposition: investigate the vendor relationship.** The institution pulls the vendor's per-call logs, cross-references the institution's `audit.routing.provider_chosen` against the vendor's record of which provider produced each call, and determines whether the discrepancy is explainable (vendor-side load-balancing the institution accepted contractually but did not document in the chain decorator) or unexplained (genuine reroute warranting escalation). For unexplained reroutes the investigation handoff goes to the institution's MRM committee — see Escalation paths below.

**Escalation paths.**

- **Failure-mode (a)** routes to the institution's deployment-control owner. The remediation is a deployment-time configuration fix; no MRM committee review is required unless the gap window overlapped a model-substitution audit or a customer dispute that consumed the routing evidence.
- **Failure-mode (b)** routes to the institution's chain-operations team and the OTel platform team jointly. The remediation is a pipeline-component fix; MRM committee notification is RECOMMENDED when the affected window overlapped any decision the MRM program flagged for reproducibility review under SR 11-7.
- **Failure-mode (c)** routes to the institution's vendor-management team for the vendor-side investigation AND to the institution's MRM committee for the model-governance assessment. The MRM committee evaluates whether the unrecorded reroute represents (i) a vendor-side operational variance the institution accepts going forward (committee documents the variance and updates the institution's CC8.1 control description), (ii) a vendor-relationship breach warranting contractual remediation (committee escalates to the institution's vendor-management framework), OR (iii) a model-substitution event with regulatory implications (committee escalates to legal counsel for FFIEC / OCC / FRB notification per the institution's standard model-substitution disclosure policy).

**MRM committee handoff for case (c).** The handoff package the IR program presents to the MRM committee includes: the chain entries showing the discrepancy (the `audit.routing.*` entry naming Provider A and the corresponding LLM-call entry naming Provider B); the vendor's per-call log for the affected calls; the institution's policy-version record from `audit.routing.policy_version` covering the window; the institution's vendor-contract clauses governing provider substitution; and the IR program's preliminary determination of which subordinate disposition (i, ii, or iii above) the committee should evaluate. The committee's written disposition becomes part of the institution's MRM evidence record and feeds the next P-33 working-paper review.

**Containment.** No immediate containment is required for failure-modes (a) or (b) — the chain continues to capture LLM calls and any new routing decisions correctly once the underlying component is remediated. For failure-mode (c) where unauthorized reroute is suspected, the institution's IR program may pause the affected provider integration pending the vendor-relationship investigation; the routing-policy override (`audit.routing.bypass_reason` per spec §4.4.1) records the pause as a manual operator action.

**Remediation.** Failure-mode-specific per the dispositions above. In all three cases the institution updates its CC8.1 control description if the remediation changes the institution's routing-decision capture posture (new chain decorator wiring, new collector configuration, new vendor-management procedure for provider substitution).

**Notification.** Per the clock-start guidance above. For regulator inquiries that triggered Scenario 15, the institution responds with the reconstruction the chain DOES support (the LLM-call entries, any routing entries that were captured, the institution's IR program's investigation findings) and explicitly names the gap and its disposition. For customer disputes, the institution's customer-communication framework determines whether the reconstruction gap is material to the customer's question; counsel decides per the institution's standard customer-dispute response procedure. For model-substitution audits, the gap itself is part of the audit's findings — the audit determines whether the institution's routing-decision capture posture meets the audit's evidence-completeness standard.

**Post-incident review.** Scenario 15 events of failure-mode (b) or (c) get a post-incident review per the standard playbook cadence. The review covers: timeline of when the gap was first observed (the original reconstruction trigger versus the SOC team's P-33 finding versus the chain operations team's monitoring alert); root cause analysis per failure-mode; remediation effectiveness (post-remediation P-33 sample shows restored capture); MRM committee disposition for failure-mode (c) cases; and any updates to the institution's CC8.1 control description, the institution's chain decorator wiring documentation, or the institution's vendor-management procedures.

## 36-hour cyber-incident notification triage matrix

The FFIEC computer-security incident notification rule requires notification of the primary federal regulator within 36 hours of *determining* that a "computer-security incident" has occurred. The determination point — not the alert point — starts the clock. The following triage matrix per chain-detected event type names what triggers the determination and therefore the clock-start:

| Scenario | Default disposition | 36-hour clock-start trigger |
|---|---|---|
| Scenario 1 (chain hash mismatch) | SDK defect (no clock) OR tampering (clock starts) | Clock starts when investigation rules out SDK defect AND confirms in-flight modification of captured events |
| Scenario 2 (Merkle root mismatch) | Storage corruption / recovery (no clock) OR tampering (clock starts) | Clock starts when investigation rules out storage / recovery causes AND confirms ledger-content alteration between capture and verification |
| Scenario 3 (signature verification failed) | Always treated as 36-hour | Clock starts at alert receipt; the determination is the alert itself, given the implication of HSM-key compromise |
| Scenario 4 (master key compromise suspected) | Always treated as 36-hour | Clock starts at the moment the institution determines the suspect signal is credible (not at first alert; intel-driven scenarios may have a longer evaluation period) |
| Scenario 5 (sealing delay > 72h) | Operational SHOULD-notify (no 36-hour clock) | Does not start the 36-hour clock UNLESS the delay is associated with a suspected security incident; then the security incident's clock starts independently |
| Scenario 6 (software-key fallback in production) | Control failure (no clock by default) | Does not start the 36-hour clock UNLESS the fallback's presence is associated with confirmed unauthorized signing |
| **Scenario 7 (key_fingerprint mismatch)** | Triage cases 1-3 (operational): no clock. Triage case 4 (unauthorized substitution): clock starts | Clock starts when the institution determines the mismatch is not attributable to documented operations (rotation-in-flight, restored-backup, tenant-row-restored) AND treats as suspected unauthorized key substitution |
| **Scenario 8 (unknown key_version)** | Cases 1-2 (operational): no clock. Case 3 (tampered key_version): clock starts | Clock starts when investigation indicates the `key_version` field was altered after capture (tampering signal) rather than missing-from-registry (operational signal) |
| **Scenario 9 (audit file truncation)** | Almost never 36-hour | Does not start the 36-hour clock unless the truncation pattern recurs in a way suggesting deliberate corruption |
| **Scenario 10 (backup integrity failure)** | Operational unless gap is unfillable | Clock starts when the institution determines the gap is unfillable AND the affected days remain unverifiable |
| **Scenario 11 (cold-DR fallback / project-side trust-anchor degradation)** | Project-side governance (no clock by default) | Does NOT start the institution's 36-hour clock by itself. Sub-variant (forged rotation notice): clock starts at the determination of forgery. |
| **Scenario 12 (co-signed seal failure, case (e))** | Branch-dependent | Branch (i) published algorithm break: no clock. Branch (ii) per-algorithm signing-key compromise: clock starts at credibility determination (Scenario 4 path). Branch (iii) under investigation: clock starts at investigation-conclusion (bounded 48-72h window). |

The matrix is institutional guidance; counsel decides per incident. The chain produces precise alerts; the institution's IR program produces the determination and the clock-start.

### Edge cases the matrix names explicitly

**Concurrent multi-scenario alerts (rollup rule).** If the institution receives a Scenario 1 (chain hash mismatch) alert AND a Scenario 7 (key_fingerprint mismatch) alert AND a Scenario 8 (unknown_key_version) alert in the same 60-minute window, the institution may be looking at a single root-cause event that triggered three detection paths. The institution's IR program SHOULD treat the constellation as one incident and start the 36-hour clock at the earliest qualifying determination across the constellation, NOT at the union of per-scenario clocks. Treating each scenario's clock independently risks both (a) double-counting incidents in the regulator notification and (b) under-counting where the constellation is the signal even though no individual scenario crossed its threshold. The IR Commander documents the rollup decision and the chosen clock-start in the incident-management system.

**Cross-tenant scope discovery during investigation (vendor-hosted topology).** Scenarios 4 and 7 may discover during investigation that the suspected compromise spans more than one tenant in a multi-tenant deployment (vendor-hosted topology, shared-cloud-HSM at community-bank tier per `00-overview.md` §6.5). The 36-hour clock applies per institution; in a vendor-hosted deployment one investigation may produce N institution-side determinations on different timelines.

**Vendor-vs-institution responsibility split.** The vendor's IR team produces the alert and conducts the cross-tenant investigation. Each institution's IR program receives the alert through the vendor's notification channel (typically a dedicated security-event feed) and produces its own determination based on the vendor's investigation findings AND any institution-specific evidence the institution gathers. The institution's clock starts when the institution determines, not when the vendor first alerted. The vendor's contractual notification SLA (typically 4-12 hours from vendor's determination) is a SEPARATE clock from the institution's 36-hour FFIEC clock — both apply, both are tracked.

**Institution-side decision procedure for cross-tenant scope discovery.** When the vendor's investigation reveals the institution is affected, the institution's IR program executes a documented decision procedure within 6 hours of vendor notification:

1. **Confirm institution scope.** The institution's IR team confirms the affected `(tenant_id, key_version, key_fingerprint)` triples appear in the institution's chain entries during the period the vendor identifies. Without confirmed institution-side evidence, the institution does NOT start its 36-hour clock — vendor notification alone is not the determination event.
2. **Pull the institution's parallel evidence.** The institution's `master.reconciliation_completed` events for the affected period; chain entries with the affected `(tenant_id, key_version)` pair; any institution-side `chain.verification_failure` events that pre-date the vendor notification.
3. **Determine institution-side credibility.** Two outcomes:
   - **Credible institution-side compromise.** Start the 36-hour clock at this determination. Activate Scenario 4 path.
   - **Vendor-side compromise that did not produce institution-side evidence.** No institution-side clock starts. The institution coordinates with the vendor on the vendor's response and notifies its primary federal regulator under the institution's standard third-party-incident-notification framework (typically 72-hour SHOULD per the institution's policy).
4. **Record the determination.** Documented in the institution's incident-management system with the vendor's investigation reference, the institution's parallel evidence references, and the disposition.

The 6-hour decision window is the institution-side floor; counsel may decide faster on clear-cut cases. Beyond 6 hours without a determination, the institution's IR Commander defaults to "credible institution-side compromise" and starts the 36-hour clock to avoid the determination becoming indefinitely deferred.

**Federal-regulator routing per institution charter.** The "primary federal regulator" the FFIEC 36-hour rule notifies depends on the institution's charter type. The institution's IR program SHOULD pre-document which regulator is the primary path:

| Institution charter | Primary federal regulator | CFR citation | State-side cadence |
|---|---|---|---|
| National bank, federal savings association | OCC | 12 CFR Part 53 | N/A (federal-only charter) |
| State member bank, bank holding company, state savings & loan holding company | Federal Reserve (FRB) | 12 CFR Part 225 App. F | State banking regulator per state-specific timing (typical: 24-72 hours; varies by state) |
| State non-member bank, state savings association | FDIC (federal); state banking regulator (state) | 12 CFR Part 304 Subpart C | State banking regulator per state-specific timing |
| Federal credit union | NCUA | 12 CFR Part 748 App. B | N/A (federal-only charter) |
| State credit union | NCUA (federal); state credit-union regulator (state) | 12 CFR Part 748 App. B | State credit-union regulator per state-specific timing |
| Insured depository subsidiary of FBO | OCC, FRB, or FDIC depending on the chartering structure | Per chartering structure (12 CFR Part 53, 225 App. F, or 304 Subpart C) | State counterpart per state of charter |

Multi-charter holding companies notify all applicable primary regulators at the affected entity's chartering structure. The IR Commander confirms the routing per institution charter at the time of the incident; the institution's standing IR documentation pre-records the primary regulator and the secondary state-side counterpart where applicable.

**Dual-algorithm verifier-version timing-overlap (institution-vs-regulator).** When the institution upgrades to a verifier supporting algorithm Y (post-quantum) but the regulator's verifier still supports only algorithm X (Ed25519), and a co-signed seal arrives at examination time, the institution's verifier reports PASS (both algorithms validate) while the regulator's verifier may report PASS-WITH-ANOMALY case (b) "partial-coverage seal" because it cannot evaluate the Y-algorithm signature. Two posture-disagreement scenarios result:

| Posture disagreement | Disposition |
|---|---|
| Institution: PASS (both algorithms valid). Regulator: PASS-WITH-ANOMALY case (b) (only X-algorithm visible). | Coordinate before submitting evidence: the institution's IR program flags the posture mismatch to the regulator's primary contact PRE-examination so the regulator obtains the Y-algorithm-aware verifier from the project supply chain (or the regulator accepts the institution's verifier output as supplementary evidence under documented working-paper procedure). |
| Institution: PASS-WITH-ANOMALY case (b) under its own posture (because the institution's posture lists Y but the SDK didn't co-sign on this seal-day). Regulator: PASS under its X-only posture. | This is a control-completeness issue the institution surfaces to the regulator anyway (the institution's anomaly is per its own posture commitment, not per the regulator's posture). |
| Institution: FAIL case (e) (Y-algorithm signature validates but X-algorithm signature does not — the X-algorithm key has been compromised or X-algorithm has been broken). Regulator: PASS under its X-only verifier (the regulator sees the X-algorithm signature as the only signature and validates it; the regulator does NOT see the Y-algorithm half because its verifier doesn't process `signatures` list). | **Load-bearing case (e) interaction.** The regulator's X-only verifier is structurally blind to the Y-algorithm validation. The institution MUST proactively notify the regulator of the case-(e) failure and supply the institution's case-(e) verifier output as supplementary evidence. This is a known regulator-side blind spot during the dual-algorithm transitional period and is the load-bearing reason the institution coordinates verifier-version compatibility BEFORE the transition begins. The regulator's IT examination program SHOULD upgrade to a Y-algorithm-aware verifier within a documented timeframe (typical: 6 months from the institution's adoption). |

The institution's IR program SHOULD pre-coordinate verifier-version compatibility with the regulator BEFORE the multi-year transitional period begins; the regulator's IT examination program SHOULD subscribe to the project's release notifications to keep verifier-version pace with the institution's adopted posture.

**CIRCIA-only triggers without FFIEC trigger.** Some chain-detected events may rise to CIRCIA's "substantial cyber incident" threshold (72-hour CISA path) without rising to the FFIEC computer-security rule's "computer-security incident affecting safety/soundness" threshold. Specifically:

| Scenario | FFIEC 36-hour trigger? | CIRCIA 72-hour trigger? |
|---|---|---|
| Scenario 5 (sealing delay associated with broader cyber incident at institution) | Maybe (depends on broader incident) | Likely (broader incident is the CIRCIA event) |
| Scenario 9 (truncation pattern across hosts suggesting deliberate corruption) | Maybe | Likely (substantial cyber incident at scale) |
| Scenario 7 (key_fingerprint mismatch — confirmed unauthorized substitution) | Yes | Yes |
| Scenario 4 (master key compromise confirmed) | Yes | Yes |
| Scenario 6 (software-key fallback in production — no unauthorized signing) | No | No |
| Scenario 11 (cold-DR fallback) | No (project-side governance) | Maybe (if dual-compromise impacts institution operations at scale) |
| Scenario 11 sub-variant (forged rotation notice confirmed) | Yes | Yes (substantial cyber incident — attacker attempted trust-anchor substitution) |
| Scenario 10 (backup integrity failure causing verifier failure — gap unfillable) | Conditional (depends on whether unfillable gap impacts safety/soundness) | Yes when ledger-content alteration is suspected; conditional otherwise |
| Scenario 12 branch (i) published algorithm break | No (cryptographic-deprecation event, regulator-coordinated) | No (typically — unless multi-institution scope makes it CIRCIA-relevant) |
| Scenario 12 branch (ii) per-algorithm signing-key compromise | Yes (key compromise) | Yes (substantial cyber incident, key compromise is per-spec definition) |
| Scenario 12 branch (iii) under investigation | Conditional | Conditional |

Counsel works the matrix per incident. The institution's IR program SHOULD have CIRCIA notification on its standard checklist for chain-detected events, even when the FFIEC clock has not started — CIRCIA's "substantial cyber incident" definition is broader than the FFIEC rule's "safety/soundness" framing.

## Communication templates

### Initial regulator notification (within 36 hours)

> [Institution name] is reporting a cybersecurity incident affecting our chain-of-custody implementation for AI agent decisions. The nature of the incident is [brief description]. We are actively investigating. Initial assessment: [scope, customer impact, regulatory implications]. We will provide a follow-up within [X] hours/days.

### Internal stakeholder communication

Concise summary: what happened, what we know, what we don't know, current status, next update time.

### Customer communication (when warranted)

Per the institution's standard customer-communication framework. Chain-detected events do not necessarily require customer communication; consult Legal Counsel.

## Post-incident review

Every chain-detected event of medium or higher severity gets a post-incident review:

1. Timeline of events
2. Root cause analysis
3. Detection effectiveness
4. Response effectiveness
5. Lessons learned and changes required

The review is documented in the institution's incident-management system and informs updates to this playbook.

## CIRCIA and additional notification frameworks

The 36-hour FFIEC computer-security incident notification rule applies to most chain-detected incidents involving suspected tampering or compromise. CIRCIA (Cyber Incident Reporting for Critical Infrastructure Act, 2022) imposes additional 72-hour reporting to CISA for covered incidents. Multiple notification paths apply:

| Path | Trigger | Timing | Recipient |
|---|---|---|---|
| FFIEC computer-security rule | Computer-security incident affecting safety/soundness | 36 hours | Primary federal regulator |
| CIRCIA | Substantial cyber incident at covered entity | 72 hours | CISA |
| State data-breach notification | Customer PII breach | State-specific | State attorneys general |
| CFPB UDAAP | AI decision affecting consumers | Per CFPB rules | CFPB if examination is open |
| Sealing-delay (chain-specific) | Daily seal delayed beyond 72 hours | At threshold | Primary regulator |

The institution's IR program tracks all applicable paths. For chain-detected events, the playbook's notification step considers each path; counsel decides which apply.

## Multi-jurisdiction clock-start handling

The 36-hour FFIEC computer-security incident notification rule and the 72-hour CIRCIA path are the U.S.-federal clocks. Institutions with multi-jurisdiction supervision face additional clocks on the same incident — most prominently, the EU's Digital Operational Resilience Act (DORA) Article 17, which requires reporting of major ICT-related incidents within 24 hours. The DORA clock is tighter than the FFIEC clock; for an incident affecting an institution supervised under both, the DORA clock dominates per incident.

For institutions with multi-jurisdiction supervision, the IR playbook's clock-start triggers apply per the tightest applicable rule:

- **DORA Article 17 — 24 hours** applies to institutions with EU-supervised entities (institutions with EU operations supervised by an EU competent authority — typically a national central bank, the ECB under SSM, or an EU-side prudential supervisor — under DORA's scope).
- **FFIEC computer-security incident notification rule — 36 hours** applies to institutions with US-federal-supervised entities (national banks under OCC, state member banks under the Federal Reserve, state non-member banks under FDIC, federal credit unions under NCUA, and the corresponding charter classes per the federal-regulator routing table above).
- **Institutions with both** operate the 24-hour DORA clock as the binding floor. The FFIEC 36-hour clock continues to run in parallel; the institution's notification path closes the DORA window first and uses the additional 12 hours between DORA closure and FFIEC closure to complete FFIEC notification with the same evidence base.

The institution's IR program pre-documents which regulatory mappings trigger which clock so that the IR Commander does not have to derive the matrix at the moment of incident. The institution's standing IR documentation lists, per institution charter and per jurisdictional footprint:

| Institution's regulatory footprint | DORA 24-hour clock applies? | FFIEC 36-hour clock applies? | Binding floor |
|---|---|---|---|
| US-only (no EU operations) | No | Yes (per charter routing table above) | FFIEC 36-hour |
| EU-only (no US-federal operations) | Yes | No | DORA 24-hour |
| Multi-jurisdiction (US-federal AND EU supervision) | Yes (institution's EU entity is in scope) | Yes (institution's US entity is in scope) | DORA 24-hour |
| US-state-only (e.g. state-chartered bank with no federal supervision and no EU operations) | No | Conditional (state-side cadence per state-banking-regulator timing; FFIEC 36-hour applies only on federal charter or membership) | State-side cadence |

The matrix is institutional guidance; counsel decides per incident. The institution's IR Commander confirms the applicable clocks at incident-determination time and starts each clock at its determination event.

### Worked example: multi-region bank, chain-integrity failure on 2026-06-15 14:00 UTC

Consider a multi-region bank operating in both the United States (national bank charter under OCC) and the European Union (EU subsidiary supervised by an EU competent authority under DORA). At 2026-06-15 14:00 UTC the institution determines that a chain-integrity failure has occurred (Scenario 7 triage case 4: suspected unauthorized key substitution, 36-hour clock starts per the matrix above; this is also a DORA major ICT-related incident under Article 17).

| Clock | Starts | Closes | Recipient |
|---|---|---|---|
| DORA Article 17 (24-hour) | 2026-06-15 14:00 UTC | 2026-06-16 14:00 UTC | EU competent authority for the EU subsidiary (per DORA Article 19 reporting hierarchy) |
| FFIEC computer-security rule (36-hour) | 2026-06-15 14:00 UTC | 2026-06-17 02:00 UTC | OCC (per the federal-regulator routing table for national banks) |
| CIRCIA (72-hour) | 2026-06-15 14:00 UTC | 2026-06-18 14:00 UTC | CISA (substantial cyber incident — key substitution at scale) |

The institution's IR Commander invokes the 24-hour DORA notification path first. The DORA notification closes by 2026-06-16 14:00 UTC. The FFIEC 36-hour path then closes by 2026-06-17 02:00 UTC, giving the institution another 12 hours after DORA closure to complete FFIEC notification with the same evidence base (chain-detection alerts, verifier output, IKM-roster reconciliation evidence, the institution's parallel investigation findings). The CIRCIA path closes 24 hours after the FFIEC path with the same evidence base extended for CISA's substantial-incident framing.

The institution's IR documentation records the three notification timestamps, the three notification reference numbers issued by the three recipients, and the cross-references between them. The institution's primary federal regulator and the EU competent authority each receive parallel notifications with consistent evidence; counsel coordinates the cross-jurisdictional messaging so the institution's narrative is consistent across both supervisors.

### Reference

DORA Article 17 (Reporting of major ICT-related incidents and significant cyber threats) is the binding text for the 24-hour clock; DORA Article 19 governs the reporting hierarchy and the EU-side authority routing. The institution's legal counsel maintains the DORA reference materials alongside the FFIEC computer-security incident notification rule (12 CFR Parts 53, 225 App. F, 304 Subpart C, 748 App. B per charter) and CIRCIA (6 USC §681b) so the IR Commander has the citation set at incident-determination time.

## HSM tamper-detection integration

When the HSM emits a tamper-detection event (vendor-specific; FIPS 140-2 L3 hardware tamper detection), the chain operations team treats the event as Scenario 3 (signature verification potentially compromised) preemptively:

1. The institution's HSM monitoring routes tamper events to the chain operations team
2. The chain operations team activates Scenario 3 — pause seal job; notify CISO
3. Investigate whether the tamper event is legitimate (operator error, legitimate maintenance) or malicious
4. If malicious, treat as HSM compromise; rotate signing key; notify regulator

The integration is institution-specific based on HSM vendor; the operator guide (`docs/operator-guide.md`) documents the institution's specific integration.

## Long-dwell adversary considerations

Long-dwell adversaries (nation-state APTs, persistent threats) may be present in the institution's infrastructure for extended periods (6–18 months) before action. The chain detects modification at action time; long-dwell observation does not produce chain-detectable signals.

The chain composes with the institution's broader detection posture for long-dwell:

- Behavioral monitoring (UEBA) on the chain access patterns
- Anomalous-access detection on the master-key custodian
- Reconciliation cadence (weekly per spec) bounds the visibility window for any compromise
- Cross-correlation with cyber-threat intelligence feeds

When long-dwell is suspected (intel feed, anomalous-access signal), the institution treats as IR Scenario 4 (master-key compromise suspected) to preempt the action phase.

## Opaque agent integration

For commercial AI agent products that don't expose the events the chain captures, the institution wraps the agent's API in a process the institution controls and that calls the SDK. The wrapper:

1. Intercepts the agent's input
2. Records the input in the chain
3. Calls the agent
4. Records the agent's output
5. Returns to the caller

The wrapper may not capture intermediate events the agent produces internally; the institution documents this gap in its control description. For most regulatory purposes, capturing input and output is sufficient. For deeper integration, the institution requests SDK-level integration from the vendor.

## Updates

This playbook is reviewed and updated:

- Annually
- After every medium-or-higher severity incident
- When the spec version is bumped
- When the institution's IR framework is materially updated
- When new threat categories are identified (e.g., quantum readiness drills)
