package product

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	restockClients    = make(map[chan []byte]bool)
	newProductClients = make(map[chan []byte]bool)
)

type ProductHandler struct {
	service Service
}

func NewProductHandler(service Service) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts()
	if err != nil {
		http.Error(w, "failed to list products", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) SSEProductRestock(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan []byte)
	restockClients[messageChan] = true

	defer func() {
		delete(restockClients, messageChan)
		close(messageChan)
	}()

	for {
		msg, open := <-messageChan
		if !open {
			break
		}
		w.Write([]byte("data: "))
		w.Write(msg)
		w.Write([]byte("\n\n"))
		flusher.Flush()
	}
}

func (h *ProductHandler) SSENewProduct(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan []byte)
	newProductClients[messageChan] = true

	defer func() {
		delete(newProductClients, messageChan)
		close(messageChan)
	}()

	for {
		msg, open := <-messageChan
		if !open {
			break
		}
		w.Write([]byte("data: "))
		w.Write(msg)
		w.Write([]byte("\n\n"))
		flusher.Flush()
	}
}

func BroadcastRestock(product *Product) {
	message, err := json.Marshal(product)
	if err != nil {
		log.Println("error marshaling restock:", err)
		return
	}
	for client := range restockClients {
		client <- message
	}
}

func BroadcastNewProduct(product *Product) {
	message, err := json.Marshal(product)
	if err != nil {
		log.Println("error marshaling new product:", err)
		return
	}
	for client := range newProductClients {
		client <- message
	}
}
