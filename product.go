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

func getProduct(id string, sku string) ([]Product, error) {
	db, err := data.GetDB()
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcProductGet] ?, ?", id, sku)
	if err != nil {
		return nil, err
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
			//&p.DescriptionTag,
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
	var product []Product
	for _, p := range productData {

		productIndex := getProductIndex(product, p.SKU)
		if productIndex == -1 {
			newProduct := Product{SKU: p.SKU, ProductTypeID: p.ProductTypeID, ProductType: p.ProductType, UPC: p.UPC, Handle: p.Handle}
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
	return product, err
}
