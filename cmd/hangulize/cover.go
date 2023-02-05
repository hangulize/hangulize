package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/pkg/hsl"
)

type cov struct {
	name   string
	step   string
	ruleID int
}

// cover stores covered rule IDs.
type cover struct {
	covered map[cov]bool
	names   map[string]bool
}

// newCover creates a cover.
func newCover() *cover {
	return &cover{make(map[cov]bool), make(map[string]bool)}
}

// Visit marks a visited spec file name.
func (c *cover) Visit(name string) {
	if c == nil {
		return
	}

	c.names[name] = true
}

// Cover marks a covered rule.
func (c *cover) Cover(name, step string, ruleID int) {
	if c == nil {
		return
	}

	c.covered[cov{name, step, ruleID}] = true
	c.Visit(name)
}

// Covered returns true if the rule has been covered within a test.
func (c *cover) Covered(name, step string, ruleID int) bool {
	if c == nil {
		return false
	}

	return c.covered[cov{name, step, ruleID}]
}

// Coverage returns the test coverage.
func (c *cover) Coverage() float64 {
	if c == nil {
		return 0.0
	}

	var total int

	for name := range c.names {
		file, err := os.Open(name)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		spec, err := hangulize.ParseSpec(file)
		if err != nil {
			panic(err)
		}

		total += len(spec.Rewrite)
		total += len(spec.Transcribe)
	}

	return float64(len(c.covered)) / float64(total)
}

// WriteProfile writes the coverage profile.
//
// The coverage profile format is similar to the Go test's one. So we can
// generate an HTML representation of coverage profile by "go tool cover":
//
//	$ hangulize test --coverprofile=cover.txt *.hsl
//	coverage: 100.0% of rules
//	$ go tool cover -html cover.txt -o cover.html
func (c *cover) WriteProfile(w io.Writer) error {
	template := "%s:%d.1,%d.%d 1 %d\n"

	if _, err := io.WriteString(w, "mode: count\n"); err != nil {
		return fmt.Errorf("cover: %w", err)
	}

	for name := range c.names {
		file, err := os.Open(name)
		if err != nil {
			return fmt.Errorf("cover: %w", err)
		}
		defer file.Close()

		// Parse as an HSL to track line numbers.
		hsl, err := hsl.Parse(file)
		if err != nil {
			return fmt.Errorf("cover: %w", err)
		}

		// Parse as a spec to collect all
		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("cover: %w", err)
		}

		spec, _ := hangulize.ParseSpec(file)

		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("cover: %w", err)
		}

		cols := make([]int, 0)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			col := len(scanner.Text()) + 1
			cols = append(cols, col)
		}

		rewritePairs := hsl["rewrite"].Pairs()
		for _, rule := range spec.Rewrite {
			line := rewritePairs[rule.ID].Line()
			col := cols[line-1]
			covered := btoi(c.Covered(name, "Rewrite", rule.ID))
			fmt.Fprintf(w, template, name, line, line, col, covered)
		}

		transcribePairs := hsl["transcribe"].Pairs()
		for _, rule := range spec.Transcribe {
			line := transcribePairs[rule.ID].Line()
			col := cols[line-1]
			covered := btoi(c.Covered(name, "Transcribe", rule.ID))
			fmt.Fprintf(w, template, name, line, line, col, covered)
		}
	}

	return nil
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
