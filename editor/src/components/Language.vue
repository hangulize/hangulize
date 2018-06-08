<template>
  <div class="language">
    <sui-dropdown
      placeholder="언어"
      selection
      :options="langs"
      v-model="value"
      @input="$emit('input', value)"
    />
  </div>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'

export default {
  name: 'Language',

  props: ['lang'],

  data: () => ({
    value: '',

    langs: _.map(H.specs, (spec, lang) => ({
      text: `${spec.lang.id} ${spec.lang.korean}`,
      value: lang
    }))
  }),

  created () {
    this.value = this.lang
  }
}
</script>
