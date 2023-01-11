CONFIG_FILES = $(shell ls *.js *.json)
SOURCE_FILES = $(shell find src public -type f)
SPEC_FILES = $(shell find ../specs -name '*.hsl')

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-10s %s\n", $$1, $$2}'

ICON_FILES = \
	public/favicon.ico \
	public/icons/favicon-16x16.png \
	public/icons/favicon-32x32.png \
	public/icons/apple-touch-icon.png \
	public/icons/safari-pinned-tab.svg
$(ICON_FILES): public/logo.svg real-favicon.json
	yarn icon
	mv public/icons/favicon.ico public/favicon.ico

src/specs.json: $(SPEC_FILES)
	go run ../cmd/hangulize-specs > $@

public/hangulize.wasm: FORCE
	$(MAKE) -C ../cmd/hangulize.wasm OUT=$(CURDIR)/$@

GEN_FILES = $(ICON_FILES) src/specs.json public/hangulize.wasm

node_modules: package.json
	yarn install

.PHONY: gen
generate: $(GEN_FILES) ## Generate Wasm, JSON, and icon files

.PHONY: lint
lint: node_modules ## Lint by ESLint
	yarn lint

.PHONY: fmt
fmt: node_modules ## Format code by ESLint and Prettier
	yarn lint:fix
	yarn prettier

build: node_modules $(GEN_FILES) ## Build a web app in ./build/
	yarn build
	cp CNAME build/CNAME

.PHONY: start
start: node_modules $(GEN_FILES) ## Start a web app for development
	yarn start

.PHONY: clean
clean: ## Clean generated files and directories
	rm -f ./public/hangulize.wasm
	rm -f ./src/specs.json
	rm -f $(ICON_FILES)
	rm -rf ./node_modules
	rm -rf ./build

FORCE:
