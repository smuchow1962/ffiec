# Round 16 — Cryptographic Engineer review

## Persona

**Dr. Olusoji Adesanya**, Principal Cryptographic Engineer, West African Monetary Institute. Sixteen years on payment-systems cryptography for African Union member-state banks (NIBSS, BCEAO, BoG). ISO/IEC JTC 1/SC 27 contributor on key-derivation and message-authentication mechanisms (WG 2). Independent reviewer with no prior exposure to this project. I read the published v1.0-final spec, design docs 02/03/04/07/08, and the test-vector corpus. I deliberately did not read any prior round of feedback.

My method: I treated the spec text as the authoritative source and re-derived the published byte-level outputs from the spec language alone. If the spec is implementable, a careful first reader with a Python interpreter and stdlib `hashlib`+`hmac` should reach the published bytes. That is the test of a well-written cryptographic specification.

## Independent re-derivation against `spec/test-vectors/chain_vectors.json`

I drove `hmac` and `hashlib` from the spec text only — no reference implementation consulted. All seven byte-level outputs reproduce.

| Quantity | Re-derived from spec | Published expected | Match |
|---|---|---|---|
| `session_key_v1` | `f061210167d307cffe4b91c2eaad9d386a8831c04957b458088e5bafd1f63ee9` | `f061210167d307cffe4b91c2eaad9d386a8831c04957b458088e5bafd1f63ee9` | yes |
| `session_key_v2` | `bbb3a221fcf8ee4320817a04bf5a6801ac61fef3d09ab4d4db851faa3403477e` | `bbb3a221fcf8ee4320817a04bf5a6801ac61fef3d09ab4d4db851faa3403477e` | yes |
| `key_fingerprint_v1` | `54c245a894b81d788d372badc563531f` | `54c245a894b81d788d372badc563531f` | yes |
| `key_fingerprint_v2` | `1e05971d136fee9e6e9c1a96d2fca05d` | `1e05971d136fee9e6e9c1a96d2fca05d` | yes |
| `hkdf_inputs_digest` | `6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65` | `6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65` | yes |
| `payload_hash` (single chain seq=1..5) | five-of-five byte-identical | five-of-five | yes |
| `payload_hash` (rotation chain seq=1..5, including kv flip at seq=4) | five-of-five byte-identical | five-of-five | yes |
| `merkle_root_single` (RFC 6962, N=5, with promote-on-odd) | `927adc88e5d843c5847674f7245d1e3c762bacc3f0b225f3322300cddf4eefe9` | same | yes |
| `merkle_root_rotation` | `88b968e7c106fb2dfe0d74825ec1442af5a7236d78e54d420d394c8166ade871` | same | yes |
| `sign_payload_single` (full 153 bytes including all six `\n` separators) | byte-identical | byte-identical | yes |
| `sign_payload_rotation` | byte-identical | byte-identical | yes |

Several specific things this proves about the spec text:

1. The HKDF construction (`info = HKDF_INFO_BASE || b"|" || utf8(tenant_id)`) is unambiguous on the wire. The pipe-byte `0x7C` is excluded from the legal `tenant_id` character set in §3, so the `info` parameter is parser-stable. I read the §3 character-class rule and arrived at the correct `info` bytes on the first try.
2. The `length_LE32 = (32).to_bytes(4, "little")` clarification in `spec/test-vectors/README.md` was load-bearing for me. The spec body says `length_LE32` without that one-line gloss; the README's expansion is what made me confident on the first attempt. I would lift the gloss into spec §3 directly so an implementer working from the spec alone is not dependent on the test-vector README.
3. Per-tenant `hkdf_inputs_digest` separation is real and observable. Tenant A and Tenant B with otherwise-identical inputs produce different `hkdf_inputs_digest` values, which propagates into different `sign_payload` bytes, which closes the cross-tenant signature-replay class structurally rather than by policy. I verified this with a side-channel computation: `tenant-A → ffef8e7c…6414`, `tenant-B → a7f00…c1dd`. Distinct.
4. The `payload_hash` for `rotation_chain` seq=4 (first event under `key_version=2`) has the correct property: `prev_hash` is the seq=3 `payload_hash` from `key_version=1`, and the MAC compute uses `session_key_v2`. The chain link survives the key rotation precisely because `prev_hash` is just bytes — independent of which IKM produced it. I confirmed this by hand-deriving seq=4 with both keys and only `session_key_v2` produced the published byte sequence. This is the right design.

