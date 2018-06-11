import _ from 'lodash'
import Vue from 'vue'
import Vuex from 'vuex'

import H from 'hangulize'

Vue.use(Vuex)

class Transcription {
  constructor (id, lang) {
    this.id = id
    this.lang = lang
    this.spec = H.specs[lang]
    this.word = ''
  }
}

export default new Vuex.Store({
  state: () => ({
    transcriptions: [],
    focusedTranscriptionID: null,
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

      const id = state.nextTranscriptionID++
      const t = new Transcription(id, lang)

      state.transcriptions.splice(index, 0, t)
    },

    removeTranscription (state, index = 0) {
      state.transcriptions.splice(index, 1)
    },

    focusTranscription (state, index = 0) {
      index = _.clamp(index, 0, state.transcriptions.length - 1)
      const id = state.transcriptions[index].id
      state.focusedTranscriptionID = id
    },

    blurTranscriptions (state) {
      state.focusedTranscriptionID = null
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
