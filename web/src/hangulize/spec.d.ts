interface Example {
  word: string
  result: string
}

interface Spec {
  lang: {
    id: string
    code2: string
    code3: string
    english: string
    korean: string
    script: string
    translit: string[]
  }

  test: Example[]
}

type Specs = { [lang: string]: Spec }

export type { Example, Spec, Specs }
