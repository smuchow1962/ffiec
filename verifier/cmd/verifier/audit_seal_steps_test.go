package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/core/signpayload"
	"github.com/mmpworks/ffiec/testkit/clitest"
)

const (
	keyVersionsFail = "Status: FAIL\nStep: 11\nReason: seal.key_versions does not match per-event key_version distribution\n"
	algorithmFail   = "Status: FAIL\nStep: 11\nReason: algorithm/key-type mismatch at signature verification\n"

	signaturesListRefused = "Status: FAIL\nStep: 11\nReason: seal carries a signatures list, and this verifier does not check one yet: dual-algorithm checks are not built\n"
)

func text(seal map[string]any, field string) string {
	s, _ := seal[field].(string)
	return s
}

// resign signs seal's fields as they now stand, under its own sign_payload_version.
func resign(t *testing.T, seal map[string]any, priv ed25519.PrivateKey) map[string]any {
	t.Helper()
	version, err := signpayload.ParseVersion(text(seal, "sign_payload_version"))
	if err != nil {
		t.Fatal(err)
	}
	devMode, _ := seal["dev_mode"].(bool)
	payload, err := signpayload.Build(signpayload.Seal{
		Version: version, Algorithm: text(seal, "algorithm"), FormatVersion: text(seal, "format_version"),
		TenantID: text(seal, "tenant_id"), SealDate: text(seal, "seal_date"), MerkleRootHex: text(seal, "merkle_root_hex"),
		HKDFInputsDigest: text(seal, "hkdf_inputs_digest_hex"), Cadence: text(seal, "cadence"), DevMode: devMode,
		KeyVersionsCanon: text(seal, "key_versions_canon"), KMSHandleURIsDigest: text(seal, "kms_handle_uris_digest_hex"),
		OperationalEventsLogRoot: text(seal, "operational_events_log_root_hex"),
	})
	if err != nil {
		t.Fatal(err)
	}
	seal["sign_payload_hex"] = hex.EncodeToString(payload)
	seal["signature_b64"] = base64.StdEncoding.EncodeToString(ed25519.Sign(priv, payload))
	return seal
}

// sampleSealAs returns the sample's seal with changes applied and signed by the institution key.
func sampleSealAs(t *testing.T, changes map[string]any, remove ...string) string {
	t.Helper()
	lines := sampleLines(t)
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	for k, v := range changes {
		seal[k] = v
	}
	for _, k := range remove {
		delete(seal, k)
	}
	return writeSeal(t, resign(t, seal, institutionKey()))
}

// goodV1_0b holds the true values for the sample: 36 entries under key version 1, all plaintext-dev.
func goodV1_0b() map[string]any {
	digest := sha256.Sum256([]byte("plaintext-dev"))
	return map[string]any{
		"sign_payload_version":       "v1.0b",
		"key_versions_canon":         "1",
		"kms_handle_uris_digest_hex": hex.EncodeToString(digest[:]),
	}
}

func with(base map[string]any, field string, value any) map[string]any {
	out := map[string]any{field: value}
	for k, v := range base {
		if k != field {
			out[k] = v
		}
	}
	return out
}

func TestSeal_UnknownVersionUsesTheSpecReason(t *testing.T) {
	lines := sampleLines(t)
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	seal["sign_payload_version"] = "v1.0d"
	withKey(sampleFile, writeSeal(t, seal), samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, "Status: FAIL\nStep: 11\nReason: sign_payload_version \"v1.0d\" not supported by this verifier (running v1.0c)\n")
}

