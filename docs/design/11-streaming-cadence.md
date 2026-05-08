# 11 — Configurable seal cadence (§10.27 design)

> **Status.** Phase 1 design. Pending Auditor review.
>
> **Spec section.** §10.27 Configurable seal cadence (normative). Default daily; configurable from per-second through weekly. Normative values: `"per_second"`, `"per_minute"`, `"per_hour"`, `"hourly"`, `"daily"`, `"weekly"`. Sub-daily values are the **streaming-mode** seal.
>
> **Scope.** This design covers the Go reference implementation (`core/`, `ledger/`, `verifier/`). The .NET (`Herald/Modules/`) and Python (`Herald.Py/src/herald/`) parallelizations follow per the [Herald rollup plan](../herald-rollup-plan.md).

---

## 1. Goals

1. Support seal-cadence values from per-second through weekly, recorded as a `cadence` field on every seal record.
2. Bind the `cadence` field under the §4.3 `sign_payload` form so a tampered cadence value is detected at signature verification.
3. Verifier confirms the seal records cover the period continuously without gaps under the institution's claimed cadence.
4. Do not break the current daily-default behavior. Daily-cadence institutions MUST see no behavioral change.
5. Stay CUPID/DRY/low-cognitive-complexity per the rollup-plan targets (§4 of `herald-rollup-plan.md`).

## 2. Non-goals

- Streaming-mode IKM rotation discipline (§10.28) — that's Phase 2.
- Streaming-mode verifier procedure (§10.29) — co-Phase 1, separate design doc.
- Trusted-time integration (§10.30) — co-Phase 1, separate design doc.
- Multi-region cadence reconciliation under §10.15 — out of Phase 1; deferred.
- HSM throughput tuning — out of Phase 1; the institution's CC8.1 names the HSM and its throughput.

## 3. Current Go reference state (bootstrap)

The current `verifier/internal/verify/` implements a simplified bootstrap of the spec primitives. As of this design:

- `core/chain/` and `core/merkle/` directories exist but contain no Go source. The chain and Merkle primitives currently live inside `verifier/internal/verify/` as `mac.go`, `merkle.go`, etc.
- `SealRecord` (in `verifier/internal/verify/ledger.go`) carries `Type`, `SealDate`, `MerkleRoot`, `LeafCount`, `SignatureEd25519`, `SigningKeyFingerprint` &mdash; **no `cadence` field**.
- `FixtureBuilder.SealEntries` (in `fixture.go`) signs the seal with a fixed Ed25519 key &mdash; **no cadence-awareness in the seal-job**.
- The verifier's chain walk in `chain.go` and `pipeline.go` does not consult cadence.
- There is no ledger server with a real seal-job timer; the bootstrap uses `FixtureBuilder` to produce synthetic ledgers.

This design extends the bootstrap minimally to add cadence support without rebuilding the primitives. A future Phase (post-rollup) consolidates the chain primitives into `core/chain/` and `core/merkle/` per the spec's `core/` design.

## 4. Design

### 4.1 The `Cadence` type

Add to `core/chain/cadence.go` (creating the file under the existing `core/chain/` directory):

```go
// Package chain implements the per-event MAC chain primitive.
package chain

// Cadence is the seal aggregation cadence per §4.2.1 + §10.27.
//
// The default is CadenceDaily. CadencePerSecond, CadencePerMinute, CadencePerHour
// are the streaming-mode cadences. CadenceHourly, CadenceWeekly are non-streaming
// non-daily values.
//
// The string form is the wire-format byte value bound under §4.3 sign_payload.
type Cadence string

// Normative cadence values per §10.27. The wire-format byte values are these
// exact strings; verifiers reject any unrecognized value at §7 step 11.
const (
	CadencePerSecond Cadence = "per_second"
	CadencePerMinute Cadence = "per_minute"
	CadencePerHour   Cadence = "per_hour"
	CadenceHourly    Cadence = "hourly"
	CadenceDaily     Cadence = "daily"
	CadenceWeekly    Cadence = "weekly"
)

// IsStreaming reports whether the cadence is sub-daily (streaming-mode per §10.27).
func (c Cadence) IsStreaming() bool {
	switch c {
	case CadencePerSecond, CadencePerMinute, CadencePerHour:
		return true
	}
	return false
}

// IsValid reports whether the cadence is one of the normative values.
// Unrecognized values are non-conformant per §10.27.
func (c Cadence) IsValid() bool {
	switch c {
	case CadencePerSecond, CadencePerMinute, CadencePerHour,
		CadenceHourly, CadenceDaily, CadenceWeekly:
		return true
	}
	return false
}

// Interval returns the time-duration interval the cadence implies. Returns 0
// for unrecognized cadences (callers should pre-check with IsValid).
func (c Cadence) Interval() time.Duration {
	switch c {
	case CadencePerSecond:
		return time.Second
	case CadencePerMinute:
		return time.Minute
	case CadencePerHour, CadenceHourly:
		return time.Hour
	case CadenceDaily:
		return 24 * time.Hour
	case CadenceWeekly:
		return 7 * 24 * time.Hour
	}
	return 0
}
```

