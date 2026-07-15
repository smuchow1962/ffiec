package verify

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
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

	// AdditionalVerifications records the §14.6/§14.7/§14.8 attribute
	// families that were present and validated on chain entries. The
	// entries are deterministic: one per distinct attribute family found
	// across all chain entries, in family-name order.
	AdditionalVerifications []AdditionalVerification

	// Supervisory is the §14.13 audit.supervisory.* provenance aggregate,
	// present only when the chain carried the family. It is presentation-
	// only: it never affects StructuralPass, MACPass, or any step. A
	// per-regulator profile layer renders it; the integrity verdict does
	// not depend on it.
	Supervisory *SupervisoryContext
}

// AdditionalVerification records one attribute-family validation
// result. The Family field is the spec family name (e.g.,
// "audit.actor", "audit.reasoning", "audit.downstream_action"); OK
// is true when every entry carrying that family passed validation.
type AdditionalVerification struct {
	Family    string
	OK        bool
	EntryHits int    // how many chain entries carried this family
	Note      string // non-empty on failure
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
	} else {
		if err := step(r, "per-event-mac", plan.StopOnFirstFailure, func() error {
			_, e := CheckChainMACs(led.Entries, plan.MasterIKM)
			return e
		}); err != nil {
			return r, err
		}
		r.MACPass = true
	}

	// §14.6/§14.7/§14.8 attribute-family validation. This is additive:
	// it does not gate structural/MAC pass. When attribute families are
	// present on chain entries, the verifier validates them and records
	// the outcome in AdditionalVerifications.
	r.AdditionalVerifications = validateChainAttributes(led.Entries)

	// §10.84 communication principal-preapproval ordering (FINRA Rule
	// 2210). Also additive and non-gating; emitted only when the chain
	// carries a retail communication event. Re-sort so the combined
	// AdditionalVerifications stays in deterministic family-name order.
	if av := validateCommunicationPreapproval(led.Entries); av != nil {
		r.AdditionalVerifications = append(r.AdditionalVerifications, *av)
		sort.Slice(r.AdditionalVerifications, func(i, j int) bool {
			return r.AdditionalVerifications[i].Family < r.AdditionalVerifications[j].Family
		})
	}

	// §14.13 supervisory-context extraction. Presentation-only: it feeds
	// the regulator profile layer and never gates the integrity verdict.
	r.Supervisory = extractSupervisoryContext(led.Entries)
	return r, nil
}

// validateChainAttributes walks every entry, decodes the event payload,
// and validates any §14.6/§14.7/§14.8 attribute families found. Returns
// one AdditionalVerification per distinct family, in family-name order.
func validateChainAttributes(entries []ChainEntry) []AdditionalVerification {
	state := map[string]*attrFamilyState{}

	for i := range entries {
		payload, err := base64.StdEncoding.DecodeString(entries[i].EventPayloadJCS)
		if err != nil {
			continue // payload decode errors are caught by other steps
		}
		attrs, err := ParseEventAttributes(payload)
		if err != nil || attrs == nil {
			continue
		}
		recordFamily(state, "audit.actor", attrs.Actor != nil, ValidateActorAttributes(attrs.Actor))
		recordFamily(state, "audit.reasoning", attrs.Reasoning != nil, ValidateReasoningAttributes(attrs.Reasoning))
		recordFamily(state, "audit.downstream_action", attrs.DownstreamAction != nil, ValidateDownstreamActionAttributes(attrs.DownstreamAction))
	}

	// Deterministic output: sorted by family name.
	families := []string{"audit.actor", "audit.downstream_action", "audit.reasoning"}
	var out []AdditionalVerification
	for _, name := range families {
		s, found := state[name]
		if !found {
			continue
		}
		av := AdditionalVerification{
			Family:    name,
			OK:        s.ok,
			EntryHits: s.hits,
			Note:      s.firstErr,
		}
		out = append(out, av)
	}
	return out
}

// attrFamilyState tracks per-family validation state across all entries.
type attrFamilyState struct {
	hits     int
	ok       bool
	firstErr string
}

// recordFamily updates the per-family state for one chain entry.
func recordFamily(state map[string]*attrFamilyState, name string, present bool, err error) {
	if !present {
		return
	}
	s, exists := state[name]
	if !exists {
		s = &attrFamilyState{ok: true}
		state[name] = s
	}
	s.hits++
	if err != nil {
		s.ok = false
		if s.firstErr == "" {
			s.firstErr = err.Error()
		}
	}
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
