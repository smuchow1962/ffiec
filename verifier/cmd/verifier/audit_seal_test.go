package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/core/constants"
	"github.com/mmpworks/ffiec/core/hkdf"
	"github.com/mmpworks/ffiec/core/jcs"
	"github.com/mmpworks/ffiec/core/signpayload"
	"github.com/mmpworks/ffiec/testkit/clitest"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// The published sample's seal. Its signing key comes from a published seed, so the
// seal proves nothing about any institution: it exists so a reader can run one
// plain command, and so these tests can sign what they need.
const (
	sampleSeal    = "testdata/simulator-custody-seal.json"
	samplePubKey  = "testdata/simulator-custody-seal-pub.pem"
	sampleSealDay = "2026-11-04"

	rootMismatch  = "Status: FAIL\nStep: 10\nReason: merkle root mismatch — ledger contents do not produce sealed root\n"
	signatureFail = "Status: FAIL\nStep: 11\nReason: signature verification failed\n"
)

var updateSealFixture = flag.Bool("update-seal", false, "rewrite the sample seal and its public key in testdata")

func seededKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("ffiec verifier test seal key: " + label))
	return ed25519.NewKeyFromSeed(seed[:])
}

func institutionKey() ed25519.PrivateKey { return seededKey("simulator-custody sample") }

// sealFor signs a seal over the entries in lines, for tenant.
func sealFor(t *testing.T, lines [][]byte, tenant string, priv ed25519.PrivateKey) map[string]any {
	t.Helper()
	return sealWithDigest(t, lines, tenant, load(t, lines).Header.HKDFInputsDigestHex, priv)
}

// sealWithDigest signs a seal that names digest as its HKDF inputs.
func sealWithDigest(t *testing.T, lines [][]byte, tenant, digest string, priv ed25519.PrivateKey) map[string]any {
	t.Helper()
	af := load(t, lines)
	version, err := signpayload.ParseVersion("v1.0a")
	if err != nil {
		t.Fatal(err)
	}
	root := rootOf(t, af)
	payload, err := signpayload.Build(signpayload.Seal{
		Version: version, Algorithm: "ed25519", FormatVersion: af.Header.FormatVersion, TenantID: tenant,
		SealDate: sampleSealDay, MerkleRootHex: root, HKDFInputsDigest: digest, Cadence: "daily", DevMode: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{
		"tenant_id":              tenant,
		"seal_date":              sampleSealDay,
		"format_version":         af.Header.FormatVersion,
		"merkle_root_hex":        root,
		"algorithm":              "ed25519",
		"hkdf_inputs_digest_hex": digest,
		"cadence":                "daily",
		"dev_mode":               true,
		"sign_payload_version":   "v1.0a",
		"sign_payload_hex":       hex.EncodeToString(payload),
		"signature_b64":          base64.StdEncoding.EncodeToString(ed25519.Sign(priv, payload)),
	}
}

func load(t *testing.T, lines [][]byte) *verify.AuditFile {
	t.Helper()
	af, out := verify.LoadAuditNDJSON(append(bytes.Join(lines, []byte("\n")), '\n'))
	if out.Status == "FAIL" {
		t.Fatal(out.Reason)
	}
	return af
}

func rootOf(t *testing.T, af *verify.AuditFile) string {
	t.Helper()
	leaves := make([][]byte, 0, len(af.Entries))
	for i := range af.Entries {
		tag, err := hex.DecodeString(af.Entries[i].PayloadHashHex)
		if err != nil {
			t.Fatal(err)
		}
		leaves = append(leaves, verify.MerkleLeafHash(tag))
	}
	return hex.EncodeToString(verify.MerkleTreeHash(leaves))
}

// rewriteWithKey does what a holder of the master key can do: it changes the text
// of one event, then makes a new tag and link for that event and each one after it.
func rewriteWithKey(t *testing.T, lines [][]byte) [][]byte {
	t.Helper()
	const editedLine = 3
	af := load(t, lines)
	keyHex, err := os.ReadFile(sampleKey)
	if err != nil {
		t.Fatal(err)
	}
	ikm, err := hex.DecodeString(strings.TrimSpace(string(keyHex)))
	if err != nil {
		t.Fatal(err)
	}
	info := []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + af.Header.TenantID)
	key := hkdf.Derive(ikm, []byte(constants.HKDFSalt), info, constants.HKDFOutputLength)

	out := append([][]byte(nil), lines...)
	prev := make([]byte, sha256.Size)
	for n := 1; n < len(lines); n++ {
		whole := decodeLine(t, lines[n])
		if n == editedLine {
			whole["message"] = "this event's text was changed by the key holder"
		}
		stamp := whole["audit_chain"].(map[string]any)
		delete(whole, "audit_chain")
		canonical, err := jcs.Canonicalize(whole)
		if err != nil {
			t.Fatal(err)
		}
		mac := hmac.New(sha256.New, key)
		mac.Write(prev)
		mac.Write(canonical)
		tag := mac.Sum(nil)
		stamp["prev_hash"] = base64.StdEncoding.EncodeToString(prev)
		stamp["payload_hash"] = base64.StdEncoding.EncodeToString(tag)
		whole["audit_chain"] = stamp
		if out[n], err = jcs.Canonicalize(whole); err != nil {
			t.Fatal(err)
		}
		prev = tag
	}
	return out
}

func decodeLine(t *testing.T, line []byte) map[string]any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(line))
	decoder.UseNumber()
	var whole map[string]any
	if err := decoder.Decode(&whole); err != nil {
		t.Fatal(err)
	}
	return whole
}

