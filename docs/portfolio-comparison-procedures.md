# Portfolio comparison procedures

> **What this doc is.** Procedures for the EIC or supervisory team comparing chain-of-custody output across multiple banks (portfolio-level), or across multiple regions of the same bank (multi-region cross-comparison). Closes the FFIEC IT Examiner and EIC partials around portfolio tooling.

## Why this matters

Regulators supervise portfolios of institutions. Cross-portfolio comparison surfaces patterns: institutions with consistently strong chain operation, institutions with concerning anomalies, institutions where examination resources warrant focus. Until the v1.1 `verifier portfolio` subcommand ships, comparison is manual; this doc articulates the manual procedure.

## Cross-bank portfolio comparison

### Inputs the EIC assembles

For a portfolio of N banks:

- N verifier reports (PDF + JSON) for the most recent examination period
- N verifier bundles (the working-paper artifacts)
- N institution control descriptions identifying claimed cadence and configuration

### Comparison metrics

The EIC computes (manually or via a portfolio-comparison spreadsheet):

| Metric | Reading |
|---|---|
| **Pass rate** | Percentage of days passed across the examination period; healthy banks ~100% |
| **Anomaly rate** | Anomalies per day on average; healthy banks low; outliers warrant follow-up |
| **Sealing-delay frequency** | Days where sealing exceeded the institution's threshold; clusters indicate operational issues |
| **Late-binding rate** | Late-binding events as fraction of total events; persistent high rate indicates operational concern |
| **Master-rotation events** | Per institution; rotations during the period |
| **Verifier-validation cadence** | How often the institution validated the verifier binary |

### Comparison patterns

The EIC identifies patterns:

- **Consistent strong performance.** Pass-rate ~100%, low anomaly rate, no notable issues. These institutions warrant typical examination depth.
- **Operational concerns.** Pass-rate ~100% with sealing-delay clusters or elevated anomalies. These institutions warrant operational follow-up at next examination.
- **Integrity concerns.** Pass-rate <100% with chain-detected events. These institutions warrant deeper investigation.
- **Outliers.** Pattern materially different from peers. These institutions warrant investigation regardless of pass/fail status.

### Cross-bank language

When the EIC observes cross-portfolio patterns:

- Internal supervisory: the EIC documents the pattern in the regulator's portfolio-management system
- Cross-agency: if the pattern affects multiple agencies' examinations, the EIC coordinates through standard supervisory MOUs
- Public disclosure: if the pattern warrants public disclosure (industry guidance, supervisory publications), the EIC engages the regulator's communications function

### Sample comparison workspace

For a portfolio of 6 banks, a comparison workspace might look like:

| Bank | Pass days | Anomaly days | Sealing >24h | Late-binding rate | Notes |
|---|---|---|---|---|---|
| Bank A | 30/30 | 2 | 0 | 0.001% | Standard healthy |
| Bank B | 30/30 | 0 | 0 | 0.0001% | Best-in-portfolio |
| Bank C | 28/30 | 5 | 3 | 0.005% | HSM cluster issue Q1; remediated |
| Bank D | 30/30 | 1 | 0 | 0.001% | Standard healthy |
| Bank E | 30/30 | 8 | 1 | 0.012% | Anomaly cluster; investigation pending |
| Bank F | 30/30 | 1 | 0 | 0.0005% | Standard healthy |

The EIC reads the table to identify which institutions warrant more depth.

## Multi-region cross-comparison (within one institution)

For institutions running per-region tenant_ids (multi-region workaround per `00-overview.md` §5.4):

### Inputs

- Verifier reports per region
- The institution's documented multi-region operation pattern

### Comparison

The EIC compares per-region:

- Did each region operate continuously?
- Were there regional outages?
- Are anomaly patterns consistent across regions?
- If active-passive deployment: did the secondary region see traffic?

### Cross-region patterns

- **Consistent.** All regions show similar pass rates and anomaly patterns; healthy multi-region operation.
- **Region-specific issue.** One region shows distinct pattern; investigate region-specific operational conditions.
- **Failover evidence.** Active-passive deployments show traffic transition during failover; the EIC verifies the failover was orderly.

