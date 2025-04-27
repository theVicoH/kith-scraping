package scheduler

import (
	"context"
	"database/sql"
	"log"
	"time"

	"backend/internal/product"
	"backend/internal/service/scraper"
)

func StartScheduler(ctx context.Context, db *sql.DB, prodServ product.Service, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	run := func() {
		log.Println("Scheduler: début de la synchronisation des produits")

		targets, err := scraper.LoadActiveTargets(db)
		if err != nil {
			log.Printf("Scheduler: erreur LoadActiveTargets: %v", err)
			return
		}

		for _, t := range targets {
			s := scraper.NewScraper(t.URL)

			scraped, err := s.FetchProducts(ctx)
			if err != nil {
				log.Printf("Scheduler: erreur FetchProducts pour URL %s: %v", t.URL, err)
				continue
			}

			if err := prodServ.SyncProducts(ctx, scraped, t.Category); err != nil {
				log.Printf("Scheduler: erreur SyncProductsFor pour URL %s: %v", t.URL, err)
			} else {
				log.Printf("Scheduler: synchronisation terminée pour URL %s", t.URL)
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler arrêté")
			return
		case <-ticker.C:
			run()
		}
	}
}
