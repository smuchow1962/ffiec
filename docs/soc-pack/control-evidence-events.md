# Control-evidence operational events

> **What this doc is.** The schema for operational events the ledger emits as control evidence. SOC and examination teams consume these as testable, mechanically comparable evidence.

## Why this matters

A SOC engagement tests that controls operated effectively over the reporting period. The evidence is the operational events the implementation emits. Without a standardized schema, every SOC engagement reinvents how to consume the events. This document standardizes the schema so consumption is mechanical.

## Schema versioning

The operational-events JSON Schema (`docs/operational-events.schema.json`) versions independently of the chain-of-custody spec. The schema's `$id` carries the schema version: `https://ffiec-chain-of-custody.org/schemas/operational-events/v{MAJOR}.{MINOR}.json`. Implementations and SOC consumers pin to a schema version in their control description.

**Versioning convention (semver, major.minor only — no patch).**

- A **minor** version increment adds new event types or adds optional fields to existing events. Backward compatibility is preserved for v1.x: a v1.0 consumer reading a v1.x event (x ≥ 0) continues to validate the events it knows about and ignores any new event types or new optional fields.
- A **major** version increment removes an event, removes a required field from an existing event, or makes a previously-optional field required. Major-version changes are NOT backward compatible. Major-version increments are coordinated with a chain-of-custody spec major-version increment (the spec gates whether a major schema break is acceptable).

**Backward-compatibility commitments for the v1.x line.**

- **No event removal.** Once an event type is in the schema at v1.0, it remains in the schema across the entire v1.x line. An event whose semantics need to change is renamed (new event type added; old event type left in the schema with a `deprecated` annotation). This holds for all 19 event definitions the v1.0 schema names — the 17 events spec §10.2 enumerates plus the two institution-defined events the SOC pack normates: `cold_dr.dryrun_attestation_consumed` (cold-DR dry-run consumption per `supply-chain.md`) and `mirror.reconciliation_completed` (re-signing mirror bridge invariant per `supply-chain.md`).
- **No field-type tightening.** A field declared `string` at v1.0 stays `string` across v1.x. A field declared optional at v1.0 stays optional across v1.x.
- **Field additions are additive.** New fields added in v1.x minor releases are optional. Producers are free to emit them; consumers are free to ignore them. The `additionalProperties: true` setting on the per-event `fields` block is the mechanism.
- **Vendor-extension stability.** The `vendor.*` namespace remains reserved for vendor-defined events across the v1.x line; the spec will not introduce a `vendor.<name>` event type that collides with an in-use vendor extension.

**Schema-version vs spec-version compatibility matrix.**

| Schema version | Compatible spec versions | Notes |
|---|---|---|
| v1.0 | v1.0 (any v1.0-* iteration) | Initial v1.0 schema; covers the 17 events spec §10.2 enumerates plus the two institution-defined events (`cold_dr.dryrun_attestation_consumed`, `mirror.reconciliation_completed`), plus the `vendor.*` extension namespace. |
| v1.1 (future) | v1.0, v1.1 | Adds events the v1.1 spec introduces; v1.0 consumers continue to operate against v1.1 producers (additive only). |
| v1.x (future) | v1.0 through v1.x | Same additive contract. |
| v2.0 (future) | v2.0 (gated on a v2.0 spec) | A breaking change requires a spec major-version increment; institutions plan migration through the spec's normal v2.0 adoption window. |

Institutions running a v1.0 SDK + v1.0 ledger + v1.x verifier consume operational events the verifier-side consumer recognizes (the v1.x consumer knows the v1.0 event set and treats unknown events as vendor extensions or future-version events per its configured policy). The reverse — v1.x SDK + v1.0 verifier — is also safe for the v1.x line: v1.0 consumers ignore the additive content. The institution's control description names the schema version it operates and the verifier-side consumer version it ingests.

## Schema

All operational events conform to the JSON Schema published at `docs/operational-events.schema.json`. The shape is:

```json
{
  "event": "<event_name>",
  "timestamp": "<RFC 3339 UTC>",
  "tenant_id": "<string>",
  "correlation_id": "<string>",
  "fields": {
    "<event-specific fields>"
  }
}
```

