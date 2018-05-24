package hangulize

func expandEdges(expr string) string {
	expr = reLeftEdge.ReplaceAllStringFunc(expr, func(e string) string {
		switch e {
		case ``:
			return ``
		case `^`:
			return `(?:^|\s+)`
		default:
			// ^^...
			return `^`
		}
	})
	expr = reRightEdge.ReplaceAllStringFunc(expr, func(e string) string {
		switch e {
		case ``:
			return ``
		case `$`:
			return `(?:$|\s+)`
		default:
			// $$...
			return `$`
		}
	})
	return expr
}
