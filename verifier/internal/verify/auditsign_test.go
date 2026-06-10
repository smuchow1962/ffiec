package verify

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/mmpworks/ffiec/core/signpayload"
)

// signedV1aSeal builds an AuditSeal carrying a real v1.0a sign_payload and
// a valid Ed25519 signature over it, signed with priv. The returned seal
// is the positive-path fixture the live step-11 walk verifies.
func signedV1aSeal(t *testing.T, priv ed25519.PrivateKey, tenant string) AuditSeal {
	t.Helper()
	seal := AuditSeal{
		Algorithm:           "ed25519",
		FormatVersion:       "v1",
		SignPayloadVersion:  "v1.0a",
		Cadence:             "daily",
		DevMode:             false,
		TenantID:            tenant,
		SealDate:            "2026-05-06",
		MerkleRootHex:       "927adc88e5d843c5847674f7245d1e3c762bacc3f0b225f3322300cddf4eefe9",
		HKDFInputsDigestHex: "6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65",
	}
	payload := mustReconstruct(t, seal)
	seal.SignPayloadHex = hex.EncodeToString(payload)
	seal.SignatureB64 = base64.StdEncoding.EncodeToString(ed25519.Sign(priv, payload))
	return seal
}

func mustReconstruct(t *testing.T, seal AuditSeal) []byte {
	t.Helper()
	payload, err := reconstructSignPayload(seal)
	if err != nil {
		t.Fatalf("reconstruct sign_payload: %v", err)
	}
	return payload
}

func TestCheckSealSignatureV1_ValidSignaturePasses(t *testing.T) {
	priv := loadTestSigningKey(t)
	pub := priv.Public().(ed25519.PublicKey)
	seal := signedV1aSeal(t, priv, "tenant-ffiec-test-1")

	got := CheckSealSignatureV1(seal, pub)
	if got.Status != "PASS" || got.ExitCode != 0 {
		t.Fatalf("valid seal signature: got %+v, want PASS/0", got)
	}
}

func TestCheckSealSignatureV1_GarbageSignatureFails(t *testing.T) {
	priv := loadTestSigningKey(t)
	pub := priv.Public().(ed25519.PublicKey)
	seal := signedV1aSeal(t, priv, "tenant-ffiec-test-1")
	// Replace the signature with 64 bytes of zero — a decodable but
	// invalid signature, the N004 "signature garbage" failure mode.
	seal.SignatureB64 = base64.StdEncoding.EncodeToString(make([]byte, ed25519.SignatureSize))

	assertStep11Fail(t, CheckSealSignatureV1(seal, pub))
}

func TestCheckSealSignatureV1_WrongTenantFails(t *testing.T) {
	priv := loadTestSigningKey(t)
	pub := priv.Public().(ed25519.PublicKey)
	// Sign for tenant-OTHER (the N005 failure mode): the signature is a
	// REAL signature, but over a payload that binds a different tenant_id.
	signed := signedV1aSeal(t, priv, "tenant-ffiec-test-OTHER")
	// Present it as the expected tenant: the published sign_payload_hex and
	// signature came from the OTHER-tenant payload, but the seal's
	// structured tenant_id now claims test-1, so the reconstruct-and-
	// compare assertion catches the mismatch before the verify.
	presented := signed
	presented.TenantID = "tenant-ffiec-test-1"

	assertStep11Fail(t, CheckSealSignatureV1(presented, pub))
}

func TestCheckSealSignatureV1_PublishedHexTamperedFails(t *testing.T) {
	priv := loadTestSigningKey(t)
	pub := priv.Public().(ed25519.PublicKey)
	seal := signedV1aSeal(t, priv, "tenant-ffiec-test-1")
	// Flip one byte of the published sign_payload_hex so it no longer
	// equals what reconstructs from the structured fields. The verify call
	// is never reached — the reconstruct-and-compare assertion fails first.
	tampered := []byte(seal.SignPayloadHex)
	if tampered[0] == 'a' {
		tampered[0] = 'b'
	} else {
		tampered[0] = 'a'
	}
	seal.SignPayloadHex = string(tampered)

	assertStep11Fail(t, CheckSealSignatureV1(seal, pub))
}

