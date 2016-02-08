package shopify

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/CodeMonk/GoSquaretoShopify/squarespace"
)

var (
	HANDLE_FIELD  = 0
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

func fieldToIndex(k string) (int, error) {
	for i, f := range ShopifyFields {
		if k == f {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Field %s not found!\n", k)
}

func indexToField(i int) (string, error) {
	if i < 0 || i > len(ShopifyFields)-1 {
		return "", fmt.Errorf("Field must be between [0 and %d]",
			len(ShopifyFields)-1)
	}
	return ShopifyFields[i], nil
}

// Return an empty CSV row
func makeRow() []string {
	return make([]string, len(ShopifyFields))
}

// setField will set the named field in the csv rows
func setField(row []string, k, v string) error {
	i, err := fieldToIndex(k)
	if err != nil {
		return err
	}

	// And, set it
	row[i] = v
	return nil
}

type Shopify struct {
	Writer *csv.Writer
}

func New(out *os.File) *Shopify {
	return &Shopify{
		Writer: csv.NewWriter(out),
	}
}

// Write out our header field
func (s *Shopify) WriteHeader() error {
	// Write header
	return s.Writer.Write(ShopifyFields)
}

// WriteRow writes the supplied row to the csv
func (s *Shopify) WriteRow(row []string) error {
	// Write out row with all the data
	return s.Writer.Write(row)
}

func (s *Shopify) CloseAll() {
	s.Writer.Flush()
}

// WriteSquareSpaceProduct will write a single product to the csv.  Note, this could cause
// multiple records to be added.
func (s *Shopify) WriteSquareSpaceProduct(item *squarespace.SquareSpaceProduct) (err error) {

	// So, let's do our first record, with our first image, and our first variant
	row := makeRow()

	err = setField(row, "Handle", item.URL.ProductPath)
	if err != nil {
		return
	}
	err = setField(row, "Title", item.Name)
	if err != nil {
		return
	}
	err = setField(row, "Body (HTML)", item.Description)
	if err != nil {
		return
	}

	// Get our categories and tags, and add them all as tags
	err = setField(row, "Tags", strings.Join(append(item.Categories, item.Tags...), ","))
	if err != nil {
		return
	}

	// Now, for each image/variant, add the fields
	extraRecords := len(item.Images)
	if len(item.Variants) > extraRecords {
		extraRecords = len(item.Variants)
	}

	for i := 0; i < extraRecords; i++ {
		if i != 0 {
			// Our first record, we write the whole thing.  For subsequent records, only
			// use the handle
			row = makeRow()
			err = setField(row, "Handle", item.URL.ProductPath)
			if err != nil {
				return
			}
		}

		// Now, include image (i) and variant (i), if they exist
		if i < len(item.Images)-1 {
			err = setField(row, "Image Src", item.Images[i].URL)
			if err != nil {
				return
			}
		}

		if i < len(item.Variants)-1 {
			// First, the hard coded ones
			err = setField(row, "Variant Inventory Tracker", "shopify")
			if err != nil {
				return
			}

			err = setField(row, "Variant Inventory Policy", "deny")
			if err != nil {
				return
			}
			err = setField(row, "Variant Fulfillment Service", "manual")
			if err != nil {
				return
			}
			// Now the important ones
			err = setField(row, "Variant Price", item.Variants[i].Price.Str)
			if err != nil {
				return
			}
			err = setField(row, "Variant SKU", item.Variants[i].SKU)
			if err != nil {
				return
			}
			err = setField(row, "Variant Inventory Qty", fmt.Sprintf("%v",
				item.Variants[i].Stock.Quantity))
			if err != nil {
				return
			}
			// Finally, add the reason for our Variant
			index := 0
			for k, v := range item.Variants[i].Attributes {
				index++
				if index > 3 {
					return errors.New("Error received too many variants - can only handle three options")
				}
				key_field := fmt.Sprintf("Option%d Name", index)
				value_field := fmt.Sprintf("Option%d Value", index)
				err = setField(row, key_field, k)
				if err != nil {
					return
				}
				err = setField(row, value_field, v)
				if err != nil {
					return
				}
			}
		}

		// And, write our row
		err = s.WriteRow(row)
		if err != nil {
			return
		}
	}

	return
}

// WriteCSV writes the provided rows out to our csv file
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
