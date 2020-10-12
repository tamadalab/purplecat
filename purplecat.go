package purplecat

import (
	"fmt"
	"io"
)

type DependencyTree struct {
	ProjectName  string
	LicenseName  string
	Dependencies []DependencyTree
}

/* ParseProject find the project and its license from given path.
 */
func ParseProject(projectPath string) *DependencyTree {

}

/* Println prints the dependency tree to given Writer.
 */
func (dt *DependencyTree) Println(out io.Writer) {
	dt.printImpl(out, "  ")
}

func (dt *DependencyTree) printImpl(out io.Writer, indent string) {
	line := fmt.Sprintf("%s%s: %s\n", indent, dt.ProjectName, dt.LicenseName)
	out.Write([]byte(line))
	for _, dependency := range dt.Dependencies {
		dependency.printImpl(out, indent+"  ")
	}
}
