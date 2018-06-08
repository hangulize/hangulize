import _ from 'lodash'
import Vue from 'vue'
import Vuex from 'vuex'

import H from 'hangulize'

Vue.use(Vuex)

export default new Vuex.Store({
  state: () => ({
    transcriptions: [],
    nextTranscriptionID: 0
  }),

  getters: {
    getTranscription: (state) => (i) => {
      return state.transcriptions[i]
    }
  },

  mutations: {
    // Inserts a transcription onto the given index.
    insertTranscription (state, index = 0) {
      let lang

      if (index === 0) {
        // Pick a random lang for intiializing.
        const langs = Object.keys(H.specs)
        const i = _.random(langs.length)
        lang = langs[i]
      } else {
        // Use lang of the prev transcription as default.
        lang = state.transcriptions[index - 1].lang
      }

      const t = {
        id: state.nextTranscriptionID++,
        lang: lang,
        spec: H.specs[lang],
        word: ''
      }

      state.transcriptions.splice(index, 0, t)
    },

    updateLang (state, {index, lang}) {
      state.transcriptions[index].lang = lang
      state.transcriptions[index].spec = H.specs[lang]
    },

    updateSpec (state, {index, spec}) {
      state.transcriptions[index].spec = spec
    },

    updateWord (state, {index, word}) {
      state.transcriptions[index].word = word
    }
  }
})
