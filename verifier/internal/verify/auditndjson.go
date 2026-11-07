package verify

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mmpworks/ffiec/core/jcs"
)

// §7 over the on-disk NDJSON audit file: one header line, then one event per
// line with its audit_chain stamp. This is the file a chain sink writes during
// capture, before any daily seal exists. auditseal.go adds the seal steps.

// WitnessStatus is the §7 witness-mode status line.
const WitnessStatus = "PASS-STRUCTURALLY, key-bound verification skipped"

// LoadAuditNDJSON reads an NDJSON audit file into an AuditFile. A file the
// verifier cannot begin to read comes back as a FAIL Outcome with exit code 2.
//
// The MAC covers bytes. A line that is not itself the canonical form (added
// whitespace, a reordered key, a duplicate key that a reader could resolve two
// ways) keeps its raw bytes as the canonical input, so step 9 refuses it with
// the normative reason.
func LoadAuditNDJSON(raw []byte) (*AuditFile, Outcome) {
	if len(raw) == 0 {
		return nil, structuralFail("pre-flight", "empty file: header missing")
	}
	af := &AuditFile{}
	if raw[len(raw)-1] != '\n' {
		marker := "file does not end in a newline"
		af.NDJSONTruncated = &marker
	}

	lines := bytes.Split(raw, []byte("\n"))
	header, out := loadNDJSONHeader(bytes.TrimSuffix(lines[0], []byte("\r")))
	if out.Status == "FAIL" {
		return nil, out
	}
	af.Header = header

	for i, line := range lines[1:] {
		line = bytes.TrimSuffix(line, []byte("\r"))
		if len(line) == 0 {
			continue
		}
		if af.NDJSONTruncated != nil && i == len(lines)-2 {
			break // The cut line. The pre-flight refuses this file before any entry matters.
		}
		entry, stampFields, err := loadNDJSONEntry(line)
		if err != nil {
			return nil, structuralFail("pre-flight", fmt.Sprintf("malformed entry at line %d: %v", i+2, err))
		}
		af.addRecord(entry, stampFields, i+2)
	}
	return af, pass()
}

// chainStampFields is the chain stamp of the unknown-kind section's rule 6. format_version is the
// dispatch field, and the other stamp table fields are optional provenance, so neither counts.
var chainStampFields = []string{"seq", "prev_hash", "payload_hash", "key_version", "key_fingerprint"}

// addRecord applies rule 6. A v1 record with none of the stamp is counted and left out of every
// step. A record with no format_version stays in the walk, where step 5 fails it, and a partial
// stamp stays in the walk, where step 6 fails it.
func (af *AuditFile) addRecord(entry AuditEntry, stampFields, line int) {
	if stampFields == 0 && entry.FormatVersion != "" && entry.FormatVersion == af.Header.FormatVersion {
		af.UnstampedRecords++
		return
	}
	if stampFields > 0 && stampFields < len(chainStampFields) {
		entry.PartialStampLine = line
	}
	af.Entries = append(af.Entries, entry)
}

func loadNDJSONHeader(line []byte) (AuditHeader, Outcome) {
	var h struct {
		FormatVersion    string `json:"format_version"`
		TenantID         string `json:"tenant_id"`
		ChainID          string `json:"chain_id"`
		GenesisHash      string `json:"genesis_hash"`
		HKDFInputsDigest string `json:"hkdf_inputs_digest"`
		OpenedAtUTC      string `json:"opened_at_utc"`
	}
	if err := json.Unmarshal(line, &h); err != nil {
		return AuditHeader{}, structuralFail("pre-flight", fmt.Sprintf("malformed header: %v", err))
	}
	// §4.4 requires these, and §10.12 makes a missing one exit 2. The header carries no MAC, so
	// the walk's later steps must never read an absent field as "nothing to check".
	for _, required := range []struct{ name, value string }{
		{"format_version", h.FormatVersion},
		{"tenant_id", h.TenantID},
		{"chain_id", h.ChainID},
		{"genesis_hash", h.GenesisHash},
		{"hkdf_inputs_digest", h.HKDFInputsDigest},
		{"opened_at_utc", h.OpenedAtUTC},
	} {
		if required.value == "" {
			return AuditHeader{}, structuralFail("pre-flight", "header missing required field "+required.name)
		}
	}
	return AuditHeader{
		FormatVersion:       h.FormatVersion,
		TenantID:            h.TenantID,
		ChainID:             h.ChainID,
		GenesisHashHex:      base64ToHex(h.GenesisHash),
		HKDFInputsDigestHex: base64ToHex(h.HKDFInputsDigest),
	}, pass()
}

