# Round 13 — Cryptographic Engineer

## Persona

**Name.** Dr. Saoirse Ní Chonchúir.
**Role.** Principal Cryptographic Engineer at Bank of Ireland's payments-cryptography group. 21 years building HMAC/HKDF systems for SEPA and T2/T2S, with three of those years on the cross-border instant-payment messaging. ETSI TC CYBER reviewer; authored two ETSI TR contributions on key-derivation reproducibility for clearing systems.
**Reading angle.** Dual-algorithm encoding (Variant B vs Variant A), per-algorithm `sign_payload` reconstruction in the verifier pseudocode, the `signatures` field's propagation through spec → design 03 → design 07 and into the test corpus, the §7 step 11 case (e) severity language, fixture-level cross-implementation reproducibility.

I have read spec/chain-of-custody-v1.md (focused on §3, §4.1, §4.1.1, §4.2 with the new `signatures` list, §4.3, §4.4, §5, §6, §7 step 11 cases (a)-(e), §10), design 02-chain-construction.md, design 03-merkle-seal.md §3.7, design 04-hsm-custody.md §3.2, design 07-verifier-design.md §4.3, design 08-test-vectors.md, spec/test-vectors/README.md + chain_vectors.json + the per-case `description.md` files including 015-dual-algorithm-cosigned-seal and N017/N018/N019. I have not been told what previous reviewer rounds said about this material; the review below reflects what I read in the shipped artifacts.

I also independently re-derived three byte-level constants from `chain_vectors.json` against my own Python implementation of the spec construction:

- `hkdf_inputs_digest = SHA-256(salt || info_for("tenant-ffiec-test-1") || (32).to_bytes(4,"little"))` reproduces `6f8a5005…727d65` byte-for-byte.
- `key_fingerprint_v1 = SHA-256(utf8("tenant-ffiec-test-1") || ikm_v1)[:16]` reproduces `54c245a8…531f` byte-for-byte.
- `merkle_root_single` over the five `single_chain.payload_hash` values, computed under the RFC 6962 leaf-prefix `0x00` and node-prefix `0x01` scheme, reproduces `927adc88…eefe9` byte-for-byte.

The byte-level conformance contract for the single-algorithm path is reproducible from the spec text alone, with no shipped reference implementation in the loop. That is the cleanest property a cross-implementation contract can have, and I want to record it explicitly because it is the property the corpus exists to provide.

The substance is strong. My questions concentrate on (1) confirming the dual-algorithm encoding choice is now nailed down end-to-end, (2) the case (e) severity text in spec §7 step 11, (3) the `algorithm/key-type mismatch` reason string's exercise in the corpus, and (4) the fixture-deferral posture for cases 015 / N017 / N018 / N019.

## Question 1 — Variant B is now stated; is it stated where every implementer will see it before they wire signing?

**What the spec says now.** Spec §4.2 `signatures` row carries the explicit text:

> Each signature in the list MUST cover its own algorithm-bound `sign_payload` (Variant B): each algorithm's `sign_payload` is constructed per §4.3 with that algorithm's identifier in the second line. A single shared `sign_payload` covering all algorithms (Variant A) is non-conformant — it leaks the algorithm-confusion defense by letting an attacker present an algorithm-X signature on a payload that names algorithm-Y.

**What design 03 §3.7 says.** The `signatures` row in the design 03 schema table:

> Each entry covers its own algorithm-bound `sign_payload` (Variant B). When present, the verifier dispatches per-algorithm and produces per-algorithm validation results.

**What design 07 §4.3 says.** The dual-algorithm pseudocode reconstructs `sp_alg` per entry with `entry.algorithm` (not `seal.algorithm`) on line 2, with the comment:

> the algorithm field is `entry.algorithm`, NOT `seal.algorithm` — Variant B is normative per spec §4.2 to close the algorithm-confusion attack class.

**What case 015's description.md says.** The recipe explicitly shows two distinct `sign_payload` text blocks differing only in the second line (`ed25519` vs `dilithium3`), with the trailer note: *"Variant B is the load-bearing property closing algorithm-confusion attacks."*

**What N019's description.md says.** Provenance section: *"This case exists to lock in Variant B per-algorithm sign_payload encoding. A verifier that uses Variant A (single shared sign_payload) would mistakenly validate the Ed25519 signature against the Dilithium3 algorithm name (the algorithm-confusion attack the spec explicitly closes)."*

**My finding.** Variant B is now declared in four places — the normative spec row, the design 03 schema, the design 07 pseudocode, and the corpus stubs (015 + N019). The four places agree on the wording and on why Variant B is the load-bearing choice. An implementer who reads any one of those four arrives at Variant B. A verifier author working from the spec, the verifier-design pseudocode, AND the corpus would have to deliberately ignore three independent statements of the same rule to ship Variant A.