The `event` field uses a dotted namespace: `ledger.startup`, `seal.job_completed`, `chain.verification_failure`, etc. The `fields` block carries per-event detail.

## Standard events

### Ledger lifecycle

```json
{
  "event": "ledger.startup",
  "timestamp": "2026-04-01T00:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "process-7f3a...",
  "fields": {
    "version": "v1.0.4",
    "config_hash": "9f86d081...",
    "hsm_cluster": "cluster-acme-prod-1",
    "spec_version": "v1.0"
  }
}
```

```json
{
  "event": "ledger.hsm_session_opened",
  "timestamp": "2026-04-01T00:00:01Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "process-7f3a...",
  "fields": {
    "hsm_cluster": "cluster-acme-prod-1",
    "key_label": "tenant-acme-seal"
  }
}
```

### Seal-job lifecycle

```json
{
  "event": "seal.job_started",
  "timestamp": "2026-04-02T01:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "seal-2026-04-01",
  "fields": {
    "seal_date": "2026-04-01",
    "trigger": "scheduled"
  }
}
```

```json
{
  "event": "seal.job_completed",
  "timestamp": "2026-04-02T01:03:42Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "seal-2026-04-01",
  "fields": {
    "seal_date": "2026-04-01",
    "event_count": 142087,
    "merkle_root": "a7f5e391...",
    "signed_at": "2026-04-02T01:03:42Z",
    "duration_ms": 222000,
    "late_binding_count": 2,
    "cadence": "daily",
    "key_versions": [3],
    "hkdf_inputs_digest": "6f8a5005...",
    "format_version": "v1",
    "algorithm": "ed25519",
    "dev_mode": false
  }
}
```

**Field vocabulary note.** Earlier rounds of this spec used `master_version` (string) as a custodian-side label for the IKM generation. The v1.0-rework normalises the per-entry vocabulary to `key_version` (integer ≥ 1; recorded on every chain entry) and the seal-record vocabulary to `key_versions` (list of integers; carries multi-version state on rotation days). Implementations migrating from the older `master_version` field SHOULD emit `key_versions` going forward; the seal record's `key_versions` is the field the verifier and auditor consume. Custodian-side operational labels remain whatever the institution's KMS / HSM operator framework uses and are not directly consumed by the chain.

```json
{
  "event": "seal.job_failed",
  "timestamp": "2026-04-02T01:04:02Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "seal-2026-04-01",
  "fields": {
    "seal_date": "2026-04-01",
    "error": "hsm_unavailable",
    "retry_count": 3,
    "next_retry_at": "2026-04-02T01:09:02Z"
  }
}
```

### Chain integrity events

```json
{
  "event": "chain.verification_failure",
  "timestamp": "2026-04-02T14:23:18Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "ingest-batch-9f3a...",
  "fields": {
    "run_id": "r_a3f29b71c",
    "seq": 17,
    "key_version": 3,
    "key_fingerprint": "b94c1a77b40bf5106c66ca6c1c1b4989",
    "format_version": "v1",
    "step": 9,
    "reason": "payload_hash_MAC_mismatch",
    "detail": "recomputed=8a4b... stored=9c2d..."
  }
}
```

The `step` field references the spec §7 procedure step that produced the failure (1-12). For `step=8` failures (key_fingerprint mismatch), the `detail` field carries `expected_fingerprint=...; recorded_fingerprint=...` so the SOC team and incident responder can map the failure directly to the IKM-roster entry that was wrong.

```json
{
  "event": "chain.verification_failure",
  "timestamp": "2026-04-15T03:14:22Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "ingest-batch-d4e7...",
  "fields": {
    "run_id": "r_b7c3f8a2d",
    "seq": 1,
    "key_version": 1,
    "key_fingerprint": "2eede65f0f764c97eaf3b3f306a48537",
    "format_version": "v1",
    "step": 8,
    "reason": "key_fingerprint_mismatch",
    "detail": "expected=b94c1a77b40bf5106c66ca6c1c1b4989 recorded=2eede65f0f764c97eaf3b3f306a48537; investigate IKM roster row (tenant=tenant_acme_prod_us_east_1, key_version=1)"
  }
}
```

### Verifier-invocation lifecycle

