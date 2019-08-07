package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"productsapi/blog"
	"strings"

	"github.com/wilsonelectronics/productsapi/auth"
	"github.com/wilsonelectronics/productsapi/cache"
	"github.com/wilsonelectronics/productsapi/category"
	"github.com/wilsonelectronics/productsapi/product"
	"github.com/wilsonelectronics/productsapi/tag"
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

// GetAccessToken . . .
func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	inputParams := strings.Split(r.URL.Path, "/")[1:]
	handle := inputParams[0]
	if r.Method == "POST" {
		auth.SetTokenData(handle, r)
	}
	token, err := auth.GetTokenData(handle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	accessToken, err := json.Marshal(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(accessToken)
}

// FlushRedisDB . . .
func FlushRedisDB(w http.ResponseWriter, r *http.Request) {
	cache.Flush()
}

// GetBlogPosts . . .
func GetBlogPosts(w http.ResponseWriter, r *http.Request) {
	client := blog.NewClient("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=", os.Getenv("hubSpotAPI"))

	response, err := client.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// GetPostWithSlug . . .
func GetPostWithSlug(w http.ResponseWriter, r *http.Request) {
	client := blog.NewClient("https://api.hubapi.com/content/api/v2/blog-posts?hapikey=", os.Getenv("hubSpotAPI"))

	slug := r.FormValue("slug")

	post, err := client.GetPost(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(post)
}

// GetPostTopics . . .
func GetPostTopics(w http.ResponseWriter, r *http.Request) {
	client := blog.NewClient("https://api.hubapi.com/blogs/v3/topics/search?hapikey=", os.Getenv("hubSpotAPI"))

	topics, err := client.GetTopics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(topics)
}
