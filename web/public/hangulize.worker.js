self.importScripts('wasm_exec.js')

let hangulizeLoaded = false
let phonemizeLoaded = false

const go = new Go()
WebAssembly.instantiateStreaming(fetch('hangulize.wasm'), go.importObject).then((wasm) => {
  go.run(wasm.instance)

  hangulizeLoaded = true

  for (const lang in hangulize.specs) {
    const spec = hangulize.specs[lang]
    if (spec.lang.phonemizer) {
      const id = spec.lang.phonemizer
      hangulize.importPhonemizer(id, async (word) => {
        if (!phonemizeLoaded) {
          return word
        }
        phonemize.postMessage({
          method: 'phonemize',
          id: id,
          word: word,
        })
        return await new Promise((resolve) => {
          phonemizerResolve = resolve
        })
      })
    }
  }

  self.postMessage({
    method: 'loaded',
    version: hangulize.version,
    specs: hangulize.specs,
  })
})

self.onmessage = async (msg) => {
  switch (msg.data.method) {
    case 'hangulize':
      let result

      try {
        result = await hangulize(msg.data.lang, msg.data.word)
      } catch (e) {
        self.postMessage({
          method: 'error',
          error: e,
          lang: msg.data.lang,
          word: msg.data.word,
          nonce: msg.data.nonce,
        })
        return
      }

      self.postMessage({
        method: 'result',
        result: result,
        lang: msg.data.lang,
        word: msg.data.word,
        nonce: msg.data.nonce,
      })
      break
  }
}

const phonemize = new Worker('phonemize.worker.js')
let phonemizerResolve = null

phonemize.addEventListener('message', (msg) => {
  console.log('phonemize', msg.data)
  switch (msg.data.method) {
    case 'loaded':
      phonemizeLoaded = true
      return

    case 'result':
      phonemizerResolve(msg.data.result)
      return
  }
})