The `verifier.run_completed` event is emitted by the institution (or by the examiner running the verifier on behalf of the institution) every time the verifier runs against the ledger. The event closes the meta-custody gap — a Daubert expert asked "how do we know the verifier output you're relying on is itself authentic and unaltered?" answers from this event stream. The event records the verifier's binary hash (so a later examiner can confirm the same binary was used), the command-line flags (so the verifier's mode — strict, posture, master-key path — is documented), the operator's identity, and the host the verifier ran on.

```json
{
  "event": "verifier.run_completed",
  "run_at": "2026-05-07T14:30:00Z",
  "verifier_version": "1.0a-amendment.0",
  "verifier_binary_sha256": "abc123...",
  "command_line_args": ["--strict", "--posture=ffiec", "--master-key=hsm:tenant-1"],
  "run_by": "operator-id-123",
  "run_on_host": "verifier-host-01.fdic-examiner.gov",
  "result": "PASS",
  "affected_tenant_ids": ["tenant-1", "tenant-2"],
  "coverage_period_start": "2026-05-01T00:00:00Z",
  "coverage_period_end": "2026-05-07T00:00:00Z"
}
```

**Field semantics.**

- `run_at` — the timestamp the verifier completed its run.
- `verifier_version` — the verifier's reported version string (the `--version` output).
- `verifier_binary_sha256` — the SHA-256 hash of the verifier executable as run. The hash supports reproducible-build verification: a later examiner can confirm the verifier binary recorded here matches the binary they have in hand.
- `command_line_args` — the flags the verifier was invoked with. The flags document the mode: `--strict` versus standard mode, the deployment posture (`--posture=ffiec`), the master-key access path (`--master-key=hsm:...` versus `--master-key=file:...` versus omitted for structural-only verification).
- `run_by` — the operator's identity. For institution-emitted events this is the operator's account identifier; for examiner-emitted events this is the examiner's certificate identifier.
- `run_on_host` — the hostname the verifier ran on. Confirms the verifier was run in the institution's custody (or in the examiner's documented environment).
- `result` — one of `PASS`, `FAIL`, `PASS_STRUCTURAL_ONLY` (key-bound verification skipped because `--master-key` was omitted), mapping to the spec's exit-code semantics per spec §10.12.
- `affected_tenant_ids` — the tenant identifiers the verifier ran against.
- `coverage_period_start` and `coverage_period_end` — the start and end of the ledger window verified.

**Description.** Institution-emitted (or examiner-emitted) operational event documenting that the verifier was run, by whom, on what hosts, with what flags, and what coverage. This event becomes the institution's audit trail for the verifier itself, closing the meta-custody gap (Daubert "how do we know the verifier output is itself authentic?"). The event is the institution's load-bearing evidence that a specific verifier output came from a specific binary, run by a specific operator, on a specific host, against a specific coverage period.

