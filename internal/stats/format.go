package stats

import (
	"fmt"
	"io"
	"strings"
)

// Format writes a human-readable summary of r to w.
func Format(w io.Writer, r Result) {
	var sb strings.Builder
	sb.WriteString("\n--- logslice stats ---\n")
	fmt.Fprintf(&sb, "  total lines   : %d\n", r.TotalLines)
	fmt.Fprintf(&sb, "  matched lines : %d\n", r.MatchedLines)
	fmt.Fprintf(&sb, "  skipped lines : %d\n", r.SkippedLines)
	fmt.Fprintf(&sb, "  filtered out  : %d\n", r.FilteredLines)

	if r.EarliestTime != nil {
		fmt.Fprintf(&sb, "  earliest match: %s\n", r.EarliestTime.Format("2006-01-02 15:04:05"))
	}
	if r.LatestTime != nil {
		fmt.Fprintf(&sb, "  latest match  : %s\n", r.LatestTime.Format("2006-01-02 15:04:05"))
	}

	fmt.Fprintf(&sb, "  elapsed       : %s\n", r.Duration.Round(1000000))
	sb.WriteString("----------------------\n")
	fmt.Fprint(w, sb.String())
}
