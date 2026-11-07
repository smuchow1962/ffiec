package main

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/core/constants"
	"github.com/mmpworks/ffiec/core/hkdf"
	"github.com/mmpworks/ffiec/core/jcs"
	"github.com/mmpworks/ffiec/core/signpayload"
	"github.com/mmpworks/ffiec/testkit/clitest"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

const (
	dayOne = "2026-11-04"
	dayTwo = "2026-11-05"
)

// stamp is what one test entry carries outside the MAC. An empty receivedAt leaves the field out.
type stamp struct {
	receivedAt  string
	lateBinding bool
}

func at(date string) stamp { return stamp{receivedAt: date + "T12:00:00Z"} }

func manyAt(n int, date string) []stamp {
	stamps := make([]stamp, n)
	for i := range stamps {
		stamps[i] = at(date)
	}
	return stamps
}

// buildRun writes a real keyed run: one header line, then one entry per stamp, each with a good
// MAC and link. It uses the sample's tenant and key, so the sample's flags verify it.
func buildRun(t *testing.T, run string, stamps []stamp) [][]byte {
	t.Helper()
	sample := sampleLines(t)
	header := decodeLine(t, sample[0])
	header["chain_id"] = run
	headerLine, err := jcs.Canonicalize(header)
	if err != nil {
		t.Fatal(err)
	}
	lines := [][]byte{headerLine}
	key := sampleSessionKey(t, load(t, sample).Header.TenantID)
	prev := make([]byte, sha256.Size)
	for i, s := range stamps {
		whole := decodeLine(t, sample[1])
		chain := whole["audit_chain"].(map[string]any)
		delete(whole, "audit_chain")
		whole["run_id"] = run
		whole["message"] = fmt.Sprintf("entry %d of %s", i+1, run)
		canonical, err := jcs.Canonicalize(whole)
		if err != nil {
			t.Fatal(err)
		}
		mac := hmac.New(sha256.New, key)
		mac.Write(prev)
		mac.Write(canonical)
		tag := mac.Sum(nil)
		chain["seq"] = i + 1
		chain["prev_hash"] = base64.StdEncoding.EncodeToString(prev)
		chain["payload_hash"] = base64.StdEncoding.EncodeToString(tag)
		if s.receivedAt != "" {
			chain["received_at"] = s.receivedAt
		}
		if s.lateBinding {
			chain["late_binding"] = true
		}
		whole["audit_chain"] = chain
		line, err := jcs.Canonicalize(whole)
		if err != nil {
			t.Fatal(err)
		}
		lines = append(lines, line)
		prev = tag
	}
	return lines
}

func sampleSessionKey(t *testing.T, tenant string) []byte {
	t.Helper()
	keyHex, err := os.ReadFile(sampleKey)
	if err != nil {
		t.Fatal(err)
	}
	ikm, err := hex.DecodeString(strings.TrimSpace(string(keyHex)))
	if err != nil {
		t.Fatal(err)
	}
	info := []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + tenant)
	return hkdf.Derive(ikm, []byte(constants.HKDFSalt), info, constants.HKDFOutputLength)
}

// sealOver signs a seal for date whose root covers the named entries of the given runs, in
// (run_id, seq) order. This is the test's own statement of step 10, written from the spec.
func sealOver(t *testing.T, date string, runs map[string][][]byte, covered func(run string, seq int) bool) map[string]any {
	t.Helper()
	names := make([]string, 0, len(runs))
	for name := range runs {
		names = append(names, name)
	}
	sort.Strings(names)
	leaves := [][]byte{}
	var digest, format, tenant string
	for _, name := range names {
		af := load(t, runs[name])
		digest, format, tenant = af.Header.HKDFInputsDigestHex, af.Header.FormatVersion, af.Header.TenantID
		for i := range af.Entries {
			if !covered(name, af.Entries[i].Seq) {
				continue
			}
			tag, err := hex.DecodeString(af.Entries[i].PayloadHashHex)
			if err != nil {
				t.Fatal(err)
			}
			leaves = append(leaves, verify.MerkleLeafHash(tag))
		}
	}
	root := hex.EncodeToString(verify.MerkleTreeHash(leaves))
	version, err := signpayload.ParseVersion("v1.0a")
	if err != nil {
		t.Fatal(err)
	}
	payload, err := signpayload.Build(signpayload.Seal{
		Version: version, Algorithm: "ed25519", FormatVersion: format, TenantID: tenant,
		SealDate: date, MerkleRootHex: root, HKDFInputsDigest: digest, Cadence: "daily",
	})
	if err != nil {
		t.Fatal(err)
	}
	var priv ed25519.PrivateKey = institutionKey()
	return map[string]any{
		"tenant_id": tenant, "seal_date": date, "format_version": format, "merkle_root_hex": root,
		"algorithm": "ed25519", "hkdf_inputs_digest_hex": digest, "cadence": "daily", "dev_mode": false,
		"sign_payload_version": "v1.0a", "sign_payload_hex": hex.EncodeToString(payload),
		"signature_b64": base64.StdEncoding.EncodeToString(ed25519.Sign(priv, payload)),
	}
}

