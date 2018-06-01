'use strict';

let H = __hangulize__;

let app = new Vue({
  el: '.console',

  data: {
    langs: [],
    selectedLang: '',

    word: '',
  },

  computed: {
    hangulized() {
      return H.Hangulize(this.selectedLang, this.word);
    },
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
