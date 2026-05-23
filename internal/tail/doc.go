// Package tail implements efficient tail-reading for large log files.
//
// It provides ReadLastN and ReadLastNFromReader which use a ring buffer
// to return the final N lines of a file or reader without loading the
// entire content into memory. This is useful for previewing the end of
// large log files before performing a full slice operation.
//
// Example usage:
//
//	lines, err := tail.ReadLastN("/var/log/app.log", 50)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, l := range lines {
//		fmt.Println(l)
//	}
package tail