func everyEntry(string, int) bool { return true }

// runSet runs the keyed walk over files and seals, with a single-algorithm posture.
func runSet(files, seals []string, extra ...string) clitest.Result {
	args := []string{"verify", "--root-key", samplePubKey, "--master-key-hex", sampleKey, "--posture", "ed25519"}
	for _, f := range files {
		args = append(args, "--audit-file", f)
	}
	for _, s := range seals {
		args = append(args, "--seal", s)
	}
	return clitest.Run(Main, append(args, extra...)...)
}

// §7 step 10: a day with 2 runs verifies when both files are supplied.
func TestFileSet_TwoRunsOneSealPass(t *testing.T) {
	a, b := buildRun(t, "RUN-A", manyAt(3, dayOne)), buildRun(t, "RUN-B", manyAt(2, dayOne))
	seal := sealOver(t, dayOne, map[string][][]byte{"RUN-A": a, "RUN-B": b}, everyEntry)
	runSet([]string{writeLines(t, a), writeLines(t, b)}, []string{writeSeal(t, seal)}).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n")
}

// The order of the flags changes nothing: the roots use (run_id, seq).
func TestFileSet_FileOrderDoesNotMatter(t *testing.T) {
	a, b := buildRun(t, "RUN-A", manyAt(3, dayOne)), buildRun(t, "RUN-B", manyAt(2, dayOne))
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": a, "RUN-B": b}, everyEntry))
	runSet([]string{writeLines(t, b), writeLines(t, a)}, []string{seal}).MustSucceed(t)
}

// The roots are the completeness check: a missing file gives the root mismatch.
func TestFileSet_MissingFileFailsAtTheRoot(t *testing.T) {
	a, b := buildRun(t, "RUN-A", manyAt(3, dayOne)), buildRun(t, "RUN-B", manyAt(2, dayOne))
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": a, "RUN-B": b}, everyEntry))
	runSet([]string{writeLines(t, a)}, []string{seal}).MustExit(t, 1).StdoutContains(t, rootMismatch)
}

// A run that crosses midnight verifies with both days' seals, and each seal covers its own day.
func TestFileSet_RunAcrossMidnightNeedsBothSeals(t *testing.T) {
	stamps := append(manyAt(2, dayOne), manyAt(2, dayTwo)...)
	run := buildRun(t, "RUN-A", stamps)
	runs := map[string][][]byte{"RUN-A": run}
	first := sealOver(t, dayOne, runs, func(_ string, seq int) bool { return seq <= 2 })
	second := sealOver(t, dayTwo, runs, func(_ string, seq int) bool { return seq > 2 })
	file := writeLines(t, run)
	runSet([]string{file}, []string{writeSeal(t, first), writeSeal(t, second)}).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS\n")
	// Codex's blocker: with one seal, the other day's entries are covered by no root.
	runSet([]string{file}, []string{writeSeal(t, first)}).
		MustExit(t, 1).
		StdoutContains(t, "Status: FAIL\nStep: 10\nReason: entry outside every supplied seal at seq 3: run_id RUN-A\n")
}

// Codex's extreme case: only entries outside D, beside a genuine empty-day seal, reads no PASS.
func TestFileSet_OnlyOutsideEntriesWithAnEmptySealFail(t *testing.T) {
	run := buildRun(t, "RUN-A", manyAt(2, dayTwo))
	empty := sealOver(t, dayOne, map[string][][]byte{"RUN-A": run}, func(string, int) bool { return false })
	runSet([]string{writeLines(t, run)}, []string{writeSeal(t, empty)}).
		MustExit(t, 1).
		StdoutContains(t, "entry outside every supplied seal at seq 1: run_id RUN-A")
}

// A true empty day: a header with no entries and the empty-tree root.
func TestFileSet_EmptyDayPasses(t *testing.T) {
	run := buildRun(t, "RUN-A", nil)
	empty := sealOver(t, dayOne, map[string][][]byte{"RUN-A": run}, everyEntry)
	runSet([]string{writeLines(t, run)}, []string{writeSeal(t, empty)}).MustSucceed(t)
}

// §4.2.2: a late-binding entry counts as the next day. Its received_at alone would say otherwise.
func TestFileSet_LateBindingMovesAnEntryOneDay(t *testing.T) {
	run := buildRun(t, "RUN-A", []stamp{at(dayTwo), {receivedAt: dayOne + "T23:59:59Z", lateBinding: true}})
	seal := sealOver(t, dayTwo, map[string][][]byte{"RUN-A": run}, everyEntry)
	runSet([]string{writeLines(t, run)}, []string{writeSeal(t, seal)}).MustSucceed(t)
}

