package shopify

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

var (
	HEADER_FIELD  = 0
	IMAGE_FIELD   = 24
	ShopifyFields = []string{
		"Handle",
		"Title",
		"Body (HTML)",
		"Vendor",
		"Type",
		"Tags",
		"Published",
		"Option1 Name",
		"Option1 Value",
		"Option2 Name",
		"Option2 Value",
		"Option3 Name",
		"Option3 Value",
		"Variant SKU",
		"Variant Grams",
		"Variant Inventory Tracker",
		"Variant Inventory Qty",
		"Variant Inventory Policy",
		"Variant Fulfillment Service",
		"Variant Price",
		"Variant Compare At Price",
		"Variant Requires Shipping",
		"Variant Taxable",
		"Variant Barcode",
		"Image Src",
		"Image Alt Text",
		"Gift Card",
		"Google Shopping / MPN",
		"Google Shopping / Age Group",
		"Google Shopping / Gender",
		"Google Shopping / Google Product Category",
		"SEO Title",
		"SEO Description",
		"Google Shopping / AdWords Grouping",
		"Google Shopping / AdWords Labels",
		"Google Shopping / Condition",
		"Google Shopping / Custom Product",
		"Google Shopping / Custom Label 0",
		"Google Shopping / Custom Label 1",
		"Google Shopping / Custom Label 2",
		"Google Shopping / Custom Label 3",
		"Google Shopping / Custom Label 4",
		"Variant Image",
		"Variant Weight Unit",
	}
)

type Shopify struct {
	Writer *csv.Writer
}

func New(out *os.File) *Shopify {
	return &Shopify{
		Writer: csv.NewWriter(out),
	}
}

func (s *Shopify) GetRow(mappings map[string]string) ([]string, error) {
	result := make([]string, len(ShopifyFields))

	for index, field := range ShopifyFields {
		if value, ok := mappings[field]; ok {
			result[index] = value
		} else {
			result[index] = ""
		}
	}

	return result, nil
}

func (s *Shopify) WriteCSV(rows [][]string) error {
	err := s.WriteHeader()
	if err != nil {
		return fmt.Errorf("Unable to write header: %v", err)
	}

	// Write All rows
	for _, row := range rows {
		err = s.WriteRow(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Shopify) WriteHeader() error {
	// Write header
	return s.Writer.Write(ShopifyFields)
}
func (s *Shopify) WriteRow(row []string) (err error) {
	// Write row -- be sure and split images back out.
	images := strings.Split(row[IMAGE_FIELD], "|")

	first_row := true
	image_row := make([]string, len(ShopifyFields))
	image_row[HEADER_FIELD] = row[HEADER_FIELD]

	for _, image := range images {
		if first_row {
			// Write out row with all the data
			first_row = false
			row[IMAGE_FIELD] = image
			err = s.Writer.Write(row)
			if err != nil {
				return
			}
		} else {
			// for subsequent rows, write a blank row except for images and header
			image_row[IMAGE_FIELD] = image
			err = s.Writer.Write(image_row)
			if err != nil {
				return
			}
		}
	}
	return
}

func (s *Shopify) CloseAll() {
	s.Writer.Flush()
}
