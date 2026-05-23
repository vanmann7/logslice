package filter_test

import (
	"testing"

	"github.com/logslice/logslice/internal/filter"
)

func TestFilter_MustContain(t *testing.T) {
	opts := filter.Options{MustContain: []string{"ERROR"}, CaseSensitive: true}

	if !filter.Filter("2024-01-01 ERROR something went wrong", opts) {
		t.Error("expected line with ERROR to pass")
	}
	if filter.Filter("2024-01-01 INFO all good", opts) {
		t.Error("expected line without ERROR to fail")
	}
}

func TestFilter_MustExclude(t *testing.T) {
	opts := filter.Options{MustExclude: []string{"DEBUG"}, CaseSensitive: true}

	if filter.Filter("2024-01-01 DEBUG verbose output", opts) {
		t.Error("expected DEBUG line to be excluded")
	}
	if !filter.Filter("2024-01-01 INFO important", opts) {
		t.Error("expected INFO line to pass")
	}
}

func TestFilter_CaseInsensitive(t *testing.T) {
	opts := filter.Options{MustContain: []string{"error"}, CaseSensitive: false}

	if !filter.Filter("2024-01-01 ERROR uppercase", opts) {
		t.Error("expected case-insensitive match to pass")
	}
	if !filter.Filter("2024-01-01 Error mixed", opts) {
		t.Error("expected mixed-case match to pass")
	}
}

func TestFilter_MultipleMustContain(t *testing.T) {
	opts := filter.Options{MustContain: []string{"ERROR", "database"}, CaseSensitive: true}

	if !filter.Filter("ERROR: database connection failed", opts) {
		t.Error("expected line with both terms to pass")
	}
	if filter.Filter("ERROR: network timeout", opts) {
		t.Error("expected line missing 'database' to fail")
	}
}

func TestFilter_EmptyOptions(t *testing.T) {
	opts := filter.Options{}
	if !filter.Filter("any line at all", opts) {
		t.Error("expected empty options to pass all lines")
	}
}

func TestApplyToLines(t *testing.T) {
	lines := []string{
		"INFO started",
		"ERROR something failed",
		"DEBUG verbose",
		"ERROR another failure",
	}
	opts := filter.Options{MustContain: []string{"ERROR"}, CaseSensitive: true}
	got := filter.ApplyToLines(lines, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
	for _, l := range got {
		if l != "ERROR something failed" && l != "ERROR another failure" {
			t.Errorf("unexpected line: %q", l)
		}
	}
}

func TestApplyToLines_Empty(t *testing.T) {
	got := filter.ApplyToLines(nil, filter.Options{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
