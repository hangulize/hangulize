import { Input } from 'semantic-ui-react'

import { Spec } from './hangulize'
import SelectLanguage from './SelectLanguage'

interface PromptProps {
  specs: Spec[]
  lang: string
  word: string
  loading: boolean
  onChange: (lang: string, word: string) => void
}

function Prompt({ specs, lang, word, loading, onChange }: PromptProps) {
  return (
    <Input
      className="word"
      fluid
      autoFocus
      loading={loading}
      actionPosition="left"
      action={
        <SelectLanguage
          specs={specs}
          value={lang}
          onChange={(newLang) => {
            onChange(newLang, word)
          }}
        />
      }
      placeholder="외래어 단어..."
      value={word}
      onChange={(e) => {
        onChange(lang, e.target.value)
      }}
    />
  )
}

export default Prompt