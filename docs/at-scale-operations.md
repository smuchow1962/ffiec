# At-scale operations

> **What this doc is.** Operational guidance for institutions running the chain at very large scale (billions of events per day, globally-systemic banks, holding-company complexity). Closes fresh-batch at-scale partials.

## Throughput planning

### Hot-path budget

Per `02-chain-construction.md` §5, the hot-path budget is <2 ms p99 per event. Within that budget, components contribute:

| Component | Budget | At-scale consideration |
|---|---|---|
| JCS encoding | <50 µs | Scales with payload size; large payloads warrant SDK-side payload simplification |
| HMAC-SHA-256 | <10 µs | Constant per event |
| SQLite insert (synchronous=FULL for audit, NORMAL for telemetry) | <500 µs | Disk-bound; SSD required |
| Serialization to OTLP | <100 µs | Scales with payload |

### Per-host throughput ceiling

Single-host audit-event throughput is fsync-bound. Reference numbers:

| Configuration | Audit events/sec/host | Telemetry events/sec/host (batched) |
|---|---|---|
| Single SSD, synchronous=FULL | ~2,000 | ~50,000 |
| NVMe SSD, synchronous=FULL | ~5,000 | ~150,000 |
| Multiple SSD, parallel SQLite | ~10,000 | ~500,000 |
| Hardware-backed key storage | Throughput limited by HSM round-trip latency, typically 1,000–5,000/sec |

### Sharding patterns

For institutions exceeding single-host throughput:

**Per-process sharding.** One SDK instance per agent process; many agent processes per host; many hosts per institution. Most common pattern.

**Per-tenant sharding.** Different tenant_ids on different hosts; useful for institutions with diverse AI use cases.

**Per-business-line sharding.** Different business lines (fraud, credit, customer service) on different infrastructure; aligns with business-line risk-isolation.

For globally-systemic banks at billions/day:

- 1 billion events/day = ~12,000 events/sec
- At 2,000 events/sec/host, 6 hosts active concurrently is sufficient
- Add capacity buffer (typical 3x) for traffic spikes; deploy 18+ hosts
- For burst-tolerance, telemetry events use synchronous=NORMAL with batching, achieving 100K+/sec aggregate

## Ledger-server scaling

The ledger server scales with:

- OTLP receiver throughput (gRPC or HTTP)
- WAL write throughput (PostgreSQL write capacity)
- Hot store query load (read replicas as needed)
- Cold store ingestion (asynchronous; not in hot path)

### Reference deployment for billions/day

| Component | Sizing |
|---|---|
| Ledger receiver pool | 8–16 instances behind a load balancer |
| WAL Postgres | Aurora or RDS with provisioned IOPS; replica per region |
| Hot store | Same Postgres instance; read replicas for query load |
| Cold store | S3 with Glue/Iceberg; ingestion async from WAL |
| HSM cluster | 4+ members per region; cross-region replicated |
| Seal-job pool | 2 instances per tenant for HA |

### Hot path optimization at scale

Beyond the spec defaults, institutions at billions/day may:

- Batch OTLP payloads (multiple events per OTLP request) for receiver efficiency
- Tune Postgres write parameters (commit_delay, wal_writer_delay) within fsync constraints
- Pre-allocate WAL space to reduce growth latency
- Use tablespaces for retention partitioning

These are operational tunings; the spec accommodates them without modification.

## Master-key rotation at scale

### Rolling restart

The standard pattern (`docs/operator-guide.md`): rotate master, force application processes to re-handshake via rolling restart. At thousands-of-processes scale, the rolling restart takes hours. During the window, both master versions are valid; the seal record carries `master_version = "v3,v4"` (`04-hsm-custody.md` §3.2.2).

### Emergency rotation

For master-key compromise scenarios, the rotation must complete faster than the rolling-restart procedure allows.

