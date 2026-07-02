package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mmpworks/ffiec/ledger/internal/examiner"
	"github.com/mmpworks/ffiec/testkit/clitest"
)

// ----- Top-level dispatch -------------------------------------------------

func TestMain_NoArgsPrintsUsage(t *testing.T) {
	clitest.Run(Main).
		MustExit(t, 2).
		StderrContains(t, "ledgerctl —").
		StderrContains(t, "examiner issue")
}

func TestMain_HelpPrintsToStdoutAndExitsZero(t *testing.T) {
	clitest.Run(Main, "--help").
		MustSucceed(t).
		StdoutContains(t, "ledgerctl —")
}

func TestMain_UnknownCommandIsUsageError(t *testing.T) {
	clitest.Run(Main, "frobnicate").
		MustExit(t, 2).
		StderrContains(t, `unknown command "frobnicate"`)
}

func TestMain_ExaminerWithoutSubcommandIsUsageError(t *testing.T) {
	clitest.Run(Main, "examiner").
		MustExit(t, 2).
		StderrContains(t, "missing subcommand")
}

func TestMain_KeyWithoutSubcommandIsUsageError(t *testing.T) {
	clitest.Run(Main, "key").
		MustExit(t, 2).
		StderrContains(t, "missing subcommand")
}

// ----- key generate -------------------------------------------------------

func TestKeyGenerate_ProducesValidPKCS8(t *testing.T) {
	keyPath := clitest.TempPath(t, "issuance.key.pem")
	clitest.Run(Main, "key", "generate", "--out", keyPath).
		MustSucceed(t).
		StdoutContains(t, "issuance_key_fingerprint: sha256:")

	raw := clitest.ReadFile(t, keyPath)
	block, _ := pem.Decode(raw)
	if block == nil || block.Type != "PRIVATE KEY" {
		t.Fatalf("expected PRIVATE KEY PEM block, got %+v", block)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("ParsePKCS8PrivateKey: %v", err)
	}
	if _, ok := key.(ed25519.PrivateKey); !ok {
		t.Fatalf("expected Ed25519 private key, got %T", key)
	}

	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("key file is empty")
	}
}

func TestKeyGenerate_RequiresOutFlag(t *testing.T) {
	clitest.Run(Main, "key", "generate").
		MustExit(t, 2).
		StderrContains(t, "--out")
}

func TestKeyGenerate_HelpExitsZero(t *testing.T) {
	clitest.Run(Main, "key", "generate", "--help").
		MustSucceed(t).
		StdoutContains(t, "Usage of ledgerctl key generate")
}

// ----- examiner issue (end-to-end) ---------------------------------------

func TestExaminerIssue_EndToEnd(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "issuance.key.pem")
	bundlePath := filepath.Join(dir, "karen.bundle")

	clitest.Run(Main, "key", "generate", "--out", keyPath).MustSucceed(t)

	r := clitest.Run(Main,
		"examiner", "issue",
		"--institution-id", "northbridge-federal",
		"--identity", "OCC-IT-EX-Cert-2026-Karen",
		"--tenant", "tenant_a",
		"--tenant", "tenant_b",
		"--period-start", "2026-01-01",
		"--period-end", "2026-06-14",
		"--expires-at", futureExpiry(t),
		"--correlation-id", "exam-2026-06-15-occ",
		"--scope", "compliance.read",
		"--issuance-key", keyPath,
		"--compliance-url", "https://compliance.northbridge.example",
		"--out", bundlePath,
	).MustSucceed(t)
	r.StderrContains(t, "issued bnd_")
	r.StderrContains(t, bundlePath)

	bundleRaw := clitest.ReadFile(t, bundlePath)
	var b examiner.Bundle
	if err := json.Unmarshal(bundleRaw, &b); err != nil {
		t.Fatalf("unmarshal bundle: %v", err)
	}
	if !strings.HasPrefix(b.BundleID, "bnd_") {
		t.Errorf("bundle id %q does not look minted", b.BundleID)
	}
	if got, want := b.Issuer.InstitutionID, "northbridge-federal"; got != want {
		t.Errorf("InstitutionID = %q, want %q", got, want)
	}
	if got, want := len(b.Scope.Tenants), 2; got != want {
		t.Errorf("scope.tenants len = %d, want %d", got, want)
	}
	if b.Scope.Surface != examiner.ScopeComplianceRead {
		t.Errorf("scope.surface = %q, want %q", b.Scope.Surface, examiner.ScopeComplianceRead)
	}
	if b.SignatureEd25519 == "" {
		t.Error("bundle has no signature")
	}

	// Stdout carries the chain events — one per tenant.
	dec := json.NewDecoder(strings.NewReader(r.Stdout))
	var events []examiner.OperationalEvent
	for {
		var e examiner.OperationalEvent
		if err := dec.Decode(&e); err != nil {
			break
		}
		events = append(events, e)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 chain events on stdout, got %d\nstdout: %s", len(events), r.Stdout)
	}
	for i, e := range events {
		if e.Event != examiner.EventHandoverReceived {
			t.Errorf("event[%d].Event = %q, want %q", i, e.Event, examiner.EventHandoverReceived)
		}
		if e.CorrelationID != "exam-2026-06-15-occ" {
			t.Errorf("event[%d].CorrelationID = %q", i, e.CorrelationID)
		}
	}
	if events[0].TenantID != "tenant_a" || events[1].TenantID != "tenant_b" {
		t.Errorf("tenant fan-out wrong: %q, %q", events[0].TenantID, events[1].TenantID)
	}

	// The bundle round-trips through Verify with the issuance public key.
	pub := publicFromPKCS8File(t, keyPath)
	if _, err := examiner.Verify(bundleRaw, pub); err != nil {
		t.Errorf("Verify: %v", err)
	}
}

