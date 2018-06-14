<template>
  <header>
    <p class="logo">
      <img src="static/logo.svg" width="196" />
    </p>

    <p>
      외국어 단어를
      <a href="https://ko.wikipedia.org/wiki/외래어_표기법">외래어 표기법</a>에
      따라 한글로 옮겨적습니다.

      <a href="https://github.com/sublee/hangulize2">
        <sui-icon name="github" />GitHub
      </a>
      <a href="https://godoc.org/github.com/sublee/hangulize2/hangulize">
        <sui-icon name="book" />GoDoc
      </a>
    </p>
  </header>
</template>

<style scoped>
p {
  text-align: center;
}

.logo {
  padding: 3rem;
  padding-bottom: 0.5rem;
}
</style>
