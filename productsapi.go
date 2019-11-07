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

type serverError struct {
	Message string
	Status  int
}

// Handler . . .
func Handler(opts *HandlerOptions) http.Handler {
	if opts == nil {
		panic("Invalid HandlerOptions provided")
	}
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					switch t := err.(type) {
					case string:
						http.Error(w, t, http.StatusInternalServerError)
					case error:
						http.Error(w, t.Error(), http.StatusInternalServerError)
					case serverError:
						http.Error(w, t.Message, t.Status)
					default:
						http.Error(w, "unknown error", http.StatusInternalServerError)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	})
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
