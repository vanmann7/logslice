package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteLines_Basic(t *testing.T) {
	lines := []string{"line one", "line two", "line three"}
	var buf bytes.Buffer

	result, err := WriteLines(lines, Options{Dest: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LinesWritten != 3 {
		t.Errorf("expected 3 lines written, got %d", result.LinesWritten)
	}

	output := buf.String()
	for _, l := range lines {
		if !strings.Contains(output, l) {
			t.Errorf("expected output to contain %q", l)
		}
	}
}

func TestWriteLines_LineNumbers(t *testing.T) {
	lines := []string{"alpha", "beta"}
	var buf bytes.Buffer

	_, err := WriteLines(lines, Options{Dest: &buf, LineNumbers: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "1\talpha") {
		t.Errorf("expected line number prefix '1\\talpha', got: %s", output)
	}
	if !strings.Contains(output, "2\tbeta") {
		t.Errorf("expected line number prefix '2\\tbeta', got: %s", output)
	}
}

func TestWriteLines_Summary(t *testing.T) {
	lines := []string{"foo", "bar", "baz"}
	var buf bytes.Buffer

	_, err := WriteLines(lines, Options{Dest: &buf, Summary: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "3 line(s) matched") {
		t.Errorf("expected summary line, got: %s", output)
	}
}

func TestWriteLines_Empty(t *testing.T) {
	var buf bytes.Buffer

	result, err := WriteLines([]string{}, Options{Dest: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LinesWritten != 0 {
		t.Errorf("expected 0 lines written, got %d", result.LinesWritten)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %q", buf.String())
	}
}
