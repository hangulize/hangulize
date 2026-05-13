import { default as _ } from 'lodash'
import { useCallback, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Container, Divider, Header, Image } from 'semantic-ui-react'

import type { Example } from '../hangulize/spec'
import { Hangulize, HangulizeState, useHangulize } from '../hangulize/useHangulize'
import Description from './Description'
import Examples from './Examples'
import Footer from './Footer'
import Prompt from './Prompt'
import Result from './Result'

function determineLoadingResult(
  hangulize: Hangulize,
  lang: string,
  word: string
): [boolean, string] {
  if (!hangulize.isValidInput(lang, word)) {
    return [false, '']
  }

  if (!hangulize.result) {
    return [true, '…']
  }

  if (hangulize.state === HangulizeState.PROCESSING_DELAYED) {
    return [true, hangulize.result]
  }

  return [false, hangulize.result]
}

export default function App() {
  const [hangulize, setHangulizeInput] = useHangulize()

  // Sync lang and word with search parameters.
  const [searchParams, setSearchParams] = useSearchParams()
  const lang = searchParams.get('lang') || (_.sample(Object.keys(hangulize.specs)) as string)
  const word = searchParams.get('word') || ''

  const spec = hangulize.specs[lang]
  if (spec === undefined) {
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

  const handleChange = useCallback(async (lang: string, word: string) => {
    searchParams.set('lang', lang)
    searchParams.set('word', word)
    setSearchParams(searchParams, { replace: true })
  }, [])

  useEffect(() => {
    setHangulizeInput(lang, word)
  }, [lang, word])

  const [loading, result] = determineLoadingResult(hangulize, lang, word)

  return (
    <Container text className="app">
      <Header as="h1">
        <Image src={process.env.PUBLIC_URL + '/logo.svg'} />
        <Header.Content>
          한글라이즈
          <Header.Subheader className="version">v{hangulize.version}</Header.Subheader>
        </Header.Content>
      </Header>

      <Prompt
        specs={hangulize.specs}
        lang={lang}
        word={word}
        loading={loading}
        onChange={handleChange}
      />
      <Examples specs={hangulize.specs} lang={lang} />
      <Result loading={loading}>{result}</Result>

      <Description />
      <Divider />
      <Footer />
    </Container>
  )
}
