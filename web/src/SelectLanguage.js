import _ from 'lodash'
import Fuse from 'fuse.js'
import Hangul from 'hangul-js'
import React, { useEffect, useRef, useState } from 'react'
import { Dropdown } from 'semantic-ui-react'
import Flags from './flags'

function SelectLanguage({
  specs = [],
  value = 'ita',
  onChange = (lang) => {},
}) {
  const [options, setOptions] = useState([])

  useEffect(() => {
    setOptions(_.sortBy(specs, (s) => s.lang.korean).map((s) => {
      return {
        key: s.lang.id,
        value: s.lang.id,
        flag: Flags[s.lang.code3],
        text: s.lang.korean,

        search: {
          // it, ka, ...
          code2: s.lang.code2,

          // ita, kat, ...
          code3: s.lang.code3,

          // Italian, Georgian (1st scheme), ...
          en: s.lang.english,

          // 이탈리아어, 조지아어(제1안), ...
          ko: s.lang.korean,

          // ㅇㅣㅌㅏㄹㄹㅣㅇㅏㅇㅓ, ㅈㅗㅈㅣㅇㅏㅇㅓ(ㅈㅔ1ㅇㅏㄴ), ...
          koSyll: Hangul.d(s.lang.korean).join(''),

          // ㅇㅌㄹㅇㅇ, ㅈㅈㅇㅇ(ㅈ1ㅇ), ...
          koInit: Hangul.d(s.lang.korean, true).map((x) => x[0]).join(''),
        },
      }
    }))
  }, [specs])

  const fuse = useRef(null)

  useEffect(() => {
    fuse.current = new Fuse(options, {
      keys: [
        'search.code2',
        'search.code3',
        'search.en',
        'search.ko',
        'search.koSyll',
        'search.koInit',
      ],
    })
  }, [options])

  const search = (_, searchQuery) => {
    searchQuery = Hangul.disassemble(searchQuery).join('')
    return fuse.current.search(searchQuery).map((x) => x.item)
  }

  return (
    <Dropdown
      className="select-language"
      placeholder="언어..."
      button basic floating compact
      search={search}
      options={options}
      value={value}
      onChange={(e, opt) => { onChange(opt.value) }}
    />
  )
}

export default SelectLanguage
