import { TextHighlightRules }
  from 'ace-code-editor/lib/ace/mode/text_highlight_rules'
import { Mode as TextMode }
  from 'ace-code-editor/lib/ace/mode/text'

class HGLHighlightRules extends TextHighlightRules {
  constructor () {
    super()
    this.$rules = {
      start: [{
        token: 'storage',
        regex: /(rewrite|transcribe)(?=:)/,
        next: 'rewrite'
      }, {
        token: 'storage',
        regex: /[a-z_]+(?=:)/
      }, {
        token: 'keyword.operator',
        regex: /=|->|,|:($|\s)/
      }, {
        token: 'variable',
        regex: /^\s*[^"#]+(?=\s*(?:=|->))/
      }, {
        token: 'string.double',
        regex: /"/,
        next: 'text'
      }, {
        token: 'string.unquoted',
        regex: /(?:=|->|,)\s*[^\s]+/
      }, {
        token: 'comment',
        regex: /#.*/
      }, {
        token: 'invalid',
        regex: /\s+$/
      }, {
        token: 'invalid',
        regex: /\S+/
      }],

      text: [{
        token: 'string.double',
        regex: /"/,
        next: 'start'
      }, {
        token: 'invalid',
        regex: /.$/,
        next: 'start'
      }, {
        defaultToken: 'string'
      }],

      rewrite: [{
        token: '',
        regex: /(?=.+:)/,
        next: 'start'
      }, {
        token: 'keyword.operator',
        regex: /=|->|,|:/
      }, {
        token: 'string.double',
        regex: /"/,
        next: 'pattern'
      }, {
        token: 'comment',
        regex: /#.*/
      }, {
        token: 'invalid',
        regex: /\s+$/
      }, {
        token: 'invalid',
        regex: /\S+/
      }],

      pattern: [{
        token: 'string.double',
        regex: /"|$/,
        next: 'rewrite'
      }, {
        token: 'invalid',
        regex: /.$/,
        next: 'start'
      }, {
        token: 'string.interpolated',
        regex: /{.+}/
      }, {
        token: 'hangul.jongseong',
        regex: /-[\u3131-\u314e]/
      }, {
        token: 'hangul.choseong',
        regex: /[\u3131-\u314e]/
      }, {
        token: 'hangul.jungseong',
        regex: /[\u314f-\u3163]/
      }, {
        defaultToken: 'string'
      }]
    }
  }
}

class HGLMode extends TextMode {
  constructor () {
    super()
    this.HighlightRules = HGLHighlightRules
    this.$id = 'ace/mode/hgl'
    this.type = 'hgl'
  }

  getNextLineIndent (state, line, tab) {
    // Indent when section opened.
    if (/:($|\s|#)/.exec(line)) {
      return '    '
    }
    return this.$getIndent(line)
  }
}

export default HGLMode
