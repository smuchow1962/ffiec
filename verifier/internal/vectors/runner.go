package vectors

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/mmpworks/ffiec/core/constants"
	"github.com/mmpworks/ffiec/core/hkdf"
	"github.com/mmpworks/ffiec/core/jcs"
)

// Check is one named conformance assertion the runner executed. The
// shape is intentionally minimal — Name pins which spec property
// the check covers, Pass is the outcome, Detail carries the
// diagnostic on failure (and is empty on pass so a green run prints
// nothing of substance).
//
// Two implementations agreeing on the set of Pass=true checks have
// the same byte-equivalence posture against the corpus. Two
// implementations diverging on any Check name means the test-vector
// consumer itself disagrees on what the corpus pins, which is a
// runner-side bug, not a verifier-side bug — so the runner's Check
// names are part of the contract.
type Check struct {
	Name   string
	Pass   bool
	Detail string
}

// Report is the structured result of running all Commit-1 checks
// against the loaded fixture. The runner is deliberately
// data-oriented so a CLI subcommand or a test file can iterate the
// Checks slice and emit them in whatever shape suits the caller.
type Report struct {
	Checks []Check
}

// Pass returns true when every check in the report passed. The
// conformance-gate test asserts this; a future CLI subcommand
// asserts this too.
func (r *Report) Pass() bool {
	for _, c := range r.Checks {
		if !c.Pass {
			return false
		}
	}
	return true
}

// PassCount returns the count of checks that passed.
func (r *Report) PassCount() int {
	n := 0
	for _, c := range r.Checks {
		if c.Pass {
			n++
		}
	}
	return n
}

// RunMasterFixture executes the Commit-1 conformance checks against
// the loaded master fixture. These checks cover the pieces of the
// chain-of-custody construction that DON'T depend on RFC 8785 JCS:
//
//   - HKDF-SHA-256 of the FFIEC constants over each IKM produces the
//     expected session key (spec §4.1; pinned in chain_vectors.json).
//   - SHA-256(utf8(tenant_id) || ikm)[:16] produces the expected
//     key fingerprint (spec §4.1; pinned).
//   - SHA-256(salt || info_for_tenant || length_LE32) produces the
//     expected hkdf_inputs_digest (spec §4.2 + §3 + §4.2 sub-
//     normative on length encoding; pinned).
//
// What this DOES check at Commit 2:
//
//   - RFC 8785 JCS self-test — confirms the in-repo JCS
//     implementation reproduces the baked-in fixture byte-for-byte.
//     A failure here means the JCS implementation drifted; vector
//     008 conformance becomes unreliable from this commit forward.
//
// What this DOES NOT yet check (lands in later commits):
//
//   - Per-event MAC bytes — depends on the canonical event bytes
//     of a §4.4-wire chain entry, which lands in Commit 3.
//   - Merkle root bytes — the leaves are the per-event payload
//     hashes which depend on the §4.4 wire entries (Commit 3).
//   - sign_payload bytes — depends on the dispatch table for v1.0a /
//     v1.0b / v1.0c byte forms (Commit 5).
//   - Verdict-object JCS bytes — depends on the verdict-object
//     writer (Commit 5).
//
// Later commits extend the runner with the dependent checks as the
// implementation lands.
func RunMasterFixture(fx *MasterFixture) *Report {
	r := &Report{}

	r.Checks = append(r.Checks, checkJCSSelfTest())
	r.Checks = append(r.Checks, checkFixtureConstants(fx))
	r.Checks = append(r.Checks, checkHKDFInputsDigest(fx))

	for _, v := range ikmVariants(fx) {
		r.Checks = append(r.Checks, checkSessionKey(fx.Inputs.TenantID, v))
		r.Checks = append(r.Checks, checkKeyFingerprint(fx.Inputs.TenantID, v))
	}

	return r
}

// ikmVariant pairs one IKM hex-decoded with the expected outputs
// the fixture pins for that IKM generation.
type ikmVariant struct {
	label                  string
	ikmHex                 string
	expectedSessionKeyHex  string
	expectedFingerprintHex string
}

// ikmVariants returns the two IKM generations (v1 + v2) the master
// fixture pins, one per slice entry. Each entry carries the IKM
// itself plus the expected session-key and fingerprint outputs the
// fixture pins for that generation.
func ikmVariants(fx *MasterFixture) []ikmVariant {
	return []ikmVariant{
		{
			label:                  "v1",
			ikmHex:                 fx.Inputs.IKMv1Hex,
			expectedSessionKeyHex:  fx.Expected.SessionKeyV1Hex,
			expectedFingerprintHex: fx.Expected.KeyFingerprintV1Hex,
		},
		{
			label:                  "v2",
			ikmHex:                 fx.Inputs.IKMv2Hex,
			expectedSessionKeyHex:  fx.Expected.SessionKeyV2Hex,
			expectedFingerprintHex: fx.Expected.KeyFingerprintV2Hex,
		},
	}
}

