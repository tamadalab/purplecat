package purplecat

import (
	"fmt"
	"io"
	"strings"
)

/*ActType shows the acton types for Context.
 */
type ActType int

/*NetworkAccessFlag is the one of ActType, which represents fo rNetworkAccess checking.
 */
const NetworkAccessFlag ActType = iota + 1

/*Project represents the software project, contains its name, licenses, and dependencies.
 */
type Project struct {
	Info         ProjectInfo
	Licenses     []*License
	Dependencies []*Project
}

/*ProjectInfo shows the project information especially its name.
 */
type ProjectInfo interface {
	Name() string
}

/*Context means application context of purplecat. */
type Context struct {
	DenyNetworkAccess bool
	Format            string
	Depth             int
}

/*UnknownLicense is the instance of license, means unknown. */
var UnknownLicense = &License{Name: "unknown", SpdxID: "unknown", URL: ""}

/*License represents the license name, and its url. */
type License struct {
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
}

/*NewContext creates the instance of Context by given arguments. */
func NewContext(denyNetworkAccess bool, format string, depth int) *Context {
	return &Context{DenyNetworkAccess: denyNetworkAccess, Format: format, Depth: depth}
}

/*Allow checks given ActType is allowed in the current context. */
func (context *Context) Allow(actType ActType) bool {
	if actType == NetworkAccessFlag {
		return !context.DenyNetworkAccess
	}
	return false
}

/*NewWriter creates an suitable Writer instance. */
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

/*ParseProject find the project and its license from given path. */
func ParseProject(projectPath string) *Project {
	return nil
}