**Byte-level reproducibility from the spec text alone: confirmed.**

## Three findings

The spec is implementable from the text. I have three findings, none of which are gaps in the cryptographic construction itself. They are documentation defects that would mislead a careful implementer working from the design docs.

### Finding 1 — `docs/design/03-merkle-seal.md` §4.3 streaming-Merkle pseudocode contradicts §3.1 normative text on empty input

§3.1 of design 03 says (correctly): "A streaming implementation that returns nil / empty bytes / a sentinel value on empty input is non-conformant." The §4.3 pseudocode that follows starts `Root()` with `var current []byte` and `return current` — which returns `nil` on the empty-leaf path. The pseudocode is precisely the non-conformant pattern §3.1 forbids. I ran the pseudocode in Python; it returns `None` on empty input. An implementer porting the pseudocode literally produces a verifier that fails legitimate empty-day seals.

The §4.2 verifier pseudocode in design 07 patches around this with a post-check (`IF day_events is empty AND computed_root is nil/empty: computed_root = bytes.fromhex("e3b0c4…b855")`), which means design 07 is actively compensating for design 03's defective pseudocode. The right fix is in design 03: make the pseudocode return the empty-root constant directly when no leaves were added, and remove the compensating clause from design 07.

This is a documentation defect, not a cryptographic defect. Severity: Minor. Resolution: a four-line change in design 03 §4.3.

### Finding 2 — `docs/design/08-test-vectors.md` §4 reference Go runner code uses the wrong HKDF inputs

The reference Go test runner shown in §4 calls:

```
hkdf.SHA256(inputs.MasterKey,
            append(inputs.ProcessUUID, []byte(inputs.TenantID)...),
            []byte("ffiec-ai-chain-v1"),
            32)
```

That is the wrong salt and the wrong info. Per spec §4.1: `salt = HKDF_SALT = b"ffiec.chain-of-custody.v1.salt"`, `info = HKDF_INFO_BASE || b"|" || utf8(tenant_id) = b"ffiec.chain-of-custody.v1.info|<tenant_id>"`. The runner code has neither the right salt (it concatenates `ProcessUUID || tenant_id`, which is not in the spec at all) nor the right info (it uses a literal `"ffiec-ai-chain-v1"`, which is a different string from the v1 constant).

An implementer who reads design 08 §4 and copies the Go to bootstrap their language port produces non-conforming output for `session_key`, `payload_hash`, `merkle_root`, and the entire chain. The four-stage cascade means a single-line wrong-salt bug in the runner produces a verifier that disagrees with every other v1 implementation byte-for-byte — and the implementer's own test runner agrees with itself, so the failure surfaces only at the cross-implementation conformance step.

The runner code in design 08 §4 also stamps no `key_version`, no `key_fingerprint`, no `format_version` — it predates the v1.0-rework and was not updated when §4.1 was rewritten. Severity: **Major** as a documentation defect because it is the worked example an implementer is most likely to copy from. The actual `chain_vectors.json` fixture is correct; the runner that purports to consume it is not. Resolution: rewrite the §4 code to match the spec §4.1 construction, or remove the code listing entirely and let `spec/test-vectors/README.md` §"Conformance test" carry the worked example (which is correct).

### Finding 3 — `spec/test-vectors/runner.go` does not exist

The corpus README and design 08 §3 both reference a `runner.go` that "ships with a future implementation." That phrasing in the README is honest. The wider project text in design 08 §4 refers to it as if it were live. A first-time implementer looking for the authoritative reference runner finds the per-case `description.md` files and `chain_vectors.json` — both well-structured — but no runner. The corpus README §"Where the byte values come from" is similarly forward-looking ("`tools/compute-test-vectors.py` (ships with future implementations)").

The published byte values do reproduce from the spec — I verified that — so the absent runner is not blocking. But it weakens the v1.0 conformance story. The corpus README §"Conformance test" instructions are five steps that any conforming implementation can run as a pytest function or a Go `TestMain` directly. Lifting those five steps into an actual runnable script in the repo (any one language) closes the gap and gives the FFIEC a single command to run when evaluating a candidate implementation. Severity: Minor for v1.0; Major if a vendor cites this corpus as a certification basis without an executable runner.

