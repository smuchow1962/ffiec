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
}

// runVerifyAuditFile walks an NDJSON audit file and prints the §7 output: the
// status, then the step and reason on a failure, or anomaly lines on a pass.
//
// The walk has 4 shapes: with or without the master key, and with or without a
// seal. Each flag pair is checked before any file is read, so a wrong command
// line never looks like a verdict on the chain.
func runVerifyAuditFile(opts auditFileOptions, stdout io.Writer) error {
	if (opts.sealPath == "") != (opts.rootKey == "") {
		return configError("--seal and --root-key: use both or neither. A seal means nothing without the public key that checks it")
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
	return printOutcome(stdout, walkSealed(af, ikms, pub), anomalies)
}

func walkUnsealed(af *verify.AuditFile, ikms verify.IKMRegistry) verify.Outcome {
	if ikms == nil {
		return verify.WalkAuditChainWitness(af)
	}
	return verify.WalkAuditChain(af, ikms)
}

func walkSealed(af *verify.AuditFile, ikms verify.IKMRegistry, pub ed25519.PublicKey) verify.Outcome {
	if ikms == nil {
		return verify.WalkAuditChainWitnessSealed(af, pub)
	}
	return verify.WalkAuditChainSealed(af, ikms, pub)
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
	if out.Status == "FAIL" {
		fmt.Fprintf(w, "Step: %s\nReason: %s\n", out.Step, out.Reason)
		return &exitError{code: out.ExitCode, message: out.Reason}
	}
	for _, a := range anomalies {
		fmt.Fprintf(w, "Anomaly: %s\n", a)
	}
	return nil
}
