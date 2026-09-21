package parser

import "testing"

func TestParseLevel(t *testing.T) {
	tests := []struct {
		line, want string
	}{
		{"2026-09-21 INFO server started", "info"},
		{"[ERROR] database failed", "error"},
		{"WARNING: disk nearly full", "warn"},
		{"plain message", "unknown"},
	}
	for _, tt := range tests {
		if got := Parse("test.log", tt.line).Level; got != tt.want {
			t.Fatalf("Parse(%q) level=%q want=%q", tt.line, got, tt.want)
		}
	}
}

func TestFilter(t *testing.T) {
	f, err := NewFilter("timeout", "ignore", []string{"error"})
	if err != nil {
		t.Fatal(err)
	}
	if !f.Match(Parse("x", "ERROR timeout contacting API")) {
		t.Fatal("expected line to match")
	}
	if f.Match(Parse("x", "ERROR timeout ignore this")) {
		t.Fatal("excluded line matched")
	}
}
