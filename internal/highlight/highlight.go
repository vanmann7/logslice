// Package highlight provides utilities for highlighting matched keywords
// within log lines using ANSI escape codes.
package highlight

import "strings"

const (
	AnsiReset  = "\033[0m"
	AnsiYellow = "\033[33m"
	AnsiRed    = "\033[31m"
	AnsiBold   = "\033[1m"
)

// Options controls how highlighting is applied.
type Options struct {
	Terms       []string
	Color       string
	CaseSensitive bool
}

// Line applies ANSI highlighting to all occurrences of each term in the
// given line. Returns the original line if no terms match.
func Line(line string, opts Options) string {
	if len(opts.Terms) == 0 {
		return line
	}

	color := opts.Color
	if color == "" {
		color = AnsiYellow
	}

	result := line
	for _, term := range opts.Terms {
		if term == "" {
			continue
		}
		result = replaceTerm(result, term, color, opts.CaseSensitive)
	}
	return result
}

// replaceTerm replaces all occurrences of term in s with a highlighted version.
func replaceTerm(s, term, color string, caseSensitive bool) string {
	if caseSensitive {
		return strings.ReplaceAll(s, term, color+AnsiBold+term+AnsiReset)
	}

	lower := strings.ToLower(s)
	lowerTerm := strings.ToLower(term)

	var result strings.Builder
	offset := 0
	for {
		idx := strings.Index(lower[offset:], lowerTerm)
		if idx == -1 {
			result.WriteString(s[offset:])
			break
		}
		abs := offset + idx
		result.WriteString(s[offset:abs])
		result.WriteString(color + AnsiBold + s[abs:abs+len(term)] + AnsiReset)
		offset = abs + len(term)
	}
	return result.String()
}

// Lines applies Line to each entry in a slice, returning a new slice.
func Lines(lines []string, opts Options) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = Line(l, opts)
	}
	return out
}
