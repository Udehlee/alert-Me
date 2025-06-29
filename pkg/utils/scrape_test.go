package utils

import (
	"testing"
)

func TestExtractProduct(t *testing.T) {
	URL := "https://www.konga.com/product/starlink-starlink-mini-6740842?cid=10154"

	product, err := ExtractProduct(URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if product.URL != URL {
		t.Errorf("Expected URL %s, got %s", URL, product.URL)
	}

	if product.Name == "" {
		t.Errorf("Expected product name, got empty string")
	}

	if product.Price <= 0 {
		t.Errorf("Expected product price > 0, got %.2f", product.Price)
	}

}
