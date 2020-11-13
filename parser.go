package purplecat

import (
	"fmt"
)

// Parser is the interface for parsing the build file of the software project.
type Parser interface {
	// parses the project located on the given path, and returns the built project instance.
	Parse(path *Path) (Project, error)
	// returns true if the project located on the given path is the target of this parser instance.
	IsTarget(path *Path, context *Context) bool
}

// GradleParser is the instance of Parser for parsing build.gradle (not yet).
type gradleParser struct {
	context *Context
}

// GenerateParser creates and returns the instance of Parser for given path.
func (context *Context) GenerateParser(givenPath string) (Parser, error) {
	parsers := []Parser{
		&mavenParser{context: context},
		&goModParser{context: context},
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
