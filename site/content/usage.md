---
title: ":runner: Usage"
---

```sh
$ purplecat -h
purplecat version 0.3.0
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

### Resultant Format in CLI mode

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

[![Docker](https://img.shields.io/badge/docker-tamadalab%2Fpurplecat%3A0.2.0-blue?logo=docker)](https://github.com/orgs/tamadalab/packages/container/package/purplecat)

* `tamadalab/purplecat`
    * `0.3.0`, `latest`
    * `0.2.0`
    * `0.1.0`

```sh
$ docker run -v /target/project/dir:/home/purplecat ghcr.io/tamadalab/purplecat pom.xml # <- CLI Mode
$ docker run -p 8080:8080 -v /target/project/dir:/home/purplecat ghcr.io/tamadalab/purplecat --server --port 8080 # <- Server Mode
```

## :bathtub: Rest API

Purplecat provides REST API server as specifying option `'-s'` or `'--server'` to `purplecat` command.

### End points

#### `/purplecat/api/licenses`

* `GET`
    * run purplecat by giving build file and returns the result as JSON format.
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
