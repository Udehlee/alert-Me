package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Udehlee/alert-Me/models"

	"github.com/gocolly/colly"
)

// it is assumed that you enter a product url, so
// scrapers maps domain names to their specific functions
// and scrape the selected product’s name and price from its URL
var (
	nameTag  = os.Getenv("NAME_SELECTOR")
	priceTag = os.Getenv("PRICE_SELECTOR")

	scrapers = map[string]func(url string) (string, string, error){
		"jumia.com.ng": func(url string) (string, string, error) {
			var name, price string
			c := colly.NewCollector(
				colly.AllowedDomains("www.jumia.com.ng", "jumia.com.ng"),
			)

			c.OnHTML(nameTag, func(e *colly.HTMLElement) {
				name = strings.TrimSpace(e.Text)
			})

			c.OnHTML(priceTag, func(e *colly.HTMLElement) {
				price = strings.TrimSpace(e.Text)
			})

			c.OnError(func(r *colly.Response, err error) {
				log.Printf("Error scraping jumia: %v, URL: %s", err, r.Request.URL)
			})

			c.Visit(url)
			return name, price, nil
		},

		"konga.com": func(url string) (string, string, error) {
			var name, price string
			c := colly.NewCollector(
				colly.AllowedDomains("www.konga.com", "konga.com"),
			)

			c.OnHTML(nameTag, func(e *colly.HTMLElement) {
				name = strings.TrimSpace(e.Text)
			})

			c.OnHTML(priceTag, func(e *colly.HTMLElement) {
				price = strings.TrimSpace(e.Text)
			})

			c.OnError(func(r *colly.Response, err error) {
				log.Printf("Error scraping konga: %v, URL: %s", err, r.Request.URL)
			})

			c.Visit(url)
			return name, price, nil
		},
	}
)

// ExtractProduct gets a product's name and price from its URL
// using the associated scraper based on its domain
func ExtractProduct(Url string) (models.Product, error) {
	product := models.Product{}

	u, err := url.Parse(Url)
	if err != nil {
		return product, fmt.Errorf("invalid url: %w", err)
	}
	domain := u.Hostname()

	for k, scraper := range scrapers {
		if strings.Contains(domain, k) {
			name, price, err := scraper(Url)
			if err != nil {
				return product, err
			}

			p, err := ToFloatPrice(price)
			if err != nil {
				return product, err
			}
			product.Name = name
			product.Price = p

			return product, nil
		}
	}

	//if there is no match
	return product, fmt.Errorf("no scraper found for domain: %s", domain)
}
