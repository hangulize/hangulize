package hangulize

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	runewidth "github.com/mattn/go-runewidth"

	"github.com/hangulize/hangulize/internal/subword"
)

// Trace is emitted when a replacement occurs. It is used for tracing of
// Hangulize pipeline internal.
type Trace struct {
	Step Step
	Word string

	Why string

	Rule    Rule
	HasRule bool
}

func newTrace(step Step, word string, why string, rule *Rule) Trace {
	// It can hold either why or rule.
	var _rule Rule
	var hasRule bool

	if rule != nil {
		_rule = *rule
		hasRule = true
	}

	return Trace{step, word, why, _rule, hasRule}
}

func (t Trace) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[%s] %#v", t.Step, t.Word)

	if t.HasRule {
		fmt.Fprintf(&buf, " | %s", t.Rule)
	} else if t.Why != "" {
		fmt.Fprintf(&buf, " | (%s)", t.Why)
	}

	return buf.String()
}

// -----------------------------------------------------------------------------

// Traces is an array of Trace.
type Traces []Trace

// Render generates a report text.
func (ts Traces) Render(w io.Writer) {
	// Detect the max rune width of the words.
	var width, maxWidth int
	widths := make([]int, len(ts))

	for i, t := range ts {
		width = runewidth.StringWidth(t.Word)
		widths[i] = width

		if maxWidth < width {
			maxWidth = width
		}
	}

	// Render the report.
	var step Step

	for i, t := range ts {
		if step != t.Step {
			step = t.Step
			fmt.Fprintf(w, "[%s]\n", step)
		}

		fmt.Fprintf(w, "  %s", t.Word)
		fmt.Fprintf(w, strings.Repeat(" ", maxWidth-widths[i]))

		if t.HasRule {
			fmt.Fprintf(w, " | %s", t.Rule)
		} else if t.Why != "" {
			fmt.Fprintf(w, " | (%s)", t.Why)
		}

		fmt.Fprintf(w, "\n")
	}
}

// -----------------------------------------------------------------------------

// tracer collects tracing logs.
//
// HangulizeTrace uses it.
//
type tracer struct {
	traces   Traces
	lastWord string
}

// Traces returns the collected tracing logs.
func (tr *tracer) Traces() Traces {
	return tr.traces
}

// trace records a tracing log.
//
// Public tracing methods are implemented by this method.
//
func (tr *tracer) trace(step Step, word string, why string, rule *Rule) {
	if tr == nil {
		return
	}

	if word == tr.lastWord {
		// No changes.
		return
	}

	trace := newTrace(step, word, why, rule)
	tr.traces = append(tr.traces, trace)
	tr.lastWord = word
}

// Trace records a tracing log with word and why.
func (tr *tracer) Trace(step Step, word, why string) {
	tr.trace(step, word, why, nil)
}

// -----------------------------------------------------------------------------

// subwordsTracer collects tracing logs for the specific subwords.
type subwordsTracer struct {
	tr        *tracer
	step      Step
	subwords  []subword.Subword
	trap      map[int][]*string
	rules     map[int]Rule
	maxRuleID int
}

// SubwordsTracer creates a subwordsTracer under the tracer.
func (tr *tracer) SubwordsTracer(
	step Step,
	subwords []subword.Subword,
) *subwordsTracer {
	if tr == nil {
		return nil
	}
	return &subwordsTracer{
		tr,
		step,
		subwords,
		make(map[int][]*string),
		make(map[int]Rule),
		-1,
	}
}

// Trace records a tracing log for a subword.
//
// It buffers the tracing logs. Call Commit to flush them.
//
func (swtr *subwordsTracer) Trace(swIndex int, word string, rule Rule) {
	if swtr == nil {
		return
	}

	if swtr.trap[rule.ID] == nil {
		swtr.trap[rule.ID] = make([]*string, len(swtr.subwords))
	}

	swtr.trap[rule.ID][swIndex] = &word
	swtr.rules[rule.ID] = rule
	swtr.maxRuleID = rule.ID
}

// Commit flushes the buffered tracing logs.
func (swtr *subwordsTracer) Commit() {
	if swtr == nil {
		return
	}

	subwords := make([]subword.Subword, len(swtr.subwords))
	copy(subwords, swtr.subwords)

	var (
		dirty bool
		rule  Rule
		words []*string
	)

	for ruleID := 0; ruleID <= swtr.maxRuleID; ruleID++ {
		dirty = false
		rule = swtr.rules[ruleID]
		words = swtr.trap[ruleID]

		for swIndex, word := range words {
			if word == nil {
				continue
			}
			subwords[swIndex] = subword.New(*word, 0)
			dirty = true
		}

		if dirty {
			b := subword.NewBuilder(subwords)
			word := b.String()
			word = strings.Replace(word, "\x00", ".", -1)

			swtr.tr.trace(swtr.step, word, "", &rule)
		}
	}
}
