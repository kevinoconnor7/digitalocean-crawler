package crawler_test

import (
	. "github.com/kevinoconnor7/digitalocean-crawler/crawler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Asset", func() {
	Describe("GenerateAsset", func() {
		It("Should initalize the asset with given parameters", func() {
			asset := GenerateAsset("foo.com", IMAGES)
			Expect(asset.URL).To(Equal("foo.com"))
			Expect(asset.Type).To(Equal(IMAGES))
			Expect(asset.Links).ToNot(BeNil())
		})
	})

	Describe("GenerateLink", func() {
		var (
			asset    *Asset
			newAsset *Asset
			link     *Link
		)

		BeforeEach(func() {
			asset = GenerateAsset("foo.com", IMAGES)
			newAsset = GenerateAsset("blah.com", PAGE)
			link = GenerateLink(asset.GetIndex(), newAsset.GetIndex())
		})

		It("Should store the correct target/source", func() {
			Expect(link.Source).To(Equal(asset.GetIndex()))
			Expect(link.Target).To(Equal(newAsset.GetIndex()))
		})
		It("Should intialize the value to one", func() {
			Expect(link.Value).To(Equal(1))
		})
	})

	Describe("AddLink", func() {
		var (
			asset    *Asset
			newAsset *Asset
		)

		BeforeEach(func() {
			asset = GenerateAsset("foo.com", IMAGES)
			newAsset = GenerateAsset("blah.com", PAGE)
			asset.AddLink(newAsset)
		})

		Context("When creating a new link", func() {
			It("Should store the link on the asset", func() {
				link, ok := asset.Links[newAsset.URL]
				Expect(ok).To(BeTrue())
				Expect(link).ToNot(BeNil())
			})

			It("Should return a pointer to the new link", func() {
				otherAsset := GenerateAsset("bar.com", PAGE)
				link := asset.AddLink(otherAsset)
				storedLink, _ := asset.Links[otherAsset.URL]
				Expect(link).To(Equal(storedLink))
			})
		})

		Context("When adding an already existing link", func() {
			It("Shouldn't replace the link", func() {
				link, _ := asset.Links[newAsset.URL]
				linkAddress := &link

				asset.AddLink(newAsset)
				link, _ = asset.Links[newAsset.URL]

				Expect(&link).To(Equal(linkAddress))
			})

			It("Should increment the value", func() {
				link, _ := asset.Links[newAsset.URL]
				value := link.Value

				asset.AddLink(newAsset)
				Expect(link.Value).To(Equal(value + 1))

				asset.AddLink(newAsset)
				Expect(link.Value).To(Equal(value + 2))
			})

			It("Shouldn't modify source/target", func() {
				link, _ := asset.Links[newAsset.URL]
				source := link.Source
				target := link.Target

				asset.AddLink(newAsset)

				Expect(link.Source).To(Equal(source))
				Expect(link.Target).To(Equal(target))
			})

			It("Should return nil", func() {
				Expect(asset.AddLink(newAsset)).To(BeNil())
			})
		})
	})
})