// Editing an unsigned field moves an entry to another seal, where a root fails.
func TestFileSet_EditedReceivedAtFailsAtARoot(t *testing.T) {
	stamps := append(manyAt(2, dayOne), manyAt(2, dayTwo)...)
	run := buildRun(t, "RUN-A", stamps)
	runs := map[string][][]byte{"RUN-A": run}
	first := writeSeal(t, sealOver(t, dayOne, runs, func(_ string, seq int) bool { return seq <= 2 }))
	second := writeSeal(t, sealOver(t, dayTwo, runs, func(_ string, seq int) bool { return seq > 2 }))
	stamps[1] = at(dayTwo)
	edited := buildRun(t, "RUN-A", stamps)
	runSet([]string{writeLines(t, edited)}, []string{first, second}).MustExit(t, 1).StdoutContains(t, rootMismatch)
}

// With no received_at anywhere, all entries count as the one seal's day, and the root decides.
func TestFileSet_NoReceivedAtUsesTheOneSeal(t *testing.T) {
	run := buildRun(t, "RUN-A", make([]stamp, 2))
	seal := sealOver(t, dayOne, map[string][][]byte{"RUN-A": run}, everyEntry)
	runSet([]string{writeLines(t, run)}, []string{writeSeal(t, seal)}).MustSucceed(t)
}

// The supplied-file pre-flight: each of these is exit 2, before step 1.
func TestFileSet_PreflightInputErrors(t *testing.T) {
	good := buildRun(t, "RUN-A", manyAt(2, dayOne))
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": good}, everyEntry))
	noDates := buildRun(t, "RUN-A", make([]stamp, 2))
	mixed := buildRun(t, "RUN-A", []stamp{at(dayOne), {}})
	badDate := buildRun(t, "RUN-A", []stamp{{receivedAt: "yesterday"}})
	otherDay := writeSeal(t, sealOver(t, dayTwo, map[string][][]byte{"RUN-A": good}, everyEntry))

	file := writeLines(t, good)
	for name, tc := range map[string]struct {
		files, seals []string
		reason       string
	}{
		"two byte-identical files":   {[]string{file, writeLines(t, good)}, []string{seal}, "duplicate audit file supplied"},
		"mixed received_at":          {[]string{writeLines(t, mixed)}, []string{seal}, "received_at present on some entries and absent on others"},
		"two seals for one date":     {[]string{file}, []string{seal, seal}, "duplicate seal date supplied"},
		"several seals, no dates":    {[]string{writeLines(t, noDates)}, []string{seal, otherDay}, "several seals supplied and no received_at to select by"},
		"a received_at with no date": {[]string{writeLines(t, badDate)}, []string{seal}, "received_at is not an RFC 3339 timestamp"},
	} {
		t.Run(name, func(t *testing.T) {
			runSet(tc.files, tc.seals).MustExit(t, 2).StdoutContains(t, "Status: FAIL\nStep: pre-flight\nReason: "+tc.reason)
		})
	}
}

// §10.25: 2 different files for one run are a fork. It does not aggregate, and in every mode
// steps 10 and 11 do not run.
func TestFileSet_ForkFailsAtStep10(t *testing.T) {
	const fork = "Status: FAIL\nStep: 10\nReason: duplicate (tenant_id, run_id) detected: two chain files claim the same run identifier — possible fork or unauthorized duplicate genesis\n"
	branchOne := buildRun(t, "RUN-A", manyAt(2, dayOne))
	branchTwo := buildRun(t, "RUN-A", manyAt(3, dayOne))
	files := []string{writeLines(t, branchOne), writeLines(t, branchTwo)}
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": branchOne}, everyEntry))
	runSet(files, []string{seal}).MustExit(t, 1).StdoutContains(t, fork)
	clitest.Run(Main, "verify", "--audit-file", files[0], "--audit-file", files[1]).MustExit(t, 1).StdoutContains(t, fork)
}

