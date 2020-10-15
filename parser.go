package purplecat

import (
	"fmt"
)

type Parser interface {
	Parse(path *Path) (*DependencyTree, error)
}

type GradleParser struct {
	context *Context
}

func (context *Context) GenerateParser(givenPath string) (Parser, error) {
	generators := []struct {
		fileName string
		parser   Parser
	}{
		{"pom.xml", &MavenParser{context: context}},
		// {"build.gradle", &GradleParser{context: context}},
		// {"go.mod", &GoModParser{context: context}},
	}
	path := NewPath(givenPath)
	for _, generator := range generators {
		if path.Base() == generator.fileName && path.Exists(context) {
			return generator.parser, nil
		}
		buildFilePath := path.Join(generator.fileName)
		if buildFilePath.Exists(context) {
			return generator.parser, nil
		}
	}
	return nil, fmt.Errorf("%s: cannot parse the project", path)
}
