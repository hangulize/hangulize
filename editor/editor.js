'use strict';

let H = __hangulize__;

let app = new Vue({
  el: '.console',

  data: {
    langs: [],
    selectedLang: '',

    word: '',
    hangulized: '',
    traces: [],

    spec: undefined,
    source: '',
  },

  watch: {
    word(word, oldWord) {
      let h = H.NewHangulizer(this.spec);

      clearTimeout(this._timeoutHangulized);

      let delay = 100;
      if (word.length < oldWord.length && word.startsWith(oldWord.substr(0, word.length))) {
        delay = 500;
      }
 
      this._timeoutHangulized = setTimeout(() => {
        let hangulizedTraces = h.HangulizeTrace(word);
        this.hangulized = hangulizedTraces[0];
        this.traces = hangulizedTraces[1];
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
      }, 500);
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
