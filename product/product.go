package product

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/wilsonelectronics/productsapi/cache"
	"github.com/wilsonelectronics/productsapi/data"

	"github.com/piotrkowalczuk/ntypes"
)

// Product . . .
type Product struct {
	GUID            string                  `json:"guid"`
	SKU             string                  `json:"sku"`
	ProductTypeID   int                     `json:"productTypeId"`
	ProductType     string                  `json:"productType"`
	UPC             string                  `json:"upc"`
	Details         *details                `json:"details"`
	Tags            []*productTag           `json:"tags"`
	Kits            []*kit                  `json:"kits"`
	Medias          []*media                `json:"media"`
	Notes           []*note                 `json:"notes"`
	Specifications  []*productSpecification `json:"specifications"`
	Vendors         []*productVendor        `json:"vendors"`
	RelatedProducts []*relatedProduct       `json:"relatedProducts"`
}

type details struct {
	Description      ntypes.String `json:"description"`
	DescriptionShort ntypes.String `json:"descriptionShort"`
	Title            string        `json:"title"`
	TitleTag         ntypes.String `json:"titleTag"`
	BodyHTML         ntypes.String `json:"body_HTML"`
	Price            float64       `json:"price"`
	ImageURL         string        `json:"imageURL"`
	Handle           string        `json:"handle"`
	ModifiedTime     string        `json:"modifiedTime"`
	IsActive         bool          `json:"isActive"`
	IsDeleted        bool          `json:"isDeleted"`
}

type kit struct {
	GUID           string        `json:"guid"`
	ProductGUID    string        `json:"productGuid"`
	KitItemName    string        `json:"kitItemName"`
	KitItemLinkURL ntypes.String `json:"kitItemLinkURL"`
	KitItemIconURL string        `json:"kitItemIconURL"`
	ItemOrder      int           `json:"itemOrder"`
	SKU            string        `json:"sku"`
}

type media struct {
	GUID         string        `json:"guid"`
	ProductGUID  string        `json:"productGuid"`
	MediaTypeID  int           `json:"mediaTypeId"`
	MediaTitle   ntypes.String `json:"mediaTitle"`
	MediaLinkURL ntypes.String `json:"mediaLinkURL"`
	MediaLogoURL ntypes.String `json:"mediaLogoURL"`
	MediaOrder   int           `json:"mediaOrder"`
	IsActive     bool          `json:"isActive"`
}

type note struct {
	GUID             string        `json:"guid"`
	ProductGUID      string        `json:"productGuid"`
	NoteTypeID       int           `json:"noteTypeId"`
	NoteText         ntypes.String `json:"noteText"`
	NoteOrder        int           `json:"noteOrder"`
	NoteTitle        ntypes.String `json:"noteTitle"`
	NoteIconImageURL ntypes.String `json:"noteIconImageUrl"`
}

type productSpecification struct {
	GUID               string `json:"guid"`
	ProductGUID        string `json:"productGuid"`
	SpecificationID    int    `json:"specificationId"`
	FieldValue         string `json:"specificationValue"`
	IsActive           bool   `json:"isActive"`
	SpecificationLabel string `json:"specificationLabel"`
}

type productVendor struct {
	GUID             string `json:"guid"`
	ProductGUID      string `json:"productGuid"`
	VendorID         int    `json:"vendorId"`
	VendorName       string `json:"vendorName"`
	VendorImageURL   string `json:"vendorImageURL"`
	ProductVendorURL string `json:"productVendorURL"`
}

type relatedProduct struct {
	GUID        string `json:"guid"`
	ProductGUID string `json:"productGuid"`
	SKU         string `json:"sku"`
	ImageURL    string `json:"imageURL"`
	Handle      string `json:"handle"`
}

type productTag struct {
	GUID        string `json:"guid"`
	ProductGUID string `json:"productGuid"`
	TagID       int    `json:"tagID"`
	Tag         string `json:"tag"`
	IsActive    bool   `json:"isActive"`
}

type chanResult struct {
	Result interface{}
	Error  error
}

func mergeChans(cs ...<-chan *chanResult) <-chan *chanResult {
	out := make(chan *chanResult)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan *chanResult) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// GetByHandle . . .
