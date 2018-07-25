package hangulize

// Analyzer is an interface for lexical analysis on the Normalize step.
//
// Some languages such as English or Japanese require external dictionaries to
// guess pronunciation from spelling. An analyzer normalizes a word based on
// the own external dictionaries. The Normalize step for that languages should
// adapt the analyzer's behavior.
//
// But external dictionaries may have large size of dataset. Therefore
// Hangulize itself doesn't include them by default due to it's lightness. A
// user has a responsibility to install and use a specific analyzer.
//
// There is the Kanji analyzer for Japanese in a separate package.
//
//   import "github.com/hangulize/hangulize"
//   import "github.com/hangulize/hangulize/analyze/kanji"
//
//   fmt.Println(hangulize.HangulizeAnalyze("jpn", kanji.Analyzer, "日本語"))
//
type Analyzer interface {
	ID() string
	Analyze(string) []string
}
