package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Udehlee/alert-Me/models"
)

type EbayClient struct {
	Client      *http.Client
	AccessToken string
	BaseURL     string
}

func NewEbayClient(token string) *EbayClient {
	return &EbayClient{
		Client:      &http.Client{},
		AccessToken: token,
		BaseURL:     os.Getenv("BASE_URL"),
	}
}

// GetProduct retrieves specific product from ebays api endpoint
func (ec *EbayClient) GetProduct(query string) (*models.Product, error) {
	url := fmt.Sprintf("%s/buy/browse/v1/item_summary/search?q=%s&limit=1", ec.BaseURL, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ec.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := ec.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("eBay API error: %s", res.Status)
	}

	var result models.SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Products[0], nil
}
