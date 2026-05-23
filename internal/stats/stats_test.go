package stats

import (
	"testing"
	"time"
)

func TestCollector_BasicCounts(t *testing.T) {
	c := NewCollector()
	c.RecordTotal()
	c.RecordTotal()
	c.RecordTotal()
	c.RecordMatched(nil)
	c.RecordSkipped()
	c.RecordFiltered()

	r := c.Finalise()
	if r.TotalLines != 3 {
		t.Errorf("expected TotalLines=3, got %d", r.TotalLines)
	}
	if r.MatchedLines != 1 {
		t.Errorf("expected MatchedLines=1, got %d", r.MatchedLines)
	}
	if r.SkippedLines != 1 {
		t.Errorf("expected SkippedLines=1, got %d", r.SkippedLines)
	}
	if r.FilteredLines != 1 {
		t.Errorf("expected FilteredLines=1, got %d", r.FilteredLines)
	}
}

func TestCollector_TimeBounds(t *testing.T) {
	c := NewCollector()

	t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)

	c.RecordMatched(&t1)
	c.RecordMatched(&t2)
	c.RecordMatched(&t3)

	r := c.Finalise()
	if r.EarliestTime == nil || !r.EarliestTime.Equal(t1) {
		t.Errorf("expected EarliestTime=%v, got %v", t1, r.EarliestTime)
	}
	if r.LatestTime == nil || !r.LatestTime.Equal(t2) {
		t.Errorf("expected LatestTime=%v, got %v", t2, r.LatestTime)
	}
}

func TestCollector_NilTimeDoesNotPanic(t *testing.T) {
	c := NewCollector()
	c.RecordMatched(nil)
	r := c.Finalise()
	if r.EarliestTime != nil || r.LatestTime != nil {
		t.Error("expected nil time bounds when no timestamps recorded")
	}
}

func TestCollector_DurationPositive(t *testing.T) {
	c := NewCollector()
	time.Sleep(time.Millisecond)
	r := c.Finalise()
	if r.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", r.Duration)
	}
}

func TestResult_ZeroValue(t *testing.T) {
	c := NewCollector()
	r := c.Finalise()
	if r.TotalLines != 0 || r.MatchedLines != 0 {
		t.Error("expected zero counts on empty collector")
	}
}