func writeSeal(t *testing.T, seal map[string]any) string {
	t.Helper()
	raw, err := json.MarshalIndent(seal, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "seal.json")
	if err := os.WriteFile(p, append(raw, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func writePub(t *testing.T, priv ed25519.PrivateKey) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "pub.pem")
	if err := cliutil.WriteEd25519PublicPEM(p, priv.Public().(ed25519.PublicKey)); err != nil {
		t.Fatal(err)
	}
	return p
}

func withKey(file, seal, pub string) clitest.Result {
	return clitest.Run(Main, "verify", "--audit-file", file, "--seal", seal, "--root-key", pub, "--master-key-hex", sampleKey)
}

// TestSampleSealIsCurrent keeps the committed sample honest: the seal in testdata
// must be the one this file's rules produce. Run with -update-seal to rewrite it.
func TestSampleSealIsCurrent(t *testing.T) {
	seal := sealFor(t, sampleLines(t), load(t, sampleLines(t)).Header.TenantID, institutionKey())
	want, err := json.MarshalIndent(seal, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want = append(want, '\n')
	if *updateSealFixture {
		if err := os.WriteFile(sampleSeal, want, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := cliutil.WriteEd25519PublicPEM(samplePubKey, institutionKey().Public().(ed25519.PublicKey)); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(sampleSeal)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s is stale; run go test -run TestSampleSealIsCurrent -args -update-seal", sampleSeal)
	}
}

func TestSeal_IntactChainPasses(t *testing.T) {
	withKey(sampleFile, sampleSeal, samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutLacks(t, "no daily seal supplied").
		StdoutContains(t, "Anomaly: seal is dev_mode")
}

// A cut at the end leaves every remaining link good. Only the seal sees it.
func TestSeal_LastEventCutFailsAtTheRoot(t *testing.T) {
	lines := sampleLines(t)
	cut := writeLines(t, lines[:len(lines)-1])
	clitest.Run(Main, "verify", "--audit-file", cut, "--master-key-hex", sampleKey).
		MustSucceed(t).
		StdoutContains(t, "Anomaly: no daily seal supplied")
	withKey(cut, sampleSeal, samplePubKey).MustExit(t, 1).StdoutContains(t, rootMismatch)
}

// The rewrite passes steps 1-9, because the key holder made each tag. The seal stops it.
func TestSeal_RewriteWithTheKeyFailsAtTheRoot(t *testing.T) {
	rewritten := writeLines(t, rewriteWithKey(t, sampleLines(t)))
	clitest.Run(Main, "verify", "--audit-file", rewritten, "--master-key-hex", sampleKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n")
	withKey(rewritten, sampleSeal, samplePubKey).MustExit(t, 1).StdoutContains(t, rootMismatch)
}

func TestSeal_ReplacedRootWithTheOldSignatureFails(t *testing.T) {
	lines := rewriteWithKey(t, sampleLines(t))
	seal := sealFor(t, sampleLines(t), load(t, lines).Header.TenantID, institutionKey())
	seal["merkle_root_hex"] = rootOf(t, load(t, lines))
	withKey(writeLines(t, lines), writeSeal(t, seal), samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)
}

// The writer can sign a good seal with a key of their own, and can even name that
// key inside the seal. The verifier checks only the key file its caller supplied.
func TestSeal_OwnSignatureFailsUnderTheInstitutionKey(t *testing.T) {
	lines := rewriteWithKey(t, sampleLines(t))
	own := seededKey("the writer's own")
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, own)
	ownPub := own.Public().(ed25519.PublicKey)
	seal["public_key_id"] = hex.EncodeToString(ownPub)
	seal["public_key_hex"] = hex.EncodeToString(ownPub)
	withKey(writeLines(t, lines), writeSeal(t, seal), samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)
}

// The limit, pinned: the same file passes if the reader trusts the writer's key.
// So the public key must reach the reader from outside the file's holder.
func TestSeal_OwnSignaturePassesUnderTheWritersKey(t *testing.T) {
	lines := rewriteWithKey(t, sampleLines(t))
	own := seededKey("the writer's own")
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, own)
	withKey(writeLines(t, lines), writeSeal(t, seal), writePub(t, own)).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n")
}

// A genuine signature over this file's root, for another tenant, does not cover this file.
func TestSeal_ForAnotherTenantFails(t *testing.T) {
	seal := sealFor(t, sampleLines(t), "another-tenant", institutionKey())
	withKey(sampleFile, writeSeal(t, seal), samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)
}

func TestSeal_ForOtherHKDFInputsFails(t *testing.T) {
	lines := sampleLines(t)
	seal := sealWithDigest(t, lines, load(t, lines).Header.TenantID, strings.Repeat("ab", sha256.Size), institutionKey())
	withKey(sampleFile, writeSeal(t, seal), samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)
}

func TestSeal_WithNoSignatureFails(t *testing.T) {
	seal := sealFor(t, sampleLines(t), load(t, sampleLines(t)).Header.TenantID, institutionKey())
	delete(seal, "signature_b64")
	withKey(sampleFile, writeSeal(t, seal), samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)
}

// Steps 10 and 11 need no master key, so a reader with no key still gets them.
func TestSeal_WitnessModeRunsTheSealSteps(t *testing.T) {
	lines := sampleLines(t)
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--seal", sampleSeal, "--root-key", samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Status: "+verify.WitnessStatus+"\n").
		StdoutLacks(t, "no daily seal supplied")
	clitest.Run(Main, "verify", "--audit-file", writeLines(t, lines[:len(lines)-1]), "--seal", sampleSeal, "--root-key", samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, rootMismatch)
}

func TestSeal_NeedsARootKeyAndARootKeyNeedsASeal(t *testing.T) {
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--seal", sampleSeal).MustExit(t, 3)
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--root-key", samplePubKey).MustExit(t, 3)
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--seal", sampleSeal, "--root-key", "no-such-key.pem").MustExit(t, 3)
}

func TestSeal_UnreadableSealIsStructural(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "seal.json")
	if err := os.WriteFile(bad, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, seal := range []string{bad, "no-such-seal.json"} {
		withKey(sampleFile, seal, samplePubKey).MustExit(t, 2).StdoutContains(t, "Step: pre-flight\n")
	}
}
