package purplecat

import (
	"fmt"
	"io"
	"strings"
)

// Writer writes given Project to given io.Writer by some format.
type Writer interface {
	Write(tree *Project) error
}

type markdownWriter struct {
	Out io.Writer
}
type csvWriter struct {
	Out io.Writer
}
type jsonWriter struct {
	Out io.Writer
}
type yamlWriter struct {
	Out io.Writer
}
type tomlWriter struct {
	Out io.Writer
}
type xmlWriter struct {
	Out io.Writer
}

func (mw *markdownWriter) Write(tree *Project) error {
	return mw.writeImpl(tree, "")
}

func (mw *markdownWriter) writeImpl(tree *Project, indent string) error {
	line := fmt.Sprintf("%s* %s: [%s]\n", indent, tree.Name(), joinLicenseNames(tree))
	mw.Out.Write([]byte(line))
	for _, dependency := range tree.Dependencies() {
		if dependency != nil {
			mw.writeImpl(dependency, indent+"    ")
		}
	}
	return nil
}

func (cw *csvWriter) Write(tree *Project) error {
	cw.Out.Write([]byte("project-name,license-name,parent-project-name\n"))
	cw.writeImpl(tree, "")
	return nil
}

func (cw *csvWriter) writeImpl(tree *Project, parent string) {
	line := fmt.Sprintf("%s,%s,%s\n", tree.Name(), joinLicenseNames(tree), parent)
	cw.Out.Write([]byte(line))
	for _, dep := range tree.Dependencies() {
		if dep != nil {
			cw.writeImpl(dep, tree.Name())
		}
	}
}

func (jw *jsonWriter) Write(tree *Project) error {
	jw.Out.Write([]byte(jw.jsonString(tree)))
	return nil
}

func (jw *jsonWriter) dependency(deps Projects) string {
	array := []string{}
	for _, dep := range deps {
		if dep != nil {
			array = append(array, jw.jsonString(dep))
		}
	}
	return fmt.Sprintf(`,"dependencies":[%s]`, strings.Join(array, ","))
}

func (jw *jsonWriter) jsonString(tree *Project) string {
	dependentString := ""
	deps := tree.Dependencies()
	if len(deps) > 0 {
		dependentString = jw.dependency(deps)
	}
	return fmt.Sprintf(`{"project-name":"%s","license-names":["%s"]%s}`, tree.Name(), joinLicenseNames(tree), dependentString)
}

func joinLicenseNames(tree *Project) string {
	licenseNames := []string{}
	for _, license := range tree.Licenses() {
		licenseNames = append(licenseNames, license.Name)
	}
	return strings.Join(licenseNames, ",")
}

func (yw *yamlWriter) Write(tree *Project) error {
	yw.Out.Write([]byte("---\n"))
	yw.Out.Write([]byte(yw.string(tree, []string{"", "", ""})))
	yw.Out.Write([]byte("\n"))
	return nil
}

func (yw *yamlWriter) deps2string(tree *Project, indents []string) []string {
	array := []string{}
	for _, dep := range tree.Dependencies() {
		if dep != nil {
			newIndents := []string{indents[0] + "  ", indents[1], indents[2]}
			array = append(array, yw.string(dep, newIndents))
		}
	}
	return array
}

func (yw *yamlWriter) string(tree *Project, indents []string) string {
	base := fmt.Sprintf(`%s%sproject-name:%s
%s%slicense-names:[%s]`, indents[0], indents[1], tree.Name(), indents[0], indents[2], joinLicenseNames(tree))
	array := yw.deps2string(tree, indents)
	if len(array) > 0 {
		base = fmt.Sprintf(`%s
%s%sdependencies:
%s`, base, indents[0], indents[2], strings.Join(array, "\n"))
	}
	return base
}

func (tw *tomlWriter) Write(tree *Project) error {
	return nil
}

func (xw *xmlWriter) Write(tree *Project) error {
	data := fmt.Sprintf(`<?xml version="1.0"?>
<purplecat>
%s
</purplecat>`, xw.string(tree, "  "))
	xw.Out.Write([]byte(data))
	return nil
}

func (xw *xmlWriter) string(tree *Project, indent string) string {
	xmlLicenses := []string{}
	for _, license := range tree.Licenses() {
		xmlLicenses = append(xmlLicenses, indent+"  <license-name>"+license.Name+"</license-name>")
	}
	project := fmt.Sprintf(`%s<project-name>%s</project-name>
%s<license-names>
%s
%s</license-names>`, indent, tree.Name(), indent, strings.Join(xmlLicenses, "\n"), indent)
	array := []string{}
	for _, dep := range tree.Dependencies() {
		if dep != nil {
			array = append(array, fmt.Sprintf(`%s  <dependency>
%s    
%s  </dependency>`, indent, xw.string(dep, indent+"    "), indent))
		}
	}
	if len(array) > 0 {
		project = fmt.Sprintf(`%s
%s<dependencies>
%s
%s</dependencies>`, project, indent, strings.Join(array, "\n"), indent)
	}
	return project
}
