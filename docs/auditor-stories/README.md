# Auditor stories

A growing set of day-in-the-life audit narratives. Same team, same playbook, different industries, different chain-of-custody postures, different countries. Read in any order. Read all eleven if you want the contrast.

## The team

The eight people who travel together when the engagement requires it:

- **Dawn** — Lead Auditor (governance + narrative)
- **Raj** — Database specialist
- **Elena** — CRM systems
- **Mike** — Application / API layer
- **Diana** — IAM & access control
- **Luis** — DevOps / logs / pipelines
- **Chen** — Data engineering / ETL
- **Tom** — Internal-audit liaison specialist (partners with the client's CAE)

They have done this work for years. They know the patterns. They know what an honest answer sounds like and what a deflection sounds like. They are not mean. They are not credulous. They have earned travel-itinerary autonomy from their firm and split themselves across cities, time zones, and language barriers when an engagement calls for it.

## The baseline

Before these engagements, the team spent a week at a mid-size financial services company that had no chain-of-custody. The customer-interaction data was overwritten in the CRM, the database admin log retained 7 days, IAM was on tickets nobody closed, the pipeline had no checksums, and engineers could delete CloudWatch log groups. The reconciliation test produced two mismatches, one missing record, and three timestamp inversions out of ten. The final theme: *"The organization cannot demonstrate that customer interaction data is complete, accurate, or unaltered."*

That week is the floor. Every story in this folder is calibrated against it.

## TesseraSeal

When the stories say **TesseraSeal**, they mean the product brand for a deployed Herald-ecosystem chain — Herald.Py instrumenting the capture surface, Herald Core operating the ledger and seal job, Herald.Compliance presenting the regulator-facing surface where examiners pull verifier output and seal records. The wire form is FFIEC chain-of-custody v1.0a. Seals are Ed25519 over an HSM-held key. Verification is the `herald-verify` CLI documented in the spec §10.12 exit-code contract: 0 PASS, 1 procedure-could-not-begin, 2 procedure-began-and-failed, 3 chain-anomaly.

Eleven distinct deployment contexts show up across these stories.

## The stories — domestic engagements

### [01 — Northbridge Federal Savings (banking)](01-northbridge-federal-savings.md)

A regional US bank, ~$45B consolidated assets, OCC-supervised. **TesseraSeal in use across all data capture for 18 months.** Every customer-facing surface is instrumented. The team walks in expecting the baseline. They finish 30% under budget time with one Nit and zero Gaps or Partials. Dawn's recurring "It never is" gets quietly disconfirmed by 11 AM.

### [02 — Mercator Health System (healthcare)](02-mercator-health-system.md)

A top-20 integrated health system. HIPAA-covered. FDA-regulated. **TesseraSeal turned on 90 days ago for the sepsis clinical decision support model only.** Everything else is still legacy plumbing. The team experiences both worlds in one day.

### [03 — Stelvio Industrial (manufacturing)](03-stelvio-industrial.md)

A specialty steel manufacturer, ~$2.1B revenue, third-generation family-owned, DoD prime supplier. **Partial TesseraSeal deployment** across three tiers — AI side fully chained, OT side mutable, IT business side on baseline-diary plumbing.

### [04 — Atrio Banking Platform (BaaS, multi-tenant)](04-atrio-banking-platform.md)

A Banking-as-a-Service infrastructure provider, ~580 employees, Charlotte-based. Operates the technology stack for **47 fintech programs running under 12 sponsor banks**. **TesseraSeal in full multi-tenant deployment for 24 months.** Three sponsor-bank state examiners + the OCC + CFPB are all in the building this week. 1,410 verifier runs over 30 days × 47 tenants — all PASS.

### [05 — Helmstad BioSciences (biopharma)](05-helmstad-biosciences.md)

A mid-size oncology biopharma, ~$1.2B revenue, Cambridge-based. **TesseraSeal on the AI clinical-trial-eligibility screening tool for 4 months.** Pre-inspection readiness audit ahead of an FDA bioresearch monitoring inspection in 6 weeks. ALCOA+ walked attribute-by-attribute on the AI side; CRO data-feed boundary is the inspection risk.

### [06 — Pacific Crescent Power & Gas (utility)](06-pacific-crescent-power.md)

An investor-owned utility serving the Pacific Northwest, ~3.1M customer accounts, multi-state. **TesseraSeal on the AI gas-pipeline leak-detection system for 9 months.** OT side legacy. The 1:00 PM scene plays out around a real Brentwood alert that turns out to be a real (small, safely sealed) leak. Public-safety stakes are the day's emotional weight.

### [07 — Olmstead University (higher-ed research)](07-olmstead-university.md)

A private R1 research university, ~28,000 students, ~$1.4B endowment. Affiliated teaching hospital. **TesseraSeal on the undergraduate admissions AI screening for 11 months** — deployed after a civil-rights firm's disparate-impact threat letter. Multi-framework readiness audit. The 3:00 PM reconciliation finds two override-down decisions where the rationale-field history is gone — both on the suspect class the threat letter named.

## The stories — international engagements

### [08 — NetiVa Intelligence (Tel Aviv)](08-netiva-tel-aviv.md)

An Israeli AI company providing financial-market intelligence and AML tooling to 23 Tier-1 banks worldwide. Founded by Unit 8200 alumni. **TesseraSeal in production for 14 months. Multi-tenant.** Day 1 of a 3-day vendor-management evaluation commissioned by Heritage Pacific Bank. **Only Dawn, Luis, and Chen travel to Tel Aviv; the other five join by video bridge from US Eastern Time when the afternoon overlap window opens.** Bank of Israel + ISA + INCD coordination posture. Nation-state threat model assumed. The 4:30 PM "Heritage Question" pivots from technical confirmation to vendor-risk decision support.

### [09 — Sun-Won Cosmetics Group (Seoul + Taipei)](09-sun-won-cosmetics-korea-taiwan.md)

A Korean K-beauty retail group with 300+ stores in Korea and 80+ in Taiwan. KOSPI parent + Taipei Exchange subsidiary. **TesseraSeal in production for 16 months across 4 use cases.** The team splits — Dawn, Raj, Diana, Tom in Seoul; Elena, Mike, Luis, Chen in Taipei. Three regulators reviewing the same chain artifacts: PIPA + PDPA + FSS for a BNPL consumer-finance arm + Taiwan FSC for the Taipei Exchange listing. The cross-border data-flow basis on the inventory tenant is the day's most interesting Partial.

### [10 — Salt Pond Toys (Rhode Island + China + Los Angeles)](10-salt-pond-toys-rhode-island.md)

A family-owned mid-size toy manufacturer, ~$320M revenue, Newport. Manufacturing partnership with three Chinese factories in Guangdong; logistics through the Port of Los Angeles. **TesseraSeal in production for 11 months across 4 use cases** including AI quality-vision at Shenzhen and customs-entry-filing AI at LA. The team splits — most in Newport, two in LA, with Salt Pond's Shenzhen GM joining by video bridge from his Guangdong evening. The 3:00 PM recall-readiness exercise on lot 25-D-0492 produces a 14-minute end-to-end trace; Mary Catherine's 2024 lead-paint scare comparison anchors the day. The bonded-carrier maritime leg is the documented out-of-chain handoff.

### [11 — Eberhardt × Lumière (Germany + France joint engagement)](11-eberhardt-lumiere-germany-france.md)

A German automotive-electronics Mittelstand company (Eberhardt Werkstoffe, Stuttgart) and its French AI consultancy partner (Lumière AI, Paris). **Both companies are TesseraSeal users**; the chain extends across the partnership boundary via cross-vendor anchors. EU AI Act enforcement readiness + a BMW joint-supplier audit + a 2024 model-drift incident in the rear view. The team splits between Stuttgart and Paris (same time zone, one-hour TGV apart). The 12:00 PM joint working lunch produces the engagement's signature moment — Chen demonstrates a live byte-equal hash-match between Eberhardt's model-handover entry and Lumière's corresponding model-build entry. Three chains, one root-cause path.

## Industry coverage

| Story | Industry | Country | Posture |
|---|---|---|---|
| (baseline) | Mid-size financial services | US | None — full legacy |
| 01 | Banking | US | Full single-tenant |
| 02 | Healthcare | US | AI only, just started |
| 03 | Manufacturing | US | Three-tier (AI / OT / IT) |
| 04 | Fintech BaaS | US | Full multi-tenant |
| 05 | Biopharma | US | AI only, recently started |
| 06 | Utility | US | AI only, post-incident |
| 07 | Higher-ed research | US | AI only, lawsuit-motivated |
| 08 | AI vendor / financial-services adjacent | Israel | Full multi-tenant, nation-state threat model |
| 09 | Retail / K-beauty | Korea + Taiwan | Full multi-jurisdiction |
| 10 | Consumer products / toys | US + China + LA | Full multi-location |
| 11 | Automotive electronics + AI consultancy | Germany + France | Full cross-vendor |

## Reading these

The stories are written as day-of narratives — timestamped scenes, dialogue, internal thought, callout boxes for findings. They are deliberately concrete: the verifier output is real to the spec, the regulators cited are real, the systems named are real product categories the team would actually encounter. Names of companies and individual people are fictional.

Use them to:

- Calibrate what a *good* TesseraSeal deployment audit looks like across single-tenant (Northbridge), multi-tenant (Atrio), nation-state-threat (NetiVa), multi-jurisdiction (Sun-Won), multi-location (Salt Pond), and cross-vendor (Eberhardt × Lumière) postures
- Calibrate what a *partial* deployment looks like at the system boundary (Mercator, Helmstad, Pacific Crescent, Olmstead) and at the maturity boundary (Stelvio)
- Understand the auditor's chair: what evidence the team asks for, what answers count as evidence, where the chain helps and where it does not
- Train internal audit, examination prep, vendor due diligence, or onboarding for a chain rollout
- Brief executives on the "ROI shape" of a chain rollout — Northbridge's debrief, Atrio's coordinated-examiner-room scene, and Eberhardt × Lumière's live hash-match are the bookends

The stories are not training material in the regulator-pack sense. The audit-procedures, examiner-quickstart, and litigation-support documents are where the prescriptive work lives. These narratives are the human side — the people, the day, the cadence.

## Recurring threads

Across the eleven stories, certain through-lines repeat:

- **Dawn's "It never is."** Calibrated to context at every kickoff. By NetiVa it has become "It never is. But under the INCD threat model, even when it is — you stress it harder." By Eberhardt × Lumière it lands as "It never is. Today the question is whether the joint chain holds at the handoff."
- **The drive-in monologue** (or elevator monologue, or video-bridge monologue, depending on geography). Dawn and Raj — or Dawn and Tom, or Dawn alone — review prior engagements at the start of each new client. The thread accumulates. By Story 11, Dawn names ten prior contexts and slots Eberhardt × Lumière as a more rigorous version of Helmstad's CRO PGP-signed-PDF approach.
- **The 4:30 PM client question.** Each engagement has a moment where the company's CAE asks Dawn a high-stakes specific question — FDA inspector reviewability (Helmstad), public-safety evidentiary defensibility (Pacific Crescent), litigation-defense scope (Olmstead), back-out cost under nation-state-coordinated incident-response (NetiVa), three-regulator simultaneity (Sun-Won), 24-hour CPSC recall (Salt Pond), root-cause attribution across two suppliers (Eberhardt × Lumière). Dawn's answer is always sober, always concrete.
- **Tom's posture shift.** In the diary baseline he was "sweating." Across these eleven, he has been "pleased" (Northbridge), "negotiating coverage extensions" (Mercator, Helmstad), "mediating the IT/OT culture gap" (Stelvio, Pacific Crescent), "navigating jurisdictional splits" (Olmstead), "moderating cross-time-zone debriefs" (NetiVa, Salt Pond, Sun-Won, Eberhardt × Lumière), and "running quiet shop in a coordinated-examiner room" (Atrio).
- **The reconciliation test.** Always 5-10 records traced end-to-end, always the load-bearing scene. The trace shape varies by engagement — single-jurisdiction (Northbridge), bifurcated (Mercator), three-tier (Stelvio), multi-tenant batch (Atrio's 1,410-run), cross-jurisdiction (Sun-Won), cross-location with bonded-carrier handoff (Salt Pond), cross-vendor anchor verification (Eberhardt × Lumière), three-day rolling reconciliation (NetiVa).
- **Cross-vendor anchors.** First introduced as Helmstad's PGP-signed CRO PDF hashed into the chain. Refined into Salt Pond's Bureau Veritas CPSIA-cert anchor. Climaxes at Eberhardt × Lumière as bidirectional hash-equality between two independent chains. The pattern matters more each engagement.
- **Time-zone management.** Sets in at NetiVa (US ↔ IL 7-8 hours), continues at Sun-Won (KST/CST 1 hour apart from each other, 13-14 hours from US), Salt Pond (US ET ↔ PT ↔ CST 12 hours), and Eberhardt × Lumière (CET, no offset between cities, but the team is up before sunrise for the early-morning kickoff). Tom's moderation skill becomes a load-bearing engagement function across these.
