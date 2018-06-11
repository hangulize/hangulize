<template>
  <div id="app">
    <Logo />
    <TranscriptionList />
  </div>
</template>

<script>
import store from './store/index'
import Logo from './components/Logo'
import Transcription from './components/Transcription'
import TranscriptionList from './components/TranscriptionList'

export default {
  name: 'App',
  store,

  components: {
    Logo,
    Transcription,
    TranscriptionList
  }
}
</script>

<style>
@import url('https://spoqa.github.io/spoqa-han-sans/css/SpoqaHanSans-kr.css');

body {
  background: #ece9da;
  padding: 3rem 0;
  padding-top: 0;
  font-family: 'Spoqa Han Sans', sans-serif;
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
