package models

type User struct {
	Username string `json:"username"`
}

type Bundle struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	BundledItems []BundledItem `json:"bundled_items"`
}

type BundlesResponse struct {
	Bundles []Bundle `json:"bundles"`
}

type BundledItem struct {
	ItemName string `json:"item_type"`
	ID       int    `json:"id"`
	ItemID   int    `json:"item_id"`
}

// BundledItemPattern describes pattern as an item of the bundle
type BundledItemPattern struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// BundledItemResponse describes response of BundleItem
type BundledItemResponse struct {
	BundledItem BundledItem          `json:"bundled_item"`
	Item        []BundledItemPattern `json:"item"`
}

type BundleContentResponse struct {
	Bundle Bundle `json:"bundle"`
}

// PatternPack describes pack for pattern
type PatternPack struct {
	YarnName   string     `json:"yarn_name"`
	Yarn       Yarn       `json:"yarn"`
	YarnWeight YarnWeight `json:"yarn_weight"`
}

// PatternNeedleSize describes pattern needle size
type PatternNeedleSize struct {
	Name string `json:"name"`
}

type PatternAttribute struct {
	Permalink string `json:"permalink"`
}

type Yarn struct {
	Permalink string `json:"permalink"`
}

type YarnWeight struct {
	Name string `json:"name"`
}

type PatternAuthor struct {
	Name      string `json:"name"`
	Permalink string `json:"permalink"`
}

type PatternResponse struct {
	Pattern Pattern `json:"pattern"`
}

// Pattern describes pattern
type Pattern struct {
	Name                  string              `json:"name"`
	GaugeDescription      string              `json:"gauge_description"`
	Permalink             string              `json:"permalink"`
	PatternAuthor         PatternAuthor       `json:"pattern_author"`
	PatternNeedleSizes    []PatternNeedleSize `json:"pattern_needle_sizes"`
	SizesAvailable        string              `json:"sizes_available"`
	PatternAttributes     []PatternAttribute  `json:"pattern_attributes"`
	Packs                 []PatternPack       `json:"packs"`
	YarnWeightDescription string              `json:"yarn_weight_description"`
}