**Forced-handshake mechanism.** The master-key custodian invalidates session-key issuances bound to the old master. Application processes attempting to use their session key for HMAC computation receive a `key_revoked` response from the custodian on their next handshake; processes that haven't reached out recently continue operating until their next session-key refresh (which would normally occur on a per-process timeout).

For emergency rotation, the institution forces the refresh:

1. Master-key custodian invalidates old-master session keys at the custodian
2. Application processes are signaled to re-handshake (via deployment-orchestration mechanism: SIGHUP, Kubernetes pod restart, etc.)
3. Processes that don't respond within 60 seconds are terminated and replaced
4. The institution monitors the proportion of events under the new master_version; the rotation is complete when 99% of events use the new version

**Parallel-rotation patterns.** For institutions with many processes, the rotation can proceed in parallel waves: 25% of processes per 5 minutes. The institution's deployment-orchestration tools (Kubernetes, Spinnaker, etc.) drive the cadence.

**Documentation.** Emergency rotation is documented in the institution's IR playbook (Scenario 4) with the procedure detailed for the institution's specific deployment.

### Normativity of the 60-second forced-handshake termination

The 60-second termination step in the forced-handshake mechanism is **operational guidance, not a normative requirement of the spec**. The spec normates that emergency rotation invalidates old-master session keys and that the institution drives processes onto the new master; it does NOT normate the wall-clock budget by which non-responsive processes MUST be terminated. The 60-second figure reflects a deployment posture available to institutions whose orchestration is uniformly synchronous (interactive-traffic AI agents, request/response workloads on Kubernetes with a healthy liveness story).

Globally-systemic banks running batch-processing or scheduled-job AI usage routinely operate processes that legitimately do not respond on a 60-second timescale. A nightly fraud-scoring batch that holds a process open for hours, a long-running compliance-review job that cannot be safely interrupted mid-decision, a regional reconciliation worker tied to a market close — these workloads are conformant under the spec and MUST NOT be forced to adopt a posture that would corrupt their work product. The institution operates **compensating controls** for the rotation window in these cases, and the institution's CC8.1 control description names the posture explicitly.

**Compensating-control posture for extended-rotation windows.** An institution whose deployment posture cannot guarantee 60-second forced-handshake operates the following compensating controls during emergency rotation:

1. **Extended rotation window with a documented bounded-compromise commitment.** The institution's IR playbook (Scenario 4 entry) names the maximum window during which both master versions are valid (e.g., 4 hours, 12 hours, 24 hours) and the rationale tied to the longest-running legitimate process class. The window MUST be bounded; "until processes happen to refresh" is not a conformant posture.

2. **Bounded-compromise commitment during the extended window.** The institution's IR plan documents the integrity claim that holds across the window. The recommended commitment shape: "no captured chain entries after rotation start are accepted under the old master beyond \[bounded time T\]; entries captured between rotation start and T under the old master are flagged in the seal record's `master_version = [old, new]` field per `04-hsm-custody.md` §3.2.2 and are subject to the institution's rotation-window IR review." T is the institution-named extended-window upper bound, explicitly tied to the longest-running batch class.

3. **Active monitoring during the extended window.** The institution monitors `master_version_active` per tenant continuously; an alert fires if the old master version is observed beyond T. The institution's SOC team treats post-T old-master observations as a rotation-window anomaly requiring IR review (a process held the old session key longer than the documented window admits).

4. **Documentation of the named workloads.** The institution enumerates the process classes that justify the extended window (batch fraud scoring, scheduled compliance jobs, regional reconciliation workers, etc.) with the maximum legitimate run time of each. The list is reviewed quarterly by the institution's MRM committee; new long-running workloads added during the period trigger a window-bound review.

5. **Trigger to escalate to forced-termination posture.** If the institution determines the extended window cannot be bounded for a specific compromise event (e.g., a confirmed master-key compromise where the old master MUST be invalidated immediately), the institution escalates to forced-termination of all processes regardless of run state. The escalation is the institution's IR-leadership decision and is logged as part of the Scenario 4 IR record.

