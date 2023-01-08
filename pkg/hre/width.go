package hre

import (
	"regexp/syntax"
)

// RegexpMaxWidth calculates the maximum width of a parsed Regexp pattern.
func RegexpMaxWidth(re *syntax.Regexp) int {
	switch re.Op {

	case syntax.OpNoMatch, syntax.OpEmptyMatch:
		// matches no strings
		// matches empty string
		return 0

	case syntax.OpLiteral:
		// matches Runes sequence
		return len(re.Rune)

	case syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		// matches Runes interpreted as range pair list
		// matches any character except newline
		// matches any character
		return 1

	case syntax.OpBeginLine, syntax.OpEndLine:
		// matches empty string at beginning of line
		// matches empty string at end of line
		fallthrough
	case syntax.OpBeginText, syntax.OpEndText:
		// matches empty string at beginning of text
		// matches empty string at end of text
		fallthrough
	case syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		// matches word boundary `\b`
		// matches word non-boundary `\B`
		return 0

	case syntax.OpCapture:
		// capturing subexpression with index Cap, optional name Name
		return RegexpMaxWidth(re.Sub0[0])

	case syntax.OpStar, syntax.OpPlus:
		// matches Sub[0] zero or more times
		// matches Sub[0] one or more times return -1
		return -1

	case syntax.OpQuest:
		// matches Sub[0] zero or one times
		return RegexpMaxWidth(re.Sub0[0])

	case syntax.OpRepeat:
		// matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		if re.Max == -1 {
			return -1
		}
		return re.Max * RegexpMaxWidth(re.Sub0[0])

	case syntax.OpConcat:
		// matches concatenation of Subs
		var total int

		for _, sub := range re.Sub {
			n := RegexpMaxWidth(sub)

			if n == -1 {
				return -1
			}

			total += n
		}

		return total

	case syntax.OpAlternate:
		// matches alternation of Subs
		var max int

		for _, sub := range re.Sub {
			n := RegexpMaxWidth(sub)

			if n == -1 {
				return -1
			}

			if n > max {
				max = n
			}
		}

		return max

	default:
		return 0
	}
}
