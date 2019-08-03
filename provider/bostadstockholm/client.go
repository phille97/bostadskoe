package bostadstockholm

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/phille97/bostadskoe/provider"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

func New(baseUrl string, httpClient *http.Client) (*Client, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL:    u,
		UserAgent:  "Mozilla/5.0 (compatible; Bostadskoe/1.0; +https://github.com/phille97/bostadskoe)",
		httpClient: httpClient,
	}, nil
}

func (c Client) CurrentResidences() (*[]provider.Residence, error) {
	u := c.BaseURL.ResolveReference(&url.URL{Path: "/Lista/AllaAnnonser"})
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bostader []Bostad
	err = json.NewDecoder(resp.Body).Decode(&bostader)
	if err != nil {
		return nil, err
	}

	residenceSlice := make([]provider.Residence, len(bostader))
	for i, b := range bostader {
		residenceSlice[i] = b
	}

	return &residenceSlice, nil
}
