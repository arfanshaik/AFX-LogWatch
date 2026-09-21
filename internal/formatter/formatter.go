package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/arfanshaik/AFX-LogWatch/internal/parser"
)

type Options struct {
	JSON  bool
	Color bool
}

type Formatter struct {
	w    io.Writer
	opts Options
}

func New(w io.Writer, opts Options) *Formatter { return &Formatter{w: w, opts: opts} }

func (f *Formatter) Write(e parser.Entry) error {
	if f.opts.JSON {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(f.w, string(b))
		return err
	}

	level := e.Level
	if f.opts.Color {
		level = colorize(e.Level)
	}
	_, err := fmt.Fprintf(f.w, "%-18s %-9s %s\n", filepath.Base(e.Source), level, e.Message)
	return err
}

func colorize(level string) string {
	switch level {
	case "trace":
		return "\x1b[90mTRACE\x1b[0m"
	case "debug":
		return "\x1b[36mDEBUG\x1b[0m"
	case "info":
		return "\x1b[32mINFO\x1b[0m"
	case "warn":
		return "\x1b[33mWARN\x1b[0m"
	case "error":
		return "\x1b[31mERROR\x1b[0m"
	case "fatal":
		return "\x1b[1;31mFATAL\x1b[0m"
	default:
		return "\x1b[90mUNKNOWN\x1b[0m"
	}
}
