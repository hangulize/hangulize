import './wasm_exec'

export function runWebAssembly(url: URL) {
  /**
   * Loads a Go WebAssembly.
   *
   * @returns A promise being resolved when loading and calling the callback is
   * done.
   *
   */
  const go = new Go()

  let markReady: () => void
  const ready = new Promise<void>((resolve) => {
    markReady = resolve as () => void
  })

  WebAssembly.instantiateStreaming(fetch(url), go.importObject).then((wasm) => {
    go.run(wasm.instance)
    markReady()
  })

  return ready
}
