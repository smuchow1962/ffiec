package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

func runVerify(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var (
		ledgerPath = fs.String("ledger", "", "Path to the tenant-day ledger file (required, unless --audit-file)")
		rootKey    = fs.String("root-key", "", "Path to the institution's seal-signing Ed25519 PUBLIC PEM (required, unless --audit-file)")
		masterKey  = fs.String("master-key", "", "Optional path to the master IKM bytes (raw 32B file). Enables full per-event verification.")
		keyHex     = fs.String("master-key-hex", "", "With --audit-file: path to the master IKM as hex text, in place of --master-key.")
		keyVersion = fs.Int("key-version", 1, "With --audit-file: the key_version the supplied IKM belongs to.")
		posture    = fs.String("posture", "", "With --seal: the institution's declared algorithm list, comma-separated, from a source outside the seal. Example: ed25519,dilithium3.")
		strict     = fs.Bool("strict", false, "With --audit-file: a posture problem fails with exit 1 and does not warn. Needs the master key, --seal, --root-key and --posture.")
		minVersion = fs.String("min-sign-payload-version", "", "With --seal: the institution's floor, from a source outside the seal. A seal below it fails.")
	)
	var auditFiles, seals cliutil.StringSlice
	fs.Var(&auditFiles, "audit-file", "Path to an on-disk NDJSON audit file, once per run. Runs the §7 walk over them. With no key, runs witness mode.")
	fs.Var(&seals, "seal", "With --audit-file: path to a day's seal record as JSON, once per day. Runs steps 10 and 11. Needs --root-key.")
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if len(auditFiles) > 0 {
		if err := emptyPolicyFlag(fs); err != nil {
			return err
		}
		if *masterKey != "" && *keyHex != "" {
			return cliutil.Usagef("--master-key and --master-key-hex: use one")
		}
		return runVerifyAuditFile(auditFileOptions{paths: auditFiles, masterKey: *masterKey, keyHexFile: *keyHex, keyVersion: *keyVersion, sealPaths: seals, rootKey: *rootKey, posture: *posture, strict: *strict, minVersion: *minVersion}, stdout, stderr)
	}
	if len(seals) > 0 {
		return cliutil.Usagef("--seal: needs --audit-file. --ledger files carry their seal inside")
	}
	if *posture != "" || *strict || *minVersion != "" {
		return cliutil.Usagef("--posture, --strict and --min-sign-payload-version: need --audit-file. The --ledger walk does not apply them yet")
	}
	if err := cliutil.RequireFlags(map[string]string{
		"ledger":   *ledgerPath,
		"root-key": *rootKey,
	}); err != nil {
		return err
	}

	led, err := verify.LoadLedger(*ledgerPath)
	if err != nil {
		return err
	}
	pub, err := cliutil.LoadEd25519PublicPEM(*rootKey)
	if err != nil {
		return fmt.Errorf("--root-key: %w", err)
	}
	plan := verify.Plan{SealingKey: pub, StopOnFirstFailure: true}
	if *masterKey != "" {
		ikm, err := os.ReadFile(*masterKey)
		if err != nil {
			return fmt.Errorf("--master-key: %w", err)
		}
		plan.MasterIKM = ikm
	}

	res, err := verify.Verify(led, plan)
	writeReport(stdout, *ledgerPath, res)
	if err != nil {
		// Print the structured report first so the failure context is
		// already visible, then surface the failure as a runtime error.
		return err
	}
	return nil
}

// writeReport prints a deterministic single-page report to w. Two verifiers
// running on identical inputs produce byte-identical reports — that is the
// property an examiner relies on when comparing notes across a team.
func writeReport(w io.Writer, path string, r *verify.Result) {
	if r == nil {
		return
	}
	printLine(w, "verifier report")
	printLine(w, "  ledger:        %s", path)
	printLine(w, "  entry_count:   %d", r.EntryCount)
	printLine(w, "  seal_date:     %s", r.SealDate)
	printLine(w, "  signing_key:   %s", r.SignerKey)
	printLine(w, "  steps:")
	for _, s := range r.Steps {
		status := "PASS"
		if !s.OK {
			status = "FAIL"
		}
		if s.Note != "" {
			printLine(w, "    [%s] %-18s — %s", status, s.Name, s.Note)
		} else {
			printLine(w, "    [%s] %s", status, s.Name)
		}
	}
	if len(r.AdditionalVerifications) > 0 {
		printLine(w, "  additional_verifications:")
		for _, av := range r.AdditionalVerifications {
			status := "PASS"
			if !av.OK {
				status = "FAIL"
			}
			if av.Note != "" {
				printLine(w, "    [%s] %-30s (%d entries) — %s", status, av.Family, av.EntryHits, av.Note)
			} else {
				printLine(w, "    [%s] %-30s (%d entries)", status, av.Family, av.EntryHits)
			}
		}
	}
	overall := "PASS (structural)"
	if r.MACPass {
		overall = "PASS (full, key-bound)"
	}
	if !r.StructuralPass {
		overall = "FAIL"
	}
	printLine(w, "  overall:       %s", overall)
}