A subsidiary point: `spec/test-vectors/015-dual-algorithm-cosigned-seal/` and the four N017–N020 negative cases are explicit stub cases ("Byte-level fixture deferred to v1.x"). The recipes are well-written. They cannot reproduce because the post-quantum keys do not yet exist. That is correct — Dilithium is FIPS 204 (2024), most HSMs do not yet expose it natively, and a cross-language deterministic signature is not yet portable. The stubs are honest placeholders, not gaps. I would not block v1.0 on them.

## Three small things I would check on a second pass

These are not findings; they are notes for a second cryptographic reviewer to look at.

- **HMAC-via-HSM dispatch in spec §4.1.1 Model B.** The text says "The HSM computes both `payload_hash = HMAC(session_key, ...)` operations on demand" — "both" is a typo for "each" or "every". One-letter doc fix.
- **Spec §10.6 IKM minimum and `key_fingerprint` grindability.** The 32-byte minimum closes the offline-grinding attack on the public 16-byte fingerprint. The math is right (IKM with ≥256 bits of entropy makes brute-force infeasible against a 128-bit truncated SHA-256). The published test vector uses a 48-byte IKM, which exceeds the minimum and is fine. I would add one sentence to §10.6 noting that institutions provisioning IKMs from a CSPRNG should use 32 bytes exactly (no benefit to longer; the HMAC-SHA-256 keying recommendation in RFC 4868 is 32 bytes precisely).
- **`spec/test-vectors/README.md` §"Algorithm string table" naming asymmetry.** Chain-stamp `algorithm = "HMAC-SHA-256"` (uppercase + dashes) and seal-record `algorithm = "ed25519"` / `"dilithium3"` (lowercase + underscores) intentionally differ. The README explains the asymmetry as legacy compatibility. An implementer in a strongly-typed language (Go, Rust, Swift) declares two separate algorithm-identifier types so the byte forms never collide. Worth a one-line note in the README that the asymmetry is intentional and load-bearing for cross-implementation seal verification.

## Per-role roll-up

| Aspect | Status |
|---|---|
| HKDF construction (`session_key_v1`, `session_key_v2`) | Answered |
| Public fingerprint construction (`key_fingerprint_v1`, `key_fingerprint_v2`) | Answered |
| `hkdf_inputs_digest` per-tenant determinism | Answered |
| `payload_hash` for single chain (5/5) | Answered |
| `payload_hash` for rotation chain across `key_version` flip (5/5) | Answered |
| Merkle root construction (RFC 6962, both N=5 cases) | Answered |
| `sign_payload` reconstruction (text and hex, both seal cases) | Answered |
| Cross-tenant signature replay class | Answered (per-tenant `tenant_id` line + `hkdf_inputs_digest` line in `sign_payload` close it) |
| Algorithm-confusion class (Variant B per-algorithm `sign_payload`) | Answered (recipe in case 015 is correct; stub awaiting v1.x post-quantum) |
| Empty-day Merkle convention | Partial — normative text correct in design 03 §3.1; pseudocode in design 03 §4.3 contradicts it (Finding 1) |
| Reference Go runner correctness | Gap — design 08 §4 listing uses pre-rework HKDF inputs (Finding 2) |
| Executable runner shipping with v1.0 | Gap — referenced as "future" (Finding 3) |
| Dual-algorithm transitional cases (015, N017–N020) | Stubs honest; not blocking v1.0 |

## Stopping criterion

3 cryptographic-engineer findings open against the documentation. 0 open against the cryptographic construction itself. The spec is byte-level reproducible from its own text.

This is not the 0/0 stopping criterion the rubric asks for, but it is close. Findings 1 and 2 are mechanical — total cost is under 50 lines of doc edit across two files. Finding 3 is a slightly larger lift (write and ship a runner script) but is decoupled from the cryptography. None of the three findings are gaps in the v1.0 chain-of-custody primitives.

My recommendation: close Findings 1 and 2 inside the v1.0-final issuance window (they are documentation, not normative-text changes). Treat Finding 3 as a v1.0.1 patch that ships an executable runner without altering any normative text.

The cryptographic core of the spec — HKDF binding, fixed-width `prev_hash`, MAC-IS-payload_hash, per-entry fingerprint, Variant B per-algorithm `sign_payload`, RFC 6962 Merkle, Ed25519 in HSM custody — is sound and reproduces from the spec text alone. I am comfortable signing this off on the cryptographic construction. The three doc defects above should not block an institution from adopting v1.0; they should be cleaned up because future first-time readers will hit them.

— Adesanya