func GetByHandle(handle string) (*Product, error) {
	bytes, err := cache.Retrieve(handle)
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return getFromDbAndCache(handle)
	}

	product := &Product{}
	err = json.Unmarshal(bytes, product)
	return product, err
}

func getFromDbAndCache(handle string) (*Product, error) {
	db, err := data.GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("set nocount on; exec [spcProductGet] ?", handle)
	product := &Product{Details: &details{}}
	if err = row.Scan(
		&product.GUID,
		&product.SKU,
		&product.ProductTypeID,
		&product.ProductType,
		&product.UPC,
		&product.Details.Description,
		&product.Details.DescriptionShort,
		&product.Details.Title,
		&product.Details.TitleTag,
		&product.Details.BodyHTML,
		&product.Details.Price,
		&product.Details.ImageURL,
		&product.Details.Handle,
		&product.Details.ModifiedTime,
		&product.Details.IsActive,
		&product.Details.IsDeleted); err != nil {
		return nil, fmt.Errorf("spcProductGet Query Scan failed: %s", err)
	}

	kitChan := make(chan *chanResult)
	vendorChan := make(chan *chanResult)
	relatedProductsChan := make(chan *chanResult)
	specificationsChan := make(chan *chanResult)
	mediasChan := make(chan *chanResult)
	notesChan := make(chan *chanResult)
	tagsChan := make(chan *chanResult)

	go getKits(product.GUID, kitChan)
	go getVendors(product.GUID, vendorChan)
	go getRelatedProducts(product.GUID, relatedProductsChan)
	go getSpecifications(product.GUID, specificationsChan)
	go getMedias(product.GUID, mediasChan)
	go getNotes(product.GUID, notesChan)
	go getTags(product.GUID, tagsChan)

	for ch := range mergeChans(kitChan, vendorChan, relatedProductsChan, specificationsChan, mediasChan, notesChan, tagsChan) {
		if ch.Error != nil {
			return nil, ch.Error
		}

		if k, ok := ch.Result.(*kit); ok {
			product.Kits = append(product.Kits, k)
		} else if v, ok := ch.Result.(*productVendor); ok {
			product.Vendors = append(product.Vendors, v)
		} else if rp, ok := ch.Result.(*relatedProduct); ok {
			product.RelatedProducts = append(product.RelatedProducts, rp)
		} else if s, ok := ch.Result.(*productSpecification); ok {
			product.Specifications = append(product.Specifications, s)
		} else if m, ok := ch.Result.(*media); ok {
			product.Medias = append(product.Medias, m)
		} else if n, ok := ch.Result.(*note); ok {
			product.Notes = append(product.Notes, n)
		} else if t, ok := ch.Result.(*productTag); ok {
			product.Tags = append(product.Tags, t)
		}
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}
	cache.Store(handle, productJSON)

	return product, nil
}

func getKits(id string, ch chan *chanResult) {
	defer close(ch)

	db, err := data.GetDB()
	if err != nil {
		ch <- &chanResult{Error: err}
		return
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductKitGet] ?", id)
	if err != nil {
		ch <- &chanResult{Error: fmt.Errorf("spcProductKitGet Query failed: %s", err)}
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &kit{}

		if err = rows.Scan(
			&r.GUID,
			&r.ProductGUID,
			&r.KitItemName,
			&r.KitItemLinkURL,
			&r.KitItemIconURL,
			&r.ItemOrder,
			&r.SKU); err != nil {
			ch <- &chanResult{Error: fmt.Errorf("spcProductKitGet Query Scan failed: %s", err)}
		}
		ch <- &chanResult{Result: r}
	}
}

func getVendors(id string, ch chan *chanResult) {
	defer close(ch)

	db, err := data.GetDB()
	if err != nil {
		ch <- &chanResult{Error: err}
		return
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductVendorGet] ?", id)
	if err != nil {
		ch <- &chanResult{Error: fmt.Errorf("spcProductVendorGet Query failed: %s", err)}
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &productVendor{}

		if err = rows.Scan(
			&r.GUID,
			&r.ProductGUID,
			&r.VendorID,
			&r.VendorName,
			&r.VendorImageURL,
			&r.ProductVendorURL,
		); err != nil {
			ch <- &chanResult{Error: fmt.Errorf("spcProductVendorGet Query Scan failed: %s", err)}
		}
		ch <- &chanResult{Result: r}
	}
}

