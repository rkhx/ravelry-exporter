package gapi

import (
	"fmt"
	"github.com/rkhx/ravelry-exporter/internal/models"
	"strings"
)

type RowData map[string]interface{}

type Link struct {
	URL    string
	Text   string
	Extras string
}

func NewRowData(p *models.Pattern) RowData {
	var yarnLinks []Link
	for _, pack := range p.Packs {
		yarnLinks = append(yarnLinks, Link{
			Text:   pack.YarnName,
			URL:    fmt.Sprintf("www.ravelry.com/yarns/library/%s", pack.Yarn.Permalink),
			Extras: pack.YarnWeight.Name,
		})
	}

	attributes := make([]string, len(p.PatternAttributes))
	for i, attribute := range p.PatternAttributes {
		attributes[i] = attribute.Permalink
	}

	sizesNames := make([]string, len(p.PatternNeedleSizes))
	for _, needleSize := range p.PatternNeedleSizes {
		sizesNames = append(sizesNames, needleSize.Name)
	}

	return RowData{
		"Pattern Name": Link{Text: p.Name, URL: fmt.Sprintf("https://www.ravelry.com/patterns/library/%s", p.Permalink)},
		"Designer":     Link{Text: p.PatternAuthor.Name, URL: fmt.Sprintf("https://www.ravelry.com/patterns/sources/%s", p.PatternAuthor.Permalink)},
		"Gauge, needle size": strings.Join([]string{
			strings.Trim(p.GaugeDescription, "\n"),
			strings.Join(sizesNames, "\n"),
			p.YarnWeightDescription,
		}, "\n\n"),
		"Sizes":            p.SizesAvailable,
		"Recommended yarn": yarnLinks,
		"Attributes":       strings.Join(attributes, ", "),
	}
}
