package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"products-api/controller"

	"github.com/codegangsta/negroni"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	_ "github.com/denisenkom/go-mssqldb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// JSONWebKeys : Key information used for JWT validation
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("Domain") + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	x5c := jwks.Keys[0].X5c
	for k, v := range x5c {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}
	return cert, nil
}

func main() {

	// Set up jwt middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			// Verify 'aud' claim
			audience := os.Getenv("Audience")
			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(audience, false)
			if !checkAudience {
				return token, errors.New("invalid audience")
			}

			// Verify 'iss' claim
			iss := os.Getenv("Domain")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err)
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	r := mux.NewRouter()

	r.Handle("/tags", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { controller.GetTags(w, r) }))),
	)

	r.HandleFunc("/tag/products/{tagId}", controller.GetTagProducts)

	r.HandleFunc("/categories", controller.GetCategories)
	// r.HandleFunc("/category/products/{categoryGuid}", controller.GetCategoryProducts)

	r.Handle("/category/products/{categoryGuid}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { controller.GetCategoryProducts(w, r) }))),
	)

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
