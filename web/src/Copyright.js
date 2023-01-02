function Copyright({ year }) {
  if (year === undefined) {
    year = new Date().getFullYear()
  }

  const brian = 'https://www.facebook.com/kkeutsori'
  const heungsub = 'https://subl.ee/'

  return (
    <p>
      &copy; 2010–{year} <a href={brian}>Brian</a> &amp;{' '}
      <a href={heungsub}>Heungsub</a>. All rights reserved.
    </p>
  )
}

export default Copyright
