package purplecat

import (
	"reflect"
	"testing"
)

func TestGenerateParser(t *testing.T) {
	testdata := []struct {
		path        string
		typeName    string
		successFlag bool
	}{
		{"./testdata/mavenproject", "mavenParser", true},
		{"./testdata/mavenproject/pom.xml", "mavenParser", true},
		{"./testdata/goproject", "goModParser", true},
		{"./testdata/goproject/go.mod", "goModParser", true},
		{"./testdata/unknownproject", "", false},
		{"./testdata/unknownproject/Makefile", "", false},
		{"./testdata/missingproject", "", false},
	}

	context := NewContext(false, "json", 1)
	for _, td := range testdata {
		parser, err := context.GenerateParser2(td.path)
		if err == nil && reflect.TypeOf(parser).Elem().Name() != td.typeName {
			t.Errorf(`GenerateParser("%s") resultant type did not match, wont %s, got %s`, td.path, td.typeName, reflect.TypeOf(parser).Elem())
		}
		if (err == nil) != td.successFlag { // isSuccess() != td.successFlag
			t.Errorf(`result of GenerateParser("%s") did not match, wont %v, got %v`, td.path, td.successFlag, !td.successFlag)
		}
	}
}
