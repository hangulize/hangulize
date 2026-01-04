import { default as _ } from 'lodash'
import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
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

export default function App({ introHTML }: { introHTML: string }) {
  const [hangulize, setHangulizeInput] = useHangulize()
  const navigate = useNavigate()
  const { lang: pathLang, word: pathWord } = useParams<{ lang?: string; word?: string }>()
  const [searchParams] = useSearchParams()
  const [userClearedInput, setUserClearedInput] = useState(false)

  // Check for legacy querystring parameters and redirect
  const queryLang = searchParams.get('lang')
  const queryWord = searchParams.get('word')

  // Handle legacy querystring redirect
  if (queryLang || queryWord) {
    const redirectLang = queryLang || pathLang || (_.sample(Object.keys(hangulize.specs)) as string)
    const redirectWord = queryWord || ''
    navigate(`/${redirectLang}${redirectWord ? '/' + redirectWord : ''}`, { replace: true })
    return null
  }

  // Determine lang and word from path parameters
  const lang = pathLang || (_.sample(Object.keys(hangulize.specs)) as string)
  const word = pathWord || ''

  // Set defaults if no lang provided in path
  if (!pathLang) {
    navigate(`/${lang}`, { replace: true })
    return null
  }

  const spec = hangulize.specs[lang]
  if (spec === undefined) {
    throw new Error(`unknown lang: ${lang}`)
  }

  useEffect(() => {
    // Set random word if none provided
    if (!word && spec.test.length !== 0 && !userClearedInput) {
      const randomWord = (_.sample(spec.test) as Example).word
      navigate(`/${lang}/${randomWord}`, { replace: true })
      return
    }

    if (word) {
      document.title = `한글라이즈: ${word}`
    } else {
      document.title = '한글라이즈'
    }
  }, [lang, word, spec.test, navigate, userClearedInput])

  useEffect(() => {
    setHangulizeInput(lang, word)
  }, [lang, word])

  // Handle data-nav links
  useEffect(() => {
    const handleNavClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement
      const link = target.closest('a[data-nav="true"]')
      if (link) {
        e.preventDefault()
        const href = link.getAttribute('href')
        if (href) {
          navigate(href)
        }
      }
    }

    document.addEventListener('click', handleNavClick)
    return () => document.removeEventListener('click', handleNavClick)
  }, [navigate])

  const handleChange = useCallback(
    async (newLang: string, newWord: string) => {
      // Track if user explicitly cleared the input
      if (word && !newWord) {
        setUserClearedInput(true)
      } else if (newWord) {
        setUserClearedInput(false)
      }
      navigate(`/${newLang}${newWord ? '/' + newWord : ''}`, { replace: true })
    },
    [navigate, word]
  )

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

      <Description lang={lang} />
      <Divider />
      <section className="intro" dangerouslySetInnerHTML={{ __html: introHTML }} />

      <Divider />
      <Footer />
    </Container>
  )
}
