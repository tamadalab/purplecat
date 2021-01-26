![build](https://github.com/tamadalab/purplecat/workflows/build/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/tamadalab/purplecat/badge.svg?branch=main)](https://coveralls.io/github/tamadalab/purplecat?branch=main)
[![codebeat badge](https://codebeat.co/badges/760a8a6f-2675-4a71-9a77-07c33a807192)](https://codebeat.co/projects/github-com-tamadalab-purplecat-main)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamadalab/purplecat)](https://goreportcard.com/report/github.com/tamadalab/purplecat)

[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/purplecat/blob/main/LICENSE)
[![Version](https://img.shields.io/badge/Version-0.3.3-yellowgreen.svg)](https://github.com/tamada/purplecat/releases/tag/v0.3.3)
[![Heroku-Deployed](https://img.shields.io/badge/Heroku-Deployed-green?logo=Heroku)](https://afternoon-wave-39227.herokuapp.com/purplecat/)
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Ftamadalab%2Fpurplecat%3A0.3.2-blue?logo=docker)](https://github.com/orgs/tamadalab/packages/container/package/purplecat)

# :cat: purplecat

Purple cat, Purple cat, What do you see?
I see the dependent libraries and their licenses!

We develop `purplecat` for detecting the license conflicts.
For this, `purplecat` finds the dependent libraries and their licenses.

## :runner: Usage

```sh
$ purplecat -h
purplecat version 0.3.2
purplecat [COMMON_OPTIONS] [CLI_MODE_OPTIONS] [SERVER_MODE_OPTIONS] <PROJECTs...|BUILD_FILEs...>
COMMON_OPTIONS
    -c, --cache-type <TYPE>        specifies the cache type. (default: default).
                                   Available values are: default, ref-only, newdb and memory.
        --cachedb-path <DBPATH>    specifies the cache database path
                                   (default: ~/.config/purplecat/cachedb.json).
    -l, --log-level <LOGLEVEL>     specifies the log level. (default: WARN).
                                   Available values are: DEBUG, INFO, WARN, and FATAL
    -h, --help                     prints this message.

CLI_MODE_OPTIONS
    -d, --depth <DEPTH>            specifies the depth for parsing (default: 1)
    -f, --format <FORMAT>          specifies the result format. Default is 'markdown'.
                                   Available values are: CSV, JSON, YAML, XML, and Markdown.
    -o, --output <FILE>            specifies the destination file (default: STDOUT).
    -N, --offline                  offline mode (no network access).

SERVER_MODE_OPTIONS
    -p, --port <PORT>              specifies the port number of REST API server. Default is 8080.
                                   If '--server' option did not specified, purplecat ignores this option.
    -s, --server                   starts REST API server. With this option, purplecat ignores
                                   CLI_MODE_OPTIONS and arguments.

PROJECT
    target project for extracting dependent libraries and their licenses.
BUILD_FILE
    build file of the project for extracting dependent libraries and their licenses

purplecat support the projects using the following build tools.
    * Maven 3 (pom.xml)
```

### Resultant Format in CLI Mode

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

[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Ftamadalab%2Fpurplecat%3A0.3.2-blue?logo=docker)](https://github.com/orgs/tamadalab/packages/container/package/purplecat)

* `tamadalab/purplecat`
    * `0.3.2`, `latest`
    * `0.3.1`
    * `0.3.0`
    * `0.2.0`
    * `0.1.0`

```sh
$ docker run -v /target/project/dir:/home/purplecat ghcr.io/tamadalab/purplecat pom.xml # <- CLI Mode
$ docker run -p 8080:8080 -v /target/project/dir:/home/purplecat ghcr.io/tamadalab/purplecat --server --port 8080 # <- Server Mode
```

## :bathtub: Rest API

Purplecat provides REST API server as specifying option `'-s'` or `'--server'` to `purplecat` command.

[![Heroku-Deployed](https://img.shields.io/badge/Heroku-Deployed-green?logo=Heroku)](https://afternoon-wave-39227.herokuapp.com/purplecat/)

### End points

#### `/purplecat/api/licenses`

* `GET`
    * run purplecat with pom url by query param and returns the result as JSON format.
    * Query params
        * `target` (required)
            * specifies the target build file url.
        * `depth`
            * specifies the depth of the parsing. Default is 1.
    * Status Codes
        * 200 OK
            * provides license data of the build files as json format.
        * 404 Not found
            * specified build file not found.
        * 500 Error
            * parsing error.
* `POST`
    * run purplecat with pom data from request body and returns the result as JSON format.
    * Query params
        * `depth`
            * specifies the depth of the parsing. Default is 1.
    * Requst body
        * plain `pom.xml` data.
    * Status Codes
        * 200 OK
            * provides license data of the build files as json format.
        * 404 Not found
            * specified build file not found.
        * 500 Error
            * parsing error.

#### `/purplecat/api/caches`

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

* [github.com/antchfx/htmlquery](https://github.com/antchfx/htmlquery)
* [github.com/antchfx/xmlquery](https://github.com/antchfx/xmlquery)
* [github.com/asaskevich/govalidator](https://github.com/asaskevich/govalidator)
* [github.com/go-resty/resty/v2](https://github.com/go-resty/resty)
* [github.com/gorilla/mux](https://github.com/gorilla/mux)
* [github.com/mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)
* [github.com/spf13/pflag](https://github.com/spf13/pflag)
* [golang.org/pkg/net/http](https://golang.org/pkg/net/http/)

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
