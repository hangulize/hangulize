# -*- coding: utf-8 -*-
import copy
# Build the dictionaries to transliterate Serbian cyrillic to latin and vice versa.

# This dictionary is to transliterate from cyrillic to latin.
SR_CYR_TO_LAT_DICT = {
    u'А': u'A', u'а': u'a',
    u'Б': u'B', u'б': u'b',
    u'В': u'V', u'в': u'v',
    u'Г': u'G', u'г': u'g',
    u'Д': u'D', u'д': u'd',
    u'Ђ': u'Đ', u'ђ': u'đ',
    u'Е': u'E', u'е': u'e',
    u'Ж': u'Ž', u'ж': u'ž',
    u'З': u'Z', u'з': u'z',
    u'И': u'I', u'и': u'i',
    u'Ј': u'J', u'ј': u'j',
    u'К': u'K', u'к': u'k',
    u'Л': u'L', u'л': u'l',
    u'Љ': u'Lj', u'љ': u'lj',
    u'М': u'M', u'м': u'm',
    u'Н': u'N', u'н': u'n',
    u'Њ': u'Nj', u'њ': u'nj',
    u'О': u'O', u'о': u'o',
    u'П': u'P', u'п': u'p',
    u'Р': u'R', u'р': u'r',
    u'С': u'S', u'с': u's',
    u'Т': u'T', u'т': u't',
    u'Ћ': u'Ć', u'ћ': u'ć',
    u'У': u'U', u'у': u'u',
    u'Ф': u'F', u'ф': u'f',
    u'Х': u'H', u'х': u'h',
    u'Ц': u'C', u'ц': u'c',
    u'Ч': u'Č', u'ч': u'č',
    u'Џ': u'Dž', u'џ': u'dž',
    u'Ш': u'Š', u'ш': u'š',
}

# This dictionary is to transliterate from Serbian latin to cyrillic.
# Let's build it by simply swapping keys and values of previous dictionary.
SR_LAT_TO_CYR_DICT = {y: x for x, y in iter(SR_CYR_TO_LAT_DICT.items())}

# Build the dictionaries to transliterate Montenegrin cyrillic to latin and vice versa.

# Montenegrin Latin is based on Serbo-Croatian Latin, with the addition of the two letters Ś and Ź,
# to replace the digraphs SJ and ZJ.
# These parallel the two letters of the Montenegrin Cyrillic alphabet not found in Serbian, С́ and З́.
# These, respectively, could also be represented in the original alphabets as šj and žj, and шj and жj.
# Source: https://en.wikipedia.org/wiki/Montenegrin_alphabet#Latin_alphabet
# Also see: http://news.bbc.co.uk/2/hi/8520466.stm
ME_CYR_TO_LAT_DICT = copy.deepcopy(SR_CYR_TO_LAT_DICT)
ME_CYR_TO_LAT_DICT.update({
    u'С́': u'Ś', u'с́': u'ś',  # Montenegrin
    u'З́': u'Ź', u'з́': u'ź'  # Montenegrin
})

# This dictionary is to transliterate from Montenegrin latin to cyrillic.
ME_LAT_TO_CYR_DICT = {y: x for x, y in iter(ME_CYR_TO_LAT_DICT.items())}

# Build the dictionaries to transliterate Macedonian cyrillic to latin and vice versa.
MK_CYR_TO_LAT_DICT = copy.deepcopy(SR_CYR_TO_LAT_DICT)

# Differences with Serbian:
# 1) Between Ze (З з) and I (И и) is the letter Dze (Ѕ ѕ), which looks like the Latin letter S and represents /d͡z/.
MK_CYR_TO_LAT_DICT[u'Ѕ'] = u'Dz'
MK_CYR_TO_LAT_DICT[u'ѕ'] = u'dz'

# 2) Dje (Ђ ђ) is replaced by Gje (Ѓ ѓ), which represents /ɟ/ (voiced palatal stop).
# In some dialects, it represents /d͡ʑ/ instead, like Dje
# It is written ⟨Ǵ ǵ⟩ in the corresponding Macedonian Latin alphabet.
del MK_CYR_TO_LAT_DICT[u'Ђ']
del MK_CYR_TO_LAT_DICT[u'ђ']
MK_CYR_TO_LAT_DICT[u'Ѓ'] = u'Ǵ'
MK_CYR_TO_LAT_DICT[u'ѓ'] = u'ǵ'

