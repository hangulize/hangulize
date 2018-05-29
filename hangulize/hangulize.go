/*
Package hangulize provides an automatic Hangul transcriber for non-Korean
words.  Transcription means that the systematic representation of language in
written form.

	https://en.wikipedia.org/wiki/Transcription_(linguistics)

Originally, the Hangulize was invented with Python in 2010
(https://github.com/sublee/hangulize).  It has been provided at
http://hangulize.org/ for Korean translators.  Brian Jongseong Park proposed
the seed idea of the Hangulize on his Blog.

	http://iceager.egloos.com/2610028

This Go re-implementation will be a reboot of the Hangulize with attractive
feature improvements.

Pipeline

When we transcribe, the word goes through Hangulize's procedural pipeline.  The
pipeline has 5 steps: "Normalize", "Group", "Rewrite", "Transcribe", and
"Compose Hangul".  If we can transcribe "Hello!" in English into "헬로!"
(actually, English is not supported yet), the pipeline world work like:

	0. Input       "Hello!"
	1. Normalize   "hello!"
	2. Group       "hello", "!"
	3. Rewrite     "heˈlō", "!"
	4. Transcribe  "ㅎㅔ-ㄹㄹㅗ", "!"
	5. Compose H.  "헬로!"

The "1. Normalize" step eliminates letter case from the word to make the next
steps work easier.  The "2. Group" step groups letters by their meaningness
into subwords.  Meaningful letter is the letter which appears in the
rewrite/transcribe rules.  The "3. Rewrite" step minimizes the gap between
pronunciation and spelling.  The "4. Transcribe" step determines Hangul
spelling for the pronunciation.  Finally, the  "5. Compose Hangul" step
converts decomposed Jamo phonemes to composed Hangul syllables.

Spec

A spec is written by the HGL format which is a configuration DSL for Hangulize
2.  One spec is for one language transcription system.  So we need to describe
about the language at the first:

	lang:
	    id      = "ita"
	    codes   = "it", "ita" # ISO 639-1 and 3 codes
	    english = "Italian"
	    korean  = "이탈리아어"
	    script  = "roman"

Then write about yourself and the stage of this spec:

	config:
	    author = "John Doe <john@example.com>"
	    stage  = "draft"

We will write many patterns in rewrite/transcribe rules soon.  Some expressions
may appear many times annoyingly.  To not repeat ourselves, we can use
variables and macros.

A variable is a combination of letters.  Variable in pattern will match with
one of the letters.  Variable "foo" can be referenced with "<foo>" in the
patterns.

	vars:
	    "vowels" = "a", "e", "i", "o", "u"

A macro expression is replaced with the target before parsing the patterns. "@"
is the common macro for "<vowels>" variable:

	macros:
	    "@" = "<vowels>"

Now we can write "rewrite" rules.  There are Pattern and RPattern.  Pattern
matches with letters in a word.  RPattern represents how the matched letters
should be replaced.  A replaced word by a rule would become as the input for
the next rule:

	rewrite:
	    "^gli$"   -> "li"
	    "^gli{@}" -> "li"
	    "{@}gli"  -> "li"
	    "gn{@}"   -> "nJ"

Pattern is based on Regular Expression but it has it's own custom syntax.  We
call it "HRE" which means "Hangulize-specific Regular Expression".  For the
detail, see the documentation of Pattern.

"transcribe" rules are exactly same with "rewrite" rules.  But it's RPatterns
represent Hangul Jamo phonemes.  In contrast to "rewrite", a replaced word
won't become as the input for the next rules:

	transcribe:
	    "b" -> "ㅂ"
	    "d" -> "ㄷ"
	    "f" -> "ㅍ"
	    "g" -> "ㄱ"

Finally, we should write expected transcription examples.  They are used for
unit testing.  Verify your spec yourself:

	test:
	    "allegretto" -> "알레그레토"
	    "gita"       -> "지타"
	    "bisnonno"   -> "비스논노"
	    "Pinocchio"  -> "피노키오"

*/
package hangulize

import (
	"strings"
)

// Hangulize is the most simple and useful API of thie package.  It transcribes
// a non-Korean word into Hangul, which is the Korean alphabet.  For example,
// it will transcribe "Владивосто́к" in Russian into "블라디보스토크".
func Hangulize(lang string, word string) string {
	spec, ok := LoadSpec(lang)
	if !ok {
		// spec not found
		return word
	}

	h := NewHangulizer(spec)
	return h.Hangulize(word)
}

// -----------------------------------------------------------------------------

// Hangulizer provides the transcription logic for the underlying spec.
type Hangulizer struct {
	spec *Spec
}

// NewHangulizer creates a Hangulizer for a spec.
func NewHangulizer(spec *Spec) *Hangulizer {
	return &Hangulizer{spec}
}

// Spec returns the underlying spec.
func (h *Hangulizer) Spec() *Spec {
	return h.spec
}

// Hangulize transcribes a loanword into Hangul.
func (h *Hangulizer) Hangulize(word string) string {
	p := pipeline{h, nil}
	return p.forward(word)
}

// HangulizeTrace transcribes a loanword into Hangul
// and returns the traced internal events too.
func (h *Hangulizer) HangulizeTrace(word string) (string, []Trace) {
	var tr tracer
	p := pipeline{h, &tr}

	word = p.forward(word)

	return word, tr.Traces()
}

// -----------------------------------------------------------------------------
// The Hangulize Pipeline

type pipeline struct {
	h  *Hangulizer
	tr *tracer
}

