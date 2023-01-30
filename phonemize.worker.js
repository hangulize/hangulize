self.importScripts('wasm_exec.js')

const go = new Go()
WebAssembly.instantiateStreaming(fetch('phonemize.wasm'), go.importObject).then((wasm) => {
  go.run(wasm.instance)
  phonemize.load()
  self.postMessage({ method: 'loaded' })
})

self.onmessage = async (msg) => {
  switch (msg.data.method) {
    case 'phonemize':
      const result = phonemize(msg.data.id, msg.data.word)
      self.postMessage({
        method: 'result',
        result: result,
        id: msg.data.id,
        word: msg.data.word,
      })
      break
  }
}
