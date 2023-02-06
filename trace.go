package hangulize

import (
	"strings"

	"github.com/hangulize/hangulize/internal/subword"
)

// Trace is a tracing event which the Hangulize procedure emits.
type Trace struct {
	Step string
	Word string
	Why  string
	Rule *Rule
}

// tracer traces each step in the Hangulize procedure.
type tracer struct {
	fn       func(Trace)
	prevWord string
}

// newTracer creates a tracer or nil depending on the given function.
func newTracer(fn func(Trace)) *tracer {
	if fn == nil {
		return nil
	}
	return &tracer{fn, ""}
}

// trace emits a tracing event if it is necessary.
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

// Input traces an "Input" step.
func (r *tracer) Input(word string) {
	r.trace(Trace{Step: "Input", Word: word})
}

// Transliterate traces a "Transliterate" step.
func (r *tracer) Transliterate(word, scheme string) {
	r.trace(Trace{Step: "Transliterate", Word: word, Why: scheme})
}

// Normalize traces a "Normalize" step.
func (r *tracer) Normalize(word, script string) {
	r.trace(Trace{Step: "Normalize", Word: word, Why: script})
}

// Rewrite creates a subwords tracer for a "Rewrite" step. Call the first
// return value to record a subword then call the second return value to merge
// and commit them.
//
//	recordSubword, commit := tracer.Rwrite(subwords)
//	defer commit()
//
//	recordSubword(0, "1st subword", rule)
//	recordSubword(1, "2nd subword", rule)
func (r *tracer) Rewrite(subwords []subword.Subword) (func(int, string, Rule), func()) {
	st := newSubwordsTracer(r, "Rewrite", subwords)
	return st.RecordSubword, st.Commit
}

// Transcribe creates a subwords tracer for a "Transcribe" step. Call the first
// return value to record a subword then call the second return value to merge
// and commit them.
//
//	recordSubword, commit := tracer.Transcribe(subwords)
//	defer commit()
//
//	recordSubword(0, "1st subword", rule)
//	recordSubword(1, "2nd subword", rule)
func (r *tracer) Transcribe(subwords []subword.Subword) (func(int, string, Rule), func()) {
	st := newSubwordsTracer(r, "Transcribe", subwords)
	return st.RecordSubword, st.Commit
}

// Syllabify traces a "Syllabify" step.
func (r *tracer) Syllabify(word string) {
	r.trace(Trace{Step: "Syllabify", Word: word})
}

// Localize traces a "Localize" step.
func (r *tracer) Localize(word, script string) {
	r.trace(Trace{Step: "Localize", Word: word, Why: script})
}

// subwordsTracer traces a step by recording updated subwords.
type subwordsTracer struct {
	tracer   *tracer
	step     string
	subwords []subword.Subword
	records  []*subwordsTracerRecord
}

// newSubwordsTracer creates a subwordsTracer.
func newSubwordsTracer(tracer *tracer, step string, subwords []subword.Subword) *subwordsTracer {
	return &subwordsTracer{tracer: tracer, step: step, subwords: subwords}
}

// subwordsTracerRecord records modified subwords by a rule.
type subwordsTracerRecord struct {
	rule Rule

	// subwords[i] is the i-th subword which has been modified. nil means that
	// this subword has not been changed.
	subwords []*string
}

// RecordSubword records an i-th subword modified by a rule.
func (r *subwordsTracer) RecordSubword(i int, word string, rule Rule) {
	if r.tracer == nil {
		return
	}

	if len(r.records) < rule.ID+1 {
		more := make([]*subwordsTracerRecord, rule.ID+1-len(r.records))
		r.records = append(r.records, more...)
	}

	if r.records[rule.ID] == nil {
		r.records[rule.ID] = &subwordsTracerRecord{rule, make([]*string, len(r.subwords))}
	}

	r.records[rule.ID].subwords[i] = &word
}

// Commit merges the recorded subwords by the rule then traces them by the
// underlying tracer.
func (r *subwordsTracer) Commit() {
	if r.tracer == nil {
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