# 3) Tshe (Ћ ћ) is replaced by Kje (Ќ ќ), which represents /c/ (voiceless palatal stop).
# In some dialects, it represents /t͡ɕ/ instead, like Tshe.
# It is written ⟨Ḱ ḱ⟩ in the corresponding Macedonian Latin alphabet.
del MK_CYR_TO_LAT_DICT[u'Ћ']
del MK_CYR_TO_LAT_DICT[u'ћ']
MK_CYR_TO_LAT_DICT[u'Ќ'] = u'Ḱ'
MK_CYR_TO_LAT_DICT[u'ќ'] = u'ḱ'

# This dictionary is to transliterate from Macedonian latin to cyrillic.
MK_LAT_TO_CYR_DICT = {y: x for x, y in iter(MK_CYR_TO_LAT_DICT.items())}

# This dictionary is to transliterate from Russian cyrillic to latin (GOST_7.79-2000 System B).
RU_CYR_TO_LAT_DICT = {
    u"А": u"A", u"а": u"a",
    u"Б": u"B", u"б": u"b",
    u"В": u"V", u"в": u"v",
    u"Г": u"G", u"г": u"g",
    u"Д": u"D", u"д": u"d",
    u"Е": u"E", u"е": u"e",
    u"Ё": u"YO", u"ё": u"yo",
    u"Ж": u"ZH", u"ж": u"zh",
    u"З": u"Z", u"з": u"z",
    u"И": u"I", u"и": u"i",
    u"Й": u"J", u"й": u"j",
    u"К": u"K", u"к": u"k",
    u"Л": u"L", u"л": u"l",
    u"М": u"M", u"м": u"m",
    u"Н": u"N", u"н": u"n",
    u"О": u"O", u"о": u"o",
    u"П": u"P", u"п": u"p",
    u"Р": u"R", u"р": u"r",
    u"С": u"S", u"с": u"s",
    u"Т": u"T", u"т": u"t",
    u"У": u"U", u"у": u"u",
    u"Ф": u"F", u"ф": u"f",
    u"Х": u"X", u"х": u"x",
    u"Ц": u"CZ", u"ц": u"cz",
    u"Ч": u"CH", u"ч": u"ch",
    u"Ш": u"SH", u"ш": u"sh",
    u"Щ": u"SHH", u"щ": u"shh",
    u"Ъ": u"''", u"ъ": u"''",
    u"Ы": u"Y'", u"ы": u"y'",
    u"Ь": u"'", u"ь": u"'",
    u"Э": u"E'", u"э": u"e'",
    u"Ю": u"YU", u"ю": u"yu",
    u"Я": u"YA", u"я": u"ya",
}

# This dictionary is to transliterate from Russian latin to cyrillic.
RU_LAT_TO_CYR_DICT = {y: x for x, y in RU_CYR_TO_LAT_DICT.items()}
RU_LAT_TO_CYR_DICT.update({
    u"''": u"ъ",
    u"'": u"ь",
    u"C": u"К", u"c": u"к",
    u"CK": u"К", u"Ck": u"К", u"ck": u"к",
    u"JA": u"ЖА", u"Ja": u"Жа", u"ja": u"жа",
    u"JE": u"ЖЕ", u"Je": u"Же", u"je": u"же",
    u"JI": u"ЖИ", u"Ji": u"Жи", u"ji": u"жи",
    u"JO": u"ЖО", u"Jo": u"Жо", u"jo": u"жо",
    u"JU": u"ЖУ", u"Ju": u"Жу", u"ju": u"жу",
    u"PH": u"Ф", u"Ph": u"Ф", u"ph": u"ф",
    u"TH": u"З", u"Th": u"З", u"th": u"з",
    u"W": u"В", u"w": u"в", u"Q": u"К", u"q": u"к",
    u"WH": u"В", u"Wh": u"В", u"wh": u"в",
    u"Y": u"И", u"y": u"и",
    u"YA": u"Я", u"Ya": u"я", u"ya": u"я",
    u"YE": u"Е", u"Ye": u"е", u"ye": u"е",
    u"YI": u"И", u"Yi": u"и", u"yi": u"и",
    u"YO": u"Ё", u"Yo": u"ё", u"yo": u"ё",
    u"YU": u"Ю", u"Yu": u"ю", u"yu": u"ю",
    u"Y'": u"ы", u"y'": u"ы",
    u"iy": u"ый", u"ij": u"ый",  # dobriy => добрый
})

