package verify

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
)

// §7 steps 10 and 11 over an on-disk NDJSON audit file. The chain sink writes the
// file during capture, and the seal job signs the day's root later, so the seal
// arrives as its own input. The public key is a third input, and no field of the
// audit file or the seal can supply it.

// NoSealAnomaly is the anomaly line for a clean walk that had no seal to check.
const NoSealAnomaly = "no daily seal supplied: steps 10 and 11 did not run, so a cut at the end of the file cannot be detected"

// LoadAuditSeal reads one seal record. Its field names are those of the audit_file.seal
// block in the conformance vectors. A seal the verifier cannot read means the §7
// procedure cannot begin, so it comes back as a FAIL Outcome with exit code 2.
func LoadAuditSeal(raw []byte) (AuditSeal, Outcome) {
	var seal AuditSeal
	if err := json.Unmarshal(raw, &seal); err != nil {
		return AuditSeal{}, structuralFail("pre-flight", fmt.Sprintf("malformed seal: %v", err))
	}
	return seal, pass()
}

// WalkAuditChainSealed runs the keyed walk and then the seal steps: the pre-flight
// and steps 1-11. af.Seal holds the seal.
func WalkAuditChainSealed(af *AuditFile, ikms IKMRegistry, pub ed25519.PublicKey) Outcome {
	if out := WalkAuditChain(af, ikms); out.Status == "FAIL" {
		return out
	}
	return checkSeal(af, pub)
}

// WalkAuditChainWitnessSealed runs the no-key walk and then the seal steps: the
// pre-flight and steps 1-6, 10 and 11. Steps 10 and 11 need no master key. Steps 7-9
// do, so a pass here stays a witness pass.
func WalkAuditChainWitnessSealed(af *AuditFile, pub ed25519.PublicKey) Outcome {
	witness := WalkAuditChainWitness(af)
	if witness.Status == "FAIL" {
		return witness
	}
	if out := checkSeal(af, pub); out.Status == "FAIL" {
		return out
	}
	return witness
}

func checkSeal(af *AuditFile, pub ed25519.PublicKey) Outcome {
	if out := checkMerkleRootAudit(af); out.Status == "FAIL" {
		return out
	}
	if out := checkSealCoversFile(af); out.Status == "FAIL" {
		return out
	}
	return CheckSealSignatureV1(af.Seal, pub)
}

// checkSealCoversFile refuses a seal signed for another tenant or other HKDF inputs.
// The signature binds the seal's own fields, and only this comparison binds those
// fields to this file. It uses the step 11 reason because the signature is then
// genuine for some other seal, the same failure negative vector N005 pins.
func checkSealCoversFile(af *AuditFile) Outcome {
	if af.Seal.TenantID != af.Header.TenantID || af.Seal.HKDFInputsDigestHex != af.Header.HKDFInputsDigestHex {
		return fail("11", "signature verification failed")
	}
	return pass()
}

// SealAnomalies lists what a clean seal check did not establish.
func SealAnomalies(seal AuditSeal) []string {
	if !seal.DevMode {
		return nil
	}
	return []string{"seal is dev_mode: a development key signed it, so it is not a production seal: see §10.7"}
}
