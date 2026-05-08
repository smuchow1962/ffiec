# 08 — Conformance test vectors

> **What this doc is.** The design of the conformance test corpus. The corpus is the discriminator between conforming and non-conforming implementations. It is the artifact that proves two implementations are interoperable.

## 1. What the corpus is

A directory `spec/test-vectors/` containing self-contained, deterministic test cases. Each test case has:

- A **description** of what the test exercises
- **Inputs** &mdash; the canonical event payloads, the (test) session key, the (test) tenant master, the (test) Ed25519 keypair
- **Expected outputs** &mdash; the `payload_hash` for each event, the daily Merkle root, the Ed25519 signature
- A **test runner** that any implementation can use to verify it produces the expected outputs

A conforming implementation produces byte-for-byte identical outputs for every test case.

## 2. Why the corpus matters

Without test vectors:

- Two vendors implement the spec, both believe they are correct, and produce different outputs for the same events. The chain breaks at vendor boundaries.
- An auditor cannot point to an authoritative reference. "What does the spec mean?" has multiple answers.
- The FFIEC cannot name a standard whose conformance is enforceable.

With test vectors:

- Conformance is mechanical. Run the tests; they pass or fail.
- Two implementations that pass the corpus produce identical bytes.
- Adding a new test vector is the way to lock in a previously-ambiguous behavior.
- Bugs that produce subtly-wrong output are caught by the test that exercises that path.

## 3. The structure of `spec/test-vectors/`

The corpus is seeded from Herald.Py's `tests/fixtures/chain_vectors.json` shape (the byte-level cross-language contract that Herald.Py and Herald.NET both reproduce verbatim). The FFIEC corpus uses FFIEC-specific HKDF constants (`HKDF_SALT = b"ffiec.chain-of-custody.v1.salt"`, `HKDF_INFO_BASE = b"ffiec.chain-of-custody.v1.info"`).

```
spec/test-vectors/
├── README.md                       # Index of test cases
├── runner.go                       # Reference Go test runner
├── runner_test.go                  # Validates the runner itself
├── chain_vectors.json              # Cross-language byte-level fixture mirror
├── 001-single-event-empty-prev/
│   ├── description.md              # What this case exercises
│   ├── inputs/
│   │   ├── ikm.hex                 # 32 bytes test IKM (clearly marked "TEST USE ONLY")
│   │   ├── tenant_id.txt
│   │   ├── key_version.txt         # Integer (1 for non-rotation cases)
│   │   ├── ed25519_private.pem     # Test signing key
│   │   ├── ed25519_public.pem      # Corresponding public key
│   │   ├── seal_date.txt           # ISO 8601 date
│   │   ├── opened_at_utc.txt       # ISO 8601 UTC for the audit-file header
│   │   ├── mac_computed_at_utc.txt # ISO 8601 UTC for every chain entry (forensic; same for all)
│   │   ├── kms_handle_uri.txt      # e.g. "plaintext-dev"
│   │   └── event-NNN.json          # One canonical event payload per file (JCS-encoded, chain-stamp fields excluded)
│   └── expected/
│       ├── session_key.hex                 # 32 bytes; HKDF output for (ikm, tenant_id)
│       ├── key_fingerprint.hex             # 16 bytes; SHA-256(utf8(tenant_id) || ikm)[:16]
│       ├── hkdf_inputs_digest.hex          # 32 bytes; SHA-256(salt || info_for_tenant || length_LE32)
│       ├── header.json                     # Full AuditFileHeader the writer emits
│       ├── event-NNN.canonical.hex         # Canonical bytes the MAC was computed over
│       ├── event-NNN.payload_hash.hex      # 32 bytes; HMAC-SHA-256 output
│       ├── event-NNN.chain_entry.json      # Full chain-stamp record (seq, prev_hash, payload_hash, key_version, key_fingerprint, format_version, mac_computed_at_utc, kms_handle_uri)
│       ├── merkle_root.hex                 # 32 bytes; RFC 6962 root
│       ├── sign_payload.txt                # The exact bytes signed by Ed25519 (v1.0a 10-line form)
│       └── signature.hex                   # 64 bytes
├── 002-multi-event-same-run/
├── 003-multi-run-same-day/
├── 008-jcs-edge-cases/
│   ├── 001-unicode-keys/
│   ├── 002-number-edge/
│   ├── 003-nested-objects/
│   └── 004-array-ordering/
├── 010-tenant-ikm-rotation-mid-day/
│   # Events 1-3 under key_version=1 with ikm_v1; events 4-5 under key_version=2 with ikm_v2.
│   # Same tenant. The seal record's key_versions = [1, 2]. Each chain entry stamps the
│   # appropriate (key_version, key_fingerprint) for its IKM generation. Verifier walks
│   # the rotation seamlessly via per-entry (tenant_id, key_version) -> IKM lookup.
├── 015-dual-algorithm-cosigned-seal/
├── 016-non-power-of-2-merkle/
└── negative/
    ├── N001-payload-hash-bit-flip/
    ├── N002-events-reordered/
    ├── N003-merkle-root-altered/
    ├── N004-signature-garbage/
    ├── N005-signature-wrong-tenant/
    ├── N006-key-fingerprint-flipped/
    ├── N007-unknown-key-version/
    ├── N008-entry-format-version-mismatch/
    ├── N009-header-format-version-v2/
    ├── N010-header-hkdf-digest-flipped/
    ├── N011-header-genesis-nonzero/
    ├── N012-cross-chain-tenant-mismatch/
    ├── N013-mid-write-truncation/
    ├── N014-botched-rotation/
    ├── N015-prev-hash-substituted/
    ├── N016-prev-hash-and-payload-recomputed-by-attacker/
    ├── N017-dual-algo-partial-coverage/
    ├── N018-dual-algo-not-in-posture/
    ├── N019-dual-algo-one-valid-one-invalid/
    ├── N020-algorithm-key-type-mismatch/
    ├── N021-routing-event-tampered/
    ├── N022-format-version-v1-1/
    └── N023-format-version-case-variant/
```

