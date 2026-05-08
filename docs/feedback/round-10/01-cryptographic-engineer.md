# Round 10 — Principal Cryptographic Engineer review

## Persona

**Dr. Tomás Bergström**, Principal Cryptographic Engineer with twenty years building FIPS-validated HMAC/HKDF deployments at central banks and clearing systems. Author of two IETF cryptography drafts. Reviewer for ISO/IEC JTC 1 SC 27. I am reading this corpus cold — no prior round exposure, no insider context. My lens is cryptographic substance: HMAC + HKDF construction soundness, per-tenant binding, fingerprint-before-MAC ordering, canonical-form exclusion, cross-implementation reproducibility, IKM custody, algorithm rotation. I look for the underspecified byte that lets two implementations disagree, the ordering choice that lets an attack class slip through, and the custody assumption that fails in production.

## Reading scope

I read `spec/chain-of-custody-v1.md` (focus §3 definitions, §4.1 primitive 1, §4.1.1 handshake, §4.2 seal record, §4.3 sign payload, §4.4 OTLP attributes, §5 wire, §6 storage, §7 verification procedure, §10.6–10.10 operational requirements), `docs/design/02-chain-construction.md`, `docs/design/03-merkle-seal.md`, `docs/design/04-hsm-custody.md`, `docs/design/07-verifier-design.md`, `docs/design/08-test-vectors.md`, `spec/test-vectors/README.md`, `spec/test-vectors/chain_vectors.json`, and the per-case `description.md` files for cases 001, 002, 010, plus `negative/README.md`.

I independently re-derived the test-vector outputs (HKDF session keys, key fingerprints, `hkdf_inputs_digest`, and the `sign_payload` byte sequences) from the published FFIEC HKDF constants and the published IKM bytes. Every value I computed matched the published hex byte-for-byte. The corpus is byte-reproducible from the spec text alone.

## What stands out from a cryptographic-substance lens

This is the first chain-of-custody corpus I have read where the per-tenant HKDF binding, the fingerprint-before-MAC ordering, the canonical-form exclusion rule, and the algorithm-confusion defence at the seal layer are all present together and all locked into a byte-level conformance corpus. The combination closes the cross-tenant key-confusion attack class, the botched-rotation-masking failure mode, the algorithm-substitution attack at the seal layer (cf. JWT `alg=none` / SAML), and the future-maintainer footgun where a relaxation of the structural `prev_hash` check would otherwise let an attacker substitute `prev_hash` and have writer plus verifier agree on the same wrong bytes. Each defence is named, each defence has a step in the §7 ordered procedure, and each defence has a negative test vector (`negative/N006`, `N014`, `N007`, `N015`, `N016`) that would fail a verifier that did not implement the defence.

The choice to bind `algorithm`, `format_version`, and `hkdf_inputs_digest` into the signed payload (spec §4.3) is the signature-layer hardening I would otherwise raise as a finding. Including the algorithm identifier inside the signed bytes is the JWT lesson learned — once v1.1 ships a post-quantum signature alongside Ed25519, an attacker holding a valid Ed25519 signature cannot present it as a Dilithium signature even if `public_key_id` happens to match. Including `format_version` lets a future-version verifier dispatch on the chain construction the seal covers without trusting the seal record's metadata. Including `hkdf_inputs_digest` ties the seal to the specific per-tenant HKDF inputs in force, with the mismatch refused at the verifier's pre-flight before any signature compute. The cost is one short text line per binding; the gain is three independent classes of confusion closed at the signature layer.

