package purplecat

import "testing"

func TestParse(t *testing.T) {
	projectName := "testdata/mavenproject"
	path := NewPath(projectName)
	parser := &MavenParser{NewContext(false, "json", 1)}

	tree, err := parser.Parse(path)
	if err != nil {
		t.Errorf("%s: license parse failed: %s", projectName, err.Error())
		return
	}
	validateDependencyTree(t, tree, "jp.ac.kyoto_su/project4test/1.0.0", "The Apache Software License, Version 2.0", 3)
	validateDependencyTree(t, tree.Dependencies[0], "args4j/args4j/2.33", "MIT License", 1)
	validateDependencyTree(t, tree.Dependencies[1], "junit/junit/4.13.1", "Eclipse Public License 1.0", 2)
}

func validateDependencyTree(t *testing.T, tree *DependencyTree, wontProjectName, wontLicense string, wontDependencyCount int) {
	if tree.ProjectInfo.Name() != wontProjectName {
		t.Errorf("project name did not match, wont %s, got %s", wontProjectName, tree.ProjectInfo.Name())
	}
	if tree.Licenses[0].Name != wontLicense {
		t.Errorf("%s: license did not match, wont %s, got %s", wontProjectName, wontLicense, tree.Licenses[0].Name)
	}
	if len(tree.Dependencies) != wontDependencyCount {
		t.Errorf("%s: dependency parse error: wont dependency count: %d, got %d", wontProjectName, wontDependencyCount, len(tree.Dependencies))
	}
}
