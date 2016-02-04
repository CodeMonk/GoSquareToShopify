package squarespace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// squarespace handles Unmarshaling of squarespace json

type SquareSpace struct {
	raw     []byte                // Our raw data
	Results []*SquareSpaceProduct `json:"results"`
	HasPrev bool                  `json:"hasPrevPage"`
	HasNext bool                  `json:"hasNextPage"`
}

func (ss *SquareSpace) String() string {
	return fmt.Sprintf("SquareSpace(%d bytes, %d items)", len(ss.raw), len(ss.Results))
}

func (ss *SquareSpace) GoString() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("SquareSpace(%d bytes, %d items)\n", len(ss.raw), len(ss.Results)))
	out.WriteString(fmt.Sprintf("  Prev/Next: %v / %v\n", ss.HasPrev, ss.HasNext))
	for item := 0; item < len(ss.Results); item++ {
		out.WriteString(fmt.Sprintf("  %4d : %#v\n", item, ss.Results[item]))
	}

	return out.String()
}

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

func (ssp *SquareSpaceProduct) String() string {
	return fmt.Sprintf("SquareSpaceProduct(%s): %s", ssp.Id, ssp.Name)
}

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

type ProductType string

type ProductURL struct {
	Path           string `json:"fullPath"`
	ProductPath    string `json:"productPath"`
	CollectionPath string `json:"collectionPath"`
}

type Visibility struct {
	State     string     `json:"state"`
	VisibleOn *time.Time `json:"visibleOn"`
}
type Image struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func (i *Image) String() string {
	return fmt.Sprintf("Image(%s-%s): %s", i.Id, i.Type, i.URL)
}
func (i *Image) GoString() string {
	return fmt.Sprintf("Image(%s-%s): %s", i.Id, i.Type, i.URL)
}

type Variants struct {
	Price *Price `json:"price"`
}

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

func New(in *os.File) (*SquareSpace, error) {

	data, err := ioutil.ReadAll(in)
	defer in.Close()
	if err != nil {
		return nil, fmt.Errorf("Unable to read file: %v", err)
	}

	ss := &SquareSpace{
		raw: data,
	}
	// And, unmarshall us
	err = json.Unmarshal(data, &ss)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal file: %v", err)
	}

	return ss, nil
}

func (sqp *SquareSpaceProduct) GetMappings(shFields []string) (map[string]string, error) {
	mapping := make(map[string]string, len(shFields))
	var err error

	mapping["Handle"] = sqp.URL.ProductPath
	mapping["Title"] = sqp.Name
	mapping["Body (HTML)"] = sqp.Description
	if len(sqp.Images) > 0 {
		images := make([]string, len(sqp.Images))
		for i, image := range sqp.Images {
			images[i] = image.URL
		}
		mapping["Image Src"] = strings.Join(images, "|")
	}
	if len(sqp.Variants) > 0 {
		mapping["Variant Price"] = sqp.Variants[0].Price.Str
	}
	mapping["Variant Inventory Policy"] = "deny"
	mapping["Variant Inventory Qty"] = "1"
	mapping["Variant Fulfillment Service"] = "manual"

	return mapping, err
}
