package models

import "time"

type Subscription struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	Plan       string    `json:"plan"`
	Duration   int       `json:"duration"`
	Price      float64   `json:"price"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}
