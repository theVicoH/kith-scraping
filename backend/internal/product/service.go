package product

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/internal/service/scraper"
)

type Service interface {
	SyncProducts(ctx context.Context, scraped []scraper.Product, category string) error
	ListProducts() ([]Product, error)
}

type productService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &productService{repo: repo}
}

func (s *productService) SyncProducts(ctx context.Context, scraped []scraper.Product, category string) error {
	existingRefs, err := s.repo.GetAllReferences(ctx)
	if err != nil {
		return fmt.Errorf("Service.GetAllReferences: %w", err)
	}

	existingMap := make(map[string]bool)
	for _, ref := range existingRefs {
		existingMap[ref] = true
	}

	var prods []Product
	var newProductRefs []string
	var potentialRestockRefs []string

	for _, sp := range scraped {
		prod := Product{
			Reference:  sp.Reference,
			Title:      sp.Name,
			Price:      sp.Price,
			ImageURL:   sp.ImageURL,
			ProductURL: sp.URL,
			Category:   category,
			EventType:  "",
			EventDate:  time.Now(),
			InStock:    sp.InStock,
		}

		prods = append(prods, prod)

		if !existingMap[sp.Reference] && sp.InStock {
			newProductRefs = append(newProductRefs, sp.Reference)
			log.Printf("New product detected: %s", prod.Title)
		} else if existingMap[sp.Reference] && sp.InStock {
			isRestock, err := s.repo.WasOutOfStock(ctx, sp.Reference)
			if err != nil {
				log.Printf("Error checking previous stock status: %v", err)
			} else if isRestock {
				potentialRestockRefs = append(potentialRestockRefs, sp.Reference)
				log.Printf("Product restocked: %s", prod.Title)
			}
		}
	}

	savedProds, err := s.repo.UpsertProducts(ctx, prods)
	if err != nil {
		return fmt.Errorf("Service.UpsertProducts: %w", err)
	}

	prodMap := make(map[string]*Product)
	for i := range savedProds {
		prodMap[savedProds[i].Reference] = &savedProds[i]
	}

	for _, ref := range newProductRefs {
		if prod, ok := prodMap[ref]; ok {
			BroadcastNewProduct(prod)
		}
	}

	for _, ref := range potentialRestockRefs {
		if prod, ok := prodMap[ref]; ok {
			BroadcastRestock(prod)
		}
	}

	found := make(map[string]bool)
	for _, p := range prods {
		found[p.Reference] = true
	}

	var missing []string
	for _, ref := range existingRefs {
		if !found[ref] {
			missing = append(missing, ref)
		}
	}

	if err := s.repo.MarkOutOfStock(ctx, missing, category); err != nil {
		return fmt.Errorf("Service.MarkOutOfStock: %w", err)
	}
	return nil
}

func (s *productService) ListProducts() ([]Product, error) {
	ctx := context.Background()
	return s.repo.GetAllProducts(ctx)
}
