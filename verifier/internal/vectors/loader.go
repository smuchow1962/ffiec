package vectors

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// VectorsDirEnvVar is the env var the loader consults to find the
// test-vector corpus. Decision D-2 (locked 2026-05-21) chose env-var
// override with a sensible default so the conformance gate works on
// any contributor's machine without hardcoding an absolute path.
const VectorsDirEnvVar = "FFIEC_PUBLIC_VECTORS_DIR"

// defaultVectorsDir is the relative path the loader falls back to
// when the env var is unset. It assumes the side-by-side repo layout
// ffiec + ffiec-public ship under. Resolution is from the test
// binary's working directory at test time (typically the package
// directory containing the *_test.go file).
const defaultVectorsDir = "../ffiec-public/spec/test-vectors"

// ErrCorpusNotFound is returned when neither the env var nor the
// default path resolves to a readable test-vector corpus. Test code
// SHOULD treat this as a t.Skip reason rather than t.Fatal — a
// contributor without ffiec-public checked out should still get a
// green local build, with the gate naturally re-arming on a machine
// where the corpus IS available.
var ErrCorpusNotFound = errors.New("test-vector corpus not found")

// MasterFixturePath is the filename of the master byte-level fixture
// inside the corpus directory. Per spec/test-vectors/README.md, this
// file carries the load-bearing pinned bytes (HKDF inputs, expected
// session keys, expected fingerprints, hkdf_inputs_digest, the two
// reference chains, the resulting Merkle roots, and the sign_payload
// byte forms).
const MasterFixturePath = "chain_vectors.json"

// CorpusDir resolves the test-vector corpus directory per the
// FFIEC_PUBLIC_VECTORS_DIR env var with the side-by-side default.
// Returns ErrCorpusNotFound (wrapped with the paths attempted) when
// no candidate location is a readable directory containing the
// master fixture.
//
// Resolution order:
//
//  1. $FFIEC_PUBLIC_VECTORS_DIR (when set and non-empty).
//  2. defaultVectorsDir resolved relative to the cwd.
//  3. defaultVectorsDir resolved relative to each ancestor of the
//     cwd up to the filesystem root. This makes `go test` from any
//     package directory inside the ffiec repo find the side-by-side
//     ffiec-public/ checkout without env-var configuration —
//     CONTRIBUTING-friendly out of the box.
func CorpusDir() (string, error) {
	tried := []string{}

	if v := os.Getenv(VectorsDirEnvVar); v != "" {
		if found, abs, ok := tryCandidate(v); ok {
			return found, nil
		} else {
			tried = append(tried, abs)
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%w (getwd failed: %v)", ErrCorpusNotFound, err)
	}
	for dir := cwd; ; {
		candidate := filepath.Join(dir, defaultVectorsDir)
		if found, abs, ok := tryCandidate(candidate); ok {
			return found, nil
		} else {
			tried = append(tried, abs)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("%w (tried: %v; set %s to override)", ErrCorpusNotFound, tried, VectorsDirEnvVar)
}

// tryCandidate stats a candidate directory's master fixture; returns
// (resolvedAbsDir, attemptedAbsPath, ok). Splits out so the caller's
// loop reads cleanly without nested error handling.
func tryCandidate(dir string) (resolved, attempted string, ok bool) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Sprintf("%s (abs failed: %v)", dir, err), false
	}
	master := filepath.Join(abs, MasterFixturePath)
	if _, err := os.Stat(master); err == nil {
		return abs, abs, true
	}
	return "", abs, false
}

// MasterFixture is the typed decoding of chain_vectors.json. Only
// the fields the Commit-1 conformance gate references are decoded;
// later commits extend the struct as their tests light up.
//
// The hex-encoded fields (HKDF byte values, expected session keys,
// fingerprints, merkle roots, sign_payloads) are kept as strings on
// the struct so equality assertions report the byte-form a reader
// can visually compare against the spec's pinned values. Helpers
// like (*MasterFixture).SessionKeyV1Bytes decode on demand.
type MasterFixture struct {
	Inputs   FixtureInputs   `json:"inputs"`
	Expected FixtureExpected `json:"expected"`
	// The two reference chains aren't decoded at Commit 1 because
	// the existing verifier walks a different schema. Commit 3 adds
	// SingleChain + RotationChain []FixtureEntry fields here, once
	// the §4.4 wire types in core/wire land.
}

// FixtureInputs mirrors chain_vectors.json's "inputs" block per
// spec/test-vectors/README.md. The pinned tenant + IKM values let
// every implementation reproduce the expected outputs deterministically.
type FixtureInputs struct {
	HKDFSalt       string `json:"HKDF_SALT"`
	HKDFInfoBase   string `json:"HKDF_INFO_BASE"`
	TenantID       string `json:"tenant_id"`
	RunID          string `json:"run_id"`
	IKMv1Hex       string `json:"ikm_v1_hex"`
	IKMv2Hex       string `json:"ikm_v2_hex"`
	SealDate       string `json:"seal_date"`
	Algorithm      string `json:"algorithm"`
	FormatVersion  string `json:"format_version"`
}

// FixtureExpected mirrors chain_vectors.json's "expected" block.
// The hex-encoded byte values are the conformance pins.
type FixtureExpected struct {
	SessionKeyV1Hex      string `json:"session_key_v1_hex"`
	SessionKeyV2Hex      string `json:"session_key_v2_hex"`
	KeyFingerprintV1Hex  string `json:"key_fingerprint_v1_hex"`
	KeyFingerprintV2Hex  string `json:"key_fingerprint_v2_hex"`
	HKDFInputsDigestHex  string `json:"hkdf_inputs_digest_hex"`
	MerkleRootSingleHex  string `json:"merkle_root_single_hex"`
	MerkleRootRotationHex string `json:"merkle_root_rotation_hex"`
	SignPayloadSingleHex string `json:"sign_payload_single_hex"`
	SignPayloadRotationHex string `json:"sign_payload_rotation_hex"`
}

// LoadMasterFixture reads and decodes chain_vectors.json from the
// resolved corpus directory. Returns ErrCorpusNotFound (wrapped)
// when the corpus is not present; returns a decode error when the
// fixture is present but malformed.
func LoadMasterFixture() (*MasterFixture, error) {
	dir, err := CorpusDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, MasterFixturePath)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var fx MasterFixture
	if err := json.Unmarshal(raw, &fx); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return &fx, nil
}

// IKMv1Bytes returns the raw IKM v1 bytes decoded from the hex
// representation in the fixture. Returns a wrapped error if the hex
// is malformed (which would be a fixture-side problem).
func (f *FixtureInputs) IKMv1Bytes() ([]byte, error) {
	return decodeHexField("ikm_v1_hex", f.IKMv1Hex)
}

// IKMv2Bytes returns the raw IKM v2 bytes decoded from hex.
func (f *FixtureInputs) IKMv2Bytes() ([]byte, error) {
	return decodeHexField("ikm_v2_hex", f.IKMv2Hex)
}

func decodeHexField(name, s string) ([]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", name, err)
	}
	return b, nil
}
