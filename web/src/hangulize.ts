import { useRef } from 'react'

interface Example {
  word: string
  transcribed: string
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

interface Data {
  method: string
}

interface InitializedData extends Data {
  method: 'initialized'
  version: string
  specs: Spec[]
}

interface HangulizedData extends Data {
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

  handleMessage(msg: { data: InitializedData | HangulizedData }) {
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

    const nonce = (Math.random() + 1).toString(36).substring(7)
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

interface useHangulizeProps {
  onInitialize: (version: string, specs: Spec[]) => void
  onTranscribe: (result: string) => void
  onSlow: () => void
}

function useHangulize({ onInitialize, onTranscribe, onSlow }: useHangulizeProps) {
  const worker = useRef<Hangulizer | null>(null)
  const slowTimeout = useRef<ReturnType<typeof setTimeout> | null>(null)
  const cancel = useRef<((reason: Error) => void) | null>(null)

  // Init a worker.
  if (worker.current === null) {
    worker.current = new Hangulizer(onInitialize)
  }

  // hangulize(lang, word) tries to transcribe after 50 ms. If transcription
  // doesn't finish of can't finish in 500 ms, it calls onSlow().
  const hangulize = async (lang: string, word: string, delay = 50) => {
    if (worker.current === null) {
      return
    }

    if (!worker.current.initialized) {
      onSlow()
      return
    }

    if (slowTimeout.current === null) {
      slowTimeout.current = setTimeout(onSlow, 500)
    }
    if (cancel.current !== null) {
      cancel.current(new Error())
    }

    try {
      await new Promise((resolve, reject: (reason: Error) => void) => {
        cancel.current = reject
        setTimeout(resolve, delay)
      })
    } catch (e) {
      return
    }

    const result = await worker.current.hangulize(lang, word)
    onTranscribe(result)

    clearTimeout(slowTimeout.current)
    slowTimeout.current = null
    cancel.current = null
  }

  return hangulize
}

export { findSpec, Hangulizer, useHangulize }
export type { Example, Spec }
