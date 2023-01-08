CONFIG_FILES = $(shell ls *.js *.json)
SOURCE_FILES = $(shell find src public -type f)

.PHONY: start clean

build: $(SOURCE_FILES) $(CONFIG_FILES) public/hangulize.wasm src/specs.json
	npm run build

PORT ?= 8080
start: $(SOURCE_FILES) $(CONFIG_FILES) public/hangulize.wasm src/specs.json
	PORT=$(PORT) npm start

public/hangulize.wasm:
	$(MAKE) -C ../cmd/hangulize.wasm OUT=$(CURDIR)/$@

src/specs.json:
	go run ../cmd/hangulize-specs > $@

clean:
	rm -rf build public/hangulize.wasm src/specs.json