# Transliterate from Tajik cyrillic to latin
TJ_CYR_TO_LAT_DICT = copy.deepcopy(RU_CYR_TO_LAT_DICT)
# Change Mapping according to ISO 9 (1995)
TJ_CYR_TO_LAT_DICT[u"Э"] = u"È"
TJ_CYR_TO_LAT_DICT[u"э"] = u"è"
TJ_CYR_TO_LAT_DICT[u"ъ"] = u"’"
TJ_CYR_TO_LAT_DICT[u"Х"] = u"H"
TJ_CYR_TO_LAT_DICT[u"х"] = u"h"
TJ_CYR_TO_LAT_DICT[u"Ч"] = u"Č"
TJ_CYR_TO_LAT_DICT[u"ч"] = u"č"
TJ_CYR_TO_LAT_DICT[u"Ж"] = u"Ž"
TJ_CYR_TO_LAT_DICT[u"ж"] = u"ž"
TJ_CYR_TO_LAT_DICT[u"Ё"] = u"Ë"
TJ_CYR_TO_LAT_DICT[u"ё"] = u"ë"
TJ_CYR_TO_LAT_DICT[u"Ш"] = u"Š"
TJ_CYR_TO_LAT_DICT[u"ш"] = u"š"
TJ_CYR_TO_LAT_DICT[u"Ю"] = u"Û"
TJ_CYR_TO_LAT_DICT[u"ю"] = u"û"
TJ_CYR_TO_LAT_DICT[u"Я"] = u"Â"
TJ_CYR_TO_LAT_DICT[u"я"] = u"â"
# delete letters not used 
del TJ_CYR_TO_LAT_DICT[u"Ц"]
del TJ_CYR_TO_LAT_DICT[u"ц"]
del TJ_CYR_TO_LAT_DICT[u"Щ"]
del TJ_CYR_TO_LAT_DICT[u"щ"]
del TJ_CYR_TO_LAT_DICT[u"Ы"]
del TJ_CYR_TO_LAT_DICT[u"ы"]

# update the dict for the additional letters in the tajik cyrillic alphabet ( Ғ, Ӣ, Қ, Ӯ, Ҳ, Ҷ )
# see https://en.wikipedia.org/wiki/Tajik_alphabet#Cyrillic
TJ_CYR_TO_LAT_DICT.update({
    u"Ғ": u"Ǧ", u"ғ": u"ǧ",
    u"Ӣ": u"Ī", u"ӣ": u"ī", 
    u"Қ": u"Q", u"қ": u"q", 
    u"Ӯ": u"Ū", u"ӯ": u"ū",
    u"Ҳ": u"Ḩ", u"ҳ": u"ḩ",
    u"Ҷ": u"Ç", u"ҷ": u"ç"
})

# transliterate from latin tajik to cyrillic
TJ_LAT_TO_CYR_DICT = {y: x for x, y in iter(TJ_CYR_TO_LAT_DICT.items())}

# Transliterate from Bulgarian cyrillic to latin
BG_CYR_TO_LAT_DICT = copy.deepcopy(RU_CYR_TO_LAT_DICT)

# There are a couple of letters that don't exist in Bulgarian:
del BG_CYR_TO_LAT_DICT[u"Ё"]
del BG_CYR_TO_LAT_DICT[u"ё"]
del BG_CYR_TO_LAT_DICT[u"Ы"]
del BG_CYR_TO_LAT_DICT[u"ы"]
del BG_CYR_TO_LAT_DICT[u"Э"]
del BG_CYR_TO_LAT_DICT[u"э"]

