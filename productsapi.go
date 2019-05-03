package productsapi

import (
	"products-api/controller"

	"github.com/gorilla/mux"
)

// NewRouter . . .
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/tags", controller.GetTags)
	r.HandleFunc("/tag/products/{tagId}", controller.GetTagProducts)

	r.HandleFunc("/categories", controller.GetCategories)
	r.HandleFunc("/category/products/{categoryGuid}", controller.GetCategoryProducts)

	r.HandleFunc("/product/{handle}", controller.GetProduct)
	return r
}
