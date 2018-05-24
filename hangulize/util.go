package hangulize

func safeSlice(s string, start int, stop int) string {
	if start < 0 || stop < 0 {
		return ""
	}
	if stop-start > 0 {
		return s[start:stop]
	}
	return ""
}
