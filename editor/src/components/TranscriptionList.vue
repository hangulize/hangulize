<template>
  <div class="transcriptions">
    <template v-for="(t, i) in transcriptions">
      <Transcription :key="t.id" :index="i" @submit="onSubmit(i)" />
    </template>
  </div>
</template>

<script>
import { mapMutations, mapState } from 'vuex'

import Transcription from './Transcription'

export default {
  name: 'TranscriptionList',

  components: {
    Transcription
  },

  computed: {
    ...mapState(['transcriptions'])
  },

  methods: {
    onSubmit (i) {
      this.insertTranscription(i + 1)
    },

    ...mapMutations(['insertTranscription'])
  },

  created () {
    this.insertTranscription()
  }
}
</script>
