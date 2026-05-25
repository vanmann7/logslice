// Package pivot provides field-based grouping of log lines.
// It extracts a named field from structured log lines and groups
// matching lines under a common key for summarisation or comparison.
package pivot

import (
	"regexp"
	"strings"
)

// Result holds grouped log lines keyed by the extracted field value.
type Result struct {
	Groups map[string][]string
	Order  []string // insertion-ordered keys
}

// Options controls how field extraction is performed.
type Options struct {
	// Pattern is a regexp with a named capture group "key".
	// Example: `level=(?P<key>\w+)`
	Pattern string
	// FallbackKey is used when a line does not match the pattern.
	FallbackKey string
}

// Apply groups lines by the value of the named capture "key" in opts.Pattern.
// Lines that do not match are grouped under opts.FallbackKey (default: "<unmatched>").
func Apply(lines []string, opts Options) (*Result, error) {
	fallback := opts.FallbackKey
	if fallback == "" {
		fallback = "<unmatched>"
	}

	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, err
	}

	names := re.SubexpNames()
	keyIdx := -1
	for i, n := range names {
		if n == "key" {
			keyIdx = i
			break
		}
	}

	res := &Result{Groups: make(map[string][]string)}

	for _, line := range lines {
		k := fallback
		if keyIdx >= 0 {
			m := re.FindStringSubmatch(line)
			if m != nil && keyIdx < len(m) {
				k = strings.TrimSpace(m[keyIdx])
			}
		}
		if _, exists := res.Groups[k]; !exists {
			res.Order = append(res.Order, k)
		}
		res.Groups[k] = append(res.Groups[k], line)
	}

	return res, nil
}

// MustApply is like Apply but panics on error.
func MustApply(lines []string, opts Options) *Result {
	res, err := Apply(lines, opts)
	if err != nil {
		panic(err)
	}
	return res
}
