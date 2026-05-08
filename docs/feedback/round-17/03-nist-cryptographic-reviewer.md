# Round 17 Outside Reviewer — NIST-Lineage Cryptographic Reviewer

**Reviewer persona:** Senior cryptographer at a national lab; contributor to NIST SP 800-108 KDF guidance and to RFC 8785 (JCS). Reviews wire-form and binding-format specs for regulated, cross-border, multi-decade evidentiary settings.
**Spec version reviewed:** v1.0a (chain-of-custody-v1.md, 1358 lines)
**Date:** 2026-05-07
**Question budget used:** 8 of 8

## Overall verdict

I would sign off, with two reservations as Partials below. The construction is conventional in the right places: HMAC-SHA-256 over JCS-canonical JSON (§4.1), an RFC 6962 Merkle root with `0x00`/`0x01` domain prefixes (§4.2), Ed25519 inside FIPS 140-2 Level 3 custody (§4.3). HKDF input binding uses an unambiguous concatenation under a constrained `tenant_id` character class (§3, §4.1). The §4.3 v1.0a `sign_payload` is byte-pinned tightly — explicit `0x0A` separators, lowercase hex, "no trailing newline on the terminal field" (line 352), `sign_payload_version` bound into line 2. The dual-algorithm Variant B per-algorithm `sign_payload` rule (§4.2, §4.3.2, §7 step 11 case (e)) closes algorithm-confusion under AND-security. **Highest-risk gap:** no per-event MAC algorithm agility for v1.0. §4.1 line 201 acknowledges this and points at the §4.3.2 90-day emergency-patch SLA — operationally tight for 7-year-retention chains.

## Findings

### Gaps

**G1 — `kms_handle_uri` is provenance-only, not chain-bound.** §5 line 653 lists `kms_handle_uri` in the canonical-form exclusion set; the field is recorded on every entry (§4.4 line 427) but does not contribute to the per-event MAC and is not in the v1.0a `sign_payload` either. §10.7 treats it as a regulator-visible line (lines 882-883), but its integrity is procedural-only — an attacker with ledger-write access who flips `"plaintext-dev"` to `"aws-kms:arn:..."` does not invalidate any MAC. **Close it by:** binding a per-day list of distinct `kms_handle_uri` values into the v1.0a `sign_payload` alongside `key_versions`-style cross-checking, or including the URI in the canonical bytes the per-event MAC covers.

**G2 — `key_versions` on the seal record is cross-checked, not signed.** §7 step 11 (line 749) says variable-length encoding would complicate signature reproducibility, and falls back on a post-validation cross-check. Under witness-verifier mode (line 762, no IKM) the §7 step 8 fingerprint check is skipped, so the cross-check becomes the only line of defense — and that cross-check itself depends on per-event `key_version` integrity, which is only guaranteed once IKM lookup succeeds. The "variable-length" rationale is weak: `sign_payload` already binds variable-length `tenant_id`. A deterministic comma-separated-decimal encoding would close it.

### Partials

**P1 — Per-event MAC algorithm agility absent for v1.0.** The seal signature has a fully-specified dual-algorithm posture (§4.3.2, §4.2 `signatures` list, §7 step 11 cases (a)-(e)). The per-event HMAC has nothing equivalent — §4.1 line 201 names this as future-version work and points at the 90-day emergency SLA. A 90-day migration across 7-year retention is operationally tight. **Close it by:** a normative dual-algorithm window for the per-event MAC analogous to `signatures` — a parallel `payload_hash_alt` under a second algorithm during a transitional period. Without it, a compromised SHA-256 leaves no in-band cryptographic option for verifying historical chains.

**P2 — JCS conformance is testable but not self-tested by the verifier.** Test vector 008 pins float canonicalization, surrogate pairs, NFC-vs-NFD non-normalization, U+007F-not-escaped, UTF-16 code-unit key ordering, NaN/Infinity rejection. §5 line 663 makes 008 mandatory. **What's missing:** the verifier does not exercise JCS conformance before declaring PASS. An implementer who shipped a chain with all-ASCII events under a non-conformant canonicaliser passes §7 silently; divergence surfaces only on the first non-ASCII event years later. A pre-flight self-test in §7 (canonicalize a baked-in fixture, assert against a baked-in expected hash) would catch this at verifier startup.

