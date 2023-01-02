import React from 'react'
import Flags from './flags'
import './Prompt.css'
import { Button, Container, Dropdown, Form, Input, Label, Loader, Message, Segment, TextArea } from 'semantic-ui-react'

function specOption(spec) {
  return {key: spec.lang.id, value: spec.lang.id, flag: Flags[spec.lang.code3], text: spec.lang.korean}
}

class Prompt extends React.Component {
  hangulizing = false

  static defaultProps = {
    hangulize: null, // (lang, word) => transcribed
    specs: [],
    lang: 'ita',
    word: '',
    cachedResult: '',
    readOnly: false,
    onChangeLang: (lang, prevLang) => {},
    onChangeWord: (word, prevWord) => {},
    onHangulize: (result) => {},
  }

  state = {
    loading: false,
    readOnly: false,
    specs: [],
    lang: '',
    word: '',
    result: '',
  }

  async componentDidUpdate(prevProps) {
    if (this.props.specs !== this.state.specs) {
      this.setState({specs: this.props.specs})
    }

    if (this.props.word !== this.state.word) {
      console.log('change word')
      this.setState({word: this.props.word})
      await this.hangulizeSoon(this.props.word, 0)
    }

    if (this.props.hangulize != prevProps.hangulize) {
      console.log('change hangulize', this.props.hangulize, prevProps.hangulize)
      await this.hangulizeSoon(this.state.word, 0)
    }
  }

  async componentDidMount() {
    this.setState({
      readOnly: this.props.readOnly,
      specs: this.props.specs,
      lang: this.props.lang,
      word: this.props.word,
    })

    if (this.props.word === '') {
      return
    }

    if (this.props.cachedResult === '') {
      this.setState({loading: true})
    } else {
      this.setState({result: this.props.cachedResult})
    }
  }

  async handleChangeLang(e, opt) {
    const lang = opt.value
    const prevLang = this.state.lang
    if (lang === prevLang) {
      return
    }

    this.setState({lang: lang})
    this.props.onChangeLang(lang, prevLang)

    await this.hangulizeSoon()
  }

  async handleChangeWord(e) {
    const word = e.target.value
    const prevWord = this.state.word
    if (word === prevWord) {
      return
    }

    this.setState({word: word})
    this.props.onChangeWord(word, prevWord)

    await this.hangulizeSoon(word)
  }

  async hangulizeSoon(word = null, delay = 50) {
    if (word === null) {
      word = this.state.word
    }

    if (this.props.hangulize === null) {
      this.setState({loading: true})
      return
    }

    if (this.rejectHangulize) {
      this.rejectHangulize(new Error())
      delete this.rejectHangulize
    }

    try {
      await new Promise((resolve, reject) => {
        this.rejectHangulize = reject
        setTimeout(resolve, delay)
      });
    } catch (e) {
      return
    }

    // Set loading=true after 100 ms
    let deferredLoading = setTimeout(() => {
      this.setState({loading: true})
    }, 100)

    let result = await this.props.hangulize(this.state.lang, word)
    this.setState({result: result, loading: false})
    this.props.onHangulize(result)

    clearTimeout(deferredLoading)
  }

  render() {
    const dropdown = (
      <Dropdown
        placeholder="언어..."
        button basic floating search compact
        value={this.props.lang}
        options={this.state.specs.map(specOption)}
        onChange={this.handleChangeLang.bind(this)}
        readOnly={this.state.readOnly}
        disabled={this.state.readOnly}
      />
    )

    return (
      <Input
        className="word"
        fluid
        readOnly={this.state.readOnly}
        transparent={this.state.readOnly}
        loading={this.state.loading}
        actionPosition="left"
        action={dropdown}
        placeholder="외래어 단어"
        value={this.props.word}
        onChange={this.handleChangeWord.bind(this)}
        size="large"
      />
    )
  }
}

export default Prompt
