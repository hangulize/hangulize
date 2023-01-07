import React from 'react'
import { Input } from 'semantic-ui-react'
import SelectLanguage from './SelectLanguage'

function Prompt({
  specs = [],
  lang = 'ita',
  word = '',
  loading = false,
  onChange = (lang, word) => {},
}) {
  return (
    <Input
      className="word"
      fluid
      loading={loading}
      actionPosition="left"
      action={(
        <SelectLanguage
          specs={specs}
          value={lang}
          onChange={(newLang) => { onChange(newLang, word) }}
        />
      )}
      placeholder="외래어 단어..."
      value={word}
      onChange={(e) => { onChange(lang, e.target.value) }}
      size="medium"
    />
  )
}

export default Prompt
