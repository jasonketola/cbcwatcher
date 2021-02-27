package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"io/ioutil"

	"github.com/gocolly/colly"
)

// Product stores information about each product
type Product struct {
	Brand       string
	Name        string
	Description string
	URL         string
}



func main() {

	j := 1
	k := 0
	n := 0
	s := 0

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("belmont.craftbeercellar.com"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./cbc_cache"),
	)

	// Create another collector to scrape product details
	detailCollector := c.Clone()

	products := make([]Product, 0, 2000)
	crap := make([]Product, 0, 2000)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})


	data, err := ioutil.ReadFile("products_prior.json")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	// On every a HTML element which has name attribute call callback
	c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the class contains "product-link product-thumbnail"
		productURL := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Contains(e.Attr("class"), "product-link product-thumbnail") {
			fmt.Println("Reviewing ", j, " of 1733")

			if strings.Contains(string(data), productURL) {
				k++
				fmt.Println("Matched ", k)
			} else {
				n++
				fmt.Println("New ", n)
				fmt.Println(productURL)
		//		time.Sleep(7 * time.Second)
				detailCollector.Visit(productURL)
			}
	





			j++
		}
	})

	// Extract details of the product
	detailCollector.OnHTML(`div[class="col-sm-8 product-right"]`, func(e *colly.HTMLElement) {

		brand := strings.TrimSpace(strings.TrimLeft(strings.TrimRight(e.ChildText("h6"), "Share: "), "Brand:"))


		product := Product{
			Brand:       brand,
			Name:        e.ChildText("h1"),
			Description: e.ChildText(`p[class="text-product-desc"]`),
			URL:         e.Request.URL.String(),
		}


		if brand == "" {
			s++
			crap = append(crap, product)
			return
		} else {
			products = append(products, product)
		}
	})

	for i := 0; i < 18; i++ {
		myURL := "https://belmont.craftbeercellar.com/store/search.asp?matchesperpage=100&start=" + strconv.Itoa(i)
		time.Sleep(7 * time.Second)
		c.Visit(myURL)
	}
	file, _ := os.Create("products.json")
	file2,_ := os.Create("crap.json")
	defer file.Close()
	defer file2.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	enc2 := json.NewEncoder(file2)
	enc2.SetIndent("", " ")

	// Dump json to the standard output
	enc.Encode(products)
	enc2.Encode(crap)
	fmt.Println("done")

	println("Reviewed: ", j)
	println("Matched: ", k)
	println("Skipped: ", s)
	println("New: ", n-s)
}
