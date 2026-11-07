package main

import (
	"fmt"
	"testing"

	"github.com/mmpworks/ffiec/testkit/clitest"
)

const (
	dualPosture   = "ed25519,dilithium3"
	partialReason = "partial-coverage seal: single-algorithm signature during institution's declared dual-algorithm posture"
)

// goodSeal is a seal for the sample file, signed by the institution's key.
func goodSeal(t *testing.T) map[string]any {
	t.Helper()
	lines := sampleLines(t)
	return sealFor(t, lines, load(t, lines).Header.TenantID, institutionKey())
}

// sealed runs the keyed walk with the given extra flags and no default posture.
func sealed(seal string, extra ...string) clitest.Result {
	args := []string{"verify", "--audit-file", sampleFile, "--seal", seal, "--root-key", samplePubKey, "--master-key-hex", sampleKey}
	return clitest.Run(Main, append(args, extra...)...)
}

func noPostureReason(seal map[string]any) string {
	return fmt.Sprintf("algorithm posture not supplied at seal_date %s: a stripped signature cannot be detected", seal["seal_date"])
}

func notOnListReason(seal map[string]any) string {
	return fmt.Sprintf("algorithm not on institution's declared posture list at seal_date %s", seal["seal_date"])
}

// passWithAnomaly is the shape vector N017 pins: the status, then the step, the reason and
// the anomaly line, all on stdout, with exit 0.
func passWithAnomaly(reason string) string {
	return "Status: PASS-WITH-ANOMALY\nStep: 11\nReason: " + reason + "\nAnomaly: " + reason + "\n"
}

// §7 step 11 case (c): non-strict with no posture source never reads a clean PASS.
func TestPosture_NoSourceNonStrictWarns(t *testing.T) {
	seal := goodSeal(t)
	sealed(writeSeal(t, seal)).
		MustSucceed(t).
		StdoutContains(t, passWithAnomaly(noPostureReason(seal))).
		StdoutLacks(t, "Status: PASS\n")
}

// A witness run keeps its own status line: PASS-WITH-ANOMALY would imply that the MACs were
// checked. The warning is still on stdout.
func TestPosture_NoSourceWitnessKeepsItsStatusAndWarns(t *testing.T) {
	seal := goodSeal(t)
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--seal", writeSeal(t, seal), "--root-key", samplePubKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS-STRUCTURALLY, key-bound verification skipped\n").
		StdoutContains(t, "Anomaly: "+noPostureReason(seal)+"\n")
}

// §10.12: these are exit 3, and each is decided before any file is read. The file paths here
// do not exist, so reading one first would be exit 2.
func TestPosture_FlagErrorsAreExit3BeforeAnyRead(t *testing.T) {
	const gone, goneSeal = "no-such-file.ndjson", "no-such-seal.json"
	withKeyFlags := []string{"--root-key", samplePubKey, "--master-key-hex", sampleKey}
	for name, args := range map[string][]string{
		"--strict with no posture source": append([]string{"--seal", goneSeal, "--strict"}, withKeyFlags...),
		"--strict with no master key":     {"--seal", goneSeal, "--root-key", samplePubKey, "--strict", "--posture", "ed25519"},
		"--strict with no seal":           {"--master-key-hex", sampleKey, "--strict", "--posture", "ed25519"},
		"--posture with no seal":          {"--master-key-hex", sampleKey, "--posture", "ed25519"},
		"--posture with an empty name":    append([]string{"--seal", goneSeal, "--posture", "ed25519,"}, withKeyFlags...),
		"--posture that names nothing":    append([]string{"--seal", goneSeal, "--posture", " "}, withKeyFlags...),
	} {
		t.Run(name, func(t *testing.T) {
			clitest.Run(Main, append([]string{"verify", "--audit-file", gone}, args...)...).MustExit(t, 3)
		})
	}
}

// §7 step 11 case (d). A source that names the seal's one algorithm gives a clean result in both modes.
func TestPosture_SingleAlgorithmIsClean(t *testing.T) {
	seal := writeSeal(t, goodSeal(t))
	for _, mode := range [][]string{{"--posture", "ed25519"}, {"--posture", "ed25519", "--strict"}} {
		sealed(seal, mode...).
			MustSucceed(t).
			StdoutContains(t, "Status: PASS\n").
			StdoutLacks(t, "algorithm posture").
			StdoutLacks(t, "PASS-WITH-ANOMALY")
	}
}

