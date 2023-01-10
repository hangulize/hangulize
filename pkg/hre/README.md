# HRE

A regular expression dialect for Hangulize

- `^` - the beginning of word not line
- `$` - the end of word not line
- `^^` - the beginning of line (`^` in the standard)
- `$$` - the end of line (`$` in the standard)
- `cat{dog}` - "cat" before "dog" (positive lookahead)
- `{dog}cat` - "cat" after "dog" (positive lookbehind)
- `cat{~dog}` - "cat" before not "dog" (negative lookahead)
- `{~dog}cat` - "cat" after not "dog" (negative lookbehind)
- `<var>` - one of letters in the variable "var"
