package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/testkit/clitest"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// makeFixtureFiles is the shared setup: a clean fixture written to temp
// files, ready for verify/walk/diff tests to consume.
func makeFixtureFiles(t *testing.T) (ledgerPath, pubPath, ikmPath string) {
	t.Helper()
	dir := t.TempDir()
	ledgerPath = filepath.Join(dir, "day.ledger")
	pubPath = filepath.Join(dir, "seal.pub.pem")
	ikmPath = filepath.Join(dir, "ikm.bin")

	clitest.Run(Main,
		"gen-fixture",
		"--out-ledger", ledgerPath,
		"--out-pub", pubPath,
		"--out-ikm", ikmPath,
		"--entries", "5",
		"--seal-date", "2026-06-15",
	).MustSucceed(t)

	return ledgerPath, pubPath, ikmPath
}

// ----- top-level dispatch -------------------------------------------------

func TestMain_NoArgsPrintsUsage(t *testing.T) {
	clitest.Run(Main).
		MustExit(t, 2).
		StderrContains(t, "verifier —")
}

func TestMain_HelpPrintsToStdout(t *testing.T) {
	clitest.Run(Main, "--help").
		MustSucceed(t).
		StdoutContains(t, "verifier —")
}

func TestMain_UnknownCommandIsUsageError(t *testing.T) {
	clitest.Run(Main, "frobnicate").
		MustExit(t, 2).
		StderrContains(t, `unknown command "frobnicate"`)
}

// ----- gen-fixture --------------------------------------------------------

func TestGenFixture_ProducesLoadableLedger(t *testing.T) {
	ledgerPath, pubPath, ikmPath := makeFixtureFiles(t)

	led, err := verify.LoadLedger(ledgerPath)
	if err != nil {
		t.Fatalf("LoadLedger: %v", err)
	}
	if got := len(led.Entries); got != 5 {
		t.Errorf("entries = %d, want 5", got)
	}
	if got := led.Seal.SealDate; got != "2026-06-15" {
		t.Errorf("seal_date = %q", got)
	}
	if _, err := os.Stat(pubPath); err != nil {
		t.Errorf("pub file missing: %v", err)
	}
	ikm, err := os.ReadFile(ikmPath)
	if err != nil {
		t.Errorf("ikm file missing: %v", err)
	} else if len(ikm) != 32 {
		t.Errorf("ikm length = %d, want 32", len(ikm))
	}
}

func TestGenFixture_RequiresFlags(t *testing.T) {
	clitest.Run(Main, "gen-fixture").
		MustExit(t, 2).
		StderrContains(t, "missing required flags")
}

func TestGenFixture_RejectsZeroEntries(t *testing.T) {
	dir := t.TempDir()
	clitest.Run(Main, "gen-fixture",
		"--out-ledger", filepath.Join(dir, "x.ledger"),
		"--out-pub", filepath.Join(dir, "x.pub.pem"),
		"--entries", "0",
	).MustExit(t, 2).StderrContains(t, "entries")
}

// ----- verify -------------------------------------------------------------

func TestVerify_StructuralOnlyOnCleanLedger(t *testing.T) {
	ledgerPath, pubPath, _ := makeFixtureFiles(t)

	r := clitest.Run(Main,
		"verify",
		"--ledger", ledgerPath,
		"--root-key", pubPath,
	).MustSucceed(t)

	r.StdoutContains(t, "verifier report")
	r.StdoutContains(t, "[PASS] chain-linkage")
	r.StdoutContains(t, "[PASS] merkle-root")
	r.StdoutContains(t, "[PASS] seal-signature")
	r.StdoutContains(t, "[FAIL] per-event-mac")
	r.StdoutContains(t, "structural-only")
	r.StdoutContains(t, "PASS (structural)")
}

func TestVerify_FullPathPassesWithMasterKey(t *testing.T) {
	ledgerPath, pubPath, ikmPath := makeFixtureFiles(t)

	r := clitest.Run(Main,
		"verify",
		"--ledger", ledgerPath,
		"--root-key", pubPath,
		"--master-key", ikmPath,
	).MustSucceed(t)

	r.StdoutContains(t, "[PASS] per-event-mac")
	r.StdoutContains(t, "PASS (full, key-bound)")
}

