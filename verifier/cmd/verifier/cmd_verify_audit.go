package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

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
}

// runVerifyAuditFile walks an NDJSON audit file and prints the §7 output: the
// status, then the step and reason on a failure, or anomaly lines on a pass.
func runVerifyAuditFile(opts auditFileOptions, stdout io.Writer) error {
	raw, err := os.ReadFile(opts.path)
	if err != nil {
		return printOutcome(stdout, structuralOutcome(fmt.Sprintf("cannot read file: %v", err)), nil)
	}
	af, loaded := verify.LoadAuditNDJSON(raw)
	if loaded.Status == "FAIL" {
		return printOutcome(stdout, loaded, nil)
	}

	var outcome verify.Outcome
	if opts.masterKey == "" && opts.keyHexFile == "" {
		outcome = verify.WalkAuditChainWitness(af)
	} else {
		ikm, keyErr := readKey(opts)
		if keyErr != nil {
			return keyErr
		}
		outcome = verify.WalkAuditChain(af, verify.IKMRegistry{opts.keyVersion: ikm})
	}
	return printOutcome(stdout, outcome, verify.NDJSONAnomalies(af, raw))
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
