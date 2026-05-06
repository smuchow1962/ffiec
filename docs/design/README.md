# Design documents

> **Purpose.** Internal design notes that justify the choices in the normative specification. The auditor (Big Four / FFIEC examiner) is the advisor whose lens the design satisfies. Where a design choice could go either way, the doc explains *why we chose* the way we did and what an auditor would say to confirm or push back.

## Reading order

1. [`00-overview.md`](00-overview.md) — system overview and data flow
2. [`01-primitives-spec.md`](01-primitives-spec.md) — the four primitives, design rationale
3. [`02-chain-construction.md`](02-chain-construction.md) — HMAC chain detail (the hot path)
4. [`03-merkle-seal.md`](03-merkle-seal.md) — daily Merkle seal detail
5. [`04-hsm-custody.md`](04-hsm-custody.md) — HSM signing and key custody
6. [`05-otlp-wire.md`](05-otlp-wire.md) — OTel wire format, attribute conventions
7. [`06-ledger-server-design.md`](06-ledger-server-design.md) — `ledger/` architecture
8. [`07-verifier-design.md`](07-verifier-design.md) — `verifier/` architecture
9. [`08-test-vectors.md`](08-test-vectors.md) — conformance corpus design
10. [`09-threat-model.md`](09-threat-model.md) — adversaries and properties

## The auditor's-lens convention

Every design doc closes with an **&ldquo;Auditor's-lens review&rdquo;** section. That section names the questions an auditor would ask, the answer the design provides, and any open issues where the auditor's concern is acknowledged but not yet fully resolved.

Open issues escalate to specification review per the process in [GOVERNANCE.md](../../GOVERNANCE.md).
