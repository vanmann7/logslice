// Package sampler provides line sampling for large log files,
// allowing every Nth line to be selected from a slice result.
package sampler

import "fmt"

// Options configures the sampling behaviour.
type Options struct {
	// Every selects every Nth matching line (1 = all lines, 2 = every other, etc.).
	Every int
	// MaxLines caps the total number of lines returned (0 = unlimited).
	MaxLines int
}

// Validate returns an error if the options are invalid.
func (o Options) Validate() error {
	if o.Every < 1 {
		return fmt.Errorf("sampler: Every must be >= 1, got %d", o.Every)
	}
	if o.MaxLines < 0 {
		return fmt.Errorf("sampler: MaxLines must be >= 0, got %d", o.MaxLines)
	}
	return nil
}

// Apply returns a sampled subset of lines according to opts.
// Lines are 0-indexed internally; the counter resets per call.
func Apply(lines []string, opts Options) ([]string, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	every := opts.Every
	if every == 0 {
		every = 1
	}

	var result []string
	for i, line := range lines {
		if i%every != 0 {
			continue
		}
		result = append(result, line)
		if opts.MaxLines > 0 && len(result) >= opts.MaxLines {
			break
		}
	}
	return result, nil
}

// MustApply is like Apply but panics on invalid options.
// Useful in tests or when options are validated upstream.
func MustApply(lines []string, opts Options) []string {
	out, err := Apply(lines, opts)
	if err != nil {
		panic(err)
	}
	return out
}
