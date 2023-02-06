import * as Comlink from 'comlink'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import type { HangulizeEndpoint } from './hangulize.worker'
import manifest from './manifest.json'
import type { Specs } from './spec'
import type { TranslitEndpoint } from './translit.worker'
import { normalizeMethod } from './translit.worker'

export const enum HangulizeState {
  INITIALIZING,
  IDLE,
  PROCESSING,
  PROCESSING_DELAYED,
}

export class Hangulize {
  state: HangulizeState
  version: string
  specs: Specs
  result: string

  constructor(state: HangulizeState, version: string, specs: Specs, result: string) {
    if (Object.keys(specs).length === 0) {
      throw new Error('no specs')
    }

    this.state = state
    this.version = version
    this.specs = specs
    this.result = result
  }

  isValidInput(inputOrLang: HangulizeInput | string, word = '') {
    let lang: string
    if (typeof inputOrLang === 'string') {
      lang = inputOrLang
    } else {
      lang = inputOrLang.lang
      word = inputOrLang.word
    }

    if (!word.trim()) {
      return false
    }

    if (this.specs[lang] === undefined) {
      return false
    }

    return true
  }
}

export interface SetHangulizeInput {
  (lang: string, word: string): void
}

interface HangulizeInput {
  lang: string
  word: string
}

function spawnWorkers(): [
  Comlink.Remote<HangulizeEndpoint>,
  { [method: string]: Comlink.Remote<TranslitEndpoint> }
] {
  const worker = Comlink.wrap<HangulizeEndpoint>(
    new Worker(new URL('hangulize.worker', import.meta.url))
  )

  const translitWorkers: { [method: string]: Comlink.Remote<TranslitEndpoint> } = {}
  manifest.translits.forEach((method) => {
    const normalMethod = normalizeMethod(method)

    if (!(normalMethod in translitWorkers)) {
      translitWorkers[normalMethod] = Comlink.wrap<TranslitEndpoint>(
        new Worker(new URL('translit.worker', import.meta.url))
      )
      translitWorkers[normalMethod].load(method)
    }

    const translit = Comlink.proxy(async (word: string) => {
      return await translitWorkers[normalMethod].transliterate(method, word)
    })
    worker.useTranslit(method, translit)
  })

  return [worker, translitWorkers]
}

export function useHangulize(): [Hangulize, SetHangulizeInput] {
  const [state, setState] = useState(HangulizeState.INITIALIZING)
  const [version, setVersion] = useState(manifest.version)
  const [specs, setSpecs] = useState<Specs>(manifest.specs)
  const [result, setResult] = useState('')
  const hangulize = new Hangulize(state, version, specs, result)

  const worker = useMemo(() => {
    const [worker] = spawnWorkers()

    const setStatesOnInit = async () => {
      const newVersion = await worker.getVersion()
      const newSpecs = await worker.getSpecs()

      setVersion(newVersion)
      setSpecs(newSpecs)
    }
    setStatesOnInit()

    return worker
  }, [])

  const [input, setInput] = useState<HangulizeInput>({ lang: '', word: '' })
  const target = useRef<HangulizeInput>({ lang: '', word: '' })
  const processing = useRef<Promise<string> | null>(null)
  const delayTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    async function hangulizeEffect() {
      if (!hangulize.isValidInput(input)) {
        setResult('')
        setState(HangulizeState.IDLE)

        if (delayTimer.current !== null) {
          clearTimeout(delayTimer.current)
          delayTimer.current = null
        }
        return
      }

      if (state === HangulizeState.INITIALIZING) {
        return
      }

      // Ignore the same target.
      if (target.current === input) {
        return
      }

      // Set the latest target.
      target.current = input

      if (delayTimer.current === null) {
        setState(HangulizeState.PROCESSING)

        delayTimer.current = setTimeout(() => {
          setState(HangulizeState.PROCESSING_DELAYED)
        }, 500)
      }

      // Sleep for 50 ms to avoid too often job.
      if (result) {
        await new Promise((resolve) => {
          setTimeout(resolve, 50)
        })
        if (target.current !== input) return
      }

      // Wait for the currently running job to not make too many jobs which
      // probably will be dropped.
      if (processing.current !== null) {
        await processing.current
      }

      // This is the last request. Hangulize it.
      processing.current = (async () => {
        try {
          return await worker.hangulize(input.lang, input.word)
        } catch (e) {
          console.error(e)
          return ''
        }
      })()

      const newResult = await processing.current
      if (target.current !== input) return

      // This is still the last request. Finally, trigger onResult().
      setResult(newResult)
      setState(HangulizeState.IDLE)

      // Don't call onSlow() if it finishes in 500ms.
      clearTimeout(delayTimer.current)
      delayTimer.current = null
    }
    hangulizeEffect()
  }, [specs, input])

  const setHangulizeInput = useCallback((lang: string, word: string) => {
    setInput({ lang: lang, word: word })
  }, [])

  return [hangulize, setHangulizeInput]
}
