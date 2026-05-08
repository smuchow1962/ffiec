package cliutil

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
)

// StringSlice is a flag.Value that accumulates repeated string flags.
//
//	var tenants StringSlice
//	fs.Var(&tenants, "tenant", "Tenant in scope (repeatable)")
type StringSlice []string

// String returns the comma-joined slice for flag's String() contract.
func (s *StringSlice) String() string { return strings.Join(*s, ",") }

// Set appends v to the slice.
func (s *StringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

// ParseFlags wraps fs.Parse with two pieces of CLI ergonomics:
//
//   - --help prints usage to stdout and signals stop=true so the caller
//     returns nil cleanly (exit 0).
//   - parse errors (unknown flag, type mismatch) become *UsageError so
//     the dispatcher maps them to exit 2.
//
// Returning stop separately keeps leaf functions free of typed-error
// inspection — they just check stop and bail.
//
// The fs is expected to be created with flag.ContinueOnError and SetOutput
// to io.Discard; ParseFlags re-points fs.Output to stdout when --help fires.
func ParseFlags(fs *flag.FlagSet, args []string, stdout io.Writer, progName string) (stop bool, err error) {
	err = fs.Parse(args)
	if err == nil {
		return false, nil
	}
	if errors.Is(err, flag.ErrHelp) {
		fs.SetOutput(stdout)
		fmt.Fprintf(stdout, "Usage of %s %s:\n", progName, fs.Name())
		fs.PrintDefaults()
		return true, nil
	}
	return false, Usagef("%s: %v", fs.Name(), err)
}

// RequireFlags returns a *UsageError listing every key in m whose value is
// the empty string. The list is sorted for stable error messages.
func RequireFlags(m map[string]string) error {
	var missing []string
	for k, v := range m {
		if v == "" {
			missing = append(missing, "--"+k)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	sort.Strings(missing)
	return Usagef("missing required flags: %s", strings.Join(missing, ", "))
}