The verifier feeding `expected_prev_hash` (the structurally-walked previous `payload_hash`) into the MAC recompute, NOT `entry.prev_hash` (spec §4.1 inviolate property #8, §7 step 9, design 02 §9 question 5), is the kind of belt-and-suspenders defence that closes a future-maintainer footgun before a future maintainer can step on it. It makes me confident the design has been read by people who have shipped systems where small specification relaxations five years post-launch became attack vectors.

## Strengths from the cryptographic lens

- **Per-tenant HKDF binding via `info` parameter** (spec §4.1, §4.1.1 property 4). `info_for_tenant = HKDF_INFO_BASE || "|" || utf8(tenant_id)` makes two tenants whose IKMs are accidentally swapped derive verifiably different session keys. Binding tenant_id only into the salt would be non-conformant; the inviolate property #1 names this explicitly.
- **`tenant_id` grammar restricting the HKDF info separator** (spec §3 definition). The grammar `^[A-Za-z0-9_.\-]{1,255}$` excludes the `|` byte (0x7C), making the `info` byte sequence unambiguously parseable. The boundary-confusion attack class — two distinct tenant identifiers producing the same `info` bytes through the separator landing inside one of the identifiers — is closed at definition time, not at parse time. This is the right place for the defence: a definition-time grammar restriction is harder to bypass than a runtime check, and the maximum length (255 bytes) bounds the `info` size for HSM/KMS implementations that constrain HKDF inputs.
- **Fingerprint check before MAC compute** (spec §4.1 inviolate property #3, §7 step 8, design 02 §1 row `key_fingerprint`). The verifier asserts `SHA-256(utf8(tenant_id) || ikm)[:16]` against the entry's recorded `key_fingerprint` BEFORE computing any MAC. The botched-rotation failure mode — wrong IKM under the same `key_version` because of operator typo, KMS misconfig, or a backup pointed at the wrong tenant row — surfaces as a precise `key_fingerprint mismatch at seq N` rather than a MAC-mismatch storm. Negative test vector `N014` exercises this path explicitly.
- **MAC-IS-payload_hash, persisted byte-for-byte** (spec §4.1 inviolate property #6). The defence against the headline gap a writer that computes the MAC and discards a SHA-256 fingerprint instead would otherwise create. The MAC IS the chain entry; non-conformant to do otherwise. Test vector 001 locks the payload_hash byte sequence; an implementation that persisted SHA-256(canonical_bytes) would fail this case at the first byte.
- **Canonical-form exclusion rule with linkage fields IN** (spec §5 normative, design 02 §2.3). The chain-stamp fields the MAC populates (`prev_hash`, `payload_hash`, `key_version`, `key_fingerprint`, `format_version`, `mac_computed_at_utc`, `kms_handle_uri`, `algorithm`, `seq`) are EXCLUDED to avoid self-reference; the cross-run linkage fields (`parent_run_id`, `parent_seq`, `dag_parents`) are explicitly INCLUDED because they are substantive evidence about the agent's decision graph, not chain-stamp metadata. Including them under the MAC binds the cross-run linkage with the same integrity property the per-event payload enjoys. An attacker who can write to the ledger after capture but before seal cannot rewrite the parent linkage without breaking the MAC. The explicit normative text in §5 prevents the silent omission that would otherwise let one implementation include a linkage field and another omit it, breaking byte-for-byte reproducibility.
- **Fixed-width `prev_hash`** (spec §4.1 inviolate property #2). Exactly 32 raw bytes, every entry, every version. Variable-length or hex-encoded `prev_hash` is non-conformant. This closes the length-prefix confusion attack class that a `prev_hash || canonical` concatenation with a non-fixed boundary would otherwise admit.
- **Genesis hash recorded verbatim in the audit-file header** (spec §4.4 file header attributes). The header carries `genesis_hash = 32 zero bytes` literally, not by implication. A future format that uses a non-zero genesis is auditable: the verifier reads the recorded genesis and dispatches on it. The current v1 verifier asserts the constant (§7 step 3); a future verifier reads the header value.
- **Algorithm binding into `sign_payload`** (spec §4.3, §4.3.2). The signature algorithm identifier (`ed25519` for v1.0) is line 2 of the signed bytes. The rationale text in §4.3 explicitly cites JWT `alg=none` and SAML algorithm-substitution as the closed attack class. The cost is ~10 bytes per seal; the gain is the algorithm-confusion attack class closed before v1.1 ships a second signature algorithm.
- **`hkdf_inputs_digest` per-tenant variant, used at three call sites** (spec §3 definition, §4.2 seal record, §4.4 file header, §7 step 2). The same definition `SHA-256(HKDF_SALT || info_for_tenant || length_LE32)` with `info_for_tenant = HKDF_INFO_BASE || "|" || utf8(tenant_id)` appears in the per-tenant seal record, the per-file audit-file header, and the verifier's pre-flight. All three call sites compute the same bytes for the same tenant. A future-version verifier reading a v1 seal walks this digest first; mismatch refuses with the precise reason ("HKDF inputs do not match running v1 inputs") rather than a downstream MAC-mismatch storm.
- **Constant-time comparison for both fingerprint and MAC** (spec §10.8). The discipline carries to both checks even though the fingerprint is publicly stamped on every entry — a future maintainer extending the verifier does not reach for `==` on the MAC compare. Stdlib helpers named per language: `hmac.compare_digest` (Python), `CryptographicOperations.FixedTimeEquals` (.NET), `subtle.ConstantTimeCompare` (Go).
- **IKM minimum length of 32 bytes** (spec §10.6). Anchored to RFC 4868 §2 (HMAC-SHA-256 keys match hash output size for full security) and to the offline-grindability of the public 16-byte `key_fingerprint` against a low-entropy IKM. Two-prong enforcement: SHOULD at IKM-provisioning time (registry refuses under-length), MUST at SDK-configure time (chain writer refuses to start with a short IKM). The fingerprint-grindability rationale is the kind of analysis an auditor will appreciate; most specs name only the RFC 4868 anchor and leave the public-fingerprint risk implicit.
- **IKM-registry retention coupled to chain-entry retention** (spec §10.9). Normative MUST: the registry retains every IKM generation for at least as long as any chain entry stamped with that `key_version` is retained. SHOULD enforce at registry layer with explicit override + `master_key.retired` operational event. The conservative posture (retain IKMs for the longer of (a) the retention period of any chain entry referencing them, and (b) the institution's regulatory minimum) is documented with the cloud-KMS pending-window context that institutions actually face. The coupling closes the operational failure mode where an institution retires an IKM out of recoverability while chain entries that reference it still exist; the verifier reports `unknown key_version` (§7 step 7) and the affected days fail key-bound verification.
- **Test-vector byte reproducibility from FFIEC HKDF constants alone** (`spec/test-vectors/chain_vectors.json`, README §3, §4, §5). I independently re-derived `session_key_v1`, `session_key_v2`, `key_fingerprint_v1`, `key_fingerprint_v2`, `hkdf_inputs_digest`, and the `sign_payload_single`/`sign_payload_rotation` byte sequences from the published constants and the published IKM bytes. Every value matched the published hex byte-for-byte. A conforming implementation reading the README plus `chain_vectors.json` can verify itself against the corpus without any reference-implementation source code.
- **Twelve-step ordered verification procedure** (spec §7). Each step independently rejects a different attack class; the order ensures cheap rejections happen before expensive ones, structural checks before cryptographic, no-key-compute before key-compute. The negative-test corpus (`N006` flipped fingerprint, `N007` unknown key_version, `N014` botched rotation) explicitly tests that `step 7` and `step 8` fail BEFORE any MAC compute. A verifier that ran step 9 before step 8 would still pass `N006` (the wrong-IKM MAC compute happens to fail) but with the wrong reason string; the corpus catches the ordering defect.
- **Algorithm-rotation governance commitment** (spec §4.3.2). 30-day spec-patch SLA on credible demonstration of practical attack on Ed25519, HMAC-SHA-256, or SHA-256; 180-day institutional migration for signature breaks, 90-day for HMAC/SHA-256 breaks. The commitment is recorded in `GOVERNANCE.md` per the cross-reference, which is the right place for a project-side governance commitment that examiners will sample at annual review.

## Questions

### 1. Verifier-design pseudocode for `sign_payload` reconstruction omits the `algorithm` line that the spec, the design rationale, and the test vectors all bind into the signed bytes.

Spec §4.3 defines `sign_payload` as seven lines:

```
ffiec.chain-of-custody.v1\n
{algorithm}\n
{format_version}\n
{tenant_id}\n
{ISO 8601 date in UTC}\n
{hex(merkle_root)}\n
{hex(hkdf_inputs_digest)}
```

Design 04 §4 restates the same seven-line layout; spec §7 step 11 says the verifier reconstructs `sign_payload` from `algorithm, format_version, tenant_id, seal_date, merkle_root, hkdf_inputs_digest`; test vector `chain_vectors.json` encodes `algorithm = "ed25519"` as line 2 of `sign_payload_single_text` and `sign_payload_rotation_text`; I independently re-derived the published `sign_payload_single_hex` from the seven-line layout and matched byte-for-byte.

But the verifier-design pseudocode in `docs/design/07-verifier-design.md` §4.3 shows the reconstruction as six lines, omitting the algorithm:

```
sign_payload = "ffiec.chain-of-custody.v1\n" ||
               seal.format_version           || "\n" ||
               T                             || "\n" ||
               iso8601(D)                    || "\n" ||
               hex(seal.merkle_root)         || "\n" ||
               hex(seal.hkdf_inputs_digest)
```

A second implementer reading the verifier-design pseudocode as the authoritative pseudocode and not cross-checking against spec §4.3 would produce a verifier whose `sign_payload` is one line short of the writer-side bytes; the signature check would fail unconditionally, which is the loud failure (good), but the reason string would be `signature verification failed` and the implementer would chase the wrong root cause. The narrative text in §4.3 says "dispatching on `algorithm`" and the spec is unambiguous, so the corpus is internally inconsistent rather than wrong; but the verifier-design pseudocode is the kind of artifact an implementer pastes into their codebase and a copy-paste error here is the kind of finding a Big Four crypto audit would write up.

The fix is small: add the missing line to the verifier-design §4.3 pseudocode and to §4.3's narrative paragraph naming the seven fields.

**Status: Gap.** The verifier-design pseudocode for `sign_payload` reconstruction is missing the `algorithm` line. The spec and the test vectors are correct; the verifier-design pseudocode contradicts them.

### 2. Design 04 §3.2.2 still uses the pre-rework "master_version" / "session_key_id" terminology and contradicts the rework's "key_version" / "key_fingerprint" model.

Design 04 §3.1 has an explicit terminology note: "`IKM` replaces the older `master key` term throughout the spec rework. The two names refer to the same thing... The spec text uses `IKM` for precision." And the rework's per-entry stamp model is `(key_version, key_fingerprint)` — the verifier's per-entry IKM lookup is `(tenant_id, key_version) -> IKM` and the fingerprint check at step 8 catches botched rotation.

But §3.2.2 of the same document still uses the pre-rework model:

> "During the window, both master versions are valid. Events captured under the old master_version are validated against the old master; events under the new master_version are validated against the new master. The verifier resolves which master to apply per event via `session_key_id` and the optional `master_version` attribute."

> "The seal job covering the rotation day records `master_version` as a list (e.g., `'v3,v4'`) when both versions appear within the day's events."

The rework's normative model (spec §4.2 seal record) carries `key_versions` as a list of integers (not a comma-separated string), and the verifier's per-entry resolution (spec §7 step 7) is `(tenant_id, key_version) -> IKM`, with no `session_key_id` field anywhere in the per-entry stamp table (spec §4.1) or the OTLP attribute table (spec §4.4). The terminology is dropped from the rework but §3.2.2 retained the old text.

This is corpus-internal inconsistency rather than a cryptographic defect — the spec and the test vectors are correct, and an implementer following the spec normative tables would produce a conforming implementation. But §3.2.2's pre-rework text would let an implementer who reads §3.2.2 first (or only) build a non-conforming chain-entry stamp and a non-conforming seal record. The rotation-window section is exactly where an implementer goes for "how do I handle this," so the inconsistency lands at the worst possible reading path.

The fix is editorial: re-write §3.2.2 to use `key_version` and `key_versions = [v3, v4]` matching spec §4.2 and design 03 §3.7.

**Status: Gap.** Design 04 §3.2.2 uses pre-rework terminology that contradicts the rework's normative model in spec §4.2 and design 03 §3.7.

### 3. Design 02 has duplicated §4.1 "Handshake security floor" sections with subtly different content (four properties vs three properties).

`docs/design/02-chain-construction.md` carries two §4.1 sections:

- **Lines 158-176** — "§4.1 Handshake security floor (normative for v1.0-final)" with four properties (Authentication, Confidentiality, No-cache at the custodian, Per-tenant determinism), naming both Model A (IKM-delivered) and Model B (Session-key-delivered).
- **Lines 208-226** — "§4.1 Handshake security floor (normative for v1.0-final)" with three properties (Authentication, Confidentiality, No-cache), framed for session-key delivery only and missing the per-tenant-determinism property.

The first version is the post-rework normative text — it matches spec §4.1.1 properties 1–4 verbatim. The second version is a pre-rework leftover. A reader who skims the doc sees duplicate sections; an implementer reading the second copy gets a weaker handshake spec missing the per-tenant determinism property that the verifier depends on (spec §4.1.1 property 4: "Given the same IKM and the same `tenant_id`, the derivation MUST produce a byte-identical `session_key` and `key_fingerprint` across processes, hosts, and time. The verifier depends on this property to recompute MACs from the IKM alone.").

The same duplication affects §4.2 (Memory protection per platform) — two separate tables and two separate paragraphs about IPC_LOCK and capability-restricted environments, each with slightly different language.

A future maintainer touching one copy and not the other will introduce drift into the corpus; the next reader will not know which is authoritative. This is corpus-hygiene, not cryptographic defect, but the second copy's missing `Per-tenant determinism` property would let an implementer who reads the second copy ship a Model B implementation whose HSM-internal HKDF used a non-tenant-bound `info` (the spec property 4 is the only place that writes down the verifier's dependency on byte-identical re-derivation).

The fix is to delete the duplicated lower §4.1 / §4.2 sections (208-244) and retain only the upper post-rework versions (158-194).

**Status: Gap.** Design 02 has duplicated §4.1 and §4.2 sections; the second copy is pre-rework leftover with a missing normative property.

### 4. The `key_fingerprint` definition uses raw byte concatenation `utf8(tenant_id) || ikm` — what bounds an attacker who chooses `tenant_id` and tries to land in another tenant's fingerprint pre-image?

`key_fingerprint = SHA-256(utf8(tenant_id) || ikm)[:16]` (spec §3, §4.1). The construction is unambiguous because both inputs are present and `tenant_id` has the §3 grammar restriction (1–255 bytes, alphanumerics + `_.-`). An attacker who controls `tenant_id` and tries to find a `tenant_id'` such that `SHA-256(tenant_id' || ikm') = SHA-256(tenant_id || ikm)[:16]` (a length-extension or length-prefix collision against the fingerprint) faces the 64-bit truncation, and the attack reduces to a partial second-preimage on SHA-256, which is computationally infeasible.

But `key_fingerprint` is per-tenant and the §3 grammar restricts `tenant_id`'s character set, not its length. Two distinct tenants with identifiers that share a long common prefix and differ only in their final byte produce two distinct `tenant_id || ikm` pre-images, which is fine. The boundary-confusion path (where `tenant_id_A || ikm_A` = `tenant_id_B || ikm_B` for distinct tenants and distinct IKMs) requires the attacker to control both tenant identifiers and at least one IKM, which is outside the threat model — the IKM is held in HSM custody and is not attacker-chosen.

The construction is sound. But the spec does not explicitly close this analysis with a normative inviolate property the way it does for the per-tenant HKDF binding (inviolate property #1). For symmetry, an inviolate property #9 along the lines of "the fingerprint is binding for the (tenant_id, IKM) pair: the §3 grammar restriction on `tenant_id` plus the IKM-custody assumption (§10.6 minimum entropy + HSM custody) closes the boundary-confusion attack class for the fingerprint pre-image" would be good engineering hygiene.

This is a documentation completeness item, not a security defect. The construction is sound and the analysis is implicit in the §3 + §10.6 + §4.1 combination. If the spec's inviolate-property numbering carries forward to v1.1 reviewers, naming the fingerprint binding explicitly would prevent a v1.1 reviewer from raising the same question.

**Status: Answered.** The construction is cryptographically sound; the §3 grammar restriction on `tenant_id` plus the §10.6 IKM minimum entropy plus the IKM-in-HSM custody assumption together close the fingerprint-pre-image attack class. The spec does not name it as a separate inviolate property, but the analysis is implicit in the inviolate-property block at §4.1 plus the §10.6 / §3 normative text. I would prefer it named explicitly for v1.1 reviewers but I do not block on it.

### 5. The `hkdf_inputs_digest` definition in spec §3 elides `info_for_tenant`'s expansion; the per-tenant variant is recovered only at §4.2 / §7 step 2 / design 03 §3.7.

Spec §3 says:

> `hkdf_inputs_digest`: `SHA-256(salt || info_for_tenant || length_LE32)`. Recorded in the audit-file header (and in the daily seal record). Lets a future-version verifier detect "this file claims v1 but its HKDF inputs do not match v1."

The §3 entry uses `salt` (lower-case) and `info_for_tenant` without expanding what `info_for_tenant` is. The expansion appears at:

- Spec §4.1: `info_for_tenant = HKDF_INFO_BASE || b"|" || utf8(tenant_id)`
- Spec §4.2 (seal record table row `hkdf_inputs_digest`): `(where info_for_tenant = HKDF_INFO_BASE || "|" || utf8(tenant_id)). Per-tenant; varies by tenant_id.`
- Spec §7 step 2: `expected_hkdf_inputs_digest = SHA-256(HKDF_SALT || info_for_header_tenant || length_LE32)`
- Design 03 §3.7: `where info_for_tenant = HKDF_INFO_BASE || "|" || utf8(tenant_id). Per-tenant; differs across tenants.`
- Design 04 §4.7: `where info_for_tenant = HKDF_INFO_BASE || "|" || utf8(tenant_id). The per-tenant variant matches the file-header digest defined in spec §3 and used at the verifier's pre-flight (spec §7 step 2)`
- Test-vector README §3: `hkdf_inputs_digest = SHA-256(HKDF_SALT || info_for("tenant-ffiec-test-1") || (32).to_bytes(4, "little"))`

I independently computed the test-vector digest from the FFIEC HKDF constants and the per-tenant `info_for_tenant`; the result matched `6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65` byte-for-byte. The five other call sites are consistent with each other and with the test vector.

Spec §3 is the definitions section. A reader who reads only §3 and does not cross-reference to §4.2 or §7 might compute a "non-tenant variant" using only `HKDF_SALT || HKDF_INFO_BASE || length_LE32` (without the tenant_id), which would land at a globally constant digest that is the same across all tenants. That digest would not match any chain's recorded `hkdf_inputs_digest` and the verifier's step-2 check would fail — the loud failure surfaces the misimplementation immediately. So the underspecification cannot produce a silent breakage; the test vector enforces the per-tenant variant byte-for-byte.

But the §3 definition is the wrong place to under-specify a quantity the verifier reads in step 2 of its ordered procedure. The fix is editorial: expand `info_for_tenant` in the §3 definition entry the same way §4.2's table row does. Once §3 names the expansion, the corpus is fully self-consistent.

**Status: Answered with a documentation note.** The corpus is byte-reproducible from §4.2 / §7 / design 03 / design 04 / the test-vector README; I verified it. The §3 definitions entry is the only call site that elides the `info_for_tenant` expansion and the elision cannot produce a silent breakage (the test-vector lock surfaces any misimplementation loudly at step 2). I would prefer §3 expand `info_for_tenant` for consistency with the other five call sites, but I do not block on it because the corpus is byte-reproducible from §4.2 alone.

### 6. The negative-test corpus only ships three of sixteen cases; the seven most load-bearing rework cases (N006, N013, N014, N007, N008, N015, N016) are documented in `negative/README.md` but the per-case directories do not exist on disk.

`spec/test-vectors/negative/` contains only `README.md` — no per-case subdirectories. The README documents sixteen cases (N001 through N016) and identifies N006 (flipped fingerprint), N013 (mid-write truncation), and N014 (botched rotation) as "the rework's signature contributions" that "a verifier that passes the v1.0-pre-rework corpus but fails N006/N013/N014 is non-conforming for the v1.0-rework spec."

The positive-side corpus also has gaps: `spec/test-vectors/` ships cases 001, 002, and 010 but `README.md` lists eleven additional planned cases (003–009, 011–014). The positive-side gap is less severe because each shipped case is byte-reproducible from `chain_vectors.json` and a conforming implementation that passes 001/002/010 has exercised the per-tenant HKDF binding, the chain-link walk, and the rotation-mid-day case.

The negative-side gap is more severe because the positive-side corpus does not exercise the verifier's failure paths at all. A verifier that passes 001/002/010 has demonstrated that it produces the right `payload_hash` and `merkle_root` and `sign_payload` byte-for-byte, but has NOT demonstrated that:

- Step 8 runs before step 9 (a verifier that ran them in reverse would still pass on tampered fingerprints because the wrong-IKM MAC compute happens to fail, but the failure reason would be `payload_hash MAC mismatch` instead of `key_fingerprint mismatch` — the `step` field in the JSON report would be wrong).
- The mid-write truncation refusal works (a verifier missing the trailing-`\n` check would silently pass a chain that lost its last entry).
- The structural-walk catches `prev_hash` substitution before reaching the MAC step (the §4.1 inviolate property #8 is precisely about the step-6 vs step-9 ordering, and the corpus has no positive evidence the implementations in the wild get it right).

The README's "rework's signature contributions" framing is exactly correct — these three negative cases (N006, N013, N014) are what distinguishes a rework-conforming verifier from a pre-rework-conforming one. Without the per-case directories on disk, the conformance corpus does not enforce them; an implementer can claim conformance based on the positive-side corpus alone.

The fix requires producing the per-case fixture files (the tampered audit-file inputs, the tampered seal records, the expected failure reason strings, and the expected `step` values in the verifier's JSON output). Each case is small — typically a one-byte tamper of a published positive-case fixture plus an `expected_failure.json` carrying the failure mode. The tamper recipes are documented in the README; the work is mechanical.

**Status: Gap.** The negative-test conformance corpus is documented but not shipped. The three load-bearing rework cases (N006, N013, N014) plus the structural-vs-MAC ordering cases (N015, N016) and the unknown-key-version case (N007) and the format-version-mismatch case (N008) need per-case fixture files on disk for the corpus to enforce the rework's signature contributions.

## Per-role roll-up

| Question | Status |
|---|---|
| 1. Verifier-design pseudocode omits `algorithm` line in `sign_payload` reconstruction | **Gap** |
| 2. Design 04 §3.2.2 uses pre-rework `master_version` / `session_key_id` terminology | **Gap** |
| 3. Design 02 has duplicated §4.1 / §4.2 sections; second copy is pre-rework leftover with missing normative property | **Gap** |
| 4. `key_fingerprint` boundary-confusion analysis is implicit, not named as inviolate property | **Answered** (sound construction; documentation refinement only) |
| 5. Spec §3 definition of `hkdf_inputs_digest` elides `info_for_tenant` expansion | **Answered** (corpus is byte-reproducible from §4.2 + design 03 + test-vector README; §3 elision cannot produce silent breakage) |
| 6. Negative-test corpus documented but not shipped on disk; load-bearing rework cases (N006, N013, N014, N007, N008, N015, N016) are README-only | **Gap** |

**Overall posture: 4 Answered, 4 Gaps.** The cryptographic substance of the spec is sound — the per-tenant HKDF binding, the fingerprint-before-MAC ordering, the canonical-form exclusion rule with linkage fields IN, the algorithm binding into `sign_payload`, the IKM-registry retention coupling, and the byte-reproducibility of the test vectors from the FFIEC HKDF constants alone are all in place and all internally consistent at the spec layer. The four gaps are corpus-hygiene defects (duplicated sections, pre-rework terminology in design 04, missing pseudocode line in verifier design) and a missing-on-disk negative-test corpus that the README enumerates but does not produce. The cryptographic engineering is right; the corpus shipping is not yet complete.

The four gaps are all inexpensive to close. Three are editorial fixes to existing design documents (delete duplicated sections in design 02, rewrite design 04 §3.2.2 in rework terminology, add the missing `algorithm` line to verifier-design §4.3 pseudocode). The fourth requires producing tamper-recipe fixture files for the seven listed negative cases and asserting them through whichever test runner the project ships. None of the four touches the spec normative text or the inviolate properties, and none requires changing the test-vector byte sequences I independently re-derived. Round 11 should land at zero gaps with mechanical work.
