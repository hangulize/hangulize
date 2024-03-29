# HSL

HSL(Hangulize Spec Language)은 한글라이즈에서 언어 별 전사 규칙을 표현하기 위한
설정 문법이다. HSL 파일의 미디어타입은 [`application/vnd.hsl`][media-type]이다.

[media-type]: https://www.iana.org/assignments/media-types/application/vnd.hsl


## 문법

HSL 문서는 여러 섹션으로 구성된다. 섹션엔 이름이 있다. 가령 이름이
`hangulize`일 경우 `hangulize:`으로 표시한다:

```
hangulize:
    ...
```

섹션은 둘 중 한 형식으로 구성된다. 하나는 사전형이다. 사전형 섹션은 `key =
value[, value2, ...]` 형식으로 키와 값의 연관 관계를 표현한다. 사전형 속에서 한
키는 단 한 번만 나타나야 하고 여러번 나타날 수 없다. 여기서 키들 사이의 순서는
중요하지 않다. 키는 하나의 문자열이고 값은 여러개의 문자열이다:

```
jaeum:
    "ㄱ" = "기역"
    "ㄴ" = "니은"

tossi:
    "은(는)" = "은", "는"
    "이(가)" = "이", "가"
```

다른 하나는 목록형이다. 목록형 섹션 역시 키와 값들의 연관을 표현하지만 그 관계는
작성한 순서대로 나열된다. 목록에서의 위치에 따라 의미가 다를 수 있어서 키는
중복돼도 괜찮다. `key -> value[, value2, ...]` 형식으로 표현한다:

```
romanize:
    "한" -> "han"
    "글" -> "gul", "geul"
```

`#`부터는 주석이다. 해석되지 않는다:

```
test:
    # 가장 단순한 테스트케이스
    "gloria" -> "글로리아"
```
