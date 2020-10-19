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

func newArtifact(node *xmlquery.Node) *artifact {
	groupID, _ := getStringByXPath("./groupId", node)
	artifactID, _ := getStringByXPath("./artifactId", node)
	version, _ := getStringByXPath("./version", node)
	props := map[string]string{}
	props["project.version"] = version
	props["project.groupId"] = groupID
	props["project.artifactId"] = artifactID
	return &artifact{groupID: groupID, artifactID: artifactID, version: version, properties: props}
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

var cache = map[string]*Project{}

/*MavenParser is the instance of Parser for parsing pom.xml.*/
type MavenParser struct {
	context *Context
}

func isPom(fileName string) bool {
	logger.Debugf("isPom(%s), %v", fileName, strings.HasSuffix(fileName, ".pom"))
	return fileName == "pom.xml" || strings.HasSuffix(fileName, ".pom")
}

/*IsTarget checks given path is the target for receiver Parser.*/
func (mp *MavenParser) IsTarget(path *Path, context *Context) bool {
	base := path.Base()
	if base == "pom.xml" || strings.HasSuffix(base, ".pom") {
		return path.Exists(context)
	}
	join := path.Join("pom.xml")
	return join.Exists(context)
}

/*Parse parses the given path as pom.xml and returns the instance of Project. */
func (mp *MavenParser) Parse(pomPath *Path) (*Project, error) {
	if !isPom(pomPath.Base()) {
		pomPath = pomPath.Join("pom.xml")
	}
	if !pomPath.Exists(mp.context) {
		return nil, fmt.Errorf("%s: not maven project (pom.xml not found)", pomPath.Path)
	}
	return parsePom(pomPath, mp.context, 0)
}

func parsePom(pomPath *Path, context *Context, currentDepth int) (*Project, error) {
	if context.Depth < currentDepth {
		return nil, fmt.Errorf("over the parsing depth limit %d, current: %d", context.Depth, currentDepth)
	}
	logger.Infof("parsePom(%s, %d)", pomPath.Path, currentDepth)
	pom, err := pomPath.Open(context)
	if err != nil {
		return nil, err
	}
	defer pom.Close()

	doc, err := xmlquery.Parse(pom)
	if err != nil {
		return nil, err
	}
	return constructDependencyTree(doc, pomPath.Dir(), context, currentDepth)
	// "ISO-8859-1" encoded xml parse error.
	// error message: xml: encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil
	// to resolve above problem, see https://stackoverflow.com/questions/6002619/unmarshal-an-iso-8859-1-xml-input-in-go
	// decoder := xml.NewDecoder(pom)
	// decoder.CharsetReader = charset.NewReaderLabel

	// root, err := xmlpath.ParseDecoder(decoder)
	// if err != nil {
	// 	return nil, err
	// }
	// return constructDependencyTree(root, pomPath.Dir(), context, currentDepth)
}

func hitCache(artifact *artifact) (*Project, bool) {
	dep, ok := cache[artifact.Name()]
	return dep, ok
}

func findParentLicense(parent *artifact, context *Context, currentDepth int) []*License {
	parentPomPath, err1 := constructPom(parent, context)
	if err1 != nil {
		return []*License{}
	}
	dep, err2 := parsePom(parentPomPath, context, currentDepth)
	if err2 != nil {
		return []*License{}
	}
	return dep.Licenses
}

func readProperties(node *xmlquery.Node, artifact *artifact) {
	list, err := xmlquery.QueryAll(node, "/project/properties/*")
	if err == nil {
		for _, property := range list {
			artifact.properties[property.Data] = property.InnerText()
		}
	}
}

func constructDependencyTree(root *xmlquery.Node, path *Path, context *Context, currentDepth int) (*Project, error) {
	projectArtifact := parseProjectInfo(root)
	if dep, ok := hitCache(projectArtifact); ok {
		return dep, nil
	}
	readProperties(root, projectArtifact)
	licenses, ok := findLicensesFromPom(projectArtifact, root)
	if !ok && projectArtifact.parent != nil {
		licenses = findParentLicense(projectArtifact.parent, context, currentDepth)
	}
	project := &Project{Info: projectArtifact, Licenses: licenses}
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

func nodeToDependencyTree(base *artifact, node *xmlquery.Node, context *Context, currentDepth int) *Project {
	artifact := newArtifact(node)
	artifact = normalizeProject(artifact, base)
	pomPath, err := constructPom(artifact, context)
	if err != nil {
		return nil
	}
	dep, _ := parsePom(pomPath, context, currentDepth+1)
	return dep
}

func parseDependency(art *artifact, project *Project, root *xmlquery.Node, context *Context, currentDepth int) (*Project, error) {
	dependencies, err := xmlquery.QueryAll(root, "/project/dependencies/dependency")
	if err != nil {
		return project, err
	}
	for _, dep := range dependencies {
		dependency := nodeToDependencyTree(art, dep, context, currentDepth)
		project.Dependencies = append(project.Dependencies, dependency)
	}
	return project, nil
}

func buildLicense(licenseNode *xmlquery.Node) *License {
	licenseName, _ := getStringByXPath("name", licenseNode)
	url, _ := getStringByXPath("url", licenseNode)
	return &License{Name: licenseName, URL: url}
}

func findLicensesFromPom(artifact *artifact, root *xmlquery.Node) ([]*License, bool) {
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
	if err != nil {
		return nil, false
	}
	parent := newArtifact(parentNode)
	return parent, parent.isValid()
}

func parseProjectInfo(root *xmlquery.Node) *artifact {
	node, err := xmlquery.Query(root, "/project")
	// node, err := xmlquery.Query(root, "/project/(groupId,artifactId,version)")
	if err != nil {
		return nil
	}
	artifact := newArtifact(node)
	if !artifact.isValid() {
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
	return base
}
