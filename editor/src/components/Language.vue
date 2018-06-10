<template>
  <div class="language">
    <Lang tag="div" :lang="lang" @click.stop="selecting = true">
      <sui-icon name="dropdown" />
    </Lang>

    <LanguageSelector
      v-if="selecting"
      :lang="lang"
      @input="(lang) => value = lang"
      @cancel="selecting = false"
    />
  </div>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'

import Lang from './Lang'
import LanguageSelector from './LanguageSelector'

export default {
  name: 'Language',

  components: {
    Lang,
    LanguageSelector
  },

  props: ['lang'],

  data: () => ({
    value: '',

    selecting: false,

    langs: _.map(H.specs, (spec, lang) => ({
      text: `${spec.lang.id} ${spec.lang.korean}`,
      value: lang
    }))
  }),

  computed: {
    id () {
      return this.value.toUpperCase()
    },

    name () {
      return H.specs[this.value].lang.korean
    }
  },

  watch: {
    value (value) {
      this.$emit('input', value)
      this.selecting = false
    },

    selecting (selecting) {
      if (selecting) {
        this.$emit('open')
      } else {
        this.$emit('close')
      }
    }
  },

  created () {
    this.value = this.lang
  }
}
</script>

<style>
.language {
  position: relative;
}

.language > .lang {
  position: relative;
  display: inline-block;
  line-height: 2.5em;
  cursor: pointer;
}

.language > .lang i {
  color: #ccc;
}

.language > .lang:hover i {
  color: #000;
}
</style>
