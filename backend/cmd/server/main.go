package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/database"

	"backend/internal/service/scheduler"

	"backend/internal/service/scraper"

	"backend/internal/router"
	"net/http"

	"context"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config.Load: %v", err)
	}
	dbConn, err := database.NewDatabaseClient(cfg.Database)
	if err != nil {
		log.Fatalf("NewDatabaseClient: %v", err)
	}
	defer dbConn.DB.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	homeScraper := scraper.NewScraper("https://eu.kith.com/")
	if err := scraper.BootstrapTargets(ctx, dbConn.DB, homeScraper); err != nil {
		log.Fatalf("BootstrapTargets: %v", err)
	}

	productService := app.InitProductService(dbConn.DB)

	go scheduler.StartScheduler(ctx, dbConn.DB, productService, 5*time.Minute)

	productHandler := app.InitProductHandler(productService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.SetupRouter(productHandler),
	}

	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}() 
	<-ctx.Done()

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
