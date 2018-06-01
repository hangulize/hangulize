'use strict';

let H = __hangulize__;

let app = new Vue({
  el: '.console',

  data: {
    langs: [],
    selectedLang: '',

    word: '',

    spec: undefined,
    source: '',
  },

  computed: {
    hangulized() {
      if (this.spec) {
        let h = H.NewHangulizer(this.spec);
        return h.Hangulize(this.word);
      } else {
        return H.Hangulize(this.selectedLang, this.word);
      }
    },
  },

  watch: {
    selectedLang(lang) {
      let specOK = H.LoadSpec(lang);
      let spec = specOK[0];

      this.source = spec.Source;
    },

    source(source) {
      clearTimeout(this._timeoutParseSpec);

      this._timeoutParseSpec = setTimeout(() => {
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
