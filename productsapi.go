package productsapi

import (
	"productsapi/controller"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NewBaseCORSOptions . . .
func NewBaseCORSOptions() []handlers.CORSOption {
	return []handlers.CORSOption{
		handlers.AllowedMethods([]string{"GET"}),
		handlers.AllowedHeaders([]string{"Content-Type", "*"}),
		handlers.AllowedOrigins([]string{"http://localhost:3000"})}
}

// NewBaseRouter . . .
func NewBaseRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/tags", controller.GetTags)
	r.HandleFunc("/tag/products/{tagId}", controller.GetTagProducts)

	r.HandleFunc("/categories", controller.GetCategories)
	r.HandleFunc("/category/products/{categoryGuid}", controller.GetCategoryProducts)

	r.HandleFunc("/product/{handle}", controller.GetProduct)
	return r
}
