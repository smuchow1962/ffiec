package verify

// EventAttributes holds the optional §14.6 / §14.7 / §14.8 attribute
// families that MAY appear on a chain entry's event payload. Each family
// is parsed only when its required fields are present; missing families
// are silently ignored (they are OPTIONAL per spec).
//
// The attribute names map 1:1 to the JCS-canonical JSON keys the
// Herald.Compliance (C#) and Herald.Py (Python) emitters produce.
// Diverging from these names would break the byte-equivalence contract.
type EventAttributes struct {
	// §14.6 — audit.actor.* family (Seam D)
	Actor *ActorAttributes `json:"audit.actor,omitempty"`

	// §14.7 — audit.reasoning.substrate_kind (Seam B)
	Reasoning *ReasoningAttributes `json:"audit.reasoning,omitempty"`

	// §14.8 — audit.downstream_action.* family
	DownstreamAction *DownstreamActionAttributes `json:"audit.downstream_action,omitempty"`

	// §14.13 — audit.supervisory.* supervisory-context provenance.
	// Presentation-only: the verifier surfaces it through the profile
	// layer and never gates the integrity verdict on it (per §14.13 and
	// the §7 determinism contract).
	Supervisory *SupervisoryAttributes `json:"audit.supervisory,omitempty"`
}

// SupervisoryAttributes is §14.13. Every field is institution-asserted
// provenance bound under the per-event MAC like any audit.* namespace;
// the verifier records it for examiner-facing rendering and applies no
// predicate to it. charter_type is present when the family is emitted;
// the supervisor fields are conditional on the charter class.
type SupervisoryAttributes struct {
	CharterType                 string `json:"charter_type"`
	PrimaryStateSupervisor      string `json:"primary_state_supervisor,omitempty"`
	FederalPrudentialSupervisor string `json:"federal_prudential_supervisor,omitempty"`
	DualSupervision             bool   `json:"dual_supervision,omitempty"`
}

// ActorAttributes is §14.6. When emitted, authenticated_user_id_hash
// and authentication_method are required; session_id is RECOMMENDED;
// delegation_chain is present when applicable.
type ActorAttributes struct {
	AuthenticatedUserIDHash string                 `json:"authenticated_user_id_hash"`
	AuthenticationMethod    string                 `json:"authentication_method"`
	SessionID               string                 `json:"session_id,omitempty"`
	DelegationChain         []DelegationChainEntry `json:"delegation_chain,omitempty"`
}

// DelegationChainEntry is one delegated-authority identity pair in the
// §14.6 delegation_chain array. Each entry is itself an
// authenticated_user_id_hash + authentication_method pair.
type DelegationChainEntry struct {
	AuthenticatedUserIDHash string `json:"authenticated_user_id_hash"`
	AuthenticationMethod    string `json:"authentication_method"`
}

// ReasoningAttributes is §14.7. When emitted, substrate_kind is
// RECOMMENDED.
type ReasoningAttributes struct {
	SubstrateKind string `json:"substrate_kind"`
}

// DownstreamActionAttributes is §14.8. When emitted, all four fields
// are required.
type DownstreamActionAttributes struct {
	ActionKind         string `json:"action_kind"`
	SystemOfRecordID   string `json:"system_of_record_id"`
	ChangeRecordIDHash string `json:"change_record_id_hash"`
	AppliedAtUTC       string `json:"applied_at_utc"`
}
