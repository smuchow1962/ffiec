// Package examiner implements credential-bundle issuance and revocation
// for examination read-only access to the Compliance read surface.
//
// See ledger/cmd/ledgerctl/README.md for design context. The chain-event
// shape that brackets every issuance/revocation matches the receipt and
// disposal events documented in docs/operator-guide.md.
package examiner

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Scope names the access surface a Bundle authorises.
type Scope string

const (
	// ScopeComplianceRead is the read-only Compliance surface case from
	// docs/auditor-stories/01-northbridge-federal-savings.md.
	ScopeComplianceRead Scope = "compliance.read"

	// ScopeMasterKeyRead is the heavier Approach C master-key.read case from
	// docs/operator-guide.md §"Examination-time master-key handover procedure".
	ScopeMasterKeyRead Scope = "master-key.read"
)

// BundleVersion is the on-disk format version. v1.0 is the only supported
// version today; the verifier rejects anything else.
const BundleVersion = "1.0"

// Bundle is the credential a ledgerctl operator hands an examiner. Once
// signed, the document is self-contained — the examiner's verifier checks
// the signature against the institution's published issuance public key
// before the bundle is trusted for any fetch operation.
type Bundle struct {
	BundleID  string    `json:"bundle_id"`
	Version   string    `json:"version"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`

	Issuer   BundleIssuer     `json:"issuer"`
	Examiner BundleExaminer   `json:"examiner"`
	Scope    BundleScope      `json:"scope"`
	Endpoint BundleEndpoint   `json:"endpoint"`
	Cred     BundleCredential `json:"credential"`

	// SignatureEd25519 is base64-encoded. The signing input excludes this
	// field; canonicalForSigning clears it before encoding.
	SignatureEd25519 string `json:"signature_ed25519,omitempty"`
}

// BundleIssuer identifies the institution and the issuance keypair under
// which the bundle was signed.
type BundleIssuer struct {
	InstitutionID          string `json:"institution_id"`
	IssuanceKeyFingerprint string `json:"issuance_key_fingerprint"`
}

// BundleExaminer identifies the receiving examiner and an optional laptop
// pin. The pin is advisory — the verifier MAY enforce it, but verification
// of pulled ledger bytes does not depend on it.
type BundleExaminer struct {
	Identity          string `json:"identity"`
	LaptopFingerprint string `json:"laptop_fingerprint,omitempty"`
}

// BundleScope declares what the bundle authorises and for how long. Surface
// names a single read surface; Tenants enumerates the tenant chains the
// examiner may fetch from.
type BundleScope struct {
	Surface                Scope    `json:"surface"`
	Tenants                []string `json:"tenants"`
	ExaminationPeriodStart string   `json:"examination_period_start"` // YYYY-MM-DD
	ExaminationPeriodEnd   string   `json:"examination_period_end"`   // YYYY-MM-DD
	CorrelationID          string   `json:"correlation_id"`
}

// BundleEndpoint is where the verifier fetches from. The optional TLS cert
// fingerprint pins the server certificate; mismatch fails closed at import.
type BundleEndpoint struct {
	ComplianceURL      string `json:"compliance_url"`
	TLSCertFingerprint string `json:"tls_cert_fingerprint,omitempty"`
}

// BundleCredential is the bytes the verifier presents at the read surface.
// The mTLS cert and encrypted private key are produced upstream by the
// operator's PKI tooling; ledgerctl just embeds them.
type BundleCredential struct {
	Type            string `json:"type"`
	CertPEM         string `json:"cert_pem,omitempty"`
	EncryptedKeyPEM string `json:"encrypted_key_pem,omitempty"`
}

// Validate checks the structural invariants the issuance flow depends on.
// Cryptographic validity is checked separately via Verify.
func (b *Bundle) Validate() error {
	if b.BundleID == "" {
		return fmt.Errorf("bundle_id is required")
	}
	if b.Version != BundleVersion {
		return fmt.Errorf("version %q not supported (expected %q)", b.Version, BundleVersion)
	}
	if b.IssuedAt.IsZero() {
		return fmt.Errorf("issued_at is required")
	}
	if b.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at is required")
	}
	if !b.ExpiresAt.After(b.IssuedAt) {
		return fmt.Errorf("expires_at must be after issued_at")
	}
	if b.Issuer.InstitutionID == "" {
		return fmt.Errorf("issuer.institution_id is required")
	}
	if b.Issuer.IssuanceKeyFingerprint == "" {
		return fmt.Errorf("issuer.issuance_key_fingerprint is required")
	}
	if b.Examiner.Identity == "" {
		return fmt.Errorf("examiner.identity is required")
	}
	if len(b.Scope.Tenants) == 0 {
		return fmt.Errorf("scope.tenants must list at least one tenant")
	}
	switch b.Scope.Surface {
	case ScopeComplianceRead, ScopeMasterKeyRead:
		// known scope
	default:
		return fmt.Errorf("scope.surface %q is not a known scope", b.Scope.Surface)
	}
	if b.Scope.CorrelationID == "" {
		return fmt.Errorf("scope.correlation_id is required")
	}
	if b.Endpoint.ComplianceURL == "" {
		return fmt.Errorf("endpoint.compliance_url is required")
	}
	if b.Cred.Type == "" {
		return fmt.Errorf("credential.type is required")
	}
	return nil
}

// NewBundleID returns a freshly random bundle id of the form "bnd_<8hex>".
// The id is opaque to the wire — collision risk for an institution issuing
// a few credentials per examination is negligible at 32 bits.
func NewBundleID() (string, error) {
	var raw [4]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("randomness: %w", err)
	}
	return "bnd_" + hex.EncodeToString(raw[:]), nil
}

// Fingerprint returns "sha256:<hex>" of the supplied bytes. Used for the
// issuance public key fingerprint that lands in Issuer.IssuanceKeyFingerprint.
func Fingerprint(b []byte) string {
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}
