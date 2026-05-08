# Cost model

> **What this doc is.** All-in cost picture for adopting the chain-of-custody implementation across institution sizes. Lets the CFO and the IT-cost-center owner make an informed decision.

## Summary

For a typical AI agent deployment:

| Tier | Annual cost (rough order of magnitude) |
|---|---|
| Community bank ($1B–$10B), low-volume AI | $25k–$60k/year |
| Mid-size bank ($10B–$50B), moderate AI use | $80k–$200k/year |
| Tier-1 bank ($50B+), high-volume AI | $300k–$1M+/year |

The dominant cost is HSM operating cost. Software (open source), compute, and storage are secondary.

## Cost components

### HSM (60–80% of total cost)

| HSM | Approximate annual cost |
|---|---|
| AWS CloudHSM Classic | ~$13,000/year per HSM cluster member; clusters typically run 2 members for HA |
| Azure Managed HSM | ~$15,000/year per HSM partition |
| Google Cloud HSM | Pay-per-use; typical $5k–$15k/year per tenant for chain-of-custody volume |
| On-prem HSM | $30k–$80k/year all-in (capital amortization + support contracts + ops) |

Multi-region deployments multiply HSM cost by the region count. The community-bank shape uses one cloud HSM; the Tier-1 shape uses HSMs in 2–4 regions plus a backup cluster.

### Compute (5–15% of total cost)

The ledger server runs on commodity instances. Reference sizing:

| Tier | Ledger compute | Annual cost (cloud) |
|---|---|---|
| Community | 1 instance, 4 vCPU, 16 GB | ~$2,500/year |
| Mid-size | 2–3 instances, 8 vCPU, 32 GB each | ~$8,000–$12,000/year |
| Tier-1 | 4–8 instances per region, 16 vCPU, 64 GB each | ~$40,000+/year |

The SDK runs in the institution's existing application processes; the chain operations are <1 ms per event and don't materially change the host's cost.

### Storage (5–15% of total cost)

| Tier | Hot store | Cold store (7-year retention) |
|---|---|---|
| Community (10 GB/year of events) | ~$300/year for hot Postgres | ~$200/year for S3 + Iceberg |
| Mid-size (200 GB/year) | ~$3,000/year | ~$1,500/year |
| Tier-1 (5 TB/year) | ~$30,000/year | ~$15,000/year |

These are cloud reference numbers; on-prem storage costs are institution-specific.

### Software (0% direct cost; opportunity cost)

The reference implementation is open-source under Apache 2.0. Direct license cost is zero. The institution pays in:

- **Implementation effort.** First-time deployment is typically a 4–12 week engineering project.
- **Operations effort.** Ongoing operations is typically 0.25–1 FTE for a Tier-1; a fraction of an FTE for smaller institutions.
- **Vendor management.** If using a vendor implementation, standard third-party-risk management applies.

### Verifier (negligible)

The verifier runs on examiner laptops; no infrastructure cost. The institution may run the verifier internally for audit purposes; one CI job per audit cycle is sufficient.

### Examination support (variable)

The institution's effort to support examinations:

- Producing the ledger snapshot for the examiner: typically a half-day of operations work
- Supporting examiner questions: typically 5–15 hours per examination

This effort is offset by reduced effort answering ad-hoc audit and regulatory questions about logging integrity, which the chain settles definitively.

## Comparing to alternative

The "do nothing" baseline (continue with vendor-supplied logging without integrity guarantees):

- Direct cost: $0
- Indirect cost: legal exposure when logs are insufficient evidence in a regulatory inquiry; vendor lock-in; inability to satisfy emerging FFIEC AI-governance expectations

The "build your own" alternative:

- Direct cost: 6–18 months of senior engineering effort to design and implement an equivalent system from scratch
- Indirect cost: ongoing maintenance burden; lack of regulator-recognized standard; cross-vendor compatibility burden

The chain-of-custody approach is structurally cheaper than either alternative for any institution that needs the integrity property.

## Cost reduction options

For institutions where cost is the binding constraint:

| Option | Saving | Trade-off |
|---|---|---|
| Relax cadence to weekly | Reduces HSM operations by 7x | Wider retroactive-tamper-detection window; requires examiner approval |
| Shared cloud HSM | Reduces HSM cost by ~50% | Multi-tenancy concerns; documentation in control description |
| Cold-only retention beyond 1 year | Reduces hot-storage cost by 80%+ | Slower query for older events |
| Single region (no multi-region resilience) | Reduces HSM cost by region count factor | Tradeoff against resilience |
| Self-host (on-prem HSM, on-prem compute) | Reduces cloud-HSM costs but adds capex | Capex amortization; ops headcount |

## Cadence vs annual cost

