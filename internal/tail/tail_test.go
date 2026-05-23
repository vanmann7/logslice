package tail_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/tail"
)

func TestReadLastNFromReader_Basic(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5\n"
	got, err := tail.ReadLastNFromReader(strings.NewReader(input), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"line3", "line4", "line5"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, got[i])
		}
	}
}

func TestReadLastNFromReader_FewerLinesThanN(t *testing.T) {
	input := "only\ntwo\n"
	got, err := tail.ReadLastNFromReader(strings.NewReader(input), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	if got[0] != "only" || got[1] != "two" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestReadLastNFromReader_Empty(t *testing.T) {
	got, err := tail.ReadLastNFromReader(strings.NewReader(""), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestReadLastNFromReader_ZeroN(t *testing.T) {
	got, err := tail.ReadLastNFromReader(strings.NewReader("a\nb\n"), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result for n=0, got %v", got)
	}
}

func TestReadLastN_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	content := "alpha\nbeta\ngamma\ndelta\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	got, err := tail.ReadLastN(path, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "gamma" || got[1] != "delta" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestReadLastN_FileNotFound(t *testing.T) {
	_, err := tail.ReadLastN("/nonexistent/path/to/file.log", 5)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
