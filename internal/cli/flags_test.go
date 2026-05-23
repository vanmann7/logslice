package cli

import (
	"testing"
)

func TestParseFlags_Basic(t *testing.T) {
	args := []string{"-start", "2024-01-15 10:00:00", "-end", "2024-01-15 11:00:00", "app.log"}
	opts, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.FilePath != "app.log" {
		t.Errorf("expected FilePath=app.log, got %q", opts.FilePath)
	}
	if opts.Start != "2024-01-15 10:00:00" {
		t.Errorf("unexpected Start: %q", opts.Start)
	}
	if opts.End != "2024-01-15 11:00:00" {
		t.Errorf("unexpected End: %q", opts.End)
	}
}

func TestParseFlags_MissingFile(t *testing.T) {
	args := []string{"-start", "2024-01-15 10:00:00"}
	_, err := ParseFlags(args)
	if err == nil {
		t.Fatal("expected error for missing file argument")
	}
}

func TestParseFlags_MissingRange(t *testing.T) {
	args := []string{"app.log"}
	_, err := ParseFlags(args)
	if err == nil {
		t.Fatal("expected error when neither -start nor -end is provided")
	}
}

func TestParseFlags_IncludeExclude(t *testing.T) {
	args := []string{
		"-start", "2024-01-15 10:00:00",
		"-include", "ERROR",
		"-include", "WARN",
		"-exclude", "DEBUG",
		"app.log",
	}
	opts, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Include) != 2 {
		t.Errorf("expected 2 include keywords, got %d", len(opts.Include))
	}
	if len(opts.Exclude) != 1 {
		t.Errorf("expected 1 exclude keyword, got %d", len(opts.Exclude))
	}
	if opts.Include[0] != "ERROR" || opts.Include[1] != "WARN" {
		t.Errorf("unexpected include values: %v", opts.Include)
	}
	if opts.Exclude[0] != "DEBUG" {
		t.Errorf("unexpected exclude value: %v", opts.Exclude)
	}
}

func TestParseFlags_BoolFlags(t *testing.T) {
	args := []string{
		"-start", "2024-01-15 10:00:00",
		"-strict",
		"-line-numbers",
		"-summary",
		"app.log",
	}
	opts, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Strict {
		t.Error("expected Strict=true")
	}
	if !opts.LineNumbers {
		t.Error("expected LineNumbers=true")
	}
	if !opts.Summary {
		t.Error("expected Summary=true")
	}
}

func TestParseFlags_OutputFile(t *testing.T) {
	args := []string{"-start", "2024-01-15 10:00:00", "-output", "out.txt", "app.log"}
	opts, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Output != "out.txt" {
		t.Errorf("expected Output=out.txt, got %q", opts.Output)
	}
}
