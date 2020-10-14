package purplecat

import (
	"fmt"
	"io"
	"strings"
)

type ActType int

const NETWORK_ACCESS ActType = iota + 1

type DependencyTree struct {
	ProjectName  string
	LicenseNames []string
	Dependencies []*DependencyTree
}

type Context struct {
	AllowNetworkAccess bool
	Format             string
	Depth              int
}

func NewContext(allowNetworkAccess bool, format string, depth int) *Context {
	return &Context{AllowNetworkAccess: allowNetworkAccess, Format: format, Depth: depth}
}

func (context *Context) Allow(actType ActType) bool {
	if actType == NETWORK_ACCESS {
		return context.AllowNetworkAccess
	}
	return false
}

func (context *Context) NewWriter(out io.Writer) (Writer, error) {
	switch strings.ToLower(context.Format) {
	case "csv":
		return &CsvWriter{Out: out}, nil
	case "json":
		return &JsonWriter{Out: out}, nil
	case "toml":
		return &TomlWriter{Out: out}, nil
	case "yaml", "yml":
		return &YamlWriter{Out: out}, nil
	case "xml":
		return &XmlWriter{Out: out}, nil
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
