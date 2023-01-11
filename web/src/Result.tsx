import './Result.css'

import _ from 'lodash'
import { useLayoutEffect, useRef, useState } from 'react'

interface ResultProps {
  children: string
  loading: boolean
}

function Result({ children, loading }: ResultProps) {
  const result = (children || '').trim()

  const minSize = 3
  const maxSize = 7
  const sizePrec = 3
  const lo = useRef(minSize)
  const hi = useRef(maxSize)
  const [size, setSize] = useState(maxSize)

  // Start a binary search to find the best zoom.
  const reset = () => {
    lo.current = minSize
    hi.current = maxSize
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

    // Try to size down on two or more lines.
    if (lines > 1) {
      hi.current = size
      setSize(_.floor((size - lo.current) / 2 + lo.current, sizePrec))
      return
    }

    // Try to size up if there's still a room.
    const room = Math.pow(0.1, sizePrec) + Math.pow(0.1, sizePrec + 1)
    if (hi.current - lo.current > room) {
      lo.current = size
      setSize(_.ceil((hi.current - size) / 2 + size, sizePrec))
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
