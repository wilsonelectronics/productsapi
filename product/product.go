package product

import (
	"encoding/json"
	"fmt"
	"products-api/cache"
	"products-api/data"

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
	GUID        string `json:"guid"`
	ProductGUID string `json:"productGuid"`
	NoteTypeID  int    `json:"noteTypeId"`
	NoteText    string `json:"noteText"`
	NoteOrder   int    `json:"noteOrder"`
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
	if db == nil || err != nil {
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

	if product.Kits, err = getKits(product.GUID); err != nil {
		return nil, err
	}
	if product.Vendors, err = getVendors(product.GUID); err != nil {
		return nil, err
	}
	if product.RelatedProducts, err = getRelatedProducts(product.GUID); err != nil {
		return nil, err
	}
	if product.Specifications, err = getSpecifications(product.GUID); err != nil {
		return nil, err
	}
	if product.Medias, err = getMedias(product.GUID); err != nil {
		return nil, err
	}
	if product.Notes, err = getNotes(product.GUID); err != nil {
		return nil, err
	}
	if product.Tags, err = getTags(product.GUID); err != nil {
		return nil, err
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}
	cache.Store(handle, productJSON)

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
		return nil, fmt.Errorf("spcProductKitGet Query failed: %s", err)
	}
	defer rows.Close()

	kits := []*kit{}

	for rows.Next() {
		k := &kit{}

		if err = rows.Scan(
			&k.GUID,
			&k.ProductGUID,
			&k.KitItemName,
			&k.KitItemLinkURL,
			&k.KitItemIconURL,
			&k.ItemOrder,
			&k.SKU); err != nil {
			return nil, fmt.Errorf("spcProductKitGet Query Scan failed: %s", err)
		}
		kits = append(kits, k)
	}
	return kits, nil
}

func getVendors(id string) ([]*productVendor, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductVendorGet] ?", id)
	if err != nil {
		return nil, fmt.Errorf("spcProductVendorGet Query failed: %s", err)
	}
	defer rows.Close()

	vendors := []*productVendor{}

	for rows.Next() {
		v := &productVendor{}

		if err = rows.Scan(
			&v.GUID,
			&v.ProductGUID,
			&v.VendorID,
			&v.VendorName,
			&v.VendorImageURL,
			&v.ProductVendorURL,
		); err != nil {
			return nil, fmt.Errorf("spcProductVendorGet Query Scan failed: %s", err)
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
		return nil, fmt.Errorf("spcProductRelatedGet Query failed: %s", err)
	}
	defer rows.Close()

	related := []*relatedProduct{}

	for rows.Next() {
		r := &relatedProduct{}

		if err = rows.Scan(
			&r.GUID,
			&r.ProductGUID,
			&r.SKU,
			&r.ImageURL,
			&r.Handle,
		); err != nil {
			return nil, fmt.Errorf("spcProductRelatedGet Query Scan failed: %s", err)
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
			&s.GUID,
			&s.ProductGUID,
			&s.SpecificationID,
			&s.FieldValue,
			&s.IsActive,
			&s.SpecificationLabel); err != nil {
			return nil, fmt.Errorf("getProductSpecifications Query Scan failed: %s", err)
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
			&m.GUID,
			&m.ProductGUID,
			&m.MediaTypeID,
			&m.MediaTitle,
			&m.MediaLinkURL,
			&m.MediaLogoURL,
			&m.MediaOrder,
			&m.IsActive); err != nil {
			return nil, fmt.Errorf("getProductMedia Query Scan failed: %s", err)
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
		return nil, fmt.Errorf("getProductNotes Query failed: %s", err)
	}
	defer rows.Close()

	notes := []*note{}

	for rows.Next() {
		n := &note{}

		if err = rows.Scan(
			&n.GUID,
			&n.ProductGUID,
			&n.NoteTypeID,
			&n.NoteText,
			&n.NoteOrder,
		); err != nil {
			return nil, fmt.Errorf("getProductNotes Query Scan failed: %s", err)
		}
		notes = append(notes, n)
	}
	return notes, err
}

func getTags(id string) ([]*productTag, error) {
	db, err := data.GetDB()
	if db == nil || err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductTagsGet] ?", id)
	if err != nil {
		return nil, fmt.Errorf("spcProductTagsGet Query failed: %s", err)
	}
	defer rows.Close()

	tags := []*productTag{}

	for rows.Next() {
		t := &productTag{}

		if err = rows.Scan(
			&t.GUID,
			&t.ProductGUID,
			&t.TagID,
			&t.Tag,
			&t.IsActive); err != nil {
			return nil, fmt.Errorf("spcProductTagsGet Query Scan failed: %s", err)
		}
		tags = append(tags, t)
	}

	return tags, err
}
