// Package serverconfig defines the ledger server's configuration schema,
// load/save/validate helpers, and the default values used by `ledger
// config init`.
//
// The schema is intentionally narrow at bootstrap. New fields land here
// rather than in flags so a config file is the single source of truth for
// what the server is going to do at startup.
package serverconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Version is the current config schema version. Loaders that see a newer
// number than they support fail closed rather than silently misinterpret.
const Version = "1.0"

// Config is the full server configuration. Stored on disk as JSON so the
// shape is human-inspectable without bringing in a YAML parser.
type Config struct {
	Version string  `json:"version"`
	OTLP    OTLP    `json:"otlp"`
	Storage Storage `json:"storage"`
	HSM     HSM     `json:"hsm"`
	Logging Logging `json:"logging"`
}

// OTLP is the on-wire receiver settings.
type OTLP struct {
	GRPCAddr string `json:"grpc_addr"` // e.g. "0.0.0.0:4317"
	HTTPAddr string `json:"http_addr"` // e.g. "0.0.0.0:4318"
}

// Storage is where the append-only ledger lives.
type Storage struct {
	// Driver is "directory" (NDJSON files under Path) or "postgres" (DSN).
	Driver string `json:"driver"`
	Path   string `json:"path,omitempty"` // for "directory"
	DSN    string `json:"dsn,omitempty"`  // for "postgres"
}

// HSM is the daily-seal-signing custody.
type HSM struct {
	// Driver is "pkcs11" (production), "cloud-kms" (managed), or "software"
	// (development only — the bootstrap server refuses to use software keys
	// when ProductionMode is true).
	Driver         string `json:"driver"`
	PKCS11Module   string `json:"pkcs11_module,omitempty"`
	CloudKMSKey    string `json:"cloud_kms_key,omitempty"`
	SoftwareKeyPEM string `json:"software_key_pem,omitempty"`
	ProductionMode bool   `json:"production_mode"`
}

// Logging controls the server's own observability output.
type Logging struct {
	Level string `json:"level"` // "debug" | "info" | "warn" | "error"
}

// Default returns the bootstrap-default config. Used by `ledger config init`.
func Default() Config {
	return Config{
		Version: Version,
		OTLP: OTLP{
			GRPCAddr: "0.0.0.0:4317",
			HTTPAddr: "0.0.0.0:4318",
		},
		Storage: Storage{
			Driver: "directory",
			Path:   "/var/lib/ffiec-ledger",
		},
		HSM: HSM{
			Driver:         "software",
			ProductionMode: false,
		},
		Logging: Logging{Level: "info"},
	}
}

// Load reads a config file from path and validates it.
func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Save writes c to path as pretty-printed JSON with mode 0640.
func Save(path string, c *Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, raw, 0o640); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// Validate checks the structural and semantic invariants. The validation
// is intentionally strict — fail closed at startup, never half-configured.
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("nil config")
	}
	if c.Version != Version {
		return fmt.Errorf("version %q not supported (expected %q)", c.Version, Version)
	}
	if err := c.OTLP.validate(); err != nil {
		return err
	}
	if err := c.Storage.validate(); err != nil {
		return err
	}
	if err := c.HSM.validate(); err != nil {
		return err
	}
	if err := c.Logging.validate(); err != nil {
		return err
	}
	return nil
}

func (o *OTLP) validate() error {
	if o.GRPCAddr == "" && o.HTTPAddr == "" {
		return errors.New("otlp: at least one of grpc_addr or http_addr must be set")
	}
	return nil
}

func (s *Storage) validate() error {
	switch s.Driver {
	case "directory":
		if s.Path == "" {
			return errors.New("storage: directory driver requires path")
		}
	case "postgres":
		if s.DSN == "" {
			return errors.New("storage: postgres driver requires dsn")
		}
	case "":
		return errors.New("storage: driver is required")
	default:
		return fmt.Errorf("storage: unknown driver %q", s.Driver)
	}
	return nil
}

func (h *HSM) validate() error {
	switch h.Driver {
	case "pkcs11":
		if h.PKCS11Module == "" {
			return errors.New("hsm: pkcs11 driver requires pkcs11_module")
		}
	case "cloud-kms":
		if h.CloudKMSKey == "" {
			return errors.New("hsm: cloud-kms driver requires cloud_kms_key")
		}
	case "software":
		if h.ProductionMode {
			return errors.New("hsm: software driver is forbidden in production_mode (use pkcs11 or cloud-kms)")
		}
	case "":
		return errors.New("hsm: driver is required")
	default:
		return fmt.Errorf("hsm: unknown driver %q", h.Driver)
	}
	return nil
}

func (l *Logging) validate() error {
	switch l.Level {
	case "debug", "info", "warn", "error":
		return nil
	case "":
		return errors.New("logging: level is required")
	default:
		return fmt.Errorf("logging: unknown level %q", l.Level)
	}
}
