// Package tail provides utilities for reading the last N lines of a log file
// efficiently without loading the entire file into memory.
package tail

import (
	"bufio"
	"io"
	"os"
)

// ReadLastN returns the last n lines from the named file.
// If the file has fewer than n lines, all lines are returned.
func ReadLastN(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadLastNFromReader(f, n)
}

// ReadLastNFromReader returns the last n lines from the given reader.
// It reads all lines into a ring buffer of size n.
func ReadLastNFromReader(r io.Reader, n int) ([]string, error) {
	if n <= 0 {
		return []string{}, nil
	}

	buf := make([]string, n)
	pos := 0
	count := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		buf[pos%n] = scanner.Text()
		pos++
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if count == 0 {
		return []string{}, nil
	}

	if count <= n {
		// All lines fit; return in order
		result := make([]string, count)
		copy(result, buf[:count])
		return result, nil
	}

	// Ring buffer wrapped; reconstruct in order
	result := make([]string, n)
	start := pos % n
	for i := 0; i < n; i++ {
		result[i] = buf[(start+i)%n]
	}
	return result, nil
}
