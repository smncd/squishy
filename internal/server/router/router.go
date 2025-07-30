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

// Router is an HTTP request router that uses http.ServeMux, as well as custom context.
type Router struct {
	ctx RouterContext
	mux *http.ServeMux
}

type HandlerFunc = func(w http.ResponseWriter, r *http.Request, ctx RouterContext)

func New(ctx RouterContext) (*Router, error) {
	validate := validator.New()
	err := validate.Struct(ctx)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	rt := &Router{
		ctx: ctx,
		mux: http.NewServeMux(),
	}

	return rt, nil
}

// Returns the http Multiplexer instance
func (rt *Router) Mux() *http.ServeMux {
	return rt.mux
}

// handleFunc registers a new request handle with the given path and method.
func (rt *Router) handleFunc(method string, path string, handler HandlerFunc) {
	rt.mux.HandleFunc(fmt.Sprintf("%s %s", method, path), func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, rt.ctx)
	})
}

// Registers a new GET request handle with the given path.
func (rt *Router) GET(path string, handler HandlerFunc) {
	rt.handleFunc("GET", path, handler)
}

// Registers a new POST request handle with the given path.
func (rt *Router) POST(path string, handler HandlerFunc) {
	rt.handleFunc("POST", path, handler)
}

// Registers a new PUT request handle with the given path.
func (rt *Router) PUT(path string, handler HandlerFunc) {
	rt.handleFunc("PUT", path, handler)
}

// Registers a new DELETE request handle with the given path.
func (rt *Router) DELETE(path string, handler HandlerFunc) {
	rt.handleFunc("DELETE", path, handler)
}

func (rt *Router) StaticFS(path string, fsys fs.FS, errorHandler func(w http.ResponseWriter, r *http.Request)) {
	path = ensureTrailingSlash(path)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.TrimPrefix(r.URL.Path, path)

		_, err := fs.Stat(fsys, filePath)
		if err != nil {
			errorHandler(w, r)
			return
		}

		http.StripPrefix(path, http.FileServer(http.FS(fsys))).ServeHTTP(w, r)
	})

	rt.mux.Handle(fmt.Sprintf("GET %s", path), handler)
}