func TestSeal_V1_0bTrueValuesPass(t *testing.T) {
	withKey(sampleFile, sampleSealAs(t, goodV1_0b(), "key_versions"), samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutLacks(t, "key_versions")
}

// A genuine signature over a false list. The institution signed it, so only the entries can show it false.
func TestSeal_V1_0bFalseKeyVersionsFail(t *testing.T) {
	withKey(sampleFile, sampleSealAs(t, with(goodV1_0b(), "key_versions_canon", "99"), "key_versions"), samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, keyVersionsFail)
}

// §4.3: a flip in the day's kms_handle_uri values "produces a signature failure".
func TestSeal_V1_0bFalseKMSDigestFails(t *testing.T) {
	withKey(sampleFile, sampleSealAs(t, with(goodV1_0b(), "kms_handle_uris_digest_hex", strings.Repeat("ab", sha256.Size)), "key_versions"), samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, signatureFail)
}

func TestSeal_KeyVersionsListIsCrossChecked(t *testing.T) {
	withKey(sampleFile, sampleSealAs(t, map[string]any{"key_versions": []int{1, 99}}), samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, keyVersionsFail)
}

func TestSeal_NoKeyVersionsSaysTheCheckDidNotRun(t *testing.T) {
	withKey(sampleFile, sampleSealAs(t, nil, "key_versions"), samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Anomaly: seal names no key_versions, so that cross-check did not run: see §7 step 11")
}

// §10.79: a standard walk records the bound root and does not walk it. So PASS stands, and the output says so.
func TestSeal_V1_0cSaysTheOperationalRootWasNotWalked(t *testing.T) {
	bogus := strings.Repeat("cd", sha256.Size)
	withKey(sampleFile, sampleSealAs(t, with(with(goodV1_0b(), "sign_payload_version", "v1.0c"), "operational_events_log_root_hex", bogus), "key_versions"), samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutContains(t, "Anomaly: operational_events_log_root "+bogus+" is bound by the signature and was not walked: see §10.79")
}

func TestSeal_NoVersionSaysWhatTheSignatureDoesNotBind(t *testing.T) {
	withKey(sampleFile, sampleSealAs(t, nil, "sign_payload_version"), samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Anomaly: seal has no sign_payload_version, so its signature does not bind cadence or dev_mode: see §4.3")
}

func TestSeal_AlgorithmItCannotVerifyFails(t *testing.T) {
	withKey(sampleFile, sampleSealAs(t, map[string]any{"algorithm": "dilithium3"}), samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, algorithmFail)
}

// This verifier checks one signature. A seal that lists several must never read PASS, because
// a forged seal could list a bad second signature, or leave one out, and nothing here would see it.
func TestSeal_CarryingASignaturesListIsRefused(t *testing.T) {
	lines := sampleLines(t)
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	second := map[string]any{"algorithm": "ml-dsa-65", "public_key_id": "k2", "signature": base64.StdEncoding.EncodeToString(make([]byte, 64))}
	first := map[string]any{"algorithm": "ed25519", "public_key_id": "k1", "signature": seal["signature_b64"]}
	for name, list := range map[string]any{
		"a bad second signature": []any{first, second},
		"one signature":          []any{first},
		"an empty list":          []any{},
		"null":                   nil,
	} {
		t.Run(name, func(t *testing.T) {
			withList := map[string]any{"signatures": list}
			for k, v := range seal {
				withList[k] = v
			}
			withKey(sampleFile, writeSeal(t, withList), samplePubKey).
				MustExit(t, 3).
				StdoutContains(t, signaturesListRefused)
		})
	}
}

// A seal holding only the list has no single signature to check, and still reads as refused.
func TestSeal_HoldingOnlyASignaturesListIsRefused(t *testing.T) {
	lines := sampleLines(t)
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	seal["signatures"] = []any{map[string]any{"algorithm": "ed25519", "signature": seal["signature_b64"]}}
	delete(seal, "signature_b64")
	withKey(sampleFile, writeSeal(t, seal), samplePubKey).MustExit(t, 3).StdoutContains(t, signaturesListRefused)
}

// Someone with no signing key cuts an event and recomputes the root. Whatever they do with the
// signature fields, the forgery stays exit 1: nothing in an unsigned file may move it to exit 3.
func TestSeal_ForgedRootStaysAFindingWithAnyList(t *testing.T) {
	lines := sampleLines(t)
	cut := lines[:len(lines)-1]
	forger := seededKey("the writer's own")
	for name, list := range map[string]any{
		"an empty list":              []any{},
		"null":                       nil,
		"a forged ed25519 entry":     []any{map[string]any{"algorithm": "ed25519", "signature": "use-forged"}},
		"no ed25519 entry":           []any{map[string]any{"algorithm": "ml-dsa-65", "signature": "AAAA"}},
		"an entry with no signature": []any{map[string]any{"algorithm": "ed25519"}},
		"not a list":                 "ed25519",
	} {
		t.Run(name, func(t *testing.T) {
			seal := sealFor(t, cut, load(t, cut).Header.TenantID, forger)
			forged := seal["signature_b64"]
			delete(seal, "signature_b64")
			if entries, ok := list.([]any); ok && len(entries) == 1 {
				if entry := entries[0].(map[string]any); entry["signature"] == "use-forged" {
					entry["signature"] = forged
				}
			}
			seal["signatures"] = list
			sealFile := writeSeal(t, seal)
			file := writeLines(t, cut)
			withKey(file, sealFile, samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)
			clitest.Run(Main, "verify", "--audit-file", file, "--seal", sealFile, "--root-key", samplePubKey).
				MustExit(t, 1).
				StdoutContains(t, signatureFail)
		})
	}
}

// A valid Ed25519 signature filed under another algorithm's name is not the ed25519 entry.
func TestSeal_ListEntryCountsOnlyUnderItsOwnAlgorithm(t *testing.T) {
	lines := sampleLines(t)
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	seal["signatures"] = []any{map[string]any{"algorithm": "ml-dsa-65", "signature": seal["signature_b64"]}}
	delete(seal, "signature_b64")
	withKey(sampleFile, writeSeal(t, seal), samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)
}

// One unsigned field must not turn an integrity finding into "the verifier cannot judge this".
func TestSeal_SignaturesListHidesNoFinding(t *testing.T) {
	lines := sampleLines(t)
	list := []any{}
	seal := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	seal["signatures"] = list
	withKey(writeLines(t, lines[:len(lines)-1]), writeSeal(t, seal), samplePubKey).MustExit(t, 1).StdoutContains(t, rootMismatch)

	badSignature := sealFor(t, lines, load(t, lines).Header.TenantID, seededKey("the writer's own"))
	badSignature["signatures"] = list
	withKey(sampleFile, writeSeal(t, badSignature), samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)

	otherTenant := sealFor(t, lines, "another-tenant", institutionKey())
	otherTenant["signatures"] = list
	withKey(sampleFile, writeSeal(t, otherTenant), samplePubKey).MustExit(t, 1).StdoutContains(t, signatureFail)

	cut := writeLines(t, lines[:len(lines)-1])
	noSingle := sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
	noSingle["signatures"] = list
	delete(noSingle, "signature_b64")
	withKey(cut, writeSeal(t, noSingle), samplePubKey).MustExit(t, 1).StdoutContains(t, rootMismatch)
}

func TestSeal_EmptySealIsStructural(t *testing.T) {
	empty := filepath.Join(t.TempDir(), "seal.json")
	if err := os.WriteFile(empty, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	withKey(sampleFile, empty, samplePubKey).MustExit(t, 2).StdoutContains(t, "Step: pre-flight\nReason: seal names no merkle_root_hex\n")
}

func TestSeal_NeedsTheAuditFile(t *testing.T) {
	clitest.Run(Main, "verify", "--ledger", "any.ledger", "--root-key", samplePubKey, "--seal", sampleSeal).
		MustFail(t).
		StderrContains(t, "--seal")
}

// The seal's date and tenant go to stderr, because §7 fixes what stdout holds.
func TestSeal_NamesWhatItCheckedOnStderr(t *testing.T) {
	withKey(sampleFile, sampleSeal, samplePubKey).
		MustSucceed(t).
		StderrContains(t, "Seal checked: tenant_id PUMP-01, seal_date "+sampleSealDay+"\n").
		StdoutLacks(t, "Seal checked")
}