// §7 step 10: a seal that is not daily gets its signature checked and no selection or root
// comparison. A valid signature is exit 3, and an invalid one is exit 1.
func TestFileSet_ASealThatIsNotDaily(t *testing.T) {
	run := buildRun(t, "RUN-A", manyAt(2, dayOne))
	hourly := func(date string, priv ed25519.PrivateKey) map[string]any {
		af := load(t, run)
		version, err := signpayload.ParseVersion("v1.0a")
		if err != nil {
			t.Fatal(err)
		}
		const notThisFilesRoot = "4491a35608a4acf8d5ab1bcdebfe5f8a3649a41d2b0048a9c38f848de2fe44df"
		payload, err := signpayload.Build(signpayload.Seal{
			Version: version, Algorithm: "ed25519", FormatVersion: "v1", TenantID: af.Header.TenantID, SealDate: date,
			MerkleRootHex: notThisFilesRoot, HKDFInputsDigest: af.Header.HKDFInputsDigestHex, Cadence: "hourly",
		})
		if err != nil {
			t.Fatal(err)
		}
		return map[string]any{
			"tenant_id": af.Header.TenantID, "seal_date": date, "format_version": "v1", "merkle_root_hex": notThisFilesRoot,
			"algorithm": "ed25519", "hkdf_inputs_digest_hex": af.Header.HKDFInputsDigestHex, "cadence": "hourly", "dev_mode": false,
			"sign_payload_version": "v1.0a", "sign_payload_hex": hex.EncodeToString(payload),
			"signature_b64": base64.StdEncoding.EncodeToString(ed25519.Sign(priv, payload)),
		}
	}
	file := writeLines(t, run)
	runSet([]string{file}, []string{writeSeal(t, hourly(dayOne, institutionKey()))}).
		MustExit(t, 3).
		StdoutContains(t, `seal cadence "hourly": step 10 selection is defined for daily seals only`)
	runSet([]string{file}, []string{writeSeal(t, hourly(dayOne, seededKey("the writer's own")))}).
		MustExit(t, 1).
		StdoutContains(t, signatureFail)
	// Beside a good daily seal, the run still never reads PASS.
	daily := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": run}, everyEntry))
	runSet([]string{file}, []string{daily, writeSeal(t, hourly(dayTwo, institutionKey()))}).MustExit(t, 3)
}

// Rule 6 of the unknown-kind section: the 5 stamp fields decide.
func TestFileSet_ChainStampRules(t *testing.T) {
	run := buildRun(t, "RUN-A", manyAt(2, dayOne))
	seal := writeSeal(t, sealOver(t, dayOne, map[string][][]byte{"RUN-A": run}, everyEntry))
	record := func(chain map[string]any) []byte {
		whole := map[string]any{"run_id": "RUN-A", "tenant_id": load(t, run).Header.TenantID, "kind": "a.later.kind"}
		if chain != nil {
			whole["audit_chain"] = chain
		}
		line, err := jcs.Canonicalize(whole)
		if err != nil {
			t.Fatal(err)
		}
		return line
	}
	withRecord := func(line []byte) string { return writeLines(t, append(append([][]byte{}, run...), line)) }

	t.Run("a record with none of the 5 has no leaf and is counted", func(t *testing.T) {
		runSet([]string{withRecord(record(map[string]any{"format_version": "v1"}))}, []string{seal}).
			MustSucceed(t).
			StdoutContains(t, "Anomaly: unstamped records: 1\n")
	})
	t.Run("a partial stamp fails at step 6", func(t *testing.T) {
		partial := map[string]any{"format_version": "v1", "seq": 3}
		runSet([]string{withRecord(record(partial))}, []string{seal}).
			MustExit(t, 1).
			StdoutContains(t, "Status: FAIL\nStep: 6\nReason: partial chain stamp at line 4\n")
	})
	t.Run("a record with no format_version fails at step 5", func(t *testing.T) {
		runSet([]string{withRecord(record(nil))}, []string{seal}).
			MustExit(t, 1).
			StdoutContains(t, "Status: FAIL\nStep: 5\nReason: format_version mismatch at seq 3\n")
	})
	t.Run("the optional provenance fields do not count toward the 5", func(t *testing.T) {
		provenance := map[string]any{"format_version": "v1", "kms_handle_uri": "x", "mac_computed_at_utc": "y"}
		runSet([]string{withRecord(record(provenance))}, []string{seal}).
			MustSucceed(t).
			StdoutContains(t, "Anomaly: unstamped records: 1\n")
	})
}

// The several-seal run names each seal's own posture warning.
func TestFileSet_NoPostureWarnsForEverySeal(t *testing.T) {
	stamps := append(manyAt(1, dayOne), manyAt(1, dayTwo)...)
	run := buildRun(t, "RUN-A", stamps)
	runs := map[string][][]byte{"RUN-A": run}
	first := writeSeal(t, sealOver(t, dayOne, runs, func(_ string, seq int) bool { return seq == 1 }))
	second := writeSeal(t, sealOver(t, dayTwo, runs, func(_ string, seq int) bool { return seq == 2 }))
	clitest.Run(Main, "verify", "--audit-file", writeLines(t, run), "--seal", first, "--seal", second,
		"--root-key", samplePubKey, "--master-key-hex", sampleKey).
		MustSucceed(t).
		StdoutContains(t, "Status: PASS-WITH-ANOMALY\n").
		StdoutContains(t, "at seal_date "+dayOne+":").
		StdoutContains(t, "at seal_date "+dayTwo+":")
}
