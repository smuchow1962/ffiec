# Disaster recovery and resilience

> **What this doc is.** Reference DR/resilience guide for the chain implementation. The institution's standard DR program covers the chain; this document supplements with chain-specific considerations.

## Reference RPO/RTO targets

| Component | RPO target | RTO target | Mechanism |
|---|---|---|---|
| WAL | < 1 second | < 15 minutes | PostgreSQL streaming replication or equivalent |
| Hot store | Derivable from WAL | < 1 hour | Re-derive from WAL on standby |
| Cold store | < 1 hour | < 4 hours | S3 cross-region replication or equivalent |
| Daily seals | Same as WAL | Same as WAL | Stored in WAL-backed table |
| HSM | < seconds (active-active cluster) | < 5 minutes | HSM cluster failover |

These are reference targets; institutions declare their own targets in their control description.

## Synchronous vs asynchronous replication

The trade-off:

| Posture | RPO | Hot-path latency cost | When appropriate |
|---|---|---|---|
| Async same-region | Seconds | None | Default; most banks |
| Sync same-region | Sub-second | +1–3 ms | High-volume institutions willing to pay the latency |
| Async cross-region | Seconds for cross-region; none for hot path | None to hot path | Cross-region resilience without hot-path cost |
| Sync cross-region | Sub-second | +10s of ms | Generally not recommended; defeats the hot-path budget |

The hot-path budget is <2 ms p99 (`docs/design/02-chain-construction.md` §5). Same-region synchronous replication fits; cross-region synchronous does not. Most institutions choose async same-region with optional async cross-region for resilience.

## Failure scenarios and recovery

### Scenario: Single ledger instance failure

- **Detection.** Health-check failure on `/readyz`
- **Response.** Load balancer routes new traffic to healthy instances; failed instance restarts or is replaced
- **Recovery time.** Seconds to minutes
- **Data loss.** None (events in flight either succeed on retry or are buffered in SDK local SQLite)

### Scenario: Ledger database failure (primary)

- **Detection.** Connection failures from ledger workload to primary database
- **Response.** Automatic or manual failover to read replica; promote replica to primary
- **Recovery time.** Minutes (depends on failover automation)
- **Data loss.** Up to RPO target (sub-second with sync replication; seconds with async)

### Scenario: HSM unavailability

- **Detection.** HSM operation failures; `seal.job_failed` events fire
- **Response.** Cluster failover to healthy HSM members; if cluster-wide, escalate to vendor support
- **Recovery time.** Sub-minutes for cluster failover; hours to days for vendor-side incidents
- **Effect on chain.** Events continue to be captured and chained (independent of HSM); seals are delayed until HSM is restored
- **Notification.** 72-hour threshold per `docs/design/04-hsm-custody.md` §5.2

### Scenario: Region-wide outage

