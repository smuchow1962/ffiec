package verify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// §10.84 — Communication principal-preapproval ordering.
//
// FINRA Rule 2210 requires a registered principal to approve a retail
// communication before it is sent (pre-approval). §10.84 proves the
// ORDERING — approval preceded the send — by composing three existing
// event shapes rather than adding a new cryptographic mechanism:
//
//   - the communication event carries audit.communication.audience;
//   - a §10.50 review event with audit.review.role = "registered_principal"
//     records the approval (its signed_at_utc is the approval time);
//   - a §14.8 downstream_action event with action_kind in
//     {communication_sent, notification_sent} records the send.
//
// The check is SOFT-ENFORCEMENT per §10.12: an ordering violation is an
// anomaly under Status: PASS, never a chain-integrity FAIL. It is scoped
// to retail communications; institutional and correspondence audiences
// are outside Rule 2210's pre-approval requirement and are not checked.

// preapprovalFamily is the AdditionalVerification family name for the
// §10.84 check. On OK it corresponds to the §10.12 closed-enum marker
// communication_principal_preapproval_verified; on failure the Note
// carries the §10.84 anomaly line.
const preapprovalFamily = "audit.communication_preapproval"

const preapprovalMarker = "communication_principal_preapproval_verified"

// communication-send action_kind values §10.84 recognizes on a §14.8
// downstream_action event. The set is small and closed for the §10.84
// dispatch; institution-named action kinds outside it are simply not
// treated as §10.84 sends.
var communicationSendActionKinds = map[string]bool{
	"communication_sent": true,
	"notification_sent":  true,
}

// preapprovalEventView is the projection of one chain entry's event
// payload that §10.84 needs. It tolerates both wire conventions the spec
// uses: the flat dotted §10.50 review keys and the nested §14.8
// downstream_action object. Absent fields decode to their zero values;
// the pointer fields distinguish "absent" from "zero".
type preapprovalEventView struct {
	RunID string `json:"run_id"`
	Seq   int    `json:"seq"`

	// §4.4 entry-level parent linkage (used by the send event).
	ParentRunID string `json:"parent_run_id"`
	ParentSeq   *int   `json:"parent_seq"`

	// §10.84 audience classifier on the communication event.
	Audience string `json:"audit.communication.audience"`

	// §10.50 review event (flat dotted keys; the approval).
	ReviewRole        string `json:"audit.review.role"`
	ReviewParentRunID string `json:"audit.review.parent_run_id"`
	ReviewParentSeq   *int   `json:"audit.review.parent_seq"`
	SignedReview      *struct {
		SignedAtUTC string `json:"audit.signed_review.signed_at_utc"`
	} `json:"audit.review.signed_review"`

	// §14.8 downstream_action event (nested object; the send).
	DownstreamAction *struct {
		ActionKind   string `json:"action_kind"`
		AppliedAtUTC string `json:"applied_at_utc"`
	} `json:"audit.downstream_action"`
}

// commRef identifies one communication by its (run_id, seq).
type commRef struct {
	runID string
	seq   int
}

// validateCommunicationPreapproval runs the §10.84 ordering check across
// all chain entries. It returns nil when the chain carries no retail
// communication event (nothing to check); otherwise a single
// AdditionalVerification whose OK reflects whether every retail
// communication's registered-principal approval preceded its send.
//
// Cognitive complexity is kept low by splitting the work into three
// passes: decode the views, index the retail communications, then check
// each one's approval-before-send ordering.
func validateCommunicationPreapproval(entries []ChainEntry) *AdditionalVerification {
	views := decodePreapprovalViews(entries)

	retail := retailCommunications(views)
	if len(retail) == 0 {
		return nil // no §10.84 subject in this chain
	}

	av := &AdditionalVerification{Family: preapprovalFamily, OK: true, EntryHits: len(retail)}
	for _, c := range retail {
		if note := checkOneCommunication(c, views); note != "" {
			av.OK = false
			if av.Note == "" {
				av.Note = note
			}
		}
	}
	return av
}

