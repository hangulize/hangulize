import Fuse from 'fuse.js'
import Hangul from 'hangul-js'
import _ from 'lodash'
import { useEffect, useRef, useState } from 'react'
import { Dropdown, DropdownItemProps } from 'semantic-ui-react'

import type { Specs } from '../hangulize/spec'
import Flags from './flags.json'

interface Option extends DropdownItemProps {
  key: string
  value: string
  flag: string
  text: string

  search: {
    // ISO 639-1: it, ka, ...
    code2: string

    // ISO 639-3: ita, kat, ...
    code3: string

    // English: Italian, Georgian (1st scheme), ...
    en: string

    // Korean: 이탈리아어, 조지아어(제1안), ...
    ko: string

    // Korean syllables: ㅇㅣㅌㅏㄹㄹㅣㅇㅏㅇㅓ, ㅈㅗㅈㅣㅇㅏㅇㅓ(ㅈㅔ1ㅇㅏㄴ), ...
    koSyll: string

    // Korean initials: ㅇㅌㄹㅇㅇ, ㅈㅈㅇㅇ(ㅈ1ㅇ), ...
    koInit: string
  }
}

interface SelectLanguageProps {
  specs: Specs
  value: string
  onChange: (lang: string) => void
}

export default function SelectLanguage({ specs, value, onChange }: SelectLanguageProps) {
  const [options, setOptions] = useState<Option[]>([])
  const fuse = useRef<Fuse<Option> | null>(null)

  useEffect(() => {
    setOptions(
      _.sortBy(specs, (s) => s.lang.korean).map((s) => {
        return {
          key: s.lang.id,
          value: s.lang.id,
          flag: Flags[s.lang.id as keyof typeof Flags],
          text: s.lang.korean,

          search: {
            code2: s.lang.code2,

            code3: s.lang.code3,

            en: s.lang.english,

            ko: s.lang.korean,

            koSyll: Hangul.d(s.lang.korean).join(''),

            koInit: Hangul.d(s.lang.korean, true)
              .map((x) => x[0])
              .join(''),
          },
        }
      })
    )
  }, [specs])

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

  const search = (query: string) => {
    if (!fuse.current) {
      return []
    }
    query = Hangul.disassemble(query).join('')
    return fuse.current.search(query).map((x) => x.item)
  }

  return (
    <Dropdown
      className="select-language"
      placeholder="언어..."
      button
      basic
      floating
      compact
      search={(_, q) => search(q)}
      options={options}
      value={value}
      onChange={(_, opt) => {
        onChange(opt.value as string)
      }}
    />
  )
}
