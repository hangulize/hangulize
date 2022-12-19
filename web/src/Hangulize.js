import React from 'react'
import Select from 'react-select'
import './App.css'

class App extends React.Component {
  worker = null
  results = {}

  // A map of {nonce: resolve} for Promise objects to await hangulize()
  // results. nonce='' is reserved for initialization event.
  resolvers = {}

  state = {
    version: '',
    specs: [],
  }

  componentDidMount() {
    this.worker = new Worker('hangulize.worker.js')
    this.worker.addEventListener('message', this.handleWorkerMessage.bind(this))
    this.workerToInit = new Promise((resolve, reject) => {
      this.resolvers[''] = resolve
    })
  }

  handleWorkerMessage(msg) {
    console.log(msg.data)

    switch (msg.data.method) {
      case 'initialized':
        this.setState({version: msg.data.version, specs: msg.data.specs})
        this.resolvers['']()
        break

      case 'hangulized':
        const nonce = msg.data.nonce

        if (!(nonce in this.resolvers)) {
          break
        }

        this.results[nonce] = msg.data.result
        this.resolvers[nonce]()
        delete this.resolvers[nonce]
    }
  }

  async hangulize(lang, word) {
    await this.workerToInit

    const nonce = (Math.random() + 1).toString(36).substring(7)
    this.worker.postMessage({method: 'hangulize', lang: lang, word: word, nonce: nonce})

    let p = new Promise((resolve, reject) => {
      this.resolvers[nonce] = resolve
    })
    await p

    let result = this.results[nonce]
    delete this.results[nonce]

    return result
  }

  render() {
    return (
      <div className="App">
        <p>Version: {this.state.version}</p>
        <Hangulize
          lang="ita"
          word="cappuccino"
          hangulize={this.hangulize.bind(this)}
          specs={this.state.specs}
        />
      </div>
    )
  }
}

class Hangulize extends React.Component {
  static defaultProps = {
    hangulize: (lang, word) => { return word },
    specs: [],
    lang: 'ita',
    word: '',
    result: '',
  }

  state = {
    lang: '',
    word: '',
    result: '',
    specs: [
      {lang: {code3: 'aze', english: 'Azerbaijani', korean: '아제르바이잔어'}},
    ],
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      specs: nextProps.specs,
      lang: nextProps.lang,
      word: nextProps.word,
    })
  }

  async componentDidMount() {
    this.setState({
      specs: this.props.specs,
      lang: this.props.lang,
      word: this.props.word,
    })

    if (this.props.result === '') {
      await this.hangulizeSoon()
    } else {
      this.setState({result: this.props.result})
    }
  }

  async handleChangeLang(opt) {
    this.setState({lang: opt.value})
    await this.hangulizeSoon()
  }

  async handleChangeWord(e) {
    this.setState({word: e.target.value})
    await this.hangulizeSoon()
  }

  async hangulizeSoon() {
    await new Promise(r => setTimeout(r, 100));
    let result = await this.props.hangulize(this.state.lang, this.state.word)
    this.setState({result: result})
  }

  render() {
    return (
      <div>
        <Select
          defaultValue={this.props.lang}
          options={this.state.specs.map((s) => {
            return {value: s.lang.code3, label: s.lang.korean}
          })}
          onChange={this.handleChangeLang.bind(this)}
        />
        <input type="text" defaultValue={this.props.word} onChange={this.handleChangeWord.bind(this)} />
        <p>{this.state.result}</p>
      </div>
    )
  }
}

export default App
