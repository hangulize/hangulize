import { default as _ } from 'lodash'
import { useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Container, Divider, Header, Image } from 'semantic-ui-react'

import Description from './Description'
import Examples from './Examples'
import Footer from './Footer'
import { Example, findSpec, Spec, useHangulize } from './hangulize'
import Prompt from './Prompt'
import Result from './Result'
import staticSpecs from './specs.json'

function App() {
  const [version, setVersion] = useState('')
  const [specs, setSpecs] = useState(staticSpecs.specs as Spec[])
  const [result, setResult] = useState('')
  const [loading, setLoading] = useState(true)

  if (specs.length === 0) {
    throw new Error('no specs')
  }

  // Sync lang and word with search parameters.
  const [searchParams, setSearchParams] = useSearchParams()
  const lang = searchParams.get('lang') || (_.sample(specs) as Spec).lang.id
  const word = (searchParams.get('word') || '').trim()

  const spec = findSpec(specs, lang)
  if (spec === null) {
    throw new Error(`unknown lang: ${lang}`)
  }

  useEffect(() => {
    let redirect = false
    if (!searchParams.has('lang')) {
      searchParams.set('lang', lang)
      redirect = true
    }
    if (!searchParams.has('word')) {
      if (spec.test.length !== 0) {
        const randomWord = (_.sample(spec.test) as Example).word
        searchParams.set('word', randomWord)
        redirect = true
      }
    }
    if (redirect) {
      setSearchParams(searchParams, { replace: true })
    }

    if (word) {
      document.title = `한글라이즈: ${word}`
    } else {
      document.title = '한글라이즈'
    }
  })

  const handleChange = async (lang: string, word: string) => {
    searchParams.set('lang', lang)
    searchParams.set('word', word)
    setSearchParams(searchParams, { replace: true })
  }

  const hangulize = useHangulize({
    onInitialize: (version: string, specs: Spec[]) => {
      setVersion('v' + version)
      setSpecs(specs)
    },
    onTranscribe: (result: string) => {
      setResult(result)
      setLoading(false)
    },
    onSlow: () => {
      setLoading(true)
    },
  })

  // Transcribe when something has been changed.
  const prevVersion = useRef<string | null>(null)
  const prevLang = useRef<string | null>(null)
  const prevWord = useRef<string | null>(null)

  useEffect(() => {
    if (prevVersion.current !== version || prevLang.current !== lang || prevWord.current !== word) {
      const delay = result ? 50 : 0
      hangulize(lang, word, delay)
    }

    prevVersion.current = version
    prevLang.current = lang
    prevWord.current = word
  })

  return (
    <Container text className="app">
      <Header as="h1">
        <Image src={process.env.PUBLIC_URL + '/logo.svg'} />
        <Header.Content>
          한글라이즈
          <Header.Subheader className="version">
            {version ? <span>{version}</span> : '불러오는 중...'}
          </Header.Subheader>
        </Header.Content>
      </Header>

      <Prompt specs={specs} lang={lang} word={word} loading={loading} onChange={handleChange} />
      <Examples specs={specs} lang={lang} />

      {word && (loading || result) ? <Result loading={loading}>{result}</Result> : <></>}

      <Description />
      <Divider />
      <Footer />
    </Container>
  )
}

export default App
