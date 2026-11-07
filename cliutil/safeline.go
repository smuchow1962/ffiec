package cliutil

import (
	"fmt"
	"strings"
	"unicode"
)

// SafeLine is the one gate for text that reaches a line of output. Values from a file under
// examination can reach a report, and many of them have no MAC and no signature. A line break
// there would forge a line of the report, and a control character could rewrite a terminal.
// The value stays readable: each such character prints as its code point.
func SafeLine(s string) string {
	var line strings.Builder
	for _, r := range s {
		if unicode.IsControl(r) || unicode.In(r, unicode.Cf, unicode.Zl, unicode.Zp) {
			fmt.Fprintf(&line, `\u%04X`, r)
			continue
		}
		line.WriteRune(r)
	}
	return line.String()
}
