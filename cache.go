package purplecat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/tamadalab/purplecat/logger"
)

// CacheType shows the cache type
type CacheType int

const (
	// NoCache is the one of CacheType, shows no cache.
	NoCache CacheType = iota + 1
	// RefOnlyCache is the one of CacheType, finds licenses from the cache DB, and no store the licenses to the cache DB.
	RefOnlyCache
	// DefaultCache is the one of CacheType, finds/stores the licenses from/to the cache DB.
	DefaultCache
)

// DefaultCacheDBPath represents the default cache database path.
const DefaultCacheDBPath = "${HOME}/.config/purplecat/cachedb.json"

// CacheDBEnvName is the environment name for the cache database path.
const CacheDBEnvName = "PURPLECAT_CACHE_DB_PATH"

// CacheContext shows the settings for cache database.
type CacheContext struct {
	cType CacheType
	path  string
}

func findCacheDBPath() string {
	path := os.Getenv(CacheDBEnvName)
	if path != "" {
		return path
	}
	home, _ := homedir.Dir()
	return strings.ReplaceAll(DefaultCacheDBPath, "${HOME}", home)
}

/*NewCacheContext creates an instance of CacheContext.*/
func NewCacheContext(cType CacheType) (*CacheContext, error) {
	availableTypes := []CacheType{NoCache, RefOnlyCache, DefaultCache}
	for _, aType := range availableTypes {
		if cType == aType {
			return &CacheContext{cType: cType, path: findCacheDBPath()}, nil
		}
	}
	return nil, fmt.Errorf("%d: unknown cache type", cType)
}

// CacheDB is an interface of the cache database.
type CacheDB interface {
	Find(projectName string) (foundLicenses []*License, found bool)
	Register(projectName string, licenses []*License) bool
	Delete(projectName string) (deletedLicenses []*License, success bool)
	Store() error
}

// NewCacheDB creates suitable CacheDB by the given cache context.
func NewCacheDB(cc *CacheContext) (CacheDB, error) {
	if cc.cType == NoCache {
		return &noCacheDB{}, nil
	}
	defaultDB, err := loadDefaultCacheDB(cc)
	if err == nil && cc.cType == RefOnlyCache {
		return &refOnlyCacheDB{cache: defaultDB}, err
	}
	return defaultDB, err
}

func loadDefaultCacheDB(cc *CacheContext) (*defaultCacheDB, error) {
	fp, err := os.Open(cc.path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	cacheDB := &defaultCacheDB{context: cc, DB: map[string][]*License{}}
	if err := json.Unmarshal(data, cacheDB); err != nil {
		return nil, err
	}
	return cacheDB, nil
}

type noCacheDB struct {
}

func (ncdb *noCacheDB) Find(projectName string) ([]*License, bool) {
	return []*License{}, false
}

func (ncdb *noCacheDB) Delete(projectName string) ([]*License, bool) {
	return []*License{}, true
}

func (ncdb *noCacheDB) Register(projectName string, licenses []*License) bool {
	return true
}

func (ncdb *noCacheDB) Store() error {
	return fmt.Errorf("NoCache does not support Store")
}

type refOnlyCacheDB struct {
	cache CacheDB
}

func (rocdb *refOnlyCacheDB) Find(projectName string) ([]*License, bool) {
	return rocdb.cache.Find(projectName)
}

func (rocdb *refOnlyCacheDB) Delete(projectName string) ([]*License, bool) {
	return []*License{}, false
}

func (rocdb *refOnlyCacheDB) Register(projectName string, licenses []*License) bool {
	return true
}

func (rocdb *refOnlyCacheDB) Store() error {
	return fmt.Errorf("RefOnlyCache did not support Store")
}

type defaultCacheDB struct {
	context *CacheContext         `json:"-"`
	DB      map[string][]*License `json:"cachedb"`
}

func (ddb *defaultCacheDB) Find(projectName string) ([]*License, bool) {
	licenses, ok := ddb.DB[projectName]
	if ok {
		logger.Infof("Cache found(%s: %v)", projectName, licenses)
	}
	return licenses, ok
}

func (ddb *defaultCacheDB) Delete(projectName string) ([]*License, bool) {
	licenses, ok := ddb.DB[projectName]
	delete(ddb.DB, projectName)
	return licenses, ok
}

func (ddb *defaultCacheDB) Register(projectName string, licenses []*License) bool {
	ddb.DB[projectName] = licenses
	return true
}

func (ddb *defaultCacheDB) Store() error {
	writer, err := os.OpenFile(ddb.context.path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	bytes, err := json.Marshal(ddb)
	if err != nil {
		return err
	}
	length, err := writer.Write(bytes)
	if length != len(bytes) {
		return fmt.Errorf("cannot write fully data")
	}
	return nil
}