// forward runs the Hangulize pipeline for a word.
func (p *pipeline) forward(word string) string {
	p.input(word)
	word = p.normalize(word)
	subwords := p.group(word)
	subwords = p.rewrite(subwords)
	subwords = p.transcribe(subwords)
	word = p.composeHangul(subwords)
	return word
}

// 0. Just recording beginning (Word)
//
func (p *pipeline) input(word string) {
	p.tr.TraceWord("input", "", word)
}

// 1. Normalize (Word -> Word)
//
// This step eliminates letter case to make the next steps work easier.
//
// For example, "Hello" in Roman script will be normalized to "hello".
//
func (p *pipeline) normalize(word string) string {
	word = p.h.spec.normReplacer.Replace(word)

	p.tr.TraceWord("normalize", "custom", word)

	norm := p.h.spec.norm
	except := p.h.spec.normLetters
	word = Normalize(word, norm, except)

	p.tr.TraceWord("normalize", p.h.spec.Lang.Script, word)

	return word
}

// 2. Group meaningful letters (Word -> Subwords[level=0 or 1])
//
// Meaningful letter is the letter which appears in the rewrite/transcribe
// rules.  This step groups letters by their meaningness into subwords.  A
// meaningful subword has level=1 meanwhile meaningless one has level=0.
//
// For example, "hello, world!" will be grouped into
// [{"hello",1}, {", ",0}, {"world",1}, {"!",0}].
//
func (p *pipeline) group(word string) []subword {
	rep := newSubwordReplacer(word, 0, 1)

	for i, ch := range word {
		let := string(ch)
		if inSet(let, p.h.spec.letters) || isSpace(let) {
			rep.Replace(i, i+len(let), let)
		}
	}

	return rep.Subwords()
}

// 3. Rewrite (Subwords -> Subwords[level=1])
//
// This step minimizes the gap between pronunciation and spelling.
//
// It replaces the word by a rule.  The replaced word will be used as the input
// for the next rule.  Each result becomes the next input.  That's why this
// step called "rewrite".
//
// For example, "hello" can be rewritten to "heˈlō".
//
// NOTE(sublee): But this step has a limitation.  It guesses a pronunciation
// from the spelling.  But it can be too hard for some script systems, such as
// English or Franch.
//
func (p *pipeline) rewrite(subwords []subword) []subword {
	var swBuf subwordsBuilder

	for _, sw := range subwords {
		word := sw.word
		level := sw.level

		rep := newSubwordReplacer(word, level, 1)

		for _, rule := range p.h.spec.Rewrite {
			repls := rule.replacements(word)
			rep.ReplaceBy(repls...)
			word = rep.String()
		}

		swBuf.Append(rep.Subwords()...)
	}

	subwords = swBuf.Subwords()

	// TODO(sublee): per-rule tracing
	p.tr.TraceSubwords("rewrite", "", subwords)

	return subwords
}

// 4. Transcribe (Subwords -> Subwords[level=2])
//
// This step determines Hangul spelling for the pronunciation.
//
// Rather than generating composed Hangul syllables, it enumerates decomposed
// Jamo phonemes, such as "ㅎㅏ-ㄴ".  In this form, a Jaeum after a hyphen
// ("-ㄴ") means that it is a Jongseong (tail).
//
// For example, "heˈlō" can be transcribed as "ㅎㅔ-ㄹㄹㅗ".
//
func (p *pipeline) transcribe(subwords []subword) []subword {
	var swBuf subwordsBuilder

	for _, sw := range subwords {
		word := sw.word
		level := sw.level

		rep := newSubwordReplacer(word, level, 2)

		// transcribe is not rewrite.  A result of a replacement is not the
		// input of the next replacement.  dummy marks the replaced subwords
		// with NULL characters.
		dummy := newSubwordReplacer(word, 0, 0)

		for _, rule := range p.h.spec.Transcribe {
			repls := rule.replacements(word)
			rep.ReplaceBy(repls...)

			for _, repl := range repls {
				nulls := strings.Repeat("\x00", len(repl.word))
				dummy.Replace(repl.start, repl.stop, nulls)
			}

			rep.flush()
			word = dummy.String()
		}

		swBuf.Append(rep.Subwords()...)
	}

	// Discard level=1 subwords.  They have been generated by "3. Rewrite" but
	// never transcribed.  They are superfluity of the internal behavior.
	subwords = swBuf.Subwords()
	swBuf.Reset()

	for _, sw := range subwords {
		if sw.level == 1 {
			if hasSpace(sw.word) {
				swBuf.Append(subword{" ", 1})
			}
			continue
		}
		swBuf.Append(sw)
	}

	subwords = swBuf.Subwords()

	// TODO(sublee): per-rule tracing
	p.tr.TraceSubwords("transcribe", "", subwords)

	return subwords
}

// 5. Compose Hangul (Subwords -> Word)
//
// This step converts decomposed Jamo phonemes to composed Hangul syllables.
//
// For example, "ㅎㅔ-ㄹㄹㅗ" will be "헬로".
//
func (p *pipeline) composeHangul(subwords []subword) string {
	var buf strings.Builder
	var jamoBuf strings.Builder

	for _, sw := range subwords {
		// Don't touch level=0 subwords.  They just have passed through the
		// pipeline, because they are meaningless.
		if sw.level == 0 {
			buf.WriteString(ComposeHangul(jamoBuf.String()))
			jamoBuf.Reset()

			buf.WriteString(sw.word)
			continue
		}
		jamoBuf.WriteString(sw.word)
	}
	buf.WriteString(ComposeHangul(jamoBuf.String()))

	word := buf.String()

	p.tr.TraceWord("compose hangul", "", word)

	return word
}
