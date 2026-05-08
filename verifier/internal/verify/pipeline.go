package verify

import (
	"crypto/ed25519"
	"errors"
	"fmt"
)

// Result is the verifier's report. Each Step records one named check and
// its outcome; the report is meant to be predictable and minimal so two
// verifiers running on identical inputs produce byte-identical reports.
type Result struct {
	Path           string
	EntryCount     int
	SealDate       string
	SignerKey      string // signing_key_fingerprint as printed
	StructuralPass bool
	MACPass        bool   // false if the master-key path was not run
	Steps          []Step // ordered, one per executed check
}

// Step is one named check's outcome.
type Step struct {
	Name string
	OK   bool
	Note string
}

// Plan controls which steps Verify runs.
//
//   - SealingKey is required. Without it, the verifier cannot validate the
//     seal — and a verifier that does not validate the seal is not telling
//     the examiner anything load-bearing about integrity.
//   - MasterIKM is optional. When supplied, the per-event MAC path runs;
//     when nil, the report records "key-bound verification skipped".
//   - StopOnFirstFailure controls whether Verify continues recording steps
//     after a failure. The CLI defaults this to true so failure surfaces
//     early; tests sometimes set it false to inspect every step.
type Plan struct {
	SealingKey         ed25519.PublicKey
	MasterIKM          []byte
	StopOnFirstFailure bool
}

// Verify orchestrates the four core checks and writes a Result. A nil
// returned error means structural verification passed; the per-event MAC
// path is reported on Result.MACPass and is independent.
func Verify(led *Ledger, plan Plan) (*Result, error) {
	if led == nil {
		return nil, errors.New("nil ledger")
	}
	if plan.SealingKey == nil {
		return nil, errors.New("plan.SealingKey is required")
	}

	r := &Result{
		EntryCount: len(led.Entries),
		SealDate:   led.Seal.SealDate,
		SignerKey:  led.Seal.SigningKeyFingerprint,
	}

	if err := step(r, "chain-linkage", plan.StopOnFirstFailure, func() error {
		_, e := CheckChain(led.Entries)
		return e
	}); err != nil {
		return r, err
	}
	if err := step(r, "merkle-root", plan.StopOnFirstFailure, func() error {
		return CheckMerkleRoot(led.Entries, led.Seal.MerkleRoot)
	}); err != nil {
		return r, err
	}
	if err := step(r, "seal-signature", plan.StopOnFirstFailure, func() error {
		return CheckSealSignature(&led.Seal, plan.SealingKey)
	}); err != nil {
		return r, err
	}
	r.StructuralPass = lastN(r.Steps, 3) == 3

	if len(plan.MasterIKM) == 0 {
		r.Steps = append(r.Steps, Step{
			Name: "per-event-mac",
			OK:   false,
			Note: "skipped: no --master-key supplied (structural-only verification)",
		})
		return r, nil
	}
	if err := step(r, "per-event-mac", plan.StopOnFirstFailure, func() error {
		_, e := CheckChainMACs(led.Entries, plan.MasterIKM)
		return e
	}); err != nil {
		return r, err
	}
	r.MACPass = true
	return r, nil
}

// step runs fn, records the outcome, and returns the error if stopOnFail
// is true. It keeps each named step in the report so an examiner sees a
// predictable list rather than a single opaque success/failure.
func step(r *Result, name string, stopOnFail bool, fn func() error) error {
	err := fn()
	if err == nil {
		r.Steps = append(r.Steps, Step{Name: name, OK: true})
		return nil
	}
	r.Steps = append(r.Steps, Step{Name: name, OK: false, Note: err.Error()})
	if stopOnFail {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

// lastN returns how many of the trailing N steps in s passed (OK==true).
// Used to roll up structural-only success without scanning the whole slice.
func lastN(s []Step, n int) int {
	if len(s) < n {
		return 0
	}
	count := 0
	for i := len(s) - n; i < len(s); i++ {
		if s[i].OK {
			count++
		}
	}
	return count
}
