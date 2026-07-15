// Package profile is the per-regulator presentation layer that sits ABOVE
// the §7 integrity verdict. A Profile reframes verifier output for a
// specific supervisory audience — which authority names to render, how to
// label the supervisory context — WITHOUT changing the integrity outcome.
//
// The load-bearing guarantee (per spec §14.13 and the §7 determinism
// contract): a Profile is presentation-only. It reads a completed
// verify.Result and never mutates it; the integrity verdict it renders is
// verify.Result.IntegrityVerdict() verbatim, identical under every
// profile. Two verifiers walking the same chain reach the same verdict
// regardless of which profile framed the report. A profile cannot accept
// or reject a chain construct, cannot gate PASS/FAIL, and cannot change
// the exit code — it only selects vocabulary and framing.
package profile

import "github.com/mmpworks/ffiec/verifier/internal/verify"

// Profile is a per-regulator presentation configuration. It is pure data:
// an identity plus a display vocabulary for supervisory authorities. No
// method on Profile touches the integrity verdict.
type Profile struct {
	ID          string
	DisplayName string
	// supervisorNames maps a stable supervisory-authority identifier
	// (the §14.13 institution-asserted value, e.g. "tx_dob", "fdic") to
	// the examiner-facing display name this profile prefers.
	supervisorNames map[string]string
}

// View is the presentation a Profile produces from a Result. Every field
// is display material; IntegrityVerdict is copied verbatim from the
// Result and is the same under any profile.
type View struct {
	ProfileID        string
	ProfileDisplay   string
	IntegrityVerdict string
	// SupervisoryLines renders the §14.13 supervisory context in this
	// profile's vocabulary. Empty when the chain carried no supervisory
	// family.
	SupervisoryLines []string
}

// federalNames are the federal prudential-supervisor display names every
// profile shares; a regulator profile layers its state-side names on top.
var federalNames = map[string]string{
	"fdic": "FDIC",
	"frb":  "Federal Reserve",
	"occ":  "OCC",
	"ncua": "NCUA",
	"none": "none",
}

// registry holds the known profiles. FFIEC is the default (authority-
// neutral framing); regulator profiles add state-side vocabulary.
var registry = map[string]Profile{
	"ffiec": {
		ID:              "ffiec",
		DisplayName:     "FFIEC (default)",
		supervisorNames: federalNames,
	},
	"tx-dob": {
		ID:          "tx-dob",
		DisplayName: "Texas Department of Banking",
		supervisorNames: mergeNames(federalNames, map[string]string{
			"tx_dob": "Texas Department of Banking",
		}),
	},
}

// Default returns the authority-neutral FFIEC profile.
func Default() Profile { return registry["ffiec"] }

// Lookup returns the profile registered under id and whether it exists.
// An unknown id is a caller error; the CLI reports it rather than
// silently falling back, so an examiner never gets an unexpected framing.
func Lookup(id string) (Profile, bool) {
	p, ok := registry[id]
	return p, ok
}

// IDs returns the registered profile identifiers, for CLI help/usage.
func IDs() []string { return []string{"ffiec", "tx-dob"} }

// supervisorDisplay maps a stable supervisor id to this profile's display
// name, falling back to the raw id so an unrecognized authority still
// renders (never panics, never drops information).
func (p Profile) supervisorDisplay(id string) string {
	if id == "" {
		return ""
	}
	if name, ok := p.supervisorNames[id]; ok {
		return name
	}
	return id
}

// Render builds the presentation View from a completed Result. It reads
// r read-only and copies the integrity verdict verbatim — this is the
// presentation-only guarantee in code. When r carries no supervisory
// context, SupervisoryLines is empty and the View differs from another
// profile's only in the profile identity fields.
func (p Profile) Render(r *verify.Result) View {
	v := View{
		ProfileID:      p.ID,
		ProfileDisplay: p.DisplayName,
	}
	if r != nil {
		v.IntegrityVerdict = r.IntegrityVerdict()
		v.SupervisoryLines = p.supervisoryLines(r.Supervisory)
	}
	return v
}

// supervisoryLines renders the §14.13 supervisory context in this
// profile's vocabulary. Returns nil when there is no context.
func (p Profile) supervisoryLines(ctx *verify.SupervisoryContext) []string {
	if ctx == nil {
		return nil
	}
	var lines []string
	if ctx.CharterType != "" {
		lines = append(lines, "charter_type: "+ctx.CharterType)
	}
	if ctx.PrimaryStateSupervisor != "" {
		lines = append(lines, "state_supervisor: "+p.supervisorDisplay(ctx.PrimaryStateSupervisor))
	}
	if ctx.FederalPrudentialSupervisor != "" {
		lines = append(lines, "federal_supervisor: "+p.supervisorDisplay(ctx.FederalPrudentialSupervisor))
	}
	if ctx.DualSupervision {
		lines = append(lines, "dual_supervision: yes")
	}
	return lines
}

// mergeNames returns a new map combining base and extra (extra wins on
// key collisions). Used to layer a regulator profile's state-side names
// over the shared federal names without mutating the shared map.
func mergeNames(base, extra map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(extra))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
