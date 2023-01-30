import './Result.css'

import _ from 'lodash'
import { useLayoutEffect, useRef, useState } from 'react'

interface ResultProps {
  children: string
  loading: boolean
}

export default function Result({ children, loading }: ResultProps) {
  const result = (children || '').trim()

  const sizes = ['size-1', 'size-2', 'size-3', 'size-4', 'size-5']
  const lo = useRef(0)
  const hi = useRef(sizes.length - 1)
  const [mid, setMid] = useState(hi.current)

  // Start a binary search to find the best zoom.
  const reset = () => {
    lo.current = 0
    hi.current = sizes.length - 1
  }

  // Do a binary search before rendering to a user.
  const ref = useRef<HTMLDivElement>(null)
  const update = () => {
    if (!ref.current || !ref.current.firstChild) {
      return
    }

    // Try to size down if the height is expanded.
    const height = (ref.current.firstChild as HTMLParagraphElement).offsetHeight
    const capacity = parseInt(getComputedStyle(ref.current).maxHeight)
    if (height > capacity) {
      hi.current = mid
      setMid(_.floor((mid - lo.current) / 2) + lo.current)
      return
    }

    // The current size is optimal.
    if (hi.current - lo.current <= 1) {
      return
    }

    // Otherwise, try to size up.
    if (mid !== hi.current) {
      lo.current = mid
      setMid(_.ceil((hi.current - mid) / 2) + mid)
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

  if (!result) {
    return <></>
  }

  return (
    <div className={`result ${sizes[mid]} ${loading ? 'loading' : ''}`} ref={ref}>
      <p>{result}</p>
    </div>
  )
}
