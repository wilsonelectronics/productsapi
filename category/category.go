package category

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wilsonelectronics/productsapi/cache"
	"github.com/wilsonelectronics/productsapi/data"

	"github.com/piotrkowalczuk/ntypes"
)

// Category . . .
type Category struct {
	GUID        string        `json:"guid"`
	Name        string        `json:"name"`
	Handle      string        `json:"handle"`
	HeaderText  ntypes.String `json:"headerText"`
	Description string        `json:"description"`
	ImageURL    ntypes.String `json:"imageURL"`
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
	ImageURL         string        `json:"imageURL"`
	Handle           string        `json:"handle"`
	ModifiedTime     string        `json:"modifiedTime"`
	IsActive         bool          `json:"isActive"`
	IsDeleted        bool          `json:"isDeleted"`
	Tags             []*tag        `json:"tags"`
}

type productTag struct {
	ProductTagGUID string
	ProductGUID    string
	TagID          int
	TagName        string
	IsActive       bool
}

type tag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

// GetAll . . .
func GetAll() ([]*Category, error) {
	bytes, err := cache.Retrieve("categories")
	if err != nil {
		return nil, err
	}

	categories := []*Category{}
	if bytes == nil {
		categories, err = getAllFromDbAndCache()
		if err != nil {
			return nil, err
		}

		categoriesJSON, err := json.Marshal(categories)
		if err != nil {
			return nil, err
		}

		cache.Store("categories", categoriesJSON)
	} else {
		err = json.Unmarshal(bytes, &categories)
	}

	return categories, err
}

// GetProducts . . .
func GetProducts(categoryGUID string) ([]*Product, error) {
	bytes, err := cache.Retrieve(categoryGUID)
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	if bytes == nil {
		products, err = getProductsFromDb(categoryGUID)
		if err != nil {
			return nil, err
		}

		productsJSON, err := json.Marshal(products)
		if err != nil {
			return nil, err
		}

		cache.Store(categoryGUID, productsJSON)
	} else {
		err = json.Unmarshal(bytes, &products)
	}

	return products, err
}

func getAllFromDbAndCache() ([]*Category, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcCategoryGet]")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}
	for rows.Next() {
		category := &Category{}
		if err = rows.Scan(
			&category.GUID,
			&category.Name,
			&category.Handle,
			&category.HeaderText,
			&category.Description,
			&category.ImageURL); err != nil {
			return nil, fmt.Errorf("Error in spcCategoryGet: %s", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func getProductsFromDb(categoryGUID string) ([]*Product, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}

	rows, err := db.Query("set nocount on; exec [spcCategoryProductsGet] ?", categoryGUID)
	if err != nil {
		return nil, err
	}

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
			return nil, fmt.Errorf("Error in spcCategoryProductsGet Scan: %s", err)
		}
		products = append(products, product)
	}
	rows.Close()
	db.Close()

	err = getAndSetTagsForProducts(products)
	return products, err
}

func getAndSetTagsForProducts(products []*Product) error {
	productGUIDs := []string{}
	for _, p := range products {
		productGUIDs = append(productGUIDs, p.GUID)
	}

	db, err := data.GetDB()
	if db == nil || err != nil {
		return err
	}

	defer db.Close()
	rows, err := db.Query("set nocount on; exec [spcProductTagsGet] ?", strings.Join(productGUIDs, ","))
	if err != nil {
		return fmt.Errorf("Error in spcProductTagsGet: %s", err)
	}

	productTags := []*productTag{}
	for rows.Next() {
		productTag := &productTag{}
		if rows.Scan(
			&productTag.ProductTagGUID,
			&productTag.ProductGUID,
			&productTag.TagID,
			&productTag.TagName,
			&productTag.IsActive); err != nil {
			return fmt.Errorf("Error in spcProductTagsGet Scan: %s", err)
		}

		productTags = append(productTags, productTag)
	}

	for _, p := range products {
		for _, pt := range productTags {
			if pt.ProductGUID == p.GUID {
				p.Tags = append(p.Tags, &tag{
					ID:       pt.TagID,
					Name:     pt.TagName,
					IsActive: pt.IsActive})
			}
		}
	}

	return nil
}
