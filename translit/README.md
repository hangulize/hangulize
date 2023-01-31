# Translit

Here are packages that implement the `hangulize.Translit` interface for
specific methods. Each package provides `T` or `Ts` as the singleton value.

Install all of the standard Translits:

```go
import "github.com/hangulize/hangulize"
import "github.com/hangulize/hangulize/translit"

func main() {
    translit.Install()
    hangulize.Hangulize("chi", "靑島")
    hangulize.Hangulize("jpn", "北海道")
}
```

Or use a specific Translit to reduce the build size:

```go
import "github.com/hangulize/hangulize"
import "github.com/hangulize/hangulize/translit/furigana"

func main() {
    hangulize.UseTranslit(&furigana.T)
    hangulize.Hangulize("jpn", "自由ヶ丘")
}
```
