package main

import (
	"testing"

	"github.com/mmpworks/ffiec/core/jcs"
	"github.com/mmpworks/ffiec/testkit/clitest"
)

// Richard's fork pair. An editor removes chain_id from one branch's header to make 2 claims
// differ. That file now stops at exit 2 in either flag order, sealed or not.
func TestFork_RemovingAClaimFromOneHeaderIsExit2(t *testing.T) {
	branchOne := buildRun(t, "RUN-A", manyAt(2, dayOne))
	branchTwo := withoutChainID(t, buildRun(t, "RUN-A", manyAt(3, dayOne)))
	one, two := writeLines(t, branchOne), writeLines(t, branchTwo)
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": branchOne}, everyEntry))
	const reason = "Status: FAIL\nStep: pre-flight\nReason: header missing required field chain_id\n"
	for name, files := range map[string][]string{"as given": {one, two}, "reversed": {two, one}} {
		t.Run(name, func(t *testing.T) {
			runSet(files, []string{seal}).MustExit(t, 2).StdoutContains(t, reason)
			clitest.Run(Main, "verify", "--audit-file", files[0], "--audit-file", files[1]).MustExit(t, 2).StdoutContains(t, reason)
		})
	}
}

// The other way to break a claim: name another run in the header. Step 4 refuses it, because each
// entry's run_id is under the MAC. So after steps 1-9 the header's claim is the entries' claim,
// and checkForks can read it from the header.
func TestFork_ANameThatIsNotTheEntriesRunFailsAtStep4(t *testing.T) {
	relabelled := buildRun(t, "RUN-A", manyAt(2, dayOne))
	header := decodeLine(t, relabelled[0])
	header["chain_id"] = "RUN-X"
	line, err := jcs.Canonicalize(header)
	if err != nil {
		t.Fatal(err)
	}
	relabelled[0] = line
	clitest.Run(Main, "verify", "--audit-file", writeLines(t, buildRun(t, "RUN-A", manyAt(3, dayOne))), "--audit-file", writeLines(t, relabelled)).
		MustExit(t, 1).
		StdoutContains(t, "Step: 4\nReason: cross-chain lift detected at seq 1 (event.run_id mismatch)\n")
}

// One file holding 2 runs fails at step 4, on the first entry of the other run.
func TestFork_OneFileWithTwoRunIDsFailsAtStep4(t *testing.T) {
	first := buildRun(t, "RUN-A", manyAt(2, dayOne))
	second := buildRun(t, "RUN-B", manyAt(1, dayOne))
	file := writeLines(t, append(append([][]byte{}, first...), second[1]))
	clitest.Run(Main, "verify", "--audit-file", file).
		MustExit(t, 1).
		StdoutContains(t, "Step: 4\nReason: cross-chain lift detected at seq 3 (event.run_id mismatch)\n")
}

func withoutChainID(t *testing.T, lines [][]byte) [][]byte {
	t.Helper()
	header := decodeLine(t, lines[0])
	delete(header, "chain_id")
	line, err := jcs.Canonicalize(header)
	if err != nil {
		t.Fatal(err)
	}
	return append([][]byte{line}, lines[1:]...)
}

// Richard's note. An unstamped record has no leaf, so a clean seal still covers the file. Under
// --strict the run stays exit 0, and only the count tells the reader. The README says so.
func TestStrict_AnAddedUnstampedRecordStaysExit0(t *testing.T) {
	run := buildRun(t, "RUN-A", manyAt(2, dayOne))
	for i := 1; i < len(run); i++ {
		whole := decodeLine(t, run[i])
		whole["audit_chain"].(map[string]any)["kms_handle_uri"] = "aws-kms:arn:test"
		line, err := jcs.Canonicalize(whole)
		if err != nil {
			t.Fatal(err)
		}
		run[i] = line
	}
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": run}, everyEntry))
	added, err := jcs.Canonicalize(map[string]any{"run_id": "RUN-A", "kind": "added", "audit_chain": map[string]any{"format_version": "v1"}})
	if err != nil {
		t.Fatal(err)
	}
	runSet([]string{writeLines(t, append(run, added))}, []string{seal}, "--strict").
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutContains(t, "Anomaly: unstamped records: 1\n")
}

// §7 and §4.2.2: the count, and each entry's (run_id, seq, received_at, captured_at).
func TestLateBinding_IsReported(t *testing.T) {
	late := dayOne + "T23:59:59Z"
	run := buildRun(t, "RUN-A", []stamp{at(dayTwo), {receivedAt: late, lateBinding: true}})
	seal := writeSeal(t, sealOver(t, dayTwo, map[string][][]byte{"RUN-A": run}, everyEntry))
	runSet([]string{writeLines(t, run)}, []string{seal}).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n").
		StdoutContains(t, "Anomaly: late-binding entries: 1\n").
		StdoutContains(t, "Anomaly: late-binding entry at seal_date "+dayTwo+": (\"RUN-A\", 2, "+late+", captured_at absent)\n")
}

func TestLateBinding_NoLineWithNoneOfThem(t *testing.T) {
	run := buildRun(t, "RUN-A", manyAt(1, dayOne))
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": run}, everyEntry))
	runSet([]string{writeLines(t, run)}, []string{seal}).MustSucceed(t).StdoutLacks(t, "late-binding")
}

// Richard's limit. A run crosses midnight. Day 2's entries are cut and day 2's seal is left out.
// Every supplied entry is covered, so nothing here can fail. The output names the days it covered,
// and says the rest cannot be vouched for, so the reader sees "day 1 only".
func TestCoverage_ACutLaterDayIsNamedAsNotCovered(t *testing.T) {
	whole := buildRun(t, "RUN-A", append(manyAt(2, dayOne), manyAt(2, dayTwo)...))
	dayOneOnly := whole[:3]
	first := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": whole}, func(_ string, seq int) bool { return seq <= 2 }))
	file := writeLines(t, dayOneOnly)
	const line = "Anomaly: seal dates covered: " + dayOne + ": a later day's entries, cut from the end of a run along with that day's seal, cannot be detected\n"
	runSet([]string{file}, []string{first}).MustSucceed(t).StdoutContains(t, line)
	clitest.Run(Main, "verify", "--audit-file", file, "--seal", first, "--root-key", samplePubKey, "--posture", "ed25519").
		MustSucceed(t).
		StdoutContains(t, "Status: PASS-STRUCTURALLY").
		StdoutContains(t, line)
}

func TestCoverage_NamesEveryDateInOrder(t *testing.T) {
	run := buildRun(t, "RUN-A", append(manyAt(1, dayOne), manyAt(1, dayTwo)...))
	runs := map[string][][]byte{"RUN-A": run}
	first := writeSeal(t, sealOver(t, dayOne, runs, func(_ string, seq int) bool { return seq == 1 }))
	second := writeSeal(t, sealOver(t, dayTwo, runs, func(_ string, seq int) bool { return seq == 2 }))
	runSet([]string{writeLines(t, run)}, []string{second, first}).
		MustSucceed(t).
		StdoutContains(t, "Anomaly: seal dates covered: "+dayOne+", "+dayTwo+":")
}
