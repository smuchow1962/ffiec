package verify

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// §10.84 event-payload builders. The communication event carries the
// audience classifier; the approval is a §10.50 review with role
// registered_principal (flat dotted keys); the send is a §14.8
// downstream_action (nested object) with §4.4 top-level parent linkage.

func retailCommunicationPayload(runID string, seq int) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(
		`{"run_id":%q,"seq":%d,"audit.communication.audience":"retail"}`, runID, seq))
}

func institutionalCommunicationPayload(runID string, seq int) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(
		`{"run_id":%q,"seq":%d,"audit.communication.audience":"institutional"}`, runID, seq))
}

func principalApprovalPayload(runID string, seq int, parentRunID string, parentSeq int, signedAtUTC string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(
		`{"run_id":%q,"seq":%d,"audit.review.role":"registered_principal","audit.review.parent_run_id":%q,"audit.review.parent_seq":%d,"audit.review.signed_review":{"audit.signed_review.signed_at_utc":%q}}`,
		runID, seq, parentRunID, parentSeq, signedAtUTC))
}

func communicationSendPayload(runID string, seq int, parentRunID string, parentSeq int, appliedAtUTC string) json.RawMessage {
	hash := strings.Repeat("cd", 32)
	return json.RawMessage(fmt.Sprintf(
		`{"run_id":%q,"seq":%d,"parent_run_id":%q,"parent_seq":%d,"audit.downstream_action":{"action_kind":"communication_sent","system_of_record_id":"comms-platform","change_record_id_hash":%q,"applied_at_utc":%q}}`,
		runID, seq, parentRunID, parentSeq, hash, appliedAtUTC))
}

func findAV(r *Result, family string) (AdditionalVerification, bool) {
	for _, av := range r.AdditionalVerifications {
		if av.Family == family {
			return av, true
		}
	}
	return AdditionalVerification{}, false
}

// -- §10.84 behavioral tests ----------------------------------------------

