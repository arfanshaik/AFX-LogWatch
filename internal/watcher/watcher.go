package watcher

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type Event struct {
	Source string
	Line   string
}

type Config struct {
	Path         string
	FromStart    bool
	PollInterval time.Duration
}

func Follow(ctx context.Context, cfg Config, out chan<- Event) error {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 250 * time.Millisecond
	}

	file, err := os.Open(cfg.Path)
	if err != nil {
		return fmt.Errorf("open %s: %w", cfg.Path, err)
	}
	defer file.Close()

	if !cfg.FromStart {
		if _, err := file.Seek(0, io.SeekEnd); err != nil {
			return err
		}
	}

	reader := bufio.NewReader(file)
	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			if line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
				if len(line) > 0 && line[len(line)-1] == '\r' {
					line = line[:len(line)-1]
				}
			}
			select {
			case out <- Event{Source: cfg.Path, Line: line}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if err == nil {
			continue
		}
		if !errors.Is(err, io.EOF) {
			return fmt.Errorf("read %s: %w", cfg.Path, err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			rotated, err := needsReopen(file, cfg.Path)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return err
			}
			if rotated {
				file.Close()
				file, err = os.Open(cfg.Path)
				if err != nil {
					continue
				}
				reader = bufio.NewReader(file)
			} else {
				reader.Reset(file)
			}
		}
	}
}

func needsReopen(file *os.File, path string) (bool, error) {
	current, err := file.Stat()
	if err != nil {
		return false, err
	}
	target, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	pos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	return !os.SameFile(current, target) || target.Size() < pos, nil
}

func ScanLines(r io.Reader, fn func(string) error) error {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		if err := fn(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}
