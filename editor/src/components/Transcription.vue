<template>
  <form
    class="transcription"

    tabindex="-1"

    :class="{ focused, selecting }"
    @focus="focused = true"
    @blur="focused = false"

    @submit.prevent="onSubmit"
  >
    <label>

      <Language
        :lang="lang"
        @input="updateLang"
        @open="selecting = true"
        @close="selecting = false"
      />

      <input
        ref="word"
        :placeholder="example.word"
        :value="word"
        :class="'script-' + spec.lang.script"
        @focus="focused = true"
        @blur="focused = false"
        @input="(e) => updateWord(e.target.value)"
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
    focused: false,
    selecting: false,

    transcribed: '',
    exampleTranscribed: true
  }),

  computed: {
    transcription () {
      return this.$store.getters.getTranscription(this.index)
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

    onSubmit () {
      this.$emit('submit', this)
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
    this.$refs.word.focus()
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
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.2);
}
form.selecting {
  box-shadow: none;
  background: #f4f4f4;
}
</style>
