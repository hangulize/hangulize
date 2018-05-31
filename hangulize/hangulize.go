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
