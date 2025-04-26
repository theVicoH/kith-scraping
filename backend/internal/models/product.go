package model

import "time"

type Product struct {
	ID         int       `json:"id"`
	Reference  string    `json:"reference"`
	Title      string    `json:"title"`
	Price      int       `json:"price"`
	ImageURL   string    `json:"image_url"`
	ProductURL string    `json:"product_url"`
	Category   string    `json:"category"`
	EventType  string    `json:"event_type"`
	EventDate  time.Time `json:"event_date"`
}