# Some letters that are pronounced differently
BG_CYR_TO_LAT_DICT[u"Й"] = u"Y"
BG_CYR_TO_LAT_DICT[u"й"] = u"y"
BG_CYR_TO_LAT_DICT[u"Х"] = u"H"
BG_CYR_TO_LAT_DICT[u"х"] = u"h"
BG_CYR_TO_LAT_DICT[u"Ц"] = u"TS"
BG_CYR_TO_LAT_DICT[u"ц"] = u"ts"
BG_CYR_TO_LAT_DICT[u"Щ"] = u"SHT"
BG_CYR_TO_LAT_DICT[u"щ"] = u"sht"
BG_CYR_TO_LAT_DICT[u"Ю"] = u"YU"
BG_CYR_TO_LAT_DICT[u"ю"] = u"yu"
BG_CYR_TO_LAT_DICT[u"Я"] = u"YA"
BG_CYR_TO_LAT_DICT[u"я"] = u"ya"
# The following letters use the pre-2012 "Andreichin" system for lettering, 
# because in the newest "Ivanov" system "a" and "y" translate to two Bulgarian
# letters and choosing to which one depends on the word and text context
# https://en.wikipedia.org/wiki/Romanization_of_Bulgarian
BG_CYR_TO_LAT_DICT[u"Ъ"] = u"Ă"
BG_CYR_TO_LAT_DICT[u"ъ"] = u"ă"
BG_CYR_TO_LAT_DICT[u"Ь"] = u"J"
BG_CYR_TO_LAT_DICT[u"ь"] = u"j"

# Transliterate from latin Bulgarian to cyrillic.
BG_LAT_TO_CYR_DICT = {y: x for x, y in iter(BG_CYR_TO_LAT_DICT.items())}
BG_LAT_TO_CYR_DICT.update({
    u"ZH": u"Ж", u"Zh": u"Ж", u"zh": u"ж",
    u"TS": u"Ц", u"Ts": u"Ц", u"ts": u"ц",
    u"CH": u"Ч", u"Ch": u"Ч", u"ch": u"ч",
    u"SH": u"Ш", u"Sh": u"Ш", u"sh": u"ш",
    u"SHT": u"Щ", u"Sht": u"Щ", u"sht": u"щ",
    u"YU": u"Ю", u"Yu": u"Ю", u"yu": u"ю",
    u"YA": u"Я", u"Ya": u"Я", u"ya": u"я",
})

# Transliterate from Ukrainian
UA_CYR_TO_LAT_DICT = copy.deepcopy(RU_CYR_TO_LAT_DICT)
# Change mapping to match with Scientific Ukrainian
UA_CYR_TO_LAT_DICT[u"Г"] = u"H"
UA_CYR_TO_LAT_DICT[u"г"] = u"h"
UA_CYR_TO_LAT_DICT[u"Ж"] = u"Ž"
UA_CYR_TO_LAT_DICT[u"ж"] = u"ž"
UA_CYR_TO_LAT_DICT[u"И"] = u"Y"
UA_CYR_TO_LAT_DICT[u"и"] = u"y"
UA_CYR_TO_LAT_DICT[u"Х"] = u"X"
UA_CYR_TO_LAT_DICT[u"х"] = u"x"
UA_CYR_TO_LAT_DICT[u"Ц"] = u"C"
UA_CYR_TO_LAT_DICT[u"ц"] = u"c"
UA_CYR_TO_LAT_DICT[u"Ч"] = u"Č"
UA_CYR_TO_LAT_DICT[u"ч"] = u"č"
UA_CYR_TO_LAT_DICT[u"Ш"] = u"Š"
UA_CYR_TO_LAT_DICT[u"ш"] = u"š"
UA_CYR_TO_LAT_DICT[u"Щ"] = u"Šč"
UA_CYR_TO_LAT_DICT[u"щ"] = u"šč"
UA_CYR_TO_LAT_DICT[u"Ю"] = u"Ju"
UA_CYR_TO_LAT_DICT[u"ю"] = u"ju"
UA_CYR_TO_LAT_DICT[u"Я"] = u"Ja"
UA_CYR_TO_LAT_DICT[u"я"] = u"ja"
# Delete unused letters
del UA_CYR_TO_LAT_DICT[u"Ё"]
del UA_CYR_TO_LAT_DICT[u"ё"]
del UA_CYR_TO_LAT_DICT[u"Ъ"]
del UA_CYR_TO_LAT_DICT[u"ъ"]
del UA_CYR_TO_LAT_DICT[u"Ы"]
del UA_CYR_TO_LAT_DICT[u"ы"]
del UA_CYR_TO_LAT_DICT[u"Э"]
del UA_CYR_TO_LAT_DICT[u"э"]

# Update for Ukrainian letters
UA_CYR_TO_LAT_DICT.update({
    u"Ґ": u"G", u"ґ": u"g",
    u"Є": u"Je", u"є": u"je",
    u"І": u"I", u"і": u"i",
    u"Ї": u"Ï", u"ї": u"ï"  
})

