---
title: ":house: Home"
date: 2020-10-12
draft: false
---

![build](https://github.com/tamadalab/purplecat/workflows/build/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/tamadalab/purplecat/badge.svg?branch=main)](https://coveralls.io/github/tamadalab/purplecat?branch=main)
[![codebeat badge](https://codebeat.co/badges/760a8a6f-2675-4a71-9a77-07c33a807192)](https://codebeat.co/projects/github-com-tamadalab-purplecat-main)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamadalab/purplecat)](https://goreportcard.com/report/github.com/tamadalab/purplecat)

[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/purplecat/blob/main/LICENSE)
[![Version](https://img.shields.io/badge/Version-0.1.0-yellowgreen.svg)](https://github.com/tamada/purplecat/releases/tag/v0.1.0)

# :cat: purplecat

Purple cat, Purple cat, What do you see?
I see the dependent libraries and their licenses!

We develop `purplecat` for detecting the license conflicts.
For this, `purplecat` finds the dependent libraries and their licenses.


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

* [github.com/antchfx/xmlquery](https://github.com/antchfx/xmlquery)
* [github.com/asaskevich/govalidator](https://github.com/asaskevich/govalidator)
* [github.com/go-resty/resty/v2](https://github.com/go-resty/resty)
* [github.com/mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)
* [github.com/spf13/pflag](https://github.com/spf13/pflag)


