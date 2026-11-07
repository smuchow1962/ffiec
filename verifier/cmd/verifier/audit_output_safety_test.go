package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/core/jcs"
	"github.com/mmpworks/ffiec/testkit/clitest"
)

// A stamp field is outside the MAC, and a signatures list is unsigned, so anyone who can edit the
// file can put a line break or a terminal sequence in them. §7 fixes what stdout means, so no such
// value may forge a line or reach a terminal.
const (
	forgedLine = "x\nStatus: PASS\nAnomaly: none"
	escape     = "x\x1b[2Jy"
)

func assertOneStatusLineNoEscape(t *testing.T, stdout string) {
	t.Helper()
	statuses := 0
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "Status:") {
			statuses++
		}
	}
	if statuses != 1 {
		t.Errorf("%d status lines:\n%s", statuses, stdout)
	}
	if strings.ContainsRune(stdout, 0x1b) {
		t.Errorf("an ESC byte reached stdout: %q", stdout)
	}
}

// Richard's probe. captured_at prints in the late-binding line.
func TestOutput_CapturedAtCannotForgeALine(t *testing.T) {
	for name, value := range map[string]string{"a line break": forgedLine, "an escape": escape} {
		t.Run(name, func(t *testing.T) {
			run := buildRun(t, "RUN-A", []stamp{at(dayTwo), {receivedAt: dayOne + "T23:59:59Z", lateBinding: true}})
			seal := writeSeal(t, sealOver(t, dayTwo, map[string][][]byte{"RUN-A": run}, everyEntry))
			whole := decodeLine(t, run[2])
			whole["audit_chain"].(map[string]any)["captured_at"] = value
			line, err := jcs.Canonicalize(whole)
			if err != nil {
				t.Fatal(err)
			}
			run[2] = line
			result := runSet([]string{writeLines(t, run)}, []string{seal}).MustSucceed(t)
			assertOneStatusLineNoEscape(t, result.Stdout)
		})
	}
}

// The same class: an algorithm name in the unsigned list prints in unknown_algorithms.
func TestOutput_AnAlgorithmNameCannotForgeALine(t *testing.T) {
	for name, value := range map[string]string{"a line break": forgedLine, "an escape": escape} {
		t.Run(name, func(t *testing.T) {
			seal := goodSeal(t)
			seal["signatures"] = []any{map[string]any{"algorithm": value, "signature": "AAAA"}}
			result := sealed(writeSeal(t, seal), "--posture", "ed25519").MustSucceed(t)
			assertOneStatusLineNoEscape(t, result.Stdout)
		})
	}
}

// The class sweep: run_id reaches a reason. A witness run has no MAC check, so an editor can set
// it and chain_id together, and step 4 still holds.
func TestOutput_ARunIDInAReasonCannotForgeALine(t *testing.T) {
	run := buildRun(t, forgedLine, manyAt(1, dayTwo))
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{forgedLine: run}, everyEntry))
	result := clitest.Run(Main, "verify", "--audit-file", writeLines(t, run), "--seal", seal, "--root-key", samplePubKey).MustExit(t, 1)
	assertOneStatusLineNoEscape(t, result.Stdout)
}

// The escape keeps the value readable, and ordinary text prints as it stands.
func TestOutput_SafeLine(t *testing.T) {
	for in, want := range map[string]string{
		"a\nb":                           `a\u000Ab`,
		"a\x1b[2J":                       `a\u001B[2J`,
		"a" + string(rune(0x202E)) + "b": fmt.Sprintf("a%cu202Eb", 92),
		"a" + string(rune(0x2028)) + "b": fmt.Sprintf("a%cu2028b", 92),
		"merkle — §7":                    "merkle — §7",
	} {
		if got := cliutil.SafeLine(in); got != want {
			t.Errorf("SafeLine(%q) = %q, want %q", in, got, want)
		}
	}
}
