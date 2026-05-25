// Package sampler provides deterministic line-sampling utilities for
// logslice output pipelines.
//
// When working with very large log extractions it is often useful to
// inspect a representative subset rather than every matched line.
// Sampler lets callers select every Nth line and optionally cap the
// total result count, making it easy to preview sliced output without
// overwhelming the terminal or downstream tooling.
//
// Basic usage:
//
//	out, err := sampler.Apply(lines, sampler.Options{
//		Every:    10,  // keep every 10th line
//		MaxLines: 100, // but no more than 100 total
//	})
package sampler