- **Detection.** Multiple component failures across the region
- **Response.** Pattern A (active-active with seal-region pinning per `docs/design/00-overview.md` §6.4): if the seal region is the affected region, promote a replication region to seal-region status per CC8.1; the promoted region must have all events for the tenant-day before producing the seal. Pattern B (per-region tenant_id): the regional tenant's chain is unaffected at other regions; events captured under the regional tenant during the outage may be lost up to RPO. Either pattern: SDK in unaffected regions continues to capture; SDK local SQLite buffers events for replay when ledger reachability is restored.
- **Recovery time.** Hours (depends on the institution's regional DR plan)
- **Data loss.** Up to RPO target for cross-region replication
- **Effect on chain.** Events captured during the outage may be lost if SDK local SQLite is also lost; institution's RPO sets the bound. For Pattern A, replication-loss is detected by the institution's per-region event-count reconciliation (`master.cross_region_replication_completed` operational event); the seal accurately seals what's in the seal region's ledger and the lost events are a regional ingest issue, not a chain-integrity issue.

### Scenario: WAL corruption (storage-level)

- **Detection.** Verifier reports merkle root mismatch; database integrity checks fire
- **Response.** Restore WAL from backup; replay any gap from SDK buffers or downstream OTLP
- **Recovery time.** Hours
- **Data loss.** Depends on backup recency; gap-fill from SDK buffers minimizes
- **Effect on chain.** Affected days unverifiable until gap is filled; if unfillable, treated as integrity-control failure

### Scenario: Backup tampering

- **Detection.** Verifier reports merkle root mismatch
- **Response.** Investigate; restore from a known-good backup; treat as denial-of-service-against-integrity
- **Recovery time.** Hours to days depending on investigation
- **Effect on chain.** Affected days unverifiable; institution treats as control failure

### Scenario: HSM signing-key compromise

- **Detection.** External intelligence; signature anomalies; reconciliation finding
- **Response.** Activate IR playbook scenario 3; rotate signing key
- **Recovery time.** Hours for rotation; days to weeks for full investigation
- **Effect on chain.** Past seals signed under compromised key require institution's compensating evidence

## Multi-region resilience (normative)

Two conformant patterns at v1.0 per spec §10.15 and `docs/design/00-overview.md` §6.4. The institution selects per tenant and documents the choice in CC8.1.

### Pattern A — Active-active with seal-region pinning (RECOMMENDED)

A single canonical `tenant_id` operates in multiple regions. Each region has its own ledger ingesting events from local SDKs. One region is the **seal region**; the others are **replication regions**. Replication regions ship their events to the seal region by seal-time. The seal region computes one Merkle root per tenant-day over events from all regions and produces one HSM-signed seal.

```
tenant_acme_prod  (single tenant_id, three regions)

  Region us-east-1   Region eu-west-1   Region ap-south-1
  (replication)      (replication)      (seal region)

       SDKs               SDKs               SDKs
        |                  |                  |
        v                  v                  v
  Local ledger       Local ledger       Aggregate ledger
        |                  |                  |
        +----- replicate ---->+----- replicate ---->+
                                                    |
                                            Daily Merkle seal
                                            HSM signature
```

Operational discipline (CC8.1):

- Designate the seal region per tenant.
- Operate cross-region replication of events to the seal region (institution's choice — Postgres streaming replication, Kafka cross-region, S3 cross-region replication).
- Replication-completion SLA aligns with the seal-job's start time per spec §4.3 (60-minute publish window after seal-window end).
- Operate per-region event-count reconciliation: each region reports event count per tenant-day; the seal region's count must equal the sum of regional counts. Mismatch is a control failure, not a chain-integrity failure (the seal accurately seals what's in the seal region's ledger).
- Document seal-region failover: if the seal region becomes unavailable before seal-time, promote a replication region. The promoted region must have all events for the tenant-day before producing the seal; the institution's tenant key registry resolves the per-day signing entity.

Run-locality (normative for v1.0): runs start and end in one region. Cross-region run continuation is deferred to v1.1. Workloads requiring cross-region run continuation route to a single region or operate under Pattern B.

### Pattern B — Per-region `tenant_id` (CONFORMANT alternative)

```
tenant_acme_prod_us_east_1   (its own IKM, seals, verifier runs)
tenant_acme_prod_us_west_2   (its own IKM, seals, verifier runs)
tenant_acme_prod_eu_west_1   (its own IKM, seals, verifier runs)
```

Each regional tenant is an independent chain. Cross-region correlation is institution-side — the institution maintains a registry mapping the regional tenants to one logical "Acme prod" deployment, consulted for human disambiguation during examiner inquiries or customer disputes.

Pattern B is heavier than Pattern A: the verifier runs O(regions) times per audit period, and cross-region comparison is institution-correlated rather than seal-aggregated. It is appropriate when:

- The institution's regional regulatory regime mandates in-region key custody (some EU banking jurisdictions; APAC data-sovereignty regimes).
- The institution's risk posture treats cross-region replication as an unacceptable trust boundary (the seal region's compromise would corrupt the seal even though events were captured securely elsewhere).
- The institution operates regional ledgers under different vendor contracts and per-region IKM custody is part of the contractual posture.

### Selection and switching

The patterns are mutually exclusive per tenant — an institution operating Pattern A for one tenant MAY operate Pattern B for another. A given tenant operates under one pattern only. Switching patterns is a chain-discontinuity event analogous to a posture change per spec §4.1.2 and is governed by the institution's documented change-management procedure.

## DR exercise procedure

The institution exercises the DR plan at least annually:

### Table-top exercise (quarterly)

- Walk through each scenario
- Confirm runbooks are current
- Confirm contact lists are current
- Update the IR playbook with any gaps identified

### Actual failover exercise (annually)

- Failover from primary to secondary region (during a maintenance window)
- Run the verifier against the secondary region's ledger
- Confirm new events are captured under the secondary region's tenant_id
- Failback to primary region
- Document the exercise in the institution's control-evidence repository

## Backup integrity testing

Backups are not backups until they're tested for restore. The institution tests:

- Streaming-replication failover to a replica
- WAL restore from S3 to a fresh Postgres instance
- Cold-store restore for old data
- Verifier validation against restored data

A backup that fails restore is a finding even if the primary is healthy.

## Capacity planning

Hot path: <2 ms p99 per event. Throughput scales with ledger instance count and Postgres write capacity. Reference sizing per Tier in `docs/cost-model.md`.

Seal job: O(N) over events for the day. Streaming Merkle is O(log N) memory. For 1B events/day, the seal job runs in tens of minutes on a single ledger instance.

Cold-store growth: ~event-volume × retention period. For 200 GB/year × 7 years = 1.4 TB at end-of-retention. Storage cost is small relative to other components.

## What this doesn't cover

The institution's broader DR program covers:

- Application-level failover for the AI agent platform itself
- Network failover and routing
- DNS failover
- Customer-facing communications during outages
- Vendor coordination during outages

The chain composes with the broader program. The chain's specific concerns (HSM, daily seals, master-key custody) are documented here; everything else is the institution's standard practice.