Cadence relaxation (daily → weekly → monthly with examiner approval per spec §4.2.1 and `regulator-pack/examiner-approval-template.md`) is the most-asked cost-reduction option. The intuition that fewer seal jobs means proportional HSM cost reduction is wrong, and institutions that adopt cadence relaxation expecting a 7× HSM-cost reduction from daily → weekly are disappointed. The savings are real but they live in operational complexity rather than direct HSM line-item cost.

The reason is the HSM pricing model. Cloud HSM providers charge for cluster-presence rather than for signing-volume at the volumes the chain produces. AWS CloudHSM Classic charges per HSM cluster member regardless of how many sign operations the cluster performs; Azure Managed HSM charges per partition; Google Cloud HSM pay-per-use is closer to volume-proportional but the chain's signing volume is a small fraction of typical pay-per-use customer profiles. A community bank's daily seal job calls the HSM once or a small number of times per tenant per day. Reducing those signs by 7× (daily → weekly) or 30× (daily → monthly) does not reduce the HSM bill by the same factor because the cluster is sized for availability and the cluster's price floor is not signing-volume-driven.

The savings are in operational complexity:

- Fewer seal-job runs means fewer alerting incidents, fewer post-midnight on-call events, fewer false-positive seal-age alerts during scheduled maintenance windows
- Fewer rotation-boundary edge cases (rotation-and-seal coordination is the most error-prone runbook step in the chain's daily operations)
- Less verifier compute on the institution's own audit cycles
- Lower SOC-evidence-archive volume, marginally
- Lower staff attention budget on routine seal monitoring, releasing the chain-operations team for higher-value work

These savings compound to a meaningful operational-complexity reduction. They do not show up as a 7× line-item drop on the HSM invoice.

### Annual cost by cadence and tier

The numbers below estimate annual cost for daily, weekly, and monthly cadence at each institution size. The cadence-relaxation tiers assume examiner approval is in place per the institution's regulatory relationship.

| Tier | Daily (baseline) | Weekly | Monthly |
|---|---|---|---|
| Community bank ($1B–$10B), low-volume AI | $25k–$60k | $24k–$57k | $23k–$55k |
| Mid-size bank ($10B–$50B), moderate AI use | $80k–$200k | $76k–$190k | $73k–$185k |
| Tier-1 bank ($50B+), high-volume AI | $300k–$1M+ | $295k–$985k | $290k–$975k |

The reduction is single-digit percent at every tier. The HSM cluster-presence cost is the dominant baseline and does not relax with cadence. The compute and storage components do not relax with cadence because the SDK still captures and persists every event regardless of seal cadence; the seal job is a small fraction of total compute load.

### Operational-complexity savings (qualitative)

The complexity savings are real but quantifying them depends on the institution's operations posture. A reasonable estimate framework:

- **On-call attention.** Institution's on-call rotation pages on seal-age threshold and seal-job-failed events. Daily cadence produces roughly 365 seal-job events per tenant per year; weekly produces 52; monthly produces 12. If each event takes the on-call engineer five minutes of attention on average (most pass-through; some require investigation), the annual attention drops from ~30 hours per tenant (daily) to ~4 hours per tenant (weekly) to ~1 hour per tenant (monthly). At a fully-loaded engineering rate of $150/hour, this is $4.5k → $0.6k → $0.15k per tenant per year. For institutions running dozens of tenants, this aggregates.
- **Runbook drift.** Less-frequently-exercised runbooks drift from documented procedure. Daily seal operations are tightly drilled because they happen every day. Weekly seal operations are drilled four times a month. Monthly seal operations are drilled once a month. The institution's chain-operations team maintains procedural readiness through scheduled drills regardless of cadence; the cost is the drill schedule, not the cadence.
- **Verifier runtime.** The verifier processes the same number of events regardless of cadence; the cadence only changes the seal-record count. Verifier runtime cost is roughly cadence-invariant.
- **SOC-evidence archive.** Operational events scale with cadence. Daily produces ~365 seal-related events per tenant per year; monthly produces ~12. The archive cost difference is negligible at the chain's typical event volumes.

### When cadence relaxation actually saves money

Cadence relaxation actually saves money when the institution operates many dormant tenants. The dormant-tenant request in `regulator-pack/examiner-approval-template.md` §3.1 covers this case explicitly. An institution with 50 tenants of which 35 are dormant (paused pilots, decommissioned business lines awaiting compliance retention expiry) experiences meaningful savings under monthly cadence for the dormant population:

- Daily seal cost on 35 dormant tenants: 35 × 365 = ~12,800 empty-day seals per year, each incurring an HSM sign operation and an operational-event capture
- Monthly seal cost on the same 35 dormant tenants: 35 × 12 = 420 empty-month seals per year

The HSM sign-operation cost is small per signature but the operational-event capture and the cluster-attention cost is meaningful at the dormant-population scale. A community bank with a high dormant-tenant ratio sees real savings; a tier-1 bank with primarily-active tenants sees the qualitative-complexity savings only.

### Recommendation for institutions seeking cost relief

Institutions exploring cadence relaxation as a cost-reduction lever SHOULD evaluate the savings as operational-complexity reduction rather than direct HSM-cost reduction. The decision framework:

1. Quantify the institution's tenant population. How many tenants? What fraction are dormant? What fraction are low-volume but active? What fraction are high-volume?
2. Identify the dominant cost driver. For most institutions, HSM cluster-presence dominates. Cadence relaxation does not relax this cost.
3. Identify the on-call attention budget the chain-operations team currently spends on seal monitoring. Cadence relaxation reduces this budget proportionally.
4. Frame the regulator-approval request honestly. The institution's request describes the operational-complexity reduction the institution seeks, not a fictitious HSM-cost reduction. Regulators evaluating cadence-relaxation requests look for the institution's understanding of the cost picture; an institution that overclaims cost savings reveals incomplete operational understanding.

For institutions whose primary cost concern is the HSM line item, the higher-leverage cost-reduction options are HSM-provider negotiation, right-sizing the HSM cluster (don't over-provision for low-volume tenants), and shared-cluster topology (per the Cost reduction options table above). Cadence relaxation is a secondary lever; it produces real but modest savings.

## Cost-benefit calibration

For most institutions, the cost is dominated by the HSM. The HSM is also the integrity anchor — relaxing it weakens the entire claim. Most cost-reduction effort focuses on:

- Choosing the right HSM provider for the institution's existing cloud posture
- Right-sizing the HSM cluster (don't over-provision for low-volume tenants)
- Aligning seal cadence with examination cycle

The institution's CFO works with the CISO to find the right operating point for the institution's risk profile.

## Worked examples

### Community bank, $5B assets, AI for customer-service routing

- One AWS CloudHSM cluster, 2 members: $26,000/year
- Single ledger server, 4 vCPU: $2,500/year
- 10 GB/year storage, 7-year retention: $500/year all-in
- Weekly seal cadence (with examiner approval): no additional cost
- Total: ~$29,000/year

### Mid-size bank, $25B assets, AI for fraud screening + customer service

- One Azure Managed HSM partition: $15,000/year
- Two ledger servers, 8 vCPU each: $8,000/year
- 200 GB/year storage: $4,500/year
- Daily seal cadence: no additional cost
- Internal-audit verifier runs: negligible
- Total: ~$27,500/year

### Tier-1 bank, $250B assets, AI across multiple business lines

- 4 cloud HSM clusters across 4 regions: $52,000/year per cluster x 4 = $208,000/year
- Multi-region ledger compute: $40,000/year
- 5 TB/year storage with multi-region replication: $50,000/year
- Daily seal cadence with hourly option for high-stakes business lines: no additional cost
- 0.5 FTE for ongoing operations: $80,000/year
- Total: ~$378,000/year

For Tier-1 institutions the cost is one or two business-line-supporting people; for community banks it's a single annual line item. Both are tractable for institutions that take logging integrity seriously.

## Cost trajectory over years

Institutions adopting the chain typically see this cost trajectory:

| Year | Stage | Cost driver | Typical cost |
|---|---|---|---|
| Year 1 | Pilot / minimum-viable | Single HSM, single region, daily cadence | $30k–$60k |
| Year 2 | Production-ready | Standard tier deployment; CUEC operation | $50k–$120k |
| Year 3 | Multi-use-case | AI use grows; tenant separation may be added | $80k–$200k |
| Year 5 | Mature operation | Multi-region or multi-line-of-business may be added | $150k–$400k |
| Year 7+ | Optimization | Cost optimization based on operational experience | Plateau |

Trajectory drivers:

- AI use case expansion (more events, more tenants)
- Resilience needs growing (multi-region)
- Risk profile change (cadence tightening or relaxation)
- HSM provider negotiation (volume discounts)
- Operational efficiency improvements (automation reduces ops headcount)

The trajectory is institution-specific. The institution's CFO forecasts based on AI program growth.

## Cost reduction over time

As the chain operation matures, costs typically decrease per-event:

- HSM costs are roughly fixed; per-event cost drops as event volume increases
- Storage costs scale with retention; cold-store optimization reduces ongoing costs
- Operational headcount stabilizes; per-event cost drops as automation matures
- SOC engagement costs may stabilize as the institution's controls mature

A mid-size bank operating the chain for 5+ years typically sees total cost stabilize at $100k–$200k/year regardless of moderate growth in AI use volume.

## Notes

- Numbers are 2026 cloud reference pricing; actual costs vary by negotiated pricing and existing volume commitments
- On-prem deployments shift costs to capex and ops headcount with different financial profiles
- Multi-tenant deployments (vendor-hosted) can substantially reduce per-tenant costs but increase vendor-management overhead
- The cost picture is most predictable under cloud HSM; on-prem options have higher variance
- SOC engagement cost is typically $50k–$200k for a chain-of-custody implementation; this is the SOC firm's pricing, not chain-specific
