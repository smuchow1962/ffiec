# Round 12 — Cryptographic Engineer

## Persona

**Name.** Dr. Klaus Reinhardt.
**Role.** Principal Cryptographic Engineer at the European Central Bank's payments-resilience cryptography group. 24 years building HMAC/HKDF systems for SWIFT, T2/T2S, and TIPS. Reviewer for ENISA cryptographic guidance.
**Reading angle.** Cryptographic substance, post-quantum dual-algorithm transitional period encoding, per-tenant binding integrity, fingerprint-before-MAC ordering, canonical-form exclusion, cross-implementation reproducibility.

I have reviewed the v1.0 spec text, design docs 02/03/04/07/08, the test-vector READMEs, the primary fixture `chain_vectors.json`, and all 16 negative-case `description.md` files. I have not been told what previous reviewer rounds said about this material. The review below reflects only what I read in the shipped artifacts.

The substance is strong. The HMAC/HKDF construction is right. The fingerprint-before-MAC ordering is right and the load-bearing test cases (N006, N014) prove it. The sign_payload now binds `algorithm`, `format_version`, and `hkdf_inputs_digest`, which closes the algorithm-confusion class properly. My questions concentrate on the *new* dual-algorithm material in spec §4.2 / §7 step 11 cases (a)-(e), one stale section number in design 02, and the consistency of the design-doc schema tables with the spec they declare to mirror.

## Question 1 — Dual-algorithm seal encoding: are `signatures[]` entries decoded against the same `sign_payload` bytes, or do they each have their own algorithm-line in their own payload?

**What the spec says.** Spec §4.2 defines `signatures` as "list of `{algorithm, signature}`". Spec §4.3 defines `sign_payload` with `algorithm` as line 2 of the seven-line text payload. Spec §7 step 11 dual-algorithm cases (a)-(e) talk about "both signature validations" without restating which payload bytes each signature covers.

**The cryptographic question.** When an institution co-signs the seal under Ed25519 *and* (say) Dilithium-2 during the transitional period, are there:

- **Variant A.** ONE `sign_payload` per seal — the `algorithm` line on the payload is the *primary* algorithm (the one in the top-level `algorithm` field), and BOTH signatures cover that one byte string?
- **Variant B.** TWO `sign_payload` strings — one with `algorithm = "ed25519"` on line 2, one with `algorithm = "ml-dsa-65"` on line 2 — and each signature in `signatures[]` covers its own algorithm-bound payload?

The two are not equivalent under the algorithm-confusion threat model. **Variant A leaks the algorithm-binding property the §4.3 design rationale promises.** A pure Variant A means an attacker who breaks Ed25519 in 2032 and recovers the institution's Ed25519 private key can also forge a "Dilithium" signature over the same `sign_payload` byte string by simply re-signing it with a stolen Dilithium key against a public-key entry whose `algorithm` field happens to read "ml-dsa-65" — the bytes the attacker signs are identical to what a legitimate Dilithium signer would sign, because the algorithm line on the payload was filled in by the *primary* algorithm not by *this* algorithm. Variant B closes the gap: each per-algorithm signature is over a per-algorithm-bound payload, so an Ed25519-line payload cannot be replayed as a Dilithium signature target.

**The §4.3 prose** ("`algorithm` is the signature algorithm identifier in force for this seal") and the **§7 step 11 case (e) reasoning** ("one of the algorithms has been broken… the un-broken algorithm's signature still provides integrity assurance") both *imply* Variant B — case (e)'s logic only holds if the un-broken algorithm's signature was over its own algorithm-bound payload, otherwise the attacker who broke algorithm X can also produce the algorithm-Y "valid" half. But the spec text does not say Variant B explicitly.

