package hangulize

// RPattern is used for dynamic replacement.  "R" of RPattern means
// "replacement" or "right-side".
//
// Some expressions in RPattern have special meaning:
//
// - "/" - zero-width edge of chunk
// - "<var>" - ...
//
type RPattern struct {
	expr string
}
