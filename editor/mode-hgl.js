define('ace/mode/hgl_highlight_rules', [
  'require',
  'exports',
  'module',
  'ace/lib/oop',
  'ace/mode/text_highlight_rules',
], function(require, exports, module) {
  'use strict';

  function HGLHighlightRules() {
    this.$rules = {
      start: [{
        token: 'storage',
        regex: /(rewrite|transcribe)(?=:)/,
        next: 'rewrite',
      }, {
        token: 'storage',
        regex: /[a-z_]+(?=:)/,
      }, {
        token: 'keyword.operator',
        regex: /=|->|,|:($|\s)/,
      }, {
        token: 'variable',
        regex: /^(?<=\s*)[^"#]+(?=\s*(?:\=|->))/,
      }, {
        token: 'string.double',
        regex: /"/,
        next: 'text',
      }, {
        token: 'string.unquoted',
        regex: /(?<=(?:\=|->|,)\s*)[^\s]+/,
      }, {
        token: 'comment',
        regex: /#.*/,
      }, {
        token: 'invalid',
        regex: /\s+$/,
      }, {
        token: 'invalid',
        regex: /\S+/,
      }],

      text: [{
        token: 'string.double',
        regex: /"/,
        next: 'start',
      }, {
        token: 'invalid',
        regex: /.$/,
        next: 'start',
      }, {
        defaultToken : 'string',
      }],

      rewrite: [{
        token: '',
        regex: /(?=.+:)/,
        next: 'start',
      }, {
        token: 'keyword.operator',
        regex: /=|->|,|:/,
      }, {
        token: 'string.double',
        regex: /"/,
        next: 'pattern',
      }, {
        token: 'comment',
        regex: /#.*/,
      }, {
        token: 'invalid',
        regex: /\s+$/,
      }, {
        token: 'invalid',
        regex: /\S+/,
      }],

      pattern: [{
        token: 'string.double',
        regex: /"|$/,
        next: 'rewrite',
      }, {
        token: 'invalid',
        regex: /.$/,
        next: 'start',
      }, {
        token: 'string.interpolated',
        regex: /{.+}/,
      }, {
        token: 'hangul.jongseong',
        regex: /-[\u3131-\u314e]/,
      }, {
        token: 'hangul.choseong',
        regex: /[\u3131-\u314e]/,
      }, {
        token: 'hangul.jungseong',
        regex: /[\u314f-\u3163]/,
      }, {
        defaultToken : 'string',
      }],
    };
  }

  HGLHighlightRules.metaData = {
    fileTypes: ['hgl'],
    name: 'HGL',
  };

  let oop = require('ace/lib/oop');
  let TextHighlightRules = require('./text_highlight_rules').TextHighlightRules;
  oop.inherits(HGLHighlightRules, TextHighlightRules);

  exports.HGLHighlightRules = HGLHighlightRules;

});

define('ace/mode/hgl', [
  'require',
  'exports',
  'module',
  'ace/lib/oop',
  'ace/mode/text',
  'ace/mode/hgl_highlight_rules',
  'ace/mode/folding/hgl',
], function(require, exports, module) {
  'use strict';

  let HGLHighlightRules = require('./hgl_highlight_rules').HGLHighlightRules;

  function HGLMode() {
    this.HighlightRules = HGLHighlightRules;
  }

  let oop = require('ace/lib/oop');
  let TextMode = require('./text').Mode;
  oop.inherits(HGLMode, TextMode);

  HGLMode.prototype.$id = 'ace/mode/hgl';
  HGLMode.prototype.type = 'hgl';

  HGLMode.prototype.getNextLineIndent = function(state, line, tab) {
    // Indent when section opened.
    if (/:($|\s|#)/.exec(line)) {
      return '    ';
    }
    return this.$getIndent(line);
  };

  exports.Mode = HGLMode;

});

(function() {
  window.require(['ace/mode/hgl'], function(m) {
    if (typeof module == 'object' && typeof exports == 'object' && module) {
      module.exports = m;
    }
  });
})();