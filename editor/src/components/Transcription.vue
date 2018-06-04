<template>
  <form class="transcription" @submit.prevent="onSubmit">
    <span class="lang">
      <code>{{ spec.info.lang.id }}</code>
      {{ spec.info.lang.korean }}
    </span>

    <input
      ref="word"
      v-model="word"
      :placeholder="example.word"
    />

    <span class="transcribed">{{ transcribed }}</span>
  </form>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'

export default {
  name: 'Transcription',

  props: ['lang', 'word'],

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
