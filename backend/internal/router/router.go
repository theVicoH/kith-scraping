package router

import (
	"net/http"

	"backend/internal/product"
)

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Update this line to allow requests from your frontend origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func SetupRouter(productHandler *product.ProductHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/products", productHandler.ListProducts)
	mux.HandleFunc("/api/sse/restocks", productHandler.SSEProductRestock)
	mux.HandleFunc("/api/sse/new-products", productHandler.SSENewProduct)

	return withCORS(mux)
}
