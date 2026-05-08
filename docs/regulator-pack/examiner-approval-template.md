# Examiner-approval template — seal cadence relaxation

> **What this doc is.** Template for an institution's request to relax seal cadence (daily → weekly, or weekly → monthly). Submitted to the institution's primary regulator; responded to by the regulator's IT examination function.

## Request format

Institutions submit the request as a structured document with these sections:

### Section 1 — Institution identification

- Institution legal name and FDIC certificate number (or equivalent)
- Primary regulator
- Tenant identifier(s) covered by this request
- Effective date of the proposed change
- Duration of the relaxation (typically until the next examination cycle)

### Section 2 — Current and proposed cadence

| | Current | Proposed |
|---|---|---|
| Cadence | Daily | Weekly (or Monthly) |
| Seal job runtime window | UTC 00:00 + 60 min | (Specify) |
| Seal-day-boundary semantics | UTC | UTC (unchanged) |

### Section 3 — Rationale

Plain-language explanation of why the relaxation is appropriate. Typical rationales:

- AI agent traffic volume is below [X] events/day; daily HSM operations are disproportionate
- Examination cycle is 18-month; daily granularity is finer than the regulator-defined need
- HSM operating cost at daily cadence is meaningful relative to institution size
- Risk profile has changed (e.g., AI agent retired or moved to lower-risk use case)

The rationale is the institution's case. The regulator weighs it against the integrity-window trade-off.

#### Section 3.1 — Dormant tenant rationale

Some tenants are not low-volume; they are dormant. A paused pilot program awaiting executive go/no-go, a decommissioned business line waiting out its compliance retention window, an acquired-bank product line in run-off — the tenant produces zero events for weeks or months but the daily seal job still runs and the HSM still signs an empty-day Merkle root each midnight. For institutions operating dozens of dormant pilot or run-off tenants, the empty-day signing cost compounds on infrastructure that is not producing examinable AI activity.

The dormant-tenant request is a specific variant of the cadence-relaxation request. Rather than relaxing cadence for an active low-volume tenant, the institution requests cadence relaxation for a tenant that is structurally inactive and expected to remain so. The institution's request MUST name the dormant-state criteria and the resumption trigger so the regulator can see the relaxation is bounded.

The request specifies:

- **Documented dormancy threshold.** The tenant has captured zero events for a documented continuous period — typically 90 days, but the institution may propose a different threshold with rationale. The institution names the metric (e.g., `events_captured_total{tenant_id} == 0` over the rolling 90-day window) and the source (the ledger's operational-event stream or the institution's observability backend).
- **Proposed dormant cadence.** Weekly or monthly. The empty seals still run; the HSM still signs each empty Merkle root; the cadence just stretches.
- **Automatic resumption to daily cadence on the first event.** This is the load-bearing safety property. The moment the tenant captures its first event after dormancy, the cadence resets to daily for that tenant. The institution's control description names the operational mechanism (the seal scheduler observes a non-empty event window and re-arms the daily timer; the institution emits a `seal.cadence_resumed` operational event for the regulator's audit trail). The institution MUST NOT require human approval for resumption — resumption is automatic, and the human-in-the-loop is the institution's notification to the regulator that resumption occurred.
- **Control-description language.** The institution's control description names the dormant-state criteria, the operational mechanism for detection and resumption, the regulatory approval reference, and the institution's monitoring of the dormant population (typically a quarterly review confirming dormant tenants are still dormant).
- **Examiner's right to revoke.** The regulator may revoke the dormant-tenant relaxation at any time without cause. The institution commits in the request to revert the affected tenant(s) to daily cadence within one business day of receiving a revocation notice.

The dormant-tenant request differs from the generic low-volume request in two ways. First, the threshold is binary (zero events) rather than rate-based (below X events/day), which is easier for the regulator to verify mechanically. Second, the request commits the institution to automatic resumption rather than requiring re-approval at first activity, which removes the operational risk of a dormant tenant becoming active and silently continuing under the relaxed cadence.

