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
[![Version](https://img.shields.io/badge/Version-0.3.3-yellowgreen.svg)](https://github.com/tamada/purplecat/releases/tag/v0.3.3)
[![Heroku-Deployed](https://img.shields.io/badge/Heroku-Deployed-green?logo=Heroku)](https://afternoon-wave-39227.herokuapp.com/purplecat/)
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Ftamadalab%2Fpurplecat%3A0.3.3-blue?logo=docker)](https://github.com/orgs/tamadalab/packages/container/package/purplecat)

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

* [github.com/antchfx/htmlquery](https://github.com/antchfx/htmlquery)
* [github.com/antchfx/xmlquery](https://github.com/antchfx/xmlquery)
* [github.com/asaskevich/govalidator](https://github.com/asaskevich/govalidator)
* [github.com/go-resty/resty/v2](https://github.com/go-resty/resty)
* [github.com/gorilla/mux](https://github.com/gorilla/mux)
* [github.com/mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)
* [github.com/spf13/pflag](https://github.com/spf13/pflag)
* [golang.org/pkg/net/http](https://golang.org/pkg/net/http/)
