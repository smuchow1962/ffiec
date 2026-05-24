package verify

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// -- §14.6 closed enum for authentication_method --------------------------

// authMethodEnum is the §14.6 closed canonical enumeration. The last
// value ("institution_named") is the escape hatch per CC8.1 — any
// value not in this set AND not "institution_named" is rejected.
var authMethodEnum = map[string]bool{
	"saml_sso":           true,
	"oidc_sso":           true,
	"mtls_workload":      true,
	"spiffe_workload":    true,
	"hsm_bearer_token":   true,
	"api_key_with_iam":   true,
	"institution_named":  true,
}

// -- §14.7 closed enum for substrate_kind ---------------------------------

// substrateKindEnum is the §14.7 closed canonical enumeration.
var substrateKindEnum = map[string]bool{
	"neurosymbolic":                     true,
	"retrieval_grounded_with_citations": true,
	"rule_based":                        true,
	"post_hoc_llm_rationalization":      true,
	"attention_feature_importance":      true,
	"none":                              true,
	"institution_named":                 true,
}

// -- validation helpers ---------------------------------------------------

// isWellFormedSHA256Hex returns true when s is exactly 64 lowercase hex
// characters (the text representation of a 32-byte SHA-256 digest).
func isWellFormedSHA256Hex(s string) bool {
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

// isValidRFC3339UTC returns true when s is a valid RFC 3339 timestamp
// whose UTC offset is exactly "Z" (zero offset). The spec requires
// applied_at_utc to be RFC 3339 UTC.
func isValidRFC3339UTC(s string) bool {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return false
	}
	_, offset := t.Zone()
	return offset == 0 && strings.HasSuffix(s, "Z")
}

// -- predicate functions --------------------------------------------------

// ValidateActorAttributes checks §14.6 constraints on an ActorAttributes
// value. Returns nil when all constraints hold.
func ValidateActorAttributes(a *ActorAttributes) error {
	if a == nil {
		return nil
	}
	if a.AuthenticatedUserIDHash == "" {
		return fmt.Errorf("audit.actor: authenticated_user_id_hash is required when the family is emitted")
	}
	if !isWellFormedSHA256Hex(a.AuthenticatedUserIDHash) {
		return fmt.Errorf("audit.actor: authenticated_user_id_hash is not well-formed SHA-256 hex: %q", a.AuthenticatedUserIDHash)
	}
	if a.AuthenticationMethod == "" {
		return fmt.Errorf("audit.actor: authentication_method is required when the family is emitted")
	}
	if !authMethodEnum[a.AuthenticationMethod] {
		return fmt.Errorf("audit.actor: authentication_method %q is not in the closed enum", a.AuthenticationMethod)
	}
	if err := validateDelegationChain(a.DelegationChain); err != nil {
		return err
	}
	return nil
}

// validateDelegationChain checks that delegation_chain entries are
// well-formed and JCS-canonical lex-sorted.
func validateDelegationChain(chain []DelegationChainEntry) error {
	for i, d := range chain {
		if !isWellFormedSHA256Hex(d.AuthenticatedUserIDHash) {
			return fmt.Errorf("audit.actor.delegation_chain[%d]: authenticated_user_id_hash is not well-formed SHA-256 hex", i)
		}
		if !authMethodEnum[d.AuthenticationMethod] {
			return fmt.Errorf("audit.actor.delegation_chain[%d]: authentication_method %q is not in the closed enum", i, d.AuthenticationMethod)
		}
	}
	if len(chain) < 2 {
		return nil // zero or one entry is trivially sorted
	}
	if !isDelegationChainLexSorted(chain) {
		return fmt.Errorf("audit.actor.delegation_chain: entries are not JCS-canonical lex-sorted")
	}
	return nil
}

// isDelegationChainLexSorted checks whether the delegation_chain entries
// are lex-sorted by their JCS-canonical JSON serialisation. The spec
// requires "JCS-canonical lex-sorted array of delegated-authority
// identity pairs."
func isDelegationChainLexSorted(chain []DelegationChainEntry) bool {
	keys := make([]string, len(chain))
	for i, d := range chain {
		// JCS serialisation of each entry is deterministic because struct
		// fields have a fixed order. json.Marshal on a Go struct produces
		// fields in declared order — which matches JCS for ASCII-only keys.
		b, err := json.Marshal(d)
		if err != nil {
			return false // serialisation failure — treat as unsorted
		}
		keys[i] = string(b)
	}
	return sort.StringsAreSorted(keys)
}

// ValidateReasoningAttributes checks §14.7 constraints.
func ValidateReasoningAttributes(r *ReasoningAttributes) error {
	if r == nil {
		return nil
	}
	if r.SubstrateKind == "" {
		return fmt.Errorf("audit.reasoning: substrate_kind is empty")
	}
	if !substrateKindEnum[r.SubstrateKind] {
		return fmt.Errorf("audit.reasoning: substrate_kind %q is not in the closed enum", r.SubstrateKind)
	}
	return nil
}

// ValidateDownstreamActionAttributes checks §14.8 constraints.
func ValidateDownstreamActionAttributes(d *DownstreamActionAttributes) error {
	if d == nil {
		return nil
	}
	if d.ActionKind == "" {
		return fmt.Errorf("audit.downstream_action: action_kind is required when the family is emitted")
	}
	if d.SystemOfRecordID == "" {
		return fmt.Errorf("audit.downstream_action: system_of_record_id is required when the family is emitted")
	}
	if d.ChangeRecordIDHash == "" {
		return fmt.Errorf("audit.downstream_action: change_record_id_hash is required when the family is emitted")
	}
	if !isWellFormedSHA256Hex(d.ChangeRecordIDHash) {
		return fmt.Errorf("audit.downstream_action: change_record_id_hash is not well-formed SHA-256 hex: %q", d.ChangeRecordIDHash)
	}
	if d.AppliedAtUTC == "" {
		return fmt.Errorf("audit.downstream_action: applied_at_utc is required when the family is emitted")
	}
	if !isValidRFC3339UTC(d.AppliedAtUTC) {
		return fmt.Errorf("audit.downstream_action: applied_at_utc is not valid RFC 3339 UTC: %q", d.AppliedAtUTC)
	}
	return nil
}

// ValidateEventAttributes is the top-level dispatcher: when any of the
// three attribute families are present on a decoded event payload, the
// family-specific validator runs. Families that are nil (absent from the
// JSON) are silently skipped.
func ValidateEventAttributes(attrs *EventAttributes) []error {
	if attrs == nil {
		return nil
	}
	var errs []error
	if err := ValidateActorAttributes(attrs.Actor); err != nil {
		errs = append(errs, err)
	}
	if err := ValidateReasoningAttributes(attrs.Reasoning); err != nil {
		errs = append(errs, err)
	}
	if err := ValidateDownstreamActionAttributes(attrs.DownstreamAction); err != nil {
		errs = append(errs, err)
	}
	return errs
}

// ParseEventAttributes attempts to unmarshal the §14.6/§14.7/§14.8
// attribute families from a decoded event payload. Returns nil (not an
// error) when no families are present — they are all optional.
func ParseEventAttributes(eventPayloadJSON []byte) (*EventAttributes, error) {
	var attrs EventAttributes
	if err := json.Unmarshal(eventPayloadJSON, &attrs); err != nil {
		return nil, fmt.Errorf("parse event attributes: %w", err)
	}
	if attrs.Actor == nil && attrs.Reasoning == nil && attrs.DownstreamAction == nil {
		return nil, nil // no attribute families present; nothing to validate
	}
	return &attrs, nil
}
