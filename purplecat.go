package purplecat

import (
	"fmt"
	"io"
	"strings"
)

// ActionType shows the acton types for Context.
type ActionType int

// NetworkAccessFlag is the one of ActType, which represents fo rNetworkAccess checking.
const NetworkAccessFlag ActionType = iota + 1

// Licenses is the slice for pointers of License.
type Licenses []*License

// Projects is the slice for pointers of Project.
type Projects []*Project

// Project shows the project information especially its name, licenses, and dependencies.
type Project struct {
	PName       string     `json:"name"`
	LicenseList []*License `json:"licenses"`
	Deps        []string   `json:"dependencies"`
	context     CacheDB    `json:"-"`
}

// NewProject creates an instance of Project.
func (context *Context) NewProject(name string, licenses Licenses) *Project {
	project := &Project{PName: name, LicenseList: licenses, context: context.Cache, Deps: []string{}}
	context.Cache.Register(project)
	return project
}

// Name returns the name of the receiver project.
func (project *Project) Name() string {
	return project.PName
}

// Licenses returns the license list of the receiver project.
func (project *Project) Licenses() Licenses {
	return project.LicenseList
}

// Dependencies returns the dependency list of the receiver project.
func (project *Project) Dependencies() Projects {
	projects := []*Project{}
	for _, dep := range project.Deps {
		depProject, ok := project.context.Find(dep)
		if ok {
			projects = append(projects, depProject)
		}
	}
	return projects
}

// AddDependency adds the given project as the dependency for the receiver project.
func (project *Project) AddDependency(p *Project) {
	if p == nil {
		return
	}
	if _, ok := project.context.Find(p.Name()); !ok {
		project.context.Register(p)
	}
	project.Deps = append(project.Deps, p.Name())
}

// UnknownLicense is the instance of license, means unknown.
var UnknownLicense = &License{Name: "unknown", SpdxID: "unknown", URL: ""}

// License represents the license name, and its url.
type License struct {
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
}

// Context means application context of purplecat.
type Context struct {
	DenyNetworkAccess bool
	Format            string
	Depth             int
	Cache             CacheDB
}

// NewContext creates the instance of Context by given arguments.
func NewContext(denyNetworkAccess bool, format string, depth int) *Context {
	context, _ := NewContextWithCache(denyNetworkAccess, format, depth, MemoryCache)
	return context
}

// NewContextWithCache creates the instance of Context with loading cache database by given arguments.
func NewContextWithCache(denyNetworkAccess bool, format string, depth int, cType CacheType) (*Context, error) {
	cache, err := NewCacheDB(cType)
	if err != nil {
		return nil, err
	}
	return &Context{DenyNetworkAccess: denyNetworkAccess, Format: format, Depth: depth, Cache: cache}, nil
}

// SearchCache searches from cache database.
func (context *Context) SearchCache(name string) (*Project, bool) {
	if context.Cache == nil {
		return nil, false
	}
	return context.Cache.Find(name)
}

// RegisterCache stores the given licenses to cache database.
func (context *Context) RegisterCache(project *Project) bool {
	if context.Cache == nil {
		return false
	}
	return context.Cache.Register(project)
}

// Allow checks given ActType is allowed in the current context.
func (context *Context) Allow(actType ActionType) bool {
	if actType == NetworkAccessFlag {
		return !context.DenyNetworkAccess
	}
	return false
}

// NewWriter creates an suitable Writer instance.
func (context *Context) NewWriter(out io.Writer) (Writer, error) {
	switch strings.ToLower(context.Format) {
	case "csv":
		return &csvWriter{Out: out}, nil
	case "json":
		return &jsonWriter{Out: out}, nil
	case "toml":
		return &tomlWriter{Out: out}, nil
	case "yaml", "yml":
		return &yamlWriter{Out: out}, nil
	case "xml":
		return &xmlWriter{Out: out}, nil
	case "markdown", "md":
		return &markdownWriter{Out: out}, nil
	default:
		return nil, fmt.Errorf("%s: unknown format", context.Format)
	}
}