# Latin to Cyrillic
UA_LAT_TO_CYR_DICT = {y: x for x, y in iter(UA_CYR_TO_LAT_DICT.items())}


# This version of Mongolian Latin <-> Cyrillic is based on  MNS 5217:2012
# as far as I know this is the latest standard. Imform me @ https://github.com/Serbipunk
# https://gogo.mn/r/101115
# https://en.wikipedia.org/wiki/Mongolian_Cyrillic_alphabet
MN_CYR_LAT_LIST = [
    u"А", u"A", u"а", u"a",
    u"Э", u"E", u"э", u"e",
    u"И", u"I", u"и", u"i",  # i
    u"О", u"O", u"о", u"o",
    u"У", u"U", u"у", u"u",
    u"Ө", u"Ö", u"ө", u"ö",
    u"Ү", u"Ü", u"ү", u"ü",
    u"Н", u"N", u"н", u"n",
    u"М", u"M", u"м", u"m",
    u"Л", u"L", u"л", u"l",
    u"В", u"V", u"в", u"v",
    u"П", u"P", u"п", u"p",
    u"Ф", u"F", u"ф", u"f",
    u"К", u"K", u"к", u"k",
    u"Х", u"Kh", u"х", u"kh",        # lat 1
    u"Г", u"G", u"г", u"g",
    u"С", u"S", u"с", u"s",
    u"Ш", u"Sh", u"ш", u"sh",  # sh  # lat2
    u"Т", u"T", u"т", u"t",
    u"Д", u"D", u"д", u"d",
    u"Ц", u"Ts", u"ц", u"ts",        # lat3
    u"Ч", u"Ch", u"ч", u"ch",        # lat4
    u"З", u"Z", u"з", u"z",
    u"Ж", u"J", u"ж", u"j",
    u"Й", u"I", u"й", u"i",  # i * 2
    u"Р", u"R", u"р", u"r",
    u"Б", u"B", u"б", u"b",
    u"Е", u"Ye", u"е", u"ye",             # lat 5
    u"Ё", u"Yo", u"ё", u"yo",             # lat 6
    u"Щ", u"Sh", u"щ", u"sh",  # sh x 2   # lat 7
    u"Ъ", u"I", u"ъ", u"i",  # i * 3
    u"Ы", u"Y", u"ы", u"y",
    u"Ь", u"I", u"ь", u"i",  # i * 4
    u"Ю", u"Yu", u"ю", u"yu",             # lat 8
    u"Я", u"Ya", u"я", u"ya",             # lat 9
]
MN_CYR_TO_LAT_DICT = dict([(c, l) for c, l in zip(MN_CYR_LAT_LIST[::2], MN_CYR_LAT_LIST[1::2])])
MN_LAT_TO_CYR_DICT = dict([(l, c) for c, l in zip(MN_CYR_LAT_LIST[-2::-2], MN_CYR_LAT_LIST[-1::-2])])

# Bundle up all the dictionaries in a lookup dictionary
TRANSLIT_DICT = {
    'sr': { # Serbia
        'tolatin': SR_CYR_TO_LAT_DICT,
        'tocyrillic': SR_LAT_TO_CYR_DICT
    },
    'me': { # Montenegro
        'tolatin': ME_CYR_TO_LAT_DICT,
        'tocyrillic': ME_LAT_TO_CYR_DICT
    },
    'mk': { # Macedonia
        'tolatin': MK_CYR_TO_LAT_DICT,
        'tocyrillic': MK_LAT_TO_CYR_DICT
    },
    'ru': { # Russian
        'tolatin': RU_CYR_TO_LAT_DICT,
        'tocyrillic': RU_LAT_TO_CYR_DICT
    },
    'tj': { # Tajik
        'tolatin': TJ_CYR_TO_LAT_DICT,
        'tocyrillic': TJ_LAT_TO_CYR_DICT
    },
    'bg': { # Bulgarian
        'tolatin': BG_CYR_TO_LAT_DICT,
        'tocyrillic': BG_LAT_TO_CYR_DICT
    },
    'ua': { # Ukrainian
        'tolatin': UA_CYR_TO_LAT_DICT,
        'tocyrillic': UA_LAT_TO_CYR_DICT
    },
    'mn': { # Mongolian
        'tolatin': MN_CYR_TO_LAT_DICT,
        'tocyrillic': MN_LAT_TO_CYR_DICT
    }
}
