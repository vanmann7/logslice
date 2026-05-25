// Package grep provides regex-based line matching for log files.
package grep

import (
	"fmt"
	"regexp"
)

// Options configures the grep matcher.
type Options struct {
	Pattern       string
	CaseSensitive bool
	Invert        bool
}

// Matcher holds a compiled regex and matching options.
type Matcher struct {
	re     *regexp.Regexp
	invert bool
}

// New compiles the pattern and returns a Matcher.
// If CaseSensitive is false, the pattern is wrapped with (?i).
func New(opts Options) (*Matcher, error) {
	if opts.Pattern == "" {
		return nil, fmt.Errorf("grep: pattern must not be empty")
	}
	pattern := opts.Pattern
	if !opts.CaseSensitive {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("grep: invalid pattern %q: %w", opts.Pattern, err)
	}
	return &Matcher{re: re, invert: opts.Invert}, nil
}

// Match reports whether line satisfies the matcher.
func (m *Matcher) Match(line string) bool {
	matched := m.re.MatchString(line)
	if m.invert {
		return !matched
	}
	return matched
}

// Apply filters lines, returning only those that satisfy the matcher.
func (m *Matcher) Apply(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if m.Match(l) {
			out = append(out, l)
		}
	}
	return out
}

// MustNew compiles the pattern and panics on error. Useful in tests.
func MustNew(opts Options) *Matcher {
	m, err := New(opts)
	if err != nil {
		panic(err)
	}
	return m
}
