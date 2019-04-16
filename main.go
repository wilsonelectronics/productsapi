package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/gorilla/handlers"

	_ "github.com/denisenkom/go-mssqldb"
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

	r.HandleFunc("/collections", Collections).Methods("GET")
	r.HandleFunc("/collections/{collectionGuid}", CollectionProducts)
	//r.HandleFunc("/product/{sku}", GetProduct)

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
		"https://wilsonpro-canada-staging.herokuapp.com/",
	})
	addr := setPort()

	log.Fatal(http.ListenAndServe(addr, handlers.CORS(methods, origins, headers)(r)))
}