func TestPreapproval_ApprovalBeforeSend_MarkerVerified(t *testing.T) {
	payloads := []json.RawMessage{
		retailCommunicationPayload("comm-1", 1),
		principalApprovalPayload("appr-1", 2, "comm-1", 1, "2026-06-01T10:00:00Z"),
		communicationSendPayload("send-1", 3, "comm-1", 1, "2026-06-01T12:00:00Z"),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	av, found := findAV(r, preapprovalFamily)
	if !found {
		t.Fatalf("expected %s in additional verifications", preapprovalFamily)
	}
	if !av.OK {
		t.Errorf("expected OK=true (approval preceded send), got note: %s", av.Note)
	}
	if av.EntryHits != 1 {
		t.Errorf("entry_hits = %d, want 1", av.EntryHits)
	}
	// Soft-enforcement: the core checks are unaffected.
	if !r.StructuralPass || !r.MACPass {
		t.Error("§10.84 is additive and must not gate structural/MAC pass")
	}
}

func TestPreapproval_ApprovalEqualsSend_MarkerVerified(t *testing.T) {
	// approval == send is conformant: the rule is approval <= send.
	payloads := []json.RawMessage{
		retailCommunicationPayload("comm-1", 1),
		principalApprovalPayload("appr-1", 2, "comm-1", 1, "2026-06-01T12:00:00Z"),
		communicationSendPayload("send-1", 3, "comm-1", 1, "2026-06-01T12:00:00Z"),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	av, found := findAV(r, preapprovalFamily)
	if !found || !av.OK {
		t.Errorf("approval == send must be conformant; found=%v note=%s", found, av.Note)
	}
}

func TestPreapproval_ApprovalAfterSend_Anomaly(t *testing.T) {
	payloads := []json.RawMessage{
		retailCommunicationPayload("comm-1", 1),
		principalApprovalPayload("appr-1", 2, "comm-1", 1, "2026-06-01T14:00:00Z"),
		communicationSendPayload("send-1", 3, "comm-1", 1, "2026-06-01T12:00:00Z"),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	av, found := findAV(r, preapprovalFamily)
	if !found {
		t.Fatalf("expected %s in additional verifications", preapprovalFamily)
	}
	if av.OK {
		t.Error("expected OK=false: approval did not precede send")
	}
	if !strings.Contains(av.Note, "approval did not precede send at seq 3") {
		t.Errorf("unexpected anomaly note: %q", av.Note)
	}
	// Anomaly is soft-enforcement: the chain still passes integrity.
	if !r.StructuralPass || !r.MACPass {
		t.Error("§10.84 anomaly must not gate structural/MAC pass")
	}
}

func TestPreapproval_RetailSendWithoutApproval_Anomaly(t *testing.T) {
	payloads := []json.RawMessage{
		retailCommunicationPayload("comm-1", 1),
		communicationSendPayload("send-1", 2, "comm-1", 1, "2026-06-01T12:00:00Z"),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	av, found := findAV(r, preapprovalFamily)
	if !found || av.OK {
		t.Fatalf("expected an anomaly for a retail send without approval; found=%v ok=%v", found, av.OK)
	}
	if !strings.Contains(av.Note, "no registered-principal approval bound at seq 2") {
		t.Errorf("unexpected anomaly note: %q", av.Note)
	}
}

func TestPreapproval_InstitutionalAudience_NotChecked(t *testing.T) {
	// Institutional communications are outside Rule 2210 pre-approval;
	// even with approval after send, no §10.84 anomaly is emitted.
	payloads := []json.RawMessage{
		institutionalCommunicationPayload("comm-1", 1),
		principalApprovalPayload("appr-1", 2, "comm-1", 1, "2026-06-01T14:00:00Z"),
		communicationSendPayload("send-1", 3, "comm-1", 1, "2026-06-01T12:00:00Z"),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if _, found := findAV(r, preapprovalFamily); found {
		t.Error("§10.84 must not emit for a non-retail audience")
	}
}

func TestPreapproval_NoCommunicationEvent_FamilyAbsent(t *testing.T) {
	payloads := []json.RawMessage{
		plainPayload(0),
		validActorPayload(),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if _, found := findAV(r, preapprovalFamily); found {
		t.Error("§10.84 must not emit when no communication event is present")
	}
}

func TestPreapproval_PendingSend_NoAnomaly(t *testing.T) {
	// A retail communication approved but not yet sent is a normal
	// pending state, not an anomaly.
	payloads := []json.RawMessage{
		retailCommunicationPayload("comm-1", 1),
		principalApprovalPayload("appr-1", 2, "comm-1", 1, "2026-06-01T10:00:00Z"),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	av, found := findAV(r, preapprovalFamily)
	if !found {
		t.Fatalf("expected %s (retail communication present)", preapprovalFamily)
	}
	if !av.OK {
		t.Errorf("a not-yet-sent retail communication must not be an anomaly; note: %s", av.Note)
	}
}

func TestPreapproval_DeterministicFamilyOrder(t *testing.T) {
	// With a §14.8 downstream_action (the send) AND the §10.84 family,
	// the combined AdditionalVerifications stay sorted by family name.
	payloads := []json.RawMessage{
		retailCommunicationPayload("comm-1", 1),
		principalApprovalPayload("appr-1", 2, "comm-1", 1, "2026-06-01T10:00:00Z"),
		communicationSendPayload("send-1", 3, "comm-1", 1, "2026-06-01T12:00:00Z"),
	}
	led, pub, ikm := makeFixtureWithPayloads(t, payloads)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	// Expect exactly: audit.communication_preapproval, audit.downstream_action.
	want := []string{"audit.communication_preapproval", "audit.downstream_action"}
	if len(r.AdditionalVerifications) != len(want) {
		t.Fatalf("got %d additional verifications, want %d", len(r.AdditionalVerifications), len(want))
	}
	for i, w := range want {
		if r.AdditionalVerifications[i].Family != w {
			t.Errorf("position %d: family = %q, want %q", i, r.AdditionalVerifications[i].Family, w)
		}
	}
}
