package vectors

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mmpworks/ffiec/core/jcs"
)

// CanonicalVector is one materialized canonical-output vector: a corpus
// directory carrying expected_canonical.txt + expected_canonical_sha256.txt.
// The Slot is the directory name (e.g. "049-unknown-wire-format-kind-fallthrough").
type CanonicalVector struct {
	Slot      string
	Dir       string
	Canonical []byte // raw bytes of expected_canonical.txt
	// PrimarySHA256 is the first row of expected_canonical_sha256.txt —
	// the sha256 of the whole expected_canonical.txt blob.
	PrimarySHA256 string
}

// canonicalExpectedFile / canonicalSHAFile name the two pins the
// canonical-output family carries.
const (
	canonicalExpectedFile = "expected_canonical.txt"
	canonicalSHAFile      = "expected_canonical_sha256.txt"
)

// DiscoverCanonicalVectors walks the corpus directory and returns every
// materialized canonical-output vector (those carrying both
// expected_canonical.txt and expected_canonical_sha256.txt), sorted by
// slot for deterministic test ordering. Directories that lack the pair
// are skipped — they belong to a different vector family (sign_payload,
// rich-expected) or are description-only stubs.
func DiscoverCanonicalVectors(corpusDir string) ([]CanonicalVector, error) {
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		return nil, fmt.Errorf("read corpus dir %s: %w", corpusDir, err)
	}

	var found []CanonicalVector
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(corpusDir, e.Name())
		expPath := filepath.Join(dir, canonicalExpectedFile)
		shaPath := filepath.Join(dir, canonicalSHAFile)
		if !fileExists(expPath) || !fileExists(shaPath) {
			continue
		}
		v, err := loadCanonicalVector(e.Name(), dir, expPath, shaPath)
		if err != nil {
			return nil, fmt.Errorf("vector %s: %w", e.Name(), err)
		}
		found = append(found, v)
	}

	sort.Slice(found, func(i, j int) bool { return found[i].Slot < found[j].Slot })
	return found, nil
}

// loadCanonicalVector reads one vector's expected_canonical.txt and the
// primary (first-row) sha256 from expected_canonical_sha256.txt.
func loadCanonicalVector(slot, dir, expPath, shaPath string) (CanonicalVector, error) {
	canonical, err := os.ReadFile(expPath)
	if err != nil {
		return CanonicalVector{}, fmt.Errorf("read %s: %w", canonicalExpectedFile, err)
	}
	primarySHA, err := readPrimarySHA(shaPath)
	if err != nil {
		return CanonicalVector{}, err
	}
	return CanonicalVector{
		Slot:          slot,
		Dir:           dir,
		Canonical:     canonical,
		PrimarySHA256: primarySHA,
	}, nil
}

// readPrimarySHA extracts the sha256 hex from the first non-empty line
// of expected_canonical_sha256.txt. Each row is "<label> <hex>" (e.g.
// "full_blob <hex>" or "structured_output <hex>"); the first row always
// pins sha256(expected_canonical.txt) per the corpus convention.
func readPrimarySHA(shaPath string) (string, error) {
	raw, err := os.ReadFile(shaPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", canonicalSHAFile, err)
	}
	for _, line := range strings.Split(string(raw), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		// Last field is the hex; a single-field line is hex-only.
		hexVal := fields[len(fields)-1]
		if isSHA256Hex(hexVal) {
			return hexVal, nil
		}
	}
	return "", fmt.Errorf("%s: no sha256 hex found", canonicalSHAFile)
}

// RunCanonicalVector executes the three-tier conformance assertion for
// one canonical-output vector and returns a Report. The checks are:
//
//  1. self-consistency: sha256(expected_canonical.txt) == primary pin.
//     Catches corpus corruption — the bytes and their declared hash
//     must agree before any verifier-side claim.
//  2. JCS idempotence: re-canonicalizing the parsed expected bytes
//     through core/jcs reproduces them byte-for-byte. THIS is the
//     verifier-side proof — it shows the Go canonicalizer agrees with
//     whatever produced the pinned bytes (the .NET / Python references).
//  3. embedded-bytes pins: each *_canonical_bytes_utf8 field in
//     input.json hashes to its declared *_canonical_sha256 sibling and
//     is itself JCS-idempotent.
func RunCanonicalVector(v CanonicalVector) *Report {
	r := &Report{}
	r.Checks = append(r.Checks, checkCanonicalSelfConsistency(v))
	r.Checks = append(r.Checks, checkCanonicalJCSIdempotent(v.Slot, "expected_canonical", v.Canonical))
	r.Checks = append(r.Checks, checkEmbeddedCanonicalBytes(v)...)
	return r
}

// checkCanonicalSelfConsistency asserts sha256(canonical) == primary pin.
func checkCanonicalSelfConsistency(v CanonicalVector) Check {
	name := fmt.Sprintf("%s/self-consistency", v.Slot)
	sum := sha256.Sum256(v.Canonical)
	got := hex.EncodeToString(sum[:])
	if got != v.PrimarySHA256 {
		return failCheck(name, "sha256(expected_canonical.txt)=%s, primary pin=%s", got, v.PrimarySHA256)
	}
	return passCheck(name)
}

