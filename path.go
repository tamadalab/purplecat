package purplecat

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/asaskevich/govalidator"
	"github.com/go-resty/resty/v2"
	"github.com/tamadalab/purplecat/logger"
)

/*
 * Path shows the path for the build files or project directory located at local file or internet url.
 */
type Path struct {
	Path      string
	url       *url.URL
	supporter pathSupporter
}

/*
 * NewPath creates an pointer of Path represents the given location.
 */
func NewPath(path string) *Path {
	if govalidator.IsURL(path) {
		if url, err := url.Parse(path); err == nil {
			if url.Host != "" && url.Scheme != "" {
				return &Path{Path: path, url: url, supporter: &urlPathSupporter{}}
			}
		}
	}
	return &Path{Path: path, supporter: &localFilePathSupporter{}}
}

/*
 * Exists checks existtence of the receiver path.
 * If the receiver path shows url, and context denies the network access,
 * this function returns false.
 */
func (path *Path) Exists(context *Context) bool {
	return path.supporter.ExistFile(path, context)
}

/*
 * Base returns the base name of the receiver path.
 * Example:
 *     path := NewPath("/some/location/of/local/file")
 *     base := path.Base() // --> base is `file`
 */
func (path *Path) Base() string {
	return path.supporter.Base(path)
}

/*
 * Join add the given string to the receiver path, and returns new pointer of Path.
 * Example:
 *     path := NewPath("/some/location/of/local/file")
 *     path2 := path.Join("subfile") // --> base is `/some/location/of/local/file/subfile`
 */
func (path *Path) Join(append string) *Path {
	return NewPath(path.supporter.Join(path, append))
}

/*
 * Open returns ReadCloser for reading the content of the receiver path.
 */
func (path *Path) Open(context *Context) (io.ReadCloser, error) {
	return path.supporter.Open(path, context)
}

/*
 * Dir returns the directory part of the receiver path.
 * Example:
 *     path := NewPath("/some/location/of/local/file")
 *     dir := path.Dir() // --> dir is `/some/location/of/local`
 */
func (path *Path) Dir() *Path {
	return NewPath(path.supporter.Dir(path))
}

type pathSupporter interface {
	Base(path *Path) string
	Join(base *Path, append string) string
	Dir(path *Path) string
	ExistFile(path *Path, context *Context) bool
	Open(path *Path, context *Context) (io.ReadCloser, error)
}

type urlPathSupporter struct {
}

func (ups *urlPathSupporter) Base(urlPath *Path) string {
	return path.Base(urlPath.Path)
}

func (ups *urlPathSupporter) Join(base *Path, append string) string {
	newUrl := *base.url
	newUrl.Path = path.Join(newUrl.Path, append)
	return newUrl.String()
}

func (ups *urlPathSupporter) Dir(urlPath *Path) string {
	newUrl := *urlPath.url
	newUrl.Path = path.Dir(newUrl.Path)
	fmt.Printf("NewUrl(%s), OldUrl(%s)\n", newUrl.String(), urlPath.url.String())
	return newUrl.String()
}

func (ups *urlPathSupporter) ExistFile(path *Path, context *Context) bool {
	if !context.Allow(NETWORK_ACCESS) {
		return false
	}
	client := resty.New()
	request := client.NewRequest()
	response, err := request.Get(path.Path)
	result := (err != nil || response.StatusCode() != 404)
	logger.Debugf("Exist(%s): %v (%d)", path.url.String(), result, response.StatusCode())
	return result
}

func (ups *urlPathSupporter) Open(path *Path, context *Context) (io.ReadCloser, error) {
	if !context.Allow(NETWORK_ACCESS) {
		return nil, fmt.Errorf("network access denied")
	}
	logger.Debugf("Open(%s)", path.url.String())
	resp, err := http.Get(path.Path)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type localFilePathSupporter struct {
}

func (lfps *localFilePathSupporter) Base(filePath *Path) string {
	return filepath.Base(filePath.Path)
}

func (lfps *localFilePathSupporter) Join(base *Path, append string) string {
	return filepath.Join(base.Path, append)
}

func (lfps *localFilePathSupporter) Dir(filePath *Path) string {
	return filepath.Dir(filePath.Path)
}

func (lfps *localFilePathSupporter) ExistFile(path *Path, context *Context) bool {
	stat, err := os.Stat(path.Path)
	result := err == nil && stat.Mode().IsRegular()
	logger.Debugf("Exist(%s): %v", path.Path, result)
	return result
}

func (lfps *localFilePathSupporter) Open(path *Path, context *Context) (io.ReadCloser, error) {
	logger.Debugf("Open(%s)", path.Path)
	pom, err := os.Open(path.Path)
	if err != nil {
		return nil, err
	}
	return pom, nil
}
