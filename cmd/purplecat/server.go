package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tamadalab/purplecat"
	"github.com/tamadalab/purplecat/logger"
)

func respond(w http.ResponseWriter, statusCode int, content []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(content)
}

func respondError(w http.ResponseWriter, statusCode int, err error) {
	jsonContent, _ := json.Marshal(map[string]string{"error": err.Error()})
	respond(w, statusCode, jsonContent)
}

func parseDepth(r *http.Request) int {
	depthString := r.FormValue("depth")
	if depthString != "" {
		depth, err := strconv.Atoi(depthString)
		if err == nil {
			return depth
		}
		logger.Warnf("depth=%s: cannot convert to integer: %s", depthString, err.Error())
	}
	return 1
}

func respondJSON(w http.ResponseWriter, context *purplecat.Context, project *purplecat.Project) {
	buffer := bytes.NewBuffer([]byte{})
	writer, err := context.NewWriter(buffer)
	if err != nil {
		respondError(w, 500, err)
		return
	}
	writer.Write(project)
	respond(w, 200, buffer.Bytes())
}

type formPathSupporter struct {
	data *bytes.Buffer
}

func (fp *formPathSupporter) Base(path *purplecat.Path) string {
	return path.Path
}

func (fp *formPathSupporter) Dir(path *purplecat.Path) string {
	return ""
}

func (fp *formPathSupporter) Join(path *purplecat.Path, append string) string {
	return append
}

func (fp *formPathSupporter) ExistFile(path *purplecat.Path, context *purplecat.Context) bool {
	return true
}

func (fp *formPathSupporter) Open(*purplecat.Path, *purplecat.Context) (io.ReadCloser, error) {
	return fp, nil
}

func (fp *formPathSupporter) Read(buffer []byte) (int, error) {
	return fp.data.Read(buffer)
}

func (fp *formPathSupporter) Close() error {
	return nil
}

func runPurplecatByPost(w http.ResponseWriter, r *http.Request, context *purplecat.Context) (*purplecat.Project, error) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return nil, fmt.Errorf("Content-Type was not set in the request header")
	} else if contentType != "application/xml" {
		return nil, fmt.Errorf("Supported Content-Type is only \"application/xml\" in the current implementation")
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}
	supporter := &formPathSupporter{data: bytes.NewBuffer(data)}
	path := purplecat.NewPathWithSupporter("pom.xml", supporter)
	parser, err := context.GenerateParser(path)
	if err != nil {
		return nil, err
	}
	return parser.Parse(path)
}

func runPurplecatByGet(w http.ResponseWriter, r *http.Request, context *purplecat.Context) (*purplecat.Project, error) {
	target := r.FormValue("target")
	if target == "" {
		return nil, fmt.Errorf(`query param "target" is mandatory`)
	}
	path := purplecat.NewPath(target)
	parser, err := context.GenerateParser(path)
	if err != nil {
		return nil, err
	}
	return parser.Parse(path)
}

func createContext(r *http.Request, cache purplecat.CacheDB) *purplecat.Context {
	depth := parseDepth(r)
	context := purplecat.NewContext(false, "json", depth)
	context.Cache = cache
	return context
}

func updateHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func runPurplecatHandler(cache purplecat.CacheDB, method string, runFunc func(http.ResponseWriter, *http.Request, *purplecat.Context) (*purplecat.Project, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("%s /purplecat/licenses", method)
		if err := r.ParseForm(); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		context := createContext(r, cache)
		project, err := runFunc(w, r, context)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
		} else {
			updateHeader(w, r)
			respondJSON(w, context, project)
			cache.Store()
		}
	}
}

func runPurplecatPostHandler(cache purplecat.CacheDB) func(http.ResponseWriter, *http.Request) {
	return runPurplecatHandler(cache, "POST", runPurplecatByPost)
}

func runPurplecatGetHandler(cache purplecat.CacheDB) func(http.ResponseWriter, *http.Request) {
	return runPurplecatHandler(cache, "GET", runPurplecatByGet)
}

func clearCacheHandler(cache purplecat.CacheDB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("DELETE caches")
		cache.Clear()
		data, _ := json.Marshal(map[string]string{"message": "ok"})
		respond(w, 200, data)
	}
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Infof("OPTIONS: %s", r.URL)
	updateHeader(w, r)
	w.Header().Set("Access-Control-Request-Method", "POST,GET,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "origin, accept, X-PINGOTHER, Content-Type")
	w.WriteHeader(200)
	w.Write([]byte{})
}

func wholeCacheHandler(cache purplecat.CacheDB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("GET caches")
		buffer := bytes.NewBuffer([]byte{})
		cache.Dump(buffer)
		respond(w, 200, buffer.Bytes())
	}
}

func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("request uri: %s", r.RequestURI)
		h.ServeHTTP(w, r)
	}
}

func createFileServer() http.Handler {
	dirs := []string{
		"docs",
		"/opt/purplecat/docs",
		"/usr/local/opt/purplecat/docs",
	}
	for _, dir := range dirs {
		if existDir(dir) {
			logger.Debugf("serve %s", dir)
			return wrapHandler(http.StripPrefix("/purplecat/", http.FileServer(http.Dir(dir))))
		}
	}
	return nil
}

func existDir(dir string) bool {
	stat, err := os.Stat(dir)
	return err == nil && stat.IsDir()
}

func createRestAPI(cache purplecat.CacheDB) *mux.Router {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/purplecat/api/").Subrouter()
	subRouter.HandleFunc("/licenses", runPurplecatGetHandler(cache)).Methods("GET")
	subRouter.HandleFunc("/licenses", runPurplecatPostHandler(cache)).Methods("POST")
	subRouter.HandleFunc("/licenses", optionsHandler).Methods("OPTIONS")
	subRouter.HandleFunc("/caches", clearCacheHandler(cache)).Methods("DELETE")
	subRouter.HandleFunc("/caches", wholeCacheHandler(cache)).Methods("GET")
	subRouter.HandleFunc("/caches", optionsHandler).Methods("OPTIONS")
	router.PathPrefix("/purplecat/").Handler(createFileServer())
	return router
}

func startServer(router *mux.Router, opts *serverOpts) int {
	portString := fmt.Sprintf(":%d", opts.port)
	logger.SetLevel(logger.INFO)
	logger.Infof("Listen server at port %d...", opts.port)
	logger.Fatalf("server shutdown: %s", http.ListenAndServe(portString, router))
	return 0
}

func (server *serverOpts) StartServer(common *commonOpts, cache purplecat.CacheDB) int {
	router := createRestAPI(cache)
	return startServer(router, server)
}
