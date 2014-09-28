package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/kevinoconnor7/digitalocean-crawler/crawler"
)

func main() {
	maxPages := flag.Int("maxPages", 10, "Max number of pages to crawl")
	flag.Parse()
	c, err := crawler.Crawl("https://www.digitalocean.com", *maxPages)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/graph.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, c.OutputResults())
	})
	fmt.Println("Starting web server on :8080")
	panic(http.ListenAndServe(":8080", nil))
}
