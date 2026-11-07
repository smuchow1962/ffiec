package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/mmpworks/ffiec/verifier/internal/exitcodes"
)

// StatusPassWithAnomaly is the §7 step 11 status for a seal whose signature verified and whose
// algorithm posture a non-strict run could not confirm. It is exit 0 (§10.12), so a reader
// must read this line.
const StatusPassWithAnomaly = "PASS-WITH-ANOMALY"

// SealPolicy is what the caller knows from outside the seal. A nil Posture means no source was
// supplied. No field of the seal can fill it: no sign_payload form binds a list of algorithms,
// so anyone who can edit the seal file could lower the posture.
type SealPolicy struct {
	Posture []string
	Strict  bool
	// MinSignPayloadVersion is the floor from the command line, source (iv) of the §7 floor check.
	// Empty means none was given.
	MinSignPayloadVersion string
}

// StrictNotRun names the --strict rules of §7 that this verifier does not apply, so a --strict
// PASS never implies them.
var StrictNotRun = []string{
	"--strict did not run step 12a: the GenAI model identifier check is not built",
	"--strict did not compare seal.cadence with the institution's claimed cadence: no input names it",
}

// signPayloadRank is the §7 floor order. A seal with no version is the pre-amendment form, below
// every floor.
var signPayloadRank = map[string]int{signPayloadPreAmendment: 0, "v1.0a": 1, signPayloadV1_0b: 2, signPayloadV1_0c: 3}

// IsSignPayloadFloor reports whether value can be a floor.
func IsSignPayloadFloor(value string) bool {
	rank, known := signPayloadRank[value]
	return known && rank > 0
}

// checkSealFloor is the §7 floor check, before signature dispatch. The seal's own field is source
// (iii). No sign_payload form binds it, so an editor can remove it, and only the flag protects
// against that. The field can only make the run fail, so honouring it lowers nothing.
func checkSealFloor(seal AuditSeal, policy SealPolicy) Outcome {
	if out := checkOneFloor(seal, policy.MinSignPayloadVersion, "--min-sign-payload-version"); out.Status == "FAIL" {
		return out
	}
	return checkOneFloor(seal, seal.DeclaredSignPayloadFloor, "declared_sign_payload_floor")
}

// checkOneFloor leaves a version it does not know to dispatch, which names it. A floor it does not
// know cannot be met, so the check fails closed.
func checkOneFloor(seal AuditSeal, floor, source string) Outcome {
	if floor == "" {
		return pass()
	}
	version, versionKnown := signPayloadRank[seal.SignPayloadVersion]
	minimum, floorKnown := signPayloadRank[floor]
	if versionKnown && (!floorKnown || version < minimum) {
		return fail("11", fmt.Sprintf("sign_payload_version %q below institution-declared floor %q (CC8.1-named source: %s)", seal.SignPayloadVersion, floor, source))
	}
	return pass()
}

// judgeSeal is the part of §7 steps 11 and 12 that reads the outside policy. It runs after the
// signature has verified.
func judgeSeal(af *AuditFile, policy SealPolicy) Outcome {
	judged := JudgeSealPosture(af.Seal, policy)
	if judged.Status != "FAIL" && policy.Strict && holdsDevelopmentKey(af) {
		return fail("12", "dev-mode seal in production verification — refused")
	}
	return judged
}

// holdsDevelopmentKey is the --strict refusal of §7 step 12, §4.2 and §10.7: a dev_mode seal, or an
// entry whose kms_handle_uri begins with "plaintext-".
func holdsDevelopmentKey(af *AuditFile) bool {
	return af.Seal.DevMode || slices.ContainsFunc(af.Entries, func(e AuditEntry) bool {
		return strings.HasPrefix(e.KMSHandleURI, "plaintext-")
	})
}

// ParseAlgorithmPosture reads the comma-separated list given on the command line, the source
// that §7 step 11 case (c) names as an inline declaration.
func ParseAlgorithmPosture(list string) ([]string, error) {
	names := strings.Split(list, ",")
	for i, name := range names {
		names[i] = strings.TrimSpace(name)
		if names[i] == "" {
			return nil, errors.New("every name in the list must be an algorithm")
		}
	}
	return names, nil
}

// JudgeSealPosture is the §7 step 11 dispatch on cases (b), (c) and (d), plus the refusal of a
// dual-signed seal. It decides only: it checks no signature, so the walk runs it after the
// seal's signature has verified, and a forged seal never reaches it.
func JudgeSealPosture(seal AuditSeal, policy SealPolicy) Outcome {
	if policy.Posture == nil {
		return judgeWithNoPosture(seal, policy.Strict)
	}
	present := sealAlgorithms(seal)
	if unknown := missingFrom(policy.Posture, present); len(unknown) > 0 {
		reason := fmt.Sprintf("algorithm not on institution's declared posture list at seal_date %s", seal.SealDate)
		out := postureFinding(reason, policy.Strict)
		out.Anomalies = append(out.Anomalies, "unknown_algorithms: ["+strings.Join(unknown, ", ")+"]")
		return out
	}
	if len(missingFrom(present, policy.Posture)) > 0 {
		return postureFinding("partial-coverage seal: single-algorithm signature during institution's declared dual-algorithm posture", policy.Strict)
	}
	return refuseSealSignaturesList(seal)
}

// judgeWithNoPosture covers the last sentences of case (c). With no list to compare, a missing
// signature cannot be told from a seal that never had one.
func judgeWithNoPosture(seal AuditSeal, strict bool) Outcome {
	if strict {
		return Outcome{Status: "FAIL", Step: "11", ExitCode: int(exitcodes.ConfigurationError),
			Reason: "--strict needs an algorithm posture source: it cannot tell case (b) from case (d)"}
	}
	if out := refuseSealSignaturesList(seal); out.Status == "FAIL" {
		return out
	}
	return postureFinding(fmt.Sprintf("algorithm posture not supplied at seal_date %s: a stripped signature cannot be detected", seal.SealDate), false)
}

// postureFinding is the one split between the modes: a FAIL under --strict, and a warning that
// never reads as a clean PASS otherwise.
func postureFinding(reason string, strict bool) Outcome {
	if strict {
		return fail("11", reason)
	}
	return Outcome{Status: StatusPassWithAnomaly, Step: "11", Reason: reason, Anomalies: []string{reason}}
}

// sealAlgorithms lists the algorithms the seal names: its own field and the list's entries.
func sealAlgorithms(seal AuditSeal) []string {
	present := []string{seal.Algorithm}
	var entries []struct {
		Algorithm string `json:"algorithm"`
	}
	if json.Unmarshal(seal.SignaturesList, &entries) != nil {
		return present
	}
	for _, entry := range entries {
		if entry.Algorithm != "" && !slices.Contains(present, entry.Algorithm) {
			present = append(present, entry.Algorithm)
		}
	}
	return present
}

// missingFrom returns the names that list does not hold, in the order given.
func missingFrom(list, names []string) []string {
	var missing []string
	for _, name := range names {
		if !slices.Contains(list, name) {
			missing = append(missing, name)
		}
	}
	return missing
}
