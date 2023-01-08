import _ from 'lodash'
import React, { useState, useEffect, useRef } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Container, Divider, Header, Image } from 'semantic-ui-react'
import Description from './Description'
import Examples from './Examples'
import Footer from './Footer'
import Hangulizer from './Hangulizer'
import Prompt from './Prompt'
import Result from './Result'
import staticSpecs from './specs'
import { getSpec } from './util'

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
  const [specs, setSpecs] = useState(staticSpecs.specs)
  const [result, setResult] = useState('')
  const [loading, setLoading] = useState(true)

  if (specs.length === 0) {
    throw new Error('no specs')
  }

  // Sync lang and word with search parameters.
  const [searchParams, setSearchParams] = useSearchParams()
  const lang = searchParams.get('lang') || _.sample(specs).lang.id
  const word = searchParams.get('word')

  useEffect(() => {
    let redirect = false
    if (!searchParams.has('lang')) {
      searchParams.set('lang', lang)
      redirect = true
    }
    if (!searchParams.has('word')) {
      const randomWord = _.sample(getSpec(specs, lang).test).word
      searchParams.set('word', randomWord)
      redirect = true
    }
    if (redirect) {
      setSearchParams(searchParams, {replace: true})
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
    setSearchParams(searchParams, {replace: true})
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
      <Header as='h1'>
        <Image src={process.env.PUBLIC_URL + '/logo.svg'} />
        <Header.Content>
          한글라이즈
          <Header.Subheader className="version">
            {version ? <span>{version}</span> : '불러오는 중...'}
          </Header.Subheader>
        </Header.Content>
      </Header>

      <Prompt
        specs={specs}
        lang={lang}
        word={word}
        loading={loading}
        onChange={handleChange}
      />
      <Examples specs={specs} lang={lang} />

      {result ? <Result>{result}</Result> : <></>}

      <Description />
      <Divider />
      <Footer />
    </Container>
  )
}

export default App
