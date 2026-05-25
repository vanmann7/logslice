package sampler_test

import (
	"testing"

	"github.com/logslice/logslice/internal/sampler"
)

func makeLines(n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("line-%d", i)
	}
	return lines
}

import "fmt"

func TestApply_Every1ReturnsAll(t *testing.T) {
	lines := makeLines(5)
	out, err := sampler.Apply(lines, sampler.Options{Every: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 5 {
		t.Errorf("expected 5 lines, got %d", len(out))
	}
}

func TestApply_Every2ReturnsHalf(t *testing.T) {
	lines := makeLines(6)
	out, err := sampler.Apply(lines, sampler.Options{Every: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 lines, got %d", len(out))
	}
	for _, l := range out {
		t.Log(l)
	}
}

func TestApply_MaxLinesCaps(t *testing.T) {
	lines := makeLines(20)
	out, err := sampler.Apply(lines, sampler.Options{Every: 1, MaxLines: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 5 {
		t.Errorf("expected 5 lines, got %d", len(out))
	}
}

func TestApply_Every3WithMax(t *testing.T) {
	lines := makeLines(30)
	out, err := sampler.Apply(lines, sampler.Options{Every: 3, MaxLines: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 4 {
		t.Errorf("expected 4 lines, got %d", len(out))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	out, err := sampler.Apply([]string{}, sampler.Options{Every: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected 0 lines, got %d", len(out))
	}
}

func TestApply_InvalidEvery(t *testing.T) {
	_, err := sampler.Apply([]string{"a"}, sampler.Options{Every: 0})
	if err == nil {
		t.Error("expected error for Every=0, got nil")
	}
}

func TestApply_InvalidMaxLines(t *testing.T) {
	_, err := sampler.Apply([]string{"a"}, sampler.Options{Every: 1, MaxLines: -1})
	if err == nil {
		t.Error("expected error for MaxLines=-1, got nil")
	}
}

func TestMustApply_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	sampler.MustApply([]string{"x"}, sampler.Options{Every: -1})
}
