package hangulize

import "errors"

// ErrSpecNotFound occurs when the spec for the given language is not found.
var ErrSpecNotFound = errors.New("spec not found")

// ErrTranslit occurs when a transliteration has been failed.
var ErrTranslit = errors.New("translit error")

// ErrTranslitNotImported occurs when the selected spec requires a Translit but
// it has not been imported yet.
var ErrTranslitNotImported = errors.New("translit not imported")
