package models

type ProductView struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
