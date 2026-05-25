package cli

import (
	"bytes"
	"testing"
)

func TestParseFlags_Basic(t *testing.T) {
	args := []string{"--file", "app.log", "--range", "2024-01-01T10:00,2024-01-01T11:00"}
	f, err := ParseFlags(args, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.File != "app.log" {
		t.Errorf("File: want app.log, got %s", f.File)
	}
}

func TestParseFlags_MissingFile(t *testing.T) {
	args := []string{"--range", "2024-01-01T10:00,2024-01-01T11:00"}
	_, err := ParseFlags(args, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing --file")
	}
}

func TestParseFlags_MissingRange(t *testing.T) {
	args := []string{"--file", "app.log"}
	_, err := ParseFlags(args, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing --range")
	}
}

func TestParseFlags_IncludeExclude(t *testing.T) {
	args := []string{"--file", "x.log", "--range", "a,b", "--include", "ERROR", "--exclude", "DEBUG"}
	f, err := ParseFlags(args, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Include != "ERROR" {
		t.Errorf("Include: want ERROR, got %s", f.Include)
	}
	if f.Exclude != "DEBUG" {
		t.Errorf("Exclude: want DEBUG, got %s", f.Exclude)
	}
}

func TestParseFlags_BoolFlags(t *testing.T) {
	args := []string{"--file", "x.log", "--range", "a,b", "--strict", "--line-numbers", "--summary"}
	f, err := ParseFlags(args, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Strict {
		t.Error("Strict should be true")
	}
	if !f.LineNumbers {
		t.Error("LineNumbers should be true")
	}
	if !f.Summary {
		t.Error("Summary should be true")
	}
}

func TestParseFlags_DedupFlags(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantConsecutive bool
		wantGlobal      bool
	}{
		{
			name:            "dedup-consecutive",
			args:            []string{"--file", "x.log", "--range", "a,b", "--dedup-consecutive"},
			wantConsecutive: true,
		},
		{
			name:       "dedup-global",
			args:       []string{"--file", "x.log", "--range", "a,b", "--dedup-global"},
			wantGlobal: true,
		},
		{
			name:            "both dedup flags",
			args:            []string{"--file", "x.log", "--range", "a,b", "--dedup-consecutive", "--dedup-global"},
			wantConsecutive: true,
			wantGlobal:      true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := ParseFlags(tc.args, &bytes.Buffer{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f.DedupConsecutive != tc.wantConsecutive {
				t.Errorf("DedupConsecutive: want %v, got %v", tc.wantConsecutive, f.DedupConsecutive)
			}
			if f.DedupGlobal != tc.wantGlobal {
				t.Errorf("DedupGlobal: want %v, got %v", tc.wantGlobal, f.DedupGlobal)
			}
		})
	}
}
