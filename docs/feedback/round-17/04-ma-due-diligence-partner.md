# Round 17 Outside Reviewer — M&A IT Due-Diligence Partner (Big-4)

**Reviewer persona:** Partner at a Big-4 transaction-services practice, 18 years on cross-border IT diligence and AI-asset acquisitions. Engagements run 3 weeks; I have killed two deals in 18 months over evidence-trail gaps that would not survive post-close regulator review.
**Spec version reviewed:** v1.0a (chain-of-custody-v1.md, 1358 lines)
**Date:** 2026-05-07
**Question budget used:** 6 of 8

## Overall verdict

I would accept v1.0a as the substrate of an AI-asset target's evidence chain, with two carve-outs for representation drafting. Cryptography is conservative (FIPS-current; three-layer §1.4; second-preimage Merkle §4.2). Post-close evidentiary posture is unusually well-developed for a v1 (§1.1 Daubert, §10.13 retention, §10.21 cross-vendor handover). The chain-coverage map at §10.19 is the single most useful section I have read in an audit-evidence spec — it lets me write the "systems NOT covered" rep as an enumerated schedule rather than a hand-wave. Carve-outs: (1) entity-change handling (renamed entities, HSM custody succession, mid-period acquirer takeover) is not normative in the spec body — §13 sends the M&A reader to `docs/m-and-a-handoff.md`, which is informative; for a binding rep I need normative language. (2) §10.17 `signatories` captures names but not the signatory's legal-entity affiliation at signature time — material when the same individual signs pre-close under Target authority and post-close under Acquirer authority.

## Findings

### Gaps (deal-killer absences)

**G1 — No normative entity-succession/rename procedure.** §3 character class and §3.1 legacy patterns are naming hygiene; they do not define a tenant-rename / entity-succession event. Nothing in the spec body says "this `tenant_id` was operated by Target through 2026-05-06 and by Acquirer from 2026-05-07; continuity preserved by mechanism X." I can only write the rep if I attach `m-and-a-handoff.md` as a schedule, which seller's counsel resists because it is informative. **Closer:** lift a §10.22 "Entity succession" into normative spec, defining a `chain.entity_succession` operational event (from-entity, to-entity, effective-UTC, dual-signature) bound under the seal of the transfer day.

**G2 — §10.21 does not bind the contractual handover instrument.** `audit.model_handover.*` captures artifact, card, fairness audit, `provider_chain_entry_id`. It does NOT capture the SHA-256 of the model-supply contract or DPA under which delivery happened. For an Acquirer replacing the model vendor post-close, "which contract version was in force at this delivery?" must be answerable from the chain alone. **Closer:** add `contract_id` / `contract_version` / `contract_hash_sha256` parallel to §4.4.1 `audit.cross_border_transfer.*`, which already establishes the precedent.

### Partials (something present but insufficient)

**P1 — §10.17 signatory affiliation not bound.** `chain.partition_ceremony_attended` records `signatories` and `witness` as `{role, name}`. When Target's CISO signs Day -1 under Target authority and the same individual signs Day +1 under Acquirer authority, the chain cannot distinguish the two. **Closer:** add REQUIRED `entity_affiliation` per object in `signatories` and `witness`.

**P2 — §10.20 commits the number, not the manifest.** `training_data_retention_floor_days` is an integer commitment; provider adherence is contractual, not cryptographic. For deal-window lookback the chain says "provider committed to 540 days" but does not bind the manifest of training-shard hashes. **Closer:** add OPTIONAL `audit.model_handover.training_shard_manifest_sha256` so the provider's enumerated hash list is bound at handover and the post-close auditor recomputes against surviving shards.

**P3 — §10.19 map not versioned or chain-anchored.** Map is a static CC8.1 enumeration; no version+effective-date stamp, no chain-anchored hash. In an 18-month lookback I need to know which map version was in force on a given date; without that, seller can produce a current map that does not describe the system as it operated 14 months ago. **Closer:** add `chain.coverage_map_published` operational event under §10.2 carrying map version, effective-UTC, and SHA-256 of the canonical map.

### Nits (cosmetic / wording / cross-reference)

**N1 — §13 lists M&A handoff only under "Bank chain adopter" and "Bank vendor-management team."** Acquirer-side reader (neither role) finds the pointer only by inference. Add a dedicated "Acquirer-side IT due-diligence" reader heading.

**N2 — §10.21 `audit_report_languages` plural is correct but the schema row offers no example.** Fresh readers without the auditor-story context miss why the field is plural-required.

## Strengths

- **§10.19 chain-coverage map.** Best section in the spec for rep drafting. Five enumerated boundary categories is exactly the schedule shape reps need; `audit.external_artifact.*` closes the boundary with hash-anchoring rather than "trust the seller."
- **§10.21 `provider_chain_entry_id`.** Bidirectional verification lets Acquirer's verifier walk Target's chain into the model-vendor's chain — the mechanism that makes a vendor-replacement rep testable rather than papery.
- **§1.1 + §1.3 + §1.4 Daubert composition.** IT witness can answer the four *Daubert* questions from shipped artifacts. Rare in v1 specs.
- **§4.4.1 `audit.cross_border_transfer.*`.** Schrems II reps survive post-close because the lawful-basis contract version is hash-bound on the chain entry.
- **§10.10 IKM rotation across seal boundary.** Post-close key transfer (Acquirer provisions its own HSM, rotates IKM) does not break chain continuity for the lookback.

## Questions to the spec authors

1. **§10.19 versioning.** Is there a normative requirement that the coverage map be version-stamped and chain-anchored? If not, see P3.
2. **§13 acquirer-side reader.** Is a dedicated Acquirer-IT-DD reader heading planned, distinct from "Bank chain adopter" and "Bank vendor-management team"?
3. **§10.21 contract binding.** Was `contract_id / contract_version / contract_hash_sha256` considered and deliberately omitted, or is this a v1.x candidate? §4.4.1 already establishes the pattern.
4. **§10.17 customer-bank entity change.** When the partition is held by a SaaS vendor on behalf of `customer_bank_id` and the customer-bank itself undergoes an entity change post-ceremony, is the chain re-attributable to the new legal entity by any normative mechanism, or only via the informative M&A handoff?
5. **§10.20 enforcement.** What is the spec's posture when a provider violates the retention-floor commitment post-handover but pre-close? Is there a `chain.retention_floor_violation_observed` event under §10.2, or is detection purely institution-side?
6. **Multi-tenant carve-out.** For a Target that is one fintech program on a multi-tenant platform (Atrio shape, 47 programs × 12 sponsor banks), §10.1 establishes per-tenant IKM isolation and §4.2 per-tenant seals — but is there a normative §10.x section for "Acquirer takes possession of one tenant's seals + IKM custody handoff while other tenants continue uninterrupted on the same HSM cluster"? Mechanics are implied by §10.5 + §10.17 `partition_handle`; I want one normative section to cite.
