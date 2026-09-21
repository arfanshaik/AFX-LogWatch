package parser

import (
	"regexp"
	"strings"
	"time"
)

type Entry struct {
	Time    time.Time `json:"time"`
	Source  string    `json:"source"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Raw     string    `json:"raw"`
}

var levelPattern = regexp.MustCompile(`(?i)(?:^|[\s\[\]():=-])(TRACE|DEBUG|INFO|WARN|WARNING|ERROR|FATAL|PANIC)(?:$|[\s\[\]():=-])`)

func Parse(source, line string) Entry {
	level := "unknown"
	if m := levelPattern.FindStringSubmatch(line); len(m) > 1 {
		level = normalizeLevel(m[1])
	}
	return Entry{
		Time:    time.Now(),
		Source:  source,
		Level:   level,
		Message: strings.TrimSpace(line),
		Raw:     line,
	}
}

func normalizeLevel(level string) string {
	switch strings.ToLower(level) {
	case "warning":
		return "warn"
	case "panic":
		return "fatal"
	default:
		return strings.ToLower(level)
	}
}

type Filter struct {
	include *regexp.Regexp
	exclude *regexp.Regexp
	levels  map[string]struct{}
}

func NewFilter(include, exclude string, levels []string) (*Filter, error) {
	f := &Filter{levels: map[string]struct{}{}}
	var err error
	if include != "" {
		f.include, err = regexp.Compile(include)
		if err != nil {
			return nil, err
		}
	}
	if exclude != "" {
		f.exclude, err = regexp.Compile(exclude)
		if err != nil {
			return nil, err
		}
	}
	for _, level := range levels {
		f.levels[normalizeLevel(strings.TrimSpace(level))] = struct{}{}
	}
	return f, nil
}

func (f *Filter) Match(e Entry) bool {
	if f.include != nil && !f.include.MatchString(e.Raw) {
		return false
	}
	if f.exclude != nil && f.exclude.MatchString(e.Raw) {
		return false
	}
	if len(f.levels) > 0 {
		_, ok := f.levels[e.Level]
		return ok
	}
	return true
}
