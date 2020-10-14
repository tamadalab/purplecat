[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/purplecat/blob/main/LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-yellowgreen.svg)](https://github.com/tamada/purplecat/releases/tag/v1.0.0)

# :cat: purplecat

Purple cat, Purple cat, What do you see?
I see the dependent libraries and their licenses!

We develop `purplecat` for detecting the license conflicts.
For this, `purplecat` finds the dependent libraries and their licenses.

## :runner: Usage 

```sh
$ purplecat -h
purplecat version 1.0.0
purplecat [OPTIONS] <PROJECTs...>
OPTIONS
    -d, --depth <DEPTH>      specifies the depth for parsing (default: 1)
    -f, --format <FORMAT>    specifies the format of the result. Default is 'markdown'.
                             Available values are: CSV, JSON, YAML, XML, and Markdown.
    -o, --output <FILE>      specifies the destination file (default: STDOUT).
    -N, --offline            offline mode (no network access).

    -h, --help               prints this message.
PROJECT
    target project for extracting related libraries and their licenses.
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

```sh
$ docker run -v /target/project/dir:/home/projects tamadalab/purplecat
```

## :anchor: Install

### :beer: Homebrew (not yet)

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

- :+1: Commercial use, 
- :+1: Modification, 
- :+1: Distribution, and 
- :+1: Private use.

### :jack_o_lantern: Logo

![purplecat](https://github.com/tamadalab/purplecat/raw/main/site/static/images/purplecat_128.png)

This image comes from https://pixy.org/4693873/ ([CC-0](https://creativecommons.org/publicdomain/zero/1.0)).

### :name_badge: Project name come from?

The project name come from the children book, titled "[Brown Bear, Brown Bear, What do you see?](https://www.amazon.com/dp/B07BZS8RS9)".

I read the book to my child many times by her request.
The purple cat is the most favorit character of her.

### :woman_office_worker: Developers :man_office_worker:

* [Haruaki TAMADA](https://github.com/tamada)

