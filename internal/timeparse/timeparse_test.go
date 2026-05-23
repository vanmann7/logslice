package timeparse_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/timeparse"
)

func TestParseTimestamp(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
		wantUTC string
	}{
		{"2024-03-15T10:22:33Z", false, "2024-03-15 10:22:33 +0000 UTC"},
		{"2024-03-15T10:22:33.123456789Z", false, "2024-03-15 10:22:33.123456789 +0000 UTC"},
		{"2024-03-15 10:22:33", false, "2024-03-15 10:22:33 +0000 UTC"},
		{"2024/03/15 10:22:33", false, "2024-03-15 10:22:33 +0000 UTC"},
		{"not-a-timestamp", true, ""},
		{"", true, ""},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, _, err := timeparse.ParseTimestamp(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.UTC().String() != tc.wantUTC {
				t.Errorf("got %q, want %q", got.UTC().String(), tc.wantUTC)
			}
		})
	}
}

func TestParseRange(t *testing.T) {
	t.Run("valid range", func(t *testing.T) {
		start, end, err := timeparse.ParseRange("2024-03-15T10:00:00Z", "2024-03-15T11:00:00Z")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !start.Before(end) {
			t.Error("expected start to be before end")
		}
	})

	t.Run("inverted range", func(t *testing.T) {
		_, _, err := timeparse.ParseRange("2024-03-15T11:00:00Z", "2024-03-15T10:00:00Z")
		if err == nil {
			t.Error("expected error for inverted range")
		}
	})

	t.Run("open start", func(t *testing.T) {
		start, end, err := timeparse.ParseRange("", "2024-03-15T11:00:00Z")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !start.IsZero() {
			t.Error("expected zero start time")
		}
		if end.IsZero() {
			t.Error("expected non-zero end time")
		}
	})

	t.Run("open end", func(t *testing.T) {
		start, end, err := timeparse.ParseRange("2024-03-15T10:00:00Z", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start.IsZero() {
			t.Error("expected non-zero start time")
		}
		_ = end
	})

	t.Run("both empty", func(t *testing.T) {
		start, end, err := timeparse.ParseRange("", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !start.IsZero() || !end.IsZero() {
			t.Error("expected both times to be zero")
		}
	})

	t.Run("equal timestamps", func(t *testing.T) {
		_, _, err := timeparse.ParseRange("2024-03-15T10:00:00Z", "2024-03-15T10:00:00Z")
		if err != nil {
			t.Errorf("unexpected error for equal timestamps: %v", err)
		}
	})
}

func TestParseTimestampReturnsFormat(t *testing.T) {
	_, fmt, err := timeparse.ParseTimestamp("2024-03-15T10:22:33Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != time.RFC3339 {
		t.Errorf("expected RFC3339 format, got %q", fmt)
	}
}
