import React, { useState, useEffect, useRef } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import {
  Container,
  Divider,
  Form,
  Header,
  Image,
  Message,
} from 'semantic-ui-react'
import _ from 'underscore'
import Copyright from './Copyright'
import Examples from './Examples'
import Hangulizer from './Hangulizer'
import Prompt from './Prompt'
import DefaultSpecs from './hangulize-specs'
import { getSpec } from './util'
import Logo from './logo.svg'

function useHangulize({
  onInitialize = (version, specs) => null,
  onTranscribe = (result) => null,
  onDelay = () => null,
}) {
  const hangulizer = useRef(null)
  const delay = useRef(null)
  const cancel = useRef(null)

  // Init a worker.
  if (hangulizer.current === null) {
    hangulizer.current = new Hangulizer(onInitialize)
  }

  // hangulize(lang, word) tries to transcribe after 50 ms. If transcription
  // doesn't finish of can't finish in 500 ms, it calls onDelay().
  const hangulize = async (lang, word) => {
    if (!hangulizer.current.initialized) {
      onDelay()
      return
    }

    if (delay.current === null) {
      delay.current = setTimeout(onDelay, 500)
    }
    if (cancel.current !== null) {
      cancel.current(new Error())
    }

    try {
      await new Promise((resolve, reject) => {
        cancel.current = reject
        setTimeout(resolve, 50)
      });
    } catch (e) {
      return
    }

    const result = await hangulizer.current.hangulize(lang, word)
    onTranscribe(result)

    clearTimeout(delay.current)
    delay.current = null
    cancel.current = null
  }

  return hangulize
}

function App() {
  const [version, setVersion] = useState('')
  const [specs, setSpecs] = useState(DefaultSpecs.specs)
  const [result, setResult] = useState('')
  const [loading, setLoading] = useState(true)

  if (specs.length === 0) {
    throw new Error('no specs')
  }

  // Sync lang and word with search parameters.
  const [searchParams, setSearchParams] = useSearchParams()
  const lang = searchParams.get('lang') || _.sample(specs, 1)[0].lang.id
  const word = searchParams.get('word')

  useEffect(() => {
    if (!searchParams.has('lang')) {
      searchParams.set('lang', lang)
      setSearchParams(searchParams)
    }

    if (!searchParams.has('word')) {
      const randomWord = _.sample(getSpec(specs, lang).test, 1)[0].word
      searchParams.set('word', randomWord)
      setSearchParams(searchParams)
    }

    if (word) {
      document.title = `한글라이즈: ${word}`
    } else {
      document.title = '한글라이즈'
    }
  })

  const handleChange = async (lang, word) => {
    searchParams.set('lang', lang)
    searchParams.set('word', word)
    setSearchParams(searchParams)
  }

  const hangulize = useHangulize({
    onInitialize: (version, specs) => {
      setVersion('v' + version)
      setSpecs(specs)
    },
    onTranscribe: (result) => {
      setResult(result)
      setLoading(false)
    },
    onDelay: () => {
      setLoading(true)
    },
  })

  // Transcribe when something has been changed.
  const prevVersion = useRef(null)
  const prevLang = useRef(null)
  const prevWord = useRef(null)

  useEffect(() => {
    if (
      prevVersion.current !== version ||
      prevLang.current !== lang ||
      prevWord.current !== word
    ) {
      hangulize(lang, word)
    }

    prevVersion.current = version
    prevLang.current = lang
    prevWord.current = word
  })

  return (
    <Container text className="Hangulize">
      <Header>
        <Image src={Logo} />
        <Header.Content>
          한글라이즈
          <Header.Subheader>
            {version ? version : '불러오는 중...'}
          </Header.Subheader>
        </Header.Content>
      </Header>

      <Form>
        <Form.Field>
          <Prompt
            specs={specs}
            lang={lang}
            word={word}
            loading={loading}
            onChange={handleChange}
          />
        </Form.Field>
        <Form.Field>
          <Examples specs={specs} lang={lang} />
        </Form.Field>
        <Form.Field>
          {result ? <Message size="massive">{result}</Message> : ''}
        </Form.Field>
      </Form>

      <Divider />
      <Copyright />
    </Container>
  )
}

export default App
