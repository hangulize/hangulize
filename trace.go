package hangulize

import (
	"fmt"
	"io"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
)

// Trace is emitted when a replacement occurs. It is used for tracing of
// Hangulize pipeline internal.
type Trace struct {
	Step Step
	Why  string
	Word string
}

func (t *Trace) String() string {
	return fmt.Sprintf("[%s] %#v %s", t.Step, t.Word, t.Why)
}

// -----------------------------------------------------------------------------

type Traces []Trace

func (ts Traces) Render(w io.Writer) {
	var width, maxWidth int
	widths := make([]int, len(ts))

	for i, t := range ts {
		width = runewidth.StringWidth(t.Word)
		widths[i] = width

		if maxWidth < width {
			maxWidth = width
		}
	}

	var step Step

	for i, t := range ts {
		if step != t.Step {
			step = t.Step
			fmt.Fprintf(w, "[%s]\n", step)
		}

		fmt.Fprintf(w, "  %s", t.Word)
		fmt.Fprintf(w, strings.Repeat(" ", maxWidth-widths[i]))
		if t.Why != "" {
			fmt.Fprintf(w, " | %s", t.Why)
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

func (tr *tracer) trace(step Step, why, word string) {
	if word == tr.lastWord {
		return
	}
	tr.traces = append(tr.traces, Trace{step, why, word})
	tr.lastWord = word
}

func (tr *tracer) TraceWord(step Step, why, word string) {
	if tr == nil {
		return
	}
	tr.trace(step, why, word)
}

func (tr *tracer) TraceSubwords(step Step, why string, subwords []subword) {
	if tr == nil {
		return
	}
	b := subwordsBuilder{subwords}
	word := b.String()
	word = strings.Replace(word, "\x00", ".", -1)
	tr.trace(step, why, word)
}

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
			rtr.tr.TraceSubwords(step, rule.String(), subwords)
		}
	}
}
