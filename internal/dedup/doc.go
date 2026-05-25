// Package dedup implements line-level deduplication for log slices.
//
// Two strategies are supported:
//
//   - Consecutive: removes adjacent duplicate lines, preserving non-adjacent
//     repetitions. Useful for collapsing repeated log bursts.
//
//   - Global: removes every line that has already appeared anywhere in the
//     input, keeping only the first occurrence. Useful for unique-line views.
//
// Example usage:
//
//	result := dedup.Apply(lines, dedup.Options{Global: true})
package dedup
