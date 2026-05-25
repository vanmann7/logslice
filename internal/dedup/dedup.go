// Package dedup provides line deduplication for log output.
// It removes consecutive or global duplicate lines based on configuration.
package dedup

import "crypto/sha256"

// Options controls deduplication behaviour.
type Options struct {
	// Consecutive removes only adjacent duplicate lines.
	Consecutive bool
	// Global removes all duplicate lines across the entire input.
	Global bool
}

// Apply filters lines according to the given Options.
// If neither Consecutive nor Global is set, lines are returned unchanged.
func Apply(lines []string, opts Options) []string {
	if !opts.Consecutive && !opts.Global {
		return lines
	}
	if opts.Global {
		return globalDedup(lines)
	}
	return consecutiveDedup(lines)
}

func consecutiveDedup(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	out := make([]string, 0, len(lines))
	prev := ""
	for i, l := range lines {
		if i == 0 || l != prev {
			out = append(out, l)
		}
		prev = l
	}
	return out
}

func globalDedup(lines []string) []string {
	seen := make(map[[32]byte]struct{}, len(lines))
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		h := sha256.Sum256([]byte(l))
		if _, exists := seen[h]; !exists {
			seen[h] = struct{}{}
			out = append(out, l)
		}
	}
	return out
}
