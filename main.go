package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"products-api/controller"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func setPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "$PORT not set"
	}
	return ":" + port
}

func main() {
	fmt.Println("Backend API!!")

	r := mux.NewRouter()

	// r.HandleFunc("/collections", GetCachedCollections).Methods("GET")
	// r.HandleFunc("/collection/{collectionGuid}", GetCachedCollectionProducts).Methods("GET")
	r.HandleFunc("/product/{guid}", controller.GetSingleProduct)

	methods := handlers.AllowedMethods([]string{"GET"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "*"})
	origins := handlers.AllowedOrigins([]string{"https://localhost:3000",
		"http://localhost:3000",
		"https://localhost:4000",
		"http://localhost:4000",
		"https://wilsonpro.ca",
		"https://wilsonpro.com",
		"https://www.weboost.com/",
		"https://staging.weboost.com/",
		"https://staging.wilsonpro.ca",
		"https://staging.wilsonpro.com/",
		"https://staging-wilsonpro-canada-api.herokuapp.com",
	})
	addr := setPort()

	log.Fatal(http.ListenAndServe(addr, handlers.CORS(methods, origins, headers)(r)))
}
