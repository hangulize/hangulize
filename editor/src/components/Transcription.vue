<template>
  <form
    class="transcription"

    tabindex="-1"
    :class="{ focused, selecting }"

    @focus="focus"
    @blur="blur"
    @submit.prevent="insert"
  >
    <label>

      <Language
        :lang="lang"
        @input="updateLang"
        @open="selecting = true"
        @close="selecting = false"
      />

      <input
        ref="input"

        :placeholder="example.word"
        :value="word"
        :class="'script-' + spec.lang.script"

        @input="(e) => updateWord(e.target.value)"
        @focus="focus"
        @blur="blur"
        @keydown.up="focusAbove"
        @keydown.down="focusBelow"
        @keydown.backspace="maybeRemove"
      />

      <span
        class="transcribed"
        :class="{ example: exampleTranscribed }"
      >
        {{ transcribed }}
      </span>

    </label>
  </form>
</template>

<script>
import _ from 'lodash'

import H from 'hangulize'

import Language from './Language'

export default {
  name: 'Transcription',

  components: {
    Language
  },

  props: ['index'],

  data: () => ({
    random: _.random(true),
    selecting: false,

    transcribed: '',

    // Whether the transcribed is from an example.
    exampleTranscribed: true
  }),

  computed: {
    input () {
      return this.$refs.input
    },

    transcription () {
      return this.$store.getters.getTranscription(this.index)
    },

    id () {
      return this.transcription.id
    },

    lang () {
      return this.transcription.lang
    },

    spec () {
      return this.transcription.spec
    },

    word () {
      return this.transcription.word
    },

    example () {
      const test = this.spec.test
      const i = _.floor(test.length * this.random)
      return test[i]
    },

    focused () {
      return this.id === this.$store.state.focusedTranscriptionID
    }
  },

  watch: {
    focused (focused) {
      if (focused) {
        this.input.select()
      }
    }
  },

  methods: {
    hangulize () {
      // Will be implemented at created().
    },

    updateLang (lang) {
      this.$store.commit('updateLang', {
        index: this.index,
        lang: lang
      })

      this.transcribed = ''
      this.hangulize()
    },

    updateWord (word) {
      this.$store.commit('updateWord', {
        index: this.index,
        word: word
      })
      this.hangulize()
    },

    insert () {
      this.$store.commit('insertTranscription', this.index + 1)
    },

    maybeRemove () {
      if (this.index === 0) {
        return
      }

      if (this.input.value !== '') {
        return
      }

      this.$store.commit('removeTranscription', this.index)
      this.focusAbove()
    },

    focus () {
      this.$store.commit('focusTranscription', this.index)
    },

    focusAbove () {
      this.$store.commit('focusTranscription', this.index - 1)
    },

    focusBelow () {
      this.$store.commit('focusTranscription', this.index + 1)
    },

    blur () {
      if (this.focused) {
        this.$store.commit('blurTranscriptions')
      }
    }
  },

  created () {
    // NOTE(sublee): hangulize() is expensive.  If we call this for every user
    // input, the user experience would be bad.  So we need to wrap it with
    // _.debounce().
    //
    // But if we define it in the methods, all Transcription component
    // instances share the same debounce schedule.  We define here instead to
    // separate schedules for each component instance.
    //
    // https://forum.vuejs.org/t/issues-with-vuejs-component-and-debounce/7224/13
    //
    this.hangulize = _.debounce(() => {
      const h = H.newHangulizer(this.spec)

      if (this.word) {
        this.exampleTranscribed = false
        this.transcribed = h.Hangulize(this.word)
        return
      }

      this.exampleTranscribed = true
      if (this.spec === H.specs[this.lang]) {
        this.transcribed = this.example.transcribed
      } else {
        this.transcribed = h.Hangulize(this.example.word)
      }
    }, 100)

    this.hangulize()
  },

  mounted () {
    this.focus()
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css?family=Noto+Sans&subset=cyrillic,greek,vietnamese');
@import url('https://fonts.googleapis.com/earlyaccess/notosansjp.css');
@import url('https://spoqa.github.io/spoqa-han-sans/css/SpoqaHanSans-kr.css');

*:focus {
  outline: none;
}

form {
  background: #fff;
  margin: 0.5em;
  display: block;
}

label {
  display: block;
  padding: 0 1em;
}

input {
  background: transparent;
  font-size: 2rem;
  font-weight: 600;
  line-height: 1;
  padding: 0;
  border: none;
  width: 20em;
}

input::placeholder {
  color: #aaa;
}

input::selection {
  background: #bdf;
}

input.script-roman, input.script-cyrillic {
  font-family: 'Noto Sans', sans-serif;
}

input.script-roman, input.script-kana {
  font-family: 'Noto Sans JP', sans-serif;
}

.transcribed {
  display: block;
  padding: 0.5em 0;
  font-family: 'Spoqa Han Sans', sans-serif;
  font-size: 1.75rem;
  font-weight: 400;
  color: #49f;
}

.transcribed.example {
  color: #abd;
}

form.focused {
  box-shadow: 0 3px 10px rgba(68, 51, 34, 0.2);
  outline: 1px solid rgba(68, 51, 34, 0.05);
}
form.selecting {
  box-shadow: none;
  background: #f4f4f4;
}
</style>
