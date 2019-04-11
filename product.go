package main

import (
	"log"
	"products-api/data"

	"github.com/piotrkowalczuk/ntypes"
)

// ProductData : Model for [Portal].[dbo].[Product]
type ProductData struct {
	ProductSiteGUID  string
	SKU              string
	SiteID           int
	WebsiteName      string
	ProductTypeID    int
	ProductType      string
	UPC              string
	Description      string
	DescriptionShort ntypes.String
	// DescriptionTag   ntypes.String
	Title        string
	TitleTag     ntypes.String
	BodyHTML     ntypes.String
	Price        float64
	ImageURL     string
	Handle       string
	ModifiedTime string
	ModifiedBy   int
	IsActive     bool
	IsDeleted    bool
}

// SingleProduct . . .
type SingleProduct struct {
	ProductSiteGUID string `json:"productSiteId,omitempty"`
	SKU             string `json:"sku,omitempty"`
	ProductTypeID   int    `json:"productTypeId,omitempty"`
	ProductType     string `json:"productType,omitempty"`
	Handle          string `json:"handle,omitempty"`
	UPC             string `json:"upc,omitempty"`
	Sites           []Site `json:"sites,omitempty"`
}

// Site . . .
type Site struct {
	ID              int                     `json:"id,omitempty"`
	WebsiteName     string                  `json:"name,omitempty"`
	ProductSiteGUID string                  `json:"productSiteID,omitempty"`
	Details         []Details               `json:"details"`
	Tags            []productTags           `json:"tags"`
	Kits            []Kit                   `json:"kit"`
	Media           []Media                 `json:"media"`
	Notes           []Notes                 `json:"notes"`
	Specification   []ProductSpecifications `json:"specifications"`
	Vendors         []ProductVendors        `json:"vendors"`
	Related         []RelatedProducts       `json:"related"`
}

// Details . . .
type Details struct {
	ProductSiteGUID  string        `json:"productSiteId,omitempty"`
	Description      string        `json:"description,omitempty"`
	DescriptionShort ntypes.String `json:"descriptionShort,omitempty"`
	// DescriptionTag   ntypes.String `json:"descriptionTag,omitempty"`
	Title     string        `json:"title,omitempty"`
	TitleTag  ntypes.String `json:"titleTag,omitempty"`
	BodyHTML  ntypes.String `json:"body_HTML,omitempty"`
	Price     float64       `json:"price,omitempty"`
	ImageURL  string        `json:"imageURL,omitempty"`
	Handle    string        `json:"handle,omitempty"`
	IsActive  bool          `json:"isActive,omitempty"`
	IsDeleted bool          `json:"isDeleted,omitempty"`
}

// Kit . . .
type Kit struct {
	ProductSiteKitGUID string        `json:"productSiteKitId"`
	ProductSiteGUID    string        `json:"productSiteId,omitempty"`
	KitItemName        string        `json:"kitItemName,omitempty"`
	KitItemLinkURL     ntypes.String `json:"kitItemLinkURL,omitempty"`
	KitItemIconURL     string        `json:"kitItemIconURL,omitempty"`
	ItemOrder          int           `json:"ItemOrder,omitempty"`
}

// Media . . .
type Media struct {
	ProductSiteMediaGUID string        `json:"productSiteMediaId"`
	ProductSiteGUID      string        `json:"productSiteId"`
	MediaTypeID          int           `json:"mediaTypeId,omitempty"`
	MediaTitle           ntypes.String `json:"mediaTitle,omitempty"`
	MediaLinkURL         ntypes.String `json:"mediaLinkURL,omitempty"`
	MediaLogoURL         ntypes.String `json:"mediaLogoURL,omitempty"`
	MediaOrder           int           `json:"mediaOrder,omitempty"`
	IsActive             bool          `json:"isActive,omitempty"`
}

// Notes . . .
type Notes struct {
	ProductSiteNoteGUID string `json:"productSiteNoteId,omitempty"`
	ProductSiteGUID     string `json:"productSiteId,omitempty"`
	NoteTypeID          int    `json:"noteTypeId,omitempty"`
	NoteText            string `json:"noteText,omitempty"`
	NoteOrder           int    `json:"noteOrder,omitempty"`
}

