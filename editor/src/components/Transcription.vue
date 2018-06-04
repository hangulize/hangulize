<template>
  <div class="transcription">
    <span class="lang">
      <code>{{ spec.info.lang.id }}</code>
      {{ spec.info.lang.korean }}
    </span>

    <input v-model="word" :placeholder="example.word" />

    <span class="transcribed">{{ transcribed }}</span>
  </div>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'

export default {
  name: 'Transcription',

  props: ['lang'],

  data () {
    return {
      word: ''
    }
  },

  computed: {
    spec () {
      return H.specs[this.lang]
    },

    example () {
      const test = this.spec.info.test
      const i = _.random(test.length)
      return test[i]
    },

    transcribed () {
      if (!this.word) {
        return this.example.transcribed
      }

      const h = H.newHangulizer(this.spec)
      return h.Hangulize(this.word)
    }
  }
}
</script>
