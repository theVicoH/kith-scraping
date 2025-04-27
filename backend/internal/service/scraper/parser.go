package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Category struct {
	URL      string
	Category string
}

func ParseCategory(e *colly.HTMLElement) (Category, error) {
	href := strings.TrimSpace(e.ChildAttr("a", "href"))
	title := strings.TrimSpace(e.ChildText("h2.h1"))
	if href == "" || title == "" {
		return Category{}, fmt.Errorf("ParseCategory: href or title missing")
	}
	return Category{
		URL:      href,
		Category: title,
	}, nil
}

type Product struct {
	Name     string
	Price     int
	URL       string
	ImageURL  string
	InStock   bool
	Reference string
}

func ParseProduct(e *colly.HTMLElement) (Product, error) {
	link := strings.TrimSpace(e.ChildAttr("a.product-card__image", "href"))
	name := strings.TrimSpace(e.ChildText("div.product-card-info__text div.text-black"))
	priceStr := strings.TrimSpace(e.ChildText("div.product-card-info__text div.text-10 span"))
	img := strings.TrimSpace(e.ChildAttr("a.product-card__image picture img", "src"))
	if link == "" || name == "" || priceStr == "" || img == "" {
		return Product{}, fmt.Errorf("ParseProduct: missing elements in product card")
	}

	reference := link
	if parts := strings.Split(link, "/"); len(parts) > 0 {
		reference = parts[len(parts)-1]
	}

	priceStr = strings.TrimPrefix(priceStr, "€")

	priceStr = strings.ReplaceAll(priceStr, ",", "")

	priceFloat, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Printf("Raw price string: %q\n", priceStr)
		return Product{}, fmt.Errorf("ParseProduct: impossible to convert price '%s': %w", priceStr, err)
	}
	priceInCents := int(priceFloat * 100)

	return Product{
		Name:      name,
		Price:     priceInCents,
		URL:       "https://eu.kith.com" + link,
		ImageURL:  "https:" + img,
		InStock:   true,
		Reference: reference,
	}, nil
}