// ProductSpecifications . . .
type ProductSpecifications struct {
	ProductSiteSpecificationsGUID string `json:"productSiteSpecificationsId,omitempty"`
	ProductSiteGUID               string `json:"productSiteId,omitempty"`
	SpecificationID               int    `json:"specificationId,omitempty"`
	SpecificationLabel            string `json:"specificationLabel,omitempty"`
	FieldValue                    string `json:"specificationValue,omitempty"`
	IsActive                      bool   `json:"isActive,omitempty"`
}

// ProductVendors . . .
type ProductVendors struct {
	ProductSiteVendorGUID string `json:"productSiteVendorId,omitempty"`
	ProductSiteGUID       string `json:"productSiteId,omitempty"`
	VendorID              int    `json:"vendorId,omitempty"`
	VendorName            string `json:"vendorName,omitempty"`
	VendorImageURL        string `json:"vendorImageURL,omitempty"`
	ProductVendorURL      string `json:"productVendorURL,omitempty"`
}

// RelatedProducts . . .
type RelatedProducts struct {
	ProductSiteRelatedGUID string `json:"productSiteRelatedId,omitempty"`
	ProductSiteGUID        string `json:"productSiteId,omitempty"`
	SKU                    string `json:"sku,omitempty"`
	ImageURL               string `json:"imageURL,omitempty"`
	OrderID                int    `json:"orderId,omitempty"`
}

type productTags struct {
	ProductSiteTagGUID string `json:"productSiteTagId"`
	ProductSiteGUID    string `json:"productSiteId"`
	TagID              int    `json:"tagID"`
	Tag                string `json:"tag"`
}

