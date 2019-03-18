package main

import (
	"log"
	"products-api/data"
)

// Meta . . .
type Meta struct {
	CollectionGUID string `json:"collectionGuid"`
	CollectionName string `json:"collectionName"`
	*Products      `json:"products,omitempty"`
}

// ProductData . . .
type ProductData struct {
	CategoryID       string `json:"categoryId"`
	ProductID        string `json:"productId"`
	SKU              string `json:"sku"`
	Title            string `json:"title"`
	DescriptionShort string `json:"descriptionShort"`
	Price            string `json:"price"`
	ImageURL         string `json:"imageURL"`
	Handle           string `json:"handle"`
	ShopifyID        string `json:"shopifyId"`
}

// Products . . .
type Products struct {
	Products []*Product `json:"product,omitempty"`
}

// Product . . .
type Product struct {
	Handle string `json:"handle,omitempty"`
	Name   string `json:"name,omitempty"`
	SKU    string `json:"sku,omitempty"`
	Images images
	HTML   html
}

type images struct {
	Main mainImages
}

type mainImages struct {
	ImageURL string `json:"imageURL,omitempty"`
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
			&c.CollectionGUID,
			&c.CollectionName,
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

	rows, err := db.Query("set nocount on; exec [spcProductGet]")
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

	var product []Product
	for _, p := range collectionProducts {
		pIndex := getProductIndex(product, p.SKU)
		if pIndex == -1 {
			//newProduct := Product{SKU: p.SKU}
		}
	}

	return nil
}

func getProductIndex(arr []Product, SKU string) int {
	for k, v := range arr {
		if v.SKU == SKU {
			return k
		}
	}
	return -1
}
