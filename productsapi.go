package productsapi

import (
	"net/http"
	"productsapi/controller"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// HandleFuncModel . . .
type HandleFuncModel struct {
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

// CORSOptions . . .
func CORSOptions(moreOptions ...handlers.CORSOption) []handlers.CORSOption {
	return append([]handlers.CORSOption{
		handlers.AllowedMethods([]string{"GET"}),
		handlers.AllowedHeaders([]string{"Content-Type", "*"}),
		handlers.AllowedOrigins([]string{"http://localhost:3000"})},
		moreOptions...)
}

// Router . . .
func Router(moreRoutes ...HandleFuncModel) *mux.Router {
	r := mux.NewRouter()
	for _, v := range append([]HandleFuncModel{
		HandleFuncModel{path: "/tags", handler: controller.GetTags},
		HandleFuncModel{path: "/tag/products/{tagId}", handler: controller.GetTagProducts},
		HandleFuncModel{path: "/categories", handler: controller.GetCategories},
		HandleFuncModel{path: "/category/products/{categoryGuid}", handler: controller.GetCategoryProducts},
		HandleFuncModel{path: "/product/{handle}", handler: controller.GetProduct}},
		moreRoutes...) {
		r.HandleFunc(v.path, v.handler)
	}
	return r
}
