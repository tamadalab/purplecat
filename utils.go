package purplecat

import "os"

func FindFile(path string) bool {
	stat, err := os.Stat(path)
	if err == nil && stat.Mode().IsRegular() {
		return true
	}
	return false
}