CUPID notes:
- **Composable**: `Cadence` is a value type with no dependencies; consumers (seal job, verifier, fixture builder) consume it without coupling to a configuration manager.
- **Predictable**: every method is deterministic; no hidden state.
- **Idiomatic Go**: a string-typed constant + methods is the idiomatic Go enum.
- **Domain-based**: the type models the spec's domain primitive directly.

Cognitive complexity: each method is a switch with ≤6 cases. Linear, ≤10 cognitive complexity.

### 4.2 SealRecord extension

Add `Cadence` field to `verifier/internal/verify/ledger.go`'s `SealRecord`:

```go
type SealRecord struct {
	Type                  string         `json:"type"`
	SealDate              string         `json:"seal_date"`
	MerkleRoot            string         `json:"merkle_root"`
	LeafCount             int            `json:"leaf_count"`
	SignatureEd25519      string         `json:"signature_ed25519"`
	SigningKeyFingerprint string         `json:"signing_key_fingerprint"`
	Cadence               chain.Cadence  `json:"cadence,omitempty"` // §10.27; default daily when absent
}
```

The `omitempty` annotation preserves backward compatibility: existing daily-cadence ledgers (no `cadence` field) deserialize cleanly; the consumer treats absent as `CadenceDaily`.

### 4.3 Cadence-aware fixture builder

Extend `FixtureBuilder` to accept a cadence:

```go
type FixtureBuilder struct {
	TenantBindingKDFLabel string
	SealDate              string
	MasterIKM             []byte
	SealingPriv           ed25519.PrivateKey
	Cadence               chain.Cadence  // default CadenceDaily when zero-value
}
```

`SealEntries` populates the new `Cadence` field on the produced `SealRecord`:

```go
cadence := fb.Cadence
if cadence == "" {
	cadence = chain.CadenceDaily
}
if !cadence.IsValid() {
	return SealRecord{}, fmt.Errorf("FixtureBuilder.Cadence is non-conformant: %q", cadence)
}
s := SealRecord{
	// ... existing fields ...
	Cadence: cadence,
}
```

### 4.4 Sign-payload binding

Per §10.27 normative text, the `cadence` field MUST be bound under the §4.3 `sign_payload` form. The current `canonicalForSealSignature` in `canonical.go` produces the seal-signature canonical bytes; extend it to include cadence on a new line.

The sign-payload byte form (in the rolled-up draft, the §4.3 form) is multiline LF-separated. The cadence field becomes a new line at a fixed position. Test vector `020-streaming-seal-cadence-1s` pins the byte values.

This is the highest-cognitive-complexity change in Phase 1: adding a line to the sign-payload form requires coordinated update of the signing path AND the verification path AND the test-vector fixtures. The Auditor review will confirm the byte-form is reproducible across the three paths.

### 4.5 Verifier cadence validation

The verifier confirms cadence-record continuity. For a tenant operating at non-default cadence:

- Read the seal record's `cadence` field.
- Compute the expected number of seal records covering the verification period: `period / cadence.Interval()`.
- Walk the seal records in `seal_date` order; confirm no gaps. Adjacent seals must have `seal_date` differing by exactly one cadence-interval.
- For default daily cadence, the existing one-seal-per-day expectation is preserved.

