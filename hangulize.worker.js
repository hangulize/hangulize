self.importScripts("wasm_exec.js");

const go = new Go();
WebAssembly.instantiateStreaming(fetch("hangulize.wasm"), go.importObject).then((result) => {
  go.run(result.instance);

  self.postMessage({
    method: "initialized",
    version: hangulize.version,
    specs: hangulize.specs
  });

  self.onmessage = (msg) => {
    switch (msg.data.method) {
      case "hangulize":
        const result = hangulize(msg.data.lang, msg.data.word);
        self.postMessage({
          method: "hangulized",
          result: result,
          lang: msg.data.lang,
          word: msg.data.word,
          nonce: msg.data.nonce
        });
        break
    }
  };
});
