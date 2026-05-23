// Package progress provides a simple progress reporter for tracking
// how many bytes have been processed from a log file.
package progress

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Reporter tracks byte-level progress through a file and can emit
// periodic updates to a writer.
type Reporter struct {
	total    int64
	read     atomic.Int64
	output   io.Writer
	enabled  bool
}

// New creates a new Reporter. If total is 0 or output is nil, the reporter
// operates in silent mode (no output is written).
func New(total int64, output io.Writer) *Reporter {
	enabled := total > 0 && output != nil
	return &Reporter{
		total:   total,
		output:  output,
		enabled: enabled,
	}
}

// Add records that n additional bytes have been processed.
func (r *Reporter) Add(n int64) {
	r.read.Add(n)
}

// Percent returns the current completion percentage (0–100).
// Returns 0 if total is unknown.
func (r *Reporter) Percent() float64 {
	if r.total <= 0 {
		return 0
	}
	v := r.read.Load()
	pct := float64(v) / float64(r.total) * 100
	if pct > 100 {
		pct = 100
	}
	return pct
}

// Print writes a single progress line to the configured output.
// It is a no-op when the reporter is in silent mode.
func (r *Reporter) Print() {
	if !r.enabled {
		return
	}
	v := r.read.Load()
	fmt.Fprintf(r.output, "progress: %s / %s (%.1f%%)\n",
		formatBytes(v), formatBytes(r.total), r.Percent())
}

// Done writes a final completion message.
func (r *Reporter) Done() {
	if !r.enabled {
		return
	}
	fmt.Fprintf(r.output, "progress: done — %s processed\n", formatBytes(r.read.Load()))
}

func formatBytes(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.2f GiB", float64(b)/float64(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.2f MiB", float64(b)/float64(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.2f KiB", float64(b)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
