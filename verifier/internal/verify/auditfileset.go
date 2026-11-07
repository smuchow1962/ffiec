package verify

import (
	"bytes"
	"cmp"
	"crypto/ed25519"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/mmpworks/ffiec/verifier/internal/exitcodes"
)

// §7 steps 10 and 11 over a set of files and seals. Steps 1-9 run per file. Steps 10 and 11 run
// once per seal, over the entries whose effective date equals the seal's seal_date, taken from
// the union of the files in (run_id, seq) order. Every supplied entry is covered by a seal or
// the run fails, so no entry can escape both roots with a PASS.

const (
	cadenceDaily = "daily"
	forkReason   = "duplicate (tenant_id, run_id) detected: two chain files claim the same run identifier — possible fork or unauthorized duplicate genesis"
)

// SuppliedFile is one audit file as read: its name for the operator, its bytes, and its parse.
type SuppliedFile struct {
	Name  string
	Raw   []byte
	Audit *AuditFile
}

// FileSet is everything one run was given. No seals means steps 10 and 11 do not run.
type FileSet struct {
	Files []SuppliedFile
	Seals []AuditSeal
}

// PreflightFileSet is the supplied-file pre-flight of §7. Each failure is an input error with
// exit 2, because §7 has not begun.
func PreflightFileSet(set FileSet) Outcome {
	for _, check := range []func(FileSet) string{duplicateFile, receivedAtProblem, duplicateSealDate, sealsWithNoDates} {
		if reason := check(set); reason != "" {
			return structuralFail("pre-flight", reason)
		}
	}
	return pass()
}

func duplicateFile(set FileSet) string {
	for i := range set.Files {
		for j := i + 1; j < len(set.Files); j++ {
			if bytes.Equal(set.Files[i].Raw, set.Files[j].Raw) {
				return "duplicate audit file supplied"
			}
		}
	}
	return ""
}

// receivedAtProblem reads every entry's received_at once, so the walk can rely on it. The spec
// names no result for a value that is not a timestamp. An input error fails closed. A record
// that step 5 or step 6 will refuse is skipped, so the walk names the real fault.
func receivedAtProblem(set FileSet) string {
	withDate, without := 0, 0
	for _, entry := range unionEntries(set) {
		if entry.FormatVersion == "" || entry.PartialStampLine > 0 {
			continue
		}
		if entry.ReceivedAt == "" {
			without++
			continue
		}
		withDate++
		if _, err := effectiveDate(entry); err != nil {
			return fmt.Sprintf("received_at is not an RFC 3339 timestamp at seq %d: run_id %s", entry.Seq, entry.RunID)
		}
	}
	if withDate > 0 && without > 0 {
		return "received_at present on some entries and absent on others"
	}
	return ""
}

func duplicateSealDate(set FileSet) string {
	seen := map[string]bool{}
	for _, seal := range set.Seals {
		if seen[seal.SealDate] {
			return "duplicate seal date supplied"
		}
		seen[seal.SealDate] = true
	}
	return ""
}

func sealsWithNoDates(set FileSet) string {
	if len(set.Seals) > 1 && !entriesCarryDates(unionEntries(set)) {
		return "several seals supplied and no received_at to select by"
	}
	return ""
}

// WalkFileSet runs §7 over the set: steps 1-9 per file, the fork check, and then the seal steps.
// With no IKM registry it is a witness walk, so a pass keeps the witness status. A posture
// warning then adds its anomaly lines, because PASS-WITH-ANOMALY would imply that the MACs were
// checked.
func WalkFileSet(set FileSet, ikms IKMRegistry, pub ed25519.PublicKey, policy SealPolicy) Outcome {
	walked := walkEachFile(set, ikms)
	if walked.Status == "FAIL" {
		return walked
	}
	if out := checkForks(set); out.Status == "FAIL" {
		return out
	}
	if len(set.Seals) == 0 {
		return walked
	}
	judged := checkSeals(set, pub, policy)
	if judged.Status == "FAIL" || walked.Status != WitnessStatus {
		return judged
	}
	walked.Anomalies = judged.Anomalies
	return walked
}

