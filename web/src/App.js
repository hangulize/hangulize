import React, { useState, useEffect, useRef } from 'react'
import _ from 'underscore'
import { useSearchParams, Link } from 'react-router-dom'
import DefaultSpecs from './hangulize-specs'
import './Hangulize.css'
import { Container, Dimmer, Divider, Grid, Header, Loader, Segment } from 'semantic-ui-react'
import Prompt from './Prompt'
import Examples from './Examples'
import Hangulizer from './hangulize'
import { getSpec } from './util'

function App() {
  const [version, setVersion] = useState('')
  const [specs, setSpecs] = useState(DefaultSpecs.specs)
  const [hangulize, setHangulize] = useState(null)

  if (specs.length === 0) {
    throw new Error('no specs')
  }

  // Initialize a Hangulize worker.
  const hangulizer = useRef(null)
  if (hangulizer.current === null) {
    hangulizer.current = new Hangulizer((hangulizer, version, specs) => {
      setVersion(version)
      setSpecs(specs)
      setHangulize(() => hangulizer.hangulize.bind(hangulizer))
    })
  }

  // Sync lang and word with search parameters.
  const [searchParams, setSearchParams] = useSearchParams()

  const lang = searchParams.get('lang') || _.sample(specs, 1)[0].lang.id
  const setLang = (newLang) => {
    searchParams.set('lang', newLang)
    setSearchParams(searchParams)
  }
  useEffect(() => {
    if (!searchParams.has('lang')) setLang(lang)
  })

  let word = searchParams.get('word')
  const setWord = (newWord) => {
    searchParams.set('word', newWord)
    setSearchParams(searchParams)
  }
  if (!searchParams.has('word')) {
    word = _.sample(getSpec(specs, lang).test, 1)[0].word
    setWord(word)
  }

  return (
    <Container fluid className="Hangulize">
      <Header>한글라이즈</Header>
      <p>{version}</p>
      <Prompt
        hangulize={hangulize}
        specs={specs}
        lang={lang}
        word={word}
        onChangeLang={setLang}
        onChangeWord={setWord}
      />
      <Examples
        specs={specs}
        lang={lang}
      />
      <Divider />
      <p>
        &copy; 2010–2018{' '}
        <a href="https://www.facebook.com/kkeutsori">Brian</a> &amp;{' '}
        <a href="https://subl.ee/">Heungsub</a>.{' '}
        All rights reserved.
      </p>
    </Container>
  )
}

export default App
