/*
Package hgl implements a parser for the HGL format which is used for
Hangulize 2.

The HGL format is a simple config language.  A HGL config has sections.  Each
section contains string-strings pairs.  A section can be one of a dictionary or
a pair list:

	lang:
		id      = "ita"
		code    = "it", "ita", "ita"
		english = "Italian"
		korean  = "이탈리아어"
		script  = "roman"

	rewrite:
		"^gli$"    -> "li"
		"^glia$"   -> "g.lia"
		"^glioma$" -> "g.lioma"
		"^gli{@}"  -> "li"
		"{@}gli"   -> "li"

Keys in a dictionary section must be unique, and the section won't keep their
order.  While a pair list section works in an inversed way.  A pair list
section just keeps described pairs in therir order.
*/
package hgl