The only place Variant B is NOT mentioned is design 04 §3.2 ("Daily seal signing key(s)"). §3.2 now says:

> Per-tenant: one signing keypair per algorithm per tenant. For dual-algorithm transitional posture (v1.x post-quantum): one Ed25519 keypair AND one post-quantum keypair (Dilithium / SLH-DSA), both per tenant. The institution's HSM holds both keypairs; the seal job signs the day's Merkle root under both keys per spec §4.2 `signatures` list (Variant B per spec §4.2: each algorithm's signature covers its own algorithm-bound `sign_payload`).

Good — design 04 §3.2 also now references Variant B by name and by spec section. So Variant B is declared in **five** places, and design 04 closes the operational picture (two keypairs per tenant, both in HSM custody).

**Status: Answered.** Variant B is locked end-to-end across spec, design 03, design 04, design 07, and the corpus stubs. The encoding variant chosen is the one that closes the algorithm-confusion attack class.

## Question 2 — Does spec §7 step 11 case (e) carry severity-and-disposition language strong enough that an examiner reading non-strict PASS-WITH-ANOMALY does not under-react?

**What case (e) says now.** The spec text I read:

> **(e) Both signatures present, one valid + one invalid.** Under `--strict`: FAIL with `co-signed seal failure: algorithm X validated, algorithm Y did not`. Under non-strict: PASS-WITH-ANOMALY with the same reason. **The severity of case (e) is Severe regardless of bracket** — the PASS-WITH-ANOMALY disposition under non-strict reflects the spec's posture that the un-broken algorithm's signature still provides integrity assurance for downstream consumers, NOT that the failure is itself low-severity. Examiners writing up case (e) findings cite `regulator-pack/finding-language.md` row "11 (dual-algo) co-signed seal failure" (Severe MRA) regardless of which bracket the verifier reported.

**The cryptographic concern this addresses.** Without the bolded sentence, an examiner reading "PASS-WITH-ANOMALY" under non-strict could plausibly conclude that case (e) is a low-severity operational anomaly comparable to `late_binding` or `master_key_rotation_observed`. It is not. Case (e) means EITHER (i) one of the two co-signing algorithms has been broken in a way that lets an attacker forge signatures, OR (ii) one of the per-algorithm signing keys has been compromised. Both are Severe.

**The severity-cross-reference machinery.** The text directs the examiner to `regulator-pack/finding-language.md` row "11 (dual-algo) co-signed seal failure (Severe MRA)" as the authoritative finding-language entry. That cross-reference is the right shape — it puts the severity classification in the finding-language registry where the bank's MRA process will pick it up, not buried in the verifier output where it can be missed.

**The two-cause framing.** The text continues:

> This is a load-bearing case: a seal where one of the two algorithms validated and the other did not indicates EITHER (i) one of the algorithms has been broken (in which case the un-broken algorithm's signature still provides integrity assurance and the institution coordinates with the regulator on the broken-algorithm migration timeline), OR (ii) one of the seals is forged under a compromised algorithm-specific signing key (in which case the un-broken algorithm's signature confirms the un-compromised half of the chain custody). The verifier does NOT attempt to interpret which case applies; the institution's IR program does the interpretation per IR Scenario 12.

This is the right separation-of-concerns: the verifier reports the cryptographic fact (one valid + one invalid); the institution's IR program does the cause attribution. The verifier should not try to guess whether it is case (i) or case (ii) — those determinations require evidence the verifier does not have access to (the institution's HSM access logs, the regulator's published-algorithm-break notifications, the per-algorithm-key compromise indicators).

**N019's IR-disposition section** mirrors this with three branches (algorithm-break / per-algorithm-key compromise / under-investigation) and ties the 36-hour clock-start to IR Scenario 4 only when the cause attribution lands on (ii). That is consistent with the spec's posture.

**Status: Answered.** Case (e) carries the severity language explicitly, points to the finding-language registry for the MRA classification, separates the cryptographic fact (verifier reports) from cause attribution (IR program decides), and N019's description echoes the IR disposition. This is the strongest version of case (e) text I have seen across the four spec rounds whose artifacts I have access to.

## Question 3 — The `algorithm/key-type mismatch at signature verification` reason string from §7 step 11 — is it exercised by the negative corpus?

**What the spec says.** §7 step 11:

> If the seal's `algorithm` does not match the public key's algorithm (e.g. the seal claims `"ed25519"` but the resolved public key is a Dilithium key), report `algorithm/key-type mismatch at signature verification` rather than the generic message.

**Why the distinction matters.** The generic `signature verification failed` message is what the verifier reports when the bytes do not match. `algorithm/key-type mismatch` is what the verifier reports when the algorithm dispatch identifies a public key whose declared algorithm differs from the seal's claimed algorithm — a key-type confusion (the `public_key_id` resolved to the wrong key entry, perhaps because the registry was tampered with or because the institution mis-registered the key). The two failure modes have different remediation paths: the generic failure suggests data tampering; the algorithm/key-type mismatch suggests registry / configuration drift.

**What the negative corpus has.** N004 (signature garbage) and N005 (signature wrong-tenant) both produce the generic `signature verification failed`. N017/N018/N019 exercise the dual-algorithm cases (b)/(c)/(e). I read all 19 negative-case `description.md` files; none of them exercise the `algorithm/key-type mismatch` path specifically. A verifier that conflates the two reason strings (always reporting `signature verification failed` and never `algorithm/key-type mismatch`) would pass the negative corpus.

**Why this matters for the conformance regime.** The `algorithm/key-type mismatch` reason is a precision-of-reporting requirement, not a security-property requirement (the verifier still rejects the seal in both cases — the question is which message it produces). Precision-of-reporting is what lets the auditor distinguish "the institution's signing pipeline produced corrupted output" from "the institution's key registry resolved to the wrong key entry." Both are findings; they are different findings.

**What I would want.** A negative case `N020-algorithm-key-type-mismatch/` with a tampering recipe along the lines of: take case 002 (single-algorithm Ed25519 seal); modify the `public_key_id` to point at a registered Dilithium public key (with a corresponding valid Dilithium signing keypair somewhere in the registry, but NOT the one the seal was actually signed under); the Ed25519 signature bytes are present and well-formed but the resolved public key is the wrong algorithm. Expected verifier output: `algorithm/key-type mismatch at signature verification` (NOT `signature verification failed`).

This case is computable today under v1.0 (single-algorithm posture) — it does not require the v1.x post-quantum keys that 015/N017/N018/N019 wait for. The fixture deferral that applies to the dual-algorithm corpus does not apply here.

**Status: Gap.** The reason string is documented in spec §7 step 11 but not exercised by any negative-corpus case. A non-conforming verifier that always returns the generic message would pass today's corpus on this point.

## Question 4 — Fixture-deferral posture for cases 015 / N017 / N018 / N019: is it explicit enough that v1.0 implementers do not block on the missing bytes?

**What the four stub files say.** All four `description.md` files for the dual-algorithm cases carry the same Status block:

> **Stub case.** Byte-level fixture deferred to v1.x. Recipe and expected verifier outcome documented here.

**What case 015's "v1.x population guidance" section adds.** A five-step recipe for the v1.x implementer:

> 1. Generates Ed25519 + Dilithium3 (or SLH-DSA) test keypairs (clearly marked TEST USE ONLY).
> 2. Computes both algorithm-bound sign_payloads.
> 3. Signs each with the corresponding algorithm.
> 4. Publishes `inputs/`, `expected/`, including the per-algorithm public keys, sign_payloads (text and hex), and signatures.
> 5. Updates `chain_vectors.json` with `dual_algorithm_chain` section mirroring `single_chain` and `rotation_chain`.

**Why this posture is right.** v1.0 ships single-algorithm Ed25519. The dual-algorithm cases exist to lock in the verifier's dual-algorithm dispatch logic so that when v1.x ships the second algorithm, the verifier authors do not have to invent the dispatch from scratch. The fixture deferral is honest about what is shippable today (the recipe and expected verifier behavior) versus what waits for the second algorithm (the byte-level signature fixtures, which require a Dilithium / SLH-DSA implementation in the test-vector generator that does not yet exist for v1.0).

**What v1.0 implementers can do today against the stub corpus.** A v1.0 verifier author can:

- Implement the dual-algorithm dispatch path per the design 07 §4.3 pseudocode.
- Test it against synthetic fixtures the implementer constructs locally from the case 015 recipe (using whatever Dilithium / SLH-DSA library is available — probably a research-grade liboqs binding or equivalent).
- Use the expected verifier outcomes documented in the four `description.md` files as the assertion targets.

The stub corpus is sufficient to ship a v1.0 verifier whose dual-algorithm dispatch is **structurally correct and ready for v1.x byte-level fixtures**. The risk that gets deferred is implementation-specific encoding drift across the two-algorithm signature wire formats — that risk does not become real until the v1.x byte-level fixtures arrive.

**Status: Answered.** The fixture-deferral posture is explicit in every stub, and case 015's v1.x population guidance documents the recipe a v1.x implementer follows to produce the byte-level fixtures. v1.0 implementers can wire the dispatch today against the documented expected outcomes; v1.x implementers will populate the byte-level fixtures from the recipe. The deferral does not block v1.0.

## Question 5 — The empty-day Merkle root convention in design 03 §3.1 — does the byte value match the wire fixture?

**What design 03 §3.1 says.**

> A tenant-day with zero captured events still gets a seal. The Merkle root for an empty leaf set is defined as `SHA-256("")` per RFC 6962 §2.1:
>
> ```
> empty_root = SHA-256("")
>            = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
> ```

**What I checked.** I independently computed `SHA-256(b'')` in Python: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`. Matches the doc.

**What I noticed in passing.** The empty-day root is documented but **not** in `chain_vectors.json`. The conformance corpus carries `single_chain` and `rotation_chain` but no empty-day fixture. Per design 08 §5.1 and the future-cases list in `spec/test-vectors/README.md`, case `004-empty-day/` is planned but not shipped. The empty-day convention is the kind of thing that reads correctly in the design doc and is then implemented inconsistently across vendors (one vendor uses `SHA-256("")`, another uses 32 zero bytes, another refuses to seal an empty day at all). The byte value sits in design 03 §3.1 but does not have a corresponding entry in `chain_vectors.json` for cross-implementation pinning.

**This is an observation, not a v1.0-blocking gap.** The empty-day case is not in the load-bearing path (any tenant whose ledger is producing chain entries has at least one event per day), but it is the kind of edge case that produces inconsistent implementations precisely because it is rare in production. Adding the empty-day root to `chain_vectors.json` (one line: `"merkle_root_empty_hex": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"`) would lock the convention without waiting for the full case 004 fixture.

**Status: Partial.** The convention is documented and the byte value reproduces. The fixture-level pin in `chain_vectors.json` is missing. The fix is one line.

## Per-role roll-up

| Item | Status | Closes when |
|---|---|---|
| Q1: Variant B per-algorithm `sign_payload` declared end-to-end | Answered | Already done — declared in spec §4.2, design 03 §3.7, design 04 §3.2, design 07 §4.3, case 015 + N019 |
| Q2: case (e) severity language and IR-attribution separation | Answered | Already done — bolded "Severe regardless of bracket" sentence + finding-language cross-reference + verifier-vs-IR-program separation |
| Q3: `algorithm/key-type mismatch` reason string exercised by negative corpus | Gap | A `N020-algorithm-key-type-mismatch/` case lands using a single-algorithm tampering recipe (does not require v1.x PQ keys) |
| Q4: fixture-deferral posture for 015 / N017 / N018 / N019 | Answered | Already done — Status block + case 015 v1.x population guidance |
| Q5: empty-day Merkle root convention pinned in `chain_vectors.json` | Partial | One line added to `expected.merkle_root_empty_hex` in chain_vectors.json |

**Stopping criterion.** 0 gaps + 0 partials. Currently at 1 Gap + 1 Partial. Both are on reason-string / fixture precision rather than on the cryptographic substance — Q3 is a precision-of-reporting test the corpus is missing, Q5 is a one-line addition to pin a documented convention.

The substance: HMAC-SHA-256 + HKDF construction with per-tenant `info` binding, fixed-width 32-byte `prev_hash`, MAC IS payload_hash, fingerprint-before-MAC ordering, expected_prev_hash (not entry.prev_hash) into MAC recompute, RFC 6962 Merkle with leaf/node domain separation, Ed25519 over the seven-line `sign_payload` with algorithm + format_version + hkdf_inputs_digest binding, dual-algorithm Variant B per-algorithm bound payload, twelve-step verifier procedure with named failure modes, mid-write truncation refusal, software-key compile-time exclusion, constant-time comparison for fingerprint and MAC, IKM minimum 32 bytes per RFC 4868, IKM registry retention coupled to chain-entry retention, dual-algorithm dispatch with case (a)/(b)/(c)/(d)/(e) handling and per-algorithm validation results recorded for the working paper. All of this reads as cryptographically sound and internally consistent. The single-algorithm verifier is in shipping shape under ENISA-equivalent review. The dual-algorithm verifier is in shipping shape pending the v1.x byte-level fixtures whose absence is honestly documented.

The byte-level reproducibility I verified independently — `hkdf_inputs_digest`, `key_fingerprint_v1`, `merkle_root_single` all reproduce from the spec text alone — is the property that lets me say "I would certify this." Cross-implementation contracts that hold against an independent re-derivation are the only kind that hold across vendor boundaries; the corpus delivers that for the load-bearing single-algorithm path.
