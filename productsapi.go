package productsapi

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type route struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

// HandlerOptions . . .
type HandlerOptions struct {
	routes      []route
	corsOptions []handlers.CORSOption
}

// Handler . . .
func Handler(opts *HandlerOptions) http.Handler {
	if opts == nil {
		panic("Invalid HandlerOptions provided")
	}
	r := mux.NewRouter()
	for _, v := range opts.routes {
		r.HandleFunc(v.Path, v.Handler)
	}
	return handlers.CORS(opts.corsOptions...)(r)
}

// AddRoute . . .
func AddRoute(opts *HandlerOptions, path string, f func(http.ResponseWriter, *http.Request)) {
	opts.routes = append(opts.routes, route{Path: path, Handler: f})
}

// AddCORSOption . . .
func AddCORSOption(opts *HandlerOptions, corsKey string, corsValues ...string) {
	var f func([]string) handlers.CORSOption
	switch corsKey {
	case "METHODS":
		f = handlers.AllowedMethods
	case "HEADERS":
		f = handlers.AllowedHeaders
	case "ORIGINS":
		f = handlers.AllowedOrigins
	default:
		log.Fatal("Invalid CORS option type: " + corsKey + ". Must be METHODS, HEADERS, or ORIGINS")
	}
	opts.corsOptions = append(opts.corsOptions, f(corsValues))
}
