VERSION = $(shell git describe --tags --match "v[0-9]*" --abbrev=7 | cut -c 2-)
GO_FILES = $(shell go list -f '{{ range .GoFiles }}{{ $$.Dir }}/{{ . }}{{ "\n" }}{{ end }}' ./../../...)
HGL_FILES = $(shell find ../../hgls -name '*.hgl')
OUT ?= hangulize.wasm

$(OUT): $(GO_FILES) $(HGL_FILES)
	GOOS=js GOARCH=wasm go build -ldflags="-X 'main.version=$(VERSION)'" -o $@
