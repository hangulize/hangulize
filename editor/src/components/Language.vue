<template>
  <div class="language">
    <sui-dropdown
      placeholder="언어"
      selection
      :options="langs"
      v-model="selected_"
    />
  </div>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'

export default {
  name: 'Language',

  props: ['selected'],

  data () {
    return {
      selected_: null,

      langs: _.map(H.specs, (spec, lang) => ({
        text: `${spec.lang.id} ${spec.lang.korean}`,
        value: lang
      }))
    }
  },

  watch: {
    selected_ (selected) {
      this.$emit('update:selected', selected)
    }
  },

  created () {
    this.selected_ = this.selected
  }
}
</script>
