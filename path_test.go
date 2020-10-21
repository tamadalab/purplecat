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

func TestOpen(t *testing.T) {
	testdata := []struct {
		path        string
		denyNetwork bool
		successFlag bool
	}{
		{"./path.go", true, true},
		{"./path.go", false, true},
		{"./unknown_file", true, false},
		{"./unknown_file", false, false},
		{"https://github.com/tamadalab/purplecat", false, true},
		{"https://github.com/tamadalab/purplecat", true, false},
		{"https://github.com/tamadalab/unknownproject", false, false},
		{"https://github.com/tamadalab/unknownproject", true, false},
	}

	for _, td := range testdata {
		context := NewContext(td.denyNetwork, "json", 1)
		path := NewPath(td.path)
		file, err := path.Open(context)
		if (err == nil) != td.successFlag {
			t.Errorf("%s: open wont %v, got %v (deny network %v)", td.path, td.successFlag, !td.successFlag, td.denyNetwork)
		}
		if err == nil {
			defer file.Close()
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

func TestJoin(t *testing.T) {
	testdata := []struct {
		basePath   string
		appendPath string
		wontResult string
	}{
		{"cmd", "purplecat/main.go", "cmd/purplecat/main.go"},
		{"https://github.com/tamadalab/", "purplecat", "https://github.com/tamadalab/purplecat"},
	}

	for _, td := range testdata {
		path := NewPath(td.basePath)
		gotPath := path.Join(td.appendPath)
		if gotPath.Path != td.wontResult {
			t.Errorf(`"%s".Join("%s") did not match, wont %s, got %s`, td.basePath, td.appendPath, td.wontResult, gotPath.Path)
		}
	}
}
