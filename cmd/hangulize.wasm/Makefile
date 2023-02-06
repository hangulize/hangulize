VERSION = $(shell git describe --tags --match "v[0-9]*" --abbrev=7 | cut -c 2-)
GO_FILES = $(shell GOOS=js GOARCH=wasm go list -f '{{ range .GoFiles }}{{ $$.Dir }}/{{ . }}{{ "\n" }}{{ end }}' ./... ./../../...)
SPEC_FILES = $(shell find ../../specs -name '*.hsl')
OUT ?= hangulize.wasm

$(OUT): $(GO_FILES) $(SPEC_FILES)
	GOOS=js GOARCH=wasm go build -ldflags="-X 'main.version=$(VERSION)'" -o $@

.PHONY: test
test:
	GOOS=js GOARCH=wasm go test -exec $$(go env GOROOT)/misc/wasm/go_js_wasm_exec .
