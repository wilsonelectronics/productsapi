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
