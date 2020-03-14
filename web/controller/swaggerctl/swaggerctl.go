package swaggerctl

import (
	"github.com/StephanHCB/go-autumn-web-swagger-ui"
	"github.com/go-chi/chi"
	"github.com/go-http-utils/headers"
	"net/http"
	"os"
	"path/filepath"
)

func SetupSwaggerRoutes(server chi.Router) {
	// 	serve swagger-ui and swagger.json
	AddStaticHttpFilesystemRoute(server, auwebswaggerui.Assets, "/swagger-ui")
	AddStaticSingleFileRoute(server, "docs", "/", "swagger.json")
}

// serve static files from an http.FileSystem instance via a chi router
//
// inspired by https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
//
// parameters:
// - uriPath under which sub-url they should be served, if you leave out trailing slash, also adds a redirect
//   example: "/swagger-ui"
func AddStaticHttpFilesystemRoute(server chi.Router, fs http.FileSystem, uriPath string) {
	strippedFs := http.StripPrefix(uriPath, http.FileServer(fs))

	if hasNoTrailingSlash(uriPath) {
		server.Get(uriPath, http.RedirectHandler(uriPath+"/", 301).ServeHTTP)
		uriPath += "/"
	}
	uriPath += "*"

	server.Get(uriPath, func(w http.ResponseWriter, r *http.Request) {
		strippedFs.ServeHTTP(w, r)
	})
}

// serve static files from a directory via a chi router
//
// inspired by https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
//
// parameters:
// - relativeFilesystemPath where to find the file(s) to serve relative to the current working directory
//   example: "third_party/swagger_ui"
//   note: do NOT add a trailing slash
// - uriPath under which sub-url they should be served, if you leave out trailing slash, also adds a redirect
//   example: "/swagger-ui"
func AddStaticDirectoryRoute(server chi.Router, relativeFilesystemPath string, uriPath string) {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, relativeFilesystemPath)
	root := http.Dir(filesDir)
	fs := http.StripPrefix(uriPath, http.FileServer(root))

	if hasNoTrailingSlash(uriPath) {
		server.Get(uriPath, http.RedirectHandler(uriPath+"/", 301).ServeHTTP)
		uriPath += "/"
	}
	uriPath += "*"

	server.Get(uriPath, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}

// serve a single static file via a chi router
//
// parameters:
// - relativeFilesystemPath where to find the file(s) to serve relative to the current working directory
//   example: "docs"
//   note: do NOT add a trailing slash
// - uriPath under which sub-url they should be served
//   example: "/"
// - filename which exact file to serve. This will be added to the route, so only exactly this file is made available
//   example: "swagger.json"
func AddStaticSingleFileRoute(server chi.Router, relativeFilesystemPath string, uriPath string, filename string) {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, relativeFilesystemPath)
	root := http.Dir(filesDir)
	fs := http.FileServer(root)

	if hasNoTrailingSlash(uriPath) {
		uriPath += "/"
	}

	server.Get(uriPath+filename, func(w http.ResponseWriter, r *http.Request) {
		// this stops browsers from caching our swagger.json
		w.Header().Set(headers.CacheControl, "no-cache")
		fs.ServeHTTP(w, r)
	})
}

func hasNoTrailingSlash(path string) bool {
	return path != "/" && path[len(path)-1] != '/'
}
