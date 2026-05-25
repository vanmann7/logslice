package pivot

import (
	"fmt"
	"strings"
)

// FormatOptions controls how a Result is rendered to text.
type FormatOptions struct {
	// ShowCounts emits a summary count line before each group.
	ShowCounts bool
	// Separator is printed between groups. Defaults to a blank line.
	Separator string
}

// Format renders a pivot Result as a human-readable string.
// Groups are emitted in insertion order.
func Format(res *Result, opts FormatOptions) string {
	if res == nil || len(res.Order) == 0 {
		return ""
	}

	sep := opts.Separator
	if sep == "" {
		sep = "\n"
	}

	var sb strings.Builder
	for i, key := range res.Order {
		lines := res.Groups[key]
		if opts.ShowCounts {
			sb.WriteString(fmt.Sprintf("=== %s (%d lines) ===\n", key, len(lines)))
		}
		for _, l := range lines {
			sb.WriteString(l)
			sb.WriteByte('\n')
		}
		if i < len(res.Order)-1 {
			sb.WriteString(sep)
		}
	}
	return sb.String()
}
