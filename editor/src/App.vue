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

    <div id="editor"></div>

    <template v-for="(_, i) in words">
      <Transcription
        :key="i"
        :lang="selectedLang"
        @submit="onSubmit"
      />
    </template>
  </div>
</template>

<script>
import _ from 'lodash'
import ace from 'brace'

import HangulizeTheme from './ace/theme-hangulize'
import HGLMode from './ace/mode-hgl'

import H from 'hangulize'
import Transcription from './components/Transcription'

export default {
  name: 'App',

  components: {
    Transcription
  },

  data () {
    return {
      langs: [],
      selectedLang: '',

      word: '',
      delayedWord: '',

      words: [''],

      spec: null,
      editor: null
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
      if (!this.editor) {
        return
      }

      const spec = H.specs[lang]
      this.editor.session.setValue(spec.info.source)
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
      this.words.push('')
    }
  },

  created () {
    _.forEach(H.specs, (spec, lang) => {
      this.langs.push({ text: spec.info.lang.korean, value: lang })

      if (!this.selectedLang) {
        this.selectedLang = lang
      }
    })
  },

  mounted () {
    const editor = ace.edit('editor')

    editor.renderer.setShowGutter(false)
    editor.setTheme(HangulizeTheme)
    editor.session.setMode(new HGLMode())

    editor.session.on('change', () => {
      this.updateSource(editor.getValue())
    })

    this.editor = editor
  }
}
</script>

<style scoped>
@import 'https://fonts.googleapis.com/css?family=IBM+Plex+Mono';
@import 'https://fonts.googleapis.com/css?family=Nanum+Gothic+Coding';

#editor {
  font-family: 'IBM Plex Mono', 'Nanum Gothic Coding', monospace;
  height: 50em;
  position: absolute;
  top: 0;
  right: 0;
  width: 50%;
  height: 100%;
  font-size: 15px;
}
</style>
