package hangulize

import (
	"bytes"
	"strings"
	"unicode"
)

type pipeline struct {
	h  *Hangulizer
	tr *tracer
}

// forward runs the Hangulize pipeline for a word.
func (p *pipeline) forward(word string) string {
	p.input(word)

	// preparing phase
	word, _ = p.pronounce(word)
	word = p.normalize(word)

	// transcribing phase
	subwords := p.group(word)
	subwords = p.rewrite(subwords)
	subwords = p.transcribe(subwords)

	// finalizing phase
	word = p.compose(subwords)
	word = p.transliterate(word)

	return word
}

// 0. Just recording beginning (Word)
//
func (p *pipeline) input(word string) {
	p.tr.TraceWord("input", "", word)
}

// 1. Pronounce (Word -> Word)
//
// This step guesses the pronunciation from the spelling based on lexical
// analysis. Most languages don't require this step. But some languages, such
// as English, just the spelling is not enough to guess the pronunciation.
//
func (p *pipeline) pronounce(word string) (string, bool) {
	id := p.h.spec.Lang.Pronouncer
	if id == "" {
		// The language doesn't require a pronouncer. It's okay.
		return word, true
	}

	pron, ok := p.h.GetPronouncer(id)
	if ok {
		goto PronouncerFound
	}

	// Fallback by the global pronouncer registry.
	pron, ok = GetPronouncer(id)
	if ok {
		goto PronouncerFound
	}

	// The language requires a pronouncer but not imported yet.
	return word, false

PronouncerFound:
	return pron.Pronounce(word), true
}

// 2. Normalize (Word -> Word)
//
// This step eliminates letter case to make the next steps work easier.
//
// For example, "Hello" in Latin script will be normalized to "hello".
//
func (p *pipeline) normalize(word string) string {
	word = p.h.spec.normReplacer.Replace(word)

	p.tr.TraceWord("normalize", "custom", word)

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

	p.tr.TraceWord("normalize", p.h.spec.Lang.Script, word)

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
func (p *pipeline) group(word string) []subword {
	rep := newSubwordReplacer(word, 0, 1)

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
// NOTE(sublee): But this step has a limitation. It guesses a pronunciation
// from the spelling. But it can be too hard for some script systems, such as
// English or Franch.
//
func (p *pipeline) rewrite(subwords []subword) []subword {
	var swBuf subwordsBuilder

	rtr := p.tr.RuleTracer(subwords)

	for i, sw := range subwords {
		word := sw.word
		level := sw.level

		rep := newSubwordReplacer(word, level, 1)

		for j, rule := range p.h.spec.Rewrite {
			repls := rule.replacements(word)
			rep.ReplaceBy(repls...)
			word = rep.String()

			rtr.Trace(j, rule, i, word)
		}

		swBuf.Append(rep.Subwords()...)
	}

	subwords = swBuf.Subwords()

	rtr.Commit("rewrite")

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
func (p *pipeline) transcribe(subwords []subword) []subword {
	var swBuf subwordsBuilder

	rtr := p.tr.RuleTracer(subwords)

	for i, sw := range subwords {
		if sw.level == 0 {
			swBuf.Append(sw)
			continue
		}

		word := sw.word
		level := sw.level

		rep := newSubwordReplacer(word, level, 2)

		// transcribe is not rewrite. A result of a replacement is not the
		// input of the next replacement. dummy marks the replaced subwords
		// with NULL characters.
		dummy := newSubwordReplacer(word, 0, 0)

		for j, rule := range p.h.spec.Transcribe {
			repls := rule.replacements(word)
			rep.ReplaceBy(repls...)

			for _, repl := range repls {
				nulls := strings.Repeat("\x00", len(repl.word))
				dummy.Replace(repl.start, repl.stop, nulls)
			}

			rep.flush()
			word = dummy.String()

			rtr.Trace(j, rule, i, rep.word)
		}

		swBuf.Append(rep.Subwords()...)
	}

	// Discard level=1 subwords. They have been generated by "3. Rewrite" but
	// never transcribed. They are superfluity of the internal behavior.
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

	rtr.Commit("transcribe")

	return subwords
}

// 6. Compose (Subwords -> Word)
//
// This step converts decomposed Jamo phonemes to composed Hangul syllables.
//
// For example, "ㅎㅔ-ㄹㄹㅗ" will be "헬로".
//
func (p *pipeline) compose(subwords []subword) string {
	var buf bytes.Buffer
	var jamoBuf bytes.Buffer

	for _, sw := range subwords {
		// Don't touch level=0 subwords. They just have passed through the
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

// 7. Transliterate (Word -> Word)
//
// Finally, this step converts foreign punctuations to fit it Korean.
//
// Korean has adapted the European punctuations. Those are the most common in
// the world. But a few langauges, such as Japanese or Chinese, use different
// punctuations with Korean. This step will reduce that kind of culture gap.
//
// For example, "「...」" will be "'...'".
//
func (p *pipeline) transliterate(word string) string {
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

	return buf.String()
}
