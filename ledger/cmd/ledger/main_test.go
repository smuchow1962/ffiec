package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmpworks/ffiec/ledger/internal/serverconfig"
	"github.com/mmpworks/ffiec/testkit/clitest"
)

// makeConfigFile is the workhorse for every CLI test that needs a valid
// config to point at. Returns the path to a freshly-written default config.
func makeConfigFile(t *testing.T) string {
	t.Helper()
	path := clitest.TempPath(t, "ledger.json")
	clitest.Run(Main, "config", "init", "--out", path).MustSucceed(t)
	return path
}

// ----- top-level dispatch -------------------------------------------------

func TestMain_NoArgsPrintsUsage(t *testing.T) {
	clitest.Run(Main).
		MustExit(t, 2).
		StderrContains(t, "ledger —")
}

func TestMain_HelpPrintsToStdout(t *testing.T) {
	clitest.Run(Main, "--help").
		MustSucceed(t).
		StdoutContains(t, "ledger —")
}

func TestMain_UnknownCommand(t *testing.T) {
	clitest.Run(Main, "frobnicate").
		MustExit(t, 2).
		StderrContains(t, `unknown command "frobnicate"`)
}

func TestMain_ConfigWithoutSubcommand(t *testing.T) {
	clitest.Run(Main, "config").
		MustExit(t, 2).
		StderrContains(t, "missing subcommand")
}

// ----- config init / validate ---------------------------------------------

func TestConfigInit_WritesValidDefault(t *testing.T) {
	path := clitest.TempPath(t, "ledger.json")
	clitest.Run(Main, "config", "init", "--out", path).
		MustSucceed(t).
		StderrContains(t, "wrote default config")

	raw := clitest.ReadFile(t, path)
	var c serverconfig.Config
	if err := json.Unmarshal(raw, &c); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}
	if err := c.Validate(); err != nil {
		t.Errorf("default config does not validate: %v", err)
	}
	if c.Version != serverconfig.Version {
		t.Errorf("version = %q, want %q", c.Version, serverconfig.Version)
	}
}

func TestConfigInit_RefusesToOverwriteByDefault(t *testing.T) {
	path := makeConfigFile(t)
	// Second run with no --overwrite should fail.
	clitest.Run(Main, "config", "init", "--out", path).
		MustExit(t, 1).
		StderrContains(t, "already exists")
}

func TestConfigInit_OverwriteAllowed(t *testing.T) {
	path := makeConfigFile(t)
	clitest.Run(Main, "config", "init", "--out", path, "--overwrite").
		MustSucceed(t)
}

func TestConfigInit_RequiresOut(t *testing.T) {
	clitest.Run(Main, "config", "init").
		MustExit(t, 2).
		StderrContains(t, "--out")
}

func TestConfigValidate_AcceptsDefault(t *testing.T) {
	path := makeConfigFile(t)
	clitest.Run(Main, "config", "validate", "--config", path).
		MustSucceed(t).
		StdoutContains(t, "config valid")
}

func TestConfigValidate_RejectsBrokenConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(path, []byte(`{"version":"0.0"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	clitest.Run(Main, "config", "validate", "--config", path).
		MustExit(t, 1).
		StderrContains(t, "version")
}

func TestConfigValidate_RequiresConfig(t *testing.T) {
	clitest.Run(Main, "config", "validate").
		MustExit(t, 2).
		StderrContains(t, "--config")
}

// ----- serve --------------------------------------------------------------

func TestServe_DryRunSucceeds(t *testing.T) {
	path := makeConfigFile(t)
	r := clitest.Run(Main, "serve", "--config", path, "--dry-run").MustSucceed(t)
	r.StdoutContains(t, "ledger readiness:")
	r.StdoutContains(t, "otlp_grpc:")
	r.StdoutContains(t, "hsm:")
}

func TestServe_RealRunStillStubbed(t *testing.T) {
	path := makeConfigFile(t)
	clitest.Run(Main, "serve", "--config", path).
		MustExit(t, 1).
		StderrContains(t, "not yet implemented")
}

func TestServe_RequiresConfig(t *testing.T) {
	clitest.Run(Main, "serve").
		MustExit(t, 2).
		StderrContains(t, "--config")
}

func TestServe_RedactsPostgresDSN(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ledger.json")
	c := serverconfig.Default()
	c.Storage = serverconfig.Storage{Driver: "postgres", DSN: "postgres://user:secret@db/ledger"}
	if err := serverconfig.Save(path, &c); err != nil {
		t.Fatal(err)
	}
	r := clitest.Run(Main, "serve", "--config", path, "--dry-run").MustSucceed(t)
	r.StdoutContains(t, `dsn="<set>"`)
	r.StdoutLacks(t, "secret")
}

// ----- health -------------------------------------------------------------

func TestHealth_AcceptsValidConfig(t *testing.T) {
	path := makeConfigFile(t)
	clitest.Run(Main, "health", "--config", path).
		MustSucceed(t).
		StdoutContains(t, "ok")
}

func TestHealth_RejectsBrokenConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(path, []byte(`{"version":"0.0"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	clitest.Run(Main, "health", "--config", path).
		MustExit(t, 1)
}

func TestHealth_RequiresConfig(t *testing.T) {
	clitest.Run(Main, "health").
		MustExit(t, 2).
		StderrContains(t, "--config")
}
