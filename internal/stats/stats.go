package stats

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

type Counter struct {
	mu     sync.Mutex
	total  uint64
	levels map[string]uint64
}

func New() *Counter { return &Counter{levels: make(map[string]uint64)} }

func (c *Counter) Add(level string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.total++
	c.levels[level]++
}

func (c *Counter) WriteSummary(w io.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]string, 0, len(c.levels))
	for k := range c.levels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "[stats] total=%d", c.total)
	for _, k := range keys {
		fmt.Fprintf(w, " %s=%d", k, c.levels[k])
	}
	fmt.Fprintln(w)
}