The 60-second termination remains the recommended posture for institutions whose workloads can adopt it. Institutions operating compensating controls are conformant under the spec; their CC8.1 control description names the extended window, the bounded-compromise commitment, the monitoring covering the window, and the named workloads.

## Reconciliation at scale

The key-fingerprint reconciliation (spec §10.1, `09-threat-model.md` §2.7) joins the IKM custodian's handshake log against the ledger's per-entry `(tenant_id, key_version, key_fingerprint)` observations. At billions-of-handshakes-per-week scale:

### Implementation pattern

1. The IKM custodian emits `handshake.completed` events (per handshake; carries `(tenant_id, key_version)` and the requesting workload identity)
2. The ledger emits per-entry stamps via the chain itself; the reconciliation job derives `(tenant_id, key_version, key_fingerprint)` triples from the events table
3. A reconciliation job runs at the institution's chosen cadence (weekly per spec; can be more frequent)
4. The job recomputes `expected_fingerprint = SHA-256(utf8(tenant_id) || ikm)[:16]` for each observed `(tenant_id, key_version)` triple against the institution's IKM-roster's IKM bytes; mismatches are alerted

### Performance considerations

For billions of handshakes per week:

- The IKM-custodian log stores `(tenant_id, key_version)` per handshake with timestamp; volume is manageable (a few GB/week for billions of handshakes)
- The ledger's observed `(tenant_id, key_version, key_fingerprint)` triples can be derived from the events table by GROUP BY
- The join is efficient if both sides are indexed by `(tenant_id, key_version)`; the fingerprint recompute is cheap (~100 ns per triple)
- Mismatch detection is the alert; the alert volume should be near-zero in normal operation

### Baseline establishment

