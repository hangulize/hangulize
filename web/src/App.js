import React, { useState, useEffect, useRef } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import {
  Form,
  Message,
  Image,
  Placeholder,
  Container,
  Dimmer,
  Divider,
  Grid,
  Header,
  Loader,
  Segment,
} from 'semantic-ui-react'
import _ from 'underscore'
import Copyright from './Copyright'
import Examples from './Examples'
import Hangulizer from './Hangulizer'
import Prompt from './Prompt'
import DefaultSpecs from './hangulize-specs'
import { getSpec } from './util'
import Logo from './logo.svg'

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
      setVersion('v' + version)
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

  const [result, setResult] = useState('')

  useEffect(() => {
    if (word) {
      document.title = `한글라이즈: ${word}`
    } else {
      document.title = '한글라이즈'
    }
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
            hangulize={hangulize}
            specs={specs}
            lang={lang}
            word={word}
            onChangeLang={setLang}
            onChangeWord={setWord}
            onHangulize={setResult}
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
