package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestLine_NoTerms(t *testing.T) {
	line := "no terms here"
	result := highlight.Line(line, highlight.Options{})
	if result != line {
		t.Errorf("expected unchanged line, got %q", result)
	}
}

func TestLine_SingleTermHighlighted(t *testing.T) {
	result := highlight.Line("hello world", highlight.Options{
		Terms: []string{"world"},
	})
	if !strings.Contains(result, "world") {
		t.Error("expected result to contain original term")
	}
	if !strings.Contains(result, highlight.AnsiYellow) {
		t.Error("expected result to contain ANSI color code")
	}
	if !strings.Contains(result, highlight.AnsiReset) {
		t.Error("expected result to contain ANSI reset code")
	}
}

func TestLine_CaseInsensitiveMatch(t *testing.T) {
	result := highlight.Line("Error occurred", highlight.Options{
		Terms:         []string{"error"},
		CaseSensitive: false,
	})
	if !strings.Contains(result, "Error") {
		t.Error("expected original casing to be preserved")
	}
	if !strings.Contains(result, highlight.AnsiYellow) {
		t.Error("expected highlight to be applied")
	}
}

func TestLine_CaseSensitiveNoMatch(t *testing.T) {
	line := "Error occurred"
	result := highlight.Line(line, highlight.Options{
		Terms:         []string{"error"},
		CaseSensitive: true,
	})
	if strings.Contains(result, highlight.AnsiYellow) {
		t.Error("expected no highlight when case does not match")
	}
}

func TestLine_MultipleTerms(t *testing.T) {
	result := highlight.Line("foo bar baz", highlight.Options{
		Terms: []string{"foo", "baz"},
	})
	if !strings.Contains(result, highlight.AnsiYellow) {
		t.Error("expected highlights to be applied")
	}
	// Both terms should appear in output
	if !strings.Contains(result, "foo") || !strings.Contains(result, "baz") {
		t.Error("expected both terms present in output")
	}
}

func TestLine_CustomColor(t *testing.T) {
	result := highlight.Line("critical failure", highlight.Options{
		Terms: []string{"critical"},
		Color: highlight.AnsiRed,
	})
	if !strings.Contains(result, highlight.AnsiRed) {
		t.Error("expected custom color to be used")
	}
}

func TestLines_AppliedToAll(t *testing.T) {
	input := []string{"foo log", "bar log", "baz entry"}
	result := highlight.Lines(input, highlight.Options{
		Terms: []string{"log"},
	})
	if len(result) != len(input) {
		t.Fatalf("expected %d lines, got %d", len(input), len(result))
	}
	for i, r := range result {
		if strings.Contains(input[i], "log") && !strings.Contains(r, highlight.AnsiYellow) {
			t.Errorf("line %d: expected highlight applied", i)
		}
	}
}

func TestLine_EmptyTerm(t *testing.T) {
	line := "some log line"
	result := highlight.Line(line, highlight.Options{Terms: []string{""}})
	if strings.Contains(result, highlight.AnsiYellow) {
		t.Error("empty term should not produce highlights")
	}
}
