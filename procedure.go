package hangulize

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/hangulize/hangulize/internal/jamo"
	"github.com/hangulize/hangulize/internal/subword"
)

// procedure implements the Hangulize procedure.
type procedure struct {
	spec      *Spec
	translits map[string]Translit
	tracer    *tracer
}

// newProcedure creates a new procedure.
func newProcedure(spec *Spec, translits map[string]Translit, traceFunc func(Trace)) *procedure {
	return &procedure{spec, translits, newTracer(traceFunc)}
}

// forward runs the Hangulize procedure for a word.
func (p procedure) forward(word string) (string, error) {
	p.tracer.Input(word)

	// phase: preparing
	word, err := p.transliterate(word)
	if err != nil {
		return "", err
	}
	word = p.normalize(word)

	// phase: transcribing
	subwords := p.partition(word)
	subwords = p.rewrite(subwords)
	subwords = p.transcribe(subwords)

	// phase: finalizing
	word = p.syllabify(subwords)
	word = p.localize(word)

	return word, nil
}

// 1. Transliterate (Word -> Word)
//
// This step converts a word from one script to another script or to the
// phonograms, usually based on lexical analysis. Most languages already use
// phonograms which are sufficient to represent the exact pronunciation. But in
// some languages, such as American English or Chinese, it's not true.
func (p procedure) transliterate(word string) (string, error) {
	for _, scheme := range p.spec.Lang.Translit {
		t, ok := p.translits[scheme]
		if !ok {
			return word, fmt.Errorf("%w: %s", ErrTranslitNotImported, scheme)
		}

		var err error
		word, err = t.Transliterate(word)
		if err != nil {
			return word, fmt.Errorf("%w: %s", ErrTranslit, scheme)
		}

		p.tracer.Transliterate(word, t.Scheme())
	}

	return word, nil
}

// 2. Normalize (Word -> Word)
//
// This step eliminates letter case to make the next steps work easier.
//
// For example, "Hello" in Latin script will be normalized to "hello".
func (p procedure) normalize(word string) string {
	// Per-spec normalization.
	word = p.spec.normReplacer.Replace(word)
	p.tracer.Normalize(word, "")

	// Per-script normalization.
	script := p.spec.script
	except := p.spec.normLetters

	var buf bytes.Buffer

	for _, let := range word {
		if except[let] || !script.Is(let) {
			buf.WriteRune(let)
		} else {
			buf.WriteRune(script.Normalize(let))
		}
	}

	word = buf.String()
	p.tracer.Normalize(word, p.spec.Lang.Script)
	return word
}

// 3. Partition (Word -> Subwords[level=0 or 1])
//
// Meaningful letter is the letter which appears in the rewrite/transcribe
// rules. This step partitions letters by their meaningness into subwords. A
// meaningful subword has level=1 meanwhile meaningless one has level=0.
//
// For example, "hello, world!" will be partitioned into
// [{"hello",1}, {", ",0}, {"world",1}, {"!",0}].
func (p procedure) partition(word string) []subword.Subword {
	rep := subword.NewReplacer(word, 0, 1)

	for i, let := range word {
		letStr := string(let)

		switch {
		case p.spec.script.Is(let):
			fallthrough
		case p.spec.puncts[let]:
			fallthrough
		case hasSpaceOnly(letStr):
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
	var swBuf subword.Builder

	traceRecordSubword, traceCommit := p.tracer.Rewrite(subwords)
	defer traceCommit()

	for i, sw := range subwords {
		word := sw.Word
		level := sw.Level

		rep := subword.NewReplacer(word, level, 1)

		for _, rule := range p.spec.Rewrite {
			repls := rule.replacements(word)
			rep.ReplaceBy(repls...)
			word = rep.String()

			traceRecordSubword(i, word, rule)
		}

		swBuf.Write(rep.Subwords()...)
	}

	return swBuf.Subwords()
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
	var swBuf subword.Builder

	traceSubword, trace := p.tracer.Transcribe(subwords)
	defer trace()

	for i, sw := range subwords {
		if sw.Level == 0 {
			swBuf.Write(sw)
			continue
		}

		word := sw.Word
		level := sw.Level

		rep := subword.NewReplacer(word, level, 2)

		// transcribe is not rewrite. A result of a replacement is not the
		// input of the next replacement. dummy masks the replaced subwords
		// with NULL characters.
		dummy := subword.NewReplacer(word, 0, 0)

		for _, rule := range p.spec.Transcribe {
			repls := rule.replacements(word)
			rep.ReplaceBy(repls...)

			for _, repl := range repls {
				nulls := strings.Repeat("\x00", len(repl.Word))
				dummy.Replace(repl.Start, repl.Stop, nulls)
			}

			word = dummy.String()
			traceSubword(i, rep.String(), rule)
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

	return swBuf.Subwords()
}

// 6. Syllabify (Subwords -> Word)
//
// This step converts decomposed Jamo phonemes to composed Hangul syllables.
//
// For example, "ㅎㅔ-ㄹㄹㅗ" becomes "헬로".
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
	p.tracer.Syllabify(word)
	return word
}

// 7. Localize (Word -> Word)
//
// Finally, this step converts foreign punctuations to fit in Korean.
//
// Korean has adapted the European punctuations. Those are the most common in
// the world. But a few languages, such as Japanese or Chinese, use different
// punctuations with Korean. This step will reduce that kind of culture gap.
//
// For example, "「...」" becomes "'...'".
func (p procedure) localize(word string) string {
	script := p.spec.script

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
	p.tracer.Localize(word, p.spec.Lang.Script)
	return word
}
