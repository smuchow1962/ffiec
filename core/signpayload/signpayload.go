package signpayload

import (
	"fmt"
	"strings"
)

// Magic is the fixed first line of every sign_payload form (§4.3). It
// is followed by a single 0x0A; every subsequent field is joined by a
// single 0x0A with NO trailing newline on the terminal field.
const Magic = "ffiec.chain-of-custody.v1"

// Version identifies the sign_payload byte-form generation. The zero
// value (VersionPreAmendment) is the pre-amendment 6-line form a seal
// record selects by OMITTING the sign_payload_version field.
type Version int

const (
	// VersionPreAmendment is the 6-line form for seals produced before
	// 2026-05-07 (no sign_payload_version field on the record).
	VersionPreAmendment Version = iota
	// VersionV1_0a is the 10-line amendment form.
	VersionV1_0a
	// VersionV1_0b is the 12-line form binding key_versions_canon +
	// kms_handle_uris_digest.
	VersionV1_0b
	// VersionV1_0c is the 13-line form additionally binding
	// operational_events_log_root (§10.79).
	VersionV1_0c
)

// versionTag is the ASCII string bound into line 2 of each amendment
// form. The pre-amendment form has no version line, so it is absent here.
var versionTag = map[Version]string{
	VersionV1_0a: "v1.0a",
	VersionV1_0b: "v1.0b",
	VersionV1_0c: "v1.0c",
}

// ParseVersion maps a seal record's sign_payload_version field value to
// a Version. An empty string is the pre-amendment form (the field was
// omitted). An unrecognized value is an error — the §7 step 11 dispatch
// rejects it with "sign_payload_version "X" not supported".
func ParseVersion(field string) (Version, error) {
	switch field {
	case "":
		return VersionPreAmendment, nil
	case "v1.0a":
		return VersionV1_0a, nil
	case "v1.0b":
		return VersionV1_0b, nil
	case "v1.0c":
		return VersionV1_0c, nil
	default:
		return 0, fmt.Errorf("sign_payload_version %q not supported by this verifier (running v1.0c)", field)
	}
}

// Seal carries the fields a sign_payload reconstruction consumes. A
// caller populates the fields its target Version needs; fields a given
// version does not bind are ignored. Hex fields are the lowercase ASCII
// hex text forms exactly as they appear in the seal record and the
// reconstructed payload (no decoding round-trip — the bytes signed are
// the hex text, per §4.3).
type Seal struct {
	Version Version

	Algorithm        string // e.g. "ed25519"
	FormatVersion    string // e.g. "v1"
	TenantID         string
	SealDate         string // YYYY-MM-DD UTC
	MerkleRootHex    string // 64 lowercase hex
	HKDFInputsDigest string // 64 lowercase hex
	Cadence          string // "hourly" | "daily" | "weekly" | "per_second" ...
	DevMode          bool

	// v1.0b additions:
	KeyVersionsCanon    string // ascending comma-separated decimals, "" for empty day
	KMSHandleURIsDigest string // 64 lowercase hex

	// v1.0c addition:
	OperationalEventsLogRoot string // 64 lowercase hex
}

// Build reconstructs the sign_payload bytes for the seal's Version. The
// dispatch is monotonic per §4.3: each form is the previous form plus
// terminal field(s), with the version tag distinguishing line 2. The
// terminal field carries NO trailing newline.
//
// Cognitive-complexity note: the form is a flat ordered field list per
// version. Building each as a []string then joining with "\n" keeps the
// byte construction obviously correct against the spec's line-by-line
// definition — a reviewer reads the slice top-to-bottom and matches it
// to the spec block.
func Build(s Seal) ([]byte, error) {
	fields, err := fieldsFor(s)
	if err != nil {
		return nil, err
	}
	return []byte(strings.Join(fields, "\n")), nil
}

// fieldsFor returns the ordered field list for the seal's Version. The
// first element is always the Magic line; subsequent elements are the
// version-specific fields in spec order.
func fieldsFor(s Seal) ([]string, error) {
	switch s.Version {
	case VersionPreAmendment:
		return preAmendmentFields(s), nil
	case VersionV1_0a:
		return v1_0aFields(s), nil
	case VersionV1_0b:
		return v1_0bFields(s), nil
	case VersionV1_0c:
		return v1_0cFields(s), nil
	default:
		return nil, fmt.Errorf("signpayload: unknown version %d", s.Version)
	}
}

// preAmendmentFields is the 6-line form (no version line, no cadence,
// no dev_mode) per §4.3 dispatch line 898.
func preAmendmentFields(s Seal) []string {
	return []string{
		Magic,
		s.Algorithm,
		s.FormatVersion,
		s.TenantID,
		s.SealDate,
		s.MerkleRootHex,
		s.HKDFInputsDigest,
	}
}

// v1_0aFields is the 10-line amendment form.
func v1_0aFields(s Seal) []string {
	return []string{
		Magic,
		versionTag[VersionV1_0a],
		s.Algorithm,
		s.FormatVersion,
		s.TenantID,
		s.SealDate,
		s.MerkleRootHex,
		s.HKDFInputsDigest,
		s.Cadence,
		devModeByte(s.DevMode),
	}
}

// v1_0bFields is the 12-line form (v1.0a + key_versions_canon +
// kms_handle_uris_digest).
func v1_0bFields(s Seal) []string {
	return append(replaceVersionTag(v1_0aFields(s), VersionV1_0b),
		s.KeyVersionsCanon,
		s.KMSHandleURIsDigest,
	)
}

// v1_0cFields is the 13-line form (v1.0b + operational_events_log_root).
func v1_0cFields(s Seal) []string {
	return append(replaceVersionTag(v1_0bFields(s), VersionV1_0c),
		s.OperationalEventsLogRoot,
	)
}

// replaceVersionTag swaps the version tag (line index 1) of a derived
// form so v1.0b/v1.0c reuse the v1.0a/v1.0b field order without
// hard-coding the lower form's tag. This keeps the monotonic extension
// honest: each form is literally the previous form's fields with the
// tag corrected and the new terminal field(s) appended.
func replaceVersionTag(fields []string, v Version) []string {
	out := make([]string, len(fields))
	copy(out, fields)
	if len(out) > 1 {
		out[1] = versionTag[v]
	}
	return out
}

// devModeByte serializes dev_mode as the single ASCII byte the spec
// mandates: "1" for true, "0" for false. The literal strings
// "true"/"false" or JSON booleans are non-conformant here.
func devModeByte(devMode bool) string {
	if devMode {
		return "1"
	}
	return "0"
}
