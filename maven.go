package purplecat

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/antchfx/xmlquery"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/tamadalab/purplecat/logger"
)

const (
	localMavenRepository   = ".m2/repository"
	mavenCentralRepository = "repo.maven.apache.org/maven2/"
)

type artifact struct {
	groupID    string
	artifactID string
	version    string
	properties map[string]string
	parent     *artifact
}

func getStringByXPath(xpath string, node *xmlquery.Node) (string, bool) {
	targetNode, err := xmlquery.Query(node, xpath)
	if err != nil || targetNode == nil {
		return "", false
	}
	return strings.TrimSpace(targetNode.InnerText()), true
}

func newArtifact(groupID, artifactID, version string) *artifact {
	props := map[string]string{}
	props["project.version"] = version
	props["project.groupId"] = groupID
	props["project.artifactId"] = artifactID
	return &artifact{groupID: groupID, artifactID: artifactID, version: version, properties: props}
}

func newArtifactXPath(node *xmlquery.Node) *artifact {
	groupID, _ := getStringByXPath("./groupId", node)
	artifactID, _ := getStringByXPath("./artifactId", node)
	version, _ := getStringByXPath("./version", node)
	return newArtifact(groupID, artifactID, version)
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

type mavenParser struct {
	context *Context
}

func isPom(fileName string) bool {
	logger.Debugf("isPom(%s), %v", fileName, strings.HasSuffix(fileName, ".pom"))
	return fileName == "pom.xml" || strings.HasSuffix(fileName, ".pom")
}

// IsTarget returns true if the project located on the given path is maven project.
func (mp *mavenParser) IsTarget(path *Path, context *Context) bool {
	base := path.Base()
	if base == "pom.xml" || strings.HasSuffix(base, ".pom") {
		return path.Exists(context)
	}
	join := path.Join("pom.xml")
	return join.Exists(context)
}

// Parse parses the given path as pom.xml and returns the instance of Project.
func (mp *mavenParser) Parse(pomPath *Path) (*Project, error) {
	if !isPom(pomPath.Base()) {
		pomPath = pomPath.Join("pom.xml")
	}
	if !pomPath.Exists(mp.context) {
		return nil, fmt.Errorf("%s: not maven project (pom.xml not found)", pomPath.Path)
	}
	return parsePom(pomPath, mp.context, 0)
}

func readXML(pomPath *Path, context *Context) (*xmlquery.Node, error) {
	pom, err := pomPath.Open(context)
	if err != nil {
		return nil, err
	}
	defer pom.Close()
	return xmlquery.Parse(pom)
}

func parsePom(pomPath *Path, context *Context, currentDepth int) (*Project, error) {
	if context.Depth < currentDepth {
		return nil, fmt.Errorf("over the parsing depth limit %d, current: %d", context.Depth, currentDepth)
	}
	logger.Infof("parsePom(%s, %d)", pomPath.Path, currentDepth)
	doc, err := readXML(pomPath, context)
	if err != nil {
		return nil, err
	}
	return parsePomProject(doc, pomPath, context, currentDepth)
}

func parsePomProject(doc *xmlquery.Node, pomPath *Path, context *Context, currentDepth int) (*Project, error) {
	project, err := constructProject(doc, pomPath.Dir(), context, currentDepth)
	if err != nil {
		return nil, err
	}
	for _, dep := range project.Deps {
		if _, ok := context.SearchCache(dep); ok {
			continue
		}
		path, err := generatePomPath(dep, context)
		if err == nil {
			parsePom(path, context, currentDepth+1)
		}
	}
	return project, nil
	// return constructDependencyTree(doc, pomPath.Dir(), context, currentDepth)
}

func hitCache(artifact *artifact, context *Context) (*Project, bool) {
	if project, ok := context.SearchCache(artifact.Name()); ok {
		return project, true
	}
	return nil, false
}

func findParentLicense(parent *artifact, context *Context, currentDepth int) []*License {
	parentPomPath, err1 := constructPomPath(parent, context)
	if err1 != nil {
		return []*License{}
	}
	dep, err2 := parsePom(parentPomPath, context, currentDepth)
	if err2 != nil {
		return []*License{}
	}
	return dep.Licenses()
}

func readProperties(node *xmlquery.Node, artifact *artifact) {
	list, err := xmlquery.QueryAll(node, "/project/properties/*")
	if err == nil {
		for _, property := range list {
			artifact.properties[property.Data] = property.InnerText()
		}
	}
}

func constructProject(root *xmlquery.Node, path *Path, context *Context, currentDepth int) (*Project, error) {
	artifact := parseProjectInfo(root)
	if dep, ok := hitCache(artifact, context); ok {
		return dep, nil
	}
	readProperties(root, artifact)
	licenses, ok := findLicensesFromPom(artifact, root)
	if !ok && artifact.parent != nil {
		licenses = findParentLicense(artifact.parent, context, currentDepth)
	}
	project := context.NewProject(artifact.Name(), licenses)
	context.RegisterCache(project)
	return readDependencies(artifact, project, root)
}

func constructCentralRepoPomPath(artifact *artifact) *Path {
	url := path.Join(mavenCentralRepository, artifact.pomPath())
	return NewPath("https://" + url)
}

func constructLocalPomPath(artifact *artifact) *Path {
	home, _ := homedir.Dir()
	return NewPath(filepath.Join(home, localMavenRepository, artifact.pomPath()))
}

func generatePomPath(name string, context *Context) (*Path, error) {
	items := strings.Split(name, "/")
	artifact := newArtifact(items[0], items[1], items[2])
	return constructPomPath(artifact, context)
}

func constructPomPath(art *artifact, context *Context) (*Path, error) {
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

func updateByProps(target string, props map[string]string) string {
	for k, v := range props {
		target = strings.ReplaceAll(target, fmt.Sprintf(`${%s}`, k), v)
	}
	return target
}

func normalizeProject(target, base *artifact) *artifact {
	target.version = updateByProps(target.version, base.properties)
	target.artifactID = updateByProps(target.artifactID, base.properties)
	target.groupID = updateByProps(target.groupID, base.properties)

	return target
}

func readDependencies(artifact *artifact, project *Project, root *xmlquery.Node) (*Project, error) {
	dependencies, err := xmlquery.QueryAll(root, "/project/dependencies/dependency")
	if err != nil {
		return project, err
	}
	for _, dep := range dependencies {
		dependency := newArtifactXPath(dep)
		normalizeProject(dependency, artifact)
		project.Deps = append(project.Deps, dependency.Name())
	}
	return project, nil
}

func buildLicense(licenseNode *xmlquery.Node) *License {
	licenseName, _ := getStringByXPath("name", licenseNode)
	url, _ := getStringByXPath("url", licenseNode)
	return &License{Name: licenseName, URL: url}
}

func findLicensesFromPom(artifact *artifact, root *xmlquery.Node) (Licenses, bool) {
	list, err := xmlquery.QueryAll(root, "/project/licenses/license")
	licenses := []*License{}
	if err != nil {
		return licenses, false
	}
	for _, license := range list {
		licenses = append(licenses, buildLicense(license))
	}
	return licenses, len(licenses) > 0
}

func parentArtifact(root *xmlquery.Node) (*artifact, bool) {
	parentNode, err := xmlquery.Query(root, "/project/parent")
	if err != nil || parentNode == nil {
		return nil, false
	}
	parent := newArtifactXPath(parentNode)
	return parent, parent.isValid()
}

func parseProjectInfo(root *xmlquery.Node) *artifact {
	node, err := xmlquery.Query(root, "/project")
	if err != nil {
		return nil
	}
	artifact := newArtifactXPath(node)
	parent, ok := parentArtifact(root)
	if ok {
		merge(artifact, parent)
		artifact.parent = parent
	}
	return artifact
}

func merge(base, append *artifact) *artifact {
	if base.groupID == "" {
		base.groupID = append.groupID
		base.properties["project.groupId"] = append.groupID
	}
	if base.version == "" {
		base.version = append.version
		base.properties["project.version"] = append.version
	}
	return base
}
