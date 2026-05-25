package dedup

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestApply_NoOptions_ReturnsUnchanged(t *testing.T) {
	lines := []string{"a", "a", "b", "b"}
	got := Apply(lines, Options{})
	if diff := cmp.Diff(lines, got); diff != "" {
		t.Errorf("unexpected diff (-want +got):\n%s", diff)
	}
}

func TestApply_Consecutive_RemovesAdjacentDupes(t *testing.T) {
	lines := []string{"a", "a", "b", "a", "a"}
	want := []string{"a", "b", "a"}
	got := Apply(lines, Options{Consecutive: true})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected diff (-want +got):\n%s", diff)
	}
}

func TestApply_Consecutive_NoDupes(t *testing.T) {
	lines := []string{"a", "b", "c"}
	got := Apply(lines, Options{Consecutive: true})
	if diff := cmp.Diff(lines, got); diff != "" {
		t.Errorf("unexpected diff (-want +got):\n%s", diff)
	}
}

func TestApply_Global_RemovesAllDupes(t *testing.T) {
	lines := []string{"a", "b", "a", "c", "b"}
	want := []string{"a", "b", "c"}
	got := Apply(lines, Options{Global: true})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected diff (-want +got):\n%s", diff)
	}
}

func TestApply_Global_PreservesOrder(t *testing.T) {
	lines := []string{"z", "a", "z", "b", "a"}
	want := []string{"z", "a", "b"}
	got := Apply(lines, Options{Global: true})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected diff (-want +got):\n%s", diff)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	for _, opts := range []Options{
		{},
		{Consecutive: true},
		{Global: true},
	} {
		got := Apply([]string{}, opts)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	}
}

func TestApply_GlobalPreferredOverConsecutive(t *testing.T) {
	// When both flags set, Global wins.
	lines := []string{"a", "b", "a"}
	want := []string{"a", "b"}
	got := Apply(lines, Options{Consecutive: true, Global: true})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected diff (-want +got):\n%s", diff)
	}
}
