package timeparse_test

import (
	"testing"

	"github.com/logslice/logslice/internal/timeparse"
)

func TestExtractFromLine(t *testing.T) {
	cases := []struct {
		name       string
		line       string
		wantRaw    string
		wantOffset int
		wantNil    bool
	}{
		{
			name:       "RFC3339 at start",
			line:       "2024-03-15T10:22:33Z INFO server started",
			wantRaw:    "2024-03-15T10:22:33Z",
			wantOffset: 0,
		},
		{
			name:       "RFC3339 with nanoseconds",
			line:       "2024-03-15T10:22:33.123456789Z ERROR disk full",
			wantRaw:    "2024-03-15T10:22:33.123456789Z",
			wantOffset: 0,
		},
		{
			name:       "space separated datetime",
			line:       "[2024-03-15 10:22:33] WARN slow query",
			wantRaw:    "2024-03-15 10:22:33",
			wantOffset: 1,
		},
		{
			name:       "slash date format",
			line:       "2024/03/15 10:22:33 DEBUG cache miss",
			wantRaw:    "2024/03/15 10:22:33",
			wantOffset: 0,
		},
		{
			name:       "apache combined log",
			line:       `192.168.1.1 - - [15/Mar/2024:10:22:33 +0000] "GET / HTTP/1.1" 200`,
			wantRaw:    "15/Mar/2024:10:22:33 +0000",
			wantOffset: 14,
		},
		{
			name:       "syslog format",
			line:       "Mar 15 10:22:33 hostname sshd: connection accepted",
			wantRaw:    "Mar 15 10:22:33",
			wantOffset: 0,
		},
		{
			name:    "no timestamp",
			line:    "this line has no timestamp at all",
			wantNil: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantNil: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := timeparse.ExtractFromLine(tc.line)
			if tc.wantNil {
				if result != nil {
					t.Errorf("expected nil, got %+v", result)
				}
				return
			}
			if result == nil {
				t.Fatalf("expected result, got nil")
			}
			if result.Raw != tc.wantRaw {
				t.Errorf("Raw: got %q, want %q", result.Raw, tc.wantRaw)
			}
			if result.Offset != tc.wantOffset {
				t.Errorf("Offset: got %d, want %d", result.Offset, tc.wantOffset)
			}
		})
	}
}
