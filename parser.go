package purplecat

import (
	"fmt"
)

/*Parser is the interface for parsing the build file of the software project.*/
type Parser interface {
	Parse(path *Path) (*Project, error)
	IsTarget(path *Path, context *Context) bool
}

/*GradleParser is the instance of Parser for parsing build.gradle (not yet).*/
type GradleParser struct {
	context *Context
}

/*GenerateParser creates and returns the instance of Parser for given path. */
func (context *Context) GenerateParser(givenPath string) (Parser, error) {
	parsers := []Parser{
		&MavenParser{context: context},
		&GoModParser{context: context},
		// {"build.gradle", &GradleParser{context: context}},
	}
	path := NewPath(givenPath)
	for _, parser := range parsers {
		if parser.IsTarget(path, context) {
			return parser, nil
		}
	}
	return nil, fmt.Errorf("%s: cannot parse the project", givenPath)
}