func TestVerify_FailsLoudOnTamperedLedger(t *testing.T) {
	ledgerPath, pubPath, _ := makeFixtureFiles(t)

	// Tamper one entry's payload in-place.
	led, err := verify.LoadLedger(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}
	led.Entries[2].EventPayloadJCS = base64.StdEncoding.EncodeToString([]byte("tampered"))
	raw, err := verify.MarshalNDJSON(led.Entries, led.Seal)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ledgerPath, raw, 0o600); err != nil {
		t.Fatal(err)
	}

	r := clitest.Run(Main,
		"verify",
		"--ledger", ledgerPath,
		"--root-key", pubPath,
	).MustExit(t, 1)
	r.StdoutContains(t, "[FAIL] chain-linkage")
}

func TestVerify_FailsOnWrongRootKey(t *testing.T) {
	ledgerPath, _, _ := makeFixtureFiles(t)
	dir := t.TempDir()
	otherPubPath := filepath.Join(dir, "other.pub.pem")
	otherLedgerPath := filepath.Join(dir, "other.ledger")
	clitest.Run(Main, "gen-fixture",
		"--out-ledger", otherLedgerPath,
		"--out-pub", otherPubPath,
	).MustSucceed(t)

	clitest.Run(Main,
		"verify",
		"--ledger", ledgerPath,
		"--root-key", otherPubPath,
	).MustExit(t, 1).StdoutContains(t, "[FAIL] seal-signature")
}

func TestVerify_RequiresFlags(t *testing.T) {
	clitest.Run(Main, "verify").
		MustExit(t, 2).
		StderrContains(t, "missing required flags").
		StderrContains(t, "--ledger").
		StderrContains(t, "--root-key")
}

func TestVerify_RejectsMissingLedgerFile(t *testing.T) {
	_, pubPath, _ := makeFixtureFiles(t)
	clitest.Run(Main, "verify",
		"--ledger", filepath.Join(t.TempDir(), "no.ledger"),
		"--root-key", pubPath,
	).MustExit(t, 1).StderrContains(t, "open ledger")
}

func TestVerify_HelpExitsZero(t *testing.T) {
	clitest.Run(Main, "verify", "--help").
		MustSucceed(t).
		StdoutContains(t, "Usage of verifier verify").
		StdoutContains(t, "ledger")
}

func TestVerify_UnknownProfileIsConfigError(t *testing.T) {
	ledgerPath, pubPath, _ := makeFixtureFiles(t)
	clitest.Run(Main,
		"verify",
		"--ledger", ledgerPath,
		"--root-key", pubPath,
		"--profile", "nonesuch",
	).MustExit(t, 1).StderrContains(t, "unknown profile")
}

// The default (ffiec) profile on a chain with no supervisory family
// reproduces the integrity core with no extra profile framing.
func TestVerify_DefaultProfileEmitsNoProfileLine(t *testing.T) {
	ledgerPath, pubPath, _ := makeFixtureFiles(t)
	r := clitest.Run(Main,
		"verify",
		"--ledger", ledgerPath,
		"--root-key", pubPath,
	).MustSucceed(t)
	r.StdoutContains(t, "PASS (structural)")
	if strings.Contains(r.Stdout, "profile:") {
		t.Errorf("default profile must not emit a profile line:\n%s", r.Stdout)
	}
}

// A non-default profile surfaces its active-profile line without changing
// the integrity verdict — presentation-only.
func TestVerify_NonDefaultProfileShowsFramingNotVerdict(t *testing.T) {
	ledgerPath, pubPath, _ := makeFixtureFiles(t)
	r := clitest.Run(Main,
		"verify",
		"--ledger", ledgerPath,
		"--root-key", pubPath,
		"--profile", "tx-dob",
	).MustSucceed(t)
	r.StdoutContains(t, "profile:       tx-dob")
	r.StdoutContains(t, "PASS (structural)") // verdict unchanged by the profile
}