Each numbered directory is self-contained. A test runner reads `inputs/`, computes outputs, and compares to `expected/`.

**Critical conformance property.** Two implementations of v1 MUST produce byte-identical `payload_hash`, `key_fingerprint`, `session_key`, `hkdf_inputs_digest`, `merkle_root`, and `signature` for the same inputs. The corpus is the binding contract.

## 4. The reference test runner

`runner.go` provides a Go test harness:

```go
package testvectors

func RunTestCase(t *testing.T, dir string) {
    inputs := loadInputs(dir + "/inputs/")
    expected := loadExpected(dir + "/expected/")

    // Post-rework v1.0 HKDF inputs per spec §4.1:
    //   salt = HKDF_SALT  = b"ffiec.chain-of-custody.v1.salt"
    //   info = HKDF_INFO_BASE || b"|" || utf8(tenant_id)
    //   length = 32
    info := append([]byte("ffiec.chain-of-custody.v1.info|"), []byte(inputs.TenantID)...)
    sessionKey := hkdf.SHA256(
        inputs.IKM,
        []byte("ffiec.chain-of-custody.v1.salt"),
        info,
        32,
    )
    if !bytes.Equal(sessionKey, expected.SessionKey) {
        t.Fatalf("session key mismatch")
    }

    var prevHash [32]byte
    for _, event := range inputs.Events {
        canonical := jcs.Encode(event.Payload)
        h := hmac.New(sha256.New, sessionKey)
        h.Write(prevHash[:])
        h.Write(canonical)
        payloadHash := h.Sum(nil)
        if !bytes.Equal(payloadHash, expected.PayloadHashes[event.RunID][event.Seq]) {
            t.Fatalf("payload_hash mismatch at run=%s seq=%d", event.RunID, event.Seq)
        }
        copy(prevHash[:], payloadHash)
    }

    merkle := merkle.NewStreaming()
    for _, ph := range orderedPayloadHashes(inputs.Events) {
        merkle.Add(ph)
    }
    if !bytes.Equal(merkle.Root(), expected.MerkleRoot) {
        t.Fatalf("merkle root mismatch")
    }

    signPayload := buildSignPayload(inputs.TenantID, inputs.SealDate, expected.MerkleRoot)
    sig := ed25519.Sign(inputs.SigningKey, signPayload)
    if !bytes.Equal(sig, expected.Signature) {
        t.Fatalf("signature mismatch")
    }

    if !ed25519.Verify(inputs.PublicKey, signPayload, expected.Signature) {
        t.Fatalf("signature verification failed (own signature)")
    }
}
```

Other languages re-implement the same harness reading from the same on-disk files. The on-disk format (hex, PEM, plain JSON) is language-neutral.

