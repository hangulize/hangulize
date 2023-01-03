import React from 'react'
import { Button, Dropdown, Input } from 'semantic-ui-react'
import _ from 'underscore'
import Flags from './flags'

function Prompt({
  specs = [],
  lang = 'ita',
  word = '',
  loading = false,
  onChange = (lang, word) => {},
}) {
  const dropdown = (
    <Dropdown
      placeholder="언어..."
      button basic floating search compact
      value={lang}
      options={_.sortBy(specs, (s) => s.lang.korean).map((s) => {
        return {
          key: s.lang.id,
          value: s.lang.id,
          flag: Flags[s.lang.code3],
          text: s.lang.korean,
        }
      })}
      onChange={(e, opt) => { onChange(opt.value, word) }}
    />
  )

  return (
    <Input
      className="word"
      fluid
      loading={loading}
      actionPosition="left"
      action={dropdown}
      placeholder="외래어 단어..."
      value={word}
      onChange={(e) => { onChange(lang, e.target.value) }}
      size="large"
    />
  )
}

export default Prompt
