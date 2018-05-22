package hangulize

// Event is emitted when a replacement occurs.  It is used for tracing of
// Hangulize pipeline internal.
type Event struct {
	from string
	to   string
	why  string
}

// Emit a replacement event to be traced.
func event(ch chan<- Event, from string, to string, why string) {
	if ch == nil {
		return
	}
	ch <- Event{from, to, why}
}
