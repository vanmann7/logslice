package pivot_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/pivot"
)

func makeLines() []string {
	return []string{
		"2024-01-01 level=info msg=started",
		"2024-01-01 level=warn msg=slow",
		"2024-01-01 level=info msg=ok",
		"2024-01-01 level=error msg=fail",
		"2024-01-01 no-level-field here",
	}
}

func TestApply_GroupsByLevel(t *testing.T) {
	res, err := pivot.Apply(makeLines(), pivot.Options{
		Pattern: `level=(?P<key>\w+)`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups["info"]) != 2 {
		t.Errorf("expected 2 info lines, got %d", len(res.Groups["info"]))
	}
	if len(res.Groups["warn"]) != 1 {
		t.Errorf("expected 1 warn line, got %d", len(res.Groups["warn"]))
	}
	if len(res.Groups["error"]) != 1 {
		t.Errorf("expected 1 error line, got %d", len(res.Groups["error"]))
	}
}

func TestApply_UnmatchedGoesToFallback(t *testing.T) {
	res, err := pivot.Apply(makeLines(), pivot.Options{
		Pattern:     `level=(?P<key>\w+)`,
		FallbackKey: "other",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups["other"]) != 1 {
		t.Errorf("expected 1 unmatched line under 'other', got %d", len(res.Groups["other"]))
	}
}

func TestApply_DefaultFallbackKey(t *testing.T) {
	res, err := pivot.Apply(makeLines(), pivot.Options{
		Pattern: `level=(?P<key>\w+)`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Groups["<unmatched>"]; !ok {
		t.Error("expected default fallback key '<unmatched>' to exist")
	}
}

func TestApply_OrderPreserved(t *testing.T) {
	res, err := pivot.Apply(makeLines(), pivot.Options{
		Pattern: `level=(?P<key>\w+)`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Order) == 0 {
		t.Fatal("expected non-empty order slice")
	}
	if res.Order[0] != "info" {
		t.Errorf("expected first key to be 'info', got %q", res.Order[0])
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	_, err := pivot.Apply(makeLines(), pivot.Options{Pattern: `(?P<key`})
	if err == nil {
		t.Error("expected error for invalid regexp pattern")
	}
}

func TestApply_EmptyLines(t *testing.T) {
	res, err := pivot.Apply([]string{}, pivot.Options{
		Pattern: `level=(?P<key>\w+)`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(res.Groups))
	}
}

func TestMustApply_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid pattern")
		}
	}()
	pivot.MustApply(makeLines(), pivot.Options{Pattern: `(?P<key`})
}
