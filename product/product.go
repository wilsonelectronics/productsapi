package product

import (
	"encoding/json"
	"fmt"
	"log"
	"products-api/cache"
	"products-api/data"

	"github.com/piotrkowalczuk/ntypes"
)

// Product . . .
type Product struct {
	ProductGUID     string                  `json:"productGuid"`
	SKU             string                  `json:"sku"`
	ProductTypeID   int                     `json:"productTypeId"`
	ProductType     string                  `json:"productType"`
	UPC             string                  `json:"upc,omitempty"`
	Details         *details                `json:"details"`
	Tags            []*productTag           `json:"tags"`
	Kits            []*kit                  `json:"kits"`
	Medias          []*media                `json:"media"`
	Notes           []*note                 `json:"notes"`
	Specifications  []*productSpecification `json:"specifications"`
	Vendors         []*productVendor        `json:"vendors"`
	RelatedProducts []*relatedProduct       `json:"relatedProducts"`
}

// details . . .
type details struct {
	Description      ntypes.String `json:"description"`
	DescriptionShort ntypes.String `json:"descriptionShort,omitempty"`
	Title            string        `json:"title"`
	TitleTag         ntypes.String `json:"titleTag,omitempty"`
	BodyHTML         ntypes.String `json:"body_HTML,omitempty"`
	Price            float64       `json:"price"`
	ImageURL         string        `json:"imageURL"`
	Handle           string        `json:"handle"`
	ModifiedTime     string        `json:"modifiedTime"`
	IsActive         bool          `json:"isActive"`
	IsDeleted        bool          `json:"isDeleted"`
}

// kit . . .
type kit struct {
	ProductSiteKitGUID string        `json:"productSiteKitId"`
	ProductGUID        string        `json:"productGuid,omitempty"`
	KitItemName        string        `json:"kitItemName,omitempty"`
	KitItemLinkURL     ntypes.String `json:"kitItemLinkURL,omitempty"`
	KitItemIconURL     string        `json:"kitItemIconURL,omitempty"`
	ItemOrder          int           `json:"ItemOrder,omitempty"`
}

// media . . .
type media struct {
	ProductSiteMediaGUID string        `json:"productSiteMediaId"`
	ProductGUID          string        `json:"productGuid"`
	MediaTypeID          int           `json:"mediaTypeId,omitempty"`
	MediaTitle           ntypes.String `json:"mediaTitle,omitempty"`
	MediaLinkURL         ntypes.String `json:"mediaLinkURL,omitempty"`
	MediaLogoURL         ntypes.String `json:"mediaLogoURL,omitempty"`
	MediaOrder           int           `json:"mediaOrder,omitempty"`
	IsActive             bool          `json:"isActive,omitempty"`
}

// note . . .
type note struct {
	ProductSiteNoteGUID string `json:"productSiteNoteId,omitempty"`
	ProductGUID         string `json:"productGuid,omitempty"`
	NoteTypeID          int    `json:"noteTypeId,omitempty"`
	NoteText            string `json:"noteText,omitempty"`
	NoteOrder           int    `json:"noteOrder,omitempty"`
}

// productSpecification . . .
type productSpecification struct {
	ProductSiteSpecificationsGUID string `json:"productSiteSpecificationsId,omitempty"`
	ProductGUID                   string `json:"productGuid,omitempty"`
	SpecificationID               int    `json:"specificationId,omitempty"`
	SpecificationLabel            string `json:"specificationLabel,omitempty"`
	FieldValue                    string `json:"specificationValue,omitempty"`
	IsActive                      bool   `json:"isActive,omitempty"`
}

// productVendor . . .
type productVendor struct {
	ProductSiteVendorGUID string `json:"productSiteVendorId,omitempty"`
	ProductGUID           string `json:"productGuid,omitempty"`
	VendorID              int    `json:"vendorId,omitempty"`
	VendorName            string `json:"vendorName,omitempty"`
	VendorImageURL        string `json:"vendorImageURL,omitempty"`
	ProductVendorURL      string `json:"productVendorURL,omitempty"`
}

// relatedProduct . . .
type relatedProduct struct {
	ProductSiteRelatedGUID string `json:"productSiteRelatedId,omitempty"`
	ProductGUID            string `json:"productGuid,omitempty"`
	SKU                    string `json:"sku,omitempty"`
	ImageURL               string `json:"imageURL,omitempty"`
	OrderID                int    `json:"orderId,omitempty"`
}

type productTag struct {
	ProductSiteTagGUID string `json:"productSiteTagId"`
	ProductGUID        string `json:"productGuid"`
	TagID              int    `json:"tagID"`
	Tag                string `json:"tag"`
}

// GetByID . . .
func GetByID(id string) (*Product, error) {
	bytes, err := cache.Retrieve(id)
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return getFromDbAndCache(id)
	}

	product := &Product{}
	err = json.Unmarshal(bytes, product)
	return product, err
}

