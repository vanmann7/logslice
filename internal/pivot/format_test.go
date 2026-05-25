package pivot_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/pivot"
)

func buildResult() *pivot.Result {
	return pivot.MustApply(
		[]string{
			"level=info msg=a",
			"level=warn msg=b",
			"level=info msg=c",
		},
		pivot.Options{Pattern: `level=(?P<key>\w+)`},
	)
}

func TestFormat_ContainsAllLines(t *testing.T) {
	out := pivot.Format(buildResult(), pivot.FormatOptions{})
	for _, want := range []string{"msg=a", "msg=b", "msg=c"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestFormat_ShowCounts(t *testing.T) {
	out := pivot.Format(buildResult(), pivot.FormatOptions{ShowCounts: true})
	if !strings.Contains(out, "=== info (2 lines) ===") {
		t.Errorf("expected count header for info group, got:\n%s", out)
	}
	if !strings.Contains(out, "=== warn (1 lines) ===") {
		t.Errorf("expected count header for warn group, got:\n%s", out)
	}
}

func TestFormat_CustomSeparator(t *testing.T) {
	out := pivot.Format(buildResult(), pivot.FormatOptions{Separator: "---\n"})
	if !strings.Contains(out, "---") {
		t.Errorf("expected custom separator in output")
	}
}

func TestFormat_NilResult(t *testing.T) {
	out := pivot.Format(nil, pivot.FormatOptions{})
	if out != "" {
		t.Errorf("expected empty string for nil result, got %q", out)
	}
}

func TestFormat_EmptyResult(t *testing.T) {
	res := &pivot.Result{Groups: map[string][]string{}, Order: []string{}}
	out := pivot.Format(res, pivot.FormatOptions{})
	if out != "" {
		t.Errorf("expected empty string for empty result, got %q", out)
	}
}

func TestFormat_OrderRespected(t *testing.T) {
	res := buildResult()
	out := pivot.Format(res, pivot.FormatOptions{})
	idxInfo := strings.Index(out, "level=info")
	idxWarn := strings.Index(out, "level=warn")
	if idxInfo < 0 || idxWarn < 0 {
		t.Fatal("expected both groups in output")
	}
	if idxInfo > idxWarn {
		t.Errorf("expected info group before warn group")
	}
}
