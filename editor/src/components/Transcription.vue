<template>
  <form
    class="transcription"
    :class="{ focused: focused }"
    @submit.prevent="onSubmit"
    tabindex="-1"
  >
    <Language
      :lang="lang"
      @input="updateLang"
    />

    <label>

      <input
        ref="word"
        :placeholder="example.word"
        :class="'script-' + spec.lang.script"
        @focus="focused = true"
        @blur="focused = false"
        @input="updateWord"
      />

      <span class="transcribed">{{ transcribed }}</span>

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

    transcribed: ''
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

    updateWord (e) {
      this.$store.commit('updateWord', {
        index: this.index,
        word: e.target.value
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
      this.transcribed = h.Hangulize(this.word || this.example.word)
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

form {
  background: #fff;
  margin: 0.5em;
  display: block;
  float: left;
  clear: left;
  opacity: 0.5;
}

form.focused {
  opacity: 1;
}

label {
  display: block;
  padding: 0 1em;
}

input {
  font-size: 2rem;
  font-weight: 600;
  line-height: 2;
  padding: 0 0.5em;
  border: none;
}

input.script-roman, input.script-cyrillic {
  font-family: 'Noto Sans', sans-serif;
}

input.script-roman, input.script-kana {
  font-family: 'Noto Sans JP', sans-serif;
}

.transcribed {
  font-family: 'Spoqa Han Sans', sans-serif;
  font-size: 1.75rem;
  font-weight: 400;
  color: #49f;
}
</style>
