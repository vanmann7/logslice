package filter

import (
	"strings"
)

// Options holds configuration for line filtering.
type Options struct {
	// MustContain requires all specified substrings to be present in a line.
	MustContain []string
	// MustExclude rejects lines containing any of the specified substrings.
	MustExclude []string
	// CaseSensitive controls whether matching is case-sensitive.
	CaseSensitive bool
}

// Filter applies inclusion and exclusion rules to a log line.
// Returns true if the line passes all filter criteria.
func Filter(line string, opts Options) bool {
	cmp := line
	if !opts.CaseSensitive {
		cmp = strings.ToLower(line)
	}

	for _, must := range opts.MustContain {
		needle := must
		if !opts.CaseSensitive {
			needle = strings.ToLower(must)
		}
		if !strings.Contains(cmp, needle) {
			return false
		}
	}

	for _, excl := range opts.MustExclude {
		needle := excl
		if !opts.CaseSensitive {
			needle = strings.ToLower(excl)
		}
		if strings.Contains(cmp, needle) {
			return false
		}
	}

	return true
}

// ApplyToLines filters a slice of lines according to opts and returns
// only those that pass all criteria.
func ApplyToLines(lines []string, opts Options) []string {
	result := make([]string, 0, len(lines))
	for _, l := range lines {
		if Filter(l, opts) {
			result = append(result, l)
		}
	}
	return result
}
