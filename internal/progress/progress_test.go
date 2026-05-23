package progress

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_SilentWhenNoOutput(t *testing.T) {
	r := New(1024, nil)
	if r.enabled {
		t.Error("expected reporter to be disabled when output is nil")
	}
}

func TestNew_SilentWhenTotalZero(t *testing.T) {
	var buf bytes.Buffer
	r := New(0, &buf)
	if r.enabled {
		t.Error("expected reporter to be disabled when total is 0")
	}
}

func TestPercent_Basic(t *testing.T) {
	r := New(200, nil)
	r.Add(50)
	got := r.Percent()
	if got != 25.0 {
		t.Errorf("expected 25.0%%, got %.1f%%", got)
	}
}

func TestPercent_ExceedsTotal(t *testing.T) {
	r := New(100, nil)
	r.Add(150)
	if r.Percent() != 100.0 {
		t.Errorf("expected percent capped at 100, got %.1f", r.Percent())
	}
}

func TestPercent_ZeroTotal(t *testing.T) {
	r := New(0, nil)
	r.Add(999)
	if r.Percent() != 0 {
		t.Errorf("expected 0 when total unknown, got %.1f", r.Percent())
	}
}

func TestPrint_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	r := New(1024, &buf)
	r.Add(512)
	r.Print()
	out := buf.String()
	if !strings.Contains(out, "progress:") {
		t.Errorf("expected 'progress:' in output, got: %q", out)
	}
	if !strings.Contains(out, "50.0%") {
		t.Errorf("expected '50.0%%' in output, got: %q", out)
	}
}

func TestPrint_SilentMode(t *testing.T) {
	var buf bytes.Buffer
	r := New(0, &buf)
	r.Print()
	if buf.Len() != 0 {
		t.Errorf("expected no output in silent mode, got: %q", buf.String())
	}
}

func TestDone_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	r := New(2048, &buf)
	r.Add(2048)
	r.Done()
	out := buf.String()
	if !strings.Contains(out, "done") {
		t.Errorf("expected 'done' in output, got: %q", out)
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input int64
		want  string
	}{
		{500, "500 B"},
		{2048, "2.00 KiB"},
		{1048576, "1.00 MiB"},
		{1073741824, "1.00 GiB"},
	}
	for _, c := range cases {
		got := formatBytes(c.input)
		if got != c.want {
			t.Errorf("formatBytes(%d) = %q, want %q", c.input, got, c.want)
		}
	}
}
