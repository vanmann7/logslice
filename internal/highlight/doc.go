// Package highlight provides ANSI terminal highlighting for log lines.
//
// It is used by the CLI output layer to visually emphasize keyword matches
// when the --highlight flag is enabled. Highlighting is applied after
// filtering and slicing, immediately before writing to stdout.
//
// Example usage:
//
//	opts := highlight.Options{
//		Terms:         []string{"error", "warn"},
//		CaseSensitive: false,
//		Color:         highlight.AnsiRed,
//	}
//	highlighted := highlight.Lines(matchedLines, opts)
package highlight
