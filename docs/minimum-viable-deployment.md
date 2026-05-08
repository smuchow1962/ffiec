# Minimum-viable deployment

> **What this doc is.** Smallest deployment shape an institution can run conformantly. For institutions with limited operational bandwidth, low AI volume, or pilot-stage adoption.

## When to use this

Three scenarios:

- **Pilot.** Bank is piloting AI for a low-stakes use case (customer-service routing, document classification). Production-grade chain is overkill for the pilot phase.
- **Community bank ($1B–$10B).** AI use is meaningful but volume is low. Cost sensitivity is high.
- **First-time adopter.** Bank has not previously operated a chain implementation; phased ramp-up reduces operational risk.

For Tier-1 banks or institutions with high-volume AI, see `docs/cost-model.md` mid-size or large-tier guidance and `docs/at-scale-operations.md`.

## What's in the minimum-viable deployment

### Core components (all required)

1. **The SDK** in the AI agent process
2. **The ledger server**, single-instance
3. **A FIPS 140-2 L3 HSM**, single cluster member acceptable for pilot/community-bank
4. **The verifier**, run on examiner laptops
5. **The institution's secret manager** (typically already operating)

### Configuration choices for minimum-viable

| Decision | Minimum-viable choice |
|---|---|
| Topology | Self-hosted in single cloud region |
| Cadence | Daily (default) |
| HSM provider | AWS CloudHSM, Azure Managed HSM, or Google Cloud HSM (pick the institution's primary cloud) |
| HSM cluster size | 1 member for pilot; 2 members for production-but-low-volume |
| Postgres | Single instance with daily backup |
| Cold store | S3 with default-tier configuration |
| Multi-region resilience | Defer; single-region is acceptable until AI use grows |
| Reconciliation cadence | Weekly (spec ceiling) |
| Verifier validation | Cosign + SHA-256 hash check; defer reproducible-build verification until year 2 |

### Minimum CUECs to operate from day-one

The full CUEC list is 21 items. For minimum-viable deployment, operate the critical subset (CUEC tiering):

**Critical (operate from day 1):**

- CUEC-IAM-01..05 (database role restrictions, HSM separation of duties — at small bank, dual control acceptable)
- CUEC-CRY-01..05 (HSM PIN management, master key custody, key rotation)
- CUEC-OPS-01 (NTP sync — typically already operating)
- CUEC-OPS-04 (seal-age monitoring)
- CUEC-OPS-05 (regulator notification on 72-hour seal delay)
- CUEC-VER-01 (run the verifier before each examination; quarterly internal review)
- CUEC-IR-01..04 (incident response playbook adaptation)

**Important (operate within first 6 months):**

- CUEC-VER-02..04 (verifier validation before use; reproducible-build verification annually)
- CUEC-IR-04 (cyber-incident notification)
- CUEC-VND-01..04 (if using a vendor implementation)
- CUEC-OPS-02..03 (operational event retention; backup procedures)

**Operational hygiene (operate within first year):**

- CUEC-OPS-* (the rest)
- CUEC-CFG-01..03 (configuration and change management)

This tiering lets institutions ramp up over 6–12 months. The control description declares which tier is operational; SOC and examination teams test against the declared scope.

## Cost target

For a community bank pilot:

| Component | Annual cost |
|---|---|
| Single AWS CloudHSM member (or equivalent) | $13,000 |
| Single ledger Postgres instance | $2,500 |
| Storage (10 GB events/year, 7-year retention) | $500 |
| Compute (small ledger instance) | $1,500 |
| Operational overhead (0.1 FTE) | $20,000 |
| **Total minimum-viable** | **~$37,000/year** |

For an institution running multiple AI use cases, costs scale; the cost model articulates.

## Operational ramp-up timeline

### Month 0–1: Pilot deployment

- Provision HSM and master-key custodian
- Deploy ledger in single region
- Integrate SDK in pilot AI agent
- Verify end-to-end flow with sample events

### Month 1–3: Initial production

- AI agent goes live
- Daily seal job runs
- Monitor seal-age metric; verify 60-min default works
- Internal audit runs the verifier monthly

### Month 3–6: Operational hardening

- Implement remaining critical CUECs
- Document control descriptions
- Establish reconciliation procedure
- Run IR playbook table-top exercise

### Month 6–12: Audit readiness

- First SOC engagement (if applicable) using deliverable layer
- First chain examination by primary regulator
- Full CUEC operation
- Cost optimization review

### Year 2+: Maturation

- Multi-region resilience (if AI use grows)
- Higher-frequency cadence (if risk profile changes)
- Cross-tenant patterns (if AI use diversifies)
- Quantum-readiness drills

## Skip what

For minimum-viable, the institution may skip until later:

- Multi-region resilience (single-region is acceptable initially)
- Hourly seal cadence (daily is the spec default)
- Full reproducible-build verification (cosign + SHA-256 is sufficient for low-risk early adoption)
- Sophisticated reconciliation (weekly per spec is sufficient; baseline establishment can be retrospective)
- Per-business-line tenant separation (single tenant_id is fine for institutions with one use case)

These are not optional for production-grade deployment at scale; they are deferable for minimum-viable.

## When to upgrade

Triggers for moving beyond minimum-viable:

- AI volume exceeds 1M events/day
- AI use case becomes customer-decisioning (credit, insurance, etc.) — higher integrity requirements
- The bank's risk profile materially increases
- A regulator finding mandates an upgrade
- Multi-region resilience becomes operationally important

The cost model and the per-tier-proportionality guidance (`00-overview.md` §5.5) drive the upgrade path.

## Common minimum-viable adopter mistakes

- **Treating the chain as just another OTel stream.** It's not; the integrity properties depend on operating the CUECs.
- **Skipping the IR playbook.** A chain-detected event requires response; without a playbook, the institution improvises.
- **Underestimating HSM operational cost.** The $13k/year HSM is the institution's standing commitment; budget accordingly.
- **Skipping verifier validation.** Running an unverified verifier risks false confidence; cosign validation is mandatory before use.
- **Deferring documentation.** Control descriptions, IR playbook, customer-correlation index are all institution-specific; deferred documentation creates risk during the first examination.

## Minimum-viable doesn't mean compromised

The chain's integrity properties are unchanged at minimum-viable scale. The chain catches tampering at 1,000 events/day exactly as it does at 1B events/day. The cryptographic substrate is identical. What scales is the operational posture, not the integrity claim.

A community bank running minimum-viable chain has the same defensible audit posture as a Tier-1 bank running at scale.

## Minimum-viable for vendor-hosted topology

Institutions choosing vendor-hosted topology for minimum-viable have an even simpler ramp:

- The vendor operates the ledger and HSM
- The institution operates the SDK in its agent processes
- The institution operates the IR playbook for chain-detected events
- The institution holds the public key
- The institution receives the verifier output and reports

Vendor-hosted minimum-viable is a fit for institutions with very limited IT-ops bandwidth. The vendor's SOC report covers the chain operations the institution would otherwise operate; the institution inherits the controls.

For vendor-hosted, see also `docs/byoc-deployment.md` for the boundary between the bank and the vendor.
