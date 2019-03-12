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

// Products . . .
type Products struct {
	Products []*Product `json:"product,omitempty"`
}

// Product . . .
type Product struct {
	Handle string `json:"handle,omitempty"`
	Name   string `json:"name,omitempty"`
	SKU    string `json:"sku,omitempty"`
	// Images image
	HTML html
}

// type image struct {
// 	Main main
// }

// type main struct {
// 	ImageURL string `json:"imageURL,omitempty"`
// }

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
