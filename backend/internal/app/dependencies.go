package app

import (
	"database/sql"

	"backend/internal/product"
)

func InitProductService(db *sql.DB) product.Service {
	repo := product.NewPostgresRepository(db)
	return product.NewService(repo)
}

func InitProductHandler(service product.Service) *product.ProductHandler {
	return product.NewProductHandler(service)
}
