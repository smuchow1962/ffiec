package verify

import "encoding/base64"

// SupervisoryContext is the aggregated §14.13 audit.supervisory.*
// provenance for a tenant-day, extracted for examiner-facing rendering.
// It is presentation-only: nothing in this file gates StructuralPass,
// MACPass, or any §7 step. The values are institution-asserted (§1.2) —
// the verifier surfaces what the institution recorded, and does not
// evaluate whether the asserted authority is correct.
//
// When the chain carries no supervisory family, the extractor returns
// nil and the report renders no supervisory block.
type SupervisoryContext struct {
	CharterType                 string
	PrimaryStateSupervisor      string
	FederalPrudentialSupervisor string
	DualSupervision             bool
	EntryHits                   int // how many entries carried the family
}

// extractSupervisoryContext walks every entry and folds any §14.13
// supervisory family into a single presentation aggregate. It takes the
// first non-empty value seen for each field — a tenant-day is expected to
// carry a consistent supervisory context, and divergence is an
// institution-side concern the SOC engagement reviews, not a verifier
// integrity check. Returns nil when no entry carries the family.
func extractSupervisoryContext(entries []ChainEntry) *SupervisoryContext {
	var ctx *SupervisoryContext
	for i := range entries {
		payload, err := base64.StdEncoding.DecodeString(entries[i].EventPayloadJCS)
		if err != nil {
			continue // payload decode errors surface in the integrity steps
		}
		attrs, err := ParseEventAttributes(payload)
		if err != nil || attrs == nil || attrs.Supervisory == nil {
			continue
		}
		if ctx == nil {
			ctx = &SupervisoryContext{}
		}
		ctx.EntryHits++
		fillFirst(&ctx.CharterType, attrs.Supervisory.CharterType)
		fillFirst(&ctx.PrimaryStateSupervisor, attrs.Supervisory.PrimaryStateSupervisor)
		fillFirst(&ctx.FederalPrudentialSupervisor, attrs.Supervisory.FederalPrudentialSupervisor)
		if attrs.Supervisory.DualSupervision {
			ctx.DualSupervision = true
		}
	}
	return ctx
}

// fillFirst copies src into *dst only when *dst is still empty, so the
// aggregate keeps the first non-empty value seen across the tenant-day.
func fillFirst(dst *string, src string) {
	if *dst == "" && src != "" {
		*dst = src
	}
}

// IntegrityVerdict is the single canonical rendering of the run's
// integrity outcome. It is the ONLY place the PASS/FAIL verdict string is
// produced, so every consumer — the default report writer and any
// regulator profile layer above it — reads the same verdict rather than
// recomputing it. That is what keeps a per-regulator profile presentation-
// only: a profile can reframe the output around this string but cannot
// change it (§14.13, §7 determinism contract).
func (r *Result) IntegrityVerdict() string {
	if r == nil || !r.StructuralPass {
		return "FAIL"
	}
	if r.MACPass {
		return "PASS (full, key-bound)"
	}
	return "PASS (structural)"
}
