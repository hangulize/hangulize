# -*- coding: utf-8 -*-
from __future__ import print_function, unicode_literals

from collections import defaultdict
import io
import sys

import babel
import hangulize


def quote(x):
    return '"%s"' % x


def stringify(x):
    if not x:
        return ''

    if isinstance(x, tuple):
        return ''.join(map(stringify, x))

    if isinstance(x, unicode):
        return x

    if isinstance(x, bytes):
        return x.decode('utf-8')

    if isinstance(x, hangulize.Choseong):
        return x.letter

    if isinstance(x, hangulize.Jungseong):
        return x.letter

    if isinstance(x, hangulize.Jongseong):
        return '-' + x.letter


class Section(object):

    def __init__(self, name):
        self.name = name
        self.pairs = []

    def put(self, k, *vs):
        pair = (k, map(stringify, vs))
        self.pairs.append(pair)

    def draw(self, sep, quote_keys=False):
        pairs = self.pairs[:]

        if not pairs:
            return ''

        if quote_keys:
            pairs = [(quote(k), vs) for k, vs in pairs]

        key_width = max(len(k) for k, vs in pairs)
        tmpl = '{0:%ds} {1} {2}' % key_width

        buf = io.StringIO()

        buf.write(self.name)
        buf.write(':\n')

        for k, vs in pairs:
            buf.write('    ')
            buf.write(tmpl.format(k, sep, quote(vs[0])))
            for v in vs[1:]:
                buf.write(', ')
                buf.write(quote(v))
            buf.write('\n')

        buf.write('\n')
        return buf.getvalue().encode('utf-8')


def main(argv):
    try:
        code = argv[1]
    except IndexError:
        print('Usage 1to2.py LANG')
        raise SystemExit(2)

    lang = hangulize.get_lang(code)
    locale = babel.Locale(lang.iso639_1)

    # detect normalize
    additional_of_normalize_roman = {}
    normalize_roman_called = []

    def hacked_normalize_roman(string, additional=None):
        if additional:
            additional_of_normalize_roman.update(additional)
        normalize_roman_called.append(1)

    normalize_f = lang.normalize.__func__
    normalize_f.__globals__['normalize_roman'] = hacked_normalize_roman
    normalize_f(lang, '')

    normalize = defaultdict(set)
    for src, dst in additional_of_normalize_roman.items():
        assert dst == dst.lower()
        if src == dst:
            continue
        normalize[dst].add(src)

    # detect script
    if normalize_roman_called:
        script = 'roman'
    else:
        script = '???'

    # find vars
    vars_ = []
    for attr in dir(lang.__class__):
        if attr.startswith('_'):
            continue
        if hasattr(lang.__class__.__bases__[0], attr):
            continue
        vars_.append(attr)
    if lang.vowels:
        vars_.append('vowels')

    # group rewrite/transcribe
    rewrite = []
    transcribe = []
    for x, rule in enumerate(lang.notation.rules):
        pattern = rule[0]
        rpattern = rule[1:]
        # some rpattern is 2d tuple redundantly.
        if isinstance(rpattern[0], tuple):
            rpattern = rpattern[0]
        if isinstance(rpattern[0], hangulize.Phoneme):
            transcribe.append((pattern, rpattern))
        else:
            rewrite.append((pattern, rpattern))

    # find test
    test_module = getattr(__import__('tests.%s' % code), code)
    for attr, val in vars(test_module).items():
        if attr.endswith('TestCase') and not attr.startswith('Hangulize'):
            break
    test_case = val
    examples = test_case.get_examples()

    # render

    sec = Section('lang')
    sec.put('id', code)
    sec.put('codes', lang.iso639_1, lang.iso639_3)
    sec.put('english', locale.get_language_name('en_US'))
    sec.put('korean', locale.get_language_name('ko_KR'))
    sec.put('script', script)
    print(sec.draw('='), end='')

    sec = Section('config')
    sec.put('author', '???')
    sec.put('stage', 'draft')
    if lang.__tmp__:
        sec.put('markers', *lang.__tmp__)
    print(sec.draw('='), end='')

    sec = Section('macros')
    if lang.vowels:
        sec.put('@', '<vowels>')
    print(sec.draw('=', quote_keys=True), end='')

    sec = Section('vars')
    for var in vars_:
        sec.put(var, *getattr(lang, var))
    print(sec.draw('=', quote_keys=True), end='')

    sec = Section('normalize')
    for to, froms in normalize.items():
        sec.put(to, *froms)
    print(sec.draw('=', quote_keys=True), end='')

    sec = Section('rewrite')
    for pattern, rpattern in rewrite:
        sec.put(pattern, rpattern)
    print(sec.draw('->', quote_keys=True), end='')

    sec = Section('transcribe')
    for pattern, rpattern in transcribe:
        sec.put(pattern, rpattern)
    print(sec.draw('->', quote_keys=True), end='')

    sec = Section('test')
    for loanword, hangul in examples.items():
        sec.put(loanword, hangul)
    print(sec.draw('->', quote_keys=True), end='')


if __name__ == '__main__':
    main(sys.argv)
