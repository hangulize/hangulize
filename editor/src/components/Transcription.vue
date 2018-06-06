<template>
  <form
    class="transcription"
    :class="{ focused: focused }"
    @submit.prevent="onSubmit"
  >
    <Language :selected.sync="lang" />

    <label>

      <input
        ref="word"
        v-model="word"
        :placeholder="example.word"
        :class="'script-' + spec.lang.script"
        @focus="focused = true"
        @blur="focused = false"
      />

      <span class="transcribed">{{ transcribed }}</span>

    </label>
  </form>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'

import Language from './Language'

export default {
  name: 'Transcription',

  components: {
    Language
  },

  props: ['spec', 'lang'],

  data () {
    return {
      word: '',
      transcribed: '',

      focused: false,

      random: _.random(true)
    }
  },

  computed: {
    example () {
      const test = this.spec.test
      const i = _.floor(test.length * this.random)
      return test[i]
    }
  },

  watch: {
    lang (lang) {
      this.spec = H.specs[lang]
      this.$emit('update:lang', lang)
    },

    spec () {
      this.hangulize()
    },

    word () {
      this.hangulize()
    }
  },

  methods: {
    hangulize () {
      const h = H.newHangulizer(this.spec)
      this.transcribed = h.Hangulize(this.word || this.example.word)
    },

    onSubmit (e) {
      this.$emit('submit')
    }
  },

  created () {
    this.hangulize()
  },

  mounted () {
    this.$refs.word.focus()
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css?family=Noto+Sans&subset=cyrillic,greek,vietnamese');
@import url('https://fonts.googleapis.com/earlyaccess/notosansjp.css');
@import url('https://spoqa.github.io/spoqa-han-sans/css/SpoqaHanSans-kr.css');

form {
  background: #fff;
  margin: 0.5em;
  display: block;
  float: left;
  clear: left;
  opacity: 0.5;
}

form.focused {
  opacity: 1;
}

label {
  display: block;
  padding: 0 1em;
}

input {
  font-size: 2rem;
  font-weight: 600;
  line-height: 2;
  padding: 0 0.5em;
  border: none;
}

input.script-roman, input.script-cyrillic {
  font-family: 'Noto Sans', sans-serif;
}

input.script-roman, input.script-kana {
  font-family: 'Noto Sans JP', sans-serif;
}

.transcribed {
  font-family: 'Spoqa Han Sans', sans-serif;
  font-size: 1.75rem;
  font-weight: 400;
  color: #49f;
}
</style>
