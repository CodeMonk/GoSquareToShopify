package shopify

import (
	"encoding/csv"
	"fmt"
	"os"
)

var (
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
	Fields map[string]string
	Writer *csv.Writer
}

func New(out *os.File) *Shopify {
	fields := make(map[string]string, len(ShopifyFields))
	for _, field := range ShopifyFields {
		fields[field] = ""
	}

	return &Shopify{
		Fields: fields,
		Writer: csv.NewWriter(out),
	}
}

func (s *Shopify) Set(key string, value string) error {
	if _, ok := s.Fields[key]; ok {
		s.Fields[key] = value
		return nil
	}

	return fmt.Errorf("Unable to add \"%s\": key not in fields!\n", key)
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
	return s.Writer.WriteAll(rows)
}

func (s *Shopify) WriteHeader() error {
	// Write header
	return s.Writer.Write(ShopifyFields)
}
func (s *Shopify) WriteRow(row []string) error {
	// Write header
	return s.Writer.Write(row)
}
func (s *Shopify) CloseAll() {
	s.Writer.Flush()
}
