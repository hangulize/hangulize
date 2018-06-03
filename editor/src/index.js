import _ from 'lodash';
import Vue from 'vue/dist/vue.common';
import ace from 'brace';

import HangulizeTheme from './ace-theme-hangulize';
import { Mode as HGLMode } from './ace-mode-hgl';

import H from '../hangulize';

const editor = ace.edit('editor');
editor.renderer.setShowGutter(false);
editor.setTheme(HangulizeTheme);
editor.session.setMode(new HGLMode());

const app = new Vue({
  el: '.console',

  data: {
    langs: [],
    selectedLang: '',

    word: '',
    delayedWord: '',

    spec: undefined,
    // source: '',
  },

  computed: {
    hangulized() {
      if (!this.spec) {
        return {};
      }

      const h = H.NewHangulizer(this.spec);
      const hangulizedTraces = h.HangulizeTrace(this.delayedWord);

      return {
        word: hangulizedTraces[0],
        traces: hangulizedTraces[1],
      };
    },
  },

  watch: {
    selectedLang(lang) {
      const specOK = H.LoadSpec(lang);
      const spec = specOK[0];

      // this.source = spec.Source;
      editor.session.setValue(spec.Source);
    },

    // source(source) {
    //   this.updateSource(source);
    // },

    word(word) {
      this.updateWord(word);
    },
  },

  methods: {
    updateSource: _.debounce(function updateSource(source) {
      this.spec = H.ParseSpec(source);
    }, 300),

    updateWord: _.debounce(function updateWord(word) {
      this.delayedWord = word;
    }, 300),
  },

});

function loadLangs() {
  H.ListLangs().forEach((langID) => {
    app.langs.push({ id: langID });

    if (!app.selectedLang) {
      app.selectedLang = langID;
    }
  });
}

loadLangs();

editor.session.on('change', () => {
  app.updateSource(editor.getValue());
});
