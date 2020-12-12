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

func createContext(r *http.Request, cache *purplecat.CacheContext) *purplecat.Context {
	depth := parseDepth(r)
	context := purplecat.NewContext(false, "json", depth)
	context.Cache = cache
	return context
}

func runPurplecatHandler(cache *purplecat.CacheContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respondError(w, err)
			return
		}
		context := createContext(r, cache)
		project, err := runPurplecat(r, context)
		if err != nil {
			respondError(w, err)
		} else {
			respondJSON(w, context, project)
			cache.Store()
		}
	}
}

func clearCacheHandler(cache *purplecat.CacheContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Warnf("DELETE /purplecat/caches: not implement yet.")
	}
}

func wholeCacheHandler(cache *purplecat.CacheContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Warnf("GET /purplecat/caches: not implement yet.")
	}
}

func createRestAPI(cache *purplecat.CacheContext) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/purplecat/licenses", runPurplecatHandler(cache)).Methods("GET")
	router.HandleFunc("/purplecat/caches", clearCacheHandler(cache)).Methods("DELETE")
	router.HandleFunc("/purplecat/caches", wholeCacheHandler(cache)).Methods("GET")
	return router
}

func createCache() (*purplecat.CacheContext, error) {
	cache, err := purplecat.NewCacheContext(purplecat.DefaultCache)
	if err != nil {
		return nil, err
	}
	if err := cache.Init(); err != nil {
		return nil, err
	}
	return cache, nil
}

func startServer(router *mux.Router) int {
	logger.SetLevel(logger.INFO)
	logger.Infof("Listen server...")
	logger.Fatalf("server shutdown: %s", http.ListenAndServe(":8080", router))
	return 0
}

func goMain(args []string) int {
	cache, err := createCache()
	if err != nil {
		logger.Warnf("cache error: %s", err.Error())
		return 1
	}
	router := createRestAPI(cache)
	return startServer(router)
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
