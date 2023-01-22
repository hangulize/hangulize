import './wasm_exec'

import { Spec } from './hangulize'

declare class Go {
  run(instance: WebAssembly.Instance): Promise<void>
  importObject: WebAssembly.Imports
}

interface Hangulize {
  (lang: string, word: string): string
  version: string
  specs: { [lang: string]: Spec }
}

declare global {
  const hangulize: Hangulize
}

const go = new Go()
WebAssembly.instantiateStreaming(
  fetch(new URL('hangulize.wasm', import.meta.url)),
  go.importObject
).then((wasm) => {
  go.run(wasm.instance)

  self.postMessage({
    method: 'loaded',
    version: hangulize.version,
    specs: hangulize.specs,
  })
})

self.onmessage = async (msg: {
  data: { method: string; lang: string; word: string; nonce: string }
}) => {
  let result

  switch (msg.data.method) {
    case 'hangulize':
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
