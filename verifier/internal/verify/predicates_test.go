package verify

import (
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/core/jcs"
)

// -- helpers ---------------------------------------------------------------

// validSHA256Hex returns a well-formed 64-char lowercase hex string.
func validSHA256Hex() string {
	return strings.Repeat("ab", 32)
}

// -- isWellFormedSHA256Hex -------------------------------------------------

func TestIsWellFormedSHA256Hex(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"valid 64-char hex", validSHA256Hex(), true},
		{"valid all-zeros", strings.Repeat("00", 32), true},
		{"too short", strings.Repeat("ab", 31), false},
		{"too long", strings.Repeat("ab", 33), false},
		{"empty", "", false},
		{"uppercase rejected", strings.Repeat("AB", 32), true}, // hex.DecodeString accepts uppercase
		{"non-hex chars", strings.Repeat("zz", 32), false},
		{"64 chars but not hex", "g" + strings.Repeat("a", 63), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isWellFormedSHA256Hex(tc.in); got != tc.want {
				t.Errorf("isWellFormedSHA256Hex(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// -- isValidRFC3339UTC -----------------------------------------------------

func TestIsValidRFC3339UTC(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"valid UTC", "2026-05-23T14:30:00Z", true},
		{"valid UTC with subseconds", "2026-05-23T14:30:00.123Z", true},
		{"non-UTC offset rejected", "2026-05-23T14:30:00+05:00", false},
		{"non-UTC negative offset", "2026-05-23T14:30:00-07:00", false},
		{"not RFC 3339 at all", "May 23, 2026", false},
		{"empty", "", false},
		{"date only", "2026-05-23", false},
		{"missing Z suffix with zero offset", "2026-05-23T14:30:00+00:00", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidRFC3339UTC(tc.in); got != tc.want {
				t.Errorf("isValidRFC3339UTC(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// -- §14.6 ValidateActorAttributes ----------------------------------------

func TestValidateActorAttributes_Nil(t *testing.T) {
	if err := ValidateActorAttributes(nil); err != nil {
		t.Errorf("nil actor should pass: %v", err)
	}
}

func TestValidateActorAttributes_Valid(t *testing.T) {
	a := &ActorAttributes{
		AuthenticatedUserIDHash: validSHA256Hex(),
		AuthenticationMethod:    "saml_sso",
	}
	if err := ValidateActorAttributes(a); err != nil {
		t.Errorf("valid actor rejected: %v", err)
	}
}

func TestValidateActorAttributes_ValidWithSessionAndDelegation(t *testing.T) {
	a := &ActorAttributes{
		AuthenticatedUserIDHash: validSHA256Hex(),
		AuthenticationMethod:    "oidc_sso",
		SessionID:               "session-12345",
		DelegationChain: []DelegationChainEntry{
			{AuthenticatedUserIDHash: strings.Repeat("00", 32), AuthenticationMethod: "api_key_with_iam"},
			{AuthenticatedUserIDHash: strings.Repeat("ff", 32), AuthenticationMethod: "saml_sso"},
		},
	}
	if err := ValidateActorAttributes(a); err != nil {
		t.Errorf("valid actor with delegation rejected: %v", err)
	}
}

func TestValidateActorAttributes_MissingHash(t *testing.T) {
	a := &ActorAttributes{
		AuthenticationMethod: "saml_sso",
	}
	err := ValidateActorAttributes(a)
	if err == nil {
		t.Fatal("expected error for missing hash")
	}
	if !strings.Contains(err.Error(), "authenticated_user_id_hash is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateActorAttributes_MalformedHash(t *testing.T) {
	a := &ActorAttributes{
		AuthenticatedUserIDHash: "not-a-hash",
		AuthenticationMethod:    "saml_sso",
	}
	err := ValidateActorAttributes(a)
	if err == nil {
		t.Fatal("expected error for malformed hash")
	}
	if !strings.Contains(err.Error(), "not well-formed SHA-256 hex") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateActorAttributes_MissingMethod(t *testing.T) {
	a := &ActorAttributes{
		AuthenticatedUserIDHash: validSHA256Hex(),
	}
	err := ValidateActorAttributes(a)
	if err == nil {
		t.Fatal("expected error for missing method")
	}
	if !strings.Contains(err.Error(), "authentication_method is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateActorAttributes_InvalidMethod(t *testing.T) {
	a := &ActorAttributes{
		AuthenticatedUserIDHash: validSHA256Hex(),
		AuthenticationMethod:    "password_login",
	}
	err := ValidateActorAttributes(a)
	if err == nil {
		t.Fatal("expected error for invalid method")
	}
	if !strings.Contains(err.Error(), "not in the closed enum") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateActorAttributes_AllEnumValues(t *testing.T) {
	methods := []string{
		"saml_sso", "oidc_sso", "mtls_workload", "spiffe_workload",
		"hsm_bearer_token", "api_key_with_iam", "institution_named",
	}
	for _, m := range methods {
		t.Run(m, func(t *testing.T) {
			a := &ActorAttributes{
				AuthenticatedUserIDHash: validSHA256Hex(),
				AuthenticationMethod:    m,
			}
			if err := ValidateActorAttributes(a); err != nil {
				t.Errorf("valid enum value %q rejected: %v", m, err)
			}
		})
	}
}

func TestValidateActorAttributes_DelegationChainUnsorted(t *testing.T) {
	// "ff..." sorts after "00..." in JCS canonical form, so reversed = unsorted.
	a := &ActorAttributes{
		AuthenticatedUserIDHash: validSHA256Hex(),
		AuthenticationMethod:    "saml_sso",
		DelegationChain: []DelegationChainEntry{
			{AuthenticatedUserIDHash: strings.Repeat("ff", 32), AuthenticationMethod: "saml_sso"},
			{AuthenticatedUserIDHash: strings.Repeat("00", 32), AuthenticationMethod: "saml_sso"},
		},
	}
	err := ValidateActorAttributes(a)
	if err == nil {
		t.Fatal("expected error for unsorted delegation chain")
	}
	if !strings.Contains(err.Error(), "not JCS-canonical lex-sorted") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateActorAttributes_DelegationChainBadHash(t *testing.T) {
	a := &ActorAttributes{
		AuthenticatedUserIDHash: validSHA256Hex(),
		AuthenticationMethod:    "saml_sso",
		DelegationChain: []DelegationChainEntry{
			{AuthenticatedUserIDHash: "bad", AuthenticationMethod: "saml_sso"},
		},
	}
	err := ValidateActorAttributes(a)
	if err == nil {
		t.Fatal("expected error for bad delegation hash")
	}
	if !strings.Contains(err.Error(), "delegation_chain[0]") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateActorAttributes_DelegationChainBadMethod(t *testing.T) {
	a := &ActorAttributes{
		AuthenticatedUserIDHash: validSHA256Hex(),
		AuthenticationMethod:    "saml_sso",
		DelegationChain: []DelegationChainEntry{
			{AuthenticatedUserIDHash: validSHA256Hex(), AuthenticationMethod: "magic"},
		},
	}
	err := ValidateActorAttributes(a)
	if err == nil {
		t.Fatal("expected error for bad delegation method")
	}
	if !strings.Contains(err.Error(), "delegation_chain[0]") && !strings.Contains(err.Error(), "not in the closed enum") {
		t.Errorf("unexpected error: %v", err)
	}
}

// -- §14.7 ValidateReasoningAttributes ------------------------------------

func TestValidateReasoningAttributes_Nil(t *testing.T) {
	if err := ValidateReasoningAttributes(nil); err != nil {
		t.Errorf("nil reasoning should pass: %v", err)
	}
}

func TestValidateReasoningAttributes_AllEnumValues(t *testing.T) {
	kinds := []string{
		"neurosymbolic", "retrieval_grounded_with_citations", "rule_based",
		"post_hoc_llm_rationalization", "attention_feature_importance",
		"none", "institution_named",
	}
	for _, k := range kinds {
		t.Run(k, func(t *testing.T) {
			r := &ReasoningAttributes{SubstrateKind: k}
			if err := ValidateReasoningAttributes(r); err != nil {
				t.Errorf("valid enum value %q rejected: %v", k, err)
			}
		})
	}
}

func TestValidateReasoningAttributes_Empty(t *testing.T) {
	r := &ReasoningAttributes{SubstrateKind: ""}
	err := ValidateReasoningAttributes(r)
	if err == nil {
		t.Fatal("expected error for empty substrate_kind")
	}
	if !strings.Contains(err.Error(), "substrate_kind is empty") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateReasoningAttributes_InvalidValue(t *testing.T) {
	r := &ReasoningAttributes{SubstrateKind: "deep_learning"}
	err := ValidateReasoningAttributes(r)
	if err == nil {
		t.Fatal("expected error for invalid substrate_kind")
	}
	if !strings.Contains(err.Error(), "not in the closed enum") {
		t.Errorf("unexpected error: %v", err)
	}
}

// -- §14.8 ValidateDownstreamActionAttributes -----------------------------

func TestValidateDownstreamActionAttributes_Nil(t *testing.T) {
	if err := ValidateDownstreamActionAttributes(nil); err != nil {
		t.Errorf("nil downstream action should pass: %v", err)
	}
}

func TestValidateDownstreamActionAttributes_Valid(t *testing.T) {
	d := &DownstreamActionAttributes{
		ActionKind:         "account_status_change",
		SystemOfRecordID:   "core-banking-v3",
		ChangeRecordIDHash: validSHA256Hex(),
		AppliedAtUTC:       "2026-05-23T14:30:00Z",
	}
	if err := ValidateDownstreamActionAttributes(d); err != nil {
		t.Errorf("valid downstream action rejected: %v", err)
	}
}

func TestValidateDownstreamActionAttributes_MissingActionKind(t *testing.T) {
	d := &DownstreamActionAttributes{
		SystemOfRecordID:   "sys",
		ChangeRecordIDHash: validSHA256Hex(),
		AppliedAtUTC:       "2026-05-23T14:30:00Z",
	}
	err := ValidateDownstreamActionAttributes(d)
	if err == nil {
		t.Fatal("expected error for missing action_kind")
	}
	if !strings.Contains(err.Error(), "action_kind is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDownstreamActionAttributes_MissingSystemOfRecordID(t *testing.T) {
	d := &DownstreamActionAttributes{
		ActionKind:         "record_updated",
		ChangeRecordIDHash: validSHA256Hex(),
		AppliedAtUTC:       "2026-05-23T14:30:00Z",
	}
	err := ValidateDownstreamActionAttributes(d)
	if err == nil {
		t.Fatal("expected error for missing system_of_record_id")
	}
	if !strings.Contains(err.Error(), "system_of_record_id is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDownstreamActionAttributes_MissingChangeRecordIDHash(t *testing.T) {
	d := &DownstreamActionAttributes{
		ActionKind:       "record_updated",
		SystemOfRecordID: "sys",
		AppliedAtUTC:     "2026-05-23T14:30:00Z",
	}
	err := ValidateDownstreamActionAttributes(d)
	if err == nil {
		t.Fatal("expected error for missing change_record_id_hash")
	}
	if !strings.Contains(err.Error(), "change_record_id_hash is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDownstreamActionAttributes_MalformedChangeRecordIDHash(t *testing.T) {
	d := &DownstreamActionAttributes{
		ActionKind:         "record_updated",
		SystemOfRecordID:   "sys",
		ChangeRecordIDHash: "not-a-sha256",
		AppliedAtUTC:       "2026-05-23T14:30:00Z",
	}
	err := ValidateDownstreamActionAttributes(d)
	if err == nil {
		t.Fatal("expected error for malformed hash")
	}
	if !strings.Contains(err.Error(), "not well-formed SHA-256 hex") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDownstreamActionAttributes_MissingAppliedAtUTC(t *testing.T) {
	d := &DownstreamActionAttributes{
		ActionKind:         "record_updated",
		SystemOfRecordID:   "sys",
		ChangeRecordIDHash: validSHA256Hex(),
	}
	err := ValidateDownstreamActionAttributes(d)
	if err == nil {
		t.Fatal("expected error for missing applied_at_utc")
	}
	if !strings.Contains(err.Error(), "applied_at_utc is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDownstreamActionAttributes_InvalidAppliedAtUTC(t *testing.T) {
	d := &DownstreamActionAttributes{
		ActionKind:         "record_updated",
		SystemOfRecordID:   "sys",
		ChangeRecordIDHash: validSHA256Hex(),
		AppliedAtUTC:       "2026-05-23T14:30:00+05:00",
	}
	err := ValidateDownstreamActionAttributes(d)
	if err == nil {
		t.Fatal("expected error for non-UTC timestamp")
	}
	if !strings.Contains(err.Error(), "not valid RFC 3339 UTC") {
		t.Errorf("unexpected error: %v", err)
	}
}

// -- ValidateEventAttributes (top-level dispatcher) -----------------------

func TestValidateEventAttributes_NilReturnsNoErrors(t *testing.T) {
	errs := ValidateEventAttributes(nil)
	if len(errs) != 0 {
		t.Errorf("nil attributes produced %d errors, want 0", len(errs))
	}
}

func TestValidateEventAttributes_AllFamiliesValid(t *testing.T) {
	attrs := &EventAttributes{
		Actor: &ActorAttributes{
			AuthenticatedUserIDHash: validSHA256Hex(),
			AuthenticationMethod:    "oidc_sso",
		},
		Reasoning: &ReasoningAttributes{
			SubstrateKind: "neurosymbolic",
		},
		DownstreamAction: &DownstreamActionAttributes{
			ActionKind:         "payment_executed",
			SystemOfRecordID:   "ledger-v2",
			ChangeRecordIDHash: validSHA256Hex(),
			AppliedAtUTC:       "2026-01-15T08:00:00Z",
		},
	}
	errs := ValidateEventAttributes(attrs)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateEventAttributes_MultipleFamiliesCanFail(t *testing.T) {
	attrs := &EventAttributes{
		Actor: &ActorAttributes{
			AuthenticatedUserIDHash: "bad",
			AuthenticationMethod:    "saml_sso",
		},
		Reasoning: &ReasoningAttributes{
			SubstrateKind: "invalid_kind",
		},
	}
	errs := ValidateEventAttributes(attrs)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

// -- ParseEventAttributes -------------------------------------------------

func TestParseEventAttributes_NoFamiliesReturnsNil(t *testing.T) {
	payload := []byte(`{"event":"decision","seq":0}`)
	attrs, err := ParseEventAttributes(payload)
	if err != nil {
		t.Fatal(err)
	}
	if attrs != nil {
		t.Error("expected nil when no families present")
	}
}

func TestParseEventAttributes_ActorFamilyParsed(t *testing.T) {
	payload := []byte(`{"event":"decision","audit.actor":{"authenticated_user_id_hash":"` + validSHA256Hex() + `","authentication_method":"saml_sso"}}`)
	attrs, err := ParseEventAttributes(payload)
	if err != nil {
		t.Fatal(err)
	}
	if attrs == nil || attrs.Actor == nil {
		t.Fatal("expected actor attributes to be parsed")
	}
	if attrs.Actor.AuthenticationMethod != "saml_sso" {
		t.Errorf("authentication_method = %q, want %q", attrs.Actor.AuthenticationMethod, "saml_sso")
	}
}

func TestParseEventAttributes_ReasoningFamilyParsed(t *testing.T) {
	payload := []byte(`{"audit.reasoning":{"substrate_kind":"rule_based"}}`)
	attrs, err := ParseEventAttributes(payload)
	if err != nil {
		t.Fatal(err)
	}
	if attrs == nil || attrs.Reasoning == nil {
		t.Fatal("expected reasoning attributes to be parsed")
	}
	if attrs.Reasoning.SubstrateKind != "rule_based" {
		t.Errorf("substrate_kind = %q, want %q", attrs.Reasoning.SubstrateKind, "rule_based")
	}
}

func TestParseEventAttributes_DownstreamActionFamilyParsed(t *testing.T) {
	h := validSHA256Hex()
	payload := []byte(`{"audit.downstream_action":{"action_kind":"record_updated","system_of_record_id":"sys","change_record_id_hash":"` + h + `","applied_at_utc":"2026-05-23T10:00:00Z"}}`)
	attrs, err := ParseEventAttributes(payload)
	if err != nil {
		t.Fatal(err)
	}
	if attrs == nil || attrs.DownstreamAction == nil {
		t.Fatal("expected downstream action attributes to be parsed")
	}
	if attrs.DownstreamAction.ActionKind != "record_updated" {
		t.Errorf("action_kind = %q, want %q", attrs.DownstreamAction.ActionKind, "record_updated")
	}
}

func TestParseEventAttributes_InvalidJSON(t *testing.T) {
	_, err := ParseEventAttributes([]byte("{not json"))
	if err == nil {
		t.Fatal("expected error on invalid JSON")
	}
}

// -- isDelegationChainLexSorted -------------------------------------------

func TestIsDelegationChainLexSorted_Empty(t *testing.T) {
	if !isDelegationChainLexSorted(nil) {
		t.Error("nil chain should be sorted")
	}
}

func TestIsDelegationChainLexSorted_Single(t *testing.T) {
	chain := []DelegationChainEntry{
		{AuthenticatedUserIDHash: validSHA256Hex(), AuthenticationMethod: "saml_sso"},
	}
	if !isDelegationChainLexSorted(chain) {
		t.Error("single-entry chain should be sorted")
	}
}

func TestIsDelegationChainLexSorted_Sorted(t *testing.T) {
	// "00..." serialises before "ff..." in JSON canonical form.
	chain := []DelegationChainEntry{
		{AuthenticatedUserIDHash: strings.Repeat("00", 32), AuthenticationMethod: "saml_sso"},
		{AuthenticatedUserIDHash: strings.Repeat("ff", 32), AuthenticationMethod: "saml_sso"},
	}
	if !isDelegationChainLexSorted(chain) {
		t.Error("sorted chain rejected")
	}
}

func TestIsDelegationChainLexSorted_Unsorted(t *testing.T) {
	chain := []DelegationChainEntry{
		{AuthenticatedUserIDHash: strings.Repeat("ff", 32), AuthenticationMethod: "saml_sso"},
		{AuthenticatedUserIDHash: strings.Repeat("00", 32), AuthenticationMethod: "saml_sso"},
	}
	if isDelegationChainLexSorted(chain) {
		t.Error("unsorted chain accepted")
	}
}

// TestIsDelegationChainLexSorted_SortKeyIsJCSBytes is the regression
// guard for the json.Marshal → core/jcs migration. The sort key MUST be
// the RFC 8785 canonical bytes of each entry (byte-ordinal key order),
// matching the .NET reference (AuditActor sorts by CanonicalizeEnvelope
// under StringComparer.Ordinal). This test pins that the canonical form
// of each entry — keys sorted authenticated_user_id_hash <
// authentication_method (byte-ordinal) — is what drives the order, so a
// future change back to a non-JCS serializer is caught.
func TestIsDelegationChainLexSorted_SortKeyIsJCSBytes(t *testing.T) {
	// Two entries identical except the hash; the canonical bytes differ
	// only at the hash value, so the entry with the lexically smaller
	// hash sorts first. A serializer that emitted keys in a different
	// order (or HTML-escaped a value) would produce a different key and
	// could mis-order entries whose discriminating field is not the
	// first declared struct field.
	lo := DelegationChainEntry{
		AuthenticatedUserIDHash: strings.Repeat("0a", 32),
		AuthenticationMethod:    "oidc_sso",
	}
	hi := DelegationChainEntry{
		AuthenticatedUserIDHash: strings.Repeat("0b", 32),
		AuthenticationMethod:    "oidc_sso",
	}

	if !isDelegationChainLexSorted([]DelegationChainEntry{lo, hi}) {
		t.Error("entries in ascending JCS-canonical order rejected")
	}
	if isDelegationChainLexSorted([]DelegationChainEntry{hi, lo}) {
		t.Error("entries in descending JCS-canonical order accepted")
	}

	// Verify the canonical bytes have keys in byte-ordinal order —
	// authenticated_user_id_hash (0x61...) sorts before
	// authentication_method only by comparing the full key bytes; both
	// share the "authentic" prefix, so the discriminating byte is
	// 'ed_' vs 'ation'. This asserts the JCS primitive is doing the
	// sort, not Go struct field order.
	canon, err := jcs.Canonicalize(delegationEntryMap(lo))
	if err != nil {
		t.Fatalf("canonicalize: %v", err)
	}
	got := string(canon)
	want := `{"authenticated_user_id_hash":"` + strings.Repeat("0a", 32) +
		`","authentication_method":"oidc_sso"}`
	if got != want {
		t.Errorf("delegation entry canonical bytes:\n got: %s\nwant: %s", got, want)
	}
}
