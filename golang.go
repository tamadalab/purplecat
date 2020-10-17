package purplecat

import "fmt"

type GoModParser struct {
	context *Context
}

const pkgGoDevURL = "https://pkg.go.dev"

func (gmp *GoModParser) IsTarget(path *Path, context *Context) bool {
	base := path.Base()
	if base == "go.mod" {
		return path.Exists(context)
	}
	join := path.Join("go.mod")
	return join.Exists(context)
}

func (gmp *GoModParser) Parse(path *Path) (*DependencyTree, error) {
	return nil, fmt.Errorf("not implement yet")
}

func FindLicenseFromPkgGoDev(path *Path, context *Context) (*DependencyTree, error) {
	if !context.Allow(NetworkAccessFlag) {
		return nil, fmt.Errorf("network access denied")
	}
	return nil, nil
}
