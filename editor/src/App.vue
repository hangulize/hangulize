<template>
  <div id="app">
    <Header />
    <TranscriptionList />
    <Footer />
  </div>
</template>

<script>
import store from './store/index'
import Header from './components/Header'
import TranscriptionList from './components/TranscriptionList'
import Footer from './components/Footer'

export default {
  name: 'App',
  store,

  components: {
    Header,
    TranscriptionList,
    Footer
  }
}
</script>

<style>
@import url('https://spoqa.github.io/spoqa-han-sans/css/SpoqaHanSans-kr.css');

body {
  background: #ece9da;
  font-family: 'Spoqa Han Sans', sans-serif;
}

header {
  padding-bottom: 1rem;
}

footer {
  padding-top: 1rem;
}
</style>

<style scoped>
#app {
  width: 90%;
  min-width: 400px;
  max-width: 600px;
  margin: 0 auto;
}
</style>
