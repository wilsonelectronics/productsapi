package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"products-api/auth"
	"products-api/controller"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	_ "github.com/denisenkom/go-mssqldb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	// Set up jwt middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			// Verify 'aud' claim
			audience := os.Getenv("AUDIENCE")
			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(audience, false)
			if !checkAudience {
				return token, errors.New("invalid audience")
			}

			// Verify 'iss' claim
			iss := os.Getenv("DOMAIN")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := auth.GetPemCert(token)
			if err != nil {
				panic(err)
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	r := mux.NewRouter()

	r.HandleFunc("/tags", controller.GetTags)
	r.HandleFunc("/tag/products/{tagId}", controller.GetTagProducts)

	r.HandleFunc("/categories", controller.GetCategories)
	r.HandleFunc("/category/products/{categoryGuid}", controller.GetCategoryProducts)

	r.HandleFunc("/product/{handle}", controller.GetProduct)

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

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT was not defined.")
	} else {
		log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(methods, origins, headers)(r)))
	}
}
