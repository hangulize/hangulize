import _ from 'lodash'
import { useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { Container, Icon, Label } from 'semantic-ui-react'

import { Example, findSpec, Spec } from './hangulize'

interface ExamplesProps {
  specs: Spec[]
  lang: string
}

function Examples({ specs, lang }: ExamplesProps) {
  const [examples, setExamples] = useState<Example[]>([])

  const shuffle = () => {
    const spec = findSpec(specs, lang)
    if (spec !== null) {
      setExamples(_.sampleSize(spec.test, 5))
    }
  }

  const prevLang = useRef('')
  useEffect(() => {
    // Check lang dependency manually to hide the dependency with specs.
    if (prevLang.current !== lang) {
      shuffle()
      prevLang.current = lang
    }
  })

  return (
    <Container className="examples">
      <label>예시</label> <Icon name="shuffle" link onClick={shuffle} />
      {examples.map((x, i) => {
        return (
          <Label as={Link} to={`?lang=${lang}&word=${x.word}`} key={i}>
            {x.word}
          </Label>
        )
      })}
    </Container>
  )
}

export default Examples
