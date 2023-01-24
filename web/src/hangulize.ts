import * as Comlink from 'comlink'
import { default as _ } from 'lodash'
import { useMemo, useRef, useState } from 'react'

import type { HangulizeEndpoint } from './hangulize.worker'

interface Example {
  word: string
  result: string
}

interface Spec {
  lang: {
    id: string
    code2: string
    code3: string
    english: string
    korean: string
    script: string
    phonemizer: string
  }

  test: Example[]
}

type Specs = { [lang: string]: Spec }

interface UseHangulizeProps {
  onInit: (version: string, specs: Specs) => void
  onResult: (result: string) => void
  onSlow: () => void
}

function useHangulize({ onInit: onInit, onResult: onResult, onSlow }: UseHangulizeProps) {
  const [loaded, setLoaded] = useState(false)

  const hangulize = useMemo(() => {
    const worker = new Worker(new URL('hangulize.worker', import.meta.url))
    const hangulize = Comlink.wrap<HangulizeEndpoint>(worker)

    const waitInit = async () => {
      const version = await hangulize.getVersion()
      const specs = await hangulize.getSpecs()
      setLoaded(true)
      onInit(version, specs)
    }
    waitInit()

    return hangulize
  }, [])

  const lastId = useRef('')
  const lastJob = useRef<Promise<string> | null>(null)
  const slow = useRef<ReturnType<typeof setTimeout> | null>(null)

  // hangulize(lang, word) tries to transcribe after 50 ms. If transcription
  // doesn't finish of can't finish in 500 ms, it calls onSlow().
  return async (lang: string, word: string, delay = 50) => {
    // Trigger onSlow() when it cannot finish in 500ms.
    if (!loaded) {
      onSlow()
      return
    }
    if (slow.current === null) {
      slow.current = setTimeout(onSlow, 500)
    }

    // Only the last request is necessary. After any "await", check lastId with
    // id to decide whether the process should be dropped.
    const id = _.uniqueId()
    lastId.current = id

    // Sleep for the given delay to avoid too often job.
    await new Promise((resolve) => {
      setTimeout(resolve, delay)
    })
    if (lastId.current !== id) return

    // Wait for the currently running job to not make too many jobs which
    // probably will be dropped.
    if (lastJob.current !== null) {
      await lastJob.current
      if (lastId.current !== id) return
    }

    // This is the last request. Hangulize it.
    lastJob.current = (async () => {
      try {
        return await hangulize.hangulize(lang, word)
      } catch (e) {
        console.error(e)
        return ''
      }
    })()

    const result = await lastJob.current
    if (lastId.current !== id) return

    // This is still the last request. Finally, trigger onResult().
    onResult(result)

    // Don't call onSlow() if it finishes in 500ms.
    clearTimeout(slow.current)
    slow.current = null
  }
}

export { useHangulize }
export type { Example, Spec, Specs }
