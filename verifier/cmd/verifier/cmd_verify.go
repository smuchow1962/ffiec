package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/verifier/internal/profile"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

func runVerify(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var (
		ledgerPath = fs.String("ledger", "", "Path to the tenant-day ledger file (required)")
		rootKey    = fs.String("root-key", "", "Path to the institution's seal-signing Ed25519 PUBLIC PEM (required)")
		masterKey  = fs.String("master-key", "", "Optional path to the master IKM bytes (raw 32B file). Enables full per-event verification.")
		profileID  = fs.String("profile", "ffiec", "Examiner presentation profile: "+strings.Join(profile.IDs(), " | ")+". Presentation-only — never changes the integrity verdict.")
	)
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if err := cliutil.RequireFlags(map[string]string{
		"ledger":   *ledgerPath,
		"root-key": *rootKey,
	}); err != nil {
		return err
	}
	prof, ok := profile.Lookup(*profileID)
	if !ok {
		// Fail rather than silently fall back — an examiner must never get
		// an unexpected framing from a mistyped profile name.
		return fmt.Errorf("--profile: unknown profile %q (known: %s)", *profileID, strings.Join(profile.IDs(), ", "))
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
	writeReport(stdout, *ledgerPath, res, prof)
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
//
// The integrity core (steps, additional_verifications, overall verdict) is
// identical under every profile. The profile adds only presentation: an
// active-profile line and the §14.13 supervisory-context block, appended
// after the core. The overall verdict is read from the profile's View,
// which copies verify.Result.IntegrityVerdict() verbatim — a profile
// reframes the report but cannot change the verdict.
func writeReport(w io.Writer, path string, r *verify.Result, prof profile.Profile) {
	if r == nil {
		return
	}
	view := prof.Render(r)
	fmt.Fprintf(w, "verifier report\n")
	fmt.Fprintf(w, "  ledger:        %s\n", path)
	fmt.Fprintf(w, "  entry_count:   %d\n", r.EntryCount)
	fmt.Fprintf(w, "  seal_date:     %s\n", r.SealDate)
	fmt.Fprintf(w, "  signing_key:   %s\n", r.SignerKey)
	fmt.Fprintln(w, "  steps:")
	for _, s := range r.Steps {
		status := "PASS"
		if !s.OK {
			status = "FAIL"
		}
		if s.Note != "" {
			fmt.Fprintf(w, "    [%s] %-18s — %s\n", status, s.Name, s.Note)
		} else {
			fmt.Fprintf(w, "    [%s] %s\n", status, s.Name)
		}
	}
	if len(r.AdditionalVerifications) > 0 {
		fmt.Fprintln(w, "  additional_verifications:")
		for _, av := range r.AdditionalVerifications {
			status := "PASS"
			if !av.OK {
				status = "FAIL"
			}
			if av.Note != "" {
				fmt.Fprintf(w, "    [%s] %-30s (%d entries) — %s\n", status, av.Family, av.EntryHits, av.Note)
			} else {
				fmt.Fprintf(w, "    [%s] %-30s (%d entries)\n", status, av.Family, av.EntryHits)
			}
		}
	}
	fmt.Fprintf(w, "  overall:       %s\n", view.IntegrityVerdict)

	// Profile framing is additive and presentation-only. It is emitted
	// after the integrity core, and only when there is something profile-
	// specific to show — a non-default profile, or supervisory context on
	// the chain. The default ffiec profile on a chain with no supervisory
	// family reproduces the core report byte-for-byte.
	if prof.ID != profile.Default().ID || len(view.SupervisoryLines) > 0 {
		fmt.Fprintf(w, "  profile:       %s (%s)\n", view.ProfileID, view.ProfileDisplay)
	}
	if len(view.SupervisoryLines) > 0 {
		fmt.Fprintln(w, "  supervisory_context (institution-asserted; presentation-only):")
		for _, line := range view.SupervisoryLines {
			fmt.Fprintf(w, "    %s\n", line)
		}
	}
}