// §7 step 11 case (b). Someone removes the second signature, or the whole list. The surviving
// signature is valid, so only a posture from outside the seal shows the loss.
func TestPosture_StrippedSignature(t *testing.T) {
	bare := goodSeal(t)
	oneEntry := goodSeal(t)
	oneEntry["signatures"] = []any{map[string]any{"algorithm": "ed25519", "signature": oneEntry["signature_b64"]}}
	for name, seal := range map[string]map[string]any{"the list removed": bare, "one entry left": oneEntry} {
		t.Run(name, func(t *testing.T) {
			file := writeSeal(t, seal)
			sealed(file, "--posture", dualPosture).
				MustSucceed(t).
				StdoutContains(t, passWithAnomaly(partialReason))
			sealed(file, "--posture", dualPosture, "--strict").
				MustExit(t, 1).
				StdoutContains(t, "Status: FAIL\nStep: 11\nReason: "+partialReason+"\n")
		})
	}
}

// §7 step 11 case (c). The seal names an algorithm the posture does not. Vector N018 pins the strict shape.
func TestPosture_AlgorithmNotOnTheList(t *testing.T) {
	extra := goodSeal(t)
	extra["signatures"] = []any{
		map[string]any{"algorithm": "ed25519", "signature": extra["signature_b64"]},
		map[string]any{"algorithm": "slh_dsa_shake_128s", "signature": "AAAA"},
	}
	file := writeSeal(t, extra)
	reason := notOnListReason(extra)
	sealed(file, "--posture", dualPosture).
		MustSucceed(t).
		StdoutContains(t, passWithAnomaly(reason)).
		StdoutContains(t, "Anomaly: unknown_algorithms: [slh_dsa_shake_128s]\n")
	sealed(file, "--posture", dualPosture, "--strict").
		MustExit(t, 1).
		StdoutContains(t, "Status: FAIL\nStep: 11\nReason: "+reason+"\nAnomaly: unknown_algorithms: [slh_dsa_shake_128s]\n")

	plain := goodSeal(t)
	sealed(writeSeal(t, plain), "--posture", "dilithium3", "--strict").
		MustExit(t, 1).
		StdoutContains(t, "Anomaly: unknown_algorithms: [ed25519]\n")
}

// §7 step 11: a field in the seal is not a source. No sign_payload form binds it.
func TestPosture_ASealFieldIsNotASource(t *testing.T) {
	for _, field := range []string{"declared_algorithms", "declared_algorithm_posture"} {
		t.Run(field, func(t *testing.T) {
			seal := goodSeal(t)
			seal[field] = []any{"ed25519"}
			file := writeSeal(t, seal)
			sealed(file).
				MustSucceed(t).
				StdoutContains(t, passWithAnomaly(noPostureReason(seal)))
			sealed(file, "--posture", dualPosture).
				MustSucceed(t).
				StdoutContains(t, passWithAnomaly(partialReason))
			sealed(file, "--posture", dualPosture, "--strict").MustExit(t, 1)
		})
	}
}

// Both signatures are named and the posture agrees, and this verifier cannot check the second
// one. That is "cannot judge" in both modes, and never a PASS.
func TestPosture_ADualSealStaysRefused(t *testing.T) {
	seal := goodSeal(t)
	seal["signatures"] = []any{
		map[string]any{"algorithm": "ed25519", "signature": seal["signature_b64"]},
		map[string]any{"algorithm": "dilithium3", "signature": "AAAA"},
	}
	file := writeSeal(t, seal)
	sealed(file, "--posture", dualPosture).MustExit(t, 3).StdoutContains(t, signaturesListRefused)
	sealed(file, "--posture", dualPosture, "--strict").MustExit(t, 3).StdoutContains(t, signaturesListRefused)
	sealed(file).MustExit(t, 3).StdoutContains(t, signaturesListRefused)
}

// The warning comes after the integrity checks. A forged seal stays exit 1 and never reaches it.
func TestPosture_AForgedSealStaysAFinding(t *testing.T) {
	lines := sampleLines(t)
	forged := sealFor(t, lines, load(t, lines).Header.TenantID, seededKey("the writer's own"))
	file := writeSeal(t, forged)
	sealed(file).MustExit(t, 1).StdoutContains(t, signatureFail)
	sealed(file, "--posture", dualPosture).MustExit(t, 1).StdoutContains(t, signatureFail)
}

// The other anomaly lines still print beside the warning.
func TestPosture_WarningKeepsTheOtherAnomalies(t *testing.T) {
	sealed(sampleSeal).MustSucceed(t).StdoutContains(t, "Anomaly: seal is dev_mode")
}
