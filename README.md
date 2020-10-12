[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/purplecat/blob/main/LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-yellowgreen.svg)](https://github.com/tamada/purplecat/releases/tag/v1.0.0)

# :cat: purplecat

Purple cat, Purple cat, What do you see?
I see the dependent libraries and their licenses.

## :runner: Usage 

```sh
$ purplecat -h
purplecat version 1.0.0
purplecat [OPTIONS] <PROJECTs...>
OPTIONS
    -d, --dest <FILE>    specifies the destination file (default is STDOUT).
    -N, --offline        offline mode (no network access).

PROJECT
    target project for extracting related libraries and their licenses.
```

## :whale: Docker

```
$ docker run -v /target/project/dir:/home/projects tamadalab/purplecat
```

## :anchor: Install

### :beer: Homebrew

```sh
$ brew tap tamadalab/brew
$ brew install purplecat
```

### :muscle: Compiling yourself

```sh
$ make
```

### Requirements

* [github.com/spf13/pflag](https://github.com/spf13/pflag)
* [github.com/tamada/lioss](https://github.com/tamada/lioss)

## :smile: About

![purplecat](https://github.com/tamadalab/purplecat/raw/main/site/static/images/purplecat_128.png)

This image comes from https://pixy.org/4693873/ ([CC-0](https://creativecommons.org/publicdomain/zero/1.0)).
