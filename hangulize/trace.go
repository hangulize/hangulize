package hangulize

import "fmt"

// Trace is emitted when a replacement occurs.  It is used for tracing of
// Hangulize pipeline internal.
type Trace struct {
	step string
	why  string
	from string
	to   string
}

func (t *Trace) String() string {
	return fmt.Sprintf("[%s] %#v %s", t.step, t.to, t.why)
}

type Tracer struct {
	traces []Trace
}

func (tr *Tracer) Traces() []Trace {
	return tr.traces
}

func (tr *Tracer) wordToWord(step, why string, from, to string) {
	if from == to {
		return
	}
	tr.traces = append(tr.traces, Trace{step, why, from, to})
}

func (tr *Tracer) WordToWord(step, why string, from, to string) {
	if tr == nil {
		return
	}
	tr.wordToWord(step, why, from, to)
}

func (tr *Tracer) WordToSubwords(step, why string, from string, to []Subword) {
	if tr == nil {
		return
	}
	toWord := NewSubwordsBuilder(to).String()
	tr.wordToWord(step, why, from, toWord)
}

func (tr *Tracer) SubwordsToWord(step, why string, from []Subword, to string) {
	if tr == nil {
		return
	}
	fromWord := NewSubwordsBuilder(from).String()
	tr.wordToWord(step, why, fromWord, to)
}

func (tr *Tracer) SubwordsToSubwords(step, why string, from, to []Subword) {
	if tr == nil {
		return
	}
	fromWord := NewSubwordsBuilder(from).String()
	toWord := NewSubwordsBuilder(to).String()
	tr.wordToWord(step, why, fromWord, toWord)
}
