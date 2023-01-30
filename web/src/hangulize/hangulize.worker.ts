import './wasm_exec'

import * as Comlink from 'comlink'

import { runWebAssembly } from './runWebAssembly'
import type { Specs } from './spec'

interface WasmHangulize {
  (lang: string, word: string): Promise<string>
  version: string
  specs: Specs
  importPhonemizer: (id: string, fn: (word: string) => Promise<string>) => void
}

declare global {
  const hangulize: WasmHangulize
}

const ready = runWebAssembly(new URL('hangulize.wasm', import.meta.url))

export interface HangulizeEndpoint {
  hangulize: (lang: string, word: string) => Promise<string>
  getVersion: () => Promise<string>
  getSpecs: () => Promise<Specs>
  importPhonemizer: (id: string, fn: (word: string) => Promise<string>) => Promise<void>
}

const endpoint: HangulizeEndpoint = {
  async hangulize(lang: string, word: string) {
    await ready
    return await hangulize(lang, word)
  },

  async getVersion() {
    await ready
    return hangulize.version
  },

  async getSpecs() {
    await ready
    return hangulize.specs
  },

  async importPhonemizer(id: string, fn: (word: string) => Promise<string>) {
    await ready
    hangulize.importPhonemizer(id, fn)
  },
}

Comlink.expose(endpoint)
