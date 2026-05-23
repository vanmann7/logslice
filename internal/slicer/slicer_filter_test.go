package slicer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/slicer"
)

func TestSliceFiltered_KeywordInclude(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T10:00:00Z INFO service started",
		"2024-01-01T10:01:00Z ERROR database timeout",
		"2024-01-01T10:02:00Z INFO request handled",
		"2024-01-01T10:03:00Z ERROR disk full",
	}, "\n")

	from := mustParseTime("2024-01-01T10:00:00Z")
	to := mustParseTime("2024-01-01T10:05:00Z")

	opts := slicer.FilteredSliceOptions{
		From: from,
		To:   to,
		FilterOpts: filter.Options{
			MustContain:   []string{"ERROR"},
			CaseSensitive: true,
		},
	}

	var out strings.Builder
	if err := slicer.SliceFiltered(strings.NewReader(input), &out, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 ERROR lines, got %d: %v", len(lines), lines)
	}
}

func TestSliceFiltered_KeywordExclude(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T10:00:00Z INFO service started",
		"2024-01-01T10:01:00Z DEBUG internal state",
		"2024-01-01T10:02:00Z INFO request handled",
	}, "\n")

	from := mustParseTime("2024-01-01T10:00:00Z")
	to := mustParseTime("2024-01-01T10:05:00Z")

	opts := slicer.FilteredSliceOptions{
		From: from,
		To:   to,
		FilterOpts: filter.Options{
			MustExclude:   []string{"DEBUG"},
			CaseSensitive: true,
		},
	}

	var out strings.Builder
	if err := slicer.SliceFiltered(strings.NewReader(input), &out, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(out.String(), "DEBUG") {
		t.Error("expected DEBUG lines to be excluded")
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestSliceFiltered_StrictDropsNoTimestamp(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T10:00:00Z INFO start",
		"continuation line without timestamp",
		"2024-01-01T10:01:00Z INFO end",
	}, "\n")

	from := mustParseTime("2024-01-01T10:00:00Z")
	to := mustParseTime("2024-01-01T10:05:00Z")

	opts := slicer.FilteredSliceOptions{
		From:   from,
		To:     to,
		Strict: true,
	}

	var out strings.Builder
	if err := slicer.SliceFiltered(strings.NewReader(input), &out, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(out.String(), "continuation") {
		t.Error("strict mode should drop lines without timestamps")
	}
}

func TestSliceFiltered_EmptyInput(t *testing.T) {
	opts := slicer.FilteredSliceOptions{
		From: time.Now().Add(-time.Hour),
		To:   time.Now(),
	}
	var out strings.Builder
	if err := slicer.SliceFiltered(strings.NewReader(""), &out, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() != 0 {
		t.Error("expected empty output for empty input")
	}
}