func TestExaminerIssue_MissingFlagsListsAllMissing(t *testing.T) {
	r := clitest.Run(Main, "examiner", "issue").MustExit(t, 2)
	r.StderrContains(t, "missing required flags")
	for _, f := range []string{
		"--institution-id", "--identity", "--period-start", "--period-end",
		"--expires-at", "--correlation-id", "--issuance-key", "--compliance-url", "--out",
	} {
		r.StderrContains(t, f)
	}
}

func TestExaminerIssue_RequiresAtLeastOneTenant(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "k.pem")
	bundlePath := filepath.Join(dir, "b.bundle")
	clitest.Run(Main, "key", "generate", "--out", keyPath).MustSucceed(t)

	clitest.Run(Main,
		"examiner", "issue",
		"--institution-id", "x",
		"--identity", "x",
		"--period-start", "2026-01-01",
		"--period-end", "2026-06-14",
		"--expires-at", futureExpiry(t),
		"--correlation-id", "c",
		"--issuance-key", keyPath,
		"--compliance-url", "https://example",
		"--out", bundlePath,
	).MustExit(t, 2).StderrContains(t, "--tenant")
}

func TestExaminerIssue_RejectsBadExpires(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "k.pem")
	bundlePath := filepath.Join(dir, "b.bundle")
	clitest.Run(Main, "key", "generate", "--out", keyPath).MustSucceed(t)

	clitest.Run(Main,
		"examiner", "issue",
		"--institution-id", "x",
		"--identity", "x",
		"--tenant", "t",
		"--period-start", "2026-01-01",
		"--period-end", "2026-06-14",
		"--expires-at", "not-a-date",
		"--correlation-id", "c",
		"--issuance-key", keyPath,
		"--compliance-url", "https://example",
		"--out", bundlePath,
	).MustExit(t, 2).StderrContains(t, "--expires-at")
}

func TestExaminerIssue_RejectsBadScope(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "k.pem")
	bundlePath := filepath.Join(dir, "b.bundle")
	clitest.Run(Main, "key", "generate", "--out", keyPath).MustSucceed(t)

	// Bad scope passes the flag-required check but fails Validate inside
	// Issue → exit 1 (runtime error), not 2.
	clitest.Run(Main,
		"examiner", "issue",
		"--institution-id", "x",
		"--identity", "x",
		"--tenant", "t",
		"--period-start", "2026-01-01",
		"--period-end", "2026-06-14",
		"--expires-at", futureExpiry(t),
		"--correlation-id", "c",
		"--scope", "bogus.scope",
		"--issuance-key", keyPath,
		"--compliance-url", "https://example",
		"--out", bundlePath,
	).MustExit(t, 1).StderrContains(t, "scope.surface")
}

func TestExaminerIssue_RejectsMissingKeyFile(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "b.bundle")

	clitest.Run(Main,
		"examiner", "issue",
		"--institution-id", "x",
		"--identity", "x",
		"--tenant", "t",
		"--period-start", "2026-01-01",
		"--period-end", "2026-06-14",
		"--expires-at", futureExpiry(t),
		"--correlation-id", "c",
		"--issuance-key", filepath.Join(dir, "does-not-exist.pem"),
		"--compliance-url", "https://example",
		"--out", bundlePath,
	).MustExit(t, 1).StderrContains(t, "--issuance-key")
}

