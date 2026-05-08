package serverconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault_Validates(t *testing.T) {
	c := Default()
	if err := c.Validate(); err != nil {
		t.Fatalf("default config does not validate: %v", err)
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	c := Default()
	c.Logging.Level = "warn"
	c.Storage.Path = "/tmp/data"
	path := filepath.Join(t.TempDir(), "ledger.json")
	if err := Save(path, &c); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Logging.Level != "warn" {
		t.Errorf("Logging.Level lost in round trip: %q", got.Logging.Level)
	}
	if got.Storage.Path != "/tmp/data" {
		t.Errorf("Storage.Path lost in round trip: %q", got.Storage.Path)
	}
}

func TestValidate_RejectsBadInputs(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*Config)
	}{
		{"unknown version", func(c *Config) { c.Version = "0.9" }},
		{"no otlp endpoints", func(c *Config) { c.OTLP = OTLP{} }},
		{"missing storage driver", func(c *Config) { c.Storage = Storage{} }},
		{"directory driver without path", func(c *Config) { c.Storage = Storage{Driver: "directory"} }},
		{"postgres driver without dsn", func(c *Config) { c.Storage = Storage{Driver: "postgres"} }},
		{"unknown storage driver", func(c *Config) { c.Storage = Storage{Driver: "voodoo"} }},
		{"missing hsm driver", func(c *Config) { c.HSM = HSM{} }},
		{"pkcs11 missing module", func(c *Config) { c.HSM = HSM{Driver: "pkcs11"} }},
		{"cloud-kms missing key", func(c *Config) { c.HSM = HSM{Driver: "cloud-kms"} }},
		{"unknown hsm driver", func(c *Config) { c.HSM = HSM{Driver: "magic"} }},
		{"software hsm forbidden in production", func(c *Config) {
			c.HSM = HSM{Driver: "software", ProductionMode: true}
		}},
		{"missing log level", func(c *Config) { c.Logging = Logging{} }},
		{"unknown log level", func(c *Config) { c.Logging = Logging{Level: "spew"} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := Default()
			tc.mut(&c)
			if err := c.Validate(); err == nil {
				t.Errorf("Validate accepted bad input")
			}
		})
	}
}

func TestValidate_AcceptsKnownLogLevels(t *testing.T) {
	for _, lvl := range []string{"debug", "info", "warn", "error"} {
		t.Run(lvl, func(t *testing.T) {
			c := Default()
			c.Logging.Level = lvl
			if err := c.Validate(); err != nil {
				t.Errorf("level %q rejected: %v", lvl, err)
			}
		})
	}
}

func TestLoad_RejectsMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "absent.json")); err == nil {
		t.Fatal("expected error on missing file")
	}
}

func TestLoad_RejectsMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "junk.json")
	if err := os.WriteFile(path, []byte("{not valid"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected error on malformed JSON")
	}
}
