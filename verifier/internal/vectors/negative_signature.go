package vectors

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// §7 step-11 live signature walk for the negative corpus.
//
// N004 (signature garbage) and N005 (signature wrong tenant) are the two
// negative vectors whose pinned failure is a §4.3 Ed25519 signature
// rejection. This file is the live driver: when the fixture carries a real
// (base64-decodable) signature AND a published public key is resolvable,
// it runs the verifier's own WalkAuditFileWithKey over the tampered
// audit_file and asserts the verifier emits the pinned FAIL/11/"signature
// verification failed"/exit-1.
//
// Live-readiness gate. Heather is materializing N004/N005 as real
// signature-bearing fixtures signed with the corpus test key. Until BOTH
// the signature decodes to a real Ed25519 form AND the public key
// resolves, the vector stays contract-only — driving a live walk over a
// placeholder signature would assert the right exit code for the WRONG
// reason (a base64-decode error rather than the signature/tenant-binding
// rejection the vector pins). sealSignatureLiveReady is that gate; a
// not-ready vector falls back to the reason-template self-consistency
// check, exactly as before, so the gate re-arms automatically when the
// real fixture lands.

// publicKeyHexFile is the corpus-published public key (the auditor-
// reproducible source). When the corpus publishes it, the verifier reads
// it from there; otherwise the live driver falls back to the public key
// under TESSERASEAL_TEST_KEY_DIR. Either way the verifier reads PUBLIC
// material only.
const publicKeyHexFile = "test-signing-key.pub.hex"

// testKeyDirEnvVar overrides the local-only test-key directory the corpus
// generators and the verifier's signature path read by convention.
const testKeyDirEnvVar = "TESSERASEAL_TEST_KEY_DIR"

// defaultTestKeyDir is the local-only key directory the public key falls
// back to when the corpus has not yet published pub.hex.
const defaultTestKeyDir = `E:\dev\testing\private-keys\tesseraseal`

// assertSealSignatureWalk drives the live §7 step-11 signature walk over
// the N004/N005 fixture when it is live-ready, or falls back to the
// contract-only reason-template check when it is not. The returned Check's
// name carries the "/§7-step-11-signature" suffix on the live path so the
// honesty guard can prove the assertion was live.
func assertSealSignatureWalk(v NegativeVector, exp expectedOutput) Check {
	name := v.Slot + "/§7-step-11-signature"

	af, ikms, err := loadAuditFile(v.Dir)
	if err != nil {
		return failCheck(name, "%v", err)
	}
	pub, keyErr := resolveCorpusPublicKey(v.Dir)
	if keyErr != nil || !sealSignatureLiveReady(af.Seal, pub) {
		// Not live-ready: keep the contract-only assertion so the gate
		// stays green and re-arms when the real fixture + key land.
		return runMaterializedNegative(v)
	}

	got := verify.WalkAuditFileWithKey(af, ikms, pub)
	return matchOutcome(name, got, exp)
}

// sealSignatureLiveReady reports whether the seal carries enough to drive
// a genuine step-11 crypto rejection: a base64-decodable signature (so the
// verify is a real Ed25519 check, not a decode error) AND a resolvable
// public key. The signature need NOT be 64 bytes — N004's garbage decodes
// to 66 bytes and is a legitimate live rejection at the Ed25519.Verify
// call. The bar is "decodable", which excludes only the textual
// placeholder (N005's pre-materialization state).
func sealSignatureLiveReady(seal verify.AuditSeal, pub ed25519.PublicKey) bool {
	if len(pub) != ed25519.PublicKeySize {
		return false
	}
	if seal.SignPayloadHex == "" {
		return false
	}
	if _, err := base64.StdEncoding.DecodeString(seal.SignatureB64); err != nil {
		return false
	}
	return true
}

// resolveCorpusPublicKey loads the institution public key the live
// signature walk verifies against. It prefers the corpus-published pub.hex
// (auditor-reproducible, travels with the vectors); when the corpus has
// not yet published it, it falls back to the public key under
// TESSERASEAL_TEST_KEY_DIR. The verifier reads PUBLIC material only — never
// the private seed.
//
// Search order:
//  1. <corpusDir>/test-signing-key.pub.hex   (published with the corpus)
//  2. <vectorDir>/test-signing-key.pub.hex   (per-vector publication)
//  3. $TESSERASEAL_TEST_KEY_DIR/test-signing-key.pub.hex (or the default)
func resolveCorpusPublicKey(vectorDir string) (ed25519.PublicKey, error) {
	candidates := publicKeyCandidatePaths(vectorDir)
	for _, p := range candidates {
		raw, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		pub, err := verify.ParsePublicKeyHex(strings.TrimSpace(string(raw)))
		if err != nil {
			return nil, fmt.Errorf("public key at %s: %w", p, err)
		}
		return pub, nil
	}
	return nil, fmt.Errorf("no published public key found in %v", candidates)
}

// publicKeyCandidatePaths is the ordered list resolveCorpusPublicKey
// searches: corpus root, then the vector directory, then the local-only
// test-key directory (env-var override or default).
func publicKeyCandidatePaths(vectorDir string) []string {
	// vectorDir is <corpus>/negative/<slot>; the corpus root is two up.
	corpusRoot := filepath.Dir(filepath.Dir(vectorDir))
	keyDir := os.Getenv(testKeyDirEnvVar)
	if keyDir == "" {
		keyDir = defaultTestKeyDir
	}
	return []string{
		filepath.Join(corpusRoot, publicKeyHexFile),
		filepath.Join(vectorDir, publicKeyHexFile),
		filepath.Join(keyDir, publicKeyHexFile),
	}
}