// decodePreapprovalViews decodes every entry's event payload into a
// preapprovalEventView. Entries whose payload cannot be decoded are
// skipped — other §7 steps own payload-integrity failures.
func decodePreapprovalViews(entries []ChainEntry) []preapprovalEventView {
	views := make([]preapprovalEventView, 0, len(entries))
	for i := range entries {
		payload, err := base64.StdEncoding.DecodeString(entries[i].EventPayloadJCS)
		if err != nil {
			continue
		}
		var v preapprovalEventView
		if err := json.Unmarshal(payload, &v); err != nil {
			continue
		}
		views = append(views, v)
	}
	return views
}

// retailCommunications returns the commRef of every communication event
// whose audience is "retail" — the only audience Rule 2210 subjects to
// registered-principal pre-approval.
func retailCommunications(views []preapprovalEventView) []commRef {
	var out []commRef
	for _, v := range views {
		if v.Audience == "retail" {
			out = append(out, commRef{runID: v.RunID, seq: v.Seq})
		}
	}
	return out
}

// checkOneCommunication returns "" when communication c's ordering is
// conformant (approval preceded send, and the marker is warranted), or a
// §10.84 anomaly line when it is not. The anomaly wording mirrors the
// spec verbatim so cross-implementation verifier output stays identical.
func checkOneCommunication(c commRef, views []preapprovalEventView) string {
	approvalAt, hasApproval := findPrincipalApproval(c, views)
	sendAt, sendSeq, hasSend := findCommunicationSend(c, views)

	if !hasSend {
		// No send bound yet — nothing to order against. Not an anomaly:
		// a communication approved but not yet sent is a normal pending
		// state, and a communication with neither is out of scope.
		return ""
	}
	if !hasApproval {
		return fmt.Sprintf("communication principal-preapproval ordering: no registered-principal approval bound at seq %d", sendSeq)
	}
	if approvalAfter(approvalAt, sendAt) {
		return fmt.Sprintf("communication principal-preapproval ordering: approval did not precede send at seq %d", sendSeq)
	}
	return ""
}

// findPrincipalApproval returns the approval timestamp of the §10.50
// review event with role "registered_principal" parent-linked to c.
func findPrincipalApproval(c commRef, views []preapprovalEventView) (string, bool) {
	for _, v := range views {
		if v.ReviewRole != "registered_principal" {
			continue
		}
		if v.ReviewParentRunID != c.runID || v.ReviewParentSeq == nil || *v.ReviewParentSeq != c.seq {
			continue
		}
		if v.SignedReview == nil || v.SignedReview.SignedAtUTC == "" {
			continue
		}
		return v.SignedReview.SignedAtUTC, true
	}
	return "", false
}

// findCommunicationSend returns the send timestamp and the send event's
// seq for the §14.8 downstream_action event parent-linked to c whose
// action_kind is a §10.84 communication-send kind.
func findCommunicationSend(c commRef, views []preapprovalEventView) (appliedAt string, sendSeq int, found bool) {
	for _, v := range views {
		if v.DownstreamAction == nil || !communicationSendActionKinds[v.DownstreamAction.ActionKind] {
			continue
		}
		if v.ParentRunID != c.runID || v.ParentSeq == nil || *v.ParentSeq != c.seq {
			continue
		}
		return v.DownstreamAction.AppliedAtUTC, v.Seq, true
	}
	return "", 0, false
}

// approvalAfter reports whether the approval timestamp is strictly after
// the send timestamp — i.e., the pre-approval ordering was violated.
// Unparseable timestamps are treated as a violation: the ordering cannot
// be proven, and §10.84 is a completeness check, so an unprovable
// ordering surfaces rather than passing silently.
func approvalAfter(approvalAtUTC, sendAtUTC string) bool {
	approval, err1 := time.Parse(time.RFC3339, approvalAtUTC)
	send, err2 := time.Parse(time.RFC3339, sendAtUTC)
	if err1 != nil || err2 != nil {
		return true
	}
	return approval.After(send)
}
