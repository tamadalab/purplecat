package purplecat

import (
	"fmt"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"

	"github.com/tamadalab/purplecat/logger"
)

type goModParser struct {
	context *Context
}

const (
	localGoModPath = "go/pkg/mod"
	pkgGoDevURL    = "https://pkg.go.dev"
)

func (gmp *goModParser) IsTarget(path *Path, context *Context) bool {
	base := path.Base()
	if base == "go.mod" {
		return path.Exists(context)
	}
	join := path.Join("go.mod")
	return join.Exists(context)
}

func (gmp *goModParser) Parse(path *Path) (*Project, error) {
	return nil, fmt.Errorf("not implement yet")
}

func findLicenseViaPkgGoDev(path *Path, context *Context, currentDepth int) (*Project, error) {
	if !context.Allow(NetworkAccessFlag) {
		return nil, fmt.Errorf("network access denied")
	}
	if context.Depth < currentDepth {
		return nil, fmt.Errorf("over the parsing depth limit %d, current: %d", context.Depth, currentDepth)
	}
	logger.Infof("findLicenseViaPkgGoDev(%s, %d)", path.Path, currentDepth)
	return findLicenseViaPkgGoDevImpl(path, context, currentDepth)
}

func findLicenseViaPkgGoDevImpl(path *Path, context *Context, currentDepth int) (*Project, error) {
	modulePage, err := path.Open(context)
	if err != nil {
		return nil, err
	}
	defer modulePage.Close()

	doc, err := htmlquery.Parse(modulePage)
	if err != nil {
		return nil, err
	}
	return parseProject(doc)
}

func parseProject(root *html.Node) (*Project, error) {
	return nil, nil
}