func TestExaminerIssue_HelpExitsZero(t *testing.T) {
	// flag.PrintDefaults emits single-dash flag names ("-issuance-key").
	// Both spellings are accepted on input; the help text uses the short form.
	clitest.Run(Main, "examiner", "issue", "--help").
		MustSucceed(t).
		StdoutContains(t, "Usage of ledgerctl examiner issue").
		StdoutContains(t, "issuance-key")
}

// ----- examiner revoke ----------------------------------------------------

func TestExaminerRevoke_EmitsReturnEventPerTenant(t *testing.T) {
	r := clitest.Run(Main,
		"examiner", "revoke",
		"--bundle-id", "bnd_3f29b71c",
		"--reason", "examination_concluded",
		"--correlation-id", "exam-2026-06-15-occ",
		"--tenant", "tenant_a",
		"--tenant", "tenant_b",
	).MustSucceed(t)
	r.StderrContains(t, "revoked bnd_3f29b71c")

	dec := json.NewDecoder(strings.NewReader(r.Stdout))
	var events []examiner.OperationalEvent
	for {
		var e examiner.OperationalEvent
		if err := dec.Decode(&e); err != nil {
			break
		}
		events = append(events, e)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 chain events, got %d\nstdout: %s", len(events), r.Stdout)
	}
	for i, e := range events {
		if e.Event != examiner.EventHandoverReturned {
			t.Errorf("event[%d].Event = %q, want %q", i, e.Event, examiner.EventHandoverReturned)
		}
		if got := e.Fields["handover_authorization_record_id"]; got != "bnd_3f29b71c" {
			t.Errorf("event[%d].Fields.handover_authorization_record_id = %v", i, got)
		}
	}
}

func TestExaminerRevoke_MissingFlags(t *testing.T) {
	clitest.Run(Main, "examiner", "revoke").
		MustExit(t, 2).
		StderrContains(t, "missing required flags").
		StderrContains(t, "--bundle-id")
}

// ----- examiner list / audit (stubbed at bootstrap) ----------------------

func TestExaminerList_StillStubbed(t *testing.T) {
	clitest.Run(Main, "examiner", "list").
		MustExit(t, 1).
		StderrContains(t, "not yet implemented")
}

func TestExaminerAudit_StillStubbed(t *testing.T) {
	clitest.Run(Main, "examiner", "audit").
		MustExit(t, 1).
		StderrContains(t, "not yet implemented")
}

// ----- fixture-lapse guard ------------------------------------------------

// futureExpiry returns an --expires-at value comfortably in the future
// relative to the real wall clock. The CLI's issuer (examiner.Issue)
// defaults Now to time.Now, so a black-box CLI fixture that hardcodes an
// absolute expiry date silently lapses the day that date passes and turns
// green tests red with no code change. Compute expiry relative to now
// instead. (The examiner unit tests may hardcode expiry because they also
// inject a fixed Now — a frozen clock plus a fixed expiry never lapses.)
func futureExpiry(t *testing.T) string {
	t.Helper()
	return time.Now().Add(365 * 24 * time.Hour).UTC().Format(time.RFC3339)
}

// TestNoHardcodedExpiryLiteral_InCLIFixtures is the class-level regression
// for the lapsed-fixture bug: it fails if any --expires-at in this file is
// followed by a hardcoded absolute-date literal instead of futureExpiry(t).
// The "not-a-date" negative fixture is not date-shaped and does not trip;
// the relative-clock form is not a string literal and does not trip. The
// needle is split around the comma so this guard cannot match its own
// regex source.
func TestNoHardcodedExpiryLiteral_InCLIFixtures(t *testing.T) {
	src, err := os.ReadFile("main_test.go")
	if err != nil {
		t.Fatalf("read own source: %v", err)
	}
	re := regexp.MustCompile(`"--expires-at",` + `\s*"\d{4}-\d{2}-\d{2}`)
	if loc := re.FindIndex(src); loc != nil {
		line := 1 + strings.Count(string(src[:loc[0]]), "\n")
		t.Fatalf("main_test.go:%d passes a hardcoded date to --expires-at; "+
			"use futureExpiry(t) so the fixture cannot silently lapse", line)
	}
}

// ----- helpers ------------------------------------------------------------

func publicFromPKCS8File(t *testing.T, path string) ed25519.PublicKey {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		t.Fatal("no PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		t.Fatal("not Ed25519")
	}
	return priv.Public().(ed25519.PublicKey)
}
