package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmpworks/ffiec/testkit/clitest"
)

// The sample holds event bytes a .NET sink wrote, from MMP.GasPumpSimulator. Its
// key is a published demo key. A clean PASS here proves the two languages agree
// on the canonical form, the HKDF inputs, and the MAC.
const (
	sampleFile = "testdata/simulator-custody.golden"
	sampleKey  = "testdata/simulator-custody-key.hex"
)

func sampleLines(t *testing.T) [][]byte {
	t.Helper()
	raw, err := os.ReadFile(sampleFile)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.Split(bytes.TrimSuffix(raw, []byte("\r\n")), []byte("\r\n"))
}

func writeLines(t *testing.T, lines [][]byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "attacked.ndjson")
	if err := os.WriteFile(p, append(bytes.Join(lines, []byte("\r\n")), '\r', '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// edit returns a copy of lines with one replacement in line n. Line n holds seq n; line 0 is the header.
func edit(t *testing.T, lines [][]byte, n int, old, replacement string) [][]byte {
	t.Helper()
	if !bytes.Contains(lines[n], []byte(old)) {
		t.Fatalf("line %d lacks %q", n, old)
	}
	out := append([][]byte(nil), lines...)
	out[n] = bytes.Replace(lines[n], []byte(old), []byte(replacement), 1)
	return out
}

func TestVerifyAuditFile_CleanSampleCrossLanguage(t *testing.T) {
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--master-key-hex", sampleKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutContains(t, "Anomaly: no daily seal supplied").
		StdoutContains(t, "Anomaly: 36 of 36 entries name a development key").
		StdoutContains(t, "Anomaly: 36 of 36 entries carry no chain_kind").
		StdoutContains(t, "Anomaly: CRLF line endings")
}

func TestVerifyAuditFile_WitnessModeNeedsNoKey(t *testing.T) {
	clitest.Run(Main, "verify", "--audit-file", sampleFile).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS-STRUCTURALLY, key-bound verification skipped\n")
}

func TestVerifyAuditFile_EachAttackFailsAtItsStep(t *testing.T) {
	lines := sampleLines(t)
	swapped := append([][]byte(nil), lines...)
	swapped[4], swapped[5] = lines[5], lines[4]
	deleted := append(append([][]byte(nil), lines[:4]...), lines[5:]...)
	inserted := append(append(append([][]byte(nil), lines[:7]...), lines[6]), lines[7:]...)

	cases := []struct {
		name   string
		lines  [][]byte
		output string
	}{
		{"edit an event field", edit(t, lines, 6, "$30.00", "$3.00"), "Step: 9\nReason: payload_hash MAC mismatch at seq 6\n"},
		{"edit a hash in the stamp", edit(t, lines, 6, `"payload_hash":"A`, `"payload_hash":"B`), "Step: 9\nReason: payload_hash MAC mismatch at seq 6\n"},
		{"add whitespace", edit(t, lines, 6, `"category":"ledger"`, `"category": "ledger"`), "Step: 9\nReason: payload_hash MAC mismatch at seq 6\n"},
		{"add a duplicate key", edit(t, lines, 6, `"category":"ledger"`, `"category":"refund","category":"ledger"`), "Step: 9\nReason: payload_hash MAC mismatch at seq 6\n"},
		{"delete an entry", deleted, "Step: 6\nReason: chain link broken at seq 4\n"},
		{"swap two entries", swapped, "Step: 6\nReason: chain link broken at seq 4\n"},
		{"replay an entry", inserted, "Step: 6\nReason: chain link broken at seq 7\n"},
		{"lift from another tenant", edit(t, lines, 3, `"tenant_id":"PUMP-01"`, `"tenant_id":"PUMP-02"`), "Step: 4\nReason: cross-chain lift detected at seq 3 (event.tenant_id mismatch)\n"},
		{"lift from another run", edit(t, lines, 3, `"run_id":"DSP-`, `"run_id":"OTHER-DSP-`), "Step: 4\nReason: cross-chain lift detected at seq 3 (event.run_id mismatch)\n"},
		{"renumber an entry", edit(t, lines, 1, `"seq":1`, `"seq":9`), "Step: 6\nReason: seq out of order at seq 1\n"},
		{"header format", edit(t, lines, 0, `"format_version":"v1"`, `"format_version":"v2"`), "Step: 1\nReason: format_version \"v2\" not supported by this verifier (running v1)\n"},
		{"header tenant", edit(t, lines, 0, `"tenant_id":"PUMP-01"`, `"tenant_id":"PUMP-09"`), "Step: 2\nReason: header HKDF inputs do not match running v1 inputs\n"},
		{"header genesis", edit(t, lines, 0, `"genesis_hash":"AAAA`, `"genesis_hash":"AQEB`), "Step: 3\nReason: header genesis_hash does not match v1 constant\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clitest.Run(Main, "verify", "--audit-file", writeLines(t, c.lines), "--master-key-hex", sampleKey).
				MustExit(t, 1).
				StdoutContains(t, "Status: FAIL\n"+c.output).
				StdoutLacks(t, "Anomaly:")
		})
	}
}

func TestVerifyAuditFile_WitnessCatchesAStructuralAttack(t *testing.T) {
	lines := sampleLines(t)
	deleted := append(append([][]byte(nil), lines[:4]...), lines[5:]...)
	clitest.Run(Main, "verify", "--audit-file", writeLines(t, deleted)).
		MustExit(t, 1).
		StdoutContains(t, "Step: 6\nReason: chain link broken at seq 4\n")
}

func TestVerifyAuditFile_WrongKeyFailsAtTheFingerprint(t *testing.T) {
	wrong := filepath.Join(t.TempDir(), "wrong.hex")
	if err := os.WriteFile(wrong, bytes.Repeat([]byte("ab"), 32), 0o600); err != nil {
		t.Fatal(err)
	}
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--master-key-hex", wrong).
		MustExit(t, 1).
		StdoutContains(t, "Step: 8\nReason: key_fingerprint mismatch at seq 1")
}

// A cut at the end of the file leaves a good chain. Only a seal can catch it, so
// the pass has to say the check did not run. This pins the anomaly line.
func TestVerifyAuditFile_TailCutPassesAndSaysWhatItCouldNotCheck(t *testing.T) {
	lines := sampleLines(t)
	clitest.Run(Main, "verify", "--audit-file", writeLines(t, lines[:len(lines)-5]), "--master-key-hex", sampleKey).
		MustSucceed(t).
		StdoutContains(t, "Anomaly: no daily seal supplied: steps 10 and 11 did not run, so a cut at the end of the file cannot be detected")
}

func TestVerifyAuditFile_InputErrors(t *testing.T) {
	dir := t.TempDir()
	empty := filepath.Join(dir, "empty.ndjson")
	if err := os.WriteFile(empty, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(sampleFile)
	if err != nil {
		t.Fatal(err)
	}
	cut := filepath.Join(dir, "cut.ndjson")
	if err := os.WriteFile(cut, raw[:len(raw)-1], 0o600); err != nil {
		t.Fatal(err)
	}
	badKey := filepath.Join(dir, "bad.hex")
	if err := os.WriteFile(badKey, []byte("not hex"), 0o600); err != nil {
		t.Fatal(err)
	}

	clitest.Run(Main, "verify", "--audit-file", empty).
		MustExit(t, 2).StdoutContains(t, "Reason: empty file: header missing\n")
	clitest.Run(Main, "verify", "--audit-file", cut, "--master-key-hex", sampleKey).
		MustExit(t, 2).StdoutContains(t, "Step: pre-flight\nReason: audit file ends mid-line")
	clitest.Run(Main, "verify", "--audit-file", filepath.Join(dir, "absent.ndjson")).
		MustExit(t, 2).StdoutContains(t, "Step: pre-flight\nReason: cannot read file")
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--master-key-hex", badKey).
		MustExit(t, 3).StderrContains(t, "--master-key-hex")
	clitest.Run(Main, "verify", "--audit-file", sampleFile, "--master-key", sampleKey, "--master-key-hex", sampleKey).
		MustExit(t, 2)
}

// The ledger path keeps its contract when the new flag is absent.
func TestVerify_LedgerStillRequiresItsFlags(t *testing.T) {
	clitest.Run(Main, "verify").MustExit(t, 2)
}