## 5. Test case taxonomy

### 5.1 Happy-path tests

- **001 &mdash; single event, empty prev_hash.** The smallest case. Verifies that `seq=1` correctly uses the all-zero `prev_hash`.
- **002 &mdash; multi-event, same run.** Verifies the chain links across multiple events.
- **003 &mdash; multi-run, same day.** Verifies cross-run isolation within one tenant-day: two runs in the same day each maintain their own `prev_hash` state (both start at the all-zero genesis), and the day's Merkle root is computed over all `payload_hash` values ordered `(run_id, seq)` ascending. A naive implementation that chains the second run's `seq=1` from the first run's terminal `payload_hash` produces a different (non-conformant) Merkle root.

### 5.2 JCS-canonicalization edge cases

- **008-001 &mdash; Unicode keys.** Object keys that contain Unicode requiring NFC normalization. Verifies JCS Unicode handling.
- **008-002 &mdash; number edges.** `1`, `1.0`, `1e0`, `0.1`, `1.5`, very large integers. Verifies JCS number canonicalization (RFC 8785 §3.2.2).
- **008-003 &mdash; nested objects.** Multi-level nested structures. Verifies that key sorting applies recursively.
- **008-004 &mdash; array ordering.** Arrays whose order is semantically meaningful. Verifies JCS preserves array order while sorting object keys.

The JCS edge-case corpus is elevated to MUST under spec §5 — implementations MUST pass every fixture in `008-jcs-edge-cases/` to be conformant.

### 5.3 Master-rotation window

- **010 — tenant-IKM rotation mid-day.** A tenant-day where events 1-3 are captured under `key_version=1` (with `ikm_v1`) and events 4-5 are captured under `key_version=2` (with `ikm_v2`) due to IKM rotation mid-day. Verifies the seal record carries `key_versions = [1, 2]` (sorted ascending), the verifier resolves which IKM per event via the per-entry `(tenant_id, key_version)` lookup at spec §7 step 7 and the per-entry `key_fingerprint` check at step 8, and the day passes.

### 5.4 Dual-algorithm cosigned seal

- **015 — dual-algorithm cosigned seal.** A seal record carrying both Ed25519 and a (stub) post-quantum signature in the `signatures` list. Verifies Variant B per-algorithm `sign_payload` reconstruction (spec §4.3.2): each algorithm is computed over its own algorithm-bound `sign_payload` with line 3 carrying that algorithm's identifier. The post-quantum entry is a structural stub — the cryptographic algorithm dispatches as configured but the underlying bytes are placeholders — so the test exercises the verifier's per-algorithm dispatch logic without depending on a NIST-finalized post-quantum HSM product. Negative companions N017-N020 exercise dual-algorithm failure modes.

### 5.5 Non-power-of-2 Merkle balancing

- **016 — non-power-of-2 Merkle balancing.** A tenant-day with a leaf count that is not a power of 2 (e.g., 3, 5, 7 events). Verifies the RFC 6962 right-most-subtree-carries-the-unpaired-leaves rule. A naive Merkle implementation that pads to the next power of 2 with zero leaves produces a different (non-conformant) root.

### 5.6 Negative tests

A separate `spec/test-vectors/negative/` directory contains test cases where the verifier MUST report failure with a specific, named failure reason mapped to the spec §7 procedure step that produced the failure.

