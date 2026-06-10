package vectors

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/mmpworks/ffiec/core/signpayload"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// signPayloadExpectedFile / signPayloadSHAFile name the sign_payload
// family's pins.
const (
	signPayloadExpectedFile = "expected_sign_payload.txt"
	signPayloadSHAFile      = "expected_sign_payload_sha256.txt"
)

// SignPayloadVector is one materialized sign_payload vector: a corpus
// directory carrying input.json (with a seal block) plus
// expected_sign_payload.txt and its sha256.
type SignPayloadVector struct {
	Slot          string
	Dir           string
	ExpectedBytes []byte // raw bytes of expected_sign_payload.txt
	ExpectedSHA   string // sha256 hex from expected_sign_payload_sha256.txt
}

// DiscoverSignPayloadVectors walks the corpus and returns every vector
// carrying the sign_payload pin pair, sorted by slot.
func DiscoverSignPayloadVectors(corpusDir string) ([]SignPayloadVector, error) {
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		return nil, fmt.Errorf("read corpus dir %s: %w", corpusDir, err)
	}

	var found []SignPayloadVector
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(corpusDir, e.Name())
		expPath := filepath.Join(dir, signPayloadExpectedFile)
		if !fileExists(expPath) {
			continue
		}
		v, err := loadSignPayloadVector(e.Name(), dir, expPath)
		if err != nil {
			return nil, fmt.Errorf("vector %s: %w", e.Name(), err)
		}
		found = append(found, v)
	}

	sort.Slice(found, func(i, j int) bool { return found[i].Slot < found[j].Slot })
	return found, nil
}

func loadSignPayloadVector(slot, dir, expPath string) (SignPayloadVector, error) {
	expected, err := os.ReadFile(expPath)
	if err != nil {
		return SignPayloadVector{}, fmt.Errorf("read %s: %w", signPayloadExpectedFile, err)
	}
	sha, err := readPrimarySHA(filepath.Join(dir, signPayloadSHAFile))
	if err != nil {
		// sha file is optional for this family in some slots; the
		// reconstruction byte-compare is the load-bearing check.
		sha = ""
	}
	return SignPayloadVector{
		Slot:          slot,
		Dir:           dir,
		ExpectedBytes: expected,
		ExpectedSHA:   sha,
	}, nil
}

// sealFields is the schema-independent set of fields the v1.0b
// reconstruction binds. Different vector layouts (the seal +
// computed_canonical_fields block of 018/019/020, the
// sign_payload_inputs block of the backfill family) decode into this
// common shape, so the reconstruction logic stays single-pathed.
type sealFields struct {
	SignPayloadVersion  string
	Algorithm           string
	FormatVersion       string
	TenantID            string
	SealDate            string
	MerkleRootHex       string
	HKDFInputsDigest    string
	Cadence             string
	DevMode             bool
	KeyVersionsCanon    string
	KMSHandleURIsDigest string
	// OperationalEventsLogRoot is the v1.0c (§10.79) sibling-log Merkle
	// root, bound as the 13th sign_payload field. Resolved from the seal
	// block's operational_events_log_root_hex; empty for v1.0a/v1.0b seals
	// where the field is absent.
	OperationalEventsLogRoot string
}

// byteFormInput is the 018/019/020 layout: a seal block carrying the
// merkle root + hkdf digest directly, plus computed_canonical_fields
// for the v1.0b additions.
type byteFormInput struct {
	Seal struct {
		SignPayloadVersion string `json:"sign_payload_version"`
		Algorithm          string `json:"algorithm"`
		FormatVersion      string `json:"format_version"`
		TenantID           string `json:"tenant_id"`
		SealDate           string `json:"seal_date"`
		MerkleRootHex      string `json:"merkle_root_hex"`
		HKDFInputsDigest   string `json:"hkdf_inputs_digest_hex"`
		Cadence            string `json:"cadence"`
		DevMode            bool   `json:"dev_mode"`
		// v1.0c (§10.79) sibling-log root. Absent on v1.0a/v1.0b seals;
		// the _hex suffix matches merkle_root_hex / hkdf_inputs_digest_hex.
		OperationalEventsLogRootHex string `json:"operational_events_log_root_hex"`
	} `json:"seal"`
	ComputedCanonicalFields struct {
		KeyVersionsCanon    string `json:"key_versions_canon"`
		KMSHandleURIsDigest string `json:"kms_handle_uris_digest_hex"`
	} `json:"computed_canonical_fields"`
}

