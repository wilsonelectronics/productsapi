package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Jwks : Slice of json web keys for validating JWT sent with requests
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

// GetAccessToken is inteaded to return a hash for the frontend sites to use as a Bearer Token
func GetAccessToken() string {

	url := "https://wilson-portal.auth0.com/oauth/token"

	payload := strings.NewReader("{\"client_id\":\"RqQGt7ncz2hP8SlsMcyghUnXrzXbfM81\",\"client_secret\":\"MhdbhSrjpK7Mqnz2O8freMsFyakusGME8fdQ9UULdDQq3Zo8b-vD1EQj9rfSelcN\",\"audience\":\"https://product.wilsonelectronics.com\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)

	return string(body)

}

// GetPemCert . . .
func GetPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("DOMAIN") + ".well-known/jwks.json")

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
