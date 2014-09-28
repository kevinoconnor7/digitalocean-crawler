package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/kevinoconnor7/digitalocean-crawler/queue"
)

type Crawler struct {
	httpClient  *http.Client
	Assets      map[string]*Asset `json:"-"`
	AssetsArray []*Asset          `json:"nodes"`
	Links       []*Link           `json:"links"`
	toProcess   queue.Queue
	baseURL     *url.URL
	matchRE     *regexp.Regexp
}

func Crawl(urlStr string, maxResults int) (*Crawler, error) {
	url, err := url.Parse(urlStr)

	if err != nil {
		return nil, err
	}

	c := GenerateCrawler(url)
	c.toProcess.Push(c.baseURL.String())
	i := 0
	for c.toProcess.Length > 0 {
		err = c.ProcessQueue()
		if err != nil {
			return nil, err
		}
		i++
		if i > maxResults {
			break
		}
	}

	return c, nil
}

func GenerateCrawler(url *url.URL) *Crawler {
	return &Crawler{
		baseURL: url,
		Assets:  make(map[string]*Asset),
	}
}

func (c *Crawler) ProcessQueue() error {
	if c.toProcess.Length <= 0 {
		return fmt.Errorf("Nothing left in queue to process")
	}

	url, ok := c.toProcess.Pop().(string)

	if !ok {
		return fmt.Errorf("Non-string type in queue")
	}

	asset, ok := c.Assets[url]

	if !ok {
		asset = GenerateAsset(url, PAGE)
		c.StoreAsset(asset)
	}

	if asset.processed == true {
		return nil
	}

	asset.processed = true

	req, err := GenerateRequest("GET", url, nil)

	if err != nil {
		return err
	}

	resp, err := c.GetClient().Do(req)

	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)

	err = c.ProcessDoc(doc, asset)

	return err
}

func (c *Crawler) GetQueue() *queue.Queue {
	return &c.toProcess
}

func (c *Crawler) ProcessDoc(doc *goquery.Document, asset *Asset) error {
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")

		if !ok {
			return
		}

		c.HandleURL(url, asset, PAGE)
	})

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")

		if !ok {
			return
		}

		c.HandleURL(url, asset, IMAGES)
	})

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")

		if !ok {
			return
		}

		c.HandleURL(url, asset, SCRIPTS)
	})

	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")

		if !ok {
			return
		}

		c.HandleURL(url, asset, STYLESHEETS)
	})
	return nil
}

func (c *Crawler) StoreAsset(asset *Asset) bool {
	_, ok := c.Assets[asset.URL]

	if !ok {
		asset.index = len(c.AssetsArray)
		c.AssetsArray = append(c.AssetsArray, asset)
		c.Assets[asset.URL] = asset
		c.toProcess.Push(asset.URL)
		return true
	}

	return false
}

func (c *Crawler) GetClient() *http.Client {
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	return c.httpClient
}

func (c *Crawler) HandleURL(url string, asset *Asset, contentType ContentType) {
	url = GetContextedURL(url, c.baseURL)
	if url == "" {
		return
	}

	newAsset := GenerateAsset(url, contentType)
	c.StoreAsset(newAsset)
	newLink := asset.AddLink(newAsset)

	if newLink != nil {
		c.Links = append(c.Links, newLink)
	}
}

func (c *Crawler) OutputResults() string {
	resp, _ := json.Marshal(c)
	return string(resp)
}

func GetContextedURL(urlStr string, baseUrl *url.URL) string {
	url, err := url.Parse(urlStr)

	if err != nil {
		return ""
	}

	if url.Host == "" {
		url.Host = baseUrl.Host
	}

	if url.Scheme == "" {
		url.Scheme = baseUrl.Scheme
	}

	if url.Host != baseUrl.Host {
		return ""
	}

	return url.String()
}

func GenerateRequest(method string, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Crawler - kevin@kevo.io")

	return req, nil
}
