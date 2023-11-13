package gbooks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klaital/library/storage/library"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	ApiKey  string
	BaseUrl string

	http http.Client
}

func New(apiKey string) *Client {
	return &Client{
		ApiKey:  apiKey,
		BaseUrl: "https://www.googleapis.com/books/v1/volumes",
		http:    http.Client{},
	}
}

type isbnLookupResponse struct {
	Kind       string `json:"kind"`
	TotalItems int    `json:"totalItems"`
	Items      []struct {
		Kind       string `json:"kind"`
		ID         string `json:"id"`
		Etag       string `json:"etag"`
		SelfLink   string `json:"selfLink"`
		VolumeInfo struct {
			Title               string   `json:"title"`
			Subtitle            string   `json:"subtitle"`
			Authors             []string `json:"authors"`
			Publisher           string   `json:"publisher"`
			PublishedDate       string   `json:"publishedDate"`
			Description         string   `json:"description"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			ReadingModes struct {
				Text  bool `json:"text"`
				Image bool `json:"image"`
			} `json:"readingModes"`
			PageCount           int      `json:"pageCount"`
			PrintType           string   `json:"printType"`
			Categories          []string `json:"categories"`
			MaturityRating      string   `json:"maturityRating"`
			AllowAnonLogging    bool     `json:"allowAnonLogging"`
			ContentVersion      string   `json:"contentVersion"`
			PanelizationSummary struct {
				ContainsEpubBubbles  bool `json:"containsEpubBubbles"`
				ContainsImageBubbles bool `json:"containsImageBubbles"`
			} `json:"panelizationSummary"`
			ImageLinks struct {
				SmallThumbnail string `json:"smallThumbnail"`
				Thumbnail      string `json:"thumbnail"`
			} `json:"imageLinks"`
			Language            string `json:"language"`
			PreviewLink         string `json:"previewLink"`
			InfoLink            string `json:"infoLink"`
			CanonicalVolumeLink string `json:"canonicalVolumeLink"`
		} `json:"volumeInfo"`
		SaleInfo struct {
			Country     string `json:"country"`
			Saleability string `json:"saleability"`
			IsEbook     bool   `json:"isEbook"`
		} `json:"saleInfo"`
		AccessInfo struct {
			Country                string `json:"country"`
			Viewability            string `json:"viewability"`
			Embeddable             bool   `json:"embeddable"`
			PublicDomain           bool   `json:"publicDomain"`
			TextToSpeechPermission string `json:"textToSpeechPermission"`
			Epub                   struct {
				IsAvailable bool `json:"isAvailable"`
			} `json:"epub"`
			Pdf struct {
				IsAvailable bool `json:"isAvailable"`
			} `json:"pdf"`
			WebReaderLink       string `json:"webReaderLink"`
			AccessViewStatus    string `json:"accessViewStatus"`
			QuoteSharingAllowed bool   `json:"quoteSharingAllowed"`
		} `json:"accessInfo"`
		SearchInfo struct {
			TextSnippet string `json:"textSnippet"`
		} `json:"searchInfo"`
	} `json:"items"`
}

func (c *Client) LookupIsbn(ctx context.Context, isbn string) (*library.Item, error) {
	req, err := http.NewRequest("GET", c.BaseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("constructing ISBN lookup request: %w", err)
	}
	q := url.Values{}
	if len(c.ApiKey) > 0 {
		q.Add("key", c.ApiKey)
	}
	q.Add("q", fmt.Sprintf("isbn:%s", isbn))
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qurying isbn: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status indicates failure: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	//fmt.Printf("Raw response:\n%s\n", string(b))
	var respData isbnLookupResponse
	err = json.Unmarshal(b, &respData)
	if err != nil {
		return nil, fmt.Errorf("parsing response body: %w", err)
	}

	if len(respData.Items) == 0 {
		return nil, errors.New("no data in response")
	}

	return &library.Item{
		ID:         0,
		LocationID: 0,
		Code:       isbn,
		CodeType:   "ISBN",
		CodeSource: respData.Items[0].SelfLink,
		Title:      respData.Items[0].VolumeInfo.Title,
	}, nil
}
