import './wasm_exec'

import * as Comlink from 'comlink'

declare global {
  const phonemize: (id: string, word: string) => string
}

const go = new Go()
const wasmURL = new URL('phonemize.wasm', import.meta.url)

WebAssembly.instantiateStreaming(fetch(wasmURL), go.importObject).then((wasm) => {
  go.run(wasm.instance)
  Comlink.expose(phonemize)
})
