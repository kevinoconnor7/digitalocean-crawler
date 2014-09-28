package crawler_test

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	. "github.com/kevinoconnor7/digitalocean-crawler/crawler"

	"net/http"
	"net/http/httptest"
	u "net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crawler", func() {
	Describe("GenerateRequest", func() {
		var (
			req *http.Request
			err error
		)

		BeforeEach(func() {
			req, err = GenerateRequest("GET", "https://foo.com", nil)
		})

		It("Should not have errored", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should return a http.Request", func() {
			Expect(req).ToNot(BeNil())
		})

		It("Should set the User-Agent header", func() {
			Expect(req.Header.Get("User-agent")).To(Equal("Crawler - kevin@kevo.io"))
		})

		It("Should set the method", func() {
			Expect(req.Method).To(Equal("GET"))
		})

		It("Should set the URL", func() {
			Expect(req.URL.String()).To(Equal("https://foo.com"))
		})
	})

	Describe("GetContextedURL", func() {
		var (
			BaseURL *u.URL
			url     string
		)

		BeforeEach(func() {
			BaseURL, _ = u.Parse("https://foo.com")
		})

		Context("For relative paths", func() {
			It("Should provide the entire URL in return", func() {
				url = "/home"
				Expect(GetContextedURL(url, BaseURL)).To(Equal("https://foo.com/home"))

				url = "/"
				Expect(GetContextedURL(url, BaseURL)).To(Equal("https://foo.com/"))
			})
		})

		Context("For absolute paths", func() {
			It("Should return empty for mismatched hosts", func() {
				url = "//bar.com"
				Expect(GetContextedURL(url, BaseURL)).To(Equal(""))

				url = "http://bar.com"
				Expect(GetContextedURL(url, BaseURL)).To(Equal(""))

				url = "https://bar.com"
				Expect(GetContextedURL(url, BaseURL)).To(Equal(""))

				url = "//bar.com/test"
				Expect(GetContextedURL(url, BaseURL)).To(Equal(""))
			})

			It("Should ignore mismatched schemes", func() {
				url = "http://foo.com"
				Expect(GetContextedURL(url, BaseURL)).To(Equal("http://foo.com"))

				url = "//foo.com"
				Expect(GetContextedURL(url, BaseURL)).To(Equal("https://foo.com"))
			})
		})

		It("Should return with proper path", func() {
			url = "/follow/me"
			Expect(GetContextedURL(url, BaseURL)).To(Equal("https://foo.com/follow/me"))

			url = "https://foo.com/follow/me"
			Expect(GetContextedURL(url, BaseURL)).To(Equal("https://foo.com/follow/me"))
		})

		It("Should return with proper query string", func() {
			url = "/?page=2"
			Expect(GetContextedURL(url, BaseURL)).To(Equal("https://foo.com/?page=2"))

			url = "https://foo.com/?page=2"
			Expect(GetContextedURL(url, BaseURL)).To(Equal("https://foo.com/?page=2"))
		})
	})

	Describe("Crawler", func() {

		var (
			c *Crawler
		)
		BeforeEach(func() {
			url, _ := u.Parse("http://foo.faketld")
			c = GenerateCrawler(url)
		})

		Describe("ProcessQueue", func() {

			var (
				ts *httptest.Server
			)

			BeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "test")
				}))
				url, _ := u.Parse(ts.URL)
				c = GenerateCrawler(url)
			})

			Context("With nothing in the queue", func() {
				It("Should produce an error", func() {
					err := c.ProcessQueue()
					Expect(err).To(HaveOccurred())
				})
			})

			Context("With a non-string item in the queue", func() {
				It("Should return an error", func() {
					c.GetQueue().Push(1)
					err := c.ProcessQueue()
					Expect(err).To(HaveOccurred())
				})
			})

			Context("With items to be processed", func() {
				It("Should process the top item", func() {
					c.StoreAsset(GenerateAsset(ts.URL, PAGE))
					err := c.ProcessQueue()

					Expect(err).ToNot(HaveOccurred())
					Expect(c.GetQueue().Length).To(Equal(0))
				})

				It("Should store an asset not in the map", func() {
					c.GetQueue().Push(ts.URL)
					err := c.ProcessQueue()
					Expect(err).ToNot(HaveOccurred())
					Expect(len(c.AssetsArray)).To(Equal(1))
				})
			})

			AfterEach(func() {
				ts.Close()
			})

		})

		Describe("ProcessDoc", func() {
			var (
				asset *Asset
			)
			BeforeEach(func() {
				asset = GenerateAsset("http://foo.faketld", PAGE)
				c.StoreAsset(asset)
			})

			It("Should match images", func() {
				buf := bytes.NewBuffer(nil)
				buf.WriteString(`
<html>
	<body>
		<img src="/some/image.jpg" />
		<img src="/some/image2.jpg" />
	</body>
</html>
				`)

				doc, _ := goquery.NewDocumentFromReader(buf)
				c.ProcessDoc(doc, asset)

				Expect(c.GetQueue().Length).To(Equal(3))
			})

			It("Should match links", func() {
				buf := bytes.NewBuffer(nil)
				buf.WriteString(`
<html>
	<body>
		<a href="/some/page"></a>
		<a href="/some/page2"></a>
	</body>
</html>
				`)

				doc, _ := goquery.NewDocumentFromReader(buf)
				c.ProcessDoc(doc, asset)

				Expect(c.GetQueue().Length).To(Equal(3))
			})

			It("Should match scripts", func() {
				buf := bytes.NewBuffer(nil)
				buf.WriteString(`
<html>
	<head>
		<script type="text/javascript" src="/main.js"></script>
	</head>
	<body>
		<script type="text/javascript" src="/src.js"></script>
	</body>
</html>
				`)

				doc, _ := goquery.NewDocumentFromReader(buf)
				c.ProcessDoc(doc, asset)

				Expect(c.GetQueue().Length).To(Equal(3))
			})

			It("Should match stylesheets", func() {
				buf := bytes.NewBuffer(nil)
				buf.WriteString(`
<html>
	<head>
		<link href="/src.css" rel="stylesheet" type="text/css">
		<link href="/src2.css" rel="stylesheet" type="text/css">
	</head>
</html>
				`)

				doc, _ := goquery.NewDocumentFromReader(buf)
				c.ProcessDoc(doc, asset)

				Expect(c.GetQueue().Length).To(Equal(3))
			})
		})

		Describe("HandleURL", func() {
			var (
				asset *Asset
			)
			BeforeEach(func() {
				asset = GenerateAsset("http://foo.faketld", PAGE)
				c.StoreAsset(asset)
			})

			Context("With unmatching base URL", func() {
				It("Should not store the asset", func() {
					c.HandleURL("http://bar.faketld/", asset, PAGE)
					Expect(c.GetQueue().Length).To(Equal(1))
					Expect(len(c.Assets)).To(Equal(1))
				})
			})

			Context("With matching base URL", func() {
				It("Should create a link if needed", func() {
					// Just to make sure this test is setup correctly we'll assume
					// that no links already exist
					Expect(len(c.Links)).To(Equal(0))

					c.HandleURL("http://foo.faketld/subpage", asset, PAGE)
					Expect(len(c.Links)).To(Equal(1))

					c.HandleURL("http://foo.faketld/subpage", asset, PAGE)
					Expect(len(c.Links)).To(Equal(1))
				})

				It("Should store the asset", func() {
					c.HandleURL("http://foo.faketld/subpage", asset, PAGE)
					Expect(c.GetQueue().Length).To(Equal(2))
				})
			})
		})

		Describe("StoreAsset", func() {
			var (
				asset *Asset
			)
			BeforeEach(func() {
				asset = GenerateAsset("http://foo.faketld", PAGE)
			})

			Context("Given an existing asset url", func() {
				It("Should return false without storing it", func() {
					Expect(c.StoreAsset(asset)).To(BeTrue())
					Expect(c.StoreAsset(asset)).To(BeFalse())
				})
			})

			Context("Given a new asset url", func() {
				It("Should return true", func() {
					Expect(c.StoreAsset(asset)).To(BeTrue())
				})

				It("Should be given an index", func() {
					c.StoreAsset(asset)
					Expect(asset.GetIndex()).To(Equal(0))

					otherAsset := GenerateAsset("http://bar.faketld", PAGE)
					c.StoreAsset(otherAsset)
					Expect(otherAsset.GetIndex()).To(Equal(1))
				})

				It("Should be stored in asset array", func() {
					c.StoreAsset(asset)
					Expect(c.AssetsArray[asset.GetIndex()]).To(Equal(asset))

					otherAsset := GenerateAsset("http://bar.faketld", PAGE)
					c.StoreAsset(otherAsset)
					Expect(c.AssetsArray[otherAsset.GetIndex()]).To(Equal(otherAsset))
				})

				It("Should be stored in asset map", func() {
					c.StoreAsset(asset)
					Expect(c.Assets[asset.URL]).To(Equal(asset))

					otherAsset := GenerateAsset("http://bar.faketld", PAGE)
					c.StoreAsset(otherAsset)
					Expect(c.Assets[otherAsset.URL]).To(Equal(otherAsset))
				})

				It("Should be added to the queue", func() {
					c.StoreAsset(asset)
					Expect(c.GetQueue().Pop().(string)).To(Equal(asset.URL))

					otherAsset := GenerateAsset("http://bar.faketld", PAGE)
					c.StoreAsset(otherAsset)
					Expect(c.GetQueue().Pop().(string)).To(Equal(otherAsset.URL))
				})
			})
		})
	})

	Describe("Crawl", func() {
		var (
			ts *httptest.Server
		)

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "test")
			}))
		})

		It("Should return a valid crawler", func() {
			c, err := Crawl(ts.URL, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(c).ToNot(BeNil())
		})

		It("Should process the queue", func() {
			c, err := Crawl(ts.URL, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(c.GetQueue().Length).To(Equal(0))
		})

		AfterEach(func() {
			ts.Close()
		})
	})
})
