package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Options controls how output is written.
type Options struct {
	// Dest is the output destination. If nil, os.Stdout is used.
	Dest io.Writer
	// LineNumbers prefixes each line with its original line number.
	LineNumbers bool
	// Summary prints a count of matched lines after writing.
	Summary bool
}

// Result holds metadata about a write operation.
type Result struct {
	LinesWritten int
}

// WriteLines writes the given lines to the destination described by opts.
// It returns a Result with metadata about what was written.
func WriteLines(lines []string, opts Options) (Result, error) {
	dest := opts.Dest
	if dest == nil {
		dest = os.Stdout
	}

	bw := bufio.NewWriter(dest)
	result := Result{}

	for i, line := range lines {
		var err error
		if opts.LineNumbers {
			_, err = fmt.Fprintf(bw, "%d\t%s\n", i+1, line)
		} else {
			_, err = fmt.Fprintln(bw, line)
		}
		if err != nil {
			return result, fmt.Errorf("output: write line %d: %w", i+1, err)
		}
		result.LinesWritten++
	}

	if opts.Summary {
		if _, err := fmt.Fprintf(bw, "--- %d line(s) matched ---\n", result.LinesWritten); err != nil {
			return result, fmt.Errorf("output: write summary: %w", err)
		}
	}

	if err := bw.Flush(); err != nil {
		return result, fmt.Errorf("output: flush: %w", err)
	}

	return result, nil
}
