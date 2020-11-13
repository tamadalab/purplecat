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

type localProject struct {
	info         projectNamer
	licenses     []*License
	dependencies []Project
}

type projectNamer interface {
	Name() string
}

func (project *localProject) Name() string {
	return project.info.Name()
}

func (project *localProject) Licenses() []*License {
	return project.licenses
}

func (project *localProject) Dependencies() []Project {
	return project.dependencies
}

func (project *localProject) PutDependency(p Project) {
	project.dependencies = append(project.dependencies, p)
}

// Project shows the project information especially its name.
type Project interface {
	Name() string
	Licenses() []*License
	Dependencies() []Project
	PutDependency(project Project)
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
	Cache             *CacheContext
	CacheDB           CacheDB
}

// NewContext creates the instance of Context by given arguments.
func NewContext(denyNetworkAccess bool, format string, depth int) *Context {
	context, _ := NewContextWithCache(denyNetworkAccess, format, depth, NoCache)
	return context
}

// NewContextWithCache creates the instance of Context with loading cache database by given arguments.
func NewContextWithCache(denyNetworkAccess bool, format string, depth int, cType CacheType) (*Context, error) {
	cacheContext, err := NewCacheContext(cType)
	if err != nil {
		return nil, err
	}
	cacheDB, err := NewCacheDB(cacheContext)
	if err != nil {
		return nil, err
	}
	return &Context{DenyNetworkAccess: denyNetworkAccess, Format: format, Depth: depth, Cache: cacheContext, CacheDB: cacheDB}, nil
}

// SearchCache searches from cache database.
func (context *Context) SearchCache(name string) ([]*License, bool) {
	if context.CacheDB == nil {
		db, err := NewCacheDB(context.Cache)
		if err != nil {
			return nil, false
		}
		context.CacheDB = db
	}
	return context.CacheDB.Find(name)
}

// RegisterCache stores the given licenses to cache database.
func (context *Context) RegisterCache(name string, licenses []*License) bool {
	if context.CacheDB == nil {
		db, err := NewCacheDB(context.Cache)
		if err != nil {
			return false
		}
		context.CacheDB = db
	}
	return context.CacheDB.Register(name, licenses)
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
