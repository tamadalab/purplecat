package purplecat

import "fmt"

type GoModParser struct {
	context *Context
}

const PKG_GO_DEV_URL = "https://pkg.go.dev"

func (gmp *GoModParser) Parse(path *Path) (*DependencyTree, error) {
	return nil, fmt.Errorf("not implement yet")
}

func FindLicenseFromPkgGoDev(path *Path, context *Context) (*DependencyTree, error) {
	if !context.Allow(NETWORK_ACCESS) {
		return nil, fmt.Errorf("network access denied")
	}
	return nil, nil
}