// signPayloadInputsLayout is the backfill family layout (035): a
// sign_payload_inputs block. The merkle_root is recomputed from the
// baseline manifest + metadata leaf by the §10.42 backfill runner
// (see backfill.go); the hkdf_inputs_digest is genuinely absent from
// 035's input.json (the _compute.py hardcodes a placeholder that the
// fixture never surfaces back to input). The byte-form sign_payload
// runner therefore cannot reconstruct 035's full sign_payload from its
// input alone — it is deferred here with the precise reason, while the
// §10.42 recompute itself is gated by the backfill runner.
type signPayloadInputsLayout struct {
	SignPayloadInputs struct {
		SignPayloadVersion       string `json:"sign_payload_version"`
		Algorithm                string `json:"algorithm"`
		FormatVersion            string `json:"format_version"`
		TenantID                 string `json:"tenant_id"`
		SealDate                 string `json:"seal_date"`
		MerkleRootHex            string `json:"merkle_root_hex"`
		HKDFInputsDigest         string `json:"hkdf_inputs_digest_hex"`
		Cadence                  string `json:"cadence"`
		DevMode                  bool   `json:"dev_mode"`
		KeyVersionsCanon         string `json:"key_versions_canon"`
		KMSHandleURIsDigest      string `json:"kms_handle_uris_digest_hex"`
		OperationalEventsLogRoot string `json:"operational_events_log_root_hex"`
	} `json:"sign_payload_inputs"`
}

// errDeferredReconstruction signals a backfill-family vector whose full
// sign_payload byte-compare cannot run from input.json alone. For 035
// the merkle_root IS recomputable (the §10.42 backfill runner gates it
// against expected_merkle_root_hex.txt), but the hkdf_inputs_digest is
// absent from input.json — the _compute.py hardcodes a placeholder it
// never writes back. The byte-form runner SKIPs the full reconstruction
// with this precise reason; the load-bearing §10.42 recompute is gated
// separately by the backfill runner.
var errDeferredReconstruction = fmt.Errorf("deferred: backfill sign_payload reconstruction needs hkdf_inputs_digest_hex, which 035's input.json does not pin (§10.42 merkle_root recompute IS gated by the backfill runner)")

// resolveSealFields decodes whichever known layout the vector uses into
// the common sealFields. Returns errDeferredReconstruction when the
// layout omits the merkle root / hkdf digest (the backfill tier).
func resolveSealFields(raw []byte) (sealFields, error) {
	// Try the byte-form layout first (018/019/020).
	var bf byteFormInput
	if err := json.Unmarshal(raw, &bf); err == nil && bf.Seal.SignPayloadVersion != "" && bf.Seal.MerkleRootHex != "" {
		return sealFields{
			SignPayloadVersion:       bf.Seal.SignPayloadVersion,
			Algorithm:                bf.Seal.Algorithm,
			FormatVersion:            bf.Seal.FormatVersion,
			TenantID:                 bf.Seal.TenantID,
			SealDate:                 bf.Seal.SealDate,
			MerkleRootHex:            bf.Seal.MerkleRootHex,
			HKDFInputsDigest:         bf.Seal.HKDFInputsDigest,
			Cadence:                  bf.Seal.Cadence,
			DevMode:                  bf.Seal.DevMode,
			KeyVersionsCanon:         bf.ComputedCanonicalFields.KeyVersionsCanon,
			KMSHandleURIsDigest:      bf.ComputedCanonicalFields.KMSHandleURIsDigest,
			OperationalEventsLogRoot: bf.Seal.OperationalEventsLogRootHex,
		}, nil
	}

	// Try the sign_payload_inputs layout (backfill family).
	var spi signPayloadInputsLayout
	if err := json.Unmarshal(raw, &spi); err == nil && spi.SignPayloadInputs.SignPayloadVersion != "" {
		in := spi.SignPayloadInputs
		// The backfill family's merkle_root is DERIVED — it is the
		// §10.42 backfill root over (baseline manifest + metadata leaf),
		// not pinned directly in sign_payload_inputs. Recompute it via
		// the same wave-2 path the backfill runner gates. Once 035's
		// input.json carries hkdf_inputs_digest_hex (Heather's fix
		// 0f6825a), the full sign_payload reconstruction is unblocked.
		if in.MerkleRootHex == "" {
			root, err := recoverBackfillRoot(raw)
			if err != nil {
				return sealFields{}, errDeferredReconstruction
			}
			in.MerkleRootHex = root
		}
		if in.HKDFInputsDigest == "" {
			return sealFields{}, errDeferredReconstruction
		}
		return sealFields{
			SignPayloadVersion:       in.SignPayloadVersion,
			Algorithm:                in.Algorithm,
			FormatVersion:            in.FormatVersion,
			TenantID:                 in.TenantID,
			SealDate:                 in.SealDate,
			MerkleRootHex:            in.MerkleRootHex,
			HKDFInputsDigest:         in.HKDFInputsDigest,
			Cadence:                  in.Cadence,
			DevMode:                  in.DevMode,
			KeyVersionsCanon:         in.KeyVersionsCanon,
			KMSHandleURIsDigest:      in.KMSHandleURIsDigest,
			OperationalEventsLogRoot: in.OperationalEventsLogRoot,
		}, nil
	}

	return sealFields{}, fmt.Errorf("unrecognized sign_payload fixture layout")
}

