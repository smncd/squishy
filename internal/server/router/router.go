package router

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"gitlab.com/smncd/squishy/internal/filesystem"
)

// The router context, shared between all methods
type RouterContext struct {
	// Squishyfile instance
	S *filesystem.SquishyFile
}

type HandlerFunc = func(w http.ResponseWriter, r *http.Request, ctx RouterContext)

// Router is an HTTP request router that uses http.ServeMux, as well as custom context.
type Router struct {
	ctx    RouterContext
	mux    *http.ServeMux
	routes map[string]http.Handler
}

func New(ctx RouterContext) (*Router, error) {
	validate := validator.New()
	err := validate.Struct(ctx)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	rt := &Router{
		ctx:    ctx,
		mux:    http.NewServeMux(),
		routes: make(map[string]http.Handler),
	}

	return rt, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

	http.NotFound(w, req)
}

// Handle allows registering a route with an http.Handler
func (r *Router) Handle(method, path string, handler http.Handler) {
	key := fmt.Sprintf("%s %s", method, path)
	r.routes[key] = handler
}

// HandleFunc allows registering a route with a function
func (r *Router) HandleFunc(method, path string, handlerFunc HandlerFunc) {
	r.Handle(method, path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlerFunc(w, req, r.ctx)
	}))
}

// Registers a new GET request handle with the given path.
func (rt *Router) GET(path string, handler HandlerFunc) {
	rt.HandleFunc("GET", path, handler)
}

// Registers a new POST request handle with the given path.
func (rt *Router) POST(path string, handler HandlerFunc) {
	rt.HandleFunc("POST", path, handler)
}

// Registers a new PUT request handle with the given path.
func (rt *Router) PUT(path string, handler HandlerFunc) {
	rt.HandleFunc("PUT", path, handler)
}

// Registers a new DELETE request handle with the given path.
func (rt *Router) DELETE(path string, handler HandlerFunc) {
	rt.HandleFunc("DELETE", path, handler)
}

func (rt *Router) StaticFS(path string, fsys fs.FS) {
	path = ensureTrailingSlash(path)

	handler := func(w http.ResponseWriter, r *http.Request, ctx RouterContext) {
		filePath := strings.TrimPrefix(r.URL.Path, path)

		_, err := fs.Stat(fsys, filePath)
		if err != nil {
			return
		}

		http.StripPrefix(path, http.FileServer(http.FS(fsys))).ServeHTTP(w, r)
	}

	rt.HandleFunc("GET", path, handler)
}
