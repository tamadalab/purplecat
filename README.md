![build](https://github.com/tamadalab/purplecat/workflows/build/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/tamadalab/purplecat/badge.svg?branch=main)](https://coveralls.io/github/tamadalab/purplecat?branch=main)
[![codebeat badge](https://codebeat.co/badges/760a8a6f-2675-4a71-9a77-07c33a807192)](https://codebeat.co/projects/github-com-tamadalab-purplecat-main)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamadalab/purplecat)](https://goreportcard.com/report/github.com/tamadalab/purplecat)

[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/purplecat/blob/main/LICENSE)
[![Version](https://img.shields.io/badge/Version-0.2.0-yellowgreen.svg)](https://github.com/tamada/purplecat/releases/tag/v0.2.0)
[![Docker](https://img.shields.io/badge/docker-tamadalab%2Fpurplecat%3A0.2.0-blue?logo=docker&style=social)](https://hub.docker.com/r/tamadalab/purplecat)

# :cat: purplecat

Purple cat, Purple cat, What do you see?
I see the dependent libraries and their licenses!

We develop `purplecat` for detecting the license conflicts.
For this, `purplecat` finds the dependent libraries and their licenses.

## :runner: Usage

```sh
$ purplecat -h
purplecat version 0.2.0
purplecat [OPTIONS] <PROJECTs...|BUILD_FILEs...>
OPTIONS
    -c, --cache-type <TYPE>        specifies the cache type. (default: default).
                                   Available values are: default, ref-only, newdb and memory.
        --cachedb-path <DBPATH>    specifies the cache database path
                                   (default: ~/.config/purplecat/cachedb.json).
    -d, --depth <DEPTH>            specifies the depth for parsing (default: 1)
    -f, --format <FORMAT>          specifies the result format. Default is 'markdown'.
                                   Available values are: CSV, JSON, YAML, XML, and Markdown.
    -l, --log-level <LOGLEVEL>     specifies the log level. (default: WARN).
                                   Available values are: DEBUG, INFO, WARN, and FATAL
    -o, --output <FILE>            specifies the destination file (default: STDOUT).
    -N, --offline                  offline mode (no network access).

    -h, --help                     prints this message.
PROJECT
    target project for extracting dependent libraries and their licenses.
BUILD_FILE
    build file of the project for extracting dependent libraries and their licenses

purplecat support the projects using the following build tools.
    * Maven 3 (pom.xml)
```

### Resultant Format

#### CSV

```csv
project-name,license-name,parent
project1,["Apache 2.0"],
dependent-project1,["Apache 2.0"],project1
dependent-project2,["BSD"],project1
```

#### Json

```json
{
    "project-name": "project1",
    "license-name": "Apache 2.0",
    "dependencies": [
        {
            "project-name": "dependent-project1",
            "license-name": "Apache 2.0"
        }, {
            "project-name": "dependent-project2",
            "license-name": "BSD"
        }
    ]
}
```

you can get formatted Json using `jq`, like the following the command.

```sh
$ purplecat .... | jq .
```

#### Yaml

```yaml
project-name: project1
license-name: Apache 2.0
dependencies:
- project-name: dependent-project1
  license-name: Apache 2.0
- project-name: dependent-project2
  license-name: BSD
```

#### Xml

```xml
<?xml version="1.0"?>
<purplecat>
  <project-name>project1</project-name>
  <license-names>
    <license-name>Apache 2.0</license-name>
  </license-names>
  <dependencies>
    <dependency>
      <project-name>dependent-project1</project-name>
      <license-names>
        <license-name>Apache 2.0</license-name>
      </license-names>
    </dependency>
    <dependency>
      <project-name>dependent-project2</project-name>
      <license-names>
        <license-name>BSD</license-name>
      </license-names>
    </dependency>
  </dependencies>
</purplecat>
```

#### Toml (not implemented yet)

```toml
project-name = "project1"
license-name = "Apache 2.0"
[[dependencies]]
project-name: dependent-project1
license-name: Apache 2.0
[[dependencies]]
project-name: dependent-project2
license-name: BSD
```

#### Markdown

```markdown
* project1: Apache 2.0
    * dependent-project1: ["Apache 2.0"]
    * dependent-project2: ["BSD"]
```

## :whale: Docker

[![Docker](https://img.shields.io/badge/docker-tamadalab%2Fpurplecat%3A0.1.0-blue?logo=docker&style=social)](https://hub.docker.com/r/tamadalab/purplecat)

* `tamadalab/purplecat`
    * `0.2.0`, `latest`
    * `0.1.0`

```sh
$ docker run -v /target/project/dir:/home/purplecat tamadalab/purplecat pom.xml
```

## :bathtub: Rest API

Purplecat provides REST API server as `pcrserver`.

```sh
pcrserver [OPTIONS]
OPTIONS
    -c, --cache-type <TYPE>        specifies the cache type. (default: default).
                                   Available values are: default, ref-only, newdb and memory.
        --cachedb-path <DBPATH>    specifies the cache database path
                                   (default: ~/.config/purplecat/cachedb.json).
    -l, --log-level <LOGLEVEL>     specifies the log level. (default: WARN).
                                   Available values are: DEBUG, INFO, WARN, and FATAL
    -p, --port <PORT>              specifies the port number, default is 8080.

    -h, --help                     print this message.
```

### End points

#### `/purplecat/licenses`

* `GET`
    * run purplecat by giving build file.
    * Query params
        * `target` (required)
            * specifies the target build file url.
        * `depth`
            * specifies the depth of the parsing.
    * Status Codes
        * 200 OK
            * provides license data of the build files as json format.
        * 404 Not found
            * specified build file not found.
        * 500 Error
            * parsing error.

#### `/purplecat/caches`

* `GET`
    * getting the whole cached data.
    * Status Codes
        * 200 OK
           * provides cache data as json format.
* `DELETE`
    * delete cache data.
        * 200 OK
            * always returns this code.


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

## :smile: About

### :scroll: License

[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/purplecat/blob/main/LICENSE)

This license permits

- :+1: Commercial use,
- :+1: Modification,
- :+1: Distribution, and
- :+1: Private use.

### :jack_o_lantern: Logo

![purplecat](https://github.com/tamadalab/purplecat/raw/main/site/static/images/purplecat_128.png)

This image comes from https://pixy.org/4693873/ ([CC-0](https://creativecommons.org/publicdomain/zero/1.0)).

### :name_badge: Project name come from?

The project name come from the children book, titled "[Brown Bear, Brown Bear, What do you see?](https://www.amazon.com/dp/B07BZS8RS9)" by Bill Martin, Jr. and Eric Carle.

I read the book to my child many times by her request.
The purple cat is the most favorit character of her.

### :woman_office_worker: Developers :man_office_worker:

* [Haruaki TAMADA](https://github.com/tamada)
