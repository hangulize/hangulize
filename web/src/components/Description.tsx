import descriptions from '../descriptions.yaml'

export default function Description({ lang }: { lang: string }) {
  const desc = descriptions[lang as keyof typeof descriptions]

  if (desc === undefined) {
    return null
  }

  return <p>{desc}</p>
}