**P3 — Partition-ceremony coupling (§10.17) is procedural.** `chain.partition_ceremony_attended` hash-anchors the scanned PDF and binds signatories under the per-event MAC. That is procedurally adequate. **What is not cryptographic:** the ceremony itself is not coupled to the chain by HSM-side attestation. There is no requirement for an HSM-emitted attestation token (Thales, Entrust, CloudHSM all expose categories of these) bound into the event. The scheme proves "the chain says the ceremony happened"; it does not prove "the HSM agrees."

### Nits

**N1** — Line 285's `signatures` row is a multi-paragraph spec inside a markdown table cell; lift to a dedicated §4.2.x subsection.

**N2** — §4.1 line 187 cites SP 800-56C §5.4. SP 800-108r1 §4.1 ("Other Input") is the more direct reference for KDF context binding; cite both.

**N3** — §4.2 line 267 correctly notes empty-day root collision is benign. Worth a one-line note that the property extends to identical-event-set-and-ordering days.

## Strengths

- **§4.1 SHA-256 length-extension audit (lines 191-199):** five distinct call sites, each individually argued safe with the reason (domain separator, truncation, or HMAC-mediation). Strong artefact.
- **§4.3 v1.0a `sign_payload` byte pinning.** Lowercase hex, single `0x0A` separators, terminal-field-no-trailing-newline, reproducible byte-length count, `sign_payload_version` bound into line 2. Forensically reproducible across implementations and across decades.
- **Test vector 008.** The seven failure-mode hints in description.md are precisely the traps an experienced reviewer would write down independently.
- **§7 step 9's `expected_prev_hash` rule (line 732):** the MAC input uses the structurally-walked value, not the entry's claimed `prev_hash`. Subtle and correct — closes a latent footgun under future structural-check relaxations.
- **§3 tenant_id character class.** The `|` byte (0x7C) cannot appear inside `tenant_id`, so the HKDF `info` concatenation is unambiguously parseable. Enforced at both ends.
- **§4.3.2 dual-algorithm Variant B / AND-security.** Each algorithm covers its own algorithm-bound `sign_payload`; case (e) is FAIL under `--strict`. Correct posture for a transitional period.

## Questions to the spec authors

1. **Per-event MAC dual-algorithm window.** Open to a `payload_hash_alt` field carrying a second-algorithm MAC during a transitional period, parallel to the seal's `signatures` list? (P1)

2. **`kms_handle_uri` binding.** Why is the field excluded from both the per-event canonical bytes and the v1.0a `sign_payload`, given §10.7 treats it as a regulator-visible line? (G1)

3. **`key_versions` signature binding.** §7 step 11 cites variable-length encoding, but `sign_payload` already binds variable-length `tenant_id`. Is there a stronger reason for cross-check rather than signature, given witness-verifier mode cannot run the fingerprint check? (G2)

4. **HMAC migration roadmap.** §4.1 line 187 names per-tenant salts as a candidate v1.x posture; line 201 names per-entry algorithm dispatch as future v1.1/v2.0. Near-term roadmap, or held until a SHA-256 weakening?

5. **Pattern A invariant 5 cache-freshness (§10.15 line 1023).** Working group view on operational cost for institutions whose replication pipeline state is not directly queryable? Push-update mechanism documented?

6. **Verifier-side JCS self-test.** Open to a normative pre-flight JCS self-test in §7 to catch canonicaliser drift at verifier startup? (P2)

7. **HSM-side attestation in §10.17.** Considering a normative HSM-emitted attestation token bound into `chain.partition_ceremony_attended` alongside the PDF hash? (P3)

8. **Forked-chain disambiguation.** Two chain files sharing `(tenant_id, run_id)` but diverging at `seq=N` each verify under §7 in isolation; fork detection appears to rely on the institution's append-only storage discipline. Is fork detection out of scope for the chain itself? Could not locate a normative section naming this explicitly.
