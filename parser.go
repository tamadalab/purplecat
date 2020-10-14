package purplecat

import (
	"fmt"
	"path/filepath"
)

type Parser interface {
	Parse(path string) (*DependencyTree, error)
}

type GradleParser struct {
	context *Context
}

type GoModParser struct {
	context *Context
}

func (context *Context) GenerateParser(path string) (Parser, error) {
	generators := []struct {
		fileName string
		parser   Parser
	}{
		{"pom.xml", &MavenParser{context: context}},
		// {"build.gradle", &GradleParser{context: context}},
		// {"go.mod", &GoModParser{context: context}},
	}
	for _, generator := range generators {
		buildFile := filepath.Join(path, generator.fileName)
		if FindFile(buildFile) {
			return generator.parser, nil
		}
	}
	return nil, fmt.Errorf("%s: cannot parse the project", path)
}

