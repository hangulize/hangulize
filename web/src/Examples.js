import _ from 'lodash'
import React, { useState, useEffect, useRef } from 'react'
import { Link } from 'react-router-dom'
import { Container, Icon, Label } from 'semantic-ui-react'
import { getSpec } from './util'

function Examples({ specs, lang, onClick }) {
  const [examples, setExamples] = useState([])

  const shuffle = () => {
    const spec = getSpec(specs, lang)
    if (spec !== null) {
      setExamples(_.sampleSize(spec.test, 5))
    }
  }

  const prevLang = useRef(null)
  useEffect(() => {
    if (prevLang.current !== lang) {
      shuffle()
    }
    prevLang.current = lang
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
