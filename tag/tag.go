package tag

import (
	"encoding/json"
	"fmt"
	"products-api/cache"
	"products-api/data"

	"github.com/piotrkowalczuk/ntypes"
)

// Tag . . .
type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CreatedTime string `json:"createdTime"`
	IsActive    string `json:"isActive"`
}

// Product . . .
type Product struct {
	GUID             string        `json:"guid"`
	SKU              string        `json:"sku"`
	ProductTypeID    int           `json:"productTypeId"`
	UPC              ntypes.String `json:"upc"`
	Description      string        `json:"description"`
	DescriptionShort ntypes.String `json:"descriptionShort"`
	Title            string        `json:"title"`
	TitleTag         ntypes.String `json:"titleTag"`
	BodyHTML         ntypes.String `json:"bodyHtml"`
	Price            float64       `json:"price"`
	ImageURL         string        `json:"imageUrl"`
	Handle           string        `json:"handle"`
	ModifiedTime     string        `json:"modifiedTime"`
	IsActive         bool          `json:"isActive"`
	IsDeleted        bool          `json:"isDeleted"`
}

// GetAll . . .
func GetAll() ([]*Tag, error) {
	bytes, err := cache.Retrieve("tags")
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return getAllFromDbAndCache()
	}

	tags := []*Tag{}
	err = json.Unmarshal(bytes, &tags)
	return tags, err
}

// GetProductsByID . . .
func GetProductsByID(tagID string) ([]*Product, error) {
	bytes, err := cache.Retrieve(tagID)
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return getProductsFromDbAndCache(tagID)
	}

	products := []*Product{}
	err = json.Unmarshal(bytes, &products)
	return products, err
}

func getAllFromDbAndCache() ([]*Tag, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcTagsGet]")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []*Tag{}
	for rows.Next() {
		tag := &Tag{}
		if err = rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.CreatedTime,
			&tag.IsActive); err != nil {
			return nil, fmt.Errorf("Error in spcTagsGet: %s", err)
		}
		tags = append(tags, tag)
	}

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}
	cache.Store("tags", tagsJSON)

	return tags, nil
}

func getProductsFromDbAndCache(tagID string) ([]*Product, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcTagProductsGet] ?", tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*Product{}
	for rows.Next() {
		product := &Product{}
		if err = rows.Scan(
			&product.GUID,
			&product.SKU,
			&product.ProductTypeID,
			&product.UPC,
			&product.Description,
			&product.DescriptionShort,
			&product.Title,
			&product.TitleTag,
			&product.BodyHTML,
			&product.Price,
			&product.ImageURL,
			&product.Handle,
			&product.ModifiedTime,
			&product.IsActive,
			&product.IsDeleted); err != nil {
			return nil, fmt.Errorf("Error in spcTagProductsGet: %s", err)
		}
		products = append(products, product)
	}

	productsJSON, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}
	cache.Store(tagID, productsJSON)

	return products, nil
}

// // Meta . . .
// type Meta struct {
// 	tagGUID     string        `json:"collectionGuid,omitempty"`
// 	tagName     string        `json:"collectionName,omitempty"`
// 	tagHandle   string        `json:"collectionHandle,omitempty"`
// 	Description      string        `json:"collectionDescription,omitempty"`
// 	tagImageURL ntypes.String `json:"collectionImage,omitempty"`
// 	*Products        `json:"products,omitempty"`
// }

// // CollectionProductData . . .
// type CollectionProductData struct {
// 	tagID       string
// 	ProductID        string
// 	SKU              string
// 	Title            string
// 	DescriptionShort string
// 	Price            string
// 	ImageURL         string
// 	Handle           string
// 	ShopifyID        string
// }

// // Products . . .
// type Products struct {
// 	tagID string    `json:"collectionGuid,omitempty"`
// 	Products   []Product `json:"products,omitempty"`
// }

