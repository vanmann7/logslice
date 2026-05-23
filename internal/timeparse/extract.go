package timeparse

import (
	"regexp"
	"strings"
)

// Common patterns used to locate a timestamp within a log line.
var timestampPatterns = []*regexp.Regexp{
	// ISO 8601 / RFC3339 with optional nanoseconds and timezone
	regexp.MustCompile(`\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?`),
	// Slash-separated date with time
	regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}`),
	// Common HTTP / Apache log format: 02/Jan/2006:15:04:05 -0700
	regexp.MustCompile(`\d{2}/[A-Za-z]{3}/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4}`),
	// Syslog: Jan 02 15:04:05
	regexp.MustCompile(`[A-Z][a-z]{2} +\d{1,2} \d{2}:\d{2}:\d{2}`),
}

// ExtractResult holds the result of extracting a timestamp from a log line.
type ExtractResult struct {
	Raw    string
	Offset int // byte offset of the match within the line
}

// ExtractFromLine scans a log line for the first recognisable timestamp
// substring and returns it together with its byte offset.
// Returns nil if no timestamp is found.
func ExtractFromLine(line string) *ExtractResult {
	for _, re := range timestampPatterns {
		loc := re.FindStringIndex(line)
		if loc == nil {
			continue
		}
		raw := strings.TrimSpace(line[loc[0]:loc[1]])
		return &ExtractResult{Raw: raw, Offset: loc[0]}
	}
	return nil
}
