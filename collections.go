package main

import (
	"fmt"
	"log"
	"products-api/data"

	"github.com/piotrkowalczuk/ntypes"
)

// Meta . . .
type Meta struct {
	CategoryGUID     string        `json:"collectionGuid,omitempty"`
	CategoryName     string        `json:"collectionName,omitempty"`
	CategoryHandle   string        `json:"collectionHandle,omitempty"`
	Description      string        `json:"collectionDescription,omitempty"`
	CategoryImageURL ntypes.String `json:"collectionImage,omitempty"`
	*Products        `json:"products,omitempty"`
}

// ProductData . . .
type ProductData struct {
	CategoryID       string
	ProductID        string
	SKU              string
	Title            string
	DescriptionShort string
	Price            string
	ImageURL         string
	Handle           string
	ShopifyID        string
}

// Products . . .
type Products struct {
	CategoryID string    `json:"collectionID,omitempty"`
	Product    []Product `json:"products,omitempty"`
}

// Product . . .
type Product struct {
	CategoryID string `json:"collectionGuid,omitempty"`
	Handle     string `json:"handle,omitempty"`
	Name       string `json:"name,omitempty"`
	SKU        string `json:"sku,omitempty"`
	Images     images
	HTML       html
}

type images struct {
	Main mainImages
}

type mainImages struct {
	Large  string `json:"large,omitempty"`
	Medium string `json:"medium,omitempty"`
	Small  string `json:"small,omitempty"`
}

type html struct {
	DescriptionShort string `json:"descriptionShort,omitempty"`
}

// GetCollections gets all collection from azure database and passes them to redis
func GetCollections() []Meta {

	db, err := data.GetDB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcCategoryGet]")
	if err != nil {
		log.Println("Collection query failed: ", err)
		return nil
	}
	defer rows.Close()

	collection := []Meta{}

	for rows.Next() {
		c := Meta{}

		if err = rows.Scan(
			&c.CategoryGUID,
			&c.CategoryName,
			&c.CategoryHandle,
			&c.Description,
			&c.CategoryImageURL,
		); err != nil {
			log.Println("Query failed 2: ", err)
			return nil
		}
		collection = append(collection, c)
	}
	return collection
}

// GetCollectionProducts gets all products for a given collection and passes them to redis
func GetCollectionProducts() []Products {

	db, err := data.GetDB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductCategoryGet]")
	if err != nil {
		log.Println("CollectionProducts query failed: ", err)
		return nil
	}
	defer rows.Close()

	collectionProducts := []ProductData{}

	for rows.Next() {
		p := ProductData{}

		if err = rows.Scan(
			&p.CategoryID,
			&p.ProductID,
			&p.SKU,
			&p.Title,
			&p.DescriptionShort,
			&p.Price,
			&p.ImageURL,
			&p.Handle,
			&p.ShopifyID,
		); err != nil {
			return nil
		}
		collectionProducts = append(collectionProducts, p)
	}

	var products []Products
	for _, p := range collectionProducts {
		cIndex := collectionIndex(products, p.CategoryID)
		if cIndex == -1 {
			newProdCollection := Products{CategoryID: p.CategoryID}
			products = append(products, newProdCollection)
			cIndex = len(products) - 1
		}

		pIndex := getProductIndex(products[cIndex].Product, p.SKU)
		if pIndex == -1 {
			newProduct := Product{CategoryID: p.CategoryID, Handle: p.Handle, Name: p.Title, SKU: p.SKU}
			products[pIndex].Product = append(products[pIndex].Product, newProduct)
			pIndex = len(products) - 1
		}
	}

	fmt.Println(products)
	return products
}

func collectionIndex(arr []Products, CategoryID string) int {
	for k, v := range arr {
		if v.CategoryID == CategoryID {
			return k
		}
	}
	return -1
}

func getProductIndex(arr []Product, SKU string) int {
	for k, v := range arr {
		if v.SKU == SKU {
			return k
		}
	}
	return -1
}