// loadNDJSONEntry returns the record and how many chain stamp fields it holds. A record with no
// audit_chain block holds none, and it has no format_version.
func loadNDJSONEntry(line []byte) (AuditEntry, int, error) {
	decoder := json.NewDecoder(bytes.NewReader(line))
	decoder.UseNumber()
	var whole map[string]any
	if err := decoder.Decode(&whole); err != nil {
		return AuditEntry{}, 0, err
	}
	stampValue, _ := whole["audit_chain"].(map[string]any)
	stampFields := 0
	for _, field := range chainStampFields {
		if _, held := stampValue[field]; held {
			stampFields++
		}
	}

	canonicalInput, err := macInput(line, whole)
	if err != nil {
		return AuditEntry{}, 0, err
	}
	lateBinding, _ := stampValue["late_binding"].(bool)
	return AuditEntry{
		Seq:               stampInt(stampValue, "seq"),
		RunID:             stringField(whole, "run_id"),
		KeyVersion:        stampInt(stampValue, "key_version"),
		KeyFingerprintHex: base64ToHex(stringField(stampValue, "key_fingerprint")),
		FormatVersion:     stringField(stampValue, "format_version"),
		KMSHandleURI:      stringField(stampValue, "kms_handle_uri"),
		Event:             whole,
		EventCanonicalHex: hex.EncodeToString(canonicalInput),
		PrevHashHex:       base64ToHex(stringField(stampValue, "prev_hash")),
		PayloadHashHex:    base64ToHex(stringField(stampValue, "payload_hash")),
		ReceivedAt:        stringField(stampValue, "received_at"),
		LateBinding:       lateBinding,
	}, stampFields, nil
}

// macInput returns the bytes step 9 recomputes over, and removes the stamp from
// whole. §4.1 inviolate property 7: the canonical form excludes the stamp.
func macInput(line []byte, whole map[string]any) ([]byte, error) {
	asWritten, err := jcs.Canonicalize(whole)
	if err != nil {
		return nil, err
	}
	delete(whole, "audit_chain")
	if !bytes.Equal(asWritten, line) {
		return line, nil
	}
	return jcs.Canonicalize(whole)
}

func stringField(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

func stampInt(m map[string]any, key string) int {
	n, ok := m[key].(json.Number)
	if !ok {
		return 0
	}
	v, err := n.Int64()
	if err != nil {
		return 0
	}
	return int(v)
}

// base64ToHex maps a stamp's base64 hash to the walk's hex form. Text that is
// not base64 stays visibly wrong, so the comparison at its step fails.
func base64ToHex(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "not-base64:" + s
	}
	return hex.EncodeToString(b)
}

// WalkAuditChain runs the keyed walk for a file with no seal: the pre-flight and
// steps 1-9. Steps 10 and 11 need a seal. Without one, a cut at the end of the
// file is undetectable, which is why the caller adds NoSealAnomaly on this path.
func WalkAuditChain(af *AuditFile, ikms IKMRegistry) Outcome {
	if out := walkHeaderSteps(af); out.Status == "FAIL" {
		return out
	}
	if out := walkEntries(af, ikms); out.Status == "FAIL" {
		return out
	}
	return pass()
}

// WalkAuditChainWitness runs the no-key walk for a file with no seal: the
// pre-flight and steps 1-6. Steps 7-9 need the IKM.
func WalkAuditChainWitness(af *AuditFile) Outcome {
	if out := walkHeaderSteps(af); out.Status == "FAIL" {
		return out
	}
	expectedPrev := genesisZeroHex
	for i := range af.Entries {
		e := &af.Entries[i]
		if out := walkEntryStructure(af.Header, e, i+1, expectedPrev); out.Status == "FAIL" {
			return out
		}
		expectedPrev = e.PayloadHashHex
	}
	return Outcome{Status: WitnessStatus, ExitCode: 0}
}

func walkHeaderSteps(af *AuditFile) Outcome {
	for _, step := range []func() Outcome{
		func() Outcome { return checkTruncation(af) },
		func() Outcome { return checkFormatVersion(af.Header) },
		func() Outcome { return checkHKDFInputsDigest(af.Header) },
		func() Outcome { return checkGenesis(af.Header) },
		func() Outcome { return checkTenantClass(af.Header) },
	} {
		if out := step(); out.Status == "FAIL" {
			return out
		}
	}
	return pass()
}

// NDJSONAnomalies lists what a clean walk of this file did not establish, one
// line each, for the §7 anomaly lines under a passing status. raw is the file's bytes.
func NDJSONAnomalies(af *AuditFile, raw []byte) []string {
	var anomalies []string
	devKeys, noKind := 0, 0
	for i := range af.Entries {
		if strings.HasPrefix(af.Entries[i].KMSHandleURI, "plaintext-") {
			devKeys++
		}
		if _, ok := af.Entries[i].Event["chain_kind"]; !ok {
			noKind++
		}
	}
	total := len(af.Entries)
	if devKeys > 0 {
		anomalies = append(anomalies, fmt.Sprintf("%d of %d entries name a development key (kms_handle_uri plaintext-): see §10.7", devKeys, total))
	}
	if noKind > 0 {
		anomalies = append(anomalies, fmt.Sprintf("%d of %d entries carry no chain_kind: see §3", noKind, total))
	}
	if bytes.Contains(raw, []byte("\r\n")) {
		anomalies = append(anomalies, "CRLF line endings: §4.1 terminates each entry with a single 0x0A")
	}
	if af.UnstampedRecords > 0 {
		anomalies = append(anomalies, fmt.Sprintf("unstamped records: %d", af.UnstampedRecords))
	}
	return anomalies
}
