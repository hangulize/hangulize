import './wasm_exec'

import * as Comlink from 'comlink'

import { runWebAssembly } from './runWebAssembly'
import type { Specs } from './spec'

export interface HangulizeTrace {
  step: string
  word: string
  why: string
}

interface WasmHangulize {
  (lang: string, word: string, traceFn?: (t: HangulizeTrace) => void): Promise<string>
  version: string
  specs: Specs
  useTranslit: (method: string, fn: (word: string) => Promise<string>) => boolean
}

declare global {
  const hangulize: WasmHangulize
}

const ready = runWebAssembly(new URL('hangulize.wasm', import.meta.url))

export interface HangulizeEndpoint {
  hangulize: (lang: string, word: string, traceFn?: (t: HangulizeTrace) => void) => Promise<string>
  getVersion: () => Promise<string>
  getSpecs: () => Promise<Specs>
  useTranslit: (method: string, fn: (word: string) => Promise<string>) => Promise<void>
}

const endpoint: HangulizeEndpoint = {
  async hangulize(lang: string, word: string, traceFn?: (t: HangulizeTrace) => void) {
    await ready
    return await hangulize(lang, word, traceFn)
  },

  async getVersion() {
    await ready
    return hangulize.version
  },

  async getSpecs() {
    await ready
    return hangulize.specs
  },

  async useTranslit(method: string, fn: (word: string) => Promise<string>) {
    await ready
    hangulize.useTranslit(method, fn)
  },
}

Comlink.expose(endpoint)
