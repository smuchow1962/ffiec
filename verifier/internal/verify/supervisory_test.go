package verify

import (
	"encoding/base64"
	"testing"
)

// entryWithPayload builds a ChainEntry whose EventPayloadJCS carries the
// given JSON object (base64 of the raw bytes, matching the wire form
// ParseEventAttributes decodes).
func entryWithPayload(jsonPayload string) ChainEntry {
	return ChainEntry{
		EventPayloadJCS: base64.StdEncoding.EncodeToString([]byte(jsonPayload)),
	}
}

func TestExtractSupervisoryContext_AbsentWhenNoFamily(t *testing.T) {
	// A plain audit event with no audit.supervisory.* yields no context.
	entries := []ChainEntry{
		entryWithPayload(`{"audit.actor":{"authenticated_user_id_hash":"` +
			"aa" + `","authentication_method":"oidc_sso"}}`),
		entryWithPayload(`{"name":"plain-event"}`),
	}
	if ctx := extractSupervisoryContext(entries); ctx != nil {
		t.Fatalf("expected nil context, got %+v", ctx)
	}
}

func TestExtractSupervisoryContext_AggregatesFields(t *testing.T) {
	entries := []ChainEntry{
		entryWithPayload(`{"audit.supervisory":{"charter_type":"state_bank",` +
			`"primary_state_supervisor":"tx_dob",` +
			`"federal_prudential_supervisor":"fdic",` +
			`"dual_supervision":true}}`),
		entryWithPayload(`{"name":"plain-event"}`),
	}
	ctx := extractSupervisoryContext(entries)
	if ctx == nil {
		t.Fatal("expected a supervisory context, got nil")
	}
	if ctx.CharterType != "state_bank" {
		t.Errorf("charter_type = %q, want state_bank", ctx.CharterType)
	}
	if ctx.PrimaryStateSupervisor != "tx_dob" {
		t.Errorf("primary_state_supervisor = %q, want tx_dob", ctx.PrimaryStateSupervisor)
	}
	if ctx.FederalPrudentialSupervisor != "fdic" {
		t.Errorf("federal_prudential_supervisor = %q, want fdic", ctx.FederalPrudentialSupervisor)
	}
	if !ctx.DualSupervision {
		t.Error("dual_supervision = false, want true")
	}
	if ctx.EntryHits != 1 {
		t.Errorf("EntryHits = %d, want 1", ctx.EntryHits)
	}
}

// Supervisory context is presentation-only: its presence or absence must
// never change the integrity verdict. A clean structural result reports
// the same verdict whether or not a supervisory family was carried.
func TestIntegrityVerdict_IndependentOfSupervisory(t *testing.T) {
	base := &Result{StructuralPass: true}
	if got := base.IntegrityVerdict(); got != "PASS (structural)" {
		t.Errorf("structural verdict = %q", got)
	}

	withCtx := &Result{StructuralPass: true, Supervisory: &SupervisoryContext{CharterType: "state_bank"}}
	if got := withCtx.IntegrityVerdict(); got != "PASS (structural)" {
		t.Errorf("verdict changed by supervisory context: %q", got)
	}

	full := &Result{StructuralPass: true, MACPass: true}
	if got := full.IntegrityVerdict(); got != "PASS (full, key-bound)" {
		t.Errorf("full verdict = %q", got)
	}

	fail := &Result{StructuralPass: false}
	if got := fail.IntegrityVerdict(); got != "FAIL" {
		t.Errorf("fail verdict = %q", got)
	}
}
