package verify

import (
	"encoding/hex"
	"testing"

	"github.com/mmpworks/ffiec/core/jcs"
)

// TestOperationalRoot_MatchesVector027 pins the leaf form to the corpus's own value, with
// no corpus on disk needed: SHA-256(0x00 || canonical bytes) over vector 027's 3 events.
func TestOperationalRoot_MatchesVector027(t *testing.T) {
	const bound = "af8557b1ea725b1fa02e7329c1149e33194276df299a9ff9c448280c964f6a20"
	events := []map[string]any{
		{"chain_kind": "operational", "event_type": "master.reconciliation_completed", "seq": 1, "tenant_id": "test-tenant-v1.0c", "ts_utc": "2026-05-11T01:00:00Z"},
		{"chain_kind": "operational", "event_type": "verifier.run_completed", "seq": 2, "tenant_id": "test-tenant-v1.0c", "ts_utc": "2026-05-11T06:30:00Z"},
		{"chain_kind": "operational", "event_type": "chain.verification_failure", "seq": 3, "tenant_id": "test-tenant-v1.0c", "ts_utc": "2026-05-11T18:45:00Z"},
	}
	af := &AuditFile{Seal: AuditSeal{SignPayloadVersion: signPayloadV1_0c, OperationalEventsLogRoot: bound}}
	for _, event := range events {
		canonical, err := jcs.Canonicalize(event)
		if err != nil {
			t.Fatal(err)
		}
		af.Entries = append(af.Entries, AuditEntry{Event: event, EventCanonicalHex: hex.EncodeToString(canonical)})
	}
	if out := checkSealOperationalRoot(af); out.Status == "FAIL" {
		t.Fatalf("got %q", out.Reason)
	}
	af.Seal.OperationalEventsLogRoot = emptyTreeRootForTest
	if out := checkSealOperationalRoot(af); out.Status != "FAIL" {
		t.Fatal("a false bound root passed")
	}
}

// §4.4 makes kms_handle_uri optional, so an entry with none adds nothing to the §4.3 digest form.
func TestKMSDigest_AnAbsentURIIsNotAValue(t *testing.T) {
	withAbsent := &AuditFile{Entries: []AuditEntry{{}, {KMSHandleURI: "aws-kms:a"}}}
	onlyPresent := &AuditFile{Entries: []AuditEntry{{KMSHandleURI: "aws-kms:a"}}}
	if kmsHandleURIsDigest(withAbsent) != kmsHandleURIsDigest(onlyPresent) {
		t.Fatal("an entry with no kms_handle_uri changed the digest")
	}
}

const emptyTreeRootForTest = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
