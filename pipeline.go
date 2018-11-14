package hangulize

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"

	"github.com/hangulize/hangulize/internal/jamo"
	"github.com/hangulize/hangulize/internal/subword"
)

var reSpace = regexp.MustCompile(`\s`)

func isSpace(word string) bool {
	return strings.TrimSpace(word) == ""
}

func hasSpace(word string) bool {
	return reSpace.MatchString(word)
}

// -----------------------------------------------------------------------------

// Step is an identifier for the each pipeline step.
type Step int

const (
	_ Step = iota * 10

	// Input step just records the beginning.
	Input

	// Phonemize step converts the spelling to the phonograms.
	Phonemize

	// Normalize step eliminates letter case to make the next steps work easier.
	Normalize

	// Group step associates meaningful letters.
	Group

	// Rewrite step minimizes the gap between pronunciation and spelling.
	Rewrite

	// Transcribe step determines Hangul spelling for the pronunciation.
	Transcribe

	// Syllabify step composes Jamo phonemes into Hangul syllabic blocks.
	Syllabify

	// Transliterate step converts foreign punctuations to fit in Korean.
	Transliterate
)

// AllSteps is the array of all steps.
var AllSteps = []Step{
	Input,
	Phonemize,
	Normalize,
	Group,
	Rewrite,
	Transcribe,
	Syllabify,
	Transliterate,
}

func (s Step) String() string {
	return map[Step]string{
		Input:         "Input",
		Phonemize:     "Phonemize",
		Normalize:     "Normalize",
		Group:         "Group",
		Rewrite:       "Rewrite",
		Transcribe:    "Transcribe",
		Syllabify:     "Syllabify",
		Transliterate: "Transliterate",
	}[s]
}

// -----------------------------------------------------------------------------

type pipeline struct {
	h  *Hangulizer
	tr *tracer
}

// forward runs the Hangulize pipeline for a word.
func (p pipeline) forward(word string) string {
	p.input(word)

	// preparing phase
	word, _ = p.phonemize(word)
	word = p.normalize(word)

	// transcribing phase
	subwords := p.group(word)
	subwords = p.rewrite(subwords)
	subwords = p.transcribe(subwords)

	// finalizing phase
	word = p.syllabify(subwords)
	word = p.transliterate(word)

	return word
}

// -----------------------------------------------------------------------------

// 0. Just recording beginning (Word)
//
func (p pipeline) input(word string) {
	p.tr.Trace(Input, word, "")
}

// 1. Phonemize (Word -> Word)
//
// This step converts the spelling to the phonograms, usually based on lexical
// analysis. Most languages already use phonograms which are sufficient to
// represent the exact pronunciation. But in some languages, such as American
// English or Chinese, it's not true.
//
func (p pipeline) phonemize(word string) (string, bool) {
	id := p.h.spec.Lang.Phonemizer
	if id == "" {
		// The language doesn't require a phonemizer. It's okay.
		return word, true
	}

	pron, ok := p.h.GetPhonemizer(id)
	if ok {
		goto PhonemizerFound
	}

	// Fallback by the global phonemizer registry.
	pron, ok = GetPhonemizer(id)
	if ok {
		goto PhonemizerFound
	}

	// The language requires a phonemizer but not imported yet.
	return word, false

PhonemizerFound:
	word = pron.Phonemize(word)

	p.tr.Trace(Phonemize, word, id)

	return word, true
}

// 2. Normalize (Word -> Word)
//
// This step eliminates letter case to make the next steps work easier.
//
// For example, "Hello" in Latin script will be normalized to "hello".
//
func (p pipeline) normalize(word string) string {
	// Per-spec normalization.
	word = p.h.spec.normReplacer.Replace(word)

	p.tr.Trace(Normalize, word, "")

	// Per-script normalization.
	script := p.h.spec.script
	except := p.h.spec.normLetters

	var buf bytes.Buffer

	for _, ch := range word {
		if except.HasRune(ch) || !script.Is(ch) {
			buf.WriteRune(ch)
		} else {
			buf.WriteRune(script.Normalize(ch))
		}
	}

	word = buf.String()

	p.tr.Trace(Normalize, word, p.h.spec.Lang.Script)

	return word
}

// 3. Group meaningful letters (Word -> Subwords[level=0 or 1])
//
// Meaningful letter is the letter which appears in the rewrite/transcribe
// rules. This step groups letters by their meaningness into subwords. A
// meaningful subword has level=1 meanwhile meaningless one has level=0.
//
// For example, "hello, world!" will be grouped into
// [{"hello",1}, {", ",0}, {"world",1}, {"!",0}].
//
func (p pipeline) group(word string) []subword.Subword {
	rep := subword.NewReplacer(word, 0, 1)

	for i, ch := range word {
		let := string(ch)

		switch {
		case p.h.spec.script.Is(ch):
			fallthrough
		case p.h.spec.puncts.HasRune(ch):
			fallthrough
		case isSpace(let):
			rep.Replace(i, i+len(let), let)
		}
	}

	return rep.Subwords()
}

