package main

import (
	"github.com/gopherjs/gopherjs/js"

	"github.com/sublee/hangulize2/hangulize"
)

func main() {
	js.Global.Set("hangulize", hangulize.Hangulize)
}
