package purplecat

import (
	"os"
	"testing"
)

func TestParseCacheType(t *testing.T) {
	testdata := []struct {
		giveString    string
		wontCacheType CacheType
	}{
		{"memory", MemoryCache},
		{"  Ref-Only", RefOnlyCache},
		{"DeFaUlT", DefaultCache},
		{"newdb", NewCache},
		{"unknown", -1},
	}
	for _, td := range testdata {
		gotCacheType := ParseCacheType(td.giveString)
		if gotCacheType != td.wontCacheType {
			t.Errorf(`ParseCacheType("%s") did not match, wont %s, got %s`, td.giveString, td.wontCacheType, gotCacheType)
		}
	}
}

func TestCacheTypeString(t *testing.T) {
	testdata := []struct {
		giveType   CacheType
		wontString string
	}{
		{MemoryCache, "memory"},
		{RefOnlyCache, "ref-only"},
		{DefaultCache, "default"},
		{NewCache, "newdb"},
	}
	for _, td := range testdata {
		gotString := td.giveType.String()
		if td.wontString != gotString {
			t.Errorf(`CacheType(%d).String() did not match, wont %s, got %s`, td.giveType, td.wontString, gotString)
		}
	}
}

func TestCacheStore(t *testing.T) {
	os.Setenv(CacheDBEnvName, "testdata/output.json")
	defer os.Unsetenv(CacheDBEnvName)
	cc, _ := NewCacheContext(DefaultCache)
	db, err := NewCacheDB(cc)
	if err != nil {
		t.Errorf("load error: %s", err.Error())
	}
	project := &Project{PName: "github.com/tamadalab/purplecat@v0.1.0", LicenseList: []*License{LicenseWTFPL}, context: cc, Deps: []string{}}
	ok := db.Register(project)
	if !ok {
		t.Errorf("register failed")
	}
	err2 := db.Store()
	if err2 != nil {
		t.Errorf("writer error: %s", err.Error())
	}
	if !existFile("testdata/output.json") {
		t.Errorf("dest file (testdata/output.json) is not exist")
	}
	defer os.Remove("testdata/output.json")
}

func TestCacheType(t *testing.T) {
	os.Setenv(CacheDBEnvName, "testdata/testcachetype.json")
	defer os.Unsetenv(CacheDBEnvName)
	testdata := []struct {
		giveType  CacheType
		wontError bool
		wontType  CacheType
	}{
		{MemoryCache, false, MemoryCache},
		{RefOnlyCache, false, RefOnlyCache},
		{DefaultCache, false, DefaultCache},
		{-1, true, -1},
	}
	for _, td := range testdata {
		context, err := NewCacheContext(td.giveType)
		if td.wontError && err == nil {
			t.Errorf("NewCacheContext(%d) wont error: %v, got %v", td.giveType, td.wontError, err)
		}
		if err == nil && context.cType != td.wontType {
			t.Errorf("NewContext(%d) result did not match, wont %d, got %d", td.giveType, td.wontType, context.cType)
		}
	}
}

var License0BSD = &License{Name: "BSD Zero Clause License", SpdxID: "0BSD", URL: "http://landley.net/toybox/license.html"}
var LicenseWTFPL = &License{Name: "Do What The F*ck You Want To Public License", SpdxID: "WTFPL", URL: "http://www.wtfpl.net/about/"}

func TestNoCacheDB(t *testing.T) {
	cc, _ := NewCacheContext(MemoryCache)
	db, err := NewCacheDB(cc)
	if err != nil {
		t.Errorf("NewCacheDB(%d) creation error: %v", MemoryCache, err)
	}
	_, found := db.Find("github.com/tamadalab/purplecat@v0.1.0")
	if found {
		t.Errorf("FindError")
	}
	project := &Project{PName: "github.com/tamadalab/purplecat@v0.1.0", LicenseList: []*License{LicenseWTFPL}, context: cc, Deps: []string{}}
	ok := db.Register(project)
	if !ok {
		t.Errorf("Register error")
	}
	_, found = db.Find("github.com/tamadalab/purplecat@v0.1.0")
	if !found {
		t.Errorf("FindError2")
	}
	_, success := db.Delete("github.com/tamadalab/purplecat@v0.1.0")
	if !success {
		t.Errorf("delete success!?")
	}
	if err := db.Store(); err != nil {
		t.Errorf("some error on store")
	}
}

func TestRefOnlyCacheDB(t *testing.T) {
	os.Setenv(CacheDBEnvName, "testdata/cachedb.json")
	defer os.Unsetenv(CacheDBEnvName)
	cc, _ := NewCacheContext(RefOnlyCache)
	key := "github.com/tamadalab/purplecat@v0.1.0"
	err := cc.Init()
	if err != nil {
		t.Errorf("NewCacheDB(%d) creation error: %v", RefOnlyCache, err)
	}
	_, found := cc.Find(key)
	if !found {
		t.Errorf("db.Find(%s) did not match, wont %v", key, true)
	}
	project := &Project{PName: key, LicenseList: []*License{License0BSD}, context: cc, Deps: []string{}}
	ok := cc.Register(project)
	if !ok {
		t.Errorf("db.Register(%s) did not match, wont %v", key, false)
	}
	_, found = cc.Find(key)
	if !found {
		t.Errorf("db.Find(%s) did not match, wont %v", key, true)
	}
	_, success := cc.Delete(key)
	if !success {
		t.Errorf("db.Delete(%s) did not match, wont %v", key, false)
	}
	if err := cc.Store(); err != nil {
		t.Errorf("db.Store() did not match, wont not nil, but nil")
	}
}