// 4. Rewrite (Subwords -> Subwords[level=1])
//
// This step minimizes the gap between pronunciation and spelling.
//
// It replaces the word by a rule. The replaced word will be used as the input
// for the next rule. Each result becomes the next input. That's why this step
// called "rewrite".
//
// For example, "hello" can be rewritten to "heˈlō".
//
func (p pipeline) rewrite(subwords []subword.Subword) []subword.Subword {
	var swBuf subword.Builder

	swtr := p.tr.SubwordsTracer(Rewrite, subwords)

	for i, sw := range subwords {
		word := sw.Word
		level := sw.Level

		rep := subword.NewReplacer(word, level, 1)

		for _, rule := range p.h.spec.Rewrite {
			repls := rule.replacements(word)
			rep.ReplaceBy(repls...)
			word = rep.String()

			swtr.Trace(i, word, rule)
		}

		swBuf.Write(rep.Subwords()...)
	}

	subwords = swBuf.Subwords()

	swtr.Commit()

	return subwords
}

// 5. Transcribe (Subwords -> Subwords[level=2])
//
// This step determines Hangul spelling for the pronunciation.
//
// Rather than generating composed Hangul syllables, it enumerates decomposed
// Jamo phonemes, such as "ㅎㅏ-ㄴ". In this form, a Jaeum after a hyphen
// ("-ㄴ") means that it is a Jongseong (tail).
//
// For example, "heˈlō" can be transcribed as "ㅎㅔ-ㄹㄹㅗ".
//
func (p pipeline) transcribe(subwords []subword.Subword) []subword.Subword {
	var swBuf subword.Builder

	swtr := p.tr.SubwordsTracer(Transcribe, subwords)

	for i, sw := range subwords {
		if sw.Level == 0 {
			swBuf.Write(sw)
			continue
		}

		word := sw.Word
		level := sw.Level

		rep := subword.NewReplacer(word, level, 2)

		// transcribe is not rewrite. A result of a replacement is not the
		// input of the next replacement. dummy marks the replaced subwords
		// with NULL characters.
		dummy := subword.NewReplacer(word, 0, 0)

		for _, rule := range p.h.spec.Transcribe {
			repls := rule.replacements(word)
			rep.ReplaceBy(repls...)

			for _, repl := range repls {
				nulls := strings.Repeat("\x00", len(repl.Word))
				dummy.Replace(repl.Start, repl.Stop, nulls)
			}

			word = dummy.String()
			swtr.Trace(i, rep.String(), rule)
		}

		swBuf.Write(rep.Subwords()...)
	}

	// Discard level=1 subwords. They have been generated by "3. Rewrite" but
	// never transcribed. They are superfluity of the internal behavior.
	subwords = swBuf.Subwords()
	swBuf.Reset()

	for _, sw := range subwords {
		if sw.Level == 1 {
			if hasSpace(sw.Word) {
				swBuf.Write(subword.New(" ", 1))
			}
			continue
		}
		swBuf.Write(sw)
	}

	subwords = swBuf.Subwords()

	swtr.Commit()

	return subwords
}

// 6. Syllabify (Subwords -> Word)
//
// This step converts decomposed Jamo phonemes to composed Hangul syllables.
//
// For example, "ㅎㅔ-ㄹㄹㅗ" will be "헬로".
//
func (p pipeline) syllabify(subwords []subword.Subword) string {
	var buf bytes.Buffer
	var jamoBuf bytes.Buffer

	for _, sw := range subwords {
		// Don't touch level=0 subwords. They just have passed through the
		// pipeline, because they are meaningless.
		if sw.Level == 0 {
			buf.WriteString(jamo.ComposeHangul(jamoBuf.String()))
			jamoBuf.Reset()

			buf.WriteString(sw.Word)
			continue
		}
		jamoBuf.WriteString(sw.Word)
	}
	buf.WriteString(jamo.ComposeHangul(jamoBuf.String()))

	word := buf.String()

	p.tr.Trace(Syllabify, word, "")

	return word
}

// 7. Transliterate (Word -> Word)
//
// Finally, this step converts foreign punctuations to fit in Korean.
//
// Korean has adapted the European punctuations. Those are the most common in
// the world. But a few languages, such as Japanese or Chinese, use different
// punctuations with Korean. This step will reduce that kind of culture gap.
//
// For example, "「...」" will be "'...'".
//
func (p pipeline) transliterate(word string) string {
	script := p.h.spec.script

	chars := []rune(word)
	last := len(chars) - 1

	// Pre-evaluate punct or space classification.
	isPunct := make(map[int]bool)
	isSpace := make(map[int]bool)
	for i, ch := range chars {
		isPunct[i] = unicode.IsPunct(ch)
		isSpace[i] = unicode.IsSpace(ch)
	}
	isSpace[-1] = true
	isSpace[last+1] = true

	var buf bytes.Buffer

	for i, ch := range chars {
		// Skip ZWSP.
		if ch == '\u200B' {
			continue
		}

		if !isPunct[i] {
			buf.WriteRune(ch)
			continue
		}

		punct := script.TransliteratePunct(ch)

		// Trim left after punct or space.
		l := i - 1
		if isPunct[l] || isSpace[l] {
			punct = strings.TrimLeftFunc(punct, unicode.IsSpace)
		}

		// Trim right before punct or space.
		r := i + 1
		if isPunct[r] || isSpace[r] {
			punct = strings.TrimRightFunc(punct, unicode.IsSpace)
		}

		buf.WriteString(punct)
	}

	word = buf.String()

	p.tr.Trace(Transliterate, word, p.h.spec.Lang.Script)

	return word
}
