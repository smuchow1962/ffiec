package profile

import (
	"reflect"
	"testing"

	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

func TestLookup_KnownAndUnknown(t *testing.T) {
	if _, ok := Lookup("ffiec"); !ok {
		t.Error("ffiec profile should be registered")
	}
	if _, ok := Lookup("tx-dob"); !ok {
		t.Error("tx-dob profile should be registered")
	}
	if _, ok := Lookup("nonesuch"); ok {
		t.Error("unknown profile should not resolve (CLI must report, not silently fall back)")
	}
	if Default().ID != "ffiec" {
		t.Errorf("default profile = %q, want ffiec", Default().ID)
	}
}

// The load-bearing guarantee: a profile is presentation-only. For the same
// Result, every profile renders the SAME integrity verdict — equal to the
// Result's own IntegrityVerdict — and Render never mutates the Result.
func TestRender_IntegrityVerdictIsProfileInvariant(t *testing.T) {
	cases := []*verify.Result{
		{StructuralPass: true},
		{StructuralPass: true, MACPass: true},
		{StructuralPass: false},
		{StructuralPass: true, Supervisory: &verify.SupervisoryContext{
			CharterType:                 "state_bank",
			PrimaryStateSupervisor:      "tx_dob",
			FederalPrudentialSupervisor: "fdic",
			DualSupervision:             true,
		}},
	}

	ffiec, _ := Lookup("ffiec")
	txdob, _ := Lookup("tx-dob")

	for _, r := range cases {
		snapshot := *r // copy the integrity-bearing fields before rendering

		vF := ffiec.Render(r)
		vT := txdob.Render(r)

		// F1: verdict is identical across profiles and equals the Result's.
		if vF.IntegrityVerdict != r.IntegrityVerdict() || vT.IntegrityVerdict != r.IntegrityVerdict() {
			t.Errorf("verdict diverged from Result: ffiec=%q txdob=%q result=%q",
				vF.IntegrityVerdict, vT.IntegrityVerdict, r.IntegrityVerdict())
		}
		if vF.IntegrityVerdict != vT.IntegrityVerdict {
			t.Errorf("verdict differs across profiles: ffiec=%q txdob=%q", vF.IntegrityVerdict, vT.IntegrityVerdict)
		}

		// F1: Render must not mutate the integrity-bearing fields.
		if r.StructuralPass != snapshot.StructuralPass || r.MACPass != snapshot.MACPass {
			t.Error("Render mutated the Result's integrity fields")
		}
		if !reflect.DeepEqual(r.Supervisory, snapshot.Supervisory) {
			t.Error("Render mutated the Result's supervisory context")
		}
	}
}

// A profile changes only VOCABULARY, not verdict. tx-dob renders the
// Texas state supervisor by its full name; ffiec renders the raw stable
// id. Same underlying data, different framing.
func TestRender_ProfileVocabularyDiffers(t *testing.T) {
	r := &verify.Result{
		StructuralPass: true,
		Supervisory: &verify.SupervisoryContext{
			CharterType:            "state_bank",
			PrimaryStateSupervisor: "tx_dob",
		},
	}
	ffiec, _ := Lookup("ffiec")
	txdob, _ := Lookup("tx-dob")

	fLines := ffiec.Render(r).SupervisoryLines
	tLines := txdob.Render(r).SupervisoryLines

	if !containsLine(tLines, "state_supervisor: Texas Department of Banking") {
		t.Errorf("tx-dob should render the full authority name, got %v", tLines)
	}
	if !containsLine(fLines, "state_supervisor: tx_dob") {
		t.Errorf("ffiec should render the raw stable id, got %v", fLines)
	}
}

func TestRender_NoSupervisoryContextYieldsNoLines(t *testing.T) {
	r := &verify.Result{StructuralPass: true}
	txdob, _ := Lookup("tx-dob")
	if lines := txdob.Render(r).SupervisoryLines; len(lines) != 0 {
		t.Errorf("expected no supervisory lines, got %v", lines)
	}
}

func containsLine(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}