// checkCanonicalJCSIdempotent re-canonicalizes the pinned bytes through
// core/jcs and asserts byte-identity. A divergence means the Go
// canonicalizer disagrees with the reference that produced the pin — a
// conformance break.
//
// The canonical files are newline-delimited JSON: a single-document
// vector is one line; a multi-record vector (e.g. a chain entry plus
// its verdict, or one verdict per declared posture) is one canonical
// document per line. We re-canonicalize each non-empty line
// independently and reassemble with the same newline framing, so the
// idempotence check holds whether the file carries one document or many.
func checkCanonicalJCSIdempotent(slot, label string, canonical []byte) Check {
	name := fmt.Sprintf("%s/%s-jcs-idempotent", slot, label)

	lines := splitNDJSON(canonical)
	if len(lines) == 0 {
		return failCheck(name, "no JSON documents in canonical bytes")
	}

	reassembled := make([][]byte, 0, len(lines))
	for i, line := range lines {
		recanon, err := recanonicalizeLine(line)
		if err != nil {
			return failCheck(name, "line %d: %v", i, err)
		}
		reassembled = append(reassembled, recanon)
	}

	got := bytes.Join(reassembled, []byte("\n"))
	if !bytes.Equal(got, bytes.Join(lines, []byte("\n"))) {
		return failCheck(name, "re-canonicalized bytes diverge\n got: %s\nwant: %s", got, canonical)
	}
	return passCheck(name)
}

// recanonicalizeLine decodes one canonical JSON document (preserving
// integral digits via UseNumber) and re-emits it through core/jcs.
func recanonicalizeLine(line []byte) ([]byte, error) {
	var value any
	dec := json.NewDecoder(bytes.NewReader(line))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return jcs.Canonicalize(value)
}

// splitNDJSON splits canonical bytes on newlines and drops empty lines
// (a trailing newline in the file is framing, not a document). It
// returns the non-empty document lines with no trailing whitespace.
func splitNDJSON(canonical []byte) [][]byte {
	var out [][]byte
	for _, line := range bytes.Split(canonical, []byte("\n")) {
		trimmed := bytes.TrimRight(line, "\r")
		if len(bytes.TrimSpace(trimmed)) == 0 {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

// embeddedCanonicalSuffix / embeddedSHASuffix name the paired-attribute
// convention the corpus uses inside input.json: a "<x>_canonical_bytes_utf8"
// field is pinned by a sibling "<x>_canonical_sha256" field.
const (
	embeddedCanonicalSuffix = "_canonical_bytes_utf8"
	embeddedSHASuffix       = "_canonical_sha256"
)

// checkEmbeddedCanonicalBytes finds every *_canonical_bytes_utf8 field
// in the vector's input.json, asserts each hashes to its paired
// *_canonical_sha256, and is JCS-idempotent. Vectors with no input.json
// (or no embedded fields) contribute zero checks — the
// expected_canonical.txt checks above already gate them.
func checkEmbeddedCanonicalBytes(v CanonicalVector) []Check {
	inputPath := filepath.Join(v.Dir, "input.json")
	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return nil // no input.json — nothing embedded to check
	}

	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return []Check{failCheck(v.Slot+"/input-parse", "parse input.json: %v", err)}
	}

	pairs := collectEmbeddedPairs(top)
	checks := make([]Check, 0, len(pairs))
	for _, p := range pairs {
		checks = append(checks, runEmbeddedPair(v.Slot, p))
	}
	return checks
}

// embeddedPair is one (canonical bytes, expected sha) pairing found in
// input.json, keyed by the shared base name.
type embeddedPair struct {
	base      string
	canonical string
	sha       string
}

// collectEmbeddedPairs scans the top-level input.json object for
// "<base>_canonical_bytes_utf8" string fields and pairs each with its
// "<base>_canonical_sha256" sibling when present.
func collectEmbeddedPairs(top map[string]json.RawMessage) []embeddedPair {
	bytesByBase := map[string]string{}
	shaByBase := map[string]string{}

	for key, rawVal := range top {
		if base, ok := strings.CutSuffix(key, embeddedCanonicalSuffix); ok {
			if s, ok := decodeStringField(rawVal); ok {
				bytesByBase[base] = s
			}
		}
		if base, ok := strings.CutSuffix(key, embeddedSHASuffix); ok {
			if s, ok := decodeStringField(rawVal); ok {
				shaByBase[base] = s
			}
		}
	}

	var pairs []embeddedPair
	bases := make([]string, 0, len(bytesByBase))
	for base := range bytesByBase {
		bases = append(bases, base)
	}
	sort.Strings(bases)
	for _, base := range bases {
		pairs = append(pairs, embeddedPair{
			base:      base,
			canonical: bytesByBase[base],
			sha:       shaByBase[base], // empty when no sibling sha
		})
	}
	return pairs
}

// runEmbeddedPair asserts one embedded canonical-bytes field hashes to
// its declared sha (when present) and is JCS-idempotent.
func runEmbeddedPair(slot string, p embeddedPair) Check {
	name := fmt.Sprintf("%s/embedded-%s", slot, p.base)

	if p.sha != "" {
		sum := sha256.Sum256([]byte(p.canonical))
		got := hex.EncodeToString(sum[:])
		if got != p.sha {
			return failCheck(name, "sha256 mismatch: got=%s want=%s", got, p.sha)
		}
	}

	// JCS idempotence on the embedded bytes.
	idem := checkCanonicalJCSIdempotent(slot, "embedded-"+p.base, []byte(p.canonical))
	if !idem.Pass {
		return idem
	}
	return passCheck(name)
}

func decodeStringField(raw json.RawMessage) (string, bool) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	return s, true
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isSHA256Hex(s string) bool {
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