func TestCheckSealSignatureV1_WrongKeyFails(t *testing.T) {
	priv := loadTestSigningKey(t)
	seal := signedV1aSeal(t, priv, "tenant-ffiec-test-1")
	// A different, valid Ed25519 public key — the verifier was pointed at
	// the wrong signer. Deterministic seed so the test is reproducible.
	otherPub := ed25519.NewKeyFromSeed(make([]byte, ed25519.SeedSize)).
		Public().(ed25519.PublicKey)

	assertStep11Fail(t, CheckSealSignatureV1(seal, otherPub))
}

func TestCheckSealSignatureV1_WrongKeyLengthFails(t *testing.T) {
	// No private material needed: a wrong-sized key is a config error the
	// step surfaces as a signature failure without touching the seed.
	seal := AuditSeal{SignPayloadVersion: "v1.0a", Algorithm: "ed25519"}
	assertStep11Fail(t, CheckSealSignatureV1(seal, ed25519.PublicKey{0x01, 0x02}))
}

func TestSealCarriesLiveSignature(t *testing.T) {
	tests := []struct {
		name string
		b64  string
		want bool
	}{
		{"valid 64-byte sig", base64.StdEncoding.EncodeToString(make([]byte, 64)), true},
		{"placeholder string", "TEST-SIGNATURE-PLACEHOLDER-not-an-ed25519-sig", false},
		{"empty", "", false},
		{"wrong length 32", base64.StdEncoding.EncodeToString(make([]byte, 32)), false},
		{"not base64", "!!!!not base64!!!!", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SealCarriesLiveSignature(AuditSeal{SignatureB64: tc.b64})
			if got != tc.want {
				t.Errorf("SealCarriesLiveSignature(%q) = %v, want %v", tc.b64, got, tc.want)
			}
		})
	}
}

func TestParsePublicKeyHex(t *testing.T) {
	tests := []struct {
		name    string
		hexStr  string
		wantErr bool
	}{
		{"published key", publishedPubHex, false},
		{"not hex", "zzzz", true},
		{"too short", "0985603b", true},
		{"empty", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pub, err := ParsePublicKeyHex(tc.hexStr)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParsePublicKeyHex(%q): want error, got nil", tc.hexStr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePublicKeyHex(%q): unexpected error %v", tc.hexStr, err)
			}
			if len(pub) != ed25519.PublicKeySize {
				t.Fatalf("parsed key has %d bytes, want %d", len(pub), ed25519.PublicKeySize)
			}
		})
	}
}

// TestReconstructSignPayload_MatchesBuilder confirms the step-11
// reconstruction is byte-identical to core/signpayload.Build for the same
// fields — the reconstruction is the spec's gated byte form, not a
// re-derivation.
func TestReconstructSignPayload_MatchesBuilder(t *testing.T) {
	seal := AuditSeal{
		Algorithm:           "ed25519",
		FormatVersion:       "v1",
		SignPayloadVersion:  "v1.0a",
		Cadence:             "daily",
		TenantID:            "tenant-ffiec-test-1",
		SealDate:            "2026-05-06",
		MerkleRootHex:       "927adc88e5d843c5847674f7245d1e3c762bacc3f0b225f3322300cddf4eefe9",
		HKDFInputsDigestHex: "6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65",
	}
	got := mustReconstruct(t, seal)
	want, err := signpayload.Build(signpayload.Seal{
		Version:          signpayload.VersionV1_0a,
		Algorithm:        "ed25519",
		FormatVersion:    "v1",
		TenantID:         "tenant-ffiec-test-1",
		SealDate:         "2026-05-06",
		MerkleRootHex:    "927adc88e5d843c5847674f7245d1e3c762bacc3f0b225f3322300cddf4eefe9",
		HKDFInputsDigest: "6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65",
		Cadence:          "daily",
	})
	if err != nil {
		t.Fatalf("builder: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("reconstruction differs from builder:\n got: %q\nwant: %q", got, want)
	}
}

// assertStep11Fail asserts the outcome is the §7 step-11 normative failure
// (FAIL / step 11 / exit 1 / the pinned reason).
func assertStep11Fail(t *testing.T, got Outcome) {
	t.Helper()
	if got.Status != "FAIL" {
		t.Fatalf("status: got %q, want FAIL", got.Status)
	}
	if got.Step != "11" {
		t.Fatalf("step: got %q, want 11", got.Step)
	}
	if got.ExitCode != 1 {
		t.Fatalf("exit code: got %d, want 1", got.ExitCode)
	}
	if got.Reason != "signature verification failed" {
		t.Fatalf("reason: got %q, want %q", got.Reason, "signature verification failed")
	}
}
