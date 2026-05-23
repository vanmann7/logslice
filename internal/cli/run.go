package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/input"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/slicer"
	"github.com/user/logslice/internal/timeparse"
)

// Run executes the logslice pipeline using the provided Options.
// It reads the log file, slices by time range, applies filters, and writes output.
func Run(opts *Options, stderr io.Writer) error {
	start, end, err := timeparse.ParseRange(opts.Start, opts.End)
	if err != nil {
		return fmt.Errorf("invalid time range: %w", err)
	}

	lines, err := input.ReadLines(opts.FilePath)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	var result []string

	if len(opts.Include) > 0 || len(opts.Exclude) > 0 {
		fo := filter.Options{
			MustContain: opts.Include,
			MustExclude: opts.Exclude,
		}
		result = slicer.SliceFiltered(lines, start, end, opts.Strict, fo)
	} else {
		result = slicer.Slice(lines, start, end, opts.Strict)
	}

	var w io.Writer = os.Stdout
	if opts.Output != "" {
		f, err := os.Create(opts.Output)
		if err != nil {
			return fmt.Errorf("opening output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	wopts := output.WriteOptions{
		LineNumbers: opts.LineNumbers,
		Summary:     opts.Summary,
	}

	if err := output.WriteLines(w, result, wopts); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
