<template>
  <div id="editor"></div>
</template>

<script>
import ace from 'brace'

import HangulizeTheme from '../ace/theme-hangulize'
import HGLMode from '../ace/mode-hgl'

export default {
  name: 'Editor',

  props: ['source'],

  watch: {
    source (source) {
      if (this.$editor) {
        this.$editor.session.setValue(source)
      }
    }
  },

  mounted () {
    const editor = ace.edit('editor')

    editor.renderer.setShowGutter(false)
    editor.setTheme(HangulizeTheme)
    editor.session.setMode(new HGLMode())

    editor.session.on('change', () => {
      this.$emit('change', editor.getValue())
    })

    this.$editor = editor
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
