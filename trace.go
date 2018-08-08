package hangulize

import (
	"fmt"
	"strings"
)

// Trace is emitted when a replacement occurs. It is used for tracing of
// Hangulize pipeline internal.
type Trace struct {
	Step string
	Why  string
	Word string
}

func (t *Trace) String() string {
	return fmt.Sprintf("[%s] %#v %s", t.Step, t.Word, t.Why)
}

type tracer struct {
	traces   []Trace
	lastWord string
}

func (tr *tracer) Traces() []Trace {
	return tr.traces
}

func (tr *tracer) trace(step, why, word string) {
	if word == tr.lastWord {
		return
	}
	tr.traces = append(tr.traces, Trace{step, why, word})
	tr.lastWord = word
}

func (tr *tracer) TraceWord(step, why, word string) {
	if tr == nil {
		return
	}
	tr.trace(step, why, word)
}

func (tr *tracer) TraceSubwords(step, why string, subwords []subword) {
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
		0,
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

func (rtr *ruleTracer) Commit(step string) {
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
