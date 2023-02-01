package hangulize

import (
	"bytes"
	"fmt"
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

// Step is an identifier for the each step in the Hangulize procedure.
type Step int

const (
	_ Step = iota * 10

	// Input step just records the beginning.
	Input

	// Transliterate step converts the spelling to the phonograms.
	Transliterate

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

	// Localize step converts foreign punctuations to fit in Korean.
	Localize
)

// AllSteps is the array of all steps.
var AllSteps = []Step{
	Input,
	Transliterate,
	Normalize,
	Group,
	Rewrite,
	Transcribe,
	Syllabify,
	Localize,
}

func (s Step) String() string {
	return map[Step]string{
		Input:         "Input",
		Transliterate: "Transliterate",
		Normalize:     "Normalize",
		Group:         "Group",
		Rewrite:       "Rewrite",
		Transcribe:    "Transcribe",
		Syllabify:     "Syllabify",
		Localize:      "Localize",
	}[s]
}

// -----------------------------------------------------------------------------

type procedure struct {
	h  Hangulizer
	tr *tracer
}

// forward runs the Hangulize procedure for a word.
func (p procedure) forward(word string) (string, error) {
	p.input(word)

	// preparing phase
	word, err := p.transliterate(word)
	if err != nil {
		return "", err
	}
	word = p.normalize(word)

	// transcribing phase
	subwords := p.group(word)
	subwords = p.rewrite(subwords)
	subwords = p.transcribe(subwords)

	// finalizing phase
	word = p.syllabify(subwords)
	word = p.localize(word)

	return word, nil
}

// -----------------------------------------------------------------------------

// 0. Just recording beginning (Word)
func (p procedure) input(word string) {
	p.tr.Trace(Input, word, "")
}

// 1. Transliterate (Word -> Word)
//
// This step converts a word from one script to another script or to the
// phonograms, usually based on lexical analysis. Most languages already use
// phonograms which are sufficient to represent the exact pronunciation. But in
// some languages, such as American English or Chinese, it's not true.
func (p procedure) transliterate(word string) (string, error) {
	spec := p.h.Spec()
	translits := p.h.Translits()

	for _, scheme := range spec.Lang.Translit {
		t, ok := translits[scheme]
		if !ok {
			return word, fmt.Errorf("%w: %s", ErrTranslitNotImported, scheme)
		}

		var err error
		word, err = t.Transliterate(word)
		if err != nil {
			return word, fmt.Errorf("%w: %s", ErrTranslit, scheme)
		}

		p.tr.Trace(Transliterate, word, t.Scheme())
	}

	return word, nil
}

// 2. Normalize (Word -> Word)
//
// This step eliminates letter case to make the next steps work easier.
//
// For example, "Hello" in Latin script will be normalized to "hello".
func (p procedure) normalize(word string) string {
	spec := p.h.Spec()

	// Per-spec normalization.
	word = spec.normReplacer.Replace(word)

	p.tr.Trace(Normalize, word, "")

	// Per-script normalization.
	script := spec.script
	except := spec.normLetters

	var buf bytes.Buffer

	for _, let := range word {
		if except[let] || !script.Is(let) {
			buf.WriteRune(let)
		} else {
			buf.WriteRune(script.Normalize(let))
		}
	}

	word = buf.String()

	p.tr.Trace(Normalize, word, spec.Lang.Script)

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
func (p procedure) group(word string) []subword.Subword {
	spec := p.h.Spec()

	rep := subword.NewReplacer(word, 0, 1)

	for i, let := range word {
		letStr := string(let)

		switch {
		case spec.script.Is(let):
			fallthrough
		case spec.puncts[let]:
			fallthrough
		case isSpace(letStr):
			rep.Replace(i, i+len(letStr), letStr)
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
func (p procedure) rewrite(subwords []subword.Subword) []subword.Subword {
	spec := p.h.Spec()

	var swBuf subword.Builder

	swtr := p.tr.SubwordsTracer(Rewrite, subwords)

	for i, sw := range subwords {
		word := sw.Word
		level := sw.Level

		rep := subword.NewReplacer(word, level, 1)

		for _, rule := range spec.Rewrite {
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
func (p procedure) transcribe(subwords []subword.Subword) []subword.Subword {
	spec := p.h.Spec()

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

		for _, rule := range spec.Transcribe {
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
func (p procedure) syllabify(subwords []subword.Subword) string {
	var buf bytes.Buffer
	var jamoBuf bytes.Buffer

	for _, sw := range subwords {
		// Don't touch level=0 subwords. They just have passed through the
		// procedure, because they are meaningless.
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
func (p procedure) localize(word string) string {
	spec := p.h.Spec()
	script := spec.script

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

		punct := script.LocalizePunct(ch)

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

	p.tr.Trace(Localize, word, spec.Lang.Script)

	return word
}
