'use strict';

let H = __hangulize__;

let app = new Vue({
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

      let h = H.NewHangulizer(this.spec);
      let hangulizedTraces = h.HangulizeTrace(this.delayedWord);

      return {
        word:   hangulizedTraces[0],
        traces: hangulizedTraces[1],
      };
    },
  },

  watch: {
    selectedLang(lang) {
      let specOK = H.LoadSpec(lang);
      let spec = specOK[0];

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
    updateSource: _.debounce(function(source) {
      this.spec = H.ParseSpec(source);
    }, 300),

    updateWord: _.debounce(function(word) {
      this.delayedWord = word;
    }, 300),
  },

});

function loadLangs() {
  H.ListLangs().forEach((langID) => {
    app.langs.push({id: langID});

    if (!app.selectedLang) {
      app.selectedLang = langID;
    }
  });
}

loadLangs();

var editor = ace.edit("editor");
editor.renderer.setShowGutter(false);
editor.setTheme("ace/theme/hangulize");
editor.session.setMode("ace/mode/hgl");
editor.session.on('change', function(delta) {
  app.updateSource(editor.getValue());
});
