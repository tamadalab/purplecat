package purplecat

import (
	"os"
	"testing"
)

func TestCacheType(t *testing.T) {
	testdata := []struct {
		giveType  CacheType
		wontError bool
		wontType  CacheType
	}{
		{NoCache, false, NoCache},
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
	cc, _ := NewCacheContext(NoCache)
	db, err := NewCacheDB(cc)
	if err != nil {
		t.Errorf("NewCacheDB(%d) creation error: %v", NoCache, err)
	}
	_, found := db.Find("github.com/tamadalab/purplecat@v0.1.0")
	if found {
		t.Errorf("FindError")
	}
	ok := db.Register("github.com/tamadalab/purplecat@v0.1.0", []*License{LicenseWTFPL})
	if !ok {
		t.Errorf("Register error")
	}
	_, found = db.Find("github.com/tamadalab/purplecat@v0.1.0")
	if found {
		t.Errorf("FindError2")
	}
	_, success := db.Delete("github.com/tamadalab/purplecat@v0.1.0")
	if !success {
		t.Errorf("delete success!?")
	}
	if err := db.Store(); err == nil {
		t.Errorf("store success!?")
	}
}

func TestRefOnlyCacheDB(t *testing.T) {
	os.Setenv(CacheDBEnvName, "testdata/cachedb.json")
	defer os.Unsetenv(CacheDBEnvName)
	cc, _ := NewCacheContext(RefOnlyCache)
	key := "github.com/tamadalab/purplecat@v0.1.0"
	db, err := NewCacheDB(cc)
	if err != nil {
		t.Errorf("NewCacheDB(%d) creation error: %v", NoCache, err)
	}
	_, found := db.Find(key)
	if !found {
		t.Errorf("db.Find(%s) did not match, wont %v", key, true)
	}
	ok := db.Register(key, []*License{License0BSD})
	if !ok {
		t.Errorf("db.Register(%s) did not match, wont %v", key, false)
	}
	_, found = db.Find(key)
	if !found {
		t.Errorf("db.Find(%s) did not match, wont %v", key, true)
	}
	_, success := db.Delete(key)
	if success {
		t.Errorf("db.Delete(%s) did not match, wont %v", key, false)
	}
	if err := db.Store(); err == nil {
		t.Errorf("db.Store() did not match, wont not nil, but nil")
	}
}
