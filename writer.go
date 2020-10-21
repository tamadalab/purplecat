package purplecat

import (
	"fmt"
	"io"
	"strings"
)

type Writer interface {
	Write(tree *Project) error
}

type MarkdownWriter struct {
	Out io.Writer
}
type CsvWriter struct {
	Out io.Writer
}
type JSONWriter struct {
	Out io.Writer
}
type YamlWriter struct {
	Out io.Writer
}
type TomlWriter struct {
	Out io.Writer
}
type XMLWriter struct {
	Out io.Writer
}

func (mw *MarkdownWriter) Write(tree *Project) error {
	return mw.writeImpl(tree, "")
}

func (mw *MarkdownWriter) writeImpl(tree *Project, indent string) error {
	licenses := []string{}
	for _, name := range tree.Licenses {
		licenses = append(licenses, fmt.Sprintf(`"%s"`, name.Name))
	}
	line := fmt.Sprintf("%s* %s: [%s]\n", indent, tree.Info.Name(), strings.Join(licenses, ","))
	mw.Out.Write([]byte(line))
	for _, dependency := range tree.Dependencies {
		if dependency != nil {
			mw.writeImpl(dependency, indent+"    ")
		}
	}
	return nil
}

func (cw *CsvWriter) Write(tree *Project) error {
	cw.Out.Write([]byte("project-name,license-name,parent-project-name\n"))
	cw.writeImpl(tree, "")
	return nil
}

func (cw *CsvWriter) writeImpl(tree *Project, parent string) {
	array := []string{}
	for _, name := range tree.Licenses {
		array = append(array, fmt.Sprintf(`"%s"`, name.Name))
	}
	line := fmt.Sprintf("%s,%s,%s\n", tree.Info.Name(), array, parent)
	cw.Out.Write([]byte(line))
	for _, dep := range tree.Dependencies {
		if dep != nil {
			cw.writeImpl(dep, tree.Info.Name())
		}
	}
}

func (jw *JSONWriter) Write(tree *Project) error {
	jw.Out.Write([]byte(jw.jsonString(tree)))
	return nil
}

func (jw *JSONWriter) dependency(deps []*Project) string {
	array := []string{}
	for _, dep := range deps {
		if dep != nil {
			array = append(array, jw.jsonString(dep))
		}
	}
	return fmt.Sprintf(`,"dependencies":[%s]`, strings.Join(array, ","))
}

func (jw *JSONWriter) jsonString(tree *Project) string {
	dependentString := ""
	if len(tree.Dependencies) > 0 {
		dependentString = jw.dependency(tree.Dependencies)
	}
	return fmt.Sprintf(`{"project-name":"%s","license-names":["%s"]%s}`, tree.Info.Name(), joinLicenseNames(tree), dependentString)
}

func joinLicenseNames(tree *Project) string {
	licenseNames := []string{}
	for _, license := range tree.Licenses {
		licenseNames = append(licenseNames, license.Name)
	}
	return strings.Join(licenseNames, ",")
}

func (yw *YamlWriter) Write(tree *Project) error {
	yw.Out.Write([]byte("---\n"))
	yw.Out.Write([]byte(yw.string(tree, "", "", "")))
	yw.Out.Write([]byte("\n"))
	return nil
}

func (yw *YamlWriter) string(tree *Project, indent, header1, header2 string) string {
	base := fmt.Sprintf(`%s%sproject-name:%s
%s%slicense-names:[%s]`, indent, header1, tree.Info.Name(), indent, header2, joinLicenseNames(tree))
	array := []string{}
	for _, dep := range tree.Dependencies {
		if dep != nil {
			array = append(array, yw.string(dep, indent+"  ", "- ", "  "))
		}
	}
	if len(array) > 0 {
		base = fmt.Sprintf(`%s
%s%sdependencies:
%s`, base, indent, header2, strings.Join(array, "\n"))
	}
	return base
}

func (tw *TomlWriter) Write(tree *Project) error {
	return nil
}

func (xw *XMLWriter) Write(tree *Project) error {
	data := fmt.Sprintf(`<?xml version="1.0"?>
<purplecat>
%s
</purplecat>`, xw.string(tree, "  "))
	xw.Out.Write([]byte(data))
	return nil
}

func (xw *XMLWriter) string(tree *Project, indent string) string {
	xmlLicenses := []string{}
	for _, license := range tree.Licenses {
		xmlLicenses = append(xmlLicenses, indent+"  <license-name>"+license.Name+"</license-name>")
	}
	project := fmt.Sprintf(`%s<project-name>%s</project-name>
%s<license-names>
%s
%s</license-names>`, indent, tree.Info.Name(), indent, strings.Join(xmlLicenses, "\n"), indent)
	array := []string{}
	for _, dep := range tree.Dependencies {
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
