package hangulize

import (
	"fmt"
)

// Event is emitted when a replacement occurs.  It is used for tracing of
// Hangulize pipeline internal.
type Event struct {
	why  string
	from string
	to   string
}

func (e *Event) String() string {
	return fmt.Sprintf("%#v %s", e.to, e.why)
}

// Emit a replacement event to be traced.
func event(ch chan<- Event, to string, from string, why string) {
	if ch == nil {
		return
	}
	if from == to {
		return
	}
	ch <- Event{why, from, to}
}