For institutions with incomplete handshake logs from day-1 (some handshakes weren't logged due to historical operational practices), the reconciliation initially produces non-zero unmatched_count. The institution baselines this:

1. Establish the unmatched_count baseline over 4 weeks of operation
2. Alert on deviation from baseline (typically: alert if unmatched_count exceeds baseline + 2 standard deviations)
3. Tune the alert threshold quarterly

The unmatched_count is a relative metric, not an absolute one. The institution tracks the baseline in its control description.

### Tier-1 reconciliation sampling variant

Spec §10.1 mandates weekly key-fingerprint reconciliation of every observed `(tenant_id, key_version, key_fingerprint)` triple against the IKM roster. At billions of events per day per tenant, the weekly reconciliation has tens of billions of triples and the exhaustive comparison runs longer than a single weekly cadence admits. The spec's intent is that every observed triple is accounted for; tier-1 institutions discharge that intent through a stratified-sampling variant rather than exhaustive enumeration.

**Sampling design.** Tier-1 institutions MAY substitute a stratified sample for the exhaustive reconciliation, subject to the controls below. The default sampling design covers every `(tenant_id, key_version)` pair observed in the period with a per-pair sample of N triples (institution-defined; **10,000 triples per `(tenant_id, key_version)` pair per week** is the recommended default unless the institution's MRM committee approves a different N tied to a documented statistical-confidence target). Strata SHOULD cover all observed pairs in the period; an unsampled pair is not a conformant outcome.

**Statistical commitment.** The institution's sampling design names:

1. The sample size per stratum (10K is the recommended default; institutions may justify a different N).
2. The strata definition (per `(tenant_id, key_version)` pair is the recommended default; finer strata are admissible if the institution's MRM committee approves).
3. The statistical confidence the sample supports (e.g., "99% confidence that the population mismatch rate is below 0.01% if the sample shows zero mismatches"); this confidence is the basis for the institution's CC8.1 control claim.
4. The period over which the sampling design is valid (typically annual, reviewed at the institution's MRM committee meeting cadence).

**Institutional governance.** The sampling variant is governance-gated, not engineering-gated:

1. **MRM committee approval.** The institution's MRM committee approves the sampling design (sample size, strata, confidence target) before it operates in production. The approval is recorded in committee minutes and named in the institution's CC8.1 control description.
2. **SOC sample-of-samples procedure.** The SOC team confirms the sampling design through a tier-1 sample-of-samples procedure: per audit period, the SOC team samples M of the institution's weekly reconciliation runs, re-runs the per-stratum sampling on the same period's data independently, and confirms the institution's run produced equivalent coverage. M is the institution's SOC-engagement size for tier-1 reconciliation testing (typically 4 of 52 weekly runs per audit year).
3. **CC8.1 control description.** The institution's CC8.1 control description names: the sampling design, the period of validity, the MRM-approval reference, the SOC sample-of-samples procedure, and the trigger that escalates to exhaustive reconciliation (below).

**Escalation trigger to exhaustive reconciliation.** A sampled triple that fails reconciliation triggers exhaustive coverage of the implicated `(tenant_id, key_version)` pair for the period. The institution's IR playbook treats the failure as a Scenario-4-class event (key-related anomaly) until the exhaustive run completes and either confirms the failure as isolated or escalates to a key-compromise IR procedure. A second sampled failure within a rolling 90-day window for the same `(tenant_id, key_version)` pair triggers exhaustive reconciliation for that pair for the next four weekly cadences while the institution's MRM committee reviews whether the sampling design remains adequate.

**What the sampling variant does NOT change.** The IKM-custodian log retention, the per-entry observation in the events table, the ledger's GROUP BY derivation, and the per-week cadence are unchanged. The sampling variant changes only the comparison step: tier-1 institutions compare a stratified sample of triples per pair rather than every triple per pair. The spec's underlying intent — "every observed pair has its fingerprint validated against the IKM roster on the cadence" — is preserved by the sampling design covering all observed pairs.

Institutions below tier-1 volumes SHOULD continue exhaustive reconciliation; the sampling variant is an at-scale accommodation, not a baseline relaxation.

## Verifier resource posture at tier-1 volumes

A tier-1 bank captures on the order of 1M events per second at peak. A 7-year retention horizon places approximately **30 trillion events** in scope per institution. The verifier — a single static binary in the spec's reference design — cannot complete a one-year-coverage examination of that population in an 8-hour examination window from a single host walking the chain serially. Tier-1 institutions operate the verifier in a **partitioned, streaming, resumable** posture documented below; the spec accommodates this posture without modification, and the institution's CC8.1 control description names it explicitly.

### Tier-1 deployment expectation

Tier-1 institutions run **partitioned verifications**, not single-binary serial walks. The partitioning dimensions, in priority order:

1. **Per business line / region.** The institution's chain is already segmented by tenant_id at the spec level; tier-1 deployments typically map tenant_ids onto business lines (fraud, credit, customer service, capital markets) and regions (NA, EMEA, APAC). Per-business-line verification runs in parallel; cross-business-line correlation is the institution's analytical task at the holding-company level.
2. **Per tenant-day.** Within a business line, the seal record's tenant-day boundary is the natural verification partition. Each tenant-day's events form a self-contained Merkle tree (per `02-chain-construction.md`); the verifier verifies one tenant-day on one core, and a host with K cores verifies K tenant-days concurrently.
3. **Parallelized across cores within a host, and across hosts within a verifier farm.** A tier-1 verifier deployment is a fleet, not a binary. The fleet's coordinator schedules tenant-day partitions onto worker cores; each worker is a static binary running the streaming verifier mode below.

The combined effect: a tier-1 examination runs hundreds to thousands of tenant-day verifications concurrently. The static binary is unchanged; the deployment topology is what scales.

### Streaming, resumable verifier mode

Tier-1 verifier deployments operate the binary in **streaming mode with resumable progress state**. The verifier serializes its progress at checkpoint intervals so an interrupted examination resumes from the last checkpoint rather than restarting the partition.

**Progress-state serialization shape.** At each checkpoint, the verifier serializes a record naming:

- `tenant_id` — the partition tenant
- `day_boundary` — the tenant-day under verification
- `last_completed_run_id` — the run within the day whose chain has been fully walked
- `last_completed_seq` — the sequence number within `last_completed_run_id` whose entry was the last verified
- `last_verified_seal_root` — the Merkle root the verifier reconstructed up to and including `last_completed_seq`
- `seal_status` — one of `not_yet_reached`, `verified`, `pending` (the seal record is the partition's authoritative checkpoint; `verified` means the reconstructed root matched the signed seal)
- `verifier_version` — the verifier binary version (resumption requires the same verifier; mid-partition verifier upgrades restart the partition from the day boundary)

The progress-state record is written to the verifier's working store at a configurable cadence (default: every 60 seconds of wall-clock time, or every 1M entries verified, whichever comes first). The record is institution-controlled; it is not part of the chain's evidence and does NOT enter the seal.

**Resumption semantics.** On resumption, the verifier reads the progress-state record, re-establishes the chain-walk state at `last_completed_seq + 1`, and continues. Resumption is a constant-time operation regardless of how many entries have already been verified; the verifier does not re-walk the verified prefix.

**Examination-window fit.** Streaming-resumable mode means an 8-hour examination window does NOT need to bound a single uninterrupted verifier run. The institution can split the window across multiple workers, restart workers on host failure without losing prior work, and pause/resume around examiner Q&A without restarting partitions.

### Bounded-memory commitment

The verifier's RAM consumption scales with **concurrent partition count**, NOT with total event count. A worker verifying one tenant-day partition holds:

- The partition's chain-walk state (the rolling `prev_payload_hash`, the running `key_version` lookup, the event-record-being-verified buffer) — bounded constant per partition (~kilobytes)
- The Merkle accumulator state for the partition's seal reconstruction — bounded constant per partition (~kilobytes for the running tree fringe; streaming Merkle accumulates without retaining leaves)
- A small look-ahead buffer for OTLP ingest pipelining — bounded constant per partition (~megabytes, configurable)

A worker host running W concurrent partitions consumes O(W × constant) RAM, NOT O(events × constant). A worker with 16 GB RAM running 64 concurrent partitions consumes ~256 MB on partition state plus the OTLP buffers and OS overhead — well within the host's budget.

The bounded-memory commitment is the load-bearing property that lets tier-1 verification scale to 30 trillion events without scaling the verifier's RAM with the population.

### Examination-window budget

The institution's examination-window budget is the product:

> **events covered = examination window × parallel cores × throughput per core**

Throughput per core for the streaming verifier (HKDF + HMAC + Merkle leaf + canonical-bytes recompute per event) is approximately **50,000 events per second per core** on commodity hardware (reference number from the project's benchmark suite; the institution validates against its own hardware during pilot).

The budget table for typical examination postures:

| Institution tier | Annual event volume | Examination window | Parallel cores | Events covered | Coverage of one-year retention |
|---|---|---|---|---|---|
| Community ($1B–$10B) | ~5B/year | 8 hours | 4 | ~5.8B | 100% (full one-year coverage) |
| Mid-size ($10B–$250B) | ~500B/year | 8 hours | 32 | ~46B | 100% via per-business-line partition (8 BLs × ~6B each) |
| Tier-1 (>$250B) | ~30T/year | 8 hours | 1,024 | ~1.5T | 5% (sample-based examination per business line; full coverage spans multi-day examination cycle) |

For tier-1 institutions, full annual coverage in a single 8-hour window is not the operational target; the examination posture is **stratified sampling within the window** (the examiner samples per-business-line per-tenant-day per-run partitions, the verifier verifies the sampled partitions exhaustively, and the institution's CC8.1 control description names the sampling posture). The verifier scales to support either full coverage (over a multi-day examination) or sampled coverage (within an 8-hour window); the institution's deployment posture and the examiner's sampling design jointly determine the operational mode.

### Step 12a budget at tier-1 volumes

Spec §7 step 12a (GenAI model identifier completeness check) runs **inline during the per-event walk** and is normative. Step 12a performs, for each chain entry:

1. **One namespace-prefix scan** of the entry's captured attributes for keys matching the `gen_ai.*` namespace
2. **Two field-presence checks** when any `gen_ai.*` attribute is present: `gen_ai.request.model` MUST be present; `gen_ai.response.model` MUST be present (the verifier asserts both; the spec prohibits deferral of either check to a second pass)

**Per-event cost of step 12a.** On commodity hardware against a typical attribute-bag size (~20 keys per entry), the namespace-prefix scan is approximately **300 nanoseconds** (a single pass over the attribute-bag's keys, byte-comparing each key's prefix against the literal `gen_ai.`). The two field-presence checks are approximately **50 nanoseconds each** (hash-table lookups by key on the same attribute-bag). The combined per-event 12a cost is approximately **400 nanoseconds**.

**Comparison to the baseline per-event verification cost.** The per-event baseline (HKDF + HMAC + Merkle-leaf hash, per `02-chain-construction.md` §5) is approximately **18 microseconds per event** on the same hardware. Step 12a adds approximately 400 nanoseconds, or **~2.2% of baseline per-event cost**. The cost fits inside the existing per-event verifier budget without measurable latency-percentile impact at tier-1 throughput; the budget table above is unchanged in event-coverage capacity by step 12a's introduction.

**Why the cost lands on the hot path.** Spec §7 step 12a's deferral-prohibition language ("MUST NOT defer it to a second pass") is not stylistic; it reflects the operational cost of deferral at scale. Deferring step 12a to a second pass would require:

1. Persisting per-entry intermediate state across passes (queue overhead growing linearly with event count)
2. Coordinating two passes' progress state (the streaming-resumable contract above would need to track both passes' positions)
3. Communicating completeness results between passes (an observability-pipeline question that introduces a new failure mode: the second-pass result is lost or arrives after the first-pass output is consumed by the examiner)

The per-event 400-nanosecond cost is an order of magnitude cheaper than the queue-and-coordinate overhead of a deferred second pass. The spec's choice to land step 12a on the hot path is the lower-cost design at tier-1 volumes; tier-1 institutions plan capacity against the combined 18.4 microseconds-per-event budget without further accommodation.

**Tier-1 capacity-planning recommendation.** Tier-1 institutions size verifier-fleet core counts against the combined per-event cost (baseline + step 12a) and the examination-window targets in the budget table above. The 2.2% delta does NOT warrant a verifier-fleet expansion; it is absorbed inside the existing per-core throughput estimate.

## Long-tail retention

For litigation hold or extended-retention scenarios (10+ years post-decision):

### Cryptographic agility

The chain's algorithm-rotation provision (`09-threat-model.md` §2.8) admits new algorithms. Past records remain verifiable under their original algorithm; new records use the new algorithm. For retention beyond a single algorithm's expected secure lifetime:

1. Past seal records carry their original algorithm identifier
2. The verifier dispatches based on the algorithm identifier per seal
3. As long as the verifier knows the original algorithm (Ed25519 is unlikely to be forgotten in any reasonable timeframe), past seals remain verifiable

For paranoid long-tail scenarios (50-year retention against quantum break), the institution can augment past seals with post-quantum signatures at the time the new algorithm is available; the seal record format supports multiple signatures per seal.

### HSM lifecycle

HSM hardware retires every 5–10 years. Past keys can move to new HSMs (key transfer is HSM-supported) or be exported in encrypted form for archival.

For verification (which uses the public key only), HSM hardware is irrelevant; the verifier uses any conforming Ed25519 implementation. Past public keys remain valid indefinitely as long as the institution retains them (in the tenant key registry).

### Long-tail DR

For litigation-hold scenarios extending beyond the institution's standard backup retention:

1. The institution promotes the affected scope to long-tail retention (separate retention policy)
2. Cold-store archives are retained for the longer period
3. Public keys covering the period are pinned in the registry
4. The institution's legal team monitors the hold's status

## Quantum-readiness drill

The 30-day spec-patch SLA is a project commitment. The institution's preparation drills its readiness:

### Annual quantum drill

1. The institution simulates a quantum-readiness trigger (NIST publishes deprecation timeline, hypothetical break demonstrated)
2. The institution's IR playbook is exercised at table-top level
3. The institution evaluates: HSM vendor quantum-readiness status, internal-process readiness for algorithm rotation, communication plan with regulator
4. Lessons learned feed updates to the institution's quantum-readiness plan

### Quantum-preparation checklist

Before a quantum trigger fires:

- [ ] The institution has a documented quantum-readiness plan
- [ ] The HSM vendor has a public quantum-readiness roadmap
- [ ] The institution has identified the post-quantum algorithm options that match its FIPS posture
- [ ] The institution has run table-top exercises of the rotation
- [ ] The institution has communicated the quantum-readiness status to its primary regulator (typically annually)

## Holding-company scale

For holding companies with complex subsidiary structures:

### Per-subsidiary tenant_ids

The standard pattern: each subsidiary has its own tenant_id. The verifier runs per subsidiary. Cross-subsidiary correlation is at the holding-company level via business identifiers.

### Cross-tenant references

When subsidiary A's AI agent calls subsidiary B's tool:

- Subsidiary A captures the tool-call as a chain event in tenant_A's chain
- Subsidiary B captures the tool-execution as a chain event in tenant_B's chain
- Correlation is via shared business-level identifier (e.g., `correlation_id` in the canonical payload)

The chain doesn't model cross-tenant references directly; the correlation is the institution's analytical responsibility.

### Holding-company examination

The Federal Reserve examines the holding company; multiple subsidiaries' chains are in scope. Examination evidence includes per-subsidiary verifier output. Cross-subsidiary aggregation is the EIC's task; the v1.1 `verifier portfolio` subcommand will streamline.

## Operational metrics for at-scale monitoring

Beyond the standard metrics in `06-ledger-server-design.md` §7.2:

| Metric | At-scale relevance |
|---|---|
| `ledger_otlp_received_total{status}` | Per-tenant per-host; aggregate to institution |
| `chain_construction_duration_seconds` | p99 per host; alert on degradation |
| `seal_event_count{tenant_id}` | Per tenant per day; spot anomalous days |
| `master_version_active` | Per tenant; should be one normally; two during rotation window |
| `reconciliation_unmatched_count` | Per tenant per period; alert on baseline deviation |
| `hsm_signing_latency_p99` | Per HSM cluster; alert on regression |

The institution's observability backend (Splunk, Datadog, etc.) ingests these and provides operational dashboards.

## Cost projection at scale

For globally-systemic banks (billions/day):

- HSM cluster (4 members) per region × 4 regions = $200k/year
- Ledger compute (16 instances + multiple regions) = $80k/year
- Storage (5+ TB/year × 7 years cold) = $50k/year
- 1 FTE chain-ops + 0.5 FTE IR specialization = $200k/year

Total: ~$500k/year for the largest institutions. Scales sub-linearly with event volume; the institution's cost-per-event drops as scale increases.

## Performance benchmarks

The reference implementation publishes benchmark data per release:

- Single-host throughput at various configurations
- Memory consumption at various event volumes
- Verifier scan time per million events
- Seal-job duration for typical day sizes

Institutions adopt the spec at scale should validate against the reference benchmarks during pilot.
