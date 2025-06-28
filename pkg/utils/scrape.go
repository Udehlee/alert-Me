package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Udehlee/alert-Me/models"

	"github.com/gocolly/colly"
)

var (
	domain   = os.Getenv("DOMAIN")
	nameTag  = os.Getenv("NAME_SELECTOR")
	priceTag = os.Getenv("PRICE_SELECTOR")
)

// Scraper scrape the selected product’s name and price from its URL
// provided you also entered the specific E-commerce domain, nameTag and priceTag
func Scraper(url string) (string, string, error) {
	var name, price string

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)
	c.SetRequestTimeout(30 * time.Second)

	c.OnHTML(nameTag, func(e *colly.HTMLElement) {
		name = strings.TrimSpace(e.Text)
		fmt.Printf("product name: %s\n", name)
	})

	c.OnHTML(priceTag, func(e *colly.HTMLElement) {
		price = strings.TrimSpace(e.Text)
		fmt.Printf("product price: %s\n", price)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping product_url %s: %v", r.Request.URL, err)
	})

	if err := c.Visit(url); err != nil {
		return "", "", fmt.Errorf("visit product_url failed: %w", err)
	}
	c.Wait()

	if name == "" || price == "" {
		return "", "", fmt.Errorf("empty scraped data: name=%q\n, price=%q", name, price)
	}

	return name, price, nil
}

// ExtractProduct extract a product's name and price from its URL
func ExtractProduct(url string) (models.Product, error) {
	product := models.Product{
		URL: url,
	}

	name, price, err := Scraper(url)
	if err != nil {
		return product, err
	}

	p, err := ToFloat(price)
	if err != nil {
		return product, err
	}

	product.Name = name
	product.Price = p

	return product, nil
}
