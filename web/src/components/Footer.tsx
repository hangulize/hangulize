import { Grid, Icon } from 'semantic-ui-react'

interface FooterProps {
  year?: number
}

export default function Footer({ year }: FooterProps) {
  if (year === undefined) {
    year = new Date().getFullYear()
  }

  const brian = 'https://www.facebook.com/kkeutsori'
  const heungsub = 'https://subl.ee/'
  const github = 'https://github.com/hangulize/hangulize'

  return (
    <Grid className="footer">
      <Grid.Column floated="left" width={14} className="copyright">
        &copy; 2010–{year}{' '}
        <a href={brian} target="_blank" rel="noreferrer">
          Brian
        </a>{' '}
        &amp;{' '}
        <a href={heungsub} target="_blank" rel="noreferrer">
          Heungsub
        </a>
      </Grid.Column>
      <Grid.Column floated="right" width={2} align="right">
        <a href={github} target="_blank" rel="noreferrer">
          <Icon name="github" size="large" link color="grey" />
        </a>
      </Grid.Column>
    </Grid>
  )
}
