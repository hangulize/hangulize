# 1to2

Migrates a language spec to Hangulize 2
from the original Hangulize.

## Installation

```console
$ git clone https://github.com/sublee/hangulize.git
$ pip install -e hangulize

$ cd /path/to/hangulize2
$ pip install -r 1to2/requirements.txt
```

## Migration

```console
$ python 1to2/1to2.py LANG > hangulize/bundle/LANG.hgl
```
