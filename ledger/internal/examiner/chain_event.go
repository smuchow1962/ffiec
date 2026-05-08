package examiner

import (
	"context"
	"fmt"
	"time"
)

// ChainWriter writes an operational event to the institution's append-only
// ledger. Issue and Revoke depend on this interface; the concrete
// implementation lives outside this package — the ledger server's append
// path will satisfy it once that path lands.
//
// For ledgerctl bootstrap usage, ledger/cmd/ledgerctl ships a stdout-printing
// implementation so the CLI is exercisable end-to-end before the real append
// path is wired up.
type ChainWriter interface {
	Write(ctx context.Context, evt OperationalEvent) error
}

// OperationalEvent matches the shape in docs/operator-guide.md §"Receipt and
// disposal evidence". The two reserved event names below bracket every
// examiner credential issuance and revocation; the schema is shared with
// the rest of the operational-events vocabulary.
type OperationalEvent struct {
	Event         string                 `json:"event"`
	Timestamp     time.Time              `json:"timestamp"`
	TenantID      string                 `json:"tenant_id"`
	CorrelationID string                 `json:"correlation_id"`
	Fields        map[string]interface{} `json:"fields"`
}

// Reserved event names.
const (
	EventHandoverReceived = "master_key.examiner_handover_received"
	EventHandoverReturned = "master_key.examiner_handover_returned"
)

// MultiTenantWrite fans an event out across every tenant in the supplied
// list so the issuance/revocation lands inside each tenant's chain. The
// event's TenantID field is rewritten per call.
func MultiTenantWrite(ctx context.Context, w ChainWriter, evt OperationalEvent, tenants []string) error {
	if w == nil {
		return fmt.Errorf("multi-tenant write: nil chain writer")
	}
	if len(tenants) == 0 {
		return fmt.Errorf("multi-tenant write: no tenants supplied")
	}
	for _, t := range tenants {
		e := evt
		e.TenantID = t
		if err := w.Write(ctx, e); err != nil {
			return fmt.Errorf("write event for tenant %q: %w", t, err)
		}
	}
	return nil
}
