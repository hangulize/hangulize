<template>
  <div class="transcriptions">
    <GlobalEvents
      @keydown.up="focusLast"
      @keydown.down="focusFirst"
    />

    <template v-for="(t, i) in transcriptions">
      <Transcription :key="t.id" :index="i" />
    </template>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import GlobalEvents from 'vue-global-events'

import Transcription from './Transcription'

export default {
  name: 'TranscriptionList',

  components: {
    Transcription,
    GlobalEvents
  },

  computed: {
    ...mapState(['transcriptions'])
  },

  methods: {
    focus (index) {
      this.$store.commit('focusTranscription', index)
    },

    focusFirst (e) {
      if (e.target === document.body) {
        this.focus(0)
      }
    },

    focusLast (e) {
      if (e.target === document.body) {
        this.focus(this.transcriptions.length - 1)
      }
    }
  },

  created () {
    if (this.transcriptions.length === 0) {
      this.$store.commit('insertTranscription')
      this.$store.commit('insertTranscription')
      this.$store.commit('insertTranscription')
      this.$nextTick(() => this.focus(0))
    }
  }
}
</script>