func getProduct(id string, sku string) []SingleProduct {
	db, err := data.GetDB()
	if err != nil {
		return nil
	}
	if db == nil {
		return nil
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductGet] ?, ?", id, sku)
	if err != nil {
		return nil
	}
	defer rows.Close()

	productData := []*ProductData{}

	for rows.Next() {
		p := &ProductData{}

		if err = rows.Scan(
			&p.ProductSiteGUID,
			&p.SKU,
			&p.SiteID,
			&p.WebsiteName,
			&p.ProductTypeID,
			&p.ProductType,
			&p.UPC,
			&p.Description,
			&p.DescriptionShort,
			&p.Title,
			&p.TitleTag,
			&p.BodyHTML,
			&p.Price,
			&p.ImageURL,
			&p.Handle,
			&p.ModifiedTime,
			&p.ModifiedBy,
			&p.IsActive,
			&p.IsDeleted,
		); err != nil {
			log.Println("Error in spcProductGet!", err)
		}
		productData = append(productData, p)
	}

	kit := func() []Kit {
		kit, err := getProductKit(id, sku)
		if err != nil {
			log.Println("Error: ", err)
			return nil
		}
		return kit
	}()

	vendor := func() []ProductVendors {
		vendor, err := getProductVendor(id, sku)
		if err != nil {
			log.Println("Error: ", err)
			return nil
		}
		return vendor
	}()

	related := func() []RelatedProducts {
		related, err := getRelatedProduct(id, sku)
		if err != nil {
			log.Println("Error: ", err)
			return nil
		}
		return related
	}()

	specs := func() []ProductSpecifications {
		spec, err := getProductSpecifications(id, sku)
		if err != nil {
			log.Println("Error: ", err)
			return nil
		}
		return spec
	}()

	media := func() []Media {
		media, err := getProductMedia(id, sku)
		if err != nil {
			log.Println("Error: ", err)
			return nil
		}
		return media
	}()

	notes := func() []Notes {
		note, err := getProductNotes(id, sku)
		if err != nil {
			log.Println("Error: ", err)
			return nil
		}
		return note
	}()

	tags := func() []productTags {
		tag, err := getProductTags(id, sku)
		if err != nil {
			log.Println("Error: ", err)
			return nil
		}
		return tag
	}()

	// time.Sleep(1 * time.Second)
	var product []SingleProduct
	for _, p := range productData {

		productIndex := getProdIndex(product, p.SKU)
		if productIndex == -1 {
			newProduct := SingleProduct{SKU: p.SKU, ProductTypeID: p.ProductTypeID, ProductType: p.ProductType, UPC: p.UPC, Handle: p.Handle}
			product = append(product, newProduct)
			productIndex = len(product) - 1
		}

		siteIndex := getSiteIndex(product[productIndex].Sites, p.SiteID)
		if siteIndex == -1 {
			newSite := Site{ID: p.SiteID, WebsiteName: p.WebsiteName, ProductSiteGUID: p.ProductSiteGUID}
			product[productIndex].Sites = append(product[productIndex].Sites, newSite)
			siteIndex = len(product[productIndex].Sites) - 1
		}

		if p.ProductSiteGUID == p.ProductSiteGUID {
			newDetail := Details{ProductSiteGUID: p.ProductSiteGUID, Description: p.Description, DescriptionShort: p.DescriptionShort, Title: p.Title, TitleTag: p.TitleTag, BodyHTML: p.BodyHTML, Price: p.Price, ImageURL: p.ImageURL, IsActive: p.IsActive, IsDeleted: p.IsDeleted}
			product[productIndex].Sites[siteIndex].Details = append(product[productIndex].Sites[siteIndex].Details, newDetail)
		}

		for _, t := range tags {
			if p.ProductSiteGUID == t.ProductSiteGUID {
				newTag := productTags{ProductSiteTagGUID: t.ProductSiteTagGUID, ProductSiteGUID: t.ProductSiteGUID, TagID: t.TagID, Tag: t.Tag}
				product[productIndex].Sites[siteIndex].Tags = append(product[productIndex].Sites[siteIndex].Tags, newTag)
			}
		}

		for _, k := range kit {
			if p.ProductSiteGUID == k.ProductSiteGUID {
				newKit := Kit{ProductSiteKitGUID: k.ProductSiteKitGUID, ProductSiteGUID: k.ProductSiteGUID, KitItemName: k.KitItemName, KitItemLinkURL: k.KitItemLinkURL, KitItemIconURL: k.KitItemIconURL, ItemOrder: k.ItemOrder}
				product[productIndex].Sites[siteIndex].Kits = append(product[productIndex].Sites[siteIndex].Kits, newKit)
			}
		}

		for _, v := range vendor {
			if p.ProductSiteGUID == v.ProductSiteGUID {
				newVendor := ProductVendors{ProductSiteVendorGUID: v.ProductSiteVendorGUID, ProductSiteGUID: v.ProductSiteGUID, VendorID: v.VendorID, VendorName: v.VendorName, VendorImageURL: v.VendorImageURL, ProductVendorURL: v.ProductVendorURL}
				product[productIndex].Sites[siteIndex].Vendors = append(product[productIndex].Sites[siteIndex].Vendors, newVendor)
			}
		}

		for _, r := range related {
			if p.ProductSiteGUID == r.ProductSiteGUID {
				newRelated := RelatedProducts{ProductSiteRelatedGUID: r.ProductSiteRelatedGUID, ProductSiteGUID: r.ProductSiteGUID, SKU: r.SKU, ImageURL: r.ImageURL, OrderID: r.OrderID}
				product[productIndex].Sites[siteIndex].Related = append(product[productIndex].Sites[siteIndex].Related, newRelated)
			}
		}

		for _, s := range specs {
			if p.ProductSiteGUID == s.ProductSiteGUID {
				newSpec := ProductSpecifications{ProductSiteSpecificationsGUID: s.ProductSiteSpecificationsGUID, ProductSiteGUID: s.ProductSiteGUID, SpecificationID: s.SpecificationID, SpecificationLabel: s.SpecificationLabel, FieldValue: s.FieldValue, IsActive: s.IsActive}
				product[productIndex].Sites[siteIndex].Specification = append(product[productIndex].Sites[siteIndex].Specification, newSpec)
			}
		}

		for _, m := range media {
			if p.ProductSiteGUID == m.ProductSiteGUID {
				newMedia := Media{ProductSiteMediaGUID: m.ProductSiteMediaGUID, ProductSiteGUID: m.ProductSiteGUID, MediaTypeID: m.MediaTypeID, MediaTitle: m.MediaTitle, MediaLinkURL: m.MediaLinkURL, MediaLogoURL: m.MediaLogoURL, MediaOrder: m.MediaOrder, IsActive: m.IsActive}
				product[productIndex].Sites[siteIndex].Media = append(product[productIndex].Sites[siteIndex].Media, newMedia)
			}
		}

		for _, n := range notes {
			if p.ProductSiteGUID == n.ProductSiteGUID {
				newNote := Notes{ProductSiteNoteGUID: n.ProductSiteNoteGUID, ProductSiteGUID: n.ProductSiteGUID, NoteTypeID: n.NoteTypeID, NoteText: n.NoteText, NoteOrder: n.NoteOrder}
				product[productIndex].Sites[siteIndex].Notes = append(product[productIndex].Sites[siteIndex].Notes, newNote)
			}
		}
	}
	return product
}

