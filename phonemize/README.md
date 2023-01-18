# Phonemize

Here are packages that implement the `hangulize.Phonemizer` interface for
specific languages. Each package provides `P` as the singleton value.

```go
import "github.com/hangulize/hangulize"
import "github.com/hangulize/hangulize/phonemize/furigana"

func main() {
    hangulize.ImportPhonemizer(&furigana.P)
    hangulize.Hangulize("jpn", "天空の城ラピュタ")
}
```