func walkEachFile(set FileSet, ikms IKMRegistry) Outcome {
	result := pass()
	for _, file := range set.Files {
		result = walkOneFile(file.Audit, ikms)
		if result.Status == "FAIL" {
			result.File = file.Name
			return result
		}
	}
	return result
}

func walkOneFile(af *AuditFile, ikms IKMRegistry) Outcome {
	if ikms == nil {
		return WalkAuditChainWitness(af)
	}
	return WalkAuditChain(af, ikms)
}

// checkForks applies §10.25. Byte-identical files stopped at the pre-flight, so 2 files with one
// claim differ here. The header's claim is the entries' claim: the loader requires chain_id, step 4
// has held every entry's run_id equal to it, and each run_id is under the MAC.
func checkForks(set FileSet) Outcome {
	seen := map[[2]string]bool{}
	for _, file := range set.Files {
		claim := [2]string{file.Audit.Header.TenantID, file.Audit.Header.ChainID}
		if seen[claim] {
			return fail("10", forkReason)
		}
		seen[claim] = true
	}
	return pass()
}

// checkSeals orders the seals by date, so the first failure and the output do not depend on the
// order of the flags. Every integrity check on every seal runs before any posture decision.
func checkSeals(set FileSet, pub ed25519.PublicKey, policy SealPolicy) Outcome {
	seals := slices.Clone(set.Seals)
	slices.SortStableFunc(seals, func(a, b AuditSeal) int { return cmp.Compare(a.SealDate, b.SealDate) })
	headers := fileHeaders(set)
	if slices.ContainsFunc(seals, namesOtherCadence) {
		return checkSignaturesOnly(seals, headers, pub, policy)
	}
	entries := unionEntries(set)
	if out := checkEveryEntryCovered(entries, seals); out.Status == "FAIL" {
		return out
	}
	days := make([]*AuditFile, len(seals))
	for i, seal := range seals {
		days[i] = &AuditFile{Header: headers[0], Entries: entriesOf(seal, entries), Seal: seal}
		if out := checkSeal(days[i], headers, pub, policy); out.Status == "FAIL" {
			return out
		}
	}
	return judgeEverySeal(days, policy)
}

func fileHeaders(set FileSet) []AuditHeader {
	headers := make([]AuditHeader, len(set.Files))
	for i, file := range set.Files {
		headers[i] = file.Audit.Header
	}
	return headers
}

// namesOtherCadence: a seal with no sign_payload_version binds no cadence, and §7 step 10
// treats it as daily.
func namesOtherCadence(seal AuditSeal) bool {
	return seal.SignPayloadVersion != signPayloadPreAmendment && seal.Cadence != cadenceDaily
}

// checkSignaturesOnly is the §7 step 10 path for a seal that is not daily. No signed field names
// that seal's window, so no selection can be made. A bad signature stays exit 1. A good one is
// exit 3, because this run cannot judge the seal, and so the run never reads PASS.
func checkSignaturesOnly(seals []AuditSeal, headers []AuditHeader, pub ed25519.PublicKey, policy SealPolicy) Outcome {
	for _, seal := range seals {
		for _, check := range []func() Outcome{
			func() Outcome { return checkSealCoversFiles(seal, headers) },
			func() Outcome { return checkSealFloor(seal, policy) },
			func() Outcome { return checkSealAlgorithm(seal) },
			func() Outcome { return checkSealSingleSignature(seal, pub) },
		} {
			if out := check(); out.Status == "FAIL" {
				return out
			}
		}
	}
	first := seals[slices.IndexFunc(seals, namesOtherCadence)]
	return Outcome{Status: "FAIL", Step: "10", ExitCode: int(exitcodes.ConfigurationError),
		Reason: fmt.Sprintf("seal cadence %q: step 10 selection is defined for daily seals only", first.Cadence)}
}

