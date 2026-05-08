package examiner

import (
	"context"
	"fmt"
	"time"
)

// RevokeRequest is what the operator's CLI flags collect for `examiner revoke`.
type RevokeRequest struct {
	BundleID      string
	Reason        string
	CorrelationID string
	Tenants       []string

	Now func() time.Time
}

// Revoke writes a `master_key.examiner_handover_returned` chain event for
// each tenant in scope. Role tear-down on the Compliance read surface is the
// responsibility of the role provisioner that lives outside this package;
// the chain event is the auditable record that revocation occurred.
func Revoke(ctx context.Context, req RevokeRequest, w ChainWriter) error {
	if req.BundleID == "" {
		return fmt.Errorf("bundle_id is required")
	}
	if req.CorrelationID == "" {
		return fmt.Errorf("correlation_id is required")
	}
	if len(req.Tenants) == 0 {
		return fmt.Errorf("at least one tenant is required")
	}
	if req.Now == nil {
		req.Now = time.Now
	}
	now := req.Now().UTC()

	evt := OperationalEvent{
		Event:         EventHandoverReturned,
		Timestamp:     now,
		CorrelationID: req.CorrelationID,
		Fields: map[string]interface{}{
			"handover_authorization_record_id": req.BundleID,
			"reason":                           req.Reason,
		},
	}
	return MultiTenantWrite(ctx, w, evt, req.Tenants)
}
