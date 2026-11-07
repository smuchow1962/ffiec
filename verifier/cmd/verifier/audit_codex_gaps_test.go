package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/mmpworks/ffiec/core/jcs"
	"github.com/mmpworks/ffiec/core/signpayload"
	"github.com/mmpworks/ffiec/testkit/clitest"
)

const devRefused = "Status: FAIL\nStep: 12\nReason: dev-mode seal in production verification — refused\n"

// §4.4 lists these header fields as required, and §10.12 makes a missing one exit 2. The header
// carries no MAC, so anyone who can edit the file can remove one. Removing chain_id used to
// switch off the run half of step 4, and the run read PASS.
func TestHeader_AMissingRequiredFieldIsExit2(t *testing.T) {
	lines := sampleLines(t)
	for _, field := range []string{"format_version", "tenant_id", "chain_id", "genesis_hash", "hkdf_inputs_digest", "opened_at_utc"} {
		t.Run(field, func(t *testing.T) {
			header := decodeLine(t, lines[0])
			delete(header, field)
			line, err := jcs.Canonicalize(header)
			if err != nil {
				t.Fatal(err)
			}
			edited := writeLines(t, append([][]byte{line}, lines[1:]...))
			withKey(edited, sampleSeal, samplePubKey).
				MustExit(t, 2).
				StdoutContains(t, "Status: FAIL\nStep: pre-flight\nReason: header missing required field "+field+"\n")
		})
	}
}

// signedSeal signs a v1.0a seal for the sample file. devMode is bound by the signature.
func signedSeal(t *testing.T, version string, devMode bool) map[string]any {
	t.Helper()
	lines := sampleLines(t)
	af := load(t, lines)
	root := rootOf(t, af)
	parsed, err := signpayload.ParseVersion(version)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := signpayload.Build(signpayload.Seal{
		Version: parsed, Algorithm: "ed25519", FormatVersion: af.Header.FormatVersion, TenantID: af.Header.TenantID,
		SealDate: sampleSealDay, MerkleRootHex: root, HKDFInputsDigest: af.Header.HKDFInputsDigestHex, Cadence: "daily", DevMode: devMode,
	})
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{
		"tenant_id": af.Header.TenantID, "seal_date": sampleSealDay, "format_version": af.Header.FormatVersion,
		"merkle_root_hex": root, "algorithm": "ed25519", "hkdf_inputs_digest_hex": af.Header.HKDFInputsDigestHex,
		"cadence": "daily", "dev_mode": devMode, "key_versions": []int{1}, "sign_payload_version": version,
		"sign_payload_hex": hex.EncodeToString(payload),
		"signature_b64":    base64.StdEncoding.EncodeToString(ed25519.Sign(institutionKey(), payload)),
	}
}

// withProductionURIs gives every entry a key handle that is not a development one. The handle is
// in the stamp, outside the MAC, so the entries stay good.
func withProductionURIs(t *testing.T) string {
	t.Helper()
	lines := sampleLines(t)
	out := [][]byte{lines[0]}
	for _, line := range lines[1:] {
		whole := decodeLine(t, line)
		whole["audit_chain"].(map[string]any)["kms_handle_uri"] = "aws-kms:arn:test"
		canonical, err := jcs.Canonicalize(whole)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, canonical)
	}
	return writeLines(t, out)
}

func strictRun(file string, seal string, extra ...string) []string {
	args := []string{"verify", "--audit-file", file, "--seal", seal, "--root-key", samplePubKey, "--master-key-hex", sampleKey, "--posture", "ed25519", "--strict"}
	return append(args, extra...)
}

// §7 step 12, §4.2 and §10.7: --strict refuses a development seal, and the same run without
// --strict warns. dev_mode is bound by the signature, so nobody can clear it with no key.
func TestStrict_RefusesADevelopmentSeal(t *testing.T) {
	seal := writeSeal(t, signedSeal(t, "v1.0a", true))
	file := withProductionURIs(t)
	clitest.Run(Main, strictRun(file, seal)...).MustExit(t, 1).StdoutContains(t, devRefused)
	withKey(file, seal, samplePubKey).MustSucceed(t).StdoutContains(t, "Status: PASS\n").StdoutContains(t, "Anomaly: seal is dev_mode")
}

// §10.7: --strict also refuses a seal over entries with a plaintext- key handle.
func TestStrict_RefusesADevelopmentKeyHandle(t *testing.T) {
	seal := writeSeal(t, signedSeal(t, "v1.0a", false))
	clitest.Run(Main, strictRun(sampleFile, seal)...).MustExit(t, 1).StdoutContains(t, devRefused)
}

