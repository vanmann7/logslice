package context_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/cli"
	runctx "github.com/user/logslice/internal/context"
)

func defaultFlags() cli.Flags {
	return cli.Flags{
		File:  "test.log",
		Range: "10:00-11:00",
		Quiet: false,
	}
}

func TestNew_ReturnsNonNilFields(t *testing.T) {
	rc := runctx.New(defaultFlags(), 1024)
	if rc == nil {
		t.Fatal("expected non-nil RunContext")
	}
	if rc.Stats == nil {
		t.Error("expected non-nil Stats collector")
	}
	if rc.Progress == nil {
		t.Error("expected non-nil Progress reporter")
	}
	if rc.Output == nil {
		t.Error("expected non-nil Output writer")
	}
}

func TestNewWithWriter_UsesProvidedWriter(t *testing.T) {
	var buf bytes.Buffer
	flags := defaultFlags()
	rc := runctx.NewWithWriter(flags, 512, &buf)
	if rc.Output != &buf {
		t.Error("expected Output to be the provided writer")
	}
}

func TestElapsed_IsPositive(t *testing.T) {
	rc := runctx.New(defaultFlags(), 0)
	time.Sleep(1 * time.Millisecond)
	if rc.Elapsed() <= 0 {
		t.Error("expected positive elapsed duration")
	}
}

func TestElapsed_IncreasesOverTime(t *testing.T) {
	rc := runctx.New(defaultFlags(), 0)
	first := rc.Elapsed()
	time.Sleep(2 * time.Millisecond)
	second := rc.Elapsed()
	if second <= first {
		t.Errorf("expected elapsed to increase: first=%v second=%v", first, second)
	}
}

func TestSummary_ContainsExpectedText(t *testing.T) {
	rc := runctx.New(defaultFlags(), 0)
	summary := rc.Summary()
	if !strings.Contains(summary, "Lines") {
		t.Errorf("expected summary to mention Lines, got: %q", summary)
	}
}

func TestNew_QuietFlagDoesNotPanic(t *testing.T) {
	flags := defaultFlags()
	flags.Quiet = true
	rc := runctx.New(flags, 256)
	if rc == nil {
		t.Fatal("expected non-nil RunContext with quiet flag")
	}
}