| Case | What's tampered | Step | Expected reason |
|---|---|---|---|
| **N001** | Flipped bit in `payload_hash` | 9 | `payload_hash MAC mismatch at seq N` |
| **N002** | Reordered events without re-chaining | 6 | `chain link broken at seq N` |
| **N003** | Merkle root altered in seal | 10 | `merkle root mismatch — ledger contents do not produce sealed root` |
| **N004** | Signature replaced with garbage | 11 | `signature verification failed` |
| **N005** | Signature for wrong tenant | 11 | `signature verification failed` (the `tenant_id` is in the signed payload) |
| **N006** | `key_fingerprint` flipped to an arbitrary 16 bytes | 8 | `key_fingerprint mismatch at seq N: looked-up IKM does not match the entry's recorded fingerprint` (NO MAC compute happens) |
| **N007** | `key_version` set to a generation not in the test IKM registry | 7 | `unknown key_version: no IKM for (tenant=T, key_version=V) at seq N` (NO MAC compute happens) |
| **N008** | `format_version` on an entry set to `"v2"` | 5 | `format_version mismatch at seq N` |
| **N009** | `header.format_version` set to `"v2"` | 1 | `format_version v2 not supported by this verifier (running v1)` |
| **N010** | `header.hkdf_inputs_digest` flipped | 2 | `header HKDF inputs do not match running v1 inputs` |
| **N011** | `header.genesis_hash` set to non-zero bytes | 3 | `header genesis_hash does not match v1 constant` |
| **N012** | `event.tenant_id` differs from `header.tenant_id` (cross-chain lift) | 4 | `cross-chain lift detected at seq N (event.tenant_id mismatch)` |
| **N013** | Audit file's last byte is not `\n` (mid-write truncation) | (file pre-flight) | `audit file ends mid-line — possible mid-write crash; the writer's last append did not complete` |
| **N014** | Botched rotation: `key_version=1` re-used for a different IKM (same tenant) | 8 | `key_fingerprint mismatch at seq N` (this is the load-bearing rotation defence — the test verifies the verifier catches it before any MAC compute) |
| **N015** | `entry.prev_hash` substituted but `payload_hash` left alone | 6 | `chain link broken at seq N` (structural walk catches it) |
| **N016** | `entry.prev_hash` substituted AND `payload_hash` recomputed by an attacker without the IKM | 6 | `chain link broken at seq N` (structural walk catches the substitution; MAC step would have caught the recompute regardless because the attacker's MAC would not match the verifier's MAC computed with the real IKM) |
| **N017** | Dual-algorithm seal with one algorithm missing during institution-declared dual-algorithm posture (partial coverage) | 11 | `partial-coverage seal: single-algorithm signature during institution's declared dual-algorithm posture` (case (b) per spec §7 step 11) |
| **N018** | Dual-algorithm seal carrying an algorithm not on the institution's declared posture list | 11 | `algorithm not on institution's declared posture list at seal_date {D}` (case (c)) |
| **N019** | Dual-algorithm seal where one signature validates and one does not (load-bearing case) | 11 | `co-signed seal failure: algorithm X validated, algorithm Y did not` (case (e)) |
| **N020** | Algorithm identifier mismatches the public key's algorithm | 11 | `algorithm/key-type mismatch at signature verification` |
| **N021** | Routing event tampered (post-capture mutation of `audit.routing.*` attributes) | 9 | `payload_hash MAC mismatch at seq N` — the routing decision is bound under the chain MAC like any other entry; tampering surfaces at the per-event MAC check. Fixture regenerated under v1.0a wire form. |
| **N022** | `header.format_version` set to `"v1.1"` (variant of the v1 family) | 1 | `format_version v1.1 not supported by this verifier (running v1)` — exact-match refusal per spec §7 step 1; `"v1"` is the only conformant value. |
| **N023** | `header.format_version` set to `"V1"` (case variant) | 1 | `format_version V1 not supported by this verifier (running v1)` — case-sensitive refusal; `"v1"` lowercase is the only conformant value. |

The negative tests guard against false-positive verifier behavior. A verifier that produces a "pass" result on any negative test is broken. Cases N006-N023 exercise the inviolate properties from spec §4.1, the twelve-step procedure from spec §7, the dual-algorithm dispatch in spec §4.3.2 and §7 step 11 cases (a)-(e), and the format-version exact-match refusal in §7 step 1.

### 5.7 Recommended future cryptographic-attack negative cases

The N001-N023 corpus covers structural attacks (reordering, bit-flips, signature-on-wrong-tenant, format-version variants). It does not yet cover cryptographic-attack scenarios that target the primitives directly. These four cases extend the cryptographic-attack test surface for academic citation and IETF CFRG review. They are recommended future cases, not v1.0-blocking — the v1.0 corpus is sufficient for the v1.0 conformance bar.

| Case | What's tampered | What it validates |
|---|---|---|
| **Merkle length-extension variant** | Replace the RFC 6962 leaf-prefix (0x00) with a null byte or omit it entirely; recompute leaves and root | The domain-separator defense in spec §4.2. A conforming verifier rejects the modified root because the recomputed root under the correct leaf-prefix does not match. The case documents that the prefix is load-bearing, not decorative. |
| **Non-strict Ed25519 signature** | Submit a signature in a non-canonical form that some non-strict implementations would accept | RFC 8032 §8.4 strict canonicalization. A conforming verifier rejects the non-canonical signature; a non-strict verifier (which the spec rejects as non-conformant) would accept it. The case is the canary for cross-implementation interop where one party uses a non-strict library. |
| **HKDF info collision** | Manually forge two entries with the same `tenant_id` pointing to different IKMs; recompute fingerprints under each | Per-tenant isolation. The verifier's per-entry `key_fingerprint` check at spec §7 step 8 catches the divergence: one entry's fingerprint matches the registered IKM, the other does not. The case validates that fingerprint-mismatch detection is the load-bearing defense for HKDF info collisions, not raw HKDF properties. |
| **Birthday-bound Merkle collision (informational)** | Documentation case showing the 2^128 attack space against SHA-256 second-preimage in the Merkle construction | No practical attack at v1.0 deployment scale. The case is informational — it confirms that birthday-bound attacks against SHA-256 require approximately 2^128 work under current cryptanalysis, well beyond practical compute even in adversarial scenarios. The case exists to close the academic-review question, not to exercise a runtime check. |

The four cases close the cryptographic-attack test surface for academic citation and IETF CFRG review. The v1.0 corpus is sufficient for the v1.0 conformance bar; these cases extend coverage for higher-assurance review contexts.

## 6. Versioning the corpus

Each test case has a `spec_version` field in `description.md`. The runner only runs tests whose `spec_version` matches the implementation&rsquo;s declared version.

The corpus grows over time:

- v1.0-final-amendment ships with positive cases 001, 002, 003, 008 (with four sub-fixtures), 010, 015, and 016, alongside negative cases N001-N023.
- Future v1.x amendments add new cases that exercise new features (e.g., a NIST-finalized post-quantum signature alongside Ed25519, hardware-attested edge-device key custody, cross-region run continuation).
- v1.0-conforming implementations continue to pass the v1.0 cases under future tooling.

The corpus is **append-only**. Existing test cases are never modified after a spec version locks. A bug in a test case is fixed by adding a new case that supersedes the buggy one; the old case is marked deprecated but not deleted.

## 7. Test-case provenance

Every test case includes:

- **The author** &mdash; who created the case
- **The reviewer** &mdash; who verified it independently
- **The motivation** &mdash; what defect or ambiguity the case exists to prevent

Bug-driven test cases (a conforming implementation produced wrong output for some real input; the case is added to prevent regression) are flagged with the original bug ID.

## 8. Cross-implementation conformance

Once two independent implementations exist (the reference Go implementation plus at least one other), the conformance regimen is:

1. Both implementations run the corpus. Both pass.
2. Each implementation produces output for a generated stream of events. The other consumes the output. Both verify.
3. Periodically, a randomly-generated event stream is used as a fuzzing input. Both implementations are required to produce identical bytes.

Step 3 is the most aggressive test. A defect that the corpus does not exercise but that produces divergent output is fuzzed out by random input.

## 9. The corpus as a regulatory artifact

The corpus is the FFIEC&rsquo;s testable basis for declaring an implementation conformant.

- The FFIEC could publish its own test vectors in a future revision and require conforming implementations to pass them.
- The corpus in this repository becomes the candidate test-vector set the FFIEC reviews and (potentially) adopts.
- A conformance certification regime is the FFIEC&rsquo;s decision; we provide the technical basis on which it can rest.

## 10. Auditor's-lens review

| Question | Answer |
|---|---|
| Can the corpus exercise every observable behavior of the spec? | Not exhaustively, no. The corpus catches the cases we&rsquo;ve identified. We add new cases as ambiguities are surfaced. The corpus is a growing artifact, not a closed set. |
| What if a test case has a bug? | Add a new case that exercises the correct behavior; deprecate the buggy one; document the change in the spec change log. Existing implementations may need to update. |
| How do we know two implementations producing identical output prove correctness rather than identical wrongness? | We don&rsquo;t, with two implementations alone. Three or more independent implementations agreeing strongly suggests correctness. The corpus plus fuzzing (Section 8 step 3) is the practical mitigation. |
| What about cryptographic breaks (e.g., SHA-256 collision)? | Out of scope of conformance testing. The spec depends on the cryptographic primitives; if those break, the spec breaks. |
| Formal model of the spec | A TLA+ or Tamarin model of the protocol for cross-checking the corpus is a v1.x roadmap candidate. The conformance corpus plus cross-implementation fuzzing carries the v1.0 conformance posture; a formal model is additive evidence that lifts the assurance bar above what testing alone can provide. |
