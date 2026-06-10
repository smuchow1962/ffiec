package vectors

import (
	"encoding/hex"
	"fmt"

	"github.com/mmpworks/ffiec/core/constants"
	"github.com/mmpworks/ffiec/core/hkdf"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// ---- 024: per-device HKDF derivation (§10.32) ----
//
// The §10.32 info extends the §4.1 info_for_tenant with a third
// `|`-separated device segment. The same (ikm, tenant_id) pair produces
// a distinct session key per device_id (and a distinct key from the §4.1
// baseline with no device segment). The vector pins the §4.1 baseline
// plus two device-bound keys.

func runPerDeviceDerivation(v RichVector) *Report {
	r := &Report{}

	var exp struct {
		Inputs struct {
			TenantID string `json:"tenant_id"`
			IKMHex   string `json:"ikm_hex"`
			DeviceA  string `json:"device_a"`
			DeviceB  string `json:"device_b"`
		} `json:"inputs"`
		Expected struct {
			BaselineKeyHex string `json:"section_4_1_baseline_session_key_hex"`
			DeviceAKeyHex  string `json:"section_10_32_device_a_session_key_hex"`
			DeviceBKeyHex  string `json:"section_10_32_device_b_session_key_hex"`
		} `json:"expected"`
	}
	if err := readVectorJSON(v.Dir, "expected.json", &exp); err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/read", "%v", err))
		return r
	}

	ikm, err := hex.DecodeString(exp.Inputs.IKMHex)
	if err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/decode-ikm", "%v", err))
		return r
	}

	r.Checks = append(r.Checks, checkDerivedKey(v.Slot+"/baseline-key", ikm, infoForTenant(exp.Inputs.TenantID), exp.Expected.BaselineKeyHex))
	r.Checks = append(r.Checks, checkDerivedKey(v.Slot+"/device-a-key", ikm, infoForDevice(exp.Inputs.TenantID, exp.Inputs.DeviceA), exp.Expected.DeviceAKeyHex))
	r.Checks = append(r.Checks, checkDerivedKey(v.Slot+"/device-b-key", ikm, infoForDevice(exp.Inputs.TenantID, exp.Inputs.DeviceB), exp.Expected.DeviceBKeyHex))
	return r
}

// infoForTenant builds the §4.1 baseline info: base || '|' || tenant_id.
func infoForTenant(tenantID string) []byte {
	return []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + tenantID)
}

// infoForDevice builds the §10.32 device-bound info: base || '|' ||
// tenant_id || '|' || device_id.
func infoForDevice(tenantID, deviceID string) []byte {
	return []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + tenantID + constants.InfoTenantSeparator + deviceID)
}

func checkDerivedKey(name string, ikm, info []byte, wantHex string) Check {
	got := hex.EncodeToString(hkdf.Derive(ikm, []byte(constants.HKDFSalt), info, constants.HKDFOutputLength))
	if got != wantHex {
		return failCheck(name, "session key mismatch: got=%s want=%s", got, wantHex)
	}
	return passCheck(name)
}

// ---- 022: streaming-verifier state machine (§10.29) ----
//
// The runner replays each scenario through the verifier's own
// verify.StreamState machine and asserts every step's verdict_after and
// the terminal verdict match the pins. No crypto — pure §10.29 dispatch.

func runStreamingVerifier(v RichVector) *Report {
	r := &Report{}

	var exp struct {
		Scenarios []streamScenarioJSON `json:"scenarios"`
	}
	if err := readVectorJSON(v.Dir, "expected.json", &exp); err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/read", "%v", err))
		return r
	}

	for _, sc := range exp.Scenarios {
		r.Checks = append(r.Checks, replayStreamingScenario(v.Slot, sc))
	}
	return r
}

type streamScenarioJSON struct {
	ScenarioID string `json:"scenario_id"`
	Steps      []struct {
		Step  int `json:"step"`
		Input struct {
			Kind   string `json:"kind"`
			Result string `json:"result"`
		} `json:"input"`
		VerdictBefore int `json:"verdict_before"`
		VerdictAfter  int `json:"verdict_after"`
	} `json:"steps"`
	FinalStreamingVerdict     int `json:"final_streaming_verdict"`
	TerminalVerdictAfterFinal int `json:"terminal_verdict_after_finalize"`
}

// replayStreamingScenario drives verify.StreamState through one
// scenario's steps and asserts each verdict_after plus the finalize
// collapse. Returns a single Check per scenario — a scenario passes only
// when every step transition AND the terminal collapse match.
func replayStreamingScenario(slot string, sc streamScenarioJSON) Check {
	name := fmt.Sprintf("%s/%s", slot, sc.ScenarioID)

	state := verify.NewStreamState()
	for _, s := range sc.Steps {
		if state.Verdict() != s.VerdictBefore {
			return failCheck(name, "step %d: verdict_before mismatch: machine=%d pinned=%d", s.Step, state.Verdict(), s.VerdictBefore)
		}
		state.Step(s.Input.Kind, s.Input.Result)
		if state.Verdict() != s.VerdictAfter {
			return failCheck(name, "step %d (%s/%s): verdict_after mismatch: machine=%d pinned=%d", s.Step, s.Input.Kind, s.Input.Result, state.Verdict(), s.VerdictAfter)
		}
	}

	if state.Verdict() != sc.FinalStreamingVerdict {
		return failCheck(name, "final streaming verdict mismatch: machine=%d pinned=%d", state.Verdict(), sc.FinalStreamingVerdict)
	}
	if got := state.Finalize(); got != sc.TerminalVerdictAfterFinal {
		return failCheck(name, "terminal verdict mismatch: machine=%d pinned=%d", got, sc.TerminalVerdictAfterFinal)
	}
	return passCheck(name)
}