func getFromDbAndCache(id string) (*Product, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("set nocount on; exec [spcProductGet] ?", id)
	product := &Product{Details: &details{}}
	if err = row.Scan(
		&product.ProductGUID,
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
		&product.Details.IsDeleted,
	); err != nil {
		return nil, err
	}

	if product.Kits, err = getKits(id); err != nil {
		return nil, err
	}
	if product.Vendors, err = getVendors(id); err != nil {
		return nil, err
	}
	if product.RelatedProducts, err = getRelatedProducts(id); err != nil {
		return nil, err
	}
	if product.Specifications, err = getSpecifications(id); err != nil {
		return nil, err
	}
	if product.Medias, err = getMedias(id); err != nil {
		return nil, err
	}
	if product.Notes, err = getNotes(id); err != nil {
		return nil, err
	}
	if product.Tags, err = getTags(id); err != nil {
		return nil, err
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}
	cache.Store(id, productJSON)

	return product, nil
}

func getKits(id string) ([]*kit, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductKitGet] ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kits := []*kit{}

	for rows.Next() {
		k := &kit{}

		if err = rows.Scan(
			&k.ProductSiteKitGUID,
			&k.ProductGUID,
			&k.KitItemName,
			&k.KitItemLinkURL,
			&k.KitItemIconURL,
			&k.ItemOrder,
		); err != nil {
			log.Println("Error in spcProductKitGet: ", err)
		}
		kits = append(kits, k)
	}
	return kits, err
}

func getVendors(id string) ([]*productVendor, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductVendorGet] ?", id)
	if err != nil {
		log.Println("Product Vendor Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	vendors := []*productVendor{}

	for rows.Next() {
		v := &productVendor{}

		if err = rows.Scan(
			&v.ProductSiteVendorGUID,
			&v.ProductGUID,
			&v.VendorID,
			&v.VendorName,
			&v.VendorImageURL,
			&v.ProductVendorURL,
		); err != nil {
			return nil, err
		}
		vendors = append(vendors, v)
	}
	return vendors, err
}

func getRelatedProducts(id string) ([]*relatedProduct, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductRelatedGet] ?", id)
	if err != nil {
		log.Println("getRelatedProduct Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	related := []*relatedProduct{}

	for rows.Next() {
		r := &relatedProduct{}

		if err = rows.Scan(
			&r.ProductSiteRelatedGUID,
			&r.ProductGUID,
			&r.SKU,
			&r.ImageURL,
			&r.OrderID,
		); err != nil {
			log.Println(err)
			return nil, err
		}
		related = append(related, r)
	}

	return related, err
}

func getSpecifications(id string) ([]*productSpecification, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductSpecificationsGet] ?", id)
	if err != nil {
		return nil, fmt.Errorf("getProductSpecifications Query failed: %s", err)
	}
	defer rows.Close()

	specs := []*productSpecification{}

	for rows.Next() {
		s := &productSpecification{}

		if err = rows.Scan(
			&s.ProductSiteSpecificationsGUID,
			&s.ProductGUID,
			&s.SpecificationID,
			&s.SpecificationLabel,
			&s.FieldValue,
			&s.IsActive,
		); err != nil {
			return nil, fmt.Errorf("getProductSpecifications 2 Query failed: %s", err)
		}
		specs = append(specs, s)
	}
	return specs, err
}

func getMedias(id string) ([]*media, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductMediaGet] ?", id)
	if err != nil {
		return nil, fmt.Errorf("getProductMedia Query failed: %s", err)
	}
	defer rows.Close()

	medias := []*media{}

	for rows.Next() {
		m := &media{}

		if err = rows.Scan(
			&m.ProductSiteMediaGUID,
			&m.ProductGUID,
			&m.MediaTypeID,
			&m.MediaTitle,
			&m.MediaLinkURL,
			&m.MediaLogoURL,
			&m.MediaOrder,
			&m.IsActive,
		); err != nil {
			log.Println("getProductMedia 2 Query failed: ", err)
			return nil, err
		}
		medias = append(medias, m)
	}
	return medias, err
}

func getNotes(id string) ([]*note, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductNotesGet] ?", id)
	if err != nil {
		log.Println("getProductNotes Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	notes := []*note{}

	for rows.Next() {
		n := &note{}

		if err = rows.Scan(
			&n.ProductSiteNoteGUID,
			&n.ProductGUID,
			&n.NoteTypeID,
			&n.NoteText,
			&n.NoteOrder,
		); err != nil {
			log.Println("getProductNotes 2 Query failed: ", err)
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, err
}

func getTags(id string) ([]*productTag, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		log.Println("error: ", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductTagsGet] ?", id)
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}
	defer rows.Close()

	tags := []*productTag{}

	for rows.Next() {
		t := &productTag{}

		if err = rows.Scan(
			&t.ProductSiteTagGUID,
			&t.ProductGUID,
			&t.TagID,
			&t.Tag,
		); err != nil {
			log.Println("Error: ", err)
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, err
}
