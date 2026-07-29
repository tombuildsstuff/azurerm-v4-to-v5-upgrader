package engine

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pmezard/go-difflib/difflib"
)

// minimizeFormatChurn reverts hclwrite's incidental, whitespace-only
// re-formatting so that only lines the transforms *semantically* changed differ
// from the original.
//
// The tool must re-serialize every touched file through hclwrite to apply its
// edits, and hclwrite always canonicalizes the whole file on output (the same
// formatting `terraform fmt` applies) - e.g. it rewrites the JSON-style object
// syntax `"k": v` to `"k" : v`. That produces diff noise on lines the upgrade
// never touched. When callers opt in (Options.PreserveFormatting), this walks a
// line diff between the original and the formatted output and, for any line that
// differs *only in whitespace*, restores the original line. Lines with a real
// content change (a rename, a new/removed line) are kept as hclwrite produced
// them, so they stay canonically formatted.
//
// If the reconciled result fails to parse as HCL (a safety net against an
// unforeseen alignment edge case), the original formatted output is returned
// unchanged - correctness never depends on this cosmetic pass.
func minimizeFormatChurn(orig, formatted []byte, filename string) []byte {
	if string(orig) == string(formatted) {
		return formatted
	}
	origLines := splitKeepEnds(string(orig))
	newLines := splitKeepEnds(string(formatted))

	// Diff on whitespace-stripped lines: an "equal" run is a set of lines that
	// match modulo whitespace - i.e. hclwrite changed only their formatting, so
	// restore the original text. "replace"/"insert" runs carry real content
	// changes and keep hclwrite's canonical output; "delete" runs are dropped.
	m := difflib.NewMatcher(stripAll(origLines), stripAll(newLines))
	var b strings.Builder
	for _, op := range m.GetOpCodes() {
		switch op.Tag {
		case 'e': // equal modulo whitespace - keep the original formatting.
			b.WriteString(strings.Join(origLines[op.I1:op.I2], ""))
		case 'r', 'i': // real change (replace) or added lines (insert).
			b.WriteString(strings.Join(newLines[op.J1:op.J2], ""))
		case 'd': // lines removed by the transform - emit nothing.
		}
	}
	out := []byte(b.String())

	if _, diags := hclsyntax.ParseConfig(out, filename, hcl.InitialPos); diags.HasErrors() {
		return formatted
	}
	return out
}

// splitKeepEnds splits s into lines, keeping the trailing "\n" on each so the
// pieces re-join to exactly s (the final line has no newline unless s ended
// with one).
func splitKeepEnds(s string) []string {
	if s == "" {
		return nil
	}
	var lines []string
	for {
		i := strings.IndexByte(s, '\n')
		if i < 0 {
			lines = append(lines, s)
			return lines
		}
		lines = append(lines, s[:i+1])
		s = s[i+1:]
		if s == "" {
			return lines
		}
	}
}

// stripAll returns each line with all whitespace removed, for whitespace-
// insensitive line matching. HCL treats inter-token whitespace as
// insignificant, so two lines equal under stripAll differ only cosmetically.
func stripAll(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = strings.Map(func(r rune) rune {
			if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
				return -1
			}
			return r
		}, l)
	}
	return out
}
