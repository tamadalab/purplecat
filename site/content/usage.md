---
title: ":runner: Usage"
---

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
$ docker run -v /target/project/dir:/home/purplecat tamadalab/purplecat:latest pom.xml
```


