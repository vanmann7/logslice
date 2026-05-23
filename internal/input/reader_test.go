package input

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadLinesFromReader_Basic(t *testing.T) {
	input := "line one\nline two\nline three\n"
	r := strings.NewReader(input)

	lines, err := ReadLinesFromReader(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line one" {
		t.Errorf("expected %q, got %q", "line one", lines[0])
	}
	if lines[2] != "line three" {
		t.Errorf("expected %q, got %q", "line three", lines[2])
	}
}

func TestReadLinesFromReader_Empty(t *testing.T) {
	r := strings.NewReader("")

	lines, err := ReadLinesFromReader(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}

func TestReadLinesFromReader_NoTrailingNewline(t *testing.T) {
	r := strings.NewReader("alpha\nbeta")

	lines, err := ReadLinesFromReader(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[1] != "beta" {
		t.Errorf("expected %q, got %q", "beta", lines[1])
	}
}

func TestReadLines_File(t *testing.T) {
	content := "2024-01-01 INFO starting\n2024-01-01 DEBUG ready\n"
	tmp := filepath.Join(t.TempDir(), "test.log")

	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	lines, err := ReadLines(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "2024-01-01 INFO starting" {
		t.Errorf("unexpected first line: %q", lines[0])
	}
}

func TestReadLines_FileNotFound(t *testing.T) {
	_, err := ReadLines("/nonexistent/path/to/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
