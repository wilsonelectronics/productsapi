package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"products-api/controller"

	"github.com/gorilla/handlers"
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

func main() {
	r := NewRouter()

	methods := handlers.AllowedMethods([]string{"GET"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "*"})
	origins := handlers.AllowedOrigins([]string{"https://localhost:3000",
		"http://localhost:3000",
		"http://localhost:3001",
		"https://localhost:4000",
		"http://localhost:4000",
		"https://wilsonpro.ca",
		"https://wilsonpro.com",
		"https://www.weboost.com/",
		"https://staging.weboost.com/",
		"https://staging.wilsonpro.ca",
		"https://staging.wilsonpro.com/",
		"https://staging-wilsonpro-canada-api.herokuapp.com",
		"https://wilsonpro-canada-staging.herokuapp.com/",
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT was not defined.")
	} else {
		fmt.Println("Listening on port: ", port)
		log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(methods, origins, headers)(r)))
	}
}