##### Example approval-letter excerpt

The institution may adapt the following text for its request and the regulator may adapt parallel language for its approval. Both excerpts are illustrative; the institution's request and the regulator's response remain authoritative documents.

> **Subject.** Dormant-tenant cadence relaxation, tenants `pilot_treasury_2024_q3`, `runoff_legacy_collections`, and `pilot_servicing_2025_q1`.
>
> **Approved cadence.** Monthly seal cadence for the named tenants while each tenant remains in dormant state.
>
> **Dormant-state definition.** A tenant is in dormant state when it has captured zero chain events for the preceding 90 continuous days as measured by the institution's `events_captured_total{tenant_id}` metric.
>
> **Automatic resumption.** Upon capture of any chain event for a dormant tenant, the institution's seal scheduler automatically resumes daily cadence for that tenant beginning with the next UTC midnight boundary. The institution emits a `seal.cadence_resumed` operational event citing this approval reference and notifies the examination contact within five business days.
>
> **Revocation.** The regulator may revoke this approval in whole or for any named tenant at any time. The institution commits to reverting the affected tenants to daily cadence within one business day of receiving notice.
>
> **Duration.** This approval is effective for 12 months from the effective date and may be renewed using the same template.
>
> **Reference.** Approval document ID `[regulator-issued-id]`, dated `[regulator-issued-date]`.

### Section 4 — Compensating controls

Cadence relaxation widens the retroactive-tamper-detection window (1 day → 7 days for weekly; 1 day → 30 days for monthly). The institution proposes compensating controls to bound the increased risk:

- More aggressive monitoring of integrity-alert events at the ledger
- Faster IR response for chain-detected events (e.g., 15-minute response instead of 1-hour)
- Per-event cryptographic timestamping (RFC 3161) as additional control layer
- Independent third-party verification at higher frequency than the regulator-mandated cycle

The institution names the controls; the regulator accepts or rejects.

### Section 5 — Operational impact

- Number of seals per year (current vs proposed)
- HSM operating cost (current vs proposed; estimated annual)
- Verifier runtime per examination (current vs proposed)
- Other operational considerations

### Section 6 — Risk assessment summary

The institution's risk function has reviewed the request and concurs that the proposed cadence is appropriate for the risk profile. Signed by the institution's:

- Chief Risk Officer (or delegate)
- Chief Information Security Officer (or delegate)
- Senior officer responsible for the affected AI agent program

### Section 7 — Reversal commitment

The institution commits to reversing the relaxation immediately if any of the following occurs:

- Material change in the AI agent program's risk profile
- Cyber-incident affecting the chain or related infrastructure
- Regulator finding that the relaxation is no longer appropriate
- Any chain-verification failure during the relaxation period

## Regulator response format

The regulator responds within 30 days of receipt with one of:

**Approval.** The request is approved as submitted. Response document records:

- Approval date and effective date
- Duration of the approval
- Conditions of approval (if any)
- Renewal procedure
- Response document ID for the institution's records

**Approval with modifications.** The regulator approves a modified version (e.g., weekly approved for 12 months instead of indefinitely). Same fields as above plus the modifications.

**Denial.** The regulator denies the request. Response document records:

- Denial date
- Specific reasons for denial
- What changes to the request would be needed for approval
- Appeal procedure if applicable

**Request for information.** The regulator requests additional information before deciding. Specifies what is needed and a response deadline.

## Where the approval document lives

The institution's control description includes:

- A reference to the approval document (regulator + document ID + date)
- The current effective cadence
- The conditions and reversal commitments

The verifier records the institution's claimed cadence in the seal record, where the examiner cross-checks at the next examination.

## Re-approval

Approvals are time-bounded. Before expiration, the institution submits a re-approval request using the same template, with an updated rationale that reflects the operating experience under the relaxation.

## Audit trail

The approval document, the institution's request, the regulator's response, and any subsequent modifications are retained for the institution's standard audit-document retention period (typically 7 years).
