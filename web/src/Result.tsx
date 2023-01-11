import './Result.css'

import _ from 'lodash'
import { useLayoutEffect, useRef, useState } from 'react'

interface ResultProps {
  children: string
  loading: boolean
}

function Result({ children, loading }: ResultProps) {
  let result = (children || '').trim()
  if (loading && result === '') {
    result = '…'
  }

  const minZoom = 0.5
  const maxZoom = 1
  const lo = useRef(minZoom)
  const hi = useRef(maxZoom)
  const [zoom, setZoom] = useState(maxZoom)

  // Start a binary search to find the best zoom.
  const reset = () => {
    lo.current = minZoom
    hi.current = maxZoom
  }

  // Do a binary search before rendering to a user.
  const p = useRef<HTMLParagraphElement>(null)
  const update = () => {
    if (!p.current || !p.current.firstChild) {
      return
    }

    // Detect the number of rendered lines.
    const textNode = p.current.firstChild
    const range = document.createRange()
    range.setStart(textNode, 0)
    range.setEnd(textNode, (textNode.textContent as string).length)
    const lines = range.getClientRects().length

    if (lines > 1) {
      // Try to size down on two or more lines.
      hi.current = zoom
      setZoom(_.floor((zoom - lo.current) / 2 + lo.current, 2))
    } else if (hi.current - lo.current > 0.011) {
      // Try to size up if there's still a room.
      lo.current = zoom
      setZoom(_.ceil((hi.current - zoom) / 2 + zoom, 2))
    }
  }

  useLayoutEffect(reset, [children])
  useLayoutEffect(() => {
    update()

    const resetAndUpdate = () => {
      reset()
      update()
    }
    window.addEventListener('resize', resetAndUpdate)
    return () => window.removeEventListener('resize', resetAndUpdate)
  })

  return (
    <div className={`result ${loading ? 'loading' : ''}`}>
      {result ? (
        <p ref={p} style={{ fontSize: `${size}rem` }}>
          {result}
        </p>
      ) : (
        <></>
      )}
    </div>
  )
}

export default Result
