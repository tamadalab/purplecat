package purplecat

import (
	"fmt"
	"io"
	"strings"
)

type ActType int

const NetworkAccessFlag ActType = iota + 1

type DependencyTree struct {
	ProjectInfo  ProjectInfo
	Licenses     []*License
	Dependencies []*DependencyTree
}

type ProjectInfo interface {
	Name() string
}

type Context struct {
	DenyNetworkAccess bool
	Format            string
	Depth             int
}

var UnknownLicense = &License{Name: "unknown", SpdxID: "unknown", URL: ""}

type License struct {
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
}

func NewContext(denyNetworkAccess bool, format string, depth int) *Context {
	return &Context{DenyNetworkAccess: denyNetworkAccess, Format: format, Depth: depth}
}

func (context *Context) Allow(actType ActType) bool {
	if actType == NetworkAccessFlag {
		return !context.DenyNetworkAccess
	}
	return false
}

func (context *Context) NewWriter(out io.Writer) (Writer, error) {
	switch strings.ToLower(context.Format) {
	case "csv":
		return &CsvWriter{Out: out}, nil
	case "json":
		return &JSONWriter{Out: out}, nil
	case "toml":
		return &TomlWriter{Out: out}, nil
	case "yaml", "yml":
		return &YamlWriter{Out: out}, nil
	case "xml":
		return &XMLWriter{Out: out}, nil
	case "markdown", "md":
		return &MarkdownWriter{Out: out}, nil
	default:
		return nil, fmt.Errorf("%s: unknown format", context.Format)
	}
}

/* ParseProject find the project and its license from given path.
 */
func ParseProject(projectPath string) *DependencyTree {
	return nil
}
