package watcher

import (
	"strings"
	"testing"
)

func TestScanLines(t *testing.T) {
	var got []string
	err := ScanLines(strings.NewReader("a\nb\nc\n"), func(line string) error {
		got = append(got, line)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Fatalf("unexpected lines: %#v", got)
	}
}
