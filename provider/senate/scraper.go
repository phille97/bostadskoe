package senate

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	httpClient *http.Client
	client     *Client
}

func NewScraper(client *Client, httpClient *http.Client) *Scraper {
	return &Scraper{
		httpClient: httpClient,
		client:     client,
	}
}

func (s Scraper) fetch(u *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")
	req.Header.Set("User-Agent", s.client.UserAgent)

	return s.httpClient.Do(req)
}

func (s Scraper) Ads(u *url.URL) (*[]url.URL, error) {
	ads := []url.URL{}

	res, err := s.fetch(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".apartment-list > .apartment").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Find("a").Attr("href")
		if !exists {
			return
		}
		u, err := url.Parse(href)
		if err != nil {
			return
		}

		ads = append(ads, *u)
	})

	return &ads, nil
}

func (s Scraper) AdDetails(u *url.URL) (*Bostad, error) {
	ad := Bostad{Url: u.String()}

	res, err := s.fetch(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".apartment-data").Each(func(i int, s *goquery.Selection) {
		parts := strings.Split(strings.Trim(s.Text(), " \n\t"), "\n")
		if len(parts) != 2 {
			return
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch key {
		case "Objektnr:":
			ad.ObjektNummer = value
		case "Antal rum:":
			antalRum, err := strconv.ParseFloat(value, 64)
			if err == nil {
				ad.AntalRum = &antalRum
			} else {
				fmt.Println(err.Error())
			}
		case "Storlek:":
			value = strings.TrimSuffix(value, " kvm")
			value = strings.TrimSuffix(value, " m2")
			value = strings.ReplaceAll(value, ".", "")
			value = strings.ReplaceAll(value, ",", "")
			yta, err := strconv.ParseFloat(value, 64)
			if err == nil {
				ad.Yta = &yta
			} else {
				fmt.Println(err.Error())
			}
		case "Hyra:":
			value = strings.TrimSuffix(value, ":-/månad")
			value = strings.TrimSuffix(value, ":-/mån")
			value = strings.TrimSuffix(value, ":-/mån (2017)")
			value = strings.ReplaceAll(value, ".", "")
			value = strings.ReplaceAll(value, ",", "")
			hyra, err := strconv.ParseFloat(value, 64)
			if err == nil {
				ad.Hyra = &hyra
			} else {
				fmt.Println(err.Error())
			}
		case "Parkering:":
			if strings.ToLower(value) == "nej" {
				ad.Parkering = false
			} else {
				ad.Parkering = true
			}
		case "Balkong:":
			if strings.ToLower(value) == "nej" {
				ad.Balkong = false
			} else {
				ad.Balkong = true
			}
		case "TV:":
			ad.TV = &value
		case "Internet:":
			ad.Internet = &value
		case "Ledig fr.o.m:":
			ad.LedigFran = &value
		}

		fmt.Println(key)
		fmt.Println(value)
	})

	if ad.ObjektNummer == "" {
		return nil, errors.New("Couldn't find ObjektNummer")
	}

	return &ad, nil
}
