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
				_, werr := writer.WriteString(line + "\n")
				if werr != nil {
					return werr
				}
			}
			continue
		}

		if (ts.Equal(opts.From) || ts.After(opts.From)) &&
			(ts.Equal(opts.To) || ts.Before(opts.To)) {
			if filter.Filter(line, opts.FilterOpts) {
				_, werr := writer.WriteString(line + "\n")
				if werr != nil {
					return werr
				}
			}
		}
	}

	return scanner.Err()
}