// ----- walk ---------------------------------------------------------------

func TestWalk_PrintsEveryEntryAndSeal(t *testing.T) {
	ledgerPath, _, _ := makeFixtureFiles(t)
	r := clitest.Run(Main,
		"walk",
		"--ledger", ledgerPath,
	).MustSucceed(t)

	for i := 0; i < 5; i++ {
		// Each entry's id should appear on its own line.
		r.StdoutContains(t, "e_000")
	}
	r.StdoutContains(t, "seal:")
	r.StdoutContains(t, "leaves=5")
	r.StdoutContains(t, "signer=sha256:")
}

func TestWalk_RequiresLedger(t *testing.T) {
	clitest.Run(Main, "walk").
		MustExit(t, 2).
		StderrContains(t, "--ledger")
}

// ----- diff --------------------------------------------------------------

func TestDiff_IdenticalLedgersReportSameness(t *testing.T) {
	ledgerPath, _, _ := makeFixtureFiles(t)
	dir := t.TempDir()
	dupPath := filepath.Join(dir, "dup.ledger")
	raw, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dupPath, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	clitest.Run(Main,
		"diff",
		"--before", ledgerPath,
		"--after", dupPath,
	).MustSucceed(t).StdoutContains(t, "ledgers identical")
}

func TestDiff_ReportsTamper(t *testing.T) {
	ledgerPath, _, _ := makeFixtureFiles(t)
	dir := t.TempDir()
	tamperedPath := filepath.Join(dir, "tampered.ledger")

	led, err := verify.LoadLedger(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}
	// Flip a single entry_hash to surface a difference.
	led.Entries[1].EntryHash = base64.StdEncoding.EncodeToString(make([]byte, 32))
	raw, err := verify.MarshalNDJSON(led.Entries, led.Seal)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tamperedPath, raw, 0o600); err != nil {
		t.Fatal(err)
	}

	clitest.Run(Main,
		"diff",
		"--before", ledgerPath,
		"--after", tamperedPath,
	).MustExit(t, 1).StdoutContains(t, "entry 1 entry_hash differs")
}

func TestDiff_RequiresFlags(t *testing.T) {
	clitest.Run(Main, "diff").
		MustExit(t, 2).
		StderrContains(t, "missing required flags")
}

// ----- determinism -------------------------------------------------------

func TestVerifyReportIsDeterministicAcrossRuns(t *testing.T) {
	ledgerPath, pubPath, _ := makeFixtureFiles(t)

	r1 := clitest.Run(Main, "verify",
		"--ledger", ledgerPath,
		"--root-key", pubPath,
	).MustSucceed(t)
	r2 := clitest.Run(Main, "verify",
		"--ledger", ledgerPath,
		"--root-key", pubPath,
	).MustSucceed(t)
	if r1.Stdout != r2.Stdout {
		t.Errorf("verify report not deterministic across runs:\n--- run 1 ---\n%s\n--- run 2 ---\n%s", r1.Stdout, r2.Stdout)
	}
}

// Confirms the fixture and the CLI's verifier produce a parseable Seal —
// catches accidental schema drift at the cmd boundary.
func TestVerify_SealIsValidJSON(t *testing.T) {
	ledgerPath, _, _ := makeFixtureFiles(t)
	raw, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}
	last := raw
	for {
		i := strings.LastIndex(string(last[:len(last)-1]), "\n")
		if i < 0 {
			break
		}
		last = last[i+1:]
		if strings.Contains(string(last), `"type":"seal"`) {
			break
		}
		last = raw[:i]
	}
	var seal verify.SealRecord
	if err := json.Unmarshal(last[:len(last)-1], &seal); err != nil {
		t.Fatalf("could not parse last line as seal: %v\n%s", err, last)
	}
	if seal.Type != "seal" {
		t.Errorf("last line is not a seal record: %+v", seal)
	}
}
