package timeparse

import (
	"fmt"
	"time"
)

// Common log timestamp formats to try when parsing.
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999Z07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// ParseTimestamp attempts to parse a timestamp string using a set of
// known log formats. Returns the parsed time and the matched format.
func ParseTimestamp(s string) (time.Time, string, error) {
	for _, fmt := range knownFormats {
		t, err := time.Parse(fmt, s)
		if err == nil {
			return t, fmt, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("timeparse: unrecognized timestamp format: %q", s)
}

// ParseRange parses a start and end timestamp string, returning both
// as time.Time values. Either may be empty to indicate an open bound.
func ParseRange(start, end string) (time.Time, time.Time, error) {
	var tStart, tEnd time.Time
	var err error

	if start != "" {
		tStart, _, err = ParseTimestamp(start)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start time: %w", err)
		}
	}

	if end != "" {
		tEnd, _, err = ParseTimestamp(end)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end time: %w", err)
		}
	}

	if !tStart.IsZero() && !tEnd.IsZero() && tEnd.Before(tStart) {
		return time.Time{}, time.Time{}, fmt.Errorf("end time %q is before start time %q", end, start)
	}

	return tStart, tEnd, nil
}
