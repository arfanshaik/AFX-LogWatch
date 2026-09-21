package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/arfanshaik/AFX-LogWatch/internal/formatter"
	"github.com/arfanshaik/AFX-LogWatch/internal/parser"
	"github.com/arfanshaik/AFX-LogWatch/internal/stats"
	"github.com/arfanshaik/AFX-LogWatch/internal/watcher"
)

const version = "1.0.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}

	switch args[0] {
	case "watch":
		return runWatch(args[1:])
	case "scan":
		return runScan(args[1:])
	case "version", "--version", "-v":
		fmt.Printf("AFX LogWatch v%s\n", version)
		return nil
	case "help", "--help", "-h":
		printHelp()
		return nil
	default:
		return fmt.Errorf("unknown command %q; use 'afx-logwatch help'", args[0])
	}
}

func runWatch(args []string) error {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)
	include := fs.String("include", "", "only show lines matching this regex")
	exclude := fs.String("exclude", "", "hide lines matching this regex")
	levels := fs.String("level", "", "comma-separated levels: trace,debug,info,warn,error,fatal")
	jsonOutput := fs.Bool("json", false, "emit JSON lines")
	noColor := fs.Bool("no-color", false, "disable ANSI colors")
	fromStart := fs.Bool("from-start", false, "read existing content before following new lines")
	poll := fs.Duration("poll", 250*time.Millisecond, "file polling interval")
	statsEvery := fs.Duration("stats-every", 0, "print stats periodically, e.g. 30s (0 disables)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	files := fs.Args()
	if len(files) == 0 {
		return errors.New("provide at least one log file: afx-logwatch watch app.log")
	}

	filter, err := parser.NewFilter(*include, *exclude, splitCSV(*levels))
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	out := formatter.New(os.Stdout, formatter.Options{JSON: *jsonOutput, Color: !*noColor})
	counter := stats.New()

	events := make(chan watcher.Event, 256)
	errs := make(chan error, len(files))

	for _, path := range files {
		cfg := watcher.Config{Path: path, FromStart: *fromStart, PollInterval: *poll}
		go func() {
			if err := watcher.Follow(ctx, cfg, events); err != nil && !errors.Is(err, context.Canceled) {
				errs <- err
			}
		}()
	}

	var ticker *time.Ticker
	var tick <-chan time.Time
	if *statsEvery > 0 {
		ticker = time.NewTicker(*statsEvery)
		defer ticker.Stop()
		tick = ticker.C
	}

	for {
		select {
		case <-ctx.Done():
			if !*jsonOutput {
				fmt.Fprintln(os.Stdout)
				counter.WriteSummary(os.Stdout)
			}
			return nil
		case err := <-errs:
			return err
		case ev := <-events:
			entry := parser.Parse(ev.Source, ev.Line)
			if !filter.Match(entry) {
				continue
			}
			counter.Add(entry.Level)
			if err := out.Write(entry); err != nil {
				return err
			}
		case <-tick:
			if !*jsonOutput {
				counter.WriteSummary(os.Stdout)
			}
		}
	}
}

func runScan(args []string) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	include := fs.String("include", "", "only show lines matching this regex")
	exclude := fs.String("exclude", "", "hide lines matching this regex")
	levels := fs.String("level", "", "comma-separated levels")
	jsonOutput := fs.Bool("json", false, "emit JSON lines")
	noColor := fs.Bool("no-color", false, "disable ANSI colors")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("scan expects exactly one file: afx-logwatch scan app.log")
	}

	filter, err := parser.NewFilter(*include, *exclude, splitCSV(*levels))
	if err != nil {
		return err
	}

	f, err := os.Open(fs.Arg(0))
	if err != nil {
		return err
	}
	defer f.Close()

	out := formatter.New(os.Stdout, formatter.Options{JSON: *jsonOutput, Color: !*noColor})
	counter := stats.New()
	return watcher.ScanLines(f, func(line string) error {
		entry := parser.Parse(fs.Arg(0), line)
		if !filter.Match(entry) {
			return nil
		}
		counter.Add(entry.Level)
		return out.Write(entry)
	})
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

func printHelp() {
	fmt.Print(`AFX LogWatch - fast real-time log monitoring in Go

Usage:
  afx-logwatch watch [options] <file> [file...]
  afx-logwatch scan  [options] <file>
  afx-logwatch version

Examples:
  afx-logwatch watch app.log
  afx-logwatch watch --level error,warn app.log
  afx-logwatch watch --include "timeout|failed" --from-start app.log
  afx-logwatch watch --json app.log > filtered.jsonl
  afx-logwatch scan --level error examples/app.log

Use "afx-logwatch watch -h" or "afx-logwatch scan -h" for options.
`)
}
