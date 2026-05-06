# 07 — Verifier design

> **What this doc is.** The architecture of `verifier/`. The standalone offline CLI that an examiner runs against a ledger to verify chain integrity, Merkle proofs, and HSM root signatures &mdash; with no dependencies on the institution or its vendor.

## 1. Design constraints

The verifier exists to satisfy one constraint: **an examiner can verify the chain without trusting the institution or its vendor**. Every design choice flows from that constraint.

| Constraint | Implication |
|---|---|
| No network calls | Pure offline binary. No DNS, no HTTPS, no telemetry. |
| Single binary | Static build with `CGO_ENABLED=0`. No dynamic linking. |
| Deterministic output | Same inputs &rarr; byte-for-byte identical report. Two examiners produce identical PDFs that can be compared by hash. |
| Fail closed | Any ambiguity is a failure. The verifier never silently passes a partial verification. |
| Read-only | The verifier never mutates the ledger or any other input. |
| Auditor-readable | Output is human-readable; the verifier explains what it checked and why each check passed or failed. |

## 2. Inputs

The verifier&rsquo;s entire trusted-computing-base for one verification run:

```
verifier verify \
  --ledger <path>           # Ledger snapshot or live read-only mount
  --root-key <path>         # Tenant public key (PEM-encoded Ed25519)
  --tenant-id <string>      # Optional; defaults to the tenant_id in the ledger
  --from <date>             # Optional; defaults to first day in ledger
  --to <date>               # Optional; defaults to last day in ledger
  --report <path>           # Output PDF path
  --json-report <path>      # Optional structured JSON output for tooling
  --strict                  # Optional; treat any warning as a failure
```

Each input is a flat file the examiner controls. The `--root-key` path is the regulator-held public key half. There is nothing else.

## 3. The verification flow

```mermaid
sequenceDiagram
    autonumber
    participant Ex as Examiner
    participant V as Verifier
    participant L as Ledger snapshot
    participant K as Tenant public key

    Ex->>V: verify --ledger L --root-key K
    V->>K: read public key
    V->>L: open snapshot

    loop for each tenant_id, day pair in range
        V->>L: stream events for (tenant, day) ordered by run_id, seq
        V->>V: walk HMAC chain<br/>(verify each prev_hash, payload_hash)
        V->>V: compute Merkle root from leaves
        V->>L: read recorded daily seal for (tenant, day)
        V->>V: verify seal.merkle_root == computed_root
        V->>V: verify Ed25519(K, sign_payload, signature)
        V-->>Ex: per-day result (pass / fail / late-binding count)
    end

    V->>Ex: write PDF report
    V->>Ex: write JSON report (if requested)
    V-->>Ex: exit 0 if all pass; exit 1 if any fail
```

Every step has a deterministic outcome. The PDF is generated from the structured results.

## 4. The verification primitives

### 4.1 Chain walk

For each `(tenant_id, run_id)` in the day:

```
events = SELECT * FROM ledger
         WHERE tenant_id=T AND DATE(captured_at, 'UTC')=D AND run_id=R
         ORDER BY seq ASC

expected_prev = 0x00...00
session_keys_seen = {}

FOR each event in events:
  IF event.session_key_id NOT IN session_keys_seen:
    # Without the master key, we cannot derive the session key.
    # We CAN verify the chain structurally even without the key by recomputing
    # in chain-walk mode (HMAC with known prev_hash and canonical payload).
    # The verifier requires the session key for HMAC equality verification.
    # If the master key is not available, the verifier marks this event as
    # "structurally consistent, key-bound verification skipped" and recommends
    # a follow-up examination with key access.
    ...

  # If session key is available:
  recomputed_hash = HMAC-SHA256(session_key, event.prev_hash || JCS(event.payload))
  IF recomputed_hash != event.payload_hash:
    FAIL: chain hash mismatch at (run_id=R, seq=event.seq)
  IF event.prev_hash != expected_prev:
    FAIL: chain link broken at (run_id=R, seq=event.seq)
  expected_prev = event.payload_hash
```

### 4.2 Merkle recomputation

```
day_events = SELECT payload_hash FROM ledger
             WHERE tenant_id=T AND DATE(captured_at, 'UTC')=D
             ORDER BY run_id ASC, seq ASC

streaming_merkle = NEW StreamingMerkle()
FOR each ph in day_events:
  streaming_merkle.add(ph)

computed_root = streaming_merkle.root()
```

### 4.3 Seal verification

```
seal = SELECT * FROM daily_seals WHERE tenant_id=T AND seal_date=D

IF seal.merkle_root != computed_root:
  FAIL: Merkle root mismatch &mdash; the ledger contents do not produce the sealed root

sign_payload = "ffiec-ai-chain-v1\n" + T + "\n" + iso8601(D) + "\n" + hex(seal.merkle_root)

IF NOT ed25519.Verify(public_key, sign_payload, seal.signature):
  FAIL: signature verification failed
```

## 5. The PDF report

The PDF is the artifact the examiner takes home. Determinism matters: two examiners running the verifier independently produce identical PDFs.

### 5.1 Structure

