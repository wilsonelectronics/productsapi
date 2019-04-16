package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"products-api/products"
	"strings"
)

// GetProduct . . .
func GetProduct(w http.ResponseWriter, r *http.Request) {
	inputParams := strings.Split(r.URL.Path, "/")[2:]
	productGUID := inputParams[0]
	product, err := products.GetByID(productGUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(productJSON)
}
