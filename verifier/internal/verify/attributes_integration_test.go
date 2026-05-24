package verify

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// makeFixtureWithPayloads is like makeFixture but lets the caller control
// the event payload JSON for each entry. The payloads slice length must
// match entryCount.
func makeFixtureWithPayloads(t *testing.T, payloads []json.RawMessage) (led *Ledger, sealingPub ed25519.PublicKey, ikm []byte) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ikm = bytes.Repeat([]byte{0xAB}, 32)
	fb := &FixtureBuilder{
		TenantBindingKDFLabel: "tenant=acme;env=prod",
		SealDate:              "2026-06-15",
		MasterIKM:             ikm,
		SealingPriv:           priv,
	}
	var entries []ChainEntry
	prev := ""
	for i, p := range payloads {
		e, err := fb.AppendEvent(prev, i, fmt.Sprintf("e_%04d", i), p)
		if err != nil {
			t.Fatal(err)
		}
		entries = append(entries, e)
		prev = e.EntryHash
	}
	seal, err := fb.SealEntries(entries)
	if err != nil {
		t.Fatal(err)
	}
	return &Ledger{Entries: entries, Seal: seal}, pub, ikm
}

// validActorPayload returns a JSON payload with a valid §14.6 actor family.
func validActorPayload() json.RawMessage {
	hash := strings.Repeat("ab", 32)
	return json.RawMessage(fmt.Sprintf(`{"event":"decision","audit.actor":{"authenticated_user_id_hash":"%s","authentication_method":"saml_sso","session_id":"sess-001"}}`, hash))
}

// validReasoningPayload returns a JSON payload with a valid §14.7 reasoning family.
func validReasoningPayload() json.RawMessage {
	return json.RawMessage(`{"event":"decision","audit.reasoning":{"substrate_kind":"neurosymbolic"}}`)
}

// validDownstreamPayload returns a JSON payload with a valid §14.8 downstream action family.
func validDownstreamPayload() json.RawMessage {
	hash := strings.Repeat("cd", 32)
	return json.RawMessage(fmt.Sprintf(`{"event":"action","audit.downstream_action":{"action_kind":"payment_executed","system_of_record_id":"core-banking","change_record_id_hash":"%s","applied_at_utc":"2026-05-23T14:00:00Z"}}`, hash))
}

