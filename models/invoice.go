package models

import (
	"time"
)

type Invoice struct {
	ID             int       `json:"id"`
	CustomerID     int       `json:"customer_id"`
	SubscriptionID int       `json:"subscription_id"`
	Amount         float64   `json:"amount"`
	CreatedAt      time.Time `json:"created_at"`
}