func getProductKit(id string, sku string) ([]Kit, error) {
	db, err := data.GetDB()
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductKitGet] ?, ?", id, sku)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kitData := []Kit{}

	for rows.Next() {
		k := Kit{}

		if err = rows.Scan(
			&k.ProductSiteKitGUID,
			&k.ProductSiteGUID,
			&k.KitItemName,
			&k.KitItemLinkURL,
			&k.KitItemIconURL,
			&k.ItemOrder,
		); err != nil {
			log.Println("Error in spcProductKitGet: ", err)
		}
		kitData = append(kitData, k)
	}
	return kitData, err
}

func getProductVendor(id string, sku string) ([]ProductVendors, error) {
	db, err := data.GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductVendorGet] ?, ?", id, sku)
	if err != nil {
		log.Println("Product Vendor Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	vendors := []ProductVendors{}

	for rows.Next() {
		v := ProductVendors{}

		if err = rows.Scan(
			&v.ProductSiteVendorGUID,
			&v.ProductSiteGUID,
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

func getRelatedProduct(id string, sku string) ([]RelatedProducts, error) {
	db, err := data.GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductRelatedGet] ?, ?", id, sku)
	if err != nil {
		log.Println("getRelatedProduct Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	related := []RelatedProducts{}

	for rows.Next() {
		r := RelatedProducts{}

		if err = rows.Scan(
			&r.ProductSiteRelatedGUID,
			&r.ProductSiteGUID,
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

func getProductSpecifications(id string, sku string) ([]ProductSpecifications, error) {
	db, err := data.GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductSpecificationsGet] ?, ?", id, sku)
	if err != nil {
		log.Println("getProductSpecifications Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	specs := []ProductSpecifications{}

	for rows.Next() {
		s := ProductSpecifications{}

		if err = rows.Scan(
			&s.ProductSiteSpecificationsGUID,
			&s.ProductSiteGUID,
			&s.SpecificationID,
			&s.SpecificationLabel,
			&s.FieldValue,
			&s.IsActive,
		); err != nil {
			log.Println("getProductSpecifications 2 Query failed: ", err)
			return nil, err
		}
		specs = append(specs, s)
	}
	return specs, err
}

func getProductMedia(id string, sku string) ([]Media, error) {
	db, err := data.GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductMediaGet] ?, ?", id, sku)
	if err != nil {
		log.Println("getProductMedia Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	media := []Media{}

	for rows.Next() {
		m := Media{}

		if err = rows.Scan(
			&m.ProductSiteMediaGUID,
			&m.ProductSiteGUID,
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
		media = append(media, m)
	}
	return media, err
}

func getProductNotes(id string, sku string) ([]Notes, error) {
	db, err := data.GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductNotesGet] ?, ?", id, sku)
	if err != nil {
		log.Println("getProductNotes Query failed: ", err)
		return nil, err
	}
	defer rows.Close()

	notes := []Notes{}

	for rows.Next() {
		n := Notes{}

		if err = rows.Scan(
			&n.ProductSiteNoteGUID,
			&n.ProductSiteGUID,
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

func getProductTags(id string, sku string) ([]productTags, error) {
	db, err := data.GetDB()
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductTagsGet] ?, ?", id, sku)
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}
	defer rows.Close()

	tags := []productTags{}

	for rows.Next() {
		t := productTags{}

		if err = rows.Scan(
			&t.ProductSiteTagGUID,
			&t.ProductSiteGUID,
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

func getSiteIndex(arr []Site, ID int) int {
	for k, v := range arr {
		if v.ID == ID {
			return k
		}
	}
	return -1
}

func getProdIndex(arr []SingleProduct, SKU string) int {
	for k, v := range arr {
		if v.SKU == SKU {
			return k
		}
	}
	return -1
}
