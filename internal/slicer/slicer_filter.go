package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/timeparse"
)

// FilteredSliceOptions extends slicing with keyword filtering.
type FilteredSliceOptions struct {
	From       time.Time
	To         time.Time
	Strict     bool
	FilterOpts filter.Options
}

// SliceFiltered reads lines from r, retains those within [From, To],
// applies keyword filter rules, and writes matching lines to w.
// If Strict is true, lines without a parseable timestamp are skipped.
// If Strict is false, lines without a timestamp are included if they
// pass the keyword filter.
func SliceFiltered(r io.Reader, w io.Writer, opts FilteredSliceOptions) error {
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()

		ts, _, err := timeparse.ExtractFromLine(line)
		if err != nil {
			if opts.Strict {
				continue
			}
			// No timestamp — include if it passes keyword filter.
			if filter.Filter(line, opts.FilterOpts) {
				if werr := writeLine(writer, line); werr != nil {
					return werr
				}
			}
			continue
		}

		if inTimeRange(ts, opts.From, opts.To) && filter.Filter(line, opts.FilterOpts) {
			if werr := writeLine(writer, line); werr != nil {
				return werr
			}
		}
	}

	return scanner.Err()
}

// inTimeRange reports whether ts falls within the inclusive range [from, to].
func inTimeRange(ts, from, to time.Time) bool {
	return (ts.Equal(from) || ts.After(from)) &&
		(ts.Equal(to) || ts.Before(to))
}

// writeLine writes a single line followed by a newline to w.
func writeLine(w *bufio.Writer, line string) error {
	_, err := w.WriteString(line + "\n")
	return err
}
