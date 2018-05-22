package hangulize

import (
	"fmt"
)

// Event is emitted when a replacement occurs.  It is used for tracing of
// Hangulize pipeline internal.
type Event struct {
	from string
	to   string
	why  string
}

func (e *Event) String() string {
	return fmt.Sprintf("[%s] %#v -> %#v", e.why, e.from, e.to)
}

// Emit a replacement event to be traced.
func event(ch chan<- Event, from string, to string, why string) {
	if ch == nil {
		return
	}
	if from == to {
		return
	}
	ch <- Event{from, to, why}
}
