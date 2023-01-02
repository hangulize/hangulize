import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Container, Button, Form, Icon, Label } from 'semantic-ui-react'
import _ from 'underscore'
import { getSpec } from './util'

function Examples({ specs, lang, onClick }) {
  const [examples, setExamples] = useState([])

  const shuffle = () => {
    const spec = getSpec(specs, lang)
    if (spec !== null) {
      setExamples(_.sample(spec.test, 5))
    }
  }
  useEffect(shuffle, [lang])

  return (
    <Container className="examples">
      <label>예시</label>
      <Button basic compact size="mini" icon="shuffle" onClick={shuffle} />
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
