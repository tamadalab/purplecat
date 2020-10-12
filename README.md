[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/purplecat/blob/main/LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-yellowgreen.svg)](https://github.com/tamada/purplecat/releases/tag/v1.0.0)

# :cat: purplecat

Purple cat, Purple cat, What do you see?
I see the dependent libraries and their licenses!

We develop `purplecat` for detecting the license conflicts.

## :runner: Usage 

```sh
$ purplecat -h
purplecat version 1.0.0
purplecat [OPTIONS] <PROJECTs...>
OPTIONS
    -d, --dest <FILE>    specifies the destination file (default is STDOUT).
    -N, --offline        offline mode (no network access).

    -h, --help           prints this message.
PROJECT
    target project for extracting related libraries and their licenses.
```


## :whale: Docker

```sh
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

### :scroll: License

[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/purplecat/blob/main/LICENSE)

This license permits

- Commercial use, 
- Modification, 
- Distribution, and 
- Private use.

### :jack_o_lantern: Logo

![purplecat](https://github.com/tamadalab/purplecat/raw/main/site/static/images/purplecat_128.png)

This image comes from https://pixy.org/4693873/ ([CC-0](https://creativecommons.org/publicdomain/zero/1.0)).

### :name_badge: Project name come from?

The project name come from the children book, titled "[Brown Bear, Brown Bear, What do you see?](https://www.amazon.com/dp/B07BZS8RS9)".

I read the book to my child many times by her request.

### :woman_office_worker: Developers :man_office_worker:

* [Haruaki TAMADA](https://github.com/tamada)

