// Package grep provides regex-based line matching for use in the logslice
// pipeline. It supports case-insensitive matching and inverted (exclude)
// semantics, complementing the keyword-based filter package with full
// regular expression power.
//
// Basic usage:
//
//	m, err := grep.New(grep.Options{
//		Pattern:       `ERROR|WARN`,
//		CaseSensitive: false,
//		Invert:        false,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	matched := m.Apply(lines)
package grep
