<template>
  <form class="transcription" @submit.prevent="onSubmit">
    <label>

      <span class="lang">
        <code>{{ spec.lang.id }}</code>
        {{ spec.lang.korean }}
      </span>

      <input
        ref="word"
        v-model="word"
        :placeholder="example.word"
        :class="'script-' + spec.lang.script"
      />

      <span class="transcribed">{{ transcribed }}</span>

    </label>
  </form>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'

export default {
  name: 'Transcription',

  props: ['spec'],

  data () {
    return {
      word: '',

      random: _.random(true)
    }
  },

  computed: {
    example () {
      const test = this.spec.test
      const i = _.floor(test.length * this.random)
      return test[i]
    },

    transcribed () {
      if (!this.word) {
        return this.example.transcribed
      }

      const h = H.newHangulizer(this.spec)
      return h.Hangulize(this.word)
    }
  },

  methods: {
    onSubmit (e) {
      this.$emit('submit')
    }
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
