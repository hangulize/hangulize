import Hangul from 'hangul-js'
import { useState, useLayoutEffect, useRef } from 'react'
import './Result.css'

function Result({
  children = '',
}) {
  const [ zoom, setZoom ] = useState(1)
  const maxZoom = 1
  const minZoom = 0.5

  const p = useRef()
  const prevLength = useRef(0)
  const prevWidth = useRef(0)

  const updateSize = () => {
    if (p.current) {
      if (p.current.childNodes.length === 0) {
        return
      }

      // Make range to detect the number of rendered lines.
      const range = document.createRange()
      const textNode = p.current.firstChild
      range.setStart(textNode, 0)
      range.setEnd(textNode, textNode.length)

      const lines = range.getClientRects().length
      const length = textNode.length
      const width = document.body.offsetWidth

      // Try to size down on two or more lines.
      if (lines > 1) {
        console.log('down')
        setZoom(Math.max(minZoom, zoom * 0.99))
        prevLength.current = length
        return
      }

      // Try to size up on the text has been shrunken.
      if (length < prevLength.current || width > prevWidth.current) {
        console.log('up')
        setZoom(Math.min(maxZoom, zoom / 0.99))
        prevWidth.current = width
        return
      }

      prevLength.current = length
      prevWidth.current = width
    }
  }

  useLayoutEffect(() => {
    updateSize()
    window.addEventListener('resize', updateSize)
    return () => window.removeEventListener('resize', updateSize)
  })

  return (
    <div className="result">
      <p ref={p} style={{zoom: zoom}}>{children}</p>
    </div>
  )
}

export default Result
