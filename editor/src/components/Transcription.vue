<template>
  <div class="transcription">
    <span class="lang">
      <code>{{ spec.info.lang.id }}</code>
      {{ spec.info.lang.korean }}
    </span>
  </div>
</template>

<script>
import H from 'hangulize'

export default {
  name: 'Transcription',

  props: ['lang'],

  computed: {
    spec () {
      console.log(H.specs[this.lang])
      return H.specs[this.lang]
    }
  }
}
</script>
