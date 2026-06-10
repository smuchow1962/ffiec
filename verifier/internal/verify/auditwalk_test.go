package verify

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/mmpworks/ffiec/core/constants"
)

// validHeader returns a §7-conformant header for the test tenant so the
// pre-flight steps 1-3 pass and the walk reaches the per-event steps.
func validHeader(tenantID string) AuditHeader {
	return AuditHeader{
		FormatVersion:       constants.FormatVersion,
		TenantID:            tenantID,
		SealDate:            "2026-05-06",
		GenesisHashHex:      genesisZeroHex,
		HKDFInputsDigestHex: digestForTenant(tenantID),
	}
}

// TestWalkAuditFile_Truncation asserts the pre-flight mid-write
// truncation refusal: a marked audit_file fails before step 1 with the
// §10.12 structural input error (exit 2).
func TestWalkAuditFile_Truncation(t *testing.T) {
	marker := "true"
	af := &AuditFile{NDJSONTruncated: &marker, Header: validHeader("tenant-x")}
	got := WalkAuditFile(af, IKMRegistry{})
	if got.Status != "FAIL" || got.Step != "pre-flight" || got.ExitCode != 2 {
		t.Fatalf("truncation: got %+v, want FAIL/pre-flight/exit2", got)
	}
	if got.Reason != "audit file ends mid-line — possible mid-write crash" {
		t.Errorf("truncation reason: %q", got.Reason)
	}
}

// TestWalkAuditFile_Step1FormatVersion asserts a non-v1 header fails at
// step 1 with the value rendered into the reason.
func TestWalkAuditFile_Step1FormatVersion(t *testing.T) {
	h := validHeader("tenant-x")
	h.FormatVersion = "v2"
	got := WalkAuditFile(&AuditFile{Header: h}, IKMRegistry{})
	if got.Status != "FAIL" || got.Step != "1" || got.ExitCode != 1 {
		t.Fatalf("step1: got %+v", got)
	}
	// Quoted rendering (%q) — matches the §7 reason-string family and the
	// uniform N009/N022/N023 fixtures. See TestCheckFormatVersion_QuotedRendering.
	want := `format_version "v2" not supported by this verifier (running v1)`
	if got.Reason != want {
		t.Errorf("step1 reason: got %q want %q", got.Reason, want)
	}
}

// TestCheckFormatVersion_QuotedRendering pins the chosen format_version
// reason rendering: the offending value is QUOTED. This is the regression
// guard the N023 tolerance-tightening leaves behind — if a future change
// reverts to the unquoted %s form, this fails loudly, and so does the
// gate's byte-exact compare against the quoted N009/N022/N023 fixtures.
func TestCheckFormatVersion_QuotedRendering(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"v2", `format_version "v2" not supported by this verifier (running v1)`},
		{"v1.1", `format_version "v1.1" not supported by this verifier (running v1)`},
		{"V1", `format_version "V1" not supported by this verifier (running v1)`},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := checkFormatVersion(AuditHeader{FormatVersion: tc.in})
			if got.Status != "FAIL" || got.Step != "1" {
				t.Fatalf("got %+v, want FAIL at step 1", got)
			}
			if got.Reason != tc.want {
				t.Errorf("reason: got %q want %q", got.Reason, tc.want)
			}
		})
	}
}

// TestWalkAuditFile_Step3Genesis asserts a non-zero genesis fails at step 3.
func TestWalkAuditFile_Step3Genesis(t *testing.T) {
	h := validHeader("tenant-x")
	h.GenesisHashHex = "01" + genesisZeroHex[2:]
	got := WalkAuditFile(&AuditFile{Header: h}, IKMRegistry{})
	if got.Status != "FAIL" || got.Step != "3" {
		t.Fatalf("step3: got %+v", got)
	}
}

// TestWalkAuditFile_Step7UnknownKeyVersion asserts an entry whose
// key_version is absent from the registry fails at step 7 with no MAC
// compute.
func TestWalkAuditFile_Step7UnknownKeyVersion(t *testing.T) {
	tenant := "tenant-x"
	h := validHeader(tenant)
	e := AuditEntry{
		Seq:           1,
		KeyVersion:    99,
		FormatVersion: constants.FormatVersion,
		Event:         map[string]interface{}{"tenant_id": tenant},
		PrevHashHex:   genesisZeroHex,
	}
	got := WalkAuditFile(&AuditFile{Header: h, Entries: []AuditEntry{e}}, IKMRegistry{1: []byte("k")})
	if got.Status != "FAIL" || got.Step != "7" {
		t.Fatalf("step7: got %+v", got)
	}
}

// TestWalkAuditFile_Step4CrossChain asserts an entry whose event tenant_id
// disagrees with the header fails at step 4 before any key work.
func TestWalkAuditFile_Step4CrossChain(t *testing.T) {
	h := validHeader("tenant-x")
	e := AuditEntry{
		Seq:           1,
		FormatVersion: constants.FormatVersion,
		Event:         map[string]interface{}{"tenant_id": "tenant-OTHER"},
		PrevHashHex:   genesisZeroHex,
	}
	got := WalkAuditFile(&AuditFile{Header: h, Entries: []AuditEntry{e}}, IKMRegistry{1: []byte("k")})
	if got.Status != "FAIL" || got.Step != "4" {
		t.Fatalf("step4: got %+v", got)
	}
}

// TestCheckAlgorithmKeyType asserts the structural §7-step-11 compare:
// a mismatch surfaces the named reason; agreement reports no mismatch.
func TestCheckAlgorithmKeyType(t *testing.T) {
	mismatch := AuditSeal{Algorithm: "ed25519", ResolvedPublicKeyType: "rsa-3072"}
	got, ok := CheckAlgorithmKeyType(mismatch)
	if !ok || got.Step != "11" || got.Reason != "algorithm/key-type mismatch at signature verification" {
		t.Fatalf("mismatch: got %+v ok=%v", got, ok)
	}
	agree := AuditSeal{Algorithm: "ed25519", ResolvedPublicKeyType: "ed25519"}
	if _, ok := CheckAlgorithmKeyType(agree); ok {
		t.Errorf("agreement should report no mismatch")
	}
	absent := AuditSeal{Algorithm: "ed25519"}
	if _, ok := CheckAlgorithmKeyType(absent); ok {
		t.Errorf("absent resolved-key-type should report no mismatch")
	}
}

// digestForTenant recomputes the §7-step-2 hkdf_inputs_digest for a tenant
// using the same construction the walk uses, so a test header is
// self-consistent. Kept test-local to avoid widening the production API.
func digestForTenant(tenantID string) string {
	info := []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + tenantID)
	le32 := make([]byte, 4)
	binary.LittleEndian.PutUint32(le32, uint32(constants.HKDFOutputLength))
	sum := sha256.New()
	sum.Write([]byte(constants.HKDFSalt))
	sum.Write(info)
	sum.Write(le32)
	return hex.EncodeToString(sum.Sum(nil))
}
