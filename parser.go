package purplecat

import (
	"fmt"
)

type Parser interface {
	Parse(path *Path) (*DependencyTree, error)
	IsTarget(path *Path, context *Context) bool
}

type GradleParser struct {
	context *Context
}

func (context *Context) GenerateParser(givenPath string) (Parser, error) {
	parsers := []Parser{
		&MavenParser{context: context},
		// {"build.gradle", &GradleParser{context: context}},
		// {"go.mod", &GoModParser{context: context}},
	}
	path := NewPath(givenPath)
	for _, parser := range parsers {
		if parser.IsTarget(path, context) {
			return parser, nil
		}
	}
	return nil, fmt.Errorf("%s: cannot parse the project", path)
}
