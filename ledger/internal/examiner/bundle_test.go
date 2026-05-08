package examiner

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestSignVerifyRoundTrip(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	b := minimalBundle(pub)

	signed, err := Sign(b, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	parsed, err := Verify(signed, pub)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if parsed.BundleID != b.BundleID {
		t.Errorf("BundleID lost in round trip: got %q want %q", parsed.BundleID, b.BundleID)
	}
	if parsed.Scope.Surface != ScopeComplianceRead {
		t.Errorf("Scope.Surface lost in round trip: got %q", parsed.Scope.Surface)
	}
}

func TestVerifyRejectsTamper(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	b := minimalBundle(pub)
	signed, err := Sign(b, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(signed, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	raw["bundle_id"] = "bnd_TAMPER"
	tampered, _ := json.Marshal(raw)

	if _, err := Verify(tampered, pub); err == nil {
		t.Fatal("Verify accepted a tampered bundle")
	}
}

func TestVerifyRejectsWrongKey(t *testing.T) {
	pub1, priv1, _ := ed25519.GenerateKey(rand.Reader)
	pub2, _, _ := ed25519.GenerateKey(rand.Reader)
	b := minimalBundle(pub1)
	signed, err := Sign(b, priv1)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if _, err := Verify(signed, pub2); err == nil {
		t.Fatal("Verify accepted a signature under the wrong public key")
	}
}

func TestValidateRejectsBadInputs(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	cases := []struct {
		name string
		mut  func(*Bundle)
	}{
		{"missing bundle_id", func(b *Bundle) { b.BundleID = "" }},
		{"unknown version", func(b *Bundle) { b.Version = "0.9" }},
		{"expires_at <= issued_at", func(b *Bundle) { b.ExpiresAt = b.IssuedAt }},
		{"empty tenants", func(b *Bundle) { b.Scope.Tenants = nil }},
		{"unknown scope", func(b *Bundle) { b.Scope.Surface = Scope("compliance.write") }},
		{"missing correlation id", func(b *Bundle) { b.Scope.CorrelationID = "" }},
		{"missing compliance url", func(b *Bundle) { b.Endpoint.ComplianceURL = "" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := minimalBundle(pub)
			tc.mut(b)
			if err := b.Validate(); err == nil {
				t.Errorf("Validate accepted bad input: %s", tc.name)
			}
		})
	}
}

type recordingWriter struct{ events []OperationalEvent }

func (w *recordingWriter) Write(_ context.Context, e OperationalEvent) error {
	w.events = append(w.events, e)
	return nil
}

func TestIssueEmitsPerTenantReceiptEvents(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	rec := &recordingWriter{}
	frozen := time.Date(2026, 6, 15, 8, 30, 0, 0, time.UTC)

	res, err := Issue(context.Background(), IssueRequest{
		InstitutionID:      "northbridge-federal",
		IssuancePrivateKey: priv,
		IssuancePublicKey:  pub,
		ExaminerIdentity:   "OCC-IT-EX-Cert-2026-Karen",
		Tenants:            []string{"tenant_a", "tenant_b"},
		ExaminationStart:   "2026-01-01",
		ExaminationEnd:     "2026-06-14",
		ExpiresAt:          time.Date(2026, 6, 22, 17, 0, 0, 0, time.UTC),
		CorrelationID:      "exam-2026-06-15-occ",
		Scope:              ScopeComplianceRead,
		ComplianceURL:      "https://compliance.northbridge.example",
		Credential:         BundleCredential{Type: "mtls_client_cert"},
		Now:                func() time.Time { return frozen },
	}, rec)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if len(rec.events) != 2 {
		t.Fatalf("want 2 chain events (one per tenant), got %d", len(rec.events))
	}
	for i, e := range rec.events {
		if e.Event != EventHandoverReceived {
			t.Errorf("event[%d].Event = %q, want %q", i, e.Event, EventHandoverReceived)
		}
		if !e.Timestamp.Equal(frozen) {
			t.Errorf("event[%d].Timestamp = %v, want %v", i, e.Timestamp, frozen)
		}
		if e.CorrelationID != "exam-2026-06-15-occ" {
			t.Errorf("event[%d].CorrelationID = %q", i, e.CorrelationID)
		}
	}
	if rec.events[0].TenantID != "tenant_a" || rec.events[1].TenantID != "tenant_b" {
		t.Errorf("tenant fan-out wrong: %q, %q", rec.events[0].TenantID, rec.events[1].TenantID)
	}
	if !strings.HasPrefix(res.Bundle.BundleID, "bnd_") {
		t.Errorf("bundle id %q does not look minted", res.Bundle.BundleID)
	}
	if res.Bundle.Issuer.IssuanceKeyFingerprint == "" {
		t.Error("issuance key fingerprint missing")
	}
}

func TestVerifyRejectsBadSignatureBytes(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	b := minimalBundle(pub)
	signed, err := Sign(b, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(signed, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	raw["signature_ed25519"] = base64.StdEncoding.EncodeToString(make([]byte, 64))
	bogus, _ := json.Marshal(raw)
	if _, err := Verify(bogus, pub); err == nil {
		t.Fatal("Verify accepted a bundle with bogus signature bytes")
	}
}

func TestVerifyRejectsMissingSignature(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	b := minimalBundle(pub)
	raw, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if _, err := Verify(raw, pub); err == nil {
		t.Fatal("Verify accepted unsigned bundle")
	}
}

func TestVerifyRejectsMalformedJSON(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	if _, err := Verify([]byte("{not json"), pub); err == nil {
		t.Fatal("Verify accepted malformed JSON")
	}
}

func TestVerifyRejectsBadBase64Signature(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	b := minimalBundle(pub)
	signed, err := Sign(b, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	var raw map[string]interface{}
	_ = json.Unmarshal(signed, &raw)
	raw["signature_ed25519"] = "not!!base64@@"
	bogus, _ := json.Marshal(raw)
	if _, err := Verify(bogus, pub); err == nil {
		t.Fatal("Verify accepted non-base64 signature")
	}
}

func TestSignDeterministicForFixedBundle(t *testing.T) {
	// Ed25519 is deterministic per RFC 8032 — identical bundles signed by
	// the same key produce byte-identical output. This guards the canonical
	// encoder against accidental nondeterminism.
	_, priv, _ := ed25519.GenerateKey(rand.Reader)
	pub := priv.Public().(ed25519.PublicKey)

	out1, err := Sign(minimalBundle(pub), priv)
	if err != nil {
		t.Fatalf("Sign #1: %v", err)
	}
	out2, err := Sign(minimalBundle(pub), priv)
	if err != nil {
		t.Fatalf("Sign #2: %v", err)
	}
	if string(out1) != string(out2) {
		t.Fatalf("Sign of identical bundles is not deterministic\n#1: %s\n#2: %s", out1, out2)
	}
}

func TestSignRejectsBadKeyLength(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	b := minimalBundle(pub)
	if _, err := Sign(b, ed25519.PrivateKey(make([]byte, 16))); err == nil {
		t.Fatal("Sign accepted a 16-byte 'private key'")
	}
}

func TestVerifyRejectsBadKeyLength(t *testing.T) {
	if _, err := Verify([]byte(`{}`), ed25519.PublicKey(make([]byte, 8))); err == nil {
		t.Fatal("Verify accepted an 8-byte 'public key'")
	}
}

type failingWriter struct{ errOn int; calls int }

func (w *failingWriter) Write(_ context.Context, _ OperationalEvent) error {
	w.calls++
	if w.calls == w.errOn {
		return errors.New("synthetic write failure")
	}
	return nil
}

func TestIssueReturnsErrorWhenChainWriteFails(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	w := &failingWriter{errOn: 1}

	_, err := Issue(context.Background(), IssueRequest{
		InstitutionID:      "x",
		IssuancePrivateKey: priv,
		IssuancePublicKey:  pub,
		ExaminerIdentity:   "x",
		Tenants:            []string{"t1", "t2"},
		ExaminationStart:   "2026-01-01",
		ExaminationEnd:     "2026-06-14",
		ExpiresAt:          time.Date(2026, 6, 22, 17, 0, 0, 0, time.UTC),
		CorrelationID:      "c",
		Scope:              ScopeComplianceRead,
		ComplianceURL:      "https://example",
		Credential:         BundleCredential{Type: "mtls_client_cert"},
		Now:                func() time.Time { return time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC) },
	}, w)
	if err == nil {
		t.Fatal("Issue returned nil despite ChainWriter failure")
	}
	if !strings.Contains(err.Error(), "synthetic write failure") {
		t.Errorf("error does not mention writer failure: %v", err)
	}
}

func TestRevokeEmitsPerTenantReturnEvents(t *testing.T) {
	rec := &recordingWriter{}
	err := Revoke(context.Background(), RevokeRequest{
		BundleID:      "bnd_3f29b71c",
		Reason:        "examination_concluded",
		CorrelationID: "exam-2026-06-15-occ",
		Tenants:       []string{"tenant_a", "tenant_b"},
		Now:           func() time.Time { return time.Date(2026, 6, 22, 15, 30, 0, 0, time.UTC) },
	}, rec)
	if err != nil {
		t.Fatalf("Revoke: %v", err)
	}
	if len(rec.events) != 2 {
		t.Fatalf("want 2 chain events, got %d", len(rec.events))
	}
	if rec.events[0].Event != EventHandoverReturned {
		t.Errorf("event = %q, want %q", rec.events[0].Event, EventHandoverReturned)
	}
}

func minimalBundle(pub ed25519.PublicKey) *Bundle {
	return &Bundle{
		BundleID:  "bnd_test1234",
		Version:   BundleVersion,
		IssuedAt:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		ExpiresAt: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		Issuer: BundleIssuer{
			InstitutionID:          "test-bank",
			IssuanceKeyFingerprint: Fingerprint([]byte(pub)),
		},
		Examiner: BundleExaminer{Identity: "test-examiner"},
		Scope: BundleScope{
			Surface:                ScopeComplianceRead,
			Tenants:                []string{"tenant_t"},
			ExaminationPeriodStart: "2026-01-01",
			ExaminationPeriodEnd:   "2026-12-30",
			CorrelationID:          "test-corr",
		},
		Endpoint: BundleEndpoint{ComplianceURL: "https://example.test"},
		Cred:     BundleCredential{Type: "mtls_client_cert"},
	}
}
