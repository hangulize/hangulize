package hangulize

// Event is used for tracing of Hangulize pipeline internal.
type Event struct {
	from string
	to   string
	why  string
}

func event(ch chan<- Event, from string, to string, why string) {
	if ch == nil {
		return
	}
	ch <- Event{from, to, why}
}
