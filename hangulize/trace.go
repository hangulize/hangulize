package hangulize

import (
	"fmt"
)

// Trace is emitted when a replacement occurs.  It is used for tracing of
// Hangulize pipeline internal.
type Trace struct {
	why  string
	from string
	to   string
}

func (e *Trace) String() string {
	return fmt.Sprintf("%#v %s", e.to, e.why)
}

// Emit a trace to the channel.
func trace(ch chan<- Trace, to string, from string, why string) {
	if ch == nil {
		return
	}
	if from == to {
		return
	}
	ch <- Trace{why, from, to}
}
