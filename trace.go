package hangulize

import (
	"strings"

	"github.com/hangulize/hangulize/internal/subword"
)

type Trace struct {
	Step string
	Word string
	Why  string
	Rule *Rule
}

type tracer struct {
	fn       func(Trace)
	prevWord string
}

func newTracer(fn func(Trace)) *tracer {
	if fn == nil {
		return nil
	}
	return &tracer{fn, ""}
}

// Trace emits a tracing event if it is necessary.
func (r *tracer) trace(t Trace) {
	if r == nil {
		return
	}

	if r.prevWord == t.Word {
		return
	}
	r.prevWord = t.Word
	r.fn(t)
}

func (r *tracer) traceSubwords(step string, subwords []subword.Subword) *subwordsTracer {
	if r == nil {
		return nil
	}
	return &subwordsTracer{step: step, tracer: r, subwords: subwords}
}

func (r *tracer) Input(word string) {
	r.trace(Trace{Step: "Input", Word: word})
}

func (r *tracer) Transliterate(word, scheme string) {
	r.trace(Trace{Step: "Transliterate", Word: word, Why: scheme})
}

func (r *tracer) Normalize(word, script string) {
	r.trace(Trace{Step: "Normalize", Word: word, Why: script})
}

func (r *tracer) Rewrite(subwords []subword.Subword) (func(int, string, Rule), func()) {
	st := r.traceSubwords("Rewrite", subwords)
	return st.RecordSubword, st.Commit
}

func (r *tracer) Transcribe(subwords []subword.Subword) (func(int, string, Rule), func()) {
	st := r.traceSubwords("Transcribe", subwords)
	return st.RecordSubword, st.Commit
}

func (r *tracer) Syllabify(word string) {
	r.trace(Trace{Step: "Syllabify", Word: word})
}

func (r *tracer) Localize(word, script string) {
	r.trace(Trace{Step: "Localize", Word: word, Why: script})
}

type subwordsTracer struct {
	step     string
	tracer   *tracer
	subwords []subword.Subword
	records  []*subwordsTracerRecord
}

type subwordsTracerRecord struct {
	rule     Rule
	subwords []*string // [i] is the changed i-th subword. nil means that this subword has not been changed.
}

func growUp[T any](slice []T, size int) []T {
	return append(slice, make([]T, size-len(slice))...)
}

func (r *subwordsTracer) RecordSubword(i int, word string, rule Rule) {
	if r == nil {
		return
	}

	r.records = growUp(r.records, rule.ID+1)
	if r.records[rule.ID] == nil {
		r.records[rule.ID] = &subwordsTracerRecord{rule, nil}
	}
	r.records[rule.ID].subwords = growUp(r.records[rule.ID].subwords, i+1)

	r.records[rule.ID].subwords[i] = &word
}

func (r *subwordsTracer) Commit() {
	if r == nil {
		return
	}

	subwords := append([]subword.Subword{}, r.subwords...)

	for _, rec := range r.records {
		dirty := false

		for i, word := range rec.subwords {
			if word == nil {
				continue
			}
			subwords[i] = subword.New(*word, 0)
			dirty = true
		}

		if dirty {
			b := subword.NewBuilder(subwords)
			word := b.String()
			word = strings.Replace(word, "\x00", ".", -1)

			r.tracer.trace(Trace{Step: r.step, Word: word, Rule: &rec.rule})
		}
	}
}
