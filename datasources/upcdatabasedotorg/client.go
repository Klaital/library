package upcdatabasedotorg

import (
	"encoding/json"
	"fmt"
	"github.com/klaital/library/storage/library"
	"io"
	"net/http"
)

type Client struct {
	BaseUrl string
	http    *http.Client
	ApiKey  string
}

func New(apiKey string) *Client {
	return &Client{
		BaseUrl: "https://api.upcdatabase.org",
		http:    http.DefaultClient,
		ApiKey:  apiKey,
	}
}

type LookupProductResponse struct {
	Error        string `json:"error"`
	Success      bool   `json:"success"`
	Barcode      string `json:"barcode"`
	Title        string `json:"title"`
	Alias        string `json:"alias"`
	Description  string `json:"description"`
	Brand        string `json:"brand"`
	Manufacturer string `json:"manufacturer"`
	Mpn          string `json:"mpn"`
	Msrp         string `json:"msrp"`
	Asin         string `json:"ASIN"`
	Category     string `json:"category"`
	Metadata     struct {
		Size       string `json:"size"`
		Color      string `json:"color"`
		Gender     string `json:"gender"`
		Age        string `json:"age"`
		Length     string `json:"length"`
		Unit       string `json:"unit"`
		Width      string `json:"width"`
		Height     string `json:"height"`
		Weight     string `json:"weight"`
		Quantity   string `json:"quantity"`
		Publisher  string `json:"publisher"`
		Genre      string `json:"genre"`
		Author     string `json:"author"`
		Relasedate string `json:"relasedate"`
	} `json:"metadata"`
	Stores []struct {
		Store string `json:"store"`
		Price string `json:"price"`
	} `json:"stores"`
	Images  []string `json:"images"`
	Reviews struct {
		Thumbsup   int `json:"thumbsup"`
		Thumbsdown int `json:"thumbsdown"`
	} `json:"reviews"`
}

func (c *Client) LookupUpc(upc string) (*library.Item, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/product/%s?apikey=%s", c.BaseUrl, upc, c.ApiKey), nil)
	if err != nil {
		return nil, fmt.Errorf("constructing upc lookup request: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending upc lookup request: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("upc lookup response status error: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading upc response body: %w", err)
	}

	var data LookupProductResponse
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Printf("Failed to deserialize:\n\n%s\n\n", string(b))
		return nil, fmt.Errorf("deserializing response body: %w", err)
	}

	if len(data.Error) > 0 {
		return nil, fmt.Errorf("Error response from upcdatabase: %s", data.Error)
	}

	item := library.Item{
		Code:       upc,
		CodeType:   "UPC",
		CodeSource: "upcdatabase.org",
		Title:      data.Title,
	}

	// Some entries have no data in the Title. Try using the Description field as a fallback.
	if len(data.Title) == 0 {
		item.Title = data.Description
	}

	return &item, nil
}
