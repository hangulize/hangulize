import { Grid, Icon } from 'semantic-ui-react'

function Copyright({ year }) {
  if (year === undefined) {
    year = new Date().getFullYear()
  }

  const brian = 'https://www.facebook.com/kkeutsori'
  const heungsub = 'https://subl.ee/'
  const github = 'https://github.com/hangulize/hangulize'

  return (
    <Grid className="copyright">
      <Grid.Column floated="left" width={14}>
        &copy; 2010–{year}{' '}
        <a href={brian} target="_blank" rel="noreferrer">Brian</a> &amp;{' '}
        <a href={heungsub} target="_blank" rel="noreferrer">Heungsub</a>.{' '}
        All rights reserved.
      </Grid.Column>
      <Grid.Column floated="right" width={2} align="right">
        <a href={github} target="_blank" rel="noreferrer">
          <Icon name="github" size="large" link color="grey" />
        </a>
      </Grid.Column>
    </Grid>
  )
}

export default Copyright