## Cross-period comparison (within one institution)

For longitudinal analysis of a single institution:

### Inputs

- Verifier reports across multiple examination periods (current + prior 1–2 examinations)
- The institution's control description history

### Comparison metrics

- Trend in pass rate
- Trend in anomaly rate
- Trend in sealing-delay frequency
- Trend in operational-events volume
- Material configuration changes (cadence, vendor, deployment topology)

### Patterns

- **Stable operation.** Pattern unchanged across periods; the institution's controls operate consistently.
- **Improvement trajectory.** Anomaly rate declining; sealing delays declining; control maturation evident.
- **Degradation trajectory.** Anomaly rate increasing; warrants follow-up examination.
- **Recovery from incident.** Period containing an incident shows distinct pattern; subsequent periods show recovery.

## Tools to support comparison

### Manual aggregation (current state)

EIC aggregates per-bank verifier outputs into a comparison spreadsheet. Standard practice for cross-portfolio examination.

### v1.1 `verifier portfolio` (planned)

The v1.1 subcommand will:

- Take multiple verifier bundles as input
- Produce a unified comparison report with the metrics above
- Output deterministic comparison output suitable for working-paper retention

Until v1.1 ships, the manual procedure produces equivalent comparison data.

### Regulator-internal tooling

Some regulators may build internal portfolio-management tooling that consumes verifier JSON output. The institution-side artifacts are unchanged; the regulator's internal aggregation is the regulator's tooling.

## Working-paper preservation

For each portfolio comparison:

- The comparison workspace is retained in the EIC's working papers
- The per-bank verifier outputs are retained per institution's standard examination practice
- The comparison conclusions feed the EIC's portfolio-management system

## Cross-agency coordination

For institutions supervised by multiple agencies (e.g., bank holding companies under Fed + bank subsidiaries under OCC):

- Each agency conducts its own examination
- Verifier output is the same; agencies may share output through standard supervisory channels
- Portfolio comparison happens within each agency's framework

The chain provides consistent input across agencies; the portfolio comparison is each agency's analytical responsibility.

## Common patterns and their implications

### "All institutions in a sub-portfolio show similar anomaly pattern"

Often indicates:
- A shared vendor experiencing operational issues
- A shared HSM provider with degradation
- A shared cloud region with infrastructure issues

Action: the supervisory function investigates the shared root cause; the institutions individually receive findings if their own monitoring should have caught the pattern earlier.

### "One institution differs materially from the portfolio average"

Often indicates:
- Misalignment between the institution's claimed cadence and its observed cadence
- Operational immaturity at the institution
- A unique configuration (cross-region, sub-tenant) producing unexpected patterns
- An issue the institution hasn't yet identified

Action: deeper examination at the next cycle; supervisory engagement on the pattern.

### "Portfolio-wide degradation over time"

Often indicates:
- An emerging operational issue affecting many institutions
- A regulatory framework change requiring institutional adaptation
- A vendor-side degradation across the portfolio

Action: the supervisory function escalates to industry-level analysis; cross-agency coordination if appropriate.

## Limitations of manual comparison

Manual comparison works for portfolios up to ~20–50 institutions. Beyond that, the comparison workspace becomes unwieldy. v1.1 `verifier portfolio` tooling will scale.

Until v1.1, regulators with very large portfolios may build internal aggregation tools that consume the verifier JSON output. The chain produces deterministic JSON; aggregation is straightforward.

## Cadence for portfolio comparison

The EIC's portfolio comparison runs:

- After each examination cycle (typical)
- Ad-hoc when patterns emerge (cross-bank issue investigation)
- Annually for portfolio-wide assessment

The cadence is the regulator's choice; the chain provides the substrate at any cadence.

## Related documents

- [`regulator-pack/sample-report.md`](regulator-pack/sample-report.md) — what individual bank reports look like
- [`regulator-pack/finding-language.md`](regulator-pack/finding-language.md) — examination-report language
- [`design/07-verifier-design.md`](design/07-verifier-design.md) — verifier output format
- [`first-engagement-guide.md`](first-engagement-guide.md) — first-engagement preparation
