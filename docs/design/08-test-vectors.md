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

```
spec/test-vectors/
├── README.md                       # Index of test cases
├── runner.go                       # Reference Go test runner
├── runner_test.go                  # Validates the runner itself
├── 001-single-event-empty-prev/
│   ├── description.md              # What this case exercises
│   ├── inputs/
│   │   ├── master.key              # 32 bytes test master (clearly marked "TEST USE ONLY")
│   │   ├── process_uuid.txt        # 16 bytes hex
│   │   ├── tenant_id.txt
│   │   ├── ed25519_private.pem     # Test signing key
│   │   ├── ed25519_public.pem      # Corresponding public key
│   │   ├── seal_date.txt           # ISO 8601 date
│   │   └── event-001.json          # Canonical event payload (JCS-encoded)
│   └── expected/
│       ├── session_key.hex
│       ├── session_key_id.txt
│       ├── event-001.payload_hash.hex
│       ├── merkle_root.hex
│       └── signature.hex
├── 002-multi-event-same-run/
├── 003-multi-run-same-day/
├── 004-empty-day/
├── 005-out-of-order-input/
├── 006-late-binding/
├── 007-ten-thousand-events-perf/
├── 008-jcs-edge-cases/
│   ├── 001-unicode-keys/
│   ├── 002-number-edge/
│   ├── 003-nested-objects/
│   └── 004-array-ordering/
├── 009-session-key-rotation/
├── 010-tenant-key-rotation/
└── 011-cross-spec-version/
```

Each numbered directory is self-contained. A test runner reads `inputs/`, computes outputs, and compares to `expected/`.

## 4. The reference test runner

`runner.go` provides a Go test harness:

```go
package testvectors

func RunTestCase(t *testing.T, dir string) {
    inputs := loadInputs(dir + "/inputs/")
    expected := loadExpected(dir + "/expected/")

    sessionKey := hkdf.SHA256(
        inputs.MasterKey,
        append(inputs.ProcessUUID, []byte(inputs.TenantID)...),
        []byte("ffiec-ai-chain-v1"),
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
- **003 &mdash; multi-run, same day.** Verifies the Merkle ordering rule (`run_id ASC, seq ASC`).

### 5.2 Edge-case tests

- **004 &mdash; empty day.** The seal for a day with zero events. Verifies the `SHA-256("")` empty-tree convention.
- **005 &mdash; out-of-order input.** Events presented to the ledger in reverse seq order. Verifies that the spec-required ordering is applied at Merkle time, not at receive time.
- **006 &mdash; late-binding event.** An event that arrives after the daily seal is computed. Verifies the late-binding flag and the next-day inclusion.
- **009 &mdash; session-key rotation mid-run.** A run that begins under one session key (process A) and continues under another (process B after restart). Verifies that the chain&rsquo;s `payload_hash` continues correctly because `prev_hash` is the only state that crosses the rotation.
- **010 &mdash; tenant-key rotation.** A tenant that rotates its master key between two days. Verifies that the seals on both days are valid against their respective public keys.

### 5.3 JCS-canonicalization edge cases

- **008-001 &mdash; Unicode keys.** Object keys that contain Unicode requiring NFC normalization. Verifies JCS Unicode handling.
- **008-002 &mdash; number edges.** `1`, `1.0`, `1e0`, `0.1`, `1.5`, very large integers. Verifies JCS number canonicalization (RFC 8785 §3.2.2).
- **008-003 &mdash; nested objects.** Multi-level nested structures. Verifies that key sorting applies recursively.
- **008-004 &mdash; array ordering.** Arrays whose order is semantically meaningful. Verifies JCS preserves array order while sorting object keys.

### 5.4 Performance tests

- **007 &mdash; ten thousand events.** A bulk test that exercises Merkle streaming. Verifies determinism at scale; serves as a basic performance benchmark (a conforming implementation MUST complete within a reasonable time bound to be considered usable).

### 5.5 Negative tests

A separate `spec/test-vectors/negative/` directory contains test cases where the verifier MUST report failure. Conforming verifiers produce specific, named failure reasons.

- **N001 &mdash; flipped bit in payload_hash.** The verifier reports `chain hash mismatch at <run, seq>`.
- **N002 &mdash; reordered events without re-chaining.** The verifier reports `chain link broken at <run, seq>`.
- **N003 &mdash; merkle root altered.** The verifier reports `merkle root mismatch &mdash; the ledger contents do not produce the sealed root`.
- **N004 &mdash; signature replaced with garbage.** The verifier reports `signature verification failed`.
- **N005 &mdash; signature for wrong tenant.** The verifier reports `signature verification failed` (the tenant_id in the payload is part of what&rsquo;s signed).

The negative tests guard against false-positive verifier behavior. A verifier that produces a "pass" result on any negative test is broken.

## 6. Versioning the corpus

Each test case has a `spec_version` field in `description.md`. The runner only runs tests whose `spec_version` matches the implementation&rsquo;s declared version.

The corpus grows over time:

- v1.0 ships with cases 001&ndash;011 plus negative N001&ndash;N005.
- v1.1 adds new cases that exercise new features (e.g., post-quantum signature alongside Ed25519).
- v1.0-conforming implementations continue to pass the v1.0 cases under v1.1 tooling.

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
| Open issue | A formal model of the spec (TLA+ or similar) for cross-checking the corpus. v1.1+ candidate. |
