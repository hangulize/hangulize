import Fuse from 'fuse.js'
import React, { useEffect, useRef } from 'react'
import { Dropdown, Input } from 'semantic-ui-react'
import _ from 'underscore'
import Flags from './flags'

function Lang({
  specs = [],
  value = 'ita',
  onChange = (lang) => null,
}) {
  const fuse = useRef(null)
  const options = useRef([])

  useEffect(() => {
    options.current = _.sortBy(specs, (s) => s.lang.korean).map((s) => {
      return {
        key: s.lang.id,
        value: s.lang.id,
        flag: Flags[s.lang.code3],
        text: s.lang.korean,

        search: {
          code2: s.lang.code2,
          code3: s.lang.code3,
          korean: s.lang.korean,
          english: s.lang.english,
        },
      }
    })

    fuse.current = new Fuse(options.current, {
      keys: ['search.code2', 'search.code3', 'search.korean', 'search.english'],
    })
  }, [specs])

  const search = (_, searchQuery) => {
    return fuse.current.search(searchQuery).map((x) => x.item)
  }

  return (
    <Dropdown
      className="lang"
      placeholder="언어..."
      button basic floating compact
      value={value}
      options={options.current}
      search={search}
      onChange={(e, opt) => { onChange(opt.value) }}
    />
  )
}

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
        <Lang
          specs={specs}
          value={lang}
          onChange={(newLang) => { onChange(newLang, word) }}
        />
      )}
      placeholder="외래어 단어..."
      value={word}
      onChange={(e) => { onChange(lang, e.target.value) }}
      size="large"
    />
  )
}

export default Prompt
