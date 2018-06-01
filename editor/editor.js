'use strict';

let H = __hangulize__;

let app = new Vue({
  el: '.console',

  data: {
    langs: [],
    selectedLang: '',

    word: '',
    hangulized: '',

    spec: undefined,
    source: '',
  },

  watch: {
    word() {
      let h = H.NewHangulizer(this.spec);

      clearTimeout(this._timeoutHangulized);

      this._timeoutHangulized = setTimeout(() => {
        this.hangulized = h.Hangulize(this.word);
      }, 100);
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
