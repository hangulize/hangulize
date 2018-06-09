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

<style scoped>
#app {
  min-width: 700px;
  width: 50%;
  margin: 0 auto;
}
</style>

<style>
@import url('https://spoqa.github.io/spoqa-han-sans/css/SpoqaHanSans-kr.css');

body {
  background: #ece9da;
  padding: 50px;
  font-family: 'Spoqa Han Sans', sans-serif;
}
</style>
