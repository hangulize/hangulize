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
    source: '',
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
    word(word, oldWord) {
      let delay = 100;
      if (word.length < oldWord.length && word.startsWith(oldWord.substr(0, word.length))) {
        delay = 300;
      }
 
      clearTimeout(this._timeoutWord);

      this._timeoutWord = setTimeout(() => {
        this.delayedWord = word;
      }, delay);
    },

    selectedLang(lang) {
      let specOK = H.LoadSpec(lang);
      let spec = specOK[0];

      this.source = spec.Source;
    },

    source(source) {
      clearTimeout(this._timeoutSpec);

      this._timeoutSpec = setTimeout(() => {
        this.spec = H.ParseSpec(source);
      }, 300);
    },
  }

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
