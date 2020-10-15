package purplecat

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/asaskevich/govalidator"
	"github.com/go-resty/resty/v2"
	"github.com/tamadalab/purplecat/logger"
)

type Path struct {
	Path      string
	Supporter PathSupporter
}

func NewPath(path string) *Path {
	var supporter PathSupporter = &LocalFilePathSupporter{}
	if govalidator.IsURL(path) {
		supporter = &UrlPathSupporter{}
	}
	return &Path{Path: path, Supporter: supporter}
}

func (path *Path) Exists(context *Context) bool {
	return path.Supporter.ExistFile(path.Path, context)
}

func (path *Path) Base() string {
	return path.Supporter.Base(path.Path)
}

func (path *Path) Join(append string) *Path {
	return &Path{Path: path.Supporter.Join(path.Path, append), Supporter: path.Supporter}
}

func (path *Path) Open(context *Context) (io.ReadCloser, error) {
	return path.Supporter.Open(path.Path, context)
}

func (path *Path) Dir() *Path {
	return &Path{Path: path.Supporter.Dir(path.Path), Supporter: path.Supporter}
}

type PathSupporter interface {
	Base(path string) string
	Join(base, append string) string
	Dir(path string) string
	ExistFile(path string, context *Context) bool
	Open(path string, context *Context) (io.ReadCloser, error)
}

type UrlPathSupporter struct {
}

func (ups *UrlPathSupporter) Base(urlPath string) string {
	return path.Base(urlPath)
}

func (ups *UrlPathSupporter) Join(base, append string) string {
	return path.Join(base, append)
}

func (ups *UrlPathSupporter) Dir(urlPath string) string {
	return path.Dir(urlPath)
}

func (ups *UrlPathSupporter) ExistFile(path string, context *Context) bool {
	if !context.Allow(NETWORK_ACCESS) {
		return false
	}
	client := resty.New()
	request := client.NewRequest()
	response, err := request.Get(path)
	result := (err != nil || response.StatusCode() != 404)
	logger.Debugf("Exist(%s): %v (%d)", path)
	return result
}

func (ups *UrlPathSupporter) Open(path string, context *Context) (io.ReadCloser, error) {
	if !context.Allow(NETWORK_ACCESS) {
		return nil, fmt.Errorf("network access denied")
	}
	logger.Debugf("Open(%s)", path)
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type LocalFilePathSupporter struct {
}

func (lfps *LocalFilePathSupporter) Base(filePath string) string {
	return filepath.Base(filePath)
}

func (lfps *LocalFilePathSupporter) Join(base, append string) string {
	return filepath.Join(base, append)
}

func (lfps *LocalFilePathSupporter) Dir(filePath string) string {
	return filepath.Dir(filePath)
}

func (lfps *LocalFilePathSupporter) ExistFile(path string, context *Context) bool {
	stat, err := os.Stat(path)
	result := err == nil && stat.Mode().IsRegular()
	logger.Debugf("Exist(%s): %v", path, result)
	return result
}

func (lfps *LocalFilePathSupporter) Open(path string, context *Context) (io.ReadCloser, error) {
	logger.Debugf("Open(%s)", path)
	pom, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return pom, nil
}