**What I need.** One sentence in §4.2 (in the `signatures` row's notes) or in §4.3 that says: *"Each entry in `signatures[]` covers an independently constructed `sign_payload` whose line-2 `algorithm` field equals that entry's `algorithm` identifier. The top-level `signature`/`algorithm` fields hold the primary algorithm's signature over its primary-algorithm-bound payload; the entry in `signatures[]` for the primary algorithm is byte-identical to the top-level pair."*

Without that sentence, two conforming implementations could ship Variant A (cheaper — one HSM call, one payload string) and Variant B (correct — N HSM calls, N payload strings), and the test vectors would not catch the divergence because there is no dual-algorithm test vector yet. This is the question I would block on at ENISA.

**Status: Gap.**

## Question 2 — `chain_vectors.json` carries no dual-algorithm fixture; negative cases N004/N005 cover only the single-algorithm step-11 path

**What the spec says.** Spec §7 step 11 case (e) is documented as load-bearing: it is the case where one algorithm validates and the other does not, and the working paper records both rows. The case (e) reasoning is the strongest argument for why the dual-algorithm transitional period exists at all — without it, the spec could just say "rotate to Dilithium when Ed25519 breaks" and skip the transitional protocol.

**What the corpus has.** Looking at `spec/test-vectors/chain_vectors.json`: `inputs.algorithm = "ed25519"`. The `expected.sign_payload_single_text` and `sign_payload_rotation_text` both have `ed25519` on line 2. There is no `signatures[]` exemplar, no co-signed payload, no Dilithium fixture. The negative directory has N001-N016, none of which exercise:

- Case (a) — both valid (PASS)
- Case (b) — single signature in declared dual posture (PASS-WITH-ANOMALY)
- Case (c) — algorithm not on declared posture list (FAIL strict / anomaly non-strict)
- Case (e) — both present, one valid + one invalid (FAIL strict / anomaly non-strict)
- The `algorithm/key-type mismatch at signature verification` reason string from §7 step 11

N004 (signature garbage) and N005 (signature wrong-tenant) both exercise the *single*-algorithm step-11 path and report `signature verification failed`. Neither exercises the dual-algorithm dispatch.

**Why this matters cryptographically.** The dual-algorithm dispatch is the most subtle part of the verifier. A verifier that reports "PASS" because it stopped after the first valid signature (Variant A again, or naïve short-circuit logic) is non-conforming under case (e) but the corpus would not catch it. The conformance contract — "two implementations passing the corpus produce identical bytes and reach identical disposition" — does not hold across the dual-algorithm dispatch today.

**What I need.**

1. A passing-case fixture (call it `015-dual-algorithm-cosigned-seal/`) carrying both an Ed25519 signature and a second-algorithm signature over the same seal_date (with both `sign_payload` strings emitted explicitly per Question 1). The verifier MUST report PASS.
2. Three negative cases:
   - **N017** — case (b): seal carries only Ed25519 in declared dual posture → PASS-WITH-ANOMALY with the `partial-coverage seal` reason string under both strict and non-strict.
   - **N018** — case (c): seal carries an algorithm not on the institution's declared posture list → FAIL strict / PASS-WITH-ANOMALY non-strict with the `algorithm not on institution's declared posture list at seal_date {D}` reason string.
   - **N019** — case (e): seal carries Ed25519 valid + Dilithium garbage → FAIL strict / PASS-WITH-ANOMALY non-strict with the `co-signed seal failure: algorithm X validated, algorithm Y did not` reason string. The working paper records both rows.

Until these ship, the dual-algorithm dispatch is documented but not testable, and that is the gap a conformance regime cannot close.

**Status: Gap.**

## Question 3 — Design 03 §3.7 schema table does not mirror spec §4.2 after the `signatures` field was added

**What design 03 §3.7 says.** The section header reads: *"3.7 Seal record schema (normative; mirrors spec §4.2)"*. The table beneath lists 14 fields (`tenant_id` through `dev_mode`). It does **not** list `signatures`.

**What spec §4.2 says.** The schema table in spec §4.2 has 15 rows. Row 11 is `signatures` — "Optional list form for the dual-algorithm transitional period (post-quantum coexistence with Ed25519). When present, the seal is co-signed under multiple algorithms…"

**The mismatch.** Design 03 explicitly claims to mirror the spec, but the spec table now has a row the design table does not. Auditors who treat the design doc as the implementation guide (a real auditor practice on the ledger-server side) will not implement the `signatures` field. The implementation that follows design 03 verbatim is non-conformant with spec §4.2 the moment dual-algorithm posture is enabled.

**Adjacent.** Design 04 §3.2 ("Daily seal signing key") talks about *one* signing keypair per tenant, even though §2.3 mentions "the spec contemplates a future spec version that adds a post-quantum signature alongside Ed25519." With the spec now formalizing `signatures[]`, the design-doc text in §3.2 is stale: institutions in dual-algorithm posture hold *two* keypairs per tenant (one Ed25519, one ML-DSA-65 or SLH-DSA), each in HSM custody, both used at seal time.

**What I need.**

1. Design 03 §3.7 table gains the `signatures` row, copied verbatim from spec §4.2. Add a "Why `signatures` is a list" subsection (parallel to "Why `key_versions` is a list" and "Why `hkdf_inputs_digest` is on the seal") explaining the dual-algorithm posture, the per-algorithm-bound payload (per Question 1), and the working-paper convention.
2. Design 04 §3.2 picks up a fourth bullet: under dual-algorithm posture, the institution operates two HSM-resident keypairs per tenant; each algorithm's seal-signing role grants `sign` only against its own keypair; the seal-job submits N independent signing requests (one per algorithm in the institution's declared posture).

Without these, design 03 is a stale mirror and design 04 has a hole where the dual-algorithm operational picture should be.

**Status: Gap.**

## Question 4 — Design 02 carries one stale subsection number that breaks the §10 numbering

**What I asked you to confirm.** §8 (What can go wrong) and §9 (Multi-process run semantics) numbering is now clean.

**What I read.** §8 is clean: §8.0, §8.1, §8.2 — no gaps, no out-of-order. §9 is clean: §9.1, §9.2, §9.3 — no gaps. **Both pass the focused check.**

**What I noticed adjacent.** §10 ("Auditor's-lens review") has one child subsection at the bottom labeled `### 9.1 Closed findings from the Herald upgrade audit`. The "9.1" label is stale — that subsection lives inside §10, so it should read `### 10.1`. The auditor reading the doc top-to-bottom hits §9.1, §9.2, §9.3, then §10, then §9.1 again, which is a parser-level inconsistency the same way a stray `;` is in mermaid: the document still renders, but the cross-reference machinery (any tool that builds a TOC, any auditor who scrolls back to "see §9.1") gets confused.

This is outside the focus the user named, but I flag it because it is on the same numbering theme and because the fix is one character.

**What I need.** The line `### 9.1 Closed findings from the Herald upgrade audit` becomes `### 10.1 Closed findings from the Herald upgrade audit`.

**Status: Partial.** The §8 and §9 question I was asked is **Answered** cleanly. The §10 sibling issue is a Partial because it is the same numbering surface and one character closes it.

## Question 5 — sign_payload algorithm-line in design 07 §4.3 pseudocode

**What I checked.** Design 07 §4.3 lines 217-224:

```
sign_payload = "ffiec.chain-of-custody.v1\n" ||
               seal.algorithm                || "\n" ||
               seal.format_version           || "\n" ||
               T                             || "\n" ||
               iso8601(D)                    || "\n" ||
               hex(seal.merkle_root)         || "\n" ||
               hex(seal.hkdf_inputs_digest)
```

**Does it match spec §4.3?** Spec §4.3 lines 175-183 layout: spec-prefix `\n`, algorithm `\n`, format_version `\n`, tenant_id `\n`, iso8601_date `\n`, hex(merkle_root) `\n`, hex(hkdf_inputs_digest). Each line ends with one `\n`. The verifier pseudocode produces byte-identical output to the spec construction. Confirmed against `chain_vectors.json` `expected.sign_payload_single_text`: `"ffiec.chain-of-custody.v1\ned25519\nv1\ntenant-ffiec-test-1\n2026-05-06\n927adc88...\n6f8a5005..."`. Seven fields, six interior `\n` bytes, no trailing `\n` after the last field. This matches the spec construction byte-for-byte.

**What's the dispatch comment doing.** The "Algorithm dispatch: resolves the public key's algorithm from public_key_id" comment is correct but applies to the *dispatch* layer (which Verify implementation runs), not to the payload construction. That distinction is right. The `IF NOT verify_for_algorithm(seal.algorithm, public_key, sign_payload, seal.signature)` line carries the algorithm into the Verify call, which is the correct binding.

**What's missing for the dual-algorithm posture.** Per Question 1 and Question 2, the §4.3 pseudocode shows only the single-signature path. There is no loop over `seal.signatures[]`, no per-algorithm payload reconstruction, no recording of per-algorithm validation results into the working paper. Under the spec's dual-algorithm posture this pseudocode is incomplete.

**What I need.** §4.3 picks up a short follow-up block:

```
# Dual-algorithm dispatch (when seal.signatures is present per spec §4.2 + §7 step 11)
IF seal.signatures is not None:
  per_algorithm_results = []
  declared_posture = institution_declared_algorithm_posture(tenant_id, seal_date)
  FOR each entry in seal.signatures:
    # Question 1 disposition: each entry covers its own algorithm-bound payload
    entry_payload = build_sign_payload(
      entry.algorithm, seal.format_version, T, iso8601(D),
      seal.merkle_root, seal.hkdf_inputs_digest
    )
    public_key_for_alg = resolve_public_key(seal.public_key_id, entry.algorithm)
    valid = verify_for_algorithm(entry.algorithm, public_key_for_alg, entry_payload, entry.signature)
    per_algorithm_results.append({"algorithm": entry.algorithm, "valid": valid})
    IF entry.algorithm not in declared_posture:
      record_anomaly("algorithm not on institution's declared posture list at seal_date {D}")
  apply_step_11_disposition(per_algorithm_results, declared_posture, strict_mode)
  # disposition implements cases (a) (b) (c) (d) (e) per spec §7
```

Once §4.2 commits to Question 1's Variant B (per-algorithm payload), this pseudocode is ready to lock; before that commitment, the `entry_payload` line is undefined and an implementer guesses.

**Status: Partial.** The single-algorithm pseudocode is correct and matches the spec. The dual-algorithm extension is missing.

## Per-role roll-up

| Item | Status | Closes when |
|---|---|---|
| Q1: dual-algorithm payload variant explicit (A vs B) | Gap | Spec §4.2 or §4.3 names per-algorithm-bound payload (Variant B) explicitly |
| Q2: corpus + negative tests for dual-algorithm cases (a) (b) (c) (e) | Gap | `015-dual-algorithm-cosigned-seal/` plus N017/N018/N019 ship with byte-level fixtures |
| Q3: design 03 §3.7 + design 04 §3.2 mirror spec §4.2 dual-algorithm | Gap | Design 03 table gains `signatures` row + rationale; design 04 §3.2 names the two-keypair operational shape |
| Q4: design 02 §8/§9 numbering clean; §10's stale `### 9.1` heading | Partial | The asked question (§8/§9) is Answered. Sibling `### 9.1` → `### 10.1` closes the partial |
| Q5: sign_payload single-algorithm pseudocode + dual-algorithm extension | Partial | Single-algorithm pseudocode is correct. Dual-algorithm extension lands once Q1 is decided |

**Stopping criterion.** 0 gaps + 0 partials. Currently at 3 Gaps + 2 Partials. The gaps cluster on dual-algorithm — Q1 (semantics), Q2 (corpus), Q3 (design-doc consistency with the new spec field) are the same finding seen from three angles. Resolving Q1 unblocks Q2 (the test fixtures are computable once the variant is chosen) and Q3 (the design-doc text follows the spec text). Q4 and Q5 close with one character and one pseudocode block respectively.

The non-dual-algorithm material in this round — HMAC/HKDF construction, fingerprint-before-MAC ordering, canonical-form exclusion, mid-write truncation refusal, the load-bearing N006/N013/N014 negative cases, the chain_vectors.json byte-level cross-implementation contract, the design 07 single-algorithm pseudocode — all reads as cryptographically sound and internally consistent. The single-algorithm verifier I would be willing to certify under ENISA review today on the strength of what is shipped. The dual-algorithm verifier I would block on Q1.
