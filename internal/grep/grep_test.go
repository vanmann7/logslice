package grep_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/grep"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := grep.New(grep.Options{Pattern: "[", CaseSensitive: true})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := grep.New(grep.Options{Pattern: ""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	m := grep.MustNew(grep.Options{Pattern: "error", CaseSensitive: false})
	if !m.Match("2024-01-01 ERROR something failed") {
		t.Error("expected case-insensitive match")
	}
}

func TestMatch_CaseSensitive_NoMatch(t *testing.T) {
	m := grep.MustNew(grep.Options{Pattern: "error", CaseSensitive: true})
	if m.Match("2024-01-01 ERROR something failed") {
		t.Error("expected no match with case-sensitive pattern")
	}
}

func TestMatch_Invert(t *testing.T) {
	m := grep.MustNew(grep.Options{Pattern: "DEBUG", CaseSensitive: true, Invert: true})
	if !m.Match("INFO server started") {
		t.Error("expected inverted match for non-DEBUG line")
	}
	if m.Match("DEBUG verbose output") {
		t.Error("expected inverted match to exclude DEBUG line")
	}
}

func TestApply_FiltersLines(t *testing.T) {
	lines := []string{
		"INFO  startup complete",
		"ERROR disk full",
		"WARN  low memory",
		"DEBUG trace data",
		"ERROR network timeout",
	}
	m := grep.MustNew(grep.Options{Pattern: "ERROR", CaseSensitive: true})
	got := m.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	for _, l := range got {
		if !m.Match(l) {
			t.Errorf("unexpected line in result: %q", l)
		}
	}
}

func TestApply_InvertFiltersLines(t *testing.T) {
	lines := []string{
		"INFO  startup complete",
		"DEBUG trace data",
		"INFO  shutdown",
	}
	m := grep.MustNew(grep.Options{Pattern: "DEBUG", CaseSensitive: true, Invert: true})
	got := m.Apply(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	m := grep.MustNew(grep.Options{Pattern: "ERROR"})
	got := m.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}

func TestApply_RegexGroups(t *testing.T) {
	lines := []string{
		"req_id=abc123 status=200",
		"req_id=xyz789 status=500",
		"req_id=def456 status=404",
	}
	m := grep.MustNew(grep.Options{Pattern: `status=5\d{2}`, CaseSensitive: true})
	got := m.Apply(lines)
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got))
	}
	if got[0] != lines[1] {
		t.Errorf("unexpected match: %q", got[0])
	}
}
