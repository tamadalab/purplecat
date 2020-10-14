package purplecat

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/net/html/charset"
	"gopkg.in/xmlpath.v2"
)

const (
	LOCAL_MAVEN_REPOSITORY   = ".m2/repository"
	MAVEN_CENTRAL_REPOSITORY = "https://repo.maven.apache.org/maven2/"
)

type artifact struct {
	groupID    string
	artifactID string
	version    string
	valid      bool
	parent     *artifact
}

func newArtifact(node *xmlpath.Node) *artifact {
	groupID, ok1 := xmlpath.MustCompile("groupId").String(node)
	artifactID, ok2 := xmlpath.MustCompile("artifactId").String(node)
	version, ok3 := xmlpath.MustCompile("version").String(node)
	return &artifact{groupID: groupID, artifactID: artifactID, version: version, valid: ok1 && ok2 && ok3}
}

func (artifact *artifact) String() string {
	return fmt.Sprintf("%s/%s/%s", artifact.groupID, artifact.artifactID, artifact.version)
}

func (artifact *artifact) repoPath() string {
	path := strings.ReplaceAll(artifact.groupID, ".", "/")
	return fmt.Sprintf("%s/%s/%s", path, artifact.artifactID, artifact.version)
}

func (artifact *artifact) pomPath() string {
	return fmt.Sprintf("%s/%s-%s.pom", artifact.repoPath(), artifact.artifactID, artifact.version)
}

var cache = map[string]*DependencyTree{}

type MavenParser struct {
	context *Context
	cache   map[string]*DependencyTree
}

func (mp *MavenParser) Parse(path string) (*DependencyTree, error) {
	pomPath := filepath.Join(path, "pom.xml")
	if !FindFile(filepath.Join(pomPath)) {
		return nil, fmt.Errorf("%s: not maven project (pom.xml not found)", path)
	}
	return parsePom(pomPath, mp.context, 0)
}

func parsePom(pomPath string, context *Context, currentDepth int) (*DependencyTree, error) {
	if context.Depth < currentDepth {
		return nil, fmt.Errorf("over the parsing depth limit %d, current: %d", context.Depth, currentDepth)
	}
	pom, err := os.Open(pomPath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer pom.Close()

	// "ISO-8859-1" encoded xml parse error.
	// error message: xml: encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil
	// to resolve above problem, see https://stackoverflow.com/questions/6002619/unmarshal-an-iso-8859-1-xml-input-in-go
	decoder := xml.NewDecoder(pom)
	decoder.CharsetReader = charset.NewReaderLabel

	root, err := xmlpath.ParseDecoder(decoder)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return constructDependencyTree(root, filepath.Dir(pomPath), context, currentDepth)
}

func constructLocalParentPomPath(artifact *artifact) string {
	home, _ := homedir.Dir()
	return filepath.Join(home, LOCAL_MAVEN_REPOSITORY, artifact.parent.pomPath())
}

func constructDependencyTree(root *xmlpath.Node, path string, context *Context, currentDepth int) (*DependencyTree, error) {
	projectArtifact := parseProjectInfo(root)
	licenseNames, ok := findLicenseNamesFromPom(root)
	if !ok && projectArtifact.parent != nil {
		parentPomPath := constructLocalParentPomPath(projectArtifact)
		dep, err := parsePom(parentPomPath, context, currentDepth)
		if err != nil {
			return nil, err
		}
		licenseNames = dep.LicenseNames
	}
	project := &DependencyTree{ProjectName: projectArtifact.String(), LicenseNames: licenseNames}
	return parseDependency(project, root, context, currentDepth)
}

func findLicensesFromRemoteRepositoryPom(artifact *artifact, context *Context, currentDepth int) (*DependencyTree, error) {
	return nil, fmt.Errorf("not implemented yet.")
}

func findDependencyTreeFromLocalRepositoryPom(artifact *artifact, context *Context, currentDepth int) (*DependencyTree, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	pomPath := filepath.Join(home, LOCAL_MAVEN_REPOSITORY, artifact.pomPath())
	if !FindFile(pomPath) {
		return nil, fmt.Errorf("%s: file not found", pomPath)
	}
	dep, err := parsePom(pomPath, context, currentDepth+1)
	if err != nil {
		return nil, err
	}
	return dep, nil
}

func findDependencyTreeFromRepository(artifact *artifact, context *Context, currentDepth int) *DependencyTree {
	dep, err := findDependencyTreeFromLocalRepositoryPom(artifact, context, currentDepth)
	if err != nil && context.Allow(NETWORK_ACCESS) {
		dep, err = findLicensesFromRemoteRepositoryPom(artifact, context, currentDepth)
	}
	return dep
}

func nodeToDependencyTree(node *xmlpath.Node, context *Context, currentDepth int) *DependencyTree {
	artifact := newArtifact(node)
	return findDependencyTreeFromRepository(artifact, context, currentDepth)
}

func parseDependency(project *DependencyTree, root *xmlpath.Node, context *Context, currentDepth int) (*DependencyTree, error) {
	dependencyPath := xmlpath.MustCompile("/project/dependencies/dependency")
	for iter := dependencyPath.Iter(root); iter.Next(); {
		dependency := nodeToDependencyTree(iter.Node(), context, currentDepth)
		project.Dependencies = append(project.Dependencies, dependency)
	}
	return project, nil
}

func findLicenseNamesFromPom(root *xmlpath.Node) ([]string, bool) {
	licenseNamePath := xmlpath.MustCompile("/project/licenses/license/name")
	licenses := []string{}
	if licenseNamePath.Exists(root) {
		for iter := licenseNamePath.Iter(root); iter.Next(); {
			licenses = append(licenses, iter.Node().String())
		}
		return licenses, true
	}
	return licenses, false
}

func parentArtifact(root *xmlpath.Node) (*artifact, bool) {
	parentPath := xmlpath.MustCompile("/project/parent")
	if !parentPath.Exists(root) {
		return nil, false
	}
	iter := parentPath.Iter(root)
	if !iter.Next() {
		return nil, false
	}
	parent := newArtifact(iter.Node())
	return parent, parent.valid
}

func parseProjectInfo(root *xmlpath.Node) *artifact {
	projectPath := xmlpath.MustCompile("/project")
	iter := projectPath.Iter(root)
	if !iter.Next() {
		return nil
	}
	artifact := newArtifact(iter.Node())
	if !artifact.valid {
		parent, ok := parentArtifact(root)
		if ok {
			merge(artifact, parent)
			artifact.parent = parent
		}
	}
	return artifact
}

func merge(base, append *artifact) *artifact {
	if base.groupID == "" {
		base.groupID = append.groupID
	}
	if base.version == "" {
		base.version = append.version
	}
	if !base.valid {
		base.valid = base.groupID != "" && base.artifactID != "" && base.version != ""
	}
	return base
}
