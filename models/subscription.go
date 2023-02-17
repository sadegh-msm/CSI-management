package models

type Subscription struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	IsActive   bool    `json:"is_active"`
	Period     int     `json:"period"`
	CustomerID int64   `json:"customer_id"`
}