// checkEveryEntryCovered runs before any root comparison. With no dates, the pre-flight has
// already left one seal, and it covers every entry.
func checkEveryEntryCovered(entries []*AuditEntry, seals []AuditSeal) Outcome {
	if !entriesCarryDates(entries) {
		return pass()
	}
	for _, entry := range entries {
		date, _ := effectiveDate(entry)
		if !slices.ContainsFunc(seals, func(seal AuditSeal) bool { return seal.SealDate == date }) {
			return fail("10", fmt.Sprintf("entry outside every supplied seal at seq %d: run_id %s", entry.Seq, entry.RunID))
		}
	}
	return pass()
}

// entriesOf is one seal's selection, in the union's order.
func entriesOf(seal AuditSeal, entries []*AuditEntry) []AuditEntry {
	selected := make([]AuditEntry, 0, len(entries))
	for _, entry := range entries {
		if date, _ := effectiveDate(entry); entry.ReceivedAt == "" || date == seal.SealDate {
			selected = append(selected, *entry)
		}
	}
	return selected
}

// judgeEverySeal returns the first FAIL. Otherwise it keeps every seal's warning, and the status
// line reports the first.
func judgeEverySeal(days []*AuditFile, policy SealPolicy) Outcome {
	result := pass()
	for _, day := range days {
		out := judgeSeal(day, policy)
		if out.Status == "FAIL" {
			return out
		}
		if out.Status != StatusPassWithAnomaly {
			continue
		}
		if result.Status != StatusPassWithAnomaly {
			result = Outcome{Status: out.Status, Step: out.Step, Reason: out.Reason}
		}
		result.Anomalies = append(result.Anomalies, out.Anomalies...)
	}
	return result
}

// LateBindingLines is the report that §7 and §4.2.2 require: the count, then each entry's
// (run_id, seq, received_at, captured_at). Late binding is normal operation, so these are
// anomaly lines and never a failure.
func LateBindingLines(set FileSet) []string {
	var entries []string
	for _, entry := range unionEntries(set) {
		if !entry.LateBinding {
			continue
		}
		date, _ := effectiveDate(entry)
		entries = append(entries, fmt.Sprintf("late-binding entry at seal_date %s: (%s, %d, %s, captured_at %s)",
			cmp.Or(date, "unknown"), entry.RunID, entry.Seq, entry.ReceivedAt, cmp.Or(entry.CapturedAt, "absent")))
	}
	if len(entries) == 0 {
		return nil
	}
	return append([]string{fmt.Sprintf("late-binding entries: %d", len(entries))}, entries...)
}

// SealDatesLine names the days this run covered. A run's last entries can be cut along with their
// day's seal: every supplied entry is then covered, so no step can fail. Only this line tells the
// reader that the result is about these days and no others.
func SealDatesLine(seals []AuditSeal) string {
	dates := make([]string, len(seals))
	for i, seal := range seals {
		dates[i] = seal.SealDate
	}
	slices.Sort(dates)
	return fmt.Sprintf("seal dates covered: %s: a later day's entries, cut from the end of a run along with that day's seal, cannot be detected", strings.Join(dates, ", "))
}

// unionEntries is the (run_id, seq) order of §7 step 10. Go compares strings by their UTF-8 bytes,
// the order the step names, and the sort is stable.
func unionEntries(set FileSet) []*AuditEntry {
	var entries []*AuditEntry
	for _, file := range set.Files {
		for i := range file.Audit.Entries {
			entries = append(entries, &file.Audit.Entries[i])
		}
	}
	slices.SortStableFunc(entries, func(a, b *AuditEntry) int {
		return cmp.Or(cmp.Compare(a.RunID, b.RunID), cmp.Compare(a.Seq, b.Seq))
	})
	return entries
}

func entriesCarryDates(entries []*AuditEntry) bool {
	return slices.ContainsFunc(entries, func(e *AuditEntry) bool { return e.ReceivedAt != "" })
}

// effectiveDate is the UTC date of received_at, plus one day for a late-binding entry (§4.2.2).
func effectiveDate(entry *AuditEntry) (string, error) {
	received, err := time.Parse(time.RFC3339Nano, entry.ReceivedAt)
	if err != nil {
		return "", err
	}
	day := received.UTC()
	if entry.LateBinding {
		day = day.AddDate(0, 0, 1)
	}
	return day.Format(time.DateOnly), nil
}
