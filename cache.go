package purplecat

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/tamadalab/purplecat/logger"
)

// CacheType shows the cache type
type CacheType int

func (ct CacheType) String() string {
	switch ct {
	case MemoryCache:
		return "memory"
	case NewCache:
		return "newdb"
	case RefOnlyCache:
		return "ref-only"
	case DefaultCache:
		return "default"
	}
	return "unknown"
}

// ParseCacheType parses the given string and returns suitable CacheType.
func ParseCacheType(cacheTypeString string) CacheType {
	switch strings.ToLower(strings.TrimSpace(cacheTypeString)) {
	case "newdb":
		return NewCache
	case "memory":
		return MemoryCache
	case "ref-only":
		return RefOnlyCache
	case "default":
		return DefaultCache
	}
	return -1
}

const (
	// MemoryCache is the one of CacheType, creates new database, and stores it into no where.
	MemoryCache CacheType = iota + 1
	// NewCache is the one of CacheType, creates new database, and stores it into the certain file.
	NewCache
	// RefOnlyCache is the one of CacheType, loads the database from the certain location, and stores it into no where.
	RefOnlyCache
	// DefaultCache is the one of CacheType, loads the database from the certain location, and stores it to the location.
	DefaultCache
)

// DefaultCacheDBPath represents the default cache database path.
const DefaultCacheDBPath = "${HOME}/.config/purplecat/cachedb.json"

// CacheDBEnvName is the environment name for the cache database path.
const CacheDBEnvName = "PURPLECAT_CACHE_DB_PATH"

// CacheContext shows the settings for cache database.
type CacheContext struct {
	cType   CacheType
	path    string
	cacheDB CacheDB
}

func findCacheDBPath() string {
	path := os.Getenv(CacheDBEnvName)
	if path == "" {
		path = DefaultCacheDBPath
	}
	logger.Debugf("findCacheDBPath: %s", path)
	return path
}

func normalizeCachePath(fromPath string) string {
	home, _ := homedir.Dir()
	return strings.ReplaceAll(fromPath, `${HOME}`, home)
}

// NewCacheContext creates an instance of CacheContext.
func NewCacheContext(cType CacheType) (*CacheContext, error) {
	return NewCacheContextWithDBPath(cType, findCacheDBPath())
}

// NewCacheContextWithDBPath creates an instance of CacheContext by specifying the cache database path.
func NewCacheContextWithDBPath(cType CacheType, cachePath string) (*CacheContext, error) {
	if cachePath == "" {
		cachePath = findCacheDBPath()
	}
	normalizedCachePath := normalizeCachePath(cachePath)
	availableTypes := []CacheType{MemoryCache, RefOnlyCache, DefaultCache}
	for _, aType := range availableTypes {
		if cType == aType {
			return &CacheContext{cType: cType, path: normalizedCachePath}, nil
		}
	}
	return nil, fmt.Errorf("%d: unknown cache type", cType)
}

// CacheDB is an interface of the cache database.
type CacheDB interface {
	Find(projectName string) (foundProject *Project, found bool)
	Register(project *Project) bool
	Delete(projectName string) (deletedProject *Project, success bool)
	Store() error
}

// Find finds the licenses of the project with the given name.
func (cc *CacheContext) Find(projectName string) (foundProject *Project, found bool) {
	return cc.cacheDB.Find(projectName)
}

// Register registers the licenses of corresponding project, it returns true in successing the registration.
func (cc *CacheContext) Register(project *Project) bool {
	return cc.cacheDB.Register(project)
}

// Delete removes the licenses and project relation from this database.
// This function returns true in successing the deletion.
func (cc *CacheContext) Delete(projectName string) (deletedProject *Project, success bool) {
	return cc.cacheDB.Delete(projectName)
}

// Store saves this database into the certain location.
func (cc *CacheContext) Store() error {
	return cc.cacheDB.Store()
}

// Init initializes the CacheContext by createing the suitable CacheDB.
func (cc *CacheContext) Init() error {
	db, err := NewCacheDB(cc)
	if err == nil {
		cc.cacheDB = db
	}
	return err
}

// NewCacheDB creates suitable CacheDB by the given cache context.
func NewCacheDB(cc *CacheContext) (CacheDB, error) {
	logger.Debugf("NewCacheDB(%d)", cc.cType)
	if cc.cType == MemoryCache {
		return &memoryCacheDB{db: map[string]*Project{}}, nil
	} else if cc.cType == NewCache {
		return &defaultCacheDB{DB: map[string]*Project{}, context: cc}, nil
	}
	defaultDB, err := loadDefaultCacheDB(cc)
	if err == nil && cc.cType == RefOnlyCache {
		return &memoryCacheDB{db: defaultDB.DB}, err
	}
	return defaultDB, err
}

func loadDefaultCacheDB(cc *CacheContext) (*defaultCacheDB, error) {
	logger.Debugf("existFile(%s): %v", cc.path, existFile(cc.path))
	if !existFile(cc.path) {
		return &defaultCacheDB{context: cc, DB: map[string]*Project{}}, nil
	}
	fp, err := os.Open(cc.path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	cacheDB := &defaultCacheDB{context: cc, DB: map[string]*Project{}}
	if err := json.Unmarshal(data, cacheDB); err != nil {
		return nil, err
	}
	return cacheDB, nil
}

type memoryCacheDB struct {
	db map[string]*Project
}

func (ncdb *memoryCacheDB) Find(projectName string) (*Project, bool) {
	value, ok := ncdb.db[projectName]
	return value, ok
}

func (ncdb *memoryCacheDB) Delete(projectName string) (*Project, bool) {
	value, ok := ncdb.Find(projectName)
	if ok {
		delete(ncdb.db, projectName)
	}
	return value, ok
}

func (ncdb *memoryCacheDB) Register(project *Project) bool {
	if project != nil {
		ncdb.db[project.Name()] = project
	}
	return project != nil
}

func (ncdb *memoryCacheDB) Store() error {
	logger.Debug("memoryCacheDB does not support Store")
	return nil
}

type defaultCacheDB struct {
	context *CacheContext       `json:"-"`
	DB      map[string]*Project `json:"cachedb"`
}

func (ddb *defaultCacheDB) Find(projectName string) (*Project, bool) {
	project, ok := ddb.DB[projectName]
	if ok {
		logger.Infof("Cache found(%s: %v)", projectName, project)
		if project.context == nil {
			project.context = ddb.context
		}
	}
	return project, ok
}

func (ddb *defaultCacheDB) Delete(projectName string) (*Project, bool) {
	project, ok := ddb.DB[projectName]
	delete(ddb.DB, projectName)
	return project, ok
}

func (ddb *defaultCacheDB) Register(project *Project) bool {
	ddb.DB[project.Name()] = project
	return true
}

func mkdirs(path string) {
	logger.Debugf("mkdir(%s)", filepath.Dir(path))
	os.MkdirAll(filepath.Dir(path), 0755)
}

func (ddb *defaultCacheDB) Store() error {
	mkdirs(ddb.context.path)
	writer, err := os.OpenFile(ddb.context.path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	return storeImpl(writer, ddb)
}

func storeImpl(writer io.Writer, ddb *defaultCacheDB) error {
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
