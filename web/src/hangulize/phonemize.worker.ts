import './wasm_exec'

import * as Comlink from 'comlink'

import { runWebAssembly } from './runWebAssembly'

declare global {
  const phonemize: (word: string) => string
}

const urls: { [id: string]: URL } = {
  furigana: new URL('furigana.phonemize.wasm', import.meta.url),
  pinyin: new URL('pinyin.phonemize.wasm', import.meta.url),
}

export interface PhonemizeEndpoint {
  load: (id: string) => void
  phonemize: (id: string, word: string) => Promise<string>
}

let ready: Promise<void> | null = null

const endpoint: PhonemizeEndpoint = {
  async load(id: string) {
    if (urls[id] === undefined) {
      throw new Error(`invalid phonemize id: ${id}`)
    }

    if (ready !== null) {
      throw new Error('load() already called')
    }

    ready = runWebAssembly(urls[id])
  },

  async phonemize(id: string, word: string) {
    if (ready === null) {
      throw new Error('load() not called yet')
    }
    await ready

    return await phonemize(word)
  },
}

Comlink.expose(endpoint)
