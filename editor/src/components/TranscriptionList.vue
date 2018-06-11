<template>
  <div class="transcriptions">
    <template v-for="(t, i) in transcriptions">
      <Transcription :key="t.id" :index="i" />
    </template>
  </div>
</template>

<script>
import { mapState } from 'vuex'

import Transcription from './Transcription'

export default {
  name: 'TranscriptionList',

  components: {
    Transcription
  },

  computed: {
    ...mapState(['transcriptions'])
  },

  created () {
    if (this.transcriptions.length === 0) {
      this.$store.commit('insertTranscription')
      this.$store.commit('insertTranscription')
      this.$store.commit('insertTranscription')
      this.$nextTick(() => this.$store.commit('focusTranscription', 0))
    }
  }
}
</script>