// recoverBackfillRoot recomputes the §10.42 backfill Merkle root from the
// raw input.json's synthetic_baseline_manifest + metadata_leaf, so the
// backfill family's sign_payload can bind the derived root it never pins
// directly. Reuses the wave-2 backfill primitives (no new crypto): the
// manifest tuples + the JCS-canonicalized metadata leaf feed
// verify.RecomputeBackfillMerkleRoot.
func recoverBackfillRoot(raw []byte) (string, error) {
	var in backfillInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", fmt.Errorf("decode backfill input: %w", err)
	}
	if len(in.SyntheticBaselineMan.Tuples) == 0 || len(in.MetadataLeaf) == 0 {
		return "", fmt.Errorf("backfill input missing manifest tuples or metadata leaf")
	}
	metadataLeafJCS, err := canonicalizeMetadataLeaf(in.MetadataLeaf)
	if err != nil {
		return "", err
	}
	return verify.RecomputeBackfillMerkleRoot(in.SyntheticBaselineMan.Tuples, metadataLeafJCS)
}

// SignPayloadResult carries the outcome of running one sign_payload
// vector. Deferred is set when the vector's reconstruction depends on a
// derived field the byte-form runner does not compute (the backfill
// tier); the gate SKIPs these rather than failing.
type SignPayloadResult struct {
	Report   *Report
	Deferred bool
	Reason   string
}

// RunSignPayloadVector reconstructs the sign_payload from the vector's
// input.json via core/signpayload and asserts byte-identity with the
// pinned expected_sign_payload.txt (and sha256 when present). This is
// the byte-equivalence proof for the sign_payload byte form against the
// .NET / Python references.
//
// Vectors whose fixture omits a directly-pinned merkle root (the
// backfill family derives it from a baseline manifest) return Deferred
// with a reason — their full reconstruction lands with the
// backfill-Merkle path, not the byte-form runner.
func RunSignPayloadVector(v SignPayloadVector) SignPayloadResult {
	r := &Report{}

	raw, err := os.ReadFile(filepath.Join(v.Dir, "input.json"))
	if err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/read-input", "read input.json: %v", err))
		return SignPayloadResult{Report: r}
	}

	fields, err := resolveSealFields(raw)
	if errors.Is(err, errDeferredReconstruction) {
		return SignPayloadResult{Report: r, Deferred: true, Reason: err.Error()}
	}
	if err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/resolve-fields", "%v", err))
		return SignPayloadResult{Report: r}
	}

	seal, err := buildSeal(fields)
	if err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/version-dispatch", "%v", err))
		return SignPayloadResult{Report: r}
	}

	got, err := signpayload.Build(seal)
	if err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/reconstruct", "%v", err))
		return SignPayloadResult{Report: r}
	}

	r.Checks = append(r.Checks, checkSignPayloadBytes(v, got))
	if v.ExpectedSHA != "" {
		r.Checks = append(r.Checks, checkSignPayloadSHA(v, got))
	}
	return SignPayloadResult{Report: r}
}

// buildSeal maps resolved sealFields to a signpayload.Seal, dispatching
// the version via ParseVersion so an unrecognized version surfaces as
// the §7-step-11 error rather than a silent wrong-form reconstruction.
func buildSeal(f sealFields) (signpayload.Seal, error) {
	version, err := signpayload.ParseVersion(f.SignPayloadVersion)
	if err != nil {
		return signpayload.Seal{}, err
	}
	return signpayload.Seal{
		Version:                  version,
		Algorithm:                f.Algorithm,
		FormatVersion:            f.FormatVersion,
		TenantID:                 f.TenantID,
		SealDate:                 f.SealDate,
		MerkleRootHex:            f.MerkleRootHex,
		HKDFInputsDigest:         f.HKDFInputsDigest,
		Cadence:                  f.Cadence,
		DevMode:                  f.DevMode,
		KeyVersionsCanon:         f.KeyVersionsCanon,
		KMSHandleURIsDigest:      f.KMSHandleURIsDigest,
		OperationalEventsLogRoot: f.OperationalEventsLogRoot,
	}, nil
}

// checkSignPayloadBytes asserts the reconstructed bytes equal the pin.
// The expected file may carry a trailing newline as file framing; the
// spec's sign_payload has NO trailing newline, so we compare against
// the pin with any single trailing newline trimmed.
func checkSignPayloadBytes(v SignPayloadVector, got []byte) Check {
	name := v.Slot + "/sign-payload-bytes"
	want := bytes.TrimSuffix(v.ExpectedBytes, []byte("\n"))
	want = bytes.TrimSuffix(want, []byte("\r")) // tolerate CRLF checkout
	if !bytes.Equal(got, want) {
		return failCheck(name, "reconstructed bytes diverge\n got: %q\nwant: %q", got, want)
	}
	return passCheck(name)
}

// checkSignPayloadSHA asserts sha256(reconstructed) == pinned sha. The
// pin is sha256 over the sign_payload bytes (no trailing newline), so
// we hash the trimmed reconstruction.
func checkSignPayloadSHA(v SignPayloadVector, got []byte) Check {
	name := v.Slot + "/sign-payload-sha256"
	sum := sha256.Sum256(got)
	gotSHA := hex.EncodeToString(sum[:])
	if gotSHA != v.ExpectedSHA {
		return failCheck(name, "sha256 mismatch: got=%s want=%s", gotSHA, v.ExpectedSHA)
	}
	return passCheck(name)
}
