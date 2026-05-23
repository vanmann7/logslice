package cli

import (
	"flag"
	"fmt"
	"os"
)

// Options holds all parsed CLI flags for logslice.
type Options struct {
	FilePath    string
	Start       string
	End         string
	Include     []string
	Exclude     []string
	Strict      bool
	LineNumbers bool
	Summary     bool
	Output      string
}

type multiFlag []string

func (m *multiFlag) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

// ParseFlags parses command-line arguments and returns an Options struct.
func ParseFlags(args []string) (*Options, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	opts := &Options{}
	var include, exclude multiFlag

	fs.StringVar(&opts.Start, "start", "", "Start of time range (e.g. '2024-01-15 10:00:00')")
	fs.StringVar(&opts.End, "end", "", "End of time range (e.g. '2024-01-15 11:00:00')")
	fs.BoolVar(&opts.Strict, "strict", false, "Exclude lines without a parseable timestamp")
	fs.BoolVar(&opts.LineNumbers, "line-numbers", false, "Prefix output lines with line numbers")
	fs.BoolVar(&opts.Summary, "summary", false, "Print a summary after output")
	fs.StringVar(&opts.Output, "output", "", "Write output to file instead of stdout")
	fs.Var(&include, "include", "Keyword that must appear in a line (repeatable)")
	fs.Var(&exclude, "exclude", "Keyword that must not appear in a line (repeatable)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() < 1 {
		return nil, fmt.Errorf("usage: logslice [flags] <logfile>")
	}

	opts.FilePath = fs.Arg(0)
	opts.Include = []string(include)
	opts.Exclude = []string(exclude)

	if opts.Start == "" && opts.End == "" {
		return nil, fmt.Errorf("at least one of -start or -end must be specified")
	}

	return opts, nil
}
