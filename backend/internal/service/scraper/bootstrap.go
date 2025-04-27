package scraper

import (
	"context"
	"database/sql"
	"fmt"
)

func BootstrapTargets(ctx context.Context, db *sql.DB, homeScraper *Scraper) error {
	fmt.Println("BootstrapTargets: démarrage")
	cats, err := homeScraper.FetchCategories(ctx)
	if err != nil {
		return fmt.Errorf("BootstrapTargets.FetchCategories: %w", err)
	}

	fmt.Println("BootstrapTargets: categories récupérées", len(cats))

	for _, c := range cats {
		existing, err := GetScrapeTargetByURL(db, c.URL)
		if err != nil {
			return fmt.Errorf("BootstrapTargets.GetByURL: %w", err)
		}

		fmt.Println("BootstrapTargets: target", c.URL, "existant:", existing != nil)

		if existing == nil {
			fmt.Println("BootstrapTargets: ajout du target", c.URL)
			if err := InsertScrapeTarget(db, c.URL, c.Category); err != nil {
				return fmt.Errorf("BootstrapTargets.Insert: %w", err)
			}
		} else {
			fmt.Println("BootstrapTargets: mise à jour du target", c.URL)

			if err := UpdateScrapeTarget(db, c.URL, c.Category); err != nil {
				return fmt.Errorf("BootstrapTargets.Update: %w", err)
			}
		}
	}
	return nil
}