// checkJCSSelfTest runs the JCS package's baked-in self-test. A
// failure here means the canonicalizer in core/jcs/ drifted from
// RFC 8785 — every downstream conformance assertion that depends
// on canonical bytes (Merkle leaves, per-event MAC, sign_payload,
// verdict-object) would be unreliable. Running the self-test FIRST
// in the conformance gate puts the most-load-bearing primitive
// first.
//
// The vector 008 conformance test in core/jcs/vector008_test.go is
// the deeper conformance bar (61 cases vs the self-test's 1 case);
// this check is the fast-fail signal that surfaces a JCS regression
// against the in-tree pin without requiring the test-vector corpus
// to be present.
func checkJCSSelfTest() Check {
	if err := jcs.SelfTest(); err != nil {
		return failCheck("jcs-self-test", "%v", err)
	}
	return passCheck("jcs-self-test")
}

// checkFixtureConstants asserts the fixture's HKDF inputs equal the
// in-repo constants. A failure here means the spec has rotated its
// constants and the in-repo constants are out of sync — a wire-form
// break that needs immediate attention rather than a verifier-side
// fix.
func checkFixtureConstants(fx *MasterFixture) Check {
	if fx.Inputs.HKDFSalt != constants.HKDFSalt {
		return failCheck("fixture-constants-hkdf-salt",
			"fixture HKDF_SALT %q does not match in-repo constants.HKDFSalt %q",
			fx.Inputs.HKDFSalt, constants.HKDFSalt)
	}
	if fx.Inputs.HKDFInfoBase != constants.HKDFInfoBase {
		return failCheck("fixture-constants-hkdf-info-base",
			"fixture HKDF_INFO_BASE %q does not match in-repo constants.HKDFInfoBase %q",
			fx.Inputs.HKDFInfoBase, constants.HKDFInfoBase)
	}
	if fx.Inputs.FormatVersion != constants.FormatVersion {
		return failCheck("fixture-constants-format-version",
			"fixture format_version %q does not match in-repo constants.FormatVersion %q",
			fx.Inputs.FormatVersion, constants.FormatVersion)
	}
	return passCheck("fixture-constants-match-in-repo")
}

// checkHKDFInputsDigest recomputes hkdf_inputs_digest per spec §3 +
// §4.2 normative-on-length-encoding ("length_LE32 is the 4-byte
// little-endian serialization of integer 32") and asserts byte-
// equality with the fixture's pin.
//
// The construction is:
//   SHA-256(HKDF_SALT || HKDFInfoBase || "|" || utf8(tenant_id) || length_LE32)
func checkHKDFInputsDigest(fx *MasterFixture) Check {
	infoForTenant := []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + fx.Inputs.TenantID)
	lengthLE32 := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthLE32, uint32(constants.HKDFOutputLength))

	h := sha256.New()
	h.Write([]byte(constants.HKDFSalt))
	h.Write(infoForTenant)
	h.Write(lengthLE32)
	got := hex.EncodeToString(h.Sum(nil))

	if got != fx.Expected.HKDFInputsDigestHex {
		return failCheck("hkdf-inputs-digest",
			"recomputed=%s, expected=%s", got, fx.Expected.HKDFInputsDigestHex)
	}
	return passCheck("hkdf-inputs-digest")
}

// checkSessionKey recomputes HKDF-SHA-256(IKM, salt, info, 32) and
// asserts byte-equality with the fixture's pin for the named IKM
// generation.
func checkSessionKey(tenantID string, v ikmVariant) Check {
	name := fmt.Sprintf("session-key-%s", v.label)
	ikm, err := hex.DecodeString(v.ikmHex)
	if err != nil {
		return failCheck(name, "decode ikm hex: %v", err)
	}
	info := []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + tenantID)
	sessionKey := hkdf.Derive(ikm, []byte(constants.HKDFSalt), info, constants.HKDFOutputLength)
	got := hex.EncodeToString(sessionKey)
	if got != v.expectedSessionKeyHex {
		return failCheck(name, "recomputed=%s, expected=%s", got, v.expectedSessionKeyHex)
	}
	return passCheck(name)
}

// checkKeyFingerprint recomputes SHA-256(utf8(tenant_id) || ikm)[:16]
// and asserts byte-equality with the fixture's pin.
func checkKeyFingerprint(tenantID string, v ikmVariant) Check {
	name := fmt.Sprintf("key-fingerprint-%s", v.label)
	ikm, err := hex.DecodeString(v.ikmHex)
	if err != nil {
		return failCheck(name, "decode ikm hex: %v", err)
	}
	h := sha256.New()
	h.Write([]byte(tenantID))
	h.Write(ikm)
	full := h.Sum(nil)
	got := hex.EncodeToString(full[:constants.KeyFingerprintLen])
	if got != v.expectedFingerprintHex {
		return failCheck(name, "recomputed=%s, expected=%s", got, v.expectedFingerprintHex)
	}
	return passCheck(name)
}

func passCheck(name string) Check {
	return Check{Name: name, Pass: true}
}

func failCheck(name, format string, args ...any) Check {
	return Check{Name: name, Pass: false, Detail: fmt.Sprintf(format, args...)}
}