```
Page 1: Cover
  - Examiner instance (from --report-author flag, optional)
  - Verifier version, spec version
  - Ledger snapshot path and SHA-256 hash
  - Public key fingerprint
  - Verification range (from..to)
  - Pass/fail summary
  - Examiner signature line

Pages 2..N: Per-day detail
  - Date
  - Event count
  - Late-binding count
  - Run count
  - Merkle root (hex)
  - Recorded seal merkle_root (hex)
  - Merkle match: pass/fail
  - Signature verification: pass/fail
  - HMAC chain verification: pass/fail
  - Per-failure detail (if any)

Page N+1: Anomalies
  - List of any anomalies that did not cause failure but warrant attention:
    * Late-binding events
    * Empty days
    * Sealing delays > 24 hours
    * Session keys not available for HMAC verification

Page N+2: Methodology
  - The exact algorithms used (with citations to the spec)
  - The order of operations
  - The pass/fail rules
```

### 5.2 PDF determinism

The PDF includes a generation timestamp by default. With `--deterministic`, the verifier emits a PDF with no timestamp and no randomized fields, producing byte-for-byte reproducible output across runs.

## 6. The JSON report

For tooling integration. Same structured data as the PDF, in machine-readable form:

```json
{
  "verifier_version": "v1.0.0",
  "spec_version": "v1.0",
  "ledger_sha256": "...",
  "public_key_fingerprint": "...",
  "range": { "from": "2026-04-01", "to": "2026-04-30" },
  "tenant_id": "tenant_acme_prod",
  "summary": {
    "days_verified": 30,
    "days_passed": 30,
    "days_failed": 0,
    "events_total": 4218391,
    "late_binding_total": 0
  },
  "days": [
    {
      "date": "2026-04-01",
      "events": 142087,
      "runs": 8312,
      "merkle_root": "...",
      "recorded_root": "...",
      "merkle_match": true,
      "signature_valid": true,
      "chain_walk": "pass",
      "anomalies": []
    },
    ...
  ]
}
```

## 7. Operational modes

### 7.1 `verifier verify` &mdash; primary mode

The full verification flow. Produces the PDF and exits with status reflecting overall pass/fail.

### 7.2 `verifier walk`

Walk one specific run. For incident-investigation use:

```
verifier walk --ledger L --run-id r_a3f29b71c
```

Output: every event in the run, in order, with chain status per event.

### 7.3 `verifier diff`

Compare two ledger snapshots:

```
verifier diff --before L1 --after L2
```

Output: events present in L2 but not in L1 (additions), events present in L1 but not in L2 (defects, since the ledger should be append-only), and any events with conflicting payload_hash between the two snapshots (definitive evidence of tampering).

### 7.4 `verifier signedroot`

Just verify a single signed root, given the root and the public key:

```
verifier signedroot --root <hex> --signature <hex> --public-key <path> \
  --tenant tenant_acme_prod --date 2026-04-01
```

Use case: an examiner has a daily seal record but not the full ledger; this command confirms the seal was signed by the expected key. The full ledger is required for chain and Merkle verification.

## 8. Build and distribution

### 8.1 Static build

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o verifier ./cmd/verifier
```

The `-s -w` flags strip debug symbols. The static build runs on any Linux, macOS, or Windows examiner laptop without dependencies.

### 8.2 Cross-compilation matrix

The release pipeline builds for:

- `linux/amd64`, `linux/arm64`
- `darwin/amd64`, `darwin/arm64`
- `windows/amd64`

All builds are reproducible (deterministic build flags, no embedded build timestamps). The release artifacts are signed with cosign.

### 8.3 Distribution

- GitHub Releases (signed binaries)
- Direct download from `verify.ffiec-ai-chain.org` (when registered)
- Reference Docker image for CI integration

The examiner is expected to obtain the binary from one of these channels, verify its signature, and run it locally.

## 9. Auditor's-lens review

| Question | Answer |
|---|---|
| Can the verifier produce a false-positive (claim a chain is valid when it isn&rsquo;t)? | Only via a defect in the verifier code. The verification logic is small, well-tested by the conformance corpus, and CodeQL-scanned. The reference implementation accepts independent re-implementation by other parties. |
| Can the verifier produce a false-negative (claim a chain is invalid when it&rsquo;s valid)? | Only via a defect or a fundamental cryptographic break. The verifier reports specific reasons for failure; the examiner can re-verify by hand or with another verifier instance. |
| Does the verifier need network access? | No. Verified by the build constraint that bans `net/http` import in the verification path. |
| Can the institution influence what the verifier checks? | No. The verifier reads only the ledger (untrusted input) and the public key (regulator-held). The institution&rsquo;s ops team has no way to inject behavior into the verifier process. |
| What if the examiner&rsquo;s laptop is compromised? | A compromised laptop can produce a false report. This is a residual risk; the mitigation is examiner laptop hygiene (separation of duties, dedicated audit machines). The verifier itself is not the right place to defend against laptop compromise. |
| What if the verifier itself is compromised (a malicious build)? | Examiners verify the binary signature with cosign before running. Reproducible builds enable third-party re-verification. The threat model accepts this is a residual risk if cosign is compromised. |
| Open issue | Master-key-less verification mode. Today the HMAC chain check requires the session key (and thus the master). If the institution will not share the master with the examiner, the verifier degrades to structural verification (prev_hash linking, Merkle, signature) without HMAC equality. v1.1 may add a zero-knowledge mode where the institution proves chain knowledge without disclosing the master. |
