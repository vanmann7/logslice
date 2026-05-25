// Package merge provides utilities for combining and interleaving
// multiple sorted log line slices into a single ordered output.
package merge

import (
	"sort"
	"time"

	"github.com/user/logslice/internal/timeparse"
)

// TimedLine associates a raw log line with its parsed timestamp.
type TimedLine struct {
	Line      string
	Timestamp *time.Time
}

// Annotate parses timestamps from each line and returns a slice of TimedLine.
// Lines that cannot be parsed retain a nil Timestamp.
func Annotate(lines []string) []TimedLine {
	result := make([]TimedLine, 0, len(lines))
	for _, l := range lines {
		tl := TimedLine{Line: l}
		if ts, _, err := timeparse.ExtractFromLine(l); err == nil {
			t := ts
			tl.Timestamp = &t
		}
		result = append(result, tl)
	}
	return result
}

// Merge combines multiple slices of log lines into a single slice ordered
// by timestamp. Lines without a parseable timestamp are appended at the end
// in their original relative order.
func Merge(sources ...[]string) []string {
	var timed []TimedLine
	var untimed []string

	for _, src := range sources {
		for _, tl := range Annotate(src) {
			if tl.Timestamp != nil {
				timed = append(timed, tl)
			} else {
				untimed = append(untimed, tl.Line)
			}
		}
	}

	sort.SliceStable(timed, func(i, j int) bool {
		return timed[i].Timestamp.Before(*timed[j].Timestamp)
	})

	result := make([]string, 0, len(timed)+len(untimed))
	for _, tl := range timed {
		result = append(result, tl.Line)
	}
	result = append(result, untimed...)
	return result
}
