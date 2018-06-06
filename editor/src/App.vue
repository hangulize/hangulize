<template>
  <div id="app">
    <sui-dropdown
      placeholder="Language"
      selection
      :options="langs"
      v-model="selectedLang"
    />

    <sui-input name="word" v-model="word" />

    <h1 is="sui-header" class="hangulized">{{ hangulized.word }}</h1>

    <sui-list>
      <sui-list-item
        v-for="(trace, i) in hangulized.traces"
        v-bind:key="i"
      >{{ trace.Word }} {{ trace.Why }}</sui-list-item>
    </sui-list>

    <template v-for="(t, i) in transcriptions">
      <Transcription :key="i" :spec="spec" :lang.sync="t.lang" @submit="onSubmit" />
    </template>

    <Editor :source="source" @change="onSourceChange" />
  </div>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'
import Transcription from './components/Transcription'
import Editor from './components/Editor'

export default {
  name: 'App',

  components: {
    Transcription,
    Editor
  },

  data () {
    return {
      langs: [],
      selectedLang: '',

      word: '',
      delayedWord: '',

      transcriptions: [{lang: 'ita'}],

      spec: null,
      source: ''
    }
  },

  computed: {
    hangulized () {
      if (!this.spec) {
        return {}
      }

      const h = H.newHangulizer(this.spec)
      const hangulizedTraces = h.HangulizeTrace(this.delayedWord)

      return {
        word: hangulizedTraces[0],
        traces: hangulizedTraces[1]
      }
    }
  },

  watch: {
    selectedLang (lang) {
      this.spec = H.specs[lang]
      this.source = this.spec.source
    },

    word (word) {
      this.updateWord(word)
    }
  },

  methods: {
    updateSource: _.debounce(function (source) {
      this.spec = H.parseSpec(source)
    }, 300),

    updateWord: _.debounce(function (word) {
      this.delayedWord = word
    }, 300),

    onSubmit () {
      this.transcriptions.push({ lang: 'ita' })
    },

    onSourceChange (source) {
      this.updateSource(source)
    }
  },

  created () {
    _.forEach(H.specs, (spec, lang) => {
      this.langs.push({ text: spec.lang.korean, value: lang })

      if (!this.selectedLang) {
        this.selectedLang = lang
      }
    })
  }
}
</script>
