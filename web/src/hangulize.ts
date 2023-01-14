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
    korean: string
    english: string
  }

  test: Example[]
}

function findSpec(specs: Spec[], lang: string) {
  for (let i = 0; i < specs.length; i++) {
    const spec = specs[i]
    if (spec.lang.id === lang) {
      return spec
    }
  }
  return null
}

interface MessagePayload {
  method: string
}

interface InitializedPayload extends MessagePayload {
  method: 'initialized'
  version: string
  specs: Spec[]
}

interface HangulizedPayload extends MessagePayload {
  method: 'hangulized'
  nonce: string
  result: string
}

class Hangulizer {
  worker: Worker
  resolvers: { [nonce: string]: (result: string) => void }
  initialized: boolean
  onInitialize?: (version: string, specs: Spec[]) => void

  constructor(onInitialize?: (version: string, specs: Spec[]) => void) {
    this.worker = new Worker(process.env.PUBLIC_URL + '/hangulize.worker.js')
    this.worker.addEventListener('message', this.handleMessage.bind(this))
    this.resolvers = {}
    this.initialized = false
    this.onInitialize = onInitialize
  }

  handleMessage(msg: { data: InitializedPayload | HangulizedPayload }) {
    switch (msg.data.method) {
      case 'initialized':
        this.initialized = true

        if (this.onInitialize !== undefined) {
          this.onInitialize(msg.data.version, msg.data.specs)
        }
        return

      case 'hangulized':
        if (!(msg.data.nonce in this.resolvers)) {
          break
        }

        this.resolvers[msg.data.nonce](msg.data.result)
        delete this.resolvers[msg.data.nonce]
        return
    }

    throw new Error(`unexpected method from worker: ${msg.data.method}`)
  }

  async hangulize(lang: string, word: string) {
    if (!word || !this.initialized) {
      return ''
    }

    const nonce = _.uniqueId()
    this.worker.postMessage({
      method: 'hangulize',
      lang,
      word,
      nonce,
    })

    return await new Promise((resolve: (result: string) => void) => {
      this.resolvers[nonce] = resolve
    })
  }
}

interface UseHangulizeProps {
  onInit: (version: string, specs: Spec[]) => void
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
    if (!worker.initialized) {
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
    lastJob.current = worker.hangulize(lang, word)
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

export { findSpec, Hangulizer, useHangulize }
export type { Example, Spec }
