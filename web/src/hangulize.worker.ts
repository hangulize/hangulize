import './wasm_exec'

import * as Comlink from 'comlink'

import { Specs } from './hangulize'

interface Hangulize {
  (lang: string, word: string): string
  version: string
  specs: Specs
  importPhonemizer: (id: string, fn: (word: string) => Promise<string>) => void
}

declare global {
  const hangulize: Hangulize
}

const go = new Go()
const wasmURL = new URL('hangulize.wasm', import.meta.url)

let initialized: (value?: any) => void
const initialization = new Promise((resolve) => {
  initialized = resolve
})

const phonemize = Comlink.wrap<(id: string, word: string) => string>(
  new Worker(new URL('phonemize.worker', import.meta.url))
)

WebAssembly.instantiateStreaming(fetch(wasmURL), go.importObject).then((wasm) => {
  go.run(wasm.instance)

  for (const lang in hangulize.specs) {
    const spec = hangulize.specs[lang]
    if (spec.lang.phonemizer) {
      const id = spec.lang.phonemizer
      hangulize.importPhonemizer(id, async (word) => {
        return await phonemize(id, word)
      })
    }
  }

  initialized()
})

interface HangulizeEndpoint {
  hangulize: (lang: string, word: string) => string
  getVersion: () => string
  getSpecs: () => Specs
}

Comlink.expose({
  async hangulize(lang: string, word: string) {
    await initialization
    return await hangulize(lang, word)
  },

  async getVersion() {
    await initialization
    return hangulize.version
  },

  async getSpecs() {
    await initialization
    return hangulize.specs
  },
})

export type { HangulizeEndpoint }
