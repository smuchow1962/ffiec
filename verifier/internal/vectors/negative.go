package vectors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// negativeDir is the corpus subdirectory holding the negative vectors
// and their authoritative INDEX.md enumeration.
const negativeDir = "negative"

// negativeIndexFile is the authoritative enumeration of negative
// vectors: each row pins (slot, target §7 step, expected reason,
// required, materialized).
const negativeIndexFile = "INDEX.md"

// NegativeVector is one row of the negative INDEX.md, joined with the
// on-disk materialization state of its directory. The verifier's
// conformance bar over negatives is: for every materialized,
// required vector, running §7 produces the expected reason string and
// the §10.12 exit code the target step implies.
type NegativeVector struct {
	Slot           string // e.g. "N001-payload-hash-bit-flip"
	Target         string // §7 step number or "pre-flight" / strict-mode
	ExpectedReason string
	Required       bool
	// Materialized is the INDEX's declared state: "yes", "stub", or
	// "deferred-v1.x". We cross-check it against the on-disk presence
	// of input.json + expected_output.txt.
	Materialized       string
	Dir                string
	OnDiskMaterialized bool // input.json AND expected_output.txt present
}

// negativeExpectedOutputFile is the per-vector pin a materialized
// negative carries: the exact Status/Step/Reason the verifier must emit.
const negativeExpectedOutputFile = "expected_output.txt"

// DiscoverNegativeVectors parses the negative INDEX.md table and joins
// each row with its directory's on-disk materialization state. The
// INDEX is the authoritative enumeration (the spec's own contract); the
// on-disk check tells the gate whether to assert (materialized) or SKIP
// with the recorded expectation (stub).
func DiscoverNegativeVectors(corpusDir string) ([]NegativeVector, error) {
	negDir := filepath.Join(corpusDir, negativeDir)
	indexPath := filepath.Join(negDir, negativeIndexFile)
	raw, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", indexPath, err)
	}

	rows := parseNegativeIndexRows(string(raw))
	vectors := make([]NegativeVector, 0, len(rows))
	for _, row := range rows {
		dir := filepath.Join(negDir, row.Slot)
		row.Dir = dir
		row.OnDiskMaterialized = fileExists(filepath.Join(dir, "input.json")) &&
			fileExists(filepath.Join(dir, negativeExpectedOutputFile))
		vectors = append(vectors, row)
	}
	return vectors, nil
}

// parseNegativeIndexRows extracts the vector rows from the INDEX.md
// markdown tables. A vector row is a pipe-delimited line whose first
// cell starts with "N" followed by digits (the slot). The header,
// separator, and conformance-discipline tables are skipped because
// their first cell does not match the slot shape.
//
// Cognitive-complexity note: the parse is a single pass over lines with
// one guard (is-this-a-vector-row). Keeping it line-oriented avoids a
// markdown-AST dependency for what is a stable, simple table shape.
func parseNegativeIndexRows(md string) []NegativeVector {
	var rows []NegativeVector
	for _, line := range strings.Split(md, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "|") {
			continue
		}
		cells := splitTableRow(line)
		if len(cells) < 5 {
			continue
		}
		slot := cells[0]
		if !looksLikeNegativeSlot(slot) {
			continue
		}
		rows = append(rows, NegativeVector{
			Slot:           slot,
			Target:         cells[1],
			ExpectedReason: unwrapCode(cells[2]),
			Required:       strings.EqualFold(cells[3], "true"),
			Materialized:   cells[4],
		})
	}
	return rows
}

// splitTableRow splits a markdown table row on '|' and trims each cell.
// Leading and trailing empty cells (from the row's outer pipes) are
// dropped.
func splitTableRow(line string) []string {
	parts := strings.Split(strings.TrimSpace(line), "|")
	cells := make([]string, 0, len(parts))
	for _, p := range parts {
		cells = append(cells, strings.TrimSpace(p))
	}
	// Drop the empty leading/trailing cells the outer pipes produce.
	if len(cells) > 0 && cells[0] == "" {
		cells = cells[1:]
	}
	if len(cells) > 0 && cells[len(cells)-1] == "" {
		cells = cells[:len(cells)-1]
	}
	return cells
}

// looksLikeNegativeSlot reports whether a cell is a vector slot like
// "N001-payload-hash-bit-flip": starts with 'N', then a digit.
func looksLikeNegativeSlot(cell string) bool {
	if len(cell) < 2 || cell[0] != 'N' {
		return false
	}
	return cell[1] >= '0' && cell[1] <= '9'
}

// unwrapCode strips a single layer of markdown inline-code backticks
// from a cell so the expected reason string is the raw text the
// verifier must emit.
func unwrapCode(cell string) string {
	return strings.Trim(cell, "`")
}

// NegativeResult is the outcome of running one negative vector. When the
// vector is not materialized on disk, the gate SKIPs it carrying the
// recorded expectation, so the conformance bar re-arms automatically
// when the spec-side fixture lands.
type NegativeResult struct {
	Slot       string
	Skipped    bool
	SkipReason string
	Report     *Report
}

// RunNegativeVector asserts a materialized negative vector, or returns a
// Skipped result carrying the recorded expectation for a stub.
//
// For a materialized vector the gate dispatches on the vector's class
// (see classifyNegative):
//
//   - live-walk classes drive the verifier's actual §7 walk / §10.42
//     recompute / §7-step-11 structural compare over the fixture and
//     assert the verifier ITSELF emits the pinned Status/Step/Reason/
//     ExitCode;
//   - contract-only vectors (whose pinned reason needs a §10.x verifier
//     path or signature crypto not built / not materialized) retain the
//     reason-template self-consistency assertion.
//
// The classification and its per-vector reasoning live in the plan doc.
// When a fixture is not materialized the gate SKIPs it, so the bar
// re-arms automatically if a fixture is removed.
func RunNegativeVector(v NegativeVector) NegativeResult {
	if !v.OnDiskMaterialized {
		return NegativeResult{
			Slot:    v.Slot,
			Skipped: true,
			SkipReason: fmt.Sprintf("not materialized (INDEX: %s); expected %q at %s",
				v.Materialized, v.ExpectedReason, v.Target),
		}
	}

	exp, err := parseExpectedOutput(v.Dir)
	if err != nil {
		r := &Report{}
		r.Checks = append(r.Checks, failCheck(v.Slot+"/parse-expected", "%v", err))
		return NegativeResult{Slot: v.Slot, Report: r}
	}

	r := &Report{Checks: runLiveNegative(v, exp)}
	return NegativeResult{Slot: v.Slot, Report: r}
}

// runMaterializedNegative compares the verifier's output for a
// materialized negative against expected_output.txt. The expected file
// pins the Status/Step/Reason triple the §7 walk must emit.
func runMaterializedNegative(v NegativeVector) Check {
	name := v.Slot + "/expected-output"
	expected, err := os.ReadFile(filepath.Join(v.Dir, negativeExpectedOutputFile))
	if err != nil {
		return failCheck(name, "read %s: %v", negativeExpectedOutputFile, err)
	}
	// The §7 walk over a real chain file is the verify package's job;
	// wiring it here lands with the materialized fixtures. Until a
	// negative is materialized this check is unreachable, so we record
	// the pin's presence and the expected reason as the assertion
	// surface the future walk plugs into.
	if !strings.Contains(string(expected), v.ExpectedReason) {
		return failCheck(name,
			"expected_output.txt does not contain the INDEX-pinned reason %q", v.ExpectedReason)
	}
	return passCheck(name)
}
