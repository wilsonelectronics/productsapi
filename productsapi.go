package productsapi

import (
	"fmt"
	"log"
	"net/http"

	"productsapi/controller"

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

func createRouter(routes ...route) *mux.Router {
	r := mux.NewRouter()
	for _, v := range routes {
		r.HandleFunc(v.Path, v.Handler)
	}
	return r
}

// Handler . . .
func Handler(opts ...*HandlerOptions) http.Handler {
	var opt *HandlerOptions
	if opts == nil {
		opt = BaseHandlerOptions()
	} else if len(opts) != 1 {
		panic(fmt.Sprintf("Invalid number of HandlerOptions provided: %d", len(opts)))
	} else {
		opt = opts[0]
	}
	return handlers.CORS(opt.corsOptions...)(createRouter(opt.routes...))
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

// BaseHandlerOptions . . .
func BaseHandlerOptions() *HandlerOptions {
	opts := &HandlerOptions{}
	AddRoute(opts, "/tags", controller.GetTags)
	AddRoute(opts, "/tag/products/{tagId}", controller.GetTagProducts)
	AddRoute(opts, "/categories", controller.GetCategories)
	AddRoute(opts, "/category/products/{categoryGuid}", controller.GetCategoryProducts)
	AddRoute(opts, "/product/{handle}", controller.GetProduct)

	AddCORSOption(opts, "METHODS", "GET")
	AddCORSOption(opts, "HEADERS", "Content-Type", "*")
	AddCORSOption(opts, "ORIGINS", "http://localhost:3000")
	return opts
}