// plainPayload returns a payload with no attribute families.
func plainPayload(seq int) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{"event":"decision","seq":%d}`, seq))
}

// -- integration tests: pipeline with attribute families ------------------

func TestVerify_NoAttributeFamilies_EmptyAdditionalVerifications(t *testing.T) {
	payloads := []json.RawMessage{
		plainPayload(0),
		plainPayload(1),
		plainPayload(2),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if len(r.AdditionalVerifications) != 0 {
		t.Errorf("expected 0 additional verifications, got %d", len(r.AdditionalVerifications))
	}
}

func TestVerify_ValidActorFamily_RecordedInAdditionalVerifications(t *testing.T) {
	payloads := []json.RawMessage{
		validActorPayload(),
		plainPayload(1),
		validActorPayload(),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if len(r.AdditionalVerifications) != 1 {
		t.Fatalf("expected 1 additional verification, got %d", len(r.AdditionalVerifications))
	}
	av := r.AdditionalVerifications[0]
	if av.Family != "audit.actor" {
		t.Errorf("family = %q, want %q", av.Family, "audit.actor")
	}
	if !av.OK {
		t.Errorf("expected OK=true, got false; note: %s", av.Note)
	}
	if av.EntryHits != 2 {
		t.Errorf("entry_hits = %d, want 2", av.EntryHits)
	}
}

func TestVerify_ValidReasoningFamily_RecordedInAdditionalVerifications(t *testing.T) {
	payloads := []json.RawMessage{
		validReasoningPayload(),
		plainPayload(1),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	found := false
	for _, av := range r.AdditionalVerifications {
		if av.Family == "audit.reasoning" {
			found = true
			if !av.OK {
				t.Errorf("reasoning family failed: %s", av.Note)
			}
			if av.EntryHits != 1 {
				t.Errorf("entry_hits = %d, want 1", av.EntryHits)
			}
		}
	}
	if !found {
		t.Error("audit.reasoning not found in additional verifications")
	}
}

func TestVerify_ValidDownstreamActionFamily_RecordedInAdditionalVerifications(t *testing.T) {
	payloads := []json.RawMessage{
		validDownstreamPayload(),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	found := false
	for _, av := range r.AdditionalVerifications {
		if av.Family == "audit.downstream_action" {
			found = true
			if !av.OK {
				t.Errorf("downstream action family failed: %s", av.Note)
			}
		}
	}
	if !found {
		t.Error("audit.downstream_action not found in additional verifications")
	}
}

func TestVerify_AllThreeFamilies_OrderedByFamilyName(t *testing.T) {
	payloads := []json.RawMessage{
		validActorPayload(),
		validReasoningPayload(),
		validDownstreamPayload(),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if len(r.AdditionalVerifications) != 3 {
		t.Fatalf("expected 3 additional verifications, got %d", len(r.AdditionalVerifications))
	}
	// Spec requires deterministic order: alphabetical by family name.
	wantOrder := []string{"audit.actor", "audit.downstream_action", "audit.reasoning"}
	for i, want := range wantOrder {
		if r.AdditionalVerifications[i].Family != want {
			t.Errorf("position %d: family = %q, want %q", i, r.AdditionalVerifications[i].Family, want)
		}
	}
}

func TestVerify_InvalidActorFamily_RecordedAsFailure(t *testing.T) {
	badPayload := json.RawMessage(`{"event":"decision","audit.actor":{"authenticated_user_id_hash":"bad","authentication_method":"saml_sso"}}`)
	payloads := []json.RawMessage{
		badPayload,
		plainPayload(1),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	// Structural and MAC pass should still be true — attribute validation
	// is additive and does not gate the core checks.
	if !r.StructuralPass {
		t.Error("StructuralPass should be true even with bad attributes")
	}
	if !r.MACPass {
		t.Error("MACPass should be true even with bad attributes")
	}
	if len(r.AdditionalVerifications) != 1 {
		t.Fatalf("expected 1 additional verification, got %d", len(r.AdditionalVerifications))
	}
	av := r.AdditionalVerifications[0]
	if av.OK {
		t.Error("expected OK=false for invalid actor attributes")
	}
	if av.Note == "" {
		t.Error("expected non-empty note on failure")
	}
}

func TestVerify_InvalidReasoningFamily_RecordedAsFailure(t *testing.T) {
	badPayload := json.RawMessage(`{"event":"decision","audit.reasoning":{"substrate_kind":"magic_crystal_ball"}}`)
	payloads := []json.RawMessage{badPayload}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	found := false
	for _, av := range r.AdditionalVerifications {
		if av.Family == "audit.reasoning" {
			found = true
			if av.OK {
				t.Error("expected OK=false for invalid reasoning")
			}
			if !strings.Contains(av.Note, "not in the closed enum") {
				t.Errorf("unexpected note: %s", av.Note)
			}
		}
	}
	if !found {
		t.Error("audit.reasoning not found in additional verifications")
	}
}

func TestVerify_InvalidDownstreamAction_RecordedAsFailure(t *testing.T) {
	badPayload := json.RawMessage(`{"event":"action","audit.downstream_action":{"action_kind":"record_updated","system_of_record_id":"sys","change_record_id_hash":"not-a-hash","applied_at_utc":"2026-05-23T14:00:00Z"}}`)
	payloads := []json.RawMessage{badPayload}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	found := false
	for _, av := range r.AdditionalVerifications {
		if av.Family == "audit.downstream_action" {
			found = true
			if av.OK {
				t.Error("expected OK=false for invalid downstream action")
			}
		}
	}
	if !found {
		t.Error("audit.downstream_action not found in additional verifications")
	}
}

func TestVerify_MixedValidAndInvalidEntries_FirstErrorCaptured(t *testing.T) {
	// Two entries with actor family: first valid, second invalid.
	goodPayload := validActorPayload()
	badPayload := json.RawMessage(`{"event":"decision","audit.actor":{"authenticated_user_id_hash":"bad","authentication_method":"saml_sso"}}`)
	payloads := []json.RawMessage{goodPayload, badPayload}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if len(r.AdditionalVerifications) != 1 {
		t.Fatalf("expected 1 additional verification, got %d", len(r.AdditionalVerifications))
	}
	av := r.AdditionalVerifications[0]
	if av.OK {
		t.Error("expected failure when one entry has invalid attributes")
	}
	if av.EntryHits != 2 {
		t.Errorf("entry_hits = %d, want 2 (both entries emitted the family)", av.EntryHits)
	}
}

func TestVerify_StructuralOnlyMode_StillValidatesAttributes(t *testing.T) {
	payloads := []json.RawMessage{
		validActorPayload(),
		validReasoningPayload(),
	}
	led, pub, _ := makeFixtureWithPayloads(t, payloads)
	// No MasterIKM → structural-only mode
	r, err := Verify(led, Plan{SealingKey: pub, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !r.StructuralPass {
		t.Error("StructuralPass should be true")
	}
	if len(r.AdditionalVerifications) != 2 {
		t.Errorf("expected 2 additional verifications in structural-only mode, got %d", len(r.AdditionalVerifications))
	}
}
