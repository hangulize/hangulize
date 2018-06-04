<template>
  <div class="transcription" tabindex="-1">
    <span class="flag-icon" :class="'flag-icon-' + country"></span>
  </div>
</template>

<script>
import 'flag-icon-css/css/flag-icon.css'

import H from 'hangulize'

export default {
  name: 'Transcription',

  props: ['lang'],

  computed: {
    country () {
      const lang = H.Langs[this.lang]
      return lang.Codes[0]
    }
  }
}
</script>
