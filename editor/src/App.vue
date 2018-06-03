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
  </div>
</template>

<script>
import _ from 'lodash'
import H from 'hangulize'
import ace from 'brace'

import HangulizeTheme from './ace/theme-hangulize'
import HGLMode from './ace/mode-hgl'

export default {
  name: 'App',

  data () {
    return {
      langs: [],
      selectedLang: '',

      word: '',
      delayedWord: '',

      spec: undefined,
      editor: undefined
    }
  },

  computed: {
    hangulized () {
      if (!this.spec) {
        return {}
      }

      const h = H.NewHangulizer(this.spec)
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

      const specOK = H.LoadSpec(lang)
      const spec = specOK[0]

      this.editor.session.setValue(spec.Source)
    },

    word (word) {
      this.updateWord(word)
    }
  },

  methods: {
    updateSource: _.debounce(function (source) {
      this.spec = H.ParseSpec(source)
    }, 300),

    updateWord: _.debounce(function (word) {
      this.delayedWord = word
    }, 300)
  },

  created () {
    H.ListLangs().forEach((langID) => {
      this.langs.push({ text: langID, value: langID })

      if (!this.selectedLang) {
        this.selectedLang = langID
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