The verifier emits an anomaly line `cadence: <value>` when the cadence is non-default, so an examiner reading the verifier output can immediately see the institution's cadence posture.

### 4.6 Sign-job timer (ledger-side)

The current bootstrap uses `FixtureBuilder` (a one-shot synthesizer), not a persistent seal-job timer. The ledger server in `ledger/cmd/ledger/cmd_serve.go` currently has no seal-job at all — the bootstrap leaves seal-record production to `FixtureBuilder` for testing.

For Phase 1, the cadence work in the ledger server is **deferred to Phase 1.5** because:
- A real cadence-aware seal-job timer requires a persistent ledger storage layer (Phase 1.5 deliverable).
- The fixture builder + verifier pair is sufficient to validate the cadence wire-format and verification path.

Phase 1.5 adds the ledger seal-job timer once the persistent storage layer lands.

## 5. Test coverage

### 5.1 Unit tests

- `core/chain/cadence_test.go`: covers `Cadence.IsStreaming()`, `IsValid()`, `Interval()` for each value.
- `verifier/internal/verify/cadence_test.go`: covers cadence-record-continuity validation for daily, hourly, per-minute, per-second cadences, plus gap-detection failure cases.

### 5.2 Test vectors

- `spec/test-vectors/020-streaming-seal-cadence-1s/` — byte-identical fixture for 1-second cadence with 5 seal records over a 5-second window.

### 5.3 Cross-implementation byte-equivalence

(Phase 5/6) — when .NET and Python land, byte-equivalence tests confirm the cadence-bearing seal records are byte-identical across all three implementations.

## 6. Migration / backward compatibility

- Default daily-cadence institutions: no behavioral change. The seal record's `cadence` field is `omitempty` and defaults to `CadenceDaily` when absent.
- Streaming-mode institutions: opt in by setting `Cadence` on `FixtureBuilder` (or, in Phase 1.5, by configuring the seal-job timer).

The wire-format identifier `"v1"` is unchanged. The cadence extension is additive within `"v1"`.

## 7. Open questions for Auditor

1. The `omitempty` posture: a default-daily seal omits `cadence`. The verifier defaults absent → daily. Is this acceptable, or should every seal carry an explicit `cadence` value? The spec is silent on the omission; the conservative posture would make `cadence` REQUIRED.
2. Cadence-interval gap-detection in the verifier: a missed cadence-interval seal MUST surface as a verifier anomaly. Is the proposed mechanism (check adjacent seal_date deltas equal cadence.Interval()) sufficient, or does it need cross-region awareness for §10.15 institutions?
3. Sign-payload-line ordering: where in the sign-payload byte-form does the `cadence` field land? The rolled-up §4.3 form already includes cadence between `dev_mode` and `key_versions_canon` (per the historical 12-line form). Confirm this position remains canonical.

## 8. Auditor review checklist

Pre-implementation review:
- [ ] §10.27 spec section accurately interpreted
- [ ] CUPID targets named per design element
- [ ] Backward compatibility preserved
- [ ] Test coverage matrix is complete
- [ ] Open questions are answerable

Post-implementation review (rounds 1+2):
- [ ] Cyclomatic complexity ≤10 per function
- [ ] No DRY violations (cadence handling not duplicated across packages)
- [ ] Constant-time comparisons preserved where applicable
- [ ] Test vectors pass byte-for-byte
- [ ] Verifier output discipline preserved (Status: PASS / Step / Reason)

## 9. Implementation order

1. `core/chain/cadence.go` — `Cadence` type with methods. Auditor pre-implementation review.
2. `core/chain/cadence_test.go` — unit tests.
3. Auditor code review round 1 (CRITICAL/MAJOR resolution).
4. `verifier/internal/verify/ledger.go` — `SealRecord.Cadence` field.
5. `verifier/internal/verify/fixture.go` — `FixtureBuilder.Cadence` field, default-handling.
6. `verifier/internal/verify/canonical.go` — sign-payload binding the cadence field.
7. `verifier/internal/verify/cadence_test.go` — verifier cadence-validation tests.
8. Spec test vector `spec/test-vectors/020-streaming-seal-cadence-1s/` with byte-stable fixture.
9. Auditor code review round 2 (polish — CUPID/DRY refinement).
10. Auditor sign-off.
