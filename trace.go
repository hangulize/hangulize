package hangulize

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
)

// Trace is emitted when a replacement occurs. It is used for tracing of
// Hangulize pipeline internal.
type Trace struct {
	Step Step
	Word string
	Why  string
	Rule *Rule
}

func (t *Trace) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[%s] %#v", t.Step, t.Word)

	if t.Rule != nil {
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
		if t.Why != "" {
			fmt.Fprintf(w, " | (%s)", t.Why)
		} else if t.Rule != nil {
			fmt.Fprintf(w, " | %s", t.Rule)
		}
		fmt.Fprintf(w, "\n")
	}
}

// -----------------------------------------------------------------------------

type tracer struct {
	traces   Traces
	lastWord string
}

func (tr *tracer) Traces() Traces {
	return tr.traces
}

func (tr *tracer) trace(
	step Step, word string,
	why string, rule *Rule,
) {
	if word == tr.lastWord {
		return
	}
	tr.traces = append(tr.traces, Trace{step, word, why, rule})
	tr.lastWord = word
}

func (tr *tracer) TraceWord(
	step Step, word string,
	why string, rule *Rule,
) {
	if tr == nil {
		return
	}
	tr.trace(step, word, why, rule)
}

func (tr *tracer) TraceSubwords(
	step Step, subwords []subword,
	why string, rule *Rule,
) {
	if tr == nil {
		return
	}
	b := subwordsBuilder{subwords}
	word := b.String()
	word = strings.Replace(word, "\x00", ".", -1)
	tr.trace(step, word, why, rule)
}

// -----------------------------------------------------------------------------

type ruleTracer struct {
	tr           *tracer
	subwords     []subword
	trap         map[int][]*string
	rules        map[int]*Rule
	maxRuleIndex int
}

func (tr *tracer) RuleTracer(subwords []subword) *ruleTracer {
	if tr == nil {
		return nil
	}
	return &ruleTracer{
		tr,
		subwords,
		make(map[int][]*string),
		make(map[int]*Rule),
		-1,
	}
}

func (rtr *ruleTracer) Trace(
	ruleIndex int, rule *Rule,
	swIndex int, word string,
) {
	if rtr == nil {
		return
	}
	if rtr.trap[ruleIndex] == nil {
		rtr.trap[ruleIndex] = make([]*string, len(rtr.subwords))
	}
	rtr.trap[ruleIndex][swIndex] = &word

	rtr.rules[ruleIndex] = rule
	rtr.maxRuleIndex = ruleIndex
}

func (rtr *ruleTracer) Commit(step Step) {
	if rtr == nil {
		return
	}
	subwords := make([]subword, len(rtr.subwords))
	copy(subwords, rtr.subwords)

	for ruleIndex := 0; ruleIndex <= rtr.maxRuleIndex; ruleIndex++ {
		dirty := false

		rule := rtr.rules[ruleIndex]
		words := rtr.trap[ruleIndex]

		for swIndex, word := range words {
			if word == nil {
				continue
			}
			subwords[swIndex] = subword{*word, 0}
			dirty = true
		}

		if dirty {
			rtr.tr.TraceSubwords(step, subwords, "", rule)
		}
	}
}
