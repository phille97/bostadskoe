package senate

import (
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

func (c Client) CurrentResidences(out chan provider.Residence, errout chan error) {
	defer close(out)

	s := NewScraper(&c, c.httpClient)

	kommuner := map[*url.URL]string{
		c.BaseURL.ResolveReference(&url.URL{Path: "/lidkoping/bostader"}): "Lidköping",
		c.BaseURL.ResolveReference(&url.URL{Path: "/falkoping/bostader"}): "Falköping",
		c.BaseURL.ResolveReference(&url.URL{Path: "/skara/bostader"}):     "Skara",
		c.BaseURL.ResolveReference(&url.URL{Path: "/gavle/bostader"}):     "Gävle",
		c.BaseURL.ResolveReference(&url.URL{Path: "/degerfors/bostader"}): "Degerfors",
		c.BaseURL.ResolveReference(&url.URL{Path: "/koping/bostader"}):    "Köping",
	}

	for listUrl, kommun := range kommuner {
		ads, err := s.Ads(listUrl)
		if err != nil {
			errout <- err
			return
		}

		for _, ad := range *ads {
			bostad, err := s.AdDetails(&ad)
			if err != nil {
				errout <- err
				return
			}
			bostad.Kommun = &kommun

			out <- bostad
		}
	}
}
