package cli

import (
	"errors"
	"flag"
	"io"
)

// Flags holds all parsed command-line options for logslice.
type Flags struct {
	File        string
	Range       string
	Include     string
	Exclude     string
	Strict      bool
	LineNumbers bool
	Summary     bool
	Tail        int
	EveryN      int
	MaxLines    int
	Highlight   string
	// Dedup options
	DedupConsecutive bool
	DedupGlobal      bool
}

// ParseFlags parses os.Args using the provided io.Writer for usage output.
// Returns an error if required flags are missing.
func ParseFlags(args []string, out io.Writer) (*Flags, error) {
	f := &Flags{}
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(out)

	fs.StringVar(&f.File, "file", "", "path to log file (required)")
	fs.StringVar(&f.Range, "range", "", "time range e.g. 2024-01-01T10:00,2024-01-01T11:00 (required)")
	fs.StringVar(&f.Include, "include", "", "only include lines containing this keyword")
	fs.StringVar(&f.Exclude, "exclude", "", "exclude lines containing this keyword")
	fs.BoolVar(&f.Strict, "strict", false, "drop lines without a recognisable timestamp")
	fs.BoolVar(&f.LineNumbers, "line-numbers", false, "prefix output lines with line numbers")
	fs.BoolVar(&f.SummaryFlag(), "summary", false, "print summary after output")
	fs.IntVar(&f.Tail, "tail", 0, "return only the last N matching lines")
	fs.IntVar(&f.EveryN, "every", 1, "sample every Nth line")
	fs.IntVar(&f.MaxLines, "max", 0, "cap output at N lines")
	fs.StringVar(&f.Highlight, "highlight", "", "comma-separated terms to highlight in output")
	fs.BoolVar(&f.DedupConsecutive, "dedup-consecutive", false, "remove consecutive duplicate lines")
	fs.BoolVar(&f.DedupGlobal, "dedup-global", false, "remove all duplicate lines globally")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if f.File == "" {
		return nil, errors.New("--file is required")
	}
	if f.Range == "" {
		return nil, errors.New("--range is required")
	}
	return f, nil
}

// SummaryFlag returns a pointer to the Summary field for flag binding.
func (f *Flags) SummaryFlag() *bool { return &f.Summary }
