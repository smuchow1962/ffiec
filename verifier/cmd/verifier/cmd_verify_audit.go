package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/verifier/internal/exitcodes"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// exitError carries a §10.12 exit code back through dispatch. Main returns the
// code as it stands, because the run already printed the §7 output.
type exitError struct {
	code    int
	message string
}

func (e *exitError) Error() string { return e.message }

// auditFileOptions are the verify flags for an on-disk NDJSON audit file.
type auditFileOptions struct {
	path       string
	masterKey  string
	keyHexFile string
	keyVersion int
	sealPath   string
	rootKey    string
	posture    string
	strict     bool
	minVersion string
}

// runVerifyAuditFile walks an NDJSON audit file and prints the §7 output: the
// status, then the step and reason on a failure, or anomaly lines on a pass.
//
// The walk has 4 shapes: with or without the master key, and with or without a
// seal. Each flag pair is checked before any file is read, so a wrong command
// line never looks like a verdict on the chain.
func runVerifyAuditFile(opts auditFileOptions, stdout, stderr io.Writer) error {
	if (opts.sealPath == "") != (opts.rootKey == "") {
		return configError("--seal and --root-key: use both or neither. A seal means nothing without the public key that checks it")
	}
	policy, flagErr := sealPolicyOf(opts)
	if flagErr != nil {
		return flagErr
	}
	raw, err := os.ReadFile(opts.path)
	if err != nil {
		return printOutcome(stdout, structuralOutcome(fmt.Sprintf("cannot read file: %v", err)), nil)
	}
	af, loaded := verify.LoadAuditNDJSON(raw)
	if loaded.Status == "FAIL" {
		return printOutcome(stdout, loaded, nil)
	}

	ikms, keyErr := readRegistry(opts)
	if keyErr != nil {
		return keyErr
	}
	if opts.sealPath == "" {
		anomalies := append([]string{verify.NoSealAnomaly}, verify.NDJSONAnomalies(af, raw)...)
		return printOutcome(stdout, walkUnsealed(af, ikms), anomalies)
	}

	pub, keyErr := cliutil.LoadEd25519PublicPEM(opts.rootKey)
	if keyErr != nil {
		return configError("--root-key: %v", keyErr)
	}
	sealRaw, err := os.ReadFile(opts.sealPath)
	if err != nil {
		return printOutcome(stdout, structuralOutcome(fmt.Sprintf("cannot read seal: %v", err)), nil)
	}
	seal, sealLoaded := verify.LoadAuditSeal(sealRaw)
	if sealLoaded.Status == "FAIL" {
		return printOutcome(stdout, sealLoaded, nil)
	}
	af.Seal = seal
	anomalies := append(verify.SealAnomalies(seal), verify.NDJSONAnomalies(af, raw)...)
	if policy.Strict {
		anomalies = append(anomalies, verify.StrictNotRun...)
	}
	outcome := walkSealed(af, ikms, pub, policy)
	if outcome.Status != "FAIL" {
		// Stderr, because §7 fixes what stdout holds.
		fmt.Fprintf(stderr, "Seal checked: tenant_id %s, seal_date %s\n", seal.TenantID, seal.SealDate)
	}
	return printOutcome(stdout, outcome, anomalies)
}

func walkUnsealed(af *verify.AuditFile, ikms verify.IKMRegistry) verify.Outcome {
	if ikms == nil {
		return verify.WalkAuditChainWitness(af)
	}
	return verify.WalkAuditChain(af, ikms)
}

func walkSealed(af *verify.AuditFile, ikms verify.IKMRegistry, pub ed25519.PublicKey, policy verify.SealPolicy) verify.Outcome {
	if ikms == nil {
		return verify.WalkAuditChainWitnessSealed(af, pub, policy)
	}
	return verify.WalkAuditChainSealed(af, ikms, pub, policy)
}

// sealPolicyOf applies §10.12's exit 3 rules for --strict and --posture. --strict requires the
// key, the seal and a posture source, because it promises every check ran, and a run without
// one of them cannot keep that promise.
func sealPolicyOf(opts auditFileOptions) (verify.SealPolicy, error) {
	policy := verify.SealPolicy{Strict: opts.strict, MinSignPayloadVersion: opts.minVersion}
	if opts.minVersion != "" && (opts.sealPath == "" || !verify.IsSignPayloadFloor(opts.minVersion)) {
		return policy, configError("--min-sign-payload-version: needs --seal, and one of v1.0a, v1.0b or v1.0c")
	}
	if opts.posture != "" {
		if opts.sealPath == "" {
			return policy, configError("--posture: needs --seal. It is a list to compare a seal with")
		}
		posture, err := verify.ParseAlgorithmPosture(opts.posture)
		if err != nil {
			return policy, configError("--posture: %v", err)
		}
		policy.Posture = posture
	}
	if !opts.strict {
		return policy, nil
	}
	switch {
	case opts.masterKey == "" && opts.keyHexFile == "":
		return policy, configError("--strict: needs --master-key or --master-key-hex")
	case opts.sealPath == "":
		return policy, configError("--strict: needs --seal and --root-key")
	case policy.Posture == nil:
		return policy, configError("--strict: needs --posture. With no posture source a verifier cannot tell a stripped signature from a single-algorithm seal")
	}
	return policy, nil
}

// readRegistry returns nil when the caller supplied no master key: witness mode.
func readRegistry(opts auditFileOptions) (verify.IKMRegistry, error) {
	if opts.masterKey == "" && opts.keyHexFile == "" {
		return nil, nil
	}
	ikm, err := readKey(opts)
	if err != nil {
		return nil, err
	}
	return verify.IKMRegistry{opts.keyVersion: ikm}, nil
}

// readKey loads the IKM. A key the verifier cannot read is a §10.12
// configuration error: the walk never began.
func readKey(opts auditFileOptions) ([]byte, error) {
	if opts.masterKey != "" {
		ikm, err := os.ReadFile(opts.masterKey)
		if err != nil {
			return nil, configError("--master-key: %v", err)
		}
		return ikm, nil
	}
	text, err := os.ReadFile(opts.keyHexFile)
	if err != nil {
		return nil, configError("--master-key-hex: %v", err)
	}
	ikm, err := hex.DecodeString(strings.TrimSpace(string(text)))
	if err != nil {
		return nil, configError("--master-key-hex: %v", err)
	}
	return ikm, nil
}

func configError(format string, a ...any) error {
	return &exitError{code: int(exitcodes.ConfigurationError), message: fmt.Sprintf(format, a...)}
}

func structuralOutcome(reason string) verify.Outcome {
	return verify.Outcome{Status: "FAIL", Step: "pre-flight", Reason: reason, ExitCode: int(exitcodes.StructuralInputError)}
}

func printOutcome(w io.Writer, out verify.Outcome, anomalies []string) error {
	fmt.Fprintf(w, "Status: %s\n", out.Status)
	if out.Step != "" {
		fmt.Fprintf(w, "Step: %s\nReason: %s\n", out.Step, out.Reason)
	}
	printAnomalies(w, out.Anomalies)
	if out.Status == "FAIL" {
		return &exitError{code: out.ExitCode, message: out.Reason}
	}
	printAnomalies(w, anomalies)
	return nil
}

func printAnomalies(w io.Writer, anomalies []string) {
	for _, a := range anomalies {
		fmt.Fprintf(w, "Anomaly: %s\n", a)
	}
}
