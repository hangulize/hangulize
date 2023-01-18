package hangulize

import "errors"

// ErrSpecNotFound occurs when the spec for the given language is not found.
var ErrSpecNotFound = errors.New("spec not found")

// ErrPhonemizerNotImported occurs when the selected spec requires a phonemizer
// but it has not been imported yet.
var ErrPhonemizerNotImported = errors.New("phonemizer not imported")
