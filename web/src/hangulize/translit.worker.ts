import './wasm_exec'

import * as Comlink from 'comlink'

import { runWebAssembly } from './runWebAssembly'

declare global {
  const translit: (method: string, word: string) => string
}

const urls: { [method: string]: URL } = {
  cyrillic: new URL('cyrillic.translit.wasm', import.meta.url),
  furigana: new URL('furigana.translit.wasm', import.meta.url),
  pinyin: new URL('pinyin.translit.wasm', import.meta.url),
}

export interface TranslitEndpoint {
  load: (method: string) => void
  transliterate: (method: string, word: string) => Promise<string>
}

export function normalizeMethod(method: string) {
  return method.replace(/\[.+\]$/, '')
}

let ready: Promise<void> | null = null

const endpoint: TranslitEndpoint = {
  async load(method: string) {
    if (ready !== null) {
      throw new Error('load() already called')
    }

    const url = urls[normalizeMethod(method)]

    if (url === undefined) {
      throw new Error(`invalid translit method: ${method}`)
    }

    ready = runWebAssembly(url)
  },

  async transliterate(method: string, word: string) {
    if (ready === null) {
      throw new Error('load() not called yet')
    }
    await ready

    return await translit(method, word)
  },
}

Comlink.expose(endpoint)
