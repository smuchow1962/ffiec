package main

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mmpworks/ffiec/ledger/internal/examiner"
)

func runExaminer(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		return usagef("examiner: missing subcommand (issue|revoke|list|audit)")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "issue":
		return runExaminerIssue(rest, stdout, stderr)
	case "revoke":
		return runExaminerRevoke(rest, stdout, stderr)
	case "list":
		return runExaminerList(rest, stdout, stderr)
	case "audit":
		return runExaminerAudit(rest, stdout, stderr)
	default:
		return usagef("examiner: unknown subcommand %q", sub)
	}
}

type stringSlice []string

func (s *stringSlice) String() string     { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error { *s = append(*s, v); return nil }

func runExaminerIssue(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("examiner issue", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var (
		institutionID   = fs.String("institution-id", "", "Institution identifier (required)")
		identity        = fs.String("identity", "", "Examiner identity claim (required)")
		laptopFP        = fs.String("laptop-fingerprint", "", "Optional examiner-laptop fingerprint pin")
		tenants         stringSlice
		periodStart     = fs.String("period-start", "", "Examination period start, YYYY-MM-DD (required)")
		periodEnd       = fs.String("period-end", "", "Examination period end, YYYY-MM-DD (required)")
		expiresAt       = fs.String("expires-at", "", "Bundle expiry, RFC3339 (required)")
		correlationID   = fs.String("correlation-id", "", "Examination correlation ID (required)")
		scope           = fs.String("scope", "compliance.read", "Scope: compliance.read|master-key.read")
		issuanceKeyPath = fs.String("issuance-key", "", "Path to Ed25519 issuance private key, PEM (required)")
		complianceURL   = fs.String("compliance-url", "", "Compliance read-surface base URL (required)")
		tlsFP           = fs.String("tls-cert-fingerprint", "", "Optional Compliance TLS cert fingerprint pin")
		credCertPath    = fs.String("cert-pem", "", "Optional examiner mTLS cert PEM to embed")
		credKeyPath     = fs.String("encrypted-key-pem", "", "Optional examiner mTLS encrypted-key PEM to embed")
		outPath         = fs.String("out", "", "Output bundle path (required)")
	)
	fs.Var(&tenants, "tenant", "Tenant in scope (repeatable, at least one required)")
	if stop, err := parseFlags(fs, args, stdout); err != nil {
		return err
	} else if stop {
		return nil
	}
	if err := requireFlags(map[string]string{
		"institution-id": *institutionID,
		"identity":       *identity,
		"period-start":   *periodStart,
		"period-end":     *periodEnd,
		"expires-at":     *expiresAt,
		"correlation-id": *correlationID,
		"issuance-key":   *issuanceKeyPath,
		"compliance-url": *complianceURL,
		"out":            *outPath,
	}); err != nil {
		return err
	}
	if len(tenants) == 0 {
		return usagef("--tenant: at least one is required")
	}

	expires, err := time.Parse(time.RFC3339, *expiresAt)
	if err != nil {
		return usagef("--expires-at: %v", err)
	}

	priv, pub, err := loadEd25519PrivatePEM(*issuanceKeyPath)
	if err != nil {
		return fmt.Errorf("--issuance-key: %w", err)
	}

	cred := examiner.BundleCredential{Type: "mtls_client_cert"}
	if *credCertPath != "" {
		raw, err := os.ReadFile(*credCertPath)
		if err != nil {
			return fmt.Errorf("--cert-pem: %w", err)
		}
		cred.CertPEM = string(raw)
	}
	if *credKeyPath != "" {
		raw, err := os.ReadFile(*credKeyPath)
		if err != nil {
			return fmt.Errorf("--encrypted-key-pem: %w", err)
		}
		cred.EncryptedKeyPEM = string(raw)
	}

	res, err := examiner.Issue(context.Background(), examiner.IssueRequest{
		InstitutionID:      *institutionID,
		IssuancePrivateKey: priv,
		IssuancePublicKey:  pub,
		ExaminerIdentity:   *identity,
		LaptopFingerprint:  *laptopFP,
		Tenants:            []string(tenants),
		ExaminationStart:   *periodStart,
		ExaminationEnd:     *periodEnd,
		ExpiresAt:          expires,
		CorrelationID:      *correlationID,
		Scope:              examiner.Scope(*scope),
		ComplianceURL:      *complianceURL,
		TLSCertFingerprint: *tlsFP,
		Credential:         cred,
	}, WriterChainWriter{W: stdout})
	if err != nil {
		return err
	}
	if err := os.WriteFile(*outPath, res.SignedJSON, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", *outPath, err)
	}
	fmt.Fprintf(stderr, "issued %s -> %s (expires %s)\n",
		res.Bundle.BundleID, *outPath, res.Bundle.ExpiresAt.Format(time.RFC3339))
	return nil
}

func runExaminerRevoke(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("examiner revoke", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var (
		bundleID      = fs.String("bundle-id", "", "Bundle id to revoke (required)")
		reason        = fs.String("reason", "", "Revocation reason (required)")
		correlationID = fs.String("correlation-id", "", "Examination correlation ID (required)")
		tenants       stringSlice
	)
	fs.Var(&tenants, "tenant", "Tenant in scope (repeatable, at least one required)")
	if stop, err := parseFlags(fs, args, stdout); err != nil {
		return err
	} else if stop {
		return nil
	}
	if err := requireFlags(map[string]string{
		"bundle-id":      *bundleID,
		"reason":         *reason,
		"correlation-id": *correlationID,
	}); err != nil {
		return err
	}
	if len(tenants) == 0 {
		return usagef("--tenant: at least one is required")
	}
	if err := examiner.Revoke(context.Background(), examiner.RevokeRequest{
		BundleID:      *bundleID,
		Reason:        *reason,
		CorrelationID: *correlationID,
		Tenants:       []string(tenants),
	}, WriterChainWriter{W: stdout}); err != nil {
		return err
	}
	fmt.Fprintf(stderr, "revoked %s\n", *bundleID)
	return nil
}

func runExaminerList(args []string, stdout, _ io.Writer) error {
	fs := flag.NewFlagSet("examiner list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Bool("active", false, "Show only active credentials")
	if stop, err := parseFlags(fs, args, stdout); err != nil {
		return err
	} else if stop {
		return nil
	}
	return errors.New("examiner list: not yet implemented (depends on Compliance read surface)")
}

func runExaminerAudit(args []string, stdout, _ io.Writer) error {
	fs := flag.NewFlagSet("examiner audit", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.String("correlation-id", "", "Filter by examination correlation ID")
	if stop, err := parseFlags(fs, args, stdout); err != nil {
		return err
	} else if stop {
		return nil
	}
	return errors.New("examiner audit: not yet implemented (depends on chain reader)")
}

// parseFlags wraps fs.Parse with two pieces of CLI ergonomics:
//
//   - --help prints usage to stdout and signals stop=true so the caller
//     returns nil cleanly (exit 0).
//   - parse errors (unknown flag, type mismatch) become usage errors so
//     Main maps them to exit 2.
//
// Returning stop separately keeps the leaf functions free of typed-error
// inspection — they just check stop and bail.
func parseFlags(fs *flag.FlagSet, args []string, stdout io.Writer) (stop bool, err error) {
	err = fs.Parse(args)
	if err == nil {
		return false, nil
	}
	if errors.Is(err, flag.ErrHelp) {
		fs.SetOutput(stdout)
		fmt.Fprintf(stdout, "Usage of ledgerctl %s:\n", fs.Name())
		fs.PrintDefaults()
		return true, nil
	}
	return false, usagef("%s: %v", fs.Name(), err)
}

func requireFlags(m map[string]string) error {
	var missing []string
	for k, v := range m {
		if v == "" {
			missing = append(missing, "--"+k)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return usagef("missing required flags: %s", strings.Join(missing, ", "))
	}
	return nil
}

// loadEd25519PrivatePEM accepts PKCS#8 ("PRIVATE KEY") — the format produced
// by `ledgerctl key generate` and by `openssl genpkey -algorithm ed25519`.
func loadEd25519PrivatePEM(path string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, nil, errors.New("no PEM block found")
	}
	if block.Type != "PRIVATE KEY" {
		return nil, nil, fmt.Errorf("unsupported PEM type %q (expected PRIVATE KEY / PKCS#8)", block.Type)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse PKCS#8: %w", err)
	}
	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("PKCS#8 key is not Ed25519")
	}
	pub, ok := priv.Public().(ed25519.PublicKey)
	if !ok {
		return nil, nil, errors.New("could not derive Ed25519 public key")
	}
	return priv, pub, nil
}