func TestStrict_PassesAProductionSeal(t *testing.T) {
	seal := writeSeal(t, signedSeal(t, "v1.0a", false))
	clitest.Run(Main, strictRun(withProductionURIs(t), seal)...).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutLacks(t, "development key").
		StdoutContains(t, "Anomaly: --strict did not run step 12a").
		StdoutContains(t, "Anomaly: --strict did not compare seal.cadence")
}

// §7's fail-closed rule: a pre-amendment seal's signature binds no dev_mode, so anyone can clear
// that field. --strict cannot confirm a production seal from it, and it refuses.
func TestStrict_RefusesASealThatBindsNoDevMode(t *testing.T) {
	file := withProductionURIs(t)
	preAmendment := writeSeal(t, preAmendmentSeal(t))
	clitest.Run(Main, strictRun(file, preAmendment)...).
		MustExit(t, 1).
		StdoutContains(t, "Status: FAIL\nStep: 12\nReason: seal binds no dev_mode: --strict cannot confirm a production seal\n")
	withKey(file, preAmendment, samplePubKey).MustSucceed(t).StdoutContains(t, "Anomaly: seal has no sign_payload_version")
}

func preAmendmentSeal(t *testing.T) map[string]any {
	t.Helper()
	af := load(t, sampleLines(t))
	root := rootOf(t, af)
	parsed, err := signpayload.ParseVersion("")
	if err != nil {
		t.Fatal(err)
	}
	payload, err := signpayload.Build(signpayload.Seal{
		Version: parsed, Algorithm: "ed25519", FormatVersion: af.Header.FormatVersion, TenantID: af.Header.TenantID,
		SealDate: sampleSealDay, MerkleRootHex: root, HKDFInputsDigest: af.Header.HKDFInputsDigestHex,
	})
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{
		"tenant_id": af.Header.TenantID, "seal_date": sampleSealDay, "format_version": af.Header.FormatVersion,
		"merkle_root_hex": root, "algorithm": "ed25519", "hkdf_inputs_digest_hex": af.Header.HKDFInputsDigestHex,
		"key_versions": []int{1}, "sign_payload_hex": hex.EncodeToString(payload),
		"signature_b64": base64.StdEncoding.EncodeToString(ed25519.Sign(institutionKey(), payload)),
	}
}

// An empty value is a command line error. It must never read as "no floor" or "no posture".
func TestFlags_AnEmptyValueIsExit3(t *testing.T) {
	for _, flag := range []string{"--min-sign-payload-version=", "--posture="} {
		t.Run(flag, func(t *testing.T) {
			sealed(sampleSeal, flag).MustExit(t, 3)
		})
	}
}

// §7 floor check: exit 1 before dispatch, from the flag or from the seal's own field.
func TestFloor_ASealBelowItFails(t *testing.T) {
	seal := signedSeal(t, "v1.0a", true)
	file := writeSeal(t, seal)
	const fromFlag = `Reason: sign_payload_version "v1.0a" below institution-declared floor "v1.0b" (CC8.1-named source: --min-sign-payload-version)`
	sealed(file, "--posture", "ed25519", "--min-sign-payload-version", "v1.0b").
		MustExit(t, 1).
		StdoutContains(t, "Status: FAIL\nStep: 11\n"+fromFlag)

	seal["declared_sign_payload_floor"] = "v1.0b"
	withKey(sampleFile, writeSeal(t, seal), samplePubKey).
		MustExit(t, 1).
		StdoutContains(t, `(CC8.1-named source: declared_sign_payload_floor)`)
}

func TestFloor_AtOrBelowTheSealPasses(t *testing.T) {
	file := writeSeal(t, signedSeal(t, "v1.0a", true))
	sealed(file, "--posture", "ed25519", "--min-sign-payload-version", "v1.0a").MustSucceed(t)
}

// A seal with no sign_payload_version is below every floor. The floor in the seal is unsigned, so
// an unknown value there fails closed. An unknown value on the command line is exit 3.
func TestFloor_EdgeValues(t *testing.T) {
	pre := goodSeal(t)
	pre["declared_sign_payload_floor"] = "v9"
	withKey(sampleFile, writeSeal(t, pre), samplePubKey).MustExit(t, 1).StdoutContains(t, `below institution-declared floor "v9"`)
	sealed(sampleSeal, "--posture", "ed25519", "--min-sign-payload-version", "v9").MustExit(t, 3)
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--min-sign-payload-version", "v1.0b").MustExit(t, 3)
}
