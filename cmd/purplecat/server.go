package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func respondError(w http.ResponseWriter, err error) {
	jsonContent, _ := json.Marshal(map[string]string{"error": err.Error()})
	respond(w, 500, jsonContent)
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
		respondError(w, err)
		return
	}
	writer.Write(project)
	respond(w, 200, buffer.Bytes())
}

func runPurplecat(r *http.Request, context *purplecat.Context) (*purplecat.Project, error) {
	target := r.FormValue("target")
	if target == "" {
		return nil, fmt.Errorf(`query param "target" is mandatory`)
	}
	parser, err := context.GenerateParser(target)
	if err != nil {
		return nil, err
	}
	path := purplecat.NewPath(target)
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
	w.Header().Set("Access-Control-Allow-Methods", "GET,DELETE")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func runPurplecatHandler(cache purplecat.CacheDB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("GET /purplecat/licenses")
		if err := r.ParseForm(); err != nil {
			respondError(w, err)
			return
		}
		context := createContext(r, cache)
		project, err := runPurplecat(r, context)
		if err != nil {
			respondError(w, err)
		} else {
			updateHeader(w, r)
			respondJSON(w, context, project)
			cache.Store()
		}
	}
}

func clearCacheHandler(cache purplecat.CacheDB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("DELETE caches")
		cache.Clear()
		data, _ := json.Marshal(map[string]string{"message": "ok"})
		respond(w, 200, data)
	}
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
	subRouter.HandleFunc("/licenses", runPurplecatHandler(cache)).Methods("GET")
	subRouter.HandleFunc("/caches", clearCacheHandler(cache)).Methods("DELETE")
	subRouter.HandleFunc("/caches", wholeCacheHandler(cache)).Methods("GET")
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
