package crawler

// ContentType Defines an enum for different asset content types
type ContentType int

const (
	PAGE ContentType = 1 + iota
	STYLESHEETS
	SCRIPTS
	IMAGES
)

// Link defines directed link between to indicies that define an asset position
type Link struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Value  int `json:"value"`
}

// GenerateLink acts as a factory for Link
func GenerateLink(source int, target int) *Link {
	return &Link{
		Source: source,
		Target: target,
		Value:  1,
	}
}

// Asset models a crawled URL
type Asset struct {
	URL       string           `json:"url"`
	Type      ContentType      `json:"type"`
	Links     map[string]*Link `json:"-"`
	processed bool
	index     int
}

// GenerateAsset is a factory for Asset
func GenerateAsset(url string, contentType ContentType) *Asset {
	return &Asset{
		URL:   url,
		Type:  contentType,
		Links: make(map[string]*Link),
	}
}

// GetIndex is an accessor for Asset.link
func (a *Asset) GetIndex() int {
	return a.index
}

// AddLink links the Asset as the parent to the given Asset
// Returns the Link if created, else returns nil
func (a *Asset) AddLink(newAsset *Asset) *Link {
	link, ok := a.Links[newAsset.URL]

	if ok {
		link.Value++
		return nil
	}

	link = GenerateLink(a.index, newAsset.index)
	a.Links[newAsset.URL] = link

	return link
}
