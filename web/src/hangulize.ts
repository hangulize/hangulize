import { default as _ } from 'lodash'
import { useMemo, useRef } from 'react'

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

interface MessagePayload {
  method: string
}

interface LoadedPayload extends MessagePayload {
  method: 'loaded'
  version: string
  specs: Specs
}

interface ResultPayload extends MessagePayload {
  method: 'result'
  nonce: string
  result: string
}

interface ErrorPayload extends MessagePayload {
  method: 'error'
  nonce: string
  error: Error
}

class Hangulizer {
  worker: Worker
  resolvers: { [nonce: string]: (resultOrError: string | Error) => void }
  loaded: boolean
  onLoad: (version: string, specs: Specs) => void
  specs: Specs

  constructor(onLoad: (version: string, specs: Specs) => void) {
    this.worker = new Worker(new URL('hangulize.worker.ts', import.meta.url))
    this.worker.addEventListener('message', this.handleMessage.bind(this))
    this.resolvers = {}
    this.loaded = false
    this.onLoad = onLoad
    this.specs = {}
  }

  handleMessage(msg: { data: LoadedPayload | ResultPayload | ErrorPayload }) {
    console.log(msg.data)

    switch (msg.data.method) {
      case 'loaded':
        this.loaded = true
        this.specs = msg.data.specs

        if (this.onLoad !== undefined) {
          this.onLoad(msg.data.version, msg.data.specs)
        }
        return

      case 'result':
        if (!(msg.data.nonce in this.resolvers)) {
          break
        }

        this.resolvers[msg.data.nonce](msg.data.result)
        delete this.resolvers[msg.data.nonce]
        return

      case 'error':
        if (!(msg.data.nonce in this.resolvers)) {
          break
        }

        this.resolvers[msg.data.nonce](msg.data.error)
        delete this.resolvers[msg.data.nonce]

        return
    }

    throw new Error(`unexpected method from worker: ${msg.data.method}`)
  }

  async hangulize(lang: string, word: string) {
    if (!word || !this.loaded) {
      return ''
    }

    const nonce = _.uniqueId()
    this.worker.postMessage({
      method: 'hangulize',
      lang,
      word,
      nonce,
    })

    const resultOrError = await new Promise((resolve) => {
      this.resolvers[nonce] = resolve
    })

    if (resultOrError instanceof Error) {
      throw resultOrError as Error
    } else {
      return resultOrError as string
    }
  }
}

interface UseHangulizeProps {
  onInit: (version: string, specs: Specs) => void
  onResult: (result: string) => void
  onSlow: () => void
}

function useHangulize({ onInit: onInit, onResult: onResult, onSlow }: UseHangulizeProps) {
  const worker = useMemo(() => new Hangulizer(onInit), [])
  const lastId = useRef('')
  const lastJob = useRef<Promise<string> | null>(null)
  const slow = useRef<ReturnType<typeof setTimeout> | null>(null)

  // hangulize(lang, word) tries to transcribe after 50 ms. If transcription
  // doesn't finish of can't finish in 500 ms, it calls onSlow().
  const hangulize = async (lang: string, word: string, delay = 50) => {
    // Trigger onSlow() when it cannot finish in 500ms.
    if (!worker.loaded) {
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
        return await worker.hangulize(lang, word)
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

  return hangulize
}

export { Hangulizer, useHangulize }
export type { Example, Spec, Specs }
