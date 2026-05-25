// Package truncate provides utilities for truncating long log lines
// to a maximum byte length, preserving a readable suffix indicator.
package truncate

import "fmt"

const defaultMaxBytes = 512

// Options controls truncation behaviour.
type Options struct {
	// MaxBytes is the maximum number of bytes allowed per line.
	// Lines longer than this will be truncated. Zero means use the default (512).
	MaxBytes int
	// Indicator is the string appended to truncated lines to signal truncation.
	// Defaults to " [truncated]" if empty.
	Indicator string
}

func (o Options) maxBytes() int {
	if o.MaxBytes <= 0 {
		return defaultMaxBytes
	}
	return o.MaxBytes
}

func (o Options) indicator() string {
	if o.Indicator == "" {
		return " [truncated]"
	}
	return o.Indicator
}

// Line truncates a single line if it exceeds Options.MaxBytes.
// The returned string is guaranteed to be at most MaxBytes bytes long
// (including the indicator suffix).
func Line(line string, opts Options) string {
	max := opts.maxBytes()
	if len(line) <= max {
		return line
	}
	indicator := opts.indicator()
	cutAt := max - len(indicator)
	if cutAt < 0 {
		cutAt = 0
	}
	// Walk back to a valid UTF-8 boundary.
	for cutAt > 0 && line[cutAt]&0xC0 == 0x80 {
		cutAt--
	}
	return fmt.Sprintf("%s%s", line[:cutAt], indicator)
}

// Lines applies Line to every element of lines, returning a new slice.
func Lines(lines []string, opts Options) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = Line(l, opts)
	}
	return out
}
