self.importScripts('wasm_exec.js')

const go = new Go()
WebAssembly.instantiateStreaming(fetch('hangulize.wasm'), go.importObject).then((wasm) => {
  go.run(wasm.instance)

  self.postMessage({
    method: 'loaded',
    version: hangulize.version,
    specs: hangulize.specs,
  })
})

self.onmessage = async (msg) => {
  switch (msg.data.method) {
    case 'hangulize':
      let result

      try {
        result = await hangulize(msg.data.lang, msg.data.word)
      } catch (e) {
        self.postMessage({
          method: 'error',
          error: e,
          lang: msg.data.lang,
          word: msg.data.word,
          nonce: msg.data.nonce,
        })
        return
      }

      self.postMessage({
        method: 'result',
        result: result,
        lang: msg.data.lang,
        word: msg.data.word,
        nonce: msg.data.nonce,
      })
      break
  }
}
