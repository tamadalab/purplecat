package purplecat

import (
	"encoding/xml"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/tamadalab/purplecat/logger"
	"golang.org/x/net/html/charset"
	"gopkg.in/xmlpath.v2"
)

const (
	localMavenRepository   = ".m2/repository"
	mavenCentralRepository = "repo.maven.apache.org/maven2/"
)

type artifact struct {
	groupID    string
	artifactID string
	version    string
	valid      bool
	parent     *artifact
}

func getStringByXPath(xpath string, node *xmlpath.Node) (string, bool) {
	str, ok := xmlpath.MustCompile(xpath).String(node)
	if ok {
		return strings.TrimSpace(str), ok
	}
	return str, ok
}

func newArtifact(node *xmlpath.Node) *artifact {
	groupID, ok1 := getStringByXPath("groupId", node)
	artifactID, ok2 := getStringByXPath("artifactId", node)
	version, ok3 := getStringByXPath("version", node)
	return &artifact{groupID: groupID, artifactID: artifactID, version: version, valid: ok1 && ok2 && ok3}
}

func (artifact *artifact) Name() string {
	return fmt.Sprintf("%s/%s/%s", artifact.groupID, artifact.artifactID, artifact.version)
}

func (artifact *artifact) repoPath() string {
	path := strings.ReplaceAll(artifact.groupID, ".", "/")
	return fmt.Sprintf("%s/%s/%s", path, artifact.artifactID, artifact.version)
}

func (artifact *artifact) pomPath() string {
	return fmt.Sprintf("%s/%s-%s.pom", artifact.repoPath(), artifact.artifactID, artifact.version)
}

func (artifact *artifact) isValid() bool {
	return artifact.groupID != "" && artifact.artifactID != "" && artifact.version != ""
}

var cache = map[string]*DependencyTree{}

type MavenParser struct {
	context *Context
}

func isPom(fileName string) bool {
	logger.Debugf("isPom(%s), %v", fileName, strings.HasSuffix(fileName, ".pom"))
	return fileName == "pom.xml" || strings.HasSuffix(fileName, ".pom")
}

func (mp *MavenParser) IsTarget(path *Path, context *Context) bool {
	base := path.Base()
	if base == "pom.xml" || strings.HasSuffix(base, ".pom") {
		return path.Exists(context)
	}
	join := path.Join("pom.xml")
	return join.Exists(context)
}

func (mp *MavenParser) Parse(pomPath *Path) (*DependencyTree, error) {
	if !isPom(pomPath.Base()) {
		pomPath = pomPath.Join("pom.xml")
	}
	if !pomPath.Exists(mp.context) {
		return nil, fmt.Errorf("%s: not maven project (pom.xml not found)", pomPath.Path)
	}
	return parsePom(pomPath, mp.context, 0)
}

func parsePom(pomPath *Path, context *Context, currentDepth int) (*DependencyTree, error) {
	if context.Depth < currentDepth {
		return nil, fmt.Errorf("over the parsing depth limit %d, current: %d", context.Depth, currentDepth)
	}
	logger.Infof("parsePom(%s, %d)", pomPath.Path, currentDepth)
	pom, err := pomPath.Open(context)
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
	return constructDependencyTree(root, pomPath.Dir(), context, currentDepth)
}

func hitCache(artifact *artifact) (*DependencyTree, bool) {
	dep, ok := cache[artifact.Name()]
	return dep, ok
}

func findParentLicense(parent *artifact, context *Context, currentDepth int) []*License {
	parentPomPath := constructLocalPomPath(parent)
	dep, err := parsePom(parentPomPath, context, currentDepth)
	if err != nil {
		return []*License{}
	}
	return dep.Licenses
}

func constructDependencyTree(root *xmlpath.Node, path *Path, context *Context, currentDepth int) (*DependencyTree, error) {
	projectArtifact := parseProjectInfo(root)
	if dep, ok := hitCache(projectArtifact); ok {
		return dep, nil
	}
	licenses, ok := findLicensesFromPom(root)
	if !ok && projectArtifact.parent != nil {
		licenses = findParentLicense(projectArtifact.parent, context, currentDepth)
	}
	project := &DependencyTree{ProjectInfo: projectArtifact, Licenses: licenses}
	cache[projectArtifact.Name()] = project
	return parseDependency(projectArtifact, project, root, context, currentDepth)
}

func constructCentralRepoPomPath(artifact *artifact) *Path {
	url := path.Join(mavenCentralRepository, artifact.pomPath())
	return NewPath("https://" + url)
}

func constructLocalPomPath(artifact *artifact) *Path {
	home, _ := homedir.Dir()
	return NewPath(filepath.Join(home, localMavenRepository, artifact.pomPath()))
}

func constructPom(art *artifact, context *Context) (*Path, error) {
	pomPathGenerators := []func(*artifact) *Path{
		constructLocalPomPath,
		constructCentralRepoPomPath,
	}
	for _, generator := range pomPathGenerators {
		pomPath := generator(art)
		if pomPath.Exists(context) {
			return pomPath, nil
		}
	}
	return nil, fmt.Errorf("%s: pom not found", art.Name())
}

func normalizeProject(target, base *artifact) *artifact {
	if target.version == "${project.version}" {
		target.version = base.version
	}
	if target.groupID == "${project.groupId}" {
		target.groupID = base.groupID
	}
	return target
}

func nodeToDependencyTree(base *artifact, node *xmlpath.Node, context *Context, currentDepth int) *DependencyTree {
	artifact := newArtifact(node)
	artifact = normalizeProject(artifact, base)
	pomPath, err := constructPom(artifact, context)
	if err != nil {
		return nil
	}
	dep, _ := parsePom(pomPath, context, currentDepth+1)
	return dep
}

func parseDependency(art *artifact, project *DependencyTree, root *xmlpath.Node, context *Context, currentDepth int) (*DependencyTree, error) {
	dependencyPath := xmlpath.MustCompile("/project/dependencies/dependency")
	for iter := dependencyPath.Iter(root); iter.Next(); {
		dependency := nodeToDependencyTree(art, iter.Node(), context, currentDepth)
		project.Dependencies = append(project.Dependencies, dependency)
	}
	return project, nil
}

func buildLicense(licenseNode *xmlpath.Node) *License {
	licenseName, _ := getStringByXPath("name", licenseNode)
	url, _ := getStringByXPath("url", licenseNode)
	return &License{Name: licenseName, URL: url}
}

func findLicensesFromPom(root *xmlpath.Node) ([]*License, bool) {
	licenseNamePath := xmlpath.MustCompile("/project/licenses/license")
	licenses := []*License{}
	if licenseNamePath.Exists(root) {
		for iter := licenseNamePath.Iter(root); iter.Next(); {
			licenses = append(licenses, buildLicense(iter.Node()))
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
		base.valid = base.isValid()
	}
	return base
}
