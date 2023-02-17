package models

import (
	"time"
)

type Invoice struct {
	ID                int64     `json:"id"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	SubscriptionID    int64     `json:"subscription_id"`
	SubscriptionPrice float64   `json:"subscription_price"`
}