// // Product . . .
// type Product struct {
// 	tagID string `json:"collectionGuid,omitempty"`
// 	Handle     string `json:"handle,omitempty"`
// 	Name       string `json:"name,omitempty"`
// 	SKU        string `json:"sku,omitempty"`
// 	Images     []Images
// 	HTML       []HTML
// }

// // Images . . .
// type Images struct {
// 	Large  string `json:"large,omitempty"`
// 	Medium string `json:"medium,omitempty"`
// 	Small  string `json:"small,omitempty"`
// }

// // HTML . . .
// type HTML struct {
// 	DescriptionShort string `json:"descriptionShort,omitempty"`
// }

// // GetCollections gets all collection from azure database and passes them to redis
// func GetCollections() []Meta {

// 	db, err := data.GetDB()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer db.Close()

// 	rows, err := db.Query("set nocount on; exec [spctagGet]")
// 	if err != nil {
// 		log.Println("Collection query failed: ", err)
// 		return nil
// 	}
// 	defer rows.Close()

// 	collection := []Meta{}

// 	for rows.Next() {
// 		c := Meta{}

// 		if err = rows.Scan(
// 			&c.tagGUID,
// 			&c.tagName,
// 			&c.tagHandle,
// 			&c.Description,
// 			&c.tagImageURL,
// 		); err != nil {
// 			log.Println("Query failed 2: ", err)
// 			return nil
// 		}
// 		collection = append(collection, c)
// 	}
// 	return collection
// }

// // GetCollectionProducts gets all products for a given collection and passes them to redis
// func GetCollectionProducts() []Products {

// 	db, err := data.GetDB()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer db.Close()

// 	rows, err := db.Query("set nocount on; exec [spcProducttagGet]")
// 	if err != nil {
// 		log.Println("CollectionProducts query failed: ", err)
// 		return nil
// 	}
// 	defer rows.Close()

// 	collectionProducts := []CollectionProductData{}
// 	for rows.Next() {
// 		p := CollectionProductData{}

// 		if err = rows.Scan(
// 			&p.tagID,
// 			&p.ProductID,
// 			&p.SKU,
// 			&p.Title,
// 			&p.DescriptionShort,
// 			&p.Price,
// 			&p.ImageURL,
// 			&p.Handle,
// 			&p.ShopifyID,
// 		); err != nil {
// 			return nil
// 		}
// 		collectionProducts = append(collectionProducts, p)
// 	}

// 	var products []Products
// 	for _, p := range collectionProducts {
// 		cIndex := collectionIndex(products, p.tagID)
// 		if cIndex == -1 {
// 			newProdCollection := Products{tagID: p.tagID}
// 			products = append(products, newProdCollection)
// 			cIndex = len(products) - 1
// 		}

// 		pIndex := getProductIndex(products[cIndex].Products, p.SKU)
// 		if pIndex == -1 {
// 			newProduct := Product{tagID: p.tagID, Handle: p.Handle, Name: p.Title, SKU: p.SKU}
// 			products[cIndex].Products = append(products[cIndex].Products, newProduct)
// 			pIndex = len(products[cIndex].Products) - 1
// 		}

// 		if p.SKU == p.SKU {
// 			newImage := Images{Large: p.ImageURL}
// 			products[cIndex].Products[pIndex].Images = append(products[cIndex].Products[pIndex].Images, newImage)
// 			newHTML := HTML{DescriptionShort: p.DescriptionShort}
// 			products[cIndex].Products[pIndex].HTML = append(products[cIndex].Products[pIndex].HTML, newHTML)
// 		}
// 	}
// 	return products
// }

// func collectionIndex(arr []Products, tagID string) int {
// 	for k, v := range arr {
// 		if v.tagID == tagID {
// 			return k
// 		}
// 	}
// 	return -1
// }

// func getProductIndex(arr []Product, SKU string) int {
// 	for k, v := range arr {
// 		if v.SKU == SKU {
// 			return k
// 		}
// 	}
// 	return -1
// }
