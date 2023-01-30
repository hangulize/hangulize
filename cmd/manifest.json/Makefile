VERSION = $(shell git describe --tags --match "v[0-9]*" --abbrev=7 | cut -c 2-)
GO_FILES = $(shell go list -f '{{ range .GoFiles }}{{ $$.Dir }}/{{ . }}{{ "\n" }}{{ end }}' ./../../...)
SPEC_FILES = $(shell find ../../specs -name '*.hsl')
OUT ?= manifest.json

$(OUT): $(GO_FILES) $(SPEC_FILES)
	go run -ldflags="-X 'main.version=$(VERSION)'" . > $@
