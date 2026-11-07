package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
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
	paths      []string
	masterKey  string
	keyHexFile string
	keyVersion int
	sealPaths  []string
	rootKey    string
	posture    string
	strict     bool
	minVersion string
}

// runVerifyAuditFile walks one or more NDJSON audit files and prints the §7 output: the
// status, then the step and reason on a failure, or anomaly lines on a pass.
//
// The walk has 4 shapes: with or without the master key, and with or without seals.
// Each flag pair is checked before any file is read, so a wrong command line never
// looks like a verdict on the chain.
func runVerifyAuditFile(opts auditFileOptions, stdout, stderr io.Writer) error {
	if (len(opts.sealPaths) == 0) != (opts.rootKey == "") {
		return configError("--seal and --root-key: use both or neither. A seal means nothing without the public key that checks it")
	}
	policy, flagErr := sealPolicyOf(opts)
	if flagErr != nil {
		return flagErr
	}
	files, loaded := loadAuditFiles(opts.paths)
	if loaded.Status == "FAIL" {
		return printOutcome(stdout, loaded, nil)
	}
	ikms, keyErr := readRegistry(opts)
	if keyErr != nil {
		return keyErr
	}
	var pub ed25519.PublicKey
	if opts.rootKey != "" {
		if pub, keyErr = cliutil.LoadEd25519PublicPEM(opts.rootKey); keyErr != nil {
			return configError("--root-key: %v", keyErr)
		}
	}
	seals, sealsLoaded := loadSeals(opts.sealPaths)
	if sealsLoaded.Status == "FAIL" {
		return printOutcome(stdout, sealsLoaded, nil)
	}
	set := verify.FileSet{Files: files, Seals: seals}
	if out := verify.PreflightFileSet(set); out.Status == "FAIL" {
		return printOutcome(stdout, out, nil)
	}
	outcome := verify.WalkFileSet(set, ikms, pub, policy)
	reportOnStderr(stderr, set, outcome)
	return printOutcome(stdout, outcome, setAnomalies(set, policy))
}

func loadAuditFiles(paths []string) ([]verify.SuppliedFile, verify.Outcome) {
	files := make([]verify.SuppliedFile, 0, len(paths))
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, structuralOutcome(fmt.Sprintf("cannot read file: %v", err))
		}
		af, loaded := verify.LoadAuditNDJSON(raw)
		if loaded.Status == "FAIL" {
			return nil, loaded
		}
		files = append(files, verify.SuppliedFile{Name: path, Raw: raw, Audit: af})
	}
	return files, verify.Outcome{Status: "PASS"}
}

func loadSeals(paths []string) ([]verify.AuditSeal, verify.Outcome) {
	seals := make([]verify.AuditSeal, 0, len(paths))
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, structuralOutcome(fmt.Sprintf("cannot read seal: %v", err))
		}
		seal, loaded := verify.LoadAuditSeal(raw)
		if loaded.Status == "FAIL" {
			return nil, loaded
		}
		seals = append(seals, seal)
	}
	return seals, verify.Outcome{Status: "PASS"}
}

// reportOnStderr carries what §7 leaves out of stdout: which seals were checked, and which of
// several files a per-file step failed in.
func reportOnStderr(stderr io.Writer, set verify.FileSet, outcome verify.Outcome) {
	if outcome.File != "" && len(set.Files) > 1 {
		fmt.Fprintf(stderr, "In file: %s\n", cliutil.SafeLine(outcome.File))
	}
	if outcome.Status == "FAIL" {
		return
	}
	for _, seal := range set.Seals {
		fmt.Fprintf(stderr, "Seal checked: tenant_id %s, seal_date %s\n", cliutil.SafeLine(seal.TenantID), cliutil.SafeLine(seal.SealDate))
	}
}

// setAnomalies names what a clean walk did not establish. With several files or seals, each
// line names the one it is about.
func setAnomalies(set verify.FileSet, policy verify.SealPolicy) []string {
	var anomalies []string
	if len(set.Seals) == 0 {
		anomalies = append(anomalies, verify.NoSealAnomaly)
	} else {
		anomalies = append(anomalies, verify.SealDatesLine(set.Seals))
	}
	anomalies = append(anomalies, verify.LateBindingLines(set)...)
	for _, seal := range set.Seals {
		anomalies = append(anomalies, labelled(seal.SealDate, len(set.Seals), verify.SealAnomalies(seal))...)
	}
	for _, file := range set.Files {
		anomalies = append(anomalies, labelled(file.Audit.Header.ChainID, len(set.Files), verify.NDJSONAnomalies(file.Audit, file.Raw))...)
	}
	if policy.Strict {
		anomalies = append(anomalies, verify.StrictNotRun...)
	}
	return anomalies
}

func labelled(label string, count int, lines []string) []string {
	if count < 2 {
		return lines
	}
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = label + ": " + line
	}
	return out
}

// emptyPolicyFlag finds --posture= or --min-sign-payload-version= with no value. Read as "not given",
// either one would silently drop a check the command line asked for.
func emptyPolicyFlag(fs *flag.FlagSet) error {
	var err error
	fs.Visit(func(f *flag.Flag) {
		if (f.Name == "posture" || f.Name == "min-sign-payload-version") && f.Value.String() == "" {
			err = configError("--%s: the value is empty", f.Name)
		}
	})
	return err
}

// sealPolicyOf applies §10.12's exit 3 rules for --strict and --posture. --strict requires the
// key, the seal and a posture source, because it promises every check ran, and a run without
// one of them cannot keep that promise.
func sealPolicyOf(opts auditFileOptions) (verify.SealPolicy, error) {
	policy := verify.SealPolicy{Strict: opts.strict, MinSignPayloadVersion: opts.minVersion}
	if opts.minVersion != "" && (len(opts.sealPaths) == 0 || !verify.IsSignPayloadFloor(opts.minVersion)) {
		return policy, configError("--min-sign-payload-version: needs --seal, and one of v1.0a, v1.0b or v1.0c")
	}
	if opts.posture != "" {
		if len(opts.sealPaths) == 0 {
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
	case len(opts.sealPaths) == 0:
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
		fmt.Fprintf(w, "Step: %s\nReason: %s\n", out.Step, cliutil.SafeLine(out.Reason))
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
		fmt.Fprintf(w, "Anomaly: %s\n", cliutil.SafeLine(a))
	}
}

// printLine writes one line through the output gate. Its format holds no line break.
func printLine(w io.Writer, format string, a ...any) {
	fmt.Fprintln(w, cliutil.SafeLine(fmt.Sprintf(format, a...)))
}
