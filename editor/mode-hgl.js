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
        token: 'comment.number-sign',
        regex: /#.*/,
      }, {
        token: 'storage.type',
        regex: /(rewrite|transcribe)(?=:)/,
        next: 'rewrite',
      }, {
        token: 'storage.type',
        regex: /.+(?=:)/,
      }, {
        token: 'keyword.operator',
        regex: /=|->|,|:/,
      }, {
        token: 'variable',
        regex: /^(?<=\s*)[^"]+(?=\s*(?:\=|->))/,
      }, {
        token: 'string.double',
        regex: /"/,
        next: 'text',
      }, {
        token: 'string.unquoted',
        regex: /(?<=(?:\=|->|,)\s*)[^\s]+/,
      }],

      text: [{
        token: 'string.double',
        regex: /"/,
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
      }],

      pattern: [{
        token: 'string.double',
        regex: /"/,
        next: 'rewrite',
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
        token: 'hangul',
        regex: /[\uac00-\ud7a3]/,
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