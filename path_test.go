package purplecat

import "testing"

func TestExists(t *testing.T) {
	testdata := []struct {
		path        string
		denyNetwork bool
		exists      bool
	}{
		{"./path.go", true, true},
		{"./path.go", false, true},
		{"https://github.com/tamadalab/purplecat", false, true},
		{"https://github.com/tamadalab/purplecat", true, false},
	}

	for _, td := range testdata {
		context := NewContext(td.denyNetwork, "json", 1)
		path := NewPath(td.path)
		if path.Exists(context) != td.exists {
			t.Errorf("%s: exists wont %v, got %v (deny network %v)", td.path, td.exists, !td.exists, td.denyNetwork)
		}
	}
}

func TestBase(t *testing.T) {
	testdata := []struct {
		path     string
		wontBase string
	}{
		{"./path.go", "path.go"},
		{"https://github.com/tamadalab/purplecat", "purplecat"},
	}

	for _, td := range testdata {
		path := NewPath(td.path)
		gotBase := path.Base()
		if gotBase != td.wontBase {
			t.Errorf("%s.Base() did not match, wont %s, got %s", td.path, td.wontBase, gotBase)
		}
	}
}

func TestDir(t *testing.T) {
	testdata := []struct {
		path    string
		wontDir string
	}{
		{"./path.go", "."},
		{"https://github.com/tamadalab/purplecat", "https://github.com/tamadalab"},
	}

	for _, td := range testdata {
		path := NewPath(td.path)
		gotDir := path.Dir()
		if gotDir.Path != td.wontDir {
			t.Errorf("%s.Dir() did not match, wont %s, got %s", td.path, td.wontDir, gotDir.Path)
		}
	}
}

