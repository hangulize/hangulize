class Hangulizer {
  constructor(onInitialize = (v, s) => {}) {
    this.worker = new Worker('hangulize.worker.js')
    this.worker.addEventListener('message', this.handleMessage.bind(this))
    this.resolvers = {}
    this.initialized = false
    this.onInitialize = onInitialize
  }

  handleMessage(msg) {
    console.log(msg.data)

    switch (msg.data.method) {
      case 'initialized':
        this.initialized = true
        this.onInitialize(msg.data.version, msg.data.specs)
        break

      case 'hangulized':
        const nonce = msg.data.nonce

        if (!(nonce in this.resolvers)) {
          break
        }

        this.resolvers[nonce](msg.data.result)
        delete this.resolvers[nonce]
        break

      default:
        throw new Error(`unexpected method from worker: ${msg.data.method}`)
    }
  }

  async hangulize(lang, word) {
    if (!word || !this.initialized) {
      return ''
    }

    const nonce = (Math.random() + 1).toString(36).substring(7)
    this.worker.postMessage({
      method: 'hangulize',
      lang: lang,
      word: word,
      nonce: nonce,
    })

    return await new Promise((resolve, _) => {
      this.resolvers[nonce] = resolve
    })
  }
}

export default Hangulizer
