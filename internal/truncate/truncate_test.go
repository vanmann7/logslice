package truncate

import (
	"strings"
	"testing"
)

func TestLine_ShortLineUnchanged(t *testing.T) {
	input := "short line"
	got := Line(input, Options{MaxBytes: 64})
	if got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestLine_ExactLengthUnchanged(t *testing.T) {
	input := strings.Repeat("a", 512)
	got := Line(input, Options{})
	if got != input {
		t.Errorf("line of exact max length should not be truncated")
	}
}

func TestLine_LongLineTruncated(t *testing.T) {
	input := strings.Repeat("x", 600)
	got := Line(input, Options{MaxBytes: 512})
	if len(got) > 512 {
		t.Errorf("truncated line length %d exceeds max 512", len(got))
	}
	if !strings.HasSuffix(got, " [truncated]") {
		t.Errorf("expected truncation indicator, got: %q", got[len(got)-20:])
	}
}

func TestLine_CustomIndicator(t *testing.T) {
	input := strings.Repeat("y", 100)
	opts := Options{MaxBytes: 20, Indicator: "..."}
	got := Line(input, opts)
	if len(got) > 20 {
		t.Errorf("line length %d exceeds max 20", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected custom indicator, got %q", got)
	}
}

func TestLine_DefaultMaxBytes(t *testing.T) {
	// Zero MaxBytes should fall back to 512.
	input := strings.Repeat("z", 1000)
	got := Line(input, Options{})
	if len(got) > 512 {
		t.Errorf("default max bytes not applied, got length %d", len(got))
	}
}

func TestLine_UTF8Boundary(t *testing.T) {
	// Build a string where truncation would fall in the middle of a multi-byte rune.
	base := strings.Repeat("a", 10)
	// Append a 3-byte UTF-8 rune (€ = 0xE2 0x82 0xAC) so the cut might land inside it.
	rune3 := "€" // 3 bytes
	input := base + rune3 + strings.Repeat("b", 200)
	opts := Options{MaxBytes: 12, Indicator: "[t]"}
	got := Line(input, opts)
	if len(got) > 12 {
		t.Errorf("result length %d exceeds max 12", len(got))
	}
	// Must be valid UTF-8 prefix (no partial rune bytes before indicator).
	result := strings.TrimSuffix(got, "[t]")
	for i, b := range []byte(result) {
		if b&0xC0 == 0x80 && i == 0 {
			t.Errorf("result starts with a continuation byte")
		}
	}
}

func TestLines_AppliedToAll(t *testing.T) {
	input := []string{
		"short",
		strings.Repeat("L", 600),
		"also short",
	}
	opts := Options{MaxBytes: 50}
	got := Lines(input, opts)
	if len(got) != len(input) {
		t.Fatalf("expected %d lines, got %d", len(input), len(got))
	}
	if got[0] != "short" {
		t.Errorf("short line should be unchanged")
	}
	if len(got[1]) > 50 {
		t.Errorf("long line not truncated, length %d", len(got[1]))
	}
	if got[2] != "also short" {
		t.Errorf("short line should be unchanged")
	}
}
