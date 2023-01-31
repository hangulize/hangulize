VERSION = $(shell git describe --tags --match "v[0-9]*" --abbrev=7 | cut -c 2-)
GO_FILES = $(shell GOOS=js GOARCH=wasm go list -f '{{ range .GoFiles }}{{ $$.Dir }}/{{ . }}{{ "\n" }}{{ end }}' ./... ./../../translit/cyrillic/...)
OUT ?= cyrillic.translit.wasm

$(OUT): $(GO_FILES)
	GOOS=js GOARCH=wasm go build -ldflags="-X 'main.version=$(VERSION)'" -o $@