func getRelatedProducts(id string, ch chan *chanResult) {
	defer close(ch)

	db, err := data.GetDB()
	if err != nil {
		ch <- &chanResult{Error: err}
		return
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductRelatedGet] ?", id)
	if err != nil {
		ch <- &chanResult{Error: fmt.Errorf("spcProductRelatedGet Query failed: %s", err)}
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &relatedProduct{}

		if err = rows.Scan(
			&r.GUID,
			&r.ProductGUID,
			&r.SKU,
			&r.ImageURL,
			&r.Handle,
		); err != nil {
			ch <- &chanResult{Error: fmt.Errorf("spcProductRelatedGet Query Scan failed: %s", err)}
		}
		ch <- &chanResult{Result: r}
	}
}

func getSpecifications(id string, ch chan *chanResult) {
	defer close(ch)

	db, err := data.GetDB()
	if err != nil {
		ch <- &chanResult{Error: err}
		return
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductSpecificationsGet] ?", id)
	if err != nil {
		ch <- &chanResult{Error: fmt.Errorf("getProductSpecifications Query failed: %s", err)}
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &productSpecification{}

		if err = rows.Scan(
			&r.GUID,
			&r.ProductGUID,
			&r.SpecificationID,
			&r.FieldValue,
			&r.IsActive,
			&r.SpecificationLabel); err != nil {
			ch <- &chanResult{Error: fmt.Errorf("getProductSpecifications Query Scan failed: %s", err)}
		}
		ch <- &chanResult{Result: r}
	}
}

func getMedias(id string, ch chan *chanResult) {
	defer close(ch)

	db, err := data.GetDB()
	if err != nil {
		ch <- &chanResult{Error: err}
		return
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductMediaGet] ?", id)
	if err != nil {
		ch <- &chanResult{Error: fmt.Errorf("getProductMedia Query failed: %s", err)}
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &media{}

		if err = rows.Scan(
			&r.GUID,
			&r.ProductGUID,
			&r.MediaTypeID,
			&r.MediaTitle,
			&r.MediaLinkURL,
			&r.MediaLogoURL,
			&r.MediaOrder,
			&r.IsActive); err != nil {
			ch <- &chanResult{Error: fmt.Errorf("getProductMedia Query Scan failed: %s", err)}
		}
		ch <- &chanResult{Result: r}
	}
}

func getNotes(id string, ch chan *chanResult) {
	defer close(ch)

	db, err := data.GetDB()
	if err != nil {
		ch <- &chanResult{Error: err}
		return
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductNotesGet] ?", id)
	if err != nil {
		ch <- &chanResult{Error: fmt.Errorf("getProductNotes Query failed: %s", err)}
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &note{}

		if err = rows.Scan(
			&r.GUID,
			&r.ProductGUID,
			&r.NoteTypeID,
			&r.NoteText,
			&r.NoteOrder,
			&r.NoteTitle,
			&r.NoteIconImageURL); err != nil {
			ch <- &chanResult{Error: fmt.Errorf("getProductNotes Query Scan failed: %s", err)}
		}
		ch <- &chanResult{Result: r}
	}
}

func getTags(id string, ch chan *chanResult) {
	defer close(ch)

	db, err := data.GetDB()
	if err != nil {
		ch <- &chanResult{Error: err}
		return
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductTagsGet] ?", id)
	if err != nil {
		ch <- &chanResult{Error: fmt.Errorf("spcProductTagsGet Query failed: %s", err)}
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &productTag{}

		if err = rows.Scan(
			&r.GUID,
			&r.ProductGUID,
			&r.TagID,
			&r.Tag,
			&r.IsActive); err != nil {
			ch <- &chanResult{Error: fmt.Errorf("spcProductTagsGet Query Scan failed: %s", err)}
		}
		ch <- &chanResult{Result: r}
	}
}
