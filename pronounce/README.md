# Pronouncers

Here are packages that implement the `hangulize.Pronouncer` interface for
specific languages. Each package provides `P` as the singleton value.

```go
import "github.com/hangulize/hangulize"
import "github.com/hangulize/hangulize/pronounce/furigana"

func main() {
    hangulize.UsePronouncer(&furigana.P)
}
```
