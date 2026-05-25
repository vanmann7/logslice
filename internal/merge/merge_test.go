package merge_test

import (
	"testing"

	"github.com/user/logslice/internal/merge"
)

func TestMerge_SingleSource(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z info start",
		"2024-01-01T10:01:00Z info middle",
		"2024-01-01T10:02:00Z info end",
	}
	got := merge.Merge(lines)
	if len(got) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(got))
	}
	for i, l := range lines {
		if got[i] != l {
			t.Errorf("line %d: expected %q, got %q", i, l, got[i])
		}
	}
}

func TestMerge_TwoSourcesInterleaved(t *testing.T) {
	a := []string{
		"2024-01-01T10:00:00Z info A first",
		"2024-01-01T10:02:00Z info A second",
	}
	b := []string{
		"2024-01-01T10:01:00Z info B first",
		"2024-01-01T10:03:00Z info B second",
	}
	got := merge.Merge(a, b)
	if len(got) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(got))
	}
	expected := []string{
		"2024-01-01T10:00:00Z info A first",
		"2024-01-01T10:01:00Z info B first",
		"2024-01-01T10:02:00Z info A second",
		"2024-01-01T10:03:00Z info B second",
	}
	for i, e := range expected {
		if got[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, got[i])
		}
	}
}

func TestMerge_UntimedLinesAppendedLast(t *testing.T) {
	a := []string{"2024-01-01T10:00:00Z info timed"}
	b := []string{"no timestamp here", "also no timestamp"}
	got := merge.Merge(a, b)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "2024-01-01T10:00:00Z info timed" {
		t.Errorf("expected timed line first, got %q", got[0])
	}
	if got[1] != "no timestamp here" {
		t.Errorf("expected untimed line at index 1, got %q", got[1])
	}
}

func TestMerge_EmptySources(t *testing.T) {
	got := merge.Merge([]string{}, []string{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d lines", len(got))
	}
}

func TestMerge_NoSources(t *testing.T) {
	got := merge.Merge()
	if got == nil {
		t.Error("expected non-nil slice for no sources")
	}
	if len(got) != 0 {
		t.Errorf("expected 0 lines, got %d", len(got))
	}
}

func TestAnnotate_WithAndWithoutTimestamp(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z info hello",
		"plain text line",
	}
	annotated := merge.Annotate(lines)
	if len(annotated) != 2 {
		t.Fatalf("expected 2 annotated lines, got %d", len(annotated))
	}
	if annotated[0].Timestamp == nil {
		t.Error("expected non-nil timestamp for first line")
	}
	if annotated[1].Timestamp != nil {
		t.Error("expected nil timestamp for second line")
	}
}
