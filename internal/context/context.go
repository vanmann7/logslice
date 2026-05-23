// Package context provides a unified execution context for a logslice run,
// bundling parsed flags, stats collection, and progress reporting together.
package context

import (
	"io"
	"os"
	"time"

	"github.com/user/logslice/internal/cli"
	"github.com/user/logslice/internal/progress"
	"github.com/user/logslice/internal/stats"
)

// RunContext holds all runtime state for a single logslice execution.
type RunContext struct {
	Flags    cli.Flags
	Stats    *stats.Collector
	Progress *progress.Reporter
	Started  time.Time
	Output   io.Writer
}

// New creates a RunContext from the given flags, using os.Stdout as the
// default output writer and wiring up stats collection and progress reporting.
func New(flags cli.Flags, totalBytes int64) *RunContext {
	return NewWithWriter(flags, totalBytes, os.Stdout)
}

// NewWithWriter creates a RunContext with an explicit output writer. Useful
// for testing or when output is redirected to a file.
func NewWithWriter(flags cli.Flags, totalBytes int64, w io.Writer) *RunContext {
	var progressOut io.Writer
	if !flags.Quiet {
		progressOut = os.Stderr
	}

	return &RunContext{
		Flags:    flags,
		Stats:    stats.NewCollector(),
		Progress: progress.New(progressOut, totalBytes),
		Started:  time.Now(),
		Output:   w,
	}
}

// Elapsed returns the duration since the RunContext was created.
func (rc *RunContext) Elapsed() time.Duration {
	return time.Since(rc.Started)
}

// Summary returns a formatted stats summary string for this run.
func (rc *RunContext) Summary() string {
	return stats.Format(rc.Stats.Result())
}
