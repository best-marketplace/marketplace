package models

import "time"

type ProductView struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Price        int       `json:"price"`
	SellerName   string    `json:"seller_name"`
	CategoryName string    `json:"category_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
