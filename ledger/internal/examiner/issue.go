package examiner

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"time"
)

// IssueRequest captures everything the operator's CLI flags collect. Field
// names mirror the README's `examiner issue` flag names.
type IssueRequest struct {
	InstitutionID      string
	IssuancePrivateKey ed25519.PrivateKey
	IssuancePublicKey  ed25519.PublicKey

	ExaminerIdentity  string
	LaptopFingerprint string

	Tenants          []string
	ExaminationStart string // YYYY-MM-DD
	ExaminationEnd   string // YYYY-MM-DD
	ExpiresAt        time.Time
	CorrelationID    string
	Scope            Scope

	ComplianceURL      string
	TLSCertFingerprint string

	// Credential is the bundle's `credential` block. The mTLS cert and
	// encrypted private key are produced upstream by the operator's PKI
	// tooling; ledgerctl just embeds them.
	Credential BundleCredential

	// Now overrides time.Now for deterministic tests.
	Now func() time.Time
}

// IssueResult is what the CLI prints / returns to the operator.
type IssueResult struct {
	Bundle     *Bundle
	SignedJSON []byte
}

// Issue mints, signs, and emits the issuance chain event. The signed JSON
// is what the operator hands to the examiner.
//
// Three things happen in order:
//
//  1. A bundle id is minted and the Bundle is populated from req.
//  2. The bundle is signed by req.IssuancePrivateKey.
//  3. A `master_key.examiner_handover_received` chain entry is fanned out
//     to every tenant in scope so the issuance is itself sealed evidence.
//
// If step 3 fails, the operator still has a signed bundle but the chain
// has not recorded the issuance. The caller decides whether to abort or
// retry; ledgerctl returns the error and leaves the bundle file unwritten
// to keep the on-disk state consistent with the chain.
func Issue(ctx context.Context, req IssueRequest, w ChainWriter) (*IssueResult, error) {
	if req.Now == nil {
		req.Now = time.Now
	}
	now := req.Now().UTC()

	bundleID, err := NewBundleID()
	if err != nil {
		return nil, err
	}

	pubBytes := []byte(req.IssuancePublicKey)
	if len(pubBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("issuance public key: %d bytes, want %d", len(pubBytes), ed25519.PublicKeySize)
	}

	b := &Bundle{
		BundleID:  bundleID,
		Version:   BundleVersion,
		IssuedAt:  now,
		ExpiresAt: req.ExpiresAt.UTC(),
		Issuer: BundleIssuer{
			InstitutionID:          req.InstitutionID,
			IssuanceKeyFingerprint: Fingerprint(pubBytes),
		},
		Examiner: BundleExaminer{
			Identity:          req.ExaminerIdentity,
			LaptopFingerprint: req.LaptopFingerprint,
		},
		Scope: BundleScope{
			Surface:                req.Scope,
			Tenants:                append([]string(nil), req.Tenants...),
			ExaminationPeriodStart: req.ExaminationStart,
			ExaminationPeriodEnd:   req.ExaminationEnd,
			CorrelationID:          req.CorrelationID,
		},
		Endpoint: BundleEndpoint{
			ComplianceURL:      req.ComplianceURL,
			TLSCertFingerprint: req.TLSCertFingerprint,
		},
		Cred: req.Credential,
	}

	signed, err := Sign(b, req.IssuancePrivateKey)
	if err != nil {
		return nil, err
	}

	evt := OperationalEvent{
		Event:         EventHandoverReceived,
		Timestamp:     now,
		CorrelationID: req.CorrelationID,
		Fields: map[string]interface{}{
			"examiner_identity":                req.ExaminerIdentity,
			"examination_period_start":         req.ExaminationStart,
			"examination_period_end":           req.ExaminationEnd,
			"handover_channel":                 string(req.Scope),
			"handover_authorization_record_id": bundleID,
			"expires_at":                       req.ExpiresAt.UTC().Format(time.RFC3339),
		},
	}
	if err := MultiTenantWrite(ctx, w, evt, req.Tenants); err != nil {
		return nil, fmt.Errorf("emit handover-received event: %w", err)
	}

	return &IssueResult{Bundle: b, SignedJSON: signed}, nil
}
