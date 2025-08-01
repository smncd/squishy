package router

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

type HandlerFunc[C any] func(http.ResponseWriter, *http.Request, C)

// Router is an HTTP request router that uses http.ServeMux, as well as custom context.
type Router[C any] struct {
	c       C
	mux     *http.ServeMux
	routes  map[string]http.Handler
	noRoute http.Handler
	logger  *log.Logger
}

func New[C any](c C) *Router[C] {
	r := &Router[C]{
		c:      c,
		mux:    http.NewServeMux(),
		routes: make(map[string]http.Handler),
		logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
	}

	r.noRoute = r.handlerLoggingWrapper(http.NotFoundHandler())

	return r
}

func (r *Router[C]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	key := fmt.Sprintf("%s %s", method, path)
	if handler, ok := r.routes[key]; ok {
		handler.ServeHTTP(w, req)
		return
	}

	for routeKey, handler := range r.routes {
		methodPart := method + " "
		if !strings.HasPrefix(routeKey, methodPart) {
			continue
		}

		routePath := strings.TrimPrefix(routeKey, methodPart)
		if strings.HasSuffix(routePath, "/") && strings.HasPrefix(path, routePath) {
			handler.ServeHTTP(w, req)
			return
		}
	}

	r.noRoute.ServeHTTP(w, req)
}

// Registers a new route with the specified HTTP method and path using an http.Handler.
func (r *Router[C]) Route(method, path string, handler http.Handler) {
	key := fmt.Sprintf("%s %s", method, path)
	r.routes[key] = r.handlerLoggingWrapper(handler)
}

// Registers a new route with the specified HTTP method and path using a custom HandlerFunc.
func (r *Router[C]) RouteFunc(method, path string, handlerFunc HandlerFunc[C]) {
	r.Route(method, path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlerFunc(w, req, r.c)
	}))
}

// Registers a new GET request handle with the given path.
func (r *Router[C]) GET(path string, handler HandlerFunc[C]) {
	r.RouteFunc("GET", path, handler)
}

// Registers a new POST request handle with the given path.
func (r *Router[C]) POST(path string, handler HandlerFunc[C]) {
	r.RouteFunc("POST", path, handler)
}

// Registers a new PUT request handle with the given path.
func (r *Router[C]) PUT(path string, handler HandlerFunc[C]) {
	r.RouteFunc("PUT", path, handler)
}

// Registers a new DELETE request handle with the given path.
func (r *Router[C]) DELETE(path string, handler HandlerFunc[C]) {
	r.RouteFunc("DELETE", path, handler)
}

// Used in case no other routes are available/applicable.
// Matches _any_ incoming requests.
// Defaults to `http.NotFoundHandler()`.
func (r *Router[C]) NoRoute(handlerFunc HandlerFunc[C]) {
	r.noRoute = r.handlerLoggingWrapper(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlerFunc(w, req, r.c)
	}))
}

// Serve files from a filesystem.
func (r *Router[C]) StaticFS(path string, fsys fs.FS) {
	path = ensureTrailingSlash(path)

	handler := func(w http.ResponseWriter, req *http.Request, c C) {
		filePath := strings.TrimPrefix(req.URL.Path, path)

		_, err := fs.Stat(fsys, filePath)
		if err != nil {
			r.noRoute.ServeHTTP(w, req)
			return
		}

		http.StripPrefix(path, http.FileServer(http.FS(fsys))).ServeHTTP(w, req)
	}

	r.RouteFunc("GET", path, handler)
}

// Wrap http handler in logging function, with custom response writer wrapper.
func (r *Router[C]) handlerLoggingWrapper(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		rww := NewResponseWriterWrapper(w)

		handler.ServeHTTP(rww, req)

		r.logger.Printf("| %s %s %s %v", req.Proto, req.Method, req.URL.Path, rww.statusCode)
	})
}
