---
title: ":runner: Usage"
---

```sh
$ purplecat -h
purplecat version 0.1.0
purplecat [OPTIONS] <PROJECTs...|BUILD_FILEs...>
OPTIONS
    -d, --depth <DEPTH>       specifies the depth for parsing (default: 1)
    -f, --format <FORMAT>     specifies the format of the result. Default is 'markdown'.
                              Available values are: CSV, JSON, YAML, XML, and Markdown.
    -l, --level <LOGLEVEL>    specifies the log level. (default: WARN).
                              Available values are: DEBUG, INFO, WARN, SEVERE
    -o, --output <FILE>       specifies the destination file (default: STDOUT).
    -N, --offline             offline mode (no network access).

    -h, --help                prints this message.
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

## :whale: Docker (not implement yet)

```sh
$ docker run -v /target/project/dir:/home/projects tamadalab/purplecat
```
