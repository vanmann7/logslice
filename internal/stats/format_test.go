package stats

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestFormat_ContainsExpectedFields(t *testing.T) {
	c := NewCollector()
	for i := 0; i < 10; i++ {
		c.RecordTotal()
	}
	for i := 0; i < 5; i++ {
		c.RecordMatched(nil)
	}
	c.RecordSkipped()
	c.RecordFiltered()

	r := c.Finalise()
	var buf bytes.Buffer
	Format(&buf, r)
	out := buf.String()

	for _, want := range []string{
		"total lines",
		"matched lines",
		"skipped lines",
		"filtered out",
		"elapsed",
		"10",
		"5",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestFormat_ShowsTimeBounds(t *testing.T) {
	c := NewCollector()
	t1 := time.Date(2024, 3, 15, 8, 30, 0, 0, time.UTC)
	t2 := time.Date(2024, 3, 15, 9, 45, 0, 0, time.UTC)
	c.RecordMatched(&t1)
	c.RecordMatched(&t2)

	r := c.Finalise()
	var buf bytes.Buffer
	Format(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "2024-03-15 08:30:00") {
		t.Errorf("expected earliest time in output, got:\n%s", out)
	}
	if !strings.Contains(out, "2024-03-15 09:45:00") {
		t.Errorf("expected latest time in output, got:\n%s", out)
	}
}

func TestFormat_NoTimesWhenNoneRecorded(t *testing.T) {
	c := NewCollector()
	r := c.Finalise()
	var buf bytes.Buffer
	Format(&buf, r)
	out := buf.String()

	if strings.Contains(out, "earliest") {
		t.Errorf("did not expect earliest time line when no timestamps recorded")
	}
	if strings.Contains(out, "latest") {
		t.Errorf("did not expect latest time line when no timestamps recorded")
	}
}
