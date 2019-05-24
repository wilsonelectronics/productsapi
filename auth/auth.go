package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wilsonelectronics/productsapi/cache"
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

// TokenData . . .
type TokenData struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
	Scope     string `json:"scope"`
	ExpiresAt string `json:"expiresAt"`
	CreateOn  string `json:"createOn"`
}

// GetTokenData is inteaded to return a token object for the frontend sites to use as a Bearer Token
func GetTokenData(handle string) (*TokenData, error) {
	bytes, err := cache.Retrieve(handle)
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, err
	}

	token := &TokenData{}
	err = json.Unmarshal(bytes, token)
	return token, err
}

//SetTokenData will set a new token if one need to be set.
func SetTokenData(handle string, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = cache.Store(handle, body)
	return err
}
