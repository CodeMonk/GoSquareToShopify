package squarespace

// squarespace handles Unmarshaling of squarespace json, and storing of the products.  It can then be
// used to output any format.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// SquareSpaceProducts Holds our product Table
type SquareSpaceProducts struct {
	raw     []byte                // Our raw data
	Results []*SquareSpaceProduct `json:"results"`
	HasPrev bool                  `json:"hasPrevPage"`
	HasNext bool                  `json:"hasNextPage"`
}

// String implements our Stringer interface
func (ss *SquareSpaceProducts) String() string {
	return fmt.Sprintf("SquareSpaceProducts(%d bytes, %d items)", len(ss.raw), len(ss.Results))
}

// GoString implements our GoStringer interface (for %#v printing)
func (ss *SquareSpaceProducts) GoString() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("SquareSpaceProducts(%d bytes, %d items)\n", len(ss.raw), len(ss.Results)))
	out.WriteString(fmt.Sprintf("  Prev/Next: %v / %v\n", ss.HasPrev, ss.HasNext))
	for item := 0; item < len(ss.Results); item++ {
		out.WriteString(fmt.Sprintf("  %4d : %#v\n", item, ss.Results[item]))
	}

	return out.String()
}

// SquareSpaceProduct holds a single product, with all images and variants
// Tags and categories are combined into a comma separated list of tags
type SquareSpaceProduct struct {
	Type              ProductType     `json:"productType"`
	Id                string          `json:"id"`
	WebsiteId         string          `json:"websiteId"`
	URL               *ProductURL     `json:"url"`
	Visibility        *Visibility     `json:"visibility"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Images            []*Image        `json:"images"`
	AdditionalInfo    *AdditionalInfo `json:"additionalInfo"`
	Featured          bool            `json:"featuredProduct"`
	Tags              []string        `json:"tags"`
	Categories        []string        `json:"categories"`
	VariantAttributes []string        `json:"variantAttributeNames"`
	Variants          []*Variants     `json:"variants"`
}

// String implements our Stringer interface
func (ssp *SquareSpaceProduct) String() string {
	return fmt.Sprintf("SquareSpaceProduct(%s): %s", ssp.Id, ssp.Name)
}

// GoString implements our GoStringer interface (for %#v printing)
func (ssp *SquareSpaceProduct) GoString() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("SquareSpaceProduct(%s): %v\n", ssp.Id, ssp.Name))
	out.WriteString(fmt.Sprintf("                  Type : %v\n", ssp.Type))
	out.WriteString(fmt.Sprintf("           Description : %v\n", ssp.Description))
	out.WriteString("                Images :\n")
	for image := range ssp.Images {
		out.WriteString(fmt.Sprintf("                   %#v\n", ssp.Images[image]))
	}

	return out.String()
}

// ProductType did not parse correctly, since it could be a string or an integer,
// so we've defined it as a type so it can be unmarshalled individually, below
type ProductType string

// ProductURL holds our product info.  The ProductPath is used as the handle in shopify
type ProductURL struct {
	Path           string `json:"fullPath"`
	ProductPath    string `json:"productPath"`
	CollectionPath string `json:"collectionPath"`
}

// Visibility is parsed, but is not currently used
type Visibility struct {
	State     string     `json:"state"`
	VisibleOn *time.Time `json:"visibleOn"`
}

// Image holds an image information.  Currently, only the URL is used.
type Image struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// String implements the Stringer interface for Images
func (i *Image) String() string {
	return fmt.Sprintf("Image(%s-%s): %s", i.Id, i.Type, i.URL)
}

// GoString implements the GoStringer interface for Images
func (i *Image) GoString() string {
	return fmt.Sprintf("Image(%s-%s): %s", i.Id, i.Type, i.URL)
}

// Variants holds our different models of the product, like different sizes or colors
// The attributes are left as a mapping of attribute name to interface{}
type Variants struct {
	SKU        string            `json:"sku"`
	Price      *Price            `json:"price"`
	Stock      *Stock            `json:"stock"`
	Attributes map[string]string `json:"attributes"`
}

// Stock contains the number of items we have
type Stock struct {
	Unlimited bool `json:"unlimited"`
	Quantity  int  `json:"quantity"`
}

// Price is the variant price
type Price struct {
	Str string `json:"decimalValue"`
}

type AdditionalInfo map[string]interface{}

// UnmarshalJSON will handle strings and floats as strings
func (pt *ProductType) UnmarshalJSON(b []byte) error {
	var unmarshalled interface{}
	if err := json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf("Unable to unmarshall %s -> interface: %v", string(b),
			err)
	}

	// Set our value
	switch t := unmarshalled.(type) {
	case string:
		*pt = ProductType(t)
	case float64:
		*pt = ProductType(fmt.Sprintf("%v", t))
	default:
		return fmt.Errorf("Unexpected type for productType: %T (%v)", t, t)
	}

	return nil

}

// New allocates our SquareSpaceProducts, or returns error.  It will unmarshal from the file Handle
// provided, and, will return any errors from the read, or from the unmarshal.
// New will also close the file provided!
func New(in *os.File) (*SquareSpaceProducts, error) {

	defer in.Close()
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file: %v", err)
	}

	ss := &SquareSpaceProducts{
		raw: data,
	}
	// And, unmarshall us
	err = json.Unmarshal(data, &ss)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal file: %v", err)
	}

	return ss, nil
}
