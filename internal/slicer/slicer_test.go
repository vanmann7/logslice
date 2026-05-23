package slicer

import (
	"strings"
	"testing"
	"time"
)

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

const rfc3339 = "2006-01-02T15:04:05Z07:00"

var sampleLog = strings.Join([]string{
	"2024-01-10T08:00:00Z INFO  server started",
	"2024-01-10T08:01:00Z DEBUG request received",
	"2024-01-10T08:02:00Z INFO  processing",
	"2024-01-10T08:03:00Z WARN  slow query detected",
	"2024-01-10T08:04:00Z ERROR disk full",
	"2024-01-10T08:05:00Z INFO  cleanup done",
}, "\n")

func TestSlice_BasicRange(t *testing.T) {
	opts := Options{
		Start:   mustParseTime(rfc3339, "2024-01-10T08:01:00Z"),
		End:     mustParseTime(rfc3339, "2024-01-10T08:03:00Z"),
		Formats: []string{rfc3339},
	}

	r := strings.NewReader(sampleLog)
	var w strings.Builder

	res, err := Slice(r, &w, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.LinesRead != 6 {
		t.Errorf("expected 6 lines read, got %d", res.LinesRead)
	}
	if res.LinesWritten != 3 {
		t.Errorf("expected 3 lines written, got %d", res.LinesWritten)
	}
	if res.LinesSkipped != 3 {
		t.Errorf("expected 3 lines skipped, got %d", res.LinesSkipped)
	}

	output := w.String()
	if !strings.Contains(output, "DEBUG request received") {
		t.Error("expected DEBUG line in output")
	}
	if strings.Contains(output, "server started") {
		t.Error("did not expect 'server started' in output")
	}
}

func TestSlice_StrictModeExcludesNoTimestamp(t *testing.T) {
	log := "2024-01-10T08:01:00Z INFO start\nno timestamp here\n2024-01-10T08:02:00Z INFO end"
	opts := Options{
		Start:      mustParseTime(rfc3339, "2024-01-10T08:00:00Z"),
		End:        mustParseTime(rfc3339, "2024-01-10T08:05:00Z"),
		Formats:    []string{rfc3339},
		StrictMode: true,
	}

	r := strings.NewReader(log)
	var w strings.Builder

	res, err := Slice(r, &w, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.LinesWritten != 2 {
		t.Errorf("strict mode: expected 2 lines written, got %d", res.LinesWritten)
	}
	if strings.Contains(w.String(), "no timestamp here") {
		t.Error("strict mode: unexpected non-timestamp line in output")
	}
}

func TestSlice_EmptyInput(t *testing.T) {
	opts := Options{
		Start:   mustParseTime(rfc3339, "2024-01-10T08:00:00Z"),
		End:     mustParseTime(rfc3339, "2024-01-10T09:00:00Z"),
		Formats: []string{rfc3339},
	}

	r := strings.NewReader("")
	var w strings.Builder

	res, err := Slice(r, &w, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesRead != 0 || res.LinesWritten != 0 {
		t.Errorf("expected zero counts for empty input, got read=%d written=%d", res.LinesRead, res.LinesWritten)
	}
}