**Cross-reference.** `docs/litigation-support.md` §10 (chain-of-custody for the chain), spec §10.12 (verifier exit codes — `result` field maps to the spec's exit-code semantics), `docs/incident-response-playbook.md` §"Forensic preservation for Critical-severity scenarios" (the preservation step that captures this event alongside the chain-state snapshot).

### HSM operations

```json
{
  "event": "hsm.operation_success",
  "timestamp": "2026-04-02T01:03:42Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "seal-2026-04-01",
  "fields": {
    "op": "sign",
    "key_label": "tenant-acme-seal",
    "latency_ms": 47
  }
}
```

```json
{
  "event": "hsm.operation_failure",
  "timestamp": "2026-04-02T01:04:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "seal-2026-04-01",
  "fields": {
    "op": "sign",
    "key_label": "tenant-acme-seal",
    "error": "session_timeout"
  }
}
```

### Configuration

```json
{
  "event": "config.reload",
  "timestamp": "2026-04-15T09:30:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "config-reload-09",
  "fields": {
    "config_hash": "a3f29b71...",
    "source": "operator_signal"
  }
}
```

### Master key rotation

```json
{
  "event": "master_key.rotated",
  "timestamp": "2026-04-15T03:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "master-rotation-2026-04-15",
  "fields": {
    "old_version": "v3",
    "new_version": "v4",
    "custodian": "aws-cloudhsm-cluster-1",
    "rotation_id": "rot-2026-04-15-001"
  }
}
```

```json
{
  "event": "master_key.rotation_observed",
  "timestamp": "2026-04-15T03:14:22Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "ledger-observation-9c2d",
  "fields": {
    "old_key_version": 3,
    "new_key_version": 4,
    "first_seen_key_fingerprint": "2eede65f0f764c97eaf3b3f306a48537"
  }
}
```

### Audit file truncation

```json
{
  "event": "audit_file.truncation_detected",
  "timestamp": "2026-04-15T14:23:18Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "verifier-run-9c2d-...",
  "fields": {
    "file_path": "audit/pid-12345/tenant_acme_prod_us_east_1/r_a3f29b71c/audit.1.ndjson",
    "writer_pid": 12345,
    "writer_host": "agent-host-7.us-east-1.acme.bank",
    "last_complete_seq": 17,
    "detected_by": "verifier",
    "detected_at": "2026-04-15T14:23:18Z",
    "recovery_attempted": true,
    "recovery_source": "sdk_local_sqlite_buffer",
    "recovery_outcome": "complete"
  }
}
```

The verifier emits this event when `read_audit_file` refuses a file whose last byte is not `\n` per spec §4.1 mid-write truncation refusal. The `recovery_*` fields document whether the institution recovered the dropped event(s) per IR Scenario 9; `recovery_outcome` is one of `complete`, `partial`, or `unrecoverable`. An `unrecoverable` outcome triggers documentation as an integrity-control failure per IR Scenario 9 notification guidance.

### Regulator-fingerprint lifecycle (institution-defined; per `09-threat-model.md` §2.9 reception procedure)

These three events are emitted by the institution when the regulator publishes a new public-key fingerprint (typically because the institution rotated its tenant signing key per `04-hsm-custody.md` §3.2 or because the regulator's fingerprint-storage rotated). The institution's reception procedure is documented in the threat-model doc; the operational events provide the audit-evidence trail.

```json
{
  "event": "regulator_fingerprint.rotation_received",
  "timestamp": "2026-04-15T09:30:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "fingerprint-rotation-2026-04-15-001",
  "fields": {
    "received_at": "2026-04-15T09:30:00Z",
    "received_via": "regulator_encrypted_messaging",
    "regulator_signing_identity": "OCC-IT-EX-Cert-2026",
    "notice_artifact_sha256": "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
  }
}
```

```json
{
  "event": "regulator_fingerprint.rotation_validated",
  "timestamp": "2026-04-15T09:45:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "fingerprint-rotation-2026-04-15-001",
  "fields": {
    "validated_at": "2026-04-15T09:45:00Z",
    "validation_paths_attempted": ["gpg_signature_verify", "cross_channel_cross_check"],
    "validation_results": {"gpg_signature_verify": "PASS", "cross_channel_cross_check": "PASS"},
    "overall_validation": "PASS"
  }
}
```

For a forged-notice failure: `overall_validation: FAIL` triggers IR Scenario 11 sub-variant (trust-anchor reception failure).

```json
{
  "event": "regulator_fingerprint.installed",
  "timestamp": "2026-04-15T10:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "fingerprint-rotation-2026-04-15-001",
  "fields": {
    "installed_at": "2026-04-15T10:00:00Z",
    "old_fingerprint_archived_at": "2026-04-15T10:00:00Z",
    "new_fingerprint_active_from": "2026-04-15T10:00:00Z",
    "change_management_record_id": "CM-2026-04-15-014"
  }
}
```

### Cold-DR-key dry-run consumption (institution-defined; per `supply-chain.md` §"Institution-side consumption of cold-DR-key dry-run attestation")

Annual event recording the institution's consumption of the project's `KEY-DR-DRYRUN-{year}.asc` attestation. Evidence the institution exercises the cold-DR fallback path; pairs with IR Scenario 11.

```json
{
  "event": "cold_dr.dryrun_attestation_consumed",
  "timestamp": "2026-09-15T14:00:00Z",
  "tenant_id": "(institution-wide; not per-tenant)",
  "correlation_id": "cold-dr-dryrun-2026",
  "fields": {
    "attestation_year": 2026,
    "attestation_artifact_sha256": "1d0d5c2cfdcec18ff8b21e9c10ffe4e7c5193623230...",
    "signing_identity_validated": "project_release_management_cosign",
    "consumption_outcome": "PASS",
    "consumed_by": "chain-ops@acme.bank",
    "archived_at": "2026-09-15T14:05:00Z"
  }
}
```

A `consumption_outcome: FAIL` (attestation signature does not validate, or institution cannot retrieve the attestation within 30 days of the project's annual dry-run window) triggers IR Scenario 11 with the appropriate trigger variant.

### Master key retired

```json
{
  "event": "master_key.retired",
  "timestamp": "2026-09-15T03:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "master-retirement-2026-09-15",
  "fields": {
    "key_version": 1,
    "retired_at": "2026-09-15T03:00:00Z",
    "retention_period_days": 2557,
    "retention_approved_by": "chain-ops@acme.bank",
    "chain_entries_referencing_remaining": 0
  }
}
```

Emitted when an IKM is removed from the registry per spec §10.9 retention rule. The `chain_entries_referencing_remaining` field MUST be `0` at the moment of retirement; a non-zero value indicates premature retirement and is a control failure (the verifier will subsequently report `unknown_key_version` for the affected events). The institution's retention-control owner approves the retirement, recorded in `retention_approved_by`.

### Reconciliation

```json
{
  "event": "master.reconciliation_completed",
  "timestamp": "2026-04-07T09:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "recon-week-14",
  "fields": {
    "period": "2026-04-01/2026-04-07",
    "key_versions_observed": [3, 4],
    "key_fingerprints_observed": [
      "b94c1a77b40bf5106c66ca6c1c1b4989",
      "2eede65f0f764c97eaf3b3f306a48537"
    ],
    "fingerprint_unmatched_count": 0
  }
}
```

The reconciliation matches every `(tenant_id, key_version, key_fingerprint)` triple observed in chain entries during the period against the institution's IKM roster's expected fingerprint for each `(tenant_id, key_version)` pair. `fingerprint_unmatched_count > 0` is a high-priority alert: the IKM behind some chain entries does not match the IKM the institution thinks is on file, which means either a botched rotation, a restored-from-wrong-backup, or a cross-tenant key swap has occurred. Spec §7 step 8 catches this at the verifier; this event is the operational evidence trail for the SOC and FFIEC examiner (audit-procedures P-6).

### Cross-region replication reconciliation (Pattern A active-active multi-region only)

The `master.cross_region_replication_completed` event is emitted by each region participating in a Pattern A active-active multi-region deployment per spec §10.15 Pattern A and `docs/design/00-overview.md` §6.4. Pattern A pins the seal for a `tenant_id` to one **seal region**; other regions are **replication regions** that ship their events to the seal region before seal-time. Each replication region emits one event per `(tenant_id, seal_date)` recording its per-region event count, the seal region the events were replicated to, the replication completion timestamp, and the discrepancy (if any). The seal region's actual event count for the tenant-day MUST equal the sum of `events_replicated_to_seal_region` across all source regions.

```json
{
  "event": "master.cross_region_replication_completed",
  "emitted_at": "2026-05-07T00:30:00Z",
  "tenant_id": "tenant-acme-prod",
  "seal_date": "2026-05-06",
  "source_region": "us-east-1",
  "seal_region": "ap-south-1",
  "events_emitted_in_source_region": 145203,
  "events_replicated_to_seal_region": 145203,
  "replication_completed_at": "2026-05-07T00:25:00Z",
  "replication_lag_max_ms": 142,
  "discrepancy": 0
}
```

**Field semantics.**

- `emitted_at` — the timestamp the source region emitted the event (typically just after the source region confirms its replication batch landed at the seal region).
- `tenant_id` — the canonical `tenant_id` the institution operates under Pattern A. Pattern A uses one `tenant_id` across all regions; Pattern B's per-region `tenant_id` shape does not emit this event.
- `seal_date` — the UTC seal-day the replication evidence covers.
- `source_region` — the region the event was captured in. The replication region emits the event; the seal region also emits one with `source_region == seal_region` to record the events captured locally.
- `seal_region` — the region designated as the seal region for this tenant. Documented in CC8.1; one region per tenant.
- `events_emitted_in_source_region` — count of chain entries the source region's local ledger captured for the tenant-day.
- `events_replicated_to_seal_region` — count of chain entries the source region successfully replicated to the seal region by seal-time.
- `replication_completed_at` — the timestamp the source region's last replication batch landed at the seal region's ledger.
- `replication_lag_max_ms` — the maximum end-to-end replication lag observed during the day (institution-defined measurement; typical: the maximum `received_at(seal region) - received_at(source region)` across the day's events). Operational signal; not load-bearing for the integrity invariant.
- `discrepancy` — `events_emitted_in_source_region - events_replicated_to_seal_region`. A non-zero value (events emitted but not replicated by seal-time) is a control-completeness finding routed to `audit-procedures.md` P-37 for sampling. The chain integrity holds — the seal accurately seals what's in the seal region's ledger — but the institution's CC8.1 multi-region replication procedure must document the disposition.

**Aggregation invariant (per spec §10.15 Pattern A integrity invariant 5).** The seal region's actual event count for `(tenant_id, seal_date)` MUST equal the sum of `events_replicated_to_seal_region` across all source regions for that pair. SOC procedure P-37 confirms the equality on a per-period sample. A non-matching aggregation is a Pattern A control-completeness finding (the seal-region count does not reconcile against the source-region replication evidence); it is NOT a chain-integrity finding (the seal still accurately seals what's in the seal region's ledger).

**Cross-references.** Spec §10.15 Pattern A integrity invariant 5; `docs/design/00-overview.md` §6.4; `docs/dr-and-resilience.md` "Multi-region resilience"; `docs/audit-procedures.md` P-37 (cross-region replication-completeness sample); `docs/operator-guide.md` "Multi-region operational guidance"; `docs/regulator-pack/finding-language.md` "Multi-region replication-completeness gap".

### Mirror-registry signature reconciliation (re-signing pattern only)

The `mirror.reconciliation_completed` event is emitted at the cadence the institution operates mirror-registry signature-validation reconciliation per `supply-chain.md` §"Container supply chain" → "Mirror continuous-monitoring control (DE.CM-09)". The event applies only to institutions operating the re-signing mirror pattern; institutions on the signature-preserving pattern do not need it (the project-side cosign signature is preserved end-to-end and validated at deployment time, so there is no bridge invariant to test). For re-signing institutions, the event is the operational evidence the SOC team consumes to test the bridge-invariant control.

```json
{
  "event": "mirror.reconciliation_completed",
  "timestamp": "2026-04-07T09:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "mirror-recon-week-14",
  "fields": {
    "pulled_image_count": 47,
    "validated_image_count": 47,
    "unvalidated_image_count": 0,
    "audit_log_uri": "s3://acme-mirror-audit-worm/mirror-validation-2026-04.log",
    "reconciliation_window_start": "2026-04-01T00:00:00Z",
    "reconciliation_window_end": "2026-04-07T00:00:00Z",
    "bridge_invariant_state": "HOLDING"
  }
}
```

**Field semantics.**

- `pulled_image_count` — number of images the mirror pulled from the project registry during the window.
- `validated_image_count` — of those, the number whose project-side cosign signature validated at pull time per the mirror's audit log.
- `unvalidated_image_count` — pulled images with no audit-log validation entry. MUST equal `pulled_image_count - validated_image_count`. Non-zero is a Scenario 3 finding (the mirror silently stopped validating project signatures, which would let the institution's deployment pipeline accept re-signed images that never had a project-side signature behind them).
- `audit_log_uri` — URI of the mirror's WORM audit log the SOC team reads. The log MUST satisfy 7-year retention and WORM integrity per `supply-chain.md` §"Mirror audit-log retention and integrity (re-signing pattern only)". The reconciliation event is metadata; the WORM log is the load-bearing evidence the SOC team samples during the audit procedure.
- `reconciliation_window_start` / `reconciliation_window_end` — the window over which the reconciliation ran. Cadences are institution-defined but typically daily or per-pull; consecutive windows MUST NOT have gaps (a gap means a window of pulls was not reconciled, which is itself a finding).
- `bridge_invariant_state` — `HOLDING` when `validated_image_count == pulled_image_count` AND the audit log integrity is intact (WORM properties confirmed by the institution's standard WORM-monitoring control); `BROKEN` otherwise. A `BROKEN` state routes to IR Scenario 13.

**Worked example: bridge invariant breaking.** A mirror configuration push silently disabled cosign-signature validation at pull time. The next reconciliation window:

```json
{
  "event": "mirror.reconciliation_completed",
  "timestamp": "2026-04-14T09:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "mirror-recon-week-15",
  "fields": {
    "pulled_image_count": 52,
    "validated_image_count": 19,
    "unvalidated_image_count": 33,
    "audit_log_uri": "s3://acme-mirror-audit-worm/mirror-validation-2026-04.log",
    "reconciliation_window_start": "2026-04-07T00:00:00Z",
    "reconciliation_window_end": "2026-04-14T00:00:00Z",
    "bridge_invariant_state": "BROKEN"
  }
}
```

The institution's IR program treats `bridge_invariant_state: BROKEN` as a Scenario 13 trigger: pause new deployments from the mirror, restore the mirror's signature-validation discipline, re-pull and re-validate every image emitted during the broken window, and document the remediation in the institution's working papers.

**Institutional governance.** The SOC team's audit procedure consumes this event as the test of the bridge-invariant control. Specifically: sample N reconciliation windows in the period, confirm each window's `bridge_invariant_state` is `HOLDING`, confirm the cited `audit_log_uri` resolves to a WORM-backed log meeting retention requirements, and confirm any `BROKEN` state is paired with documented IR Scenario 13 remediation. Without the normative event schema, every SOC engagement reinvented how to consume the bridge evidence; with it, consumption is mechanical and consistent across institutions running the re-signing pattern. CUEC-VND-01 (the institution's mirror-registry control) is the institution-side control description that names the cadence and the bridge invariant; this event is the evidence the description points at.

### `chain_kind` cross-reference for operational events

Spec §3 normates a closed `chain_kind` enumeration covering chain-of-custody integrity events and operational events. Operational events emitted under this schema carry `chain_kind = "operational"` to disambiguate from chain-of-custody integrity events. The full enumeration:

| `chain_kind` | Used for |
|---|---|
| `audit` | Application audit events (the default chain-of-custody integrity entry) |
| `model_call` | Chain entry representing an LLM invocation |
| `tool_call` | Chain entry representing a tool invocation |
| `routing` | Chain entry from the routing-event surface (per spec §4.4.1) |
| `translation` | Chain entry from the ECOA translation step (per spec §10.11) |
| `operational` | Chain entry for control-evidence operational events documented in this file |

Every operational event documented in this file emits `chain_kind = "operational"` so SOC sample-comparison procedures, dashboards, and alerting rules can mechanically separate operational evidence from integrity-bearing chain content. The verifier MUST reject any `chain_kind` value not in the enumerated set with `chain_kind out of v1 enumeration at seq N`. Operational events are still chained — the institution's evidence trail benefits from the same MAC + Merkle + HSM coverage that integrity events receive — but the `chain_kind` discriminator names the event class so consumers know which evidence layer they are reading. Cross-reference spec §3 for the full normative text and `operator-guide.md` "`chain_kind` operational note" for the operations-team-facing summary.

### `routing.refused_event` (institution-defined; per spec §4.4.1 new event type)

The institution's ledger emits this operational event when a chain entry under `chain_kind = "routing"` carries `audit.routing.event_type = "audit.routing.refused"` (per spec §4.4.1). The event provides SOC sampling traction on the new refused-routing case without requiring the SOC team to walk every routing chain entry — the operational event is the index. Useful for SOC sampling per `audit-procedures.md` P-33's refused sub-sample.

```json
{
  "event": "routing.refused_event",
  "timestamp": "2026-04-15T14:23:18Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "routing-decision-9f3a-...",
  "fields": {
    "run_id": "r_a3f29b71c",
    "seq": 17,
    "chain_kind": "routing",
    "refusal_reason": "all_circuits_open",
    "policy_version": "routing-policy-v3.2.1",
    "providers_evaluated": ["openai-gpt-4o", "anthropic-claude-sonnet", "google-gemini-pro"],
    "decision_at": "2026-04-15T14:23:18.234567Z"
  }
}
```

The `refusal_reason` field carries one of the five enumerated values per spec §4.4.1: `all_circuits_open` | `no_provider_in_policy` | `quota_exhausted` | `cost_threshold_at_capacity` | `policy_override`. SOC teams sampling P-33's refused sub-sample query the institution's log store for this event, then cross-reference each event back to the underlying chain entry to confirm the chain entry's `audit.routing.refusal_reason` matches the operational event's value. A non-conforming `refusal_reason` value (outside the enumeration) is a control-completeness gap the SOC team escalates per P-33's failure-mode disposition.

### `seal.dev_mode_published_in_production` (verifier-emitted; per spec §10.7 + §4.3 dev_mode binding)

The verifier emits this operational event when it observes `dev_mode = true` on a seal record under `--strict` mode (per spec §10.7 the verifier MUST refuse to validate a `dev_mode = true` seal as a production seal under `--strict`; per spec §4.3 the `dev_mode` field is bound under the HSM signature so the value cannot be silently flipped). The event names the seal that triggered the refusal so the institution's IR program can route to Scenario 6 (software-key fallback used in production).

```json
{
  "event": "seal.dev_mode_published_in_production",
  "timestamp": "2026-04-15T03:14:22Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "verifier-run-9c2d-...",
  "fields": {
    "seal_date": "2026-04-14",
    "verifier_invocation": "verifier walk --strict --posture=ffiec",
    "verifier_version": "v1.0.4",
    "kms_handle_uri_observed": "plaintext-dev",
    "refusal_reason": "dev-mode seal in production verification — refused",
    "step": 12
  }
}
```

The `step` field references spec §7 step 12 (cadence and dev-mode check). The institution's IR Scenario 6 response includes (a) investigating how `dev_mode = true` reached production (a configuration-drift gap, an unfiltered build artifact, or a deliberate bypass), (b) restoring HSM-backed signing for the affected period, (c) re-issuing affected seals from an HSM-backed key. The `dev_mode` field's binding under `sign_payload` (spec §4.3) means the value at the seal record IS the value the HSM signed — an attacker cannot rewrite `dev_mode = true → false` after the fact, so a `seal.dev_mode_published_in_production` event reflects the institution's actual production posture at the moment of seal signing, not a downstream tampering.

The institution emits this event from the verifier toolchain as part of the verifier run; the institution's log-routing pipeline collects verifier-side operational events alongside ledger-side operational events so SOC teams have a unified evidence stream. Cross-reference `regulator-pack/finding-language.md` for the corresponding examination-finding language.

## Event coverage requirements

A conforming implementation MUST emit every event listed above when the corresponding action occurs. Missing events are a defect.

A conforming implementation MAY emit additional events under its own namespace (e.g., `vendor.<event>`), provided the additional events do not collide with the standard namespaces (`ledger`, `seal`, `chain`, `hsm`, `config`, `master_key`, `master`, `routing`).

## Retention

Operational events SHOULD be retained at least as long as the chain events they relate to. If the institution retains chain events for 7 years, operational events SHOULD also retain 7 years.

## Routing

Operational events are routed to the institution's standard observability stack (Splunk, Sentinel, Datadog, etc.). The events contain no event-payload data and no key material; routing them to any backend the institution operates is safe.

## How SOC teams use this

For each control claim in the SOC opinion, the SOC team identifies which operational events provide the evidence:

| Control claim | Evidence events |
|---|---|
| The seal job operated daily | `seal.job_started` and `seal.job_completed` for every day in scope |
| HSM signing succeeded for every seal | `hsm.operation_success` with `op=sign` matching every seal |
| Chain integrity was verified at ingest | Absence of `chain.verification_failure` events; or, when present, evidence of remediation |
| Master key rotation followed procedure | `master_key.rotated` paired with `master_key.rotation_observed` and a documented change-management record |
| Reconciliation operated at the documented cadence | `master.reconciliation_completed` events at the expected interval |

The SOC team queries the institution's log store, filters by event name, and confirms the count and content match the institution's claims.
