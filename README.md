<a href="https://hangulize.org/">
  <img src="web/public/logo.svg" height="64" align="right" />
</a>

# 한글라이즈

[![GoDoc](https://godoc.org/github.com/hangulize/hangulize?status.svg)](https://godoc.org/github.com/hangulize/hangulize)
[![Go Report Card](https://goreportcard.com/badge/github.com/hangulize/hangulize)](https://goreportcard.com/report/github.com/hangulize/hangulize)
[![Build Status](https://travis-ci.org/hangulize/hangulize.svg?branch=develop)](https://travis-ci.org/hangulize/hangulize)
[![Coverage Status](https://coveralls.io/repos/github/hangulize/hangulize/badge.svg?branch=develop)](https://coveralls.io/github/hangulize/hangulize)

(WIP: 아직 개발 중, API가 임의로 바뀔 수 있어요!)

> 외국어의 한글 표기 체계가 제대로 서려면 일반인이 외국어를 한글로 표기하고
> 싶을 때 바로바로 쉽게 용례를 찾을 수 있어야 한다. 정기적으로 회의를 열어
> 용례를 정하는 것으로는 한계가 있다. 외래어 표기 심의 방식이 자동화되어 한글로
> 표기하고 싶은 외국어를 입력하자마자 한글 표기가 나와야 한다. 이미 용례가
> 정해진 것은 그것을 따르고 용례에 없는 것이라도 각 언어의 표기 규칙에 따라
> 권장 표기를 표시해야 한다. 프로그래머들과 언어학자들이 손잡고 연구한다면 이게
> 공상으로만 그치지 않을 것이다.
>
> Brian Jongseong Park (http://iceager.egloos.com/2610028)

한글라이즈는 외래어를 한글로 변환하는 도구입니다.

```console
$ go get -u github.com/hangulize/hangulize
```

```go
import "github.com/hangulize/hangulize"

hangulize.Hangulize("ita", "Cappuccino")
// output: "카푸치노"
```

## 지원하는 언어

```
LANG     STAGE    ENG                      KOR
aze      draft    Azerbaijani              아제르바이잔어
bel      draft    Belarusian               벨라루스어
bul      draft    Bulgarian                불가리아어
cat      draft    Catalan                  카탈로니아어
ces      draft    Czech                    체코어
chi      draft    Chinese                  중국어
cym      draft    Welsh                    웨일스어
deu      draft    German                   독일어
ell      draft    Greek                    그리스어
epo      draft    Esperanto                에스페란토어
est      draft    Estonian                 에스토니아어
fin      draft    Finnish                  핀란드어
grc      draft    Ancient Greek            고대 그리스어
hbs      draft    Serbo-Croatian           세르보크로아트어
hun      draft    Hungarian                헝가리어
isl      draft    Icelandic                아이슬란드어
ita      draft    Italian                  이탈리아어
jpn      draft    Japanese                 일본어
jpn-ck   draft    Japanese (C.K.)          일본어(최영애-김용옥)
kat-1    draft    Georgian (1st scheme)    조지아어(제1안)
kat-2    draft    Georgian (2nd scheme)    조지아어(제2안)
lat      draft    Latin                    라틴어
lav      draft    Latvian                  라트비아어
lit      draft    Lithuanian               리투아니아어
mkd      draft    Macedonian               마케도니아어
nld      draft    Dutch                    네덜란드어
pol      draft    Polish                   폴란드어
por      draft    Portuguese               포르투갈어
por-br   draft    Brazilian Portuguese     브라질 포르투갈어
ron      draft    Romanian                 루마니아어
rus      draft    Russian                  러시아어
slk      draft    Slovak                   슬로바키아어
slv      draft    Slovenian                슬로베니아어
spa      draft    Spanish                  스페인어
sqi      draft    Albanian                 알바니아어
swe      draft    Swedish                  스웨덴어
tur      draft    Turkish                  터키어
ukr      draft    Ukrainian                우크라이나어
vie      draft    Vietnamese               베트남어
wlm      draft    Middle Welsh             웨일스어(중세)
```

## 읽을거리

- [한글라이즈 재제작기][remake-of-hangulize](이흥섭, 고랭코리아 2018년 8월 밋업)

[remake-of-hangulize]: https://subl.ee/~gokr1808

## 만든이

- 이흥섭, Heungsub Lee <<heungsub@subl.ee>>
- 박종성, Brian Jongseong Park <<iceager@gmail.com>>

## 라이선스

한글라이즈는 MIT 라이선스 하에 공개되어 있습니다. 소스코드를 사용할 경우
라이선스 내용을 준수해주세요. 라이선스 전문은 `LICENSE` 파일에서 확인하실 수
있습니다.
