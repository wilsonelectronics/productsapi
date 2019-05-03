package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"productsapi/category"
	"productsapi/product"
	"productsapi/tag"
	"strings"
)

// GetProduct . . .
func GetProduct(w http.ResponseWriter, r *http.Request) {
	inputParams := strings.Split(r.URL.Path, "/")[2:]
	handle := inputParams[0]

	product, err := product.GetByHandle(handle)
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

// GetCategories . . .
func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := category.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	categoriesJSON, err := json.Marshal(categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(categoriesJSON)
}

// GetCategoryProducts . . .
func GetCategoryProducts(w http.ResponseWriter, r *http.Request) {
	inputParams := strings.Split(r.URL.Path, "/")[3:]
	categoryID := inputParams[0]

	products, err := category.GetProducts(categoryID)
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
