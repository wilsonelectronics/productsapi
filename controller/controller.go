package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"products-api/product"
	"products-api/tag"
	"strings"
)

// GetProduct . . .
func GetProduct(w http.ResponseWriter, r *http.Request) {
	inputParams := strings.Split(r.URL.Path, "/")[2:]
	productGUID := inputParams[0]

	product, err := product.GetByID(productGUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(productJSON)
}

// GetTags . . .
func GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := tag.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tagsJSON)
}

// GetTagProducts . . .
func GetTagProducts(w http.ResponseWriter, r *http.Request) {
	inputParams := strings.Split(r.URL.Path, "/")[3:]
	tagID := inputParams[0]

	products, err := tag.GetProductsByID(tagID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	productsJSON, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(productsJSON)
}
