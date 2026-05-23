package slicer

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/logslice/logslice/internal/timeparse"
)

// Options configures the slicing behavior.
type Options struct {
	Start     time.Time
	End       time.Time
	Formats   []string
	StrictMode bool // if true, lines without a timestamp are excluded
}

// Result holds statistics from a slice operation.
type Result struct {
	LinesRead    int
	LinesWritten int
	LinesSkipped int
}

// Slice reads log lines from r, writing lines whose timestamps fall within
// [opts.Start, opts.End] (inclusive) to w.
func Slice(r io.Reader, w io.Writer, opts Options) (Result, error) {
	var res Result

	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		res.LinesRead++

		t, _, err := timeparse.ExtractFromLine(line, opts.Formats)
		if err != nil {
			// No timestamp found in this line.
			if opts.StrictMode {
				res.LinesSkipped++
				continue
			}
			// In non-strict mode, pass through lines without timestamps.
			_, werr := fmt.Fprintln(writer, line)
			if werr != nil {
				return res, werr
			}
			res.LinesWritten++
			continue
		}

		if (t.Equal(opts.Start) || t.After(opts.Start)) &&
			(t.Equal(opts.End) || t.Before(opts.End)) {
			_, werr := fmt.Fprintln(writer, line)
			if werr != nil {
				return res, werr
			}
			res.LinesWritten++
		} else {
			res.LinesSkipped++
		}
	}

	if err := scanner.Err(); err != nil {
		return res, fmt.Errorf("slicer: scanner error: %w", err)
	}

	return res, nil
}
