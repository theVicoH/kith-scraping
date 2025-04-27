package scraper

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	collector *colly.Collector
	baseURL   string
}

func NewScraper(baseURL string) *Scraper {
	c := colly.NewCollector(
		colly.AllowedDomains("eu.kith.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; KithScraper/1.0)"),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})
	return &Scraper{collector: c, baseURL: baseURL}
}

func (s *Scraper) FetchCategories(ctx context.Context) ([]Category, error) {
	var cats []Category

	s.collector.OnHTML("button[js-open-mobile-nav]", func(e *colly.HTMLElement) {
		_ = e.Request.Visit(s.baseURL + "collections")
	})

	s.collector.OnHTML("a[href^=\"/collections/\"]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		fullURL := e.Request.AbsoluteURL(href)
		name := strings.TrimSpace(e.Text)
		if name == "" {
			return
		}
		cats = append(cats, Category{URL: fullURL, Category: name})
	})

	s.collector.OnError(func(_ *colly.Response, err error) {
		log.Printf("FetchCategories error: %v", err)
	})

	if err := s.collector.Visit(s.baseURL); err != nil {
		return nil, err
	}
	s.collector.Wait()
	return cats, nil
}

func (s *Scraper) FetchProducts(ctx context.Context) ([]Product, error) {
	var prods []Product

	s.collector.OnHTML("li.collection-grid__product", func(e *colly.HTMLElement) {
		if p, err := ParseProduct(e); err == nil {
			prods = append(prods, p)
		} else {
			log.Printf("FetchProducts.ParseProduct: %v", err)
		}
	})
	s.collector.OnError(func(_ *colly.Response, err error) {
		log.Printf("FetchProducts error: %v", err)
	})

	if err := s.collector.Visit(s.baseURL); err != nil {
		return nil, err
	}
	s.collector.Wait()
	return prods, nil
}
