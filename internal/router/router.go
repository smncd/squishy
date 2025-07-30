package router

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

type HandlerFunc[C any] func(http.ResponseWriter, *http.Request, C)

// Router is an HTTP request router that uses http.ServeMux, as well as custom context.
type Router[C any] struct {
	ctx           C
	mux           *http.ServeMux
	routes        map[string]http.Handler
	fallbackRoute http.Handler
}

func New[C any](ctx C) *Router[C] {
	return &Router[C]{
		ctx:           ctx,
		mux:           http.NewServeMux(),
		routes:        make(map[string]http.Handler),
		fallbackRoute: http.NotFoundHandler(),
	}
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

	r.fallbackRoute.ServeHTTP(w, req)
}

// Handle allows registering a route with an http.Handler
func (r *Router[C]) Handle(method, path string, handler http.Handler) {
	key := fmt.Sprintf("%s %s", method, path)
	r.routes[key] = handler
}

// HandleFunc allows registering a route with a function
func (r *Router[C]) HandleFunc(method, path string, handlerFunc HandlerFunc[C]) {
	r.Handle(method, path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlerFunc(w, req, r.ctx)
	}))
}

// Registers a new GET request handle with the given path.
func (r *Router[C]) GET(path string, handler HandlerFunc[C]) {
	r.HandleFunc("GET", path, handler)
}

// Registers a new POST request handle with the given path.
func (r *Router[C]) POST(path string, handler HandlerFunc[C]) {
	r.HandleFunc("POST", path, handler)
}

// Registers a new PUT request handle with the given path.
func (r *Router[C]) PUT(path string, handler HandlerFunc[C]) {
	r.HandleFunc("PUT", path, handler)
}

// Registers a new DELETE request handle with the given path.
func (r *Router[C]) DELETE(path string, handler HandlerFunc[C]) {
	r.HandleFunc("DELETE", path, handler)
}

func (r *Router[C]) Fallback(handlerFunc HandlerFunc[C]) {
	r.fallbackRoute = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlerFunc(w, req, r.ctx)
	})
}

func (r *Router[C]) StaticFS(path string, fsys fs.FS) {
	path = ensureTrailingSlash(path)

	handler := func(w http.ResponseWriter, req *http.Request, ctx C) {
		filePath := strings.TrimPrefix(req.URL.Path, path)

		_, err := fs.Stat(fsys, filePath)
		if err != nil {
			r.fallbackRoute.ServeHTTP(w, req)
			return
		}

		http.StripPrefix(path, http.FileServer(http.FS(fsys))).ServeHTTP(w, req)
	}

	r.HandleFunc("GET", path, handler)
}
