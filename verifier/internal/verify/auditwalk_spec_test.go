package verify

import (
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/core/constants"
)

// §7 step 3a. A rewritten header tenant_id with a character outside the §3 class
// fails with its own reason, before the walk reaches any entry.
func TestWalkAuditFile_Step3aTenantClass(t *testing.T) {
	h := validHeader("tenant|x")
	got := WalkAuditFile(&AuditFile{Header: h}, IKMRegistry{})
	if got.Status != "FAIL" || got.Step != "3a" || got.Reason != "tenant_id violates §3 character class" {
		t.Fatalf("step3a: got %+v", got)
	}
}

// §7 step 4 binds the run as well as the tenant: event.run_id == header.chain_id.
// An entry lifted from another run of the same tenant fails here, by name.
func TestWalkAuditFile_Step4RunBinding(t *testing.T) {
	h := validHeader("tenant-x")
	h.ChainID = "RUN-A"
	e := AuditEntry{
		Seq:           1,
		FormatVersion: constants.FormatVersion,
		Event:         map[string]interface{}{"tenant_id": "tenant-x", "run_id": "RUN-B"},
		PrevHashHex:   genesisZeroHex,
	}
	got := WalkAuditFile(&AuditFile{Header: h, Entries: []AuditEntry{e}}, IKMRegistry{1: []byte("k")})
	if got.Status != "FAIL" || got.Step != "4" || !strings.HasPrefix(got.Reason, "cross-chain lift detected at seq 1") {
		t.Fatalf("step4 run binding: got %+v", got)
	}
}

// §7 step 6 names two failures. A good prev_hash with the wrong seq is the second one.
func TestWalkAuditFile_Step6SeqOutOfOrder(t *testing.T) {
	h := validHeader("tenant-x")
	e := AuditEntry{
		Seq:           5,
		FormatVersion: constants.FormatVersion,
		Event:         map[string]interface{}{"tenant_id": "tenant-x"},
		PrevHashHex:   genesisZeroHex,
	}
	got := WalkAuditFile(&AuditFile{Header: h, Entries: []AuditEntry{e}}, IKMRegistry{1: []byte("k")})
	if got.Status != "FAIL" || got.Step != "6" || got.Reason != "seq out of order at seq 1" {
		t.Fatalf("step6 seq: got %+v", got)
	}
}
