package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ReadLines reads all lines from the given file path and returns them as a slice of strings.
// It returns an error if the file cannot be opened or read.
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("input: open file %q: %w", path, err)
	}
	defer f.Close()

	return readLines(f)
}

// ReadLinesFromReader reads all lines from an io.Reader and returns them as a slice of strings.
func ReadLinesFromReader(r io.Reader) ([]string, error) {
	return readLines(r)
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)

	// Support large lines (up to 1 MB)
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("input: scanning lines: %w", err)
	}

	return lines, nil
}
