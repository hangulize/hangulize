CONFIG_FILES = $(shell ls *.js *.json)
SOURCE_FILES = $(shell find src public -type f)
ICON_FILES = \
	public/favicon.ico \
	public/icons/favicon-16x16.png \
	public/icons/favicon-32x32.png \
	public/icons/apple-touch-icon.png \
	public/icons/safari-pinned-tab.svg

.PHONY: start clean

build: $(SOURCE_FILES) $(CONFIG_FILES) $(ICON_FILES) public/hangulize.wasm src/specs.json CNAME
	npm run build
	cp CNAME build/CNAME

PORT ?= 8080
start: $(SOURCE_FILES) $(CONFIG_FILES) $(ICON_FILES) public/hangulize.wasm src/specs.json
	PORT=$(PORT) npm start

src/specs.json: FORCE
	go run ../cmd/hangulize-specs > $@

public/hangulize.wasm: FORCE
	$(MAKE) -C ../cmd/hangulize.wasm OUT=$(CURDIR)/$@

$(ICON_FILES): public/logo.svg real-favicon.json
	npx -y cli-real-favicon generate real-favicon.json /dev/null public/icons
	mv public/icons/favicon.ico public/favicon.ico

clean: FORCE
	rm -rf build public/hangulize.wasm src/specs.json

FORCE:
