package purplecat

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// defaultCacheDBPath represents the default cache database path.
const defaultCacheDBPath = "${HOME}/.config/purplecat/cachedb.json"

// CacheDBEnvName is the environment name for the cache database path.
const CacheDBEnvName = "PURPLECAT_CACHE_DB_PATH"

func findCacheDBPath() string {
	path := os.Getenv(CacheDBEnvName)
	if path == "" {
		path = defaultCacheDBPath
	}
	logger.Debugf("findCacheDBPath: %s", path)
	return path
}

// DefaultCacheDBPath returns the path of the default cache database.
func DefaultCacheDBPath() string {
	return normalizeCachePath(defaultCacheDBPath)
}

func normalizeCachePath(fromPath string) string {
	home, _ := homedir.Dir()
	return strings.ReplaceAll(fromPath, `${HOME}`, home)
}

// CacheDB is an interface of the cache database.
type CacheDB interface {
	Type() CacheType
	Find(projectName string) (foundProject *Project, found bool)
	Register(project *Project) bool
	Delete(projectName string) (deletedProject *Project, success bool)
	Store() error
	Clear() error
	Dump(io.Writer) error
}

// NewCacheDB creates suitable CacheDB by the given cache context.
func NewCacheDB(cType CacheType) (CacheDB, error) {
	return NewCacheDBWithPath(cType, "")
}

func newMemoryCacheDB(cType CacheType) CacheDB {
	return &memoryCacheDB{db: map[string]*Project{}, cType: cType}
}

func newDefaultCacheDB(dbpath string, cType CacheType) CacheDB {
	return &defaultCacheDB{DB: map[string]*Project{}, path: dbpath, cType: cType}
}

func backup(dbpath string) {
	now := time.Now()
	newpath := fmt.Sprintf("%s.%d", dbpath, now.Unix())
	os.Rename(dbpath, newpath)
}

func NewCacheDBWithPath(cType CacheType, dbpath string) (CacheDB, error) {
	if dbpath == "" {
		dbpath = findCacheDBPath()
	}
	normalizedCachePath := normalizeCachePath(dbpath)
	switch cType {
	case NewCache:
		backup(normalizedCachePath)
		return newDefaultCacheDB(normalizedCachePath, cType), nil
	default:
		fallthrough
	case MemoryCache:
		return newMemoryCacheDB(MemoryCache), nil
	case RefOnlyCache:
		defaultDB, err := loadDefaultCacheDB(normalizedCachePath)
		return &memoryCacheDB{db: defaultDB.DB, cType: RefOnlyCache}, err
	case DefaultCache:
		return loadDefaultCacheDB(normalizedCachePath)
	}
}

func loadDefaultCacheDB(path string) (*defaultCacheDB, error) {
	logger.Debugf("existFile(%s): %v", path, existFile(path))
	if !existFile(path) {
		return &defaultCacheDB{path: path, DB: map[string]*Project{}, cType: DefaultCache}, nil
	}
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return loadImpl(path, fp)
}

func loadImpl(path string, reader io.Reader) (*defaultCacheDB, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	cacheDB := &defaultCacheDB{path: path, DB: map[string]*Project{}}
	if err := json.Unmarshal(data, cacheDB); err != nil {
		return nil, err
	}
	return cacheDB, nil
}

type memoryCacheDB struct {
	cType CacheType
	db    map[string]*Project
}

func (ncdb *memoryCacheDB) Type() CacheType {
	return ncdb.cType
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

func (ncdb *memoryCacheDB) Dump(writer io.Writer) error {
	logger.Debug("memoryCacheDB does not support Dump")
	return nil
}

func (ncdb *memoryCacheDB) Clear() error {
	logger.Debug("clear memoryCacheDB")
	ncdb.db = map[string]*Project{}
	return nil
}

type defaultCacheDB struct {
	cType CacheType           `json:"-"`
	path  string              `json:"-"`
	DB    map[string]*Project `json:"cachedb"`
}

func (ddb *defaultCacheDB) Type() CacheType {
	return ddb.cType
}

func (ddb *defaultCacheDB) Find(projectName string) (*Project, bool) {
	project, ok := ddb.DB[projectName]
	if ok {
		logger.Infof("Cache found(%s: %v)", projectName, project)
		if project.context == nil {
			project.context = ddb
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

func (ddb *defaultCacheDB) Clear() error {
	logger.Infof("defaultCacheDB clear database")
	ddb.DB = map[string]*Project{}
	return nil
}

func (ddb *defaultCacheDB) Store() error {
	mkdirs(ddb.path)
	writer, err := os.OpenFile(ddb.path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	return ddb.Dump(writer)
}
func (ddb *defaultCacheDB) Dump(writer io.Writer) error {
	bytes, err := json.Marshal(ddb)
	if err != nil {
		return err
	}
	length, err := writer.Write(bytes)
	if err != nil {
		return err
	}
	if length != len(bytes) {
		return fmt.Errorf("cannot write fully data (write %d (%d) bytes)", length, len(bytes))
	}
	return nil
}
