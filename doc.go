/*
Package hangulize transcribes non-Korean words into Hangul.

	"Hello!" -> "헬로!"

Hangulize was inspired by Brian Jongseong Park
(http://iceager.egloos.com/2610028). Based on this idea,
the original Hangulize was developed in Python and went out in 2010
(https://github.com/sublee/hangulize). Since then, serving as a web application
on https://hangulize.org/, it has been of great help for Korean translators.

This Go re-implementation is a reboot of Hangulize with feature improvements.

Pipeline

Basically, Hangulize transcribes with 5 steps. These steps include "Normalize",
"Group", "Rewrite", "Transcribe", and "Syllabify". To clarify these concepts,
let's consider an imaginary example of "Hello!" in English into "헬로!"
(actually, English is not supported yet).

First, Hangulize normalizes letter cases:

	"Hello" -> "hello!"

And then, it groups letters by meanings:

	"hello!" -> "hello", "!"

After that, grouped chunks are rewritten as source language-specific rules.
This step is usually for minimizing the differences between pronunciation
and spelling:

	"hello", "!" -> "heˈlō", "!"

And it transcribes rewritten chunks into Hangul Jamo phonemes.

	"heˈlō", "!" -> "ㅎㅔ-ㄹㄹㅗ", "!"

Finally, it composes Jamo phonemes into Hangul syllabic blocks and joins all
groups.

	"ㅎㅔ-ㄹㄹㅗ", "!" -> "헬로!"

Extended Pipeline

Some languages, such as Japanese, may require 2 more steps: "Phonemize" and
"Transliterate". The prior is before the Normalize step, and the latter is
after the Syllabify step.

Japanese uses Kanji which is an ideogram. There is the Kanji-to-Kana mapping
called Furigana. To get Furigana from Kanji, we need a lexical analysis based
on several dictionaries. The Phonemize step guesses the phonograms from a
spelling based on lexical analysis.

	"日本語" -> "ニホンゴ"

Furthermore, Japanese uses the full-width characters for puctuations while
Korean and European languages use the half-width. The full-width puctuations
need to be replaced with the half-width and a space to generate a comfortable
Korean word. The Transliterate step replaces them.

	"이마、아이니유키마스" -> "이마, 아이니유키마스"

Spec

A spec is written by the HGL format which is a configuration DSL for Hangulize
2. One spec is for one language transcription system. So we need to describe
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

We will write many patterns in rewrite/transcribe rules soon. Some expressions
may appear many times annoyingly. To not repeat ourselves, we can use
variables and macros.

A variable is a combination of letters. Variable in pattern will match with
one of the letters. Variable "foo" can be referenced with "<foo>" in the
patterns.

	vars:
	    "vowels" = "a", "e", "i", "o", "u"

A macro expression is replaced with the target before parsing the patterns. "@"
is the common macro for "<vowels>" variable:

	macros:
	    "@" = "<vowels>"

Now we can write "rewrite" rules. There are Pattern and RPattern. Pattern
matches with letters in a word. RPattern represents how the matched letters
should be replaced. A replaced word by a rule would become as the input for
the next rule:

	rewrite:
	    "^gli$"   -> "li"
	    "^gli{@}" -> "li"
	    "{@}gli"  -> "li"
	    "gn{@}"   -> "nJ"

Pattern is based on Regular Expression but it has it's own custom syntax. We
call it "HRE" which means "Hangulize-specific Regular Expression". For the
detail, see the documentation of "github.com/hangulize/hre".

"transcribe" rules are exactly same with "rewrite" rules. But it's RPatterns
represent Hangul Jamo phonemes. In contrast to "rewrite", a replaced word
won't become as the input for the next rules:

	transcribe:
	    "b" -> "ㅂ"
	    "d" -> "ㄷ"
	    "f" -> "ㅍ"
	    "g" -> "ㄱ"

Finally, we should write expected transcription examples. They are used for
unit testing. Verify your spec yourself:

	test:
	    "allegretto" -> "알레그레토"
	    "gita"       -> "지타"
	    "bisnonno"   -> "비스논노"
	    "Pinocchio"  -> "피노키오"

*/
package hangulize
